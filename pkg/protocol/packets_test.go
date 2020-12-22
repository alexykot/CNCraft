package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSType(t *testing.T) {
	assert.Equal(t, SHandshake, MakeSType(Handshake, protocolSHandshake))
	assert.Equal(t, SPing, MakeSType(Status, protocolSPing))
	assert.Equal(t, SRequest, MakeSType(Status, protocolSRequest))
}
