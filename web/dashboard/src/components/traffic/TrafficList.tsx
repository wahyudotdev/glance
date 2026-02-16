import React from 'react';
import { ChevronRight } from 'lucide-react';
import dayjs from 'dayjs';
import type { TrafficEntry } from '../../types/traffic';

interface TrafficListProps {
  entries: TrafficEntry[];
  selectedEntry: TrafficEntry | null;
  onSelect: (entry: TrafficEntry) => void;
}

const getStatusColor = (status: number) => {
  if (status >= 200 && status < 300) return 'text-emerald-600 bg-emerald-50 border-emerald-100';
  if (status >= 300 && status < 400) return 'text-amber-600 bg-amber-50 border-amber-100';
  if (status >= 400) return 'text-rose-600 bg-rose-50 border-rose-100';
  return 'text-slate-600 bg-slate-50 border-slate-100';
};

export const TrafficList: React.FC<TrafficListProps> = ({ entries, selectedEntry, onSelect }) => {
  return (
    <div className="flex-1 overflow-y-auto bg-white">
      <table className="w-full text-left border-separate border-spacing-0">
        <thead className="sticky top-0 bg-white/80 backdrop-blur-md z-10 shadow-sm">
          <tr className="text-[11px] uppercase tracking-widest text-slate-400 font-bold border-b border-slate-100">
            <th className="pl-8 pr-4 py-4 font-bold">Method</th>
            <th className="px-4 py-4">Status</th>
            <th className="px-4 py-4">Path & Host</th>
            <th className="px-4 py-4">Size</th>
            <th className="pr-8 pl-4 py-4 text-right">Time</th>
          </tr>
        </thead>
        <tbody className="text-[13px]">
          {entries.map((entry) => (
            <tr
              key={entry.id}
              onClick={() => onSelect(entry)}
              className={`group cursor-pointer border-b border-slate-50 transition-all duration-150 hover:bg-blue-50/50 ${
                selectedEntry?.id === entry.id ? 'bg-blue-50 border-blue-100' : ''
              }`}
            >
              <td className="pl-8 pr-4 py-3.5">
                <span className="font-mono font-bold text-blue-600 tracking-tighter">
                  {entry.method}
                </span>
              </td>
              <td className="px-4 py-3.5">
                <span className={`px-2.5 py-1 rounded-md text-[11px] font-bold border tabular-nums ${getStatusColor(entry.status)}`}>
                  {entry.status || '---'}
                </span>
              </td>
              <td className="px-4 py-3.5 max-w-xl">
                <div className="flex flex-col">
                  <span className="text-slate-700 font-medium truncate font-mono">
                    {new URL(entry.url).pathname}{new URL(entry.url).search}
                  </span>
                  <span className="text-slate-400 text-[11px] truncate">
                    {new URL(entry.url).host}
                  </span>
                </div>
              </td>
              <td className="px-4 py-3.5 text-slate-400 tabular-nums">
                {entry.response_body ? `${(entry.response_body.length / 1024).toFixed(1)} KB` : '-'}
              </td>
              <td className="pr-8 pl-4 py-3.5 text-right text-slate-400 tabular-nums group-hover:text-slate-600 flex items-center justify-end gap-2">
                {dayjs(entry.start_time).format('HH:mm:ss.SSS')}
                <ChevronRight size={14} className={`opacity-0 group-hover:opacity-100 transition-opacity ${selectedEntry?.id === entry.id ? 'opacity-100' : ''}`} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
