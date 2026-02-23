package client

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"glance/internal/model"
	"net"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/elazarl/goproxy"
)

var newDockerClient = func() (client.APIClient, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// ListDockerContainers returns a list of running Docker containers.
func ListDockerContainers() ([]model.DockerContainer, error) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		return nil, err
	}
	defer func() { _ = cli.Close() }()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: false})
	if err != nil {
		return nil, err
	}

	var result []model.DockerContainer
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		ip := ""
		if c.NetworkSettings != nil && len(c.NetworkSettings.Networks) > 0 {
			// Get IP from the first network
			for _, net := range c.NetworkSettings.Networks {
				if net.IPAddress != "" {
					ip = net.IPAddress
					break
				}
			}
		}

		result = append(result, model.DockerContainer{
			ID:          c.ID[:12], // Short ID
			Name:        name,
			Image:       c.Image,
			State:       c.State,
			IPAddress:   ip,
			Intercepted: isIntercepted(c),
		})
	}

	return result, nil
}

func isIntercepted(c container.Summary) bool {
	// Look for proxy env vars in the container info
	// Actually, the 'c' object (types.Container) doesn't have the Env list.
	// We might need to inspect it if we want accurate state in the list.
	// For now, check if the name suggests it's an intercepted version or just check Labels if we added any.
	for _, name := range c.Names {
		if strings.HasSuffix(name, "-glance-backup") {
			return false // This is the old one
		}
	}

	// Check if this container was created by Glance (we can add a label for this)
	if c.Labels != nil && c.Labels["glance.interception"] == "active" {
		return true
	}

	return false
}

// InterceptDocker starts intercepting a container by recreating it with proxy environment variables.
func InterceptDocker(containerID string, proxyAddr string) error {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	// 1. Inspect original container
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %v", err)
	}

	// 2. Determine Host IP for the proxy
	// We try to find a numeric IP first (Gateway). If it resolves to 127.0.0.1
	// or fails, we use 'host.docker.internal' as the last resort.
	hostIP, err := findHostIP(ctx, cli, containerID)
	if err != nil || hostIP == "127.0.0.1" || hostIP == "localhost" {
		hostIP = "host.docker.internal"
	}

	_, port, err := net.SplitHostPort(proxyAddr)
	if err != nil {
		port = "15500" // Default
	}
	proxyURL := fmt.Sprintf("http://%s:%s", hostIP, port)

	// 3. Prepare updated configuration
	newConfig := *inspect.Config
	if newConfig.Labels == nil {
		newConfig.Labels = make(map[string]string)
	}
	newConfig.Labels["glance.interception"] = "active"
	newConfig.Labels["glance.original_id"] = containerID

	// Standard Env Vars
	newConfig.Env = append(newConfig.Env,
		fmt.Sprintf("HTTP_PROXY=%s", proxyURL),
		fmt.Sprintf("HTTPS_PROXY=%s", proxyURL),
		fmt.Sprintf("http_proxy=%s", proxyURL),
		fmt.Sprintf("https_proxy=%s", proxyURL),
		"NO_PROXY=localhost,127.0.0.1",
	)

	// Java Specific (JVM properties)
	javaOpts := fmt.Sprintf("-Dhttp.proxyHost=%s -Dhttp.proxyPort=%s -Dhttps.proxyHost=%s -Dhttps.proxyPort=%s", hostIP, port, hostIP, port)

	// Check if JAVA_TOOL_OPTIONS already exists and append
	foundJava := false
	for i, env := range newConfig.Env {
		if strings.HasPrefix(env, "JAVA_TOOL_OPTIONS=") {
			newConfig.Env[i] = env + " " + javaOpts
			foundJava = true
			break
		}
	}
	if !foundJava {
		newConfig.Env = append(newConfig.Env, "JAVA_TOOL_OPTIONS="+javaOpts)
	}

	// 4. Stop and Rename old container
	timeout := 10
	if errStop := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); errStop != nil {
		return fmt.Errorf("failed to stop container: %v", errStop)
	}

	oldName := strings.TrimPrefix(inspect.Name, "/")
	backupName := oldName + "-glance-backup"
	if errRename := cli.ContainerRename(ctx, containerID, backupName); errRename != nil {
		return fmt.Errorf("failed to rename container: %v", errRename)
	}

	// 5. Create and Start new container with original name
	// Note: We copy HostConfig and NetworkConfig to keep everything else the same.
	newHostConfig := *inspect.HostConfig
	newHostConfig.ExtraHosts = append(newHostConfig.ExtraHosts, "host.docker.internal:host-gateway")

	// We must clear dynamic fields like MacAddress and IPAddress to maintain compatibility with older API versions.
	endpoints := make(map[string]*network.EndpointSettings)
	for netName, settings := range inspect.NetworkSettings.Networks {
		// Create a copy and clear dynamic fields
		cp := *settings
		cp.MacAddress = ""
		cp.IPAddress = ""
		cp.GlobalIPv6Address = ""
		if cp.IPAMConfig != nil {
			cp.IPAMConfig.IPv4Address = ""
			cp.IPAMConfig.IPv6Address = ""
		}
		endpoints[netName] = &cp
	}

	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
	newContainer, err := cli.ContainerCreate(ctx, &newConfig, &newHostConfig, networkingConfig, nil, oldName)
	if err != nil {
		// Rollback rename if create fails
		_ = cli.ContainerRename(ctx, containerID, oldName)
		return fmt.Errorf("failed to create new container: %v", err)
	}

	if err := cli.ContainerStart(ctx, newContainer.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start new container: %v", err)
	}

	// 6. Inject and Trust CA Cert (since we can't easily mount a dynamic byte slice without a physical file)
	// We'll use our existing injectCACert helper on the NEW container.
	if err := injectCACert(ctx, cli, newContainer.ID); err != nil {
		fmt.Printf("Warning: Failed to inject CA cert: %v\n", err)
	}

	return nil
}

