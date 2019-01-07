package session

import (
	"errors"
	"github.com/kulinacs/cast/agent"
	log "github.com/sirupsen/logrus"
)

var errUnknownShell = errors.New("unknown shell type")

type Shell interface {
	Enumerate()
	Type() string
}

func UpgradeShell(s *agent.Shell) (Shell, error) {
	s.Write("uname -s")
	os, err := s.Read()
	if err != nil {
		return nil, err
	}
	if os == "Linux" {
		return &Sh{agent: s}, nil
	} else {
		log.WithFields(log.Fields{"os": os}).Error("unknown shell type")
		return nil, errUnknownShell
	}
}
