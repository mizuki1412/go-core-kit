package filekit

import "strings"

// TransferName 转化文件中的 /
func TransferName(name string) string {
	return strings.ReplaceAll(name, "/", "-")
}
