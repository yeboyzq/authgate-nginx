// @作者: yeboyzq
// @更新: 2026-06-26
// @备注:

package log

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
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
	debug      bool   // 调试模式
	level      string // 日志等级
	filePath   string // 文件存放路径
	maxSize    int    // 单个文件最大大小(MB)
	maxBackups int    // 保留旧文件最大数量
	maxAge     int    // 保留旧文件最大天数
	compress   bool   // 是否压缩旧文件
	access     bool   // 是否输出访问日志
}

// Init 初始化日志系统
func Init(e *echo.Echo) echo.MiddlewareFunc {
	// 加载配置
	Conf := &LogConfig{
		debug:      config.Cfg.GetBool("base.debug"),
		level:      config.Cfg.GetString("base.log.level"),
		filePath:   config.Cfg.GetString("base.log.path"),
		maxSize:    config.Cfg.GetInt("base.log.maxsize"),
		maxBackups: config.Cfg.GetInt("base.log.maxbackups"),
		maxAge:     config.Cfg.GetInt("base.log.maxage"),
		compress:   config.Cfg.GetBool("base.log.compress"),
		access:     config.Cfg.GetBool("base.log.access"),
	}
	if Conf.filePath == "" {
		panic("日志组件初始化失败: 日志存储路径为空, 请配置base.log.path设置项.")
	}
	// 设置日志级别
	var level slog.Level
	var addSource bool
	if Conf.debug {
		level = slog.LevelDebug
		addSource = false
	} else {
		switch Conf.level {
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

	// 设置业务日志轮转(app.log)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(Conf.filePath, utils.AppFileName()+".log"),
		MaxSize:    Conf.maxSize,
		MaxBackups: Conf.maxBackups,
		MaxAge:     Conf.maxAge,
		Compress:   Conf.compress,
	}

	// 设置访问日志轮转(access.log)
	accessLumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(Conf.filePath, "access.log"),
		MaxSize:    Conf.maxSize,
		MaxBackups: Conf.maxBackups,
		MaxAge:     Conf.maxAge,
		Compress:   Conf.compress,
	}

	// 创建业务日志 slog handler
	var handler slog.Handler
	if Conf.debug {
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

	// 创建访问日志 slog handler(与业务日志保持一致的格式；debug时输出到Stdout+access.log)
	var accessHandler slog.Handler
	if Conf.debug {
		accessMultiWriter := io.MultiWriter(os.Stdout, accessLumberjackLogger)
		accessHandler = slog.NewTextHandler(accessMultiWriter, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	} else {
		accessMultiWriter := io.MultiWriter(accessLumberjackLogger)
		accessHandler = slog.NewJSONHandler(accessMultiWriter, &slog.HandlerOptions{
			AddSource: addSource,
			Level:     level,
		})
	}

	// 创建业务logger
	Logger = slog.New(handler)
	// 设置为默认logger
	slog.SetDefault(Logger)
	e.Logger = Logger
	Logger.Info("日志中间件初始化完成.")

	// 访问日志开关
	if !Conf.access {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				return next(c)
			}
		}
	}
	accessLogger := slog.New(accessHandler)
	return LogMiddleware(accessLogger)
}

// LogMiddleware 自定义中间件, 使用slog记录请求日志
func LogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			// 处理请求
			err := next(c)
			if err != nil {
				c.Echo().HTTPErrorHandler(c, err)
			}
			// 获取参数
			latency := time.Since(start)
			request := c.Request()
			response := c.Response()
			requestID := getRequestID(c)
			resp, err := getResponse(response)
			if err != nil {
				c.Logger().Error(err.Error())
			}
			// 记录日志
			logger.Info("http_request",
				slog.String("request_id", requestID),
				slog.String("method", request.Method),
				slog.String("uri", request.URL.Path),
				slog.Int("status", resp.Status),
				slog.String("latency", latency.String()),
				slog.String("host", request.Host),
				slog.String("bytes_in", request.Header.Get(echo.HeaderContentLength)),
				slog.Int64("bytes_out", resp.Size),
				slog.String("remote_ip", c.RealIP()),
				slog.String("user_agent", request.UserAgent()),
			)

			return nil
		}
	}
}

// getResponse 提取 echo.Response, 如果未实现 unwrap 接口则包装构造
func getResponse(resw http.ResponseWriter) (*echo.Response, error) {
	var resp *echo.Response
	var err error
	if r, unwrapErr := echo.UnwrapResponse(resw); unwrapErr != nil {
		err = errors.New("上下文中的ResponseWriter未实现unwrapper接口: " + unwrapErr.Error())
		resp = new(echo.Response)
		resp.ResponseWriter = resw
		resp.Status = -1
		resp.Size = -1
	} else {
		resp = r
	}
	return resp, err
}

// getRequestID 获取请求ID
func getRequestID(c *echo.Context) string {
	requestID := c.Request().Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return requestID
}
