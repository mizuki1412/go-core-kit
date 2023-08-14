package cmd

import (
	"github.com/mizuki1412/go-core-kit/cli/commandkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cobra"
	"net/http"
)

func WebStaticServerCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web-server",
		Short: "本地静态web服务器",
		Run: func(cmd *cobra.Command, args []string) {
			commandkit.BindFlags(cmd)
			if configkit.GetStringD("port") == "" {
				logkit.Fatal("port参数缺失")
			}
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, r.URL.Path[1:])
			})
			_ = http.ListenAndServe(":"+configkit.GetStringD("port"), nil)
		},
	}
	cmd.Flags().StringP("port", "", "", "端口")
	return cmd
}
