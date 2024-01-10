package serialkit

import (
	"github.com/albenik/go-serial/v2"
	"github.com/mizuki1412/go-core-kit/v2/library/bytekit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
)

func Demo() {
	config := Config{
		BaudRate: 115200,
		StopBits: serial.OneStopBit,
		DataBits: 8,
		Parity:   serial.NoParity,
		COMName:  "COM2",
	}
	Open(config)
	ch := Receive(func(all []byte, buff []byte) ([]byte, bool) {
		arr := append(all, buff...)
		ok := false
		for i, e := range arr {
			if i < len(arr)-1 && e == 0x0d && arr[i+1] == 0x0a {
				ok = true
			}
		}
		return arr, ok
	}, 10000)
	Send([]byte{0x00, 0x01, 0x022, 0x0d, 0x0a})
	ret := <-ch
	if ret != nil {
		logkit.Info("receive: " + bytekit.Bytes2HexArray(ret))
	} else {
		logkit.Info("nil")
	}
}
