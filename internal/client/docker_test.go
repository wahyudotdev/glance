package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// mockDockerClient implements client.APIClient for testing
type mockDockerClient struct {
	client.APIClient
	listFunc        func() ([]container.Summary, error)
	inspectFunc     func(id string) (container.InspectResponse, error)
	createFunc      func(hostConfig *container.HostConfig) (container.CreateResponse, error)
	execCreateFunc  func() (types.IDResponse, error) // nolint:staticcheck
	execAttachFunc  func() (types.HijackedResponse, error)
	execInspectFunc func() (container.ExecInspect, error)
	copyToFunc      func() error
	stopFunc        func() error
	renameFunc      func() error
	removeFunc      func() error
	startFunc       func() error
	closeFunc       func() error
}

func (m *mockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return nil, nil
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	if m.inspectFunc != nil {
		return m.inspectFunc(containerID)
	}
	return container.InspectResponse{}, nil
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(hostConfig)
	}
	return container.CreateResponse{ID: "new-id"}, nil
}

func (m *mockDockerClient) ContainerExecCreate(ctx context.Context, container string, config container.ExecOptions) (types.IDResponse, error) { // nolint:staticcheck
	if m.execCreateFunc != nil {
		return m.execCreateFunc()
	}
	return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
}

func (m *mockDockerClient) ContainerExecAttach(ctx context.Context, execID string, check container.ExecStartOptions) (types.HijackedResponse, error) {
	if m.execAttachFunc != nil {
		return m.execAttachFunc()
	}
	return types.HijackedResponse{}, nil
}

func (m *mockDockerClient) ContainerExecInspect(ctx context.Context, execID string) (container.ExecInspect, error) {
	if m.execInspectFunc != nil {
		return m.execInspectFunc()
	}
	return container.ExecInspect{ExitCode: 0}, nil
}

func (m *mockDockerClient) CopyToContainer(ctx context.Context, container, path string, content io.Reader, options container.CopyToContainerOptions) error {
	if m.copyToFunc != nil {
		return m.copyToFunc()
	}
	return nil
}

func (m *mockDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.stopFunc != nil {
		return m.stopFunc()
	}
	return nil
}

func (m *mockDockerClient) ContainerRename(ctx context.Context, containerID, newContainerName string) error {
	if m.renameFunc != nil {
		return m.renameFunc()
	}
	return nil
}

func (m *mockDockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	if m.removeFunc != nil {
		return m.removeFunc()
	}
	return nil
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	if m.startFunc != nil {
		return m.startFunc()
	}
	return nil
}

func (m *mockDockerClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestListDockerContainers(t *testing.T) {
	oldNewClient := newDockerClient
	defer func() { newDockerClient = oldNewClient }()

	mock := &mockDockerClient{
		listFunc: func() ([]container.Summary, error) {
			return []container.Summary{
				{
					ID:    "123456789012345",
					Names: []string{"/test-container"},
					Image: "test-image",
					State: "running",
					NetworkSettings: &container.NetworkSettingsSummary{
						Networks: map[string]*network.EndpointSettings{
							"bridge": {IPAddress: "172.17.0.2"},
						},
					},
					Labels: map[string]string{"glance.interception": "active"},
				},
			}, nil
		},
	}
	newDockerClient = func() (client.APIClient, error) { return mock, nil }

	containers, err := ListDockerContainers()
	if err != nil {
		t.Fatalf("ListDockerContainers failed: %v", err)
	}

	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}

	c := containers[0]
	if c.ID != "123456789012" {
		t.Errorf("Expected short ID 123456789012, got %s", c.ID)
	}
	if c.Name != "test-container" {
		t.Errorf("Expected name test-container, got %s", c.Name)
	}
	if !c.Intercepted {
		t.Error("Expected container to be marked as intercepted")
	}
}

