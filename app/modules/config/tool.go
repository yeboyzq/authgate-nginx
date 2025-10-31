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

package config

import (
	"fmt"
	"time"
)

// GetSecretKey 获取加密密钥并截取
func GetSecretKey() string {
	key := Cfg.GetString("base.jwt.secret")
	if len(key) < 32 {
		panic("加密密钥少于32字符(base.jwt.secret)")
	}
	n := 32
	if n > len(key) {
		n = len(key)
	}
	subkey := key[:n]
	return subkey
}

// GetSiteName 获取站点名称
func GetSiteName() string {
	return Cfg.GetString("base.siteName")
}

// GetCopyright 获取版权信息
func GetCopyright() string {
	copyright := Cfg.GetString("base.copyright")
	if copyright == "default" {
		year := time.Now().Year()
		copyright = fmt.Sprintf(`Copyright © %v <a href="https://github.com/yeboyzq/authgate-nginx" target="_blank">AuthGate</a>. All rights reserved.`, year)
	}
	return copyright
}
