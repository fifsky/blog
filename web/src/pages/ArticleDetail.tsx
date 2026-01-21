import { useState } from "react";
import { useParams, Link } from "react-router";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { PageTransition } from "@/components/PageTransition";
import { SkeletonArticle } from "@/components/Skeleton";
import { articleDetailApi, prevnextArticleApi, settingApi } from "@/service";
import { ArticleItem, PrevNextItem, Options } from "@/types/openapi";
import { useAsyncEffect } from "@/hooks";

export default function ArticleDetail() {
  const [article, setArticle] = useState<ArticleItem>();
  const [data, setData] = useState<{ prev?: PrevNextItem; next?: PrevNextItem }>({});
  const [settings, setSettings] = useState<Options>();
  const [loading, setLoading] = useState(true);
  const params = useParams();

  useAsyncEffect(async () => {
    const id = params.id ? parseInt(params.id) : undefined;
    setLoading(true);
    const [a, s] = await Promise.all([articleDetailApi({ id }), settingApi()]);
    setArticle(a);
    setSettings(s);
    if (id) {
      const pn = await prevnextArticleApi({ id });
      setData(pn);
    }
    setLoading(false);
  }, [params.id]);

  const siteName = settings?.kv?.site_name || "無處告別";
  const pageTitle = `${article?.title ? article?.title + " - " : ""}${siteName}`;

  // Show skeleton during initial load
  if (!article?.id) {
    return <SkeletonArticle />;
  }

  return (
    <>
      <title>{pageTitle}</title>
      <meta name="description" content={settings?.kv?.site_desc || ""} />
      <meta name="keywords" content={settings?.kv?.site_keyword || ""} />
      <PageTransition loading={loading}>
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
      </PageTransition>
    </>
  );
}
