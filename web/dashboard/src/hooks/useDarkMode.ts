import { useState, useEffect } from 'react';

export const useDarkMode = () => {
  const [isDark, setIsDark] = useState(() => {
    try {
      const saved = localStorage.getItem('glance-dark-mode');
      if (saved !== null) {
        return saved === 'true';
      }
    } catch (e) {
      console.warn('LocalStorage not accessible', e);
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });

  useEffect(() => {
    const root = window.document.documentElement;
    if (isDark) {
      root.classList.add('dark');
      root.classList.remove('light');
      root.style.colorScheme = 'dark';
    } else {
      root.classList.remove('dark');
      root.classList.add('light');
      root.style.colorScheme = 'light';
    }
    
    try {
      localStorage.setItem('glance-dark-mode', isDark.toString());
    } catch (e) {
      // Ignore
    }
  }, [isDark]);

  const toggleDarkMode = () => setIsDark(prev => !prev);

  return { isDark, toggleDarkMode };
};
