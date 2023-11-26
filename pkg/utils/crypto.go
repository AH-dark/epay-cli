package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/AH-dark/bytestring"
)

func MD5(str string) []byte {
	m := md5.New()
	m.Write([]byte(str))
	dst := make([]byte, 32)
	hex.Encode(dst, m.Sum(nil))
	return dst
}

func MD5String(str string) string {
	return bytestring.BytesToString(MD5(str))
}
