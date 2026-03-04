# AuthGate-Nginx

> **Guard the Data Gateway. Define the Security Perimeter.**  
> **守护数据入口，定义安全边界**

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-Mulan&ensp;PSL&ensp;v2-green)]((https://opensource.org/license/mulanpsl-2-0))

AuthGate-Nginx 是一个基于 Go 语言构建的轻量级、高性能身份验证网关。它作为 Nginx 的 `auth_request` 模块的后端，为您的 Web 应用和服务提供集中的、基于 LDAP 的身份验证层。

## ✨ 特性

* 🔐 **集中式 LDAP 认证**：无缝集成您的企业 LDAP/Active Directory，实现统一的用户登录。
* 🚀 **高性能**：得益于 Go 语言的并发特性，提供极低的认证延迟。
* 🔒 **安全第一**：专注于安全边界，保护您的后端应用和数据入口。
* 🛠️ **与 Nginx 无缝集成**：通过标准的 `ngx_http_auth_request_module` 工作，配置简单。
* ⚙️ **易于配置**：通过简单的 YAML 配置文件和环境变量进行管理。
* 📦 **轻量级部署**：可编译为单一静态二进制文件，便于容器化（Docker）部署。

## 🚀 快速开始

### 前提条件

* Go 1.25 或更高版本
* 一个运行的 LDAP 服务器（例如 OpenLDAP, Active Directory）
* Nginx（已启用 `--with-http_auth_request_module`）

### 安装

1. **使用 Go 安装**

    ```bash
    go install github.com/yeboyzq/authgate-nginx@latest
    ```

2. **从源码构建**

    ```bash
    git clone https://github.com/yeboyzq/authgate-nginx.git
    cd authgate-nginx
    env GOOS=linux GOARCH=amd64 go build -o ./build/authgate-nginx ./app
    # 或
    env GOOS=windows GOARCH=amd64 go build -o ./build/authgate-nginx.exe ./app
    ```

3. **使用 Docker**

    ```bash
    docker run -d \
      -v $(pwd)/config.yaml:/data/apps/conf.d/config.yaml \
      -p 8000:8000 \
      yeboyzq/authgate-nginx:latest
    ```

### 基础配置

1. **创建配置文件** `config.yaml`：

    ```yaml
    base:
        server:
            protocol: http
            domain: 127.0.0.1
            port: 8000
        jwt:
            secret: xxxx
        ldap:
            url: ldap://localhost:389
            bindDn: cn=admin,dc=example,dc=com
            bindPassword: admin_password
            userBaseDn: ou=users,dc=example,dc=com
            filter: (uid=%s)
    ......
    ```

2. **配置 Nginx**：

    在您需要保护的 `location` 块中，添加 `auth_request` 指令。

    ```nginx
    server {
        listen  443 ssl;
        server_name  www.xxxx.com;
        access_log  /data/logs/nginx/www_xxxx_com-access.log main;
        ......

        # 认证设置
        location  = /nginx/api/verify {
            internal;
            proxy_pass  http://127.0.0.1:8000;
            proxy_pass_request_body  off;
            proxy_set_header  Content-Length "";
            proxy_set_header  Cookie $http_cookie;
            proxy_set_header  X-Original-URI $request_uri;
            proxy_set_header  X-Original-Method $request_method;
        }
        location  = /nginx/login {
            proxy_pass  http://127.0.0.1:8000;
            proxy_set_header  Host $host;
            proxy_set_header  X-Real-IP $remote_addr;
            proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header  X-Forwarded-Proto $scheme;
        }
        location  ~ ^/nginx {
            proxy_pass  http://127.0.0.1:8000;
            proxy_pass_header  Set-Cookie;
            proxy_set_header  Cookie $http_cookie;
            proxy_set_header  Host $host;
            proxy_set_header  X-Real-IP $remote_addr;
            proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header  X-Forwarded-Proto $scheme;
        }

        location  / {
            # 指定认证URL
            auth_request  /nginx/api/verify;
            auth_request_set  $auth_user $upstream_http_x_auth_user;
            proxy_set_header  Cookie $http_cookie;
            proxy_set_header  X-Auth-User $auth_user;
            # 如果身份验证服务返回401则重定向到登录页面
            error_page  401 =302 /nginx/login?redirect=$request_uri;

            proxy_pass  http://127.0.0.1:8080/;
            ......
        }
    }
    ```

3. **运行 AuthGate-Nginx**：

    ```bash
    ./authgate-nginx start -config config.yaml
    ```

现在，访问 `your-app.example.com` 将会通过 AuthGate-Nginx 进行 LDAP 身份验证。

## ⚙️ 详细配置

AuthGate-Nginx 支持通过配置文件和环境变量进行灵活配置。

### 配置文件示例 (YAML)

[config-example.yaml](./config-example.yaml)

### 环境变量

所有配置都可以通过环境变量覆盖，格式为 `APP_BASE_SITENAME`。

```bash
export APP_BASE_SITENAME="AuthGate-Nginx"
export APP_BASE_SERVER_PORT=8000
```

## 🏗️ 架构概述

```text
+----------+      +-------------+      +-----------------+      +-------------+
|  Client  | ---> |   Nginx     | ---> | AuthGate-Nginx  | ---> |  LDAP/AD    |
| (Browser)|      | (Frontend)  |      |  (Auth Gateway) |      |   Server    |
+----------+      +-------------+      +-----------------+      +-------------+
                       ^                                              |
                       | (X-User Header)                              |
                       |                                              |
                +-------------+                              (验证用户凭证)
                | Backend App |
                +-------------+
```

1. 用户请求访问受保护的资源。
2. Nginx 的 `auth_request` 模块向 AuthGate 的 `/validate` 端点发起子请求。
3. AuthGate 从请求头（如 `Authorization: Basic ...`）中提取凭据。
4. AuthGate 使用这些凭据绑定到 LDAP 服务器进行验证。
5. 如果认证成功，AuthGate 返回 `200 OK` 并附带用户信息头；失败则返回 `401 Unauthorized`。
6. Nginx 根据响应决定是继续代理请求到后端应用还是拒绝访问。

## 🤝 如何贡献

我们欢迎并感谢所有的贡献！请随时提交 Issue 或 Pull Request。

1. Fork 本仓库
2. 创建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目基于 `木兰宽松许可证， 第2版` 开源 - 查看 LICENSE 文件了解详情。

___

**AuthGate-Nginx** - Guard the Data Gateway. Define the Security Perimeter.
