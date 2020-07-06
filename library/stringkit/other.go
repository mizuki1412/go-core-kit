package stringkit

// If the amount is quoted, strip the quotes
func UnquoteIfQuoted(bytes []byte) string {
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes)
}
