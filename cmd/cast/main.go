package main

import (
	"fmt"
	"github.com/kulinacs/cast/handler"
)

func main() {
	handle := handler.TCPHandler{}
	go handle.Handle(1337)
	for {
		fmt.Println(len(handle.Sessions))
	}
}
