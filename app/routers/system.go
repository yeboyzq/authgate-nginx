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

package routers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/utils"
)

// GetSysStatus 获取系统运行状态信息
type SysStatus struct {
	StartTime    string
	NumGoroutine int

	// 一般统计数据
	MemAllocated string // 已分配且仍在使用的字节
	MemTotal     string // 分配的字节数(即使已释放)
	MemSys       string // 从系统获取的字节(下面的XxxSys之和)
	Lookups      uint64 // 指针查找次数
	MemMallocs   uint64 // malloc数量
	MemFrees     uint64 // 释放次数

	// 主要分配堆统计信息
	HeapAlloc    string // 已分配且仍在使用的字节
	HeapSys      string // 从系统获取的字节
	HeapIdle     string // 空闲范围内的字节
	HeapInuse    string // 非空闲跨度中的字节
	HeapReleased string // 释放给操作系统的字节
	HeapObjects  uint64 // 分配的对象总数

	// Low-level fixed-size结构分配器统计信息
	// Inuse是现在使用的字节数
	// Sys是从系统获取的字节
	StackInuse  string // 引导堆栈
	StackSys    string
	MSpanInuse  string // mspan结构
	MSpanSys    string
	MCacheInuse string // mcache结构
	MCacheSys   string
	BuckHashSys string // 分析bucket哈希表
	GCSys       string // GC元数据
	OtherSys    string // 其他系统分配

	// 垃圾收集器统计信息
	NextGC       string // 下一次运行在HeapAlloc的时间(字节)
	LastGC       string // 最后一次运行的绝对时间(ns)
	PauseTotalNs string
	PauseNs      string // 最近GC暂停时间的循环缓冲区, 最近位于[(NumGC+255)%256]
	NumGC        uint32
}

// GetSystemStatus 更新获取系统运行状态
func GetSystemStatus(c echo.Context) (sysStatus SysStatus) {
	localLocation := time.Now().Location()
	AppStartTime := utils.AppStartTime.In(localLocation)
	sysStatus.StartTime = AppStartTime.Format(time.RFC3339)

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = utils.FileSize(int64(m.Alloc))
	sysStatus.MemTotal = utils.FileSize(int64(m.TotalAlloc))
	sysStatus.MemSys = utils.FileSize(int64(m.Sys))
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = utils.FileSize(int64(m.HeapAlloc))
	sysStatus.HeapSys = utils.FileSize(int64(m.HeapSys))
	sysStatus.HeapIdle = utils.FileSize(int64(m.HeapIdle))
	sysStatus.HeapInuse = utils.FileSize(int64(m.HeapInuse))
	sysStatus.HeapReleased = utils.FileSize(int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = utils.FileSize(int64(m.StackInuse))
	sysStatus.StackSys = utils.FileSize(int64(m.StackSys))
	sysStatus.MSpanInuse = utils.FileSize(int64(m.MSpanInuse))
	sysStatus.MSpanSys = utils.FileSize(int64(m.MSpanSys))
	sysStatus.MCacheInuse = utils.FileSize(int64(m.MCacheInuse))
	sysStatus.MCacheSys = utils.FileSize(int64(m.MCacheSys))
	sysStatus.BuckHashSys = utils.FileSize(int64(m.BuckHashSys))
	sysStatus.GCSys = utils.FileSize(int64(m.GCSys))
	sysStatus.OtherSys = utils.FileSize(int64(m.OtherSys))

	sysStatus.NextGC = utils.FileSize(int64(m.NextGC))
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC

	return sysStatus
}

func GetSysStatus(c echo.Context) error {
	rspData := GetSystemStatus(c)
	return c.JSON(http.StatusOK, rspData)
}
