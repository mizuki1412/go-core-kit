package ssehelper

import (
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"sync"
)

// 在线客户端
var clientChannels map[string]chan string

var once sync.Once

func getChanMap() map[string]chan string {
	if clientChannels == nil {
		once.Do(func() {
			clientChannels = make(map[string]chan string)
		})
	}
	return clientChannels
}

func AddClient(clientId string, ctx *context.Context) {
	cc := getChanMap()
	if _, ok := cc[clientId]; !ok {
		cc[clientId] = make(chan string)
		logkit.Info("SSE Client add: " + clientId)
		closeNotify := ctx.Proxy.Request.Context().Done()
		go func() {
			<-closeNotify
			RemoveClient(clientId)
		}()
	} else {
		logkit.Info("SSE Client already add: " + clientId)
	}
}

func ServiceClient(clientId string, ctx *context.Context) {
	AddClient(clientId, ctx)
	for msg := range clientChannels[clientId] {
		ctx.SendSSE(msg)
	}
}

func RemoveClient(clientId string) {
	cc := getChanMap()
	if _, ok := cc[clientId]; ok {
		delete(cc, clientId)
		logkit.Info("SSE Client close: " + clientId)
	}
}

func ToSend(clientId string, msg string) {
	if c, ok := clientChannels[clientId]; ok {
		c <- msg
	}
}
