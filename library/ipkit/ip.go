package ipkit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/bytekit"
	"github.com/mizuki1412/go-core-kit/v2/library/regexkit"
	"github.com/spf13/cast"
	"strings"
)

// ComputeNetworkBytes 计算网段
func ComputeNetworkBytes(ip, mask string) (uint32, int) {
	if !regexkit.IsIP(ip) {
		panic(exception.New("ip不合法"))
	}
	if !regexkit.IsIP(mask) {
		panic(exception.New("mask不合法"))
	}
	ipInt := ipString2Int(ip)
	maskInt := ipString2Int(mask)
	num := 0
	for i := 31; i >= 0; i-- {
		if maskInt>>i&0x01 == 1 {
			num++
		} else {
			break
		}
	}
	return ipInt & ((maskInt >> (32 - num)) << (32 - num)), num
}

func ComputeNetworkString(ip, mask string) (string, int) {
	netInt, num := ComputeNetworkBytes(ip, mask)
	return ipInt2String(netInt), num
}

// ComputeSubNetString 计算指定的子网范围
func ComputeSubNetString(ip, mask string, startIndex, endIndex int) []string {
	netStart, maskNum := ComputeNetworkBytes(ip, mask)
	netEnd := netStart | 0xffffffff>>maskNum
	if endIndex-startIndex > 0 && cast.ToUint32(endIndex-startIndex) > netEnd-netStart {
		panic(exception.New("子网范围超出"))
	}
	return []string{ipInt2String(netStart + cast.ToUint32(startIndex)), ipInt2String(netStart + cast.ToUint32(endIndex))}
}

func ipString2Int(ip string) uint32 {
	ss := strings.Split(ip, ".")
	bytes := make([]byte, 0, len(ss))
	for _, e := range ss {
		bytes = append(bytes, cast.ToUint8(e))
	}
	return bytekit.Bytes2Uint32(bytes)
}

func ipInt2String(val uint32) string {
	ss := make([]string, 0, 4)
	for _, e := range bytekit.Int32ToBytes(val) {
		ss = append(ss, cast.ToString(e))
	}
	return strings.Join(ss, ".")
}
