package agent

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

// NewShell returns a new Shell with created channels for non-blocking reads and writes
func NewShell(conn io.ReadWriter, buffer int) *Shell {
	newShell := &Shell{active: true, reader: bufio.NewScanner(conn), writer: conn,
		readInternal: make(chan string, buffer), ReadInteractive: make(chan string, buffer),
		WriteInteractive: make(chan string, buffer)}
	go newShell.startReader()
	return newShell
}

// Shell wraps a io.ReadWriter in a way that allows it to handle a remote shell
type Shell struct {
	active           bool
	interactive      bool
	readMutex        sync.Mutex
	writeMutex       sync.Mutex
	reader           *bufio.Scanner
	writer           io.Writer
	readInternal     chan string
	ReadInteractive  chan string
	WriteInteractive chan string
}

// startReader starts the process that reads from the scanner and puts it on the readInternal channel
func (s *Shell) startReader() {
	for s.reader.Scan() {
		s.readInternal <- s.reader.Text()
	}
}

// read ignores the read mutex
func (s *Shell) read() (string, error) {
	val := <-s.readInternal
	log.WithFields(log.Fields{"msg": val}).Debug("message received")
	return val, s.reader.Err()
}

// write ignores the write mutex
func (s *Shell) write(val string) {
	log.WithFields(log.Fields{"msg": val}).Debug("message sent")
	fmt.Fprint(s.writer, val+"\n")
}

// Read acquires the readMutex and reads from the underlying reader
func (s *Shell) Read() (string, error) {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	return s.read()
}

// ReadAll acquires the readMutex and reads from the underlying reader with a timeout
func (s *Shell) ReadAll() ([]string, error) {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	result := make([]string, 0)
	timeout := time.NewTimer(50 * time.Millisecond)
ReadLoop:
	for {
		select {
		case <-timeout.C:
			break ReadLoop
		case val := <-s.readInternal:
			result = append(result, val)
			timeout.Reset(50 * time.Millisecond)
		default:
		}
	}
	return result, s.reader.Err()
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
