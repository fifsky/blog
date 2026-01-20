# My Blog

![Status](https://github.com/fifsky/blog/workflows/blog/badge.svg) [![codecov](https://codecov.io/gh/fifsky/blog/branch/master/graph/badge.svg?token=MG1D2J86R6)](https://codecov.io/gh/fifsky/blog)

https://fifsky.com/

## Feature
- Without ORM, without framework, using native net/http and database/sql
- Based on buf [protovalidate](https://github.com/bufbuild/protovalidate) validation
- Generate an http.Handler based on protobuf and [googleapis](https://buf.build/googleapis/googleapis)
- Generate an OpenAPI description
- Use slog to record request logs

## Development

### Web

- Node.js: >= 20 (see `.nvmrc`)
- Package manager: pnpm (see `web/package.json`)

```bash
cd web
pnpm install
pnpm dev
```
