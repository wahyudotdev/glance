import React, { useState, useEffect } from 'react';
import { X, Play, Plus, Trash2, Activity } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';

interface RequestEditorProps {
  isOpen: boolean;
  onClose: () => void;
  initialRequest?: TrafficEntry | null;
  onExecute: (request: Partial<TrafficEntry>) => Promise<void>;
  isIntercept?: boolean;
  onAbort?: (id: string) => Promise<void>;
}

export const RequestEditor: React.FC<RequestEditorProps> = ({ 
  isOpen, onClose, initialRequest, onExecute, isIntercept, onAbort
}) => {
  const [method, setMethod] = useState('GET');
  const [url, setUrl] = useState('');
  const [headers, setHeaders] = useState<{key: string, value: string}[]>([]);
  const [body, setBody] = useState('');
  const [isExecuting, setIsExecuting] = useState(false);

  useEffect(() => {
    if (initialRequest) {
      setMethod(initialRequest.method);
      setUrl(initialRequest.url);
      setBody(initialRequest.request_body || '');
      
      const h: {key: string, value: string}[] = [];
      Object.entries(initialRequest.request_headers).forEach(([k, vs]) => {
        vs.forEach(v => h.push({key: k, value: v}));
      });
      setHeaders(h);
    } else {
      setMethod('GET');
      setUrl('');
      setHeaders([{key: 'Content-Type', value: 'application/json'}]);
      setBody('');
    }
  }, [initialRequest, isOpen]);

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
      // Reconstruct header object
      const hMap: Record<string, string[]> = {};
      headers.forEach(h => {
        if (h.key) {
          if (!hMap[h.key]) hMap[h.key] = [];
          hMap[h.key].push(h.value);
        }
      });

      await onExecute({
        method,
        url,
        request_headers: hMap,
        request_body: body
      });
      onClose();
    } finally {
      setIsExecuting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-end p-0">
      <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300" onClick={onClose} />
      <div className="relative bg-white w-full max-w-2xl h-full shadow-2xl flex flex-col animate-in slide-in-from-right duration-300">
        {/* Header */}
        <div className="h-16 border-b border-slate-100 flex items-center justify-between px-6 bg-slate-50/50">
          <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
            <Play size={18} className={isIntercept ? "text-amber-500" : "text-blue-600"} />
            {isIntercept ? 'PAUSED: Intercepted Request' : (initialRequest ? 'Edit & Resend Request' : 'New Request')}
          </h2>
          <button onClick={onClose} className="p-2 hover:bg-white rounded-lg text-slate-400 transition-all">
            <X size={20} />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-8 space-y-8">
          {/* Method & URL */}
          <div className="flex gap-3">
            <select 
              value={method} 
              onChange={(e) => setMethod(e.target.value)}
              className="w-32 px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold text-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              {['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'].map(m => (
                <option key={m} value={m}>{m}</option>
              ))}
            </select>
            <input 
              type="text" 
              placeholder="https://api.example.com/endpoint"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className="flex-1 px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            />
          </div>

          {/* Headers */}
          <section>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-[10px] font-black uppercase text-slate-400 tracking-[0.2em]">Headers</h3>
              <button 
                onClick={handleAddHeader}
                className="flex items-center gap-1.5 px-3 py-1 bg-blue-50 text-blue-600 rounded-lg text-[10px] font-bold hover:bg-blue-100 transition-all"
              >
                <Plus size={12} /> Add Header
              </button>
            </div>
            <div className="space-y-2">
              {headers.map((h, i) => (
                <div key={i} className="flex gap-2 group">
                  <input 
                    type="text" 
                    placeholder="Key"
                    value={h.key}
                    onChange={(e) => handleHeaderChange(i, 'key', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-mono"
                  />
                  <input 
                    type="text" 
                    placeholder="Value"
                    value={h.value}
                    onChange={(e) => handleHeaderChange(i, 'value', e.target.value)}
                    className="flex-1 px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-mono"
                  />
                  <button 
                    onClick={() => handleRemoveHeader(i)}
                    className="p-2 text-slate-300 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                  >
                    <Trash2 size={14} />
                  </button>
                </div>
              ))}
            </div>
          </section>

          {/* Body */}
          {(method !== 'GET' && method !== 'HEAD') && (
            <section className="flex-1 min-h-[200px] flex flex-col">
              <h3 className="text-[10px] font-black uppercase text-slate-400 mb-4 tracking-[0.2em]">Request Body</h3>
              <textarea 
                value={body}
                onChange={(e) => setBody(e.target.value)}
                placeholder='{"example": "data"}'
                className="flex-1 w-full p-4 bg-slate-900 text-emerald-400 rounded-2xl font-mono text-xs focus:outline-none focus:ring-4 focus:ring-blue-500/10 min-h-[200px] resize-none"
              />
            </section>
          )}
        </div>

        {/* Footer */}
        <div className="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end gap-3">
          {isIntercept ? (
            <button 
              onClick={() => initialRequest && onAbort?.(initialRequest.id)}
              className="px-6 py-2.5 text-sm font-bold text-rose-500 hover:bg-rose-50 rounded-xl transition-all"
            >
              Abort / Discard
            </button>
          ) : (
            <button 
              onClick={onClose}
              className="px-6 py-2.5 text-sm font-bold text-slate-500 hover:bg-white rounded-xl transition-all"
            >
              Cancel
            </button>
          )}
          <button 
            onClick={handleSubmit}
            disabled={!url || isExecuting}
            className={`px-8 py-2.5 ${isIntercept ? "bg-amber-500 hover:bg-amber-600 shadow-amber-200" : "bg-blue-600 hover:bg-blue-700 shadow-blue-200"} text-white rounded-xl font-bold text-sm shadow-lg active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2`}
          >
            {isExecuting ? (
              <Activity className="animate-spin" size={16} />
            ) : (
              <Play size={16} fill="currentColor" />
            )}
            {isIntercept ? 'Resume Request' : 'Execute Request'}
          </button>
        </div>
      </div>
    </div>
  );
};
