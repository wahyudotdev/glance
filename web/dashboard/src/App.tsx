import React, { useEffect, useState, useMemo, useRef } from 'react';
import { Trash2, ChevronLeft, ChevronRight } from 'lucide-react';
import type { TrafficEntry, Config, JavaProcess, AndroidDevice } from './types/traffic';

// Layout Components
import { Sidebar } from './components/layout/Sidebar';
import { Header } from './components/layout/Header';

// Feature Components
import { TrafficList } from './components/traffic/TrafficList';
import { DetailsPanel } from './components/traffic/DetailsPanel';
import { RequestEditor } from './components/traffic/RequestEditor';
import { ResponseEditor } from './components/traffic/ResponseEditor';
import { IntegrationsView } from './components/integrations/IntegrationsView';
import { SettingsView } from './components/settings/SettingsView';
import { RulesView } from './components/settings/RulesView';
import { RuleEditor } from './components/settings/RuleEditor';
import type { Rule } from './components/settings/RulesView';

// UI Components
import { Modal } from './components/ui/Modal';
import { Toast } from './components/ui/Toast';
import type { ToastMessage } from './components/ui/Toast';

const App: React.FC = () => {
  // Navigation & UI State
  const [currentView, setCurrentView] = useState<'traffic' | 'integrations' | 'settings' | 'rules'>('traffic');
  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [isRequestEditorOpen, setIsRequestEditorOpen] = useState(false);
  const [isResponseEditorOpen, setIsResponseEditorOpen] = useState(false);
  const [isRuleEditorOpen, setIsRuleEditorOpen] = useState(false);
  
  // Toast State
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  const toast = (type: 'success' | 'error' | 'info', title: string, message: string) => {
    const id = Math.random().toString(36).substring(2, 9);
    setToasts((prev) => [...prev, { id, type, title, message }]);
  };

  const removeToast = (id: string) => {
    setToasts((prev) => prev.filter(t => t.id !== id));
  };

  const [detailsWidth, setDetailsWidth] = useState(() => {
    const saved = localStorage.getItem('agent-proxy-details-width');
    return saved ? parseFloat(saved) : 70;
  });
  const [isResizing, setIsResizing] = useState(false);

  useEffect(() => {
    localStorage.setItem('agent-proxy-details-width', detailsWidth.toString());
  }, [detailsWidth]);

  const startResizing = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
  };

  const stopResizing = () => {
    setIsResizing(false);
  };

  const resize = (e: MouseEvent) => {
    if (isResizing) {
      // Sidebar is 256px (w-64)
      const sidebarWidth = 256;
      const availableWidth = window.innerWidth - sidebarWidth;
      const mouseXFromRight = window.innerWidth - e.clientX;
      const percentage = (mouseXFromRight / availableWidth) * 100;
      
      if (percentage > 20 && percentage < 80) {
        setDetailsWidth(percentage);
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
  
  // Data State
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [selectedEntry, setSelectedEntry] = useState<TrafficEntry | null>(null);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
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
  
  // Rules State
  const [rules, setRules] = useState<Rule[]>([]);
  const [isLoadingRules, setIsLoadingRules] = useState(false);
  
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
      toast('success', 'Traffic Cleared', 'All intercepted requests have been deleted.');
    } catch (error) {
      toast('error', 'Clear Failed', String(error));
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
      toast('error', 'ADB Error', 'Could not list Android devices. Ensure adb is installed.');
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
      toast('success', 'Interception Active', `Successfully injected proxy into PID ${pid}.`);
    } catch (error) {
      toast('error', 'Interception Failed', String(error));
    }
  };

  const handleInterceptAndroid = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/intercept/${id}`, { method: 'POST' });
      toast('success', 'Proxy Configured', `Android device ${id} is now routing traffic through this proxy.`);
    } catch (error) {
      toast('error', 'Configuration Failed', String(error));
    }
  };

  const handleClearAndroid = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/clear/${id}`, { method: 'POST' });
      toast('success', 'Proxy Cleared', `Android device ${id} proxy settings have been reset.`);
    } catch (error) {
      toast('error', 'Reset Failed', String(error));
    }
  };

  const handlePushAndroidCert = async (id: string) => {
    try {
      await apiFetch(`/api/client/android/push-cert/${id}`, { method: 'POST' });
      toast('success', 'CA Cert Pushed', 'Certificate pushed to /sdcard/ and install settings opened on device.');
    } catch (error) {
      toast('error', 'Push Failed', String(error));
    }
  };

  const fetchRules = async () => {
    setIsLoadingRules(true);
    try {
      const data = await apiFetch('/api/rules');
      setRules(data || []);
    } catch (error) {
      toast('error', 'Fetch Rules Failed', String(error));
    } finally {
      setIsLoadingRules(false);
    }
  };

  const handleCreateRule = async (rule: Partial<Rule>) => {
    try {
      await apiFetch('/api/rules', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rule)
      });
      fetchRules();
      toast('success', 'Rule Created', 'The new rule has been added.');
    } catch (error) {
      toast('error', 'Create Rule Failed', String(error));
    }
  };

  const handleUpdateRule = async (id: string, rule: Partial<Rule>) => {
    try {
      await apiFetch(`/api/rules/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rule)
      });
      fetchRules();
      toast('success', 'Rule Updated', 'The rule has been modified.');
    } catch (error) {
      toast('error', 'Update Rule Failed', String(error));
    }
  };

  const handleDeleteRule = async (id: string) => {
    try {
      await apiFetch(`/api/rules/${id}`, { method: 'DELETE' });
      setRules(prev => prev.filter(r => r.id !== id));
      toast('success', 'Rule Deleted', 'The rule has been removed.');
    } catch (error) {
      toast('error', 'Delete Rule Failed', String(error));
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
      toast('success', 'Settings Saved', 'Your configuration has been updated successfully.');
    } catch (error) {
      toast('error', 'Save Failed', String(error));
    }
  };

  const handleExecuteRequest = async (req: Partial<TrafficEntry>) => {
    try {
      const isResume = !!req.id;
      const endpoint = isResume ? `/api/intercept/continue/${req.id}` : '/api/request/execute';
      
      await apiFetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          method: req.method,
          url: req.url,
          headers: req.request_headers,
          body: req.request_body
        })
      });
      toast('success', isResume ? 'Request Resumed' : 'Request Executed', isResume ? 'The modified request has been sent to the server.' : 'The custom request has been sent successfully.');
    } catch (error) {
      toast('error', 'Execution Failed', String(error));
    }
  };

  const handleContinueResponse = async (status: number, headers: Record<string, string[]>, body: string) => {
    if (!selectedEntry) return;
    try {
      await apiFetch(`/api/intercept/response/continue/${selectedEntry.id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status, headers, body })
      });
      toast('success', 'Response Resumed', 'The modified response has been sent to the client.');
    } catch (error) {
      toast('error', 'Resume Failed', String(error));
    }
  };

  const handleAbortIntercept = async (id: string) => {
    try {
      await apiFetch(`/api/intercept/abort/${id}`, { method: 'POST' });
      setIsRequestEditorOpen(false);
      setIsResponseEditorOpen(false);
      toast('info', 'Request Aborted', 'The intercepted request was discarded.');
    } catch (error) {
      toast('error', 'Abort Failed', String(error));
    }
  };

  const handleCreateBreakpoint = async (entry: TrafficEntry) => {
    try {
      await handleCreateRule({
        type: 'breakpoint',
        url_pattern: entry.url,
        method: entry.method,
        strategy: 'both'
      });
    } catch (error) {
      toast('error', 'Failed to Add Breakpoint', String(error));
    }
  };

  const handleCreateMock = async (entry: TrafficEntry) => {
    try {
      await handleCreateRule({
        type: 'mock',
        url_pattern: entry.url,
        method: entry.method,
        response: {
          status: entry.status || 200,
          body: entry.response_body || '',
          headers: { 'Content-Type': 'application/json' }
        }
      });
    } catch (error) {
      toast('error', 'Failed to Add Mock', String(error));
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

                      const msg = JSON.parse(event.data);

              

                              if (msg.type === 'intercepted') {

              

                                setSelectedEntry(msg.entry);

              

                                if (msg.intercept_type === 'response') {

              

                                  setIsResponseEditorOpen(true);

              

                                  toast('info', 'Response Paused', `${msg.entry.url} - Ready for edit.`);

              

                                } else {

              

                                  setIsRequestEditorOpen(true);

              

                                  toast('info', 'Request Paused', `${msg.entry.url} - Ready for edit.`);

              

                                }

              

                                return;

              

                              }

              

                      

              

                      const entry: TrafficEntry = msg;

              

                  

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
    if (currentView === 'rules') {
      fetchRules();
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

      <div className="flex-1 flex flex-col min-w-0" onClick={(e) => e.stopPropagation()}>
        <Header 
          proxyAddr={proxyAddr} 
          filter={filter} 
          setFilter={setFilter} 
          onClearTraffic={() => setIsClearModalOpen(true)} 
          onNewRequest={() => {
            setSelectedEntry(null);
            setIsRequestEditorOpen(true);
          }}
        />

        <main 
          className="flex-1 flex overflow-hidden"
          onClick={() => setSelectedEntry(null)}
        >
          {currentView === 'traffic' && (
            <>
              <div 
                className="flex-1 flex flex-col min-w-0 bg-white"
              >
                <TrafficList 
                  entries={filteredEntries} 
                  selectedEntry={selectedEntry} 
                  onSelect={setSelectedEntry} 
                />
                
                {/* Pagination Controls */}
                <div 
                  className="h-12 border-t border-slate-100 flex items-center justify-between px-6 bg-slate-50/50"
                  onClick={(e) => e.stopPropagation()}
                >
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
                    onClick={(e) => e.stopPropagation()}
                  />
                  <div 
                    className="flex-shrink-0 h-full overflow-hidden flex flex-col"
                    style={{ width: `${detailsWidth}%` }}
                    onClick={(e) => e.stopPropagation()}
                  >
                                      <DetailsPanel 
                                        entry={selectedEntry} 
                                        onEdit={() => setIsRequestEditorOpen(true)}
                                        onClose={() => setSelectedEntry(null)}
                                        onBreak={handleCreateBreakpoint}
                                        onMock={handleCreateMock}
                                      />
                    
                  </div>
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

          {currentView === 'rules' && (
            <RulesView 
              rules={rules}
              isLoading={isLoadingRules}
              onDelete={handleDeleteRule}
              onCreate={handleCreateRule}
              onEdit={(rule) => {
                setSelectedRule(rule);
                setIsRuleEditorOpen(true);
              }}
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

      <RequestEditor 
        isOpen={isRequestEditorOpen}
        onClose={() => setIsRequestEditorOpen(false)}
        initialRequest={selectedEntry}
        onExecute={handleExecuteRequest}
        isIntercept={!!selectedEntry && !entries.find(e => e.id === selectedEntry.id)}
        onAbort={handleAbortIntercept}
      />

      {selectedEntry && (
        <ResponseEditor 
          isOpen={isResponseEditorOpen}
          onClose={() => setIsResponseEditorOpen(false)}
          entry={selectedEntry}
          onResume={handleContinueResponse}
          onAbort={handleAbortIntercept}
        />
      )}

      <RuleEditor 
        isOpen={isRuleEditorOpen}
        onClose={() => setIsRuleEditorOpen(false)}
        rule={selectedRule}
        onSave={handleUpdateRule}
      />

      <Toast toasts={toasts} onClose={removeToast} />
    </div>
  );
};

export default App;
