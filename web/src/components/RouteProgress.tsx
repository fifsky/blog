import { useEffect } from "react";
import { useLocation } from "react-router";
// @ts-ignore
import NProgress from "nprogress";
import "nprogress/nprogress.css";
NProgress.configure({ showSpinner: true, trickleSpeed: 200, minimum: 0.12 });

export function RouteProgress() {
  const location = useLocation();
  useEffect(() => {
    NProgress.start();
    const t = setTimeout(() => NProgress.done(), 200);
    return () => clearTimeout(t);
  }, [location.key]);
  return null;
}
