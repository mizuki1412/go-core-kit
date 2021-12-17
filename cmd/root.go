package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func init() {
	initkit.DefFlags(rootCmd)
	rootCmd.AddCommand(PGSqlToStructCMD("", ""))
	rootCmd.AddCommand(MarkdownDocCMD("go-core-kit 说明文档"))
	rootCmd.AddCommand(WebStaticServerCMD())
}

var rootCmd = &cobra.Command{
	Use: "go-core-kit",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		//go func() {
		//	ch := make(chan int)
		//	go func(c chan int) {
		//		defer func() {
		//			println("close")
		//		}()
		//		timekit.Sleep(10000)
		//		println("here")
		//		c <- 1
		//	}(ch)
		//	println("start")
		//	select {
		//	case <-ch:
		//		println("ok")
		//	case <-time.After(time.Duration(5) * time.Second):
		//		println("timeout")
		//	}
		//}()
		//select {}
		ch := make(chan string, 10)
		defer close(ch)

		go print_input(ch)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in := scanner.Text()
			ch <- in
			if in == "EOF" {
				fmt.Printf("exit\n")
				return
			}
		}
	},
}

func print_input(ch chan string) {
	fmt.Printf("Begin to print input.\n")
	for {
		select {
		case input_ := <-ch:
			fmt.Printf("get %v from standard input.\n", input_)
		case <-time.After(3 * time.Second): // 超时3秒没有获得数据，则退出程序，如果只是退出循环，可以return改为continue
			fmt.Printf("More than 3 second no input, return\n")
			return
		}
	}
}

func tcpClient() {
	var buf [2048]byte
	// 绑定
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":10005")
	// 连接
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	rAddr := conn.RemoteAddr()
	//for {
	// 发送
	n, _ := conn.Write([]byte{0x00, 0x04, 0x10, 0x10, 0x03, 0x01, 0x0d, 0x0a})
	//timekit.Sleep(100)
	//log.Println(2)
	//n, _ = conn.Write([]byte(" Hello server2"))
	// 接收
	n, _ = conn.Read(buf[0:])
	fmt.Println("Reply form server", rAddr.String(), bytekit.Bytes2HexArray(buf[0:n]))
	time.Sleep(time.Second * 1)
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
	b, err := io.ReadAll(r)
	log.Println(err)
	log.Println(string(b))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
