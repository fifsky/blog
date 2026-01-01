import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { archiveApi, cateAllApi, linkAllApi, newCommentApi } from '@/service'

const apiMap = { cateAllApi, newCommentApi, archiveApi, linkAllApi }

export function SidebarList({ title, api }: { title: string; api: keyof typeof apiMap }) {
  const [items, setItems] = useState<{ url: string; content: string }[]>([])
  useEffect(() => {
    ;(async () => {
      const data = await apiMap[api]()
      setItems(data.list || [])
    })()
  }, [api])
  return (
    <div className="sect">
      <h2>{title}</h2>
      <ul className="tlist">
        {items.map((v, k) => (
          <li key={k}>
            {v.url && v.url.startsWith('http') ? (
              <a target="_blank" href={v.url} rel="noreferrer">{v.content}</a>
            ) : (
              <Link to={v.url}>{v.content}</Link>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}
