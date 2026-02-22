# Android Configuration

Glance provides ADB-based support for intercepting traffic from Android devices and emulators.

## Prerequisites

- ADB (Android Debug Bridge) installed
- USB debugging enabled on device
- Android 7.0 (API 24) or higher

## Quick Setup

Glance can automatically configure your Android device:

1. Connect device via USB or start emulator
2. Run Glance with Android support:
   ```bash
   glance --android
   ```
3. Glance will:
   - Detect connected devices
   - Install CA certificate
   - Configure proxy settings
   - Set up port forwarding

## Manual Setup

### Step 1: Enable USB Debugging

On your Android device:
1. Go to **Settings** → **About Phone**
2. Tap **Build Number** 7 times to enable Developer Options
3. Go to **Settings** → **Developer Options**
4. Enable **USB Debugging**

### Step 2: Connect Device

```bash
# List connected devices
adb devices

# Should show:
# List of devices attached
# ABC123456    device
```

### Step 3: Install CA Certificate

```bash
# Export Glance CA certificate
curl http://localhost:15501/ca.crt -o glance-ca.crt

# Push to device
adb push glance-ca.crt /sdcard/

# Install certificate
# On device: Settings → Security → Install from storage
# Select glance-ca.crt
```

### Step 4: Configure Proxy

#### Method 1: ADB Reverse

```bash
# Forward device's requests to Glance proxy
adb reverse tcp:15500 tcp:15500
```

#### Method 2: WiFi Proxy (Android 10+)

1. Go to **Settings** → **Network & Internet** → **WiFi**
2. Long press your network → **Modify Network**
3. **Advanced Options** → **Proxy** → **Manual**
4. Set:
   - **Proxy hostname**: `localhost` (if using adb reverse) or your computer's IP
   - **Proxy port**: `15500`
5. Save

#### Method 3: Global Proxy (Requires Root)

```bash
# Set global proxy
adb shell settings put global http_proxy localhost:15500

# Remove proxy
adb shell settings put global http_proxy :0
```

## Certificate Installation

### User Certificate (Android 7-10)

For apps targeting API 24-29, user certificates work by default:

1. Download CA certificate on device
2. **Settings** → **Security** → **Install from storage**
3. Select certificate file
4. Name it "Glance CA"
5. Select "VPN and apps" usage

### System Certificate (Android 11+)

Apps targeting API 30+ only trust system certificates. Two options:

#### Option 1: Network Security Config (Developers)

Add to your app's `network_security_config.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
    <debug-overrides>
        <trust-anchors>
            <!-- Trust user certificates for debugging -->
            <certificates src="user" />
            <certificates src="system" />
        </trust-anchors>
    </debug-overrides>
</network-security-config>
```

Reference in `AndroidManifest.xml`:

```xml
<application
    android:networkSecurityConfig="@xml/network_security_config">
</application>
```

#### Option 2: Install as System Certificate (Requires Root)

```bash
# Get certificate hash
export CERT_HASH=$(openssl x509 -inform PEM -subject_hash_old -in glance-ca.crt | head -1)

# Remount system as writable
adb root
adb remount

# Push certificate
adb push glance-ca.crt /system/etc/security/cacerts/$CERT_HASH.0

# Set permissions
adb shell chmod 644 /system/etc/security/cacerts/$CERT_HASH.0

# Reboot
adb reboot
```

## Emulator Configuration

### Android Studio Emulator

Emulators are easier to configure:

```bash
# Start emulator with writable system
emulator -avd YOUR_AVD_NAME -writable-system

# Or configure proxy at startup
emulator -avd YOUR_AVD_NAME -http-proxy localhost:15500
```

### Certificate for Emulator

```bash
# Emulator with root access
adb root
adb remount

# Push certificate to system
curl http://localhost:15501/ca.crt -o glance-ca.crt
adb push glance-ca.crt /system/etc/security/cacerts/$(openssl x509 -inform PEM -subject_hash_old -in glance-ca.crt | head -1).0

# Restart
adb reboot
```

## Testing

Verify configuration works:

```bash
# On device, any app should route through proxy
# Open Chrome and visit any HTTPS site
# Traffic should appear in Glance dashboard
```

Or use ADB:

```bash
# Make request from device
adb shell curl https://api.github.com/users

# Should appear in Glance
```

## Framework-Specific Configuration

### OkHttp

```kotlin
import okhttp3.*
import java.net.InetSocketAddress
import java.net.Proxy

val proxy = Proxy(Proxy.Type.HTTP, InetSocketAddress("localhost", 15500))

val client = OkHttpClient.Builder()
    .proxy(proxy)
    .build()
```

### Retrofit

```kotlin
val proxy = Proxy(Proxy.Type.HTTP, InetSocketAddress("localhost", 15500))

val okHttpClient = OkHttpClient.Builder()
    .proxy(proxy)
    .build()

val retrofit = Retrofit.Builder()
    .baseUrl("https://api.example.com")
    .client(okHttpClient)
    .build()
```

### Volley

```kotlin
val proxy = Proxy(Proxy.Type.HTTP, InetSocketAddress("localhost", 15500))

// Volley doesn't support proxy directly
// Use OkHttp or configure system-wide
```

## Troubleshooting

### Certificate Not Trusted

- **Check API Level**: API 30+ requires system certificate or network security config
- **Verify Installation**: Settings → Security → Trusted Credentials → User
- **App-Specific**: Some apps may have certificate pinning

### No Traffic Visible

- **Verify Proxy**: Check WiFi proxy settings
- **ADB Reverse**: Ensure `adb reverse` command succeeded
- **Port Conflict**: Check port 15500 isn't used on device

### HTTPS Errors

- **Certificate Mismatch**: Ensure using latest Glance CA
- **Pinning**: Some apps use certificate pinning (can't be intercepted)
- **Old Cache**: Clear app cache and restart

### Certificate Pinning

Some apps (banking, security-focused) use certificate pinning and can't be intercepted:

- **Root + Frida**: Advanced users can bypass pinning
- **Not Recommended**: Bypassing pinning may violate ToS
- **Alternative**: Use API documentation instead

## Production Apps vs Debug Builds

| Type | User Cert | System Cert | Network Config |
|------|-----------|-------------|----------------|
| **Debug Build** | ✅ Yes | ✅ Yes | ✅ Configurable |
| **Release Build** | ❌ No (API 30+) | ✅ Yes | ⚠️ If included |
| **3rd Party App** | ❌ No (API 30+) | ✅ Yes | ❌ Can't modify |

## Best Practices

1. **Use Emulator**: Easier to configure and reset
2. **Debug Builds**: Always test with debug builds during development
3. **Root When Needed**: Consider rooted emulator for testing production apps
4. **Network Security Config**: Always include for debug builds

## Next Steps

- [Client Configuration](/clients.md) - Other platforms
- [Troubleshooting](/troubleshooting.md) - Common issues
- [MCP Integration](/mcp/) - Analyze mobile traffic with AI
