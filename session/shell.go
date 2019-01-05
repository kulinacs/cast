package session

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

// NewShell returns a new Shell with created channels for non-blocking reads and writes
func NewShell(conn io.ReadWriter, buffer int) *Shell {
	return &Shell{active: true, reader: bufio.NewReader(conn), writer: conn,
		ReadInteractive: make(chan string, buffer), WriteInteractive: make(chan string, buffer)}
}

// Shell wraps a io.ReadWriter in a way that allows it to handle a remote shell
type Shell struct {
	active bool
	interactive bool
	readMutex sync.Mutex
	writeMutex sync.Mutex
	reader *bufio.Reader
	writer io.Writer
	ReadInteractive   chan string
	WriteInteractive chan string
}

// read ignores the read mutex
func (s *Shell) read() (string, error) {
	val, err := s.reader.ReadString('\n')
	log.WithFields(log.Fields{"msg": val}).Debug("message received")
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("error occured, closing")
		s.active = false
		return "", err
	}
	return val, nil
}

// write ignores the write mutex
func (s *Shell) write(val string) {
	log.WithFields(log.Fields{"msg": val}).Debug("message sent")
	fmt.Fprint(s.writer, val)
}

// Read acquires the readMutex and reads from the underlying reader
func (s *Shell) Read() (string, error) {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	return s.read()
}

// Write acquires the writeMutex and writes to the underlying writer
func (s *Shell) Write(val string) {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	s.write(val)
}

// readInteractive reads from the underlying buffer and returns the output to the ReadInteractive channel
func (s *Shell) readInteractive() {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	for s.active && s.interactive {
		val, _ := s.read()
		s.ReadInteractive <- val
	}
}

// writeInteractive reads from the WriteInteractive channel and writes it to the underlying writer
func (s *Shell) writeInteractive() {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	for s.active && s.interactive {
		select {
		case val := <-s.WriteInteractive:
			s.write(val)
		default:
		}
	}
}

// Interactive enables the interactive read and write channels
func (s *Shell) Interactive() {
	s.interactive = true
	go s.readInteractive()
	go s.writeInteractive()
}

// Detach disables the interactive read and write channels
func (s *Shell) Detach() {
	s.interactive = false
}
