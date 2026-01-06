import React, { useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import { useStore } from "@/store/context";
import { SidebarList } from "./SidebarList";
import { Calendar } from "./Calendar";
import {
  InputGroup,
  InputGroupInput,
  InputGroupAddon,
} from "@/components/ui/input-group";
import { Search } from "lucide-react";

export function Sidebar() {
  const keyword = useStore((s) => s.keyword);
  const setKeyword = useStore((s) => s.setKeyword);
  const navigate = useNavigate();
  const location = useLocation();
  const changeKeyword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setKeyword(e.target.value);
  };
  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      navigate({
        pathname: "/search",
        search: `?keyword=${encodeURIComponent(keyword)}`,
      });
    }
  };
  useEffect(() => {
    const params = new URLSearchParams(location.search);
    setKeyword(params.get("keyword") || "");
  }, [location.search]);
  return (
    <div className="p-[15px] border border-[#89d5ef] bg-white">
      <div className="mb-5">
        <InputGroup>
          <InputGroupInput
            placeholder="搜索..."
            name="keyword"
            value={keyword}
            onChange={changeKeyword}
            onKeyDown={onKeyDown}
          />
          <InputGroupAddon>
            <Search />
          </InputGroupAddon>
        </InputGroup>
      </div>
      <Calendar />
      <SidebarList title="文章分类" api="cateAllApi" />
      <SidebarList title="历史存档" api="archiveApi" />
      <SidebarList title="我关注的" api="linkAllApi" />
      <div className="mt-5">
        <i className="iconfont icon-rss" style={{ color: "orange" }}></i>
        <a
          className="pl-[5px]"
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
