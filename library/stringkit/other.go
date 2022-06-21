package stringkit

import "unsafe"

// UnquoteIfQuoted If the amount is quoted, strip the quotes
func UnquoteIfQuoted(bytes []byte) string {
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes)
}

// StringToBytes converts string to byte slice without a memory allocation.
// https://github.com/gin-gonic/gin/blob/master/internal/bytesconv/bytesconv.go
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
