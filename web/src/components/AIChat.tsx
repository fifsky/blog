import { useState, useRef, useEffect, useCallback } from "react";
import { MessageCircle, X, Send, Copy, Check, Trash2, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Viewer } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import { highlightPlugin } from "@/lib/highlight-plugin";
import { getApiUrl } from "@/utils/common";
import { Spinner } from "@/components/ui/spinner";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import {
  getAllMessages,
  addMessagePair,
  updateAssistantMessage,
  deleteMessagePair,
  clearAllMessages,
} from "@/lib/chat-db";

// ByteMD plugins for rendering
const plugins = [gfm(), highlightPlugin()];

interface DisplayMessage {
  id: number;
  pairId: string;
  role: "user" | "assistant";
  content: string;
  isStreaming?: boolean;
}

export function AIChat() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [copiedId, setCopiedId] = useState<number | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

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

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split("\n");

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const data = line.slice(6);
            if (data === "[DONE]") continue;
            const content = data.replace(/\\n/g, "\n");
            accumulatedContent += content;

            setMessages((prev) =>
              prev.map((msg) =>
                msg.id === assistantMsg.id ? { ...msg, content: accumulatedContent } : msg,
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

  // Handle Enter key
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <>
      {/* Floating Button */}
      <Tooltip>
        <TooltipTrigger asChild>
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="fixed bottom-6 right-20 z-50 w-12 h-12 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-110 flex items-center justify-center"
            aria-label="AI Chat"
          >
            {isOpen ? <X className="w-5 h-5" /> : <MessageCircle className="w-5 h-5" />}
          </button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{isOpen ? "收起聊天" : "展开聊天"}</p>
        </TooltipContent>
      </Tooltip>

      {/* Chat Window */}
      {isOpen && (
        <div className="fixed bottom-24 right-20 z-50 w-[520px] h-[600px] bg-white rounded-xl shadow-2xl border border-gray-200 flex flex-col overflow-hidden animate-in slide-in-from-bottom-4 duration-300">
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
            <button
              onClick={() => setIsOpen(false)}
              className="w-8 h-8 rounded-full hover:bg-white/20 flex items-center justify-center transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
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
                      <div className="markdown-body text-sm prose prose-sm max-w-none [&_pre]:bg-gray-100 [&_pre]:p-2 [&_pre]:rounded">
                        {!message.content ? (
                          <Spinner />
                        ) : (
                          <Viewer value={message.content || ""} plugins={plugins} />
                        )}
                      </div>
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
      )}
    </>
  );
}
