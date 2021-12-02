package serialkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"go.bug.st/serial"
	"time"
)

type Config struct {
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
	COMName  string
}

var connect serial.Port

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
	mode := &serial.Mode{
		BaudRate: config0.BaudRate,
		StopBits: config0.StopBits,
		DataBits: config0.DataBits,
		Parity:   config0.Parity,
	}
	var err error
	connect, err = serial.Open(config0.COMName, mode)
	if err != nil {
		panic(exception.New("open serial error"))
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
func Receive(handle func([]byte, []byte) ([]byte, bool), timeoutMill int64) chan []byte {
	data := make(chan []byte)
	go func() {
		// 实际执行
		chRun := make(chan []byte)
		go func(ch chan []byte) {
			all := make([]byte, 0)
			buff := make([]byte, 100)
			for {
				n, err := connect.Read(buff)
				if err != nil {
					logkit.Error(exception.New(err.Error()))
					ch <- nil
					close(ch)
				}
				if n == 0 {
					ch <- nil
					close(ch)
				}
				var ok bool
				all, ok = handle(all, buff[:n])
				if ok {
					ch <- all
					close(ch)
					break
				}
			}
		}(chRun)
		select {
		case re := <-chRun:
			data <- re
			close(data)
		case <-time.After(time.Duration(timeoutMill) * time.Millisecond):
			data <- nil
			close(data)
		}
	}()
	return data
}

func Close() {
	_ = connect.Close()
	connect = nil
}
