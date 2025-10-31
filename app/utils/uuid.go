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

import "github.com/google/uuid"

// 常用uuid变量
const (
	DefaultUUIDv4 = "00000000-0000-4000-9999-000000000000"
	DefaultUUIDv7 = "00000000-0000-7000-9999-000000000000"
	NilUUID       = "00000000-0000-0000-0000-000000000000"
	MaxUUID       = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	MinUUID       = NilUUID
)

// NewDbUUID 生成数据库uuid
func NewDbUUID() (id string) {
	NewUuid, err := uuid.NewV7()
	if err != nil {
		panic("生成数据库uuid失败: " + err.Error())
	}
	return NewUuid.String()
}

// NewRequestID 生成请求uuid
func NewRequestID() (id string) {
	NewUuid, err := uuid.NewV7()
	if err != nil {
		panic("生成请求uuid失败: " + err.Error())
	}
	return NewUuid.String()
}

// IsValidID 判断ID是否有效
func IsValidID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ParseID 验证并返回UUID
func ParseID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
