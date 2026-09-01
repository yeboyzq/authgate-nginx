# syntax=docker/dockerfile:1
# Copyright (c) 2026 authgate-nginx
# authgate-nginx is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#         http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.


# 如何构建
# amd64(本机, 最快)
# DOCKER_BUILDKIT=1 docker build -t <tag> .
# arm64 交叉(本机 load)
# docker buildx build --builder multiarch --platform linux/arm64 --load -t <tag> .
# 双架构 manifest(需推到 registry)
# docker buildx build --builder multiarch --platform linux/amd64,linux/arm64 -t <registry>/<image>:<tag> --push .


# 阶段一: 编译镜像

# builder 跑在构建机架构(BUILDPLATFORM), 通过 TARGETARCH 交叉编译, 兼容 buildx 多架构
FROM --platform=$BUILDPLATFORM golang:1.26 AS builder

LABEL stage=gobuilder

# 依赖包安装: CGO编译 + 交叉编译
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc libc6-dev \
        gcc-aarch64-linux-gnu libc6-dev-arm64-cross \
        git \
    && rm -rf /var/lib/apt/lists/*

# 目标平台参数(BuildKit 自动注入, 此处仅给默认值兜底非 buildx 场景)
ARG TARGETOS=linux
ARG TARGETARCH
# 应用名与模块路径(取自 go.mod)
ARG APP="authgate-nginx"
ARG MODULE="github.com/yeboyzq/authgate-nginx"

# 编译环境变量
ENV CGO_ENABLED=1 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /data/builds

# 优先拷贝依赖清单并下载, 利用层缓存加速重复构建
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 拷贝源码(含 .git, 用于读取版本信息)
COPY . .

# 编译: 执行编译脚本(参考 run_build.sh, 注入版本信息并按目标架构选择 CGO 交叉编译器)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH} MODULE=${MODULE} APP=${APP} \
    OUTPUT_FILE=/out/${APP} CGO_ENABLED=1 \
    bash docker/build.sh



# 阶段二: 运行镜像

FROM debian:stable-slim AS runtime

LABEL author="yeboyzq"

# 应用名
ARG APP="authgate-nginx"

# 运行时环境变量
# APP_CONF_FILE: 配置文件路径, 由 configmap 挂载注入, 镜像内不含配置文件
# TZ: 系统时区, 配合 tzdata 生效
ENV LANG=C.UTF-8 \
    TZ=Asia/Shanghai \
    APP_CONF_FILE=/data/apps/custom/conf/config.yaml

# 单层 RUN: 安装基础工具 + 创建非 root 用户 + 创建运行目录 + 设置时区 + 授权
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        curl telnet iputils-ping ca-certificates tzdata; \
    rm -rf /var/lib/apt/lists/*; \
    groupadd -r -g 1000 appuser; \
    useradd -r -u 1000 -g appuser -d /data/apps -s /usr/sbin/nologin appuser; \
    mkdir -p /data/apps/logs /data/apps/custom/conf /data/apps/custom/db; \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime; \
    echo 'Asia/Shanghai' > /etc/timezone; \
    chown -R appuser:appuser /data/apps

WORKDIR /data/apps

# 拷贝可执行文件与入口脚本(--chown 免去单独授权层)
COPY --from=builder --chown=appuser:appuser /out/${APP} /data/apps/authgate-nginx
COPY --chmod=0755 --chown=appuser:appuser docker/entrypoint.sh /data/apps/entrypoint.sh

# 应用端口: 8000=主服务, 48000=API文档服务
EXPOSE 8000 48000

# 以非 root 用户运行
USER appuser

# 健康检查: 主服务 /-/healthy 端点(已在鉴权白名单内), 使用 curl
# HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
#     CMD curl -fsS http://127.0.0.1:8000/-/healthy || exit 1

# 启动应用
ENTRYPOINT ["/data/apps/entrypoint.sh"]
