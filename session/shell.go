package session

import (
	"errors"
	"github.com/kulinacs/cast/agent"
	log "github.com/sirupsen/logrus"
	"strings"
)

var errUnknownShell = errors.New("unknown shell type")

// Shell type session, to be implemented by for example, sh or Powershell
type Shell interface {
	Type() string
	Execute(command string) ([]string, error)
	OS() string
}

// UpgradeShell takes and incoming shell agent and upgrades it to a shell session
func UpgradeShell(s *agent.Shell) (Shell, error) {
	shellLines, err := s.Execute("echo $SHELL")
	if err != nil {
		return nil, err
	}
	if len(shellLines) == 1 && strings.Contains(shellLines[0], "sh") {
		return &Posix{agent: s}, nil
	}
	log.Error("unknown shell type")
	return nil, errUnknownShell
}
