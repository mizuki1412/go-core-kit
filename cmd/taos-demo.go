package cmd

import (
	"github.com/spf13/cobra"
)

func TaosDemo() *cobra.Command {
	mqttCmd := &cobra.Command{
		Use: "mqtt",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return mqttCmd
}
