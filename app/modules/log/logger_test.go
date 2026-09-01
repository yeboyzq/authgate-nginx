/*
Copyright (c) 2026 authgate-nginx
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
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
)

type plainResponseWriter struct {
	header http.Header
	status int
}

func (w *plainResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *plainResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (w *plainResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

// newBufferLogger 创建缓冲日志器。
func newBufferLogger() (*bytes.Buffer, *slog.Logger) {
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return buf, logger
}

// prepareLogConfig 准备日志配置。
func prepareLogConfig(t *testing.T, debug bool, level string, access bool) {
	t.Helper()
	cfg := viper.New()
	cfg.Set("base.debug", debug)
	cfg.Set("base.log.level", level)
	cfg.Set("base.log.path", t.TempDir())
	cfg.Set("base.log.maxsize", 10)
	cfg.Set("base.log.maxbackups", 1)
	cfg.Set("base.log.maxage", 1)
	cfg.Set("base.log.compress", false)
	cfg.Set("base.log.access", access)
	config.Cfg = cfg
}

// TestInit 测试日志初始化。
func TestInit(t *testing.T) {
	tests := []struct {
		name   string
		debug  bool
		level  string
		access bool
	}{
		{name: "debug模式", debug: true, level: "debug", access: false},
		{name: "非debug-debug级别", debug: false, level: "debug", access: false},
		{name: "非debug-warn级别", debug: false, level: "warn", access: false},
		{name: "非debug-error级别", debug: false, level: "error", access: false},
		{name: "非debug-默认级别", debug: false, level: "unknown", access: false},
		{name: "开启访问日志", debug: false, level: "info", access: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareLogConfig(t, tt.debug, tt.level, tt.access)
			e := echo.New()
			middleware := Init(e)
			require.NotNil(t, middleware)
			require.NotNil(t, Logger)

			called := false
			handler := middleware(func(c *echo.Context) error {
				called = true
				c.Response().Header().Set(echo.HeaderXRequestID, "rid-init")
				return c.String(http.StatusOK, "ok")
			})

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			require.NoError(t, handler(ctx))
			assert.True(t, called)
		})
	}
}

// TestLogMiddleware 测试访问日志中间件。
func TestLogMiddleware(t *testing.T) {
	buf, logger := newBufferLogger()
	e := echo.New()
	httpErrHandled := false
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		httpErrHandled = true
		_ = c.String(http.StatusBadRequest, err.Error())
	}

	middleware := LogMiddleware(logger)
	handler := middleware(func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderXRequestID, "rid-2")
		return assert.AnError
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set(echo.HeaderXRequestID, "rid-1")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	require.NoError(t, handler(ctx))
	assert.True(t, httpErrHandled)
	assert.Contains(t, buf.String(), "http_request")
	assert.Contains(t, buf.String(), "rid-1")
	assert.Contains(t, buf.String(), "/hello")
}

// TestGetResponseAndGetRequestID 测试响应与请求 ID 提取。
func TestGetResponseAndGetRequestID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	resp, err := getResponse(ctx.Response())
	require.NoError(t, err)
	assert.NotNil(t, resp)

	ctx.Response().Header().Set(echo.HeaderXRequestID, "rid-response")
	assert.Equal(t, "rid-response", getRequestID(ctx))

	ctx.Request().Header.Set(echo.HeaderXRequestID, "rid-request")
	assert.Equal(t, "rid-request", getRequestID(ctx))

	resp, err = getResponse(&plainResponseWriter{})
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, -1, resp.Status)
}

// TestGetContextAndRequestIDLocal 测试本地下文与请求 ID 获取。
func TestGetContextAndRequestIDLocal(t *testing.T) {
	assert.Equal(t, context.Background(), getContext(nil))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req = req.WithContext(context.WithValue(req.Context(), "key", "value"))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	assert.Equal(t, "value", getContext(ctx).Value("key"))

	ctx.Response().Header().Set(echo.HeaderXRequestID, "rid-local")
	assert.Equal(t, "rid-local", getRequestIDLocal(ctx))
}

// TestQuoteHelpers 测试日志便捷方法。
func TestQuoteHelpers(t *testing.T) {
	buf, logger := newBufferLogger()
	Logger = logger
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set(echo.HeaderXRequestID, "rid-helper")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	Info("info-message", ctx, slog.String("key", "value"))
	Debug("debug-message")
	Warn("warn-message")
	Error("error-message")
	Log(slog.LevelInfo, "log-message", ctx)
	assert.Contains(t, buf.String(), "info-message")
	assert.Contains(t, buf.String(), "rid-helper")
	assert.Contains(t, buf.String(), "log-message")
}

// TestExtractCtxAndArgsAndWithRequestID 测试参数提取与 request_id 注入。
func TestExtractCtxAndArgsAndWithRequestID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req = req.WithContext(context.WithValue(req.Context(), "ctx", "value"))
	req.Header.Set(echo.HeaderXRequestID, "rid-1")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	gotCtx, gotEchoCtx, args := extractCtxAndArgs([]any{ctx, slog.String("k", "v")})
	assert.Equal(t, "value", gotCtx.Value("ctx"))
	assert.Equal(t, ctx, gotEchoCtx)
	assert.Len(t, args, 1)

	gotCtx, gotEchoCtx, args = extractCtxAndArgs([]any{context.Background(), slog.String("k", "v")})
	assert.NotNil(t, gotCtx)
	assert.Nil(t, gotEchoCtx)
	assert.Len(t, args, 1)

	result := withRequestID(ctx, []any{slog.String("k", "v")})
	assert.Len(t, result, 2)

	result = withRequestID(ctx, []any{slog.String("request_id", "rid-1")})
	assert.Len(t, result, 1)
}

// TestLoggerOutputGuard 测试日志器保护条件。
func TestLoggerOutputGuard(t *testing.T) {
	buf, logger := newBufferLogger()
	Logger = logger
	Debug("guard-debug")
	assert.Contains(t, buf.String(), "guard-debug")

	Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	assert.NotPanics(t, func() {
		Debug("hidden-debug")
	})
}
