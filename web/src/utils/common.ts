export const getApiUrl = (url: string) => {
  return import.meta.env.PROD ? "https://api.fifsky.com" + url : url;
};

export const getAccessToken = () => {
  const token = localStorage.getItem("access_token");
  return token || "";
};
