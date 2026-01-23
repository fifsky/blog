import { useEffect, useState } from "react";
import { photoDeleteApi, photoListApi, photoCreateApi, photoUpdateApi } from "@/service";
import { Pagination } from "@/components/Pagination";
import { Button } from "@/components/ui/button";
import useDialog from "@/hooks/useDialog";
import { AdminPhotoDialog } from "@/components/AdminPhotoDialog";
import { PhotoItem } from "@/types/openapi";
import { CTable, Column } from "@/components/CTable";
import { cn } from "@/lib/utils";
import { dialog } from "@/utils/dialog";

export default function AdminPhoto() {
  const [list, setList] = useState<PhotoItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<PhotoItem>();
  const { isOpen, open: openDialog, close: closeDialog } = useDialog(false);

  const loadList = async () => {
    const ret = await photoListApi({ page });
    setList(ret.list || []);
    setTotal(ret.total || 0);
  };

  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it);
    openDialog();
  };

  const deleteItem = (id: number) => {
    dialog.confirm("确认要删除这张照片？删除后无法恢复", {
      onOk: async () => {
        await photoDeleteApi({ id });
        loadList();
      },
    });
  };

  const handleSubmit = async (values: any) => {
    const { id } = item || ({} as PhotoItem);
    if (id) {
      await photoUpdateApi({
        id,
        title: values.title,
        description: values.description,
        province: values.province,
        city: values.city,
      });
    } else {
      await photoCreateApi({
        title: values.title,
        description: values.description,
        srcs: values.srcs,
        province: values.province,
        city: values.city,
      });
    }
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
  const columns: Column<PhotoItem>[] = [
    {
      title: <div style={{ width: 120 }}>标题</div>,
      key: "title",
    },
    {
      title: <div style={{ width: 200 }}>描述</div>,
      key: "description",
      render: (value) => <span className="line-clamp-2 text-gray-600">{value || "-"}</span>,
    },
    {
      title: <div style={{ width: 80 }}>照片</div>,
      key: "thumbnail",
      render: (value) => (
        <a href={value?.replace("!photothumb", "")} target="_blank" rel="noopener noreferrer">
          <img
            src={value}
            alt="照片"
            className="w-16 h-12 object-cover rounded cursor-pointer hover:opacity-80"
          />
        </a>
      ),
    },
    {
      title: <div style={{ width: 80 }}>省份</div>,
      key: "province_name",
      render: (value) => <span>{value || "-"}</span>,
    },
    {
      title: <div style={{ width: 80 }}>城市</div>,
      key: "city_name",
      render: (value) => <span>{value || "-"}</span>,
    },
    {
      title: <div style={{ width: 140 }}>创建时间</div>,
      key: "created_at",
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
      <title>管理相册 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理相册
        <Button
          variant={"link"}
          onClick={handleOpenDialog}
          className={cn("p-0 m-0 ml-3 leading-[21px] h-auto text-[14px] gap-0")}
        >
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          <span>新增照片</span>
        </Button>
      </h2>
      <div className="flex">
        <div className="w-full">
          <div className="my-[10px]">
            {/* 使用自定义表格组件 */}
            <CTable data={list} columns={columns} />
          </div>
          <div className="my-[10px] flex items-center justify-end">
            <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
          </div>
        </div>
      </div>

      {/* 照片对话框 */}
      <AdminPhotoDialog isOpen={isOpen} onClose={closeDialog} item={item} onSubmit={handleSubmit} />
    </div>
  );
}
