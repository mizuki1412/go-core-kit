package main

import (
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/spf13/cobra"
	"log"
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
			dao := userdao.New()
			u := &model.User{}
			u.Name.Set("test33")
			u.Off.Set(1)
			re := dao.List(userdao.ListParam{
				RoleId:      12,
				Departments: []int32{1, 2, 3},
			})
			log.Println(jsonkit.ToString(re))
		},
	})
	cli.Execute()
}
