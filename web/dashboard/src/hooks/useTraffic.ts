import { useState, useEffect, useRef, useMemo } from 'react';
import type { TrafficEntry, Config } from '../types/traffic';

export const useTraffic = (config: Config, toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [entries, setEntries] = useState<TrafficEntry[]>([]);
  const [totalEntries, setTotalEntries] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [proxyAddr, setProxyAddr] = useState(':8000');
  const [mcpSessions, setMcpSessions] = useState(0);
  const [mcpEnabled, setMcpEnabled] = useState(false);
  const [filter, setFilter] = useState('');

  const currentPageRef = useRef(currentPage);
  useEffect(() => { currentPageRef.current = currentPage; }, [currentPage]);

  const pageSizeRef = useRef(config.default_page_size);
  useEffect(() => { pageSizeRef.current = config.default_page_size; }, [config.default_page_size]);

  const apiFetch = async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  };

  const fetchTraffic = async (page: number = 1, pageSize?: number) => {
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
  };

  const fetchStatus = async () => {
    try {
      const data = await apiFetch('/api/status');
      setProxyAddr(data.proxy_addr);
      setMcpSessions(data.mcp_sessions || 0);
      setMcpEnabled(data.mcp_enabled || false);
    } catch (error) {
      console.error('Error fetching status:', error);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 10000); // Poll every 10s
    return () => clearInterval(interval);
  }, []);

  const clearTraffic = async () => {
    try {
      await apiFetch('/api/traffic', { method: 'DELETE' });
      setEntries([]);
      toast('success', 'Traffic Cleared', 'All intercepted requests have been deleted.');
      return true;
    } catch (error) {
      toast('error', 'Clear Failed', String(error));
      return false;
    }
  };

  const filteredEntries = useMemo(() => {
    return entries.filter(e => 
      e.url.toLowerCase().includes(filter.toLowerCase()) || 
      e.method.toLowerCase().includes(filter.toLowerCase())
    ).reverse();
  }, [entries, filter]);

  return {
    entries,
    totalEntries,
    currentPage,
    proxyAddr,
    mcpSessions,
    mcpEnabled,
    filter,
    setFilter,
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
