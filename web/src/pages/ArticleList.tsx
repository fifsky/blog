import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router";
import { FileText } from "lucide-react";
import { CArticle } from "@/components/CArticle";
import { Pagination } from "@/components/Pagination";
import { Empty } from "@/components/Empty";
import { articleListApi, settingApi } from "@/service";
import { useStore } from "@/store/context";
import { ArticleItem, ArticleListRequest, Options } from "@/types/openapi";

export default function ArticleList() {
  const [list, setList] = useState<ArticleItem[]>();
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [settings, setSettings] = useState<Options>();
  const location = useLocation();
  const navigate = useNavigate();
  const params = useParams();
  const setKeyword = useStore((s) => s.setKeyword);
  const loadList = async () => {
    const q = new URLSearchParams(location.search);
    const currentPage = q.get("page") ? parseInt(q.get("page")!) : 1;
    setPage(currentPage);
    const data: ArticleListRequest = {
      ...params,
      keyword: q.get("keyword") || "",
      page: currentPage,
      type: 1,
    };
    const [ret, s] = await Promise.all([articleListApi(data), settingApi()]);
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
    setSettings(s);
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

  if (!list) return;
  const siteName = settings?.kv?.site_name || "無處告別";
  return (
    <div>
      <title>{siteName}</title>
      <meta name="description" content={settings?.kv?.site_desc || ""} />
      <meta name="keywords" content={settings?.kv?.site_keyword || ""} />
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
          <Pagination page={page} pageTotal={pageTotal} onChange={changePage} />
        </>
      )}
    </div>
  );
}
