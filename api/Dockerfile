FROM alpine:3.11
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && \
  apk add --no-cache ca-certificates \
   tzdata

WORKDIR /blog/

COPY ./blogapi ./config.yml ./.env ./

ENV TZ=Asia/Shanghai
EXPOSE 80

ENTRYPOINT ["./blogapi","http","--addr=:80"]