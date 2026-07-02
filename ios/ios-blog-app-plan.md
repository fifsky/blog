# iOS 博客 App 开发计划

## 概述

基于现有 Go 后端博客系统，开发 iOS 原生 App（Swift + SwiftUI），支持博文、心情、提醒、足迹四大模块。仅供个人使用，通过 Xcode 免费侧载到个人 iPhone。

## 现状分析

### 后端 API 现有资源

现有后端提供完整的 Admin API（JWT 认证），可直接复用：

| 模块 | 已有接口 | 说明 |
|------|---------|------|
| 认证 | `POST /blog/login` | 返回 JWT token，有效期 24h |
| 文章 | `/blog/admin/article/create\|update\|delete\|list\|detail` | 完整 CRUD |
| 评论 | `/blog/admin/comment/list\|delete` | 列表 + 删除 |
| 心情 | `/blog/admin/mood/create\|update\|delete` | 完整 CRUD |
| 提醒 | `/blog/admin/remind/list\|create\|update\|delete` | 完整 CRUD |
| 足迹 | `/blog/admin/footprint/list\|create\|update\|delete` | 完整 CRUD |
| 图片上传 | `/blog/admin/upload` | multipart/form-data |
| OSS 预签名 | `/blog/admin/oss/presign` | 返回 PUT URL + CDN URL |
| 分类 | `POST /blog/cate/all` | 公开 API，无需认证 |

### 免费账号可行性

**完全可行。** Apple 官方支持免费 Apple ID 通过 Xcode 侧载。

可用功能：相机、相册、定位、MapKit 地图、蓝牙、Core Data
不可用功能：远程推送（APNs）、iCloud、内购、Sign in with Apple

**主要限制：** 每 7 天需连接 Xcode 重新编译安装（可通过 AltStore 自动化续签）。

**提醒说明：** 提醒模块仅为管理 CRUD 操作，实际提醒通知由博客后端内置定时任务通过飞书发送，App 端不需要实现本地通知。

### 需要后端修改的地方

1. **RemindUpdateRequest 增加 status 字段** — 当前 proto 定义只有 `id`、`cron`、`content`，App 需要修改提醒状态（标记完成/重新激活）
2. **OSS 预签名接口 content-type 问题** — 当前后端 presign 写死了 `Content-Type: text/plain;charset=utf8`，App 端上传图片需要 `image/jpeg`，应改为由客户端指定或根据扩展名自动判断

---

## 技术方案

### 架构：MVVM + Service Layer

```
ios/
├── BlogApp.swift
├── ContentView.swift                    // TabView 主界面
│
├── App/                                 // 应用配置
│   └── Config.swift                    // API Base URL
│
├── Core/                                // 核心基础设施
│   ├── Network/
│   │   ├── APIClient.swift              // URLSession 封装（POST + JSON + Bearer）
│   │   ├── APIError.swift               // 统一错误类型
│   │   └── MultipartFormData.swift      // 文件上传
│   ├── Auth/
│   │   ├── AuthManager.swift            // JWT token 管理（@Observable）
│   │   └── KeychainService.swift        // Keychain 封装
│   └── Extensions/
│       ├── String+JSON.swift
│       └── Date+Format.swift
│
├── Services/                             // API 服务层
│   ├── AuthService.swift
│   ├── ArticleService.swift
│   ├── MoodService.swift
│   ├── RemindService.swift
│   ├── FootprintService.swift
│   └── UploadService.swift              // 直传 + OSS 预签名
│
├── Models/                              // Codable 数据模型
│   ├── User.swift, Article.swift
│   ├── Comment.swift, Mood.swift
│   ├── Remind.swift, Footprint.swift
│   └── APIResponse.swift
│
├── Features/                            // 功能模块
│   ├── Login/
│   ├── Article/                         // 列表 + 详情(Markdown渲染) + 编辑器 + 评论
│   ├── Mood/                            // 列表 + 编辑器 + EmojiPicker
│   ├── Remind/                          // 列表 + 编辑器（纯 CRUD，无通知）
│   └── Footprint/                       // 地图(MapKit) + 列表 + 编辑器(选点+上传)
│
└── Resources/
    └── Assets.xcassets
```

### 网络层设计

后端所有接口均为 `POST + JSON body`，成功返回 protojson 编码的 protobuf message，失败返回 `{"code": "...", "message": "..."}`。

`APIClient` 核心设计：
- 统一 POST 方法，支持泛型解码 `Decodable`
- `auth` 参数控制是否注入 `Bearer token`
- JWT token 存储在 Keychain，`AuthManager` 为 `@Observable` 类型驱动登录状态
- 图片上传使用 `multipart/form-data`（field: `uploadFile`）或 OSS 预签名 PUT

### 关键技术选型

