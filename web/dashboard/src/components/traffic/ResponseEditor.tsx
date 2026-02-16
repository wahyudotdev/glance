import React, { useState, useEffect } from 'react';
import { X, Play, Plus, Trash2, Activity } from 'lucide-react';
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
      setBody(entry.response_body || '');
      
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
      <div className="relative bg-white w-full max-w-2xl h-full shadow-2xl flex flex-col animate-in slide-in-from-right duration-300">
        <div className="h-16 border-b border-slate-100 flex items-center justify-between px-6 bg-slate-50/50">
          <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
            <Play size={18} className="text-emerald-500" />
            PAUSED: Intercepted Response
          </h2>
          <button onClick={onClose} className="p-2 hover:bg-white rounded-lg text-slate-400 transition-all">
            <X size={20} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-8 space-y-8">
          <div className="bg-blue-50 p-4 rounded-xl border border-blue-100">
            <p className="text-xs text-blue-700 font-medium">Original Request:</p>
            <p className="text-sm font-mono text-blue-900 mt-1">{entry.method} {entry.url}</p>
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">Response Status</label>
            <input 
              type="number" 
              value={status}
              onChange={(e) => setStatus(parseInt(e.target.value) || 200)}
              className="w-32 px-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold text-emerald-600 focus:outline-none"
            />
          </div>

          <section>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">Response Headers</h3>
              <button onClick={handleAddHeader} className="flex items-center gap-1.5 px-3 py-1 bg-blue-50 text-blue-600 rounded-lg text-[10px] font-bold">
                <Plus size={12} /> Add Header
              </button>
            </div>
            <div className="space-y-2">
              {headers.map((h, i) => (
                <div key={i} className="flex gap-2 group">
                  <input 
                    type="text" placeholder="Key" value={h.key}
                    onChange={(e) => handleHeaderChange(i, 'key', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-mono"
                  />
                  <input 
                    type="text" placeholder="Value" value={h.value}
                    onChange={(e) => handleHeaderChange(i, 'value', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-mono"
                  />
                  <button onClick={() => handleRemoveHeader(i)} className="p-2 text-slate-300 hover:text-rose-500">
                    <Trash2 size={14} />
                  </button>
                </div>
              ))}
            </div>
          </section>

          <section className="flex-1 min-h-[200px] flex flex-col">
            <h3 className="text-[10px] font-black uppercase text-slate-400 mb-4 tracking-[0.2em]">Response Body</h3>
            <textarea 
              value={body}
              onChange={(e) => setBody(e.target.value)}
              className="flex-1 w-full p-4 bg-slate-900 text-emerald-400 rounded-2xl font-mono text-xs min-h-[200px] resize-none"
            />
          </section>
        </div>

        <div className="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end gap-3">
          <button 
            onClick={() => onAbort(entry.id)}
            className="px-6 py-2.5 text-sm font-bold text-rose-500 hover:bg-rose-50 rounded-xl transition-all"
          >
            Abort / Discard
          </button>
          <button 
            onClick={handleSubmit}
            disabled={isExecuting}
            className="px-8 py-2.5 bg-emerald-500 hover:bg-emerald-600 text-white rounded-xl font-bold text-sm shadow-lg shadow-emerald-200 active:scale-95 transition-all flex items-center gap-2"
          >
            {isExecuting ? <Activity className="animate-spin" size={16} /> : <Play size={16} fill="currentColor" />}
            Resume Response
          </button>
        </div>
      </div>
    </div>
  );
};
