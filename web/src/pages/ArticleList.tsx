import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router";
import { CArticle } from "@/components/CArticle";
import { Pagination } from "@/components/Pagination";
import { articleListApi } from "@/service";
import { useStore } from "@/store/context";
import { usePrismHighlight } from "@/hooks";
import { ArticleItem, ArticleListRequest } from "@/types/openapi";

export default function ArticleList() {
  const [list, setList] = useState<ArticleItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
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
    const ret = await articleListApi(data);
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
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

  usePrismHighlight([list]);

  return (
    <div>
      {list.map((v, k) => (
        <div className="articles" key={k}>
          <CArticle article={v} />
          <div className="border-t border-t-dashed border-t-[#dbdbdb] mt-5 pt-2.5 pb-2.5 text-right"></div>
        </div>
      ))}
      <Pagination page={page} pageTotal={pageTotal} onChange={changePage} />
    </div>
  );
}
