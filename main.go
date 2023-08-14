package main

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/cmd"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/spf13/cobra"
)

func main() {
	cli.RootCMD(&cobra.Command{
		Use: "main",
		Run: func(cmd *cobra.Command, args []string) {
			restkit.AddActions(user.All()...)
			restkit.AddActions(download.Init)
			_ = restkit.Run()
		},
	})
	cli.AddChildCMD(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(configkit.GetString(configkey.ProjectDir))
		},
	})
	cli.AddChildCMD(cmd.TCPServerCMD())
	cli.Execute()
}
