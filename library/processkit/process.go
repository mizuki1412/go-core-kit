package processkit

import (
	"fmt"
	"os"
)

func Cmd(program string, argv ...string) (bool, string) {
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	list := []string{program}
	for _, s := range argv {
		list = append(list, s)
	}
	process, err := os.StartProcess(program, list, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!", err) //
		os.Exit(1)
	}
	state, err := process.Wait()
	if err != nil {
		return false, err.Error()
	}
	return state.Success(), state.String()
}
