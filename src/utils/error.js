import { get } from './mapping'

export class Err {
  data

  static instance(e) {
    if (e instanceof Err) {
      return e
    }
    return new Err(e)
  }

  constructor(data) {
    this.data = data
  }

  getMsg(){
    if(get(this.data,"stack")){
      return this.data.message
    }

    return get(this.data,"message","未知错误"+JSON.stringify(this.data))
  }
}

export const errors = function (e) {
  return Err.instance(e)
}
