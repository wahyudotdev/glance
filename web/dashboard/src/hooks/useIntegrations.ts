import { useState, useCallback } from 'react';
import type { JavaProcess, AndroidDevice } from '../types/traffic';

export const useIntegrations = (toast: (type: 'success' | 'error' | 'info', title: string, message: string) => void) => {
  const [javaProcesses, setJavaProcesses] = useState<JavaProcess[]>([]);
  const [androidDevices, setAndroidDevices] = useState<AndroidDevice[]>([]);
  const [isLoadingJava, setIsLoadingJava] = useState(false);
  const [isLoadingAndroid, setIsLoadingAndroid] = useState(false);
  const [terminalScript, setTerminalScript] = useState('');

  const apiFetch = useCallback(async (url: string, options?: RequestInit) => {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(await res.text());
    if (res.status === 204) return null;
    return res.json();
  }, []);

  const fetchJavaProcesses = useCallback(async () => {
    setIsLoadingJava(true);
    try {
      const data = await apiFetch('/api/client/java/processes');
      setJavaProcesses(data || []);
    } catch {
      console.warn('Error fetching Java processes');
    } finally {
      setIsLoadingJava(false);
    }
  }, [apiFetch]);

  const fetchAndroidDevices = useCallback(async () => {
    setIsLoadingAndroid(true);
    try {
      const data = await apiFetch('/api/client/android/devices');
      setAndroidDevices(data || []);
    } catch {
      toast('error', 'ADB Error', 'Could not list Android devices. Ensure adb is installed.');
    } finally {
      setIsLoadingAndroid(false);
    }
  }, [apiFetch, toast]);

  const fetchTerminalScript = useCallback(async () => {
    try {
      const res = await fetch('/api/client/terminal/setup');
      const text = await res.text();
      setTerminalScript(text);
    } catch {
      console.warn('Error fetching terminal script');
    }
  }, []);

  const interceptJava = useCallback(async (pid: string) => {
    try {
      await apiFetch(`/api/client/java/intercept/${pid}`, { method: 'POST' });
      toast('success', 'Interception Active', `Successfully injected proxy into PID ${pid}.`);
    } catch (error) {
      toast('error', 'Interception Failed', String(error));
    }
  }, [apiFetch, toast]);

  const interceptAndroid = useCallback(async (id: string) => {
    try {
      await apiFetch(`/api/client/android/intercept/${id}`, { method: 'POST' });
      toast('success', 'Proxy Configured', `Android device ${id} is now routing traffic through this proxy.`);
    } catch (error) {
      toast('error', 'Configuration Failed', String(error));
    }
  }, [apiFetch, toast]);

  const clearAndroid = useCallback(async (id: string) => {
    try {
      await apiFetch(`/api/client/android/clear/${id}`, { method: 'POST' });
      toast('success', 'Proxy Cleared', `Android device ${id} proxy settings have been reset.`);
    } catch (error) {
      toast('error', 'Reset Failed', String(error));
    }
  }, [apiFetch, toast]);

  const pushAndroidCert = useCallback(async (id: string) => {
    try {
      await apiFetch(`/api/client/android/push-cert/${id}`, { method: 'POST' });
      toast('success', 'CA Cert Pushed', 'Certificate pushed to /sdcard/ and install settings opened on device.');
    } catch (error) {
      toast('error', 'Push Failed', String(error));
    }
  }, [apiFetch, toast]);

  return {
    javaProcesses,
    androidDevices,
    isLoadingJava,
    isLoadingAndroid,
    terminalScript,
    fetchJavaProcesses,
    fetchAndroidDevices,
    fetchTerminalScript,
    interceptJava,
    interceptAndroid,
    clearAndroid,
    pushAndroidCert,
  };
};
