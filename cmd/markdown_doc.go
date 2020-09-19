package cmd

import (
	"github.com/mizuki1412/go-core-kit/tool-local/markdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mdocCmd)
}

var mdocCmd = &cobra.Command{
	Use:   "mdoc",
	Short: "gen readme md doc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		markdown.Run("go-core-kit 说明文档")
	},
}
