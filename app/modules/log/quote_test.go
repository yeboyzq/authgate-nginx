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
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// TestPanicGuard 测试 Panic 在不同级别下的行为。
func TestPanicGuard(t *testing.T) {
	buf, logger := newBufferLogger()
	Logger = logger
	assert.PanicsWithValue(t, "panic-message", func() {
		Panic("panic-message")
	})
	assert.Contains(t, buf.String(), "panic-message")

	buf, logger = newBufferLogger()
	Logger = slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelError + 1}))
	assert.NotPanics(t, func() {
		Panic("hidden-panic")
	})
	assert.NotContains(t, buf.String(), "hidden-panic")
}

// TestExtractCtxAndArgsDefault 测试默认参数提取分支。
func TestExtractCtxAndArgsDefault(t *testing.T) {
	ctx, echoCtx, args := extractCtxAndArgs([]any{"plain", 1})
	assert.Equal(t, context.Background(), ctx)
	assert.Nil(t, echoCtx)
	assert.Len(t, args, 2)

	args = withRequestID(nil, []any{slog.String("k", "v")})
	assert.Len(t, args, 1)
}

// TestWithRequestIDFallback 测试 request_id 回退与去重。
func TestWithRequestIDFallback(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/quote", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Response().Header().Set(echo.HeaderXRequestID, "rid-fallback")

	args := withRequestID(ctx, []any{slog.String("k", "v")})
	assert.Len(t, args, 2)
	assert.Contains(t, fmt.Sprint(args...), "rid-fallback")

	args = withRequestID(ctx, []any{slog.String("request_id", "rid-fallback")})
	assert.Len(t, args, 1)
}

// TestFatalSubprocess 测试 Fatal 会退出子进程。
func TestFatalSubprocess(t *testing.T) {
	if os.Getenv("TEST_LOG_FATAL_SUBPROCESS") == "1" {
		_, logger := newBufferLogger()
		Logger = logger
		Fatal("fatal-message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run", "TestFatalSubprocess")
	cmd.Env = append(os.Environ(), "TEST_LOG_FATAL_SUBPROCESS=1")
	err := cmd.Run()
	assert.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode())
}
