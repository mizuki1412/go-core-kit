package main

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/mizuki1412/go-core-kit/v2/cmd"
	"github.com/mizuki1412/go-core-kit/v2/library/timekit"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/admindivision"
	"github.com/mizuki1412/go-core-kit/v2/mod/common/download"
	"github.com/mizuki1412/go-core-kit/v2/mod/user"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"github.com/spf13/cobra"
	"log"
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
	cli.AddChildCMD(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			restkit.AddActions(user.All()...)
			restkit.AddActions(Init)
			_ = restkit.Run()
		},
	})
	cli.AddChildCMD(cmd.FrontDaoCMDNext("http://localhost:10000/v3/api-docs"))
	cli.Execute()
}

func Init(router *router.Router) {
	r := router.Group("/sse")
	//r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Get("/test/:id", sse)
	}
}

type Params struct {
	Id string
}

func sse(ctx *context.Context) {
	params := Params{}
	ctx.BindForm(&params)
	w := ctx.Proxy.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	closeNotify := ctx.Proxy.Request.Context().Done()
	go func() {
		<-closeNotify
		log.Println("SSE close for user = ", params.Id)
		return
	}()
	for {
		log.Println("sse send")
		_, err := w.WriteString(fmt.Sprintf("event: message\ndata: %s\n\n", "abv"))
		if err != nil {
			log.Println(err.Error())
		}
		w.Flush()
		timekit.Sleep(1000)
	}
}
