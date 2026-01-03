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
    <div id="header" className="flex items-center justify-between">
      <h1>
        <Link to="/">
          <img
            className="inline-block"
            alt="首页"
            src="/assets/images/logo.png"
          />
          <span>首页</span>
        </Link>
      </h1>
      <div id="navigation" className="inline-flex items-center">
        <ul className="flex items-center">
          <li>
            <Link to="/">首页</Link>
          </li>
          <li>
            <a
              href="https://github.com/fifsky"
              target="_blank"
              rel="noreferrer"
            >
              Github
            </a>
          </li>
          <li>
            <a href="https://gist.github.com/fifsky">技术</a>
          </li>
          <li>
            <Link to="/about">关于</Link>
          </li>
          <li>
            <a href="https://caixudong.com">简历</a>
          </li>
          {isLogin && (
            <li>
              <Link to="/admin/index">管理中心</Link>
            </li>
          )}
          {isLogin && (
            <li>
              <a
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  logOut();
                }}
              >
                退出
              </a>
            </li>
          )}
          {!isLogin && (
            <li>
              <Link to="/login">登录</Link>
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}
