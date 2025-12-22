package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"regexp"

	"golang.org/x/crypto/scrypt"
)

// MD5 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// MD5Byte 带byte的MD5
func MD5Byte(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

// NewPass 通过scrypt生成密码
func NewPass(passwd, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(dk), nil
}

// ComputeHash1 计算hash1
func ComputeHash1(message string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))
}

// ComputeHmacSha256 计算HmacSha256
func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))

}

// EncodeBase64 编码Base64
func EncodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Pkcs7Pad 使用 PKCS#7 填充
func Pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// EncryptAESECB 使用 AES ECB 模式加密
func EncryptAESECB(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddedPlaintext := Pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(paddedPlaintext))
	for bs, be := 0, block.BlockSize(); bs < len(paddedPlaintext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(ciphertext[bs:be], paddedPlaintext[bs:be])
	}
	return ciphertext, nil
}

// Pkcs7Unpad 移除 PKCS#7 填充
func Pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: data is empty")
	}
	// 最后一个字节是填充的长度
	padLength := int(data[length-1])
	// 填充长度不能大于总长度
	if padLength > length {
		return nil, errors.New("pkcs7: invalid padding size")
	}
	// 返回移除填充后的部分
	return data[:(length - padLength)], nil
}

// DecryptAESECB 使用 AES ECB 模式解密
func DecryptAESECB(ciphertext, key []byte) ([]byte, error) {
	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 检查密文长度是否是块大小的整数倍
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	decrypted := make([]byte, len(ciphertext))

	// ECB 模式是逐块解密
	for bs, be := 0, block.BlockSize(); bs < len(ciphertext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], ciphertext[bs:be])
	}

	return decrypted, nil
}

// ParseBase64 快速判断并解码Base64
func ParseBase64(s string) (bool, string) {
	if s == "" || len(s)%4 != 0 {
		return false, ""
	}

	// 检查字符集并确定类型
	// 标准Base64模式正则表达式：允许大小写字母、数字、+、/、=，长度为4的倍数
	// ^[A-Za-z0-9+/]*=?=?$
	stdPattern := `^[A-Za-z0-9+/]*={0,2}$`
	// URL安全Base64模式（包含-_）
	urlPattern := `^[A-Za-z0-9_-]*={0,2}$`

	stdMatch, _ := regexp.MatchString(stdPattern, s)
	urlMatch, _ := regexp.MatchString(urlPattern, s)
	if !stdMatch && !urlMatch {
		return false, ""
	}
	// 检查填充字符
	/*if strings.Contains(s, "=") {
		lastIndex := strings.LastIndex(s, "=")
		if lastIndex != len(s)-strings.Count(s, "=") {
			return false, ""
		}
		if strings.Count(s, "=") > 2 {
			return false, ""
		}
	}*/
	// 快速字符集检查
	/*for _, r := range s {
		if !((r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '+' || r == '/' || r == '=' ||
			r == '-' || r == '_') {
			return false, ""
		}
	}*/
	// 尝试标准Base64
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err == nil && len(decoded) > 0 {
		reencoded := base64.StdEncoding.EncodeToString(decoded)
		return reencoded == s, string(decoded)
	}
	// 尝试URL安全的Base64
	decoded, err = base64.URLEncoding.DecodeString(s)
	if err == nil && len(decoded) > 0 {
		reencoded := base64.URLEncoding.EncodeToString(decoded)
		return reencoded == s, string(decoded)
	}
	return false, ""
}
