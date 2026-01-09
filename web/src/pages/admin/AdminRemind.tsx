import { useEffect, useState } from "react";
import { remindDeleteApi, remindListApi, remindCreateApi, remindUpdateApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Pagination } from "@/components/Pagination";
import { Button } from "@/components/ui/button";
import useDialog from "@/hooks/useDialog";
import { AdminRemindDialog } from "@/components/AdminRemindDialog";
import { remindTimeFormat, remindType } from "@/utils/remind_date";
import { RemindItem } from "@/types/openapi";
import { CTable, Column } from "@/components/CTable";
import { cn } from "@/lib/utils";

export default function AdminRemind() {
  const [list, setList] = useState<RemindItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<RemindItem>();
  const { isOpen, open: openDialog, close: closeDialog } = useDialog(false);

  const loadList = async () => {
    const ret = await remindListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it);
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
    setItem(undefined);
    loadList();
  };

  const handleOpenDialog = () => {
    setItem(undefined);
    openDialog();
  };
  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<RemindItem>[] = [
    {
      title: <div style={{ width: 20 }}>&nbsp;</div>,
      key: "id",
      render: (_, record) => <input type="checkbox" name="ids" value={record.id} />,
    },
    {
      title: <div style={{ width: 80 }}>提醒类别</div>,
      key: "type",
      render: (value) => <>{remindType[value as keyof typeof remindType]}</>,
    },
    {
      title: <div style={{ width: 180 }}>时间</div>,
      key: "id",
      render: (_, record) => <>{remindTimeFormat(record)}</>,
    },
    {
      title: "内容",
      key: "content",
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_, record) => (
        <>
          <Button
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              editItem(record.id);
            }}
          >
            编辑
          </Button>
          <span className="px-1.5 text-[#ccc]">|</span>
          <Button
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              deleteItem(record.id);
            }}
          >
            删除
          </Button>
        </>
      ),
    },
  ];

  return (
    <div>
      <title>管理提醒 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理提醒
        <Button
          variant={"link"}
          onClick={handleOpenDialog}
          className={cn("p-0 m-0 ml-3 leading-[21px] h-auto text-[14px] gap-0")}
        >
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          <span>新增提醒</span>
        </Button>
      </h2>
      <div className="flex">
        <div className="w-full">
          <div className="my-[10px] flex items-center">
            <BatchHandle />
          </div>
          {/* 使用自定义表格组件 */}
          <CTable data={list} columns={columns} />
          <div className="my-[10px] flex items-center justify-between">
            <BatchHandle />
            <Pagination page={page} pageTotal={pageTotal} onChange={setPage} />
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
