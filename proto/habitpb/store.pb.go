// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: store.proto

package habitpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Habits struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Store map[string]*Habits_Habit `protobuf:"bytes,1,rep,name=store,proto3" json:"store,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Habits) Reset() {
	*x = Habits{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Habits) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Habits) ProtoMessage() {}

func (x *Habits) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Habits.ProtoReflect.Descriptor instead.
func (*Habits) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0}
}

func (x *Habits) GetStore() map[string]*Habits_Habit {
	if x != nil {
		return x.Store
	}
	return nil
}

type Habits_Habit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Streak        int32 `protobuf:"varint,1,opt,name=streak,proto3" json:"streak,omitempty"`
	LastPerformed int64 `protobuf:"varint,2,opt,name=last_performed,json=lastPerformed,proto3" json:"last_performed,omitempty"`
}

func (x *Habits_Habit) Reset() {
	*x = Habits_Habit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Habits_Habit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Habits_Habit) ProtoMessage() {}

func (x *Habits_Habit) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Habits_Habit.ProtoReflect.Descriptor instead.
func (*Habits_Habit) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Habits_Habit) GetStreak() int32 {
	if x != nil {
		return x.Streak
	}
	return 0
}

func (x *Habits_Habit) GetLastPerformed() int64 {
	if x != nil {
		return x.LastPerformed
	}
	return 0
}

var File_store_proto protoreflect.FileDescriptor

var file_store_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x01,
	0x0a, 0x06, 0x48, 0x61, 0x62, 0x69, 0x74, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x48, 0x61, 0x62, 0x69, 0x74, 0x73,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x1a, 0x46, 0x0a, 0x05, 0x48, 0x61, 0x62, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6b, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x70, 0x65, 0x72, 0x66,
	0x6f, 0x72, 0x6d, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x6c, 0x61, 0x73,
	0x74, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x64, 0x1a, 0x47, 0x0a, 0x0a, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x48, 0x61, 0x62, 0x69,
	0x74, 0x73, 0x2e, 0x48, 0x61, 0x62, 0x69, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x68, 0x61, 0x62, 0x69, 0x74, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_proto_rawDescOnce sync.Once
	file_store_proto_rawDescData = file_store_proto_rawDesc
)

func file_store_proto_rawDescGZIP() []byte {
	file_store_proto_rawDescOnce.Do(func() {
		file_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_proto_rawDescData)
	})
	return file_store_proto_rawDescData
}

var file_store_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_store_proto_goTypes = []interface{}{
	(*Habits)(nil),       // 0: Habits
	(*Habits_Habit)(nil), // 1: Habits.Habit
	nil,                  // 2: Habits.StoreEntry
}
var file_store_proto_depIdxs = []int32{
	2, // 0: Habits.store:type_name -> Habits.StoreEntry
	1, // 1: Habits.StoreEntry.value:type_name -> Habits.Habit
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_store_proto_init() }
func file_store_proto_init() {
	if File_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Habits); i {
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
		file_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Habits_Habit); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_proto_goTypes,
		DependencyIndexes: file_store_proto_depIdxs,
		MessageInfos:      file_store_proto_msgTypes,
	}.Build()
	File_store_proto = out.File
	file_store_proto_rawDesc = nil
	file_store_proto_goTypes = nil
	file_store_proto_depIdxs = nil
}
