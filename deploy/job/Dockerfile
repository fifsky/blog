FROM registry.verystar.net/library/golang:1.25 AS builder
ARG PROJECT_NAME
WORKDIR /build
ENV GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE="" \
    GONOSUMDB="git.verystar.cn" \
    GOSUMDB="sum.golang.google.cn"
COPY go.mod go.sum ./
COPY api/go.mod api/go.sum ./api/
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./internal/${PROJECT_NAME}


FROM registry.verystar.net/library/debian:12
WORKDIR /app
COPY --from=builder /build/app ./
EXPOSE 80
ENTRYPOINT ["/app/app", "job"]
