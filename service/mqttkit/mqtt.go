package mqttkit

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"mizuki/framework/core-kit/service/configkit"
	"mizuki/framework/core-kit/service/logkit"
)

var client MQTT.Client

func New() *MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(configkit.GetStringD(ConfigKeyMQTTBroker))
	opts.SetClientID(configkit.GetStringD(ConfigKeyMQTTClientID))
	opts.SetUsername(configkit.GetStringD(ConfigKeyMQTTUsername)).SetPassword(configkit.GetStringD(ConfigKeyMQTTPwd))
	//opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
	//})
	//create and start a client using the above ClientOptions
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logkit.Fatal(token.Error().Error())
	}
	logkit.Info("mqtt connect success")
	return &client
}

func Subscribe(topic string, qos byte, callback MQTT.MessageHandler) {
	if client == nil {
		New()
	}
	if token := client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		logkit.Error(token.Error().Error())
	}
	logkit.Info("mqtt subscribe success: " + topic)
}

func Publish(topic string, qos byte, retained bool, payload interface{}) {
	token := client.Publish(topic, qos, retained, payload)
	token.Wait()
	if token.Error() != nil {
		logkit.Error(token.Error().Error())
	}
}
