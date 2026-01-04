import { useEffect, useRef, useState } from "react";
export function AppMessage() {
  // 是否处于可见状态（用于控制淡入淡出）
  const [visible, setVisible] = useState(false);
  // 是否挂载到页面（淡出结束后再卸载）
  const [mounted, setMounted] = useState(false);
  const [msg, setMsg] = useState("");
  const timerRef = useRef<number | undefined>(undefined);
  const hideRef = useRef<number | undefined>(undefined);
  const DURATION = 300; // 动画时长（ms）
  const SHOW_MS = 3000; // 展示时长（ms）
  useEffect(() => {
    const handler = (e: Event) => {
      const d = (e as CustomEvent).detail;
      setMsg(d?.msg || String(d));
      // 挂载并淡入
      setMounted(true);
      requestAnimationFrame(() => setVisible(true));
      // 清理上一次的定时器
      if (timerRef.current) clearTimeout(timerRef.current);
      if (hideRef.current) clearTimeout(hideRef.current);
      // 展示一段时间后淡出
      timerRef.current = window.setTimeout(() => {
        setVisible(false);
        // 等淡出动画结束后卸载
        hideRef.current = window.setTimeout(() => setMounted(false), DURATION);
      }, SHOW_MS);
    };
    window.addEventListener("app-alert", handler as EventListener);
    return () =>
      window.removeEventListener("app-alert", handler as EventListener);
  }, []);
  if (!mounted) return null;
  return (
    <div
      className="fixed top-4 left-1/2 -translate-x-1/2 z-50 transition-opacity duration-300 ease-out"
      style={{ opacity: visible ? 1 : 0 }}
    >
      <div className="px-4 py-2 rounded bg-[#ffffe0] border border-[#e6db55] text-[#333] shadow">
        <i className="iconfont icon-info text-[#e6db55] mr-2"></i>
        {msg}
      </div>
    </div>
  );
}
