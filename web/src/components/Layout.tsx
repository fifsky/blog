import { useEffect } from "react";
import { Outlet, NavLink } from "react-router";
import { CHeader } from "./CHeader";
import { Sidebar } from "./Sidebar";
import { Mood } from "./Mood";
import { CFooter } from "./CFooter";
import { useStore } from "@/store/context";
import { getAccessToken } from "@/utils/common";
import { cn } from "@/lib/utils";

function NavItem({ to, children, end }: { to: string; children: React.ReactNode; end?: boolean }) {
  return (
    <li className="ml-1.5">
      <NavLink
        to={to}
        end={end}
        className={({ isActive }) =>
          cn(
            "no-underline inline-flex text-gray-800 border border-[#89d5ef]",
            isActive
              ? "mt-0 px-4 py-1.5 pb-1 border-b-[#fff] bg-white"
              : "mt-1.5 px-3.5 py-0.5 bg-[#89d5ef] hover:bg-white hover:text-[#ff7031]",
          )
        }
      >
        {children}
      </NavLink>
    </li>
  );
}

export function Layout() {
  const currentUserAction = useStore((s) => s.currentUserAction);

  useEffect(() => {
    if (getAccessToken()) {
      currentUserAction();
    }
  }, []);

  return (
    <div className="w-[1024px] mt-4 mx-auto min-h-[500px]">
      <CHeader />
      <div className="flex justify-between items-start">
        <div className="w-[778px] overflow-hidden">
          <Mood />
          <div className="tabs relative top-px">
            <ul className="flex justify-end list-none">
              <NavItem to="/travel">旅途</NavItem>
              <NavItem to="/about">关于我</NavItem>
              <NavItem to="/" end>
                所有文章
              </NavItem>
            </ul>
          </div>
          <div className="p-5 border border-[#89d5ef] bg-white min-h-100">
            <Outlet />
          </div>
        </div>
        <Sidebar />
      </div>
      <CFooter />
    </div>
  );
}
