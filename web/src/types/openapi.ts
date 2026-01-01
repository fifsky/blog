export type IDRequest = { id: number }
export type IDResponse = { id: number }

export type LoginRequest = { user_name: string; password: string }
export type LoginResponse = { access_token: string; user: UserItem }

export type UserItem = {
  id: number
  name: string
  nick_name: string
  email: string
  status: number
  type: number
  created_at: string
  updated_at: string
}

export type User = {
  id: number
  name: string
  nick_name: string
  email: string
  status: number
  type: number
  created_at: string
  updated_at: string
}

export type PageRequest = { page?: number }

export type UserCreateRequest = { name: string; password: string; nick_name: string; email?: string; type: number }
export type UserUpdateRequest = { id: number; name?: string; password?: string; nick_name?: string; email?: string; type?: number }
export type UserListResponse = { list: UserItem[]; page_total: number }
export type GetUserRequest = { id: number }

export type Options = { kv: Record<string, string> }

export type CateMenuItem = { url: string; content: string }
export type CateMenuResponse = { list: CateMenuItem[] }
export type LinkMenuItem = { url: string; content: string }
export type LinkMenuResponse = { list: LinkMenuItem[] }
export type DateArchiveItem = { url: string; content: string }
export type ArchiveResponse = { list: DateArchiveItem[] }

export type CateListItem = {
  id: number
  name: string
  desc: string
  domain: string
  created_at: string
  updated_at: string
  num: number
}
export type CateListResponse = { list: CateListItem[]; page_total: number }
export type CateCreateRequest = { name: string; domain: string; desc?: string }
export type CateUpdateRequest = { id: number; name?: string; desc?: string; domain?: string }

export type UserSummary = { id: number; name: string; nick_name: string }
export type CateSummary = { id: number; name: string; domain: string }

export type ArticleItem = {
  id: number
  cate_id: number
  type: number
  user_id: number
  title: string
  url?: string
  content: string
  status: number
  created_at: string
  updated_at: string
  user: UserSummary
  cate: CateSummary
}
export type ArticleListRequest = {
  year?: string
  month?: string
  domain?: string
  keyword?: string
  page?: number
  type?: number
}
export type ArticleListResponse = { list: ArticleItem[]; page_total: number }
export type ArticleCreateRequest = { cate_id: number; type: number; title: string; url?: string; content: string }
export type ArticleUpdateRequest = { id: number; cate_id?: number; type?: number; title?: string; url?: string; content?: string; status?: number }

export type PrevNextItem = { id: number; title: string }
export type PrevNextResponse = { prev?: PrevNextItem; next?: PrevNextItem }

export type MoodItem = { id: number; content: string; user: UserSummary; created_at: string }
export type MoodListResponse = { list: MoodItem[]; page_total: number }
export type MoodCreateRequest = { content: string }
export type MoodUpdateRequest = { id: number; content?: string }

export type LinkCreateRequest = {
    name: string;
    url: string;
    desc?: string;
}

export type LinkUpdateRequest = {
    id: number;
    name?: string;
    url?: string;
    desc?: string;
}
