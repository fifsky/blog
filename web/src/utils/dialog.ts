import { toast } from "sonner";

export const dialog = {
  message(msg: string, fn: ((ok?: boolean) => void) | null = null) {
    toast.error(msg);
    fn && fn();
  },
  confirm(msg: string, fn: ((ok: boolean) => void) | null = null) {
    const ok = window.confirm(msg);
    fn && fn(ok);
  },
};
