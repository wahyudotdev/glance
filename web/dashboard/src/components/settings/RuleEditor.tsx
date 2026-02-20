import React, { useState, useEffect } from 'react';
import { X, Save, ShieldAlert, Eye, Edit2, Maximize2, Minimize2 } from 'lucide-react';
import type { Rule } from './RulesView';
import { JSONTreeEditor } from '../ui/JSONTreeEditor';

interface RuleEditorProps {
  isOpen: boolean;
  onClose: () => void;
  rule: Rule | null;
  onSave: (id: string, rule: Partial<Rule>) => void;
}

export const RuleEditor: React.FC<RuleEditorProps> = ({ 
  isOpen, onClose, rule, onSave 
}) => {
  const [type, setType] = useState<'breakpoint' | 'mock'>('breakpoint');
  const [pattern, setPattern] = useState('');
  const [method, setMethod] = useState('');
  const [strategy, setStrategy] = useState('both');
  const [mockStatus, setMockStatus] = useState(200);
  const [mockBody, setMockBody] = useState('');
  const [isFullScreen, setIsFullScreen] = useState(false);

  useEffect(() => {
    if (rule) {
      setType(rule.type);
      setPattern(rule.url_pattern);
      setMethod(rule.method || 'ANY');
      setStrategy(rule.strategy || 'both');
      setMockStatus(rule.response?.status || 200);
      
      let body = rule.response?.body || '';
      try {
        const parsed = JSON.parse(body);
        body = JSON.stringify(parsed, null, 2);
      } catch (e) {}
      setMockBody(body);
    }
  }, [rule, isOpen]);

  if (!isOpen || !rule) return null;

  const handleSave = () => {
    const updated: Partial<Rule> = {
      type,
      url_pattern: pattern,
      method: method === 'ANY' ? '' : method,
    };

    if (type === 'breakpoint') {
      updated.strategy = strategy;
    } else {
      updated.response = {
        status: mockStatus,
        body: mockBody,
        headers: { 'Content-Type': 'application/json' }
      };
    }

    onSave(rule.id, updated);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-[120] flex items-center justify-end p-0">
      <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-300" onClick={onClose} />
      <div className={`relative bg-white dark:bg-slate-900 shadow-2xl flex flex-col transition-all duration-300 ${isFullScreen ? 'w-full h-full' : 'w-full max-w-2xl h-full animate-in slide-in-from-right'}`}>
        {/* Header */}
        <div className="h-16 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between px-6 bg-slate-50/50 dark:bg-slate-950/50 transition-colors shrink-0">
          <div className="flex items-center gap-4">
            <h2 className="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
              <Edit2 size={18} className="text-blue-600 dark:text-blue-400" />
              {isFullScreen ? 'Full Screen JSON Editor' : 'Edit Rule'}
            </h2>
          </div>
          <div className="flex items-center gap-2">
            {type === 'mock' && (
              <button 
                onClick={() => setIsFullScreen(!isFullScreen)}
                className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all"
                title={isFullScreen ? "Exit Full Screen" : "Edit in Full Screen"}
              >
                {isFullScreen ? <Minimize2 size={20} /> : <Maximize2 size={20} />}
              </button>
            )}
            <button onClick={onClose} className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all">
              <X size={20} />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className={`flex-1 overflow-y-auto transition-all ${isFullScreen ? 'p-0 flex flex-col bg-slate-900' : 'p-8 space-y-8'}`}>
          {!isFullScreen && (
            <>
              {/* Action Toggle */}
              <div className="flex p-1 bg-slate-100 dark:bg-slate-800 rounded-xl w-fit transition-colors">
                <button 
                  onClick={() => setType('breakpoint')}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-bold transition-all ${type === 'breakpoint' ? 'bg-white dark:bg-slate-700 text-amber-600 dark:text-amber-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}`}
                >
                  <ShieldAlert size={16} /> Pause
                </button>
                <button 
                  onClick={() => setType('mock')}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-bold transition-all ${type === 'mock' ? 'bg-white dark:bg-slate-700 text-emerald-600 dark:text-emerald-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}`}
                >
                  <Eye size={16} /> Mock
                </button>
              </div>

              <div className="grid grid-cols-1 gap-6">
                <div className="flex gap-3">
                  <div className="w-32 space-y-1.5">
                    <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-wider">Method</label>
                    <select 
                      value={method} 
                      onChange={(e) => setMethod(e.target.value)}
                      className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold dark:text-slate-200 transition-colors"
                    >
                      {['ANY', 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => (
                        <option key={m} value={m}>{m}</option>
                      ))}
                    </select>
                  </div>
                  <div className="flex-1 space-y-1.5">
                    <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-wider">URL Pattern</label>
                    <input 
                      type="text" 
                      value={pattern}
                      onChange={(e) => setPattern(e.target.value)}
                      className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-mono dark:text-slate-200 transition-colors"
                    />
                  </div>
                </div>

                {type === 'breakpoint' ? (
                  <div className="space-y-3">
                    <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-wider">Interception Strategy</label>
                    <div className="grid grid-cols-3 gap-3">
                      {['request', 'response', 'both'].map((s) => (
                        <button
                          key={s}
                          onClick={() => setStrategy(s)}
                          className={`py-3 px-4 rounded-xl border-2 text-xs font-bold capitalize transition-all ${strategy === s ? 'border-amber-500 bg-amber-50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-400' : 'border-slate-100 dark:border-slate-800 bg-white dark:bg-slate-900 text-slate-400 dark:text-slate-500 hover:border-slate-200 dark:hover:border-slate-700'}`}
                        >
                          {s}
                        </button>
                      ))}
                    </div>
                  </div>
                ) : (
                  <div className="space-y-6 animate-in slide-in-from-top-2 duration-300">
                    <div className="flex items-center justify-between">
                      <div className="space-y-1.5">
                        <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-wider">Response Status</label>
                        <input 
                          type="number" 
                          value={mockStatus}
                          onChange={(e) => setMockStatus(parseInt(e.target.value) || 200)}
                          className="w-24 px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-emerald-600 dark:text-emerald-400 transition-colors"
                        />
                      </div>
                    </div>
                    <div className="space-y-1.5 flex flex-col h-[400px]">
                      <label className="text-[10px] font-black uppercase text-slate-400 dark:text-slate-500 tracking-wider">Response Body</label>
                      <JSONTreeEditor 
                        value={mockBody}
                        onChange={setMockBody}
                        className="flex-1 min-h-0 shadow-inner"
                      />
                    </div>
                  </div>
                )}
              </div>
            </>
          )}

          {isFullScreen && (
            <JSONTreeEditor 
              value={mockBody}
              onChange={setMockBody}
              className="flex-1"
              isFullScreen={true}
            />
          )}
        </div>

        {/* Footer */}
        <div className="p-6 bg-slate-50/50 dark:bg-slate-950/50 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-3 transition-colors shrink-0">
          <button 
            onClick={onClose}
            className="px-6 py-2.5 text-sm font-bold text-slate-500 dark:text-slate-400 hover:bg-white dark:hover:bg-slate-800 rounded-xl transition-all"
          >
            Cancel
          </button>
          <button 
            onClick={handleSave}
            className="px-8 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-200 dark:shadow-none active:scale-95 transition-all flex items-center gap-2"
          >
            <Save size={16} />
            Save Changes
          </button>
        </div>
      </div>
    </div>
  );
};
