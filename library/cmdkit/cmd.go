package cmdkit

import (
	"bufio"
	"errors"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/spf13/cast"
	"io"
	"os/exec"
	"time"
)

type RunParams struct {
	Timeout int  `comment:"超时时间s"`
	Async   bool `comment:"异步处理返回值"`
}

// Run
// example: []string{"/bin/bash", "-c", "xxx xxx"}, []string{"/bin/sh", "-c", "xxx.sh xxx"}, []string{"xxx","-h"}
// example: []string{"cmd", "/C", "xxx xxx"},
func Run(command []string, params ...RunParams) (string, error) {
	if len(command) == 0 {
		panic(exception.New("cmd need command"))
	}
	var param RunParams
	if len(params) == 0 {
		param = RunParams{}
	} else {
		param = params[0]
	}
	//var cmdName string
	//var arg1 string
	//switch runtime.GOOS {
	//case "darwin", "linux":
	//	cmdName = "/bin/sh"
	//	arg1 = "-c"
	//case "windows":
	//	cmdName = "cmd"
	//	arg1 = "/C"
	//}
	name := command[0]
	var args []string
	if len(command) > 1 {
		args = command[1:]
	}
	cmd := exec.Command(name, args...)
	// 程序退出时Kill子进程
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if !param.Async {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return "", err
		}
		if err = cmd.Start(); err != nil {
			return "", err
		}
		if param.Timeout > 0 {
			to := make(chan map[string]any)
			go func() {
				ret0, err2 := getRet(stdout, stderr, cmd)
				to <- map[string]any{"ret": ret0, "err": err2}
			}()
			select {
			case <-time.After(time.Duration(param.Timeout) * time.Second):
				return "", errors.New("cmd timeout:" + name)
			case m := <-to:
				ret := m["ret"].(string)
				var err error
				if m["err"] != nil {
					err = m["err"].(error)
				}
				return cast.ToString(ret), err
			}
		} else {
			ret, err := getRet(stdout, stderr, cmd)
			return ret, err
		}
	} else {
		if err := cmd.Start(); err != nil {
			return "", err
		}
	}
	return "", nil
}

func getRet(stdout io.ReadCloser, stderr io.ReadCloser, cmd *exec.Cmd) (string, error) {
	ret := ""
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		ret += line
	}
	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return ret, err
	}
	if len(bytesErr) != 0 {
		return ret, errors.New(string(bytesErr))
	}
	if err = cmd.Wait(); err != nil {
		return ret, err
	}
	return ret, nil
}
