import { createApi } from "@/utils/request";
import type { AppError } from "@/utils/error";
import type {
  LoginRequest,
  LoginResponse,
  ArticleListRequest,
  ArticleListResponse,
  ArticleItem,
  PrevNextResponse,
  MoodListRequest,
  MoodListResponse,
  Options,
  CateMenuResponse,
  ArchiveResponse,
  LinkMenuResponse,
  CateListResponse,
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
  RemindListResponse,
  AdminArticleListRequest,
  AdminArticleListResponse,
  ArticleDeleteRequest,
  ArticleRestoreRequest,
  UserListRequest,
  UserStatusRequest,
  RemindListRequest,
  RemindDeleteRequest,
  MoodDeleteRequest,
  CateDeleteRequest,
  LinkDeleteRequest,
} from "@/types/openapi";

export const loginApi = (data: LoginRequest, errorHandler?: (e: AppError) => void) =>
  createApi<LoginResponse>("/api/login", data, errorHandler);
export const articleListApi = (data: ArticleListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<ArticleListResponse>("/api/article/list", data, errorHandler);
export const articleDetailApi = (
  data: { id?: number; url?: string },
  errorHandler?: (e: AppError) => void,
) => createApi<ArticleItem>("/api/article/detail", data, errorHandler);
export const prevnextArticleApi = (data: { id: number }, errorHandler?: (e: AppError) => void) =>
  createApi<PrevNextResponse>("/api/article/prevnext", data, errorHandler);
export const moodListApi = (data: MoodListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<MoodListResponse>("/api/mood/list", data, errorHandler);
export const commentListApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/api/comment/list", data, errorHandler);
export const commentPostApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/api/comment/post", data, errorHandler);
export const settingApi = (errorHandler?: (e: AppError) => void) =>
  createApi<Options>("/api/setting", undefined, errorHandler);
export const cateAllApi = (errorHandler?: (e: AppError) => void) =>
  createApi<CateMenuResponse>("/api/cate/all", undefined, errorHandler);
export const archiveApi = (errorHandler?: (e: AppError) => void) =>
  createApi<ArchiveResponse>("/api/article/archive", undefined, errorHandler);
export const newCommentApi = (errorHandler?: (e: AppError) => void) =>
  createApi<any>("/api/comment/new", undefined, errorHandler);
export const linkAllApi = (errorHandler?: (e: AppError) => void) =>
  createApi<LinkMenuResponse>("/api/link/all", undefined, errorHandler);

export const loginUserApi = (errorHandler?: (e: AppError) => void) =>
  createApi<User>("/api/admin/loginuser", undefined, errorHandler);
export const settingUpdateApi = (data: Options, errorHandler?: (e: AppError) => void) =>
  createApi<Options>("/api/admin/setting/update", data, errorHandler);
export const articleCreateApi = (
  data: ArticleCreateRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/api/admin/article/create", data, errorHandler);
export const articleUpdateApi = (
  data: ArticleUpdateRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/api/admin/article/update", data, errorHandler);
export const articleDeleteApi = (
  data: ArticleDeleteRequest,
  errorHandler?: (e: AppError) => void,
) => createApi("/api/admin/article/delete", data, errorHandler);
export const articleRestoreApi = (
  data: ArticleRestoreRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/api/admin/article/restore", data, errorHandler);
export const articleListAdminApi = (
  data: AdminArticleListRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<AdminArticleListResponse>("/api/admin/article/list", data, errorHandler);
export const uploadApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/upload", data, errorHandler);
export const commentAdminListApi = (data: UserListRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/comment/list", data, errorHandler);
export const commentDeleteApi = (data: { id: number }, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/comment/delete", data, errorHandler);
export const moodCreateApi = (data: MoodCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/mood/create", data, errorHandler);
export const moodUpdateApi = (data: MoodUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/mood/update", data, errorHandler);
export const moodDeleteApi = (data: MoodDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/mood/delete", data, errorHandler);
export const cateListApi = (data?: any, errorHandler?: (e: AppError) => void) =>
  createApi<CateListResponse>("/api/admin/cate/list", data, errorHandler);
export const cateCreateApi = (data: CateCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/cate/create", data, errorHandler);
export const cateUpdateApi = (data: CateUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/cate/update", data, errorHandler);
export const cateDeleteApi = (data: CateDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/cate/delete", data, errorHandler);
export const linkListApi = (data?: any, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/link/list", data, errorHandler);
export const linkCreateApi = (data: LinkCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/link/create", data, errorHandler);
export const linkUpdateApi = (data: LinkUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/link/update", data, errorHandler);
export const linkDeleteApi = (data: LinkDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/link/delete", data, errorHandler);
export const remindListApi = (data: RemindListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<RemindListResponse>("/api/admin/remind/list", data, errorHandler);
export const remindCreateApi = (data: RemindCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/remind/create", data, errorHandler);
export const remindUpdateApi = (data: RemindUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/remind/update", data, errorHandler);
export const remindDeleteApi = (data: RemindDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/remind/delete", data, errorHandler);
export const userListApi = (data: UserListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<UserListResponse>("/api/admin/user/list", data, errorHandler);
export const userCreateApi = (data: UserCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/user/create", data, errorHandler);
export const userUpdateApi = (data: UserUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/api/admin/user/update", data, errorHandler);
export const userGetApi = (data: GetUserRequest, errorHandler?: (e: AppError) => void) =>
  createApi<User>("/api/admin/user/get", data, errorHandler);
export const userStatusApi = (data: UserStatusRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/api/admin/user/status", data, errorHandler);
