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
  Sparkles,
  ChevronDown,
  ChevronRight,
} from "lucide-react";
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
  deleteMessagePair,
  clearAllMessages,
} from "@/lib/chat-db";
import { ToolCallCard } from "./ToolCallCard";
import type { ToolCall, DisplayMessage, MessageBlock } from "./types";
import { AgentChatIndicator } from "@/components/agents-ui/agent-chat-indicator";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { aiListSkillsApi } from "@/service";

// ByteMD plugins for rendering
const plugins = [gfm(), breaks(), highlightPlugin()];

const STORAGE_KEY = "ai-chat-button-position";
const DEFAULT_POSITION = { bottom: 24, right: 80 };

interface ButtonPosition {
  bottom: number;
  right: number;
}

export function AIChat() {
  const [isOpen, setIsOpen] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [copiedId, setCopiedId] = useState<number | null>(null);

  // Track which thinking blocks are expanded
  const [expandedThinking, setExpandedThinking] = useState<Record<number, boolean>>({});

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const [skills, setSkills] = useState<{ name: string; description: string }[]>([]);
  const [showSkillPopover, setShowSkillPopover] = useState(false);
  const [skillSearchQuery, setSkillSearchQuery] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  // Fetch skills on mount
  useEffect(() => {
    aiListSkillsApi()
      .then((res: any) => {
        if (res?.skills) {
          setSkills(res.skills);
        }
      })
      .catch((err: any) => {
        console.error("Failed to fetch skills:", err);
      });
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setInput(val);

    // Check for skill trigger '/' at the end of the input or after a space
    if (val.endsWith("/")) {
      setShowSkillPopover(true);
      setSkillSearchQuery("");
    } else if (showSkillPopover) {
      const lastSlashIndex = val.lastIndexOf("/");
      if (lastSlashIndex !== -1) {
        setSkillSearchQuery(val.slice(lastSlashIndex + 1));
      } else {
        setShowSkillPopover(false);
      }
    }
  };

  const handleSkillSelect = (skillName: string) => {
    const lastSlashIndex = input.lastIndexOf("/");
    if (lastSlashIndex !== -1) {
      const newInput = input.slice(0, lastSlashIndex) + skillName + " ";
      setInput(newInput);
    } else {
      setInput(input + skillName + " ");
    }
    setShowSkillPopover(false);
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  // Handle Enter key
  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (showSkillPopover) {
      if (e.key === "Escape") {
        setShowSkillPopover(false);
        e.preventDefault();
        return;
      }
    }

    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (!showSkillPopover) {
        sendMessage();
      }
    }
  };

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
      dbMessages.map((m) => ({
        id: m.id!,
        pairId: m.pairId,
        role: m.role,
        content: m.content,
        isStreaming: false,
      })),
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
      blocks: [],
    };

    setMessages((prev) => [...prev, newUserMessage, newAssistantMessage]);
    setInput("");
    setIsLoading(true);

    try {
      const controller = new AbortController();
      abortControllerRef.current = controller;

      // Build messages array for API (history + current)
      const historyMessages = messages
        .filter((m) => m.content.trim() !== "")
        .map((m) => ({ role: m.role, content: m.content }));

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
      let accumulatedContent = "";
      let accumulatedThinking = "";
      let isThinkingActive = false;
      let thinkingDuration = "";
      let currentToolCalls: ToolCall[] = [];
      let currentBlocks: MessageBlock[] = [];
      // Track the current text block content
      let currentTextBlock = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split("\n");

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const data = line.slice(6);

            if (data === "[DONE]") continue;

            // Handle thinking event
            if (data.startsWith("[THINKING] ")) {
              const thinkData = JSON.parse(data.slice(11));
              if (thinkData.content) {
                const c = thinkData.content.replace(/\\n/g, "\n");
                accumulatedThinking += c;
              }
              isThinkingActive = thinkData.thinking;
              if (thinkData.duration) {
                thinkingDuration = thinkData.duration;
              }
              setMessages((prev) =>
                prev.map((msg) =>
                  msg.id === assistantMsg.id
                    ? {
                        ...msg,
                        thinking: {
                          content: accumulatedThinking,
                          isThinking: isThinkingActive,
                          duration: thinkingDuration,
                        },
                      }
                    : msg,
                ),
              );
              continue;
            }

            // Handle tool start event
            if (data.startsWith("[TOOL_START] ")) {
              const toolData = JSON.parse(data.slice(13));
              const newToolCall: ToolCall = {
                id: toolData.id,
                name: toolData.name,
                mcpName: toolData.mcpName,
                arguments: toolData.arguments,
                isLoading: true,
              };
              currentToolCalls = [...currentToolCalls, newToolCall];

              // If there's accumulated text, push it as a block before the tool call
              if (currentTextBlock) {
                currentBlocks = [...currentBlocks, { type: "text", content: currentTextBlock }];
                currentTextBlock = ""; // Reset for next text segment
              }
              currentBlocks = [...currentBlocks, { type: "tool", toolCall: newToolCall }];

              setMessages((prev) =>
                prev.map((msg) =>
                  msg.id === assistantMsg.id
                    ? { ...msg, toolCalls: currentToolCalls, blocks: currentBlocks }
                    : msg,
                ),
              );
              continue;
            }

            // Handle tool end event
            if (data.startsWith("[TOOL_END] ")) {
              const toolData = JSON.parse(data.slice(11));
              currentToolCalls = currentToolCalls.map((tc) =>
                tc.id === toolData.id ? { ...tc, result: toolData.result, isLoading: false } : tc,
              );

              // Update the block as well
              currentBlocks = currentBlocks.map((block) =>
                block.type === "tool" && block.toolCall.id === toolData.id
                  ? {
                      type: "tool",
                      toolCall: { ...block.toolCall, result: toolData.result, isLoading: false },
                    }
                  : block,
              );

              setMessages((prev) =>
                prev.map((msg) =>
                  msg.id === assistantMsg.id
                    ? { ...msg, toolCalls: currentToolCalls, blocks: currentBlocks }
                    : msg,
                ),
              );
              continue;
            }

            // Handle regular content
            const content = data.replace(/\\n/g, "\n");
            accumulatedContent += content;
            currentTextBlock += content;

            // We update the last text block if it exists and is at the end, otherwise we don't add it yet
            // Wait, to keep it simple and real-time:
            let displayBlocks = [...currentBlocks];
            if (currentTextBlock) {
              displayBlocks.push({ type: "text", content: currentTextBlock });
            }

            setMessages((prev) =>
              prev.map((msg) =>
                msg.id === assistantMsg.id
                  ? { ...msg, content: accumulatedContent, blocks: displayBlocks }
                  : msg,
              ),
            );
          }
        }
      }

      // Save final content to Dexie
      await updateAssistantMessage(assistantMsg.id!, accumulatedContent);

      // Mark streaming as complete
      setMessages((prev) =>
        prev.map((msg) => (msg.id === assistantMsg.id ? { ...msg, isStreaming: false } : msg)),
      );
    } catch (err) {
      if ((err as Error).name !== "AbortError") {
        const errorContent = "抱歉，请求失败，请稍后重试。";
        await updateAssistantMessage(assistantMsg.id!, errorContent);
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

  const toggleThinking = (id: number) => {
    setExpandedThinking((prev) => ({
      ...prev,
      [id]: !prev[id],
    }));
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
                        {/* Thinking Process UI */}
                        {message.thinking && message.thinking.content && (
                          <div className="mb-3">
                            <button
                              onClick={() => toggleThinking(message.id)}
                              className="flex items-center gap-1.5 text-xs text-gray-500 hover:text-gray-700 transition-colors font-medium select-none"
                            >
                              {expandedThinking[message.id] || message.thinking.isThinking ? (
                                <ChevronDown className="w-3.5 h-3.5" />
                              ) : (
                                <ChevronRight className="w-3.5 h-3.5" />
                              )}
                              <Sparkles className="w-3.5 h-3.5" />
                              <span>
                                {message.thinking.isThinking
                                  ? "思考中..."
                                  : message.thinking.duration
                                    ? `思考完成 (${message.thinking.duration}s)`
                                    : "思考完成"}
                              </span>
                            </button>

                            {(expandedThinking[message.id] || message.thinking.isThinking) && (
                              <div className="mt-1.5 pl-3 border-l-2 border-gray-200">
                                <div className="text-xs text-gray-500 whitespace-pre-wrap font-mono">
                                  {message.thinking.content}
                                </div>
                              </div>
                            )}
                          </div>
                        )}

                        {/* Render interleaved blocks if they exist */}
                        {message.blocks && message.blocks.length > 0 ? (
                          <div className="flex flex-col gap-2">
                            {message.blocks.map((block, idx) => {
                              if (block.type === "tool") {
                                return (
                                  <ToolCallCard key={block.toolCall.id} toolCall={block.toolCall} />
                                );
                              } else {
                                return (
                                  <div
                                    key={`text-${idx}`}
                                    className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded"
                                  >
                                    <Viewer value={block.content || ""} plugins={plugins} />
                                  </div>
                                );
                              }
                            })}
                            {!message.content &&
                              !message.toolCalls?.length &&
                              message.isStreaming && (
                                <div className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded">
                                  <AgentChatIndicator size={"sm"} />
                                </div>
                              )}
                          </div>
                        ) : (
                          // Fallback to legacy rendering
                          <>
                            {/* Tool Calls UI */}
                            {message.toolCalls && message.toolCalls.length > 0 && (
                              <div className="mb-2">
                                {message.toolCalls.map((toolCall) => (
                                  <ToolCallCard key={toolCall.id} toolCall={toolCall} />
                                ))}
                              </div>
                            )}
                            <div className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded">
                              {!message.content && !message.toolCalls?.length ? (
                                <AgentChatIndicator size={"sm"} />
                              ) : message.content ? (
                                <Viewer value={message.content || ""} plugins={plugins} />
                              ) : null}
                            </div>
                          </>
                        )}

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
                <Popover open={showSkillPopover} onOpenChange={setShowSkillPopover}>
                  <PopoverTrigger asChild>
                    <div className="flex-1 relative">
                      <input
                        ref={inputRef}
                        type="text"
                        value={input}
                        onChange={handleInputChange}
                        onKeyDown={handleInputKeyDown}
                        placeholder="输入您的问题..."
                        disabled={isLoading}
                        className="w-full px-4 py-2 bg-gray-100 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white transition-all disabled:opacity-50"
                      />
                    </div>
                  </PopoverTrigger>
                  <PopoverContent
                    className="w-64 p-0"
                    align="start"
                    onOpenAutoFocus={(e) => e.preventDefault()}
                  >
                    <Command>
                      <CommandList>
                        <CommandEmpty>No skill found.</CommandEmpty>
                        <CommandGroup heading="Skills">
                          {skills
                            .filter((skill) =>
                              skill.name.toLowerCase().includes(skillSearchQuery.toLowerCase()),
                            )
                            .map((skill) => (
                              <CommandItem
                                key={skill.name}
                                onSelect={() => handleSkillSelect(skill.name)}
                                className="flex flex-col items-start px-2 py-1 cursor-pointer"
                              >
                                <div className="font-medium text-sm">{skill.name}</div>
                                <div className="text-xs text-gray-500 line-clamp-1">
                                  {skill.description}
                                </div>
                              </CommandItem>
                            ))}
                        </CommandGroup>
                      </CommandList>
                    </Command>
                  </PopoverContent>
                </Popover>
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
