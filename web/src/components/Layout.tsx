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
          <div className="tabs relative top-px">
            <ul className="flex justify-end list-none">
              <li className="ml-1.5">
                <NavLink
                  to="/about"
                  className={({ isActive }) =>
                    isActive
                      ? "mt-0 px-4 py-1.5 pb-1 border border-[#89d5ef] border-b-[#fff] bg-white no-underline inline-flex text-gray-800"
                      : "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] border border-[#89d5ef] text-gray-800 no-underline inline-flex hover:bg-white hover:text-[#ff7031]"
                  }
                >
                  关于我
                </NavLink>
              </li>
              <li className="ml-1.5">
                <NavLink
                  to="/"
                  end
                  className={({ isActive }) =>
                    isActive
                      ? "mt-0 px-4 py-1.5 pb-1 border border-[#89d5ef] border-b-[#fff] bg-white no-underline inline-flex text-gray-800"
                      : "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] border border-[#89d5ef] text-gray-800 no-underline inline-flex hover:bg-white hover:text-[#ff7031]"
                  }
                >
                  所有文章
                </NavLink>
              </li>
            </ul>
          </div>
          <div className="p-5 border border-[#89d5ef] bg-white">
            <Outlet />
          </div>
        </div>
        <Sidebar />
      </div>
      <CFooter />
    </div>
  );
}
