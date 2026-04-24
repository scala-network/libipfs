package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIP(t *testing.T) {
	assert.True(t, IsValidIP("127.0.0.1"))
	assert.True(t, IsValidIP("0.0.0.0"))
	assert.True(t, IsValidIP("255.255.255.255"))
	assert.True(t, IsValidIP("::1"))
	assert.True(t, IsValidIP("2001:db8::1"))

	assert.False(t, IsValidIP(""))
	assert.False(t, IsValidIP("not-an-ip"))
	assert.False(t, IsValidIP("256.0.0.1"))
	assert.False(t, IsValidIP("1.2.3"))
}

func TestIsValidPort(t *testing.T) {
	assert.True(t, IsValidPort(1))
	assert.True(t, IsValidPort(80))
	assert.True(t, IsValidPort(4001))
	assert.True(t, IsValidPort(65535))

	assert.False(t, IsValidPort(0))
	assert.False(t, IsValidPort(-1))
	assert.False(t, IsValidPort(65536))
	assert.False(t, IsValidPort(100000))
}

func TestCheckBind(t *testing.T) {
	assert.False(t, CheckBind("not-an-ip", 4001))
	assert.False(t, CheckBind("127.0.0.1", 0))
	assert.False(t, CheckBind("127.0.0.1", -1))
	assert.False(t, CheckBind("", 4001))
	assert.True(t, CheckBind("127.0.0.1", 49200))
}
