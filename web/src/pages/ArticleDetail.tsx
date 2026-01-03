import { useEffect, useState } from "react";
import { useParams, Link } from "react-router";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { articleDetailApi, prevnextArticleApi } from "@/service";

export default function ArticleDetail() {
  const [article, setArticle] = useState<any>({});
  const [data, setData] = useState<{ prev?: any; next?: any }>({});
  const params = useParams();
  useEffect(() => {
    (async () => {
      const id = params.id ? parseInt(params.id) : undefined;
      const a = await articleDetailApi({ id });
      setArticle(a);
      if (id) {
        const pn = await prevnextArticleApi({ id });
        setData(pn);
        document.title = `${a.title} - 無處告別`;
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.id]);
  if (!article.id) return null;
  return (
    <div>
      <div className="article-single">
        <CArticle article={article} />
        <div className="post-navi flex justify-between">
          <div className="prev">
            <strong>上一篇：</strong>
            {data.prev && data.prev.id ? (
              <Link to={`/article/${data.prev.id}`}>{data.prev.title}</Link>
            ) : (
              <span>嘿，这已经是最新的文章啦</span>
            )}
          </div>
          <div className="next text-right">
            <strong>下一篇：</strong>
            {data.next && data.next.id ? (
              <Link to={`/article/${data.next.id}`}>{data.next.title}</Link>
            ) : (
              <span>嘿，这已经是最后的文章啦</span>
            )}
          </div>
        </div>
      </div>
      <Comment postId={article.id} />
    </div>
  );
}
