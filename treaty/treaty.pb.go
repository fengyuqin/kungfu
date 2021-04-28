// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: treaty.proto

package treaty

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

type CodeType int32

const (
	CodeType_CodeSuccess CodeType = 0 //成功
	CodeType_CodeFailed  CodeType = 1 //失败
)

// Enum value maps for CodeType.
var (
	CodeType_name = map[int32]string{
		0: "CodeSuccess",
		1: "CodeFailed",
	}
	CodeType_value = map[string]int32{
		"CodeSuccess": 0,
		"CodeFailed":  1,
	}
)

func (x CodeType) Enum() *CodeType {
	p := new(CodeType)
	*p = x
	return p
}

func (x CodeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CodeType) Descriptor() protoreflect.EnumDescriptor {
	return file_treaty_proto_enumTypes[0].Descriptor()
}

func (CodeType) Type() protoreflect.EnumType {
	return &file_treaty_proto_enumTypes[0]
}

func (x CodeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CodeType.Descriptor instead.
func (CodeType) EnumDescriptor() ([]byte, []int) {
	return file_treaty_proto_rawDescGZIP(), []int{0}
}

type MsgId int32

const (
	MsgId_Msg_None MsgId = 0
)

// Enum value maps for MsgId.
var (
	MsgId_name = map[int32]string{
		0: "Msg_None",
	}
	MsgId_value = map[string]int32{
		"Msg_None": 0,
	}
)

func (x MsgId) Enum() *MsgId {
	p := new(MsgId)
	*p = x
	return p
}

func (x MsgId) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgId) Descriptor() protoreflect.EnumDescriptor {
	return file_treaty_proto_enumTypes[1].Descriptor()
}

func (MsgId) Type() protoreflect.EnumType {
	return &file_treaty_proto_enumTypes[1]
}

func (x MsgId) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgId.Descriptor instead.
func (MsgId) EnumDescriptor() ([]byte, []int) {
	return file_treaty_proto_rawDescGZIP(), []int{1}
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerType string `protobuf:"bytes,1,opt,name=server_type,json=serverType,proto3" json:"server_type,omitempty"`  //服务器类型
	ServerId   string `protobuf:"bytes,2,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`        //服务器ID
	ServerName string `protobuf:"bytes,3,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`  //服务器名字
	ServerIp   string `protobuf:"bytes,4,opt,name=server_ip,json=serverIp,proto3" json:"server_ip,omitempty"`        //服务器IP
	ClientPort int32  `protobuf:"varint,5,opt,name=client_port,json=clientPort,proto3" json:"client_port,omitempty"` //客户端端口
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_treaty_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_treaty_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_treaty_proto_rawDescGZIP(), []int{0}
}

func (x *Server) GetServerType() string {
	if x != nil {
		return x.ServerType
	}
	return ""
}

func (x *Server) GetServerId() string {
	if x != nil {
		return x.ServerId
	}
	return ""
}

func (x *Server) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *Server) GetServerIp() string {
	if x != nil {
		return x.ServerIp
	}
	return ""
}

func (x *Server) GetClientPort() int32 {
	if x != nil {
		return x.ClientPort
	}
	return 0
}

type BalanceResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   CodeType `protobuf:"varint,1,opt,name=code,proto3,enum=treaty.CodeType" json:"code,omitempty"` //0-成功 1-失败
	Server *Server  `protobuf:"bytes,2,opt,name=server,proto3" json:"server,omitempty"`
}

func (x *BalanceResult) Reset() {
	*x = BalanceResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_treaty_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BalanceResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceResult) ProtoMessage() {}

func (x *BalanceResult) ProtoReflect() protoreflect.Message {
	mi := &file_treaty_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceResult.ProtoReflect.Descriptor instead.
func (*BalanceResult) Descriptor() ([]byte, []int) {
	return file_treaty_proto_rawDescGZIP(), []int{1}
}

func (x *BalanceResult) GetCode() CodeType {
	if x != nil {
		return x.Code
	}
	return CodeType_CodeSuccess
}

func (x *BalanceResult) GetServer() *Server {
	if x != nil {
		return x.Server
	}
	return nil
}

var File_treaty_proto protoreflect.FileDescriptor

var file_treaty_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x74, 0x72, 0x65, 0x61, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x74, 0x72, 0x65, 0x61, 0x74, 0x79, 0x22, 0xa5, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x70, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x5d,
	0x0a, 0x0d, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x24, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e,
	0x74, 0x72, 0x65, 0x61, 0x74, 0x79, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x72, 0x65, 0x61, 0x74, 0x79, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2a, 0x2b, 0x0a,
	0x08, 0x43, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x6f, 0x64,
	0x65, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x6f,
	0x64, 0x65, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x2a, 0x15, 0x0a, 0x05, 0x4d, 0x73,
	0x67, 0x49, 0x64, 0x12, 0x0c, 0x0a, 0x08, 0x4d, 0x73, 0x67, 0x5f, 0x4e, 0x6f, 0x6e, 0x65, 0x10,
	0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x2e, 0x3b, 0x74, 0x72, 0x65, 0x61, 0x74, 0x79, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_treaty_proto_rawDescOnce sync.Once
	file_treaty_proto_rawDescData = file_treaty_proto_rawDesc
)

func file_treaty_proto_rawDescGZIP() []byte {
	file_treaty_proto_rawDescOnce.Do(func() {
		file_treaty_proto_rawDescData = protoimpl.X.CompressGZIP(file_treaty_proto_rawDescData)
	})
	return file_treaty_proto_rawDescData
}

var file_treaty_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_treaty_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_treaty_proto_goTypes = []interface{}{
	(CodeType)(0),         // 0: treaty.CodeType
	(MsgId)(0),            // 1: treaty.MsgId
	(*Server)(nil),        // 2: treaty.Server
	(*BalanceResult)(nil), // 3: treaty.BalanceResult
}
var file_treaty_proto_depIdxs = []int32{
	0, // 0: treaty.BalanceResult.code:type_name -> treaty.CodeType
	2, // 1: treaty.BalanceResult.server:type_name -> treaty.Server
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_treaty_proto_init() }
func file_treaty_proto_init() {
	if File_treaty_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_treaty_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
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
		file_treaty_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BalanceResult); i {
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
			RawDescriptor: file_treaty_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_treaty_proto_goTypes,
		DependencyIndexes: file_treaty_proto_depIdxs,
		EnumInfos:         file_treaty_proto_enumTypes,
		MessageInfos:      file_treaty_proto_msgTypes,
	}.Build()
	File_treaty_proto = out.File
	file_treaty_proto_rawDesc = nil
	file_treaty_proto_goTypes = nil
	file_treaty_proto_depIdxs = nil
}
