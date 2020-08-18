package bridge

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"mime"
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
	var err error
	Server, err = socketio.NewServer(nil)
	if err != nil {
		logkit.Fatal(err)
	}
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

// 开启websocket server并配置http todo
func Start() {
	go Server.Serve()
	// defer Server.Close()
	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		Server.ServeHTTP(w, r)
	})
	// local win ui web, 默认在ui todo
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	_ = mime.AddExtensionType(".js", "text/javascript")
}
