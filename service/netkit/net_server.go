package netkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/panjf2000/gnet"
)

type NetServer struct {
	Port      int32
	OnConnect func(c gnet.Conn) (out []byte, action gnet.Action)
	OnMessage func(frame []byte, c gnet.Conn) (out []byte, action gnet.Action)
	OnClose   func(c gnet.Conn, err error) (action gnet.Action)
	// 数据message前的组包拆包
	//UnPacket func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte)
	//Packet   func(c *connection.Connection, data []byte) []byte
}

type handler struct {
	*gnet.EventServer
	OnConnectFunc func(c gnet.Conn) (out []byte, action gnet.Action)
	OnMessageFunc func(frame []byte, c gnet.Conn) (out []byte, action gnet.Action)
	OnCloseFunc   func(c gnet.Conn, err error) (action gnet.Action)
}

func (th *handler) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	logkit.Info(fmt.Sprintf("net server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop))
	return
}

func (th *handler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	if th.OnConnectFunc == nil {
		return
	}
	return th.OnConnectFunc(c)
}

func (th *handler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	if th.OnCloseFunc == nil {
		return
	}
	return th.OnCloseFunc(c, err)
}

func (th *handler) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if th.OnMessageFunc == nil {
		return
	}
	return th.OnMessageFunc(frame, c)
}

func (th *NetServer) Run() {
	handler0 := &handler{
		OnConnectFunc: th.OnConnect,
		OnMessageFunc: th.OnMessage,
		OnCloseFunc:   th.OnClose,
	}
	//options := []gev.Option{
	//	gev.Address(":" + cast.ToString(th.Port)),
	//	gev.NumLoops(-1),
	//}
	err := gnet.Serve(
		handler0,
		fmt.Sprintf("tcp://:%d", th.Port),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true))
	if err != nil {
		panic(exception.New(err.Error()))
	}
}

//type DefaultProtocol struct {
//	UnPacketFunc func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte)
//	PacketFunc   func(c *connection.Connection, data []byte) []byte
//}
//
//func (d *DefaultProtocol) UnPacket(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
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
