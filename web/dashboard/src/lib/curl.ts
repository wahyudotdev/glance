import type { TrafficEntry } from '../types/traffic';

export function generateCurl(entry: TrafficEntry): string {
  let curl = `curl -X ${entry.method} '${entry.url}'`;

  // Add headers
  Object.entries(entry.request_headers).forEach(([key, values]) => {
    values.forEach(value => {
      curl += ` 
  -H '${key}: ${value}'`;
    });
  });

  // Add body if present
  if (entry.request_body && entry.request_body.length > 0) {
    // Basic check if it's JSON-like or just string
    curl += ` 
  --data '${entry.request_body}'`;
  }

  return curl;
}
