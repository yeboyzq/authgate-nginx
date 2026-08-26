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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	AppStartTime time.Time                         // 应用启动时间
	version      string    = "dev"                 // 应用版本号("v0.0.1")
	build_info   string    = ""                    // 应用编译信息
	build_time   string    = "2001-01-01T00:00:00" // 应用编译时间(编译不传值时故意写错赋值为系统时间, 实际"2001-01-01T00:00:00+08:00")
)

// AppVersionInfo 定义版本信息
type AppVersionInfo struct {
	AppVersion string
	BuildInfo  string
	BuildTime  time.Time
	Compiler   string
}

var VersionInfo = &AppVersionInfo{
	AppVersion: version,
	BuildInfo:  build_info,
	BuildTime:  timeFormat(build_time),
	Compiler:   runtime.Version(),
}

// PrintAppVersionInfo 打印应用版本信息
func PrintAppVersionInfo() {
	fmt.Printf("Version: %s, Build Info: %s, Build Time: %s, Compiler By: %s\n", VersionInfo.AppVersion, VersionInfo.BuildInfo, VersionInfo.BuildTime, VersionInfo.Compiler)
}

// AppFileName 获取可执行文件名(不含扩展名)
func AppFileName() string {
	executablePath, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	fileName := filepath.Base(executablePath)
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

// timeFormat 格式化时间字符串
func timeFormat(str string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Now().Local()
	}
	return parsedTime
}
