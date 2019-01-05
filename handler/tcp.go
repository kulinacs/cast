package handler

import (
	"github.com/kulinacs/cast/session"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type TCPHandler struct {
	Sessions []*session.Shell
}

func (handler *TCPHandler) Handle(port int) {
	soc, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer soc.Close()
	log.WithFields(log.Fields{"port": port}).Error("starting handler")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := soc.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handler.Sessions = append(handler.Sessions, session.NewShell(conn))
	}
}
