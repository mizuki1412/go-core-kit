package serialkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"go.bug.st/serial"
)

type Config struct {
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
	COMName  string
}

var connect serial.Port

//var config Config

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

/// 一次数据返回，receive写在send之前，用channel实现数据异步返回
// handle是判断数据是否接收完成的函数，参数1表示全部数据，参数2表示收到的一段数据
// timeoutMill 超时时间，0表示不处理超时
func Receive(handle func([]byte, []byte) ([]byte, bool), timeoutMill int64) chan []byte {
	data := make(chan []byte)
	if timeoutMill > 0 {
		go func() {
			timekit.Sleep(timeoutMill)
			logkit.Error("serial read timeout")
			data <- nil
		}()
	}
	go func() {
		all := make([]byte, 0)
		buff := make([]byte, 100)
		for {
			n, err := connect.Read(buff)
			if err != nil {
				logkit.Error("serial read error: " + err.Error())
				data <- nil
			}
			if n == 0 {
				data <- nil
			}
			var ok bool
			all, ok = handle(all, buff[:n])
			if ok {
				data <- all
			}
		}
	}()
	return data
}

func Close() {
	_ = connect.Close()
	connect = nil
}
