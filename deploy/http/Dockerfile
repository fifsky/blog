FROM registry.cn-shanghai.aliyuncs.com/fifsky/golang:1.25 AS builder
WORKDIR /build
ENV GOPROXY=https://goproxy.verystar.net,https://goproxy.cn,direct \
    GOPRIVATE="" \
    GONOSUMDB="git.verystar.cn" \
    GOSUMDB="sum.golang.google.cn"
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM registry.cn-shanghai.aliyuncs.com/fifsky/alpine:3.23
WORKDIR /app
COPY --from=builder /build/app ./
EXPOSE 80
ENTRYPOINT ["/app/app", "openapi", "--addr=:80"]
