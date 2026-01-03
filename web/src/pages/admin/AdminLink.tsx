import React, { useEffect, useState } from "react";
import {
  linkDeleteApi,
  linkListApi,
  linkCreateApi,
  linkUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";

export default function AdminLink() {
  const [list, setList] = useState<any[]>([]);
  const [item, setItem] = useState<any>({});
  const loadList = async () => {
    const ret = await linkListApi({});
    setList(ret.list || []);
  };
  const editItem = (id: number) => setItem(list.find((i) => i.id === id));
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await linkDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () => setItem({});
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const { id, name, url, desc } = item;
    if (id) await linkUpdateApi({ id, name, url, desc });
    else await linkCreateApi({ name, url, desc });
    setItem({});
    loadList();
  };
  useEffect(() => {
    loadList();
  }, []);
  return (
    <div>
      <h2>管理链接</h2>
      <div className="flex justify-between">
        <div className="w-[700px]">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          <table className="list">
            <tbody>
              <tr>
                <th style={{ width: 20 }}>&nbsp;</th>
                <th style={{ width: 120 }}>连接名</th>
                <th>地址</th>
                <th style={{ width: 90 }}>操作</th>
              </tr>
              {list.length === 0 && (
                <tr>
                  <td colSpan={7} align="center">
                    还没有链接！
                  </td>
                </tr>
              )}
              {list.length > 0 &&
                list.map((v) => (
                  <tr key={v.id}>
                    <td>
                      <input type="checkbox" name="ids" value={v.id} />
                    </td>
                    <td>
                      <a href={v.url} target="_blank" rel="noreferrer">
                        {v.name}
                      </a>
                    </td>
                    <td>
                      <a href={v.url}>{v.url}</a>
                    </td>
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
          <div className="my-[10px] flex items-center justify-between">
            <BatchHandle />
          </div>
        </div>
        <div className="w-[250px]" style={{ paddingTop: 31 }}>
          <form
            className="vf"
            method="post"
            autoComplete="off"
            onSubmit={submit}
          >
            <p>
              <label className="label_input">链接名称</label>
              <input
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
              <label className="label_input">链接地址</label>
              <input
                className="input_text"
                size={30}
                name="url"
                value={item.url || ""}
                onChange={(e) =>
                  setItem((prev: any) => ({ ...prev, url: e.target.value }))
                }
              />
              <span className="hint">例如：http://fifsky.com/</span>
            </p>
            <p>
              <label className="label_input">链接描述</label>
              <textarea
                name="desc"
                rows={5}
                cols={30}
                value={item.desc || ""}
                onChange={(e) =>
                  setItem((prev: any) => ({ ...prev, desc: e.target.value }))
                }
              ></textarea>
            </p>
            <p className="act">
              <button className="formbutton" type="submit">
                {item.id ? "修改" : "添加"}
              </button>
              {item.id && (
                <a
                  className="ml-2.5"
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
    </div>
  );
}
