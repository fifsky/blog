import { getApiUrl } from "./common";
import { dialog } from "./dialog";
import { AppError } from "./error";

type RequestOptions = {
  url: string;
  method?: string;
  headers?: Record<string, string>;
  data?: any;
};

export async function request<T = any>(
  option: RequestOptions,
  errorHandler?: (e: AppError) => void,
) {
  const { url, method = "GET", headers = {}, data } = option;

  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 60000); // 60s 超时

  const init: RequestInit = {
    method,
    headers: {
      Accept: "application/json",
      ...headers,
    },
    signal: controller.signal,
  };

  if (data !== undefined && method.toUpperCase() !== "GET") {
    if (data instanceof FormData) {
      init.body = data;
      // 让浏览器自动设置 multipart/form-data 边界
    } else {
      init.body = JSON.stringify(data);
      (init.headers as Record<string, string>)["Content-Type"] = "application/json";
    }
  }
  try {
    const resp = await fetch(url, init);
    clearTimeout(timeoutId);

    const contentType = resp.headers.get("content-type") || "";
    const isJSON = contentType.includes("application/json");
    const payload = isJSON ? await resp.json() : await resp.text();

    if (resp.ok) {
      return payload as T;
    }

    if (isJSON && payload && typeof payload === "object") {
      const p: any = payload;
      if (typeof p.code === "string" && typeof p.message === "string") {
        throw new AppError(p.code, p.message, p.details);
      }
    }
    throw new AppError(String(resp.status), getErrorMessage(payload));
  } catch (e: any) {
    let err: AppError;
    if (e instanceof AppError) {
      err = e;
    } else if (e instanceof Error && e.name === "AbortError") {
      err = new AppError("REQUEST_TIMEOUT", "请求超时，请稍后重试");
    } else {
      err = new AppError(getErrorCode(e.code), getErrorMessage(e));
    }
    if (errorHandler) {
      errorHandler(err);
      // 返回一个 rejected promise 让调用方自行停止后续逻辑或继续链式处理
      throw err;
    } else {
      if (err.code !== "UNAUTHORIZED") {
        dialog.message(err.message);
      }
      throw err;
    }
  }
}

function getErrorCode(code: unknown): string {
  if (typeof code === "string") {
    return code;
  }
  if (typeof code === "number") {
    return String(code);
  }
  return "UNKNOWN";
}

function getErrorMessage(error: unknown, fallback = "Unknown error"): string {
  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === "string") {
    return error;
  }

  if (error && typeof error === "object" && "message" in error) {
    return String(error.message);
  }

  if (error && typeof error === "object" && "msg" in error) {
    return String(error.msg);
  }

  return fallback;
}

export const createApi = <TResp = any>(
  url: string,
  data?: any,
  errorHandler?: (e: AppError) => void,
) => {
  const headers = {
    "Access-Token": localStorage.getItem("access_token") || "",
  };
  const param: RequestOptions = {
    url: getApiUrl(url),
    data,
    method: "post",
    headers,
  };
  return request<TResp>(param, errorHandler);
};
