package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/netkit"
	"github.com/panjf2000/gnet/v2"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"log"
	"net"
)

func TCPServerCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tcp-server",
		Short: "本地tcp服务器",
		Run: func(cmd *cobra.Command, args []string) {
			server := &netkit.NetServer{
				Port: cast.ToInt32(configkit.GetString("port")),
				//OnConnect: func(c gnet.Conn) (out []byte, action gnet.Action) {
				//	log.Println("OnConnect： ", c.RemoteAddr())
				//	return
				//},
				TrafficHandler: func(c gnet.Conn) {
					buf, _ := c.Next(-1)
					log.Println("recv：" + bytekit.Bytes2HexArray(buf))
					return
				},
				//OnClose: func(c gnet.Conn, err error) (action gnet.Action) {
				//	log.Println("OnClose： ", c.RemoteAddr())
				//	return
				//},
			}
			server.Run()
		},
	}
	cmd.Flags().String("port", "", "端口")
	_ = cmd.MarkFlagRequired("port")
	return cmd
}

func tcpServer(port string) {
	service := ":" + port
	// 绑定
	tcpAddr, _ := net.ResolveTCPAddr("tcp", service)
	// 监听
	listener, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		// 接受
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// 创建 Goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	// 逆序调用 Close() 保证连接能正常关闭
	defer conn.Close()
	var buf [512]byte
	for {
		// 接收数据
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		_, err2 := conn.Write([]byte("Welcome client"))
		if err2 != nil {
			return
		}
	}
}
