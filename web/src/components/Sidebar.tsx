import React, { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useStore } from '@/store/context'
import { SidebarList } from './SidebarList'
import { Calendar } from './Calendar'

export function Sidebar() {
  const { state, dispatch } = useStore()
  const navigate = useNavigate()
  const location = useLocation()
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    navigate({ pathname: '/search', search: `?keyword=${encodeURIComponent(state.keyword)}` })
  }
  const changeKeyword = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch({ type: 'setKeyword', payload: e.target.value })
  }
  useEffect(() => {
    const params = new URLSearchParams(location.search)
    dispatch({ type: 'setKeyword', payload: params.get('keyword') || '' })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.search])
  return (
    <div id="sidebar">
      <div className="sect" id="search">
        <form id="searchpanel" method="get" onSubmit={submit}>
          <p>
            <input className="input_text" type="text" name="keyword" onChange={changeKeyword} value={state.keyword} />&nbsp;
            <input className="formbutton" type="submit" value="搜索" />
          </p>
        </form>
      </div>
      <Calendar />
      <SidebarList title="文章分类" api="cateAllApi" />
      <SidebarList title="历史存档" api="archiveApi" />
      <SidebarList title="我关注的" api="linkAllApi" />
      <div className="rss">
          <svg d="1767272071495" className="icon" viewBox="0 0 1024 1024" version="1.1"
               xmlns="http://www.w3.org/2000/svg" p-id="9471" width="22" height="22">
              <path
                  d="M768 192H256a64 64 0 0 0-64 64v512a64 64 0 0 0 64 64h512a64 64 0 0 0 64-64V256a64 64 0 0 0-64-64z m32 565.344c0 23.552-19.104 42.656-42.656 42.656H266.656A42.656 42.656 0 0 1 224 757.344V266.656C224 243.104 243.104 224 266.656 224h490.656C780.896 224 800 243.104 800 266.656v490.688z"
                  fill="#E47D33" p-id="9472"></path>
              <path
                  d="M800 266.656v490.656c0 23.584-19.104 42.688-42.656 42.688H266.656A42.656 42.656 0 0 1 224 757.344V266.656C224 243.104 243.104 224 266.656 224h490.656C780.896 224 800 243.104 800 266.656zM415.008 667.168c0-16.16-5.664-29.92-16.96-41.216s-25.056-16.96-41.216-16.96-29.888 5.664-41.216 16.96-16.96 25.056-16.96 41.216 5.664 29.888 16.96 41.184 25.056 16.992 41.216 16.992 29.888-5.696 41.216-16.992 16.96-25.024 16.96-41.184z m155.136 37.248a271.2 271.2 0 0 0-24.416-92.672 268.128 268.128 0 0 0-54.976-78.528 268.256 268.256 0 0 0-78.464-54.976 271.936 271.936 0 0 0-92.704-24.416h-1.536a17.6 17.6 0 0 0-13.024 5.152 17.92 17.92 0 0 0-6.368 14.24v40.896c0 5.056 1.664 9.376 4.992 13.024a18.272 18.272 0 0 0 12.576 6.048 188.352 188.352 0 0 1 118.624 55.904 188.64 188.64 0 0 1 55.904 118.656 18.688 18.688 0 0 0 6.048 12.608 18.56 18.56 0 0 0 13.024 4.992h40.896a17.728 17.728 0 0 0 14.208-6.368 18.272 18.272 0 0 0 5.216-14.56z m155.136 0.608a420.288 420.288 0 0 0-36.384-151.968 425.728 425.728 0 0 0-89.056-128.928c-37.6-37.792-80.544-67.488-128.928-89.088s-99.04-33.728-151.968-36.352h-0.896a18.208 18.208 0 0 0-13.344 5.44 17.92 17.92 0 0 0-6.048 13.92v43.328c0 5.056 1.76 9.44 5.312 13.184 3.52 3.744 7.84 5.696 12.864 5.92 43.424 2.624 84.672 12.928 123.776 30.912a348.16 348.16 0 0 1 101.824 70.144c28.768 28.768 52.16 62.72 70.144 101.792a339.2 339.2 0 0 1 30.592 123.808 18.048 18.048 0 0 0 5.92 12.864c3.744 3.52 8.224 5.312 13.472 5.312h43.328a17.952 17.952 0 0 0 13.952-6.048 17.696 17.696 0 0 0 5.44-14.24z m0 0"
                  fill="#FE9832" p-id="9473"></path>
          </svg><a href="https://api.fifsky.com/feed.xml" target="_blank" rel="noreferrer">订阅我的消息</a>
      </div>
    </div>
  )
}

