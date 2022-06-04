package test

import (
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/spf13/cobra"
)

func init() {
	initkit.DefFlags(webCMD)
	rootCmd.AddCommand(webCMD)
}

var webCMD = &cobra.Command{
	Use: "web",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		restkit.AddActions(user.All()...)
		restkit.AddActions(download.Init)
		_ = restkit.Run()
	},
}
