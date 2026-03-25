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
}

export interface DisplayMessage {
  id: number;
  pairId: string;
  role: "user" | "assistant";
  content: string;
  isStreaming?: boolean;
  toolCalls?: ToolCall[];
  thinking?: ThinkingState;
  // messages array allows interleaving content and tool calls
  blocks?: MessageBlock[];
}

export type MessageBlock = { type: "text"; content: string } | { type: "tool"; toolCall: ToolCall };
