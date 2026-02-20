import { useState, useCallback } from 'react';
import type { Config } from '../types/traffic';

export const useConfig = (toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [config, setConfig] = useState<Config>({
    proxy_addr: ':15500',
    api_addr: ':15501',
    mcp_addr: ':15502',
    mcp_enabled: false,
    history_limit: 500,
    max_response_size: 1048576,
    default_page_size: 50
  });
  const [originalConfig, setOriginalConfig] = useState<Config | null>(null);

  const fetchConfig = useCallback(async () => {
    try {
      const res = await fetch('/api/config');
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      setConfig(data);
      setOriginalConfig(data);
      return data as Config;
    } catch (error) {
      console.error('Error fetching config:', error);
    }
  }, []);

  const saveConfig = useCallback(async (newConfig: Config) => {
    try {
      const res = await fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newConfig)
      });
      if (!res.ok) throw new Error(await res.text());
      
      const needsRestart = originalConfig && (
        newConfig.proxy_addr !== originalConfig.proxy_addr ||
        newConfig.api_addr !== originalConfig.api_addr ||
        newConfig.mcp_addr !== originalConfig.mcp_addr ||
        newConfig.mcp_enabled !== originalConfig.mcp_enabled
      );

      setConfig(newConfig);
      setOriginalConfig(newConfig);

      if (needsRestart) {
        toast('info', 'Restart Required', 'Some changes will take effect only after you restart the Glance.');
      } else {
        toast('success', 'Settings Saved', 'Your configuration has been updated successfully.');
      }
    } catch (error) {
      toast('error', 'Save Failed', String(error));
    }
  }, [originalConfig, toast]);

  return {
    config,
    setConfig,
    fetchConfig,
    saveConfig,
  };
};