// StopInterceptDocker stops intercepting a container by reverting to the backup.
func StopInterceptDocker(containerID string) error {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	// 1. Inspect current container to find its name
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}

	currentName := strings.TrimPrefix(inspect.Name, "/")
	if !strings.HasSuffix(currentName, "-glance-backup") {
		// Try to find the backup if we are currently running the intercepted one
		backupName := currentName + "-glance-backup"

		// Stop and Remove intercepted container
		timeout := 10
		_ = cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
		if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
			return fmt.Errorf("failed to remove intercepted container: %v", err)
		}

		// Rename backup back to original name and start it
		backupContainer, err := cli.ContainerInspect(ctx, backupName)
		if err != nil {
			return fmt.Errorf("could not find backup container %s: %v", backupName, err)
		}

		if err := cli.ContainerRename(ctx, backupContainer.ID, currentName); err != nil {
			return fmt.Errorf("failed to restore container name: %v", err)
		}

		if err := cli.ContainerStart(ctx, backupContainer.ID, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to restart original container: %v", err)
		}
	}

	return nil
}

func findHostIP(ctx context.Context, cli client.APIClient, containerID string) (string, error) {
	// 1. Try to resolve host.docker.internal first
	// This is the most reliable way for Docker Desktop (macOS/Windows)
	_, out, _ := execInContainer(ctx, cli, containerID, []string{"getent", "hosts", "host.docker.internal"})
	if out != "" {
		parts := strings.Fields(out)
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if !strings.Contains(ip, ":") { // Prefer IPv4
				return ip, nil
			}
		}
	}

	// 2. Try to get the Gateway from the container's network settings
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err == nil {
		for _, net := range inspect.NetworkSettings.Networks {
			if net.Gateway != "" {
				return net.Gateway, nil
			}
		}
	}

	// 3. Fallback: try to find the default route gateway via shell
	_, out, err = execInContainer(ctx, cli, containerID, []string{"sh", "-c", "ip route show | grep default | awk '{print $3}'"})
	if err == nil && out != "" {
		ip := strings.TrimSpace(out)
		if ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("could not determine host IP")
}

