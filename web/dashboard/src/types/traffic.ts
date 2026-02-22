export interface TrafficEntry {
  id: string;
  method: string;
  url: string;
  request_headers: Record<string, string[]>;
  request_body: string;
  response_headers?: Record<string, string[]>;
  response_body?: string;
  status: number;
  start_time: string;
  duration: number;
  modified_by?: 'mock' | 'breakpoint' | 'editor';
}

export interface Config {
  proxy_addr: string;
  api_addr: string;
  mcp_addr: string;
  mcp_enabled: boolean;
  history_limit: number;
  max_response_size: number;
  default_page_size: number;
}

export interface JavaProcess {
  pid: string;
  name: string;
}

export interface AndroidDevice {
  id: string;
  model: string;
  name: string;
}

export interface VariableMapping {
  name: string;
  source_entry_id: string;
  source_path: string;
  target_json_path: string;
}

export interface MockResponse {
  status: number;
  headers: Record<string, string>;
  body: string;
}

export interface Rule {
  id: string;
  enabled: boolean;
  type: 'mock' | 'breakpoint';
  url_pattern: string;
  method: string;
  strategy?: string;
  response?: MockResponse;
}

export interface ScenarioStep {
  id: string;
  traffic_entry_id: string;
  order: number;
  notes?: string;
  traffic_entry?: TrafficEntry;
}

export interface Scenario {
  id: string;
  name: string;
  description: string;
  steps: ScenarioStep[];
  variable_mappings: VariableMapping[];
  created_at: string;
}
