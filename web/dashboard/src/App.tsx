import React, { useEffect, useState, useMemo, useRef } from 'react';
import { Trash2, Check, XCircle, Info } from 'lucide-react';
import type { TrafficEntry, Config, JavaProcess } from './types/traffic';

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

  const notify = (type: 'success' | 'error' | 'info', title: string, message: string) => {
    setNotification({ isOpen: true, type, title, message });
  };
  
  // Data State
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [selectedEntry, setSelectedEntry] = useState<TrafficEntry | null>(null);
  const [filter, setFilter] = useState('');
  const [proxyAddr, setProxyAddr] = useState(':8000');
  
  // Integration State
  const [javaProcesses, setJavaProcesses] = useState<JavaProcess[]>([]);
  const [isLoadingJava, setIsLoadingJava] = useState(false);
  const [terminalScript, setTerminalScript] = useState('');
  
  // Settings State
  const [config, setConfig] = useState<Config>({
    proxy_addr: ':8000',
    api_addr: ':8081',
    mcp_addr: ':8082',
    mcp_enabled: false,
    history_limit: 500,
    max_response_size: 1048576
  });

  const historyLimitRef = useRef(config.history_limit);
  useEffect(() => {
    historyLimitRef.current = config.history_limit;
  }, [config.history_limit]);

  // --- API Actions (Using Native Fetch) ---

  const apiFetch = async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  };

  const fetchTraffic = async () => {
    try {
      const data: TrafficEntry[] = await apiFetch('/api/traffic');
      const limit = historyLimitRef.current;
      if (limit > 0 && data.length > limit) {
        setEntries(data.slice(data.length - limit));
      } else {
        setEntries(data || []);
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
    fetchStatus();
    fetchConfig();
    fetchTraffic(); // Initial fetch

    // WebSocket setup
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/traffic`;
    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      try {
        const entry: TrafficEntry = JSON.parse(event.data);
        setEntries((prev) => {
          const updated = [...prev, entry];
          const limit = historyLimitRef.current;
          if (limit > 0 && updated.length > limit) {
            return updated.slice(updated.length - limit);
          }
          return updated;
        });
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected. Reconnecting in 3s...');
      setTimeout(() => {
        // Simple reconnection logic could go here if needed
      }, 3000);
    };

    return () => ws.close();
  }, []);

  useEffect(() => {
    if (currentView === 'integrations') {
      fetchJavaProcesses();
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
              <TrafficList 
                entries={filteredEntries} 
                selectedEntry={selectedEntry} 
                onSelect={setSelectedEntry} 
              />
              {selectedEntry && <DetailsPanel entry={selectedEntry} />}
            </>
          )}

          {currentView === 'integrations' && (
            <IntegrationsView 
              javaProcesses={javaProcesses}
              isLoadingJava={isLoadingJava}
              terminalScript={terminalScript}
              onFetchJava={fetchJavaProcesses}
              onInterceptJava={handleInterceptJava}
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
