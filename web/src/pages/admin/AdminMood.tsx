import React, { useEffect, useState } from "react";
import {
  moodDeleteApi,
  moodListApi,
  moodCreateApi,
  moodUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";

export default function AdminMood() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await moodListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const editItem = (id: number) => {
    setItem(list.find((i) => i.id === id));
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await moodDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => setItem({});
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const { id, content } = item;
    if (id) await moodUpdateApi({ id, content });
    else await moodCreateApi({ content });
    setItem({});
    loadList();
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div className="clearfix">
      <h2>管理心情</h2>
      <div className="col-left">
        <div className="operate clearfix">
          <BatchHandle />
        </div>
        <table className="list">
          <tbody>
            <tr>
              <th style={{ width: 20 }}>&nbsp;</th>
              <th style={{ width: 80 }}>作者</th>
              <th>心情</th>
              <th style={{ width: 90 }}>日期</th>
              <th style={{ width: 80 }}>操作</th>
            </tr>
            {list.length === 0 && (
              <tr>
                <td colSpan={7} align="center">
                  还没有心情！
                </td>
              </tr>
            )}
            {list.length > 0 &&
              list.map((v) => (
                <tr key={v.id}>
                  <td>
                    <input type="checkbox" name="ids" value={v.id} />
                  </td>
                  <td>{v.user.name}</td>
                  <td>{v.content}</td>
                  <td>{v.created_at}</td>
                  <td>
                    <a
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        editItem(v.id);
                      }}
                    >
                      编辑
                    </a>
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
      <div className="col-right" style={{ width: 250, paddingTop: 31 }}>
        <form className="vf" method="post" autoComplete="off" onSubmit={submit}>
          <p>
            <label className="label_input">发表心情</label>
            <textarea
              name="content"
              rows={5}
              cols={30}
              value={item.content || ""}
              onChange={(e) =>
                setItem((prev: any) => ({ ...prev, content: e.target.value }))
              }
            ></textarea>
          </p>
          <p className="act">
            <button className="formbutton" type="submit">
              {item.id ? "修改" : "添加"}
            </button>
            {item.id && (
              <a
                className="ml10"
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  cancel();
                }}
              >
                取消
              </a>
            )}
          </p>
        </form>
      </div>
    </div>
  );
}
