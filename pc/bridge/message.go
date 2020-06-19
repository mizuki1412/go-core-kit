package bridge

import "github.com/mizuki1412/go-core-kit/library/jsonkit"

/***
消息格式：
-> code:string, data:map
<- result:boolean, data:map, message:string
*/
type MsgReq struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
}

type MsgRes struct {
	Result  bool        `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// handle public event return
func RetErr(msg string) string {
	return jsonkit.ToString(MsgRes{
		Result:  false,
		Message: msg,
	})
}

// handle public event return
func RetSuccess(data interface{}) string {
	return jsonkit.ToString(MsgRes{
		Result: true,
		Data:   data,
	})
}

// send broadcast message
func Send(req MsgReq) {
	Server.BroadcastToRoom("/", RoomPublic, EventPublic, jsonkit.ToString(req))
}
