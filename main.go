package main

import (
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/service/netkit"
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
			server := &netkit.NetServer{
				Port: 12345,
			}
			server.Run()
		},
	})
	cli.Execute()
}
