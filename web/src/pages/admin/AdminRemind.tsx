import React, { useEffect, useState } from "react";
import {
  remindDeleteApi,
  remindListApi,
  remindCreateApi,
  remindUpdateApi,
} from "@/service";
import dayjs from "dayjs";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";

export default function AdminRemind() {
  const [list, setList] = useState<any[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<any>({
    type: 0,
    month: 1,
    week: 1,
    day: 1,
    hour: 0,
    minute: 0,
    content: "",
  });
  const remindType: Record<number, string> = {
    0: "固定",
    1: "每分钟",
    2: "每小时",
    3: "每天",
    4: "每周",
    5: "每月",
    6: "每年",
  };
  const monthFormat: Record<number, string> = {
    1: "01",
    2: "02",
    3: "03",
    4: "04",
    5: "05",
    6: "06",
    7: "07",
    8: "08",
    9: "09",
    10: "10",
    11: "11",
    12: "12",
  };
  const weekFormat: Record<number, string> = {
    1: "一",
    2: "二",
    3: "三",
    4: "四",
    5: "五",
    6: "六",
    7: "日",
  };
  const intRemindType = Number(item.type);
  const numFormat = (n: number) => (n < 10 ? "0" + n : String(n));
  const remindTimeFormat = (v: any) => {
    let str = "";
    switch (v.type) {
      case 0:
        str =
          dayjs(v.created_at).year() +
          "年" +
          monthFormat[v.month] +
          "月" +
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 3:
        str = numFormat(v.hour) + "时" + numFormat(v.minute) + "分";
        break;
      case 4:
        str =
          "周" +
          weekFormat[v.week] +
          " " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 5:
        str =
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
      case 6:
        str =
          monthFormat[v.month] +
          "月" +
          numFormat(v.day) +
          "日 " +
          numFormat(v.hour) +
          "时" +
          numFormat(v.minute) +
          "分";
        break;
    }
    return str;
  };
  const loadList = async () => {
    const ret = await remindListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.pageTotal || 0);
  };
  const editItem = (id: number) => setItem(list.find((i) => i.id === id));
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await remindDeleteApi({ id });
      loadList();
    }
  };
  const cancel = () =>
    setItem({
      type: 0,
      month: 1,
      week: 1,
      day: 1,
      hour: 0,
      minute: 0,
      content: "",
    });
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const { id, type, content, month, week, day, hour, minute } = item;
    const data = {
      id,
      type: Number(type),
      content,
      month,
      week,
      day,
      hour,
      minute,
    };
    if (id) await remindUpdateApi(data);
    else await remindCreateApi(data);
    cancel();
    loadList();
  };
  useEffect(() => {
    loadList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);
  return (
    <div className="clearfix">
      <h2>管理提醒</h2>
      <div className="col-left">
        <div className="operate clearfix">
          <BatchHandle />
        </div>
        <table className="list">
          <tbody>
            <tr>
              <th style={{ width: 20 }}>&nbsp;</th>
              <th style={{ width: 60 }}>提醒类别</th>
              <th style={{ width: 180 }}>时间</th>
              <th>内容</th>
              <th style={{ width: 80 }}>操作</th>
            </tr>
            {list.length === 0 && (
              <tr>
                <td colSpan={7} align="center">
                  还没有提醒！
                </td>
              </tr>
            )}
            {list.length > 0 &&
              list.map((v) => (
                <tr key={v.id}>
                  <td>
                    <input type="checkbox" name="ids" value={v.id} />
                  </td>
                  <td>{remindType[v.type]}</td>
                  <td>{remindTimeFormat(v)}</td>
                  <td>{v.content}</td>
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
            <label className="label_input">提醒类别</label>
            <select
              name="type"
              value={item.type}
              onChange={(e) =>
                setItem((prev: any) => ({
                  ...prev,
                  type: Number(e.target.value),
                }))
              }
            >
              {Object.entries(remindType).map(([k, v]) => (
                <option key={k} value={k}>
                  {v}
                </option>
              ))}
            </select>
          </p>
          <p>
            <label className="label_input">提醒时间</label>
            {[0, 6].includes(intRemindType) && (
              <select
                value={item.month}
                onChange={(e) =>
                  setItem((prev: any) => ({
                    ...prev,
                    month: Number(e.target.value),
                  }))
                }
              >
                {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
                  <option key={m} value={m}>
                    {monthFormat[m]}月
                  </option>
                ))}
              </select>
            )}
            {[4].includes(intRemindType) && (
              <select
                value={item.week}
                onChange={(e) =>
                  setItem((prev: any) => ({
                    ...prev,
                    week: Number(e.target.value),
                  }))
                }
              >
                {Array.from({ length: 7 }, (_, i) => i + 1).map((d) => (
                  <option key={d} value={d}>
                    周{weekFormat[d]}
                  </option>
                ))}
              </select>
            )}
            {[0, 5, 6].includes(intRemindType) && (
              <select
                value={item.day}
                onChange={(e) =>
                  setItem((prev: any) => ({
                    ...prev,
                    day: Number(e.target.value),
                  }))
                }
              >
                {Array.from({ length: 31 }, (_, i) => i + 1).map((d) => (
                  <option key={d} value={d}>
                    {numFormat(d)}日
                  </option>
                ))}
              </select>
            )}
            {[0, 3, 4, 5, 6].includes(intRemindType) && (
              <select
                value={item.hour}
                onChange={(e) =>
                  setItem((prev: any) => ({
                    ...prev,
                    hour: Number(e.target.value),
                  }))
                }
              >
                {Array.from({ length: 24 }, (_, i) => i).map((d) => (
                  <option key={d} value={d}>
                    {numFormat(d)}时
                  </option>
                ))}
              </select>
            )}
            {[0, 2, 3, 4, 5, 6].includes(intRemindType) && (
              <select
                value={item.minute}
                onChange={(e) =>
                  setItem((prev: any) => ({
                    ...prev,
                    minute: Number(e.target.value),
                  }))
                }
              >
                {Array.from({ length: 60 }, (_, i) => i).map((d) => (
                  <option key={d} value={d}>
                    {numFormat(d)}分
                  </option>
                ))}
              </select>
            )}
          </p>
          <p>
            <label className="label_input">提醒内容</label>
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
