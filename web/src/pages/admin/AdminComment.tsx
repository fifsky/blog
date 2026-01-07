import { useEffect, useState } from "react";
import { commentAdminListApi, commentDeleteApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
import { CTable, Column } from "@/components/CTable";
import { Button } from "@/components/ui/button";
import { CommentItem } from "@/types/openapi";

export default function AdminComment() {
  const [list, setList] = useState<CommentItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await commentAdminListApi({ page });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await commentDeleteApi({ id });
      loadList();
    }
  };
  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<CommentItem>[] = [
    {
      title: <div style={{ width: 20 }}>&nbsp;</div>,
      key: "id",
      render: (_, record) => (
        <input type="checkbox" name="ids" value={record.id} />
      )
    },
    {
      title: <div style={{ width: 150 }}>文章</div>,
      key: "article_title",
      render: (value, record) => (
        <a
          href={`${record.type === 2 ? record.url : "/article" + record.id}#comments`}
          target="_blank"
          rel="noreferrer"
        >
          {value}
        </a>
      )
    },
    {
      title: <div style={{ width: 60 }}>昵称</div>,
      key: "name"
    },
    {
      title: "评论",
      key: "content"
    },
    {
      title: <div style={{ width: 80 }}>IP</div>,
      key: "ip"
    },
    {
      title: <div style={{ width: 130 }}>日期</div>,
      key: "created_at",
      render: (value) => (
        <>{new Date(value).toLocaleString()}</>
      )
    },
    {
      title: <div style={{ width: 80 }}>操作</div>,
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
      )
    }
  ];

  return (
    <div>
      <h2 className="border-b border-b-[#cccccc] text-base">管理评论</h2>
      <div className="my-[10px] flex items-center">
        <BatchHandle />
      </div>
      {/* 使用自定义表格组件 */}
      <CTable data={list} columns={columns} />
      <div className="my-[10px] flex items-center justify-between">
        <BatchHandle />
        <Paginate page={page} pageTotal={pageTotal} onChange={setPage} />
      </div>
    </div>
  );
}
