package agent

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
	"time"
)

func mockAddr() *net.IPAddr {
	testAddr, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	return testAddr
}

// TestNewShell tests creating a new shell
func TestNewShell(t *testing.T) {
	testShell := NewShell(bytes.NewBuffer(nil), 10, mockAddr())
	assert.Equal(t, true, testShell.active, "new shell not active")
}

// TestNewSplitShell tests creating a new shell
func TestNewSplitShell(t *testing.T) {
	testShell := NewSplitShell(bytes.NewBuffer(nil), bytes.NewBuffer(nil), 10, mockAddr())
	assert.Equal(t, true, testShell.active, "new shell not active")
}

// TestRead tests a single non-interactive read
func TestRead(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	fmt.Fprintf(&testBuffer, testVal+"\n")
	recvVal, err := testShell.Read()
	assert.Equal(t, testVal, recvVal, "read value incorrect")
	assert.Nil(t, err)
}

// TestReadAll tests a multiline non-interactive read
func TestReadAll(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	fmt.Fprintf(&testBuffer, strings.Repeat(testVal+"\n", 3))
	recvVal, err := testShell.ReadAll()
	assert.Nil(t, err)
	for _, val := range recvVal {
		assert.Equal(t, testVal, val, "read all value incorrect")
	}
}

// TestWrite tests a single non-interactive write
func TestWrite(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	testShell.Write(testVal)
	recvVal := testBuffer.String()
	assert.Equal(t, testVal+"\n", recvVal, "write value incorrect")
}

// TestHandleReadInteractive tests a single interactive read
func TestHandleReadInteractive(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	fmt.Fprintf(&testBuffer, testVal+"\n")
	testShell.Interactive()
	recvVal := <-testShell.ReadInteractive
	assert.Equal(t, testVal, recvVal, "interactive read value incorrect")
}

// TestHandleWriteInteractive tests a single write from a pipe
func TestHandleWriteInteractive(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	testShell.Interactive()
	testShell.WriteInteractive <- testVal
	// Wait for the message to propagate
	time.Sleep(10 * time.Millisecond)
	recvVal := testBuffer.String()
	assert.Equal(t, testVal+"\n", recvVal, "interactive write value incorrect")
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
