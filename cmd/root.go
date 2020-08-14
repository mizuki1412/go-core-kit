package cmd

import (
	"bytes"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
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

		//log.Println(time.Now().)
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

func decodeGBK() {
	data := []byte{0xB2, 0xBB, 0xCA, 0xC7, 0xC4, 0xDA, 0xB2, 0xBF, 0xBB, 0xF2, 0xCD, 0xE2, 0xB2, 0xBF, 0xC3, 0xFC, 0xC1, 0xEE, 0xA3, 0xAC, 0xD2, 0xB2, 0xB2, 0xBB, 0xCA, 0xC7, 0xBF, 0xC9, 0xD4, 0xCB, 0xD0, 0xD0, 0xB5, 0xC4, 0xB3, 0xCC, 0xD0, 0xF2, 0xBB, 0xF2, 0xC5, 0xFA, 0xB4, 0xA6, 0xC0, 0xED, 0xCE, 0xC4, 0xBC, 0xFE, 0xA1, 0xA3}
	r := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	b, err := ioutil.ReadAll(r)
	log.Println(err)
	log.Println(string(b))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
