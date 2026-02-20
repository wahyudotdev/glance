import { useState, useEffect, useRef, useMemo, useCallback } from 'react';
import type { TrafficEntry, Config } from '../types/traffic';

export const useTraffic = (config: Config, toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [totalEntries, setTotalEntries] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [proxyAddr, setProxyAddr] = useState(':8000');
  const [version, setVersion] = useState('dev');
  const [mcpSessions, setMcpSessions] = useState(0);
  const [mcpEnabled, setMcpEnabled] = useState(false);
  const [filter, setFilter] = useState('');
  const [methodFilter, setMethodFilter] = useState('ALL');

  const currentPageRef = useRef(currentPage);
  useEffect(() => { currentPageRef.current = currentPage; }, [currentPage]);

  const pageSizeRef = useRef(config.default_page_size);
  useEffect(() => { pageSizeRef.current = config.default_page_size; }, [config.default_page_size]);

  const apiFetch = useCallback(async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  }, []);

  const fetchTraffic = useCallback(async (page: number = 1, pageSize?: number) => {
    try {
      const size = pageSize || config.default_page_size;
      const data = await apiFetch(`/api/traffic?page=${page}&pageSize=${size}`);
      if (data && data.entries) {
        setEntries(data.entries.reverse());
        setTotalEntries(data.total);
        setCurrentPage(page);
      }
    } catch (error) {
      console.error('Error fetching traffic:', error);
    }
  }, [config.default_page_size, apiFetch]);

  const fetchStatus = useCallback(async () => {
    try {
      const data = await apiFetch('/api/status');
      setProxyAddr(data.proxy_addr);
      setVersion(data.version || 'dev');
      setMcpSessions(data.mcp_sessions || 0);
      setMcpEnabled(data.mcp_enabled || false);
    } catch (error) {
      console.error('Error fetching status:', error);
    }
  }, [apiFetch]);

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 10000); // Poll every 10s
    return () => clearInterval(interval);
  }, [fetchStatus]);

  const clearTraffic = useCallback(async () => {
    try {
      await apiFetch('/api/traffic', { method: 'DELETE' });
      setEntries([]);
      toast('success', 'Traffic Cleared', 'All intercepted requests have been deleted.');
      return true;
    } catch (error) {
      toast('error', 'Clear Failed', String(error));
      return false;
    }
  }, [apiFetch, toast]);

  const filteredEntries = useMemo(() => {
    return entries.filter(e => {
      const matchesText = e.url.toLowerCase().includes(filter.toLowerCase()) || 
                         e.method.toLowerCase().includes(filter.toLowerCase());
      const matchesMethod = methodFilter === 'ALL' || e.method.toUpperCase() === methodFilter;
      return matchesText && matchesMethod;
    }).reverse();
  }, [entries, filter, methodFilter]);

  return {
    entries,
    totalEntries,
    currentPage,
    proxyAddr,
    version,
    mcpSessions,
    mcpEnabled,
    filter,
    setFilter,
    methodFilter,
    setMethodFilter,
    filteredEntries,
    fetchTraffic,
    fetchStatus,
    clearTraffic,
    setEntries,
    setTotalEntries,
    currentPageRef,
    pageSizeRef
  };
};
