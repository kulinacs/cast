package main

import (
	"bufio"
	"github.com/abiosoft/ishell"
	"github.com/kulinacs/cast/handler"
	//	"github.com/kulinacs/cast/session"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var handlers []handler.Handler
var shell *ishell.Shell

func sessionInteract(c *ishell.Context, handlerIndex int, sessionIndex int) {
	shell.Stop()
	defer shell.Run()
	reader := bufio.NewReader(os.Stdin)
	sessionAgent := handlers[handlerIndex].Session(sessionIndex).Agent()
	for {
		// Read the keyboad input
		input, _ := reader.ReadString('\n')
		sessionAgent.Write(input)
		if input == "background\n" {
			break
		}
		output, _ := sessionAgent.Read()
		c.Printf("%s", output)
	}
}

func getCmd() *ishell.Cmd {
	get := ishell.Cmd{
		Name: "get",
		Help: "show resources",
	}
	get.AddCmd(&ishell.Cmd{
		Name: "session",
		Help: "get session",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				for index, element := range handlers {
					for subindex, subelement := range element.Sessions() {
						c.Printf("%d - %d - %s\n", index, subindex, subelement)
					}
				}
			}
			if len(c.Args) == 2 {
				handlerIndex, _ := strconv.Atoi(c.Args[0])
				sessionIndex, _ := strconv.Atoi(c.Args[1])
				sessionInteract(c, handlerIndex, sessionIndex)
			}
		},
	})
	getHandler := ishell.Cmd{
		Name: "handler",
		Help: "get handler",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				for index, element := range handlers {
					c.Printf("%d - %s\n", index, element)
				}
			}
		},
	}
	get.AddCmd(&getHandler)
	return &get
}

func createCmd() *ishell.Cmd {
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
			handlers = append(handlers, &exampleHandler)
		},
	}
	createHandler.AddCmd(&createTCPHandler)
	create.AddCmd(&createHandler)
	return &create
}

func init() {
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	shell = ishell.New()
	shell.AddCmd(getCmd())
	shell.AddCmd(createCmd())

	shell.Run()
}
