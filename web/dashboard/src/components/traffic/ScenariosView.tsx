import React from 'react';
import { Plus, Trash2, Play, Calendar, ListOrdered } from 'lucide-react';
import type { Scenario } from '../../types/traffic';

interface ScenariosViewProps {
  scenarios: Scenario[];
  isLoading: boolean;
  onSelect: (scenario: Scenario) => void;
  onDelete: (id: string) => void;
  onCreateNew: () => void;
}

export const ScenariosView: React.FC<ScenariosViewProps> = ({ 
  scenarios, 
  isLoading, 
  onSelect, 
  onDelete, 
  onCreateNew 
}) => {
  const handleDelete = (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    onDelete(id);
  };

  return (
    <div className="flex-1 flex flex-col min-w-0 bg-slate-50 dark:bg-slate-950 overflow-y-auto p-8 transition-colors">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-black text-slate-800 dark:text-slate-100 tracking-tight">Traffic Scenarios</h1>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1 font-medium">Record and organize request sequences for AI-assisted test generation.</p>
        </div>
        
        <button 
          onClick={onCreateNew}
          className="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold text-sm shadow-lg shadow-blue-200 dark:shadow-none transition-all"
        >
          <Plus size={18} /> New Scenario
        </button>
      </div>

      {isLoading ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      ) : scenarios.length === 0 ? (
        <div className="flex-1 flex flex-col items-center justify-center text-center p-12 bg-white dark:bg-slate-900 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-800 transition-colors">
          <div className="w-16 h-16 bg-slate-50 dark:bg-slate-800 rounded-2xl flex items-center justify-center mb-4 text-slate-400">
            <ListOrdered size={32} />
          </div>
          <h3 className="text-lg font-bold text-slate-700 dark:text-slate-200">No scenarios yet</h3>
          <p className="text-sm text-slate-400 dark:text-slate-500 max-w-xs mt-2 font-medium">
            Start recording traffic or select multiple requests from the inspector to create your first scenario.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {scenarios.map((scenario) => (
            <div 
              key={scenario.id}
              onClick={() => onSelect(scenario)}
              className="group bg-white dark:bg-slate-900 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-blue-400 dark:hover:border-blue-600 hover:shadow-xl hover:shadow-blue-500/5 cursor-pointer transition-all duration-300 relative overflow-hidden"
            >
              <div className="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
                <button 
                  onClick={(e) => handleDelete(e, scenario.id)}
                  className="p-2 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-all"
                >
                  <Trash2 size={16} />
                </button>
              </div>

              <div className="flex flex-col h-full">
                <div className="flex items-center gap-3 mb-4">
                  <div className="p-2 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded-lg">
                    <Play size={20} />
                  </div>
                  <h3 className="font-bold text-slate-800 dark:text-slate-100 truncate pr-8">{scenario.name || 'Untitled Scenario'}</h3>
                </div>

                <p className="text-sm text-slate-500 dark:text-slate-400 line-clamp-2 mb-6 font-medium flex-1">
                  {scenario.description || 'No description provided.'}
                </p>

                <div className="flex items-center justify-between pt-4 border-t border-slate-50 dark:border-slate-800">
                  <div className="flex items-center gap-1.5 text-xs font-bold text-slate-400">
                    <ListOrdered size={14} />
                    {scenario.steps?.length || 0} STEPS
                  </div>
                  <div className="flex items-center gap-1.5 text-xs font-bold text-slate-400">
                    <Calendar size={14} />
                    {new Date(scenario.created_at).toLocaleDateString()}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
