package cmd

import (
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/spf13/cobra"
	"net/http"
)

func WebStaticServerCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web-server",
		Short: "本地静态web服务器",
		Run: func(cmd *cobra.Command, args []string) {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, r.URL.Path[1:])
			})
			_ = http.ListenAndServe(":"+configkit.GetString("port"), nil)
		},
	}
	cmd.Flags().String("port", "", "端口")
	_ = cmd.MarkFlagRequired("port")
	return cmd
}
