package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	testHandler := TCPHandler{}
	go testHandler.Handle(1337)
	defer testHandler.Stop()
	for !testHandler.active {
	}
	assert.Equal(t, "TCP Handler - 0.0.0.0:1337", testHandler.String(), "to string incorrect")
}

func TestType(t *testing.T) {
	testHandler := TCPHandler{}
	assert.Equal(t, "TCP", testHandler.Type(), "type incorrect")
}

