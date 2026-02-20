import { useState, useEffect } from 'react';

export const useDarkMode = () => {
  const [isDark, setIsDark] = useState(() => {
    try {
      const saved = localStorage.getItem('glance-dark-mode');
      if (saved !== null) {
        return saved === 'true';
      }
    } catch {
      // ignore
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  });

  useEffect(() => {
    const root = window.document.documentElement;
    
    // Add a class to disable transitions temporarily
    root.classList.add('no-transitions');
    
    if (isDark) {
      root.classList.add('dark');
      root.classList.remove('light');
      root.style.colorScheme = 'dark';
    } else {
      root.classList.remove('dark');
      root.classList.add('light');
      root.style.colorScheme = 'light';
    }
    
    // Force a reflow to ensure the classes are applied without transition
    void window.getComputedStyle(root).opacity;
    
    // Remove the class after a short delay
    const timer = setTimeout(() => {
      root.classList.remove('no-transitions');
    }, 0);
    
    try {
      localStorage.setItem('glance-dark-mode', isDark.toString());
    } catch {
      // Ignore
    }

    return () => clearTimeout(timer);
  }, [isDark]);

  const toggleDarkMode = () => setIsDark(prev => !prev);

  return { isDark, toggleDarkMode };
};
