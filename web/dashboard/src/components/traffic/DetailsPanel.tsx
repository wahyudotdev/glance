import React, { useState } from 'react';
import { FileText, Copy, Check } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';
import { generateCurl } from '../../lib/curl';

interface DetailsPanelProps {
  entry: TrafficEntry;
  width?: number;
}

export const DetailsPanel: React.FC<DetailsPanelProps> = ({ entry, width = 450 }) => {
  const [activeTab, setActiveTab] = useState<'headers' | 'body' | 'curl'>('headers');
  const [copied, setCopied] = useState(false);

  const handleCopyCurl = () => {
    const curl = generateCurl(entry);
    navigator.clipboard.writeText(curl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div 
      className="bg-white border-l border-slate-200 flex flex-col shadow-2xl z-20 flex-shrink-0"
      style={{ width: `${width}px` }}
    >
      <div className="p-6 border-b border-slate-100 bg-slate-50/50">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-bold text-slate-800 flex items-center gap-2 uppercase tracking-tight">
            <FileText size={16} className="text-blue-500" /> Request Details
            {entry.duration > 0 && (
              <span className="ml-2 px-2 py-0.5 bg-slate-100 text-slate-500 rounded text-[10px] font-mono">
                {(entry.duration / 1000000).toFixed(1)}ms
              </span>
            )}
          </h2>
          <button 
            onClick={handleCopyCurl}
            className="flex items-center gap-1.5 px-3 py-1.5 bg-white border border-slate-200 rounded-lg text-xs font-semibold text-slate-600 hover:border-blue-500 hover:text-blue-600 transition-all shadow-sm active:scale-95"
          >
            {copied ? <Check size={14} className="text-emerald-500" /> : <Copy size={14} />}
            {copied ? 'Copied!' : 'Copy cURL'}
          </button>
        </div>
        <div className="font-mono text-[12px] bg-slate-900 text-slate-300 p-3 rounded-lg break-all leading-relaxed shadow-inner border border-slate-800">
          <span className="text-blue-400 font-bold uppercase mr-2">{entry.method}</span>
          {entry.url}
        </div>
      </div>

      <div className="flex px-6 pt-2 border-b border-slate-100 bg-white">
        {(['headers', 'body', 'curl'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-3 text-[11px] font-bold uppercase tracking-wider transition-all border-b-2 -mb-[1px] ${
              activeTab === tab 
              ? 'border-blue-600 text-blue-600' 
              : 'border-transparent text-slate-400 hover:text-slate-600'
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {activeTab === 'headers' && (
          <div className="space-y-6">
            <section>
              <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Request Headers</h3>
              <div className="space-y-1.5">
                {Object.entries(entry.request_headers).map(([key, values]) => (
                  <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 last:border-0">
                    <span className="font-bold text-slate-600 min-w-[120px] shrink-0">{key}</span>
                    <span className="text-slate-500 break-all font-mono">{values.join(', ')}</span>
                  </div>
                ))}
              </div>
            </section>
            
            {entry.response_headers && (
              <section>
                <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Response Headers</h3>
                <div className="space-y-1.5">
                  {Object.entries(entry.response_headers).map(([key, values]) => (
                    <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 last:border-0">
                      <span className="font-bold text-slate-600 min-w-[120px] shrink-0">{key}</span>
                      <span className="text-slate-500 break-all font-mono">{values.join(', ')}</span>
                    </div>
                  ))}
                </div>
              </section>
            )}
          </div>
        )}

        {activeTab === 'body' && (
          <div className="h-full flex flex-col gap-4">
            <section className="flex-1 min-h-0 flex flex-col">
              <h3 className="text-[10px] font-black uppercase text-slate-400 mb-3 tracking-[0.2em]">Response Body</h3>
              <div className="flex-1 bg-slate-50 rounded-xl p-4 font-mono text-[12px] overflow-auto border border-slate-100 shadow-inner">
                {entry.response_body ? (
                  <pre className="text-slate-600 whitespace-pre-wrap leading-relaxed">
                    {(() => {
                      try {
                        return JSON.stringify(JSON.parse(entry.response_body), null, 2);
                      } catch {
                        return entry.response_body;
                      }
                    })()}
                  </pre>
                ) : (
                  <span className="text-slate-300 italic">No response body captured</span>
                )}
              </div>
            </section>
          </div>
        )}

        {activeTab === 'curl' && (
          <div className="space-y-4">
            <h3 className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">cURL Command</h3>
            <pre className="bg-slate-900 text-blue-300 p-5 rounded-xl text-[12px] font-mono whitespace-pre-wrap leading-relaxed shadow-xl">
              {generateCurl(entry)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};
