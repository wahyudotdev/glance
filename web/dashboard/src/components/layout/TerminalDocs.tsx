import React from 'react';
import { Terminal, X, Code } from 'lucide-react';

interface TerminalDocsProps {
  isOpen: boolean;
  onClose: () => void;
  terminalScript: string;
}

export const TerminalDocs: React.FC<TerminalDocsProps> = ({ isOpen, onClose, terminalScript }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[150] flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300">
      <div className="bg-white dark:bg-slate-900 rounded-3xl shadow-2xl w-full max-w-lg overflow-hidden animate-in zoom-in-95 duration-300 flex flex-col border border-transparent dark:border-slate-800 transition-colors">
        <div className="p-6 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-950/50 transition-colors">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-xl">
              <Terminal size={20} />
            </div>
            <div>
              <h3 className="font-bold text-slate-800 dark:text-slate-100">Manual Setup</h3>
              <p className="text-xs text-slate-500 dark:text-slate-400">Environment Variables</p>
            </div>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all">
            <X size={20} />
          </button>
        </div>

        <div className="p-8 space-y-6">
          <section className="space-y-3">
            <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200 flex items-center gap-2">
              <Code size={16} className="text-blue-500" />
              Manual Configuration
            </h4>
            <p className="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
              If the one-liner doesn't work or you prefer manual setup, export these variables in your shell:
            </p>
            <div className="relative group">
              <pre className="bg-slate-900 text-indigo-200 p-4 rounded-xl text-[10px] font-mono overflow-x-auto border border-slate-800 max-h-64">
                {terminalScript || '# Fetching setup script...'}
              </pre>
            </div>
          </section>
          
          <div className="bg-amber-50 dark:bg-amber-900/10 border border-amber-100 dark:border-amber-900/30 rounded-xl p-4 transition-colors">
            <p className="text-[11px] text-amber-700 dark:text-amber-500 leading-relaxed italic">
              Note: These variables only affect the current terminal session.
            </p>
          </div>
        </div>

        <div className="p-4 bg-slate-50 dark:bg-slate-950 border-t border-slate-100 dark:border-slate-800 flex justify-center transition-colors">
          <button 
            onClick={onClose}
            className="px-8 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 rounded-xl text-sm font-bold hover:bg-slate-50 dark:hover:bg-slate-700 transition-all"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
