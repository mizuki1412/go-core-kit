package main

import (
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
	"github.com/mizuki1412/go-core-kit/mod/user/dao/userdao"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service/restkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
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
			//testSQL()
			list, total := userdao.New().Select().OrderBy("id").Page(sqlkit.Page{
				PageNum: 2, PageSize: 2,
			})
			log.Println(jsonkit.ToString(list), total)
		},
	})
	cli.Execute()
}

func testSQL() {
	dao := userdao.New()
	u := &model.User{}
	u.Name.Set("test33")
	u.Role = &model.Role{Id: 1}
	dao.InsertObj(u)
	u.Extend.Put("abc", 1)
	dao.UpdateObj(u)
	log.Println(jsonkit.ToString(dao.SelectOneById(u.Id)))
	log.Println(jsonkit.ToString(dao.Select().Where("name=?", u.Name).List()))
}
