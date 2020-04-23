package exception

import "runtime"

type Exception struct {
	Msg  string
	File string
	Line int
}

func New(msg string) Exception {
	_, file, line, _ := runtime.Caller(1)
	return Exception{
		Msg:  msg,
		File: file,
		Line: line,
	}
}

func (th Exception) Error() string {
	return th.Msg
}
