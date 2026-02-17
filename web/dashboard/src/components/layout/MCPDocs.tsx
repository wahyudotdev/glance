import React from 'react';
import { Bot, Zap, Code, ShieldAlert, X } from 'lucide-react';

interface MCPDocsProps {
  isOpen: boolean;
  onClose: () => void;
  mcpUrl: string;
}

export const MCPDocs: React.FC<MCPDocsProps> = ({ isOpen, onClose, mcpUrl }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[150] flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div className="bg-white dark:bg-slate-900 rounded-3xl shadow-2xl w-full max-w-2xl overflow-hidden animate-in zoom-in-95 duration-300 flex flex-col max-h-[90vh] border border-transparent dark:border-slate-800 transition-colors">
        <div className="p-6 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-950/50 transition-colors">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-xl">
              <Bot size={20} />
            </div>
            <div>
              <h3 className="font-bold text-slate-800 dark:text-slate-100">Model Context Protocol</h3>
              <p className="text-xs text-slate-500 dark:text-slate-400">Connect AI Agents to this Proxy</p>
            </div>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-white dark:hover:bg-slate-800 rounded-lg text-slate-400 dark:text-slate-500 transition-all">
            <X size={20} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-8 space-y-8">
          <section className="space-y-3">
            <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200 flex items-center gap-2">
              <Zap size={16} className="text-amber-500" />
              What is MCP?
            </h4>
            <p className="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
              Glance implements the <strong>Model Context Protocol (MCP)</strong>. This allows AI tools like Claude Desktop or other agents to directly see your traffic, create mocks, and execute requests on your behalf.
            </p>
          </section>

          <section className="space-y-4">
            <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200 flex items-center gap-2">
              <Code size={16} className="text-blue-500" />
              Claude Desktop Configuration
            </h4>
            <p className="text-xs text-slate-600 dark:text-slate-400">Add this to your <code>claude_desktop_config.json</code>:</p>
            <pre className="bg-slate-900 text-blue-300 p-4 rounded-xl text-[10px] font-mono overflow-x-auto border border-slate-800">
{`{
  "mcpServers": {
    "glance": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-sse", "${mcpUrl}"]
    }
  }
}`}
            </pre>
          </section>

          <section className="space-y-4">
            <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200 flex items-center gap-2">
              <ShieldAlert size={16} className="text-amber-500" />
              Available Tools
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {[
                { name: 'list_traffic', desc: 'Read intercepted logs' },
                { name: 'get_traffic_details', desc: 'Read headers & bodies' },
                { name: 'execute_request', desc: 'Run custom HTTP requests' },
                { name: 'list_rules', desc: 'See all active mocks/pauses' },
                { name: 'add_mock_rule', desc: 'Create static responses' },
                { name: 'add_breakpoint_rule', desc: 'Pause traffic for edit' },
                { name: 'delete_rule', desc: 'Remove an active rule' },
                { name: 'clear_traffic', desc: 'Reset history' }
              ].map(tool => (
                <div key={tool.name} className="p-3 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl transition-colors">
                  <code className="text-[10px] font-bold text-indigo-600 dark:text-indigo-400">{tool.name}</code>
                  <p className="text-[10px] text-slate-500 dark:text-slate-400 mt-1">{tool.desc}</p>
                </div>
              ))}
            </div>
          </section>
        </div>

        <div className="p-4 bg-slate-50 dark:bg-slate-950 border-t border-slate-100 dark:border-slate-800 flex justify-center transition-colors">
          <button 
            onClick={onClose}
            className="px-8 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 rounded-xl text-sm font-bold hover:bg-slate-50 dark:hover:bg-slate-700 transition-all"
          >
            Close Documentation
          </button>
        </div>
      </div>
    </div>
  );
};
