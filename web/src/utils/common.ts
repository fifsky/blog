export const getApiUrl = (url: string) => {
  return import.meta.env.PROD ? "https://api.fifsky.com" + url : url;
};

export const getAccessToken = () => {
  const token = localStorage.getItem("access_token");
  return token || "";
};

// clearAuth 清除本地认证信息并派发登出事件，监听方据此更新 UI 状态
export const clearAuth = () => {
  localStorage.removeItem("access_token");
  localStorage.removeItem("expires_at");
  window.dispatchEvent(new Event("auth:logout"));
};

// getTokenExpiry 获取 token 过期时间戳（Unix 秒）
export const getTokenExpiry = (): number => {
  const expiresAt = localStorage.getItem("expires_at");
  return expiresAt ? parseInt(expiresAt, 10) : 0;
};

// TOKEN_EXPIRY_MARGIN token 过期提前量（秒），提前 10 分钟视为过期以便前端主动登出
const TOKEN_EXPIRY_MARGIN = 600;

// isTokenExpired 判断 token 是否已过期或不存在（提前 10 分钟视为过期）
export const isTokenExpired = (): boolean => {
  const expiresAt = getTokenExpiry();
  if (!expiresAt) return true;
  return Date.now() / 1000 >= expiresAt - TOKEN_EXPIRY_MARGIN;
};

// 定义 sleep 函数：接收毫秒数，返回一个 Promise
export const sleep = (ms: number) => {
  // 校验参数，确保传入的是合法数字
  if (typeof ms !== "number" || ms < 0) {
    throw new Error("sleep 函数需要传入一个非负数字作为毫秒数");
  }

  return new Promise((resolve) => {
    // 使用 setTimeout 延迟执行 resolve，实现暂停效果
    setTimeout(resolve, ms);
  });
};
