import { useEffect, useState } from "react";
import { Link } from "react-router";
import { userListApi, userStatusApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
export default function AdminUser() {
  const [list, setList] = useState<any[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await userListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要操作？")) {
      await userStatusApi({ id });
      loadList();
    }
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div>
      <h2>
        管理用户
        <Link to="/admin/post/user" className="add">
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          新增用户
        </Link>
      </h2>
      <div className="operate clearfix">
        <BatchHandle />
      </div>
      <table className="list">
        <tbody>
          <tr>
            <th style={{ width: 20 }}>&nbsp;</th>
            <th style={{ width: 80 }}>用户名</th>
            <th style={{ width: 80 }}>昵称</th>
            <th>邮箱</th>
            <th style={{ width: 60 }}>角色</th>
            <th style={{ width: 60 }}>状态</th>
            <th style={{ width: 90 }}>操作</th>
          </tr>
          {list.length === 0 && (
            <tr>
              <td colSpan={7} align="center">
                还没有用户！
              </td>
            </tr>
          )}
          {list.length > 0 &&
            list.map((v) => (
              <tr key={v.id}>
                <td>
                  <input type="checkbox" name="ids" value={v.id} />
                </td>
                <td>{v.name}</td>
                <td>{v.nick_name}</td>
                <td>{v.email}</td>
                <td>{v.type === 1 ? "管理员" : "编辑"}</td>
                <td>{v.status === 1 ? "启用" : "停用"}</td>
                <td>
                  <Link to={`/admin/post/user?id=${v.id}`}>编辑</Link>
                  <span className="line">|</span>
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      deleteItem(v.id);
                    }}
                  >
                    {v.status === 1 ? "停用" : "启用"}
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