| 需求 | 方案 | 说明 |
|------|------|------|
| Markdown 渲染 | [swift-markdown-ui](https://github.com/gonzalezreal/swift-markdown-ui) | SPM 引入，GitHub Flavored Markdown |
| Markdown 编辑 | `UITextView` + 自定义工具栏 | 通过 UIViewRepresentable 包装，底部工具栏提供快捷按钮 |
| Emoji 选择 | 自定义 EmojiPicker | 分组 emoji 网格，点击插入文本 |
| 地图展示 | MapKit（iOS 17+ API） | `Annotation` + `Marker`，自定义标注样式 |
| 图片选择 | PhotosPicker（系统） | 支持 SwiftUI 原生多选 |
| 图片缓存 | Kingfisher（可选） | 或 AsyncImage + NSCache 手动实现 |
| Token 存储 | Keychain（Security framework） | 比 UserDefaults 更安全 |

### 外部依赖

| 库 | 用途 | 必须 |
|---|---|---|
| swift-markdown-ui | Markdown 渲染 | 是 |
| Kingfisher | 图片加载缓存 | 可选 |

---

## 分步实施计划

### 第一阶段：基础设施

搭建项目骨架、网络层、认证模块。

1. 创建 Xcode 项目（SwiftUI, iOS 17+, Personal Team 签名）
2. 搭建目录结构（Core/Services/Models/Features）
3. 实现 `APIClient`（POST + JSON + Bearer token 注入）
4. 实现 `KeychainService` + `AuthManager`（登录状态管理）
5. 实现 `APIError` 统一错误处理
6. 定义所有数据模型（Codable struct，对应 proto message）
7. 实现 `ContentView`（TabView 4 Tab 骨架：博文/心情/提醒/足迹）
8. 实现 `LoginView` + ViewModel（用户名密码登录）
9. App 启动流程：检查 token 有效性 -> 登录页或主页

### 第二阶段：博文模块

博文列表、详情（Markdown 渲染）、评论、编辑器。

1. 实现 `ArticleService` + `CategoryService`
2. 集成 swift-markdown-ui（SPM）
3. 文章列表页：`List` + 下拉刷新 + 分页加载
4. 文章详情页：Markdown 渲染 + 底部评论区
5. 评论列表 + 回复评论输入框
6. 文章编辑器：UITextView + Markdown 工具栏（标题/加粗/斜体/代码块/链接/图片）
7. 编辑器内图片上传：PhotosPicker -> UploadService -> 插入 `![](cdn_url)`

### 第三阶段：心情模块

心情列表、发布、Emoji 选择。

1. 实现 `MoodService`
2. 心情列表页：时间线样式，显示内容 + 时间
3. 自定义 `EmojiPicker`（分组 emoji 网格，点击插入）
4. 心情编辑页：多行文本输入 + emoji 面板

### 第四阶段：提醒模块

提醒管理 CRUD（纯管理操作，通知由后端飞书发送）。

**后端修改（并行进行）：**
- `RemindUpdateRequest` proto 添加 `status` 字段
- `RemindUpdate` service 添加 status 更新逻辑
- OSS presign 接口 content-type 修复

1. 实现 `RemindService`
2. 提醒列表页：按状态分组展示（ACTIVE/PENDING/DONE），左滑删除
3. 提醒编辑页：内容输入 + cron 表达式（提供预设快捷选项 + 手动输入）
4. 状态管理：支持 ACTIVE -> PENDING -> DONE 的状态流转操作

### 第五阶段：足迹模块

地图展示、足迹 CRUD、照片上传。

1. 实现 `FootprintService` + `UploadService`
2. 足迹地图页：MapKit 全屏地图 + 自定义 Annotation（marker_color）
3. 底部抽屉/列表页：足迹卡片（缩略图 + 名称 + 日期）
4. 地图选点组件：长按地图选择坐标 + 反向地理编码
5. 照片选择（PhotosPicker 多选）+ OSS 预签名并行上传
6. 足迹编辑页：整合选点 + 照片上传 + 信息填写

### 第六阶段：打磨优化

UI 统一、空状态、加载态、错误处理、图片缓存。

1. 统一 UI 风格（颜色主题、字体、间距）
2. 空状态页面（列表为空时引导）
3. 加载状态 + 骨架屏
4. 错误处理统一化（Alert 展示）
5. 图片缓存优化（Kingfisher 或 NSCache）
6. 下拉刷新 / 上拉加载完善
7. HTTPS 安全配置检查

---

## 注意事项

1. **protojson 编码：** 后端使用 `UseProtoNames: true, EmitUnpopulated: true`，Swift 模型 CodingKeys 必须用 snake_case，字段类型需处理零值
2. **OSS 预签名 Content-Type：** 后端当前写死 `text/plain;charset=utf8`，上传图片会签名校验失败，必须后端配合修改
3. **7 天重签：** 每 7 天需连接 Xcode 重新安装，可设置日历提醒或使用 AltStore 自动化
4. **提醒通知：** 由后端定时任务通过飞书发送，App 端仅做 CRUD 管理
5. **分页：** 后端统一使用 `page` 参数（从 1 开始），每页固定条数
