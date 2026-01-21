import { ArticleItem } from "@/types/openapi";
import { useMemo } from "react";
import { Link, useLocation } from "react-router";
import { Viewer } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import breaks from "@bytemd/plugin-breaks";
import { highlightPlugin } from "@/lib/highlight-plugin";
import mediumZoom from "@bytemd/plugin-medium-zoom";

// ByteMD 插件配置
const plugins = [gfm(), breaks(), highlightPlugin(), mediumZoom()];

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
          <h2 className="text-[16px] font-medium">
            <Link
              to={`/article/${article.id}`}
              onClick={() => window.scrollTo({ top: 0 })}
              dangerouslySetInnerHTML={{
                __html: markHigh(article.title, keyword),
              }}
            />
          </h2>
          <div className="flex items-center gap-2 text-[12px] text-[#999]">
            <span>by {article.user.nick_name}</span>
            <span>/</span>
            <Link
              to={`/category/${article.cate.domain}`}
              rel="category tag"
              title={`查看 ${article.cate.name} 中的全部文章`}
            >
              {article.cate.name}
            </Link>
            <span>/</span>
            <span>{article.created_at}</span>
          </div>
        </div>
      </div>
      <div className="mt-2">
        <Viewer value={article.content} plugins={plugins} />
      </div>
    </div>
  );
}
