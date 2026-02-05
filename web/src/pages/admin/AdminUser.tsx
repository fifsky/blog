import { useEffect, useState } from "react";
import { Link } from "react-router";
import { userListApi, userStatusApi } from "@/service";
import { Pagination } from "@/components/Pagination";
import { CTable, Column } from "@/components/CTable";
import { Button } from "@/components/ui/button";
import { UserItem } from "@/types/openapi";
import { dialog } from "@/utils/dialog";

export default function AdminUser() {
  const [list, setList] = useState<UserItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const loadList = async () => {
    const ret = await userListApi({ page });
    setList(ret.list || []);
    setTotal(ret.total || 0);
  };
  const deleteItem = (id: number) => {
    dialog.confirm("确认要操作？", {
      onOk: async () => {
        await userStatusApi({ id });
        loadList();
      },
    });
  };
  useEffect(() => {
    loadList();
  }, [page]);

  // 定义表格列配置
  const columns: Column<UserItem>[] = [
    {
      title: <div style={{ width: 120 }}>用户名</div>,
      key: "name",
    },
    {
      title: <div style={{ width: 120 }}>昵称</div>,
      key: "nick_name",
    },
    {
      title: "邮箱",
      key: "email",
    },
    {
      title: <div style={{ width: 100 }}>角色</div>,
      key: "type",
      render: (value) => <>{value === 1 ? "管理员" : "编辑"}</>,
    },
    {
      title: <div style={{ width: 100 }}>状态</div>,
      key: "status",
      render: (value) => <>{value === 1 ? "启用" : "停用"}</>,
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_, record) => (
        <>
          <Link to={`/admin/post/user?id=${record.id}`}>编辑</Link>
          <span className="px-1.5 text-[#ccc]">|</span>
          <Button
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              deleteItem(record.id);
            }}
          >
            {record.status === 1 ? "停用" : "启用"}
          </Button>
        </>
      ),
    },
  ];

  return (
    <div>
      <title>管理用户 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理用户
        <Link to="/admin/post/user" className="ml-3 text-[14px]">
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          新增用户
        </Link>
      </h2>
      <div className="w-full mt-3">
        <CTable data={list} columns={columns} />
        <div className="my-[10px] flex items-center justify-between">
          <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
        </div>
      </div>
    </div>
  );
}
