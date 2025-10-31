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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
)

func TestCache(t *testing.T) {
	key := "test"
	value := "123456"
	config.Init()
	CacheInit()

	// 测试代码
	err := Cache.Del(key)
	assert.Error(t, err, "期望值为: Error, 实际值为: %v", err)

	err = Cache.Set(key, value, DefaultExpirAt)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)

	time.Sleep(3 * time.Second)
	data, ttl, err := Cache.GetEx(key)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)
	assert.Equal(t, value, data, "期望值为: %s, 实际值为: %s", value, data)
	assert.Greater(t, DefaultExpirAt, ttl, "期望值为: %v, 实际值为: %v, 期望值应大于实际值", DefaultExpirAt, ttl)

	tempTime := 50 * time.Minute
	err = Cache.Expire(key, tempTime)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)
	time.Sleep(3 * time.Second)
	data, ttl, err = Cache.GetEx(key)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)
	assert.Equal(t, value, data, "期望值为: %s, 实际值为: %s", value, data)
	assert.Greater(t, tempTime, ttl, "期望值为: %v, 实际值为: %v, 期望值应大于实际值", tempTime, ttl)

	err = Cache.Del(key)
	assert.NoError(t, err, "期望值为: nil, 实际值为: %v", err)
}
