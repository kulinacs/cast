package main

import (
	"bufio"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/kulinacs/cast/handler"
	"github.com/kulinacs/cast/session"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var handlers []handler.Handler
var sessions []session.Shell
var shell *ishell.Shell

func appendSession(sess session.Shell) {
	sessions = append(sessions, sess)
}

func sessionInteract(c *ishell.Context, sessionIndex int) {
	shell.Stop()
	defer shell.Run()
	reader := bufio.NewReader(os.Stdin)
	selectedSession := sessions[sessionIndex]
	for {
		// Read the keyboad input
		input, _ := reader.ReadString('\n')
		if input == "background\n" {
			break
		}
		output, err := selectedSession.Execute(input)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("closing session")
			break
		}
		for _, element := range output {
			fmt.Printf("%s\n", element)
		}
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
				for index, element := range sessions {
					c.Printf("%d - %s\n", index, element)
				}
			}
			if len(c.Args) == 1 {
				sessionIndex, _ := strconv.Atoi(c.Args[0])
				sessionInteract(c, sessionIndex)
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
			exampleHandler := handler.TCPHandler{SessionCallback: appendSession}
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	shell = ishell.New()
	shell.AddCmd(getCmd())
	shell.AddCmd(createCmd())

	shell.Run()
}
