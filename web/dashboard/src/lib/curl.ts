import type { TrafficEntry } from '../types/traffic';

export function generateCurl(entry: TrafficEntry): string {
  let parts = [`curl -X ${entry.method} '${entry.url}'` ];

  // Add headers
  Object.entries(entry.request_headers).forEach(([key, values]) => {
    values.forEach(value => {
      parts.push(`  -H '${key}: ${value}'`);
    });
  });

  // Add body if present
  if (entry.request_body && entry.request_body.length > 0) {
    parts.push(`  --data '${entry.request_body}'`);
  }

  return parts.join(' \\\n');
}
