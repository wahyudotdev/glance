import React, { useEffect, useState, useCallback } from 'react';
import { Trash2, ChevronLeft, ChevronRight, Play } from 'lucide-react';
import type { TrafficEntry, Rule } from './types/traffic';

// Layout Components
import { Sidebar } from './components/layout/Sidebar';
import { Header } from './components/layout/Header';
import { MCPDocs } from './components/layout/MCPDocs';
import { TerminalDocs } from './components/layout/TerminalDocs';
import { ChangelogModal } from './components/layout/ChangelogModal';
import { AboutView } from './components/layout/AboutView';

// Feature Components
import { TrafficList } from './components/traffic/TrafficList';
import { DetailsPanel } from './components/traffic/DetailsPanel';
import { RequestEditor } from './components/traffic/RequestEditor';
import { ResponseEditor } from './components/traffic/ResponseEditor';
import { IntegrationsView } from './components/integrations/IntegrationsView';
import { SettingsView } from './components/settings/SettingsView';
import { RulesView } from './components/settings/RulesView';
import { RuleEditor } from './components/settings/RuleEditor';
import { ScenariosView } from './components/traffic/ScenariosView';
import { ScenarioEditor } from './components/traffic/ScenarioEditor';
import type { Scenario } from './types/traffic';

// UI Components
import { Modal } from './components/ui/Modal';
import { Toast } from './components/ui/Toast';

// hooks
import { useDarkMode } from './hooks/useDarkMode';
import { useTraffic } from './hooks/useTraffic';
import { useToasts } from './hooks/useToasts';
import { useRules } from './hooks/useRules';
import { useScenarios } from './hooks/useScenarios';
import { useIntegrations } from './hooks/useIntegrations';
import { useConfig } from './hooks/useConfig';

