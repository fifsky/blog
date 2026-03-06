import { useState } from "react";
import { useParams, Link } from "react-router";
import { CArticle } from "@/components/CArticle";
import { Comment } from "@/components/Comment";
import { PageTransition } from "@/components/PageTransition";
import { SkeletonArticle } from "@/components/Skeleton";
import { articleDetailApi, prevnextArticleApi, settingApi } from "@/service";
import { ArticleItem, PrevNextItem, Options } from "@/types/openapi";
import { useAsyncEffect } from "@/hooks";

// 生成文章链接：有自定义路径则使用 /${url}，否则 /article/${id}
function articleLink(item: PrevNextItem): string {
  return item.url ? `/${item.url}` : `/article/${item.id}`;
}

export default function ArticleDetail() {
  const [article, setArticle] = useState<ArticleItem>();
  const [data, setData] = useState<{ prev?: PrevNextItem; next?: PrevNextItem }>({});
  const [settings, setSettings] = useState<Options>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<any>(null);
  const params = useParams();

  useAsyncEffect(async () => {
    setLoading(true);
    // 支持 /article/:id 和 /:path 两种路由
    const id = params.id ? parseInt(params.id) : undefined;
    const url = params.path || undefined;

    // 传入 errorHandler 将异常保存到 state 以便在渲染阶段抛出
    const errorHandler = (e: any) => { setError(e); };
    const [a, s] = await Promise.all([articleDetailApi({ id, url }, errorHandler), settingApi()]);
    setArticle(a);
    setSettings(s);
    // 使用文章实际 ID 获取上下篇
    if (a?.id) {
      const pn = await prevnextArticleApi({ id: a.id });
      setData(pn);
    }
    setLoading(false);
  }, [params.id, params.path]);

  if (error) {
    throw error;
  }

  const siteName = settings?.kv?.site_name || "無處告別";
  const pageTitle = `${article?.title ? article?.title + " - " : ""}${siteName}`;

  // Show skeleton during initial load
  if (!article?.id) {
    return <SkeletonArticle className={"h-200"} />;
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
                <Link to={articleLink(data.prev)}>{data.prev.title}</Link>
              ) : (
                <span>嘿，这已经是最新的文章啦</span>
              )}
            </div>
            <div className="w-[400px] overflow-hidden text-ellipsis whitespace-nowrap text-right">
              <strong>下一篇：</strong>
              {data.next?.id ? (
                <Link to={articleLink(data.next)}>{data.next.title}</Link>
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
