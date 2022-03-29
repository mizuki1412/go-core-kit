package regexkit

import "regexp"

func IsPhone(val string) bool {
	ok, _ := regexp.Match("^1[34578]\\d{9}$", []byte(val))
	return ok
}

func isIP(val string) bool {
	ok, _ := regexp.Match("^((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}$", []byte(val))
	return ok
}
