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

func TestPacketTypeToProtocolID(t *testing.T) {
	assert.Equal(t, protocolSHandshake, MakeSType(Handshake, protocolSHandshake).ProtocolID())
	assert.Equal(t, protocolSPing, MakeSType(Status, protocolSPing).ProtocolID())
	assert.Equal(t, protocolCPong, MakeCType(Status, protocolCPong).ProtocolID())
	assert.Equal(t, protocolCEncryptionRequest, MakeCType(Login, protocolCEncryptionRequest).ProtocolID())
	assert.Equal(t, protocolSEncryptionResponse, MakeSType(Login, protocolSEncryptionResponse).ProtocolID())
	assert.Equal(t, protocolCPluginMessage, MakeCType(Play, protocolCPluginMessage).ProtocolID())
	assert.Equal(t, protocolSPluginMessage, MakeSType(Play, protocolSPluginMessage).ProtocolID())
}
