export const dialog = {
  message(msg: string, fn: ((ok?: boolean) => void) | null = null) {
    window.dispatchEvent(new CustomEvent("app-alert", { detail: { msg } }));
    fn && fn();
  },
  confirm(msg: string, fn: ((ok: boolean) => void) | null = null) {
    const ok = window.confirm(msg);
    fn && fn(ok);
  },
};
