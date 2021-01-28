package cmd

import (
	"github.com/mizuki1412/go-core-kit/tool-local/frontdao"
	"github.com/spf13/cobra"
)

func FrontDaoCMD(url string) *cobra.Command {
	return &cobra.Command{
		Use:   "genDao",
		Short: "generate front dao",
		Run: func(cmd *cobra.Command, args []string) {
			frontdao.Gen(url)
		},
	}
}
