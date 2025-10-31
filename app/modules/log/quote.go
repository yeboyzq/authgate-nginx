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
	"context"
	"log/slog"
)

// Log 提供简洁日志记录接口
func Log(level slog.Level, msg string, args ...any) {
	Logger.Log(context.Background(), level, msg, args...)
}

// Error 提供简洁日志记录接口
func Error(msg string, args ...any) {
	if Logger.Enabled(context.Background(), slog.LevelError) {
		Logger.Error(msg, args...)
	}
}

// Warn 提供简洁日志记录接口
func Warn(msg string, args ...any) {
	if Logger.Enabled(context.Background(), slog.LevelWarn) {
		Logger.Warn(msg, args...)
	}
}

// Info 提供简洁日志记录接口
func Info(msg string, args ...any) {
	if Logger.Enabled(context.Background(), slog.LevelInfo) {
		Logger.Info(msg, args...)
	}
}

// Debug 提供简洁日志记录接口
func Debug(msg string, args ...any) {
	if Logger.Enabled(context.Background(), slog.LevelDebug) {
		Logger.Debug(msg, args...)
	}
}
