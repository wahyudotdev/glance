import React from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description: string;
  icon: React.ReactNode;
  iconBgColor: string;
  confirmLabel: string;
  confirmColor: string;
  onConfirm: () => void;
  showCancel?: boolean;
}

export const Modal: React.FC<ModalProps> = ({ 
  isOpen, onClose, title, description, icon, 
  iconBgColor, confirmLabel, confirmColor, onConfirm, showCancel = true 
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-slate-900/40 backdrop-blur-md animate-in fade-in duration-300" onClick={onClose} />
      <div className="relative bg-white dark:bg-slate-900 rounded-3xl shadow-2xl shadow-slate-400/20 w-full max-w-sm overflow-hidden animate-in zoom-in-[0.98] duration-300 border border-white/50 dark:border-slate-800/50 transition-colors">
        <div className="p-8 text-center">
          <div className={`w-16 h-16 ${iconBgColor} rounded-2xl flex items-center justify-center mx-auto mb-6 shadow-sm ring-4 ring-white dark:ring-slate-800`}>
            {icon}
          </div>
          <h3 className="text-xl font-bold text-slate-800 dark:text-slate-100 mb-2 tracking-tight">{title}</h3>
          <p className="text-slate-500 dark:text-slate-400 text-sm leading-relaxed">{description}</p>
        </div>
        <div className="flex border-t border-slate-100 dark:border-slate-800 p-4 gap-3 bg-slate-50/50 dark:bg-slate-950/50 transition-colors">
          {showCancel && (
            <button 
              onClick={onClose}
              className="flex-1 px-4 py-3 text-sm font-bold text-slate-600 dark:text-slate-400 hover:bg-white dark:hover:bg-slate-800 hover:text-slate-800 dark:hover:text-slate-200 rounded-xl transition-all border border-transparent hover:border-slate-200 dark:hover:border-slate-700"
            >
              Cancel
            </button>
          )}
          <button 
            onClick={onConfirm}
            className={`flex-1 px-4 py-3 text-sm font-bold text-white ${confirmColor} rounded-xl transition-all shadow-lg active:scale-[0.98] hover:shadow-xl hover:-translate-y-0.5`}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
};
