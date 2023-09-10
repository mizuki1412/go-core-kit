package serialkit

import (
	"github.com/albenik/go-serial/v2"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"time"
)

type Config struct {
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
	COMName  string
}

var connect *serial.Port

// 外部控制receive逻辑中断
var interrupt bool

func ListPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		panic(exception.New("list serial port error"))
	}
	list := make([]string, 0, len(ports))
	for _, port := range ports {
		list = append(list, port)
	}
	return list
}

func Open(config0 Config) {
	//config = config0
	if connect == nil {
		var err error
		connect, err = serial.Open(
			config0.COMName,
			serial.WithBaudrate(config0.BaudRate),
			serial.WithStopBits(config0.StopBits),
			serial.WithDataBits(config0.DataBits),
			serial.WithParity(config0.Parity),
			serial.WithReadTimeout(0))
		if err != nil {
			panic(exception.New("open serial error"))
		}
	}
}

func Send(data []byte) {
	if connect == nil {
		panic(exception.New("please open serial first"))
	}
	_, err := connect.Write(data)
	if err != nil {
		panic(exception.New("serial data send error"))
	}
}

// Receive 一次数据返回，receive写在send之前，用channel实现数据异步返回
// handle是判断数据是否接收完成的函数，参数1表示全部数据，参数2表示收到的一段数据
// timeoutMill 超时时间，0表示不处理超时
func Receive(handle func([]byte, []byte) ([]byte, bool), timeoutMill int) chan []byte {
	if connect == nil {
		panic(exception.New("please open serial first"))
	}
	interrupt = false
	chRun := make(chan []byte)
	now := time.Now()
	go func(ch chan []byte) {
		all := make([]byte, 0)
		buff := make([]byte, 100)
		// non-block
		for {
			if interrupt {
				logkit.Error("serial interrupt")
				ch <- nil
				close(ch)
				break
			}
			n, err := connect.Read(buff)
			if err != nil {
				logkit.Error(err.Error())
				ch <- nil
				close(ch)
				break
			}
			if n == 0 {
				// timeout
				if time.Now().After(now.Add(time.Duration(timeoutMill) * time.Millisecond)) {
					ch <- nil
					close(ch)
					break
				}
				timekit.Sleep(500)
				continue
			}
			var ok bool
			// 分段加入all，并判断是否接收完成
			all, ok = handle(all, buff[:n])
			if ok {
				ch <- all
				close(ch)
				break
			}
		}
	}(chRun)
	return chRun
}

func Close() {
	_ = connect.Close()
	connect = nil
}

func Interrupt() {
	interrupt = true
}
