export interface ToolCall {
  id: string;
  name: string;
  mcpName: string;
  arguments: string;
  result?: string;
  isLoading: boolean;
}

export interface DisplayMessage {
  id: number;
  pairId: string;
  role: "user" | "assistant";
  content: string;
  reasoningContent?: string;
  contextMessages?: Array<Record<string, unknown>>;
  isStreaming?: boolean;
  toolCalls?: ToolCall[];
}

export interface StreamEvent {
  type: "content" | "reasoning" | "tool_start" | "tool_end" | "context" | "error" | "done";
  content?: string;
  data?: any;
}
