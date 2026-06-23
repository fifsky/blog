import React, { useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router";
import { useStore } from "@/store/context";
import { SidebarList } from "./SidebarList";
import { Calendar } from "./Calendar";
import { InputGroup, InputGroupInput, InputGroupAddon } from "@/components/ui/input-group";
import { Search, Github, Mail, Send, FileText } from "lucide-react";

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
  const [weather, setWeather] = useState("");
  useEffect(() => {
    fetch("https://wttr.in/Shanghai?format=3&AT")
      .then((r) => r.text())
      .then((t) => setWeather(t.trim()))
      .catch(() => {});
  }, []);
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
      {weather && <p className="mb-4 text-center text-xs text-[#8c8c8c]">{weather}</p>}
      <Calendar />
      {/* 联系方式图标 */}
      <div className="mb-5 flex items-center justify-center gap-3">
        <a
          href="https://github.com/fifsky"
          target="_blank"
          rel="noreferrer"
          title="GitHub"
          className="grid place-items-center leading-none w-7 h-7 rounded-full border border-gray-300 text-gray-500 hover:border-[#06c] hover:text-[#06c] transition-colors [&_svg]:block [&_svg]:m-auto"
        >
          <Github size={14} />
        </a>
        <a
          href="mailto:fifsky@gmail.com"
          title="邮箱"
          className="grid place-items-center leading-none w-7 h-7 rounded-full border border-gray-300 text-gray-500 hover:border-[#06c] hover:text-[#06c] transition-colors [&_svg]:block [&_svg]:m-auto"
        >
          <Mail size={14} />
        </a>
        <a
          href="https://t.me/+_Yk3wbIcqM82YTI1"
          target="_blank"
          rel="noreferrer"
          title="Telegram"
          className="grid place-items-center leading-none w-7 h-7 rounded-full border border-gray-300 text-gray-500 hover:border-[#06c] hover:text-[#06c] transition-colors [&_svg]:block [&_svg]:m-auto"
        >
          <Send size={14} />
        </a>
        <a
          href="https://caixudong.com"
          target="_blank"
          rel="noreferrer"
          title="简历"
          className="grid place-items-center leading-none w-7 h-7 rounded-full border border-gray-300 text-gray-500 hover:border-[#06c] hover:text-[#06c] transition-colors [&_svg]:block [&_svg]:m-auto"
        >
          <FileText size={14} />
        </a>
      </div>
      <SidebarList title="文章分类" api="cateAllApi" />
      <SidebarList title="历史存档" api="archiveApi" />
      <div className="mt-5 flex items-center gap-4">
        <div>
          <i className="iconfont icon-rss" style={{ color: "orange" }}></i>
          <a
            className="pl-[5px]"
            href="https://api.fifsky.com/blog/feed.xml"
            target="_blank"
            rel="noreferrer"
          >
            订阅我的消息
          </a>
        </div>
        <a
          className="text-gray-500 hover:text-blue-500 hover:underline"
          href="/llms.txt"
          target="_blank"
          rel="noreferrer"
          title="For AI Agents"
        >
          llms.txt
        </a>
      </div>
    </div>
  );
}
