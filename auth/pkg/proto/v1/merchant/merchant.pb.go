// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/v1/merchant/merchant.proto

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

type MerchantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FullName    string `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Email       string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	PhoneNumber string `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Password    string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	IsVerified  bool   `protobuf:"varint,6,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
}

func (x *MerchantRequest) Reset() {
	*x = MerchantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_merchant_merchant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerchantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerchantRequest) ProtoMessage() {}

func (x *MerchantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_merchant_merchant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerchantRequest.ProtoReflect.Descriptor instead.
func (*MerchantRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_merchant_merchant_proto_rawDescGZIP(), []int{0}
}

func (x *MerchantRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *MerchantRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *MerchantRequest) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *MerchantRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *MerchantRequest) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

type MerchantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	IsVerified bool   `protobuf:"varint,2,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
}

func (x *MerchantResponse) Reset() {
	*x = MerchantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_merchant_merchant_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerchantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerchantResponse) ProtoMessage() {}

func (x *MerchantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_merchant_merchant_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerchantResponse.ProtoReflect.Descriptor instead.
func (*MerchantResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_merchant_merchant_proto_rawDescGZIP(), []int{1}
}

func (x *MerchantResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MerchantResponse) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

type UpdateVerifiedMerchantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *UpdateVerifiedMerchantRequest) Reset() {
	*x = UpdateVerifiedMerchantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_merchant_merchant_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateVerifiedMerchantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateVerifiedMerchantRequest) ProtoMessage() {}

func (x *UpdateVerifiedMerchantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_merchant_merchant_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateVerifiedMerchantRequest.ProtoReflect.Descriptor instead.
func (*UpdateVerifiedMerchantRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_merchant_merchant_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateVerifiedMerchantRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type UpdateVerifiedMerchantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsVerified bool `protobuf:"varint,1,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
}

func (x *UpdateVerifiedMerchantResponse) Reset() {
	*x = UpdateVerifiedMerchantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_merchant_merchant_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateVerifiedMerchantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateVerifiedMerchantResponse) ProtoMessage() {}

func (x *UpdateVerifiedMerchantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_merchant_merchant_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateVerifiedMerchantResponse.ProtoReflect.Descriptor instead.
func (*UpdateVerifiedMerchantResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_merchant_merchant_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateVerifiedMerchantResponse) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

var File_proto_v1_merchant_merchant_proto protoreflect.FileDescriptor

var file_proto_v1_merchant_merchant_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x72, 0x63, 0x68,
	0x61, 0x6e, 0x74, 0x2f, 0x6d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x65, 0x72,
	0x63, 0x68, 0x61, 0x6e, 0x74, 0x22, 0xa4, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c,
	0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75,
	0x6c, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x21, 0x0a, 0x0c,
	0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x69,
	0x73, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22, 0x43, 0x0a, 0x10,
	0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x22, 0x35, 0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x41, 0x0a, 0x1e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73,
	0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x32, 0xeb, 0x01, 0x0a, 0x0f,
	0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x59, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e,
	0x74, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x65, 0x72,
	0x63, 0x68, 0x61, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x6d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7d, 0x0a, 0x16, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x63,
	0x68, 0x61, 0x6e, 0x74, 0x12, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e,
	0x6d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x6d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_merchant_merchant_proto_rawDescOnce sync.Once
	file_proto_v1_merchant_merchant_proto_rawDescData = file_proto_v1_merchant_merchant_proto_rawDesc
)

func file_proto_v1_merchant_merchant_proto_rawDescGZIP() []byte {
	file_proto_v1_merchant_merchant_proto_rawDescOnce.Do(func() {
		file_proto_v1_merchant_merchant_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_merchant_merchant_proto_rawDescData)
	})
	return file_proto_v1_merchant_merchant_proto_rawDescData
}

var file_proto_v1_merchant_merchant_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_v1_merchant_merchant_proto_goTypes = []interface{}{
	(*MerchantRequest)(nil),                // 0: proto.v1.merchant.MerchantRequest
	(*MerchantResponse)(nil),               // 1: proto.v1.merchant.MerchantResponse
	(*UpdateVerifiedMerchantRequest)(nil),  // 2: proto.v1.merchant.UpdateVerifiedMerchantRequest
	(*UpdateVerifiedMerchantResponse)(nil), // 3: proto.v1.merchant.UpdateVerifiedMerchantResponse
}
var file_proto_v1_merchant_merchant_proto_depIdxs = []int32{
	0, // 0: proto.v1.merchant.MerchantService.CreateMerchant:input_type -> proto.v1.merchant.MerchantRequest
	2, // 1: proto.v1.merchant.MerchantService.UpdateVerifiedMerchant:input_type -> proto.v1.merchant.UpdateVerifiedMerchantRequest
	1, // 2: proto.v1.merchant.MerchantService.CreateMerchant:output_type -> proto.v1.merchant.MerchantResponse
	3, // 3: proto.v1.merchant.MerchantService.UpdateVerifiedMerchant:output_type -> proto.v1.merchant.UpdateVerifiedMerchantResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_v1_merchant_merchant_proto_init() }
func file_proto_v1_merchant_merchant_proto_init() {
	if File_proto_v1_merchant_merchant_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_merchant_merchant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerchantRequest); i {
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
		file_proto_v1_merchant_merchant_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerchantResponse); i {
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
		file_proto_v1_merchant_merchant_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateVerifiedMerchantRequest); i {
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
		file_proto_v1_merchant_merchant_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateVerifiedMerchantResponse); i {
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
			RawDescriptor: file_proto_v1_merchant_merchant_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_merchant_merchant_proto_goTypes,
		DependencyIndexes: file_proto_v1_merchant_merchant_proto_depIdxs,
		MessageInfos:      file_proto_v1_merchant_merchant_proto_msgTypes,
	}.Build()
	File_proto_v1_merchant_merchant_proto = out.File
	file_proto_v1_merchant_merchant_proto_rawDesc = nil
	file_proto_v1_merchant_merchant_proto_goTypes = nil
	file_proto_v1_merchant_merchant_proto_depIdxs = nil
}