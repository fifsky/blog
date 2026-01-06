import { useEffect } from "react";
import { Outlet, useLocation, Link, useNavigate } from "react-router";
import { CHeader } from "./CHeader";
import { CFooter } from "./CFooter";
import { useStore } from "@/store/context";
import { getAccessToken } from "@/utils/common";

export function AdminLayout() {
  const userInfo = useStore((s) => s.userInfo);
  const currentUserAction = useStore((s) => s.currentUserAction);
  const isLogin = !!userInfo.id;
  const location = useLocation();
  const navigate = useNavigate();

  const isPage = (...paths: string[]) =>
    paths.some((p) => location.pathname === p);

  // 导航链接样式变量
  const activeNavClass = "mt-0 px-4 py-1.5 pb-1 border border-[#89d5ef] border-b-[#fff] bg-white no-underline inline-flex text-gray-800";
  const navClass = "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] border border-[#89d5ef] text-gray-800 no-underline inline-flex hover:bg-white hover:text-[#ff7031]";

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
  }, []);

  useEffect(() => {
    if (!getAccessToken()) navigate("/login");
  }, [location.pathname]);

  return (
    <div id="container">
      <CHeader />
      {isLogin && (
        <div className="admin">
          <div className="tabs relative top-px">
            <ul className="flex justify-start list-none">
              <li className="ml-1.5">
                <Link
                  to="/admin/index"
                  className={isPage("/admin/index") ? activeNavClass : navClass}
                >
                  设置
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/articles"
                  className={isPage("/admin/articles", "/admin/post/article") ? activeNavClass : navClass}
                >
                  文章
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/comments"
                  className={isPage("/admin/comments") ? activeNavClass : navClass}
                >
                  评论
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/moods"
                  className={isPage("/admin/moods") ? activeNavClass : navClass}
                >
                  心情
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/cates"
                  className={isPage("/admin/cates") ? activeNavClass : navClass}
                >
                  分类
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/links"
                  className={isPage("/admin/links") ? activeNavClass : navClass}
                >
                  链接
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/remind"
                  className={isPage("/admin/remind") ? activeNavClass : navClass}
                >
                  提醒
                </Link>
              </li>
              <li className="ml-1.5">
                <Link
                  to="/admin/users"
                  className={isPage("/admin/users") ? activeNavClass : navClass}
                >
                  用户
                </Link>
              </li>
            </ul>
          </div>
          <div className="p-5 border border-[#89d5ef] bg-white">
            <Outlet />
          </div>
        </div>
      )}
      <CFooter />
    </div>
  );
}
