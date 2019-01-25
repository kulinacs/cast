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
	recvVal, err := testShell.Read(10 * time.Millisecond)
	assert.Equal(t, testVal, recvVal, "read value incorrect")
	assert.Equal(t, ErrShellClosed, err, "shell not closed")
}

// TestReadError tests a failed interactive read
func TestReadError(t *testing.T) {
	var testBuffer bytes.Buffer
	testShell := NewShell(&testBuffer, 10, mockAddr())
	recvVal, err := testShell.Read(time.Nanosecond)
	assert.Equal(t, "", recvVal, "read value incorrect")
	assert.Equal(t, errReadTimeout, err, "error incorrect")
}

// TestReadAll tests a multiline non-interactive read
func TestReadAll(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	fmt.Fprintf(&testBuffer, strings.Repeat(testVal+"\n", 3))
	recvVal, err := testShell.ReadAll()
	assert.Equal(t, ErrShellClosed, err, "shell not closed")
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

// TestExecute tests a single execution
func TestExecute(t *testing.T) {
	var testBuffer bytes.Buffer
	testVal := "test text"
	testShell := NewShell(&testBuffer, 10, mockAddr())
	recvVals, err := testShell.Execute(testVal)
	assert.Equal(t, ErrShellClosed, err, "shell not closed")
	assert.Equal(t, testVal, recvVals[0], "execute value incorrect")
}
