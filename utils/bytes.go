package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// 包含辅助方法和常量，用于转换为人类可读的字节格式。
//
//	bytefmt.ByteSize(100.5*bytefmt.MEGABYTE) // "100.5M"
//	bytefmt.ByteSize(uint64(1024)) // "1K"
// https://github.com/cloudfoundry/bytefmt/blob/master/bytes.go

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

var invalidByteQuantityError = errors.New("字节数量必须是一个正整数，其单位为测量单位 M, MB, MiB, G, GiB, or GB")

// ByteSize 返回10M，12.5K等形式的人类可读字节串。以下单位可供选择：
//
//	E: Exabyte
//	P: Petabyte
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
//
// 始终选择导致最小数量大于或等于1的单位。
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "EB"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "PB"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "TB"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// ToMegabytes 将ByteSize格式化的字符串解析为兆字节。
func ToMegabytes(s string) (uint64, error) {
	bytes, err := ToBytes(s)
	if err != nil {
		return 0, err
	}

	return bytes / MEGABYTE, nil
}

// ToBytes 将ByteSize格式化的字符串解析为字节。注意二进制前缀和SI前缀单位均表示基数为2的单位
// KB = K = KiB	= 1024
// MB = M = MiB = 1024 * K
// GB = G = GiB = 1024 * M
// TB = T = TiB = 1024 * G
// PB = P = PiB = 1024 * T
// EB = E = EiB = 1024 * P
func ToBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, invalidByteQuantityError
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes <= 0 {
		return 0, invalidByteQuantityError
	}

	switch multiple {
	case "E", "EB", "EIB":
		return uint64(bytes * EXABYTE), nil
	case "P", "PB", "PIB":
		return uint64(bytes * PETABYTE), nil
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, invalidByteQuantityError
	}
}

// BytesToStringByBuffer io.ReadCloser类型转换为string
func BytesToStringByBuffer(body io.Reader) string {
	buf := new(bytes.Buffer) // io.ReadCloser类型转换为string
	buf.ReadFrom(body)
	//b := *(*string)(unsafe.Pointer(&body))
	/*_, err := io.Copy(buf, body)
	b :=buf.String()*/
	return buf.String()
}

// BytesToStringByIo io.ReadCloser类型转换为string
func BytesToStringByIo(body io.Reader) (string, error) {
	bd, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	b := string(bd)
	//b :=fmt.Sprintf("%s", body)
	return b, err
}

// Int2Byte 把int的每个字节取出来放入byte数组中，存储采用Littledian
func Int2Byte(data int) (ret []byte) {
	var len_ uintptr = unsafe.Sizeof(data)
	ret = make([]byte, len_)
	var tmp int = 0xff
	var index uint = 0
	for index = 0; index < uint(len_); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

// Byte2Int 把byte　Slice 中的每个字节取出来，　按Littledian端拼成一个int
func Byte2Int(data []byte) int {
	var ret int = 0
	var _len int = len(data)
	var i uint = 0
	for i = 0; i < uint(_len); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
}

// BytesToIntU 字节数(大端)组转成int(无符号的)
func BytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

// BytesToIntS 字节数(大端)组转成int(有符号)
func BytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

// IntToBytes 整形转换成字节
func IntToBytes(n int, b byte) ([]byte, error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		err := binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		err := binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		return bytesBuffer.Bytes(), nil
	case 3, 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		err := binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}
