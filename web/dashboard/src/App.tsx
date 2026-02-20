import React, { useEffect, useState } from 'react';
import { Trash2, ChevronLeft, ChevronRight, Play } from 'lucide-react';
import type { TrafficEntry, Config, JavaProcess, AndroidDevice } from './types/traffic';

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
import type { Rule } from './components/settings/RulesView';
import type { Scenario } from './types/traffic';

// UI Components
import { Modal } from './components/ui/Modal';
import { Toast } from './components/ui/Toast';
import type { ToastMessage } from './components/ui/Toast';

// hooks
import { useDarkMode } from './hooks/useDarkMode';
import { useTraffic } from './hooks/useTraffic';

const App: React.FC = () => {
  const { isDark, toggleDarkMode } = useDarkMode();
  
  // Toast State
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const toast = (type: 'success' | 'error' | 'info', title: string, message: string) => {
    const id = Math.random().toString(36).substring(2, 9);
    setToasts((prev) => [...prev, { id, type, title, message }]);
  };
  const removeToast = (id: string) => {
    setToasts((prev) => prev.filter(t => t.id !== id));
  };

  // Settings State
  const [config, setConfig] = useState<Config>({
    proxy_addr: ':15500',
    api_addr: ':15501',
    mcp_addr: ':15502',
    mcp_enabled: false,
    history_limit: 500,
    max_response_size: 1048576,
    default_page_size: 50
  });
  const [originalConfig, setOriginalConfig] = useState<Config | null>(null);

  // Traffic Hook
  const {
    entries, totalEntries, currentPage, proxyAddr, version, mcpSessions, mcpEnabled, filter, setFilter,
    filteredEntries, fetchTraffic, fetchStatus, clearTraffic,
    setEntries, setTotalEntries, currentPageRef, pageSizeRef
  } = useTraffic(config, toast);

  // UI State
  const [currentView, setCurrentView] = useState<'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about'>(() => {
    const saved = localStorage.getItem('glance-current-view');
    return (saved as any) || 'traffic';
  });

  useEffect(() => {
    localStorage.setItem('glance-current-view', currentView);
  }, [currentView]);

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
  const [recordingFilter, setRecordingFilter] = useState('');
  const isRecordingRef = React.useRef(isRecording);
  const recordingFilterRef = React.useRef(recordingFilter);

  useEffect(() => {
    isRecordingRef.current = isRecording;
  }, [isRecording]);

  useEffect(() => {
    recordingFilterRef.current = recordingFilter;
  }, [recordingFilter]);

  const handleToggleRecording = () => {
    if (isRecording) {
      // STOP recording -> Open editor with current recorded steps
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
      // START recording
      setRecordedEntries([]);
      setIsRecording(true);
      toast('info', 'Recording Started', 'All incoming requests will be added to the new scenario sequence.');
    }
  };

  const fetchScenarios = async () => {
    setIsLoadingScenarios(true);
    try {
      const data = await apiFetch('/api/scenarios');
      setScenarios(data || []);
    } catch (error) {
      toast('error', 'Fetch Scenarios Failed', String(error));
    } finally {
      setIsLoadingScenarios(false);
    }
  };

  const handleSaveScenario = async (scenario: Scenario) => {
    try {
      const isUpdate = !!scenario.id;
      const url = isUpdate ? `/api/scenarios/${scenario.id}` : '/api/scenarios';
      const method = isUpdate ? 'PUT' : 'POST';

      await apiFetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(scenario)
      });

      await fetchScenarios();
      setIsScenarioEditorOpen(false);
      setSelectedScenario(null);
      setCurrentView('scenarios');
      toast('success', isUpdate ? 'Scenario Updated' : 'Scenario Created', `Successfully saved "${scenario.name}".`);
    } catch (error) {
      toast('error', 'Save Failed', String(error));
    }
  };

  const handleDeleteScenario = (id: string) => {
    const scenario = scenarios.find(s => s.id === id);
    if (scenario) {
      setScenarioToDelete(scenario);
      setIsDeleteScenarioModalOpen(true);
    }
  };

  const handleConfirmDeleteScenario = async () => {
    if (!scenarioToDelete) return;
    try {
      await apiFetch(`/api/scenarios/${scenarioToDelete.id}`, { method: 'DELETE' });
      setScenarios(prev => prev.filter(s => s.id !== scenarioToDelete.id));
      setIsDeleteScenarioModalOpen(false);
      setScenarioToDelete(null);
      toast('success', 'Scenario Deleted', 'The scenario has been removed.');
    } catch (error) {
      toast('error', 'Delete Failed', String(error));
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

  const stopResizing = () => setIsResizing(false);

  const resize = (e: MouseEvent) => {
    if (isResizing) {
      const sidebarWidth = isSidebarCollapsed ? 80 : 256;
      const availableWidth = window.innerWidth - sidebarWidth;
      const mouseXFromRight = window.innerWidth - e.clientX;
      const percentage = (mouseXFromRight / availableWidth) * 100;
      if (percentage > 20 && percentage < 80) setDetailsWidth(percentage);
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
  
  // Integration State
  const [javaProcesses, setJavaProcesses] = useState<JavaProcess[]>([]);
  const [androidDevices, setAndroidDevices] = useState<AndroidDevice[]>([]);
  const [isLoadingJava, setIsLoadingJava] = useState(false);
  const [isLoadingAndroid, setIsLoadingAndroid] = useState(false);
  const [terminalScript, setTerminalScript] = useState('');
  
  // Rules State
  const [rules, setRules] = useState<Rule[]>([]);
  const [isLoadingRules, setIsLoadingRules] = useState(false);

  // Scenarios State
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [isLoadingScenarios, setIsLoadingScenarios] = useState(false);

  // --- API Actions ---

  const apiFetch = async (url: string, options?: RequestInit) => {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000); // 10s timeout

    try {
      const res = await fetch(url, {
        ...options,
        signal: controller.signal
      });
      clearTimeout(timeoutId);
      if (!res.ok) throw new Error(await res.text());
      if (res.status === 204) return null;
      return res.json();
    } catch (error) {
      clearTimeout(timeoutId);
      throw error;
    }
  };

  const handleClear = async () => {
    const success = await clearTraffic();
    if (success) {
      setSelectedEntry(null);
      setIsClearModalOpen(false);
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

  const handleAddToScenario = async (entry: TrafficEntry, scenarioId: string | 'new') => {
    if (scenarioId === 'new') {
      setPendingEntry(entry);
      setQuickScenarioName(`New Scenario ${new Date().toLocaleTimeString()}`);
      setQuickScenarioDesc('');
      setIsQuickCreateModalOpen(true);
      return;
    }

    // Add to existing scenario
    const existing = scenarios.find(s => s.id === scenarioId);
    if (!existing) return;

    const updated: Scenario = {
      ...existing,
      steps: [
        ...(existing.steps || []),
        {
          id: '',
          traffic_entry_id: entry.id,
          order: (existing.steps?.length || 0) + 1,
          notes: ''
        }
      ]
    };

    try {
      await apiFetch(`/api/scenarios/${scenarioId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updated)
      });
      await fetchScenarios();
      toast('success', 'Added to Scenario', `Successfully added to "${existing.name}".`);
    } catch (error) {
      toast('error', 'Add Failed', String(error));
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
        notes: ''
      }],
      variable_mappings: [],
      created_at: new Date().toISOString()
    };

    try {
      await apiFetch('/api/scenarios', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newScenario)
      });
      await fetchScenarios();
      setIsQuickCreateModalOpen(false);
      setPendingEntry(null);
      toast('success', 'Scenario Created', `Successfully created "${quickScenarioName}".`);
    } catch (error) {
      toast('error', 'Creation Failed', String(error));
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

  const fetchConfig = async () => {
    try {
      const data = await apiFetch('/api/config');
      setConfig(data);
      setOriginalConfig(data);
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
      
      const needsRestart = originalConfig && (
        newConfig.proxy_addr !== originalConfig.proxy_addr ||
        newConfig.api_addr !== originalConfig.api_addr ||
        newConfig.mcp_addr !== originalConfig.mcp_addr ||
        newConfig.mcp_enabled !== originalConfig.mcp_enabled
      );

      setConfig(newConfig);
      setOriginalConfig(newConfig);

      if (needsRestart) {
        toast('info', 'Restart Required', 'Some changes will take effect only after you restart the Glance.');
      } else {
        toast('success', 'Settings Saved', 'Your configuration has been updated successfully.');
      }
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
              toast('info', 'Response Paused', `${entry.url} - Ready for edit.`);
            } else {
              setIsRequestEditorOpen(true);
              toast('info', 'Request Paused', `${entry.url} - Ready for edit.`);
            }
          } else {
            entry = msg;
          }

          setEntries((prev) => {
            // Check if entry already exists (e.g. was intercepted as request, now finishing as response)
            const exists = prev.some(e => e.id === entry.id);
            if (exists) {
              return prev.map(e => e.id === entry.id ? entry : e);
            }

            // Only increment total for new entries
            setTotalEntries(curr => curr + 1);

            if (currentPageRef.current === 1) {
              const updated = [...prev, entry];
              const pageSize = pageSizeRef.current;
              if (updated.length > pageSize) return updated.slice(1);
              return updated;
            }
            return prev;
          });

          if (isRecordingRef.current) {
            const filter = recordingFilterRef.current.toLowerCase();
            if (!filter || entry.url.toLowerCase().includes(filter)) {
              setRecordedEntries(prev => [...prev, entry]);
            }
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      ws.onclose = () => console.log('WebSocket disconnected.');
    };

    init();
    return () => { if (ws) ws.close(); };
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
    if (currentView === 'scenarios') {
      fetchScenarios();
    }
  }, [currentView]);

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
          onClearTraffic={() => setIsClearModalOpen(true)} 
          isDark={isDark}
          toggleDarkMode={toggleDarkMode}
          isRecording={isRecording}
          onToggleRecording={handleToggleRecording}
          recordedCount={recordedEntries.length}
          recordingFilter={recordingFilter}
          setRecordingFilter={setRecordingFilter}
          onShowMCP={() => setIsMCPDocsOpen(true)}
        />

        <main className="flex-1 flex overflow-hidden" onClick={() => setSelectedEntry(null)}>
          {currentView === 'traffic' && (
            <>
              <div className="flex-1 flex flex-col min-w-0 bg-white dark:bg-slate-900 transition-colors">
                <TrafficList entries={filteredEntries} selectedEntry={selectedEntry} onSelect={setSelectedEntry} />
                
                {/* Pagination Controls */}
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
                      onAddToScenario={handleAddToScenario} 
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
              javaProcesses={javaProcesses} androidDevices={androidDevices} isLoadingJava={isLoadingJava} isLoadingAndroid={isLoadingAndroid}
              onFetchJava={fetchJavaProcesses} onFetchAndroid={fetchAndroidDevices} onInterceptJava={handleInterceptJava} onInterceptAndroid={handleInterceptAndroid} onClearAndroid={handleClearAndroid} onPushAndroidCert={handlePushAndroidCert}
              onShowTerminalDocs={() => setIsTerminalDocsOpen(true)}
            />
          )}

          {currentView === 'rules' && (
            <RulesView 
              rules={rules} isLoading={isLoadingRules} onDelete={handleDeleteRule} onCreate={handleCreateRule}
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
              onDelete={handleDeleteScenario}
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
          <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-md" onClick={() => setIsQuickCreateModalOpen(false)} />
          <div className="relative bg-white dark:bg-slate-900 rounded-3xl shadow-2xl w-full max-w-md overflow-hidden border border-white/50 dark:border-slate-800/50">
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

      <RuleEditor isOpen={isRuleEditorOpen} onClose={() => setIsRuleEditorOpen(false)} rule={selectedRule} onSave={handleUpdateRule} />

      <ScenarioEditor 
        isOpen={isScenarioEditorOpen} 
        onClose={() => { setIsScenarioEditorOpen(false); setSelectedScenario(null); }}
        scenario={selectedScenario}
        onSave={handleSaveScenario}
        availableTraffic={entries} // Simplification: we might need a way to fetch missing entries
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
