import { useEffect } from "react";
import { Outlet, useLocation, Link, useNavigate } from "react-router";
import { CHeader } from "./CHeader";
import { CFooter } from "./CFooter";
import { AIChat } from "./ai/AIChat";
import { useStore } from "@/store/context";
import { getAccessToken } from "@/utils/common";
import { cn } from "@/lib/utils";

function AdminNavItem({
  to,
  children,
  isActive,
}: {
  to: string;
  children: React.ReactNode;
  isActive?: boolean;
}) {
  return (
    <li className="ml-1.5">
      <Link
        to={to}
        className={cn(
          "no-underline inline-flex text-gray-800 border border-[#89d5ef]",
          isActive
            ? "mt-0 px-4 py-1.5 pb-1 border-b-[#fff] bg-white"
            : "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] hover:bg-white hover:text-[#ff7031]",
        )}
      >
        {children}
      </Link>
    </li>
  );
}

export function AdminLayout() {
  const userInfo = useStore((s) => s.userInfo);
  const currentUserAction = useStore((s) => s.currentUserAction);
  const isLogin = !!userInfo.id;
  const location = useLocation();
  const navigate = useNavigate();

  const isPage = (...paths: string[]) => paths.some((p) => location.pathname === p);

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
    <div className="w-[1024px] mt-4 mx-auto min-h-[500px]">
      <CHeader />
      {isLogin && (
        <div>
          <div className="tabs relative top-px">
            <ul className="flex justify-start list-none">
              <AdminNavItem to="/admin/index" isActive={isPage("/admin/index")}>
                设置
              </AdminNavItem>
              <AdminNavItem
                to="/admin/articles"
                isActive={isPage("/admin/articles", "/admin/post/article")}
              >
                文章
              </AdminNavItem>
              <AdminNavItem to="/admin/photos" isActive={isPage("/admin/photos")}>
                相册
              </AdminNavItem>
              <li className="ml-1.5">
                <a
                  href="https://github.com/fifsky/blog/discussions/categories/comment"
                  target="_blank"
                  rel="noopener noreferrer"
                  className={cn(
                    "no-underline inline-flex text-gray-800 border border-[#89d5ef]",
                    "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] hover:bg-white hover:text-[#ff7031]",
                  )}
                >
                  评论
                </a>
              </li>
              <AdminNavItem to="/admin/moods" isActive={isPage("/admin/moods")}>
                心情
              </AdminNavItem>
              <AdminNavItem to="/admin/cates" isActive={isPage("/admin/cates")}>
                分类
              </AdminNavItem>
              <AdminNavItem to="/admin/links" isActive={isPage("/admin/links")}>
                链接
              </AdminNavItem>
              <AdminNavItem to="/admin/remind" isActive={isPage("/admin/remind")}>
                提醒
              </AdminNavItem>
              <AdminNavItem to="/admin/users" isActive={isPage("/admin/users")}>
                用户
              </AdminNavItem>
            </ul>
          </div>
          <div className="p-5 border border-[#89d5ef] bg-white">
            <Outlet />
          </div>
        </div>
      )}
      <CFooter />
      {isLogin && <AIChat />}
    </div>
  );
}
