import { getApiUrl } from "./common";
import { dialog } from "./dialog";

type RequestOptions = {
  url: string;
  method?: string;
  headers?: Record<string, string>;
  data?: any;
};

export async function request<T = any>(
  option: RequestOptions,
  errorHandler?: (e: any) => void
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
      (init.headers as Record<string, string>)["Content-Type"] =
        "application/json";
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

    if (
      isJSON &&
      payload &&
      typeof payload === "object" &&
      "code" in payload &&
      "msg" in payload
    ) {
      throw payload;
    }

    throw {
      code: resp.status,
      msg:
        typeof payload === "string"
          ? payload
          : resp.statusText || "Request error",
    };
  } catch (e: any) {
    const err = { code: e?.code || 500, msg: e?.msg || "Network error" };
    if (errorHandler) {
      errorHandler(err);
      // 返回一个 rejected promise 让调用方自行停止后续逻辑或继续链式处理
      throw err;
    } else {
      dialog.message(err.msg);
      throw err;
    }
  }
}

export const createApi = <TResp = any>(
  url: string,
  data?: any,
  errorHandler?: (e: any) => void
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
