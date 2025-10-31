/*
Copyright (c) 2025 authgate-nginx
authgate-nginx is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
        http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/

package routers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

// Init 初始化路由
func Init(app *echo.Echo) {
	app.AddRoute(echo.Route{Method: http.MethodGet, Path: "/nginx/-/healthy", Handler: GetSysStatus, Name: "系统模块-健康检查-访问"})

	// nginx
	app.AddRoute(echo.Route{Method: http.MethodGet, Path: "/nginx/login", Handler: LoginPageHandler, Name: "认证模块-nginx-登录页"})
	app.AddRoute(echo.Route{Method: http.MethodPost, Path: "/nginx/api/login", Handler: LoginHandler, Name: "认证模块-nginx-登录"})
	app.Any("/nginx/api/verify", VerifyHandler)
	// app.AddRoute(echo.Route{Method: http.MethodGet, Path: "/nginx/api/verify", Handler: VerifyHandler, Name: "认证模块-nginx-认证"})
	// app.AddRoute(echo.Route{Method: http.MethodPost, Path: "/nginx/api/verify", Handler: VerifyHandler, Name: "认证模块-nginx-认证"})

	log.Info("应用路由初始化完成.")
}

// ResponseJSON 自定义响应的数据结构
type ResponseJSON struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data,omitempty"`
	RequestId string `json:"requestId"`
}
