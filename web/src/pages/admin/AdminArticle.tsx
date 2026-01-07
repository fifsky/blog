import { useEffect, useState } from "react";
import { Link } from "react-router";
import { articleDeleteApi, articleListApi } from "@/service";
import { BatchHandle } from "@/components/BatchHandle";
import { Paginate } from "@/components/Paginate";
import { CTable, Column } from "@/components/CTable";
import { Badge } from "@/components/ui/badge";
import { ArticleItem } from "@/types/openapi";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function AdminArticle() {
  const [list, setList] = useState<ArticleItem[]>([]);
  const [pageTotal, setPageTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await articleListApi({ page, type: 1 });
    setList(ret.list || []);
    setPageTotal(ret.page_total || 0);
  };
  const deleteItem = async (id: number) => {
    if (confirm("确认要删除？")) {
      await articleDeleteApi({ id });
      loadList();
    }
  };
  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<ArticleItem>[] = [
    {
      title: <div style={{ width: 20 }}>&nbsp;</div>,
      key: "id",
      render: (value) => (
        <input type="checkbox" name="ids" value={value} />
      )
    },
    {
      title: (
        <div style={{ width: 20 }}>
          <i className="iconfont icon-comment text-[12px]"></i>
        </div>
      ),
      key: "id",
      render: () => (
        <Badge variant="secondary">0</Badge>
      )
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
      )
    },
    {
      title: <div style={{ width: 60 }}>作者</div>,
      key: "user.nick_name"
    },
    {
      title: <div style={{ width: 80 }}>分类</div>,
      key: "cate.name",
      render: (value, record) => (
        <a
          href={`/category/${record.cate.domain}`}
          target="_blank"
          rel="noreferrer"
        >
          {value}
        </a>
      )
    },
    {
      title: <div style={{ width: 80 }}>类型</div>,
      key: "type",
      render: (value) => (
        value === 1 ? "文章" : "页面"
      )
    },
    {
      title: <div style={{ width: 180 }}>日期</div>,
      key: "updated_at"
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_,record) => (
        <>
          <Link to={`/admin/post/article?id=${record.id}`}>编辑</Link>
          <span className="px-1.5 text-[#ccc]">|</span>
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
        </>
      )
    }
  ];

  return (
    <div>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理文章
        <Link to="/admin/post/article" className="ml-3 text-[14px]">
          <i className="iconfont icon-edit" style={{ color: "#444" }}></i>写文章
        </Link>
      </h2>
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
