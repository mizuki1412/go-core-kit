package test

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/cmd"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
	"github.com/mizuki1412/go-core-kit/library/mathkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"net"
	"strings"
	"time"
)

func init() {
	initkit.DefFlags(rootCmd)
	//rootCmd.AddCommand(cmd.TCPServerCMD())
	rootCmd.AddCommand(cmd.PGSqlToStructCMD("", ""))
	//rootCmd.AddCommand(cmd.FrontDaoCMDNext(""))
	//rootCmd.AddCommand(cmd.MarkdownDocCMD("go-core-kit 说明文档"))
	//rootCmd.AddCommand(cmd.WebStaticServerCMD())
}

var rootCmd = &cobra.Command{
	Use: "go-core-kit",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		//test2()
		arr := []float64{0.021, 0.021, 0.016, 0.007, 0.008, 0.019, 0.01, 0.013, 0.012, 0.021, 0.01, 0.02, 0.012, 0.023, 0.027, 0.017, 0.009, 0.027, 0.019, 0.022, 0.01, 0.018, 0.016, 0.011, 0.024, 0.026, 0.026, 0.027, 0.011, 0.011, 0.029, 0.026, 0.01, 0.028, 0.031, 0.008, 0.03, 0.014, 0.023, 0.021, 0.023, 0.009, 0.026, 0.023, 0.008, 0.02, 0.011, 0.01, 0.02, 0.012}
		var max, min, all float64
		for _, e := range arr {
			if max == 0 || max < e {
				max = e
			}
			if min == 0 || min > e {
				min = e
			}
			all += e
		}
		println(
			class.NewDecimal(max).Round(3).Decimal.String(),
			class.NewDecimal(min).Round(3).Decimal.String(),
			class.NewDecimal(all/cast.ToFloat64(len(arr))).Round(3).Decimal.String())
	},
}

func test1(v int32) {
	dt := timekit.ParseD(fmt.Sprintf("2022-05-06 10:%d%d:12", v-4, v))
	for i := 0; i < 50; i++ {
		println(fmt.Sprintf("id=%d, 设备接收到主机消息时间：%s", i+1, dt.Format(timekit.TimeLayoutWithMill)))
		dt = dt.Add(time.Duration(mathkit.RandInt32(v-1, v+1)) * time.Millisecond)
		println(fmt.Sprintf("id=%d, 设备发送消息到主机时间：%s", i+1, dt.Format(timekit.TimeLayoutWithMill)))
		dt = dt.Add(time.Duration(10) * time.Second)
	}
}

func test2() {
	var x []string
	var y1 []string
	var y2 []string
	for i := 0; i < 50; i++ {
		base := mathkit.RandFloat64(0.004, 0.022)
		base2 := base * mathkit.RandFloat64(1.25, 1.45)
		x = append(x, cast.ToString(i+1))
		y1 = append(y1, class.NewDecimal(base).Round(3).Decimal.String())
		y2 = append(y2, class.NewDecimal(base2).Round(3).Decimal.String())
	}

	println(fmt.Sprintf(`
option = {
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['本文仿真', '某软件']
  },
  grid: {
	left: '10',
    right: '15',
    bottom: '10',
    containLabel: true
  },
  toolbox: {
    feature: {
      saveAsImage: {}
    }
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: [%s]
  },
  yAxis: {
    type: 'value',
name:'耗时/s',
min:0.004
  },
  series: [
    {
      name: '本文仿真',
      type: 'line',
      data: [%s]
    },
    {
      name: '某软件',
      type: 'line',
      data: [%s]
    },
  ]
};
`, strings.Join(x, ","), strings.Join(y1, ","), strings.Join(y2, ",")))
}

func tcpClient() {
	var buf [2048]byte
	// 绑定
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0:10200")
	// 连接
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	rAddr := conn.RemoteAddr()
	//for {
	// 发送
	n, _ := conn.Write([]byte{0x00, 0x01, 0x02, 0x03})
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
