import { toast } from "sonner";
import { createRoot } from "react-dom/client";
import { createElement } from "react";
import { ConfirmDialog, type ConfirmOptions } from "@/components/ConfirmDialog";

export const dialog = {
  message(msg: string, fn: ((ok?: boolean) => void) | null = null) {
    toast.error(msg);
    fn?.();
  },

  confirm(
    msg: string,
    options?: {
      title?: string;
      confirmText?: string;
      cancelText?: string;
      onOk?: () => void;
      onCancel?: () => void;
    },
  ) {
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);

    const cleanup = () => {
      root.unmount();
      container.remove();
    };

    const dialogOptions: ConfirmOptions = {
      title: options?.title || "确认",
      description: msg,
      confirmText: options?.confirmText || "确定",
      cancelText: options?.cancelText || "取消",
      onConfirm: () => {
        options?.onOk?.();
        cleanup();
      },
      onCancel: () => {
        options?.onCancel?.();
        cleanup();
      },
    };

    root.render(createElement(ConfirmDialog, dialogOptions));
  },
};
