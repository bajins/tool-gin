package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Get(org []int, i int, args ...int) (r int) {
	if i >= 0 && i < len(org) {
		r = org[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

// 将任何类型转换为字符串
func ToString(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', Get(args, 0, -1), Get(args, 1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', Get(args, 0, -1), Get(args, 1, 64))
	case int:
		s = strconv.FormatInt(int64(v), Get(args, 0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), Get(args, 0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), Get(args, 0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), Get(args, 0, 10))
	case int64:
		s = strconv.FormatInt(v, Get(args, 0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), Get(args, 0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), Get(args, 0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), Get(args, 0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), Get(args, 0, 10))
	case uint64:
		s = strconv.FormatUint(v, Get(args, 0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

/**
 * 驼峰转下划线
 * // 1. 普通使用
 * log.Println(CamelCase("AAAA"))
 * log.Println(CamelCase("IconUrl"))
 * log.Println(CamelCase("iconUrl"))
 * log.Println(CamelCase("parentId"))
 * log.Println(CamelCase("a9b9Ba"))
 * log.Println(CamelCase("_An"))
 * // s输出
 * //2019/03/20 16:34:25 a_a_a_a
 * //2019/03/20 16:34:25 icon_url
 * //2019/03/20 16:34:25 icon_url
 * //2019/03/20 16:34:25 parent_id
 * //2019/03/20 16:34:25 a9b9ba
 * //2019/03/20 16:34:25 Xan
 *
 * @param s string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:24
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

/**
 * 判断为ASCII编码大写
 *
 * @param c byte
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:26
 */
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

/**
 * 判断为ASCII编码数字
 *
 * @param c byte
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:27
 */
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

/**
 * 转换为snake
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:27
 */
func ToSnakeCase(str string) string {
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	fmt.Println(snake)
	return strings.ToLower(snake)
}

/**
 * 转换为驼峰
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:28
 */
func ToCamelCase(str string) string {
	temp := strings.Split(str, "-")
	for i, r := range temp {
		if i > 0 {
			temp[i] = strings.Title(r)
		}
	}
	return strings.Join(temp, "")
}

/**
 * 转换为驼峰，使用正则
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:30
 */
func ToCamelCaseRegexp(str string) string {
	var reg = regexp.MustCompile("(_|-)([a-zA-Z]+)")
	camel := reg.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

/**
 * 驼峰式写法转为下划线写法
 *
 * @param name string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:30
 */
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

/**
 * 下划线写法转为驼峰写法
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:31
 */
func CamelName(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

/**
 * 搜索字符串数组中是否存在指定字符串
 *
 * @param slice []string
 * @param s string
 * @return int 返回-1为未搜寻到
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:32
 */
func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}
	return -1
}

/**
 * 蛇形字符串
 * snake string, XxYy to xx_yy , XxYY to xx_yy
 *
 * @param s string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:33
 */
func SnakeString(s string) string {
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

/**
 * 驼峰字符串转换
 *
 * @param s string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:35
 */
func CamelString(s string) string {
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
 * 判断字符串是否为空
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:36
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

/**
 * 字符串截取
 *
 * @param str 字符串
 * @param pos 开始位置
 * @param length 结束位置
 * @return
 * @author claer woytu.com
 * @date 2019/6/29 3:27
 */
func Substring(str string, pos, length int) string {
	runes := []rune(str)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

/**
 * 首字母转大写
 *
 * @param null
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:19
 */
func ToUpper(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

/**
 * 首字母转小写
 *
 * @param null
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:19
 */
func ToLower(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

/**
 * Unicode转汉字
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/17 11:44
 */
func UnicodeToChinese(str string) string {
	buf := bytes.NewBuffer(nil)

	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			// 将字符串转换为uint类型整数
			// base：进位制（2 进制到 36 进制）
			// bitSize：指定整数类型（0:int、8:int8、16:int16、32:int32、64:int64）
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}
