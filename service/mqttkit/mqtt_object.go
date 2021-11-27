package mqttkit

import (
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"sync"
	"time"
)

// Client 适用于自行管理mqtt连接对象的场景
type Client struct {
	C             MQTT.Client
	SubscribeList []func()
	First         bool // 第一次连接
	Id            string
}

type ConnectParam struct {
	Broker   string
	Id       string
	Username string
	Pwd      string
}

// 用于记录创建过的clients
var allClients = map[MQTT.Client]*Client{}
var allClientsMux sync.RWMutex

func NewClient(param ConnectParam) *Client {
	newClient := &Client{First: true, Id: param.Id}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(param.Broker)
	opts.SetKeepAlive(time.Duration(1) * time.Minute)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Duration(5) * time.Second)
	opts.SetClientID(param.Id)
	opts.SetUsername(param.Username).SetPassword(param.Pwd)
	var lostHan MQTT.OnConnectHandler = func(c MQTT.Client) {
		// todo 测试c是否前后一致
		allClientsMux.RLock()
		defer allClientsMux.RUnlock()
		if cl, ok := allClients[c]; ok {
			// 第一次连接不处理
			if cl.First {
				cl.First = false
				return
			}
			// 重连后重新订阅
			logkit.Info(fmt.Sprintf("mqtt reconnect: %s, subs:%s", cl.Id, cast.ToString(len(cl.SubscribeList))))
			for _, sub := range cl.SubscribeList {
				sub()
			}
		}
	}
	opts.SetOnConnectHandler(lostHan)
	newClient.C = MQTT.NewClient(opts)
	if token := newClient.C.Connect(); token.Wait() && token.Error() != nil {
		panic(exception.New(token.Error().Error()))
	}
	logkit.Info("mqtt connect success")
	allClientsMux.Lock()
	defer allClientsMux.Unlock()
	allClients[newClient.C] = newClient
	return newClient
}

func (th *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) {
	if th.C == nil {
		logkit.Error("please newClient first")
		return
	}
	f := func() {
		if token := th.C.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
			logkit.Error(token.Error().Error())
		} else {
			logkit.Info("mqtt subscribe success: " + topic)
		}
	}
	th.SubscribeList = append(th.SubscribeList, f)
	f()
}

func (th *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if th.C == nil {
		err := errors.New("please newClient first")
		logkit.Error(err.Error())
		return err
	}
	token := th.C.Publish(topic, qos, retained, payload)
	token.Wait()
	if token.Error() != nil {
		logkit.Error(token.Error().Error())
		return token.Error()
	}
	return nil
}
