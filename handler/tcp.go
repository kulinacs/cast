package handler

import (
	"fmt"
	"github.com/kulinacs/cast/agent"
	"github.com/kulinacs/cast/session"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

// TCPHandler is a TCP reverse shell handler
type TCPHandler struct {
	sessions []session.Shell
}

func (handler *TCPHandler) String() string {
	return fmt.Sprintf("%s - %d sessions", handler.Type(), len(handler.sessions))
}

// Handle listens for and creates incoming sessions
func (handler *TCPHandler) Handle(port int) {
	soc, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.WithFields(log.Fields{"port": port, "err": err}).Error("failed to start tcp handler")
		return
	}
	log.WithFields(log.Fields{"port": port}).Info("starting handler")
	defer soc.Close()
	for {
		conn, err := soc.Accept()
		if err != nil {
			log.WithFields(log.Fields{"port": port, "err": err}).Error("failed to accept incoming connection")
			continue
		}
		shellSession, err := session.UpgradeShell(agent.NewShell(conn, 20, conn.RemoteAddr()))
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("failed to upgrade session")
			continue
		}
		shellSession.Enumerate()
		log.WithFields(log.Fields{"session": shellSession}).Info("new session")
		handler.sessions = append(handler.sessions, shellSession)
	}
}

// Sessions returns the sessions slice
func (handler *TCPHandler) Sessions() []session.Shell {
	return handler.sessions
}

// Session returns the session at the specified index
func (handler *TCPHandler) Session(index int) session.Shell {
	return handler.sessions[index]
}

// Type returns the sessions type
func (handler *TCPHandler) Type() string {
	return "TCP"
}
