import axios from 'axios'
import {getApiUrl} from "./common";

export const getAccessToken = function () {
  let token = localStorage.getItem('access_token')
  if(!token){
    return ""
  }
  return token
}

export async function request(option) {
  let response
  try {
    response = await axios(option)
  } catch (e) {
    throw {code: 500, msg: e.message}
  }

  const ret = response.data
  if (ret && ret.code !== 200) {
    throw ret
  }
  return ret.data
}

export const createApi = (url, data) => {
  let headers = {
    'Access-Token': getAccessToken(),
  }

  let param = {
    url: getApiUrl(url),
    data,
    method: 'post',
    headers
  }

  return request(param)
}