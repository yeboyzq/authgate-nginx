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

package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	DefaultBuildTime time.Time = timeFormat("2000-01-01")
	// 应用启动时间
	AppStartTime time.Time
)

// AppVersionInfo 定义版本信息
type AppVersionInfo struct {
	AppVersion      string
	DatabaseVersion int32
	BuildTime       time.Time
	BuiltBy         string
}

var VersionInfo = &AppVersionInfo{
	AppVersion:      "1.0.0",
	DatabaseVersion: 0,
	BuildTime:       timeFormat("2025-11-01"),
	BuiltBy:         runtime.Version(),
}

// AppFileName 获取可执行文件名
func AppFileName() string {
	// 获取编译后的文件路径
	executablePath, err := os.Executable()
	if err != nil {
		return "无法获取可执行文件名"
	}
	// 从路径中提取文件名
	fileName := filepath.Base(executablePath)
	return fileName
}

// timeFormat 格式化时间字符串
func timeFormat(str string) time.Time {
	layout := "2006-01-02 Z0700 MST"
	parsedTime, err := time.Parse(layout, str+" +0800 CST")
	if err != nil {
		return time.Now().Local()
	}
	return parsedTime
}
