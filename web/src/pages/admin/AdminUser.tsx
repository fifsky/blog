import { useEffect, useState } from "react";
import { Link } from "react-router";
import { userListApi, userStatusApi, userGenerate2FAApi, userBind2FAApi } from "@/service";
import { Pagination } from "@/components/Pagination";
import { CTable, Column } from "@/components/CTable";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { UserItem } from "@/types/openapi";
import { dialog } from "@/utils/dialog";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { QRCodeCanvas } from "qrcode.react";

export default function AdminUser() {
  const [list, setList] = useState<UserItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [totpModal, setTotpModal] = useState(false);
  const [totpQrCode, setTotpQrCode] = useState("");
  const [totpSecret, setTotpSecret] = useState("");
  const [totpCode, setTotpCode] = useState("");
  const [currentUserId, setCurrentUserId] = useState(0);

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

  const openTotpModal = async (id: number) => {
    try {
      const res = await userGenerate2FAApi({ id });
      setCurrentUserId(id);
      setTotpSecret(res.secret);
      setTotpQrCode(res.qr_code_uri);
      setTotpCode("");
      setTotpModal(true);
    } catch (e: any) {
      dialog.message(e.message || "生成2FA失败");
    }
  };

  const submitTotp = async () => {
    if (!totpCode) {
      dialog.message("请输入验证码");
      return;
    }
    try {
      await userBind2FAApi({ id: currentUserId, secret: totpSecret, code: totpCode });
      dialog.message("绑定成功");
      setTotpModal(false);
      loadList();
    } catch (e: any) {
      dialog.message(e.message || "绑定失败");
    }
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
      title: <div style={{ width: 140 }}>操作</div>,
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
          <span className="px-1.5 text-[#ccc]">|</span>
          <Button
            variant={"link"}
            className="p-0 m-0 h-auto text-[13px]"
            onClick={(e) => {
              e.preventDefault();
              openTotpModal(record.id);
            }}
          >
            {record.has_totp ? "重置2FA" : "开启2FA"}
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
      <Dialog open={totpModal} onOpenChange={setTotpModal}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>绑定 2FA (双因素认证)</DialogTitle>
          </DialogHeader>
          <div className="flex flex-col items-center gap-4 py-4">
            <p className="text-sm text-gray-500">请使用身份验证器（如 Google Authenticator）扫描下方二维码</p>
            {totpQrCode && (
              <div className="p-2 bg-white rounded-lg border">
                <QRCodeCanvas value={totpQrCode} size={200} />
              </div>
            )}
            <div className="w-full max-w-xs space-y-2">
              <p className="text-sm font-medium">验证码</p>
              <Input
                placeholder="请输入 6 位验证码"
                value={totpCode}
                onChange={(e) => setTotpCode(e.target.value)}
                maxLength={6}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setTotpModal(false)}>
              取消
            </Button>
            <Button onClick={submitTotp}>确认绑定</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
