
.PHONY: build
build:
	npm run build

.PHONY: upload
upload:
	ossutil --recursive  cp ./dist/assets oss://fifsky/assets -f -e oss-cn-shanghai.aliyuncs.com -i ${ALIYUN_KEY} -k ${ALIYUN_SECRET}

.PHONY: docker
docker: build upload
	docker login --username=${ALIYUN_DOCKER_USERNAME} --password=${ALIYUN_DOCKER_PASSWORD} registry.cn-shanghai.aliyuncs.com
	docker build . --file Dockerfile --tag registry.cn-shanghai.aliyuncs.com/fifsky/blog-web
	docker push registry.cn-shanghai.aliyuncs.com/fifsky/blog-web