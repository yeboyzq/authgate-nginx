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

package log

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

// EchoLoggerAdapter 包装slog.Logger以实现Echo的Logger接口
type EchoLoggerAdapter struct {
	Logger *slog.Logger
}

// Write 实现 io.Writer 接口
func (s *EchoLoggerAdapter) Write(p []byte) (n int, err error) {
	s.Logger.Info(string(p))
	return len(p), nil
}

func (s *EchoLoggerAdapter) Error(err error) {
	s.Logger.Error(err.Error())
}

// CustomErrorHandler 自定义错误处理
func CustomErrorHandler(c echo.Context, err error) {
	switch e := err.(type) {
	case *echo.HTTPError:
		Warn(e.Error())
		c.JSON(e.Code, e.Message)
	default:
		Error(e.Error())
		c.JSON(http.StatusInternalServerError, e.Error())
	}
}
