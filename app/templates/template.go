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

package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

//go:embed *.tmpl
var TemplateFS embed.FS

type Template struct {
	Templates *template.Template
}

// 模板初始化
func Init(e *echo.Echo) {
	e.Renderer = &Template{
		Templates: template.Must(template.ParseFS(TemplateFS, "*.tmpl")),
	}

	log.Info("应用模板初始化完成.")
}

// 模板操作
func (s *Template) Render(c *echo.Context, w io.Writer, name string, data interface{}) error {
	return s.Templates.ExecuteTemplate(w, name, data)
}
