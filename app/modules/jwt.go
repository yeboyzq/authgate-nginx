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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
	"github.com/yeboyzq/authgate-nginx/app/utils"
)

var Jwt *JwtAuth

type JwtAuth struct {
	secret            string
	expiration        time.Duration
	signingMethodHMAC *jwt.SigningMethodHMAC
	storageName       string
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JwtInit 实例化jwt
func JwtInit() {
	Jwt = &JwtAuth{
		secret:            config.Cfg.GetString("base.jwt.secret"),
		expiration:        time.Duration(config.Cfg.GetInt("base.jwt.expiry")) * time.Hour,
		signingMethodHMAC: jwt.SigningMethodHS512,
		storageName:       config.Cfg.GetString("base.jwt.storageName"),
	}

	log.Info("认证组件初始化完成.")
}

// CreateToken 创建token
func (s *JwtAuth) CreateToken(username string) (string, *Claims, error) {
	expirationTime := time.Now().Add(s.expiration)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    utils.AppFileName(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        utils.NewDbUUID(),
		},
	}

	token := jwt.NewWithClaims(s.signingMethodHMAC, claims)
	signedString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", nil, err
	}
	return signedString, claims, nil
}

// VerifyToken 验证token
func (s *JwtAuth) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (s *JwtAuth) RefreshToken(tokenString string) (string, *Claims, error) {
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return "", nil, err
	}

	// 确保token在可刷新的时间范围内
	if time.Until(claims.ExpiresAt.Time) > s.expiration/2 {
		return tokenString, nil, nil // 不需要刷新
	}

	return s.CreateToken(claims.Username)
}

// ExtractToken 提取token
func (s *JwtAuth) ExtractToken(c echo.Context) string {
	// 从Authorization header提取
	authHeader := c.Request().Header.Get(s.storageName)
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 从cookie提取
	cookie, err := c.Cookie(s.storageName)
	if err == nil {
		return cookie.Value
	}

	// 从query参数提取
	return c.QueryParam(s.storageName)
}
