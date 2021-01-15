package modbus

import "sync"

var (
	once     sync.Once
	crtTable []uint16
)

// CRC16 Calculate Cyclical Redundancy Checking
// 使用CRC码需要高低字节对换，byte+crc16 完整字节数组做CRC16，结果应该是0
func CRC16(bs []byte) uint16 {
	once.Do(initCrcTable)

	val := uint16(0xFFFF)
	for _, v := range bs {
		val = (val >> 8) ^ crtTable[(val^uint16(v))&0x00FF]
	}
	return val
}

// crc16 计算完后转成字节数组
func CRC16Bytes(crc16 uint16) []byte {
	// 高低字节对调
	return []byte{byte(crc16 & 0xffff), byte(crc16 >> 8 & 0xffff)}
}

func CheckCRC16(all []byte) bool {
	return CRC16(all) == 0
}

func initCrcTable() {
	crcPoly16 := uint16(0xa001)
	crtTable = make([]uint16, 256)
	for i := uint16(0); i < 256; i++ {
		crc := uint16(0)
		b := i
		for j := uint16(0); j < 8; j++ {
			if ((crc ^ b) & 0x0001) > 0 {
				crc = (crc >> 1) ^ crcPoly16
			} else {
				crc = crc >> 1
			}
			b = b >> 1
		}
		crtTable[i] = crc
	}
}
