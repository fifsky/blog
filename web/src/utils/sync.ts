import { Err } from './error'
import { get } from './mapping'
import { dialog } from './dialog'

export const sync = async <T>(fn: () => Promise<T>, options: { errHandle?: (e: any) => void } = {}) => {
  try {
    return await fn()
  } catch (e: any) {
    if (get(e, 'stack')) {
      console.error(e)
    }
    if (options.errHandle) {
      options.errHandle(e)
    } else {
      dialog.alert(Err.instance(e).getMsg())
    }
  }
}

