import { create } from "zustand";
import { loginApi, loginUserApi } from "@/service";
import { LoginRequest } from "@/types/openapi";
import { clearAuth, isTokenExpired } from "@/utils/common";

export type UserInfo = Partial<{
  id: number;
  name: string;
  nickName: string;
  email: string;
  status: string;
  type: number;
}>;

export type Store = {
  userInfo: UserInfo;
  keyword: string;
  setUserInfo: (u: UserInfo) => void;
  setKeyword: (k: string) => void;
  loginAction: (data: LoginRequest) => Promise<any>;
  currentUserAction: () => Promise<void>;
};

export const useStore = create<Store>((set) => ({
  userInfo: {},
  keyword: "",
  setUserInfo: (u) => set({ userInfo: u }),
  setKeyword: (k) => set({ keyword: k }),
  loginAction: async (data: LoginRequest) => {
    const ret = await loginApi(data);
    if (ret.require_totp) {
      return ret;
    }
    if (!ret.access_token) throw "登录失败";
    localStorage.setItem("access_token", ret.access_token);
    if (ret.expires_at) {
      localStorage.setItem("expires_at", String(ret.expires_at));
    }
    set({ userInfo: ret.user });
    return ret;
  },
  currentUserAction: async () => {
    // token 已过期则静默清除，不发请求
    if (isTokenExpired()) {
      clearAuth();
      return;
    }
    try {
      const ret = await loginUserApi();
      set({ userInfo: ret });
    } catch {
      clearAuth();
    }
  },
}));

// 监听认证登出事件，统一清除 userInfo
if (typeof window !== "undefined") {
  window.addEventListener("auth:logout", () => {
    useStore.setState({ userInfo: {} });
  });
}
