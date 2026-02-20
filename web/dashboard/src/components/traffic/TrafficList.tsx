import React from 'react';
import { ChevronRight, ShieldAlert, Eye, Edit3 } from 'lucide-react';
import type { TrafficEntry } from '../../types/traffic';

interface TrafficListProps {
  entries: TrafficEntry[];
  selectedEntry: TrafficEntry | null;
  onSelect: (entry: TrafficEntry) => void;
}

const formatTime = (isoString: string) => {
  return new Intl.DateTimeFormat('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    fractionalSecondDigits: 3,
    hour12: false
  }).format(new Date(isoString));
};

const parseUrl = (urlString: string) => {
  try {
    const url = new URL(urlString);
    return {
      path: url.pathname + url.search,
      host: url.host
    };
  } catch {
    return {
      path: urlString,
      host: 'unknown'
    };
  }
};

const getStatusColor = (status: number) => {
  if (status >= 200 && status < 300) return 'text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 border-emerald-100 dark:border-emerald-800/30';
  if (status >= 300 && status < 400) return 'text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border-amber-100 dark:border-amber-800/30';
  if (status >= 400) return 'text-rose-600 dark:text-rose-400 bg-rose-50 dark:bg-rose-900/20 border-rose-100 dark:border-rose-800/30';
  return 'text-slate-600 dark:text-slate-400 bg-slate-50 dark:bg-slate-800 border-slate-100 dark:border-slate-700';
};

export const TrafficList: React.FC<TrafficListProps> = ({ entries, selectedEntry, onSelect }) => {
  return (
    <div className="flex-1 overflow-y-auto bg-white dark:bg-slate-900 transition-colors">
      <table className="w-full text-left border-separate border-spacing-0">
        <thead className="sticky top-0 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md z-10 shadow-sm transition-colors">
          <tr className="text-[11px] uppercase tracking-widest text-slate-400 dark:text-slate-500 font-bold border-b border-slate-100 dark:border-slate-800">
            <th className="pl-8 pr-4 py-4 font-bold">Method</th>
            <th className="px-4 py-4">Status</th>
            <th className="px-4 py-4">Path & Host</th>
            <th className="px-4 py-4 text-right">Latency</th>
            <th className="px-4 py-4">Size</th>
            <th className="pr-8 pl-4 py-4 text-right">Time</th>
          </tr>
        </thead>
        <tbody className="text-[13px]">
          {entries.map((entry) => (
            <tr
              key={entry.id}
              onClick={(e) => {
                e.stopPropagation();
                onSelect(entry);
              }}
              className={`group cursor-pointer border-b border-slate-50 dark:border-slate-800 transition-all duration-150 hover:bg-blue-50/50 dark:hover:bg-blue-900/10 ${
                selectedEntry?.id === entry.id ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-100 dark:border-blue-800' : ''
              }`}
            >
                                      <td className="pl-8 pr-4 py-3.5">
                                        <div className="flex flex-col gap-1">
                                          <span className="font-mono font-bold text-blue-600 dark:text-blue-400 tracking-tighter">
                                            {entry.method}
                                          </span>
                                          {entry.modified_by === 'mock' && (
                                            <span className="flex items-center gap-1 text-[9px] font-bold text-emerald-600 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 px-1.5 py-0.5 rounded border border-emerald-100 dark:border-emerald-800/30 w-fit">
                                              <Eye size={10} /> MOCKED
                                            </span>
                                          )}
                                          {entry.modified_by === 'breakpoint' && (
                                            <span className="flex items-center gap-1 text-[9px] font-bold text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 px-1.5 py-0.5 rounded border border-amber-100 dark:border-amber-800/30 w-fit">
                                              <ShieldAlert size={10} /> PAUSED
                                            </span>
                                          )}
                                          {entry.modified_by === 'editor' && (
                                            <span className="flex items-center gap-1 text-[9px] font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/20 px-1.5 py-0.5 rounded border border-indigo-100 dark:border-indigo-800/30 w-fit">
                                              <Edit3 size={10} /> EDITOR
                                            </span>
                                          )}
                                        </div>
                                      </td>
              
              <td className="px-4 py-3.5">
                <span className={`px-2.5 py-1 rounded-md text-[11px] font-bold border tabular-nums ${getStatusColor(entry.status)}`}>
                  {entry.status || '---'}
                </span>
              </td>
              <td className="px-4 py-3.5 max-w-xl">
                <div className="flex flex-col">
                  <span className="text-slate-700 dark:text-slate-200 font-medium truncate font-mono">
                    {parseUrl(entry.url).path}
                  </span>
                  <span className="text-slate-400 dark:text-slate-500 text-[11px] truncate">
                    {parseUrl(entry.url).host}
                  </span>
                </div>
              </td>
                                      <td className="px-4 py-3.5 text-right font-mono text-[11px] tabular-nums">
                                        <span className={`${entry.duration / 1000000 > 500 ? 'text-amber-600 dark:text-amber-400' : 'text-slate-400 dark:text-slate-500'}`}>
                                          {entry.duration > 0 ? `${(entry.duration / 1000000).toFixed(0)}ms` : '-'}
                                        </span>
                                      </td>
                                      <td className="px-4 py-3.5 text-slate-400 dark:text-slate-500 tabular-nums text-center">
                                        {entry.response_body ? `${(entry.response_body.length / 1024).toFixed(1)} KB` : '-'}
                                      </td>
              
                                      <td className="pr-8 pl-4 py-3.5 text-right text-slate-400 dark:text-slate-500 tabular-nums group-hover:text-slate-600 dark:group-hover:text-slate-300 flex items-center justify-end gap-2">
                                        {formatTime(entry.start_time)}
                                        <ChevronRight size={14} className={`opacity-0 group-hover:opacity-100 transition-opacity ${selectedEntry?.id === entry.id ? 'opacity-100' : ''}`} />
                                      </td>
              
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
