import React, { useState } from 'react';
import { Globe, Terminal, Code, Activity, Copy, Check, Shield, Smartphone, ChevronRight, XCircle, QrCode } from 'lucide-react';
import { QRCodeSVG } from 'qrcode.react';
import type { JavaProcess, AndroidDevice } from '../../types/traffic';

interface IntegrationsViewProps {
  javaProcesses: JavaProcess[];
  androidDevices: AndroidDevice[];
  isLoadingJava: boolean;
  isLoadingAndroid: boolean;
  terminalScript: string;
  onFetchJava: () => void;
  onFetchAndroid: () => void;
  onInterceptJava: (pid: string) => void;
  onInterceptAndroid: (id: string) => void;
  onClearAndroid: (id: string) => void;
  onPushAndroidCert: (id: string) => void;
}

export const IntegrationsView: React.FC<IntegrationsViewProps> = ({ 
  javaProcesses, androidDevices, isLoadingJava, isLoadingAndroid, terminalScript, 
  onFetchJava, onFetchAndroid, onInterceptJava, onInterceptAndroid, onClearAndroid,
  onPushAndroidCert
}) => {
  const [scriptCopied, setScriptCopied] = useState(false);

  return (
    <div className="flex-1 p-12 bg-slate-50 overflow-y-auto">
      <div className="max-w-4xl mx-auto space-y-12">
        <section>
          <h2 className="text-2xl font-bold text-slate-800 mb-6">Client Integrations</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
              <div className="w-12 h-12 bg-blue-50 rounded-xl flex items-center justify-center mb-6">
                <Globe className="text-blue-600" size={24} />
              </div>
              <h3 className="text-lg font-bold mb-2">Chromium / Chrome</h3>
              <p className="text-sm text-slate-500 mb-6 leading-relaxed">
                Launch a fresh browser instance pre-configured to route all traffic through this proxy and ignore certificate errors.
              </p>
              <button 
                onClick={async () => {
                  try { await fetch('/api/client/chromium', { method: 'POST' }); }
                  catch (e) { alert('Failed to launch Chromium: ' + e); }
                }}
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-bold text-sm hover:bg-blue-700 active:scale-95 transition-all shadow-lg shadow-blue-200"
              >
                Launch Browser
              </button>
            </div>

            <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
              <div className="flex items-start justify-between mb-6">
                <div className="w-12 h-12 bg-emerald-50 rounded-xl flex items-center justify-center">
                  <Smartphone className="text-emerald-600" size={24} />
                </div>
                <button 
                  onClick={onFetchAndroid}
                  className={`p-2 hover:bg-slate-100 rounded-lg transition-all ${isLoadingAndroid ? 'animate-spin text-emerald-600' : 'text-slate-400'}`}
                  title="Refresh Devices"
                >
                  <Activity size={16} />
                </button>
              </div>
              <h3 className="text-lg font-bold mb-2">Android (ADB)</h3>
              <p className="text-sm text-slate-500 mb-6 leading-relaxed">
                Connect devices via USB/WiFi to intercept traffic. Requires 'adb' in PATH.
              </p>
              
              <div className="space-y-4">
                <div className="bg-slate-50 rounded-xl border border-slate-100 overflow-hidden min-h-[100px]">
                  {androidDevices && androidDevices.length > 0 ? (
                    <div className="divide-y divide-slate-100 max-h-48 overflow-y-auto">
                      {androidDevices.map(device => (
                        <div key={device.id} className="px-4 py-3 flex flex-col gap-3 hover:bg-white transition-colors">
                          <div className="flex items-center justify-between">
                            <div className="flex flex-col">
                              <span className="text-xs font-bold text-slate-700 font-mono">{device.model}</span>
                              <span className="text-[10px] text-slate-400 font-mono">{device.id}</span>
                            </div>
                            <div className="flex gap-2">
                              <button 
                                onClick={() => onPushAndroidCert(device.id)}
                                className="px-3 py-1.5 bg-blue-50 border border-blue-100 text-blue-600 rounded-lg text-[10px] font-bold hover:bg-blue-600 hover:text-white transition-all active:scale-95"
                                title="Push CA Certificate to Device"
                              >
                                Push CA
                              </button>
                              <button 
                                onClick={() => onInterceptAndroid(device.id)}
                                className="px-3 py-1.5 bg-emerald-500 text-white rounded-lg text-[10px] font-bold shadow-sm hover:bg-emerald-600 transition-all active:scale-95"
                              >
                                Intercept
                              </button>
                              <button 
                                onClick={() => onClearAndroid(device.id)}
                                className="px-3 py-1.5 bg-white border border-slate-200 text-slate-500 rounded-lg text-[10px] font-bold hover:text-rose-500 hover:border-rose-200 transition-all"
                              >
                                Reset
                              </button>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="flex flex-col items-center justify-center h-full py-8 gap-2">
                      <Smartphone size={24} className="text-slate-300" />
                      <span className="text-xs text-slate-400 italic">No devices detected via ADB</span>
                    </div>
                  )}
                </div>
                
                <div className="bg-emerald-50/50 border border-emerald-100 rounded-xl p-3">
                  <p className="text-[10px] text-emerald-800 leading-relaxed">
                    <strong>Note:</strong> Interception uses <code>adb reverse</code>. Ensure your device can connect to the host machine.
                  </p>
                </div>

                <div className="space-y-3">
                  <details className="group border border-slate-200 rounded-xl overflow-hidden bg-white">
                    <summary className="px-4 py-3 text-xs font-bold text-slate-700 cursor-pointer hover:bg-slate-50 flex items-center justify-between list-none">
                      <div className="flex items-center gap-2">
                        <Shield size={14} className="text-blue-500" />
                        <span>Trust Guide: Option A (Rooted / Emulator)</span>
                      </div>
                      <ChevronRight size={14} className="group-open:rotate-90 transition-transform" />
                    </summary>
                    <div className="px-4 py-4 bg-slate-50 border-t border-slate-100 space-y-3">
                      <p className="text-[11px] text-slate-600">Run these commands to move the CA to the System Store (requires root & writable /system):</p>
                      <pre className="bg-slate-900 text-blue-300 p-3 rounded-lg text-[9px] font-mono whitespace-pre-wrap leading-relaxed">
                        # Download CA first, then in your terminal:{"\n"}
                        HASH=$(openssl x509 -inform PEM -subject_hash_old -in agent-proxy-ca.crt | head -1){"\n"}
                        adb push agent-proxy-ca.crt /sdcard/$HASH.0{"\n"}
                        adb shell "su -c 'mount -o rw,remount /system && cp /sdcard/$HASH.0 /system/etc/security/cacerts/ && chmod 644 /system/etc/security/cacerts/$HASH.0'"
                      </pre>
                    </div>
                  </details>

                  <details className="group border border-slate-200 rounded-xl overflow-hidden bg-white">
                    <summary className="px-4 py-3 text-xs font-bold text-slate-700 cursor-pointer hover:bg-slate-50 flex items-center justify-between list-none">
                      <div className="flex items-center gap-2">
                        <Code size={14} className="text-amber-500" />
                        <span>Trust Guide: Option B (App Development)</span>
                      </div>
                      <ChevronRight size={14} className="group-open:rotate-90 transition-transform" />
                    </summary>
                    <div className="px-4 py-4 bg-slate-50 border-t border-slate-100 space-y-3">
                      <p className="text-[11px] text-slate-600">Add this to your app's <code>res/xml/network_security_config.xml</code>:</p>
                      <pre className="bg-slate-900 text-amber-200 p-3 rounded-lg text-[9px] font-mono whitespace-pre-wrap leading-relaxed">
                        &lt;network-security-config&gt;{"\n"}
                        {"  "}&lt;base-config&gt;{"\n"}
                        {"    "}&lt;trust-anchors&gt;{"\n"}
                        {"      "}&lt;certificates src="system" /&gt;{"\n"}
                        {"      "}&lt;certificates src="user" /&gt;{"\n"}
                        {"    "}&lt;/trust-anchors&gt;{"\n"}
                        {"  "}&lt;/base-config&gt;{"\n"}
                        &lt;/network-security-config&gt;
                      </pre>
                      <p className="text-[11px] text-slate-600 italic">And reference it in <code>AndroidManifest.xml</code> under &lt;application android:networkSecurityConfig="..."&gt;</p>
                    </div>
                  </details>

                  <details className="group border border-slate-200 rounded-xl overflow-hidden bg-white">
                    <summary className="px-4 py-3 text-xs font-bold text-slate-700 cursor-pointer hover:bg-slate-50 flex items-center justify-between list-none">
                      <div className="flex items-center gap-2">
                        <XCircle size={14} className="text-rose-500" />
                        <span>Handshake Still Failing? (Troubleshooting)</span>
                      </div>
                      <ChevronRight size={14} className="group-open:rotate-90 transition-transform" />
                    </summary>
                    <div className="px-4 py-4 bg-slate-50 border-t border-slate-100 space-y-3">
                      <ul className="text-[11px] text-slate-600 list-disc ml-4 space-y-2">
                        <li><strong>Check Manifest:</strong> Ensure <code>android:networkSecurityConfig="@xml/network_security_config"</code> is inside the <code>&lt;application&gt;</code> tag in <code>AndroidManifest.xml</code>.</li>
                        <li><strong>Verify Installation:</strong> Go to Settings &rarr; Security &rarr; User Credentials. Ensure "Agent Proxy CA" is listed.</li>
                        <li><strong>System Time:</strong> Ensure the Android device's date/time is correct. If the device time is in the past, the CA cert will be rejected.</li>
                        <li><strong>SSL Pinning:</strong> If the app uses Certificate Pinning (common in apps like Facebook, banking, or some custom OkHttp setups), standard trust configurations <strong>will not work</strong>. You must disable pinning in the source code or use a tool like <strong>Frida</strong> or <strong>Xposed</strong> to bypass it.</li>
                        <li><strong>Restart App:</strong> Fully kill and restart the Android app after applying any security config changes.</li>
                      </ul>
                    </div>
                  </details>
                </div>
              </div>
            </div>

            <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
              <div className="w-12 h-12 bg-indigo-50 rounded-xl flex items-center justify-center mb-6">
                <Terminal className="text-indigo-600" size={24} />
              </div>
              <h3 className="text-lg font-bold mb-2">Existing Terminal</h3>
              <p className="text-sm text-slate-500 mb-6 leading-relaxed">
                Run this one-liner in any terminal to instantly enable interception.
              </p>
              <div className="relative group mb-4">
                <pre className="bg-slate-900 text-indigo-200 p-4 rounded-xl text-[10px] font-mono overflow-x-auto">
                  eval "$(curl -s {window.location.origin}/setup)"
                </pre>
                <button 
                  onClick={() => {
                    navigator.clipboard.writeText(`eval "$(curl -s ${window.location.origin}/setup)"`);
                    setScriptCopied(true);
                    setTimeout(() => setScriptCopied(false), 2000);
                  }}
                  className="absolute top-2 right-2 p-2 bg-slate-800 text-slate-400 hover:text-white rounded-lg transition-all"
                >
                  {scriptCopied ? <Check size={14} className="text-emerald-500" /> : <Copy size={14} />}
                </button>
              </div>
              <details className="text-[10px] text-slate-400 cursor-pointer">
                <summary className="hover:text-slate-600 transition-colors">Alternative: Manual Setup</summary>
                <div className="mt-2 relative group">
                  <pre className="bg-slate-900 text-indigo-200 p-4 rounded-xl text-[9px] font-mono overflow-x-auto max-h-32">
                    {terminalScript || '# Fetching setup script...'}
                  </pre>
                </div>
              </details>
            </div>

            <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow md:col-span-2">
              <div className="flex items-start justify-between mb-6">
                <div className="flex gap-4">
                  <div className="w-12 h-12 bg-amber-50 rounded-xl flex items-center justify-center">
                    <Code className="text-amber-600" size={24} />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold">CA Certificate</h3>
                    <p className="text-sm text-slate-500">Required for HTTPS interception</p>
                  </div>
                </div>
                <div className="flex gap-3">
                  <details className="relative">
                    <summary className="flex items-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-lg text-xs font-bold text-slate-600 transition-all cursor-pointer list-none">
                      <QrCode size={14} />
                      Scan QR
                    </summary>
                    <div className="absolute right-0 mt-2 p-4 bg-white border border-slate-200 rounded-2xl shadow-xl z-50 animate-in fade-in zoom-in-95 duration-200">
                      <div className="bg-white p-2 rounded-xl border border-slate-100 mb-2">
                        <QRCodeSVG 
                          value={window.location.href.replace(window.location.pathname, '') + '/api/ca/cert'}
                          size={160}
                          level="H"
                          includeMargin={true}
                        />
                      </div>
                      <p className="text-[10px] text-center text-slate-400 font-medium">Scan to download CA cert{"\n"}directly on mobile</p>
                    </div>
                  </details>
                  <a 
                    href="/api/ca/cert" 
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-xs font-bold text-white transition-all shadow-lg shadow-blue-200"
                  >
                    Download .crt
                  </a>
                </div>
              </div>
              <h3 className="text-lg font-bold mb-2">Java / JVM Applications</h3>
              <p className="text-sm text-slate-500 mb-6 leading-relaxed">
                Detect running Java applications and get interception instructions.
              </p>
              
              <div className="space-y-6">
                <div className="bg-slate-50 rounded-xl border border-slate-100 overflow-hidden">
                  <div className="px-4 py-3 border-b border-slate-200 bg-white flex items-center justify-between">
                    <span className="text-[10px] font-black uppercase text-slate-400 tracking-widest">Running Java Processes</span>
                    <button 
                      onClick={onFetchJava}
                      className={`p-1.5 hover:bg-slate-100 rounded-md transition-all ${isLoadingJava ? 'animate-spin text-blue-600' : 'text-slate-400'}`}
                    >
                      <Activity size={14} />
                    </button>
                  </div>
                  <div className="divide-y divide-slate-100 max-h-48 overflow-y-auto">
                    {javaProcesses.length > 0 ? javaProcesses.map(proc => (
                      <div key={proc.pid} className="px-4 py-3 flex items-center justify-between hover:bg-blue-50/50 transition-colors group cursor-default">
                        <div className="flex flex-col">
                          <span className="text-xs font-bold text-slate-700 font-mono group-hover:text-blue-700 transition-colors">{proc.name}</span>
                          <span className="text-[10px] text-slate-400 font-mono">PID: {proc.pid}</span>
                        </div>
                        <button 
                          onClick={() => onInterceptJava(proc.pid)}
                          className="px-3 py-1 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-600 opacity-0 group-hover:opacity-100 transition-all hover:border-blue-500 hover:text-blue-600 hover:shadow-sm active:scale-95 cursor-pointer"
                        >
                          Intercept
                        </button>
                      </div>
                    )) : (
                      <div className="px-4 py-8 text-center text-slate-400 text-xs italic">
                        No Java processes detected. Make sure 'jps' is in your PATH.
                      </div>
                    )}
                  </div>
                </div>

                <div>
                  <label className="text-[10px] font-black uppercase text-slate-400 mb-2 block tracking-widest">JVM Arguments</label>
                  <div className="relative group">
                    <pre className="bg-slate-900 text-amber-200 p-4 rounded-xl text-xs font-mono overflow-x-auto">
                      -Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=8080 \<br/>
                      -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=8080
                    </pre>
                  </div>
                </div>

                <div className="bg-amber-50 border border-amber-100 rounded-xl p-4">
                  <h4 className="text-xs font-bold text-amber-800 mb-1 flex items-center gap-2">
                    <Shield size={14} /> HTTPS Note
                  </h4>
                  <p className="text-[11px] text-amber-700 leading-relaxed">
                    For HTTPS, you must import the CA certificate into your Java keystore or use <code className="bg-amber-100 px-1 rounded">-Djavax.net.ssl.trustStore</code> pointing to a keystore containing the CA.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
};
