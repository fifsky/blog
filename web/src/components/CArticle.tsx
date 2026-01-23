import { ArticleItem } from "@/types/openapi";
import { useMemo } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { Viewer } from "@bytemd/react";
import gfm from "@bytemd/plugin-gfm";
import breaks from "@bytemd/plugin-breaks";
import { highlightPlugin } from "@/lib/highlight-plugin";
import mediumZoom from "@bytemd/plugin-medium-zoom";
import { Badge } from "@/components/ui/badge";

// ByteMD 插件配置
const plugins = [gfm(), breaks(), highlightPlugin(), mediumZoom()];

export function CArticle({ article }: { article: ArticleItem }) {
  const location = useLocation();
  const navigate = useNavigate();
  const keyword = useMemo(
    () => new URLSearchParams(location.search).get("keyword") || "",
    [location.search],
  );
  const activeTag = useMemo(
    () => new URLSearchParams(location.search).get("tag") || "",
    [location.search],
  );
  const markHigh = (content: string, k: string) => {
    if (!k) return content;
    return content.replace(k, `<mark>${k}</mark>`);
  };
  const tags = useMemo(() => (article.tags || []).filter(Boolean), [article.tags]);

  const filterByTag = (tag: string) => {
    const q = new URLSearchParams(location.search);
    q.set("tag", tag);
    q.delete("page");
    const pathname = location.pathname.startsWith("/article/") ? "/" : location.pathname;
    navigate({ pathname, search: q.toString() });
    window.scrollTo({ top: 0 });
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
      {tags.length > 0 ? (
        <div className="mt-4 flex w-full flex-wrap justify-end gap-2">
          {tags.map((t) => (
            <Badge
              key={t}
              variant={t === activeTag ? "default" : "secondary"}
              className="cursor-pointer select-none"
              role="button"
              tabIndex={0}
              onClick={() => filterByTag(t)}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  filterByTag(t);
                }
              }}
            >
              {t}
            </Badge>
          ))}
        </div>
      ) : null}
    </div>
  );
}
