package conf

import (
	"fmt"
	"log"
	"os/exec"
)

func Exec(command string) error {
	log.Printf("begin command: [%s]", command)
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	fmt.Println(string(out))
	log.Printf("finish command: [%s]", command)

	return err
}
