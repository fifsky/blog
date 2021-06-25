FROM golang:1.16 as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct GOPRIVATE="" GOSUMDB="sum.golang.google.cn"
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

FROM alpine:3.11
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && \
  apk add --no-cache ca-certificates \
   tzdata

WORKDIR /app

COPY --from=builder /build/app ./
COPY --from=builder /build/storage/ ./storage/
COPY --from=builder /build/config.yml ./

ENV TZ=Asia/Shanghai
EXPOSE 80

ENTRYPOINT ["./app","http","--addr=:80","--env=prod"]