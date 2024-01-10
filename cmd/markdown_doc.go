package cmd

import (
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/mizuki1412/go-core-kit/v2/tool-local/markdown"
	"github.com/spf13/cobra"
)

func MarkdownDocCMD(title string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mdoc",
		Short: `gen readme md doc`,
		Run: func(cmd *cobra.Command, args []string) {
			markdown.Run(title, configkit.GetString("dest"))
		},
	}
	cmd.Flags().String("dest", "", "生成目标路径 /xx/xx")
	_ = cmd.MarkFlagRequired("dest")
	return cmd
}
