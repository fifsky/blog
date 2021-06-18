.PHONY: build
build:
	CGO_ENABLED=0 go build -o blogapi main.go

.PHONY: run
run:
	go run main.go http --addr=:"8889"

.PHONY: test
test:
	go test ./...

.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blogapi main.go

.PHONY: docker
docker: linux
	docker login --username=${ALIYUN_DOCKER_USERNAME} --password=${ALIYUN_DOCKER_PASSWORD} registry.cn-shanghai.aliyuncs.com
	docker build . --file Dockerfile --tag registry.cn-shanghai.aliyuncs.com/fifsky/blog
	docker push registry.cn-shanghai.aliyuncs.com/fifsky/blog