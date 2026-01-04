import { createApi } from "@/utils/request";
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
  RemindCreateRequest,
  RemindUpdateRequest,
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
  ArticleUpdateRequest,
} from "@/types/openapi";

export const loginApi = (data: LoginRequest, errorHandler?: (e: any) => void) =>
  createApi<LoginResponse>("/api/login", data, errorHandler);
export const articleListApi = (
  data: ArticleListRequest,
  errorHandler?: (e: any) => void
) => createApi<ArticleListResponse>("/api/article/list", data, errorHandler);
export const articleDetailApi = (
  data: { id?: number; url?: string },
  errorHandler?: (e: any) => void
) => createApi<ArticleItem>("/api/article/detail", data, errorHandler);
export const prevnextArticleApi = (
  data: { id: number },
  errorHandler?: (e: any) => void
) => createApi<PrevNextResponse>("/api/article/prevnext", data, errorHandler);
export const moodListApi = (
  data: PageRequest,
  errorHandler?: (e: any) => void
) => createApi<MoodListResponse>("/api/mood/list", data, errorHandler);
export const commentListApi = (data: any, errorHandler?: (e: any) => void) =>
  createApi("/api/comment/list", data, errorHandler);
export const commentPostApi = (data: any, errorHandler?: (e: any) => void) =>
  createApi("/api/comment/post", data, errorHandler);
export const settingApi = (errorHandler?: (e: any) => void) =>
  createApi<Options>("/api/setting", undefined, errorHandler);
export const cateAllApi = (errorHandler?: (e: any) => void) =>
  createApi<CateMenuResponse>("/api/cate/all", undefined, errorHandler);
export const archiveApi = (errorHandler?: (e: any) => void) =>
  createApi<ArchiveResponse>("/api/article/archive", undefined, errorHandler);
export const newCommentApi = (errorHandler?: (e: any) => void) =>
  createApi<any>("/api/comment/new", undefined, errorHandler);
export const linkAllApi = (errorHandler?: (e: any) => void) =>
  createApi<LinkMenuResponse>("/api/link/all", undefined, errorHandler);

export const loginUserApi = (errorHandler?: (e: any) => void) =>
  createApi<User>("/api/admin/loginuser", undefined, errorHandler);
export const settingUpdateApi = (
  data: Options,
  errorHandler?: (e: any) => void
) => createApi<Options>("/api/admin/setting/update", data, errorHandler);
export const articleCreateApi = (
  data: ArticleCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/article/create", data, errorHandler);
export const articleUpdateApi = (
  data: ArticleUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/article/update", data, errorHandler);
export const articleDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/article/delete", data, errorHandler);
export const uploadApi = (data: any, errorHandler?: (e: any) => void) =>
  createApi("/api/admin/upload", data, errorHandler);
export const commentAdminListApi = (
  data: PageRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/comment/list", data, errorHandler);
export const commentDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/comment/delete", data, errorHandler);
export const moodCreateApi = (
  data: MoodCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/mood/create", data, errorHandler);
export const moodUpdateApi = (
  data: MoodUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/mood/update", data, errorHandler);
export const moodDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/mood/delete", data, errorHandler);
export const cateListApi = (data?: any, errorHandler?: (e: any) => void) =>
  createApi<CateListResponse>("/api/admin/cate/list", data, errorHandler);
export const cateCreateApi = (
  data: CateCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/cate/create", data, errorHandler);
export const cateUpdateApi = (
  data: CateUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/cate/update", data, errorHandler);
export const cateDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/cate/delete", data, errorHandler);
export const linkListApi = (data?: any, errorHandler?: (e: any) => void) =>
  createApi("/api/admin/link/list", data, errorHandler);
export const linkCreateApi = (
  data: LinkCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/link/create", data, errorHandler);
export const linkUpdateApi = (
  data: LinkUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/link/update", data, errorHandler);
export const linkDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/link/delete", data, errorHandler);
export const remindListApi = (
  data: PageRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/remind/list", data, errorHandler);
export const remindCreateApi = (
  data: RemindCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/remind/create", data, errorHandler);
export const remindUpdateApi = (
  data: RemindUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/remind/update", data, errorHandler);
export const remindDeleteApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/remind/delete", data, errorHandler);
export const userListApi = (
  data: PageRequest,
  errorHandler?: (e: any) => void
) => createApi<UserListResponse>("/api/admin/user/list", data, errorHandler);
export const userCreateApi = (
  data: UserCreateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/user/create", data, errorHandler);
export const userUpdateApi = (
  data: UserUpdateRequest,
  errorHandler?: (e: any) => void
) => createApi<IDResponse>("/api/admin/user/update", data, errorHandler);
export const userGetApi = (
  data: GetUserRequest,
  errorHandler?: (e: any) => void
) => createApi<User>("/api/admin/user/get", data, errorHandler);
export const userStatusApi = (
  data: IDRequest,
  errorHandler?: (e: any) => void
) => createApi("/api/admin/user/status", data, errorHandler);
