GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)

.PHONY: build
build: fmt-check
	$(GO) build -o app

.PHONY: run
run: build
	./app http --addr=:8080

.PHONY: generate
generate:
	$(GO) generate

.PHONY: test
test:
	$(GO) test -short -v -coverprofile=cover.out ./...
	grep -v "proto/gen" cover.out > cover.out.tmp && mv cover.out.tmp cover.out

.PHONY: cover
cover:
	$(GO) tool cover -func=cover.out -o cover_total.out
	$(GO) tool cover -html=cover.out -o cover.html

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi; \
	echo "\033[34m[Code] format perfect!\033[0m";

vet:
	$(GO) vet $(VETPACKAGES)

.PHONY: deploy
deploy:
	@script/make.sh deploy

.PHONY: proto
proto:
	cd proto && buf generate

.PHONY: docker
docker:
	docker login --username=${ALIYUN_DOCKER_USERNAME} --password=${ALIYUN_DOCKER_PASSWORD} registry.cn-shanghai.aliyuncs.com
	docker build . --file Dockerfile --tag registry.cn-shanghai.aliyuncs.com/fifsky/blog
	docker push registry.cn-shanghai.aliyuncs.com/fifsky/blog

.PHONY: lint
lint:
	golangci-lint run

.PHONY: buildui
buildui:
	cd web && pnpm run build && ossutil --recursive  cp ./dist/assets oss://fifsky/assets -f -e oss-cn-shanghai.aliyuncs.com -i ${FIFSKY_ALIYUN_KEY} -k ${FIFSKY_ALIYUN_SECRET}
