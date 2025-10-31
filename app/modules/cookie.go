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
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
)

type CookieConfig struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	MaxAge   int
}

func getCookieName() string {
	return config.Cfg.GetString("base.jwt.storageName")
}

func getCookieDomain(c echo.Context) string {
	host := c.Request().Host
	_, err := net.LookupHost(host)
	if err != nil {
		return "localhost"
	}
	return host
}

func getCookieMaxAge() int {
	expiry := config.Cfg.GetInt("base.jwt.expiry") * 60 * 60
	return expiry
}

func getCookieConfig(c echo.Context) CookieConfig {
	// 根据请求判断是否使用Secure Cookie
	isSecure := c.Scheme() == "https"

	return CookieConfig{
		Name:     getCookieName(),
		Domain:   getCookieDomain(c),
		Path:     "/",
		Secure:   isSecure,
		HttpOnly: true,                 // 防止XSS攻击
		SameSite: http.SameSiteLaxMode, // 防止CSRF攻击
		MaxAge:   getCookieMaxAge(),
	}
}

// SetAuthCookie 设置Cookie
func SetAuthCookie(c echo.Context, token string, expires time.Time) {
	cookieConfig := getCookieConfig(c)

	cookie := new(http.Cookie)
	cookie.Name = cookieConfig.Name
	cookie.Value = token
	cookie.Path = cookieConfig.Path
	cookie.Domain = cookieConfig.Domain
	cookie.MaxAge = cookieConfig.MaxAge
	cookie.Secure = cookieConfig.Secure
	cookie.HttpOnly = cookieConfig.HttpOnly
	cookie.SameSite = cookieConfig.SameSite
	cookie.Expires = expires

	c.SetCookie(cookie)
}

// ClearAuthCookie 清除Cookie
func ClearAuthCookie(c echo.Context) {
	cookieConfig := getCookieConfig(c)

	cookie := new(http.Cookie)
	cookie.Name = cookieConfig.Name
	cookie.Value = ""
	cookie.Path = cookieConfig.Path
	cookie.Domain = cookieConfig.Domain
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0)
	cookie.Secure = cookieConfig.Secure
	cookie.HttpOnly = cookieConfig.HttpOnly
	cookie.SameSite = cookieConfig.SameSite

	c.SetCookie(cookie)
}
