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

func TestString(t *testing.T) {
	testShell, _ := mockShell("")
	testPosix := Posix{agent: testShell}
	assert.Equal(t, "Posix - 127.0.0.1", testPosix.String())
}

func TestType(t *testing.T) {
	testPosix := Posix{}
	assert.Equal(t, "Posix", testPosix.Type())
}

func TestExecute(t *testing.T) {
	testAddr, _ := net.ResolveIPAddr("ip", "127.0.0.1")
	var buffer bytes.Buffer
	testShell := agent.NewShell(&buffer, 10, testAddr)
	testPosix := Posix{agent: testShell}
	testVal := "test"
	recvVals, err := testPosix.Execute(testVal)
	assert.Nil(t, err)
	assert.Equal(t, testVal, recvVals[0], "execute value incorrect")
}

func TestOS(t *testing.T) {
	testVal := "test"
	testShell, _ := mockShell(testVal)
	testPosix := Posix{agent: testShell}
	recvVal := testPosix.OS()
	assert.Equal(t, testVal, recvVal, "os enumeration failed")
}

func TestOSMultiLine(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testVal := "test\ntest\n"
	testShell, _ := mockShell(testVal)
	testPosix := Posix{agent: testShell}
	recvVal := testPosix.OS()
	assert.Equal(t, "unknown", recvVal, "os enumeration failed")
}
