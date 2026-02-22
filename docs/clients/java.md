# Java/JVM Configuration

Glance provides native support for Java applications with automatic `HttpsURLConnection` overrides.

## Automatic Setup (Java 8+)

Glance includes a specialized Java agent that automatically configures `HttpsURLConnection` to:
- Trust Glance's CA certificate
- Route traffic through the proxy
- Work with already-running applications

### How It Works

1. Glance detects running Java processes
2. Attaches bytecode instrumentation
3. Overrides `HttpsURLConnection` behavior
4. No code changes required

## Manual Configuration

If automatic setup doesn't work, you can manually configure the JVM.

### Method 1: JVM Arguments

Start your Java application with these arguments:

```bash
java -Dhttp.proxyHost=localhost \
     -Dhttp.proxyPort=15500 \
     -Dhttps.proxyHost=localhost \
     -Dhttps.proxyPort=15500 \
     -jar your-application.jar
```

### Method 2: System Properties

Set properties programmatically:

```java
System.setProperty("http.proxyHost", "localhost");
System.setProperty("http.proxyPort", "15500");
System.setProperty("https.proxyHost", "localhost");
System.setProperty("https.proxyPort", "15500");
```

### Method 3: Proxy Class

Use Java's `Proxy` class:

```java
import java.net.*;

Proxy proxy = new Proxy(
    Proxy.Type.HTTP,
    new InetSocketAddress("localhost", 15500)
);

URL url = new URL("https://api.example.com/users");
HttpURLConnection conn = (HttpURLConnection) url.openConnection(proxy);
```

## Certificate Trust

### Option 1: Import to JVM Keystore

```bash
# Export Glance CA certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Import to JVM keystore (replace JDK_PATH)
keytool -import -trustcacerts -alias glance \
        -file glance-ca.crt \
        -keystore $JDK_PATH/lib/security/cacerts \
        -storepass changeit
```

### Option 2: Custom TrustManager

Create a custom `TrustManager` that trusts Glance's certificate:

```java
import javax.net.ssl.*;
import java.security.cert.*;

// For development only - do not use in production!
TrustManager[] trustAllCerts = new TrustManager[] {
    new X509TrustManager() {
        public X509Certificate[] getAcceptedIssuers() { return null; }
        public void checkClientTrusted(X509Certificate[] certs, String authType) { }
        public void checkServerTrusted(X509Certificate[] certs, String authType) { }
    }
};

SSLContext sc = SSLContext.getInstance("SSL");
sc.init(null, trustAllCerts, new java.security.SecureRandom());
HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
```

⚠️ **Warning**: Only use this in development. Never disable certificate validation in production!

## Framework-Specific Configuration

### Spring Boot

Add to `application.properties`:

```properties
http.proxy.host=localhost
http.proxy.port=15500
https.proxy.host=localhost
https.proxy.port=15500
```

Or configure via `RestTemplate`:

```java
@Configuration
public class ProxyConfig {
    @Bean
    public RestTemplate restTemplate() {
        SimpleClientHttpRequestFactory requestFactory =
            new SimpleClientHttpRequestFactory();

        Proxy proxy = new Proxy(
            Proxy.Type.HTTP,
            new InetSocketAddress("localhost", 15500)
        );
        requestFactory.setProxy(proxy);

        return new RestTemplate(requestFactory);
    }
}
```

### Apache HttpClient

```java
import org.apache.http.HttpHost;
import org.apache.http.impl.client.HttpClients;

HttpHost proxy = new HttpHost("localhost", 15500);
CloseableHttpClient httpClient = HttpClients.custom()
    .setProxy(proxy)
    .build();
```

### OkHttp

```java
import okhttp3.*;

Proxy proxy = new Proxy(
    Proxy.Type.HTTP,
    new InetSocketAddress("localhost", 15500)
);

OkHttpClient client = new OkHttpClient.Builder()
    .proxy(proxy)
    .build();
```

## Testing

Verify your Java app is using the proxy:

```java
import java.net.*;

public class ProxyTest {
    public static void main(String[] args) throws Exception {
        URL url = new URL("https://api.github.com/users");
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();

        int responseCode = conn.getResponseCode();
        System.out.println("Response Code: " + responseCode);

        // This request should appear in Glance dashboard
    }
}
```

## Troubleshooting

### Certificate Errors

If you see `SSLHandshakeException`:
1. Verify CA certificate is imported
2. Check certificate alias exists
3. Ensure using correct keystore

### Proxy Not Used

If traffic doesn't appear in Glance:
1. Check proxy properties are set
2. Verify no other proxy configured
3. Some libraries ignore system properties

### ClassLoader Issues

If automatic agent fails:
1. Check Java version (8+ required)
2. Try manual configuration instead
3. Check for SecurityManager restrictions

## Next Steps

- [Client Configuration](/clients.md) - Other platforms
- [Troubleshooting](/troubleshooting.md) - Common issues
- [Examples](https://github.com/wahyudotdev/glance/tree/main/examples/java) - Sample code
