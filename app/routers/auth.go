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
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/modules"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

type RequestUserInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type ResponseLoginInfo struct {
	Username  string    `json:"username"`
	Token     string    `json:"token"`
	ExpiresIn time.Time `json:"expires_in"`
	Redirect  string    `json:"redirect"`
}

// LoginPageHandler 显示登录页面
func LoginPageHandler(c echo.Context) error {
	// 如果已经登录，重定向到原始页面或首页
	tokenString := modules.Jwt.ExtractToken(c)
	redirectURL := c.QueryParam("redirect")
	log.Info("请求登录: " + c.Request().Host + redirectURL)
	if tokenString != "" {
		ok, _, err := GetTokenClaims(c, tokenString)
		if ok && err == nil {
			if redirectURL == "" {
				redirectURL = "/"
			}
			return c.Redirect(http.StatusFound, redirectURL)
		}
	}

	// 渲染登录页面
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"SiteName":    config.GetSiteName(),
		"Copyright":   template.HTML(config.GetCopyright()),
		"RedirectURL": redirectURL,
	})
}

// LoginHandler 处理LDAP认证并生成Token
func LoginHandler(c echo.Context) error {
	var req RequestUserInfo
	redirectURL := c.QueryParam("redirect")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	log.Info("登录用户: " + req.Username + ", 访问地址: " + redirectURL)

	// LDAP认证
	ok, userInfo, err := modules.Ldap.Authenticate(req.Username, req.Password)
	if err != nil || !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication failed")
	}

	// 生成JWT Token
	token, claims, err := modules.Jwt.CreateToken(userInfo.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "token generation failed")
	}

	// 缓存token信息
	expiresAt := time.Until(claims.ExpiresAt.Time)
	modules.Cache.Set(token, claims, expiresAt)

	response := &ResponseJSON{
		Code: 200,
		Data: ResponseLoginInfo{
			Username:  claims.Username,
			Token:     token,
			ExpiresIn: claims.ExpiresAt.Time,
			Redirect:  redirectURL,
		},
	}

	// 设置认证Cookie
	modules.SetAuthCookie(c, token, claims.ExpiresAt.Time)

	return c.JSON(http.StatusOK, response)

	// // 如果是AJAX请求，返回JSON；否则重定向
	// if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
	// 	return c.JSON(http.StatusOK, response)
	// }
	// // 普通表单提交，重定向到原始页面或首页
	// if redirectURL == "" {
	// 	redirectURL = "/"
	// }
	// return c.Redirect(http.StatusFound, redirectURL)
}

// GetTokenClaims 验证token并返回信息
func GetTokenClaims(c echo.Context, token string) (bool, *modules.Claims, error) {
	claims := &modules.Claims{}
	claimsMap, err := modules.Cache.Get(token)
	if err != nil && err != modules.ErrKeyNotFound {
		return false, nil, err
	}
	if m, ok := claimsMap.(map[string]any); ok {
		if username, ok := m["username"].(string); ok {
			claims.Username = username
		}
		if iss, ok := m["iss"].(string); ok {
			claims.Issuer = iss
		}
		if jti, ok := m["jti"].(string); ok {
			claims.ID = jti
		}
		// 处理时间字段
		if exp, ok := m["exp"].(float64); ok {
			claims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(exp), 0))
		}
		if nbf, ok := m["nbf"].(float64); ok {
			claims.NotBefore = jwt.NewNumericDate(time.Unix(int64(nbf), 0))
		}
		if iat, ok := m["iat"].(float64); ok {
			claims.IssuedAt = jwt.NewNumericDate(time.Unix(int64(iat), 0))
		}
	} else {
		claims, err = modules.Jwt.VerifyToken(token)
		if err != nil {
			return false, nil, err
		}
	}

	if claims.Username == "" {
		return false, nil, nil
	}
	return true, claims, nil
}

// VerifyHandler 供nginx auth_request调用的验证接口
func VerifyHandler(c echo.Context) error {
	redirectURL := c.Request().Header.Get("X-Original-URI")
	tokenString := modules.Jwt.ExtractToken(c)
	// 检查白名单
	ok := modules.CheckUrlWhiteList(redirectURL)
	if ok {
		log.Debug("认证通过", "auth", "Whitelist", "uri", redirectURL)
		response := &ResponseJSON{
			Code: 200,
			// Msg:  "white list",
			// Data: ResponseLoginInfo{
			// 	Token:    tokenString,
			// 	Redirect: redirectURL,
			// },
		}
		return c.JSON(http.StatusOK, response)
	}

	// 检查认证
	if tokenString == "" {
		log.Warn("认证失败", "uri", redirectURL)
		return echo.NewHTTPError(http.StatusUnauthorized, "token required")
	}
	ok, claims, err := GetTokenClaims(c, tokenString)
	if !ok || err != nil {
		log.Warn("认证失败", "uri", redirectURL)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	// 认证通过
	log.Debug("认证通过", "auth", claims.Username, "uri", redirectURL)
	// 设置用户信息到header，供后端应用使用
	c.Response().Header().Set("X-Auth-User", claims.Username)
	response := &ResponseJSON{
		Code: 200,
		Msg:  "valid",
		Data: ResponseLoginInfo{
			Username:  claims.Username,
			Token:     tokenString,
			ExpiresIn: claims.ExpiresAt.Time,
			Redirect:  redirectURL,
		},
	}
	return c.JSON(http.StatusOK, response)
}
