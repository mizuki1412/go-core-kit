package cmd

import (
	"github.com/mizuki1412/go-core-kit/tool-local/pgsql"
	"github.com/spf13/cobra"
)

func PGSqlToStructCMD(sqlFile, destFile string) *cobra.Command {
	if sqlFile == "" {
		sqlFile = "/Users/ycj/Downloads/init.sql"
	}
	if destFile == "" {
		destFile = "/Users/ycj/Downloads/init.go"
	}
	return &cobra.Command{
		Use:   "pgss",
		Short: "pg sql to struct",
		Run: func(cmd *cobra.Command, args []string) {
			pgsql.SQL2Struct(sqlFile, destFile)
		},
	}
}
