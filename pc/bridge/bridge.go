package bridge

import (
	socketio "github.com/googollee/go-socket.io"
	"log"
	"mime"
	"net/http"
)

const EventPublic = "event:public"
const RoomPublic = "public"
var Server *socketio.Server
var connectFun func(s socketio.Conn)

func init() {
	var err error
	Server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("websocket connected:", s.ID())
		// room name
		s.Join(RoomPublic)
		connectFun(s)
		return nil
	})
	Server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("websocket socket error:", e)
	})
	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("websocket closed: ", reason)
	})
}

/**
 *  同一个socket的同一个事件消息，是同步顺序执行的。
 * 约定：event:public 作为通用的消息通道
 */
func SetEventPublicHandle(fun func(s socketio.Conn, msg string)) *socketio.Server {
	//Server.OnEvent("/", EventPublic, func(s socketio.Conn, msg string) string{
	//	if !gjson.Valid(msg) {
	//		return jsonkit.ParseString(MsgRes{Result:false, Message:"json error"})
	//	}
	//	req := &MsgReq{}
	//	jsonkit.ParseObj(msg, req)
	//	return jsonkit.ParseString(HandlePublicMsg(*req))
	//})
	Server.OnEvent("/", EventPublic, fun)
	return Server
}
func SetConnectHandle(fun func(s socketio.Conn)){
	connectFun = fun
}
// 开启websocket server并配置http todo
func Start() {
	go Server.Serve()
	defer Server.Close()
	http.HandleFunc("/socket.io/", func (w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		Server.ServeHTTP(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("./static")))
	_ = mime.AddExtensionType(".js", "text/javascript")
}

/***
消息格式：
-> code:string, data:map
<- result:boolean, data:map, message:string
*/
type MsgReq struct {
	Code string		 `json:"code"`
	Data interface{} `json:"data"`
}

type MsgRes struct {
	Result bool		 `json:"result"`
	Message string	 `json:"message"`
	Data interface{} `json:"data"`
}