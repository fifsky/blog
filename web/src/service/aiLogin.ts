import { createApi } from "@/utils/request";

export interface AILoginInitRequest {
  user_name: string;
}

export interface AILoginInitResponse {
  session_id: string;
  welcome_message: string;
}

export interface AILoginChatRequest {
  session_id: string;
  message: string;
}

export interface AILoginChatResponse {
  content: string;
  score: number;
  verified: boolean;
  access_token?: string;
  failed?: boolean;
  error_message?: string;
}

export interface ProfileSetupRequest {
  identity_description: string;
  verification_threshold?: number;
  max_attempts?: number;
}

export interface ProfileResponse {
  id: number;
  identity_description: string;
  verification_threshold: number;
  max_attempts: number;
  created_at: string;
  updated_at: string;
}

export const aiLoginInitApi = (data: AILoginInitRequest) =>
  createApi<AILoginInitResponse>("/blog/ai-login/init", data);

export const aiLoginChatApi = (data: AILoginChatRequest) =>
  createApi<AILoginChatResponse>("/blog/ai-login/chat", data);

export const aiAuthSetupProfileApi = (data: ProfileSetupRequest) =>
  createApi<{ id: number }>("/blog/admin/ai-auth/profile", data);

export const aiAuthGetProfileApi = () =>
  createApi<ProfileResponse>("/blog/admin/ai-auth/profile/get", {});
