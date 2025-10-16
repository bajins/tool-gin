package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// StringJoiner 用于构建由分隔符分隔的字符串，并可选择添加前缀和后缀。
type StringJoiner struct {
	builder    strings.Builder
	delimiter  string
	prefix     string
	suffix     string
	emptyValue string
	isFirst    bool
}

// NewStringJoiner 创建一个新的 StringJoiner。
func NewStringJoiner(delimiter string) *StringJoiner {
	return &StringJoiner{
		delimiter: delimiter,
		isFirst:   true,
	}
}

// SetPrefix 设置前缀
func (sj *StringJoiner) SetPrefix(prefix string) *StringJoiner {
	sj.prefix = prefix
	return sj
}

// SetSuffix 设置后缀
func (sj *StringJoiner) SetSuffix(suffix string) *StringJoiner {
	sj.suffix = suffix
	return sj
}

// SetEmptyValue 设置当没有添加任何元素时的默认返回值。
func (sj *StringJoiner) SetEmptyValue(emptyValue string) *StringJoiner {
	sj.emptyValue = emptyValue
	return sj
}

// Add 添加一个新的元素到 StringJoiner。
// 它使用 type switch 为常见类型提供高性能转换。
func (sj *StringJoiner) Add(val interface{}) *StringJoiner {
	if sj.isFirst {
		sj.builder.WriteString(sj.prefix)
		sj.isFirst = false
	} else {
		sj.builder.WriteString(sj.delimiter)
	}

	// 使用 Type Switch 为特定类型提供高效的字符串转换
	switch v := val.(type) {
	case string:
		sj.builder.WriteString(v)
	case int:
		sj.builder.WriteString(strconv.Itoa(v))
	case int64:
		sj.builder.WriteString(strconv.FormatInt(v, 10))
	case uint:
		sj.builder.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint64:
		sj.builder.WriteString(strconv.FormatUint(v, 10))
	case float32:
		sj.builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		sj.builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		sj.builder.WriteString(strconv.FormatBool(v))
	case []byte:
		sj.builder.Write(v)
	case rune:
		sj.builder.WriteRune(v)
	default:
		// 对于其他所有类型
		sj.builder.WriteString(fmt.Sprintf("%v", v))
	}

	return sj
}

// AddInt 添加一个新的整数元素到 StringJoiner。
func (sj *StringJoiner) AddInt(value int) *StringJoiner {
	return sj.Add(fmt.Sprintf("%d", value))
}

// AddFloat 添加一个新的浮点数元素到 StringJoiner。
func (sj *StringJoiner) AddFloat(value float64) *StringJoiner {
	return sj.Add(fmt.Sprintf("%f", value))
}

// AddBool 添加一个新的布尔元素到 StringJoiner。
func (sj *StringJoiner) AddBool(value bool) *StringJoiner {
	return sj.Add(fmt.Sprintf("%t", value))
}

// Merge 合并另一个 StringJoiner
func (sj *StringJoiner) Merge(other *StringJoiner) *StringJoiner {
	if other.builder.Len() > 0 {
		if sj.builder.Len() > 0 {
			sj.builder.WriteString(sj.delimiter)
		}
		sj.builder.WriteString(other.builder.String())
	}
	return sj
}

// Length 返回当前内容的长度
func (sj *StringJoiner) Length() int {
	if sj.builder.Len() <= 0 {
		return len(sj.emptyValue)
	}
	return len(sj.prefix) + sj.builder.Len() + len(sj.suffix)
}

// Empty 检查是否为空
func (sj *StringJoiner) Empty() bool {
	return sj.builder.Len() <= 0
}

// String 返回最终的字符串。
func (sj *StringJoiner) String() string {
	if sj.builder.Len() == 0 && sj.emptyValue != "" {
		return sj.emptyValue
	}
	if sj.builder.Len() == 0 {
		return sj.prefix + sj.suffix
	}
	sj.builder.WriteString(sj.suffix)
	return sj.builder.String()
}
