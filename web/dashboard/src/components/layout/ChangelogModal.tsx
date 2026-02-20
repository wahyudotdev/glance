import React from 'react';
import { ScrollText, X, ShieldCheck, Layout, Info, Sparkles } from 'lucide-react';

interface ChangelogModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const ChangelogModal: React.FC<ChangelogModalProps> = ({ isOpen, onClose }) => {
  if (!isOpen) return null;

  const changes = [
    {
      date: 'v0.1.3 - 2026-02-20',
      items: [
        { icon: <Layout className="text-blue-500" size={14} />, title: 'UI Fixes', description: 'Fixed cURL command overflow in the details panel for better readability.' },
        { icon: <Info className="text-indigo-500" size={14} />, title: 'About Page', description: 'Introduced a dedicated About page with project overview and library attributions.' },
        { icon: <Sparkles className="text-amber-500" size={14} />, title: 'UI Cleanup', description: 'Streamlined the Integrations view by moving secondary content to the About page.' },
        { icon: <ShieldCheck className="text-emerald-500" size={14} />, title: 'Documentation', description: 'Added Brew installation guide and compatibility matrix to README and Dashboard.' },
        { icon: <ScrollText className="text-indigo-500" size={14} />, title: 'Project Hygiene', description: 'Established a dedicated CHANGELOG.md and integrated it into the sidebar.' }
      ]
    }
  ];

  return (
    <div className="fixed inset-0 z-[150] flex items-center justify-center p-4">
      <div 
        className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300" 
        onClick={onClose} 
      />
      <div className="relative bg-white dark:bg-slate-900 rounded-3xl shadow-2xl w-full max-w-lg overflow-hidden animate-in zoom-in-95 duration-300 flex flex-col border border-slate-100 dark:border-slate-800">
        {/* Header */}
        <div className="p-6 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-950/50 transition-colors">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-xl">
              <ScrollText size={20} />
            </div>
            <div>
              <h3 className="font-bold text-slate-800 dark:text-slate-100 uppercase tracking-tight">What's New</h3>
              <p className="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest">Recent Updates & Milestones</p>
            </div>
          </div>
          <button 
            onClick={onClose} 
            className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all hover:text-slate-600 dark:hover:text-slate-300"
          >
            <X size={20} />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-8 space-y-10 max-h-[60vh] custom-scrollbar">
          {changes.map((group, groupIdx) => (
            <section key={groupIdx} className="relative">
              <div className="flex items-center gap-4 mb-6">
                <div className="h-px flex-1 bg-slate-100 dark:bg-slate-800" />
                <span className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-widest bg-white dark:bg-slate-900 px-3 py-1 rounded-full border border-slate-100 dark:border-slate-800">
                  {group.date}
                </span>
                <div className="h-px flex-1 bg-slate-100 dark:bg-slate-800" />
              </div>

              <div className="space-y-6">
                {group.items.map((item, itemIdx) => (
                  <div key={itemIdx} className="flex gap-4 group">
                    <div className="mt-1 w-8 h-8 rounded-lg bg-slate-50 dark:bg-slate-800 flex items-center justify-center shrink-0 border border-slate-100 dark:border-slate-700/50 group-hover:border-blue-200 dark:group-hover:border-blue-900/50 transition-colors shadow-sm">
                      {item.icon}
                    </div>
                    <div className="space-y-1">
                      <h4 className="text-sm font-bold text-slate-800 dark:text-slate-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                        {item.title}
                      </h4>
                      <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed">
                        {item.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          ))}
        </div>

        {/* Footer */}
        <div className="p-6 bg-slate-50 dark:bg-slate-950 border-t border-slate-100 dark:border-slate-800 flex justify-center transition-colors">
          <button 
            onClick={onClose}
            className="px-10 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-bold shadow-lg shadow-blue-200 dark:shadow-none transition-all active:scale-95"
          >
            Got it!
          </button>
        </div>
      </div>
    </div>
  );
};
