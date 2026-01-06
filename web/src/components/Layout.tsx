import { useEffect } from "react";
import { Outlet, NavLink } from "react-router";
import { CHeader } from "./CHeader";
import { Sidebar } from "./Sidebar";
import { Mood } from "./Mood";
import { CFooter } from "./CFooter";
import { useStore } from "@/store/context";
import { getAccessToken } from "@/utils/common";

export function Layout() {
  const currentUserAction = useStore((s) => s.currentUserAction);

  useEffect(() => {
    if (getAccessToken()) {
      currentUserAction();
    }
  }, []);

  return (
    <div id="container">
      <CHeader />
      <div className="flex justify-between items-start">
        <div id="main">
          <Mood />
          <div className="tabs">
            <ul className="flex justify-end list-none">
              <li>
                <NavLink
                  to="/about"
                  className={({ isActive }) => (isActive ? "active" : "")}
                >
                  关于我
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/"
                  end
                  className={({ isActive }) => (isActive ? "active" : "")}
                >
                  所有文章
                </NavLink>
              </li>
            </ul>
          </div>
          <div id="content">
            <Outlet />
          </div>
        </div>
        <Sidebar />
      </div>
      <CFooter />
    </div>
  );
}
