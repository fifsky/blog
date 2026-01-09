import { useEffect, useState } from "react";
import { useParams, Link } from "react-router";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { articleDetailApi, prevnextArticleApi } from "@/service";
import { ArticleItem, PrevNextItem } from "@/types/openapi";
import { Spinner } from "@/components/ui/spinner";

export default function ArticleDetail() {
  const [article, setArticle] = useState<ArticleItem>();
  const [data, setData] = useState<{ prev?: PrevNextItem; next?: PrevNextItem }>({});
  const params = useParams();
  useEffect(() => {
    (async () => {
      const id = params.id ? parseInt(params.id) : undefined;
      setArticle(undefined);
      const a = await articleDetailApi({ id });
      setArticle(a);
      if (id) {
        const pn = await prevnextArticleApi({ id });
        setData(pn);
      }
    })();
  }, [params.id]);
  const pageTitle = `${article?.title ? article?.title + " - " : ""}無處告別`;
  if (!article?.id) {
    return (
      <>
        <title>{pageTitle}</title>
        <div className="flex items-center justify-center py-20 text-gray-600">
          <Spinner className="size-6" />
          <span className="ml-2 text-sm">加载中...</span>
        </div>
      </>
    );
  }
  return (
    <>
      <title>{pageTitle}</title>
      <div className="mb-[10px]">
        <CArticle article={article} />
        <div className="my-5 flex justify-between">
          <div className="w-[400px] overflow-hidden text-ellipsis whitespace-nowrap">
            <strong>上一篇：</strong>
            {data.prev?.id ? (
              <Link to={`/article/${data.prev.id}`}>{data.prev.title}</Link>
            ) : (
              <span>嘿，这已经是最新的文章啦</span>
            )}
          </div>
          <div className="w-[400px] overflow-hidden text-ellipsis whitespace-nowrap text-right">
            <strong>下一篇：</strong>
            {data.next?.id ? (
              <Link to={`/article/${data.next.id}`}>{data.next.title}</Link>
            ) : (
              <span>嘿，这已经是最后的文章啦</span>
            )}
          </div>
        </div>
      </div>
      <Comment />
    </>
  );
}
