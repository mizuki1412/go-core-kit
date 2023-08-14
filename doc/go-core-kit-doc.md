[TOC]

# cli

框架的CMD入口配置，包含一些常量和可配置项。

基于 [cobra](https://cobra.dev/) 和 viper。

对 cobra 简单封装，能够在 main.go 中直接设置 rootCMD 和 childCMD 。

demo:

```go
cli.RootCMD(&cobra.Command{
  Use: "main",
  Run: func(cmd *cobra.Command, args []string) {
    restkit.AddActions(user.All()...)
    restkit.AddActions(download.Init)
    _ = restkit.Run()
  },
})
cli.AddChildCMD(&cobra.Command{
  Use: "test",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("123")
  },
})
cli.AddChildCMD(cmd.TCPServerCMD())
cli.Execute()

// 额外自定义参数
cmd.Flags().String("port", "", "端口")
// 设置必填
_ = cmd.MarkFlagRequired("port")
```

## 配置文件

默认的配置文件为当前目录的`./config.yaml`。

有个全局参数 `-c` 或 `--config` 可以指定配置文件的具体路径。

命令行的参数将覆盖配置文件中相同的参数。

## cmd 例子

`/cmd`下包括了一些场景下使用的工具：

- FrontDaoCMDNext：将 swagger 接口导出成前端 dao 文件。
- File2LineCli: 配置文件转命令行文字
- MarkdownDocCMD：markdown 文件导出
- MQTTTestCMD：mqtt demo
- PGSqlToStructCMD：通过 sql 生成 model
- TCPServerCMD：tcp server
- WebStaticServerCMD：静态文件服务器

# class 封装类

# library

## jsonkit

封装 sonic：https://github.com/bytedance/sonic

## httpkit

http client

## cmdkit

调用系统 cmd。

参考：https://colobu.com/2020/12/27/go-with-os-exec/

## concurrentkit

异步等待。

## framekit

应用于数据流帧的拆包粘包处理。

## ftpkit

ftp 相关的封装

## inikit

note: https://ini.unknwon.cn/docs/intro/getting_started

```go
cfg, err := ini.Load(
    []byte("raw data"), // 原始数据
    "filename",         // 文件路径
    io.NopCloser(bytes.NewReader([]byte("some other data"))),
)

// 典型读取操作，默认分区可以使用空字符串表示
fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

// 试一试自动类型转换
fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

// 差不多了，修改某个值然后进行保存
cfg.Section("").Key("app_mode").SetValue("production")
cfg.SaveTo("my.ini.local")
```

## ipkit

ip 的处理





# iot

针对 IoT 相关的处理库。

暂不更新。
