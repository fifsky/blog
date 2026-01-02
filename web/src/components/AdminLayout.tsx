import { useEffect } from "react";
import { Outlet, useLocation, Link, useNavigate } from "react-router";
import { CHeader } from "./CHeader";
import { CFooter } from "./CFooter";
import { useStore } from "@/store/context";
import { getAccessToken } from "@/utils/common";

export function AdminLayout() {
  const { state, currentUserAction } = useStore();
  const isLogin = !!state.userInfo.id;
  const location = useLocation();
  const navigate = useNavigate();

  const isPage = (...paths: string[]) =>
    paths.some((p) => location.pathname === p);

  useEffect(() => {
    if (getAccessToken()) {
      (async () => {
        try {
          await currentUserAction();
        } catch {
          navigate("/login");
        }
      })();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!getAccessToken()) navigate("/login");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname]);

  return (
    <div id="container">
      <CHeader />
      {isLogin && (
        <div className="admin">
          <div className="tabs">
            <ul>
              <li>
                <Link
                  to="/admin/index"
                  className={isPage("/admin/index") ? "active" : ""}
                >
                  设置
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/articles"
                  className={
                    isPage("/admin/articles", "/admin/post/article")
                      ? "active"
                      : ""
                  }
                >
                  文章
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/comments"
                  className={isPage("/admin/comments") ? "active" : ""}
                >
                  评论
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/moods"
                  className={isPage("/admin/moods") ? "active" : ""}
                >
                  心情
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/cates"
                  className={isPage("/admin/cates") ? "active" : ""}
                >
                  分类
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/links"
                  className={isPage("/admin/links") ? "active" : ""}
                >
                  链接
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/remind"
                  className={isPage("/admin/remind") ? "active" : ""}
                >
                  提醒
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/users"
                  className={isPage("/admin/users") ? "active" : ""}
                >
                  用户
                </Link>
              </li>
            </ul>
          </div>
          <div id="content">
            <Outlet />
          </div>
        </div>
      )}
      <CFooter />
    </div>
  );
}
