// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.15.8
// source: common.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ConnState int32

const (
	ConnState_HANDSHAKE ConnState = 0
	ConnState_STATUS    ConnState = 1
	ConnState_LOGIN     ConnState = 2
	ConnState_PLAY      ConnState = 3
)

// Enum value maps for ConnState.
var (
	ConnState_name = map[int32]string{
		0: "HANDSHAKE",
		1: "STATUS",
		2: "LOGIN",
		3: "PLAY",
	}
	ConnState_value = map[string]int32{
		"HANDSHAKE": 0,
		"STATUS":    1,
		"LOGIN":     2,
		"PLAY":      3,
	}
)

func (x ConnState) Enum() *ConnState {
	p := new(ConnState)
	*p = x
	return p
}

func (x ConnState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConnState) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[0].Descriptor()
}

func (ConnState) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[0]
}

func (x ConnState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConnState.Descriptor instead.
func (ConnState) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meta map[string]string `protobuf:"bytes,1,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Types that are assignable to Message:
	//	*Envelope_Cpacket
	//	*Envelope_Spacket
	//	*Envelope_CloseConn
	//	*Envelope_PlayerLoading
	//	*Envelope_PlayerJoined
	//	*Envelope_PlayerSpatial
	//	*Envelope_NewPlayer
	//	*Envelope_PlayerLeft
	Message isEnvelope_Message `protobuf_oneof:"message"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetMeta() map[string]string {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (m *Envelope) GetMessage() isEnvelope_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Envelope) GetCpacket() *CPacket {
	if x, ok := x.GetMessage().(*Envelope_Cpacket); ok {
		return x.Cpacket
	}
	return nil
}

func (x *Envelope) GetSpacket() *SPacket {
	if x, ok := x.GetMessage().(*Envelope_Spacket); ok {
		return x.Spacket
	}
	return nil
}

func (x *Envelope) GetCloseConn() *CloseConn {
	if x, ok := x.GetMessage().(*Envelope_CloseConn); ok {
		return x.CloseConn
	}
	return nil
}

func (x *Envelope) GetPlayerLoading() *PlayerLoading {
	if x, ok := x.GetMessage().(*Envelope_PlayerLoading); ok {
		return x.PlayerLoading
	}
	return nil
}

func (x *Envelope) GetPlayerJoined() *PlayerJoined {
	if x, ok := x.GetMessage().(*Envelope_PlayerJoined); ok {
		return x.PlayerJoined
	}
	return nil
}

func (x *Envelope) GetPlayerSpatial() *PlayerSpatialUpdate {
	if x, ok := x.GetMessage().(*Envelope_PlayerSpatial); ok {
		return x.PlayerSpatial
	}
	return nil
}

func (x *Envelope) GetNewPlayer() *NewPlayerJoined {
	if x, ok := x.GetMessage().(*Envelope_NewPlayer); ok {
		return x.NewPlayer
	}
	return nil
}

func (x *Envelope) GetPlayerLeft() *PlayerLeft {
	if x, ok := x.GetMessage().(*Envelope_PlayerLeft); ok {
		return x.PlayerLeft
	}
	return nil
}

type isEnvelope_Message interface {
	isEnvelope_Message()
}

type Envelope_Cpacket struct {
	Cpacket *CPacket `protobuf:"bytes,2,opt,name=cpacket,proto3,oneof"`
}

type Envelope_Spacket struct {
	Spacket *SPacket `protobuf:"bytes,3,opt,name=spacket,proto3,oneof"`
}

type Envelope_CloseConn struct {
	CloseConn *CloseConn `protobuf:"bytes,4,opt,name=close_conn,json=closeConn,proto3,oneof"`
}

type Envelope_PlayerLoading struct {
	PlayerLoading *PlayerLoading `protobuf:"bytes,5,opt,name=player_loading,json=playerLoading,proto3,oneof"`
}

type Envelope_PlayerJoined struct {
	PlayerJoined *PlayerJoined `protobuf:"bytes,6,opt,name=player_joined,json=playerJoined,proto3,oneof"`
}

type Envelope_PlayerSpatial struct {
	PlayerSpatial *PlayerSpatialUpdate `protobuf:"bytes,7,opt,name=player_spatial,json=playerSpatial,proto3,oneof"`
}

type Envelope_NewPlayer struct {
	NewPlayer *NewPlayerJoined `protobuf:"bytes,8,opt,name=new_player,json=newPlayer,proto3,oneof"`
}

type Envelope_PlayerLeft struct {
	PlayerLeft *PlayerLeft `protobuf:"bytes,9,opt,name=player_left,json=playerLeft,proto3,oneof"`
}

func (*Envelope_Cpacket) isEnvelope_Message() {}

func (*Envelope_Spacket) isEnvelope_Message() {}

func (*Envelope_CloseConn) isEnvelope_Message() {}

func (*Envelope_PlayerLoading) isEnvelope_Message() {}

func (*Envelope_PlayerJoined) isEnvelope_Message() {}

func (*Envelope_PlayerSpatial) isEnvelope_Message() {}

func (*Envelope_NewPlayer) isEnvelope_Message() {}

func (*Envelope_PlayerLeft) isEnvelope_Message() {}

type CPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bytes      []byte `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
	PacketType int32  `protobuf:"varint,2,opt,name=packet_type,json=packetType,proto3" json:"packet_type,omitempty"`
}

func (x *CPacket) Reset() {
	*x = CPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPacket) ProtoMessage() {}

func (x *CPacket) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPacket.ProtoReflect.Descriptor instead.
func (*CPacket) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *CPacket) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *CPacket) GetPacketType() int32 {
	if x != nil {
		return x.PacketType
	}
	return 0
}

type SPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bytes []byte    `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
	State ConnState `protobuf:"varint,2,opt,name=state,proto3,enum=grpc.health.v1.ConnState" json:"state,omitempty"`
}

