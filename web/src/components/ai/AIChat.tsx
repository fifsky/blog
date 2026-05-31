import { useState, useRef, useEffect, useCallback } from "react";
import {
  MessageCircle,
  X,
  Send,
  Copy,
  Check,
  Trash2,
  RotateCcw,
  Maximize2,
  Minimize2,
  BrainCircuit,
  ChevronRight,
} from "lucide-react";
import { AgentChatIndicator } from "@/components/agents-ui/agent-chat-indicator";
import { Button } from "@/components/ui/button";
import { Viewer } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import breaks from "@bytemd/plugin-breaks";
import { highlightPlugin } from "@/lib/highlight-plugin";
import { getApiUrl } from "@/utils/common";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import {
  getAllMessages,
  addMessagePair,
  updateAssistantMessage,
  updateAssistantContextMessages,
  deleteMessagePair,
  clearAllMessages,
} from "@/lib/chat-db";
import { ToolCallCard } from "./ToolCallCard";
import type { ToolCall, DisplayMessage, StreamEvent } from "./types";

// ByteMD plugins for rendering
const plugins = [gfm(), breaks(), highlightPlugin()];

const sanitize = (schema: any) => {
  const tags = schema.tagNames || [];
  if (!tags.includes("iframe")) {
    tags.push("iframe");
  }
  const attributes = schema.attributes || {};
  attributes["iframe"] = [
    "src",
    "width",
    "height",
    "frameborder",
    "allow",
    "allowfullscreen",
    "scrolling",
    "style",
    "className",
  ];
  return { ...schema, tagNames: tags, attributes };
};

const STORAGE_KEY = "ai-chat-button-position";
const DEFAULT_POSITION = { bottom: 24, right: 80 };

interface ButtonPosition {
  bottom: number;
  right: number;
}

function ReasoningBlock({
  content,
  isFinished,
  isStreaming,
}: {
  content: string;
  isFinished: boolean;
  isStreaming?: boolean;
}) {
  const [isOpen, setIsOpen] = useState(!isFinished);
  const [hasAutoClosed, setHasAutoClosed] = useState(false);

  useEffect(() => {
    if (isFinished && !hasAutoClosed) {
      setIsOpen(false);
      setHasAutoClosed(true);
    }
  }, [isFinished, hasAutoClosed]);

  return (
    <div className="mb-3 border border-gray-200 rounded-lg bg-gray-50 overflow-hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-3 py-2 flex items-center gap-2 text-sm cursor-pointer hover:bg-gray-100 text-gray-600 font-medium select-none"
      >
        <BrainCircuit className="w-4 h-4 text-purple-500" />
        深度思考
        {!isFinished && isStreaming && <AgentChatIndicator size="sm" className="ml-1" />}
        <ChevronRight
          className={`w-4 h-4 ml-auto text-gray-400 transition-transform ${isOpen ? "rotate-90" : ""}`}
        />
      </button>
      {isOpen && (
        <div className="px-3 py-2 border-t border-gray-200 bg-white text-xs text-gray-500 whitespace-pre-wrap font-mono leading-relaxed">
          {content}
        </div>
      )}
    </div>
  );
}

