build:export CGO_ENABLED=0
build:
	go build -o blogapi main.go

install:
	go run main.go install

stop:
	cat blogapi.pid | xargs kill

start:
	./blogapi > /dev/null 2>&1 &

restart:
	cat blogapi.pid | xargs kill
	./blogapi > /dev/null 2>&1 &
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