# Docker Container Interception

Glance provides a seamless, one-click interception method for Docker containers. This allows you to inspect and modify traffic from any containerized application without manual configuration or restarts.

## How Interception Works

Glance uses a **Container Recreation** strategy to ensure that all processes within a container correctly pick up proxy settings. When you click **"Intercept Traffic"** in the dashboard, Glance performs the following steps:

1.  **Configuration Capture**: Glance inspects the running container to read its current environment variables, volume mounts, network settings, and labels.
2.  **Container Backup**: The original container is stopped and renamed to `{original_name}-glance-backup`.
3.  **Environment Injection**: Glance creates a new container with the same name and configuration as the original, but injects the following environment variables:
    *   `HTTP_PROXY` and `HTTPS_PROXY`: Pointed to the host's Glance proxy.
    *   `JAVA_TOOL_OPTIONS`: Injected with JVM system properties (`-Dhttp.proxyHost`, etc.) to ensure Java applications are intercepted.
    *   `glance.interception=active`: A label used to track the interception state.
4.  **CA Certificate Injection**: Glance automatically detects the container's OS (Alpine, Debian, or RHEL) and injects its CA certificate into the system trust store and the Java Trust Store (`cacerts`).
5.  **Start**: The new intercepted container is started, and traffic begins flowing through Glance.

When you click **"Stop Intercept"**, Glance removes the intercepted container and restores the original backup to its original name and state.

## Features

### Universal Compatibility
By using environment variables (`HTTP_PROXY`) and JVM options (`JAVA_TOOL_OPTIONS`), Glance can intercept traffic from almost any runtime, including:
*   **Java / Spring Boot / JVM**
*   **Node.js / Python / Go**
*   **CLI tools like cURL and wget**

### Automatic HTTPS Decryption
Glance handles the complexity of SSL trust for you:
*   **System Trust**: Automatically runs `update-ca-certificates` or `update-ca-trust` inside the container.
*   **Java KeyStore**: Automatically locates the `cacerts` file and uses `keytool` to import the Glance CA, preventing `SSLHandshakeException` in JVM apps.

### No Special Capabilities Required
Unlike `iptables`-based methods, this approach does **not** require the `--cap-add=NET_ADMIN` capability, making it compatible with more restricted container environments.

## Smart Host Detection

In Docker environments, reaching the host's "localhost" from inside a container can be tricky. Glance automatically handles this:
*   **Automatic Host Mapping**: Glance injects `host.docker.internal:host-gateway` into the container's `extra_hosts` configuration. This allows the container to resolve `host.docker.internal` to the host machine even on standard Linux environments where this is not provided by default.
*   **Docker Desktop (macOS/Windows)**: Uses `host.docker.internal` to route traffic to the host.
*   **Linux**: Automatically detects the default bridge gateway IP and uses it as the proxy host.

## Usage Guide

1.  Open the **Glance Dashboard** (`http://localhost:15501`).
2.  Navigate to the **Integrations** tab.
3.  Find your container in the **Docker Containers** list.
4.  Click **"Intercept Traffic"**.
5.  The container will briefly restart with the new proxy settings.
6.  Start making requests from your containerized app—they will now appear in the **Traffic** tab.

## Troubleshooting

### Container Fails to Start
*   Ensure Glance is running and reachable from your Docker network.
*   Check the container logs using `docker logs {container_name}` for any proxy-related startup errors.

### HTTPS Traffic Not Intercepted
*   If your application uses a custom keystore or hardcoded certificates, automatic injection may not work. You may need to manually trust the Glance CA within your application code.
*   Ensure the container has the `keytool` utility installed if it's a Java application.

### Container Name Conflicts
*   If a previous interception session was interrupted, you might have a container named `{name}-glance-backup`. You may need to manually remove it or the intercepted container to resolve the conflict.

## Next Steps

- [Java/JVM Guide](java.md) - More details on Java-specific interception
- [Troubleshooting](../troubleshooting.md) - Common issues and solutions
- [Dashboard Overview](../features/traffic-inspection.md) - Learn how to inspect the captured traffic
