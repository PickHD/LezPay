// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/v1/wallet/wallet.proto

package __

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

type WalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId uint64 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *WalletRequest) Reset() {
	*x = WalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_wallet_wallet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WalletRequest) ProtoMessage() {}

func (x *WalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_wallet_wallet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WalletRequest.ProtoReflect.Descriptor instead.
func (*WalletRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_wallet_wallet_proto_rawDescGZIP(), []int{0}
}

func (x *WalletRequest) GetCustomerId() uint64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type WalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *WalletResponse) Reset() {
	*x = WalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_wallet_wallet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WalletResponse) ProtoMessage() {}

func (x *WalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_wallet_wallet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WalletResponse.ProtoReflect.Descriptor instead.
func (*WalletResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_wallet_wallet_proto_rawDescGZIP(), []int{1}
}

func (x *WalletResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetCustomerWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId uint64 `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *GetCustomerWalletRequest) Reset() {
	*x = GetCustomerWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_wallet_wallet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCustomerWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCustomerWalletRequest) ProtoMessage() {}

func (x *GetCustomerWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_wallet_wallet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCustomerWalletRequest.ProtoReflect.Descriptor instead.
func (*GetCustomerWalletRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_wallet_wallet_proto_rawDescGZIP(), []int{2}
}

func (x *GetCustomerWalletRequest) GetCustomerId() uint64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

type GetCustomerWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Balance int64  `protobuf:"varint,2,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (x *GetCustomerWalletResponse) Reset() {
	*x = GetCustomerWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_wallet_wallet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCustomerWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCustomerWalletResponse) ProtoMessage() {}

func (x *GetCustomerWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_wallet_wallet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCustomerWalletResponse.ProtoReflect.Descriptor instead.
func (*GetCustomerWalletResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_wallet_wallet_proto_rawDescGZIP(), []int{3}
}

func (x *GetCustomerWalletResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetCustomerWalletResponse) GetBalance() int64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

var File_proto_v1_wallet_wallet_proto protoreflect.FileDescriptor

var file_proto_v1_wallet_wallet_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x22,
	0x30, 0x0a, 0x0d, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x20, 0x0a, 0x0e, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x3b, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x45, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x57,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x32, 0xcc, 0x01, 0x0a, 0x0d, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6a, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x12,
	0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x57, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_wallet_wallet_proto_rawDescOnce sync.Once
	file_proto_v1_wallet_wallet_proto_rawDescData = file_proto_v1_wallet_wallet_proto_rawDesc
)

func file_proto_v1_wallet_wallet_proto_rawDescGZIP() []byte {
	file_proto_v1_wallet_wallet_proto_rawDescOnce.Do(func() {
		file_proto_v1_wallet_wallet_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_wallet_wallet_proto_rawDescData)
	})
	return file_proto_v1_wallet_wallet_proto_rawDescData
}

var file_proto_v1_wallet_wallet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_v1_wallet_wallet_proto_goTypes = []interface{}{
	(*WalletRequest)(nil),             // 0: proto.v1.wallet.WalletRequest
	(*WalletResponse)(nil),            // 1: proto.v1.wallet.WalletResponse
	(*GetCustomerWalletRequest)(nil),  // 2: proto.v1.wallet.GetCustomerWalletRequest
	(*GetCustomerWalletResponse)(nil), // 3: proto.v1.wallet.GetCustomerWalletResponse
}
var file_proto_v1_wallet_wallet_proto_depIdxs = []int32{
	0, // 0: proto.v1.wallet.WalletService.CreateWallet:input_type -> proto.v1.wallet.WalletRequest
	2, // 1: proto.v1.wallet.WalletService.GetCustomerWallet:input_type -> proto.v1.wallet.GetCustomerWalletRequest
	1, // 2: proto.v1.wallet.WalletService.CreateWallet:output_type -> proto.v1.wallet.WalletResponse
	3, // 3: proto.v1.wallet.WalletService.GetCustomerWallet:output_type -> proto.v1.wallet.GetCustomerWalletResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_v1_wallet_wallet_proto_init() }
func file_proto_v1_wallet_wallet_proto_init() {
	if File_proto_v1_wallet_wallet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_wallet_wallet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WalletRequest); i {
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
		file_proto_v1_wallet_wallet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WalletResponse); i {
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
		file_proto_v1_wallet_wallet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCustomerWalletRequest); i {
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
		file_proto_v1_wallet_wallet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCustomerWalletResponse); i {
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
			RawDescriptor: file_proto_v1_wallet_wallet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_wallet_wallet_proto_goTypes,
		DependencyIndexes: file_proto_v1_wallet_wallet_proto_depIdxs,
		MessageInfos:      file_proto_v1_wallet_wallet_proto_msgTypes,
	}.Build()
	File_proto_v1_wallet_wallet_proto = out.File
	file_proto_v1_wallet_wallet_proto_rawDesc = nil
	file_proto_v1_wallet_wallet_proto_goTypes = nil
	file_proto_v1_wallet_wallet_proto_depIdxs = nil
}