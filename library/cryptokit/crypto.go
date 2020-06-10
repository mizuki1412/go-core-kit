package cryptokit

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(o string) string {
	m := md5.Sum([]byte(o))
	return hex.EncodeToString(m[:])
}

func URLEncode(s string) string {
	//println(base64.RawStdEncoding.EncodeToString([]byte(s))) 会去掉==
	return base64.URLEncoding.EncodeToString([]byte(s))
}
