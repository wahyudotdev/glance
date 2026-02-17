import React, { useState, useEffect } from 'react';
import { X, Play, Plus, Trash2, Activity, AlignLeft } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';

interface ResponseEditorProps {
  isOpen: boolean;
  onClose: () => void;
  entry: TrafficEntry;
  onResume: (status: number, headers: Record<string, string[]>, body: string) => Promise<void>;
  onAbort: (id: string) => Promise<void>;
}

export const ResponseEditor: React.FC<ResponseEditorProps> = ({ 
  isOpen, onClose, entry, onResume, onAbort 
}) => {
  const [status, setStatus] = useState(200);
  const [headers, setHeaders] = useState<{key: string, value: string}[]>([]);
  const [body, setBody] = useState('');
  const [isExecuting, setIsExecuting] = useState(false);

  useEffect(() => {
    if (entry) {
      setStatus(entry.status || 200);
      
      // Auto-prettify body if it's JSON
      let finalBody = entry.response_body || '';
      try {
        const parsed = JSON.parse(finalBody);
        finalBody = JSON.stringify(parsed, null, 2);
      } catch (e) { /* Not JSON */ }
      setBody(finalBody);
      
      const h: {key: string, value: string}[] = [];
      if (entry.response_headers) {
        Object.entries(entry.response_headers).forEach(([k, vs]) => {
          vs.forEach(v => h.push({key: k, value: v}));
        });
      }
      setHeaders(h);
    }
  }, [entry, isOpen]);

  if (!isOpen) return null;

  const handleAddHeader = () => {
    setHeaders([...headers, {key: '', value: ''}]);
  };

  const handleRemoveHeader = (index: number) => {
    setHeaders(headers.filter((_, i) => i !== index));
  };

  const handleHeaderChange = (index: number, field: 'key' | 'value', val: string) => {
    const newHeaders = [...headers];
    newHeaders[index][field] = val;
    setHeaders(newHeaders);
  };

  const prettifyJson = () => {
    try {
      const parsed = JSON.parse(body);
      setBody(JSON.stringify(parsed, null, 2));
    } catch (e) {
      // Not valid JSON, ignore or show hint
    }
  };

  const handleSubmit = async () => {
    setIsExecuting(true);
    try {
      const hMap: Record<string, string[]> = {};
      headers.forEach(h => {
        if (h.key) {
          if (!hMap[h.key]) hMap[h.key] = [];
          hMap[h.key].push(h.value);
        }
      });

      await onResume(status, hMap, body);
      onClose();
    } finally {
      setIsExecuting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-end p-0">
      <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300" onClick={onClose} />
      <div className="relative bg-white dark:bg-slate-900 w-full max-w-2xl h-full shadow-2xl flex flex-col animate-in slide-in-from-right duration-300 transition-colors">
        <div className="h-16 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between px-6 bg-slate-50/50 dark:bg-slate-950/50 transition-colors">
          <h2 className="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <Play size={18} className="text-emerald-500" />
            PAUSED: Intercepted Response
          </h2>
          <button onClick={onClose} className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all">
            <X size={20} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-8 space-y-8">
          <div className="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-xl border border-blue-100 dark:border-blue-800/30 transition-colors">
            <p className="text-xs text-blue-700 dark:text-blue-400 font-medium">Original Request:</p>
            <p className="text-sm font-mono text-blue-900 dark:text-blue-200 mt-1">{entry.method} {entry.url}</p>
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">Response Status</label>
            <input 
              type="number" 
              value={status}
              onChange={(e) => setStatus(parseInt(e.target.value) || 200)}
              className="w-32 px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-emerald-600 dark:text-emerald-400 focus:outline-none transition-colors"
            />
          </div>

          <section>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">Response Headers</h3>
              <button onClick={handleAddHeader} className="flex items-center gap-1.5 px-3 py-1 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg text-[10px] font-bold transition-all hover:bg-blue-100 dark:hover:bg-blue-900/40">
                <Plus size={12} /> Add Header
              </button>
            </div>
            <div className="space-y-2">
              {headers.map((h, i) => (
                <div key={i} className="flex gap-2 group">
                  <input 
                    type="text" placeholder="Key" value={h.key}
                    onChange={(e) => handleHeaderChange(i, 'key', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs font-mono dark:text-slate-300 transition-colors"
                  />
                  <input 
                    type="text" placeholder="Value" value={h.value}
                    onChange={(e) => handleHeaderChange(i, 'value', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs font-mono dark:text-slate-300 transition-colors"
                  />
                  <button onClick={() => handleRemoveHeader(i)} className="p-2 text-slate-300 dark:text-slate-600 hover:text-rose-500 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-all">
                    <Trash2 size={14} />
                  </button>
                </div>
              ))}
            </div>
          </section>

          <section className="flex-1 min-h-[200px] flex flex-col">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-[0.2em]">Response Body</h3>
              <button 
                onClick={prettifyJson}
                className="flex items-center gap-1.5 px-3 py-1 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 rounded-lg text-[10px] font-bold hover:bg-emerald-100 dark:hover:bg-emerald-900/40 transition-all"
                title="Format JSON"
              >
                <AlignLeft size={12} />
                Prettify JSON
              </button>
            </div>
            <textarea 
              value={body}
              onChange={(e) => setBody(e.target.value)}
              className="flex-1 w-full p-4 bg-slate-900 text-emerald-400 rounded-2xl font-mono text-xs min-h-[200px] resize-none focus:outline-none border border-slate-800"
            />
          </section>
        </div>

        <div className="p-6 bg-slate-50/50 dark:bg-slate-950/50 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-3 transition-colors">
          <button 
            onClick={() => onAbort(entry.id)}
            className="px-6 py-2.5 text-sm font-bold text-rose-500 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-xl transition-all"
          >
            Abort / Discard
          </button>
          <button 
            onClick={handleSubmit}
            disabled={isExecuting}
            className="px-8 py-2.5 bg-emerald-500 hover:bg-emerald-600 text-white rounded-xl font-bold text-sm shadow-lg shadow-emerald-200 dark:shadow-none active:scale-95 transition-all flex items-center gap-2"
          >
            {isExecuting ? <Activity className="animate-spin" size={16} /> : <Play size={16} fill="currentColor" />}
            Resume Response
          </button>
        </div>
      </div>
    </div>
  );
};
