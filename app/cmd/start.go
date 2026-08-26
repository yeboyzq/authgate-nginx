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
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5/middleware"
	"github.com/yeboyzq/authgate-nginx/app/modules"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
	"github.com/yeboyzq/authgate-nginx/app/public"
	"github.com/yeboyzq/authgate-nginx/app/routers"
	"github.com/yeboyzq/authgate-nginx/app/templates"
	"github.com/yeboyzq/authgate-nginx/app/utils"
	"github.com/yeboyzq/authgate-nginx/app/utils/id"

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
	utils.PrintAppVersionInfo()

	// 组件初始化
	config.Init()
	app.Use(log.Init(app))
	modules.CacheInit()
	modules.JwtInit()
	modules.LdapInit()
	modules.WhiteListInit()

	// 全局中间件初始化
	app.Use(middleware.Recover())
	app.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: id.NewRequestID,
	}))
	app.Use(middleware.Gzip())

	// 加载路由
	routers.Init(app)
	// 加载静态文件目录
	public.Init(app)
	// 加载模板
	templates.Init(app)

	// 启动服务
	log.Info("初始化完成, 启动中...")
	// 创建优雅停机上下文
	gracefulCtx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	// 启动配置
	startConfig := echo.StartConfig{
		Address:         ":" + config.Cfg.GetString("base.server.port"),
		GracefulTimeout: 12 * time.Second,
		OnShutdownError: func(err error) {
			log.Error("关闭应用程序时出错", "error", err)
		},
	}
	// 在优雅停机上下文中启动一个goroutine来处理自定义清理工作
	go func() {
		<-gracefulCtx.Done()
		log.Info("开始优雅停机...")

		// 清理资源准备关闭
		log.Info("完成应用程序关闭前资源清理.")
	}()
	// log.Info(config.Cfg.GetString("base.appname") + ": " + config.Cfg.GetString("base.description") + "(" + config.Cfg.GetString("base.server.root_url") + ")")
	// 启动服务器
	if err := startConfig.Start(gracefulCtx, app); err != nil && err != http.ErrServerClosed {
		log.Fatal("无法启动服务器", err)
	}
}
