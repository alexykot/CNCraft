package state

import (
	"bytes"
	"fmt"
	"github.com/golangmc/minecraft-server/impl/protocol"

	"github.com/golangmc/minecraft-server/apis/data/chat"
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/uuid"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/game/auth"
	"github.com/golangmc/minecraft-server/impl/game/ents"
	"github.com/golangmc/minecraft-server/impl/protocol/client"
	"github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

/**
 * login
 */

func HandleState2(watcher pubsub.PubSub, join chan base.PlayerAndConnection) {

	watcher.Subscribe(func(packet *server.SPacketLoginStart, conn base.Connection) {
		conn.CertifyValues(packet.PlayerName)

		_, public := auth.NewCrypt()

		response := client.CPacketEncryptionRequest{
			Server: "",
			Public: public,
			Verify: conn.CertifyData(),
		}

		conn.SendPacket(&response)
	})

	watcher.Subscribe(func(packet *server.SPacketEncryptionResponse, conn base.Connection) {
		defer func() {
			if err := recover(); err != nil {
				conn.SendPacket(&client.CPacketDisconnect{
					Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
				})
			}
		}()

		ver, err := auth.Decrypt(packet.Verify)
		if err != nil {
			panic(fmt.Errorf("failed to decrypt token: %s\n%v\n", conn.CertifyName(), err))
		}

		if !bytes.Equal(ver, conn.CertifyData()) {
			panic(fmt.Errorf("encryption failed, tokens are different: %s\n%v | %v", conn.CertifyName(), ver, conn.CertifyData()))
		}

		sec, err := auth.Decrypt(packet.Secret)
		if err != nil {
			panic(fmt.Errorf("failed to decrypt secret: %s\n%v\n", conn.CertifyName(), err))
		}

		conn.CertifyUpdate(sec) // enable encryption on the connection

		auth.RunAuthGet(sec, conn.CertifyName(), func(auth *auth.Auth, err error) {
			defer func() {
				if err := recover(); err != nil {
					conn.SendPacket(&client.CPacketDisconnect{
						Reason: *msgs.New(fmt.Sprintf("Authentication failed: %v", err)).SetColor(chat.Red),
					})
				}
			}()

			if err != nil {
				panic(fmt.Errorf("failed to authenticate: %s\n%v\n", conn.CertifyName(), err))
			}

			uuid, err := uuid.TextToUUID(auth.UUID)
			if err != nil {
				panic(fmt.Errorf("failed to decode uuid for %s: %s\n%v\n", conn.CertifyName(), auth.UUID, err))
			}

			prof := game.Profile{
				UUID: uuid,
				Name: auth.Name,
			}

			for _, prop := range auth.Prop {
				prof.Properties = append(prof.Properties, &game.ProfileProperty{
					Name:      prop.Name,
					Value:     prop.Data,
					Signature: prop.Sign,
				})
			}

			player := ents.NewPlayer(&prof, conn)

			conn.SendPacket(&client.CPacketLoginSuccess{
				PlayerName: player.Name(),
				PlayerUUID: player.UUID().String(),
			})

			conn.SetState(protocol.Play)

			join <- base.PlayerAndConnection{
				Player:     player,
				Connection: conn,
			}
		})

	})

}
