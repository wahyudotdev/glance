import React from 'react';
import { Search, Trash2, Zap, Sun, Moon } from 'lucide-react';

interface HeaderProps {
  proxyAddr: string;
  filter: string;
  setFilter: (filter: string) => void;
  onClearTraffic: () => void;
  isDark: boolean;
  toggleDarkMode: () => void;
}

export const Header: React.FC<HeaderProps> = ({ 
  proxyAddr, 
  filter, 
  setFilter, 
  onClearTraffic, 
  isDark,
  toggleDarkMode
}) => {
  return (
    <header className="h-16 flex items-center justify-between px-8 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 shadow-sm z-10 transition-colors">
      <div className="flex items-center gap-4">
        <h1 className="text-lg font-bold tracking-tight text-slate-800 dark:text-slate-100">Glance</h1>
        <div className="flex items-center gap-3 px-3 py-1.5 bg-slate-100 dark:bg-slate-800 rounded-full text-[11px] font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider transition-colors">
          <span className="flex items-center gap-1.5"><div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" /> Proxy {proxyAddr}</span>
          <span className="w-px h-3 bg-slate-300 dark:bg-slate-700" />
          <span className="flex items-center gap-1.5 text-blue-600 dark:text-blue-400"><Zap size={12} /> MITM Active</span>
        </div>
      </div>
      
      <div className="flex items-center gap-3">
        <button 
          onClick={toggleDarkMode}
          className="p-2 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-all"
          title={isDark ? "Switch to Light Mode" : "Switch to Night Mode"}
        >
          {isDark ? <Sun size={18} /> : <Moon size={18} />}
        </button>

        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
          <input 
            type="text" 
            placeholder="Filter requests..." 
            className="pl-10 pr-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 dark:text-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all w-64"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
          />
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
