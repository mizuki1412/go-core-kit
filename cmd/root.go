package cmd

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/iot/modbus"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
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

		val := modbus.CRC16([]byte{01, 0x06, 0x09, 0xc5, 0x00, 0x01})
		log.Println(bytekit.Bytes2HexArray([]byte{byte(val), byte(val >> 8)}))

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
