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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
)

func TestJwt(t *testing.T) {
	config.Init()
	JwtInit()
	// 设置认证信息
	userID := "test"
	token, claims, err := Jwt.CreateToken(userID)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)

	// 获取token信息
	resultclaims, err := Jwt.VerifyToken(token)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)
	// 判断认证信息是否一致
	assert.Equal(t, claims, resultclaims, "期望值为: %v, 实际值为: %v", claims, resultclaims)
	assert.Equal(t, userID, resultclaims.Username, "期望值为: %v, 实际值为: %v", userID, resultclaims.Username)
}
