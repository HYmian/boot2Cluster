package conf

import (
	"os/exec"
	"strings"
)

func Exec(command string) error {
	cs := strings.Split(command, " ")

	var cmd *exec.Cmd
	if len(cs) > 1 {
		cmd = exec.Command(cs[0], cs[1:]...)
	} else {
		cmd = exec.Command(cs[0])
	}
	return cmd.Run()
}
