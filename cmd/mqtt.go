package cmd

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/mqttkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"log"
)

func MQTTTestCMD() *cobra.Command {
	mqttCmd := &cobra.Command{
		Use: "mqtt",
		Run: func(cmd *cobra.Command, args []string) {
			initkit.BindFlags(cmd)
			if configkit.GetStringD("topic.sub") != "" {
				mqttkit.Subscribe(configkit.GetStringD("topic.sub"), 2, func(client MQTT.Client, message MQTT.Message) {
					log.Println(message.Topic(), string(message.Payload()))
				})
			}
			if configkit.GetStringD("topic.pub") != "" {
				for {
					err := mqttkit.Publish(configkit.GetStringD("topic.pub"), 2, false, configkit.GetStringD("send"))
					if err != nil {
						log.Println(err.Error())
					}
					if configkit.GetInt("freq", 0) > 0 {
						timekit.Sleep(cast.ToInt64(configkit.GetInt("freq", 0)) * 1000)
					} else {
						break
					}
				}
			}
			select {}

		},
	}
	mqttCmd.Flags().String("topic.sub", "", "")
	mqttCmd.Flags().String("topic.pub", "", "发送的topic")
	mqttCmd.Flags().String("send", "", "发送的数据")
	mqttCmd.Flags().String("freq", "", "发送的频次/s")
	return mqttCmd
}
