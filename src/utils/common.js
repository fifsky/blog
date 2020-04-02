export const getApiUrl = (url) => {
  return process.env.NODE_ENV === 'production' ? 'https://api.fifsky.com' + url : url
}