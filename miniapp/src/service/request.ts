import Taro from "@tarojs/taro";
import { getApiBaseUrl } from "../config/env";

export type RequestOptions = {
  url: string;
  method?: "GET" | "POST";
  headers?: Record<string, string>;
  data?: any;
};

export async function request<T = any>(option: RequestOptions): Promise<T> {
  const { url, method = "GET", headers = {}, data } = option;
  const token = Taro.getStorageSync("access_token");

  const resp = await Taro.request<T>({
    url,
    method,
    data,
    header: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "Access-Token": token || "",
      ...headers,
    },
  });

  if (resp.statusCode >= 200 && resp.statusCode < 300) {
    return resp.data as T;
  }

  const payload: any = resp.data;
  if (
    payload &&
    typeof payload === "object" &&
    typeof payload.code === "string" &&
    typeof payload.message === "string"
  ) {
    throw new Error(payload.message);
  }

  throw new Error(String(resp.statusCode));
}

export const createApi = async <TResp = any>(url: string, data?: any): Promise<TResp> => {
  const base = getApiBaseUrl();
  return request<TResp>({
    url: base + url,
    method: "POST",
    data,
  });
};
