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
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm animate-in fade-in duration-200">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden animate-in zoom-in-95 duration-200">
        <div className="p-8 text-center">
          <div className={`w-16 h-16 ${iconBgColor} rounded-full flex items-center justify-center mx-auto mb-6`}>
            {icon}
          </div>
          <h3 className="text-xl font-bold text-slate-800 mb-2">{title}</h3>
          <p className="text-slate-500 text-sm leading-relaxed">{description}</p>
        </div>
        <div className="flex border-t border-slate-100 p-4 gap-3 bg-slate-50/50">
          {showCancel && (
            <button 
              onClick={onClose}
              className="flex-1 px-4 py-3 text-sm font-bold text-slate-600 hover:bg-white rounded-xl transition-all"
            >
              Cancel
            </button>
          )}
          <button 
            onClick={onConfirm}
            className={`flex-1 px-4 py-3 text-sm font-bold text-white ${confirmColor} rounded-xl transition-all shadow-lg active:scale-95`}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
};
