FROM registry.cn-shanghai.aliyuncs.com/fifsky/golang:1.25 AS builder
ARG PROJECT_NAME
WORKDIR /build
ENV GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE="" \
    GOSUMDB="sum.golang.google.cn"
COPY go.mod go.sum ./
COPY api/go.mod api/go.sum ./api/
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM registry.cn-shanghai.aliyuncs.com/fifsky/alpine:3.23
WORKDIR /app
COPY --from=builder /build/app ./
ENV TZ=Asia/Shanghai
ENTRYPOINT ["/app/app", "job"]
