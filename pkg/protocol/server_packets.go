package protocol

import (
	"fmt"

	"github.com/alexykot/cncraft/pkg/buffer"
	"github.com/alexykot/cncraft/pkg/game"
	"github.com/alexykot/cncraft/pkg/game/data"
	"github.com/alexykot/cncraft/pkg/game/items"
	"github.com/alexykot/cncraft/pkg/game/player"
	pItems "github.com/alexykot/cncraft/pkg/protocol/items"
	"github.com/alexykot/cncraft/pkg/protocol/plugin"
)

// HANDSHAKE STATE PACKETS
type SPacketHandshake struct {
	Version int32

	Host string
	Port uint16

	NextState State
}

func (p *SPacketHandshake) ProtocolID() ProtocolPacketID { return protocolSHandshake }
func (p *SPacketHandshake) Type() PacketType             { return SHandshake }
func (p *SPacketHandshake) Pull(reader *buffer.Buffer) error {
	var err error

	p.Version = reader.PullVarInt()
	p.Host = reader.PullString()
	p.Port = reader.PullUint16()

	nextState := reader.PullVarInt()

	if p.NextState, err = IntToState(int(nextState)); err != nil {
		return fmt.Errorf("failed to parse handshake  next state: %w", err)
	}

	return nil
}

func (p *SPacketHandshake) Push(writer *buffer.Buffer) {
	writer.PushVarInt(int32(protocolSHandshake))

	writer.PushVarInt(p.Version)
	writer.PushString(p.Host)
	writer.PushUint16(p.Port)

	writer.PushVarInt(int32(p.NextState))
}

// STATUS STATE PACKETS
type SPacketRequest struct{}

func (p *SPacketRequest) ProtocolID() ProtocolPacketID { return protocolSRequest }
func (p *SPacketRequest) Type() PacketType             { return SRequest }
func (p *SPacketRequest) Pull(reader *buffer.Buffer) error {
	// no fields
	return nil
}

type SPacketPing struct {
	Payload int64
}

func (p *SPacketPing) ProtocolID() ProtocolPacketID { return protocolSPing }
func (p *SPacketPing) Type() PacketType             { return SPing }
func (p *SPacketPing) Pull(reader *buffer.Buffer) error {
	p.Payload = reader.PullInt64()
	return nil // DEBT actually check for errors
}

// LOGIN STATE PACKETS
type SPacketLoginStart struct {
	Username string
}

func (p *SPacketLoginStart) ProtocolID() ProtocolPacketID { return protocolSLoginStart }
func (p *SPacketLoginStart) Type() PacketType             { return SLoginStart }
func (p *SPacketLoginStart) Pull(reader *buffer.Buffer) error {
	p.Username = reader.PullString()
	return nil // DEBT actually check for errors
}

func (p *SPacketLoginStart) Push(writer *buffer.Buffer) {
	writer.PushVarInt(int32(protocolSLoginStart))
	writer.PushString(p.Username)
}

type SPacketEncryptionResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func (p *SPacketEncryptionResponse) ProtocolID() ProtocolPacketID { return protocolSEncryptionResponse }
func (p *SPacketEncryptionResponse) Type() PacketType             { return SEncryptionResponse }
func (p *SPacketEncryptionResponse) Pull(reader *buffer.Buffer) error {
	p.SharedSecret = reader.PullBytes()
	p.VerifyToken = reader.PullBytes()
	return nil // DEBT actually check for errors
}

type SPacketLoginPluginResponse struct {
	Message int32
	Success bool
	OptData []byte
}

func (p *SPacketLoginPluginResponse) ProtocolID() ProtocolPacketID {
	return protocolSLoginPluginResponse
}
func (p *SPacketLoginPluginResponse) Type() PacketType { return SLoginPluginResponse }
func (p *SPacketLoginPluginResponse) Pull(reader *buffer.Buffer) error {
	p.Message = reader.PullVarInt()
	p.Success = reader.PullBool()
	p.OptData = reader.Bytes()[reader.IndexI():reader.Len()]
	return nil // DEBT actually check for errors
}

