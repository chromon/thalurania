package scrypt

import (
	"chalurania/service/log"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

// 密码加密
func Crypto(password string) string {
	// 密码盐，重要，随机编的
	salt := []byte{0xe8, 0x9d, 0xb2, 0x77, 0xc3, 0xfe, 0xa5, 0xde}

	// N is a CPU/memory cost parameter, which must be a power of two greater than 1.
	// r and p must satisfy r * p < 2 ^ 30.
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 5, 3, 32)
	if err != nil {
		log.Error.Println("scrypt password err:", err)
	}
	// base64, 设字符串长度为 n, 长度为 ⌈n/3⌉*4 向上取整 = 44
	return base64.StdEncoding.EncodeToString(dk)
}