package exception

import (
	"fmt"
	"runtime"
	"strconv"
)

type Exception struct {
	Msg   string
	File  string
	Line  int
	Stack []string
}

func New(msg string, skip1 ...int) Exception {
	var skip = 1
	if skip1 != nil && len(skip1) > 0 {
		skip = skip1[0]
	}
	_, file, line, _ := runtime.Caller(skip)
	stack := make([]string, 0, 4)
	stack = append(stack, getStackInfo(file, line))
	for i := 1; i < 4; i++ {
		_, file1, line1, ok := runtime.Caller(skip + i)
		if ok {
			stack = append(stack, getStackInfo(file1, line1))
		} else {
			stack = append(stack, "")
		}
	}
	return Exception{
		Msg:   msg,
		File:  file,
		Line:  line,
		Stack: stack,
	}
}

func getStackInfo(file string, line int) string {
	return file + ":" + strconv.Itoa(line)
}

func (th Exception) Error() string {
	return fmt.Sprintf(`exception: %s
Exception StackTrace: 
%s
%s
%s`, th.Msg, th.Stack[0], th.Stack[1], th.Stack[2])
}
