import React, { useState } from 'react';
import { Trash2, Plus, Activity, Edit2, Eye, ShieldAlert, AlignLeft } from 'lucide-react';
import type { Rule } from '../../types/traffic';

interface RulesViewProps {
  rules: Rule[];
  onDelete: (id: string) => void;
  onCreate: (rule: Partial<Rule>) => void;
  onEdit: (rule: Rule) => void;
  isLoading: boolean;
}

export const RulesView: React.FC<RulesViewProps> = ({ rules, onDelete, onCreate, onEdit, isLoading }) => {
  const [newPattern, setNewPattern] = useState('');
  const [newMethod, setNewMethod] = useState('ANY');
  const [newType, setNewType] = useState<'breakpoint' | 'mock'>('breakpoint');
  const [newStrategy, setNewStrategy] = useState('both');
  const [newMockStatus, setNewMockStatus] = useState(200);
  const [newMockBody, setNewMockBody] = useState('');

  const prettifyNewMock = () => {
    try {
      const parsed = JSON.parse(newMockBody);
      setNewMockBody(JSON.stringify(parsed, null, 2));
    } catch {
      // Not valid JSON
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newPattern) {
      const rule: Partial<Rule> = {
        type: newType,
        url_pattern: newPattern,
        method: newMethod === 'ANY' ? '' : newMethod,
      };

      if (newType === 'breakpoint') {
        rule.strategy = newStrategy;
      } else {
        rule.response = {
          status: newMockStatus,
          body: newMockBody,
          headers: { 'Content-Type': 'application/json' }
        };
      }

      onCreate(rule);
      setNewPattern('');
      setNewMockBody('');
    }
  };

  return (
    <div className="flex-1 p-12 bg-slate-50 dark:bg-slate-950 overflow-y-auto transition-colors">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-800 dark:text-slate-100">Traffic Rules</h2>
            <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Define patterns to automatically pause or mock traffic.</p>
          </div>
          {isLoading && <Activity className="animate-spin text-blue-600 dark:text-blue-400" size={20} />}
        </div>

        {/* Create Rule Form */}
        <div className="bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm space-y-4">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="flex gap-3">
              <select 
                value={newType}
                onChange={(e) => setNewType(e.target.value as 'breakpoint' | 'mock')}
                className={`w-40 px-4 py-2 border rounded-xl text-sm font-bold ${newType === 'mock' ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-800/30' : 'bg-amber-50 dark:bg-amber-900/20 text-amber-600 dark:text-amber-400 border-amber-200 dark:border-amber-800/30'}`}
              >
                <option value="breakpoint">Action: Pause</option>
                <option value="mock">Action: Mock</option>
              </select>

              <select 
                value={newMethod}
                onChange={(e) => setNewMethod(e.target.value)}
                className="w-32 px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-slate-600 dark:text-slate-300 transition-colors"
              >
                {['ANY', 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => (
                  <option key={m} value={m}>{m}</option>
                ))}
              </select>

              <input 
                type="text" 
                placeholder="URL Pattern (e.g. /api/users)"
                value={newPattern}
                onChange={(e) => setNewPattern(e.target.value)}
                className="flex-1 px-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-mono dark:text-slate-200"
              />
            </div>

            {newType === 'breakpoint' ? (
              <div className="flex items-center gap-3 bg-amber-50/50 dark:bg-amber-900/10 p-3 rounded-xl border border-amber-100 dark:border-amber-900/30 transition-colors">
                <label className="text-[10px] font-black uppercase text-amber-600 dark:text-amber-400 tracking-wider">Pause Strategy:</label>
                <select 
                  value={newStrategy}
                  onChange={(e) => setNewStrategy(e.target.value)}
                  className="bg-transparent text-sm font-bold text-amber-700 dark:text-amber-400 focus:outline-none"
                >
                  <option value="request">Request Only</option>
                  <option value="response">Response Only</option>
                  <option value="both">Both (Request & Response)</option>
                </select>
              </div>
            ) : (
              <div className="space-y-3 bg-emerald-50/50 dark:bg-emerald-900/10 p-4 rounded-xl border border-emerald-100 dark:border-emerald-900/30 animate-in slide-in-from-top-2 duration-200 transition-colors">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <label className="text-[10px] font-black uppercase text-emerald-600 dark:text-emerald-400 tracking-wider">Mock Status:</label>
                    <input 
                      type="number" 
                      value={newMockStatus}
                      onChange={(e) => setNewMockStatus(parseInt(e.target.value) || 200)}
                      className="w-20 bg-white dark:bg-slate-800 border border-emerald-200 dark:border-emerald-800/30 px-2 py-1 rounded-lg text-sm font-bold text-emerald-700 dark:text-emerald-400 transition-colors"
                    />
                  </div>
                  <button 
                    type="button" 
                    onClick={prettifyNewMock}
                    className="flex items-center gap-1.5 px-3 py-1 bg-emerald-100 dark:bg-emerald-900/40 text-emerald-700 dark:text-emerald-400 rounded-lg text-[10px] font-bold hover:bg-emerald-200 dark:hover:bg-emerald-900/60 transition-all active:scale-95"
                  >
                    <AlignLeft size={12} />
                    Prettify JSON
                  </button>
                </div>
                <textarea 
                  placeholder="Mock Response Body (JSON)"
                  value={newMockBody}
                  onChange={(e) => setNewMockBody(e.target.value)}
                  className="w-full h-24 bg-slate-900 text-emerald-400 p-3 rounded-xl font-mono text-xs focus:ring-2 focus:ring-emerald-500/20 border border-slate-800"
                />
              </div>
            )}

            <div className="flex justify-end">
              <button 
                type="submit"
                disabled={!newPattern}
                className="px-8 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-bold shadow-lg shadow-blue-200 dark:shadow-none transition-all active:scale-95 disabled:opacity-50 flex items-center gap-2"
              >
                <Plus size={18} />
                Create Rule
              </button>
            </div>
          </form>
        </div>

        {/* Rules List */}
        <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm overflow-hidden transition-colors">
          <table className="w-full text-left border-separate border-spacing-0">
            <thead className="bg-slate-50/50 dark:bg-slate-950/50">
              <tr className="text-[10px] uppercase tracking-widest text-slate-400 dark:text-slate-500 font-bold border-b border-slate-100 dark:border-slate-800">
                <th className="px-6 py-4">Action</th>
                <th className="px-6 py-4">Method</th>
                <th className="px-6 py-4">URL Pattern</th>
                <th className="px-6 py-4">Configuration</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
                        <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
                          {rules.length > 0 ? (
                            rules.map((rule) => (
                              <tr key={rule.id} className="group hover:bg-slate-50/50 dark:hover:bg-blue-900/10 transition-colors">
                                <td className="px-6 py-4">
                                  <span className={`px-2 py-1 rounded text-[10px] font-bold border flex items-center gap-1.5 w-fit ${rule.type === 'mock' ? 'text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 border-emerald-100 dark:border-emerald-800/30' : 'text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border-amber-100 dark:border-amber-800/30'}`}>
                                    {rule.type === 'mock' ? <Eye size={12} /> : <ShieldAlert size={12} />}
                                    {rule.type === 'mock' ? 'MOCK' : 'PAUSE'}
                                  </span>
                                </td>
                                <td className="px-6 py-4">
                                  <span className={`px-2 py-1 rounded text-[10px] font-bold border ${rule.method ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20 border-blue-100 dark:border-blue-800/30' : 'text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 border-slate-200 dark:border-slate-700'}`}>
                                    {rule.method || 'ANY'}
                                  </span>
                                </td>
                                <td className="px-6 py-4 max-w-xs truncate">
                                  <span className="text-sm font-mono text-slate-600 dark:text-slate-300">{rule.url_pattern}</span>
                                </td>
                                <td className="px-6 py-4">
                                  <span className="text-[10px] text-slate-500 dark:text-slate-400 font-medium">
                                    {rule.type === 'breakpoint' ? `Strategy: ${rule.strategy || 'both'}` : `Returns ${rule.response?.status || 200}`}
                                  </span>
                                </td>
                                <td className="px-6 py-4 text-right">
                                  <div className="flex justify-end gap-1">
                                    <button onClick={() => onEdit(rule)} className="p-2 text-slate-300 dark:text-slate-600 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all"><Edit2 size={14} /></button>
                                    <button onClick={() => onDelete(rule.id)} className="p-2 text-slate-300 dark:text-slate-600 hover:text-rose-500 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-all"><Trash2 size={14} /></button>
                                  </div>
                                </td>
                              </tr>
                            ))
                          ) : (
                            <tr>
                              <td colSpan={5} className="px-6 py-12 text-center text-slate-400 dark:text-slate-600 italic text-sm">No active rules defined.</td>
                            </tr>
                          )}
                        </tbody>
            
          </table>
        </div>
      </div>
    </div>
  );
};
