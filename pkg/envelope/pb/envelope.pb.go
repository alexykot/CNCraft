// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.17.3
// source: envelope.proto

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

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meta       map[string]string `protobuf:"bytes,1,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ShardEvent *ShardEvent       `protobuf:"bytes,2,opt,name=shard_event,json=shardEvent,proto3" json:"shard_event,omitempty"`
	// Types that are assignable to Message:
	//	*Envelope_Cpacket
	//	*Envelope_Spacket
	//	*Envelope_CloseConn
	//	*Envelope_PlayerLoading
	//	*Envelope_NewPlayer
	//	*Envelope_PlayerJoined
	//	*Envelope_PlayerLeft
	//	*Envelope_PlayerSpatial
	//	*Envelope_PlayerInventory
	Message isEnvelope_Message `protobuf_oneof:"message"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envelope_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_envelope_proto_msgTypes[0]
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
	return file_envelope_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetMeta() map[string]string {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *Envelope) GetShardEvent() *ShardEvent {
	if x != nil {
		return x.ShardEvent
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

func (x *Envelope) GetNewPlayer() *NewPlayerJoined {
	if x, ok := x.GetMessage().(*Envelope_NewPlayer); ok {
		return x.NewPlayer
	}
	return nil
}

func (x *Envelope) GetPlayerJoined() *PlayerJoined {
	if x, ok := x.GetMessage().(*Envelope_PlayerJoined); ok {
		return x.PlayerJoined
	}
	return nil
}

func (x *Envelope) GetPlayerLeft() *PlayerLeft {
	if x, ok := x.GetMessage().(*Envelope_PlayerLeft); ok {
		return x.PlayerLeft
	}
	return nil
}

func (x *Envelope) GetPlayerSpatial() *PlayerSpatialUpdate {
	if x, ok := x.GetMessage().(*Envelope_PlayerSpatial); ok {
		return x.PlayerSpatial
	}
	return nil
}

func (x *Envelope) GetPlayerInventory() *PlayerInventoryUpdate {
	if x, ok := x.GetMessage().(*Envelope_PlayerInventory); ok {
		return x.PlayerInventory
	}
	return nil
}

type isEnvelope_Message interface {
	isEnvelope_Message()
}

type Envelope_Cpacket struct {
	Cpacket *CPacket `protobuf:"bytes,3,opt,name=cpacket,proto3,oneof"`
}

type Envelope_Spacket struct {
	Spacket *SPacket `protobuf:"bytes,4,opt,name=spacket,proto3,oneof"`
}

type Envelope_CloseConn struct {
	CloseConn *CloseConn `protobuf:"bytes,5,opt,name=close_conn,json=closeConn,proto3,oneof"`
}

type Envelope_PlayerLoading struct {
	PlayerLoading *PlayerLoading `protobuf:"bytes,6,opt,name=player_loading,json=playerLoading,proto3,oneof"`
}

type Envelope_NewPlayer struct {
	NewPlayer *NewPlayerJoined `protobuf:"bytes,7,opt,name=new_player,json=newPlayer,proto3,oneof"`
}

type Envelope_PlayerJoined struct {
	PlayerJoined *PlayerJoined `protobuf:"bytes,8,opt,name=player_joined,json=playerJoined,proto3,oneof"`
}

type Envelope_PlayerLeft struct {
	PlayerLeft *PlayerLeft `protobuf:"bytes,9,opt,name=player_left,json=playerLeft,proto3,oneof"`
}

type Envelope_PlayerSpatial struct {
	PlayerSpatial *PlayerSpatialUpdate `protobuf:"bytes,10,opt,name=player_spatial,json=playerSpatial,proto3,oneof"`
}

type Envelope_PlayerInventory struct {
	PlayerInventory *PlayerInventoryUpdate `protobuf:"bytes,11,opt,name=player_inventory,json=playerInventory,proto3,oneof"`
}

func (*Envelope_Cpacket) isEnvelope_Message() {}

func (*Envelope_Spacket) isEnvelope_Message() {}

func (*Envelope_CloseConn) isEnvelope_Message() {}

func (*Envelope_PlayerLoading) isEnvelope_Message() {}

func (*Envelope_NewPlayer) isEnvelope_Message() {}

func (*Envelope_PlayerJoined) isEnvelope_Message() {}

func (*Envelope_PlayerLeft) isEnvelope_Message() {}

func (*Envelope_PlayerSpatial) isEnvelope_Message() {}

func (*Envelope_PlayerInventory) isEnvelope_Message() {}

type SampleMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to TestOneof:
	//	*SampleMessage_Name
	//	*SampleMessage_SubMessage
	TestOneof isSampleMessage_TestOneof `protobuf_oneof:"test_oneof"`
}

func (x *SampleMessage) Reset() {
	*x = SampleMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envelope_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SampleMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SampleMessage) ProtoMessage() {}

func (x *SampleMessage) ProtoReflect() protoreflect.Message {
	mi := &file_envelope_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SampleMessage.ProtoReflect.Descriptor instead.
func (*SampleMessage) Descriptor() ([]byte, []int) {
	return file_envelope_proto_rawDescGZIP(), []int{1}
}

func (m *SampleMessage) GetTestOneof() isSampleMessage_TestOneof {
	if m != nil {
		return m.TestOneof
	}
	return nil
}

func (x *SampleMessage) GetName() string {
	if x, ok := x.GetTestOneof().(*SampleMessage_Name); ok {
		return x.Name
	}
	return ""
}

func (x *SampleMessage) GetSubMessage() string {
	if x, ok := x.GetTestOneof().(*SampleMessage_SubMessage); ok {
		return x.SubMessage
	}
	return ""
}

type isSampleMessage_TestOneof interface {
	isSampleMessage_TestOneof()
}

type SampleMessage_Name struct {
	Name string `protobuf:"bytes,4,opt,name=name,proto3,oneof"`
}

type SampleMessage_SubMessage struct {
	SubMessage string `protobuf:"bytes,9,opt,name=sub_message,json=subMessage,proto3,oneof"`
}

func (*SampleMessage_Name) isSampleMessage_TestOneof() {}

func (*SampleMessage_SubMessage) isSampleMessage_TestOneof() {}

var File_envelope_proto protoreflect.FileDescriptor

var file_envelope_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x1a, 0x12, 0x73, 0x68, 0x61, 0x72, 0x64,
	0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcc, 0x05,
	0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x2f, 0x0a, 0x04, 0x6d, 0x65,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61,
	0x66, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x34, 0x0a, 0x0b, 0x73,
	0x68, 0x61, 0x72, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x72, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x2c, 0x0a, 0x07, 0x63, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x48, 0x00, 0x52, 0x07, 0x63, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x2c, 0x0a, 0x07, 0x73, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x53, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x48, 0x00, 0x52, 0x07, 0x73, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x33, 0x0a,
	0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x48, 0x00, 0x52, 0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f,
	0x6e, 0x6e, 0x12, 0x3f, 0x0a, 0x0e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6e, 0x63,
	0x72, 0x61, 0x66, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x6f, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x39, 0x0a, 0x0a, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66,
	0x74, 0x2e, 0x4e, 0x65, 0x77, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65,
	0x64, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x65, 0x77, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x3c,
	0x0a, 0x0d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x48, 0x00, 0x52, 0x0c,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x36, 0x0a, 0x0b,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6c, 0x65, 0x66, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x4c, 0x65, 0x66, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x4c, 0x65, 0x66, 0x74, 0x12, 0x45, 0x0a, 0x0e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x73,
	0x70, 0x61, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63,
	0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x70, 0x61,
	0x74, 0x69, 0x61, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x53, 0x70, 0x61, 0x74, 0x69, 0x61, 0x6c, 0x12, 0x4b, 0x0a, 0x10, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x37, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x56, 0x0a, 0x0d,
	0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0b, 0x73, 0x75, 0x62, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x73, 0x75, 0x62, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6f,
	0x6e, 0x65, 0x6f, 0x66, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x65, 0x78, 0x79, 0x6b, 0x6f, 0x74, 0x2f, 0x63, 0x6e, 0x63, 0x72,
	0x61, 0x66, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envelope_proto_rawDescOnce sync.Once
	file_envelope_proto_rawDescData = file_envelope_proto_rawDesc
)

func file_envelope_proto_rawDescGZIP() []byte {
	file_envelope_proto_rawDescOnce.Do(func() {
		file_envelope_proto_rawDescData = protoimpl.X.CompressGZIP(file_envelope_proto_rawDescData)
	})
	return file_envelope_proto_rawDescData
}

var file_envelope_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_envelope_proto_goTypes = []interface{}{
	(*Envelope)(nil),              // 0: cncraft.Envelope
	(*SampleMessage)(nil),         // 1: cncraft.SampleMessage
	nil,                           // 2: cncraft.Envelope.MetaEntry
	(*ShardEvent)(nil),            // 3: cncraft.ShardEvent
	(*CPacket)(nil),               // 4: cncraft.CPacket
	(*SPacket)(nil),               // 5: cncraft.SPacket
	(*CloseConn)(nil),             // 6: cncraft.CloseConn
	(*PlayerLoading)(nil),         // 7: cncraft.PlayerLoading
	(*NewPlayerJoined)(nil),       // 8: cncraft.NewPlayerJoined
	(*PlayerJoined)(nil),          // 9: cncraft.PlayerJoined
	(*PlayerLeft)(nil),            // 10: cncraft.PlayerLeft
	(*PlayerSpatialUpdate)(nil),   // 11: cncraft.PlayerSpatialUpdate
	(*PlayerInventoryUpdate)(nil), // 12: cncraft.PlayerInventoryUpdate
}
var file_envelope_proto_depIdxs = []int32{
	2,  // 0: cncraft.Envelope.meta:type_name -> cncraft.Envelope.MetaEntry
	3,  // 1: cncraft.Envelope.shard_event:type_name -> cncraft.ShardEvent
	4,  // 2: cncraft.Envelope.cpacket:type_name -> cncraft.CPacket
	5,  // 3: cncraft.Envelope.spacket:type_name -> cncraft.SPacket
	6,  // 4: cncraft.Envelope.close_conn:type_name -> cncraft.CloseConn
	7,  // 5: cncraft.Envelope.player_loading:type_name -> cncraft.PlayerLoading
	8,  // 6: cncraft.Envelope.new_player:type_name -> cncraft.NewPlayerJoined
	9,  // 7: cncraft.Envelope.player_joined:type_name -> cncraft.PlayerJoined
	10, // 8: cncraft.Envelope.player_left:type_name -> cncraft.PlayerLeft
	11, // 9: cncraft.Envelope.player_spatial:type_name -> cncraft.PlayerSpatialUpdate
	12, // 10: cncraft.Envelope.player_inventory:type_name -> cncraft.PlayerInventoryUpdate
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_envelope_proto_init() }
func file_envelope_proto_init() {
	if File_envelope_proto != nil {
		return
	}
	file_shard_events_proto_init()
	file_messages_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_envelope_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_envelope_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SampleMessage); i {
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
	file_envelope_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Envelope_Cpacket)(nil),
		(*Envelope_Spacket)(nil),
		(*Envelope_CloseConn)(nil),
		(*Envelope_PlayerLoading)(nil),
		(*Envelope_NewPlayer)(nil),
		(*Envelope_PlayerJoined)(nil),
		(*Envelope_PlayerLeft)(nil),
		(*Envelope_PlayerSpatial)(nil),
		(*Envelope_PlayerInventory)(nil),
	}
	file_envelope_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*SampleMessage_Name)(nil),
		(*SampleMessage_SubMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envelope_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envelope_proto_goTypes,
		DependencyIndexes: file_envelope_proto_depIdxs,
		MessageInfos:      file_envelope_proto_msgTypes,
	}.Build()
	File_envelope_proto = out.File
	file_envelope_proto_rawDesc = nil
	file_envelope_proto_goTypes = nil
	file_envelope_proto_depIdxs = nil
}
