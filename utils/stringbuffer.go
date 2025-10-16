package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// StringBuffer 是一个线程安全的、可变的字符串缓冲区。
// 它通过在 strings.Builder 周围封装一个互斥锁来实现。
type StringBuffer struct {
	builder strings.Builder
	mu      sync.Mutex
}

// NewStringBuffer 创建并返回一个新的 StringBuffer 实例。
func NewStringBuffer() *StringBuffer {
	return &StringBuffer{}
}

// Append 将任意类型的值转换为字符串并追加到缓冲区。
// 它使用 type switch 为常见类型提供高性能转换。
// 这是一个线程安全的操作，并支持链式调用。
func (sb *StringBuffer) Append(val interface{}) *StringBuffer {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	// 使用 Type Switch 为特定类型提供高效的字符串转换
	switch v := val.(type) {
	case string:
		sb.builder.WriteString(v)
	case int:
		sb.builder.WriteString(strconv.Itoa(v))
	case int64:
		sb.builder.WriteString(strconv.FormatInt(v, 10))
	case uint:
		sb.builder.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint64:
		sb.builder.WriteString(strconv.FormatUint(v, 10))
	case float32:
		sb.builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		sb.builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		sb.builder.WriteString(strconv.FormatBool(v))
	case []byte:
		sb.builder.Write(v)
	case rune:
		sb.builder.WriteRune(v)
	default:
		// 对于其他所有类型，回退到 fmt.Sprintf
		sb.builder.WriteString(fmt.Sprintf("%v", v))
	}

	return sb
}

// Len 返回缓冲区中当前存储的字节数。
// 这是一个线程安全的操作。
func (sb *StringBuffer) Len() int {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.builder.Len()
}

// String 返回缓冲区内容的字符串表示。
// 这是一个线程安全的操作。
func (sb *StringBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.builder.String()
}

// Reset 清空缓冲区，使其可以被重用。
// 这是一个线程安全的操作。
func (sb *StringBuffer) Reset() {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.builder.Reset()
}
