import { Link, useNavigate } from "react-router";
import { useStore } from "@/store/context";

export function CHeader() {
  const userInfo = useStore((s) => s.userInfo);
  const setUserInfo = useStore((s) => s.setUserInfo);
  const isLogin = !!userInfo.id;
  const navigate = useNavigate();
  const logOut = () => {
    localStorage.removeItem("access_token");
    setUserInfo({});
    navigate("/");
  };
  return (
    <div className="flex items-center justify-between">
      <div className="group pt-1 pb-5 px-0">
        <Link to="/" className="no-underline flex items-baseline gap-1 drop-shadow-md">
          <span className="text-white text-4xl font-black tracking-wider">
            你好
          </span>
          <span className="text-white text-3xl font-bold">
            。
          </span>
          <span className="text-white text-2xl font-normal italic tracking-widest">
            旧时光
          </span>
        </Link>
      </div>
      <div className="inline-flex items-center h-[35px] my-2 px-3 bg-white rounded-lg whitespace-nowrap">
        <ul className="flex items-center list-none">
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
              href="https://github.com/fifsky"
              target="_blank"
              rel="noreferrer"
              className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
            >
              Github
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
              to="/about"
              className="px-2.5 py-0.5 hover:bg-[#0066cc] hover:text-white hover:no-underline"
            >
              关于
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
  );
}
