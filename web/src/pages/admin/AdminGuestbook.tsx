import { useEffect, useState, useCallback } from "react";
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
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());

  const loadList = async () => {
    const ret = await guestbookListApi({ page });
    setList(ret.list || []);
    setTotal(ret.total || 0);
    setSelectedIds(new Set()); // 重置选择
  };

  const deleteItem = (id: number) => {
    dialog.confirm("确认要删除这条留言？", {
      onOk: async () => {
        await guestbookDeleteApi({ ids: [id] });
        loadList();
      },
    });
  };

  // 切换单个选中状态
  const toggleSelect = useCallback((id: number) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  // 全选
  const handleSelectAll = useCallback(() => {
    setSelectedIds(new Set(list.map((item) => item.id)));
  }, [list]);

  // 反选
  const handleInverseSelect = useCallback(() => {
    setSelectedIds((prev) => {
      const next = new Set<number>();
      list.forEach((item) => {
        if (!prev.has(item.id)) {
          next.add(item.id);
        }
      });
      return next;
    });
  }, [list]);

  // 批量操作
  const handleBatchOperation = useCallback(
    async (operation: string) => {
      if (selectedIds.size === 0) return;

      if (operation === "2") {
        // 删除
        dialog.confirm(`确认要删除选中的 ${selectedIds.size} 条留言？`, {
          onOk: async () => {
            await guestbookDeleteApi({ ids: Array.from(selectedIds) });
            loadList();
          },
        });
      }
    },
    [selectedIds],
  );

  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<GuestbookItem>[] = [
    {
      title: <div style={{ width: 20 }}></div>,
      key: "id",
      render: (_, record) => (
        <input
          type="checkbox"
          checked={selectedIds.has(record.id)}
          onChange={() => toggleSelect(record.id)}
        />
      ),
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
        <BatchHandle
          selectedCount={selectedIds.size}
          totalCount={list.length}
          onSelectAll={handleSelectAll}
          onInverseSelected={handleInverseSelect}
          onBatchOperation={handleBatchOperation}
        />
      </div>
      <CTable data={list} columns={columns} />
      <div className="my-2.5 flex items-center justify-between">
        <BatchHandle
          selectedCount={selectedIds.size}
          totalCount={list.length}
          onSelectAll={handleSelectAll}
          onInverseSelected={handleInverseSelect}
          onBatchOperation={handleBatchOperation}
        />
        <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
      </div>
    </div>
  );
}
