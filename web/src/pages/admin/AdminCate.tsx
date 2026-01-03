import React, { useEffect, useState } from "react";
import {
  cateDeleteApi,
  cateListApi,
  cateCreateApi,
  cateUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";

export default function AdminCate() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const loadList = async () => {
    const ret = await cateListApi({});
    setList(ret.list || []);
  };
  const editItem = (id: number) => setItem(list.find((i) => i.id === id));
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await cateDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => setItem({});
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const { id, name, domain, desc } = item;
    if (id) await cateUpdateApi({ id, name, domain, desc });
    else await cateCreateApi({ name, domain, desc });
    setItem({});
    loadList();
  };
  useEffect(() => {
    loadList();
  }, []);
  return (
    <div className="clearfix">
      <h2>管理分类</h2>
      <div className="col-left">
        <div className="operate clearfix">
          <BatchHandle />
        </div>
        <table className="list">
          <tbody>
            <tr>
              <th style={{ width: 20 }}>&nbsp;</th>
              <th>分类名</th>
              <th style={{ width: 90 }}>缩略名</th>
              <th style={{ width: 60 }}>文章数</th>
              <th style={{ width: 90 }}>操作</th>
            </tr>
            {list.length === 0 && (
              <tr>
                <td colSpan={7} align="center">
                  还没有分类！
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
                  <td>{v.domain}</td>
                  <td className="art-num">{v.num}</td>
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
        </div>
      </div>
      <div className="col-right" style={{ width: 250, paddingTop: 31 }}>
        <form className="vf" method="post" autoComplete="off" onSubmit={submit}>
          <p>
            <label className="label_input">分类名称</label>
            <input
              type="text"
              className="input_text"
              size={30}
              name="name"
              value={item.name || ""}
              onChange={(e) =>
                setItem((prev: any) => ({ ...prev, name: e.target.value }))
              }
            />
          </p>
          <p>
            <label className="label_input">分类缩略名</label>
            <input
              type="text"
              className="input_text"
              size={30}
              name="domain"
              value={item.domain || ""}
              onChange={(e) =>
                setItem((prev: any) => ({ ...prev, domain: e.target.value }))
              }
            />
            <span className="hint">缩略名，使用字母开头([a-z][0-9]-)</span>
          </p>
          <p>
            <label className="label_input">分类描述</label>
            <textarea
              name="desc"
              rows={5}
              cols={30}
              value={item.desc || ""}
              onChange={(e) =>
                setItem((prev: any) => ({ ...prev, desc: e.target.value }))
              }
            ></textarea>
            <span className="hint">描述将在分类meta中显示</span>
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
