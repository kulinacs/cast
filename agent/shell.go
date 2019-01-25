package agent

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

var errReadTimeout = errors.New("timeout attempting to read shell")

// ErrShellClosed occurs when the shell agent is no longer active
var ErrShellClosed = errors.New("shell is no longer active")

// NewShell returns a new Shell with created channels for non-blocking reads and writes
func NewShell(conn io.ReadWriter, buffer int, address net.Addr) *Shell {
	newShell := &Shell{active: true, reader: bufio.NewScanner(conn), writer: conn,
		readInternal: make(chan string, buffer), Addr: address}
	go newShell.startReader()
	return newShell
}

// NewSplitShell shell returns a new shell with a seperate reader and writer
func NewSplitShell(reader io.Reader, writer io.Writer, buffer int, address net.Addr) *Shell {
	newShell := &Shell{active: true, reader: bufio.NewScanner(reader), writer: writer,
		readInternal: make(chan string, buffer), Addr: address}
	go newShell.startReader()
	return newShell
}

// Shell wraps a io.ReadWriter in a way that allows it to handle a remote shell
type Shell struct {
	active           bool
	readMutex        sync.Mutex
	writeMutex       sync.Mutex
	reader           *bufio.Scanner
	writer           io.Writer
	readInternal     chan string
	Addr             net.Addr
}

// startReader starts the process that reads from the scanner and puts it on the readInternal channel
func (s *Shell) startReader() {
	for s.reader.Scan() {
		s.readInternal <- s.reader.Text()
	}
	log.Trace("closing session")
	s.active = false
}

// write ignores the write mutex
func (s *Shell) write(val string) {
	log.WithFields(log.Fields{"msg": val}).Trace("message sent")
	fmt.Fprint(s.writer, val+"\n")
}

// Read acquires the readMutex and reads from the underlying channel, with the given timeout
func (s *Shell) Read(timeout time.Duration) (string, error) {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		log.Trace("read timeout")
		return "", errReadTimeout
	case val := <-s.readInternal:
		log.WithFields(log.Fields{"msg": val}).Trace("message received")
		if !s.active {
			return val, ErrShellClosed
		}
		return val, s.reader.Err()
	}
}

// ReadAll acquires the readMutex and reads from the underlying reader with a timeout
func (s *Shell) ReadAll() ([]string, error) {
	s.readMutex.Lock()
	defer s.readMutex.Unlock()
	result := make([]string, 0)
	timeout := time.NewTimer(25 * time.Millisecond)
ReadLoop:
	for {
		select {
		case <-timeout.C:
			log.Trace("read all timeout")
			break ReadLoop
		case val := <-s.readInternal:
			log.WithFields(log.Fields{"msg": val}).Trace("message received, resetting timeout")
			result = append(result, val)
			timeout.Reset(25 * time.Millisecond)
		}
	}
	if !s.active {
		return result, ErrShellClosed
	}
	return result, s.reader.Err()
}

// Write acquires the writeMutex and writes to the underlying writer
func (s *Shell) Write(val string) {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	s.write(val)
}

// Execute executes a command on the agent and returns the result as a string slice
func (s *Shell) Execute(val string) ([]string, error) {
	s.Write(val)
	return s.ReadAll()
}
