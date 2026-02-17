import React, { useState } from 'react';
import { FileText, Copy, Check, Eye, Code, Play, X, ShieldAlert, Edit3, ListPlus } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';
import { generateCurl } from '../../lib/curl';

interface DetailsPanelProps {
  entry: TrafficEntry;
  onEdit?: (entry: TrafficEntry) => void;
  onClose?: () => void;
  onBreak?: (entry: TrafficEntry) => void;
  onMock?: (entry: TrafficEntry) => void;
  onAddToScenario?: (entry: TrafficEntry) => void;
}

export const DetailsPanel: React.FC<DetailsPanelProps> = ({ entry, onEdit, onClose, onBreak, onMock, onAddToScenario }) => {
  const [activeTab, setActiveTab] = useState<'headers' | 'body' | 'curl'>('headers');
  const [viewMode, setViewMode] = useState<'preview' | 'raw'>('preview');
  const [copied, setCopied] = useState(false);
  const [copiedRequest, setCopiedRequest] = useState(false);
  const [copiedResponse, setCopiedResponse] = useState(false);

  const isModified = entry.modified_by === 'mock' || entry.modified_by === 'breakpoint';

  const getContentType = () => {
    if (!entry.response_headers) return '';
    const ct = entry.response_headers['Content-Type'] || entry.response_headers['content-type'] || [];
    return ct[0] || '';
  };

  const renderPreview = () => {
    if (!entry.response_body) return <span className="text-slate-300 dark:text-slate-600 italic">No response body captured</span>;
    
    const contentType = getContentType().toLowerCase();

    if (contentType.includes('json')) {
      try {
        const parsed = JSON.parse(entry.response_body);
        return (
          <pre className="text-emerald-700 dark:text-emerald-400 whitespace-pre-wrap leading-relaxed">
            {JSON.stringify(parsed, null, 2)}
          </pre>
        );
      } catch {
        return <pre className="text-slate-600 dark:text-slate-400 whitespace-pre-wrap">{entry.response_body}</pre>;
      }
    }

    if (contentType.includes('image/')) {
      return (
        <div className="flex flex-col items-center gap-4 py-4">
          <img 
            src={entry.response_body} 
            className="max-w-full h-auto rounded-lg shadow-md border border-slate-200 dark:border-slate-700 bg-[url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAApSURBVHgB7YwxCgAwDMCK/z96p9S6ZAsG6m6ZAnpZAnpZAnpZAnpZAnoZMgX0MnpsmY8AAAAASUVORK5CYII=')] bg-repeat" 
            alt="Response Preview" 
          />
          <div className="text-[10px] text-slate-400 dark:text-slate-500 font-mono bg-white dark:bg-slate-800 px-2 py-1 rounded border border-slate-100 dark:border-slate-700 transition-colors">
            {contentType}
          </div>
        </div>
      );
    }

    if (contentType.includes('html')) {
      return (
        <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden h-full flex flex-col">
          <div className="bg-slate-50 dark:bg-slate-900 border-b border-slate-100 dark:border-slate-800 px-3 py-1 text-[10px] text-slate-400 font-bold uppercase tracking-wider">HTML Preview</div>
          <iframe 
            srcDoc={entry.response_body} 
            className="w-full flex-1 border-0 min-h-[300px]"
            title="HTML Preview"
            sandbox=""
          />
        </div>
      );
    }

    return <pre className="text-slate-600 dark:text-slate-400 whitespace-pre-wrap leading-relaxed">{entry.response_body}</pre>;
  };

  const renderRequestBody = () => {
    if (!entry.request_body) return null;
    try {
      const parsed = JSON.parse(entry.request_body);
      return (
        <pre className="text-blue-700 dark:text-blue-400 whitespace-pre-wrap leading-relaxed">
          {JSON.stringify(parsed, null, 2)}
        </pre>
      );
    } catch {
      return <pre className="text-slate-600 dark:text-slate-400 whitespace-pre-wrap leading-relaxed">{entry.request_body}</pre>;
    }
  };

  const handleCopyRequestBody = () => {
    if (entry.request_body) {
      navigator.clipboard.writeText(entry.request_body);
      setCopiedRequest(true);
      setTimeout(() => setCopiedRequest(false), 2000);
    }
  };

  const handleCopyResponseBody = () => {
    if (entry.response_body) {
      navigator.clipboard.writeText(entry.response_body);
      setCopiedResponse(true);
      setTimeout(() => setCopiedResponse(false), 2000);
    }
  };

  const handleCopyCurl = () => {
    const curl = generateCurl(entry);
    navigator.clipboard.writeText(curl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div 
      className="bg-white dark:bg-slate-900 border-l border-slate-200 dark:border-slate-800 flex flex-col shadow-2xl z-20 flex-shrink-0 h-full w-full transition-colors"
    >
      <div className="p-6 border-b border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50 transition-colors">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <h2 className="text-sm font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2 uppercase tracking-tight">
              <FileText size={16} className="text-blue-500" /> Request Details
            </h2>
            {entry.modified_by === 'mock' && (
              <span className="flex items-center gap-1 text-[9px] font-black text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 px-2 py-0.5 rounded-full border border-emerald-100 dark:border-emerald-800/30 uppercase tracking-tighter">
                <Eye size={10} /> Mocked
              </span>
            )}
            {entry.modified_by === 'breakpoint' && (
              <span className="flex items-center gap-1 text-[9px] font-black text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 px-2 py-0.5 rounded-full border border-amber-100 dark:border-amber-800/30 uppercase tracking-tighter">
                <ShieldAlert size={10} /> Paused
              </span>
            )}
            {entry.modified_by === 'editor' && (
              <span className="flex items-center gap-1 text-[9px] font-black text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/20 px-2 py-0.5 rounded-full border border-indigo-100 dark:border-indigo-800/30 uppercase tracking-tighter">
                <Edit3 size={10} /> Editor
              </span>
            )}
            {entry.duration > 0 && (
              <span className="px-2 py-0.5 bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400 rounded-full text-[9px] font-mono border border-slate-200 dark:border-slate-700">
                {(entry.duration / 1000000).toFixed(1)}ms
              </span>
            )}
          </div>
          <div className="flex gap-2">
            <button 
              onClick={() => onEdit?.(entry)}
              className="flex items-center gap-1.5 px-3 py-1.5 bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/30 rounded-lg text-xs font-semibold text-blue-600 dark:text-blue-400 hover:bg-blue-600 dark:hover:bg-blue-500 hover:text-white dark:hover:text-white transition-all shadow-sm active:scale-95"
            >
              <Play size={14} fill="currentColor" />
              Edit & Resend
            </button>
            {!isModified && (
              <>
                <button 
                  onClick={() => onBreak?.(entry)}
                  className="flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 dark:bg-amber-900/20 border border-amber-100 dark:border-amber-800/30 rounded-lg text-xs font-semibold text-amber-600 dark:text-amber-400 hover:bg-amber-600 dark:hover:bg-amber-500 hover:text-white dark:hover:text-white transition-all shadow-sm active:scale-95"
                  title="Pause on the next request matching this URL/Method"
                >
                  <ShieldAlert size={14} />
                  Break on next
                </button>
                <button 
                  onClick={() => onMock?.(entry)}
                  className="flex items-center gap-1.5 px-3 py-1.5 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-800/30 rounded-lg text-xs font-semibold text-emerald-600 dark:text-emerald-400 hover:bg-emerald-600 dark:hover:bg-emerald-500 hover:text-white dark:hover:text-white transition-all shadow-sm active:scale-95"
                  title="Return this response automatically for future requests"
                >
                  <Eye size={14} />
                  Mock this
                </button>
              </>
            )}
            {onAddToScenario && (
              <button 
                onClick={() => onAddToScenario(entry)}
                className="flex items-center gap-1.5 px-3 py-1.5 bg-indigo-50 dark:bg-indigo-900/20 border border-indigo-100 dark:border-indigo-800/30 rounded-lg text-xs font-semibold text-indigo-600 dark:text-indigo-400 hover:bg-indigo-600 dark:hover:bg-indigo-500 hover:text-white dark:hover:text-white transition-all shadow-sm active:scale-95"
                title="Add this request to a new or existing scenario"
              >
                <ListPlus size={14} />
                Add to Scenario
              </button>
            )}
            <button 
              onClick={handleCopyCurl}
              className="flex items-center gap-1.5 px-3 py-1.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs font-semibold text-slate-600 dark:text-slate-300 hover:border-blue-500 dark:hover:border-blue-400 hover:text-blue-600 dark:hover:text-blue-400 transition-all shadow-sm active:scale-95"
            >
              {copied ? <Check size={14} className="text-emerald-500" /> : <Copy size={14} />}
              {copied ? 'Copied!' : 'Copy cURL'}
            </button>
            <button 
              onClick={onClose}
              className="p-1.5 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300 transition-all"
              title="Close Details"
            >
              <X size={18} />
            </button>
          </div>
        </div>
        <div className="font-mono text-[12px] bg-slate-900 text-slate-300 p-3 rounded-lg break-all leading-relaxed shadow-inner border border-slate-800">
          <span className="text-blue-400 font-bold uppercase mr-2">{entry.method}</span>
          {entry.url}
        </div>
      </div>

      <div className="flex px-6 pt-2 border-b border-slate-100 dark:border-slate-800 bg-white dark:bg-slate-900 transition-colors">
        {(['headers', 'body', 'curl'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-3 text-[11px] font-bold uppercase tracking-wider transition-all border-b-2 -mb-[1px] ${
              activeTab === tab 
              ? 'border-blue-600 dark:border-blue-400 text-blue-600 dark:text-blue-400' 
              : 'border-transparent text-slate-400 dark:text-slate-500 hover:text-slate-600 dark:hover:text-slate-300'
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      <div className="flex-1 overflow-y-auto p-6 transition-colors">
        {activeTab === 'headers' && (
          <div className="space-y-6">
            <section>
              <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 mb-3 tracking-[0.2em]">Request Headers</h3>
              <div className="space-y-1.5">
                {Object.entries(entry.request_headers).map(([key, values]) => (
                  <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 dark:border-slate-800/50 last:border-0">
                    <span className="font-bold text-slate-600 dark:text-slate-300 min-w-[120px] shrink-0">{key}</span>
                    <span className="text-slate-500 dark:text-slate-400 break-all font-mono">{values.join(', ')}</span>
                  </div>
                ))}
              </div>
            </section>
            
            {entry.response_headers && (
              <section>
                <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 mb-3 tracking-[0.2em]">Response Headers</h3>
                <div className="space-y-1.5">
                  {Object.entries(entry.response_headers).map(([key, values]) => (
                    <div key={key} className="text-[12px] flex items-start gap-2 py-1 group border-b border-slate-50 dark:border-slate-800/50 last:border-0">
                      <span className="font-bold text-slate-600 dark:text-slate-300 min-w-[120px] shrink-0">{key}</span>
                      <span className="text-slate-500 dark:text-slate-400 break-all font-mono">{values.join(', ')}</span>
                    </div>
                  ))}
                </div>
              </section>
            )}
          </div>
        )}

        {activeTab === 'body' && (
          <div className="h-full flex flex-col gap-6">
            {entry.request_body && (
              <section className="flex-shrink-0 flex flex-col max-h-[40%]">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">Request Body</h3>
                  <button 
                    onClick={handleCopyRequestBody}
                    className="flex items-center gap-1.5 px-2 py-1 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-[10px] font-bold text-slate-600 dark:text-slate-300 hover:border-blue-500 dark:hover:border-blue-400 hover:text-blue-600 dark:hover:text-blue-400 transition-all shadow-sm active:scale-95"
                  >
                    {copiedRequest ? <Check size={12} className="text-emerald-500" /> : <Copy size={12} />}
                    {copiedRequest ? 'Copied!' : 'Copy'}
                  </button>
                </div>
                <div className="flex-1 bg-slate-50 dark:bg-slate-900/50 rounded-xl p-4 font-mono text-[12px] overflow-auto border border-slate-100 dark:border-slate-800 shadow-inner">
                  {renderRequestBody()}
                </div>
              </section>
            )}

            <section className="flex-1 min-h-0 flex flex-col">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">Response Body</h3>
                  <button 
                    onClick={handleCopyResponseBody}
                    className="flex items-center gap-1.5 px-2 py-1 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-[10px] font-bold text-slate-600 dark:text-slate-300 hover:border-blue-500 dark:hover:border-blue-400 hover:text-blue-600 dark:hover:text-blue-400 transition-all shadow-sm active:scale-95"
                  >
                    {copiedResponse ? <Check size={12} className="text-emerald-500" /> : <Copy size={12} />}
                    {copiedResponse ? 'Copied!' : 'Copy'}
                  </button>
                </div>
                <div className="flex bg-slate-100 dark:bg-slate-800 p-0.5 rounded-lg border border-slate-200 dark:border-slate-700 transition-colors">
                  <button 
                    onClick={() => setViewMode('preview')}
                    className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${viewMode === 'preview' ? 'bg-white dark:bg-slate-700 text-blue-600 dark:text-blue-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}`}
                  >
                    <Eye size={12} /> Preview
                  </button>
                  <button 
                    onClick={() => setViewMode('raw')}
                    className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${viewMode === 'raw' ? 'bg-white dark:bg-slate-700 text-blue-600 dark:text-blue-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}`}
                  >
                    <Code size={12} /> Raw
                  </button>
                </div>
              </div>
              <div className="flex-1 bg-slate-50 dark:bg-slate-900/50 rounded-xl p-4 font-mono text-[12px] overflow-auto border border-slate-100 dark:border-slate-800 shadow-inner">
                {viewMode === 'preview' ? renderPreview() : (
                  <pre className="text-slate-600 dark:text-slate-400 whitespace-pre-wrap leading-relaxed">
                    {entry.response_body || <span className="text-slate-300 dark:text-slate-600 italic">No response body captured</span>}
                  </pre>
                )}
              </div>
            </section>
          </div>
        )}

        {activeTab === 'curl' && (
          <div className="space-y-4">
            <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">cURL Command</h3>
            <pre className="bg-slate-900 text-blue-300 p-5 rounded-xl text-[12px] font-mono whitespace-pre-wrap leading-relaxed shadow-xl border border-slate-800">
              {generateCurl(entry)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};
