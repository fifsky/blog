import { RouterProvider } from "react-router-dom";
import { router } from "./router";
import { StoreProvider } from "./store/context";
// @ts-ignore
import NProgress from "nprogress";
import "nprogress/nprogress.css"; // 必须导入样式
import { useEffect } from "react";

export default function App() {
  // 2. 监听 router 导航变化
  useEffect(() => {
    // 订阅 router 状态变化
    const unsubscribe = router.subscribe((state) => {
      // state.navigation.state：导航状态（loading/idle
      console.log(state);
      if (state.navigation.state === "loading") {
        NProgress.start(); // 导航开始，启动进度条
      } else {
        NProgress.done(); // 导航结束，关闭进度条
      }
    });

    // 组件卸载时取消订阅
    return () => {
      unsubscribe();
      NProgress.done(); // 确保组件卸载时进度条关闭
    };
  }, []); // 仅挂载时订阅一次

  return (
    <StoreProvider>
      <RouterProvider router={router}></RouterProvider>
    </StoreProvider>
  );
}
