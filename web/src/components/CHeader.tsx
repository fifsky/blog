import { Link, useNavigate } from "react-router";
import { useStore } from "@/store/context";

export function CHeader() {
  const { state, dispatch } = useStore();
  const isLogin = !!state.userInfo.id;
  const navigate = useNavigate();
  const logOut = () => {
    localStorage.removeItem("access_token");
    dispatch({ type: "setUserInfo", payload: {} });
    navigate("/");
  };
  return (
    <div id="header" className="clearfix">
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
      <div id="navigation">
        <ul>
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
              <a onClick={logOut}>退出</a>
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
