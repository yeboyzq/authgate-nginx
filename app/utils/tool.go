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

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dustin/go-humanize"
	"go.yaml.in/yaml/v4"
)

// FileSize 计算文件大小并生成用户友好的字符串
func FileSize(s int64) string {
	return humanize.IBytes(uint64(s))
}

// FindFirstNonEmptyValue map返回数组中第一个非空值
func FindFirstNonEmptyValue(keys []string, data map[string]string) (string, bool) {
	for _, v := range keys {
		if value, exists := data[v]; exists && value != "" {
			return value, true
		}
	}
	return "", false
}

// JsonToString 将任意数据转换为JSON字符串
func JsonToString(v any) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// JsonToYaml 将JSON数据转换为YAML
func JsonToYaml(j []byte) ([]byte, error) {
	var jsonObj interface{}

	err := yaml.Unmarshal(j, &jsonObj)
	if err != nil {
		return nil, err
	}

	yamlBytes, err := yaml.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}

// StructElement 结构体元素
type StructElement struct {
	FieldName string
	FieldType string
	Anonymous bool
}

// AnalyzeStruct 分析结构体获取结构体元素数量及元素
func AnalyzeStruct(name string, s any) (num int, elements []StructElement, err error) {
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("%s不是结构体, 无法分析", name)
		return num, elements, err
	}

	numFields := t.NumField()

	elements = make([]StructElement, numFields)
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldType := field.Type.String()
		// 简化类型名称显示
		fieldType = strings.ReplaceAll(fieldType, "main.", "")
		elements[i] = StructElement{
			FieldName: field.Name,
			FieldType: fieldType,
			Anonymous: field.Anonymous,
		}
	}
	return numFields, elements, nil
}
