import Message from 'vue-m-message'

export const dialog = {
  alert(msg,fn = null,title="提示"){
    Message.warning(msg);
    fn && fn()
  },
  confirm(msg,fn=null,title="提示"){
    if(window.confirm(msg)){
      fn && fn(true)
    }else{
      fn && fn(false)
    }
  },
}
