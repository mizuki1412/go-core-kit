package bridge

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/spf13/cast"
	"net/http"
)

const EventPublic = "event:public"
const RoomPublic = "public"

var Server *socketio.Server
var ConnectFun = func(s socketio.Conn) {
	Clients = append(Clients, s)
}
var Clients = make([]socketio.Conn, 0, 5)

func init() {
	Server = socketio.NewServer(nil)
	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logkit.Info("websocket connected:" + s.ID())
		// room name
		s.Join(RoomPublic)
		ConnectFun(s)
		return nil
	})
	Server.OnError("/", func(s socketio.Conn, e error) {
		logkit.Info("websocket socket error:" + e.Error())
	})
	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logkit.Info("websocket closed: " + reason)
	})
}

/**
 *  同一个socket的同一个事件消息，是同步顺序执行的。
 * 	约定：event:public 作为通用的消息通道
 *	需要在start前设置
 */
func SetEventPublicHandle(fun func(req *MsgReq) string) *socketio.Server {
	Server.OnEvent("/", EventPublic, func(s socketio.Conn, msg string) (ret string) {
		defer func() {
			if err := recover(); err != nil {
				var msg string
				if e, ok := err.(exception.Exception); ok {
					msg = e.Msg
					// 带代码位置信息
					logkit.Error(e.Error())
					ret = RetErr(e.Error())
				} else {
					msg = cast.ToString(err)
					logkit.Error(msg)
					ret = msg
				}
			}
		}()
		req := &MsgReq{}
		err := jsonkit.ParseObj(msg, req)
		if err != nil {
			return RetErr("json error")
		}
		return fun(req)
	})
	return Server
}

// Start 开启websocket server并配置http
func Start() {
	go Server.Serve()
	// defer Server.Close()
	// restkit方式：根目录下，而非proxyGroup下； 和其他action共存
	socketHandle := func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		Server.ServeHTTP(w, r)
	}
	// 和rest base地址区分开; POST和GET都可能
	restkit.GetRouter().Any("/socket.io/**", func(ctx *context.Context) {
		socketHandle(ctx.Proxy.Writer, ctx.Proxy.Request)
	})
}
