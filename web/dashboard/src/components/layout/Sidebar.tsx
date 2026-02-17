import React from 'react';
import { Globe, Sparkles, Settings, Code, ShieldAlert } from 'lucide-react';

interface SidebarProps {
  currentView: 'traffic' | 'integrations' | 'settings' | 'rules';
  setCurrentView: (view: 'traffic' | 'integrations' | 'settings' | 'rules') => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ currentView, setCurrentView }) => {
  return (
    <aside className="w-64 flex flex-col py-6 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transition-colors">
      <div 
        className="px-6 mb-10 flex items-center gap-3 cursor-pointer" 
        onClick={() => setCurrentView('traffic')}
      >
        <div className="p-2 bg-blue-600 rounded-xl shadow-lg shadow-blue-200 dark:shadow-none">
          <Sparkles className="text-white" size={24} />
        </div>
        <span className="font-bold text-lg tracking-tight text-slate-800 dark:text-slate-100">Glance</span>
      </div>
      
      <nav className="flex flex-col gap-1 px-4">
        <button 
          onClick={() => setCurrentView('traffic')}
          className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm ${
            currentView === 'traffic' 
            ? 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' 
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800'
          }`}
        >
          <Globe size={20} />
          <span>Traffic Inspector</span>
        </button>
        
        <button 
          onClick={() => setCurrentView('integrations')}
          className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm ${
            currentView === 'integrations' 
            ? 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' 
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800'
          }`}
        >
          <Code size={20} />
          <span>Integrations</span>
        </button>

        <button 
          onClick={() => setCurrentView('rules')}
          className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm ${
            currentView === 'rules' 
            ? 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' 
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800'
          }`}
        >
          <ShieldAlert size={20} />
          <span>Breakpoint Rules</span>
        </button>
        
        <button 
          onClick={() => setCurrentView('settings')}
          className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm ${
            currentView === 'settings' 
            ? 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' 
            : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800'
          }`}
        >
          <Settings size={20} />
          <span>System Settings</span>
        </button>
      </nav>

      <div className="mt-auto px-6 pt-6 border-t border-slate-100 dark:border-slate-800">
        <div className="flex flex-col gap-1">
          <span className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-widest">Version</span>
          <span className="text-xs font-mono text-slate-500 dark:text-slate-400">v0.1.0-alpha</span>
        </div>
      </div>
    </aside>
  );
};
