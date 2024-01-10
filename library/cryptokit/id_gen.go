package cryptokit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/v2/library/mathkit"
	"github.com/rs/xid"
	"strings"
)

func ID() string {
	id := xid.New()
	return id.String()
}

// MacRand 随机mac地址
func MacRand() string {
	arr := make([]string, 0, 6)
	for i := 0; i < 6; i++ {
		arr = append(arr, fmt.Sprintf("%02x", mathkit.RandInt32(0, 256)))
	}
	return strings.Join(arr, ":")
}
