export const getApiUrl = (url: string) => {
  return import.meta.env.PROD ? "https://api.fifsky.com" + url : url;
};

export const getAccessToken = () => {
  const token = localStorage.getItem("access_token");
  return token || "";
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
