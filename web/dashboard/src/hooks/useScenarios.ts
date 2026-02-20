import { useState, useCallback } from 'react';
import type { Scenario, TrafficEntry } from '../types/traffic';

export const useScenarios = (toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [isLoadingScenarios, setIsLoadingScenarios] = useState(false);

  const fetchScenarios = useCallback(async () => {
    setIsLoadingScenarios(true);
    try {
      const res = await fetch('/api/scenarios');
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      setScenarios(data || []);
    } catch (error) {
      toast('error', 'Fetch Scenarios Failed', String(error));
    } finally {
      setIsLoadingScenarios(false);
    }
  }, [toast]);

  const saveScenario = useCallback(async (scenario: Scenario) => {
    try {
      const isUpdate = !!scenario.id;
      const url = isUpdate ? `/api/scenarios/${scenario.id}` : '/api/scenarios';
      const method = isUpdate ? 'PUT' : 'POST';

      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(scenario),
      });
      if (!res.ok) throw new Error(await res.text());

      await fetchScenarios();
      toast('success', isUpdate ? 'Scenario Updated' : 'Scenario Created', `Successfully saved "${scenario.name}".`);
      return true;
    } catch (error) {
      toast('error', 'Save Failed', String(error));
      return false;
    }
  }, [fetchScenarios, toast]);

  const deleteScenario = useCallback(async (id: string) => {
    try {
      const res = await fetch(`/api/scenarios/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error(await res.text());
      setScenarios((prev) => prev.filter((s) => s.id !== id));
      toast('success', 'Scenario Deleted', 'The scenario has been removed.');
      return true;
    } catch (error) {
      toast('error', 'Delete Failed', String(error));
      return false;
    }
  }, [toast]);

  const addToScenario = useCallback(async (entry: TrafficEntry, scenarioId: string | 'new', onNewScenario: (entry: TrafficEntry) => void) => {
    if (scenarioId === 'new') {
      onNewScenario(entry);
      return;
    }

    const existing = scenarios.find(s => s.id === scenarioId);
    if (!existing) return;

    const updated: Scenario = {
      ...existing,
      steps: [
        ...(existing.steps || []),
        {
          id: '',
          traffic_entry_id: entry.id,
          order: (existing.steps?.length || 0) + 1,
          notes: ''
        }
      ]
    };

    try {
      const res = await fetch(`/api/scenarios/${scenarioId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updated)
      });
      if (!res.ok) throw new Error(await res.text());
      await fetchScenarios();
      toast('success', 'Added to Scenario', `Successfully added to "${existing.name}".`);
    } catch (error) {
      toast('error', 'Add Failed', String(error));
    }
  }, [scenarios, fetchScenarios, toast]);

  return {
    scenarios,
    isLoadingScenarios,
    fetchScenarios,
    saveScenario,
    deleteScenario,
    addToScenario,
  };
};
