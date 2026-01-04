import { create } from "zustand";
import { loginApi, loginUserApi } from "@/service";
import { LoginRequest } from "@/types/openapi";

export type UserInfo = Partial<{
  id: number;
  name: string;
  nickName: string;
  email: string;
  status: number;
  type: number;
}>;

export type Store = {
  userInfo: UserInfo;
  keyword: string;
  setUserInfo: (u: UserInfo) => void;
  setKeyword: (k: string) => void;
  loginAction: (data: LoginRequest) => Promise<void>;
  currentUserAction: () => Promise<void>;
};

export const useStore = create<Store>((set) => ({
  userInfo: {},
  keyword: "",
  setUserInfo: (u) => set({ userInfo: u }),
  setKeyword: (k) => set({ keyword: k }),
  loginAction: async (data: LoginRequest) => {
    const ret = await loginApi(data);
    if (!ret.access_token) throw "登录失败";
    localStorage.setItem("access_token", ret.access_token);
    set({ userInfo: ret.user });
  },
  currentUserAction: async () => {
    try {
      const ret = await loginUserApi();
      set({ userInfo: ret });
    } catch {
      localStorage.removeItem("access_token");
    }
  },
}));
