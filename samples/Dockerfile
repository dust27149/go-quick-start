# syntax=docker/dockerfile:1.7

# Go 版本与 go.mod 保持一致
ARG GO_VERSION=1.26.5

# 第一阶段：在构建镜像中编译二进制
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

# 先复制依赖描述并下载依赖，提升缓存命中率
COPY go.mod ./
RUN go mod download

# 复制源码并编译静态二进制（可通过构建参数覆盖目标平台）
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/app .

# 第二阶段：最小运行时镜像，仅携带运行所需文件
FROM alpine:3.20 AS runtime
WORKDIR /app

# 设置容器时区为上海
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata \
	&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
	&& echo $TZ > /etc/timezone

# 创建非 root 用户，降低容器运行权限
RUN addgroup -S app && adduser -S app -G app
# 拷贝可执行文件与默认配置
COPY --from=builder /out/app /app/app
COPY config.json /app/config.json
# 预创建日志目录并授予应用用户权限
RUN mkdir -p /app/logs && chown -R app:app /app

# 使用 root 用户运行服务
USER root
# 对外暴露服务端口
EXPOSE 8080
# 容器启动命令
ENTRYPOINT ["/app/app"]
