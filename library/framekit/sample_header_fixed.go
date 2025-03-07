package framekit

// NewHeaderFixedDecoder sample：带header字节和固定长度的帧
func NewHeaderFixedDecoder(capacity int, header []byte, fixedLength int) *Decoder {
	return NewDecoder(capacity, func(bytes []byte) ([]byte, []byte, bool) {
		if len(bytes) < fixedLength {
			return bytes, nil, false
		}
		for i := 0; i < fixedLength-len(header); i++ {
			count := 0
			for j := 0; j < len(header); j++ {
				if bytes[i+j] == header[j] {
					count++
				} else {
					break
				}
			}
			if count == len(header) {
				// 找到帧头 todo 没处理最长
				res := bytes[i : i+fixedLength]
				return bytes[i+fixedLength:], res, false
			}
		}
		return bytes, nil, false
	})
}
