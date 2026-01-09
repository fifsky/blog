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
  const init: RequestInit = {
    method,
    headers: {
      Accept: "application/json",
      ...headers,
    },
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
    const contentType = resp.headers.get("content-type") || "";
    const isJSON = contentType.includes("application/json");
    const payload = isJSON ? await resp.json() : await resp.text();

    if (resp.ok) {
      return payload as T;
    }

    if (isJSON && payload && typeof payload === "object") {
      const p: any = payload;
      if (typeof p.code === "number" && typeof p.msg === "string") {
        throw new AppError(p.code, p.msg);
      }
    }
    throw new AppError(resp.status, getErrorMessage(payload));
  } catch (e: any) {
    const err: AppError =
      e instanceof AppError ? e : new AppError(getErrorCode(e.code), getErrorMessage(e));
    if (errorHandler) {
      errorHandler(err);
      // 返回一个 rejected promise 让调用方自行停止后续逻辑或继续链式处理
      throw err;
    } else {
      dialog.message(err.message);
      throw err;
    }
  }
}

function getErrorCode(code: unknown) {
  if (typeof code === "number") {
    return code;
  }
  return 500;
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
