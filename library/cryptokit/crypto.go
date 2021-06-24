package cryptokit

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func MD5(o string) string {
	m := md5.Sum([]byte(o))
	return hex.EncodeToString(m[:])
}

// URLEncode 这里是base64的，中文url用url.QueryEscape
func URLEncode(s string) string {
	//println(base64.RawStdEncoding.EncodeToString([]byte(s))) 会去掉==
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func BytesEncode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func BytesDecode(str string) []byte {
	bytes, _ := base64.StdEncoding.DecodeString(str)
	return bytes
}

func HmacSha256(message []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(message)
	//sha := hex.EncodeToString()
	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
