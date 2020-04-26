package exception

import "runtime"

type Exception struct {
	Msg  string
	File string
	Line int
}

func New(msg string, skip1 ...int) Exception {
	var skip = 1
	if skip1 != nil && len(skip1) > 0 {
		skip = skip1[0]
	}
	_, file, line, _ := runtime.Caller(skip)
	return Exception{
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func (th Exception) Error() string {
	return th.Msg
}
