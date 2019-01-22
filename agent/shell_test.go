package agent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// shellFixture returns a pipe to write to for the shell to read, a pipe to read from to get shell output, and the shell attached
func shellFixture() (*io.PipeWriter, *io.PipeReader, *Shell) {
	readPipeReader, readPipeWriter := io.Pipe()
	writePipeReader, writePipeWriter := io.Pipe()
	testShell := &Shell{active: true, reader: bufio.NewScanner(readPipeReader),
		readInternal:    make(chan string, 10),
		ReadInteractive: make(chan string, 10),
		writer:          writePipeWriter, WriteInteractive: make(chan string, 10)}
	return readPipeWriter, writePipeReader, testShell
}

// TestNewShell tests creating a new shell
func TestNewShell(t *testing.T) {
	testShell := NewShell(bytes.NewBuffer(nil), 10)
	if testShell.active != true {
		t.Errorf("Test shell inactive")
	}
}

// TestRead tests a single non-interactive read
func TestRead(t *testing.T) {
	testVal := "test text"
	testShell := &Shell{active: true, reader: bufio.NewScanner(strings.NewReader(testVal + "\n")),
		readInternal: make(chan string, 10)}
	go testShell.startReader()
	recvVal, _ := testShell.Read()
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
	if !testShell.active {
		t.Errorf("Shell incorrectly marked inactive")
	}
}

// TestReadAll tests a multiline non-interactive read
func TestReadAll(t *testing.T) {
	testVal := "test text"
	testShell := &Shell{active: true, reader: bufio.NewScanner(strings.NewReader(strings.Repeat(testVal+"\n", 5))),
		readInternal: make(chan string, 10)}
	go testShell.startReader()
	recvVal, err := testShell.ReadAll()
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	for _, val := range recvVal {
		if val != testVal {
			t.Errorf("Received value was incorrect, got: %s, want: %s", val, testVal)
		}
	}
}

// TestWrite tests a single non-interactive write
func TestWrite(t *testing.T) {
	testVal := "test text~"
	_, pipeReader, testShell := shellFixture()
	go func() {
		testShell.Write(testVal)
	}()
	bufferedReader := bufio.NewReader(pipeReader)
	recvVal, err := bufferedReader.ReadString('~')
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
}

// TestHandleReadInteractive tests a single interactive read
func TestHandleReadInteractive(t *testing.T) {
	testVal := "test text"
	pipeWriter, _, testShell := shellFixture()
	go testShell.startReader()
	go func() {
		fmt.Fprint(pipeWriter, testVal+"\n")
	}()
	testShell.interactive = true
	go testShell.readInteractive()
	recvVal := <-testShell.ReadInteractive
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
}

// TestHandleWriteInteractive tests a single write from a pipe
func TestHandleWriteInteractive(t *testing.T) {
	testVal := "test text~"
	_, pipeReader, testShell := shellFixture()
	testShell.WriteInteractive <- testVal
	// Go interactive to test writeInteractive
	testShell.interactive = true
	go testShell.writeInteractive()
	bufferedReader := bufio.NewReader(pipeReader)
	recvVal, err := bufferedReader.ReadString('~')
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
}

// TestDetach tests that the interactive channels detach
func TestDetach(t *testing.T) {
	testShell := Shell{active: true}
	testShell.Interactive()
	testShell.Detach()
	// If we can lock the mutex, the session has detached
	testShell.readMutex.Lock()
	testShell.writeMutex.Lock()
}
