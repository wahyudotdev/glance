import React from 'react';
import { Heart, Sparkles, Shield, Info } from 'lucide-react';

export const AboutView: React.FC = () => {
  return (
    <div className="flex-1 p-12 bg-slate-50 dark:bg-slate-950 overflow-y-auto transition-colors">
      <div className="max-w-4xl mx-auto space-y-12">
        <section>
          <div className="flex items-center gap-4 mb-6">
            <div className="w-12 h-12 bg-blue-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-200 dark:shadow-none">
              <Sparkles className="text-white" size={24} />
            </div>
            <div>
              <h2 className="text-3xl font-black text-slate-800 dark:text-slate-100 tracking-tight">Glance</h2>
              <p className="text-sm text-slate-500 dark:text-slate-400 font-medium">Let Your AI Understand Every Request at a Glance.</p>
            </div>
          </div>
          
          <div className="bg-white dark:bg-slate-900 p-8 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-6">
            <p className="text-slate-600 dark:text-slate-300 leading-relaxed">
              Glance is a specialized MITM (Man-in-the-Middle) proxy designed for <strong>AI Agents</strong> and developers. 
              It provides a real-time, high-fidelity view of network activity, allowing both humans and AI to inspect, 
              intercept, and mock HTTP/HTTPS traffic with ease.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="flex items-start gap-4">
                <div className="mt-1 p-2 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-xl">
                  <Shield size={18} />
                </div>
                <div>
                  <h4 className="font-bold text-slate-800 dark:text-slate-100">AI-First Design</h4>
                  <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Integrated with Model Context Protocol (MCP) to give AI agents direct access to network state.</p>
                </div>
              </div>
              <div className="flex items-start gap-4">
                <div className="mt-1 p-2 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 rounded-xl">
                  <Info size={18} />
                </div>
                <div>
                  <h4 className="font-bold text-slate-800 dark:text-slate-100">Deep Inspection</h4>
                  <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Automatic HTTPS decryption, structured body parsing, and one-click replay tools.</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="pt-8 border-t border-slate-200 dark:border-slate-800">
          <div className="flex items-center gap-3 mb-8">
            <div className="p-2 bg-rose-50 dark:bg-rose-900/20 rounded-xl text-rose-500">
              <Heart size={20} fill="currentColor" />
            </div>
            <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100">Built With Open Source</h2>
          </div>

          <div className="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 divide-y divide-slate-100 dark:divide-slate-800 overflow-hidden">
            {[
              { name: 'GoProxy', desc: 'The core MITM proxy engine.', url: 'https://github.com/elazarl/goproxy' },
              { name: 'GoFiber', desc: 'High-performance web framework for the API.', url: 'https://gofiber.io' },
              { name: 'React', desc: 'Modern library for the dashboard UI.', url: 'https://react.dev' },
              { name: 'Tailwind CSS', desc: 'Utility-first CSS framework for styling.', url: 'https://tailwindcss.com' },
              { name: 'Lucide Icons', desc: 'Beautiful & consistent iconography.', url: 'https://lucide.dev' },
              { name: 'SQLite', desc: 'Lightweight & high-performance persistence.', url: 'https://sqlite.org' },
            ].map((lib) => (
              <a 
                key={lib.name}
                href={lib.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-between p-5 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors group"
              >
                <div className="flex flex-col gap-0.5">
                  <span className="text-sm font-bold text-slate-800 dark:text-slate-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">{lib.name}</span>
                  <span className="text-xs text-slate-500 dark:text-slate-400">{lib.desc}</span>
                </div>
                <div className="text-[10px] font-mono text-slate-400 dark:text-slate-600 group-hover:text-slate-600 dark:group-hover:text-slate-400 transition-colors underline decoration-slate-200 dark:decoration-slate-800 underline-offset-4">
                  Visit Project
                </div>
              </a>
            ))}
          </div>
        </section>

        <div className="flex flex-col items-center justify-center pt-12 pb-8 gap-4 opacity-50">
          <div className="flex items-center gap-2 text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
            <span>Made with</span>
            <Heart size={10} className="text-rose-500" fill="currentColor" />
            <span>by wahyudotdev</span>
          </div>
          <div className="text-[9px] font-mono text-slate-400">Licensed under MIT</div>
        </div>
      </div>
    </div>
  );
};
