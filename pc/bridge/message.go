package bridge

import (
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
)

/*
**
消息格式：
-> code:string, data:map
<- result:boolean, data:map, message:string
*/
type MsgReq struct {
	Code string `json:"code"`
	Data any    `json:"data"`
}

type MsgRes struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// handle public event return
func RetErr(msg string) string {
	return jsonkit.ToString(MsgRes{
		Result:  false,
		Message: msg,
	})
}

// handle public event return
func RetSuccess(data any) string {
	return jsonkit.ToString(MsgRes{
		Result: true,
		Data:   data,
	})
}

// send broadcast message， server主动发送给client
func Send(req MsgReq) {
	// todo client callback未成功，socket.io-client不支持？目前只能server单向推送
	//for _,c:=range Clients{
	//	c.Emit(EventPublic, jsonkit.ToString(req), func (so socketio.Conn, data string) {
	//		log.Println("Client ACK with data: ", data)
	//	})
	//}
	Server.BroadcastToRoom("/", RoomPublic, EventPublic, jsonkit.ToString(req))
}
