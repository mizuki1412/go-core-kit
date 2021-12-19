package test

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
	"github.com/spf13/cobra"
	"net"
	"time"
)

func init() {
	//initkit.DefFlags(rootCmd)
	//rootCmd.AddCommand(cmd.PGSqlToStructCMD("", ""))
	//rootCmd.AddCommand(cmd.MarkdownDocCMD("go-core-kit 说明文档"))
	//rootCmd.AddCommand(cmd.WebStaticServerCMD())
}

var rootCmd = &cobra.Command{
	Use: "go-core-kit",
	Run: func(cmd *cobra.Command, args []string) {
		//initkit.BindFlags(cmd)
	},
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
