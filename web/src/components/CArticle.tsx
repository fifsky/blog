import { useEffect, useMemo, useRef } from "react";
import hljs from "highlight.js/lib/core";
import { Link, useLocation } from "react-router";
export function CArticle({ article }: { article: any }) {
  const rootRef = useRef<HTMLDivElement>(null);
  const location = useLocation();
  const keyword = useMemo(
    () => new URLSearchParams(location.search).get("keyword") || "",
    [location.search]
  );
  const markHigh = (content: string, k: string) => {
    if (!k) return content;
    return content.replace(k, `<mark>${k}</mark>`);
  };
  useEffect(() => {
    const root = rootRef.current;
    if (!root) return;
    root.querySelectorAll("pre code").forEach((block) => {
      hljs.highlightElement(block as HTMLElement);
    });
  }, [article]);
  if (!article) return null;
  return (
    <div className="article" ref={rootRef}>
      <div className="flex justify-between items-center h-[54px] overflow-hidden">
        <img
          className="p-[2px] w-[40px] h-[40px]"
          src="/assets/images/avatar.jpg"
          alt=""
        />
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
              to={`/categroy/${article.cate.domain}`}
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
      <div
        className="entry"
        dangerouslySetInnerHTML={{ __html: article.content }}
      />
    </div>
  );
}
