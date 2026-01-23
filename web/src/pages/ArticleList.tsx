import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router";
import { FileText } from "lucide-react";
import { CArticle } from "@/components/CArticle";
import { Pagination } from "@/components/Pagination";
import { Empty } from "@/components/Empty";
import { PageTransition } from "@/components/PageTransition";
import { SkeletonArticleList } from "@/components/Skeleton";
import { articleListApi, settingApi } from "@/service";
import { useStore } from "@/store/context";
import { ArticleItem, ArticleListRequest, Options } from "@/types/openapi";

export default function ArticleList() {
  const [list, setList] = useState<ArticleItem[]>();
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [settings, setSettings] = useState<Options>();
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const navigate = useNavigate();
  const params = useParams();
  const setKeyword = useStore((s) => s.setKeyword);

  const loadList = async () => {
    setLoading(true);
    const q = new URLSearchParams(location.search);
    const currentPage = q.get("page") ? parseInt(q.get("page")!) : 1;
    setPage(currentPage);
    const data: ArticleListRequest = {
      ...params,
      keyword: q.get("keyword") || "",
      tag: q.get("tag") || "",
      page: currentPage,
      type: 1,
    };
    const [ret, s] = await Promise.all([articleListApi(data), settingApi()]);
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
    if (location.pathname !== "/search") {
      setKeyword("");
    }
    loadList();
  }, [location.pathname, location.search]);

  const siteName = settings?.kv?.site_name || "無處告別";

  // Show skeleton during initial load
  if (!list) {
    return <SkeletonArticleList />;
  }

  return (
    <div>
      <title>{siteName}</title>
      <meta name="description" content={settings?.kv?.site_desc || ""} />
      <meta name="keywords" content={settings?.kv?.site_keyword || ""} />
      <PageTransition loading={loading}>
        {list.length === 0 ? (
          <Empty icon={<FileText />} title="暂无文章" content="当前没有可显示的文章内容" />
        ) : (
          <>
            {list.map((v, k) => (
              <div className="articles" key={k}>
                <CArticle article={v} />
                <div className="border-t border-dashed border-t-[#dbdbdb] mt-5 pt-2.5 pb-2.5 text-right"></div>
              </div>
            ))}
            <Pagination
              page={page}
              total={total}
              pageSize={parseInt(settings?.kv?.post_num || "10") || 10}
              onChange={changePage}
            />
          </>
        )}
      </PageTransition>
    </div>
  );
}
