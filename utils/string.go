package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

/**
驼峰转下划线
// 1. 普通使用
log.Println(CamelCase("AAAA"))
log.Println(CamelCase("IconUrl"))
log.Println(CamelCase("iconUrl"))
log.Println(CamelCase("parentId"))
log.Println(CamelCase("a9b9Ba"))
log.Println(CamelCase("_An"))
// s输出
//2019/03/20 16:34:25 a_a_a_a
//2019/03/20 16:34:25 icon_url
//2019/03/20 16:34:25 icon_url
//2019/03/20 16:34:25 parent_id
//2019/03/20 16:34:25 a9b9ba
//2019/03/20 16:34:25 Xan
*/
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		t = append(t, 'X')
		i++
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIIUpper(s[i+1]) {
			continue
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}

		if isASCIIUpper(c) {
			c ^= ' '
		}
		t = append(t, c)

		for i+1 < len(s) && isASCIIUpper(s[i+1]) {
			i++
			t = append(t, '_')
			t = append(t, bytes.ToLower([]byte{s[i]})[0])
		}
	}
	return string(t)
}
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

/**
// 转换为snake
*/
func ToSnakeCase(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	fmt.Println(snake)
	return strings.ToLower(snake)
}

// 转换为驼峰
func ToCamelCase(str string) string {
	temp := strings.Split(str, "-")
	for i, r := range temp {
		if i > 0 {
			temp[i] = strings.Title(r)
		}
	}
	return strings.Join(temp, "")
}

var re = regexp.MustCompile("(_|-)([a-zA-Z]+)")

// 转驼峰 优化版
func ToCamelCaseOptimization(str string) string {
	camel := re.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

// 驼峰式写法转为下划线写法
// https://github.com/polaris1119/goutils/blob/master/string.go
func UnderscoreName(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}
	return -1
}

// 蛇形字符串、驼峰字符串转换
// snake string, XxYy to xx_yy , XxYY to xx_yy
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 蛇形字符串、驼峰字符串转换
// camel string, xx_yy to XxYy
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

/**
判断字符串是否为空
*/
func IsStringEmpty(str string) bool {
	if str == "" || len(str) == 0 || strings.TrimSpace(str) == "" {
		return true
	}
	//isNil := reflect.ValueOf(str).IsNil()
	//if isNil {
	//	return true
	//}
	return false
}
