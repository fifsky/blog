import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { articleDeleteApi, articleListApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";

export default function AdminArticle() {
  const [list, setList] = useState<any[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await articleListApi({ page, type: 1 });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await articleDeleteApi({ id });
      loadList();
    }
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div id="articles">
      <h2>
        管理文章
        <Link to="/admin/post/article" className="add">
          <i className="iconfont icon-edit" style={{ color: "#444" }}></i>写文章
        </Link>
      </h2>
      <div className="operate clearfix">
        <BatchHandle />
      </div>
      <table className="list">
        <tbody>
          <tr>
            <th style={{ width: 20 }}>&nbsp;</th>
            <th style={{ width: 20 }}>
              <i className="iconfont icon-comment fs-12"></i>
            </th>
            <th>标题</th>
            <th style={{ width: 60 }}>作者</th>
            <th style={{ width: 80 }}>分类</th>
            <th style={{ width: 80 }}>类型</th>
            <th style={{ width: 90 }}>日期</th>
            <th style={{ width: 80 }}>操作</th>
          </tr>
          {list.length === 0 && (
            <tr>
              <td colSpan={7} align="center">
                还没有文章，来 <Link to="/admin/post/article">创建一篇</Link>{" "}
                文章吧！
              </td>
            </tr>
          )}
          {list.length > 0 &&
            list.map((v) => (
              <tr key={v.id}>
                <td>
                  <input type="checkbox" name="ids" value={v.id} />
                </td>
                <td className="comment-num">
                  <a
                    href={`${
                      v.type === 2 ? v.url : "/article" + v.id
                    }#comments`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {v.comment_num}
                  </a>
                </td>
                <td>
                  <a
                    href={v.type === 2 ? v.url : "/article/" + v.id}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {v.title}
                  </a>
                </td>
                <td>{v.user.nick_name}</td>
                <td>
                  <a
                    href={`/category/${v.cate.domain}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {v.cate.name}
                  </a>
                </td>
                <td>{v.type === 1 ? "文章" : "页面"}</td>
                <td>{v.updated_at}</td>
                <td>
                  <Link to={`/admin/post/article?id=${v.id}`}>编辑</Link>
                  <span className="line">|</span>
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
      <div className="operate clearfix">
        <BatchHandle />
        <Paginate page={page} pageTotal={pageTotal} onChange={setPage} />
      </div>
    </div>
  );
}
