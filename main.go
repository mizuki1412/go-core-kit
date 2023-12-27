package main

import (
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-core-kit/snippet"
	"github.com/spf13/cobra"
	"reflect"
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
			us := []*model.User{}
			rt := reflect.TypeOf(us)
			println(rt.Kind().String())
			println(rt.Elem().Name(), rt.Elem().Kind().String())
			restkit.AddActions(snippet.Init)
			_ = restkit.Run()
		},
	})
	cli.Execute()
}
