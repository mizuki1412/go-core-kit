package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/netkit"
	"github.com/panjf2000/gnet"
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
			initkit.BindFlags(cmd)
			if configkit.GetStringD("port") == "" {
				logkit.Fatal("port参数缺失")
			}
			server := netkit.NetServer{
				Port: cast.ToInt32(configkit.GetStringD("port")),
				OnConnect: func(c gnet.Conn) (out []byte, action gnet.Action) {
					log.Println(" OnConnect ： ", c.RemoteAddr())
					return
				},
				OnMessage: func(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
					log.Println("收到：" + string(frame))
					return
				},
				OnClose: func(c gnet.Conn, err error) (action gnet.Action) {
					log.Println(" close ： ")
					return
				},
			}
			server.Run()
		},
	}
	cmd.Flags().StringP("port", "", "", "端口")
	return cmd
}

func tcpServer() {
	service := ":5000"
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
