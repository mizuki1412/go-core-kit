package mqttkit

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"time"
)

var client MQTT.Client

var subscribeList []func()

// 第一次连接
var first = true

func New() *MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(configkit.GetStringD(ConfigKeyMQTTBroker))
	opts.SetKeepAlive(time.Duration(1) * time.Minute)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Duration(5) * time.Second)
	opts.SetClientID(configkit.GetStringD(ConfigKeyMQTTClientID))
	opts.SetUsername(configkit.GetStringD(ConfigKeyMQTTUsername)).SetPassword(configkit.GetStringD(ConfigKeyMQTTPwd))
	var lostHan MQTT.OnConnectHandler = func(c MQTT.Client) {
		if first {
			first = false
			return
		}
		// 重连后重新订阅
		logkit.Info("mqtt reconnect, subs:" + cast.ToString(len(subscribeList)))
		for _, sub := range subscribeList {
			sub()
		}
	}
	opts.SetOnConnectHandler(lostHan)
	//opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
	//})
	//create and start a client using the above ClientOptions
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(exception.New(token.Error().Error()))
	}
	logkit.Info("mqtt connect success")
	return &client
}

func Subscribe(topic string, qos byte, callback MQTT.MessageHandler) {
	if client == nil {
		New()
	}
	f := func() {
		if token := client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
			logkit.Error(token.Error().Error())
		} else {
			logkit.Info("mqtt subscribe success: " + topic)
		}
	}
	subscribeList = append(subscribeList, f)
	f()
}

func Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if client == nil {
		New()
	}
	token := client.Publish(topic, qos, retained, payload)
	token.Wait()
	if token.Error() != nil {
		logkit.Error(token.Error().Error())
		return token.Error()
	}
	return nil
}
