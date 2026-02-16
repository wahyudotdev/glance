import React, { useState } from 'react';
import axios from 'axios';
import { Globe, Terminal, Code, Activity, Copy, Check, Shield } from 'lucide-react';
import type { JavaProcess } from '../../types/traffic';

interface IntegrationsViewProps {
  javaProcesses: JavaProcess[];
  isLoadingJava: boolean;
  terminalScript: string;
  onFetchJava: () => void;
  onInterceptJava: (pid: string) => void;
}

export const IntegrationsView: React.FC<IntegrationsViewProps> = ({ 
  javaProcesses, isLoadingJava, terminalScript, onFetchJava, onInterceptJava 
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
                  try { await axios.post('/api/client/chromium'); }
                  catch (e) { alert('Failed to launch Chromium: ' + e); }
                }}
                className="w-full py-3 bg-blue-600 text-white rounded-xl font-bold text-sm hover:bg-blue-700 active:scale-95 transition-all shadow-lg shadow-blue-200"
              >
                Launch Browser
              </button>
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
                <div className="w-12 h-12 bg-amber-50 rounded-xl flex items-center justify-center">
                  <Code className="text-amber-600" size={24} />
                </div>
                <a 
                  href="/api/ca/cert" 
                  className="flex items-center gap-2 px-4 py-2 bg-slate-100 hover:bg-slate-200 rounded-lg text-xs font-bold text-slate-600 transition-all"
                >
                  Download CA Certificate
                </a>
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
                      <div key={proc.pid} className="px-4 py-3 flex items-center justify-between hover:bg-white transition-colors group">
                        <div className="flex flex-col">
                          <span className="text-xs font-bold text-slate-700 font-mono">{proc.name}</span>
                          <span className="text-[10px] text-slate-400 font-mono">PID: {proc.pid}</span>
                        </div>
                        <button 
                          onClick={() => onInterceptJava(proc.pid)}
                          className="px-3 py-1 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-600 opacity-0 group-hover:opacity-100 transition-all hover:border-blue-500 hover:text-blue-600"
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
