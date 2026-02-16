import React from 'react';
import type { Config } from '../../types/traffic';

interface SettingsViewProps {
  config: Config;
  setConfig: (config: Config) => void;
  onSave: (config: Config) => void;
  onReset: () => void;
}

export const SettingsView: React.FC<SettingsViewProps> = ({ config, setConfig, onSave, onReset }) => {
  return (
    <div className="flex-1 p-12 bg-slate-50 overflow-y-auto">
      <div className="max-w-2xl mx-auto">
        <h2 className="text-2xl font-bold text-slate-800 mb-8">System Settings</h2>
        <div className="space-y-6">
          <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
            <h3 className="text-sm font-bold text-slate-800 mb-4 uppercase tracking-wider">Network Ports</h3>
            <div className="space-y-4">
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 uppercase">Proxy Address</label>
                <input 
                  type="text" 
                  value={config.proxy_addr}
                  onChange={(e) => setConfig({...config, proxy_addr: e.target.value})}
                  className="px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono"
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <label className="text-[11px] font-bold text-slate-500 uppercase">API / Dashboard Address</label>
                <input 
                  type="text" 
                  value={config.api_addr}
                  onChange={(e) => setConfig({...config, api_addr: e.target.value})}
                  className="px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono"
                />
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="text-sm font-bold text-slate-800 uppercase tracking-wider">MCP Server</h3>
                <p className="text-xs text-slate-400 mt-1">Model Context Protocol for AI Agent integration</p>
              </div>
              <button 
                onClick={() => setConfig({...config, mcp_enabled: !config.mcp_enabled})}
                className={`w-12 h-6 rounded-full transition-all relative ${config.mcp_enabled ? 'bg-blue-600' : 'bg-slate-200'}`}
              >
                <div className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-all ${config.mcp_enabled ? 'left-7' : 'left-1'}`} />
              </button>
            </div>
            {config.mcp_enabled && (
              <div className="mt-4 flex flex-col gap-1.5 animate-in slide-in-from-top-2 duration-200">
                <label className="text-[11px] font-bold text-slate-500 uppercase">MCP Address (SSE)</label>
                <input 
                  type="text" 
                  value={config.mcp_addr}
                  onChange={(e) => setConfig({...config, mcp_addr: e.target.value})}
                  className="px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm font-mono"
                />
              </div>
            )}
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <button 
              onClick={onReset}
              className="px-6 py-2.5 text-sm font-bold text-slate-500 hover:bg-white rounded-xl transition-all"
            >
              Reset
            </button>
            <button 
              onClick={() => onSave(config)}
              className="px-8 py-2.5 bg-blue-600 text-white rounded-xl font-bold text-sm hover:bg-blue-700 shadow-lg shadow-blue-200 active:scale-95 transition-all"
            >
              Save Settings
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
