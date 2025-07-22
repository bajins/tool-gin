package utils

import (
	"fmt"
	"github.com/axgle/mahonia"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// ParseBool 返回字符串表示的布尔值。
// 它接受1,1.0，t，T，TRUE，true，True，YES，yes，Yes，Y，y，ON，on，On，0,0.0，f，F，FALSE，false，False，NO，no，No，N，n，OFF，off，Off。
// 任何其他值都会返回错误。
func ParseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%s", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1 {
				return true, nil
			} else if v == 0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

// Type 返回参数的类型
func Type(v interface{}) string {
	t := reflect.TypeOf(v)
	k := t.Kind()
	return k.String()
}

// InArray 判断是否在数组中
func InArray(in interface{}, list interface{}) bool {
	ret := false
	if in == nil {
		in = ""
	}

	// 判断list是否slice
	l := reflect.TypeOf(list).String()
	t := Type(in)
	if "[]"+t != l {
		return false
	}

	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		for _, l := range list.([]string) {
			v := reflect.ValueOf(l)
			if tv == v.String() {
				ret = true
				break
			}
		}

	case "int":
		tv := reflect.ValueOf(in).Int()
		for _, l := range list.([]int) {
			v := reflect.ValueOf(l)
			if tv == v.Int() {
				ret = true
				break
			}
		}
	}

	return ret
}

// GBK2UTF gbk转换utf-8
func GBK2UTF(text string) string {
	enc := mahonia.NewDecoder("GB18030")

	text = enc.ConvertString(text)

	return strings.ReplaceAll(text, "聽", "&nbsp;")
}

// Pagination 分页
// page 页数
// rows 取多少条数据
// total 数据总条数
// 返回 起始-结束
func Pagination(page, rows, total int) (int, int) {
	offset := (page - 1) * rows
	limit := offset + rows
	if limit > total {
		return offset, total
	}
	return offset, limit
}

// LogWithLocation 是一个封装好的日志工具函数
func LogWithLocation(message string) {
	// 调用栈:
	// 0: runtime.Caller
	// 1: LogWithLocation (当前函数)
	// 2: main (我们想知道的位置!)
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		// 处理错误
		log.Printf("无法获取调用信息: %s", message)
		return
	}

	funcName := runtime.FuncForPC(pc).Name()

	// 格式化输出
	log.Printf("%s:%d [%s] - %s", file, line, funcName, message)
}
