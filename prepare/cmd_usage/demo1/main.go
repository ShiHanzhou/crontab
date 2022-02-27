package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)
	cmd = exec.Command("D:\\cygwin64\\bin\\bash", "-c", "echo 1;echo 2")

	err = cmd.Run()

	fmt.Println(err)
}
