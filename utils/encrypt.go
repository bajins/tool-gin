package utils

import (
	"crypto/md5"
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
