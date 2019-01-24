package session

import (
	"fmt"
	"github.com/kulinacs/cast/agent"
	log "github.com/sirupsen/logrus"
)

// Posix is POSIX shell
type Posix struct {
	agent *agent.Shell
	os    string
}

func (s *Posix) String() string {
	return fmt.Sprintf("%s - %s", s.Type(), s.agent.Addr)
}

// Execute runs a command on the underlying agent
func (s *Posix) Execute(command string) ([]string, error) {
	return s.agent.Execute(command)
}

// Type returns the session type
func (s *Posix) Type() string {
	return "Posix"
}

// OS identifies the underlying Operating System
func (s *Posix) OS() string {
	log.Trace("enumerating OS")
	if s.os == "" {
		outputLines, err := s.agent.Execute("uname -s")
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("error occurred identifying the operation system")
			s.os = "unknown"
		} else if len(outputLines) != 1 {
			log.WithFields(log.Fields{"output": outputLines}).Error("unknown response for OS received")
			s.os = "unknown"
		} else {
			s.os = outputLines[0]
		}
	}
	return s.os
}
