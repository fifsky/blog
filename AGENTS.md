## Project Overview

Full-stack blog application with Go backend, React frontend, and native iOS (SwiftUI) app.

## Backend (Go)

### Tech Stack

- **Go 1.26.0** - No framework, uses native `net/http`
- **Protobuf** - API definitions with buf (googleapis)
- **Validation** - buf protovalidate
- **Database** - SQLite (`modernc.org/sqlite`，纯 Go 无 CGO) with native `database/sql` (no ORM)
- **Backup** - Litestream（Go library 模式嵌入，实时流式备份到阿里云 OSS）
- **Logging** - slog
- **Testing** - standard Go testing + dbunit（SQLite 适配版）

### Commands

```bash
# Build backend
make build

# Run all tests
make test

# Run single test
go test -v -run TestUser_Login ./service

# Check test coverage
make cover

# Format code
make fmt

# Lint code
make lint

# Generate protobuf code
make proto

# Build frontend
make buildui

# Type check frontend (ignore node_modules)
cd web && pnpm tsc --noEmit --skipLibCheck

# Development
pnpm dev

# Run backend
make run
```

### Code Style - Backend

**Imports:** Standard library first, third-party packages last (grouped with blank line)

```go
import (
    "context"
    "fmt"

    "app/config"
    "app/store"

    "github.com/golang-jwt/jwt/v5"
)
```

**Error Handling:**

- Use `fmt.Errorf()` for wrapping errors
- Use `response.Fail(w, code, msg)` for HTTP error responses
- Return `sql.ErrNoRows` for not found cases

**Naming Conventions:**

- Interfaces: `Store`, `UserServiceServer`
- Structs: `User`, `Post`
- Methods: `GetPost(ctx, id)`, `Login(ctx, req)`
- Test functions: `TestUser_Login(t *testing.T)`
- Private members: lowercase, public: PascalCase

**Database:**

- Always use `QueryContext(ctx, query, args...)` or `QueryRowContext(ctx, query, args...)` for queries
- Always `defer rows.Close()` after creating rows
- Use context throughout: `ctx context.Context` as first param
- Return errors directly, don't panic
- **SQLite DSN**: `file:storage/blog.db?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)` — PRAGMA 必须通过 DSN 参数设置，不能用 `ExecContext`（连接池下只影响随机连接）
- **SQLite SQL 语法注意事项**:
  - UPSERT: `INSERT INTO ... ON CONFLICT(col) DO UPDATE SET col = excluded.col`
  - 日期格式化: `strftime('%Y/%m', substr(col, 1, 19))`（含时区的 DATETIME 需用 `substr` 截取前 19 字符）
  - JSON 数组查询: `EXISTS (SELECT 1 FROM json_each(col) WHERE value = ?)`
  - 随机排序: `ORDER BY RANDOM()`
  - 分页: `LIMIT ? OFFSET ?`
  - JSON 列扫描: 用 `[]byte` 接收（SQLite 驱动返回 string，不能直接 scan 到 `json.RawMessage`）
- **Litestream 备份**: 作为 Go library 嵌入（`pkg/litestream/manager.go`），启动时自动从 OSS 恢复（若本地无 DB），关闭时优雅停止。备份路径按环境区分：开发 `blog/sqlite/local`，线上 `blog/sqlite/prod`
- **`updated_at` 自动更新**: 通过 `testdata/schema.sql` 中的 SQLite 触发器实现，`PRAGMA recursive_triggers` 默认 OFF 不会递归

**Protobuf:**

- Generate with `make proto`
- Generated code excluded from golangci-lint (see .golangci.yml)

## Frontend (React + TypeScript)

### Tech Stack

- **React 19**
- **TypeScript 5.9**
- **TailwindCSS 4**
- **React Router 7** (NEW: use `react-router`, NOT `react-router-dom`)
- **React Hook Form + Zod** (forms & validation)
- **Zustand** (state management)
- **Shadcn UI** (UI components)
- **dayjs** (dates)
- **pnpm 10.27.0** (package manager)

### Commands

```bash
cd web

# Development server
pnpm dev

# Build for production
pnpm build

# Format code
pnpm format

# Type check
pnpm tsc --noEmit --skipLibCheck
```

### Code Style - Frontend

**Imports:**

- React/Router hooks first, then internal imports with `@` alias
- No `react-router-dom` - use `react-router` instead

```typescript
import { lazy, Suspense } from "react";
import { useNavigate } from "react-router";
import { CArticle } from "@/components/CArticle";
import { articleListApi } from "@/service";
```

**Component Structure:**

- Functional components with hooks
- Use TypeScript interfaces/types for props
- Lazy load route components with `lazy()` + `Suspense`

**State Management:**

- Local state: `useState`, `useReducer`
- Global state: Zustand stores in `@/store`
- Forms: React Hook Form with Zod schemas

**API Requests:**

- Use `request<T>()` wrapper from `@/utils/request`
- Use `createApi<T>()` for authenticated requests
- Errors: `AppError` class with code/message
- Use `dialog.message()` for error notifications

**Styling:**

