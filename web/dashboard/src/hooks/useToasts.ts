import { useState, useCallback } from 'react';
import type { ToastMessage } from '../components/ui/Toast';

export const useToasts = () => {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  const toast = useCallback((type: 'success' | 'error' | 'info', title: string, message: string) => {
    const id = Math.random().toString(36).substring(2, 9);
    setToasts((prev) => [...prev, { id, type, title, message }]);
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  return {
    toasts,
    toast,
    removeToast,
  };
};
