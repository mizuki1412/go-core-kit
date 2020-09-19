package cmd

import (
	"github.com/mizuki1412/go-core-kit/tool-local/pgsql"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	rootCmd.AddCommand(pgssCmd)
}

var pgssCmd = &cobra.Command{
	Use:   "pgss",
	Short: "postgres sql to struct",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		pgsql.SQL2Struct("/Users/ycj/Downloads/demo.sql", "/Users/ycj/Downloads/dest.go")
		println(time.Since(t).Milliseconds())
	},
}
