# Chromium Browser Configuration

Glance supports traffic interception for Chromium-based browsers via standard command-line flags.

## Supported Browsers

- Google Chrome
- Microsoft Edge
- Brave Browser
- Chromium

## One-Click Launch (Dashboard)

The easiest way to use Glance with Chromium browsers is via the dashboard:

1. Open Glance Dashboard at `http://localhost:15501`
2. Go to the **Integrations** tab
3. Click **"Launch Browser"** under the Chromium / Chrome section

This will automatically:
- Launch your default Chromium browser
- Configure it to use Glance proxy
- Ignore certificate errors (for development convenience)
- Use a temporary, isolated user profile (doesn't affect your main browser)

## Manual Configuration

You can launch your browser with the following flags to route traffic through Glance:

### Chrome/Chromium

```bash
# macOS
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  --proxy-server="http://localhost:15500" \
  --ignore-certificate-errors \
  --user-data-dir="/tmp/glance-chrome"

# Linux
google-chrome \
  --proxy-server="http://localhost:15500" \
  --ignore-certificate-errors \
  --user-data-dir="/tmp/glance-chrome"

# Windows
"C:\Program Files\Google\Chrome\Application\chrome.exe" ^
  --proxy-server="http://localhost:15500" ^
  --ignore-certificate-errors ^
  --user-data-dir="C:\Temp\glance-chrome"
```

### Edge

```bash
# macOS
"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge" \
  --proxy-server="http://localhost:15500" \
  --ignore-certificate-errors \
  --user-data-dir="/tmp/glance-edge"

# Windows
"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" ^
  --proxy-server="http://localhost:15500" ^
  --ignore-certificate-errors ^
  --user-data-dir="C:\Temp\glance-edge"
```

> **Note**: We use `--user-data-dir` to create a temporary, isolated profile. This ensures your main browser settings and extensions don't interfere with the proxy, and prevents "Profile in use" errors.

## System Proxy (All Apps)

To route all browser traffic through Glance without command-line flags:

### macOS

1. **System Preferences** → **Network**
2. Select your active network
3. **Advanced** → **Proxies**
4. Enable **Web Proxy (HTTP)** and **Secure Web Proxy (HTTPS)**
5. Set both to:
   - **Server**: `localhost`
   - **Port**: `15500`
6. Click **OK** → **Apply**

### Windows

1. **Settings** → **Network & Internet** → **Proxy**
2. Enable **Use a proxy server**
3. Set:
   - **Address**: `localhost`
   - **Port**: `15500`
4. Save

### Linux

```bash
# GNOME
gsettings set org.gnome.system.proxy mode 'manual'
gsettings set org.gnome.system.proxy.http host 'localhost'
gsettings set org.gnome.system.proxy.http port 15500
gsettings set org.gnome.system.proxy.https host 'localhost'
gsettings set org.gnome.system.proxy.https port 15500

# Or use environment variables
export HTTP_PROXY=http://localhost:15500
export HTTPS_PROXY=http://localhost:15500
```

## Certificate Trust

For production use (not recommended for development), you can import Glance's CA certificate:

### Chrome Certificate Manager

1. Open Chrome → **Settings** → **Privacy and Security** → **Security**
2. Click **Manage Certificates**
3. Go to **Authorities** tab
4. Click **Import**
5. Select the Glance CA certificate from `http://localhost:15501/ca.crt`
6. Check "Trust this certificate for identifying websites"

However, for development, using `--ignore-certificate-errors` is easier and doesn't require import.

## Extension Development

For testing browser extensions:

```bash
chrome --proxy-server="localhost:15500" \
       --ignore-certificate-errors \
       --load-extension="/path/to/extension" \
       --user-data-dir="/tmp/glance-chrome"
```

## DevTools Integration

Glance works alongside Chrome DevTools:

1. Open DevTools (F12 or Cmd+Option+I)
2. Go to **Network** tab
3. You'll see requests in both:
   - Chrome DevTools (client-side view)
   - Glance Dashboard (proxy-side view)

This gives you two perspectives:
- **DevTools**: Browser's view, timing, caching
- **Glance**: Network view, headers, ability to mock

## Headless Mode

For automated testing:

```bash
chrome --headless \
       --proxy-server="localhost:15500" \
       --ignore-certificate-errors \
       --dump-dom https://example.com
```

Useful for:
- CI/CD testing
- Automated screenshot capture
- Web scraping with proxy

## Puppeteer Integration

If you're using Puppeteer:

```javascript
const puppeteer = require('puppeteer');

const browser = await puppeteer.launch({
  headless: false,
  args: [
    '--proxy-server=localhost:15500',
    '--ignore-certificate-errors',
  ]
});

const page = await browser.newPage();
await page.goto('https://example.com');

// All requests will go through Glance
```

## Playwright Integration

For Playwright:

```javascript
const { chromium } = require('playwright');

const browser = await chromium.launch({
  headless: false,
  proxy: {
    server: 'http://localhost:15500'
  },
  ignoreHTTPSErrors: true
});

const context = await browser.newContext();
const page = await context.newPage();
await page.goto('https://example.com');
```

## Selenium WebDriver

For Selenium:

```python
from selenium import webdriver
from selenium.webdriver.chrome.options import Options

chrome_options = Options()
chrome_options.add_argument('--proxy-server=localhost:15500')
chrome_options.add_argument('--ignore-certificate-errors')

driver = webdriver.Chrome(options=chrome_options)
driver.get('https://example.com')
```

## Testing

Verify browser is using proxy:

1. Launch browser with proxy configured
2. Navigate to any HTTPS website
3. Check Glance dashboard - requests should appear
4. Verify URL and status code match

## Troubleshooting

### Proxy Not Working

- **Check Browser Flags**: Ensure `--proxy-server` is set correctly
- **Port Conflict**: Verify Glance is running on port 15500
- **System Proxy**: Disable system proxy to avoid conflicts

### Certificate Errors Still Showing

- **Flag Missing**: Ensure `--ignore-certificate-errors` is present
- **Extension Conflict**: Some security extensions may block
- **HSTS**: Sites with HSTS may still show warnings

### User Data Directory Issues

- **Permission Denied**: Use a writable directory like `/tmp/`
- **Directory in Use**: Close other browser instances
- **Disk Space**: Ensure enough space for profile

### Performance Issues

- **Disable Extensions**: Use clean profile with `--user-data-dir`
- **GPU Acceleration**: Add `--disable-gpu` if having issues
- **Memory**: Chrome with DevTools + Glance can use significant RAM

## Production vs Development

| Use Case | Recommended Approach |
|----------|---------------------|
| **Development** | Auto-launch with `--ignore-certificate-errors` |
| **Testing** | System proxy + imported certificate |
| **CI/CD** | Headless with proxy flags |
| **Production** | Never use MITM proxy |

## Best Practices

1. **Use Separate Profile**: Always use `--user-data-dir` to avoid affecting main browser
2. **Don't Sign In**: Don't sign into Google/Microsoft accounts in proxy-configured browser
3. **DevTools Open**: Keep DevTools open for debugging
4. **Clear State**: Delete user data directory between tests

## Next Steps

- [Client Configuration](/clients.md) - Other platforms
- [Mocking & Breakpoints](/features/mocking.md) - Modify browser traffic
- [Scenarios](/features/scenarios.md) - Record browser workflows
