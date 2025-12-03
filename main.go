package main

import (
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/cmd"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/admindivision"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/download"
	"github.com/mizuki1412/go-core-kit/v2/mod/user"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit"
	"github.com/spf13/cobra"
)

func main() {
	cli.RootCMD(&cobra.Command{
		Use: "main",
		Run: func(cmd *cobra.Command, args []string) {
			restkit.AddActions(user.All()...)
			restkit.AddActions(download.Init)
			restkit.AddActions(admindivision.Init)
			_ = restkit.Run()
		},
	})
	c1 := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	c1.Flags().String("test", "", "")
	c1.Flags().String("test1", "", "")
	cli.AddChildCMD(c1)

	cli.AddChildCMD(cmd.FrontDaoCMDNext("http://localhost:10000/v3/api-docs"))
	cli.Execute()
}
