package state

import (
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/pkg/bus"
	"github.com/alexykot/cncraft/pkg/protocol"
)

// RegisterHandlersState2 registers handlers for packets transmitted/received in the Login connection state.
func RegisterHandlersState2(ps bus.PubSub, logger *zap.Logger) {
	// TODO replace `join chan base.PlayerAndConnection` with pubsub

	{ // server bound packets
		ps.Subscribe(protocol.MakePacketTopic(protocol.SLoginStart), func(envelopeIn bus.Envelope) {
			loginStartPack, ok := envelopeIn.GetMessage().(protocol.SPacketLoginStart)
			if !ok {
				// DEBT figure out logging here
				return
			}
			loginSuccessPack := protocol.CPacketLoginSuccess{
				PlayerUUID: "",
				PlayerName: loginStartPack.PlayerName,
			}

			ps.Publish(protocol.MakePacketTopic(protocol.CLoginSuccess),
				bus.NewEnvelope(loginSuccessPack, envelopeIn.GetAllMeta()))
		})

		ps.Subscribe(protocol.MakePacketTopic(protocol.SPing), func(envelopeIn bus.Envelope) {
			packet, ok := envelopeIn.GetMessage().(protocol.SPacketPing)
			if !ok {
				// DEBT figure out logging here
				return
			}

			ps.Publish(protocol.MakePacketTopic(protocol.CPong),
				bus.NewEnvelope(protocol.CPacketPong{Ping: packet.Ping}, envelopeIn.GetAllMeta()))
		})
	}



	// DEBT The authentication and encryption can be skipped in offline mode. Will get back to it later.
	//watcher.Subscribe(protocol.MakePacketTopic(protocol.SLoginStart),
	//	func(packet *protocol.SPacketLoginStart, conn base.Connection) {
	//	conn.CertifyValues(packet.PlayerName)
	//
	//	_, public := auth.NewCrypt()
	//
	//	response := protocol.CPacketEncryptionRequest{
	//		Server: "",
	//		Public: public,
	//		Verify: conn.CertifyData(),
	//	}
	//
	//	conn.SendPacket(&response)
	//})
	//
	//watcher.Subscribe(protocol.MakePacketTopic(protocol.SEncryptionResponse),
	//	func(packet *protocol.SPacketEncryptionResponse, conn base.Connection) {
	//	defer func() {
	//		// DEBT this is a fucking mess, panics are used profusely instead of proper error handling
	//		if err := recover(); err != nil {
	//			conn.SendPacket(&protocol.CPacketDisconnect{
	//				Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
	//			})
	//		}
	//	}()
	//
	//	ver, err := auth.Decrypt(packet.Verify)
	//	if err != nil {
	//		panic(fmt.Errorf("failed to decrypt token: %s\n%v\n", conn.CertifyName(), err))
	//	}
	//
	//	if !bytes.Equal(ver, conn.CertifyData()) {
	//		panic(fmt.Errorf("encryption failed, tokens are different: %s\n%v | %v", conn.CertifyName(), ver, conn.CertifyData()))
	//	}
	//
	//	sec, err := auth.Decrypt(packet.Secret)
	//	if err != nil {
	//		panic(fmt.Errorf("failed to decrypt secret: %s\n%v\n", conn.CertifyName(), err))
	//	}
	//
	//	conn.CertifyUpdate(sec) // enable encryption on the connection
	//
	//	auth.RunAuthGet(sec, conn.CertifyName(), func(auth *auth.Auth, err error) {
	//		defer func() {
	//			if err := recover(); err != nil {
	//				conn.SendPacket(&protocol.CPacketDisconnect{
	//					Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
	//				})
	//			}
	//		}()
	//
	//		if err != nil {
	//			panic(fmt.Errorf("failed to authenticate: %s\n%v\n", conn.CertifyName(), err))
	//		}
	//
	//		uuid, err := uuid.TextToUUID(auth.UUID)
	//		if err != nil {
	//			panic(fmt.Errorf("failed to decode uuid for %s: %s\n%v\n", conn.CertifyName(), auth.UUID, err))
	//		}
	//
	//		prof := game.Profile{
	//			UUID: uuid,
	//			Name: auth.Name,
	//		}
	//
	//		for _, prop := range auth.Prop {
	//			prof.Properties = append(prof.Properties, &game.ProfileProperty{
	//				Name:      prop.Name,
	//				Value:     prop.Data,
	//				Signature: prop.Sign,
	//			})
	//		}
	//
	//		player := ents.NewPlayer(&prof, conn)
	//
	//		conn.SendPacket(&protocol.CPacketLoginSuccess{
	//			PlayerName: player.Name(),
	//			PlayerUUID: player.UUID().String(),
	//		})
	//
	//		conn.SetState(protocol.Play)
	//
	//		join <- base.ServerPlayer{
	//			PlayerCharacter:     player,
	//			Connection: conn,
	//		}
	//	})
	//
	//})

}
