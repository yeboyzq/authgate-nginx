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

package modules

import (
	"strings"

	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

var WhiteList *[]string

var sysWhiteList = []string{"/favicon.ico", "/-/healthy"}

// WhiteListInit 实例化白名单
func WhiteListInit() {
	var whiteList []string
	customWhiteList := config.Cfg.GetStringSlice("advanced.whitelist")
	if customWhiteList != nil || len(customWhiteList) != 0 {
		whiteList = append(sysWhiteList, customWhiteList...)
	} else {
		whiteList = sysWhiteList
	}
	WhiteList = &whiteList

	log.Info("白名单组件初始化完成.")
}

// CheckUrlWhiteList Url白名单管理
func CheckUrlWhiteList(path string) bool {
	if WhiteList != nil {
		for _, v := range *WhiteList {
			if strings.HasPrefix(path, strings.TrimRight(v, "*")) && strings.HasSuffix(v, "/*") {
				return true
			}
			if path == v {
				return true
			}
		}
		return false
	}
	return false
}
