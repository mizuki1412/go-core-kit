package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cobra"
	"time"
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
		str := "{ \"key\": 123.00020010002224567, \"key2\":123, \"key3\":1592286542000  }"
		dt := class.Time{}
		dt.Set(time.Now())
		d := Data{}
		jsonkit.ParseObj(str, &d)
		fmt.Println(d)
		fmt.Println(jsonkit.ToString(d))
		fmt.Println(d.Key.Decimal.Float64())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
