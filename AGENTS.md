## Project Overview

Full-stack blog application with Go backend, React frontend, and native iOS (SwiftUI) app.

## Backend (Go)

### Tech Stack

- **Go 1.26.0** - No framework, uses native `net/http`
- **Protobuf** - API definitions with buf (googleapis)
- **Validation** - buf protovalidate
- **Database** - MySQL with native `database/sql` (no ORM)
- **Logging** - slog
- **Testing** - standard Go testing + dbunit

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
- **When you need to adjust or query the database**, you can use the `mysql` command by connecting with the database DSN found in `config.yml`. **However, any SQL execution MUST be approved by the user first.**

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
- **iOS Deployment Target** - 17.0
- **Xcode 26.5** - 构建/归档工具
- **无第三方依赖** - 纯原生 Foundation + SwiftUI
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
- 地图相关使用 GCJ-02 坐标，通过 `Core/Extensions/CoordinateTransform.swift` 与 WGS-84 转换

## Important Notes

### Backend

- No ORM - use raw SQL with `database/sql`
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
- Use dbunit fixtures in `testdata/` directory
- Linter configuration: `.golangci.yml`

#### Table-Driven Tests

All Golang unit tests should prioritize the table-driven style to enhance readability and maintainability:
**Principles:**

- Each test case should have a clear `name` to describe the scenario
- Test cases should cover: success paths, failure paths, and boundary conditions
- Use `t.Run()` to make each sub-test independently executable

### General

- Commit messages should follow conventional commit format
- Corresponding comments should be added to functions and structure fields. Comments are required for complex logic, while they can be omitted for simple logic. Only Chinese should be used for comments.
- Keep functions small and focused
- Handle errors explicitly, don't swallow them
