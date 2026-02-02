import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { Archive as ArchiveIcon } from "lucide-react";
import { Pagination } from "@/components/Pagination";
import { Empty } from "@/components/Empty";
import { PageTransition } from "@/components/PageTransition";
import { articleListApi, settingApi } from "@/service";
import { ArticleItem, Options } from "@/types/openapi";
import { SkeletonArchive } from "@/components/Skeleton";

type GroupedArticles = {
  yearMonth: string;
  articles: ArticleItem[];
};

export default function Archive() {
  const [list, setList] = useState<ArticleItem[]>();
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [settings, setSettings] = useState<Options>();
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const navigate = useNavigate();
  const pageSize = 50;

  const loadList = async () => {
    setLoading(true);
    const q = new URLSearchParams(location.search);
    const currentPage = q.get("page") ? parseInt(q.get("page")!) : 1;
    setPage(currentPage);
    const [ret, s] = await Promise.all([
      articleListApi({ page: currentPage, type: 1, page_size: pageSize }),
      settingApi(),
    ]);
    setList(ret.list || []);
    setSettings(s);
    setTotal(ret.total || 0);
    setLoading(false);
  };

  const changePage = (p: number) => {
    setPage(p);
    const q = new URLSearchParams(location.search);
    q.set("page", String(p));
    navigate({ pathname: location.pathname, search: q.toString() });
  };

  useEffect(() => {
    loadList();
  }, [location.search]);

  // Group articles by year-month
  const groupedArticles = (): GroupedArticles[] => {
    if (!list) return [];
    const groups: Record<string, ArticleItem[]> = {};
    list.forEach((article) => {
      const date = new Date(article.created_at);
      const yearMonth = `${date.getFullYear()}年${String(date.getMonth() + 1).padStart(2, "0")}月`;
      if (!groups[yearMonth]) {
        groups[yearMonth] = [];
      }
      groups[yearMonth].push(article);
    });
    return Object.entries(groups).map(([yearMonth, articles]) => ({
      yearMonth,
      articles,
    }));
  };

  // Show skeleton during initial load
  if (!list) {
    return <SkeletonArchive />;
  }
  const siteName = settings?.kv?.site_name || "無處告別";
  const pageTitle = `文章归档 - ${siteName}`;
  return (
    <div>
      <title>{pageTitle}</title>
      <meta name="description" content={settings?.kv?.site_desc || ""} />
      <meta name="keywords" content={settings?.kv?.site_keyword || ""} />

      <h2 className="text-xl font-bold mb-6">文章归档</h2>

      {list.length === 0 ? (
        <Empty icon={<ArchiveIcon />} title="暂无文章" content="当前没有可显示的文章内容" />
      ) : (
        <PageTransition loading={loading}>
          <div className="archive-timeline">
            {groupedArticles().map((group) => (
              <div key={group.yearMonth} className="mb-6">
                <h3 className="text-base font-bold mb-3">{group.yearMonth}</h3>
                <div className="pl-4 border-l-2 border-gray-200">
                  {group.articles.map((article) => (
                    <div key={article.id} className="mb-2 flex items-baseline gap-3">
                      <span className="text-gray-400 text-sm whitespace-nowrap">
                        {article.created_at.split(" ")[0]}
                      </span>
                      <span className="text-gray-400 text-sm">{article.cate?.name}</span>
                      <span className="text-gray-400">·</span>
                      <Link to={`/article/${article.id}`} className="transition-colors">
                        {article.title}
                      </Link>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
          <Pagination page={page} total={total} pageSize={pageSize} onChange={changePage} />
        </PageTransition>
      )}
    </div>
  );
}
