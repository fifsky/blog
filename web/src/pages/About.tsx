import { useMemo, useState } from "react";
import { Comment } from "@/components/Comment";
import { PageTransition } from "@/components/PageTransition";
import { SkeletonArticle } from "@/components/Skeleton";
import { articleDetailApi, settingApi } from "@/service";
import { ArticleItem, Setting } from "@/types/openapi";
import { useAsyncEffect } from "@/hooks";

interface TimelineEntry {
  date: string;
  description: string;
}

/** 解析时间轴内容：按行分割，第一个空格前为日期，其后为描述 */
function parseTimeline(content: string): TimelineEntry[] {
  return content
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0)
    .map((line) => {
      const spaceIndex = line.indexOf(" ");
      if (spaceIndex === -1) return null;
      const date = line.substring(0, spaceIndex);
      const description = line.substring(spaceIndex + 1).trim();
      // 日期必须以4位数字开头
      if (!/^\d{4}/.test(date) || !description) return null;
      return { date, description } as TimelineEntry;
    })
    .filter((entry): entry is TimelineEntry => entry !== null);
}

/** 关于页面：通过固定 path "about" 获取文章数据，以时间轴样式展示 */
export default function About() {
  const [article, setArticle] = useState<ArticleItem>();
  const [settings, setSettings] = useState<Setting>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<any>(null);

  useAsyncEffect(async () => {
    setLoading(true);
    const errorHandler = (e: any) => setError(e);
    const [a, s] = await Promise.all([
      articleDetailApi({ url: "about" }, errorHandler),
      settingApi(),
    ]);
    setArticle(a);
    setSettings(s);
    setLoading(false);
  }, []);

  if (error) {
    throw error;
  }

  const siteName = settings?.site_name || "無處告別";
  const pageTitle = `${article?.title ? article.title + " - " : ""}${siteName}`;

  const entries = useMemo(
    () => (article ? parseTimeline(article.content) : []),
    [article],
  );

  if (!article?.id) {
    return <SkeletonArticle className={"h-200"} />;
  }

  return (
    <>
      <title>{pageTitle}</title>
      <meta name="description" content={settings?.site_desc || ""} />
      <meta name="keywords" content={settings?.site_keyword || ""} />
      <PageTransition loading={loading}>
        <div className="mb-[10px] mx-auto max-w-3xl">
          {/* 标题区域 */}
          <header className="mb-10 text-center">
            <h1 className="text-2xl font-bold text-foreground">{article.title}</h1>
            <div className="mt-3 mx-auto h-px w-20 bg-gradient-to-r from-transparent via-primary/50 to-transparent" />
            <div className="mt-3 flex items-center justify-center gap-2 text-xs text-muted-foreground">
              <span>{article.user?.nick_name}</span>
              <span>·</span>
              <span>{article.created_at}</span>
            </div>
          </header>

          {/* 时间轴 */}
          <div className="relative">
            {entries.map((entry, index) => {
              const isFirst = index === 0;
              const isLast = index === entries.length - 1;

              return (
                <div key={index} className="relative flex gap-5 group">
                  {/* 时间轴标记列：圆点 + 连接线 */}
                  <div className="relative flex flex-col items-center flex-shrink-0 w-4">
                    <div
                      className={`relative z-10 rounded-full bg-primary ring-4 ring-white transition-all duration-300 group-hover:scale-125 ${
                        isFirst
                          ? "w-4 h-4 mt-1.5"
                          : isLast
                            ? "w-4 h-4 mt-1.5 animate-pulse"
                            : "w-2.5 h-2.5 mt-2.5"
                      }`}
                    />
                    {!isLast && (
                      <div className="flex-1 w-[2px] bg-gradient-to-b from-primary/25 to-primary/10 mt-2" />
                    )}
                  </div>

                  {/* 内容区域 */}
                  <div className={`flex-1 ${isLast ? "pb-2" : "pb-7"}`}>
                    <div className="inline-block text-sm font-semibold text-primary bg-primary/5 px-3 py-0.5 rounded-full mb-2 transition-colors duration-300 group-hover:bg-primary/10">
                      {entry.date}
                    </div>
                    <p className="text-foreground/80 leading-[1.8] transition-colors duration-300 group-hover:text-foreground">
                      {entry.description}
                    </p>
                  </div>
                </div>
              );
            })}

            {/* 结尾标记 */}
            <div className="relative flex gap-5">
              <div className="relative flex flex-col items-center flex-shrink-0 w-4">
                <div className="w-3 h-3 rounded-full border-2 border-primary/30 bg-white mt-2" />
              </div>
              <div className="flex-1 pt-1.5">
                <p className="text-sm text-muted-foreground/60 italic">未完待续...</p>
              </div>
            </div>
          </div>
        </div>
        <Comment postId={article.id} />
      </PageTransition>
    </>
  );
}
