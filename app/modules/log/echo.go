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
	"net/http"

	"github.com/labstack/echo/v5"
)

// CustomErrorHandler 自定义错误处理
func CustomErrorHandler(c *echo.Context, err error) {
	switch e := err.(type) {
	case *echo.HTTPError:
		Warn(e.Error())
		c.JSON(e.Code, e.Message)
	default:
		Error(e.Error())
		c.JSON(http.StatusInternalServerError, e.Error())
	}
}
