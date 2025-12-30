import { createApi } from "../utils";

export const loginApi = (data) => createApi("/api/login", data);
export const articleListApi = (data) => createApi("/api/article/list", data);
export const articleDetailApi = (data) =>
  createApi("/api/article/detail", data);
export const prevnextArticleApi = (data) =>
  createApi("/api/article/prevnext", data);
export const moodListApi = (data) => createApi("/api/mood/list", data);
export const commentListApi = (data) => createApi("/api/comment/list", data);
export const commentPostApi = (data) => createApi("/api/comment/post", data);
export const settingApi = () => createApi("/api/setting");
export const cateAllApi = () => createApi("/api/cate/all");
export const archiveApi = () => createApi("/api/article/archive");
export const newCommentApi = () => createApi("/api/comment/new");
export const linkAllApi = () => createApi("/api/link/all");

//需要登录
export const loginUserApi = () => createApi("/api/admin/loginuser");
export const settingUpdateApi = (data) =>
  createApi("/api/admin/setting/update", data);
export const articleCreateApi = (data) =>
  createApi("/api/admin/article/create", data);
export const articleUpdateApi = (data) =>
  createApi("/api/admin/article/update", data);
export const articleDeleteApi = (data) =>
  createApi("/api/admin/article/delete", data);
export const uploadApi = (data) => createApi("/api/admin/upload", data);
export const commentAdminListApi = (data) =>
  createApi("/api/admin/comment/list", data);
export const commentDeleteApi = (data) =>
  createApi("/api/admin/comment/delete", data);
export const moodCreateApi = (data) =>
  createApi("/api/admin/mood/create", data);
export const moodUpdateApi = (data) =>
  createApi("/api/admin/mood/update", data);
export const moodDeleteApi = (data) =>
  createApi("/api/admin/mood/delete", data);
export const cateListApi = (data) => createApi("/api/admin/cate/list", data);
export const cateCreateApi = (data) =>
  createApi("/api/admin/cate/create", data);
export const cateUpdateApi = (data) =>
  createApi("/api/admin/cate/update", data);
export const cateDeleteApi = (data) =>
  createApi("/api/admin/cate/delete", data);
export const linkListApi = (data) => createApi("/api/admin/link/list", data);
export const linkCreateApi = (data) =>
  createApi("/api/admin/link/create", data);
export const linkUpdateApi = (data) =>
  createApi("/api/admin/link/update", data);
export const linkDeleteApi = (data) =>
  createApi("/api/admin/link/delete", data);
export const remindListApi = (data) =>
  createApi("/api/admin/remind/list", data);
export const remindCreateApi = (data) =>
  createApi("/api/admin/remind/create", data);
export const remindUpdateApi = (data) =>
  createApi("/api/admin/remind/update", data);
export const remindDeleteApi = (data) =>
  createApi("/api/admin/remind/delete", data);
export const userListApi = (data) => createApi("/api/admin/user/list", data);
export const userCreateApi = (data) =>
  createApi("/api/admin/user/create", data);
export const userUpdateApi = (data) =>
  createApi("/api/admin/user/update", data);
export const userGetApi = (data) => createApi("/api/admin/user/get", data);
export const userStatusApi = (data) =>
  createApi("/api/admin/user/status", data);
