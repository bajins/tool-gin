package utils

import (
	"encoding/json"
	"reflect"
)

// StructToMapByReflect 结构体转map
func StructToMapByReflect(s any) map[string]any {
	elem := reflect.ValueOf(&s).Elem()
	type_ := elem.Type()

	map_ := map[string]any{}

	for i := 0; i < type_.NumField(); i++ {
		map_[type_.Field(i).Name] = elem.Field(i).Interface()
	}
	return map_
}

// StructToMapByJson 结构体转map
func StructToMapByJson(s any) map[string]any {
	m := make(map[string]any)
	// 对象转换为JSON
	j, _ := json.Marshal(&s)
	// JSON 转换回对象
	json.Unmarshal(j, &m)
	return m
}