func (x *SPacket) Reset() {
	*x = SPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SPacket) ProtoMessage() {}

func (x *SPacket) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SPacket.ProtoReflect.Descriptor instead.
func (*SPacket) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

func (x *SPacket) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *SPacket) GetState() ConnState {
	if x != nil {
		return x.State
	}
	return ConnState_HANDSHAKE
}

type CloseConn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnId string    `protobuf:"bytes,1,opt,name=conn_id,json=connId,proto3" json:"conn_id,omitempty"`
	State  ConnState `protobuf:"varint,2,opt,name=state,proto3,enum=grpc.health.v1.ConnState" json:"state,omitempty"`
}

func (x *CloseConn) Reset() {
	*x = CloseConn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseConn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseConn) ProtoMessage() {}

func (x *CloseConn) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseConn.ProtoReflect.Descriptor instead.
func (*CloseConn) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

func (x *CloseConn) GetConnId() string {
	if x != nil {
		return x.ConnId
	}
	return ""
}

func (x *CloseConn) GetState() ConnState {
	if x != nil {
		return x.State
	}
	return ConnState_HANDSHAKE
}

// A player is about to join, but not spawned yet.
type PlayerLoading struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnId    string `protobuf:"bytes,1,opt,name=conn_id,json=connId,proto3" json:"conn_id,omitempty"`
	ProfileId string `protobuf:"bytes,2,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	Username  string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"` // TODO also publish skin data
}

func (x *PlayerLoading) Reset() {
	*x = PlayerLoading{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerLoading) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerLoading) ProtoMessage() {}

func (x *PlayerLoading) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerLoading.ProtoReflect.Descriptor instead.
func (*PlayerLoading) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{4}
}

func (x *PlayerLoading) GetConnId() string {
	if x != nil {
		return x.ConnId
	}
	return ""
}

func (x *PlayerLoading) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

func (x *PlayerLoading) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// A new player joined for the first time.
type NewPlayerJoined struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId  string    `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ConnId    string    `protobuf:"bytes,2,opt,name=conn_id,json=connId,proto3" json:"conn_id,omitempty"`
	ProfileId string    `protobuf:"bytes,3,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	Username  string    `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Pos       *Position `protobuf:"bytes,5,opt,name=pos,proto3" json:"pos,omitempty"`
}

func (x *NewPlayerJoined) Reset() {
	*x = NewPlayerJoined{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewPlayerJoined) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewPlayerJoined) ProtoMessage() {}

func (x *NewPlayerJoined) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewPlayerJoined.ProtoReflect.Descriptor instead.
func (*NewPlayerJoined) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{5}
}

func (x *NewPlayerJoined) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *NewPlayerJoined) GetConnId() string {
	if x != nil {
		return x.ConnId
	}
	return ""
}

func (x *NewPlayerJoined) GetProfileId() string {
	if x != nil {
		return x.ProfileId
	}
	return ""
}

func (x *NewPlayerJoined) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *NewPlayerJoined) GetPos() *Position {
	if x != nil {
		return x.Pos
	}
	return nil
}

// A player has fully joined and spawned.
type PlayerJoined struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ConnId   string `protobuf:"bytes,2,opt,name=conn_id,json=connId,proto3" json:"conn_id,omitempty"`
}

func (x *PlayerJoined) Reset() {
	*x = PlayerJoined{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerJoined) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerJoined) ProtoMessage() {}

