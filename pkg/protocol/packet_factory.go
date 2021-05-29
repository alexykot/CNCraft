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
		SKeepAlive:       func() SPacket { return &SPacketKeepAlive{} },
		SClientStatus:    func() SPacket { return &SPacketClientStatus{} },
		SClientSettings:  func() SPacket { return &SPacketClientSettings{} },
		SSetDifficulty:   func() SPacket { return &SPacketSetDifficulty{} },
		SPlayerAbilities: func() SPacket { return &SPacketPlayerAbilities{} },

		SPluginMessage:   func() SPacket { return &SPacketPluginMessage{} },
		STeleportConfirm: func() SPacket { return &SPacketTeleportConfirm{} },
		SQueryBlockNBT:   func() SPacket { return &SPacketQueryBlockNBT{} },
		SChatMessage:     func() SPacket { return &SPacketChatMessage{} },

		SHeldItemChange:     func() SPacket { return &SPacketHeldItemChange{} },
		SEntityAction:       func() SPacket { return &SPacketEntityAction{} },
		SAnimation:          func() SPacket { return &SPacketAnimation{} },
		SClickWindow:        func() SPacket { return &SPacketClickWindow{} },
		SCloseWindow:        func() SPacket { return &SPacketCloseWindow{} },
		SWindowConfirmation: func() SPacket { return &SPacketWindowConfirmation{} },

		SPlayerMovement:       func() SPacket { return &SPacketPlayerMovement{} },
		SPlayerPosition:       func() SPacket { return &SPacketPlayerPosition{} },
		SPlayerPosAndRotation: func() SPacket { return &SPacketPlayerPosAndRotation{} },
		SPlayerRotation:       func() SPacket { return &SPacketPlayerRotation{} },
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
