package main

import (
	"context"
	"fmt"
	"os/exec"
)

type result struct {
	err	   error
	output []byte
}

func main() {
	var (
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		res *result
	)

	//创建一个结果队列
	resultChan = make(chan *result, 1000)

	//执行一个cmd 让它在一个协成里执行2秒
	ctx,cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var(
			output []byte
			err error
		)

		exec.CommandContext(ctx,"bash","la -al")
		output,err = cmd.CombinedOutput()

		resultChan <- &result{
			err:err,
			output:output,
		}
	}()

	cancelFunc()
	res =<- resultChan
	fmt.Println(res.err)
}
