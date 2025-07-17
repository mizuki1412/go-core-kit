package mqttkit

import (
	"crypto/tls"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/c"
	"github.com/mizuki1412/go-core-kit/v2/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
	"github.com/spf13/cast"
	"strings"
	"sync"
	"time"
)

// Client 适用于自行管理mqtt连接对象的场景
type Client struct {
	C             MQTT.Client
	SubscribeList []func()
	Id            string
	Param         ConnectParam
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
	opts := MQTT.NewClientOptions()
	if param.Broker == "" {
		panic(exception.New("请填写broker"))
	}
	if param.Id == "" {
		param.Id = cryptokit.ID()
	}
	newClient := &Client{Id: param.Id, Param: param}
	opts.AddBroker(param.Broker)
	opts.SetKeepAlive(time.Duration(5) * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Duration(5) * time.Second)
	opts.SetClientID(param.Id)
	opts.SetUsername(param.Username).SetPassword(param.Pwd)
	var lostHan MQTT.OnConnectHandler = func(c MQTT.Client) {
		allClientsMux.RLock()
		defer allClientsMux.RUnlock()
		// 连接成功后才会
		if cl, ok := allClients[c]; ok {
			// 重连后重新订阅
			logkit.Info(fmt.Sprintf("mqtt reconnect: %s, subs:%s", cl.Id, cast.ToString(len(cl.SubscribeList))))
			for _, sub := range cl.SubscribeList {
				sub()
			}
		}
	}
	opts.SetOnConnectHandler(lostHan)
	if strings.Index(param.Broker, "ssl:") == 0 {
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		})
	}
	newClient.C = MQTT.NewClient(opts)
	if token := newClient.C.Connect(); token.Wait() && token.Error() != nil {
		panic(exception.New(token.Error().Error()))
	}
	logkit.Info(fmt.Sprintf("mqtt connect success: %s, clientId:%s", param.Broker, param.Id))
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
	callback1 := func(client MQTT.Client, message MQTT.Message) {
		_ = c.RecoverFuncWrapper(func() {
			callback(client, message)
		})
	}
	f := func() {
		if token := th.C.Subscribe(topic, qos, callback1); token.Wait() && token.Error() != nil {
			logkit.Error(token.Error().Error())
		} else {
			logkit.Info(fmt.Sprintf("client<%s> subscribe: %s", th.Id, topic))
		}
	}
	th.SubscribeList = append(th.SubscribeList, f)
	f()
}

func (th *Client) Publish(topic string, qos byte, retained bool, payload any) error {
	if th.C == nil {
		err := errors.New("please newClient first")
		logkit.Error(err.Error())
		return err
	}
	token := th.C.Publish(topic, qos, retained, payload)
	token.WaitTimeout(time.Duration(1) * time.Minute)
	if token.Error() != nil {
		logkit.Error(token.Error().Error())
		return token.Error()
	}
	return nil
}
