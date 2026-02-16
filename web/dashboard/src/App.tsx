import React, { useEffect, useState, useMemo } from 'react';
import { Trash2, Check } from 'lucide-react';
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
  const [isSettingsSavedModalOpen, setIsSettingsSavedModalOpen] = useState(false);
  
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
    mcp_enabled: false
  });

  // --- API Actions (Using Native Fetch) ---

  const apiFetch = async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  };

  const fetchTraffic = async () => {
    try {
      const data = await apiFetch('/api/traffic');
      setEntries(data || []);
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
    } catch (error) {
      alert('Failed to clear traffic: ' + error);
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
      alert(`Successfully injected proxy into PID ${pid}! Traffic should start appearing.`);
    } catch (error) {
      alert('Failed to intercept: ' + error);
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
      setIsSettingsSavedModalOpen(true);
    } catch (error) {
      alert('Failed to save settings: ' + error);
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
        setEntries((prev) => [...prev, entry]);
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

      {/* Modals */}
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

      <Modal 
        isOpen={isSettingsSavedModalOpen}
        onClose={() => setIsSettingsSavedModalOpen(false)}
        title="Settings Saved"
        description="Your configuration has been updated successfully. Some changes (like port updates) will take effect after restarting the application."
        icon={<Check size={32} />}
        iconBgColor="bg-emerald-50 text-emerald-500"
        confirmLabel="Got it"
        confirmColor="bg-emerald-500 hover:bg-emerald-600 shadow-emerald-200"
        onConfirm={() => setIsSettingsSavedModalOpen(false)}
        showCancel={false}
      />
    </div>
  );
};

export default App;
