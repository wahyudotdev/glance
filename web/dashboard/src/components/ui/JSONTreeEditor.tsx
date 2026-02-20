import React, { useState, useEffect, useRef } from 'react';
import { ChevronRight, ChevronDown, Type, Braces, Edit3, Check, X, Search, ChevronUp, AlignLeft } from 'lucide-react';

interface JSONTreeEditorProps {
  value: string;
  onChange: (newValue: string) => void;
  className?: string;
  isFullScreen?: boolean;
  readOnly?: boolean;
}

export const JSONTreeEditor: React.FC<JSONTreeEditorProps> = ({ value, onChange, className, isFullScreen, readOnly }) => {
  const [mode, setMode] = useState<'code' | 'tree'>('code');
  const [parsedData, setParsedData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  
  // Search state
  const [searchQuery, setSearchQuery] = useState('');
  const [searchIndex, setSearchIndex] = useState(-1);
  const [searchResults, setSearchResults] = useState<number[]>([]);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-format logic
  const [internalValue, setInternalValue] = useState(value);

  useEffect(() => {
    try {
      const data = JSON.parse(value);
      setParsedData(data);
      setError(null);

      // If it's valid JSON and we are in readOnly mode, ensure it's formatted
      if (readOnly) {
        const formatted = JSON.stringify(data, null, 2);
        if (formatted !== value) {
          setInternalValue(formatted);
        } else {
          setInternalValue(value);
        }
      } else {
        setInternalValue(value);
      }
    } catch (e) {
      setError('Invalid JSON');
      setInternalValue(value);
      if (mode === 'tree') setMode('code');
    }
  }, [value, readOnly]);

  // Search logic for code mode
  useEffect(() => {
    if (!searchQuery || mode !== 'code') {
      setSearchResults([]);
      setSearchIndex(-1);
      return;
    }

    const results: number[] = [];
    const lowerValue = internalValue.toLowerCase();
    const lowerQuery = searchQuery.toLowerCase();
    let pos = lowerValue.indexOf(lowerQuery);
    while (pos !== -1) {
      results.push(pos);
      pos = lowerValue.indexOf(lowerQuery, pos + 1);
    }
    setSearchResults(results);
    setSearchIndex(results.length > 0 ? 0 : -1);
  }, [searchQuery, internalValue, mode]);

  useEffect(() => {
    if (searchIndex !== -1 && searchResults[searchIndex] !== undefined && textareaRef.current && mode === 'code') {
      const pos = searchResults[searchIndex];
      textareaRef.current.focus();
      textareaRef.current.setSelectionRange(pos, pos + searchQuery.length);
      
      const lineHeight = 16;
      const line = internalValue.substring(0, pos).split('\n').length;
      textareaRef.current.scrollTop = (line - 5) * lineHeight;
    }
  }, [searchIndex, searchResults, mode, internalValue]);

  const handleTreeChange = (newData: any) => {
    setParsedData(newData);
    onChange(JSON.stringify(newData, null, 2));
  };

  const prettifyJson = () => {
    try {
      const data = JSON.parse(value);
      onChange(JSON.stringify(data, null, 2));
    } catch (e) {}
  };

  const nextResult = (e: React.MouseEvent) => {
    e.preventDefault();
    if (searchResults.length > 0) {
      setSearchIndex((searchIndex + 1) % searchResults.length);
    }
  };

  const prevResult = (e: React.MouseEvent) => {
    e.preventDefault();
    if (searchResults.length > 0) {
      setSearchIndex((searchIndex - 1 + searchResults.length) % searchResults.length);
    }
  };

  return (
    <div className={`flex flex-col border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden bg-slate-900 ${className} ${isFullScreen ? 'rounded-none border-none' : ''}`}>
      <div className="flex items-center justify-between px-4 py-2 bg-slate-50 dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 shrink-0">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1 p-1 bg-slate-200 dark:bg-slate-800 rounded-lg">
            <button 
              onClick={() => setMode('code')}
              className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${mode === 'code' ? 'bg-white dark:bg-slate-700 text-blue-600 dark:text-blue-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}`}
            >
              <Braces size={12} /> Code
            </button>
            <button 
              onClick={() => !error && setMode('tree')}
              disabled={!!error}
              className={`flex items-center gap-1.5 px-3 py-1 rounded-md text-[10px] font-bold transition-all ${mode === 'tree' ? 'bg-white dark:bg-slate-700 text-emerald-600 dark:text-emerald-400 shadow-sm' : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 disabled:opacity-50 disabled:cursor-not-allowed'}`}
            >
              <Type size={12} /> Tree
            </button>
          </div>

          {mode === 'code' && (
            <div className="flex items-center gap-2 px-3 py-1 bg-white dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-800 shadow-sm">
              <Search size={12} className="text-slate-400" />
              <input 
                type="text" 
                placeholder="Find..."
                className="bg-transparent border-none text-[10px] outline-none w-24 dark:text-slate-200 font-bold"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              {searchResults.length > 0 && (
                <div className="flex items-center gap-1 ml-1 pl-2 border-l border-slate-100 dark:border-slate-800">
                  <span className="text-[9px] font-mono text-slate-500 tabular-nums">{searchIndex + 1}/{searchResults.length}</span>
                  <button onClick={prevResult} className="p-0.5 hover:bg-slate-100 dark:hover:bg-slate-800 rounded"><ChevronUp size={10} /></button>
                  <button onClick={nextResult} className="p-0.5 hover:bg-slate-100 dark:hover:bg-slate-800 rounded"><ChevronDown size={10} /></button>
                </div>
              )}
            </div>
          )}

          {!readOnly && mode === 'code' && !error && (
            <button 
              onClick={prettifyJson}
              className="flex items-center gap-1.5 px-3 py-1 bg-emerald-50 dark:bg-emerald-900/40 text-emerald-600 dark:text-emerald-400 rounded-lg text-[10px] font-bold hover:bg-emerald-100 dark:hover:bg-emerald-900/60 transition-all"
            >
              <AlignLeft size={12} /> Prettify
            </button>
          )}
        </div>
        {error ? (
          <span className="text-[10px] font-bold text-rose-500 uppercase tracking-tight">{error}</span>
        ) : mode === 'tree' && !readOnly && (
          <span className="text-[10px] font-bold text-emerald-500 uppercase tracking-tight">Tree Editor Active</span>
        )}
      </div>

      <div className="flex-1 min-h-0 text-slate-300 font-mono text-xs overflow-auto custom-scrollbar">
        {mode === 'code' ? (
          <textarea 
            ref={textareaRef}
            value={internalValue}
            onChange={(e) => !readOnly && onChange(e.target.value)}
            readOnly={readOnly}
            className="w-full h-full p-6 bg-transparent text-emerald-400 focus:outline-none resize-none leading-relaxed"
            spellCheck={false}
          />
        ) : (
          <div className="p-4">
            <TreeNode 
              label="root" 
              value={parsedData} 
              onUpdate={(val) => !readOnly && handleTreeChange(val)}
              isLast={true}
              readOnly={readOnly}
            />
          </div>
        )}
      </div>
    </div>
  );
};

interface TreeNodeProps {
  label: string | number;
  value: any;
  onUpdate: (val: any) => void;
  isLast: boolean;
  depth?: number;
  readOnly?: boolean;
}

const TreeNode: React.FC<TreeNodeProps> = ({ label, value, onUpdate, isLast, depth = 0, readOnly }) => {
  const [isExpanded, setIsExpanded] = useState(depth < 2);
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState('');

  const isObject = value !== null && typeof value === 'object';
  const isArray = Array.isArray(value);

  const startEdit = () => {
    if (readOnly) return;
    setEditValue(typeof value === 'string' ? value : JSON.stringify(value));
    setIsEditing(true);
  };

  const saveEdit = () => {
    try {
      // Try to parse as JSON first (handles numbers, booleans, null, and objects if pasted)
      const parsed = JSON.parse(editValue);
      onUpdate(parsed);
      setIsEditing(false);
    } catch (e) {
      // If parsing fails, treat as raw string (unquoted)
      onUpdate(editValue);
      setIsEditing(false);
    }
  };

  if (isObject) {
    const keys = Object.keys(value);
    const isEmpty = keys.length === 0;

    return (
      <div className="flex flex-col">
        <div className="flex items-center gap-1 group py-0.5">
          <button 
            onClick={() => setIsExpanded(!isExpanded)}
            className={`p-0.5 hover:bg-slate-800 rounded transition-colors ${isEmpty ? 'invisible' : 'text-slate-500'}`}
          >
            {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
          </button>
          
          <span className="text-blue-400 font-bold">{label}:</span>
          <span className="text-slate-500 ml-1">{isArray ? '[' : '{'}</span>
          {!isExpanded && <span className="text-slate-600 mx-1 italic text-[10px]">... {isArray ? value.length : keys.length} items ...</span>}
          {!isExpanded && <span className="text-slate-500">{isArray ? ']' : '}'}{!isLast && ','}</span>}
        </div>

        {isExpanded && (
          <div className="border-l border-slate-800/50 ml-2 pl-4 flex flex-col my-0.5">
            {keys.map((key, i) => (
              <TreeNode 
                key={key}
                label={isArray ? i : key}
                value={value[key]}
                onUpdate={(newVal) => {
                  if (isArray) {
                    const next = [...value];
                    next[i] = newVal;
                    onUpdate(next);
                  } else {
                    onUpdate({ ...value, [key]: newVal });
                  }
                }}
                isLast={i === keys.length - 1}
                depth={depth + 1}
                readOnly={readOnly}
              />
            ))}
          </div>
        )}

        {isExpanded && (
          <div className="flex items-center py-0.5">
            <span className="p-0.5 invisible"><ChevronRight size={14} /></span>
            <span className="text-slate-500">{isArray ? ']' : '}'}{!isLast && ','}</span>
          </div>
        )}
      </div>
    );
  }

  // Primitive value
  return (
    <div className="flex items-center gap-2 group py-0.5 min-h-[24px]">
      <span className="p-0.5 invisible"><ChevronRight size={14} /></span>
      <span className="text-blue-400 font-bold">{label}:</span>
      
      {isEditing && !readOnly ? (
        <div className="flex items-center gap-1 flex-1 max-w-md">
          <input 
            autoFocus
            type="text"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') saveEdit();
              if (e.key === 'Escape') setIsEditing(false);
            }}
            className="flex-1 bg-slate-800 border border-blue-500/50 rounded px-2 py-0.5 text-emerald-400 outline-none font-mono text-xs shadow-inner"
          />
          <button onClick={saveEdit} className="p-1 text-emerald-500 hover:bg-emerald-500/10 rounded transition-colors"><Check size={14} /></button>
          <button onClick={() => setIsEditing(false)} className="p-1 text-rose-500 hover:bg-rose-500/10 rounded transition-colors"><X size={14} /></button>
        </div>
      ) : (
        <div className={`flex items-center gap-2 ${!readOnly ? 'cursor-pointer' : ''}`} onClick={startEdit}>
          <span className={`
            ${typeof value === 'string' ? 'text-emerald-400' : ''}
            ${typeof value === 'number' ? 'text-amber-400' : ''}
            ${typeof value === 'boolean' ? 'text-purple-400' : ''}
            ${value === null ? 'text-slate-500' : ''}
          `}>
            {typeof value === 'string' ? `"${value}"` : String(value)}
            {!isLast && <span className="text-slate-500">,</span>}
          </span>
          {!readOnly && (
            <button 
              className="opacity-0 group-hover:opacity-100 p-1 text-slate-500 hover:text-blue-400 transition-all"
            >
              <Edit3 size={12} />
            </button>
          )}
        </div>
      )}
    </div>
  );
};
