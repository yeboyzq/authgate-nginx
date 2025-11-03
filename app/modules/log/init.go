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
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/utils"

	"github.com/labstack/echo/v5"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *slog.Logger

// LogConfig 日志配置
type LogConfig struct {
	Debug      bool   // 调试模式
	Level      string // 日志等级
	FilePath   string // 文件存放路径
	MaxSize    int    // 单个文件最大大小(MB)
	MaxBackups int    // 保留的旧文件最大数量
	MaxAge     int    // 保留旧文件的最大天数
	Compress   bool   // 是否压缩旧文件
}

// Init 初始化日志系统
func Init(e *echo.Echo) echo.MiddlewareFunc {
	// 加载配置
	Conf := &LogConfig{
		Debug:      config.Cfg.GetBool("base.debug"),
		Level:      config.Cfg.GetString("base.log.level"),
		FilePath:   config.Cfg.GetString("base.log.path"),
		MaxSize:    config.Cfg.GetInt("base.log.maxsize"),
		MaxBackups: config.Cfg.GetInt("base.log.maxbackups"),
		MaxAge:     config.Cfg.GetInt("base.log.maxage"),
		Compress:   config.Cfg.GetBool("base.log.compress"),
	}
	// 设置日志级别
	var level slog.Level
	var addSource bool
	if Conf.Debug {
		level = slog.LevelDebug
		addSource = false
	} else {
		switch Conf.Level {
		case "debug":
			level = slog.LevelDebug
			addSource = false
		case "warn":
			level = slog.LevelWarn
			addSource = false
		case "error":
			level = slog.LevelError
			addSource = false
		default:
			level = slog.LevelInfo
			addSource = false
		}
	}

	// 设置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(Conf.FilePath, "authgate-nginx.log"),
		MaxSize:    Conf.MaxSize,
		MaxBackups: Conf.MaxBackups,
		MaxAge:     Conf.MaxAge,
		Compress:   Conf.Compress,
	}

	// 创建 slog handler
	var handler slog.Handler
	if Conf.Debug {
		multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	} else {
		multiWriter := io.MultiWriter(lumberjackLogger)
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	}

	// 创建 logger
	Logger = slog.New(handler)

	// 设置为默认 logger
	slog.SetDefault(Logger)
	e.Logger = &EchoLoggerAdapter{Logger: Logger}
	e.HTTPErrorHandler = CustomErrorHandler
	Logger.Info("日志中间件初始化完成.")

	return LogMiddleware(Logger)
}

// LogMiddleware 自定义中间件，使用 slog 记录请求日志
func LogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// 处理请求
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// 记录日志
			latency := time.Since(start)
			request := c.Request()
			response := c.Response()
			requestID := utils.GetRequestID(c)

			logger.Info("http_request",
				"request_id", requestID,
				"method", request.Method,
				"uri", request.URL.Path,
				"status", response.Status,
				"latency", latency.String(),
				"host", request.Host,
				"remote_ip", c.RealIP(),
				"user_agent", request.UserAgent(),
			)

			return nil
		}
	}
}
