import React, { useState, useEffect } from 'react';
import { X, Save, Plus, ArrowRight, Trash2, GripVertical, Info, Link as LinkIcon } from 'lucide-react';
import type { Scenario, ScenarioStep, TrafficEntry, VariableMapping } from '../../types/traffic';

interface ScenarioEditorProps {
  isOpen: boolean;
  onClose: () => void;
  scenario: Scenario | null;
  onSave: (scenario: Scenario) => void;
  availableTraffic: TrafficEntry[]; // To show details for existing steps
}

export const ScenarioEditor: React.FC<ScenarioEditorProps> = ({ 
  isOpen, 
  onClose, 
  scenario: initialScenario, 
  onSave,
  availableTraffic 
}) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [steps, setSteps] = useState<ScenarioStep[]>([]);
  const [mappings, setMappings] = useState<VariableMapping[]>([]);
  const [activeTab, setActiveTab] = useState<'steps' | 'mappings'>('steps');

  useEffect(() => {
    if (initialScenario) {
      setName(initialScenario.name || '');
      setDescription(initialScenario.description || '');
      setSteps(initialScenario.steps || []);
      setMappings(initialScenario.variable_mappings || []);
    } else {
      setName('');
      setDescription('');
      setSteps([]);
      setMappings([]);
    }
  }, [initialScenario, isOpen]);

  if (!isOpen) return null;

  const handleSave = () => {
    onSave({
      id: initialScenario?.id || '',
      name,
      description,
      steps,
      variable_mappings: mappings,
      created_at: initialScenario?.created_at || new Date().toISOString(),
    });
  };

  const removeStep = (index: number) => {
    setSteps(prev => prev.filter((_, i) => i !== index).map((s, i) => ({ ...s, order: i + 1 })));
  };

  const updateStepNote = (index: number, notes: string) => {
    setSteps(prev => prev.map((s, i) => i === index ? { ...s, notes } : s));
  };

  const addMapping = () => {
    setMappings(prev => [...prev, { name: '', source_entry_id: '', source_path: '', target_json_path: '' }]);
  };

  const updateMapping = (index: number, field: keyof VariableMapping, value: string) => {
    setMappings(prev => prev.map((m, i) => i === index ? { ...m, [field]: value } : m));
  };

  const removeMapping = (index: number) => {
    setMappings(prev => prev.filter((_, i) => i !== index));
  };

  const getEntry = (id: string) => availableTraffic.find(e => e.id === id);

  return (
    <div className="fixed inset-0 z-[100] flex flex-col bg-slate-50 dark:bg-slate-950 transition-colors animate-in fade-in duration-200">
      <header className="h-16 flex items-center justify-between px-8 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 shrink-0">
        <div className="flex items-center gap-4">
          <button 
            onClick={onClose}
            className="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-all"
          >
            <X size={20} />
          </button>
          <div className="flex items-center gap-3">
            <h2 className="text-lg font-black text-slate-800 dark:text-slate-100 tracking-tight">
              {initialScenario ? 'Edit Scenario' : 'New Scenario'}
            </h2>
            {name && <span className="text-slate-300 dark:text-slate-700">|</span>}
            <span className="text-sm font-medium text-slate-500 truncate max-w-xs">{name}</span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button 
            onClick={onClose}
            className="px-4 py-2 text-sm font-bold text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 transition-colors"
          >
            Cancel
          </button>
          <button 
            onClick={handleSave}
            disabled={!name || steps.length === 0}
            className="flex items-center gap-2 px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-200 dark:disabled:bg-slate-800 disabled:text-slate-400 dark:disabled:text-slate-600 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-200 dark:shadow-none transition-all"
          >
            <Save size={18} /> Save Scenario
          </button>
        </div>
      </header>

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar: Metadata */}
        <div className="w-80 border-r border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-8 flex flex-col gap-6 shrink-0 overflow-y-auto">
          <div>
            <label className="block text-[11px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-2">Scenario Name</label>
            <input 
              type="text" 
              placeholder="e.g., User Signup Flow"
              className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all dark:text-slate-100"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          <div>
            <label className="block text-[11px] font-black uppercase tracking-widest text-slate-400 dark:text-slate-500 mb-2">Description</label>
            <textarea 
              rows={4}
              placeholder="Describe the workflow being tested..."
              className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all resize-none dark:text-slate-100"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          <div className="mt-4 p-4 bg-blue-50 dark:bg-blue-900/10 rounded-2xl border border-blue-100 dark:border-blue-800/30">
            <div className="flex gap-3 text-blue-600 dark:text-blue-400">
              <Info size={20} className="shrink-0" />
              <div className="text-xs leading-relaxed font-medium">
                <p className="font-bold mb-1">AI Tip</p>
                Include accurate names and notes. This metadata helps AI agents generate better test scripts.
              </div>
            </div>
          </div>
        </div>

        {/* Main Content: Steps & Mappings */}
        <div className="flex-1 flex flex-col min-w-0">
          <div className="flex border-b border-slate-200 dark:border-slate-800 px-8 bg-white dark:bg-slate-900">
            <button 
              onClick={() => setActiveTab('steps')}
              className={`px-6 py-4 text-xs font-black uppercase tracking-widest border-b-2 transition-all ${
                activeTab === 'steps' ? 'border-blue-600 text-blue-600' : 'border-transparent text-slate-400 hover:text-slate-600'
              }`}
            >
              Sequence Steps ({steps.length})
            </button>
            <button 
              onClick={() => setActiveTab('mappings')}
              className={`px-6 py-4 text-xs font-black uppercase tracking-widest border-b-2 transition-all ${
                activeTab === 'mappings' ? 'border-blue-600 text-blue-600' : 'border-transparent text-slate-400 hover:text-slate-600'
              }`}
            >
              Variable Mappings ({mappings.length})
            </button>
          </div>

          <div className="flex-1 overflow-y-auto p-8">
            {activeTab === 'steps' ? (
              <div className="max-w-4xl mx-auto flex flex-col gap-4">
                {steps.map((step, index) => {
                  const entry = getEntry(step.traffic_entry_id);
                  return (
                    <div 
                      key={step.id || index}
                      className="group flex gap-4 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-blue-200 dark:hover:border-blue-800 transition-all"
                    >
                      <div className="flex flex-col items-center pt-1">
                        <div className="w-6 h-6 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-[10px] font-black text-slate-500">
                          {index + 1}
                        </div>
                        <div className="w-px flex-1 bg-slate-100 dark:bg-slate-800 my-2" />
                        <GripVertical size={16} className="text-slate-300 dark:text-slate-700 cursor-grab active:cursor-grabbing" />
                      </div>

                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center gap-3">
                            <span className="font-mono font-bold text-xs text-blue-600 dark:text-blue-400">{entry?.method || 'GET'}</span>
                            <span className="text-sm font-medium text-slate-700 dark:text-slate-200 truncate max-w-md font-mono">{entry?.url || 'Unknown URL'}</span>
                            {entry?.status && (
                              <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-500">
                                {entry.status}
                              </span>
                            )}
                          </div>
                          <button 
                            onClick={() => removeStep(index)}
                            className="p-1.5 text-slate-400 hover:text-rose-500 rounded-lg opacity-0 group-hover:opacity-100 transition-all"
                          >
                            <Trash2 size={16} />
                          </button>
                        </div>
                        
                        <input 
                          type="text" 
                          placeholder="Add a note for this step (e.g., 'Expect login to succeed')"
                          className="w-full px-3 py-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg text-xs italic focus:outline-none focus:ring-1 focus:ring-blue-500/30 transition-all dark:text-slate-300"
                          value={step.notes || ''}
                          onChange={(e) => updateStepNote(index, e.target.value)}
                        />
                      </div>
                    </div>
                  );
                })}
                
                {steps.length === 0 && (
                  <div className="text-center py-12 bg-white dark:bg-slate-900 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-800 text-slate-400 font-medium">
                    No steps in this scenario. Add them from the Traffic Inspector.
                  </div>
                )}
              </div>
            ) : (
              <div className="max-w-5xl mx-auto">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-sm font-bold text-slate-700 dark:text-slate-200 uppercase tracking-wider">Dynamic Parameters</h3>
                  <button 
                    onClick={addMapping}
                    className="flex items-center gap-2 px-4 py-2 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 rounded-xl font-bold text-xs hover:bg-slate-200 dark:hover:bg-slate-700 transition-all"
                  >
                    <Plus size={14} /> Add Mapping
                  </button>
                </div>

                <div className="flex flex-col gap-3">
                  {mappings.map((mapping, index) => (
                    <div key={index} className="flex items-center gap-3 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800 animate-in slide-in-from-top-2 duration-200">
                      <div className="flex-1">
                        <label className="block text-[9px] font-black uppercase text-slate-400 mb-1">Var Name</label>
                        <input 
                          type="text" 
                          placeholder="sessionToken"
                          className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg text-xs font-mono dark:text-slate-200"
                          value={mapping.name}
                          onChange={(e) => updateMapping(index, 'name', e.target.value)}
                        />
                      </div>

                      <div className="shrink-0 pt-4 text-slate-300 dark:text-slate-700">
                        <LinkIcon size={16} />
                      </div>

                      <div className="flex-[2]">
                        <label className="block text-[9px] font-black uppercase text-slate-400 mb-1">Source Path (from previous response)</label>
                        <input 
                          type="text" 
                          placeholder="body.data.token"
                          className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg text-xs font-mono dark:text-slate-200"
                          value={mapping.source_path}
                          onChange={(e) => updateMapping(index, 'source_path', e.target.value)}
                        />
                      </div>

                      <div className="shrink-0 pt-4 text-slate-300 dark:text-slate-700">
                        <ArrowRight size={16} />
                      </div>

                      <div className="flex-[2]">
                        <label className="block text-[9px] font-black uppercase text-slate-400 mb-1">Target JSON Path (in next request)</label>
                        <input 
                          type="text" 
                          placeholder="header.Authorization"
                          className="w-full px-3 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg text-xs font-mono dark:text-slate-200"
                          value={mapping.target_json_path}
                          onChange={(e) => updateMapping(index, 'target_json_path', e.target.value)}
                        />
                      </div>

                      <button 
                        onClick={() => removeMapping(index)}
                        className="p-2 text-slate-400 hover:text-rose-500 rounded-lg transition-all pt-4"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  ))}
                  
                  {mappings.length === 0 && (
                    <div className="text-center py-12 bg-white dark:bg-slate-900 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-800 text-slate-400 font-medium">
                      No variable mappings defined. Mappings help the AI understand how data flows between steps.
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
