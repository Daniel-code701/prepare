package main

import (
	"fmt"
	"os/exec"
)



func main() {
	//使用var定义变量比:=好 因为使用:=来定义变量 goto会报错
	var (
		cmd *exec.Cmd
		output []byte
		err error
	)

	cmd = exec.Command("/bin/bash","-c", "ls -al;echo hello")

	if output,err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}else {
		fmt.Println(string(output))
	}


}
