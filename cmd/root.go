package cmd

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/excelkit"
	"github.com/spf13/cobra"
)

type Data struct {
	Key  class.Decimal
	Key2 class.Int64
	Key3 class.Time
}

var rootCmd = &cobra.Command{
	Use:   "go-core-kit",
	Short: "go core kit test",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		excelkit.Export(excelkit.Param{
			Title: "报警信息表",
			Keys:  []string{"key1:val1:10", "key2:val2:10", "key3:val3:20", "key4:val4:10", "key5:val5:10", "key6:val6:10"},
			Data: []map[string]interface{}{
				{"key1": "xcscs"},
				{"key2": 12},
				{"key3": 12.4444},
			},
		}, nil)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
