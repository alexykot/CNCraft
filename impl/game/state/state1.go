package state

import (
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/status"
	"github.com/golangmc/minecraft-server/impl/protocol/client"
	"github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

/**
 * status
 */

func HandleState1(watcher pubsub.PubSub) {

	watcher.Subscribe(func(packet *server.SPacketRequest, conn base.Connection) {
		response := client.CPacketResponse{Status: status.DefaultResponse()}
		conn.SendPacket(&response)
	})

	watcher.Subscribe(func(packet *server.SPacketPing, conn base.Connection) {
		response := client.CPacketPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
