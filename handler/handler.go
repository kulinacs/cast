package handler

import (
	"github.com/kulinacs/cast/session"
)

// Handler to receive incoming connections to establish sessions
type Handler interface {
	Handle(port int)
	Sessions() []session.Shell
	Session(index int) session.Shell
}
