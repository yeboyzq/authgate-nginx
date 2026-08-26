// @作者: yeboyzq
// @更新: 2026-06-30
// @备注:

package id

import (
	"fmt"

	"github.com/google/uuid"
)

// 常用uuid变量
const (
	DefaultUUIDv4 = "00000000-0000-4000-9999-000000000000"
	DefaultUUIDv7 = "00000000-0000-7000-9999-000000000000"
	NilUUID       = "00000000-0000-0000-0000-000000000000"
	MaxUUID       = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	MinUUID       = NilUUID
)

// NewDbUUID 生成数据库uuid
func NewDbUUID() string {
	NewUuid, err := uuid.NewV7()
	if err != nil {
		panic("生成数据库uuid失败: " + err.Error())
	}
	return NewUuid.String()
}

// NewRequestID 生成请求uuid
func NewRequestID() string {
	NewUuid, err := uuid.NewV7()
	if err != nil {
		panic("生成请求uuid失败: " + err.Error())
	}
	return NewUuid.String()
}

// IsUUID 判断UUID是否有效
func IsUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func IsUUIDv7(s string) bool {
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == 7
}

// ParseUUID 验证并返回UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// ParseUUIDv7 解析字符串并校验是否为 UUID 版本 7。
func ParseUUIDv7(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}
	if ver := u.Version(); ver != 7 {
		return uuid.Nil, fmt.Errorf("无效的UUID版本: 预期为7, 实际为%d", ver)
	}
	return u, nil
}
