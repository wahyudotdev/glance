import { useState, useCallback } from 'react';
import type { Rule } from '../types/traffic';

export const useRules = (toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [rules, setRules] = useState<Rule[]>([]);
  const [isLoadingRules, setIsLoadingRules] = useState(false);

  const fetchRules = useCallback(async () => {
    setIsLoadingRules(true);
    try {
      const res = await fetch('/api/rules');
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      setRules(data || []);
    } catch (error) {
      toast('error', 'Fetch Rules Failed', String(error));
    } finally {
      setIsLoadingRules(false);
    }
  }, [toast]);

  const createRule = useCallback(async (rule: Partial<Rule>) => {
    try {
      const res = await fetch('/api/rules', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rule),
      });
      if (!res.ok) throw new Error(await res.text());
      await fetchRules();
      toast('success', 'Rule Created', 'The new rule has been added.');
    } catch (error) {
      toast('error', 'Create Rule Failed', String(error));
    }
  }, [fetchRules, toast]);

  const updateRule = useCallback(async (id: string, rule: Partial<Rule>) => {
    try {
      const res = await fetch(`/api/rules/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rule),
      });
      if (!res.ok) throw new Error(await res.text());
      await fetchRules();
      toast('success', 'Rule Updated', 'The rule has been modified.');
    } catch (error) {
      toast('error', 'Update Rule Failed', String(error));
    }
  }, [fetchRules, toast]);

  const deleteRule = useCallback(async (id: string) => {
    try {
      const res = await fetch(`/api/rules/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error(await res.text());
      setRules((prev) => prev.filter((r) => r.id !== id));
      toast('success', 'Rule Deleted', 'The rule has been removed.');
    } catch (error) {
      toast('error', 'Delete Rule Failed', String(error));
    }
  }, [toast]);

  return {
    rules,
    isLoadingRules,
    fetchRules,
    createRule,
    updateRule,
    deleteRule,
  };
};
