package client

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func BuildAndAttachAgent(pid string, proxyAddr string) error {
	// 1. Setup temp directory for compilation
	tmpDir, err := os.MkdirTemp("", "agent-proxy-java")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	javaFile := filepath.Join(tmpDir, "ProxyAgent.java")
	manifestFile := filepath.Join(tmpDir, "MANIFEST.MF")
	jarFile := filepath.Join(os.TempDir(), "agent-proxy-injector.jar")

	// Write Java code
	javaCode := `
import java.lang.instrument.Instrumentation;
public class ProxyAgent {
    public static void agentmain(String agentArgs, Instrumentation inst) {
        try {
            String[] args = agentArgs.split(":");
            System.setProperty("http.proxyHost", args[0]);
            System.setProperty("http.proxyPort", args[1]);
            System.setProperty("https.proxyHost", args[0]);
            System.setProperty("https.proxyPort", args[1]);
            System.setProperty("http.nonProxyHosts", "");
            System.out.println("[AgentProxy] Injected " + args[0] + ":" + args[1]);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`
	if err := os.WriteFile(javaFile, []byte(javaCode), 0644); err != nil {
		return err
	}

	// 2. Compile
	if out, err := exec.Command("javac", javaFile).CombinedOutput(); err != nil {
		return fmt.Errorf("javac failed: %s %v", string(out), err)
	}

	// 3. Create Manifest
	manifest := "Agent-Class: ProxyAgent\nCan-Retransform-Classes: true\n"
	if err := os.WriteFile(manifestFile, []byte(manifest), 0644); err != nil {
		return err
	}

	// 4. Create JAR
	if out, err := exec.Command("jar", "cmf", manifestFile, jarFile, "-C", tmpDir, "ProxyAgent.class").CombinedOutput(); err != nil {
		return fmt.Errorf("jar creation failed: %s %v", string(out), err)
	}

	// 5. Attach using jcmd
	host, port := splitAddr(proxyAddr)
	agentArgs := fmt.Sprintf("%s:%s", host, port)

	if out, err := exec.Command("jcmd", pid, "JVMTI.agent_load", jarFile, agentArgs).CombinedOutput(); err != nil {
		return fmt.Errorf("jcmd attachment failed: %s %v", string(out), err)
	}

	return nil
}

func splitAddr(addr string) (string, string) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "127.0.0.1", strings.TrimPrefix(addr, ":")
	}
	if host == "" || host == "::" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	return host, port
}
