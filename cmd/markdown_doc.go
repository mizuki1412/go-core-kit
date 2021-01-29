package cmd

import (
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/tool-local/markdown"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MarkdownDocCMD(title string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mdoc",
		Short: `gen readme md doc`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
			if configkit.GetStringD("dest") == "" {
				logkit.Fatal("参数未指定")
			}
			markdown.Run(title, configkit.GetStringD("dest"))
		},
	}
	cmd.Flags().StringP("dest", "", "", "生成目标路径 /xx/xx")
	return cmd
}