func TestInterceptDocker(t *testing.T) {
	oldNewClient := newDockerClient
	defer func() { newDockerClient = oldNewClient }()

	var capturedHostConfig *container.HostConfig

	mock := &mockDockerClient{
		inspectFunc: func(id string) (container.InspectResponse, error) {
			return container.InspectResponse{
				Config: &container.Config{
					Image: "test-image",
					Env:   []string{"FOO=BAR"},
				},
				ContainerJSONBase: &container.ContainerJSONBase{
					Name: "/test-container",
					HostConfig: &container.HostConfig{
						ExtraHosts: []string{"existing:1.2.3.4"},
					},
				},
				NetworkSettings: &container.NetworkSettings{
					Networks: map[string]*network.EndpointSettings{
						"bridge": {Gateway: "172.17.0.1"},
					},
				},
			}, nil
		},
		createFunc: func(hostConfig *container.HostConfig) (container.CreateResponse, error) {
			capturedHostConfig = hostConfig
			return container.CreateResponse{ID: "new-id"}, nil
		},
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
		},
		execAttachFunc: func() (types.HijackedResponse, error) {
			r, w := net.Pipe()
			bufr := bufio.NewReader(r)
			go func() {
				header := []byte{1, 0, 0, 0, 0, 0, 0, 6}
				_, _ = w.Write(append(header, []byte("alpine")...))
				_ = w.Close()
			}()
			return types.HijackedResponse{Reader: bufr, Conn: w}, nil
		},
		execInspectFunc: func() (container.ExecInspect, error) {
			return container.ExecInspect{ExitCode: 0}, nil
		},
	}
	newDockerClient = func() (client.APIClient, error) { return mock, nil }

	err := InterceptDocker("cont1", "localhost:15500")
	if err != nil {
		t.Fatalf("InterceptDocker failed: %v", err)
	}

	if capturedHostConfig == nil {
		t.Fatal("HostConfig was not captured")
	}

	found := false
	for _, host := range capturedHostConfig.ExtraHosts {
		if host == "host.docker.internal:host-gateway" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected host.docker.internal:host-gateway in ExtraHosts, got %v", capturedHostConfig.ExtraHosts)
	}
}

func TestStopInterceptDocker(t *testing.T) {
	oldNewClient := newDockerClient
	defer func() { newDockerClient = oldNewClient }()

	mock := &mockDockerClient{
		inspectFunc: func(id string) (container.InspectResponse, error) {
			if strings.HasSuffix(id, "-glance-backup") {
				return container.InspectResponse{
					ContainerJSONBase: &container.ContainerJSONBase{ID: "backup-id", Name: "/test-backup"},
				}, nil
			}
			return container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{Name: "/test-container"},
			}, nil
		},
	}
	newDockerClient = func() (client.APIClient, error) { return mock, nil }

	err := StopInterceptDocker("cont1")
	if err != nil {
		t.Fatalf("StopInterceptDocker failed: %v", err)
	}
}

func TestFindHostIP(t *testing.T) {
	ctx := context.Background()
	mock := &mockDockerClient{
		inspectFunc: func(id string) (container.InspectResponse, error) {
			return container.InspectResponse{
				NetworkSettings: &container.NetworkSettings{
					Networks: map[string]*network.EndpointSettings{
						"bridge": {Gateway: "172.17.0.1"},
					},
				},
			}, nil
		},
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
		},
		execAttachFunc: func() (types.HijackedResponse, error) {
			r, w := net.Pipe()
			bufr := bufio.NewReader(r)
			go func() {
				// Empty output to trigger fallback
				_ = w.Close()
			}()
			return types.HijackedResponse{Reader: bufr, Conn: w}, nil
		},
		execInspectFunc: func() (container.ExecInspect, error) {
			return container.ExecInspect{ExitCode: 1}, nil
		},
	}

	ip, err := findHostIP(ctx, mock, "cont1")
	if err != nil {
		t.Fatalf("findHostIP failed: %v", err)
	}
	if ip != "172.17.0.1" {
		t.Errorf("Expected 172.17.0.1, got %s", ip)
	}
}

func TestFindHostIP_Fallback(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	mock := &mockDockerClient{
		inspectFunc: func(id string) (container.InspectResponse, error) {
			// Return empty networks to trigger fallback
			return container.InspectResponse{
				NetworkSettings: &container.NetworkSettings{
					Networks: map[string]*network.EndpointSettings{},
				},
			}, nil
		},
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
		},
		execAttachFunc: func() (types.HijackedResponse, error) {
			r, w := net.Pipe()
			bufr := bufio.NewReader(r)
			callCount++
			go func() {
				header := []byte{1, 0, 0, 0, 0, 0, 0, 10}
				if callCount == 1 {
					// First call: getent hosts (fail it)
					_ = w.Close()
				} else {
					// Second call: ip route show (success)
					_, _ = w.Write(append(header, []byte("172.18.0.1\n")...))
					_ = w.Close()
				}
			}()
			return types.HijackedResponse{Reader: bufr, Conn: w}, nil
		},
		execInspectFunc: func() (container.ExecInspect, error) {
			if callCount == 1 {
				return container.ExecInspect{ExitCode: 1}, nil
			}
			return container.ExecInspect{ExitCode: 0}, nil
		},
	}

	ip, err := findHostIP(ctx, mock, "cont1")
	if err != nil {
		t.Fatalf("findHostIP failed: %v", err)
	}
	if ip != "172.18.0.1" {
		t.Errorf("Expected 172.18.0.1, got %s", ip)
	}
}

