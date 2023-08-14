package mqttkit

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"sync"
)

var client *Client

var _once sync.Once

func New() *Client {
	_once.Do(func() {
		if configkit.GetStringD(configkey.MQTTBroker) == "" {
			panic(exception.New("请填写broker"))
		}
		client = NewClient(ConnectParam{
			Broker:   configkit.GetStringD(configkey.MQTTBroker),
			Id:       configkit.GetString(configkey.MQTTClientID, cryptokit.ID()),
			Username: configkit.GetStringD(configkey.MQTTUsername),
			Pwd:      configkit.GetStringD(configkey.MQTTPwd),
		})
	})
	return client
}

func Subscribe(topic string, qos byte, callback MQTT.MessageHandler) {
	if client == nil {
		New()
	}
	client.Subscribe(topic, qos, callback)
}

func Publish(topic string, qos byte, retained bool, payload any) error {
	if client == nil {
		New()
	}
	return client.Publish(topic, qos, retained, payload)
}

func GetClient() MQTT.Client {
	return client.C
}
