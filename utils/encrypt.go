package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func MD5Byte(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

// 通过scrypt生成密码
func NewPass(passwd, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(dk), nil
}

// 计算hash1
func ComputeHash1(message string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))
}

// 计算HmacSha256
func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	// 转成十六进制
	return hex.EncodeToString(h.Sum(nil))

}

// 编码Base64
func EncodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
