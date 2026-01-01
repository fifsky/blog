export const get = (obj: any, key: string, defaultValue?: any) => {
  try {
    const parts = key.replace(/\[(\w+)\]/g, '.$1').split('.')
    let curr = obj
    for (const p of parts) {
      if (curr == null) return defaultValue
      curr = curr[p]
    }
    return curr ?? defaultValue
  } catch {
    return defaultValue
  }
}

