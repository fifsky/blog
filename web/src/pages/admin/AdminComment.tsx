import { useEffect, useState, useCallback } from "react";
import { commentAdminListApi, commentDeleteApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Pagination } from "@/components/Pagination";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { CTable, Column } from "@/components/CTable";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { AdminCommentItem } from "@/types/openapi";
import { dialog } from "@/utils/dialog";

// 生成文章访问路径：有自定义路径则使用 /${url}，否则 /article/${id}
function postLink(item: AdminCommentItem): string {
  return item.post_url ? `/${item.post_url}` : `/article/${item.post_id}`;
}

export default function AdminComment() {
  const [list, setList] = useState<AdminCommentItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  // 详情弹框当前查看的评论
  const [detailItem, setDetailItem] = useState<AdminCommentItem | null>(null);

  const loadList = async () => {
    const ret = await commentAdminListApi({ page, keyword });
    setList(ret.list || []);
    setTotal(ret.total || 0);
    setSelectedIds(new Set()); // 重置选择
  };

  const deleteItem = (id: number) => {
    dialog.confirm("确认要删除这条评论？", {
      onOk: async () => {
        await commentDeleteApi({ ids: [id] });
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
        dialog.confirm(`确认要删除选中的 ${selectedIds.size} 条评论？`, {
          onOk: async () => {
            await commentDeleteApi({ ids: Array.from(selectedIds) });
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

  const onSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      setPage(1);
      loadList();
    }
  };

  // 定义表格列配置
  const columns: Column<AdminCommentItem>[] = [
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
      title: <div style={{ width: 100 }}>所属文章</div>,
      key: "post_title",
      render: (_, record) =>
        record.post_title ? (
          <a href={postLink(record)} target="_blank" rel="noreferrer" className="hover:underline">
            {record.post_title}
          </a>
        ) : (
          "-"
        ),
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
      title: <div style={{ width: 100 }}>操作</div>,
      key: "id",
      render: (_, record) => (
        <>
          <Button
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              setDetailItem(record);
            }}
          >
            详情
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
      <title>管理评论 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">管理评论</h2>
      <div className="my-[10px] flex items-center justify-between">
        <BatchHandle
          selectedCount={selectedIds.size}
          totalCount={list.length}
          onSelectAll={handleSelectAll}
          onInverseSelected={handleInverseSelect}
          onBatchOperation={handleBatchOperation}
        />
        <Input
          placeholder="搜索昵称或内容，回车搜索"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onKeyDown={onSearch}
          className="max-w-[280px] h-8"
        />
      </div>
      <CTable data={list} columns={columns} />
      <div className="my-2.5 flex items-center justify-end">
        <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
      </div>

      {/* 评论详情弹框 */}
      <Dialog open={!!detailItem} onOpenChange={(open) => !open && setDetailItem(null)}>
        <DialogContent className="sm:max-w-[520px]">
          <DialogHeader>
            <DialogTitle>评论详情</DialogTitle>
          </DialogHeader>
          {detailItem && (
            <div className="space-y-3 text-sm">
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">评论ID</span>
                <span className="text-[#374151]">{detailItem.id}</span>
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">昵称</span>
                <span className="text-[#374151]">{detailItem.name}</span>
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">邮箱</span>
                <span className="text-[#374151] break-all">{detailItem.email || "-"}</span>
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">网址</span>
                {detailItem.website ? (
                  <a
                    href={detailItem.website}
                    target="_blank"
                    rel="noreferrer"
                    className="text-[#0066cc] hover:underline break-all"
                  >
                    {detailItem.website}
                  </a>
                ) : (
                  <span className="text-[#374151]">-</span>
                )}
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">IP</span>
                <span className="text-[#374151]">{detailItem.ip || "-"}</span>
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">时间</span>
                <span className="text-[#374151]">{detailItem.created_at}</span>
              </div>
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">所属文章</span>
                {detailItem.post_title ? (
                  <a
                    href={postLink(detailItem)}
                    target="_blank"
                    rel="noreferrer"
                    className="text-[#0066cc] hover:underline"
                  >
                    {detailItem.post_title}
                  </a>
                ) : (
                  <span className="text-[#374151]">-</span>
                )}
              </div>
              {detailItem.reply_name && (
                <div className="flex gap-2">
                  <span className="w-20 shrink-0 text-[#9ca3af]">回复对象</span>
                  <span className="text-[#374151]">@{detailItem.reply_name}</span>
                </div>
              )}
              <div className="flex gap-2">
                <span className="w-20 shrink-0 text-[#9ca3af]">内容</span>
                <span className="text-[#374151] whitespace-pre-wrap break-words flex-1">
                  {detailItem.content}
                </span>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
