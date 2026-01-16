FROM registry.cn-shanghai.aliyuncs.com/fifsky/node:24 AS builder
WORKDIR /build
COPY ./web/ ./
ENV CI=true
RUN pnpm install && pnpm build

FROM registry.cn-shanghai.aliyuncs.com/fifsky/nginx:1.29-alpine
COPY --from=builder /build/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build/dist/ /usr/share/nginx/html/
