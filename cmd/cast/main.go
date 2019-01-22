package main

import (
	"github.com/abiosoft/ishell"
	"github.com/kulinacs/cast/handler"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func init() {
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	shell := ishell.New()
	get := ishell.Cmd{
		Name: "get",
		Help: "show resources",
	}
	get.AddCmd(&ishell.Cmd{
		Name: "sessions",
		Help: "get sessions",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})
	shell.AddCmd(&get)
	create := ishell.Cmd{
		Name: "create",
		Help: "create resources",
	}
	createHandler := ishell.Cmd{
		Name: "handler",
		Help: "create handler",
	}
	createTCPHandler := ishell.Cmd{
		Name: "tcp",
		Help: "create tcp handler",
		Func: func(c *ishell.Context) {
			exampleHandler := handler.TCPHandler{}
			port, _ := strconv.Atoi(c.Args[0])
			go exampleHandler.Handle(port)
		},
	}
	createHandler.AddCmd(&createTCPHandler)
	create.AddCmd(&createHandler)
	shell.AddCmd(&create)

	shell.Run()
}
