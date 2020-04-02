import { loginApi,loginUserApi } from '../service'
import {sync} from "../utils";

export default {
  loginAction({commit},data){
    return sync(async () => {
      let ret = await loginApi(data)

      if(!ret.access_token) {
        throw "登录失败"
      }

      localStorage.setItem('access_token', ret.access_token)
      commit("setUserInfo",ret.user)
    })
  },
  currentUserAction({commit}){
    return sync(async () => {
      let ret = await loginUserApi()
      commit("setUserInfo",ret)
    })
  },
}
