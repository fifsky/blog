import { useState } from "react";
import { Wrench, ChevronDown, ChevronRight } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import type { ToolCall } from "./types";

interface ToolCallCardProps {
  toolCall: ToolCall;
}

export function ToolCallCard({ toolCall }: ToolCallCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const formatJson = (str: string) => {
    try {
      return JSON.stringify(JSON.parse(str), null, 2);
    } catch {
      return str;
    }
  };

  return (
    <div className="my-2 border border-gray-200 rounded-lg bg-gray-50 overflow-hidden">
      <button
        onClick={() => !toolCall.isLoading && setIsExpanded(!isExpanded)}
        disabled={toolCall.isLoading}
        className={`w-full px-3 py-2 flex items-center gap-2 text-left text-sm ${
          toolCall.isLoading ? "cursor-wait" : "hover:bg-gray-100 cursor-pointer"
        }`}
      >
        {toolCall.isLoading ? (
          <Spinner className="w-4 h-4" />
        ) : (
          <Wrench className="w-4 h-4 text-blue-500" />
        )}
        <span className="font-medium text-gray-700">{toolCall.mcpName}</span>
        <span className="text-xs text-gray-400">({toolCall.name})</span>
        {!toolCall.isLoading && (
          <span className="ml-auto text-gray-400">
            {isExpanded ? (
              <ChevronDown className="w-4 h-4" />
            ) : (
              <ChevronRight className="w-4 h-4" />
            )}
          </span>
        )}
        {toolCall.isLoading && <span className="ml-auto text-xs text-gray-400">调用中...</span>}
      </button>
      {isExpanded && !toolCall.isLoading && (
        <div className="px-3 py-2 border-t border-gray-200 bg-white">
          <div className="mb-2">
            <div className="text-xs font-medium text-gray-500 mb-1">参数</div>
            <pre className="text-xs bg-gray-100 p-2 rounded overflow-x-auto">
              {formatJson(toolCall.arguments)}
            </pre>
          </div>
          {toolCall.result && (
            <div>
              <div className="text-xs font-medium text-gray-500 mb-1">结果</div>
              <pre className="text-xs bg-gray-100 p-2 rounded overflow-x-auto max-h-40 overflow-y-auto">
                {formatJson(toolCall.result)}
              </pre>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
