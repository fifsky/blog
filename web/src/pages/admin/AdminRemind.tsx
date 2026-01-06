import { useEffect, useState } from "react";
import {
  remindDeleteApi,
  remindListApi,
  remindCreateApi,
  remindUpdateApi,
} from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
import { Button } from "@/components/ui/button";
import useDialog from "@/hooks/useDialog";
import { AdminRemindDialog } from "@/components/AdminRemindDialog";
import { remindTimeFormat, remindType } from "@/utils/remind_date";
import { RemindItem } from "@/types/openapi";

export default function AdminRemind() {
  const [list, setList] = useState<RemindItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<RemindItem>({} as RemindItem);
  const { isOpen, open: openDialog, close: closeDialog } = useDialog(false);
  
  const loadList = async () => {
    const ret = await remindListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it || ({} as RemindItem));
    openDialog();
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await remindDeleteApi({ id });
      loadList();
    }
  };
  
  const handleSubmit = async (values: any) => {
    const { id } = item || ({} as RemindItem);
    const data = {
      id,
      type: Number(values.type),
      content: values.content,
      month: values.month,
      week: values.week,
      day: values.day,
      hour: values.hour,
      minute: values.minute,
    };
    if (id) await remindUpdateApi(data);
    else await remindCreateApi(data);
    setItem({} as RemindItem);
    loadList();
  };
  
  const handleOpenDialog = () => {
    setItem({} as RemindItem);
    openDialog();
  };
  useEffect(() => {
    loadList();
  }, [page]);
  return (
    <div>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理提醒
        <Button variant={"link"} onClick={handleOpenDialog}>
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          新增提醒
        </Button>
      </h2>
      <div className="flex">
        <div className="w-full">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          <table className="list">
            <tbody>
              <tr>
                <th style={{ width: 20 }}>&nbsp;</th>
                <th style={{ width: 80 }}>提醒类别</th>
                <th style={{ width: 180 }}>时间</th>
                <th>内容</th>
                <th style={{ width: 90 }}>操作</th>
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
                      <span className="px-1.5 text-[#ccc]">|</span>
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
      </div>
      
      {/* 提醒对话框 */}
      <AdminRemindDialog
        isOpen={isOpen}
        onClose={closeDialog}
        item={item}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
