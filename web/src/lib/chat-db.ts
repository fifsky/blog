// Dexie.js database for AI chat message persistence
import Dexie, { type Table } from "dexie";

export interface ChatMessage {
  id?: number; // auto-increment primary key
  pairId: string; // links user message with assistant response
  role: "user" | "assistant";
  content: string;
  createdAt: Date;
}

class ChatDatabase extends Dexie {
  messages!: Table<ChatMessage>;

  constructor() {
    super("AIChatDB");
    this.version(1).stores({
      messages: "++id, pairId, role, createdAt",
    });
  }
}

export const chatDB = new ChatDatabase();

// Helper functions for chat operations

/**
 * Get all messages ordered by creation time
 */
export async function getAllMessages(): Promise<ChatMessage[]> {
  return chatDB.messages.orderBy("createdAt").toArray();
}

/**
 * Add a message pair (user question + assistant response placeholder)
 */
export async function addMessagePair(
  pairId: string,
  userContent: string,
): Promise<{ userMsg: ChatMessage; assistantMsg: ChatMessage }> {
  const now = new Date();

  const userMsg: ChatMessage = {
    pairId,
    role: "user",
    content: userContent,
    createdAt: now,
  };

  const assistantMsg: ChatMessage = {
    pairId,
    role: "assistant",
    content: "",
    createdAt: new Date(now.getTime() + 1), // slightly later
  };

  await chatDB.messages.bulkAdd([userMsg, assistantMsg]);

  // Re-fetch to get the IDs
  const messages = await chatDB.messages.where("pairId").equals(pairId).toArray();
  return {
    userMsg: messages.find((m) => m.role === "user")!,
    assistantMsg: messages.find((m) => m.role === "assistant")!,
  };
}

/**
 * Update assistant message content
 */
export async function updateAssistantMessage(id: number, content: string): Promise<void> {
  await chatDB.messages.update(id, { content });
}

/**
 * Delete a message pair by pairId
 */
export async function deleteMessagePair(pairId: string): Promise<void> {
  await chatDB.messages.where("pairId").equals(pairId).delete();
}

/**
 * Clear all messages
 */
export async function clearAllMessages(): Promise<void> {
  await chatDB.messages.clear();
}

/**
 * Convert stored messages to API format
 */
export function messagesToApiFormat(
  messages: ChatMessage[],
): Array<{ role: string; content: string }> {
  return messages
    .filter((m) => m.content.trim() !== "") // exclude empty messages
    .map((m) => ({
      role: m.role,
      content: m.content,
    }));
}
