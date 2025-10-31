# Copyright (c) 2025 authgate-nginx
# authgate-nginx is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#         http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.


# 编译镜像

# 依赖镜像
FROM golang:1.24-bookworm AS builder

# 设置标签
LABEL stage=gobuilder

# 定义变量
ARG APP="authgate-nginx"


# 设置编译环境变量
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

# 设置工作目录
WORKDIR /data/builds

# 处理依赖
COPY . .
RUN go mod download

# 运行编译
RUN go build -v -o ./build/${APP} ./app


# ----------------
# 运行镜像

# 依赖镜像
FROM debian:12-slim

# 设置标签
LABEL author="yeboyzq"

# 设置环境变量
ENV LANG=C.UTF-8
ENV RUN_USER="appuser"

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 安装必备工具
RUN apt-get update
RUN apt-get install -y curl telnet iputils-ping
RUN rm -rf /var/lib/apt/lists/*

# 创建运行用户
RUN groupadd -r $RUN_USER && useradd -r -g $RUN_USER $RUN_USER

# 设置工作目录
RUN mkdir -p /data/apps
WORKDIR /data/apps

# 挂载容器目录
VOLUME ["/data/apps/conf.d"]

# 拷贝可执行文件
COPY --from=builder /data/builds/build/${APP} /data/apps/${APP}

# 拷贝配置文件
COPY --from=builder /data/builds/config-example.yaml /data/apps/conf.d/config.yaml

# 设置目录权限
RUN chown -R $RUN_USER:$RUN_USER /data/apps && chmod -R 775 /data/apps

# 设置运行用户
USER $RUN_USER

# 设置健康检查
HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
    CMD curl -sSI 'http://localhost:$RUN_PORT/-/healthy' || exit 1

# 运行golang程序的命令
ENTRYPOINT ["/data/apps/${APP}", "start"]