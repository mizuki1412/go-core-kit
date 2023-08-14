package cli

import "github.com/spf13/cobra"

type CMD struct {
	Root *cobra.Command
}

var rootCmd *CMD

func RootCMD(command *cobra.Command) {
	rootCmd = &CMD{Root: command}
	bindDefaultFlags(rootCmd.Root)
}

func AddChildCMD(command *cobra.Command) {
	if rootCmd.Root == nil {
		panic("root cmd not config")
	}
	rootCmd.Root.AddCommand(command)
}

func Execute() {
	if rootCmd.Root == nil {
		panic("root cmd not config")
	}
	if err := rootCmd.Root.Execute(); err != nil {
		panic(err.Error())
	}
}
