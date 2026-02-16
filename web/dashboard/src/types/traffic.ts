export interface TrafficEntry {
  id: string;
  method: string;
  url: string;
  request_headers: Record<string, string[]>;
  request_body: string; // Base64 or string depending on content
  status: number;
  response_headers: Record<string, string[]>;
  response_body: string;
  start_time: string;
  duration: number; // in nanoseconds
}
