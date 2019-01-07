package session

import (
	"fmt"
	"github.com/kulinacs/cast/agent"
	log "github.com/sirupsen/logrus"
)

type Sh struct {
	agent         *agent.Shell
	kernelVersion string
}

func (s *Sh) String() string {
	return fmt.Sprintf("%s - %s", s.Type(), s.KernelVersion())
}

// Type returns the session type
func (s *Sh) Type() string {
	return "Linux"
}

// Enumerate gathers basic information about the system
func (s *Sh) Enumerate() {
	s.KernelVersion()
	s.OSRelease()
}

// KernelVersion returns the kernel version of the system, enumerating it if necessary
func (s *Sh) KernelVersion() string {
	if s.kernelVersion == "" {
		var err error
		s.agent.Write("uname -r")
		s.kernelVersion, err = s.agent.Read()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("error occured reading the version")
		}
	}
	return s.kernelVersion
}

func (s *Sh) OSRelease() {
	s.agent.Write("cat /etc/os-release")
	result, err := s.agent.ReadAll()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("error occured reading os release")
	}
	log.WithFields(log.Fields{"result": result}).Error("os-release")
}
