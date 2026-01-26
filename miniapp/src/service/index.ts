import type {
  MiniAppLoginRequest,
  MiniAppLoginResponse,
  MoodCreateRequest,
  MoodListRequest,
  MoodListResponse,
  OSSPresignRequest,
  OSSPresignResponse,
  PhotoCreateRequest,
  PhotoListRequest,
  PhotoListResponse,
  RegionListRequest,
  RegionListResponse,
  RemindCreateRequest,
  RemindDeleteRequest,
  RemindListRequest,
  RemindListResponse,
  IDResponse,
} from "../types/openapi";
import { getApiBaseUrl } from "../config/env";
import { createApi, request } from "./request";

export const miniappLoginApi = (data: MiniAppLoginRequest) =>
  createApi<MiniAppLoginResponse>("/blog/miniapp/login", data);

export const moodListApi = (data: MoodListRequest) => createApi<MoodListResponse>("/blog/mood/list", data);
export const moodCreateApi = (data: MoodCreateRequest) =>
  createApi<IDResponse>("/blog/admin/mood/create", data);

export const regionListApi = (data: RegionListRequest) =>
  createApi<RegionListResponse>("/blog/admin/region/list", data);

export const nearestRegionApi = (latitude: number, longitude: number) => {
  const base = getApiBaseUrl();
  return request<{ province_id: number; province_name: string; city_id: number; city_name: string }>({
    url: `${base}/blog/region/nearest?latitude=${latitude}&longitude=${longitude}`,
    method: "GET",
  });
};

export const ossPresignApi = (data: OSSPresignRequest) =>
  createApi<OSSPresignResponse>("/blog/admin/oss/presign", data);

export const photoListApi = (data: PhotoListRequest) =>
  createApi<PhotoListResponse>("/blog/admin/photo/list", data);
export const photoCreateApi = (data: PhotoCreateRequest) =>
  createApi<IDResponse>("/blog/admin/photo/create", data);

export const remindListApi = (data: RemindListRequest) =>
  createApi<RemindListResponse>("/blog/admin/remind/list", data);
export const remindCreateApi = (data: RemindCreateRequest) =>
  createApi<IDResponse>("/blog/admin/remind/create", data);
export const remindDeleteApi = (data: RemindDeleteRequest) =>
  createApi("/blog/admin/remind/delete", data);
