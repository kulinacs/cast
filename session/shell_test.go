package session

import (
	"bytes"
	"github.com/kulinacs/cast/agent"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
)

func mockShell(output string) (*agent.Shell, *bytes.Buffer) {
	testAddr, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	var inputBuffer bytes.Buffer
	testShell := agent.NewSplitShell(bytes.NewBufferString(output), &inputBuffer, 10, testAddr)
	return testShell, &inputBuffer
}

var upgradetests = []struct {
	in  string
	out string
}{
	{"/bin/sh", "Posix"},
}

func TestUpgradeShellValid(t *testing.T) {
	for _, tt := range upgradetests {
		t.Run(tt.in, func(t *testing.T) {
			testShell, readBuffer := mockShell(tt.in)
			testSession, err := UpgradeShell(testShell)
			assert.Nil(t, err)
			assert.Equal(t, "echo $SHELL\n", readBuffer.String(), "correct command not called")
			assert.Equal(t, tt.out, testSession.Type(), "shell not correctly identified")
		})
	}
}

func TestUpgradeShellInvalid(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testShell, readBuffer := mockShell("failure\n")
	_, err := UpgradeShell(testShell)
	assert.Equal(t, "echo $SHELL\n", readBuffer.String(), "correct command not called")
	assert.Equal(t, errUnknownShell, err, "error not returned")
}
