import React from 'react';
import { Globe, Sparkles, Settings, Code, ShieldAlert, ChevronLeft, ChevronRight, ListPlus, Info } from 'lucide-react';
import { GlanceLogo } from '../ui/GlanceLogo';

interface SidebarProps {
  currentView: 'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about';
  setCurrentView: (view: 'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about') => void;
  isCollapsed: boolean;
  onToggleCollapse: () => void;
  version: string;
  onShowChangelog: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ 
  currentView, 
  setCurrentView, 
  isCollapsed, 
  onToggleCollapse,
  version,
  onShowChangelog
}) => {
  return (
    <aside className={`flex flex-col py-6 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transition-all duration-300 relative ${isCollapsed ? 'w-20' : 'w-64'}`}>
      {/* Collapse Toggle Button */}
      <button 
        onClick={onToggleCollapse}
        className="absolute -right-3 top-12 w-6 h-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-full flex items-center justify-center text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 shadow-sm z-50 transition-all"
      >
        {isCollapsed ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
      </button>

      <div 
        className={`px-6 mb-10 flex items-center gap-3 cursor-pointer overflow-hidden ${isCollapsed ? 'justify-center px-0' : ''}`} 
        onClick={() => setCurrentView('traffic')}
      >
        <GlanceLogo size={40} className="shrink-0 shadow-lg shadow-blue-200 dark:shadow-none rounded-full" />
        {!isCollapsed && (
          <div className="flex flex-col min-w-0">
            <span className="font-bold text-lg tracking-tight text-slate-800 dark:text-slate-100 leading-tight">Glance</span>
            <span className="text-[9px] text-slate-400 dark:text-slate-500 font-medium leading-tight mt-0.5">Let Your AI Understand Every Request at a Glance.</span>
          </div>
        )}
      </div>
      
      <nav className={`flex flex-col gap-1 px-4 ${isCollapsed ? 'px-2' : ''}`}>
        {[
          { id: 'traffic', label: 'Traffic Inspector', icon: <Globe size={20} /> },
          { id: 'integrations', label: 'Integrations', icon: <Code size={20} /> },
          { id: 'scenarios', label: 'Traffic Scenarios', icon: <ListPlus size={20} /> },
          { id: 'rules', label: 'Breakpoint Rules', icon: <ShieldAlert size={20} /> },
          { id: 'settings', label: 'System Settings', icon: <Settings size={20} /> },
          { id: 'about', label: 'About Glance', icon: <Info size={20} /> },
        ].map((item) => (
          <button 
            key={item.id}
            onClick={() => setCurrentView(item.id as 'traffic' | 'integrations' | 'settings' | 'rules' | 'scenarios' | 'about')}
            className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm whitespace-nowrap ${
              currentView === item.id 
              ? 'text-blue-600 bg-blue-50 dark:bg-blue-900/20' 
              : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800'
            } ${isCollapsed ? 'justify-center px-0' : ''}`}
            title={isCollapsed ? item.label : ''}
          >
            <div className="shrink-0">{item.icon}</div>
            {!isCollapsed && <span>{item.label}</span>}
          </button>
        ))}
      </nav>

      <div className={`mt-auto px-6 pt-6 border-t border-slate-100 dark:border-slate-800 overflow-hidden ${isCollapsed ? 'px-0 flex flex-col items-center' : ''}`}>
        <div className="flex flex-col gap-1">
          <span className={`text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-widest ${isCollapsed ? 'hidden' : ''}`}>Version</span>
          <div className={`flex items-center justify-between ${isCollapsed ? 'justify-center' : ''}`}>
            <span className={`text-xs font-mono text-slate-500 dark:text-slate-400 ${isCollapsed ? 'scale-75' : ''}`}>{version.startsWith('v') ? version : `v${version}`}</span>
            {!isCollapsed && (
              <button 
                onClick={onShowChangelog}
                className="text-[10px] font-black text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/20 px-2 py-0.5 rounded border border-indigo-100 dark:border-indigo-800/30 uppercase tracking-tighter hover:bg-indigo-600 hover:text-white dark:hover:bg-indigo-500 transition-all"
              >
                Changelogs
              </button>
            )}
          </div>
        </div>
        {isCollapsed && (
          <button 
            onClick={onShowChangelog}
            className="mt-2 p-1.5 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 rounded-lg transition-all"
            title="Changelogs"
          >
            <Sparkles size={14} />
          </button>
        )}
      </div>
    </aside>
  );
};

