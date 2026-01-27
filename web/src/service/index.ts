import { createApi, request } from "@/utils/request";
import { getApiUrl } from "@/utils/common";
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
  MoodItem,
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
  ArticleCalendarRequest,
  ArticleCalendarResponse,
  PhotoListRequest,
  PhotoListResponse,
  PhotoCreateRequest,
  PhotoUpdateRequest,
  PhotoDeleteRequest,
  RegionListRequest,
  RegionListResponse,
  OSSPresignRequest,
  OSSPresignResponse,
  FootprintsResponse,
  CityPhotosRequest,
  CityPhotosResponse,
  GenerateTagsRequest,
  GenerateTagsResponse,
} from "@/types/openapi";

export const loginApi = (data: LoginRequest, errorHandler?: (e: AppError) => void) =>
  createApi<LoginResponse>("/blog/login", data, errorHandler);
export const articleListApi = (data: ArticleListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<ArticleListResponse>("/blog/article/list", data, errorHandler);
export const articleCalendarApi = (
  data: ArticleCalendarRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<ArticleCalendarResponse>("/blog/article/calendar", data, errorHandler);
export const articleDetailApi = (
  data: { id?: number; url?: string },
  errorHandler?: (e: AppError) => void,
) => createApi<ArticleItem>("/blog/article/detail", data, errorHandler);
export const prevnextArticleApi = (data: { id: number }, errorHandler?: (e: AppError) => void) =>
  createApi<PrevNextResponse>("/blog/article/prevnext", data, errorHandler);
export const moodListApi = (data: MoodListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<MoodListResponse>("/blog/mood/list", data, errorHandler);
export const moodRandomApi = (errorHandler?: (e: AppError) => void) =>
  createApi<MoodItem>("/blog/mood/random", undefined, errorHandler);
export const commentListApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/comment/list", data, errorHandler);
export const commentPostApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/comment/post", data, errorHandler);
export const settingApi = (errorHandler?: (e: AppError) => void) =>
  createApi<Options>("/blog/setting", undefined, errorHandler);
export const settingChinaMapApi = (errorHandler?: (e: AppError) => void) =>
  request<any>({ url: getApiUrl("/blog/china_map") }, errorHandler);
export const cateAllApi = (errorHandler?: (e: AppError) => void) =>
  createApi<CateMenuResponse>("/blog/cate/all", undefined, errorHandler);
export const archiveApi = (errorHandler?: (e: AppError) => void) =>
  createApi<ArchiveResponse>("/blog/article/archive", undefined, errorHandler);
export const newCommentApi = (errorHandler?: (e: AppError) => void) =>
  createApi<any>("/blog/comment/new", undefined, errorHandler);
export const linkAllApi = (errorHandler?: (e: AppError) => void) =>
  createApi<LinkMenuResponse>("/blog/link/all", undefined, errorHandler);

export const loginUserApi = (errorHandler?: (e: AppError) => void) =>
  createApi<User>("/blog/admin/loginuser", undefined, errorHandler);
export const aiGenerateTagsApi = (
  data: GenerateTagsRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<GenerateTagsResponse>("/blog/admin/ai/tags", data, errorHandler);
export const settingUpdateApi = (data: Options, errorHandler?: (e: AppError) => void) =>
  createApi<Options>("/blog/admin/setting/update", data, errorHandler);
export const articleCreateApi = (
  data: ArticleCreateRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/blog/admin/article/create", data, errorHandler);
export const articleUpdateApi = (
  data: ArticleUpdateRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/blog/admin/article/update", data, errorHandler);
export const articleDeleteApi = (
  data: ArticleDeleteRequest,
  errorHandler?: (e: AppError) => void,
) => createApi("/blog/admin/article/delete", data, errorHandler);
export const articleRestoreApi = (
  data: ArticleRestoreRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<IDResponse>("/blog/admin/article/restore", data, errorHandler);
export const articleListAdminApi = (
  data: AdminArticleListRequest,
  errorHandler?: (e: AppError) => void,
) => createApi<AdminArticleListResponse>("/blog/admin/article/list", data, errorHandler);
export const adminArticleDetailApi = (data: { id: number }, errorHandler?: (e: AppError) => void) =>
  createApi<ArticleItem>("/blog/admin/article/detail", data, errorHandler);
export const uploadApi = (data: any, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/upload", data, errorHandler);
export const commentAdminListApi = (data: UserListRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/comment/list", data, errorHandler);
export const commentDeleteApi = (data: { id: number }, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/comment/delete", data, errorHandler);
export const moodCreateApi = (data: MoodCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/mood/create", data, errorHandler);
export const moodUpdateApi = (data: MoodUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/mood/update", data, errorHandler);
export const moodDeleteApi = (data: MoodDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/mood/delete", data, errorHandler);
export const cateListApi = (data?: any, errorHandler?: (e: AppError) => void) =>
  createApi<CateListResponse>("/blog/admin/cate/list", data, errorHandler);
export const cateCreateApi = (data: CateCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/cate/create", data, errorHandler);
export const cateUpdateApi = (data: CateUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/cate/update", data, errorHandler);
export const cateDeleteApi = (data: CateDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/cate/delete", data, errorHandler);
export const linkListApi = (data?: any, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/link/list", data, errorHandler);
export const linkCreateApi = (data: LinkCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/link/create", data, errorHandler);
export const linkUpdateApi = (data: LinkUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/link/update", data, errorHandler);
export const linkDeleteApi = (data: LinkDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/link/delete", data, errorHandler);
export const remindListApi = (data: RemindListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<RemindListResponse>("/blog/admin/remind/list", data, errorHandler);
export const remindCreateApi = (data: RemindCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/remind/create", data, errorHandler);
export const remindUpdateApi = (data: RemindUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/remind/update", data, errorHandler);
export const remindDeleteApi = (data: RemindDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/remind/delete", data, errorHandler);
export const userListApi = (data: UserListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<UserListResponse>("/blog/admin/user/list", data, errorHandler);
export const userCreateApi = (data: UserCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/user/create", data, errorHandler);
export const userUpdateApi = (data: UserUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/user/update", data, errorHandler);
export const userGetApi = (data: GetUserRequest, errorHandler?: (e: AppError) => void) =>
  createApi<User>("/blog/admin/user/get", data, errorHandler);
export const userStatusApi = (data: UserStatusRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/user/status", data, errorHandler);

// Photo APIs
export const photoListApi = (data: PhotoListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<PhotoListResponse>("/blog/admin/photo/list", data, errorHandler);
export const photoCreateApi = (data: PhotoCreateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/photo/create", data, errorHandler);
export const photoUpdateApi = (data: PhotoUpdateRequest, errorHandler?: (e: AppError) => void) =>
  createApi<IDResponse>("/blog/admin/photo/update", data, errorHandler);
export const photoDeleteApi = (data: PhotoDeleteRequest, errorHandler?: (e: AppError) => void) =>
  createApi("/blog/admin/photo/delete", data, errorHandler);

// Region APIs
export const regionListApi = (data: RegionListRequest, errorHandler?: (e: AppError) => void) =>
  createApi<RegionListResponse>("/blog/admin/region/list", data, errorHandler);

// OSS APIs
export const ossPresignApi = (data: OSSPresignRequest, errorHandler?: (e: AppError) => void) =>
  createApi<OSSPresignResponse>("/blog/admin/oss/presign", data, errorHandler);

// Travel APIs (public)
export const footprintsApi = (errorHandler?: (e: AppError) => void) =>
  createApi<FootprintsResponse>("/blog/travel/footprints", undefined, errorHandler);
export const cityPhotosApi = (data: CityPhotosRequest, errorHandler?: (e: AppError) => void) =>
  createApi<CityPhotosResponse>("/blog/travel/photos", data, errorHandler);
