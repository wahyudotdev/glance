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
