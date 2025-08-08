package utils

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToString 将任何类型转换为字符串
func ToString(value interface{}) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(v, 'f', -1, 32)
	case int:
		s = strconv.FormatInt(int64(v), 10)
	case int8:
		s = strconv.FormatInt(int64(v), 10)
	case int16:
		s = strconv.FormatInt(int64(v), 10)
	case int32:
		s = strconv.FormatInt(int64(v), 10)
	case int64:
		s = strconv.FormatInt(v, 10)
	case uint:
		s = strconv.FormatUint(uint64(v), 10)
	case uint8:
		s = strconv.FormatUint(uint64(v), 10)
	case uint16:
		s = strconv.FormatUint(uint64(v), 10)
	case uint32:
		s = strconv.FormatUint(uint64(v), 10)
	case uint64:
		s = strconv.FormatUint(v, 10)
	case string:
		s = v
	case []byte:
		s = string(v)
	case time.Time:
		s = TimeToString(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

// CamelCase 驼峰转下划线
// 1. 普通使用
// log.Println(CamelCase("AAAA"))
// log.Println(CamelCase("IconUrl"))
// log.Println(CamelCase("iconUrl"))
// log.Println(CamelCase("parentId"))
// log.Println(CamelCase("a9b9Ba"))
// log.Println(CamelCase("_An"))
// s输出
// 2019/03/20 16:34:25 a_a_a_a
// 2019/03/20 16:34:25 icon_url
// 2019/03/20 16:34:25 icon_url
// 2019/03/20 16:34:25 parent_id
// 2019/03/20 16:34:25 a9b9ba
// 2019/03/20 16:34:25 Xan
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

// 判断为ASCII编码大写
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// 判断为ASCII编码数字
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// ToSnakeCase 转换为snake
func ToSnakeCase(str string) string {
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	fmt.Println(snake)
	return strings.ToLower(snake)
}

// ToCamelCase 转换为驼峰
func ToCamelCase(str string) string {
	temp := strings.Split(str, "-")
	for i, r := range temp {
		if i > 0 {
			temp[i] = cases.Title(language.English).String(r)
		}
	}
	return strings.Join(temp, "")
}

// ToCamelCaseRegexp 转换为驼峰，使用正则
func ToCamelCaseRegexp(str string) string {
	var reg = regexp.MustCompile("([_\\-])([a-zA-Z]+)")
	camel := reg.ReplaceAllString(str, " $2")
	camel = cases.Title(language.English).String(camel)
	camel = strings.ReplaceAll(camel, " ", "")
	return camel
}

// UnderscoreName 驼峰式写法转为下划线写法
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

// CamelName 下划线写法转为驼峰写法
func CamelName(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = cases.Title(language.English).String(str)
	return strings.ReplaceAll(str, " ", "")
}

// SearchString 搜索字符串数组中是否存在指定字符串
// 返回-1为未搜寻到
func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}
	return -1
}

// SnakeString 蛇形字符串
// snake string, XxYy to xx_yy , XxYY to xx_yy
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

// CamelString 驼峰字符串转换
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

// IsStringEmpty 判断字符串是否为空
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

// Substring 字符串截取
func Substring(str string, pos, length int) string {
	runes := []rune(str)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// ToUpper 首字母转大写
func ToUpper(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// ToLower 首字母转小写
func ToLower(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// UnicodeToChinese Unicode转汉字
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

// RandString 生成指定长度大写字母随机字符串
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	by := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		by[i] = byte(b)
	}
	return string(by)
}

// RandomString 生成指定长度数字、小写字母、大写字母随机字符串
// 随机字符串生成库 https://github.com/lifei6671/gorand
func RandomString(len int) (s string, err error) {
	var container string
	var str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := bytes.NewBufferString(str).Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, err := crand.Int(crand.Reader, bigInt)
		if err != nil {
			return "", err
		}
		container += string(str[randomInt.Int64()])
	}
	return container, nil
}

// RandomLowercaseAlphanumeric 生成指定长度数字、小写字母随机字符串
func RandomLowercaseAlphanumeric(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	by := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, by[r.Intn(len(by))])
	}
	return string(result)
}

// RandomNumber 指定长度随机数字符串
func RandomNumber(len int) (s int, err error) {
	numbers := []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()
	bigInt := big.NewInt(int64(length))
	for i := 1; i <= len; i++ {
		random, err := crand.Int(crand.Reader, bigInt)
		if err != nil {
			return 0, err
		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	// 字符串转数字
	number, err := strconv.Atoi(container)
	if err != nil {
		return 0, err
	}
	return number, nil
}

// RandomCustomizeNumber 指定长度随机自定义数字符串
func RandomCustomizeNumber(len int, numbers []byte) (s int, err error) {
	var container string
	length := bytes.NewReader(numbers).Len()
	bigInt := big.NewInt(int64(length))
	for i := 1; i <= len; i++ {
		random, err := crand.Int(crand.Reader, bigInt)
		if err != nil {
			return 0, err
		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	// 字符串转数字
	number, err := strconv.Atoi(container)
	if err != nil {
		return 0, err
	}
	return number, nil
}

// JsonToMap 解析json为map
func JsonToMap(data string) (map[string]interface{}, error) {
	str := []byte(data)
	stu := make(map[string]interface{})

	err := json.Unmarshal(str, &stu)
	if err != nil {
		return nil, err
	}
	return stu, nil
}

// ParseJsonReader 反序列化为map
func ParseJsonReader(input io.Reader) (map[string]interface{}, error) {
	var m map[string]interface{}
	decoder := json.NewDecoder(input)
	err := decoder.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// RandomLower 纯小写
func RandomLower(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandomMixed 大小写混合
func RandomMixed(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RemoveDuplicateLines 删除重复行
func RemoveDuplicateLines(text string) string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		if line != "" && !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// RemoveDuplicateLinesKeepEmpty 删除重复行并保留空行
func RemoveDuplicateLinesKeepEmpty(text string) string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// RemoveDuplicateLinesOrdered 删除重复行并使用结构体保持插入顺序（更严格的顺序保证）
func RemoveDuplicateLinesOrdered(text string) string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]struct{})
	var result []string

	for _, line := range lines {
		if line != "" {
			if _, exists := seen[line]; !exists {
				seen[line] = struct{}{}
				result = append(result, line)
			}
		} else {
			// 保留空行
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// RemoveDuplicateLinesUniversal 删除重复行并使用通用方法，处理 Windows 和 Unix 换行符
func RemoveDuplicateLinesUniversal(text string) string {
	// 使用正则表达式分割，支持不同的换行符
	lines := regexp.MustCompile(`\r?\n`).Split(text, -1)
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		if line != "" && !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	// 使用 \n 作为连接符
	return strings.Join(result, "\n")
}

// RemoveDuplicateLinesIgnoreCase 删除重复行并忽略大小写
func RemoveDuplicateLinesIgnoreCase(text string) string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		if line != "" && !seen[lowerLine] {
			seen[lowerLine] = true
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// ExtractContent 提取正则()内的内容（不包括()前后内容）
func ExtractContent(rex *regexp.Regexp, text string) []string {
	matches := rex.FindAllStringSubmatch(text, -1)
	var contents []string

	for _, match := range matches {
		if len(match) > 1 {
			contents = append(contents, match[1])
		}
	}
	return contents
}
