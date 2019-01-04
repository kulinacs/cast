package session

import (
	"bufio"
	"fmt"
	"io"
	"testing"
)

// TestHandleReadOpen tests a single read from a pipe left open
func TestHandleReadOpen(t *testing.T) {
	testVal := "test text\n"
	pipeReader, pipeWriter := io.Pipe()
	testShell := Shell{active: true, reader: bufio.NewReader(pipeReader), Read: make(chan string, 10)}
	go func() {
		fmt.Fprint(pipeWriter, testVal)
	}()
	go testShell.handleRead()
	recvVal := <-testShell.Read
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
	if !testShell.active {
		t.Errorf("Shell incorrectly marked inactive")
	}
}

// TestHandleReadOpen tests a single read from a pipe that is closed
func TestHandleReadClose(t *testing.T) {
	testVal := "test text\n"
	pipeReader, pipeWriter := io.Pipe()
	testShell := Shell{active: true, reader: bufio.NewReader(pipeReader), Read: make(chan string, 10)}
	go func() {
		fmt.Fprint(pipeWriter, testVal)
		pipeWriter.Close()
	}()
	go testShell.handleRead()
	recvVal := <-testShell.Read
	// Receive twice to ensure we've exhausted the buffer
	<-testShell.Read
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
	if testShell.active {
		t.Errorf("Shell incorrectly left active")
	}
}

// TestHandleWriteActive test a single write from a pipe
func TestHandleWriteActive(t *testing.T) {
	testVal := "test text\n"
	pipeReader, pipeWriter := io.Pipe()
	testShell := Shell{active: true, writer: pipeWriter, Write: make(chan string, 10)}
	testShell.Write <- testVal
	go testShell.handleWrite()
	bufferedReader := bufio.NewReader(pipeReader)
	recvVal, err := bufferedReader.ReadString('\n')
	if err != nil {
		t.Errorf("an error occurred reading the pipe")
	}
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", testVal, recvVal)
	}
}

// TestHandleWriteInactive tests handleWrite returns when inactive
func TestHandleWriteInactive(t *testing.T) {
	testShell := Shell{active: false}
	testShell.handleWrite()
}
