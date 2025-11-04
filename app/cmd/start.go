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

package cmd

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5/middleware"
	"github.com/yeboyzq/authgate-nginx/app/modules"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
	"github.com/yeboyzq/authgate-nginx/app/public"
	"github.com/yeboyzq/authgate-nginx/app/routers"
	"github.com/yeboyzq/authgate-nginx/app/templates"
	"github.com/yeboyzq/authgate-nginx/app/utils"

	"github.com/labstack/echo/v5"
	"github.com/spf13/cobra"
)

// 标志解析
var (
	debugFlag      bool
	serverPortFlag int
)

// web命令入口
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动" + utils.AppFileName() + "应用",
	Long:  "启动" + utils.AppFileName() + "应用",
	Run: func(cmd *cobra.Command, args []string) {
		StartMain()
	},
}

// StartMain 应用初始化
// main()
func StartMain() {
	// 实例初始化
	app := echo.New()
	utils.AppStartTime = time.Now().UTC()

	// 组件初始化
	// config.Init()
	app.Use(log.Init(app))
	modules.CacheInit()
	modules.JwtInit()
	modules.LdapInit()
	modules.WhiteListInit()
	app.Use(middleware.RequestID())

	// 加载路由
	routers.Init(app)
	// 加载静态文件目录
	public.Init(app)
	// 加载模板
	templates.Init(app)

	// 启动服务
	log.Info("初始化完成, 启动中...")
	if err := app.Start(":" + config.Cfg.GetString("base.server.port")); err != http.ErrServerClosed {
		log.Error(err.Error())
	}
}
