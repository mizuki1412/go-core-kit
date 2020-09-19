package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(json2CliCmd)
}

var json2CliCmd = &cobra.Command{
	Use:  "json2cli",
	Long: `json config file转命令行的参数形式`,
	Run: func(cmd *cobra.Command, args []string) {
		json, _ := filekit.ReadString("config.deploy.json")
		config := jsonkit.ParseMap(json)
		ret := handleMap("", config)
		fmt.Println(strings.Join(ret, "\n"))
	},
}

func handleMap(parentKey string, data map[string]interface{}) []string {
	var ret []string
	for key, val := range data {
		k := key
		if parentKey != "" {
			k = parentKey + "." + key
		}
		if v, ok := val.(map[string]interface{}); ok {
			ret = append(ret, handleMap(k, v)...)
		} else {
			ret = append(ret, "--"+k+"="+cast.ToString(val))
		}
	}
	return ret
}
