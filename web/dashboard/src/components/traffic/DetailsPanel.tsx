import React, { useState } from 'react';
import { FileText, Copy, Check, Eye, Code } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';
import { generateCurl } from '../../lib/curl';

interface DetailsPanelProps {
  entry: TrafficEntry;
  width?: number;
}

export const DetailsPanel: React.FC<DetailsPanelProps> = ({ entry, width = 450 }) => {
  const [activeTab, setActiveTab] = useState<'headers' | 'body' | 'curl'>('headers');
  const [viewMode, setViewMode] = useState<'preview' | 'raw'>('preview');
  const [copied, setCopied] = useState(false);

  const getContentType = () => {
    if (!entry.response_headers) return '';
    const ct = entry.response_headers['Content-Type'] || entry.response_headers['content-type'] || [];
    return ct[0] || '';
  };

  const renderPreview = () => {
    if (!entry.response_body) return <span className="text-slate-300 italic">No response body captured</span>;
    
    const contentType = getContentType().toLowerCase();

    if (contentType.includes('json')) {
      try {
        const parsed = JSON.parse(entry.response_body);
        return (
          <pre className="text-emerald-700 whitespace-pre-wrap leading-relaxed">
            {JSON.stringify(parsed, null, 2)}
          </pre>
        );
      } catch {
        return <pre className="text-slate-600 whitespace-pre-wrap">{entry.response_body}</pre>;
      }
    }

    if (contentType.includes('image/')) {
      return (
        <div className="flex flex-col items-center gap-4 py-4">
          <img 
            src={entry.response_body} 
            className="max-w-full h-auto rounded-lg shadow-md border border-slate-200 bg-[url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAApSURBVHgB7YwxCgAwDMCK/z96p9S6ZAsG6m6ZAnpZAnpZAnpZAnpZAnoZMgX0MnpsmY8AAAAASUVORK5CYII=')] bg-repeat" 
            alt="Response Preview" 
          />
          <div className="text-[10px] text-slate-400 font-mono bg-white px-2 py-1 rounded border border-slate-100">
            {contentType}
          </div>
        </div>
      );
    }

    if (contentType.includes('html')) {
      return (
        <div className="bg-white border border-slate-200 rounded-lg overflow-hidden h-full flex flex-col">
          <div className="bg-slate-50 border-b border-slate-100 px-3 py-1 text-[10px] text-slate-400 font-bold uppercase tracking-wider">HTML Preview</div>
          <iframe 
            srcDoc={entry.response_body} 
            className="w-full flex-1 border-0 min-h-[300px]"
            title="HTML Preview"
            sandbox=""
          />
        </div>
      );
    }

    return <pre className="text-slate-600 whitespace-pre-wrap leading-relaxed">{entry.response_body}</pre>;
  };

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
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">Response Body</h3>
                <div className="flex bg-slate-100 p-0.5 rounded-lg border border-slate-200">
                  <button 
                    onClick={() => setViewMode('preview')}
                    className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${viewMode === 'preview' ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    <Eye size={12} /> Preview
                  </button>
                  <button 
                    onClick={() => setViewMode('raw')}
                    className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${viewMode === 'raw' ? 'bg-white text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
                  >
                    <Code size={12} /> Raw
                  </button>
                </div>
              </div>
              <div className="flex-1 bg-slate-50 rounded-xl p-4 font-mono text-[12px] overflow-auto border border-slate-100 shadow-inner">
                {viewMode === 'preview' ? renderPreview() : (
                  <pre className="text-slate-600 whitespace-pre-wrap leading-relaxed">
                    {entry.response_body || <span className="text-slate-300 italic">No response body captured</span>}
                  </pre>
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
