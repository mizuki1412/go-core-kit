package framekit

type Decoder struct {
	Bytes       []byte
	Header      []byte
	FixedLength int
}

func NewDecoder(capacity int, header []byte, fixedLength int) *Decoder {
	return &Decoder{
		Bytes:       make([]byte, 0, capacity),
		Header:      header,
		FixedLength: fixedLength,
	}
}

func (th *Decoder) Take() []byte {
	if len(th.Bytes) < th.FixedLength {
		return nil
	}
	for i := 0; i < th.FixedLength-len(th.Header); i++ {
		count := 0
		for j := 0; j < len(th.Header); j++ {
			if th.Bytes[i+j] == th.Header[j] {
				count++
			} else {
				break
			}
		}
		if count == len(th.Header) {
			// 找到帧头 todo 没处理最长
			res := th.Bytes[i : i+th.FixedLength]
			th.Bytes = th.Bytes[i+th.FixedLength:]
			return res
		}
	}
	return nil
}

func (th *Decoder) Put(data []byte) {
	if len(data) > 0 {
		th.Bytes = append(th.Bytes, data...)
	}
}
