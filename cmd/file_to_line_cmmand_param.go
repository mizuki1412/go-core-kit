package cmd

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func File2LineCli() *cobra.Command {
	var json2CliCmd = &cobra.Command{
		Use:   "config2line",
		Short: `config file转命令行的参数形式`,
		Run: func(cmd *cobra.Command, args []string) {
			// 用 cli 中的 loadconfig
			var res []string
			for _, e := range viper.AllKeys() {
				if viper.IsSet(e) && e != "config" {
					res = append(res, fmt.Sprintf("--%s=%s", e, viper.GetString(e)))
				}
			}
			fmt.Println(strings.Join(res, " "))
		},
	}
	return json2CliCmd
}

func handleMap(parentKey string, data map[string]any) []string {
	var ret []string
	for key, val := range data {
		k := key
		if parentKey != "" {
			k = parentKey + "." + key
		}
		if v, ok := val.(map[string]any); ok {
			ret = append(ret, handleMap(k, v)...)
		} else {
			ret = append(ret, "--"+k+"="+cast.ToString(val))
		}
	}
	return ret
}
