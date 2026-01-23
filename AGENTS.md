## Project Overview

Full-stack blog application with Go backend and React frontend.

## Backend (Go)

### Tech Stack

- **Go 1.25.5** - No framework, uses native `net/http`
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
- Use dbunit fixtures in `testdata/` directory
- Linter configuration: `.golangci.yml`

### General

- Commit messages should follow conventional commit format
- Corresponding comments should be added to functions and structure fields. Comments are required for complex logic, while they can be omitted for simple logic. Only Chinese should be used for comments.
- Keep functions small and focused
- Handle errors explicitly, don't swallow them
