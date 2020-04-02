import { Err } from './error'
import { get } from './mapping'
import { dialog } from './dialog'

export const sync = (fn,options = {errHandle:null}) => {
  return fn().catch((e)=>{
    if(get(e,"stack")){
      console.error(e)
    }

    if (options.errHandle) {
        options.errHandle(e)
    }else{
      dialog.alert(Err.instance(e).getMsg())
    }
  })
}
