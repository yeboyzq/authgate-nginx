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

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	DefaultHttpTimeout        time.Duration = 5 * time.Second
	DefaultHttpRetryCount     int           = 3
	DefaultTransactionTimeout time.Duration = 30 * time.Second
)

var UserAgentHeader = fmt.Sprintf("go-%s/%s", AppFileName(), VersionInfo.AppVersion)

// GetEchoContext 通过echo.Echo获取echo.Context, 禁止传递到异步任务中使用
func GetEchoContext(e *echo.Echo, method string, path string, body any) (c echo.Context) {
	// 创建新的Echo上下文并设置请求上下文
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	c = e.NewContext(req, rec)
	defer e.ReleaseContext(c)
	return c
}

// CopyEchoContext 通过echo.Context复制获取echo.Context, 主要用于异步任务中使用
func CopyEchoContext(c echo.Context) (ctx echo.Context) {
	if c == nil || c.Echo() == nil {
		panic("CopyEchoContext: echo.Context对象不能为空")
	}

	// 复制原始请求
	originalReq := c.Request()
	// 创建请求的浅拷贝, 避免修改原始请求
	copiedReq := originalReq.Clone(originalReq.Context())
	// 创建新的上下文, 使用复制的请求
	ctx = c.Echo().NewContext(copiedReq, c.Response().Writer)

	// 复制RequestID
	if requestID := GetRequestID(c); requestID != "" {
		ctx.Response().Header().Set(echo.HeaderXRequestID, requestID)
		ctx.Request().Header.Set(echo.HeaderXRequestID, requestID)
	}
	// 复制token
	if token := c.Get("token"); token != nil {
		ctx.Set("token", token)
	}
	// 复制ram
	if token := c.Get("RamDataPermissions"); token != nil {
		ctx.Set("RamDataPermissions", token)
	}
	return ctx
}

// GetContext 通过echo.Context获取context.Context
func GetContext(c echo.Context) (ctx context.Context) {
	if c == nil {
		return context.Background()
	}
	if c.Request() == nil {
		return context.Background()
	}
	if c.Request().Context() == nil {
		return context.Background()
	}
	return c.Request().Context()
}

// GetRequestID 获取请求ID
func GetRequestID(c echo.Context) string {
	requestID := c.Request().Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return requestID
}

// GetUriFormURL 通过url获取uri
func GetUriFormURL(urlStr string) (uri string, path string, err error) {
	if urlStr == "" {
		return "", "", errors.New("the parameter is empty")
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	if parsedURL.RawQuery == "" {
		uri = parsedURL.Path
	} else {
		uri = parsedURL.Path + "?" + parsedURL.RawQuery
	}
	return uri, parsedURL.Path, nil
}
