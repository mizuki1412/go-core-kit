package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
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
		var arr []interface{}
		jsonkit.JSON().Unmarshal([]byte("[1,2,3]"), &arr)
		fmt.Printf("%T\n", arr[0])
		log.Println(arraykit.Delete(arr, 2))
		//math.Dim()
		//log.Println(arraykit.Delete(arr, cast.ToInt32(2)))
		//log.Println(arraykit.DeleteAt(arr, 2))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
