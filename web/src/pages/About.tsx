import { useMemo, useState } from "react";
import { CalendarDays, Clock } from "lucide-react";
import { Comment } from "@/components/Comment";
import { PageTransition } from "@/components/PageTransition";
import { SkeletonArticle } from "@/components/Skeleton";
import { articleDetailApi, settingApi } from "@/service";
import { ArticleItem, Setting } from "@/types/openapi";
import { useAsyncEffect } from "@/hooks";
import { cn } from "@/lib/utils";

interface TimelineEntry {
  date: string;
  description: string;
}

interface TimelineDisplayEntry extends TimelineEntry {
  year: string;
  detail: string;
  sequence: string;
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

/** 拆分日期用于时间轴的年份和短日期展示 */
function splitTimelineDate(date: string): Pick<TimelineDisplayEntry, "year" | "detail"> {
  const match = date.match(/^(\d{4})(.*)$/);
  if (!match) {
    return { year: date, detail: "纪事" };
  }

  const [, year, rest] = match;
  const detail = rest
    .trim()
    .replace(/^[年./-]+/, "")
    .replace(/[/-]/g, ".")
    .replace(/年/g, ".")
    .replace(/月/g, ".")
    .replace(/日/g, "")
    .replace(/\.+/g, ".")
    .replace(/^\./, "")
    .replace(/\.$/, "");

  return { year, detail: detail || "纪事" };
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

  const entries = useMemo<TimelineDisplayEntry[]>(
    () =>
      article
        ? parseTimeline(article.content).map((entry, index) => ({
            ...entry,
            ...splitTimelineDate(entry.date),
            sequence: String(index + 1).padStart(2, "0"),
          }))
        : [],
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
        <div className="mx-auto mb-[10px] max-w-[720px]">
          {/* 标题区域 */}
          <header className="relative mb-9 min-h-[58px] overflow-hidden border-b border-primary/15 pb-7">
            <div className="pointer-events-none absolute right-0 top-0 hidden select-none text-[58px] font-black leading-none tracking-normal text-primary/5 sm:block">
              ABOUT
            </div>
          </header>

          {/* 时间轴 */}
          <div className="relative">
            <div
              className="absolute bottom-14 left-[70px] top-1 w-px bg-gradient-to-b from-primary/10 via-primary/35 to-primary/10 sm:left-1/2 sm:-translate-x-1/2"
              aria-hidden="true"
            />
            <ol aria-label="关于我的时间线">
              {entries.map((entry, index) => {
                const cardOnRight = index % 2 === 1;

                return (
                  <li
                    key={`${entry.date}-${index}`}
                    className="group relative grid grid-cols-[140px_minmax(0,1fr)] items-center gap-6 pb-6 sm:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] sm:gap-0"
                  >
                    {/* 日期节点 */}
                    <div className="relative z-10 col-start-1 row-start-1 flex justify-center sm:col-start-2">
                      <span className="inline-flex items-center gap-1 whitespace-nowrap rounded-full bg-muted px-2.5 py-1 text-xs font-semibold text-muted-foreground">
                        <CalendarDays className="size-3.5" aria-hidden="true" />
                        {entry.date}
                      </span>
                    </div>

                    {/* 内容区域 */}
                    <article
                      className={cn(
                        "col-start-2 row-start-1 min-w-0",
                        cardOnRight
                          ? "sm:col-start-3 sm:pl-6"
                          : "sm:col-start-1 sm:pr-6 sm:text-right",
                      )}
                    >
                      <p className="text-[15px] leading-[1.9] text-foreground/85 transition-colors duration-300 group-hover:text-foreground">
                        {entry.description}
                      </p>
                    </article>
                  </li>
                );
              })}
            </ol>

            {/* 结尾标记 */}
            <div className="relative grid grid-cols-[140px_minmax(0,1fr)] gap-6 sm:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] sm:gap-y-1">
              <div className="relative z-10 col-start-1 flex justify-center sm:col-start-2">
                <div className="flex size-10 items-center justify-center text-primary/60">
                  <Clock className="size-5" aria-hidden="true" />
                </div>
              </div>
              <div className="col-span-2 text-center sm:col-span-3">
                <p className="text-sm font-medium text-muted-foreground">生活还在继续...</p>
              </div>
            </div>
          </div>
        </div>
        <Comment postId={article.id} />
      </PageTransition>
    </>
  );
}
