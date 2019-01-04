package session

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

// NewShell initializes a new active
func NewShell(conn io.ReadWriter) *Shell {
	shell := &Shell{active: true, reader: bufio.NewReader(conn), writer: conn, Read: make(chan string, 10), Write: make(chan string, 10)}
	go shell.handleRead()
	go shell.handleWrite()
	return shell
}

// Shell wraps a io.ReadWriter in a way that allows it to handle a remote shell
type Shell struct {
	active bool
	reader *bufio.Reader
	writer io.Writer
	Read   chan string
	Write  chan string
}

func (s *Shell) handleRead() {
	for {
		val, err := s.reader.ReadString('\n')
		log.WithFields(log.Fields{"msg": val}).Debug("message received")
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("error occured, closing")
			s.active = false
			close(s.Read)
			return
		}
		s.Read <- val
	}
}

func (s *Shell) handleWrite() {
	for {
		select {
		case val := <-s.Write:
			log.WithFields(log.Fields{"msg": val}).Debug("writing message")
			fmt.Fprint(s.writer, val)
		default:
			if !s.active {
				return
			}
		}
	}
}
