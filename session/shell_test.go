package session

import (
	"bufio"
	"fmt"
	"io"
	"testing"
)

// shellFixture returns a pipe to write to for the shell to read, a pipe to read from to get shell output, and the shell attached
func shellFixture() (*io.PipeWriter, *io.PipeReader, *Shell){
	readPipeReader, readPipeWriter := io.Pipe()
	writePipeReader, writePipeWriter := io.Pipe()
	testShell := &Shell{active: true, reader: bufio.NewReader(readPipeReader), ReadInteractive: make(chan string, 10),
		writer: writePipeWriter, WriteInteractive: make(chan string, 10)}
	return readPipeWriter, writePipeReader, testShell
}

// TestRead tests a single non-interactive read
func TestRead(t *testing.T) {
	testVal := "test text\n"
	pipeWriter, _, testShell := shellFixture()
	go func() {
		fmt.Fprint(pipeWriter, testVal)
	}()
	recvVal, _ := testShell.Read()
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
	if !testShell.active {
		t.Errorf("Shell incorrectly marked inactive")
	}
}

// TestReadError test a single non-interactive read with an EOF
func TestReadError(t *testing.T) {
	testVal := "test text\n"
	pipeWriter, _, testShell := shellFixture()
	go func() {
		fmt.Fprint(pipeWriter, testVal)
		pipeWriter.Close()
	}()
	testShell.Read()
	_, err := testShell.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF")
	}
	if testShell.active {
		t.Errorf("Shell incorrectly left active")
	}
}


// TestWrite tests a single non-interactive write
func TestWrite(t *testing.T) {
	testVal := "test text\n"
	_, pipeReader, testShell := shellFixture()
	go func() {
		testShell.Write(testVal)
	}()
	bufferedReader := bufio.NewReader(pipeReader)
	recvVal, err := bufferedReader.ReadString('\n')
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
}

// TestHandleReadInteractive tests a single interactive read
func TestHandleReadInteractive(t *testing.T) {
	testVal := "test text\n"
	pipeWriter, _, testShell := shellFixture()
	go func() {
		fmt.Fprint(pipeWriter, testVal)
	}()
	testShell.interactive = true
	go testShell.readInteractive()
	recvVal := <-testShell.ReadInteractive
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
}

// TestHandleWriteInteractive tests a single write from a pipe
func TestHandleWriteInteractive(t *testing.T) {
	testVal := "test text\n"
	_, pipeReader, testShell := shellFixture()
	testShell.WriteInteractive <- testVal
	// Go interactive to test writeInteractive
	testShell.interactive = true
	go testShell.writeInteractive()
	bufferedReader := bufio.NewReader(pipeReader)
	recvVal, err := bufferedReader.ReadString('\n')
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
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
