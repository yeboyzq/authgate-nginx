// @作者: yeboyzq
// @更新: 2026-03-30
// @备注:

package log

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v5"
)

// extractCtxAndArgs 从可变参数中提取*echo.Context或context.Context并返回ctx、echo上下文以及清洗后的args
func extractCtxAndArgs(args []any) (context.Context, *echo.Context, []any) {
	ctx := context.Background()
	cleanArgs := args
	if len(args) == 0 {
		return ctx, nil, cleanArgs
	}
	// 仅识别第一个参数为上下文, 避免误判
	switch v := args[0].(type) {
	case *echo.Context:
		ctx = getContext(v)
		cleanArgs = args[1:]
		return ctx, v, cleanArgs
	case context.Context:
		ctx = v
		cleanArgs = args[1:]
		return ctx, nil, cleanArgs
	default:
		// 不处理
	}
	return ctx, nil, cleanArgs
}

// withRequestID 如果提供了echo上下文, 则自动附加request_id到参数列表
func withRequestID(c *echo.Context, args []any) []any {
	if c == nil {
		return args
	}
	requestID := getRequestIDLocal(c)
	if requestID == "" {
		return args
	}
	// 避免重复添加request_id
	for _, a := range args {
		if attr, ok := a.(slog.Attr); ok {
			if attr.Key == "request_id" {
				return args
			}
		}
	}
	return append(args, slog.String("request_id", requestID))
}

// Log 提供简洁日志记录接口
func Log(level slog.Level, msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	Logger.Log(ctx, level, msg, finalArgs...)
}

// Fatal 提供简洁日志记录接口
func Fatal(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelError) {
		Logger.ErrorContext(ctx, msg, finalArgs...)
		os.Exit(1)
	}
}

// Panic 提供简洁日志记录接口
func Panic(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelError) {
		Logger.ErrorContext(ctx, msg, finalArgs...)
		panic(msg)
	}
}

// Error 提供简洁日志记录接口
func Error(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelError) {
		Logger.ErrorContext(ctx, msg, finalArgs...)
	}
}

// Warn 提供简洁日志记录接口
func Warn(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelWarn) {
		Logger.WarnContext(ctx, msg, finalArgs...)
	}
}

// Info 提供简洁日志记录接口
func Info(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelInfo) {
		Logger.InfoContext(ctx, msg, finalArgs...)
	}
}

// Debug 提供简洁日志记录接口
func Debug(msg string, args ...any) {
	ctx, c, cleanArgs := extractCtxAndArgs(args)
	finalArgs := withRequestID(c, cleanArgs)
	if Logger.Enabled(ctx, slog.LevelDebug) {
		Logger.DebugContext(ctx, msg, finalArgs...)
	}
}

// getContext 本地上下文获取, 避免 shared/web 反向依赖导致的导入循环
func getContext(c *echo.Context) context.Context {
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

// getRequestIDLocal 本地获取 request_id, 避免 shared/web 反向依赖导致的导入循环
func getRequestIDLocal(c *echo.Context) string {
	requestID := c.Request().Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return requestID
}
