package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	//生成Cmd
	cmd = exec.Command("D:\\cygwin64\\bin\\bash.exe", "-c", "ls -l")

	//执行命令，捕获子进程的输出
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}

	fmt.Println(string(output))
}
