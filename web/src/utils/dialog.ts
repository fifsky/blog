export const dialog = {
  alert(msg: string, fn: ((ok?: boolean) => void) | null = null) {
    window.alert(msg)
    fn && fn()
  },
  confirm(msg: string, fn: ((ok: boolean) => void) | null = null) {
    const ok = window.confirm(msg)
    fn && fn(ok)
  }
}

