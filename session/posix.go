package session

import (
	"fmt"
	"github.com/kulinacs/cast/agent"
	"github.com/kulinacs/linenum/release"
	log "github.com/sirupsen/logrus"
	"time"
)

// Posix is POSIX shell
type Posix struct {
	agent         *agent.Shell
	kernelVersion string
	osRelease     *release.OSRelease
}

func (s *Posix) String() string {
	return fmt.Sprintf("%s - %s", s.Type(), s.agent.Addr)
}

// Agent returns the underlying agent
func (s *Posix) Agent() *agent.Shell {
	return s.agent
}

// Type returns the session type
func (s *Posix) Type() string {
	return "Posix"
}

// Enumerate gathers basic information about the system
func (s *Posix) Enumerate() {
	s.KernelVersion()
	s.OSRelease()
}

// KernelVersion returns the kernel version of the system, enumerating it if necessary
func (s *Posix) KernelVersion() string {
	log.Debug("getting kernel version")
	if s.kernelVersion == "" {
		var err error
		s.agent.Write("uname -r")
		s.kernelVersion, err = s.agent.Read(time.Millisecond * 25)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("error occurred reading the version")
		}
	}
	return s.kernelVersion
}

// OSRelease returns the parsed contents of /etc/os-release
func (s *Posix) OSRelease() *release.OSRelease {
	log.Debug("getting /etc/os-release")
	if s.osRelease == nil {
		s.agent.Write("cat /etc/os-release")
		log.Info("getting os-release")
		result, err := s.agent.ReadAll()
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("error occurred reading os release")
		}
		s.osRelease = release.ParseOSRelease(result)
	}
	return s.osRelease
}
