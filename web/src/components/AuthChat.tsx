import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { aiLoginChatApi } from "@/service/aiLogin";

interface Message {
  role: "user" | "assistant";
  content: string;
}

interface AuthChatProps {
  sessionId: string;
  welcomeMessage: string;
  onSuccess: (token: string) => void;
  onFailed: () => void;
}

export function AuthChat({ sessionId, welcomeMessage, onSuccess, onFailed }: AuthChatProps) {
  const [messages, setMessages] = useState<Message[]>([
    { role: "assistant", content: welcomeMessage },
  ]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [score, setScore] = useState(0);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim() || loading) return;

    const userMessage = input.trim();
    setInput("");
    setMessages((prev) => [...prev, { role: "user", content: userMessage }]);
    setLoading(true);

    try {
      const response = await aiLoginChatApi({
        session_id: sessionId,
        message: userMessage,
      });

      setMessages((prev) => [...prev, { role: "assistant", content: response.content }]);
      setScore(response.score);

      if (response.verified && response.access_token) {
        setTimeout(() => onSuccess(response.access_token!), 1000);
      } else if (response.failed) {
        setTimeout(onFailed, 1000);
      }
    } catch {
      setMessages((prev) => [...prev, { role: "assistant", content: "发生错误，请重试。" }]);
    } finally {
      setLoading(false);
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  };

  return (
    <div className="flex flex-col h-[400px]">
      <div className="mb-2">
        <div className="text-sm text-gray-500 mb-1">验证进度：已答对 {Math.round(score * 3)} / 3 条</div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className="bg-blue-500 h-2 rounded-full transition-all duration-300"
            style={{ width: `${score * 100}%` }}
          />
        </div>
      </div>

      <div className="flex-1 overflow-y-auto border rounded p-4 mb-4 space-y-3">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`p-3 rounded-lg whitespace-pre-wrap ${
              msg.role === "user" ? "bg-blue-100 ml-12" : "bg-gray-100 mr-12"
            }`}
          >
            {msg.content}
          </div>
        ))}
        {loading && <div className="bg-gray-100 mr-12 p-3 rounded-lg">正在思考...</div>}
        <div ref={messagesEndRef} />
      </div>

      <div className="flex gap-2">
        <Input
          ref={inputRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
          placeholder="输入消息..."
          disabled={loading}
        />
        <Button size="sm" onClick={handleSend} loading={loading}>
          发送
        </Button>
      </div>
    </div>
  );
}
