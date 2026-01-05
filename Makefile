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
	cd web && npm run build

.PHONY: upload
upload:
	cd web && ossutil --recursive  cp ./dist/assets oss://fifsky/assets -f -e oss-cn-shanghai.aliyuncs.com -i ${FIFSKY_ALIYUN_KEY} -k ${FIFSKY_ALIYUN_SECRET}
.PHONY: dockerui
dockerui: buildui upload
	docker login --username=${ALIYUN_DOCKER_USERNAME} --password=${ALIYUN_DOCKER_PASSWORD} registry.cn-shanghai.aliyuncs.com
	cd web && docker build . --file Dockerfile --tag registry.cn-shanghai.aliyuncs.com/fifsky/blog-web
	docker push registry.cn-shanghai.aliyuncs.com/fifsky/blog-web

.PHONY: shadui
shadui:
	cd web && pnpm dlx shadcn@latest create --preset "https://ui.shadcn.com/init?base=radix&style=nova&baseColor=neutral&theme=indigo&iconLibrary=lucide&font=inter&menuAccent=subtle&menuColor=default&radius=none&template=vite" --template vite