package main

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cli.RootCMD(&cobra.Command{
		Use: "main",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(viper.GetString(configkey.TimeLocation))
			fmt.Println(viper.IsSet(configkey.TimeLocation))
		},
	})
	cli.AddChildCMD(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("123")
		},
	})
	cli.AddChildCMD(&cobra.Command{
		Use: "web",
		Run: func(cmd *cobra.Command, args []string) {
			restkit.AddActions(user.All()...)
			restkit.AddActions(download.Init)
			_ = restkit.Run()
		},
	})
	cli.Execute()
}