const App: React.FC = () => {
  const { isDark, toggleDarkMode } = useDarkMode();
  const { toasts, toast, removeToast } = useToasts();
  const { config, setConfig, fetchConfig, saveConfig } = useConfig(toast);
  
  // Traffic Hook
  const {
    entries, totalEntries, currentPage, proxyAddr, version, mcpSessions, mcpEnabled, filter, setFilter,
    methodFilter, setMethodFilter,
    filteredEntries, fetchTraffic, fetchStatus, clearTraffic,
    setEntries, setTotalEntries, currentPageRef, pageSizeRef
  } = useTraffic(config, toast);

  const { rules, isLoadingRules, fetchRules, createRule, updateRule, deleteRule } = useRules(toast);
  const { scenarios, isLoadingScenarios, fetchScenarios, saveScenario, deleteScenario, addToScenario } = useScenarios(toast);
  const { 
    javaProcesses, androidDevices, dockerContainers, 
    isLoadingJava, isLoadingAndroid, isLoadingDocker, terminalScript,
    fetchJavaProcesses, fetchAndroidDevices, fetchDockerContainers, fetchTerminalScript,
    interceptJava, interceptAndroid, clearAndroid, pushAndroidCert,
    interceptDocker, stopInterceptDocker
  } = useIntegrations(toast);

  // UI State
  const [currentView, setCurrentView] = useState<'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about'>(() => {
    const hash = window.location.hash.replace('#/', '');
    const validViews = ['traffic', 'integrations', 'settings', 'rules', 'scenarios', 'about'];
    if (validViews.includes(hash)) return hash as any;

    const saved = localStorage.getItem('glance-current-view');
    return (saved as 'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about') || 'traffic';
  });

  useEffect(() => {
    localStorage.setItem('glance-current-view', currentView);
    // Sync with hash
    const hash = window.location.hash.replace('#/', '');
    if (hash && hash !== currentView) {
      window.location.hash = `#/${currentView}`;
    }
  }, [currentView]);

  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.replace('#/', '');
      const validViews = ['traffic', 'integrations', 'settings', 'rules', 'scenarios', 'about'];
      if (validViews.includes(hash)) {
        setCurrentView(hash as any);
      }
    };

    window.addEventListener('hashchange', handleHashChange);
    // Initialize from hash if present
    handleHashChange();
    
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  const [isClearModalOpen, setIsClearModalOpen] = useState(false);
  const [isRequestEditorOpen, setIsRequestEditorOpen] = useState(false);
  const [isResponseEditorOpen, setIsResponseEditorOpen] = useState(false);
  const [isRuleEditorOpen, setIsRuleEditorOpen] = useState(false);
  const [isScenarioEditorOpen, setIsScenarioEditorOpen] = useState(false);
  const [isDeleteScenarioModalOpen, setIsDeleteScenarioModalOpen] = useState(false);
  const [isMCPDocsOpen, setIsMCPDocsOpen] = useState(false);
  const [isTerminalDocsOpen, setIsTerminalDocsOpen] = useState(false);
  const [isChangelogOpen, setIsChangelogOpen] = useState(false);
  const [selectedEntry, setSelectedEntry] = useState<TrafficEntry | null>(null);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
  const [selectedScenario, setSelectedScenario] = useState<Scenario | null>(null);
  const [scenarioToDelete, setScenarioToDelete] = useState<Scenario | null>(null);

  // Quick Create Scenario Modal State
  const [isQuickCreateModalOpen, setIsQuickCreateModalOpen] = useState(false);
  const [quickScenarioName, setQuickScenarioName] = useState('');
  const [quickScenarioDesc, setQuickScenarioDesc] = useState('');
  const [pendingEntry, setPendingEntry] = useState<TrafficEntry | null>(null);

  // Recording State
  const [isRecording, setIsRecording] = useState(false);
  const [recordedEntries, setRecordedEntries] = useState<TrafficEntry[]>([]);
  const isRecordingRef = React.useRef(isRecording);
  const filterRef = React.useRef(filter);
  const methodFilterRef = React.useRef(methodFilter);

  useEffect(() => {
    isRecordingRef.current = isRecording;
  }, [isRecording]);

  useEffect(() => {
    filterRef.current = filter;
  }, [filter]);

  useEffect(() => {
    methodFilterRef.current = methodFilter;
  }, [methodFilter]);

  const handleToggleRecording = () => {
    if (isRecording) {
      const newScenario: Scenario = {
        id: '',
        name: `Recording ${new Date().toLocaleString()}`,
        description: '',
        steps: recordedEntries.map((e, i) => ({
          id: '',
          traffic_entry_id: e.id,
          order: i + 1,
          notes: ''
        })),
        variable_mappings: [],
        created_at: new Date().toISOString()
      };
      setSelectedScenario(newScenario);
      setIsScenarioEditorOpen(true);
      setIsRecording(false);
      setRecordedEntries([]);
    } else {
      setRecordedEntries([]);
      setIsRecording(true);
      
      let msg = 'All incoming requests will be added to the scenario.';
      if (filter || methodFilter !== 'ALL') {
        msg = `Recording active with filters: ${methodFilter !== 'ALL' ? '['+methodFilter+'] ' : ''}${filter}`;
      }
      toast('info', 'Recording Started', msg);
    }
  };

  const handleConfirmDeleteScenario = async () => {
    if (!scenarioToDelete) return;
    const success = await deleteScenario(scenarioToDelete.id);
    if (success) {
      setIsDeleteScenarioModalOpen(false);
      setScenarioToDelete(null);
    }
  };

  const [detailsWidth, setDetailsWidth] = useState(() => {
    const saved = localStorage.getItem('glance-details-width');
    return saved ? parseFloat(saved) : 70;
  });
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(() => {
    return localStorage.getItem('glance-sidebar-collapsed') === 'true';
  });
  const [isResizing, setIsResizing] = useState(false);
  const [isDetailsFullScreen, setIsDetailsFullScreen] = useState(false);
  const [preFullScreenWidth, setPreFullScreenWidth] = useState(70);

  const isResizingRef = React.useRef(isResizing);
  useEffect(() => { isResizingRef.current = isResizing; }, [isResizing]);

  const isSidebarCollapsedRef = React.useRef(isSidebarCollapsed);
  useEffect(() => { isSidebarCollapsedRef.current = isSidebarCollapsed; }, [isSidebarCollapsed]);

  useEffect(() => {
    localStorage.setItem('glance-details-width', detailsWidth.toString());
  }, [detailsWidth]);

  const handleToggleDetailsFullScreen = () => {
    if (isDetailsFullScreen) {
      setDetailsWidth(preFullScreenWidth);
      setIsDetailsFullScreen(false);
    } else {
      setPreFullScreenWidth(detailsWidth);
      setDetailsWidth(100);
      setIsDetailsFullScreen(true);
    }
  };

  useEffect(() => {
    localStorage.setItem('glance-sidebar-collapsed', isSidebarCollapsed.toString());
  }, [isSidebarCollapsed]);

  const startResizing = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
  };

  const stopResizing = useCallback(() => setIsResizing(false), []);

  const resize = useCallback((e: MouseEvent) => {
    if (isResizingRef.current) {
      const sidebarWidth = isSidebarCollapsedRef.current ? 80 : 256;
      const availableWidth = window.innerWidth - sidebarWidth;
      const mouseXFromRight = window.innerWidth - e.clientX;
      const percentage = (mouseXFromRight / availableWidth) * 100;
      if (percentage > 20 && percentage < 80) setDetailsWidth(percentage);
    }
  }, []);

  useEffect(() => {
    window.addEventListener('mousemove', resize);
    window.addEventListener('mouseup', stopResizing);
    return () => {
      window.removeEventListener('mousemove', resize);
      window.removeEventListener('mouseup', stopResizing);
    };
  }, [resize, stopResizing]);
  
  // --- Actions ---

  const handleClear = async () => {
    const success = await clearTraffic();
    if (success) {
      setSelectedEntry(null);
      setIsClearModalOpen(false);
    }
  };

  const handleConfirmCreateScenario = async () => {
    if (!pendingEntry || !quickScenarioName) return;

    const newScenario: Scenario = {
      id: '',
      name: quickScenarioName,
      description: quickScenarioDesc,
      steps: [{
        id: '',
        traffic_entry_id: pendingEntry.id,
        order: 1,
        notes: '',
        traffic_entry: pendingEntry
      }],
      variable_mappings: [],
      created_at: new Date().toISOString()
    };

    const success = await saveScenario(newScenario);
    if (success) {
      setIsQuickCreateModalOpen(false);
      setPendingEntry(null);
    }
  };

  const handleExecuteRequest = async (req: Partial<TrafficEntry>) => {
    try {
      const isResume = !!req.id;
      const endpoint = isResume ? `/api/intercept/continue/${req.id}` : '/api/request/execute';
      
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          method: req.method,
          url: req.url,
          headers: req.request_headers,
          body: req.request_body
        })
      });
      if (!res.ok) throw new Error(await res.text());
      toast('success', isResume ? 'Request Resumed' : 'Request Executed', isResume ? 'The modified request has been sent to the server.' : 'The custom request has been sent successfully.');
    } catch (error) {
      toast('error', 'Execution Failed', String(error));
    }
  };

  const handleContinueResponse = async (status: number, headers: Record<string, string[]>, body: string) => {
    if (!selectedEntry) return;
    try {
      const res = await fetch(`/api/intercept/response/continue/${selectedEntry.id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status, headers, body })
      });
      if (!res.ok) throw new Error(await res.text());
      toast('success', 'Response Resumed', 'The modified response has been sent to the client.');
    } catch (error) {
      toast('error', 'Resume Failed', String(error));
    }
  };

  const handleAbortIntercept = async (id: string) => {
    try {
      const res = await fetch(`/api/intercept/abort/${id}`, { method: 'POST' });
      if (!res.ok) throw new Error(await res.text());
      setIsRequestEditorOpen(false);
      setIsResponseEditorOpen(false);
      toast('info', 'Request Aborted', 'The intercepted request was discarded.');
    } catch (error) {
      toast('error', 'Abort Failed', String(error));
    }
  };

  const handleCreateBreakpoint = async (entry: TrafficEntry) => {
    await createRule({
      type: 'breakpoint',
      url_pattern: entry.url,
      method: entry.method,
      strategy: 'both'
    });
  };

  const handleCreateMock = async (entry: TrafficEntry) => {
    await createRule({
      type: 'mock',
      url_pattern: entry.url,
      method: entry.method,
      response: {
        status: entry.status || 200,
        body: entry.response_body || '',
        headers: { 'Content-Type': 'application/json' }
      }
    });
  };

  // --- Refs for WebSocket Handlers ---
  const toastRef = React.useRef(toast);
  const setEntriesRef = React.useRef(setEntries);
  const setTotalEntriesRef = React.useRef(setTotalEntries);

  useEffect(() => { toastRef.current = toast; }, [toast]);
  useEffect(() => { setEntriesRef.current = setEntries; }, [setEntries]);
  useEffect(() => { setTotalEntriesRef.current = setTotalEntries; }, [setTotalEntries]);

  // --- Initial Data Loading ---
  useEffect(() => {
    const initData = async () => {
      await fetchStatus();
      const cfg = await fetchConfig();
      if (cfg) {
        await fetchTraffic(1, cfg.default_page_size);
      }
    };
    initData();
  }, [fetchConfig, fetchStatus, fetchTraffic]);

  // --- WebSocket Connection ---
  useEffect(() => {
    let ws: WebSocket | null = null;
    let reconnectTimer: number | null = null;

    const connect = () => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/ws/traffic`;
      ws = new WebSocket(wsUrl);

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);
          let entry: TrafficEntry;

          if (msg.type === 'intercepted') {
            entry = msg.entry;
            setSelectedEntry(entry);
            if (msg.intercept_type === 'response') {
              setIsResponseEditorOpen(true);
              toastRef.current('info', 'Response Paused', `${entry.url} - Ready for edit.`);
            } else {
              setIsRequestEditorOpen(true);
              toastRef.current('info', 'Request Paused', `${entry.url} - Ready for edit.`);
            }
          } else {
            entry = msg;
          }

          setEntriesRef.current((prev) => {
            const exists = prev.some(e => e.id === entry.id);
            if (exists) {
              return prev.map(e => e.id === entry.id ? entry : e);
            }

            setTotalEntriesRef.current(curr => curr + 1);

            if (currentPageRef.current === 1) {
              const updated = [...prev, entry];
              const pageSize = pageSizeRef.current;
              if (updated.length > pageSize) return updated.slice(1);
              return updated;
            }
            return prev;
          });

          if (isRecordingRef.current) {
            const f = filterRef.current.toLowerCase();
            const mf = methodFilterRef.current;
            const matchesText = !f || entry.url.toLowerCase().includes(f) || entry.method.toLowerCase().includes(f);
            const matchesMethod = mf === 'ALL' || entry.method.toUpperCase() === mf;

            if (matchesText && matchesMethod) {
              setRecordedEntries(prev => [...prev, entry]);
            }
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected. Reconnecting in 3s...');
        reconnectTimer = window.setTimeout(connect, 3000);
      };
      
      ws.onerror = (err) => {
        console.error('WebSocket error:', err);
        ws?.close();
      };
    };

    connect();

    return () => {
      if (ws) {
        ws.onclose = null;
        ws.close();
      }
      if (reconnectTimer) clearTimeout(reconnectTimer);
    };
  }, [currentPageRef, pageSizeRef]);

  useEffect(() => {
    if (currentView === 'integrations') {
      fetchJavaProcesses();
      fetchAndroidDevices();
      fetchDockerContainers();
      fetchTerminalScript();
    }
    if (currentView === 'rules') {
      fetchRules();
    }
    if (currentView === 'scenarios') {
      fetchScenarios();
    }
  }, [currentView, fetchAndroidDevices, fetchJavaProcesses, fetchDockerContainers, fetchRules, fetchScenarios, fetchTerminalScript]);

  return (
    <div className="flex h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 font-sans selection:bg-blue-100 dark:selection:bg-blue-900 transition-colors">
      <Sidebar 
        currentView={currentView} 
        setCurrentView={setCurrentView} 
        isCollapsed={isSidebarCollapsed}
        onToggleCollapse={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
        version={version}
        onShowChangelog={() => setIsChangelogOpen(true)}
      />

      <div className="flex-1 flex flex-col min-w-0" onClick={(e) => e.stopPropagation()}>
        <Header 
          proxyAddr={proxyAddr} 
          mcpSessions={mcpSessions}
          mcpEnabled={mcpEnabled}
          filter={filter} 
          setFilter={setFilter} 
          methodFilter={methodFilter}
          setMethodFilter={setMethodFilter}
          onClearTraffic={() => setIsClearModalOpen(true)} 
          isDark={isDark}
          toggleDarkMode={toggleDarkMode}
          isRecording={isRecording}
          onToggleRecording={handleToggleRecording}
          recordedCount={recordedEntries.length}
          onShowMCP={() => setIsMCPDocsOpen(true)}
        />

        <main className="flex-1 flex overflow-hidden" onClick={() => setSelectedEntry(null)}>
          {currentView === 'traffic' && (
            <>
              <div className="flex-1 flex flex-col min-w-0 bg-white dark:bg-slate-900 transition-colors">
                <TrafficList entries={filteredEntries} selectedEntry={selectedEntry} onSelect={setSelectedEntry} />
                
                <div className="h-12 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between px-6 bg-slate-50/50 dark:bg-slate-950/50 transition-colors" onClick={(e) => e.stopPropagation()}>
                  <div className="text-[11px] font-medium text-slate-400 dark:text-slate-500">
                    Showing <span className="text-slate-600 dark:text-slate-300 font-bold">{entries.length}</span> of <span className="text-slate-600 dark:text-slate-300 font-bold">{totalEntries}</span> requests
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <button 
                      onClick={() => fetchTraffic(currentPage - 1)}
                      disabled={currentPage <= 1}
                      className="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronLeft size={16} />
                    </button>
                    <span className="text-xs font-bold text-slate-600 dark:text-slate-300 min-w-[3rem] text-center">Page {currentPage}</span>
                    <button 
                      onClick={() => fetchTraffic(currentPage + 1)}
                      disabled={currentPage * config.default_page_size >= totalEntries}
                      className="p-1.5 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronRight size={16} />
                    </button>
                  </div>
                </div>
              </div>
              
              {selectedEntry && (
                <>
                  <div 
                    className={`w-1.5 h-full cursor-col-resize hover:bg-blue-500/30 transition-colors flex-shrink-0 z-30 ${isResizing ? 'bg-blue-500/50' : 'bg-transparent'}`}
                    onMouseDown={startResizing}
                    onClick={(e) => e.stopPropagation()}
                  />
                  <div className="flex-shrink-0 h-full overflow-hidden flex flex-col" style={{ width: `${detailsWidth}%` }} onClick={(e) => e.stopPropagation()}>
                    <DetailsPanel 
                      entry={selectedEntry} 
                      scenarios={scenarios}
                      onEdit={() => setIsRequestEditorOpen(true)} 
                      onClose={() => setSelectedEntry(null)} 
                      onBreak={handleCreateBreakpoint} 
                      onMock={handleCreateMock} 
                      onAddToScenario={(entry, scenarioId) => addToScenario(entry, scenarioId, (e) => {
                        setPendingEntry(e);
                        setQuickScenarioName(`New Scenario ${new Date().toLocaleTimeString()}`);
                        setQuickScenarioDesc('');
                        setIsQuickCreateModalOpen(true);
                      })} 
                      onToggleFullScreen={handleToggleDetailsFullScreen}
                      isPanelFullScreen={isDetailsFullScreen}
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
              dockerContainers={dockerContainers}
              isLoadingJava={isLoadingJava} 
              isLoadingAndroid={isLoadingAndroid}
              isLoadingDocker={isLoadingDocker}
              onFetchJava={fetchJavaProcesses} 
              onFetchAndroid={fetchAndroidDevices} 
              onFetchDocker={fetchDockerContainers}
              onInterceptJava={interceptJava} 
              onInterceptAndroid={interceptAndroid} 
              onClearAndroid={clearAndroid} 
              onPushAndroidCert={pushAndroidCert}
              onInterceptDocker={interceptDocker}
              onStopInterceptDocker={stopInterceptDocker}
              onShowTerminalDocs={() => setIsTerminalDocsOpen(true)}
            />
          )}

          {currentView === 'rules' && (
            <RulesView 
              rules={rules} isLoading={isLoadingRules} onDelete={deleteRule} onCreate={createRule}
              onUpdate={updateRule}
              onEdit={(rule) => { setSelectedRule(rule); setIsRuleEditorOpen(true); }}
            />
          )}

          {currentView === 'settings' && (
            <SettingsView config={config} setConfig={setConfig} onSave={saveConfig} onReset={fetchConfig} onShowMCP={() => setIsMCPDocsOpen(true)} />
          )}

          {currentView === 'scenarios' && (
            <ScenariosView 
              scenarios={scenarios}
              isLoading={isLoadingScenarios}
              onSelect={(s) => { setSelectedScenario(s); setIsScenarioEditorOpen(true); }}
              onDelete={(id) => {
                const scenario = scenarios.find(s => s.id === id);
                if (scenario) {
                  setScenarioToDelete(scenario);
                  setIsDeleteScenarioModalOpen(true);
                }
              }}
              onCreateNew={() => { setSelectedScenario(null); setIsScenarioEditorOpen(true); }}
            />
          )}

          {currentView === 'about' && (
            <AboutView />
          )}
        </main>
      </div>

      <Modal 
        isOpen={isClearModalOpen} onClose={() => setIsClearModalOpen(false)} title="Clear Traffic Logs?" description="This will permanently delete all captured requests from the current session. This action cannot be undone."
        icon={<Trash2 size={32} />} iconBgColor="bg-rose-50 dark:bg-rose-900/20 text-rose-500" confirmLabel="Clear All" confirmColor="bg-rose-500 hover:bg-rose-600 shadow-rose-200" onConfirm={handleClear}
      />

      <Modal 
        isOpen={isDeleteScenarioModalOpen} 
        onClose={() => { setIsDeleteScenarioModalOpen(false); setScenarioToDelete(null); }} 
        title="Delete Scenario?" 
        description={`Are you sure you want to delete "${scenarioToDelete?.name}"? This action cannot be undone.`}
        icon={<Trash2 size={32} />} 
        iconBgColor="bg-rose-50 dark:bg-rose-900/20 text-rose-500" 
        confirmLabel="Delete Scenario" 
        confirmColor="bg-rose-500 hover:bg-rose-600 shadow-rose-200" 
        onConfirm={handleConfirmDeleteScenario}
      />

      {/* Quick Create Scenario Modal */}
      {isQuickCreateModalOpen && (
        <div className="fixed inset-0 z-[110] flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300" onClick={() => setIsQuickCreateModalOpen(false)} />
          <div className="relative bg-white dark:bg-slate-900 rounded-3xl shadow-2xl w-full max-w-md overflow-hidden border border-white/50 dark:border-slate-800/50 animate-in zoom-in-95 duration-300">
            <div className="p-8">
              <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 mb-6 tracking-tight flex items-center gap-3">
                <div className="p-2 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-lg">
                  <Play size={20} />
                </div>
                Create New Scenario
              </h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-2">Scenario Name</label>
                  <input 
                    type="text" 
                    className="w-full px-4 py-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all dark:text-slate-100 font-bold"
                    value={quickScenarioName}
                    onChange={(e) => setQuickScenarioName(e.target.value)}
                    autoFocus
                  />
                </div>
                <div>
                  <label className="block text-[10px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-2">Description (Optional)</label>
                  <textarea 
                    className="w-full px-4 py-3 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all dark:text-slate-100 resize-none h-24"
                    value={quickScenarioDesc}
                    onChange={(e) => setQuickScenarioDesc(e.target.value)}
                  />
                </div>
              </div>
            </div>
            <div className="flex border-t border-slate-100 dark:border-slate-800 p-4 gap-3 bg-slate-50/50 dark:bg-slate-950/50">
              <button 
                onClick={() => setIsQuickCreateModalOpen(false)}
                className="flex-1 px-4 py-3 text-sm font-bold text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 rounded-xl transition-all"
              >
                Cancel
              </button>
              <button 
                onClick={handleConfirmCreateScenario}
                disabled={!quickScenarioName}
                className="flex-1 px-4 py-3 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-200 dark:disabled:bg-slate-800 disabled:text-slate-400 dark:disabled:text-slate-600 rounded-xl transition-all shadow-lg active:scale-[0.98]"
              >
                Create & Add
              </button>
            </div>
          </div>
        </div>
      )}

      <RequestEditor 
        isOpen={isRequestEditorOpen} onClose={() => setIsRequestEditorOpen(false)} initialRequest={selectedEntry} onExecute={handleExecuteRequest}
        isIntercept={!!selectedEntry && !entries.find(e => e.id === selectedEntry.id)} onAbort={handleAbortIntercept}
      />

      {selectedEntry && (
        <ResponseEditor isOpen={isResponseEditorOpen} onClose={() => setIsResponseEditorOpen(false)} entry={selectedEntry} onResume={handleContinueResponse} onAbort={handleAbortIntercept} />
      )}

      <RuleEditor isOpen={isRuleEditorOpen} onClose={() => setIsRuleEditorOpen(false)} rule={selectedRule} onSave={updateRule} />

      <ScenarioEditor 
        isOpen={isScenarioEditorOpen} 
        onClose={() => { setIsScenarioEditorOpen(false); setSelectedScenario(null); }}
        scenario={selectedScenario}
        onSave={(s) => {
          saveScenario(s);
          setIsScenarioEditorOpen(false);
          setSelectedScenario(null);
          setCurrentView('scenarios');
        }}
        availableTraffic={entries}
      />

      <Toast toasts={toasts} onClose={removeToast} />

      <MCPDocs 
        isOpen={isMCPDocsOpen} onClose={() => setIsMCPDocsOpen(false)}
        mcpUrl={`${window.location.protocol}//${window.location.hostname}${config.mcp_addr}/mcp`}
      />

      <TerminalDocs 
        isOpen={isTerminalDocsOpen}
        onClose={() => setIsTerminalDocsOpen(false)}
        terminalScript={terminalScript}
      />

      <ChangelogModal 
        isOpen={isChangelogOpen}
        onClose={() => setIsChangelogOpen(false)}
      />
    </div>
  );
};

export default App;
