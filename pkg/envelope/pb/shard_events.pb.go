// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.17.3
// source: shard_events.proto

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

type BlockFace int32

const (
	BlockFace_BOTTOM BlockFace = 0 // -Y
	BlockFace_TOP    BlockFace = 1 // +Y
	BlockFace_NORTH  BlockFace = 2 // -Z
	BlockFace_SOUTH  BlockFace = 3 // +Z
	BlockFace_WEST   BlockFace = 4 // -X
	BlockFace_EAST   BlockFace = 5 // +
)

// Enum value maps for BlockFace.
var (
	BlockFace_name = map[int32]string{
		0: "BOTTOM",
		1: "TOP",
		2: "NORTH",
		3: "SOUTH",
		4: "WEST",
		5: "EAST",
	}
	BlockFace_value = map[string]int32{
		"BOTTOM": 0,
		"TOP":    1,
		"NORTH":  2,
		"SOUTH":  3,
		"WEST":   4,
		"EAST":   5,
	}
)

func (x BlockFace) Enum() *BlockFace {
	p := new(BlockFace)
	*p = x
	return p
}

func (x BlockFace) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BlockFace) Descriptor() protoreflect.EnumDescriptor {
	return file_shard_events_proto_enumTypes[0].Descriptor()
}

func (BlockFace) Type() protoreflect.EnumType {
	return &file_shard_events_proto_enumTypes[0]
}

func (x BlockFace) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BlockFace.Descriptor instead.
func (BlockFace) EnumDescriptor() ([]byte, []int) {
	return file_shard_events_proto_rawDescGZIP(), []int{0}
}

type PlayerDigging_Action int32

const (
	PlayerDigging_STARTED_DIGGING           PlayerDigging_Action = 0
	PlayerDigging_CANCELLED_DIGGING         PlayerDigging_Action = 1
	PlayerDigging_FINISHED_DIGGING          PlayerDigging_Action = 2
	PlayerDigging_DROP_ITEM_STACK           PlayerDigging_Action = 3
	PlayerDigging_DROP_ITEM                 PlayerDigging_Action = 4
	PlayerDigging_SHOOT_ARROW_FINISH_EATING PlayerDigging_Action = 5
	PlayerDigging_SWAP_ITEM_IN_HAND         PlayerDigging_Action = 6
)

// Enum value maps for PlayerDigging_Action.
var (
	PlayerDigging_Action_name = map[int32]string{
		0: "STARTED_DIGGING",
		1: "CANCELLED_DIGGING",
		2: "FINISHED_DIGGING",
		3: "DROP_ITEM_STACK",
		4: "DROP_ITEM",
		5: "SHOOT_ARROW_FINISH_EATING",
		6: "SWAP_ITEM_IN_HAND",
	}
	PlayerDigging_Action_value = map[string]int32{
		"STARTED_DIGGING":           0,
		"CANCELLED_DIGGING":         1,
		"FINISHED_DIGGING":          2,
		"DROP_ITEM_STACK":           3,
		"DROP_ITEM":                 4,
		"SHOOT_ARROW_FINISH_EATING": 5,
		"SWAP_ITEM_IN_HAND":         6,
	}
)

func (x PlayerDigging_Action) Enum() *PlayerDigging_Action {
	p := new(PlayerDigging_Action)
	*p = x
	return p
}

func (x PlayerDigging_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PlayerDigging_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_shard_events_proto_enumTypes[1].Descriptor()
}

func (PlayerDigging_Action) Type() protoreflect.EnumType {
	return &file_shard_events_proto_enumTypes[1]
}

func (x PlayerDigging_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PlayerDigging_Action.Descriptor instead.
func (PlayerDigging_Action) EnumDescriptor() ([]byte, []int) {
	return file_shard_events_proto_rawDescGZIP(), []int{1, 0}
}

type ShardEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//	*ShardEvent_PlayerDigging
	Event isShardEvent_Event `protobuf_oneof:"event"`
}

func (x *ShardEvent) Reset() {
	*x = ShardEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shard_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardEvent) ProtoMessage() {}

func (x *ShardEvent) ProtoReflect() protoreflect.Message {
	mi := &file_shard_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardEvent.ProtoReflect.Descriptor instead.
func (*ShardEvent) Descriptor() ([]byte, []int) {
	return file_shard_events_proto_rawDescGZIP(), []int{0}
}

func (m *ShardEvent) GetEvent() isShardEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *ShardEvent) GetPlayerDigging() *PlayerDigging {
	if x, ok := x.GetEvent().(*ShardEvent_PlayerDigging); ok {
		return x.PlayerDigging
	}
	return nil
}

type isShardEvent_Event interface {
	isShardEvent_Event()
}

type ShardEvent_PlayerDigging struct {
	PlayerDigging *PlayerDigging `protobuf:"bytes,1,opt,name=player_digging,json=playerDigging,proto3,oneof"`
}

func (*ShardEvent_PlayerDigging) isShardEvent_Event() {}

