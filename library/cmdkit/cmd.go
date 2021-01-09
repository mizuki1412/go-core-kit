package cmdkit

import (
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"os/exec"
	"path/filepath"
)

func LinuxCmd(name string, args ...string) error {
	command := exec.Command(name, args...)
	return command.Run()
}

func WinCmd(arg ...string) error {
	args := make([]string, 0, 3)
	args = append(args, "/C")
	args = append(args, arg...)
	command := &exec.Cmd{
		Path: "cmd",
		Args: args,
	}
	if filepath.Base("cmd") == "cmd" {
		if lp, err := exec.LookPath("cmd"); err != nil {
			logkit.Error("filePathErr")
		} else {
			command.Path = lp
		}
	}
	return command.Run()
}
