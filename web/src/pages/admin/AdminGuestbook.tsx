import { useEffect, useState } from "react";
import { guestbookDeleteApi, guestbookListApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Pagination } from "@/components/Pagination";
import { Button } from "@/components/ui/button";
import { CTable, Column } from "@/components/CTable";
import { GuestbookItem } from "@/types/openapi";
import { dialog } from "@/utils/dialog";

export default function AdminGuestbook() {
  const [list, setList] = useState<GuestbookItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const loadList = async () => {
    const ret = await guestbookListApi({ page });
    setList(ret.list || []);
    setTotal(ret.total || 0);
  };

  const deleteItem = (id: number) => {
    dialog.confirm("确认要删除这条留言？", {
      onOk: async () => {
        await guestbookDeleteApi({ id });
        loadList();
      },
    });
  };

  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<GuestbookItem>[] = [
    {
      title: <div style={{ width: 20 }}></div>,
      key: "id",
      render: (_, record) => <input type="checkbox" name="ids" value={record.id} />,
    },
    {
      title: <div style={{ width: 80 }}>昵称</div>,
      key: "name",
    },
    {
      title: "内容",
      key: "content",
    },
    {
      title: <div style={{ width: 120 }}>IP</div>,
      key: "ip",
    },
    {
      title: <div style={{ width: 150 }}>日期</div>,
      key: "created_at",
    },
    {
      title: <div style={{ width: 60 }}>操作</div>,
      key: "id",
      render: (_, record) => (
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
      ),
    },
  ];

  return (
    <div>
      <title>管理留言 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">管理留言</h2>
      <div className="my-[10px] flex items-center">
        <BatchHandle />
      </div>
      <CTable data={list} columns={columns} />
      <div className="my-2.5 flex items-center justify-between">
        <BatchHandle />
        <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
      </div>
    </div>
  );
}
