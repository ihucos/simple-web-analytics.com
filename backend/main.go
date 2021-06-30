package main

import (
	"github.com/ihucos/counter.dev/lib"
	_ "github.com/ihucos/counter.dev/endpoints"
	"syscall"
	"fmt"
)




func main() {

	// HOTFIX
	var rLimit syscall.Rlimit
	rLimit.Max = 100307
	rLimit.Cur = 100307
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}

	app := lib.NewApp()
	app.ConnectEndpoints()
	app.Logger.Println("Start")
	app.Serve()
}
