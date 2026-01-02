import React, { useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import { useStore } from "@/store/context";
import { SidebarList } from "./SidebarList";
import { Calendar } from "./Calendar";

export function Sidebar() {
  const { state, dispatch } = useStore();
  const navigate = useNavigate();
  const location = useLocation();
  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    navigate({
      pathname: "/search",
      search: `?keyword=${encodeURIComponent(state.keyword)}`,
    });
  };
  const changeKeyword = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch({ type: "setKeyword", payload: e.target.value });
  };
  useEffect(() => {
    const params = new URLSearchParams(location.search);
    dispatch({ type: "setKeyword", payload: params.get("keyword") || "" });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.search]);
  return (
    <div id="sidebar">
      <div className="sect" id="search">
        <form id="searchpanel" method="get" onSubmit={submit}>
          <p>
            <input
              className="input_text"
              type="text"
              name="keyword"
              onChange={changeKeyword}
              value={state.keyword}
            />
            &nbsp;
            <input className="formbutton" type="submit" value="搜索" />
          </p>
        </form>
      </div>
      <Calendar />
      <SidebarList title="文章分类" api="cateAllApi" />
      <SidebarList title="历史存档" api="archiveApi" />
      <SidebarList title="我关注的" api="linkAllApi" />
      <div className="rss">
        <i className="iconfont icon-rss" style={{ color: "orange" }}></i>
        <a
          href="https://api.fifsky.com/feed.xml"
          target="_blank"
          rel="noreferrer"
        >
          订阅我的消息
        </a>
      </div>
    </div>
  );
}
