export interface ToolCall {
  id: string;
  name: string;
  mcpName: string;
  arguments: string;
  result?: string;
  isLoading?: boolean;
}

export interface ThinkingState {
  content: string;
  isThinking: boolean;
  duration?: string;
  id?: string; // Optional unique ID to distinguish multiple thinking blocks
}

export interface DisplayMessage {
  id: number;
  pairId: string;
  role: "user" | "assistant";
  content: string;
  isStreaming?: boolean;
  toolCalls?: ToolCall[];
  // Legacy single thinking state kept for backwards compatibility if needed
  thinking?: ThinkingState;
  // messages array allows interleaving content, thinking, and tool calls
  blocks?: MessageBlock[];
}

export type MessageBlock =
  | { type: "text"; content: string }
  | { type: "tool"; toolCall: ToolCall }
  | { type: "thinking"; thinking: ThinkingState };
