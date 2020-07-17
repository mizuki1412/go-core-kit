package cmd

import (
	"github.com/mizuki1412/go-core-kit/tool-local/markdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mdcmd)
}

var mdcmd = &cobra.Command{
	Use:   "md",
	Short: "handle readme",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		markdown.Test()
	},
}
