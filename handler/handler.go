package handler

// Handler to receive incoming connections to establish sessions
type Handler interface {
	Handle(port int)
}
