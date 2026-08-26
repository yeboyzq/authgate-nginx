// @作者: yeboyzq
// @更新: 2026-06-30
// @备注:

package id

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewDbUUID 测试生成数据库 UUID。
func TestNewDbUUID(t *testing.T) {
	value := NewDbUUID()
	require.NotEmpty(t, value)
	parsed, err := uuid.Parse(value)
	require.NoError(t, err)
	assert.Equal(t, uuid.Version(7), parsed.Version())
}

// TestNewRequestID 测试生成请求 UUID。
func TestNewRequestID(t *testing.T) {
	value := NewRequestID()
	require.NotEmpty(t, value)
	parsed, err := uuid.Parse(value)
	require.NoError(t, err)
	assert.Equal(t, uuid.Version(7), parsed.Version())
}

// TestIsUUID 测试 UUID 合法性判断。
func TestIsUUID(t *testing.T) {
	assert.True(t, IsUUID(NewDbUUID()))
	assert.True(t, IsUUID(DefaultUUIDv7))
	assert.False(t, IsUUID(""))
	assert.False(t, IsUUID("not-a-uuid"))
	assert.False(t, IsUUID("00000000-0000-7000-9999-00000000000"))
}

// TestIsUUID 测试 UUID 合法性判断。
func TestIsUUIDv7(t *testing.T) {
	assert.True(t, IsUUIDv7(NewDbUUID()))
	assert.True(t, IsUUIDv7(DefaultUUIDv7))
	assert.False(t, IsUUIDv7(""))
	assert.False(t, IsUUIDv7("not-a-uuid"))
	assert.False(t, IsUUIDv7("00000000-0000-7000-9999-00000000000"))
}

// TestParseUUID 测试解析 UUID。
func TestParseUUID(t *testing.T) {
	generated := NewRequestID()
	parsed, err := ParseUUID(generated)
	require.NoError(t, err)
	assert.Equal(t, generated, parsed.String())

	parsed, err = ParseUUID(DefaultUUIDv7)
	require.NoError(t, err)
	assert.Equal(t, DefaultUUIDv7, parsed.String())

	parsed, err = ParseUUID("bad-id")
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsed)
}

func TestParseUUIDv7(t *testing.T) {
	generated := NewRequestID()
	parsed, err := ParseUUIDv7(generated)
	require.NoError(t, err)
	assert.Equal(t, generated, parsed.String())

	parsed, err = ParseUUIDv7(DefaultUUIDv7)
	require.NoError(t, err)
	assert.Equal(t, DefaultUUIDv7, parsed.String())

	parsed, err = ParseUUIDv7("bad-id")
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsed)
}
