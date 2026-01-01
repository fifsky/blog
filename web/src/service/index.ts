import { createApi } from '@/utils/request'
import type {
  LoginRequest,
  LoginResponse,
  ArticleListRequest,
  ArticleListResponse,
  ArticleItem,
  PrevNextResponse,
  PageRequest,
  MoodListResponse,
  Options,
  CateMenuResponse,
  ArchiveResponse,
  LinkMenuResponse,
  CateListResponse,
  IDRequest,
  IDResponse,
  MoodCreateRequest,
  MoodUpdateRequest,
  CateCreateRequest,
  CateUpdateRequest,
  LinkCreateRequest,
  LinkUpdateRequest,
  User,
  UserListResponse,
  UserCreateRequest,
  UserUpdateRequest,
  GetUserRequest,
  ArticleCreateRequest,
  ArticleUpdateRequest
} from '@/types/openapi'

export const loginApi = (data: LoginRequest) => createApi<LoginResponse>('/api/login', data)
export const articleListApi = (data: ArticleListRequest) => createApi<ArticleListResponse>('/api/article/list', data)
export const articleDetailApi = (data: { id?: number; url?: string }) => createApi<ArticleItem>('/api/article/detail', data)
export const prevnextArticleApi = (data: { id: number }) => createApi<PrevNextResponse>('/api/article/prevnext', data)
export const moodListApi = (data: PageRequest) => createApi<MoodListResponse>('/api/mood/list', data)
export const commentListApi = (data: any) => createApi('/api/comment/list', data)
export const commentPostApi = (data: any) => createApi('/api/comment/post', data)
export const settingApi = () => createApi<Options>('/api/setting')
export const cateAllApi = () => createApi<CateMenuResponse>('/api/cate/all')
export const archiveApi = () => createApi<ArchiveResponse>('/api/article/archive')
export const newCommentApi = () => createApi<any>('/api/comment/new')
export const linkAllApi = () => createApi<LinkMenuResponse>('/api/link/all')

export const loginUserApi = () => createApi<User>('/api/admin/loginuser')
export const settingUpdateApi = (data: Options) => createApi<Options>('/api/admin/setting/update', data)
export const articleCreateApi = (data: ArticleCreateRequest) => createApi<IDResponse>('/api/admin/article/create', data)
export const articleUpdateApi = (data: ArticleUpdateRequest) => createApi<IDResponse>('/api/admin/article/update', data)
export const articleDeleteApi = (data: IDRequest) => createApi('/api/admin/article/delete', data)
export const uploadApi = (data: any) => createApi('/api/admin/upload', data)
export const commentAdminListApi = (data: PageRequest) => createApi('/api/admin/comment/list', data)
export const commentDeleteApi = (data: IDRequest) => createApi('/api/admin/comment/delete', data)
export const moodCreateApi = (data: MoodCreateRequest) => createApi<IDResponse>('/api/admin/mood/create', data)
export const moodUpdateApi = (data: MoodUpdateRequest) => createApi<IDResponse>('/api/admin/mood/update', data)
export const moodDeleteApi = (data: IDRequest) => createApi('/api/admin/mood/delete', data)
export const cateListApi = (data?: any) => createApi<CateListResponse>('/api/admin/cate/list', data)
export const cateCreateApi = (data: CateCreateRequest) => createApi<IDResponse>('/api/admin/cate/create', data)
export const cateUpdateApi = (data: CateUpdateRequest) => createApi<IDResponse>('/api/admin/cate/update', data)
export const cateDeleteApi = (data: IDRequest) => createApi('/api/admin/cate/delete', data)
export const linkListApi = (data?: any) => createApi('/api/admin/link/list', data)
export const linkCreateApi = (data: LinkCreateRequest) => createApi<IDResponse>('/api/admin/link/create', data)
export const linkUpdateApi = (data: LinkUpdateRequest) => createApi<IDResponse>('/api/admin/link/update', data)
export const linkDeleteApi = (data: IDRequest) => createApi('/api/admin/link/delete', data)
export const remindListApi = (data: PageRequest) => createApi('/api/admin/remind/list', data)
export const remindCreateApi = (data: any) => createApi<IDResponse>('/api/admin/remind/create', data)
export const remindUpdateApi = (data: any) => createApi<IDResponse>('/api/admin/remind/update', data)
export const remindDeleteApi = (data: IDRequest) => createApi('/api/admin/remind/delete', data)
export const userListApi = (data: PageRequest) => createApi<UserListResponse>('/api/admin/user/list', data)
export const userCreateApi = (data: UserCreateRequest) => createApi<IDResponse>('/api/admin/user/create', data)
export const userUpdateApi = (data: UserUpdateRequest) => createApi<IDResponse>('/api/admin/user/update', data)
export const userGetApi = (data: GetUserRequest) => createApi<User>('/api/admin/user/get', data)
export const userStatusApi = (data: IDRequest) => createApi('/api/admin/user/status', data)

