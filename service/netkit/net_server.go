package netkit

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/spf13/cast"
)

type NetServer struct {
	Port      int32
	OnConnect func(c *connection.Connection)
	OnMessage func(c *connection.Connection, ctx interface{}, data []byte) (out []byte)
	OnClose   func(c *connection.Connection)
	// 数据message前的组包拆包
	UnPacket func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte)
	Packet   func(c *connection.Connection, data []byte) []byte
}

func (th *NetServer) Run() {
	handler0 := &handler{
		OnConnectFunc: th.OnConnect,
		OnMessageFunc: th.OnMessage,
		OnCloseFunc:   th.OnClose,
	}
	var s *gev.Server
	var err error
	options := []gev.Option{
		gev.Address(":" + cast.ToString(th.Port)),
		gev.NumLoops(-1),
	}
	if th.UnPacket != nil || th.Packet != nil {
		options = append(options, gev.Protocol(&DefaultProtocol{
			UnPacketFunc: th.UnPacket,
			PacketFunc:   th.Packet,
		}))
	}
	s, err = gev.NewServer(
		handler0,
		options...,
	)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	s.Start()
}

type handler struct {
	OnConnectFunc func(c *connection.Connection)
	OnMessageFunc func(c *connection.Connection, ctx interface{}, data []byte) (out []byte)
	OnCloseFunc   func(c *connection.Connection)
}

func (s *handler) OnConnect(c *connection.Connection) {
	if s.OnConnectFunc == nil {
		return
	}
	s.OnConnectFunc(c)
}
func (s *handler) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	// out 用于返回值
	if s.OnMessageFunc == nil {
		return
	}
	return s.OnMessageFunc(c, ctx, data)
}
func (s *handler) OnClose(c *connection.Connection) {
	if s.OnCloseFunc == nil {
		return
	}
	s.OnCloseFunc(c)
}

type DefaultProtocol struct {
	UnPacketFunc func(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte)
	PacketFunc   func(c *connection.Connection, data []byte) []byte
}

func (d *DefaultProtocol) UnPacket(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if d.UnPacketFunc == nil {
		defer buffer.RetrieveAll()
		return nil, buffer.Bytes()
	}
	// todo 粘包？
	return d.UnPacketFunc(c, buffer)
}
func (d *DefaultProtocol) Packet(c *connection.Connection, data []byte) []byte {
	if d.PacketFunc == nil {
		return data
	}
	return d.PacketFunc(c, data)
}
