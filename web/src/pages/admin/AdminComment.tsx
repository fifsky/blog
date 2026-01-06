import { useEffect, useState } from "react";
import { commentAdminListApi, commentDeleteApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";

export default function AdminComment() {
  const [list, setList] = useState<any[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await commentAdminListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await commentDeleteApi({ id });
      loadList();
    }
  };
  useEffect(() => {
    loadList();
  }, [page]);
  return (
    <div>
      <h2>管理评论</h2>
      <div className="my-[10px] flex items-center">
        <BatchHandle />
      </div>
      <table className="list">
        <tbody>
          <tr>
            <th style={{ width: 20 }}>&nbsp;</th>
            <th style={{ width: 150 }}>文章</th>
            <th style={{ width: 60 }}>昵称</th>
            <th>评论</th>
            <th style={{ width: 80 }}>IP</th>
            <th style={{ width: 130 }}>日期</th>
            <th style={{ width: 80 }}>操作</th>
          </tr>
          {list.length === 0 && (
            <tr>
              <td colSpan={7} align="center">
                还没有评论！
              </td>
            </tr>
          )}
          {list.length > 0 &&
            list.map((v: any) => (
              <tr key={v.id}>
                <td>
                  <input type="checkbox" name="ids" value={v.id} />
                </td>
                <td>
                  <a
                    href={`${
                      v.type === 2 ? v.url : "/article" + v.id
                    }#comments`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {v.article_title}
                  </a>
                </td>
                <td>{v.name}</td>
                <td>{v.content}</td>
                <td>{v.ip}</td>
                <td>{new Date(v.created_at).toLocaleString()}</td>
                <td>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      deleteItem(v.id);
                    }}
                  >
                    删除
                  </a>
                </td>
              </tr>
            ))}
        </tbody>
      </table>
      <div className="my-[10px] flex items-center justify-between">
        <BatchHandle />
        <Paginate page={page} pageTotal={pageTotal} onChange={setPage} />
      </div>
    </div>
  );
}
