package cmd

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"github.com/spf13/cobra"
	"log"
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
		//excelkit.Export(excelkit.Param{
		//	Title: "报警信息表",
		//	Keys:  []string{"key1:val1:10", "key2:val2:10", "key3:val3:20", "key4:val4:10", "key5:val5:10", "key6:val6:10"},
		//	Data: []map[string]interface{}{
		//		{"key1": "xcscs"},
		//		{"key2": 12},
		//		{"key3": 12.4444},
		//	},
		//}, nil)
	},
}

func testMapMerge() {
	map1 := map[string]interface{}{
		"key1": 11,
		"key2": "string",
		"key3": map[string]interface{}{
			"k1": 111,
			"k2": nil,
			"k3": map[string]interface{}{
				"k11": 1111,
			},
		},
	}
	map2 := map[string]interface{}{
		"key1": 22,
		"key2": nil,
		"key4": map[string]interface{}{
			"k4": 222,
			"k3": map[string]interface{}{
				"k11": 2222,
			},
		},
	}
	mapkit.Merge(map1, map2)
	log.Println(jsonkit.ToString(map1))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
