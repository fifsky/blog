import { Component, PropsWithChildren } from 'react'
import Taro from '@tarojs/taro'

import './app.scss'


  class App extends Component<PropsWithChildren> {

  componentDidMount () {}

  componentDidShow () {
    const token = Taro.getStorageSync('access_token')
    const pages = Taro.getCurrentPages()
    const current = pages[pages.length - 1]
    const currentPath = current?.route ? `/${current.route}` : ''

    if (!token && currentPath !== '/pages/login/index') {
      Taro.reLaunch({ url: '/pages/login/index' })
    }
  }

  componentDidHide () {}

  // this.props.children 是将要会渲染的页面
  render () {
    return this.props.children
  }
}


export default App