// PLAY STATE PACKETS

type SPacketTeleportConfirm struct {
	TeleportID int32
}

func (p *SPacketTeleportConfirm) ProtocolID() ProtocolPacketID { return protocolSTeleportConfirm }
func (p *SPacketTeleportConfirm) Type() PacketType             { return STeleportConfirm }
func (p *SPacketTeleportConfirm) Pull(reader *buffer.Buffer) error {
	p.TeleportID = reader.PullVarInt()
	return nil // DEBT actually check for errors
}

type SPacketQueryBlockNBT struct {
	TransactionID int32
	Location      data.PositionI
}

func (p *SPacketQueryBlockNBT) ProtocolID() ProtocolPacketID { return protocolSQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Type() PacketType             { return SQueryBlockNBT }
func (p *SPacketQueryBlockNBT) Pull(reader *buffer.Buffer) error {
	p.TransactionID = reader.PullVarInt()
	p.Location.Pull(reader)
	return nil // DEBT actually check for errors
}

type SPacketQueryEntityNBT struct{}

func (p *SPacketQueryEntityNBT) ProtocolID() ProtocolPacketID { return protocolSQueryEntityNBT }
func (p *SPacketQueryEntityNBT) Type() PacketType             { return SQueryEntityNBT }
func (p *SPacketQueryEntityNBT) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSetDifficulty struct {
	Difficulty game.Difficulty
}

func (p *SPacketSetDifficulty) ProtocolID() ProtocolPacketID { return protocolSSetDifficulty }
func (p *SPacketSetDifficulty) Type() PacketType             { return SSetDifficulty }
func (p *SPacketSetDifficulty) Pull(reader *buffer.Buffer) error {
	p.Difficulty = game.DifficultyValueOf(reader.PullByte())
	return nil // DEBT actually check for errors
}

type SPacketChatMessage struct {
	Message string
}

func (p *SPacketChatMessage) ProtocolID() ProtocolPacketID { return protocolSChatMessage }
func (p *SPacketChatMessage) Type() PacketType             { return SChatMessage }
func (p *SPacketChatMessage) Pull(reader *buffer.Buffer) error {
	p.Message = reader.PullString()
	return nil // DEBT actually check for errors
}

type SPacketClientStatus struct {
	Action player.StatusAction
}

func (p *SPacketClientStatus) ProtocolID() ProtocolPacketID { return protocolSClientStatus }
func (p *SPacketClientStatus) Type() PacketType             { return SClientStatus }
func (p *SPacketClientStatus) Pull(reader *buffer.Buffer) error {
	p.Action = player.StatusAction(reader.PullVarInt())
	return nil // DEBT actually check for errors
}

type SPacketClientSettings struct {
	Locale       string
	ViewDistance byte
	ChatMode     player.ChatMode
	ChatColors   bool // if false, strip messages of colors before sending
	SkinParts    player.SkinParts
	MainHand     player.MainHand
}

func (p *SPacketClientSettings) ProtocolID() ProtocolPacketID { return protocolSClientSettings }
func (p *SPacketClientSettings) Type() PacketType             { return SClientSettings }
func (p *SPacketClientSettings) Pull(reader *buffer.Buffer) error {
	p.Locale = reader.PullString()
	p.ViewDistance = reader.PullByte()
	p.ChatMode = player.ChatMode(reader.PullVarInt())
	p.ChatColors = reader.PullBool()

	parts := player.SkinParts{}
	parts.Pull(reader)

	p.SkinParts = parts
	p.MainHand = player.MainHand(reader.PullVarInt())
	return nil // DEBT actually check for errors
}

type SPacketTabComplete struct{}

func (p *SPacketTabComplete) ProtocolID() ProtocolPacketID { return protocolSTabComplete }
func (p *SPacketTabComplete) Type() PacketType             { return STabComplete }
func (p *SPacketTabComplete) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketWindowConfirmation struct {
	WindowID items.WindowID
	ActionID int16
	Accepted bool
}

func (p *SPacketWindowConfirmation) ProtocolID() ProtocolPacketID { return protocolSWindowConfirmation }
func (p *SPacketWindowConfirmation) Type() PacketType             { return SWindowConfirmation }
func (p *SPacketWindowConfirmation) Pull(reader *buffer.Buffer) error {
	p.WindowID = items.WindowID(reader.PullByte())
	p.ActionID = reader.PullInt16()
	p.Accepted = reader.PullBool()
	return nil
}

type SPacketClickWindowButton struct{}

func (p *SPacketClickWindowButton) ProtocolID() ProtocolPacketID { return protocolSClickWindowButton }
func (p *SPacketClickWindowButton) Type() PacketType             { return SClickWindowButton }
func (p *SPacketClickWindowButton) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketClickWindow struct {
	WindowID    items.WindowID
	SlotID      int16
	Button      uint8
	ActionID    int16
	Mode        int16
	ClickedItem items.Slot
}

func (p *SPacketClickWindow) ProtocolID() ProtocolPacketID { return protocolSClickWindow }
func (p *SPacketClickWindow) Type() PacketType             { return SClickWindow }
func (p *SPacketClickWindow) Pull(reader *buffer.Buffer) error {
	p.WindowID = items.WindowID(reader.PullByte())
	p.SlotID = reader.PullInt16()
	p.Button = reader.PullByte()
	p.ActionID = reader.PullInt16()
	p.Mode = int16(reader.PullVarInt())
	p.ClickedItem.IsPresent = reader.PullBool()
	if p.ClickedItem.IsPresent {
		p.ClickedItem.ItemID = pItems.ItemID(reader.PullVarInt())
		p.ClickedItem.ItemCount = int16(reader.PullByte())
	}
	return nil
}

type SPacketCloseWindow struct {
	WindowID items.WindowID
}

func (p *SPacketCloseWindow) ProtocolID() ProtocolPacketID { return protocolSCloseWindow }
func (p *SPacketCloseWindow) Type() PacketType             { return SCloseWindow }
func (p *SPacketCloseWindow) Pull(reader *buffer.Buffer) error {
	p.WindowID = items.WindowID(reader.PullByte())
	return nil
}

type SPacketPluginMessage struct {
	Message plugin.Message
}

func (p *SPacketPluginMessage) ProtocolID() ProtocolPacketID { return protocolSPluginMessage }
func (p *SPacketPluginMessage) Type() PacketType             { return SPluginMessage }
func (p *SPacketPluginMessage) Pull(reader *buffer.Buffer) error {
	channel := reader.PullString()
	message := plugin.GetMessageForChannel(plugin.Channel(channel))
	if message == nil {
		return fmt.Errorf("channel `%s` not found", channel)
	}

	message.Pull(reader)

	p.Message = message

	return nil
}

type SPacketEditBook struct{}

func (p *SPacketEditBook) ProtocolID() ProtocolPacketID { return protocolSEditBook }
func (p *SPacketEditBook) Type() PacketType             { return SEditBook }
func (p *SPacketEditBook) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketInteractEntity struct{}

func (p *SPacketInteractEntity) ProtocolID() ProtocolPacketID { return protocolSInteractEntity }
func (p *SPacketInteractEntity) Type() PacketType             { return SInteractEntity }
func (p *SPacketInteractEntity) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketGenerateStructure struct{}

func (p *SPacketGenerateStructure) ProtocolID() ProtocolPacketID { return protocolSGenerateStructure }
func (p *SPacketGenerateStructure) Type() PacketType             { return SGenerateStructure }
func (p *SPacketGenerateStructure) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketKeepAlive struct {
	KeepAliveID int64
}

func (p *SPacketKeepAlive) ProtocolID() ProtocolPacketID { return protocolSKeepAlive }
func (p *SPacketKeepAlive) Type() PacketType             { return SKeepAlive }
func (p *SPacketKeepAlive) Pull(reader *buffer.Buffer) error {
	p.KeepAliveID = reader.PullInt64()
	return nil // DEBT actually check for errors
}

type SPacketLockDifficulty struct{}

func (p *SPacketLockDifficulty) ProtocolID() ProtocolPacketID { return protocolSLockDifficulty }
func (p *SPacketLockDifficulty) Type() PacketType             { return SLockDifficulty }
func (p *SPacketLockDifficulty) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketPlayerPosition struct {
	Position data.PositionF
	OnGround bool
}

func (p *SPacketPlayerPosition) ProtocolID() ProtocolPacketID { return protocolSPlayerPosition }
func (p *SPacketPlayerPosition) Type() PacketType             { return SPlayerPosition }
func (p *SPacketPlayerPosition) Pull(reader *buffer.Buffer) error {
	p.Position = data.PositionF{
		X: reader.PullFloat64(),
		Y: reader.PullFloat64(),
		Z: reader.PullFloat64(),
	}

	p.OnGround = reader.PullBool()
	return nil // DEBT actually check for errors
}

type SPacketPlayerPosAndRotation struct {
	Location data.Location
	OnGround bool
}

func (p *SPacketPlayerPosAndRotation) ProtocolID() ProtocolPacketID {
	return protocolSPlayerPosAndRotation
}
func (p *SPacketPlayerPosAndRotation) Type() PacketType { return SPlayerPosAndRotation }
func (p *SPacketPlayerPosAndRotation) Pull(reader *buffer.Buffer) error {
	p.Location = data.Location{
		PositionF: data.PositionF{
			X: reader.PullFloat64(),
			Y: reader.PullFloat64(),
			Z: reader.PullFloat64(),
		},
		RotationF: data.RotationF{
			Yaw:   reader.PullFloat32(),
			Pitch: reader.PullFloat32(),
		},
	}

	p.OnGround = reader.PullBool()
	return nil // DEBT actually check for errors
}

type SPacketPlayerRotation struct {
	Rotation data.RotationF
	OnGround bool
}

func (p *SPacketPlayerRotation) ProtocolID() ProtocolPacketID { return protocolSPlayerRotation }
func (p *SPacketPlayerRotation) Type() PacketType             { return SPlayerRotation }
func (p *SPacketPlayerRotation) Pull(reader *buffer.Buffer) error {
	p.Rotation = data.RotationF{
		Yaw:   reader.PullFloat32(),
		Pitch: reader.PullFloat32(),
	}

	p.OnGround = reader.PullBool()
	return nil // DEBT actually check for errors
}

type SPacketPlayerMovement struct {
	OnGround bool
}

func (p *SPacketPlayerMovement) ProtocolID() ProtocolPacketID { return protocolSPlayerMovement }
func (p *SPacketPlayerMovement) Type() PacketType             { return SPlayerMovement }
func (p *SPacketPlayerMovement) Pull(reader *buffer.Buffer) error {
	p.OnGround = reader.PullBool()
	return nil
}

type SPacketVehicleMove struct{}

func (p *SPacketVehicleMove) ProtocolID() ProtocolPacketID { return protocolSVehicleMove }
func (p *SPacketVehicleMove) Type() PacketType             { return SVehicleMove }
func (p *SPacketVehicleMove) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSteerBoat struct{}

func (p *SPacketSteerBoat) ProtocolID() ProtocolPacketID { return protocolSSteerBoat }
func (p *SPacketSteerBoat) Type() PacketType             { return SSteerBoat }
func (p *SPacketSteerBoat) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketPickItem struct{}

func (p *SPacketPickItem) ProtocolID() ProtocolPacketID { return protocolSPickItem }
func (p *SPacketPickItem) Type() PacketType             { return SPickItem }
func (p *SPacketPickItem) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketCraftRecipeRequest struct{}

func (p *SPacketCraftRecipeRequest) ProtocolID() ProtocolPacketID { return protocolSCraftRecipeRequest }
func (p *SPacketCraftRecipeRequest) Type() PacketType             { return SCraftRecipeRequest }
func (p *SPacketCraftRecipeRequest) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketPlayerAbilities struct {
	Abilities   player.Abilities
	FlightSpeed float32
	GroundSpeed float32
}

func (p *SPacketPlayerAbilities) ProtocolID() ProtocolPacketID { return protocolSPlayerAbilities }
func (p *SPacketPlayerAbilities) Type() PacketType             { return SPlayerAbilities }
func (p *SPacketPlayerAbilities) Pull(reader *buffer.Buffer) error {
	panic("changes in 1.16.4 need to be implemented")

	abilities := player.Abilities{}
	abilities.Pull(reader)

	p.Abilities = abilities

	p.FlightSpeed = reader.PullFloat32()
	p.GroundSpeed = reader.PullFloat32()
	return nil // DEBT actually check for errors
}

type SPacketPlayerDigging struct {
	Status   int32
	Position *data.PositionI
	Face     byte
}

func (p *SPacketPlayerDigging) ProtocolID() ProtocolPacketID { return protocolSPlayerDigging }
func (p *SPacketPlayerDigging) Type() PacketType             { return SPlayerDigging }
func (p *SPacketPlayerDigging) Pull(reader *buffer.Buffer) error {
	p.Status = reader.PullVarInt()
	p.Position.Pull(reader)
	p.Face = reader.PullByte()
	return nil
}

type SPacketEntityAction struct {
	EntityID, ActionID, JumpBoost int32
}

func (p *SPacketEntityAction) ProtocolID() ProtocolPacketID { return protocolSEntityAction }
func (p *SPacketEntityAction) Type() PacketType             { return SEntityAction }
func (p *SPacketEntityAction) Pull(reader *buffer.Buffer) error {
	p.EntityID = reader.PullVarInt()
	p.ActionID = reader.PullVarInt()
	p.JumpBoost = reader.PullVarInt()
	return nil
}

type SPacketSteerVehicle struct{}

func (p *SPacketSteerVehicle) ProtocolID() ProtocolPacketID { return protocolSSteerVehicle }
func (p *SPacketSteerVehicle) Type() PacketType             { return SSteerVehicle }
func (p *SPacketSteerVehicle) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSetDisplayedRecipe struct{}

func (p *SPacketSetDisplayedRecipe) ProtocolID() ProtocolPacketID { return protocolSSetDisplayedRecipe }
func (p *SPacketSetDisplayedRecipe) Type() PacketType             { return SSetDisplayedRecipe }
func (p *SPacketSetDisplayedRecipe) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSetRecipeBookState struct{}

func (p *SPacketSetRecipeBookState) ProtocolID() ProtocolPacketID { return protocolSSetRecipeBookState }
func (p *SPacketSetRecipeBookState) Type() PacketType             { return SSetRecipeBookState }
func (p *SPacketSetRecipeBookState) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketNameItem struct{}

func (p *SPacketNameItem) ProtocolID() ProtocolPacketID { return protocolSNameItem }
func (p *SPacketNameItem) Type() PacketType             { return SNameItem }
func (p *SPacketNameItem) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketResourcePackStatus struct{}

func (p *SPacketResourcePackStatus) ProtocolID() ProtocolPacketID { return protocolSResourcePackStatus }
func (p *SPacketResourcePackStatus) Type() PacketType             { return SResourcePackStatus }
func (p *SPacketResourcePackStatus) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketAdvancementTab struct{}

func (p *SPacketAdvancementTab) ProtocolID() ProtocolPacketID { return protocolSAdvancementTab }
func (p *SPacketAdvancementTab) Type() PacketType             { return SAdvancementTab }
func (p *SPacketAdvancementTab) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSelectTrade struct{}

func (p *SPacketSelectTrade) ProtocolID() ProtocolPacketID { return protocolSSelectTrade }
func (p *SPacketSelectTrade) Type() PacketType             { return SSelectTrade }
func (p *SPacketSelectTrade) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketSetBeaconEffect struct{}

func (p *SPacketSetBeaconEffect) ProtocolID() ProtocolPacketID { return protocolSSetBeaconEffect }
func (p *SPacketSetBeaconEffect) Type() PacketType             { return SSetBeaconEffect }
func (p *SPacketSetBeaconEffect) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketHeldItemChange struct {
	Slot uint8
}

func (p *SPacketHeldItemChange) ProtocolID() ProtocolPacketID { return protocolSHeldItemChange }
func (p *SPacketHeldItemChange) Type() PacketType             { return SHeldItemChange }
func (p *SPacketHeldItemChange) Pull(reader *buffer.Buffer) error {
	p.Slot = uint8(reader.PullInt16())
	return nil
}

type SPacketUpdateCommandBlock struct{}

func (p *SPacketUpdateCommandBlock) ProtocolID() ProtocolPacketID { return protocolSUpdateCommandBlock }
func (p *SPacketUpdateCommandBlock) Type() PacketType             { return SUpdateCommandBlock }
func (p *SPacketUpdateCommandBlock) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketUpdateCommandBlockMinecart struct{}

func (p *SPacketUpdateCommandBlockMinecart) ProtocolID() ProtocolPacketID {
	return protocolSUpdateCommandBlockMinecart
}
func (p *SPacketUpdateCommandBlockMinecart) Type() PacketType { return SUpdateCommandBlockMinecart }
func (p *SPacketUpdateCommandBlockMinecart) Pull(reader *buffer.Buffer) {
	panic("packet not implemented")
}

type SPacketCreativeInventoryAction struct{}

func (p *SPacketCreativeInventoryAction) ProtocolID() ProtocolPacketID {
	return protocolSCreativeInventoryAction
}
func (p *SPacketCreativeInventoryAction) Type() PacketType           { return SCreativeInventoryAction }
func (p *SPacketCreativeInventoryAction) Pull(reader *buffer.Buffer) { panic("packet not implemented") }

type SPacketUpdateJigsawBlock struct{}

func (p *SPacketUpdateJigsawBlock) ProtocolID() ProtocolPacketID { return protocolSUpdateJigsawBlock }
func (p *SPacketUpdateJigsawBlock) Type() PacketType             { return SUpdateJigsawBlock }
func (p *SPacketUpdateJigsawBlock) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketUpdateStructureBlock struct{}

func (p *SPacketUpdateStructureBlock) ProtocolID() ProtocolPacketID {
	return protocolSUpdateStructureBlock
}
func (p *SPacketUpdateStructureBlock) Type() PacketType           { return SUpdateStructureBlock }
func (p *SPacketUpdateStructureBlock) Pull(reader *buffer.Buffer) { panic("packet not implemented") }

type SPacketUpdateSign struct{}

func (p *SPacketUpdateSign) ProtocolID() ProtocolPacketID { return protocolSUpdateSign }
func (p *SPacketUpdateSign) Type() PacketType             { return SUpdateSign }
func (p *SPacketUpdateSign) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketAnimation struct {
	Hand uint8
}

func (p *SPacketAnimation) ProtocolID() ProtocolPacketID { return protocolSAnimation }
func (p *SPacketAnimation) Type() PacketType             { return SAnimation }
func (p *SPacketAnimation) Pull(reader *buffer.Buffer) error {
	p.Hand = uint8(reader.PullVarInt())
	return nil
}

type SPacketSpectate struct{}

func (p *SPacketSpectate) ProtocolID() ProtocolPacketID { return protocolSSpectate }
func (p *SPacketSpectate) Type() PacketType             { return SSpectate }
func (p *SPacketSpectate) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }

type SPacketPlayerBlockPlacement struct{}

func (p *SPacketPlayerBlockPlacement) ProtocolID() ProtocolPacketID {
	return protocolSPlayerBlockPlacement
}
func (p *SPacketPlayerBlockPlacement) Type() PacketType           { return SPlayerBlockPlacement }
func (p *SPacketPlayerBlockPlacement) Pull(reader *buffer.Buffer) { panic("packet not implemented") }

type SPacketUseItem struct{}

func (p *SPacketUseItem) ProtocolID() ProtocolPacketID { return protocolSUseItem }
func (p *SPacketUseItem) Type() PacketType             { return SUseItem }
func (p *SPacketUseItem) Pull(reader *buffer.Buffer)   { panic("packet not implemented") }
