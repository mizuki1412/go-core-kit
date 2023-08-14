package main

import (
	"github.com/bytedance/sonic"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/cli"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/mod/common/download"
	"github.com/mizuki1412/go-core-kit/mod/user"
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
			t := map[string]any{
				"id":  1.2,
				"abc": class.ArrInt{},
				"ccd": map[string]any{
					"aa": 1,
				},
			}
			log.Println(jsonkit.ToString(t))
			r, _ := sonic.GetFromString(jsonkit.ToString(t), "ccd", "aa")
			log.Println(r.Float64())
			log.Println(jsonkit.ToString(nil))
		},
	})
	cli.Execute()
}
