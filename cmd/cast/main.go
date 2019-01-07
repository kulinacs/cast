package main

import (
	//	"fmt"
	"github.com/kulinacs/cast/handler"
)

func main() {
	handle := handler.TCPHandler{}
	handle.Handle(1337)
}
