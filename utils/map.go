package utils

import (
	"encoding/json"
	"reflect"
)

func StructToMapByReflect(s interface{}) map[string]interface{} {
	elem := reflect.ValueOf(&s).Elem()
	type_ := elem.Type()

	map_ := map[string]interface{}{}

	for i := 0; i < type_.NumField(); i++ {
		map_[type_.Field(i).Name] = elem.Field(i).Interface()
	}
	return map_
}

func StructToMapByJson(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	// 对象转换为JSON
	j, _ := json.Marshal(&s)
	// JSON 转换回对象
	json.Unmarshal(j, &m)
	return m
}
