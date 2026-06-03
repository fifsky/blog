import { useEffect, useState } from "react";
import { footprintDeleteApi, footprintListApi, footprintCreateApi, footprintUpdateApi } from "@/service";
import { Pagination } from "@/components/Pagination";
import { Button } from "@/components/ui/button";
import useDialog from "@/hooks/useDialog";
import { AdminFootprintDialog, type FormValues } from "@/components/AdminFootprintDialog";
import { FootprintItem } from "@/types/openapi";
import { CTable, Column } from "@/components/CTable";
import { cn } from "@/lib/utils";
import { dialog } from "@/utils/dialog";

export default function AdminFootprint() {
  const [list, setList] = useState<FootprintItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [item, setItem] = useState<FootprintItem>();
  const { isOpen, open: openDialog, close: closeDialog } = useDialog(false);

  const loadList = async () => {
    const ret = await footprintListApi({ page });
    setList(ret.list || []);
    setTotal(ret.total || 0);
  };

  const editItem = (id: number) => {
    const it = list.find((i) => i.id === id);
    setItem(it);
    openDialog();
  };

  const deleteItem = (id: number) => {
    dialog.confirm("确认要删除这条足迹？删除后无法恢复", {
      onOk: async () => {
        await footprintDeleteApi({ id });
        loadList();
      },
    });
  };

  const handleSubmit = async (values: FormValues) => {
    const { id } = item || ({} as FootprintItem);
    if (id) {
      await footprintUpdateApi({
        id,
        name: values.name,
        description: values.description,
        longitude: values.longitude,
        latitude: values.latitude,
        date: values.date,
        marker_color: values.marker_color,
        categories: values.categories,
        url: values.url,
        url_label: values.url_label,
        photo_urls: values.photo_urls,
      });
    } else {
      await footprintCreateApi({
        name: values.name,
        description: values.description,
        longitude: values.longitude,
        latitude: values.latitude,
        date: values.date,
        marker_color: values.marker_color,
        categories: values.categories,
        url: values.url,
        url_label: values.url_label,
        photo_urls: values.photo_urls,
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

  const columns: Column<FootprintItem>[] = [
    {
      title: <div style={{ width: 120 }}>地点名称</div>,
      key: "name",
    },
    {
      title: <div style={{ width: 100 }}>日期</div>,
      key: "date",
      render: (value) => <span>{value || "-"}</span>,
    },
    {
      title: <div style={{ width: 140 }}>坐标</div>,
      key: "longitude",
      render: (_, record) => (
        <span className="text-gray-500 text-xs">
          {record.longitude}, {record.latitude}
        </span>
      ),
    },
    {
      title: <div style={{ width: 120 }}>分类</div>,
      key: "categories",
      render: (value) => (
        <div className="flex gap-1 flex-wrap">
          {(value as string[] || []).map((c: string) => (
            <span key={c} className="px-1.5 py-0.5 bg-gray-100 rounded text-xs">
              {c}
            </span>
          ))}
        </div>
      ),
    },
    {
      title: <div style={{ width: 60 }}>照片</div>,
      key: "photos",
      render: (value) => {
        const photos = value as { src: string; thumbnail: string }[];
        return (
          <span className="text-gray-600">{photos?.length || 0} 张</span>
        );
      },
    },
    {
      title: <div style={{ width: 90 }}>操作</div>,
      key: "id",
      render: (_, record) => (
        <>
          <Button
            variant="link"
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
            variant="link"
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
      <title>管理足迹 - 無處告別</title>
      <h2 className="border-b border-b-[#cccccc] text-base">
        管理足迹
        <Button
          variant="link"
          onClick={handleOpenDialog}
          className={cn("p-0 m-0 ml-3 leading-[21px] h-auto text-[14px] gap-0")}
        >
          <i className="iconfont icon-add" style={{ color: "#444" }}></i>
          <span>新增足迹</span>
        </Button>
      </h2>
      <div className="w-full mt-3">
        <CTable data={list} columns={columns} />
        <div className="my-[10px] flex items-center justify-end">
          <Pagination page={page} total={total} pageSize={10} onChange={setPage} />
        </div>
      </div>

      <AdminFootprintDialog isOpen={isOpen} onClose={closeDialog} item={item} onSubmit={handleSubmit} />
    </div>
  );
}
