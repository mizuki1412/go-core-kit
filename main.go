package main

import (
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/cmd"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/admindivision"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/download"
	"github.com/mizuki1412/go-core-kit/v2/mod/user"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/cachekit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
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
		Use: "test1",
		Run: func(cmd *cobra.Command, args []string) {
			println(configkit.GetString("test"))
			println(configkit.GetString("test1"))
		},
	}
	c1.Flags().String("test", "", "")
	c1.Flags().String("test1", "", "")
	cli.AddChildCMD(c1)

	cli.AddChildCMD(cmd.FrontDaoCMDNext("http://localhost:10000/v3/api-docs"))
	cli.Execute()
}

func test() []*model.User {
	return cachekit.Wrapper(cachekit.WrapParam{
		Key: "abc",
		Ttl: 0,
	}, func() any {
		dao := userdao.New(userdao.ResultDefault)
		return dao.List(userdao.ListParam{})
	}).([]*model.User)
}
