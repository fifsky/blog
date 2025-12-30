GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)

.PHONY: build
build: fmt-check
	$(GO) build -o app

.PHONY: generate
generate:
	$(GO) generate

.PHONY: test
test:
	$(GO) test -short -v -coverprofile=cover.out ./...

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
	skaffold run -n pay --tail

.PHONY: docker
docker:
	docker login --username=${ALIYUN_DOCKER_USERNAME} --password=${ALIYUN_DOCKER_PASSWORD} registry.cn-shanghai.aliyuncs.com
	docker build . --file Dockerfile --tag registry.cn-shanghai.aliyuncs.com/fifsky/blog
	docker push registry.cn-shanghai.aliyuncs.com/fifsky/blog

.PHONY: lint
lint:
	golangci-lint run
