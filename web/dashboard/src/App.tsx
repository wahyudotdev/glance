import React, { useEffect, useState, useMemo, useRef } from 'react';
import { Trash2, Check, XCircle, Info, ChevronLeft, ChevronRight } from 'lucide-react';
import type { TrafficEntry, Config, JavaProcess, AndroidDevice } from './types/traffic';

// Layout Components
import { Sidebar } from './components/layout/Sidebar';
import { Header } from './components/layout/Header';

// Feature Components
import { TrafficList } from './components/traffic/TrafficList';
import { DetailsPanel } from './components/traffic/DetailsPanel';
import { IntegrationsView } from './components/integrations/IntegrationsView';
import { SettingsView } from './components/settings/SettingsView';

// UI Components
import { Modal } from './components/ui/Modal';

const App: React.FC = () => {
  // Navigation & UI State
  const [currentView, setCurrentView] = useState<'traffic' | 'integrations' | 'settings'>('traffic');
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  
  // Unified Notification State
  const [notification, setNotification] = useState<{
    isOpen: boolean;
    type: 'success' | 'error' | 'info';
    title: string;
    message: string;
  }>({
    isOpen: false,
    type: 'info',
    title: '',
    message: ''
  });

  const [detailsWidth, setDetailsWidth] = useState(450);
  const [isResizing, setIsResizing] = useState(false);

  const startResizing = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
  };

  const stopResizing = () => {
    setIsResizing(false);
  };

  const resize = (e: MouseEvent) => {
    if (isResizing) {
      const newWidth = window.innerWidth - e.clientX;
      if (newWidth > 300 && newWidth < window.innerWidth - 400) {
        setDetailsWidth(newWidth);
      }
    }
  };

  useEffect(() => {
    window.addEventListener('mousemove', resize);
    window.addEventListener('mouseup', stopResizing);
    return () => {
      window.removeEventListener('mousemove', resize);
      window.removeEventListener('mouseup', stopResizing);
    };
  }, [isResizing]);

  const notify = (type: 'success' | 'error' | 'info', title: string, message: string) => {
    setNotification({ isOpen: true, type, title, message });
  };
  
  // Data State
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [selectedEntry, setSelectedEntry] = useState<TrafficEntry | null>(null);
  const [filter, setFilter] = useState('');
  const [proxyAddr, setProxyAddr] = useState(':8000');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalEntries, setTotalEntries] = useState(0);
  
  // Integration State
  const [javaProcesses, setJavaProcesses] = useState<JavaProcess[]>([]);
  const [androidDevices, setAndroidDevices] = useState<AndroidDevice[]>([]);
  const [isLoadingJava, setIsLoadingJava] = useState(false);
  const [isLoadingAndroid, setIsLoadingAndroid] = useState(false);
  const [terminalScript, setTerminalScript] = useState('');
  
  // Settings State
  const [config, setConfig] = useState<Config>({
    proxy_addr: ':8000',
    api_addr: ':8081',
    mcp_addr: ':8082',
    mcp_enabled: false,
    history_limit: 500,
    max_response_size: 1048576,
    default_page_size: 50
  });

  const historyLimitRef = useRef(config.history_limit);
  useEffect(() => {
    historyLimitRef.current = config.history_limit;
  }, [config.history_limit]);

  const currentPageRef = useRef(currentPage);
  useEffect(() => {
    currentPageRef.current = currentPage;
  }, [currentPage]);

  const pageSizeRef = useRef(config.default_page_size);
  useEffect(() => {
    pageSizeRef.current = config.default_page_size;
  }, [config.default_page_size]);

  // --- API Actions (Using Native Fetch) ---

  const apiFetch = async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  };

  const fetchTraffic = async (page: number = 1, pageSize?: number) => {
    try {
      const size = pageSize || config.default_page_size;
      const data = await apiFetch(`/api/traffic?page=${page}&pageSize=${size}`);
      
      if (data && data.entries) {
        setEntries(data.entries.reverse());
        setTotalEntries(data.total);
        setCurrentPage(page);
      }
    } catch (error) {
      console.error('Error fetching traffic:', error);
    }
  };

  const handleClear = async () => {
    try {
      await apiFetch('/api/traffic', { method: 'DELETE' });
      setEntries([]);
      setSelectedEntry(null);
      setIsClearModalOpen(false);
      notify('success', 'Traffic Cleared', 'All intercepted requests have been deleted.');
    } catch (error) {
      notify('error', 'Clear Failed', String(error));
    }
  };

  const fetchJavaProcesses = async () => {
    setIsLoadingJava(true);
    try {
      const data = await apiFetch('/api/client/java/processes');
      setJavaProcesses(data || []);
    } catch (error) {
      console.error('Error fetching Java processes:', error);
    } finally {
      setIsLoadingJava(false);
    }
  };

  const fetchAndroidDevices = async () => {
    setIsLoadingAndroid(true);
    try {
      const data = await apiFetch('/api/client/android/devices');
      setAndroidDevices(data || []);
    } catch (error) {
      notify('error', 'ADB Error', 'Could not list Android devices. Ensure adb is installed and devices are connected.');
    } finally {
      setIsLoadingAndroid(false);
    }
  };

  const fetchTerminalScript = async () => {
    try {
      const res = await fetch('/api/client/terminal/setup');
      const text = await res.text();
      setTerminalScript(text);
    } catch (error) {
      console.error('Error fetching terminal script:', error);
    }
  };

  const handleInterceptJava = async (pid: string) => {
    try {
      await apiFetch(`/api/client/java/intercept/${pid}`, { method: 'POST' });
      notify('success', 'Interception Active', `Successfully injected proxy into PID ${pid}.`);
    } catch (error) {
      notify('error', 'Interception Failed', String(error));
    }
  };

  const handleInterceptAndroid = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/intercept/${id}`, { method: 'POST' });
      notify('success', 'Proxy Configured', `Android device ${id} is now routing traffic through this proxy.`);
    } catch (error) {
      notify('error', 'Configuration Failed', String(error));
    }
  };

  const handleClearAndroid = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/clear/${id}`, { method: 'POST' });
      notify('success', 'Proxy Cleared', `Android device ${id} proxy settings have been reset.`);
    } catch (error) {
      notify('error', 'Reset Failed', String(error));
    }
  };

  const handlePushAndroidCert = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/push-cert/${id}`, { method: 'POST' });
      notify('success', 'CA Cert Pushed', 'Certificate pushed to /sdcard/ and install settings opened on device.');
    } catch (error) {
      notify('error', 'Push Failed', String(error));
    }
  };

  const fetchStatus = async () => {
    try {
      const data = await apiFetch('/api/status');
      setProxyAddr(data.proxy_addr);
    } catch (error) {
      console.error('Error fetching status:', error);
    }
  };

  const fetchConfig = async () => {
    try {
      const data = await apiFetch('/api/config');
      setConfig(data);
      return data as Config;
    } catch (error) {
      console.error('Error fetching config:', error);
    }
  };

  const saveConfig = async (newConfig: Config) => {
    try {
      await apiFetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newConfig)
      });
      setConfig(newConfig);
      notify('success', 'Settings Saved', 'Your configuration has been updated successfully.');
    } catch (error) {
      notify('error', 'Save Failed', String(error));
    }
  };

  // --- Effects ---

  useEffect(() => {
    let ws: WebSocket | null = null;

    const init = async () => {
      await fetchStatus();
      const cfg = await fetchConfig();
      if (cfg) {
        await fetchTraffic(1, cfg.default_page_size);
      }

      // WebSocket setup
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/ws/traffic`;
      ws = new WebSocket(wsUrl);

              ws.onmessage = (event) => {

                try {

                  const entry: TrafficEntry = JSON.parse(event.data);

                  

                  setTotalEntries(prev => prev + 1);

          

                  setEntries((prev) => {

                    // Only add to the visible list if we are on the first page

                    if (currentPageRef.current === 1) {

                      const updated = [...prev, entry];

                      const pageSize = pageSizeRef.current;

                      // Maintain page size: remove the oldest item (first in our ASC entries list)

                      if (updated.length > pageSize) {

                        return updated.slice(1);

                      }

                      return updated;

                    }

                    return prev;

                  });

                } catch (error) {

                  console.error('Error parsing WebSocket message:', error);

                }

              };

          
            ws.onclose = () => {
        console.log('WebSocket disconnected.');
      };
    };

    init();

    return () => {
      if (ws) ws.close();
    };
  }, []);

  useEffect(() => {
    if (currentView === 'integrations') {
      fetchJavaProcesses();
      fetchAndroidDevices();
      fetchTerminalScript();
    }
  }, [currentView]);

  // --- Memos ---

  const filteredEntries = useMemo(() => {
    return entries.filter(e => 
      e.url.toLowerCase().includes(filter.toLowerCase()) || 
      e.method.toLowerCase().includes(filter.toLowerCase())
    ).reverse();
  }, [entries, filter]);

  return (
    <div className="flex h-screen bg-slate-50 text-slate-900 font-sans selection:bg-blue-100">
      <Sidebar currentView={currentView} setCurrentView={setCurrentView} />

      <div className="flex-1 flex flex-col min-w-0">
        <Header 
          proxyAddr={proxyAddr} 
          filter={filter} 
          setFilter={setFilter} 
          onClearTraffic={() => setIsClearModalOpen(true)} 
        />

        <main className="flex-1 flex overflow-hidden">
          {currentView === 'traffic' && (
            <>
              <div className="flex-1 flex flex-col min-w-0 bg-white">
                <TrafficList 
                  entries={filteredEntries} 
                  selectedEntry={selectedEntry} 
                  onSelect={setSelectedEntry} 
                />
                
                {/* Pagination Controls */}
                <div className="h-12 border-t border-slate-100 flex items-center justify-between px-6 bg-slate-50/50">
                  <div className="text-[11px] font-medium text-slate-400">
                    Showing <span className="text-slate-600 font-bold">{entries.length}</span> of <span className="text-slate-600 font-bold">{totalEntries}</span> requests
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <button 
                      onClick={() => fetchTraffic(currentPage - 1)}
                      disabled={currentPage <= 1}
                      className="p-1.5 rounded-lg border border-slate-200 bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronLeft size={16} />
                    </button>
                    <span className="text-xs font-bold text-slate-600 min-w-[3rem] text-center">
                      Page {currentPage}
                    </span>
                    <button 
                      onClick={() => fetchTraffic(currentPage + 1)}
                      disabled={currentPage * config.default_page_size >= totalEntries}
                      className="p-1.5 rounded-lg border border-slate-200 bg-white text-slate-600 hover:bg-slate-50 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronRight size={16} />
                    </button>
                  </div>
                </div>
              </div>
              
              {selectedEntry && (
                <>
                  {/* Resize Handle */}
                  <div 
                    className={`w-1.5 h-full cursor-col-resize hover:bg-blue-500/30 transition-colors flex-shrink-0 z-30 ${isResizing ? 'bg-blue-500/50' : 'bg-transparent'}`}
                    onMouseDown={startResizing}
                  />
                  <DetailsPanel entry={selectedEntry} width={detailsWidth} />
                </>
              )}
            </>
          )}

          {currentView === 'integrations' && (
            <IntegrationsView 
              javaProcesses={javaProcesses}
              androidDevices={androidDevices}
              isLoadingJava={isLoadingJava}
              isLoadingAndroid={isLoadingAndroid}
              terminalScript={terminalScript}
              onFetchJava={fetchJavaProcesses}
              onFetchAndroid={fetchAndroidDevices}
              onInterceptJava={handleInterceptJava}
              onInterceptAndroid={handleInterceptAndroid}
              onClearAndroid={handleClearAndroid}
              onPushAndroidCert={handlePushAndroidCert}
            />
          )}

          {currentView === 'settings' && (
            <SettingsView 
              config={config}
              setConfig={setConfig}
              onSave={saveConfig}
              onReset={fetchConfig}
            />
          )}
        </main>
      </div>

      {/* Unified Notification Modal */}
      <Modal 
        isOpen={notification.isOpen}
        onClose={() => setNotification({ ...notification, isOpen: false })}
        title={notification.title}
        description={notification.message}
        icon={
          notification.type === 'success' ? <Check size={32} /> :
          notification.type === 'error' ? <XCircle size={32} /> :
          <Info size={32} />
        }
        iconBgColor={
          notification.type === 'success' ? 'bg-emerald-50 text-emerald-500' :
          notification.type === 'error' ? 'bg-rose-50 text-rose-500' :
          'bg-blue-50 text-blue-500'
        }
        confirmLabel="Got it"
        confirmColor={
          notification.type === 'success' ? 'bg-emerald-500 hover:bg-emerald-600 shadow-emerald-200' :
          notification.type === 'error' ? 'bg-rose-500 hover:bg-rose-600 shadow-rose-200' :
          'bg-blue-500 hover:bg-blue-600 shadow-blue-200'
        }
        onConfirm={() => setNotification({ ...notification, isOpen: false })}
        showCancel={false}
      />

      {/* Confirmation Modal */}
      <Modal 
        isOpen={isClearModalOpen}
        onClose={() => setIsClearModalOpen(false)}
        title="Clear Traffic Logs?"
        description="This will permanently delete all captured requests from the current session. This action cannot be undone."
        icon={<Trash2 size={32} />}
        iconBgColor="bg-rose-50 text-rose-500"
        confirmLabel="Clear All"
        confirmColor="bg-rose-500 hover:bg-rose-600 shadow-rose-200"
        onConfirm={handleClear}
      />
    </div>
  );
};

export default App;
