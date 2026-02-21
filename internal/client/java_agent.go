// Package client handles external integrations with Android, Java, and other clients.
package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/elazarl/goproxy"
)

var readlink = os.Readlink
var stat = os.Stat

func findToolsJar(pid string) string {
	// 1. Check JAVA_HOME if set
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		tj := filepath.Join(jh, "lib", "tools.jar")
		// #nosec G703
		if _, err := stat(tj); err == nil {
			return tj
		}
	}

	// 2. Try to find from the target process itself (Linux only)
	if exe, err := readlink(fmt.Sprintf("/proc/%s/exe", pid)); err == nil {
		// exe is usually .../jre/bin/java or .../bin/java
		// tools.jar is usually in ../lib/tools.jar (relative to bin) or ../../lib/tools.jar
		base := filepath.Dir(filepath.Dir(exe))
		tj := filepath.Join(base, "lib", "tools.jar")
		if _, err := stat(tj); err == nil {
			return tj
		}
		// Try one level up (if it was in jre/bin)
		base2 := filepath.Dir(base)
		tj2 := filepath.Join(base2, "lib", "tools.jar")
		if _, err := stat(tj2); err == nil {
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
		if _, err := stat(p); err == nil {
			return p
		}
	}

	return ""
}

// BuildAndAttachAgent dynamically compiles a Java agent and attaches it to a running JVM.
func BuildAndAttachAgent(pid string, proxyAddr string) error {
	// ... (compilation setup remains same until step 5)
	tmpDir, err := os.MkdirTemp("", "glance-java")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck

	javaFile := filepath.Join(tmpDir, "ProxyAgent.java")
	manifestFile := filepath.Join(tmpDir, "MANIFEST.MF")
	jarFile := filepath.Join(os.TempDir(), "glance-injector.jar")

	caCertBase64 := base64.StdEncoding.EncodeToString([]byte(goproxy.CA_CERT))

	// Write Java code — uses a named static inner class instead of anonymous classes
	// to avoid ProxyAgent$1.class packaging issues in the JAR
	javaCode := fmt.Sprintf(`
import java.lang.instrument.Instrumentation;
import java.io.ByteArrayInputStream;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.TrustManagerFactory;
import javax.net.ssl.X509TrustManager;
import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLSession;

public class ProxyAgent {

    static class TrustAllManager implements X509TrustManager {
        public X509Certificate[] getAcceptedIssuers() { return new X509Certificate[0]; }
        public void checkClientTrusted(X509Certificate[] certs, String authType) {}
        public void checkServerTrusted(X509Certificate[] certs, String authType) {}
    }

    static class TrustAllHostname implements HostnameVerifier {
        public boolean verify(String hostname, SSLSession session) { return true; }
    }

    public static void agentmain(String agentArgs, Instrumentation inst) {
        try {
            String[] args = agentArgs.split(":");
            String host = args[0];
            String port = args[1];
            String caCertB64 = "%s";

            // 1. Set proxy system properties
            System.setProperty("http.proxyHost", host);
            System.setProperty("http.proxyPort", port);
            System.setProperty("https.proxyHost", host);
            System.setProperty("https.proxyPort", port);
            System.setProperty("http.nonProxyHosts", "");

            // 2. Set trust-all as default SSLContext (covers HttpsURLConnection)
            TrustManager[] trustAll = new TrustManager[]{ new TrustAllManager() };
            SSLContext sc = SSLContext.getInstance("TLS");
            sc.init(null, trustAll, new java.security.SecureRandom());
            SSLContext.setDefault(sc);
            javax.net.ssl.HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
            javax.net.ssl.HttpsURLConnection.setDefaultHostnameVerifier(new TrustAllHostname());
            System.out.println("[Glance] Trust-all SSLContext set as default");

            // 3. Write CA cert to a JKS truststore file and set system properties
            //    This ensures libraries that create their own SSLContext (Apache HttpClient,
            //    OkHttp, Spring RestTemplate) will also trust our proxy CA
            try {
                byte[] pemBytes = java.util.Base64.getDecoder().decode(caCertB64);
                CertificateFactory cf = CertificateFactory.getInstance("X.509");
                Certificate caCert = cf.generateCertificate(new ByteArrayInputStream(pemBytes));

                // Load default cacerts and add our CA
                KeyStore ks = KeyStore.getInstance("JKS");
                String cacertsPath = System.getProperty("java.home")
                    + java.io.File.separator + "lib"
                    + java.io.File.separator + "security"
                    + java.io.File.separator + "cacerts";
                try {
                    FileInputStream fis = new FileInputStream(cacertsPath);
                    try { ks.load(fis, "changeit".toCharArray()); } finally { fis.close(); }
                } catch (Exception e) {
                    ks.load(null, null);
                }
                ks.setCertificateEntry("glance-proxy-ca", caCert);

                // Write modified truststore to temp file
                String trustStorePath = System.getProperty("java.io.tmpdir")
                    + java.io.File.separator + "glance-truststore.jks";
                FileOutputStream fos = new FileOutputStream(trustStorePath);
                try { ks.store(fos, "changeit".toCharArray()); } finally { fos.close(); }

                // Point JVM to use this truststore for any new SSLContext creation
                System.setProperty("javax.net.ssl.trustStore", trustStorePath);
                System.setProperty("javax.net.ssl.trustStorePassword", "changeit");
                System.out.println("[Glance] Custom truststore written to " + trustStorePath);
            } catch (Exception e) {
                System.err.println("[Glance] Truststore setup failed (trust-all still active): " + e.getMessage());
            }

            System.out.println("[Glance] HTTPS interception enabled for " + host + ":" + port);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}`, caCertBase64)
	err = os.WriteFile(javaFile, []byte(javaCode), 0600)
	if err != nil {
		return err
	}

	// 2. Compile Agent
	// Use --release 8 (modern JDKs) or -source/-target (older JDKs) to ensure compatibility with Java 8 (version 52.0)
	// #nosec G204
	cmdCompile := execCommand("javac", "--release", "8", javaFile)
	if _, err = cmdCompile.CombinedOutput(); err != nil {
		// Fallback for older javac that doesn't support --release
		// #nosec G204
		cmdFallback := execCommand("javac", "-source", "1.8", "-target", "1.8", javaFile)
		if out2, err2 := cmdFallback.CombinedOutput(); err2 != nil {
			return fmt.Errorf("javac failed to compile for Java 8 compatibility: %s", string(out2))
		}
	}

	// 3. Create Manifest
	manifest := "Agent-Class: ProxyAgent\nCan-Retransform-Classes: true\n"
	err = os.WriteFile(manifestFile, []byte(manifest), 0600)
	if err != nil {
		return err
	}

	// 4. Create JAR (include all ProxyAgent*.class files — anonymous inner classes compile to ProxyAgent$1.class, etc.)
	jarArgs := []string{"cmf", manifestFile, jarFile}
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read temp dir: %v", err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "ProxyAgent") && strings.HasSuffix(e.Name(), ".class") {
			jarArgs = append(jarArgs, "-C", tmpDir, e.Name())
		}
	}
	// #nosec G204
	if out, err := execCommand("jar", jarArgs...).CombinedOutput(); err != nil {
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
	if _, err := execCommand("javac", "--add-modules", "jdk.attach", injectorFile).CombinedOutput(); err != nil {
		// Fallback for Java 8
		if toolsJar != "" {
			// #nosec G204
			if out, err := execCommand("javac", "-cp", toolsJar, injectorFile).CombinedOutput(); err != nil {
				log.Printf("Warning: Injector compilation with tools.jar failed: %s", string(out))
			}
		}
	}

	// 6. Attempt attachment
	host, port := splitAddr(proxyAddr)
	agentArgs := fmt.Sprintf("%s:%s", host, port)

	// Try the Injector (Java 9+)
	// #nosec G204
	cmd := execCommand("java", "--add-modules", "jdk.attach", "-cp", tmpDir, "Injector", pid, jarFile, agentArgs)
	if out, err := cmd.CombinedOutput(); err == nil && strings.Contains(string(out), "Successfully attached") {
		return nil
	} else if toolsJar != "" {
		// Try Java 8 fallback for Injector
		// #nosec G204
		cmd8 := execCommand("java", "-cp", tmpDir+":"+toolsJar, "Injector", pid, jarFile, agentArgs)
		if out8, err8 := cmd8.CombinedOutput(); err8 == nil && strings.Contains(string(out8), "Successfully attached") {
			return nil
		}
	}

	// 7. Last resort: jcmd
	// #nosec G204
	cmdLast := execCommand("jcmd", pid, "JVMTI.agent_load", jarFile, agentArgs)
	if _, err := cmdLast.CombinedOutput(); err != nil {
		// #nosec G204
		cmdLast2 := execCommand("jcmd", pid, "JVMTI.agent_load", jarFile+"="+agentArgs)
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
