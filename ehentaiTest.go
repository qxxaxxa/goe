package main

import (
	"fmt"
	"runtime"
)

func main() {

	sysType := runtime.GOOS
	if sysType == "windows" {
		fmt.Println("win")
		// windows系统
	}
	if sysType == "darwin" {
		// LINUX系统
		fmt.Println("mac")
	}

}
