// Article 批量操作
export type ArticleCalendarRequest = {
  year: number;
  month: number;
};

export type ArticleCalendarResponse = {
  days: number[];
};

export type ArticleDeleteRequest = { ids: number[] };
export type ArticleRestoreRequest = { ids: number[] };

// 通用响应
export type IDResponse = { id: number };

export type LoginRequest = { user_name: string; password: string };
export type LoginResponse = { access_token: string; user: UserItem };

export type UserItem = {
  id: number;
  name: string;
  nick_name: string;
  email: string;
  status: number;
  type: number;
  created_at: string;
  updated_at: string;
};

export type User = {
  id: number;
  name: string;
  nick_name: string;
  email: string;
  status: number;
  type: number;
  created_at: string;
  updated_at: string;
};

// 各模块独立的请求类型
export type UserListRequest = { page?: number };
export type UserStatusRequest = { id: number };
export type RemindListRequest = { page?: number };
export type RemindDeleteRequest = { id: number };
export type MoodListRequest = { page?: number };
export type MoodDeleteRequest = { id: number };
export type CateDeleteRequest = { id: number };
export type LinkDeleteRequest = { id: number };

export type UserUpdateRequest = {
  id: number;
  name?: string;
  password?: string;
  nick_name?: string;
  email?: string;
  type?: number;
};
export type GetUserRequest = { id: number };
// 站点设置
// site_name: 站点名称
// site_desc: 站点描述
// site_keyword: 站点关键字
// post_num: 每页显示文章数
export type Options = { kv: Record<string, string> };

export type CateMenuItem = { url: string; content: string };
export type CateMenuResponse = { list: CateMenuItem[] };
export type LinkMenuItem = { url: string; content: string };
export type LinkMenuResponse = { list: LinkMenuItem[] };
export type DateArchiveItem = { url: string; content: string };
export type ArchiveResponse = { list: DateArchiveItem[] };

export type CateListItem = {
  id: number;
  name: string;
  desc: string;
  domain: string;
  created_at: string;
  updated_at: string;
  num: number;
};
export type CateListResponse = { list: CateListItem[]; total: number };
export type CateCreateRequest = { name: string; domain: string; desc?: string };
export type CateUpdateRequest = {
  id: number;
  name?: string;
  desc?: string;
  domain?: string;
};

export type UserSummary = { id: number; name: string; nick_name: string };
export type CateSummary = { id: number; name: string; domain: string };

export type ArticleItem = {
  id: number;
  cate_id: number;
  type: number;
  user_id: number;
  title: string;
  url?: string;
  content: string;
  tags?: string[];
  status: number;
  view_num: number;
  created_at: string;
  updated_at: string;
  user: UserSummary;
  cate: CateSummary;
};
export type ArticleListRequest = {
  year?: string;
  month?: string;
  domain?: string;
  keyword?: string;
  page?: number;
  type?: number;
  day?: string;
  page_size?: number;
  tag?: string;
};
export type ArticleListResponse = { list: ArticleItem[]; total: number };
export type ArticleCreateRequest = {
  cate_id: number;
  type: number;
  title: string;
  url?: string;
  content: string;
  status?: number;
  tags?: string[];
};
export type ArticleUpdateRequest = {
  id: number;
  cate_id?: number;
  type?: number;
  title?: string;
  url?: string;
  content?: string;
  status?: number;
  tags?: string[];
};

export type AdminArticleListRequest = {
  page?: number;
  type?: number;
  status?: number;
};

export type AdminArticleListResponse = {
  list: ArticleItem[];
  total: number;
};

export type PrevNextItem = { id: number; title: string };
export type PrevNextResponse = { prev?: PrevNextItem; next?: PrevNextItem };

export type MoodItem = {
  id: number;
  content: string;
  user: UserSummary;
  created_at: string;
};
export type MoodListResponse = { list: MoodItem[]; total: number };
export type MoodCreateRequest = { content: string };
export type MoodUpdateRequest = { id: number; content?: string };

export type CommentItem = {
  id: number;
  article_title: string;
  name: string;
  content: string;
  ip: string;
  created_at: string;
  type: number;
  url?: string;
};
export type CommentListResponse = { list: CommentItem[]; total: number };

// 缺失的类型定义
export type ArticleDetailRequest = { id?: number; url?: string };
export type GoogleProtobufAny = { "@type"?: string } & Record<string, any>;
export type LinkItem = {
  id: number;
  name: string;
  url: string;
  desc?: string;
  created_at: string;
};
export type LinkListResponse = { list: LinkItem[]; total: number };
export type PrevNextRequest = { id: number };
export type RemindItem = {
  id: number;
  type: number;
  content: string;
  month?: number;
  week?: number;
  day?: number;
  hour?: number;
  minute?: number;
  status: number;
  next_time: string;
  created_at: string;
};
export type RemindListResponse = { list: RemindItem[]; total: number };
export type Status = {
  code: number;
  message: string;
  details?: GoogleProtobufAny[];
};
export type TextResponse = { text: string };
export type UserCreateRequest = {
  name: string;
  password: string;
  nick_name: string;
  email?: string;
  type: number;
};
export type UserListResponse = { list: UserItem[]; total: number };

export type RemindChangeRequest = { token: string };
export type RemindDelayRequest = { token: string };
export type LinkCreateRequest = {
  name: string;
  url: string;
  desc?: string;
};

export type LinkUpdateRequest = {
  id: number;
  name?: string;
  url?: string;
  desc?: string;
};

export type RemindCreateRequest = {
  type: number;
  month?: number;
  week?: number;
  day?: number;
  hour?: number;
  minute?: number;
  content: string;
};

export type RemindUpdateRequest = {
  id: number;
  type?: number;
  month?: number;
  week?: number;
  day?: number;
  hour?: number;
  minute?: number;
  content?: string;
};

// Photo types
export type PhotoItem = {
  id: number;
  title: string;
  description: string;
  src: string;
  thumbnail: string;
  province: number;
  province_name: string;
  city: number;
  city_name: string;
  created_at: string;
};
export type PhotoListRequest = { page: number };
export type PhotoListResponse = { list: PhotoItem[]; total: number };
export type PhotoCreateRequest = {
  title: string;
  description?: string;
  srcs: string[]; // 支持多个图片地址
  province: number;
  city: number;
};
export type PhotoUpdateRequest = {
  id: number;
  title?: string;
  description?: string;
  province?: number;
  city?: number;
};
export type PhotoDeleteRequest = { id: number };

// Region types
export type RegionItem = {
  region_id: number;
  parent_id: number;
  level: number;
  region_name: string;
  longitude: string;
  latitude: string;
  pinyin: string;
  az_no: string;
};
export type RegionListRequest = { parent_id: number };
export type RegionListResponse = { list: RegionItem[] };

// OSS types
export type OSSPresignRequest = { filename: string };
export type OSSPresignResponse = {
  url: string; // 预签名上传URL
  cdn_url: string; // CDN访问地址
};

// Travel types
export type FootprintRegion = {
  region_id: number;
  name: string;
  longitude: string;
  latitude: string;
};
export type FootprintsResponse = {
  provinces: FootprintRegion[];
  cities: FootprintRegion[];
};
export type TravelPhoto = {
  title: string;
  description: string;
  src: string;
  thumbnail: string;
};
export type CityPhotosRequest = { region_id: number };
export type CityPhotosResponse = { photos: TravelPhoto[] };
