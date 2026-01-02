import React, { createContext, useContext, useReducer } from "react";
import { loginApi, loginUserApi } from "@/service";
import { sync } from "@/utils/sync";
import { LoginRequest } from "@/types/openapi";

export type UserInfo = Partial<{
  id: number;
  name: string;
  nickName: string;
  email: string;
  status: number;
  type: number;
}>;

type State = {
  userInfo: UserInfo;
  keyword: string;
};

const initialState: State = { userInfo: {}, keyword: "" };

type Action =
  | { type: "setUserInfo"; payload: UserInfo }
  | { type: "setKeyword"; payload: string };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case "setUserInfo":
      return { ...state, userInfo: action.payload };
    case "setKeyword":
      return { ...state, keyword: action.payload };
    default:
      return state;
  }
}

type StoreCtx = {
  state: State;
  dispatch: React.Dispatch<Action>;
  loginAction: (data: LoginRequest) => Promise<void>;
  currentUserAction: () => Promise<void>;
};

const Ctx = createContext<StoreCtx | null>(null);

export function StoreProvider({ children }: { children: React.ReactNode }) {
  const [state, dispatch] = useReducer(reducer, initialState);

  const loginAction = async (data: LoginRequest) =>
    sync(async () => {
      const ret = await loginApi(data);
      if (!ret.access_token) throw "登录失败";
      localStorage.setItem("access_token", ret.access_token);
      dispatch({ type: "setUserInfo", payload: ret.user });
    });

  const currentUserAction = async () =>
    sync(async () => {
      try {
        const ret = await loginUserApi();
        dispatch({ type: "setUserInfo", payload: ret });
      } catch {
        localStorage.removeItem("access_token");
      }
    });

  return (
    <Ctx.Provider value={{ state, dispatch, loginAction, currentUserAction }}>
      {children}
    </Ctx.Provider>
  );
}

export function useStore() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error("Store not initialized");
  return ctx;
}