- TailwindCSS 4 utility classes
- Use `@tailwindcss/vite` plugin
- Import from `@/components/ui` for Shadcn components
- Custom styles in `@layer base` in `index.css`

**Path Aliases:**

- `@/*` → `./src/*`
- `@/components` → `./src/components`
- `@/lib` → `./src/lib`
- `@/hooks` → `./src/hooks`
- `@/utils` → `./src/utils`

**Naming Conventions:**

- Components: PascalCase (e.g., `ArticleList`, `CArticle`)
- Hooks: camelCase with `use` prefix (e.g., `useDialog`)
- Utils: camelCase (e.g., `getApiUrl`)
- Types/Interfaces: PascalCase (e.g., `ArticleItem`, `ArticleListRequest`)
- Constants: UPPER_SNAKE_CASE

## iOS (SwiftUI)

### Tech Stack

- **Swift 5.0** / **SwiftUI** - 原生 iOS App
- **iOS Deployment Target** - 18.0
- **Xcode 26.5** - 构建/归档工具
- **Textual** - Markdown 渲染（[gonzalezreal/textual](https://github.com/gonzalezreal/textual)，通过 SPM 集成）
- **架构** - MVVM（View + ViewModel + Service）
- **网络** - 原生 `URLSession`，`actor APIClient` 单例封装
- **鉴权** - Keychain 存储 token，`AuthManager` 管理
- **数据编码** - protojson（snake_case），`Codable` 直接映射

### Project Info

- **工程** - `ios/BlogApp.xcodeproj`（scheme: `BlogApp`）
- **Bundle ID** - `com.fifsky.blog`
- **版本** - 1.0 (1)
- **设备** - `TARGETED_DEVICE_FAMILY = "1,2"`（iPhone + iPad）
- **签名** - `CODE_SIGN_STYLE = Automatic`，`DEVELOPMENT_TEAM = ""`（无签名，配合 SideStore 侧载部署）
- **API 地址** - `https://api.fifsky.com`（见 `App/Config.swift`）

### Directory Structure

```
ios/
├── App/                  # 全局配置 Config.swift（所有 API 路径常量）
├── BlogApp.swift         # @main 入口
├── ContentView.swift     # 根视图（登录态切换）
├── Core/
│   ├── Auth/             # AuthManager（登录态）、KeychainService
│   ├── Network/          # APIClient 单例、APIError、MultipartFormData
│   └── Extensions/       # Date+Format、CoordinateTransform（GCJ-02/WGS-84 转换）
├── Models/               # Codable 数据模型（对应后端 proto 消息）
├── Services/             # 业务 API 封装（按模块拆分）
├── Features/             # 功能模块（MVVM）
│   ├── Article/          # 文章（列表/详情/编辑/评论）
│   ├── Mood/             # 心情
│   ├── Remind/           # 提醒
│   ├── Footprint/        # 足迹（含地图）
│   └── Login/            # 登录
├── Components/           # 通用组件（Theme、PageBackground、PhotoBrowserView 等）
├── Assets.xcassets/      # 背景图等资源
└── AppIconAssets.xcassets/ # App 图标
```

### Commands

```bash
cd ios

# 模拟器构建并运行（通过 MCP ios-simulator 工具）
# - mcp__ios_simulator__ios_build_and_run
# - mcp__ios_simulator__ios_screenshot（截图验证）

# 真机 IPA 归档（无签名，适配 SideStore）
xcrun xcodebuild archive \
  -project BlogApp.xcodeproj \
  -scheme BlogApp \
  -configuration Release \
  -destination "generic/platform=iOS" \
  -archivePath build/BlogApp.xcarchive \
  CODE_SIGN_IDENTITY="" \
  CODE_SIGNING_REQUIRED=NO \
  CODE_SIGNING_ALLOWED=NO \
  DEVELOPMENT_TEAM="" \
  AD_HOC_CODE_SIGNING_ALLOWED=NO

# 从 xcarchive 打包为 IPA
rm -rf build/Payload && mkdir -p build/Payload
cp -R build/BlogApp.xcarchive/Products/Applications/BlogApp.app build/Payload/
cd build && zip -qry BlogApp.ipa Payload
# 产物：ios/build/BlogApp.ipa（arm64，未签名）
```

### Code Style - iOS

**API 请求:**

- 所有接口统一走 `APIClient.shared`，POST + JSON Body
- 路径常量集中定义在 `App/Config.swift`（如 `Config.loginPath`）
- 业务封装在 `Services/` 下各 Service 类，返回 `async throws` 模型
- protojson 编码：模型属性直接用 `snake_case`，`int64` 字段为 `String` 类型

**架构约定 (MVVM):**

- `XxxView.swift` - SwiftUI 视图
- `XxxViewModel.swift` - `@MainActor` ObservableObject，业务逻辑
- `XxxService.swift` - API 调用封装
- 模型放 `Models/`，实现 `Codable` / `Identifiable`

**注释规范:**

- 注释仅使用中文
- 复杂逻辑必须注释，简单逻辑可省略
- 结构体字段、函数均需文档注释（`///`）

**命名约定:**

- 类型/协议：PascalCase（`AuthManager`、`ArticleListViewModel`）
- 属性/方法：camelCase（`baseURL`、`request()`）
- 常量：camelCase（Config 内的 `static let`）
- 私有成员：前置 `_` 或 camelCase

### Important Notes

- 修改 Swift 文件后需通过 Xcode/MCP 重新构建生效
- `Models/` 的字段须与后端 proto 定义（snake_case）严格对应
- 默认无签名，仅供 SideStore 侧载；如需 App Store/TestFlight，需配置 DEVELOPMENT_TEAM 和证书
- 文章详情正文通过 [Textual](https://github.com/gonzalezreal/textual) 的 `StructuredText(markdown:)` 渲染（GitHub Flavored Markdown，支持标题/列表/代码块/引用/表格），样式深浅色自动适配
- 地图相关使用 GCJ-02 坐标，通过 `Core/Extensions/CoordinateTransform.swift` 与 WGS-84 转换

## Important Notes

### Backend

- No ORM - use raw SQL with `database/sql`
- Database: SQLite (`modernc.org/sqlite`)，DSN 见 `config.yml`，Schema 定义在 `testdata/schema.sql`
- Litestream 实时备份到阿里云 OSS（`fifsky-backup` bucket），Go library 模式嵌入，无需独立进程
- K8s 部署使用 PVC 持久化 `/app/storage` 目录（SQLite 文件存储）
- Always pass context through the call chain
- Use `make fmt` before committing
- Generated protobuf code is in `proto/gen/` (excluded from lint),use `make proto` to generate

### Frontend

- Do NOT run `pnpm build` after each task - dev mode shows real-time compilation
- Always use `react-router` (v7), never `react-router-dom` (old version)
- Install Shadcn components: `pnpm dlx shadcn@latest add [component-name]`
- Use `pnpm` for all package operations (not npm/yarn)
- When adding UI components (visual/styling changes), delegate to frontend-ui-ux-engineer
- When adding logic to frontend files, handle directly

### Testing

- Backend: Must `make test` runs all tests, Use -short to skip some tests
- Single test: `go test -v -run TestName ./path/to/package`
- **Environment Variables**: Environment variables required for unit tests are defined in `.envrc`. When running unit tests from the command line, you must load these variables first:
  - `export $(cat .envrc | xargs) && go test ./...`
- Use dbunit fixtures in `testdata/` directory（SQLite 临时文件模式，每个测试用例独立数据库）
- Linter configuration: `.golangci.yml`

#### Table-Driven Tests

All Golang unit tests should prioritize the table-driven style to enhance readability and maintainability:
**Principles:**

- Each test case should have a clear `name` to describe the scenario
- Test cases should cover: success paths, failure paths, and boundary conditions
- Use `t.Run()` to make each sub-test independently executable

#### Assertion Principles

单元测试断言应避免对 fixture 数据的绝对值产生耦合，确保 fixture 变化时测试不会脆弱地失败：

**Count 断言:**

- 有数据场景：使用 `assert.Greater(t, count, 0)` 而非 `assert.Equal(t, N, count)`
- 无数据场景（空表、无匹配）：使用 `assert.Equal(t, 0, count)` — 这是确定性的
- 增减场景（创建/删除后）：使用相对比较 `assert.Equal(t, beforeTotal+1, afterTotal)` 或 `assert.Less(t, afterTotal, beforeTotal)` — 不依赖 fixture 绝对数量

**列表长度断言:**

- 有数据场景：使用 `assert.NotEmpty(t, list)` 而非 `assert.Len(t, list, N)`
- 分页场景：使用 `assert.LessOrEqual(t, len(ret), num)` 验证不超过每页数量，而非断言精确条数
- 空结果场景：使用 `assert.Empty(t, list)`

**分页边界用例:**

- 禁止编写 "第N页无数据" 这类依赖精确总数的用例（如 fixture 有 3 条，断言第 4 页返回 0）
- 分页测试仅验证有效页返回数据 + 每页数量上限

**排序断言:**

- 禁止断言特定 ID 在首位（如 `assert.Equal(t, 10, list[0].Id)`）— 依赖 fixture 具体数据
- 改为验证排序方向：遍历比较相邻元素的 `CreatedAt`（正序/倒序）

**字段值断言:**

- 按 ID 查询后验证字段值（如 `assert.Equal(t, "test", ret.Name)`）是合理的 — 测试读取正确性
- 但避免对列表中特定位置的元素断言具体值（依赖 fixture 排序），应先按 ID 查找再断言

**自建数据测试:**

- 测试中自行创建数据（空表 + Create）后断言精确值是合理的 — 数据完全可控
- 相对比较（`beforeTotal+1`）优于绝对值（`assert.Equal(t, 3, total)`）

### General

- Commit messages should follow conventional commit format
- Corresponding comments should be added to functions and structure fields. Comments are required for complex logic, while they can be omitted for simple logic. Only Chinese should be used for comments.
- Keep functions small and focused
- Handle errors explicitly, don't swallow them
