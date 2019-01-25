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
	user  string
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
		if len(outputLines) != 1 {
			log.WithFields(log.Fields{"output": outputLines, "err": err}).Error("unknown response for OS received")
			s.os = "unknown"
		} else {
			s.os = outputLines[0]
		}
	}
	return s.os
}

// User identifies the current user
func (s *Posix) User() string {
	log.Trace("enumerating user")
	outputLines, err := s.agent.Execute("id")
	if len(outputLines) != 1 {
		log.WithFields(log.Fields{"output": outputLines, "err": err}).Error("unknown response for user received")
		s.user = "unknown"
	} else {
		s.user = outputLines[0]
	}
	if s.user == "" {
	}
	return s.user
}
