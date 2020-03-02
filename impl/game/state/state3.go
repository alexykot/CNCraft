package state

import (
	"fmt"
	"time"

	"github.com/golangmc/minecraft-server/apis"
	"github.com/golangmc/minecraft-server/apis/data"
	"github.com/golangmc/minecraft-server/apis/data/chat"
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/client"
	"github.com/golangmc/minecraft-server/impl/data/plugin"
	"github.com/golangmc/minecraft-server/impl/data/values"
	implEvent "github.com/golangmc/minecraft-server/impl/game/event"
	implLevel "github.com/golangmc/minecraft-server/impl/game/level"
	clientPacket "github.com/golangmc/minecraft-server/impl/protocol/client"
	serverPacket "github.com/golangmc/minecraft-server/impl/protocol/server"
	"github.com/golangmc/minecraft-server/pkg/pubsub"
)

func HandleState3(watcher pubsub.PubSub, logger *logs.Logging, tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) {

	tasking.EveryTime(10, time.Second, func(task *task.Task) {

		api := apis.MinecraftServer()

		// I hate this, add a functional method for player iterating
		for _, player := range api.Players() {

			// also probably add one that returns both the player and their connection
			conn := api.ConnByUUID(player.UUID())

			// keep player connection alive via keep alive
			conn.SendPacket(&clientPacket.CPacketKeepAlive{KeepAliveID: time.Now().UnixNano() / 1e6})
		}
	})

	watcher.Subscribe(func(packet *serverPacket.SPacketKeepAlive, conn base.Connection) {
		logger.DebugF("player %s is being kept alive", conn.Address())
	})

	watcher.Subscribe(func(packet *serverPacket.SPacketPluginMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		player := api.PlayerByConn(conn)
		if player == nil {
			return // log no player found?
		}

		api.Watcher().Publish(implEvent.PlayerPluginMessagePullEvent{
			Conn: base.PlayerAndConnection{
				Connection: conn,
				Player:     player,
			},
			Channel: packet.Message.Chan(),
			Message: packet.Message,
		})
	})

	watcher.Subscribe(func(packet *serverPacket.SPacketChatMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		who := api.PlayerByConn(conn)
		out := msgs.
			New(who.Name()).SetColor(chat.White).
			Add(":").SetColor(chat.Gray).
			Add(" ").
			Add(chat.Translate(packet.Message)).SetColor(chat.White).
			AsText() // why not just use translate?

		api.Broadcast(out)
	})

	go func() {
		for conn := range join {
			apis.MinecraftServer().Watcher().Publish(implEvent.PlayerConnJoinEvent{Conn: conn})

			conn.SendPacket(&clientPacket.CPacketJoinGame{
				EntityID:      int32(conn.EntityUUID()),
				Hardcore:      false,
				GameMode:      game.CREATIVE,
				Dimension:     game.OVERWORLD,
				HashedSeed:    values.DefaultWorldHashedSeed,
				MaxPlayers:    10,
				LevelType:     game.DEFAULT,
				ViewDistance:  12,
				ReduceDebug:   false,
				RespawnScreen: false,
			})

			conn.SendPacket(&clientPacket.CPacketPluginMessage{
				Message: &plugin.Brand{
					Name: chat.Translate(fmt.Sprintf("&b%s&r &a%s&r", "GoLangMc", apis.MinecraftServer().ServerVersion())),
				},
			})

			conn.SendPacket(&clientPacket.CPacketServerDifficulty{
				Difficulty: game.PEACEFUL,
				Locked:     true,
			})

			conn.SendPacket(&clientPacket.CPacketPlayerAbilities{
				Abilities: client.PlayerAbilities{
					Invulnerable: true,
					Flying:       true,
					AllowFlight:  true,
					InstantBuild: false,
				},
				FlyingSpeed: 0.05, // default value
				FieldOfView: 0.1,  // default value
			})

			conn.SendPacket(&clientPacket.CPacketHeldItemChange{
				Slot: client.SLOT_0,
			})

			conn.SendPacket(&clientPacket.CPacketDeclareRecipes{})

			conn.SendPacket(&clientPacket.CPacketPlayerLocation{
				SomeID: 0,
				Location: data.Location{
					PositionF: data.PositionF{
						X: 0,
						Y: 10,
						Z: 0,
					},
					RotationF: data.RotationF{
						AxisX: 0,
						AxisY: 0,
					},
				},
				Relative: client.Relativity{},
			})

			conn.SendPacket(&clientPacket.CPacketPlayerInfo{
				Action: client.AddPlayer,
				Values: []client.PlayerInfo{
					&client.PlayerInfoAddPlayer{Player: conn.Player},
				},
			})

			conn.SendPacket(&clientPacket.CPacketEntityMetadata{Entity: conn.Player})

			level := implLevel.NewLevel("test")
			implLevel.GenSuperFlat(level, 6)

			for _, chunk := range level.Chunks() {
				conn.SendPacket(&clientPacket.CPacketChunkData{Chunk: chunk})
			}

			logger.DebugF("chunks sent to player: %s", conn.Player.Name())

			conn.SendPacket(&clientPacket.CPacketPlayerLocation{
				SomeID: 1,
				Location: data.Location{
					PositionF: data.PositionF{
						X: 0,
						Y: 10,
						Z: 0,
					},
					RotationF: data.RotationF{
						AxisX: 0,
						AxisY: 0,
					},
				},
				Relative: client.Relativity{},
			})
		}
	}()

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().Publish(implEvent.PlayerConnQuitEvent{Conn: conn})
		}
	}()
}
