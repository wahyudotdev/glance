// Package client handles external integrations with Android, Java, and other clients.
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

// BuildAndAttachAgent dynamically compiles a Java agent and attaches it to a running JVM.
func BuildAndAttachAgent(pid string, proxyAddr string) error {
	// ... (compilation setup remains same until step 5)
	tmpDir, err := os.MkdirTemp("", "agent-proxy-java")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck

	javaFile := filepath.Join(tmpDir, "ProxyAgent.java")
	manifestFile := filepath.Join(tmpDir, "MANIFEST.MF")
	jarFile := filepath.Join(os.TempDir(), "agent-proxy-injector.jar")

	caCertBase64 := base64.StdEncoding.EncodeToString([]byte(goproxy.CA_CERT))

	// Write Java code
	javaCode := fmt.Sprintf(`
import java.lang.instrument.Instrumentation;
import java.io.ByteArrayInputStream;
import java.security.cert.X509Certificate;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

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

            // 2. Aggressively force trust (for HTTPS)
            try {
                TrustManager[] trustAllCerts = new TrustManager[]{
                    new X509TrustManager() {
                        public X509Certificate[] getAcceptedIssuers() { return null; }
                        public void checkClientTrusted(X509Certificate[] certs, String authType) {}
                        public void checkServerTrusted(X509Certificate[] certs, String authType) {}
                    }
                };

                SSLContext sc = SSLContext.getInstance("SSL");
                sc.init(null, trustAllCerts, new java.security.SecureRandom());
                
                // Set as default for the entire JVM
                SSLContext.setDefault(sc);
                
                // Also set for HttpsURLConnection specifically
                javax.net.ssl.HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
                javax.net.ssl.HttpsURLConnection.setDefaultHostnameVerifier(new javax.net.ssl.HostnameVerifier() {
                    public boolean verify(String hostname, javax.net.ssl.SSLSession session) {
                        return true;
                    }
                });

                System.out.println("[AgentProxy] Aggressive TLS trust-all injected");
            } catch (Exception e) {
                System.err.println("[AgentProxy] Failed to inject aggressive trust: " + e.getMessage());
            }

            System.out.println("[AgentProxy] Interception enabled for " + host + ":" + port);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`, caCertBase64)
	if err := os.WriteFile(javaFile, []byte(javaCode), 0600); err != nil {
		return err
	}

	// 2. Compile Agent
	// Use --release 8 (modern JDKs) or -source/-target (older JDKs) to ensure compatibility with Java 8 (version 52.0)
	// #nosec G204
	cmdCompile := exec.Command("javac", "--release", "8", javaFile)
	if _, err := cmdCompile.CombinedOutput(); err != nil {
		// Fallback for older javac that doesn't support --release
		// #nosec G204
		cmdFallback := exec.Command("javac", "-source", "1.8", "-target", "1.8", javaFile)
		if out2, err2 := cmdFallback.CombinedOutput(); err2 != nil {
			return fmt.Errorf("javac failed to compile for Java 8 compatibility: %s", string(out2))
		}
	}

	// 3. Create Manifest
	manifest := "Agent-Class: ProxyAgent\nCan-Retransform-Classes: true\n"
	if err := os.WriteFile(manifestFile, []byte(manifest), 0600); err != nil {
		return err
	}

	// 4. Create JAR
	// #nosec G204
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
	if err := os.WriteFile(injectorFile, []byte(injectorCode), 0600); err != nil {
		return err
	}

	// Compile Injector
	// Try Java 9+ first
	// #nosec G204
	if _, err := exec.Command("javac", "--add-modules", "jdk.attach", injectorFile).CombinedOutput(); err != nil {
		// Fallback for Java 8
		if toolsJar != "" {
			// #nosec G204
			if out, err := exec.Command("javac", "-cp", toolsJar, injectorFile).CombinedOutput(); err != nil {
				log.Printf("Warning: Injector compilation with tools.jar failed: %s", string(out))
			}
		}
	}

	// 6. Attempt attachment
	host, port := splitAddr(proxyAddr)
	agentArgs := fmt.Sprintf("%s:%s", host, port)

	// Try the Injector (Java 9+)
	// #nosec G204
	cmd := exec.Command("java", "--add-modules", "jdk.attach", "-cp", tmpDir, "Injector", pid, jarFile, agentArgs)
	if out, err := cmd.CombinedOutput(); err == nil && strings.Contains(string(out), "Successfully attached") {
		return nil
	} else if toolsJar != "" {
		// Try Java 8 fallback for Injector
		// #nosec G204
		cmd8 := exec.Command("java", "-cp", tmpDir+":"+toolsJar, "Injector", pid, jarFile, agentArgs)
		if out8, err8 := cmd8.CombinedOutput(); err8 == nil && strings.Contains(string(out8), "Successfully attached") {
			return nil
		}
	}

	// 7. Last resort: jcmd
	// #nosec G204
	cmdLast := exec.Command("jcmd", pid, "JVMTI.agent_load", jarFile, agentArgs)
	if _, err := cmdLast.CombinedOutput(); err != nil {
		// #nosec G204
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