func TestDetectOS(t *testing.T) {
	ctx := context.Background()

	t.Run("Alpine", func(t *testing.T) {
		mock := &mockDockerClient{
			execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
				return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
			},
			execAttachFunc: func() (types.HijackedResponse, error) {
				r, w := net.Pipe()
				bufr := bufio.NewReader(r)
				go func() {
					header := []byte{1, 0, 0, 0, 0, 0, 0, 6}
					_, _ = w.Write(append(header, []byte("alpine")...))
					_ = w.Close()
				}()
				return types.HijackedResponse{Reader: bufr, Conn: w}, nil
			},
			execInspectFunc: func() (container.ExecInspect, error) {
				return container.ExecInspect{ExitCode: 0}, nil
			},
		}
		os := detectOS(ctx, mock, "cont1")
		if os != "alpine" {
			t.Errorf("Expected alpine, got %s", os)
		}
	})

	t.Run("RHEL", func(t *testing.T) {
		mock := &mockDockerClient{
			execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
				return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
			},
			execAttachFunc: func() (types.HijackedResponse, error) {
				r, w := net.Pipe()
				bufr := bufio.NewReader(r)
				go func() {
					header := []byte{1, 0, 0, 0, 0, 0, 0, 4}
					_, _ = w.Write(append(header, []byte("rhel")...))
					_ = w.Close()
				}()
				return types.HijackedResponse{Reader: bufr, Conn: w}, nil
			},
			execInspectFunc: func() (container.ExecInspect, error) {
				return container.ExecInspect{ExitCode: 0}, nil
			},
		}
		os := detectOS(ctx, mock, "cont1")
		if os != "rhel" {
			t.Errorf("Expected rhel, got %s", os)
		}
	})
}

func TestInjectCACert(t *testing.T) {
	ctx := context.Background()
	mock := &mockDockerClient{
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
		},
		execAttachFunc: func() (types.HijackedResponse, error) {
			r, w := net.Pipe()
			bufr := bufio.NewReader(r)
			go func() {
				header := []byte{1, 0, 0, 0, 0, 0, 0, 6}
				_, _ = w.Write(append(header, []byte("debian")...))
				_ = w.Close()
			}()
			return types.HijackedResponse{Reader: bufr, Conn: w}, nil
		},
		execInspectFunc: func() (container.ExecInspect, error) {
			return container.ExecInspect{ExitCode: 0}, nil
		},
	}

	err := injectCACert(ctx, mock, "cont1")
	if err != nil {
		t.Fatalf("injectCACert failed: %v", err)
	}
}

func TestExecInContainer_Error(t *testing.T) {
	ctx := context.Background()
	mock := &mockDockerClient{
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{}, fmt.Errorf("create fail") // nolint:staticcheck
		},
	}

	code, _, err := execInContainer(ctx, mock, "cont1", []string{"ls"})
	if err == nil {
		t.Error("Expected error")
	}
	if code != -1 {
		t.Errorf("Expected -1 code, got %d", code)
	}
}

func TestInjectCACert_JavaLogic(t *testing.T) {
	ctx := context.Background()

	mock := &mockDockerClient{
		execCreateFunc: func() (types.IDResponse, error) { // nolint:staticcheck
			return types.IDResponse{ID: "exec-id"}, nil // nolint:staticcheck
		},
		execAttachFunc: func() (types.HijackedResponse, error) {
			r, w := net.Pipe()
			bufr := bufio.NewReader(r)
			go func() {
				header := []byte{1, 0, 0, 0, 0, 0, 0, 10}
				_, _ = w.Write(append(header, []byte("success\n")...))
				_ = w.Close()
			}()
			return types.HijackedResponse{Reader: bufr, Conn: w}, nil
		},
		execInspectFunc: func() (container.ExecInspect, error) {
			return container.ExecInspect{ExitCode: 0}, nil
		},
	}

	// We want to test that it tries to find cacerts
	err := injectCACert(ctx, mock, "cont1")
	if err != nil {
		t.Fatalf("injectCACert failed: %v", err)
	}
}
