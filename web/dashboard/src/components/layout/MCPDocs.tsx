import React, { useState } from 'react';
import { Bot, Zap, Code, ShieldAlert, X, Copy, Check } from 'lucide-react';

interface MCPDocsProps {
  isOpen: boolean;
  onClose: () => void;
  mcpUrl: string;
}

export const MCPDocs: React.FC<MCPDocsProps> = ({ isOpen, onClose, mcpUrl }) => {
  const [copiedUrl, setCopiedUrl] = useState(false);
  const [copiedConfig, setCopiedConfig] = useState(false);
  const [copiedCLI, setCopiedCLI] = useState(false);

  if (!isOpen) return null;

  const handleCopyUrl = () => {
    navigator.clipboard.writeText(mcpUrl);
    setCopiedUrl(true);
    setTimeout(() => setCopiedUrl(false), 2000);
  };

  const handleCopyCLI = () => {
    const cmd = `claude mcp add --transport http glance ${mcpUrl}`;
    navigator.clipboard.writeText(cmd);
    setCopiedCLI(true);
    setTimeout(() => setCopiedCLI(false), 2000);
  };

  const handleCopyConfig = () => {
    const config = JSON.stringify({
      mcpServers: {
        glance: {
          type: "http",
          url: mcpUrl
        }
      }
    }, null, 2);
    navigator.clipboard.writeText(config);
    setCopiedConfig(true);
    setTimeout(() => setCopiedConfig(false), 2000);
  };

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
              Connection Guides
            </h4>
            
            <div className="space-y-6">
              <div>
                <div className="flex items-center justify-between mb-2">
                  <p className="text-[11px] font-black uppercase text-slate-400 tracking-widest">Claude Code (CLI)</p>
                  <button 
                    onClick={handleCopyCLI}
                    className="flex items-center gap-1.5 px-2 py-1 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-[10px] font-bold text-slate-600 dark:text-slate-300 hover:border-indigo-500 transition-all active:scale-95"
                  >
                    {copiedCLI ? <Check size={12} className="text-emerald-500" /> : <Copy size={12} />}
                    {copiedCLI ? 'Copied!' : 'Copy Command'}
                  </button>
                </div>
                <p className="text-xs text-slate-600 dark:text-slate-400 mb-2">Run this in your terminal:</p>
                <code className="text-[10px] bg-slate-900 text-blue-300 p-3 rounded-xl font-mono block border border-slate-800 select-all">
                  claude mcp add --transport http glance {mcpUrl}
                </code>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <p className="text-[11px] font-black uppercase text-slate-400 tracking-widest">Claude Desktop</p>
                  <button 
                    onClick={handleCopyConfig}
                    className="flex items-center gap-1.5 px-2 py-1 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-[10px] font-bold text-slate-600 dark:text-slate-300 hover:border-indigo-500 transition-all active:scale-95"
                  >
                    {copiedConfig ? <Check size={12} className="text-emerald-500" /> : <Copy size={12} />}
                    {copiedConfig ? 'Copied!' : 'Copy Config'}
                  </button>
                </div>
                <p className="text-xs text-slate-600 dark:text-slate-400 mb-2">Add to <code>claude_desktop_config.json</code>:</p>
                <pre className="bg-slate-900 text-blue-300 p-4 rounded-xl text-[10px] font-mono overflow-x-auto border border-slate-800">
{`{
  "mcpServers": {
    "glance": {
      "type": "http",
      "url": "${mcpUrl}"
    }
  }
}`}
                </pre>
              </div>

              <div>
                <p className="text-[11px] font-black uppercase text-slate-400 mb-2 tracking-widest">Cline / Roo Code / VS Code</p>
                <p className="text-xs text-slate-600 dark:text-slate-400 mb-2">For Gemini, GPT-4o, or Claude via VS Code extensions:</p>
                <div className="bg-slate-50 dark:bg-slate-800 p-4 rounded-xl border border-slate-100 dark:border-slate-700">
                  <ol className="text-xs text-slate-600 dark:text-slate-400 list-decimal ml-4 space-y-1">
                    <li>Open <strong>MCP Settings</strong> in the extension.</li>
                    <li>Add a new server with type <code>SSE</code>.</li>
                    <li>Paste the URL: <code className="text-indigo-600 dark:text-indigo-400 font-bold select-all">{mcpUrl}</code></li>
                  </ol>
                </div>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <p className="text-[11px] font-black uppercase text-slate-400 tracking-widest">Raw SSE Endpoint</p>
                  <button 
                    onClick={handleCopyUrl}
                    className="flex items-center gap-1.5 px-2 py-1 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-[10px] font-bold text-slate-600 dark:text-slate-300 hover:border-indigo-500 transition-all active:scale-95"
                  >
                    {copiedUrl ? <Check size={12} className="text-emerald-500" /> : <Copy size={12} />}
                    {copiedUrl ? 'Copied!' : 'Copy URL'}
                  </button>
                </div>
                <p className="text-xs text-slate-600 dark:text-slate-400 mb-1">For custom implementations or other agents:</p>
                <code className="text-[10px] bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded font-mono text-slate-500 select-all">{mcpUrl}</code>
              </div>
            </div>
          </section>

          <section className="space-y-4">
            <h4 className="text-sm font-bold text-slate-800 dark:text-slate-200 flex items-center gap-2">
              <ShieldAlert size={16} className="text-amber-500" />
              Available Tools
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {[
                { name: 'inspect_network_traffic', desc: 'PRIMARY: Read logs (configurable limit)' },
                { name: 'inspect_request_details', desc: 'MANDATORY: Read headers & bodies' },
                { name: 'execute_request', desc: 'Run custom HTTP requests' },
                { name: 'list_rules', desc: 'See all active mocks/pauses' },
                { name: 'add_mock_rule', desc: 'Create static responses' },
                { name: 'add_breakpoint_rule', desc: 'Pause traffic for edit' },
                { name: 'delete_rule', desc: 'Remove an active rule' },
                { name: 'list_scenarios', desc: 'Read saved sequences' },
                { name: 'get_scenario', desc: 'Full sequence context' },
                { name: 'add_scenario', desc: 'Create new sequence' },
                { name: 'update_scenario', desc: 'Edit steps & mappings' },
                { name: 'delete_scenario', desc: 'Remove a scenario' },
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