func (x *PlayerJoined) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerJoined.ProtoReflect.Descriptor instead.
func (*PlayerJoined) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{6}
}

func (x *PlayerJoined) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *PlayerJoined) GetConnId() string {
	if x != nil {
		return x.ConnId
	}
	return ""
}

// A player has left and disconnected.
type PlayerLeft struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
}

func (x *PlayerLeft) Reset() {
	*x = PlayerLeft{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerLeft) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerLeft) ProtoMessage() {}

func (x *PlayerLeft) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerLeft.ProtoReflect.Descriptor instead.
func (*PlayerLeft) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{7}
}

func (x *PlayerLeft) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

type Position struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z float64 `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
}

func (x *Position) Reset() {
	*x = Position{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Position) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Position) ProtoMessage() {}

func (x *Position) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Position.ProtoReflect.Descriptor instead.
func (*Position) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{8}
}

func (x *Position) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Position) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Position) GetZ() float64 {
	if x != nil {
		return x.Z
	}
	return 0
}

type Rotation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Yaw   float32 `protobuf:"fixed32,1,opt,name=yaw,proto3" json:"yaw,omitempty"`
	Pitch float32 `protobuf:"fixed32,2,opt,name=pitch,proto3" json:"pitch,omitempty"`
}

func (x *Rotation) Reset() {
	*x = Rotation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rotation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rotation) ProtoMessage() {}

func (x *Rotation) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rotation.ProtoReflect.Descriptor instead.
func (*Rotation) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{9}
}

func (x *Rotation) GetYaw() float32 {
	if x != nil {
		return x.Yaw
	}
	return 0
}

func (x *Rotation) GetPitch() float32 {
	if x != nil {
		return x.Pitch
	}
	return 0
}

// Updates position of the player
type PlayerSpatialUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string    `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Pos      *Position `protobuf:"bytes,2,opt,name=pos,proto3" json:"pos,omitempty"`
	Rot      *Rotation `protobuf:"bytes,3,opt,name=rot,proto3" json:"rot,omitempty"`
	OnGround bool      `protobuf:"varint,4,opt,name=on_ground,json=onGround,proto3" json:"on_ground,omitempty"`
}

func (x *PlayerSpatialUpdate) Reset() {
	*x = PlayerSpatialUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerSpatialUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerSpatialUpdate) ProtoMessage() {}

func (x *PlayerSpatialUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerSpatialUpdate.ProtoReflect.Descriptor instead.
func (*PlayerSpatialUpdate) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{10}
}

func (x *PlayerSpatialUpdate) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *PlayerSpatialUpdate) GetPos() *Position {
	if x != nil {
		return x.Pos
	}
	return nil
}

func (x *PlayerSpatialUpdate) GetRot() *Rotation {
	if x != nil {
		return x.Rot
	}
	return nil
}

