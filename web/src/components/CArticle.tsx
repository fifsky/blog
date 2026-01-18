import { ArticleItem } from "@/types/openapi";
import { useMemo } from "react";
import { Link, useLocation } from "react-router";
import { Viewer } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import highlight from "@bytemd/plugin-highlight";
import mediumZoom from "@bytemd/plugin-medium-zoom";

// ByteMD 插件配置
const plugins = [gfm(), highlight(), mediumZoom()];

export function CArticle({ article }: { article: ArticleItem }) {
  const location = useLocation();
  const keyword = useMemo(
    () => new URLSearchParams(location.search).get("keyword") || "",
    [location.search],
  );
  const markHigh = (content: string, k: string) => {
    if (!k) return content;
    return content.replace(k, `<mark>${k}</mark>`);
  };
  if (!article) return null;
  return (
    <div>
      <div className="flex justify-between items-center h-[54px] overflow-hidden">
        <img className="p-[2px] w-[40px] h-[40px]" src="/assets/images/avatar.jpg" alt="" />
        <div className="flex-1 ml-4">
          <h2 className="text-[14px]">
            <Link
              to={`/article/${article.id}`}
              dangerouslySetInnerHTML={{
                __html: markHigh(article.title, keyword),
              }}
            />
          </h2>
          <div className="text-[12px] text-[#999]">
            by&nbsp;{article.user.nick_name}&nbsp;&nbsp;/&nbsp;&nbsp;
            <Link
              to={`/category/${article.cate.domain}`}
              rel="category tag"
              title={`查看 ${article.cate.name} 中的全部文章`}
            >
              {article.cate.name}
            </Link>
            &nbsp;&nbsp;/&nbsp;&nbsp;
            {article.created_at}
          </div>
        </div>
      </div>
      <div className="article">
        <Viewer value={article.content} plugins={plugins} />
      </div>
    </div>
  );
}