export function AIChat() {
  const [isOpen, setIsOpen] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [copiedId, setCopiedId] = useState<number | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Draggable button state
  const [buttonPosition, setButtonPosition] = useState<ButtonPosition>(DEFAULT_POSITION);
  const [isDragging, setIsDragging] = useState(false);
  const dragStartRef = useRef<{ x: number; y: number; bottom: number; right: number } | null>(null);
  const hasMovedRef = useRef(false);
  const buttonRef = useRef<HTMLButtonElement>(null);

  // Load button position from localStorage on mount
  useEffect(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        const parsed = JSON.parse(saved) as ButtonPosition;
        // Validate the saved position is within viewport
        const maxRight = window.innerWidth - 48; // button width
        const maxBottom = window.innerHeight - 48; // button height
        setButtonPosition({
          bottom: Math.max(0, Math.min(parsed.bottom, maxBottom)),
          right: Math.max(0, Math.min(parsed.right, maxRight)),
        });
      } catch {
        // Use default position if parsing fails
      }
    }
  }, []);

  // Handle drag start
  const handleDragStart = useCallback(
    (e: React.MouseEvent | React.TouchEvent) => {
      e.preventDefault();
      const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
      const clientY = "touches" in e ? e.touches[0].clientY : e.clientY;

      dragStartRef.current = {
        x: clientX,
        y: clientY,
        bottom: buttonPosition.bottom,
        right: buttonPosition.right,
      };
      setIsDragging(true);
      hasMovedRef.current = false;
    },
    [buttonPosition],
  );

  // Handle drag move and end
  useEffect(() => {
    if (!isDragging) return;

    const handleMove = (e: MouseEvent | TouchEvent) => {
      if (!dragStartRef.current) return;

      const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
      const clientY = "touches" in e ? e.touches[0].clientY : e.clientY;

      const deltaX = dragStartRef.current.x - clientX;
      const deltaY = dragStartRef.current.y - clientY;

      const maxRight = window.innerWidth - 48;
      const maxBottom = window.innerHeight - 48;

      const newRight = Math.max(0, Math.min(dragStartRef.current.right + deltaX, maxRight));
      const newBottom = Math.max(0, Math.min(dragStartRef.current.bottom + deltaY, maxBottom));

      setButtonPosition({ bottom: newBottom, right: newRight });
      hasMovedRef.current = true;
    };

    const handleEnd = () => {
      if (dragStartRef.current && hasMovedRef.current) {
        // Save position to localStorage only if actually moved
        localStorage.setItem(STORAGE_KEY, JSON.stringify(buttonPosition));
      }
      setIsDragging(false);
      dragStartRef.current = null;
      // Reset hasMoved after a short delay to allow onClick to check it
      setTimeout(() => {
        hasMovedRef.current = false;
      }, 100);
    };

    document.addEventListener("mousemove", handleMove);
    document.addEventListener("mouseup", handleEnd);
    document.addEventListener("touchmove", handleMove);
    document.addEventListener("touchend", handleEnd);

    return () => {
      document.removeEventListener("mousemove", handleMove);
      document.removeEventListener("mouseup", handleEnd);
      document.removeEventListener("touchmove", handleMove);
      document.removeEventListener("touchend", handleEnd);
    };
  }, [isDragging, buttonPosition]);

  // Load messages from Dexie on mount
  useEffect(() => {
    loadMessages();
  }, []);

  const loadMessages = async () => {
    const dbMessages = await getAllMessages();
    setMessages(
      dbMessages.map((m) => {
        let content = m.content;
        if (m.reasoningContent) {
          content = `<think>\n${m.reasoningContent}\n</think>\n` + content;
        }
        return {
          id: m.id!,
          pairId: m.pairId,
          role: m.role,
          content: content,
          contextMessages: m.contextMessages,
          toolCalls: m.toolCalls,
          isStreaming: false,
        };
      }),
    );
  };

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Focus input when chat opens
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  }, [isOpen]);

  // Copy message content to clipboard
  const copyToClipboard = async (content: string, id: number) => {
    try {
      await navigator.clipboard.writeText(content);
      setCopiedId(id);
      setTimeout(() => setCopiedId(null), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  // Delete a message pair
  const handleDeletePair = async (pairId: string) => {
    await deleteMessagePair(pairId);
    setMessages((prev) => prev.filter((m) => m.pairId !== pairId));
  };

  // Clear all messages
  const handleClearAll = async () => {
    await clearAllMessages();
    setMessages([]);
  };

  // Send message to AI
  const sendMessage = useCallback(async () => {
    if (!input.trim() || isLoading) return;

    const pairId = Date.now().toString();
    const userContent = input.trim();

    // Add to Dexie and get IDs
    const { userMsg, assistantMsg } = await addMessagePair(pairId, userContent);

    // Add to local state
    const newUserMessage: DisplayMessage = {
      id: userMsg.id!,
      pairId,
      role: "user",
      content: userContent,
    };

    const newAssistantMessage: DisplayMessage = {
      id: assistantMsg.id!,
      pairId,
      role: "assistant",
      content: "",
      isStreaming: true,
      toolCalls: [],
    };

    setMessages((prev) => [...prev, newUserMessage, newAssistantMessage]);
    setInput("");
    setIsLoading(true);

    let accumulatedContent = "";
    let currentToolCalls: ToolCall[] = [];

    try {
      const controller = new AbortController();
      abortControllerRef.current = controller;

      // Build messages array for API (history + current)
      const historyMessages = messages
        .filter((m) => m.content.trim() !== "")
        .map((m) => ({
          role: m.role,
          content: m.content
            .replace(/<tool_call id="[^"]+"><\/tool_call>/g, "")
            .replace(/<think>[\s\S]*?<\/think>/g, ""),
          contextMessages: m.contextMessages,
        }));

      // Add the new user message
      const apiMessages = [...historyMessages, { role: "user", content: userContent }];

      const response = await fetch(getApiUrl("/blog/admin/ai/chat"), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Access-Token": localStorage.getItem("access_token") || "",
        },
        body: JSON.stringify({ messages: apiMessages }),
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error("Request failed");
      }

      const reader = response.body?.getReader();
      if (!reader) throw new Error("No reader available");

      const decoder = new TextDecoder();
      let responseContextMessages: Array<Record<string, unknown>> = [];
      let isThinking = false;

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split("\n");

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const dataStr = line.slice(6);
            if (!dataStr.trim()) continue;

            try {
              const event = JSON.parse(dataStr) as StreamEvent;

              if (event.type === "done") {
                if (isThinking) {
                  isThinking = false;
                  accumulatedContent += "</think>";
                  setMessages((prev) =>
                    prev.map((msg) =>
                      msg.id === assistantMsg.id ? { ...msg, content: accumulatedContent } : msg,
                    ),
                  );
                }
                continue;
              }

              if (event.type === "tool_start") {
                if (isThinking) {
                  isThinking = false;
                  accumulatedContent += "</think>";
                }
                const toolData = event.data;
                const newToolCall: ToolCall = {
                  id: toolData.id,
                  name: toolData.name,
                  mcpName: toolData.mcpName,
                  arguments: toolData.arguments,
                  isLoading: true,
                };
                currentToolCalls = [...currentToolCalls, newToolCall];
                accumulatedContent += `\n<tool_call id="${toolData.id}"></tool_call>\n`;
                setMessages((prev) =>
                  prev.map((msg) =>
                    msg.id === assistantMsg.id
                      ? { ...msg, content: accumulatedContent, toolCalls: currentToolCalls }
                      : msg,
                  ),
                );
                continue;
              }

              if (event.type === "tool_end") {
                const toolData = event.data;
                currentToolCalls = currentToolCalls.map((tc) =>
                  tc.id === toolData.id ? { ...tc, result: toolData.result, isLoading: false } : tc,
                );
                setMessages((prev) =>
                  prev.map((msg) =>
                    msg.id === assistantMsg.id ? { ...msg, toolCalls: currentToolCalls } : msg,
                  ),
                );
                continue;
              }

              if (event.type === "context") {
                responseContextMessages = event.data;
                continue;
              }

              if (event.type === "error") {
                throw new Error(event.content || "Unknown error from server");
              }

              if (event.type === "content") {
                if (isThinking) {
                  isThinking = false;
                  accumulatedContent += "</think>";
                }
                const content = event.content || "";
                accumulatedContent += content;
                setMessages((prev) =>
                  prev.map((msg) =>
                    msg.id === assistantMsg.id ? { ...msg, content: accumulatedContent } : msg,
                  ),
                );
                continue;
              }

              if (event.type === "reasoning") {
                if (!isThinking) {
                  isThinking = true;
                  accumulatedContent += "<think>";
                }
                const reasoning = event.content || "";
                accumulatedContent += reasoning;
                setMessages((prev) =>
                  prev.map((msg) =>
                    msg.id === assistantMsg.id ? { ...msg, content: accumulatedContent } : msg,
                  ),
                );
                continue;
              }
            } catch (e) {
              if (e instanceof SyntaxError) {
                // Ignore JSON parse errors (e.g. from partial chunks)
                console.error("Failed to parse SSE event:", e);
                continue;
              }
              throw e;
            }
          }
        }
      }

      // Save final content to Dexie
      await updateAssistantMessage(assistantMsg.id!, accumulatedContent, currentToolCalls);
      if (responseContextMessages.length > 0) {
        await updateAssistantContextMessages(assistantMsg.id!, responseContextMessages);
      }

      // Mark streaming as complete
      setMessages((prev) =>
        prev.map((msg) =>
          msg.id === assistantMsg.id
            ? { ...msg, contextMessages: responseContextMessages, isStreaming: false }
            : msg,
        ),
      );
    } catch (err) {
      if ((err as Error).name !== "AbortError") {
        const errorMessage = err instanceof Error ? err.message : "请求失败，请稍后重试。";
        const errorContent = accumulatedContent
          ? accumulatedContent + `\n\n> **错误**: ${errorMessage}`
          : `抱歉，${errorMessage}`;

        await updateAssistantMessage(assistantMsg.id!, errorContent, currentToolCalls);
        setMessages((prev) =>
          prev.map((msg) =>
            msg.id === assistantMsg.id
              ? { ...msg, content: errorContent, isStreaming: false }
              : msg,
          ),
        );
      }
    } finally {
      setIsLoading(false);
      abortControllerRef.current = null;
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  }, [input, isLoading, messages]);

  // Handle Enter key
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const renderAssistantMessage = (message: DisplayMessage) => {
    if (!message.content && !message.toolCalls?.length) {
      return (
        <div className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded">
          <AgentChatIndicator size={"sm"} />
        </div>
      );
    }

    const content = message.content || "";
    // Regex to split by <tool_call id="..."></tool_call> or <think>...</think>
    const parts = content.split(
      /(<tool_call id="[^"]+"><\/tool_call>|<think>[\s\S]*?<\/think>|<think>[\s\S]*$)/,
    );

    // Keep track of which tool calls have been rendered inline
    const renderedToolCallIds = new Set<string>();

    const renderedParts = parts.map((part, index) => {
      const matchTool = part.match(/<tool_call id="([^"]+)"><\/tool_call>/);
      if (matchTool) {
        const id = matchTool[1];
        const toolCall = message.toolCalls?.find((tc) => tc.id === id);
        if (toolCall) {
          renderedToolCallIds.add(id);
          return <ToolCallCard key={`tc-${id}-${index}`} toolCall={toolCall} />;
        }
        return null;
      }

      const matchThink = part.match(/<think>([\s\S]*?)(?:<\/think>|$)/);
      if (matchThink) {
        const thinkContent = matchThink[1];
        const isFinished = part.endsWith("</think>");
        return (
          <ReasoningBlock
            key={`think-${index}`}
            content={thinkContent}
            isFinished={isFinished}
            isStreaming={message.isStreaming}
          />
        );
      }

      // Render text part with Viewer if it's not empty
      if (part.trim()) {
        return (
          <div
            key={`text-${index}`}
            className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded mb-2 last:mb-0"
          >
            <Viewer value={part} plugins={plugins} sanitize={sanitize} />
          </div>
        );
      }
      return null;
    });

    // Render any remaining tool calls at the top (for backwards compatibility)
    const unrenderedToolCalls =
      message.toolCalls?.filter((tc) => !renderedToolCallIds.has(tc.id)) || [];

    return (
      <>
        {unrenderedToolCalls.length > 0 && (
          <div className="mb-2">
            {unrenderedToolCalls.map((toolCall) => (
              <ToolCallCard key={toolCall.id} toolCall={toolCall} />
            ))}
          </div>
        )}
        {renderedParts}
      </>
    );
  };

  return (
    <>
      {/* Floating Button - Draggable */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            ref={buttonRef}
            onClick={() => {
              // Only toggle if not actually moved (prevent click after drag)
              if (!hasMovedRef.current) {
                setIsOpen(!isOpen);
              }
            }}
            onMouseDown={handleDragStart}
            onTouchStart={handleDragStart}
            style={{
              bottom: `${buttonPosition.bottom}px`,
              right: `${buttonPosition.right}px`,
            }}
            className={`fixed z-50 w-12 h-12 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-lg hover:shadow-xl transition-shadow duration-300 flex items-center justify-center ${
              isDragging ? "cursor-grabbing scale-110" : "cursor-grab hover:scale-110"
            }`}
            aria-label="AI Chat"
          >
            {isOpen ? <X className="w-5 h-5" /> : <MessageCircle className="w-5 h-5" />}
          </button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{isOpen ? "收起聊天" : "展开聊天（可拖动）"}</p>
        </TooltipContent>
      </Tooltip>

      {/* Chat Window */}
      {isOpen && (
        <div
          style={
            isFullscreen
              ? undefined
              : {
                  bottom: `${buttonPosition.bottom + 60}px`,
                  right: `${buttonPosition.right}px`,
                }
          }
          className={`fixed z-50 bg-white shadow-2xl border border-gray-200 flex flex-col overflow-hidden animate-in duration-300 ${
            isFullscreen
              ? "inset-0 rounded-none"
              : "w-[520px] h-[600px] rounded-xl slide-in-from-bottom-4"
          }`}
        >
          {/* Fullscreen wrapper for centering content */}
          <div
            className={`flex flex-col h-full ${
              isFullscreen ? "w-2/3 mx-auto border-x border-gray-200" : "w-full"
            }`}
          >
            {/* Header */}
            <div className="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-4 py-2 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
                  <MessageCircle className="w-4 h-4" />
                </div>
                <div>
                  <h3 className="font-semibold">AI 助手</h3>
                  <p className="text-xs text-white/80">随时为您解答问题</p>
                </div>
              </div>
              <div className="flex items-center gap-1">
                <button
                  onClick={() => setIsFullscreen(!isFullscreen)}
                  className="w-8 h-8 rounded-full hover:bg-white/20 flex items-center justify-center transition-colors"
                  title={isFullscreen ? "退出全屏" : "全屏"}
                >
                  {isFullscreen ? (
                    <Minimize2 className="w-4 h-4" />
                  ) : (
                    <Maximize2 className="w-4 h-4" />
                  )}
                </button>
                <button
                  onClick={() => {
                    setIsOpen(false);
                    setIsFullscreen(false);
                  }}
                  className="w-8 h-8 rounded-full hover:bg-white/20 flex items-center justify-center transition-colors"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>
            </div>

            {/* Messages Area */}
            <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
              {messages.length === 0 && (
                <div className="text-center text-gray-500 mt-20">
                  <MessageCircle className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                  <p>有什么我可以帮助您的吗？</p>
                </div>
              )}
              {messages.map((message) => (
                <div
                  key={message.id}
                  className={`flex ${message.role === "user" ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-[85%] rounded-xl px-3 py-2 relative group ${
                      message.role === "user"
                        ? "bg-gradient-to-r from-blue-500 to-purple-600 text-white"
                        : "bg-white border border-gray-200"
                    }`}
                  >
                    {message.role === "assistant" ? (
                      <div className="relative">
                        {renderAssistantMessage(message)}
                        {/* Copy button for assistant */}
                        {!message.isStreaming && message.content && (
                          <div className="absolute -bottom-11 -left-4 flex items-center gap-1">
                            <button
                              onClick={() => copyToClipboard(message.content, message.id)}
                              className="opacity-0 group-hover:opacity-100 transition-opacity rounded-full p-1.5 hover:bg-gray-200"
                              title="复制"
                            >
                              {copiedId === message.id ? (
                                <Check className="w-3.5 h-3.5 text-green-500" />
                              ) : (
                                <Copy className="w-3.5 h-3.5 text-gray-500" />
                              )}
                            </button>
                            <button
                              onClick={() => handleDeletePair(message.pairId)}
                              className="opacity-0 group-hover:opacity-100 transition-opacity rounded-full p-1.5 hover:bg-gray-200"
                              title="删除对话"
                            >
                              <Trash2 className="w-3.5 h-3.5 text-gray-500" />
                            </button>
                          </div>
                        )}
                      </div>
                    ) : (
                      <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                    )}
                    {/* Delete button for messages (shown on user messages, deletes pair) */}
                    {message.role === "user" && !isLoading && (
                      <button
                        onClick={() => handleDeletePair(message.pairId)}
                        className="absolute -top-0 -left-8 opacity-0 group-hover:opacity-100 transition-opacity rounded-full p-1.5 hover:bg-gray-200"
                        title="删除对话"
                      >
                        <Trash2 className="w-3.5 h-3.5 text-gray-500" />
                      </button>
                    )}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <div className="px-4 py-3 bg-white border-t border-gray-200">
              <div className="flex items-center gap-2">
                <input
                  ref={inputRef}
                  type="text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="输入您的问题..."
                  disabled={isLoading}
                  className="flex-1 px-4 py-2 bg-gray-100 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white transition-all disabled:opacity-50"
                />
                <Button
                  onClick={sendMessage}
                  disabled={!input.trim() || isLoading}
                  className="w-9 h-9 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 p-0 flex items-center justify-center"
                >
                  <Send className="w-5 h-5" />
                </Button>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      onClick={handleClearAll}
                      disabled={messages.length === 0 || isLoading}
                      variant="outline"
                      className="w-9 h-9 rounded-full p-0 flex items-center justify-center"
                    >
                      <RotateCcw className="w-4 h-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>清空聊天</p>
                  </TooltipContent>
                </Tooltip>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
