import React, { useState } from 'react';
import { Search, Trash2, Zap, Sun, Moon, Circle, Square, ListFilter } from 'lucide-react';

interface HeaderProps {
  proxyAddr: string;
  mcpSessions: number;
  mcpEnabled: boolean;
  filter: string;
  setFilter: (filter: string) => void;
  methodFilter: string;
  setMethodFilter: (method: string) => void;
  onClearTraffic: () => void;
  isDark: boolean;
  toggleDarkMode: () => void;
  isRecording: boolean;
  onToggleRecording: () => void;
  recordedCount: number;
  onShowMCP: () => void;
}

export const Header: React.FC<HeaderProps> = ({ 
  proxyAddr, 
  mcpSessions,
  mcpEnabled,
  filter, 
  setFilter, 
  methodFilter,
  setMethodFilter,
  onClearTraffic, 
  isDark,
  toggleDarkMode,
  isRecording,
  onToggleRecording,
  recordedCount,
  onShowMCP
}) => {
  const methods = ['ALL', 'GET', 'POST', 'PUT', 'DELETE'];
  const [isFilterMenuOpen, setIsFilterMenuOpen] = useState(false);

  return (
    <header className="h-16 flex items-center justify-between px-8 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 shadow-sm z-10 transition-colors">
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-3 px-3 py-1.5 bg-slate-100 dark:bg-slate-800 rounded-full text-[11px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider transition-colors">
          <span className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /> Proxy {proxyAddr}</span>
          {mcpEnabled && (
            <>
              <span className="w-px h-3 bg-slate-300 dark:bg-slate-700" />
              <button 
                onClick={onShowMCP}
                className="flex items-center gap-1.5 text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300 transition-colors group"
                title="View MCP Documentation & Status"
              >
                <Zap size={12} fill="currentColor" className="group-hover:scale-110 transition-transform" /> 
                <span className="font-bold">MCP Server ON</span>
                {mcpSessions > 0 && (
                  <span className="ml-1 px-1.5 py-0.5 bg-indigo-600 text-white rounded-md text-[9px] font-black animate-in zoom-in duration-300 tabular-nums">
                    {mcpSessions} {mcpSessions === 1 ? 'SESSION' : 'SESSIONS'}
                  </span>
                )}
              </button>
            </>
          )}
        </div>

        <button
          onClick={onToggleRecording}
          className={`flex items-center gap-2 px-4 py-1.5 rounded-full text-[11px] font-bold uppercase tracking-wider transition-all ${
            isRecording 
              ? 'bg-rose-500 text-white shadow-lg shadow-rose-200 dark:shadow-none animate-pulse' 
              : 'bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 hover:bg-slate-200 dark:hover:bg-slate-700'
          }`}
        >
          {isRecording ? (
            <>
              <Square size={10} fill="currentColor" />
              Stop Recording ({recordedCount})
            </>
          ) : (
            <>
              <Circle size={10} fill="currentColor" className="text-rose-500" />
              Start Recording
            </>
          )}
        </button>
      </div>
      
      <div className="flex items-center gap-3">
        <button 
          onClick={toggleDarkMode}
          className="p-2 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-all"
          title={isDark ? "Switch to Light Mode" : "Switch to Night Mode"}
        >
          {isDark ? <Sun size={18} /> : <Moon size={18} />}
        </button>

        <div className="flex items-center gap-2">
          <div className="relative flex items-center group">
            <Search className="absolute left-3 text-slate-400 group-focus-within:text-blue-500 transition-colors" size={16} />
            <input 
              type="text" 
              placeholder="Filter requests..." 
              className="pl-10 pr-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 dark:text-slate-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all w-64 font-medium"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>

          <div className="relative">
            <button 
              onClick={() => setIsFilterMenuOpen(!isFilterMenuOpen)}
              className={`p-2 rounded-xl border transition-all relative ${
                methodFilter !== 'ALL' || isFilterMenuOpen
                ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800 text-blue-600 dark:text-blue-400' 
                : 'bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300'
              }`}
              title="Configure Filters"
            >
              <ListFilter size={18} />
              {methodFilter !== 'ALL' && (
                <span className="absolute -top-1 -right-1 w-2.5 h-2.5 bg-blue-600 rounded-full border-2 border-white dark:border-slate-900" />
              )}
            </button>

            {isFilterMenuOpen && (
              <>
                <div className="fixed inset-0 z-20" onClick={() => setIsFilterMenuOpen(false)} />
                <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-xl z-30 py-2 animate-in zoom-in-95 duration-100 origin-top-right">
                  <div className="px-4 py-2 border-b border-slate-50 dark:border-slate-700/50 mb-1">
                    <span className="text-[10px] font-black uppercase text-slate-400 tracking-widest">Request Method</span>
                  </div>
                  <div className="px-2 space-y-0.5">
                    {methods.map(m => (
                      <button
                        key={m}
                        onClick={() => {
                          setMethodFilter(m);
                          setIsFilterMenuOpen(false);
                        }}
                        className={`w-full text-left px-3 py-2 rounded-lg text-[11px] font-bold transition-all ${
                          methodFilter === m 
                          ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400' 
                          : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-700/50 hover:text-slate-700 dark:hover:text-slate-200'
                        }`}
                      >
                        {m}
                      </button>
                    ))}
                  </div>
                </div>
              </>
            )}
          </div>
        </div>

        <button 
          onClick={onClearTraffic}
          className="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-all" 
          title="Clear Logs"
        >
          <Trash2 size={18} />
        </button>
      </div>
    </header>
  );
};
