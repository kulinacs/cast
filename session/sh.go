package session

import (
	"fmt"
	"github.com/kulinacs/cast/agent"
	"github.com/kulinacs/linenum/release"
	log "github.com/sirupsen/logrus"
	"time"
)

// Sh is a Linux POSIX shell
type Sh struct {
	agent         *agent.Shell
	kernelVersion string
	osRelease     *release.OSRelease
}

func (s *Sh) String() string {
	return fmt.Sprintf("%s - %s", s.Type(), s.agent.Addr)
}

// Agent returns the underlying agent
func (s *Sh) Agent() *agent.Shell {
	return s.agent
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
func (s *Sh) OSRelease() *release.OSRelease {
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
