import React, { useState } from 'react';
import { Trash2, Plus, Globe, Activity, Edit2, Save, X } from 'lucide-react';

export interface Rule {
  id: string;
  type: string;
  url_pattern: string;
  method: string;
  strategy: string;
}

interface RulesViewProps {
  rules: Rule[];
  onDelete: (id: string) => void;
  onCreate: (pattern: string, method: string, strategy: string) => void;
  onUpdate: (id: string, pattern: string, method: string, strategy: string) => void;
  isLoading: boolean;
}

export const RulesView: React.FC<RulesViewProps> = ({ rules, onDelete, onCreate, onUpdate, isLoading }) => {
  const [newPattern, setNewPattern] = useState('');
  const [newMethod, setNewMethod] = useState('ANY');
  const [newStrategy, setNewMethodStrategy] = useState('request');

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editPattern, setEditPattern] = useState('');
  const [editMethod, setEditMethod] = useState('');
  const [editStrategy, setEditStrategy] = useState('');

  const startEdit = (rule: Rule) => {
    setEditingId(rule.id);
    setEditPattern(rule.url_pattern);
    setEditMethod(rule.method || 'ANY');
    setEditStrategy(rule.strategy || 'request');
  };

  const cancelEdit = () => {
    setEditingId(null);
  };

  const saveEdit = () => {
    if (editingId) {
      onUpdate(editingId, editPattern, editMethod === 'ANY' ? '' : editMethod, editStrategy);
      setEditingId(null);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newPattern) {
      onCreate(newPattern, newMethod === 'ANY' ? '' : newMethod, newStrategy);
      setNewPattern('');
    }
  };

  return (
    <div className="flex-1 p-12 bg-slate-50 overflow-y-auto">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-slate-800">Breakpoint Rules</h2>
            <p className="text-sm text-slate-500 mt-1">Manage filters that automatically pause requests for editing.</p>
          </div>
          {isLoading && <Activity className="animate-spin text-blue-600" size={20} />}
        </div>

        {/* Create Rule Form */}
        <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
          <form onSubmit={handleSubmit} className="flex gap-3">
            <select 
              value={newMethod}
              onChange={(e) => setNewMethod(e.target.value)}
              className="w-32 px-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold text-slate-600"
            >
              {['ANY', 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => (
                <option key={m} value={m}>{m}</option>
              ))}
            </select>
            <select 
              value={newStrategy}
              onChange={(e) => setNewMethodStrategy(e.target.value)}
              className="w-32 px-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold text-amber-600"
            >
              <option value="request">Request</option>
              <option value="response">Response</option>
              <option value="both">Both</option>
            </select>
            <input 
              type="text" 
              placeholder="URL Pattern (e.g. /api/users or google.com)"
              value={newPattern}
              onChange={(e) => setNewPattern(e.target.value)}
              className="flex-1 px-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm font-mono"
            />
            <button 
              type="submit"
              disabled={!newPattern}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-sm font-bold shadow-lg shadow-blue-200 transition-all active:scale-95 disabled:opacity-50 flex items-center gap-2"
            >
              <Plus size={16} />
              Add Rule
            </button>
          </form>
        </div>

        {/* Rules List */}
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
          <table className="w-full text-left border-separate border-spacing-0">
            <thead className="bg-slate-50/50">
              <tr className="text-[10px] uppercase tracking-widest text-slate-400 font-bold border-b border-slate-100">
                <th className="px-6 py-4">Method</th>
                <th className="px-6 py-4">Pause</th>
                <th className="px-6 py-4">URL Pattern</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {rules.length > 0 ? rules.filter(r => r.type === 'breakpoint').map((rule) => (
                <tr key={rule.id} className="group hover:bg-slate-50/50 transition-colors">
                  {editingId === rule.id ? (
                    <>
                      <td className="px-6 py-3">
                        <select 
                          value={editMethod}
                          onChange={(e) => setEditMethod(e.target.value)}
                          className="w-full px-2 py-1 bg-white border border-slate-200 rounded text-[10px] font-bold"
                        >
                          {['ANY', 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => (
                            <option key={m} value={m}>{m}</option>
                          ))}
                        </select>
                      </td>
                      <td className="px-6 py-3">
                        <select 
                          value={editStrategy}
                          onChange={(e) => setEditStrategy(e.target.value)}
                          className="w-full px-2 py-1 bg-white border border-slate-200 rounded text-[10px] font-bold text-amber-600"
                        >
                          <option value="request">Request</option>
                          <option value="response">Response</option>
                          <option value="both">Both</option>
                        </select>
                      </td>
                      <td className="px-6 py-3">
                        <input 
                          type="text"
                          value={editPattern}
                          onChange={(e) => setEditPattern(e.target.value)}
                          className="w-full px-2 py-1 bg-white border border-slate-200 rounded text-[11px] font-mono"
                        />
                      </td>
                      <td className="px-6 py-3 text-right">
                        <div className="flex justify-end gap-1">
                          <button onClick={saveEdit} className="p-1.5 text-emerald-600 hover:bg-emerald-50 rounded transition-all">
                            <Save size={14} />
                          </button>
                          <button onClick={cancelEdit} className="p-1.5 text-slate-400 hover:bg-slate-100 rounded transition-all">
                            <X size={14} />
                          </button>
                        </div>
                      </td>
                    </>
                  ) : (
                    <>
                      <td className="px-6 py-4">
                        <span className={`px-2 py-1 rounded text-[10px] font-bold border ${rule.method ? 'text-blue-600 bg-blue-50 border-blue-100' : 'text-slate-500 bg-slate-100 border-slate-200'}`}>
                          {rule.method || 'ANY'}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <span className="px-2 py-1 rounded text-[10px] font-bold border text-amber-600 bg-amber-50 border-amber-100 uppercase">
                          {rule.strategy || 'request'}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <Globe size={14} className="text-slate-300" />
                          <span className="text-sm font-mono text-slate-600">{rule.url_pattern}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex justify-end gap-1">
                          <button 
                            onClick={() => startEdit(rule)}
                            className="p-2 text-slate-300 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                          >
                            <Edit2 size={14} />
                          </button>
                          <button 
                            onClick={() => onDelete(rule.id)}
                            className="p-2 text-slate-300 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all"
                          >
                            <Trash2 size={14} />
                          </button>
                        </div>
                      </td>
                    </>
                  )}
                </tr>
              )) : (
                <tr>
                  <td colSpan={3} className="px-6 py-12 text-center text-slate-400 italic text-sm">
                    No active breakpoints. New requests will pass through normally.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
