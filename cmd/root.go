package cmd

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"log"
	"strings"
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
		//val := modbus.CRC16([]byte{01, 0x01, 0x01, 0x00})
		//logkit.Error()(bytekit.Bytes2HexArray([]byte{byte(val), byte(val >> 8)}))
		defer func() {
			if err := recover(); err != nil {
				var msg string
				if e, ok := err.(exception.Exception); ok {
					msg = e.Msg
					logkit.Error(e.Error())
				} else {
					msg = cast.ToString(err)
					excep := exception.New(msg, 3)
					logkit.Error(excep.Error())
				}
			}
		}()

		log.Print(strings.TrimSpace("  dsjd   "))
	},
}

func testJsonArr() {
	arr := &class.MapStringArr{}
	jsonkit.ParseObj(`[{"key":["a","b"]}]`, arr)
	log.Println(arr.Arr)
	for _, e := range arr.Arr {
		log.Println(e["key"].([]interface{}))
	}
	test()
}

func TestCatch() (ret int) {
	defer func() {
		if err := recover(); err != nil {
			ret = 11
		}
	}()
	panic(exception.New("123"))
	return 10
}

func test() {
	//panic(exception.New("test"))
	arr := []string{}
	arr[1] = ""
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
	logkit.Error(jsonkit.ToString(map1))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
