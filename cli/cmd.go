package cli

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cobra"
)

type CMD struct {
	Root *cobra.Command
}

var rootCmd *CMD

// 包装加载 config，从--config 中读取配置文件地址
func decoCmd(cmd *cobra.Command) {
	run := cmd.Run
	cmd.Run = func(cmd *cobra.Command, args []string) {
		loadConfig()
		logkit.Init()
		run(cmd, args)
	}
}

func RootCMD(command *cobra.Command) {
	decoCmd(command)
	rootCmd = &CMD{Root: command}
	bindDefaultFlags(rootCmd.Root)
}

func AddChildCMD(command *cobra.Command) {
	if rootCmd.Root == nil {
		panic("root cmd not config")
	}
	decoCmd(command)
	rootCmd.Root.AddCommand(command)
	bind(command)
}

func Execute() {
	if rootCmd.Root == nil {
		panic("root cmd not config")
	}
	if err := rootCmd.Root.Execute(); err != nil {
		panic(err.Error())
	}
}
