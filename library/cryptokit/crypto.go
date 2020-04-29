package cryptokit

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(o string) string {
	m := md5.Sum([]byte(o))
	return hex.EncodeToString(m[:])
}
