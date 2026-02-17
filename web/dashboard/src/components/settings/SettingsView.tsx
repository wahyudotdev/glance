import React from 'react';
import { HelpCircle } from 'lucide-react';
import type { Config } from '../../types/traffic';

interface SettingsViewProps {
  config: Config;
  setConfig: (config: Config) => void;
  onSave: (config: Config) => void;
  onReset: () => void;
  onShowMCP: () => void;
}

export const SettingsView: React.FC<SettingsViewProps> = ({ config, setConfig, onSave, onReset, onShowMCP }) => {
  return (
    <div className="flex-1 p-12 bg-slate-50 dark:bg-slate-950 overflow-y-auto transition-colors">
      <div className="max-w-2xl mx-auto">
        <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100 mb-8">System Settings</h2>
        <div className="space-y-6">
          <div className="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm transition-colors">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider">Network Ports</h3>
              <span className="text-[9px] font-bold text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 px-2 py-0.5 rounded border border-amber-100 dark:border-amber-800/30 uppercase tracking-tighter">Requires Restart</span>
            </div>
            <div className="space-y-4">
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">Proxy Address</label>
                <input 
                  type="text" 
                  value={config.proxy_addr}
                  onChange={(e) => setConfig({...config, proxy_addr: e.target.value})}
                  className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors"
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">API / Dashboard Address</label>
                <input 
                  type="text" 
                  value={config.api_addr}
                  onChange={(e) => setConfig({...config, api_addr: e.target.value})}
                  className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors"
                />
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm transition-colors">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="flex flex-col">
                  <h3 className="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wider">MCP Server</h3>
                  <p className="text-xs text-slate-400 dark:text-slate-500 mt-1">Model Context Protocol for AI Agent integration</p>
                </div>
                <button 
                  onClick={onShowMCP}
                  className="p-1.5 text-slate-400 dark:text-slate-500 hover:text-indigo-600 dark:hover:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 rounded-lg transition-all"
                  title="View Documentation"
                >
                  <HelpCircle size={16} />
                </button>
                <span className="text-[9px] font-bold text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 px-2 py-0.5 rounded border border-amber-100 dark:border-amber-800/30 uppercase tracking-tighter">Requires Restart</span>
              </div>
              <button 
                onClick={() => setConfig({...config, mcp_enabled: !config.mcp_enabled})}
                className={`w-12 h-6 rounded-full transition-all relative ${config.mcp_enabled ? 'bg-blue-600' : 'bg-slate-200 dark:bg-slate-700'}`}
              >
                <div className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-all ${config.mcp_enabled ? 'left-7' : 'left-1'}`} />
              </button>
            </div>
            <div className="flex flex-col gap-1.5 mt-4">
              <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">MCP Server Address (SSE)</label>
              <input 
                type="text" 
                value={config.mcp_addr}
                onChange={(e) => setConfig({...config, mcp_addr: e.target.value})}
                disabled={!config.mcp_enabled}
                className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors disabled:opacity-50"
              />
            </div>
          </div>

          <div className="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm transition-colors">
            <h3 className="text-sm font-bold text-slate-800 dark:text-slate-200 mb-4 uppercase tracking-wider">Data Management</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">History Limit (Max Entries)</label>
                <input 
                  type="number" 
                  value={config.history_limit}
                  onChange={(e) => setConfig({...config, history_limit: parseInt(e.target.value) || 0})}
                  className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors"
                  placeholder="500"
                />
                <p className="text-[10px] text-slate-400 dark:text-slate-500 italic">Auto-remove oldest entries when limit is reached.</p>
              </div>
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">Max Response Size (Bytes)</label>
                <input 
                  type="number" 
                  value={config.max_response_size}
                  onChange={(e) => setConfig({...config, max_response_size: parseInt(e.target.value) || 0})}
                  className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors"
                  placeholder="1048576"
                />
                <p className="text-[10px] text-slate-400 dark:text-slate-500 italic">Default: 1,048,576 (1 MB). 0 to disable limit.</p>
              </div>
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 dark:text-slate-400 uppercase">Default Page Size</label>
                <input 
                  type="number" 
                  value={config.default_page_size}
                  onChange={(e) => setConfig({...config, default_page_size: parseInt(e.target.value) || 0})}
                  className="px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-mono dark:text-slate-200 transition-colors"
                  placeholder="50"
                />
                <p className="text-[10px] text-slate-400 dark:text-slate-500 italic">Number of entries to load per page.</p>
              </div>
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <button 
              onClick={onReset}
              className="px-6 py-2.5 text-sm font-bold text-slate-500 dark:text-slate-400 hover:bg-white dark:hover:bg-slate-800 rounded-xl transition-all"
            >
              Reset
            </button>
            <button 
              onClick={() => onSave(config)}
              className="px-8 py-2.5 bg-blue-600 text-white rounded-xl font-bold text-sm hover:bg-blue-700 shadow-lg shadow-blue-200 dark:shadow-none active:scale-95 transition-all"
            >
              Save Settings
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
