package conf

import "os/exec"

func Exec(command string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", command)

	return cmd.Run()
}
