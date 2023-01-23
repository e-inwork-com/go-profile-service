// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.12
// source: profiles.proto

package profiles

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

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{0}
}

func (x *Profile) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfileEntry *Profile `protobuf:"bytes,1,opt,name=profileEntry,proto3" json:"profileEntry,omitempty"`
}

func (x *ProfileRequest) Reset() {
	*x = ProfileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileRequest) ProtoMessage() {}

func (x *ProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileRequest.ProtoReflect.Descriptor instead.
func (*ProfileRequest) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{1}
}

func (x *ProfileRequest) GetProfileEntry() *Profile {
	if x != nil {
		return x.ProfileEntry
	}
	return nil
}

type ProfileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *ProfileResponse) Reset() {
	*x = ProfileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profiles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileResponse) ProtoMessage() {}

func (x *ProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_profiles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileResponse.ProtoReflect.Descriptor instead.
func (*ProfileResponse) Descriptor() ([]byte, []int) {
	return file_profiles_proto_rawDescGZIP(), []int{2}
}

func (x *ProfileResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_profiles_proto protoreflect.FileDescriptor

var file_profiles_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x19, 0x0a, 0x07, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x47, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x52, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x29,
	0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x55, 0x0a, 0x0e, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0c, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_profiles_proto_rawDescOnce sync.Once
	file_profiles_proto_rawDescData = file_profiles_proto_rawDesc
)

func file_profiles_proto_rawDescGZIP() []byte {
	file_profiles_proto_rawDescOnce.Do(func() {
		file_profiles_proto_rawDescData = protoimpl.X.CompressGZIP(file_profiles_proto_rawDescData)
	})
	return file_profiles_proto_rawDescData
}

var file_profiles_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_profiles_proto_goTypes = []interface{}{
	(*Profile)(nil),         // 0: profiles.Profile
	(*ProfileRequest)(nil),  // 1: profiles.ProfileRequest
	(*ProfileResponse)(nil), // 2: profiles.ProfileResponse
}
var file_profiles_proto_depIdxs = []int32{
	0, // 0: profiles.ProfileRequest.profileEntry:type_name -> profiles.Profile
	1, // 1: profiles.ProfileService.WriteProfile:input_type -> profiles.ProfileRequest
	2, // 2: profiles.ProfileService.WriteProfile:output_type -> profiles.ProfileResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_profiles_proto_init() }
func file_profiles_proto_init() {
	if File_profiles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_profiles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
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
		file_profiles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileRequest); i {
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
		file_profiles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileResponse); i {
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
			RawDescriptor: file_profiles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_profiles_proto_goTypes,
		DependencyIndexes: file_profiles_proto_depIdxs,
		MessageInfos:      file_profiles_proto_msgTypes,
	}.Build()
	File_profiles_proto = out.File
	file_profiles_proto_rawDesc = nil
	file_profiles_proto_goTypes = nil
	file_profiles_proto_depIdxs = nil
}