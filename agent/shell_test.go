package agent

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// TestNewShell tests creating a new shell
func TestNewShell(t *testing.T) {
	testShell := NewShell(bytes.NewBuffer(nil), 10)
	if testShell.active != true {
		t.Errorf("Test shell inactive")
	}
}

// TestRead tests a single non-interactive read
func TestRead(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10)
	fmt.Fprintf(&testBuffer, testVal+"\n")
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
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10)
	fmt.Fprintf(&testBuffer, strings.Repeat(testVal+"\n", 3))
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
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10)
	testShell.Write(testVal)
	recvVal := testBuffer.String()
	if recvVal != testVal+"\n" {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
}

// TestHandleReadInteractive tests a single interactive read
func TestHandleReadInteractive(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10)
	fmt.Fprintf(&testBuffer, testVal+"\n")
	testShell.Interactive()
	recvVal := <-testShell.ReadInteractive
	if recvVal != testVal {
		t.Errorf("Received value was incorrect, got: %s, want: %s", recvVal, testVal)
	}
}

// TestHandleWriteInteractive tests a single write from a pipe
func TestHandleWriteInteractive(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10)
	testShell.WriteInteractive <- testVal
	testShell.Interactive()
	recvVal := ""
	for recvVal == "" {
		recvVal = testBuffer.String()
	}
	if recvVal != testVal+"\n" {
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
