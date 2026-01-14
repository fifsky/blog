import { useEffect, useState } from "react";
import { Link } from "react-router";
import { articleDeleteApi, articleListAdminApi, articleRestoreApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Pagination } from "@/components/Pagination";
import { CTable, Column } from "@/components/CTable";
import { Badge } from "@/components/ui/badge";
import { ArticleItem } from "@/types/openapi";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const STATUS_MAP: Record<
  number,
  { label: string; variant: "default" | "secondary" | "destructive" | "outline" }
> = {
  1: { label: "已发布", variant: "default" },
  2: { label: "已删除", variant: "destructive" },
  3: { label: "草稿", variant: "secondary" },
};

export default function AdminArticle() {
  const [list, setList] = useState<ArticleItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set());
  const [statusFilter, setStatusFilter] = useState<number | undefined>(undefined);
  const [batchLoading, setBatchLoading] = useState(false);

  const loadList = async () => {
    const ret = await articleListAdminApi({ page, type: 1, status: statusFilter });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await articleDeleteApi({ id });
      loadList();
    }
  };
  const restoreItem = async (id: number) => {
    if (confirm("确认要恢复为草稿？")) {
      await articleRestoreApi({ id });
      loadList();
    }
  };

  // 批量删除
  const batchDelete = async () => {
    if (selectedIds.size === 0) return;
    if (!confirm(`确认要删除选中的 ${selectedIds.size} 篇文章？`)) return;

    setBatchLoading(true);
    try {
      await articleDeleteApi({ ids: Array.from(selectedIds) });
      setSelectedIds(new Set());
      await loadList();
    } finally {
      setBatchLoading(false);
    }
  };

  // 批量操作处理
  const handleBatchOperation = (operation: string) => {
    if (operation === "2") {
      // 删除
      batchDelete();
    } else if (operation === "1") {
      // 置顶 - 暂未实现
      alert("置顶功能暂未实现");
    }
  };

  const handleSelectAll = () => {
    if (selectedIds.size === list.length && list.length > 0) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(list.map((item) => item.id)));
    }
  };

  const handleInverseSelected = () => {
    const allIds = new Set(list.map((item) => item.id));
    const newSelectedIds = new Set<number>();
    allIds.forEach((id) => {
      if (!selectedIds.has(id)) {
        newSelectedIds.add(id);
      }
    });
    setSelectedIds(newSelectedIds);
  };

  const handleToggleSelect = (id: number) => {
    const newSelectedIds = new Set(selectedIds);
    if (newSelectedIds.has(id)) {
      newSelectedIds.delete(id);
    } else {
      newSelectedIds.add(id);
    }
    setSelectedIds(newSelectedIds);
  };

  useEffect(() => {
    loadList();
  }, [page, statusFilter]);

  useEffect(() => {
    setSelectedIds(new Set());
  }, [page]);

  // 定义表格列配置
  const columns: Column<ArticleItem>[] = [
    {
      title: (
        <input
          type="checkbox"
          checked={selectedIds.size === list.length && list.length > 0}
          onChange={handleSelectAll}
        />
      ),
      key: "id",
      render: (value) => (
        <input
          type="checkbox"
          checked={selectedIds.has(value)}
          onChange={() => handleToggleSelect(value)}
        />
      ),
    },
    {
      title: (
        <div style={{ width: 20 }}>
          <i className="iconfont icon-comment text-[12px]"></i>
        </div>
      ),
      key: "id",
      render: () => <Badge variant="secondary">0</Badge>,
    },
    {
      title: "标题",
      key: "title",
      render: (value, record) => (
        <a
          href={record.type === 2 ? record.url : "/article/" + record.id}
          target="_blank"
          rel="noreferrer"
        >
          {value}
        </a>
      ),
    },
    {
      title: <div style={{ width: 60 }}>作者</div>,
      key: "user.nick_name",
    },
    {
      title: <div style={{ width: 80 }}>分类</div>,
      key: "cate.name",
      render: (value, record) => (
        <a href={`/category/${record.cate.domain}`} target="_blank" rel="noreferrer">
          {value}
        </a>
      ),
    },
    {
      title: <div style={{ width: 80 }}>类型</div>,
      key: "type",
      render: (value) => (value === 1 ? "文章" : "页面"),
    },
    {
      title: <div style={{ width: 80 }}>状态</div>,
      key: "status",
      render: (value) => {
        const statusInfo = STATUS_MAP[value] || STATUS_MAP[1];
        return <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>;
      },
    },
    {
      title: <div style={{ width: 180 }}>日期</div>,
      key: "updated_at",
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_, record) => (
        <>
          {record.status !== 2 && (
            <>
              <Link to={`/admin/post/article?id=${record.id}`}>编辑</Link>
              <span className="px-1.5 text-[#ccc]">|</span>
            </>
          )}
          {record.status === 2 ? (
            <Button
              variant={"link"}
              className={cn("p-0 m-0 h-auto text-[13px]")}
              onClick={async () => {
                if (confirm("确认要恢复为草稿？")) {
                  await restoreItem(record.id);
                  loadList();
                }
              }}
            >
              恢复
            </Button>
          ) : (
            <Button
              variant={"link"}
              className={cn("p-0 m-0 h-auto text-[13px]")}
              onClick={(e) => {
                e.preventDefault();
                deleteItem(record.id);
              }}
            >
              删除
            </Button>
          )}
        </>
      ),
    },
  ];

  return (
    <div>
      <title>管理文章 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理文章
        <Link to="/admin/post/article" className="ml-3 text-[14px]">
          <i className="iconfont icon-edit" style={{ color: "#444" }}></i>写文章
        </Link>
      </h2>
      <div className="my-[10px] flex items-center justify-between">
        <BatchHandle
          selectedCount={selectedIds.size}
          totalCount={list.length}
          onSelectAll={handleSelectAll}
          onInverseSelected={handleInverseSelected}
          onBatchOperation={handleBatchOperation}
          disabled={batchLoading}
        />
        <Select
          value={statusFilter?.toString() || "all"}
          onValueChange={(value) => setStatusFilter(value === "all" ? undefined : parseInt(value))}
        >
          <SelectTrigger className="w-[120px]" size={"sm"}>
            <SelectValue placeholder="全部状态" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部状态</SelectItem>
            <SelectItem value="1">已发布</SelectItem>
            <SelectItem value="3">草稿</SelectItem>
            <SelectItem value="2">已删除</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* 使用自定义表格组件 */}
      <CTable data={list} columns={columns} />

      <div className="my-[10px] flex items-center justify-between">
        <BatchHandle
          selectedCount={selectedIds.size}
          totalCount={list.length}
          onSelectAll={handleSelectAll}
          onInverseSelected={handleInverseSelected}
          onBatchOperation={handleBatchOperation}
          disabled={batchLoading}
        />
        <Pagination page={page} pageTotal={pageTotal} onChange={setPage} />
      </div>
    </div>
  );
}
