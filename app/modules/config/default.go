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
	"github.com/spf13/viper"
)

// defaultConfig 设置默认配置
func defaultConfig(config *viper.Viper) {
	// base
	config.SetDefault("base.debug", false)
	config.SetDefault("base.siteName", "认证网关")
	config.SetDefault("base.copyright", "default")
	// config.SetDefault("base.root_url", "http://127.0.0.1:8023")
	// server
	config.SetDefault("base.server.protocol", "http")
	config.SetDefault("base.server.domain", "127.0.0.1")
	config.SetDefault("base.server.port", 8000)
	// jwt
	config.SetDefault("base.jwt.secret", "ECHO22bfkZ3tkYRvAW9eCpuou22FRAME")
	config.SetDefault("base.jwt.expiry", 8)
	config.SetDefault("base.jwt.storageName", "NginxAuthToken")
	// cache
	config.SetDefault("base.cache.maxsize", 128)
	// ldap
	config.SetDefault("base.ldap.url", "ldap://localhost:389")
	config.SetDefault("base.ldap.skipVerify", false)
	config.SetDefault("base.ldap.bindDn", "cn=admin,dc=example,dc=com")
	config.SetDefault("base.ldap.bindPassword", "admin_password")
	config.SetDefault("base.ldap.userBaseDn", "ou=users,dc=example,dc=com")
	config.SetDefault("base.ldap.filter", "(uid=%s)")
	// log
	config.SetDefault("base.log.path", "./logs")
	config.SetDefault("base.log.maxsize", 10)
	config.SetDefault("base.log.maxage", 7)
	config.SetDefault("base.log.maxbackups", 30)
	config.SetDefault("base.log.compress", true)
	config.SetDefault("base.log.level", "info")

	// advanced
	config.SetDefault("advanced.whitelist", "")
}
