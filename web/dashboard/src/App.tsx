import React, { useEffect, useState, useMemo } from 'react';
import axios from 'axios';
import type { TrafficEntry } from './types/traffic';
import { 
  Activity, Globe, Shield, Search, Trash2, 
  Copy, Check, ChevronRight, FileText, Settings, Code
} from 'lucide-react';
import dayjs from 'dayjs';
import { generateCurl } from './lib/curl';

const App: React.FC = () => {
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [selectedEntry, setSelectedEntry] = useState<TrafficEntry | null>(null);
  const [filter, setFilter] = useState('');
  const [copied, setCopied] = useState(false);
  const [activeTab, setActiveTab] = useState<'headers' | 'body' | 'curl'>('headers');
  const [currentView, setCurrentView] = useState<'traffic' | 'integrations'>('traffic');
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [javaProcesses, setJavaProcesses] = useState<{pid: string, name: string}[]>([]);
  const [isLoadingJava, setIsLoadingJava] = useState(false);
  const [proxyAddr, setProxyAddr] = useState(':8000');

  const fetchTraffic = async () => {
    try {
      const response = await axios.get('http://localhost:8081/api/traffic');
      setEntries(response.data || []);
    } catch (error) {
      console.error('Error fetching traffic:', error);
    }
  };

  const handleClear = async () => {
    try {
      await axios.delete('http://localhost:8081/api/traffic');
      setEntries([]);
      setSelectedEntry(null);
      setIsClearModalOpen(false);
    } catch (error) {
      alert('Failed to clear traffic: ' + error);
    }
  };

  const fetchJavaProcesses = async () => {
    setIsLoadingJava(true);
    try {
      const response = await axios.get('http://localhost:8081/api/client/java/processes');
      setJavaProcesses(response.data || []);
    } catch (error) {
      console.error('Error fetching Java processes:', error);
    } finally {
      setIsLoadingJava(false);
    }
  };

  useEffect(() => {
    if (currentView === 'integrations') {
      fetchJavaProcesses();
    }
  }, [currentView]);

  const fetchStatus = async () => {
    try {
      const response = await axios.get('http://localhost:8081/api/status');
      setProxyAddr(response.data.proxy_addr);
    } catch (error) {
      console.error('Error fetching status:', error);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchTraffic, 1500);
    return () => clearInterval(interval);
  }, []);

  const filteredEntries = useMemo(() => {
    return entries.filter(e => 
      e.url.toLowerCase().includes(filter.toLowerCase()) || 
      e.method.toLowerCase().includes(filter.toLowerCase())
    ).reverse();
  }, [entries, filter]);

  const handleCopyCurl = (entry: TrafficEntry) => {
    const curl = generateCurl(entry);
    navigator.clipboard.writeText(curl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const getStatusColor = (status: number) => {
    if (status >= 200 && status < 300) return 'text-emerald-600 bg-emerald-50 border-emerald-100';
    if (status >= 300 && status < 400) return 'text-amber-600 bg-amber-50 border-amber-100';
    if (status >= 400) return 'text-rose-600 bg-rose-50 border-rose-100';
    return 'text-slate-600 bg-slate-50 border-slate-100';
  };

  return (
    <div className="flex h-screen bg-slate-50 text-slate-900 font-sans selection:bg-blue-100">
      {/* Sidebar */}
      <aside className="w-16 flex flex-col items-center py-6 bg-white border-r border-slate-200 gap-8">
        <div className="p-2 bg-blue-600 rounded-xl shadow-lg shadow-blue-200 cursor-pointer" onClick={() => setCurrentView('traffic')}>
          <Activity className="text-white" size={24} />
        </div>
        <nav className="flex flex-col gap-4">
          <button 
            onClick={() => setCurrentView('traffic')}
            className={`p-2 rounded-lg transition-all ${currentView === 'traffic' ? 'text-blue-600 bg-blue-50' : 'text-slate-400 hover:text-slate-600'}`}
          >
            <Globe size={20} />
          </button>
          <button 
            onClick={() => setCurrentView('integrations')}
            className={`p-2 rounded-lg transition-all ${currentView === 'integrations' ? 'text-blue-600 bg-blue-50' : 'text-slate-400 hover:text-slate-600'}`}
          >
            <Settings size={20} />
          </button>
        </nav>
      </aside>

      {/* Main Container */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="h-16 flex items-center justify-between px-8 bg-white border-b border-slate-200 shadow-sm z-10">
          <div className="flex items-center gap-4">
            <h1 className="text-lg font-bold tracking-tight text-slate-800">Traffic Inspector</h1>
            <div className="flex items-center gap-3 px-3 py-1.5 bg-slate-100 rounded-full text-[11px] font-semibold text-slate-500 uppercase tracking-wider">
              <span className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /> Proxy {proxyAddr}</span>
              <span className="w-px h-3 bg-slate-300" />
              <span className="flex items-center gap-1.5 text-blue-600"><Shield size={12} /> MITM Active</span>
            </div>
          </div>
          
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
              <input 
                type="text" 
                placeholder="Filter requests..." 
                className="pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all w-64"
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
              />
            </div>
            <button 
              onClick={() => setIsClearModalOpen(true)}
              className="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all" 
              title="Clear Logs"
            >
              <Trash2 size={18} />
            </button>
          </div>
        </header>

        {/* Content Area */}
        <main className="flex-1 flex overflow-hidden">
          {currentView === 'traffic' ? (
            <>
              {/* Traffic List */}
              <div className="flex-1 overflow-y-auto bg-white">
                <table className="w-full text-left border-separate border-spacing-0">
                  <thead className="sticky top-0 bg-white/80 backdrop-blur-md z-10 shadow-sm">
                    <tr className="text-[11px] uppercase tracking-widest text-slate-400 font-bold border-b border-slate-100">
                      <th className="pl-8 pr-4 py-4 font-bold">Method</th>
                      <th className="px-4 py-4">Status</th>
                      <th className="px-4 py-4">Path & Host</th>
                      <th className="px-4 py-4">Size</th>
                      <th className="pr-8 pl-4 py-4 text-right">Time</th>
                    </tr>
                  </thead>
                  <tbody className="text-[13px]">
                    {filteredEntries.map((entry) => (
                      <tr
                        key={entry.id}
                        onClick={() => setSelectedEntry(entry)}
                        className={`group cursor-pointer border-b border-slate-50 transition-all duration-150 hover:bg-blue-50/50 ${
                          selectedEntry?.id === entry.id ? 'bg-blue-50 border-blue-100' : ''
                        }`}
                      >
                        <td className="pl-8 pr-4 py-3.5">
                          <span className="font-mono font-bold text-blue-600 tracking-tighter">
                            {entry.method}
                          </span>
                        </td>
                        <td className="px-4 py-3.5">
                          <span className={`px-2.5 py-1 rounded-md text-[11px] font-bold border tabular-nums ${getStatusColor(entry.status)}`}>
                            {entry.status || '---'}
                          </span>
                        </td>
                        <td className="px-4 py-3.5 max-w-xl">
                          <div className="flex flex-col">
                            <span className="text-slate-700 font-medium truncate font-mono">
                              {new URL(entry.url).pathname}{new URL(entry.url).search}
                            </span>
                            <span className="text-slate-400 text-[11px] truncate">
                              {new URL(entry.url).host}
                            </span>
                          </div>
                        </td>
                        <td className="px-4 py-3.5 text-slate-400 tabular-nums">
                          {entry.response_body ? `${(entry.response_body.length / 1024).toFixed(1)} KB` : '-'}
                        </td>
                        <td className="pr-8 pl-4 py-3.5 text-right text-slate-400 tabular-nums group-hover:text-slate-600 flex items-center justify-end gap-2">
                          {dayjs(entry.start_time).format('HH:mm:ss.SSS')}
                          <ChevronRight size={14} className={`opacity-0 group-hover:opacity-100 transition-opacity ${selectedEntry?.id === entry.id ? 'opacity-100' : ''}`} />
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Details Panel */}
              {selectedEntry && (
                <div className="w-[450px] bg-white border-l border-slate-200 flex flex-col shadow-2xl z-20">
                  <div className="p-6 border-b border-slate-100 bg-slate-50/50">
                    <div className="flex items-center justify-between mb-4">
                      <h2 className="text-sm font-bold text-slate-800 flex items-center gap-2 uppercase tracking-tight">
                        <FileText size={16} className="text-blue-500" /> Request Details
                      </h2>
                      <button 
                        onClick={() => handleCopyCurl(selectedEntry)}
                        className="flex items-center gap-1.5 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-semibold text-slate-600 hover:border-blue-500 hover:text-blue-600 transition-all shadow-sm active:scale-95"
                      >
                        {copied ? <Check size={14} className="text-emerald-500" /> : <Copy size={14} />}
                        {copied ? 'Copied!' : 'Copy cURL'}
                      </button>
                    </div>
                    <div className="font-mono text-[12px] bg-slate-900 text-slate-300 p-3 rounded-lg break-all leading-relaxed shadow-inner border border-slate-800">
                      <span className="text-blue-400 font-bold uppercase mr-2">{selectedEntry.method}</span>
                      {selectedEntry.url}
                    </div>
                  </div>

                  {/* Tabs */}
                  <div className="flex px-6 pt-2 border-b border-slate-100 bg-white">
                    {(['headers', 'body', 'curl'] as const).map((tab) => (
                      <button
                        key={tab}
                        onClick={() => setActiveTab(tab)}
                        className={`px-4 py-3 text-[11px] font-bold uppercase tracking-wider transition-all border-b-2 -mb-[1px] ${
                          activeTab === tab 
                          ? 'border-blue-600 text-blue-600' 
                          : 'border-transparent text-slate-400 hover:text-slate-600'
                        }`}
                      >
                        {tab}
                      </button>
                    ))}
                  </div>

                  <div className="flex-1 overflow-y-auto p-6">
                    {activeTab === 'headers' && (
                      <div className="space-y-6">
                        <section>
                          <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Request Headers</h3>
                          <div className="space-y-1.5">
                            {Object.entries(selectedEntry.request_headers).map(([key, values]) => (
                              <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 last:border-0">
                                <span className="font-bold text-slate-600 min-w-[120px] shrink-0">{key}</span>
                                <span className="text-slate-500 break-all font-mono">{values.join(', ')}</span>
                              </div>
                            ))}
                          </div>
                        </section>
                        
                        {selectedEntry.response_headers && (
                          <section>
                            <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Response Headers</h3>
                            <div className="space-y-1.5">
                              {Object.entries(selectedEntry.response_headers).map(([key, values]) => (
                                <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 last:border-0">
                                  <span className="font-bold text-slate-600 min-w-[120px] shrink-0">{key}</span>
                                  <span className="text-slate-500 break-all font-mono">{values.join(', ')}</span>
                                </div>
                              ))}
                            </div>
                          </section>
                        )}
                      </div>
                    )}

                    {activeTab === 'body' && (
                      <div className="h-full flex flex-col gap-4">
                        <section className="flex-1 min-h-0 flex flex-col">
                          <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Response Body</h3>
                          <div className="flex-1 bg-slate-50 rounded-xl p-4 font-mono text-[12px] overflow-auto border border-slate-100 shadow-inner">
                            {selectedEntry.response_body ? (
                              <pre className="text-slate-600 whitespace-pre-wrap leading-relaxed">
                                {(() => {
                                  try {
                                    return JSON.stringify(JSON.parse(selectedEntry.response_body), null, 2);
                                  } catch {
                                    return selectedEntry.response_body;
                                  }
                                })()}
                              </pre>
                            ) : (
                              <span className="text-slate-300 italic">No response body captured</span>
                            )}
                          </div>
                        </section>
                      </div>
                    )}

                    {activeTab === 'curl' && (
                      <div className="space-y-4">
                        <h3 className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">cURL Command</h3>
                        <pre className="bg-slate-900 text-blue-300 p-5 rounded-xl text-[12px] font-mono whitespace-pre-wrap leading-relaxed shadow-xl">
                          {generateCurl(selectedEntry)}
                        </pre>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </>
          ) : (
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
                          try { await axios.post('http://localhost:8081/api/client/chromium'); }
                          catch (e) { alert('Failed to launch Chromium: ' + e); }
                        }}
                        className="w-full py-3 bg-blue-600 text-white rounded-xl font-bold text-sm hover:bg-blue-700 active:scale-95 transition-all shadow-lg shadow-blue-200"
                      >
                        Launch Browser
                      </button>
                    </div>

                    <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm opacity-60">
                      <div className="w-12 h-12 bg-slate-50 rounded-xl flex items-center justify-center mb-6">
                        <Activity className="text-slate-400" size={24} />
                      </div>
                      <h3 className="text-lg font-bold mb-2">Android (ADB)</h3>
                      <p className="text-sm text-slate-500 mb-6 leading-relaxed">
                        Configure connected Android devices or emulators to use the proxy via ADB commands. (Coming Soon)
                      </p>
                      <button disabled className="w-full py-3 bg-slate-200 text-slate-400 rounded-xl font-bold text-sm cursor-not-allowed">
                        Coming Soon
                      </button>
                    </div>

                    <div className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow md:col-span-2">
                      <div className="flex items-start justify-between mb-6">
                        <div className="w-12 h-12 bg-amber-50 rounded-xl flex items-center justify-center">
                          <Code className="text-amber-600" size={24} />
                        </div>
                        <a 
                          href="http://localhost:8081/api/ca/cert" 
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
                        {/* Process List */}
                        <div className="bg-slate-50 rounded-xl border border-slate-100 overflow-hidden">
                          <div className="px-4 py-3 border-b border-slate-200 bg-white flex items-center justify-between">
                            <span className="text-[10px] font-black uppercase text-slate-400 tracking-widest">Running Java Processes</span>
                            <button 
                              onClick={fetchJavaProcesses}
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
                                  onClick={() => alert(`To intercept PID ${proc.pid}, restart the application with the JVM arguments below.`)}
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
                            <button 
                              onClick={() => {
                                navigator.clipboard.writeText('-Dhttp.proxyHost=127.0.0.1 -Dhttp.proxyPort=8080 -Dhttps.proxyHost=127.0.0.1 -Dhttps.proxyPort=8080');
                              }}
                              className="absolute top-2 right-2 p-2 bg-slate-800 text-slate-400 hover:text-white rounded-lg opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                              <Copy size={14} />
                            </button>
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
          )}
        </main>
      </div>

      {/* Confirmation Modal */}
      {isClearModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-white rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden animate-in zoom-in-95 duration-200">
            <div className="p-8 text-center">
              <div className="w-16 h-16 bg-rose-50 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-6">
                <Trash2 size={32} />
              </div>
              <h3 className="text-xl font-bold text-slate-800 mb-2">Clear Traffic Logs?</h3>
              <p className="text-slate-500 text-sm leading-relaxed">
                This will permanently delete all captured requests from the current session. This action cannot be undone.
              </p>
            </div>
            <div className="flex border-t border-slate-100 p-4 gap-3 bg-slate-50/50">
              <button 
                onClick={() => setIsClearModalOpen(false)}
                className="flex-1 px-4 py-3 text-sm font-bold text-slate-600 hover:bg-white rounded-xl transition-all"
              >
                Cancel
              </button>
              <button 
                onClick={handleClear}
                className="flex-1 px-4 py-3 text-sm font-bold text-white bg-rose-500 hover:bg-rose-600 rounded-xl transition-all shadow-lg shadow-rose-200"
              >
                Clear All
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default App;
