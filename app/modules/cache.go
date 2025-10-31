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
	"encoding/json"
	"errors"
	"time"

	"github.com/coocood/freecache"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
	"github.com/yeboyzq/authgate-nginx/app/modules/log"
)

const DefaultExpirAt time.Duration = time.Duration(8 * time.Hour)

var (
	Cache *fcacheStore
)

// fcacheStore 定义缓存客户端结构体
type fcacheStore struct {
	client *freecache.Cache
}

// CacheInit 实例化freecache
func CacheInit() {
	cacheSize := config.Cfg.GetInt("base.cache.maxsize")
	fcache := freecache.NewCache(cacheSize * 1024)
	Cache = &fcacheStore{
		client: fcache,
	}
	log.Info("缓存组件初始化完成.")
}

// 错误定义
var (
	ErrKeyEncodeFail            = errors.New("key encode fail")
	ErrKeyDecoderFail           = errors.New("key decoder fail")
	ErrKeyNotFound              = errors.New("key not found")
	ErrKeyCacheTypeNotSupported = errors.New("cache type not supported")
)

// Encoder 编码, 将interface转byte
func Encoder(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, ErrKeyEncodeFail
	}
	return data, nil
}

// Decoder 解码, 将byte转interface
func Decoder(binary []byte) (any, error) {
	var data any
	err := json.Unmarshal(binary, &data)
	if err != nil {
		return nil, ErrKeyDecoderFail
	}
	return data, nil
}

// ************************************************************
// 以下为方法实现
// Get 获取缓存值
func (s *fcacheStore) Get(key string) (any, error) {
	v, err := s.client.Get([]byte(key))
	// 错误处理
	switch {
	case errors.Is(err, freecache.ErrNotFound):
		return nil, ErrKeyNotFound
	case err != nil:
		return nil, err
	}

	resVal, err := Decoder(v)
	return resVal, err
}

// GetEx 获取缓存值, 过期时间
func (s *fcacheStore) GetEx(key string) (any, time.Duration, error) {
	v, t, err := s.client.GetWithExpiration([]byte(key))
	// 错误处理
	switch {
	case errors.Is(err, freecache.ErrNotFound):
		return nil, 0 * time.Second, ErrKeyNotFound
	case err != nil:
		return nil, 0 * time.Second, err
	}

	resVal, err := Decoder(v)
	currentTime := time.Now().Unix() // current time in Unix timestamp (seconds)
	resTtl := time.Duration(int64(t)-currentTime) * time.Second
	return resVal, resTtl, err
}

// TTL 获取缓存过期时间
func (s *fcacheStore) TTL(key string) (time.Duration, error) {
	ttl, err := s.client.TTL([]byte(key))
	// 错误处理
	switch {
	case errors.Is(err, freecache.ErrNotFound):
		return 0 * time.Second, ErrKeyNotFound
	case err != nil:
		return 0 * time.Second, err
	}

	resTtl := time.Duration(ttl) * time.Second
	return resTtl, err
}

// Set 设置值
func (s *fcacheStore) Set(key string, value any, lifetime time.Duration) error {
	reqKey := []byte(key)
	reqVal, err := Encoder(value)
	if err != nil {
		return err
	}
	lifetimeq := int(lifetime.Seconds())
	err = s.client.Set(reqKey, reqVal, lifetimeq)
	return err
}

// Expire 更新过期时间
func (s *fcacheStore) Expire(key string, lifetime time.Duration) error {
	reqKey := []byte(key)
	lifetimeq := int(lifetime.Seconds())
	err := s.client.Touch(reqKey, lifetimeq)
	if err != nil {
		return err
	}
	return nil
}

// Del 删除值
func (s *fcacheStore) Del(key string) error {
	_, err := s.TTL(key)
	switch {
	case errors.Is(err, freecache.ErrNotFound):
		return nil
	case err != nil:
		return err
	default:
		v := s.client.Del([]byte(key))
		if !v {
			return err
		}
		return nil
	}
}

// BatchDel 批量删除
func (s *fcacheStore) BatchDel(key string, count int64) error {
	s.client.Clear()
	return nil
}

// EntryCount 返回当前缓存中的项目数
func (s *fcacheStore) EntryCount() int64 {
	return s.client.EntryCount()
}
