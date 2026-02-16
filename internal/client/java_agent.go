package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elazarl/goproxy"
)

func findToolsJar(pid string) string {
	// 1. Check JAVA_HOME if set
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		tj := filepath.Join(jh, "lib", "tools.jar")
		if _, err := os.Stat(tj); err == nil {
			return tj
		}
	}

	// 2. Try to find from the target process itself (Linux only)
	if exe, err := os.Readlink(fmt.Sprintf("/proc/%s/exe", pid)); err == nil {
		// exe is usually .../jre/bin/java or .../bin/java
		// tools.jar is usually in ../lib/tools.jar (relative to bin) or ../../lib/tools.jar
		base := filepath.Dir(filepath.Dir(exe))
		tj := filepath.Join(base, "lib", "tools.jar")
		if _, err := os.Stat(tj); err == nil {
			return tj
		}
		// Try one level up (if it was in jre/bin)
		base2 := filepath.Dir(base)
		tj2 := filepath.Join(base2, "lib", "tools.jar")
		if _, err := os.Stat(tj2); err == nil {
			return tj2
		}
	}

	// 3. Common Linux paths
	paths := []string{
		"/usr/lib/jvm/java-8-openjdk-amd64/lib/tools.jar",
		"/usr/lib/jvm/java-8-oracle/lib/tools.jar",
		"/usr/lib/jvm/default-java/lib/tools.jar",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func BuildAndAttachAgent(pid string, proxyAddr string) error {
	// ... (compilation setup remains same until step 5)
	tmpDir, err := os.MkdirTemp("", "agent-proxy-java")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	javaFile := filepath.Join(tmpDir, "ProxyAgent.java")
	manifestFile := filepath.Join(tmpDir, "MANIFEST.MF")
	jarFile := filepath.Join(os.TempDir(), "agent-proxy-injector.jar")

	caCertBase64 := base64.StdEncoding.EncodeToString([]byte(goproxy.CA_CERT))

	// Write Java code
	javaCode := fmt.Sprintf(`
import java.lang.instrument.Instrumentation;
import java.io.ByteArrayInputStream;
import java.security.KeyStore;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManagerFactory;

public class ProxyAgent {
    public static void agentmain(String agentArgs, Instrumentation inst) {
        try {
            String[] args = agentArgs.split(":");
            String host = args[0];
            String port = args[1];
            String caCertB64 = "%s";

            // 1. Set System Properties for standard HttpURLConnection
            System.setProperty("http.proxyHost", host);
            System.setProperty("http.proxyPort", port);
            System.setProperty("https.proxyHost", host);
            System.setProperty("https.proxyPort", port);
            System.setProperty("http.nonProxyHosts", "");

            // 2. Trust the Agent Proxy CA Cert (for HTTPS)
            try {
                byte[] decoded = java.util.Base64.getDecoder().decode(caCertB64);
                CertificateFactory cf = CertificateFactory.getInstance("X.509");
                X509Certificate caCert = (X509Certificate) cf.generateCertificate(new ByteArrayInputStream(decoded));

                KeyStore ks = KeyStore.getInstance(KeyStore.getDefaultType());
                ks.load(null, null);
                ks.setCertificateEntry("agent-proxy-ca", caCert);

                TrustManagerFactory tmf = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
                tmf.init(ks);

                SSLContext sc = SSLContext.getInstance("TLS");
                sc.init(null, tmf.getTrustManagers(), new java.security.SecureRandom());
                SSLContext.setDefault(sc);
                
                System.out.println("[AgentProxy] HTTPS CA certificate injected successfully");
            } catch (Exception e) {
                System.err.println("[AgentProxy] Failed to inject CA certificate: " + e.getMessage());
            }

            System.out.println("[AgentProxy] Interception enabled for " + host + ":" + port);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`, caCertBase64)
	if err := os.WriteFile(javaFile, []byte(javaCode), 0644); err != nil {
		return err
	}

	// 2. Compile Agent
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

	// 5. Create and Compile Injector
	toolsJar := findToolsJar(pid)
	injectorFile := filepath.Join(tmpDir, "Injector.java")
	injectorCode := `
import java.io.File;
import java.lang.reflect.Method;

public class Injector {
    public static void main(String[] args) throws Exception {
        String pid = args[0];
        String jarPath = args[1];
        String options = args[2];
        try {
            Class<?> vmClass;
            try {
                vmClass = Class.forName("com.sun.tools.attach.VirtualMachine");
            } catch (ClassNotFoundException e) {
                // Try to load it from a provided path if needed, though usually handled by classpath
                throw new Exception("VirtualMachine class not found. Ensure tools.jar is on classpath.");
            }
            Method attachMethod = vmClass.getMethod("attach", String.class);
            Method loadAgentMethod = vmClass.getMethod("loadAgent", String.class, String.class);
            Method detachMethod = vmClass.getMethod("detach");
            Object vm = attachMethod.invoke(null, pid);
            loadAgentMethod.invoke(vm, jarPath, options);
            detachMethod.invoke(vm);
            System.out.println("Successfully attached");
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }
    }
}`
	if err := os.WriteFile(injectorFile, []byte(injectorCode), 0644); err != nil {
		return err
	}

	// Compile Injector
	// Try Java 9+ first
	if _, err := exec.Command("javac", "--add-modules", "jdk.attach", injectorFile).CombinedOutput(); err != nil {
		// Fallback for Java 8
		if toolsJar != "" {
			if out, err := exec.Command("javac", "-cp", toolsJar, injectorFile).CombinedOutput(); err != nil {
				log.Printf("Warning: Injector compilation with tools.jar failed: %s", string(out))
			}
		}
	}

	// 6. Attempt attachment
	host, port := splitAddr(proxyAddr)
	agentArgs := fmt.Sprintf("%s:%s", host, port)

	// Try the Injector (Java 9+)
	cmd := exec.Command("java", "--add-modules", "jdk.attach", "-cp", tmpDir, "Injector", pid, jarFile, agentArgs)
	if out, err := cmd.CombinedOutput(); err == nil && strings.Contains(string(out), "Successfully attached") {
		return nil
	} else if toolsJar != "" {
		// Try Java 8 fallback for Injector
		cmd8 := exec.Command("java", "-cp", tmpDir+":"+toolsJar, "Injector", pid, jarFile, agentArgs)
		if out8, err8 := cmd8.CombinedOutput(); err8 == nil && strings.Contains(string(out8), "Successfully attached") {
			return nil
		}
	}

	// 7. Last resort: jcmd
	cmdLast := exec.Command("jcmd", pid, "JVMTI.agent_load", jarFile, agentArgs)
	if _, err := cmdLast.CombinedOutput(); err != nil {
		cmdLast2 := exec.Command("jcmd", pid, "JVMTI.agent_load", jarFile+"="+agentArgs)
		if out2, err2 := cmdLast2.CombinedOutput(); err2 != nil {
			return fmt.Errorf("attachment failed after multiple attempts.\nPID: %s\nTools.jar found: %v\njcmd error: %s",
				pid, toolsJar != "", strings.TrimSpace(string(out2)))
		}
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
