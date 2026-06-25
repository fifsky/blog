import { useCallback, useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router";
import { useStore } from "@/store/context";
import { clearAuth } from "@/utils/common";

const SCROLL_THRESHOLD = 150;
const HEADER_HEIGHT = 80;

export function CHeader() {
  const userInfo = useStore((s) => s.userInfo);
  const isLogin = !!userInfo.id;
  const navigate = useNavigate();
  const headerRef = useRef<HTMLDivElement>(null);
  const [progress, setProgress] = useState(0);
  const rafRef = useRef<number | null>(null);

  const handleScroll = useCallback(() => {
    if (rafRef.current !== null) cancelAnimationFrame(rafRef.current);
    rafRef.current = requestAnimationFrame(() => {
      const p = Math.min(window.scrollY / SCROLL_THRESHOLD, 1);
      setProgress(p);
    });
  }, []);

  useEffect(() => {
    handleScroll();
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
      if (rafRef.current !== null) cancelAnimationFrame(rafRef.current);
    };
  }, [handleScroll]);

  const logOut = () => {
    clearAuth();
    navigate("/");
  };

  const logoScale = 1 - progress * 0.3;
  const logoTranslateX = -progress * 24;
  const menuTranslateX = progress * 32;
  const logoR = Math.round(255 - progress * 221);
  const logoG = Math.round(255 - progress * 221);
  const logoB = Math.round(255 - progress * 221);
  const logoColor = `rgb(${logoR},${logoG},${logoB})`;
  const textShadow =
    progress > 0.5 ? `0 1px 2px rgba(0,0,0,${(progress - 0.5) * 0.15})` : undefined;
  const bgOpacity = progress;
  const shadowOpacity = Math.max(0, (progress - 0.5) * 2);
  const headerPaddingY = 24 - progress * 20;

  return (
    <>
      {HEADER_HEIGHT > 0 && <div style={{ height: HEADER_HEIGHT }} />}
      <div
        ref={headerRef}
        className="w-full"
        style={{
          position: "fixed",
          top: 0,
          left: 0,
          zIndex: 50,
        }}
      >
        <div
          className="absolute inset-0 bg-white pointer-events-none"
          style={{
            opacity: bgOpacity,
            boxShadow: shadowOpacity > 0 ? `0 1px 8px rgba(0,0,0,${shadowOpacity * 0.1})` : "none",
          }}
        />
        <div
          className="flex items-center justify-between relative"
          style={{
            maxWidth: 1024,
            margin: "0 auto",
            paddingTop: headerPaddingY,
            paddingBottom: headerPaddingY,
          }}
        >
          <div
            style={{
              transform: `scale(${logoScale}) translateX(${logoTranslateX}px)`,
              transformOrigin: "left center",
              willChange: "transform",
            }}
          >
            <Link
              to="/"
              className="no-underline flex items-baseline gap-1 drop-shadow-md"
              style={{
                color: logoColor,
                textShadow: textShadow,
              }}
            >
              <span className="text-3xl font-normal italic tracking-widest">無處告別</span>
              <i className="iconfont icon-zhifeiji text-4xl" />
            </Link>
          </div>
          <div
            className="inline-flex items-center h-[35px] bg-white rounded-lg whitespace-nowrap"
            style={{
              transform: `translateX(${menuTranslateX}px)`,
              willChange: "transform",
            }}
          >
            <ul className="flex items-center list-none px-4">
              <li className="bg-white">
                <Link
                  to="/"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                >
                  首页
                </Link>
              </li>
              <li className="bg-white">
                <a
                  href="https://windiness.fifsky.com"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                >
                  有风
                </a>
              </li>
              <li className="bg-white">
                <Link
                  to="/archive"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                >
                  归档
                </Link>
              </li>
              <li className="bg-white">
                <Link
                  to="/links"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                >
                  友链
                </Link>
              </li>
              <li className="bg-white">
                <Link
                  to="/about"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                >
                  关于
                </Link>
              </li>
              <li className="bg-white">
                <a
                  href="https://www.travellings.cn/go"
                  className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                  target="_blank"
                >
                  开往
                </a>
              </li>
              {isLogin && (
                <li className="bg-white">
                  <Link
                    to="/admin/index"
                    className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                  >
                    管理中心
                  </Link>
                </li>
              )}
              {isLogin && (
                <li className="bg-white">
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      logOut();
                    }}
                    className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                  >
                    退出
                  </a>
                </li>
              )}
              {!isLogin && (
                <li className="bg-white">
                  <Link
                    to="/login"
                    className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
                  >
                    登录
                  </Link>
                </li>
              )}
            </ul>
          </div>
        </div>
      </div>
    </>
  );
}