// Updates position of the player
type PlayerDigging struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId  string               `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Action    PlayerDigging_Action `protobuf:"varint,2,opt,name=action,proto3,enum=cncraft.PlayerDigging_Action" json:"action,omitempty"`
	Pos       *Position            `protobuf:"bytes,3,opt,name=pos,proto3" json:"pos,omitempty"`
	BlockFace BlockFace            `protobuf:"varint,4,opt,name=block_face,json=blockFace,proto3,enum=cncraft.BlockFace" json:"block_face,omitempty"`
}

func (x *PlayerDigging) Reset() {
	*x = PlayerDigging{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shard_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerDigging) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerDigging) ProtoMessage() {}

func (x *PlayerDigging) ProtoReflect() protoreflect.Message {
	mi := &file_shard_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerDigging.ProtoReflect.Descriptor instead.
func (*PlayerDigging) Descriptor() ([]byte, []int) {
	return file_shard_events_proto_rawDescGZIP(), []int{1}
}

func (x *PlayerDigging) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *PlayerDigging) GetAction() PlayerDigging_Action {
	if x != nil {
		return x.Action
	}
	return PlayerDigging_STARTED_DIGGING
}

func (x *PlayerDigging) GetPos() *Position {
	if x != nil {
		return x.Pos
	}
	return nil
}

func (x *PlayerDigging) GetBlockFace() BlockFace {
	if x != nil {
		return x.BlockFace
	}
	return BlockFace_BOTTOM
}

var File_shard_events_proto protoreflect.FileDescriptor

var file_shard_events_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x68, 0x61, 0x72, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x1a, 0x0c, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x0a, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x3f, 0x0a, 0x0e, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x5f, 0x64, 0x69, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x44, 0x69, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x44, 0x69, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x22, 0xe2, 0x02, 0x0a, 0x0d, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x44, 0x69,
	0x67, 0x67, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x35, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x44, 0x69, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x03, 0x70, 0x6f, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74,
	0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x70, 0x6f, 0x73, 0x12, 0x31,
	0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x66, 0x61, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x12, 0x2e, 0x63, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x46, 0x61, 0x63, 0x65, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x46, 0x61, 0x63,
	0x65, 0x22, 0xa4, 0x01, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x0f,
	0x53, 0x54, 0x41, 0x52, 0x54, 0x45, 0x44, 0x5f, 0x44, 0x49, 0x47, 0x47, 0x49, 0x4e, 0x47, 0x10,
	0x00, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x45, 0x44, 0x5f, 0x44,
	0x49, 0x47, 0x47, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x46, 0x49, 0x4e, 0x49,
	0x53, 0x48, 0x45, 0x44, 0x5f, 0x44, 0x49, 0x47, 0x47, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x13,
	0x0a, 0x0f, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x53, 0x54, 0x41, 0x43,
	0x4b, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x52, 0x4f, 0x50, 0x5f, 0x49, 0x54, 0x45, 0x4d,
	0x10, 0x04, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x48, 0x4f, 0x4f, 0x54, 0x5f, 0x41, 0x52, 0x52, 0x4f,
	0x57, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x5f, 0x45, 0x41, 0x54, 0x49, 0x4e, 0x47, 0x10,
	0x05, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x57, 0x41, 0x50, 0x5f, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x49,
	0x4e, 0x5f, 0x48, 0x41, 0x4e, 0x44, 0x10, 0x06, 0x2a, 0x4a, 0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x46, 0x61, 0x63, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x42, 0x4f, 0x54, 0x54, 0x4f, 0x4d, 0x10,
	0x00, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x4f, 0x50, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4e, 0x4f,
	0x52, 0x54, 0x48, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x4f, 0x55, 0x54, 0x48, 0x10, 0x03,
	0x12, 0x08, 0x0a, 0x04, 0x57, 0x45, 0x53, 0x54, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x45, 0x41,
	0x53, 0x54, 0x10, 0x05, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x65, 0x78, 0x79, 0x6b, 0x6f, 0x74, 0x2f, 0x63, 0x6e, 0x63, 0x72,
	0x61, 0x66, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shard_events_proto_rawDescOnce sync.Once
	file_shard_events_proto_rawDescData = file_shard_events_proto_rawDesc
)

func file_shard_events_proto_rawDescGZIP() []byte {
	file_shard_events_proto_rawDescOnce.Do(func() {
		file_shard_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_shard_events_proto_rawDescData)
	})
	return file_shard_events_proto_rawDescData
}

var file_shard_events_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_shard_events_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_shard_events_proto_goTypes = []interface{}{
	(BlockFace)(0),            // 0: cncraft.BlockFace
	(PlayerDigging_Action)(0), // 1: cncraft.PlayerDigging.Action
	(*ShardEvent)(nil),        // 2: cncraft.ShardEvent
	(*PlayerDigging)(nil),     // 3: cncraft.PlayerDigging
	(*Position)(nil),          // 4: cncraft.Position
}
var file_shard_events_proto_depIdxs = []int32{
	3, // 0: cncraft.ShardEvent.player_digging:type_name -> cncraft.PlayerDigging
	1, // 1: cncraft.PlayerDigging.action:type_name -> cncraft.PlayerDigging.Action
	4, // 2: cncraft.PlayerDigging.pos:type_name -> cncraft.Position
	0, // 3: cncraft.PlayerDigging.block_face:type_name -> cncraft.BlockFace
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_shard_events_proto_init() }
func file_shard_events_proto_init() {
	if File_shard_events_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_shard_events_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardEvent); i {
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
		file_shard_events_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerDigging); i {
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
	file_shard_events_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ShardEvent_PlayerDigging)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shard_events_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shard_events_proto_goTypes,
		DependencyIndexes: file_shard_events_proto_depIdxs,
		EnumInfos:         file_shard_events_proto_enumTypes,
		MessageInfos:      file_shard_events_proto_msgTypes,
	}.Build()
	File_shard_events_proto = out.File
	file_shard_events_proto_rawDesc = nil
	file_shard_events_proto_goTypes = nil
	file_shard_events_proto_depIdxs = nil
}
