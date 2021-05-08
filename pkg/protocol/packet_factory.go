package protocol

import (
	"fmt"
)

type PacketFactory interface {
	MakeSPacket(id PacketType) (SPacket, error)
	MakeCPacket(id PacketType) (CPacket, error)
}

type packetFactory struct {
	sPackets map[PacketType]func() SPacket
	cPackets map[PacketType]func() CPacket
}

var factorySingleton packetFactory

func init() {
	factorySingleton = packetFactory{
		sPackets: createSPacketsMap(),
		cPackets: createCPacketsMap(),
	}
}

// GetPacketFactory returns the singleton of the packet factory object.
func GetPacketFactory() PacketFactory {
	return &factorySingleton
}

func (p *packetFactory) MakeCPacket(packetType PacketType) (CPacket, error) {
	creator, ok := p.cPackets[packetType]
	if !ok {
		return nil, fmt.Errorf("packet Type %d is not recognised", packetType)
	}

	return creator(), nil
}

func (p *packetFactory) MakeSPacket(packetType PacketType) (SPacket, error) {
	creator, ok := p.sPackets[packetType]
	if !ok {
		return nil, fmt.Errorf("packet Type %d is not recognised", packetType)
	}

	return creator(), nil
}

func createSPacketsMap() map[PacketType]func() SPacket {
	return map[PacketType]func() SPacket{
		// Handshake state
		SHandshake: func() SPacket { return &SPacketHandshake{} },

		// Status state
		SRequest: func() SPacket { return &SPacketRequest{} },
		SPing:    func() SPacket { return &SPacketPing{} },

		// Login state
		SLoginStart:          func() SPacket { return &SPacketLoginStart{} },
		SEncryptionResponse:  func() SPacket { return &SPacketEncryptionResponse{} },
		SLoginPluginResponse: func() SPacket { return &SPacketLoginPluginResponse{} },

		// Play state
		STeleportConfirm:      func() SPacket { return &SPacketTeleportConfirm{} },
		SPlayerMovement:       func() SPacket { return &SPacketPlayerMovement{} },
		SCloseWindow:          func() SPacket { return &SPacketCloseWindow{} },
		SQueryBlockNBT:        func() SPacket { return &SPacketQueryBlockNBT{} },
		SSetDifficulty:        func() SPacket { return &SPacketSetDifficulty{} },
		SChatMessage:          func() SPacket { return &SPacketChatMessage{} },
		SClientStatus:         func() SPacket { return &SPacketClientStatus{} },
		SClientSettings:       func() SPacket { return &SPacketClientSettings{} },
		SPluginMessage:        func() SPacket { return &SPacketPluginMessage{} },
		SKeepAlive:            func() SPacket { return &SPacketKeepAlive{} },
		SPlayerPosition:       func() SPacket { return &SPacketPlayerPosition{} },
		SEntityAction:         func() SPacket { return &SPacketEntityAction{} },
		SPlayerPosAndRotation: func() SPacket { return &SPacketPlayerPosAndRotation{} },
		SPlayerRotation:       func() SPacket { return &SPacketPlayerRotation{} },
		SPlayerAbilities:      func() SPacket { return &SPacketPlayerAbilities{} },
		SHeldItemChange:       func() SPacket { return &SPacketHeldItemChange{} },
		SAnimation:            func() SPacket { return &SPacketAnimation{} },
	}
}

func createCPacketsMap() map[PacketType]func() CPacket {
	return map[PacketType]func() CPacket{
		// Status state packets
		CResponse: func() CPacket { return &CPacketResponse{} },
		CPong:     func() CPacket { return &CPacketPong{} },

		// Login state packets
		CDisconnectLogin:    func() CPacket { return &CPacketDisconnectLogin{} },
		CEncryptionRequest:  func() CPacket { return &CPacketEncryptionRequest{} },
		CLoginSuccess:       func() CPacket { return &CPacketLoginSuccess{} },
		CSetCompression:     func() CPacket { return &CPacketSetCompression{} },
		CLoginPluginRequest: func() CPacket { return &CPacketLoginPluginRequest{} },

		// Play state packets
		CDisconnectPlay:        func() CPacket { return &CPacketDisconnectPlay{} },
		CChatMessage:           func() CPacket { return &CPacketChatMessage{} },
		CWindowConfirmation:    func() CPacket { return &CPacketWindowConfirmation{} },
		CJoinGame:              func() CPacket { return &CPacketJoinGame{} },
		CPluginMessage:         func() CPacket { return &CPacketPluginMessage{} },
		CPlayerPositionAndLook: func() CPacket { return &CPacketPlayerPositionAndLook{} },
		CKeepAlive:             func() CPacket { return &CPacketKeepAlive{} },
		CServerDifficulty:      func() CPacket { return &CPacketServerDifficulty{} },
		CPlayerAbilities:       func() CPacket { return &CPacketPlayerAbilities{} },
		CHeldItemChange:        func() CPacket { return &CPacketHeldItemChange{} },
		CWindowItems:           func() CPacket { return &CPacketWindowItems{} },
		CDeclareRecipes:        func() CPacket { return &CPacketDeclareRecipes{} },
		CChunkData:             func() CPacket { return &CPacketChunkData{} },
		CPlayerInfo:            func() CPacket { return &CPacketPlayerInfo{} },
		CEntityMetadata:        func() CPacket { return &CPacketEntityMetadata{} },
	}
}