func injectCACert(ctx context.Context, cli client.APIClient, containerID string) error {
	// 1. Detect OS
	osType := detectOS(ctx, cli, containerID)

	// 2. Prepare CA cert as tar
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	certBytes := goproxy.CA_CERT
	hdr := &tar.Header{
		Name: "glance-ca.crt",
		Mode: 0644,
		Size: int64(len(certBytes)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := tw.Write(certBytes); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	// 3. Copy to appropriate location based on OS
	destDir := "/usr/local/share/ca-certificates"
	trustCmd := "update-ca-certificates"

	if osType == "rhel" {
		destDir = "/etc/pki/ca-trust/source/anchors"
		trustCmd = "update-ca-trust"
	}

	// Ensure directory exists
	_, _, _ = execInContainer(ctx, cli, containerID, []string{"mkdir", "-p", destDir})

	// Copy cert to OS store
	if err := cli.CopyToContainer(ctx, containerID, destDir, &buf, container.CopyToContainerOptions{}); err != nil {
		return err
	}

	// 4. Update OS trust store
	_, _, _ = execInContainer(ctx, cli, containerID, []string{trustCmd})

	// 5. Java Specific: Import into Java Keystore (cacerts)
	// Try to find keytool and cacerts
	_, javaHome, _ := execInContainer(ctx, cli, containerID, []string{"sh", "-c", "echo $JAVA_HOME"})
	javaHome = strings.TrimSpace(javaHome)

	// Search for cacerts in common locations if JAVA_HOME is not set
	cacertsPaths := []string{
		"/etc/ssl/certs/java/cacerts", // Debian/Ubuntu
		javaHome + "/lib/security/cacerts",
		javaHome + "/jre/lib/security/cacerts",
		"/usr/lib/jvm/default-jvm/jre/lib/security/cacerts", // Alpine
	}

	var foundCacerts string
	for _, path := range cacertsPaths {
		if path == "" || path == "/lib/security/cacerts" || path == "/jre/lib/security/cacerts" {
			continue
		}
		exitCode, _, _ := execInContainer(ctx, cli, containerID, []string{"test", "-f", path})
		if exitCode == 0 {
			foundCacerts = path
			break
		}
	}

	if foundCacerts != "" {
		// Prepare cert for keytool (copy to a temp location in container)
		var buf2 bytes.Buffer
		tw2 := tar.NewWriter(&buf2)
		hdr2 := &tar.Header{Name: "glance-ca-java.crt", Mode: 0644, Size: int64(len(certBytes))}
		_ = tw2.WriteHeader(hdr2)
		_, _ = tw2.Write(certBytes)
		_ = tw2.Close()
		_ = cli.CopyToContainer(ctx, containerID, "/tmp", &buf2, container.CopyToContainerOptions{})

		// Import using keytool
		// Password is 'changeit' by default in most JVMs
		importCmd := []string{"keytool", "-importcert", "-trustcacerts", "-file", "/tmp/glance-ca-java.crt", "-keystore", foundCacerts, "-storepass", "changeit", "-noprompt", "-alias", "glance-ca"}
		exitCode, out, err := execInContainer(ctx, cli, containerID, importCmd)
		if exitCode != 0 {
			fmt.Printf("Warning: Failed to import CA into Java keystore (%s): %s %v\n", foundCacerts, out, err)
		} else {
			fmt.Printf("Successfully imported CA into Java keystore: %s\n", foundCacerts)
		}
	}

	return nil
}

func detectOS(ctx context.Context, cli client.APIClient, containerID string) string {
	_, out, err := execInContainer(ctx, cli, containerID, []string{"cat", "/etc/os-release"})
	if err != nil {
		return "debian" // Default fallback
	}

	content := strings.ToLower(out)
	if strings.Contains(content, "alpine") {
		return "alpine"
	}
	if strings.Contains(content, "rhel") || strings.Contains(content, "centos") || strings.Contains(content, "fedora") {
		return "rhel"
	}
	return "debian"
}

func execInContainer(ctx context.Context, cli client.APIClient, containerID string, cmd []string) (int, string, error) {
	config := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := cli.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return -1, "", err
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return -1, "", err
	}
	defer resp.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, resp.Reader)
	if err != nil {
		return -1, "", err
	}

	// Inspect to get the exit code
	inspect, err := cli.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return -1, stdout.String(), err
	}

	return inspect.ExitCode, stdout.String(), nil
}
