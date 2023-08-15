package netkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/panjf2000/gnet/v2"
	"github.com/spf13/cast"
)

type NetServer struct {
	gnet.BuiltinEventEngine
	eng         gnet.Engine
	ProtoSchema string `description:"tcp/udp"`
	Port        int32

	TrafficHandler func(c gnet.Conn)
	//OnConnect func(c gnet.Conn) (out []byte, action gnet.Action)
	//OnMessage func(frame []byte, c gnet.Conn) (out []byte, action gnet.Action)
	//OnClose   func(c gnet.Conn, err error) (action gnet.Action)
	// 数据message前的组包拆包
	//UnPacket func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (any, []byte)
	//Packet   func(c *connection.Connection, data []byte) []byte
}

func (th *NetServer) OnBoot(eng gnet.Engine) gnet.Action {
	th.eng = eng
	logkit.Info("net server is listening on " + cast.ToString(th.Port))
	return gnet.None
}

func (th *NetServer) OnTraffic(c gnet.Conn) gnet.Action {
	if th.TrafficHandler != nil {
		th.TrafficHandler(c)
	}
	return gnet.None
}

func (th *NetServer) Run() {
	if th.ProtoSchema == "" {
		th.ProtoSchema = "tcp"
	}
	err := gnet.Run(
		th,
		fmt.Sprintf("%s://:%d", th.ProtoSchema, th.Port),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true))
	if err != nil {
		panic(exception.New(err.Error()))
	}
}

//type DefaultProtocol struct {
//	UnPacketFunc func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (any, []byte)
//	PacketFunc   func(c *connection.Connection, data []byte) []byte
//}
//
//func (d *DefaultProtocol) UnPacket(c *connection.Connection, buffer *ringbuffer.RingBuffer) (any, []byte) {
//	if d.UnPacketFunc == nil {
//		defer buffer.RetrieveAll()
//		return nil, buffer.Bytes()
//	}
//	// todo 粘包？
//	return d.UnPacketFunc(c, buffer)
//}
//func (d *DefaultProtocol) Packet(c *connection.Connection, data []byte) []byte {
//	if d.PacketFunc == nil {
//		return data
//	}
//	return d.PacketFunc(c, data)
//}