func (x *PlayerSpatialUpdate) GetOnGround() bool {
	if x != nil {
		return x.OnGround
	}
	return false
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x22, 0x88,
	0x05, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x36, 0x0a, 0x04, 0x6d,
	0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x6d,
	0x65, 0x74, 0x61, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52,
	0x07, 0x63, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x73, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x48, 0x00, 0x52, 0x07, 0x73, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x3a, 0x0a,
	0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x48, 0x00, 0x52, 0x09,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x46, 0x0a, 0x0e, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67,
	0x48, 0x00, 0x52, 0x0d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x43, 0x0a, 0x0d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6a, 0x6f, 0x69, 0x6e,
	0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x4c, 0x0a, 0x0e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x5f, 0x73, 0x70, 0x61, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x70, 0x61, 0x74, 0x69, 0x61, 0x6c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x70, 0x61,
	0x74, 0x69, 0x61, 0x6c, 0x12, 0x40, 0x0a, 0x0a, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x77, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x65, 0x77,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x5f, 0x6c, 0x65, 0x66, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x4c, 0x65, 0x66, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x4c, 0x65, 0x66, 0x74, 0x1a, 0x37, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x09,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x40, 0x0a, 0x07, 0x43, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0x50, 0x0a, 0x07, 0x53,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x55, 0x0a,
	0x09, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e,
	0x6e, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x19, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x22, 0x63, 0x0a, 0x0d, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x6e, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xae, 0x01, 0x0a, 0x0f, 0x4e, 0x65,
	0x77, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e,
	0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2a,
	0x0a, 0x03, 0x70, 0x6f, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x70, 0x6f, 0x73, 0x22, 0x44, 0x0a, 0x0c, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x6e, 0x49, 0x64,
	0x22, 0x29, 0x0a, 0x0a, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x65, 0x66, 0x74, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x34, 0x0a, 0x08, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x01, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x7a, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x01,
	0x7a, 0x22, 0x32, 0x0a, 0x08, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x79, 0x61, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x79, 0x61, 0x77, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x69, 0x74, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05,
	0x70, 0x69, 0x74, 0x63, 0x68, 0x22, 0xa7, 0x01, 0x0a, 0x13, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x53, 0x70, 0x61, 0x74, 0x69, 0x61, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x03, 0x70, 0x6f,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x03, 0x70, 0x6f, 0x73, 0x12, 0x2a, 0x0a, 0x03, 0x72, 0x6f, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x72,
	0x6f, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x2a,
	0x3b, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0d, 0x0a, 0x09,
	0x48, 0x41, 0x4e, 0x44, 0x53, 0x48, 0x41, 0x4b, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x4f, 0x47, 0x49, 0x4e,
	0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4c, 0x41, 0x59, 0x10, 0x03, 0x42, 0x2d, 0x5a, 0x2b,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x65, 0x78, 0x79,
	0x6b, 0x6f, 0x74, 0x2f, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_common_proto_goTypes = []interface{}{
	(ConnState)(0),              // 0: grpc.health.v1.ConnState
	(*Envelope)(nil),            // 1: grpc.health.v1.Envelope
	(*CPacket)(nil),             // 2: grpc.health.v1.CPacket
	(*SPacket)(nil),             // 3: grpc.health.v1.SPacket
	(*CloseConn)(nil),           // 4: grpc.health.v1.CloseConn
	(*PlayerLoading)(nil),       // 5: grpc.health.v1.PlayerLoading
	(*NewPlayerJoined)(nil),     // 6: grpc.health.v1.NewPlayerJoined
	(*PlayerJoined)(nil),        // 7: grpc.health.v1.PlayerJoined
	(*PlayerLeft)(nil),          // 8: grpc.health.v1.PlayerLeft
	(*Position)(nil),            // 9: grpc.health.v1.Position
	(*Rotation)(nil),            // 10: grpc.health.v1.Rotation
	(*PlayerSpatialUpdate)(nil), // 11: grpc.health.v1.PlayerSpatialUpdate
	nil,                         // 12: grpc.health.v1.Envelope.MetaEntry
}
var file_common_proto_depIdxs = []int32{
	12, // 0: grpc.health.v1.Envelope.meta:type_name -> grpc.health.v1.Envelope.MetaEntry
	2,  // 1: grpc.health.v1.Envelope.cpacket:type_name -> grpc.health.v1.CPacket
	3,  // 2: grpc.health.v1.Envelope.spacket:type_name -> grpc.health.v1.SPacket
	4,  // 3: grpc.health.v1.Envelope.close_conn:type_name -> grpc.health.v1.CloseConn
	5,  // 4: grpc.health.v1.Envelope.player_loading:type_name -> grpc.health.v1.PlayerLoading
	7,  // 5: grpc.health.v1.Envelope.player_joined:type_name -> grpc.health.v1.PlayerJoined
	11, // 6: grpc.health.v1.Envelope.player_spatial:type_name -> grpc.health.v1.PlayerSpatialUpdate
	6,  // 7: grpc.health.v1.Envelope.new_player:type_name -> grpc.health.v1.NewPlayerJoined
	8,  // 8: grpc.health.v1.Envelope.player_left:type_name -> grpc.health.v1.PlayerLeft
	0,  // 9: grpc.health.v1.SPacket.state:type_name -> grpc.health.v1.ConnState
	0,  // 10: grpc.health.v1.CloseConn.state:type_name -> grpc.health.v1.ConnState
	9,  // 11: grpc.health.v1.NewPlayerJoined.pos:type_name -> grpc.health.v1.Position
	9,  // 12: grpc.health.v1.PlayerSpatialUpdate.pos:type_name -> grpc.health.v1.Position
	10, // 13: grpc.health.v1.PlayerSpatialUpdate.rot:type_name -> grpc.health.v1.Rotation
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPacket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SPacket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloseConn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerLoading); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewPlayerJoined); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerJoined); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerLeft); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Position); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rotation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerSpatialUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_common_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Envelope_Cpacket)(nil),
		(*Envelope_Spacket)(nil),
		(*Envelope_CloseConn)(nil),
		(*Envelope_PlayerLoading)(nil),
		(*Envelope_PlayerJoined)(nil),
		(*Envelope_PlayerSpatial)(nil),
		(*Envelope_NewPlayer)(nil),
		(*Envelope_PlayerLeft)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		EnumInfos:         file_common_proto_enumTypes,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
