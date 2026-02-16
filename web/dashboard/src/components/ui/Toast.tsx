import React, { useEffect } from 'react';
import { Check, XCircle, Info, X } from 'lucide-react';

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'info';
  title: string;
  message: string;
}

interface ToastProps {
  toasts: ToastMessage[];
  onClose: (id: string) => void;
}

export const Toast: React.FC<ToastProps> = ({ toasts, onClose }) => {
  return (
    <div className="fixed top-6 right-6 z-[200] flex flex-col gap-3 pointer-events-none">
      {toasts.map((t) => (
        <ToastItem key={t.id} toast={t} onClose={() => onClose(t.id)} />
      ))}
    </div>
  );
};

const ToastItem: React.FC<{ toast: ToastMessage; onClose: () => void }> = ({ toast, onClose }) => {
  useEffect(() => {
    const timer = setTimeout(onClose, 3000);
    return () => clearTimeout(timer);
  }, [onClose]);

  const icons = {
    success: <Check size={18} className="text-emerald-500" />,
    error: <XCircle size={18} className="text-rose-500" />,
    info: <Info size={18} className="text-blue-500" />,
  };

  const bgs = {
    success: 'bg-emerald-50 border-emerald-100',
    error: 'bg-rose-50 border-rose-100',
    info: 'bg-blue-50 border-blue-100',
  };

  return (
    <div className={`pointer-events-auto min-w-[300px] max-w-md p-4 rounded-2xl border shadow-xl flex gap-3 animate-in slide-in-from-top-4 duration-300 ${bgs[toast.type]}`}>
      <div className="mt-0.5 shrink-0">{icons[toast.type]}</div>
      <div className="flex-1 min-w-0">
        <h4 className="text-sm font-bold text-slate-800 leading-tight">{toast.title}</h4>
        <p className="text-xs text-slate-600 mt-1 leading-relaxed line-clamp-2">{toast.message}</p>
      </div>
      <button onClick={onClose} className="shrink-0 p-1 hover:bg-black/5 rounded-lg h-fit transition-all">
        <X size={14} className="text-slate-400" />
      </button>
    </div>
  );
};
