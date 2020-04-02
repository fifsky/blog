import {createApi} from "../utils";

export const loginApi = data => createApi("/api/login", data);
export const articleListApi = data => createApi("/api/article/list", data);
export const articleDetailApi = data => createApi("/api/article/detail", data);
export const prevnextArticleApi = data => createApi("/api/article/prevnext", data);
export const moodListApi = data => createApi("/api/mood/list",data);
export const commentListApi = data => createApi("/api/comment/list", data);
export const commentPostApi = data => createApi("/api/comment/post", data);
export const settingApi = () => createApi("/api/setting");
export const cateAllApi = () => createApi("/api/cate/all");
export const archiveApi = () => createApi("/api/article/archive");
export const newCommentApi = () => createApi("/api/comment/new");
export const linkAllApi = () => createApi("/api/link/all");


//需要登录
export const loginUserApi = () => createApi("/api/admin/loginuser");
export const settingPostApi = data => createApi("/api/admin/setting/post", data);
export const articlePostApi = data => createApi("/api/admin/article/post", data);
export const articleDeleteApi = data => createApi("/api/admin/article/delete", data);
export const uploadApi = data => createApi("/api/admin/upload", data);
export const commentAdminListApi = data => createApi("/api/admin/comment/list", data);
export const commentDeleteApi = data => createApi("/api/admin/comment/delete", data);
export const moodPostApi = data => createApi("/api/admin/mood/post", data);
export const moodDeleteApi = data => createApi("/api/admin/mood/delete", data);
export const cateListApi = data => createApi("/api/admin/cate/list", data);
export const catePostApi = data => createApi("/api/admin/cate/post", data);
export const cateDeleteApi = data => createApi("/api/admin/cate/delete", data);
export const linkListApi = data => createApi("/api/admin/link/list", data);
export const linkPostApi = data => createApi("/api/admin/link/post", data);
export const linkDeleteApi = data => createApi("/api/admin/link/delete", data);
export const remindListApi = data => createApi("/api/admin/remind/list", data);
export const remindPostApi = data => createApi("/api/admin/remind/post", data);
export const remindDeleteApi = data => createApi("/api/admin/remind/delete", data);
export const userListApi = data => createApi("/api/admin/user/list", data);
export const userPostApi = data => createApi("/api/admin/user/post", data);
export const userGetApi = data => createApi("/api/admin/user/get", data);
export const userStatusApi = data => createApi("/api/admin/user/status", data);