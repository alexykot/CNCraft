package state

import (
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/nats"
)

// RegisterHandlersState3 registers handlers for packets transmitted/received in the Play connection state.
func RegisterHandlersState3(ps nats.PubSub, logger *zap.Logger) {
	// TODO replace `join chan base.PlayerAndConnection`, `quit chan base.PlayerAndConnection` params with pubsub
	// TODO figure out what `tasking *task.Tasking` is for

	//tasking.EveryTime(10, time.Second, func(task *task.Task) {
	//
	//	api := apis.MinecraftServer()
	//
	//	// I hate this, add a functional method for player iterating
	//	for _, player := range api.Players() {
	//
	//		// also probably add one that returns both the player and their connection
	//		conn := api.ConnByUUID(player.UUID())
	//
	//		// keep player connection alive via keep alive
	//		conn.SendPacket(&protocol.CPacketKeepAlive{KeepAliveID: time.Now().UnixNano() / 1e6})
	//	}
	//})
	//
	//ps.Subscribe(func(packet *protocol.SPacketKeepAlive, conn base.Connection) {
	//	logger.DebugF("player %s is being kept alive", conn.Address())
	//})
	//
	//ps.Subscribe(func(packet *protocol.SPacketPluginMessage, conn base.Connection) {
	//	api := apis.MinecraftServer()
	//
	//	player := api.PlayerByConn(conn)
	//	if player == nil {
	//		return // log no player found?
	//	}
	//
	//	api.Watcher().Publish(implEvent.PlayerPluginMessagePullEvent{
	//		Conn: base.PlayerAndConnection{
	//			Connection: conn,
	//			Player:     player,
	//		},
	//		Channel: packet.Message.Chan(),
	//		Message: packet.Message,
	//	})
	//})
	//
	//ps.Subscribe(func(packet *protocol.SPacketChatMessage, conn base.Connection) {
	//	api := apis.MinecraftServer()
	//
	//	who := api.PlayerByConn(conn)
	//	out := msgs.
	//		New(who.Name()).SetColor(chat.White).
	//		Add(":").SetColor(chat.Gray).
	//		Add(" ").
	//		Add(chat.Translate(packet.Message)).SetColor(chat.White).
	//		AsText() // why not just use translate?
	//
	//	api.Broadcast(out)
	//})
	//
	//go func() {
	//	for conn := range join {
	//		apis.MinecraftServer().Watcher().Publish(implEvent.PlayerConnJoinEvent{Conn: conn})
	//
	//		conn.SendPacket(&protocol.CPacketJoinGame{
	//			EntityID:      int32(conn.EntityUUID()),
	//			Hardcore:      false,
	//			GameMode:      game.CREATIVE,
	//			Dimension:     game.OVERWORLD,
	//			HashedSeed:    values.DefaultWorldHashedSeed,
	//			MaxPlayers:    10,
	//			LevelType:     game.DEFAULT,
	//			ViewDistance:  12,
	//			ReduceDebug:   false,
	//			RespawnScreen: false,
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketPluginMessage{
	//			Message: &plugin.Brand{
	//				Name: chat.Translate(fmt.Sprintf("&b%s&r &a%s&r", "GoLangMc", apis.MinecraftServer().ServerVersion())),
	//			},
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketServerDifficulty{
	//			Difficulty: game.PEACEFUL,
	//			Locked:     true,
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketPlayerAbilities{
	//			Abilities: client.PlayerAbilities{
	//				Invulnerable: true,
	//				Flying:       true,
	//				AllowFlight:  true,
	//				InstantBuild: false,
	//			},
	//			FlyingSpeed: 0.05, // default value
	//			FieldOfView: 0.1,  // default value
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketHeldItemChange{
	//			Slot: client.SLOT_0,
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketDeclareRecipes{})
	//
	//		conn.SendPacket(&protocol.CPacketPlayerLocation{
	//			SomeID: 0,
	//			Location: data.Location{
	//				PositionF: data.PositionF{
	//					X: 0,
	//					Y: 10,
	//					Z: 0,
	//				},
	//				RotationF: data.RotationF{
	//					AxisX: 0,
	//					AxisY: 0,
	//				},
	//			},
	//			Relative: client.Relativity{},
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketPlayerInfo{
	//			Action: client.AddPlayer,
	//			Values: []client.PlayerInfo{
	//				&client.PlayerInfoAddPlayer{Player: conn.Player},
	//			},
	//		})
	//
	//		conn.SendPacket(&protocol.CPacketEntityMetadata{Entity: conn.Player})
	//
	//		level := implLevel.NewLevel("test")
	//		implLevel.GenSuperFlat(level, 6)
	//
	//		for _, chunk := range level.Chunks() {
	//			conn.SendPacket(&protocol.CPacketChunkData{Chunk: chunk})
	//		}
	//
	//		logger.DebugF("chunks sent to player: %s", conn.Player.Name())
	//
	//		conn.SendPacket(&protocol.CPacketPlayerLocation{
	//			SomeID: 1,
	//			Location: data.Location{
	//				PositionF: data.PositionF{
	//					X: 0,
	//					Y: 10,
	//					Z: 0,
	//				},
	//				RotationF: data.RotationF{
	//					AxisX: 0,
	//					AxisY: 0,
	//				},
	//			},
	//			Relative: client.Relativity{},
	//		})
	//	}
	//}()
	//
	//go func() {
	//	for conn := range quit {
	//		apis.MinecraftServer().Watcher().Publish(implEvent.PlayerConnQuitEvent{Conn: conn})
	//	}
	//}()
}
