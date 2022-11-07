// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: pbschema/itemfilterRPC/universal.proto

package itemfilterRPC_pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

type RedisFilterType int32

const (
	RedisFilterType_SET    RedisFilterType = 0
	RedisFilterType_BLOOM  RedisFilterType = 1
	RedisFilterType_CUCKOO RedisFilterType = 2
)

// Enum value maps for RedisFilterType.
var (
	RedisFilterType_name = map[int32]string{
		0: "SET",
		1: "BLOOM",
		2: "CUCKOO",
	}
	RedisFilterType_value = map[string]int32{
		"SET":    0,
		"BLOOM":  1,
		"CUCKOO": 2,
	}
)

func (x RedisFilterType) Enum() *RedisFilterType {
	p := new(RedisFilterType)
	*p = x
	return p
}

func (x RedisFilterType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RedisFilterType) Descriptor() protoreflect.EnumDescriptor {
	return file_pbschema_itemfilterRPC_universal_proto_enumTypes[0].Descriptor()
}

func (RedisFilterType) Type() protoreflect.EnumType {
	return &file_pbschema_itemfilterRPC_universal_proto_enumTypes[0]
}

func (x RedisFilterType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RedisFilterType.Descriptor instead.
func (RedisFilterType) EnumDescriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{0}
}

type RedisFilterMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID               string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	EntitySourceType string `protobuf:"bytes,2,opt,name=EntitySourceType,proto3" json:"EntitySourceType,omitempty"`
	Name             string `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	Desc             string `protobuf:"bytes,4,opt,name=Desc,proto3" json:"Desc,omitempty"`
	Key              string `protobuf:"bytes,5,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *RedisFilterMeta) Reset() {
	*x = RedisFilterMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisFilterMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisFilterMeta) ProtoMessage() {}

func (x *RedisFilterMeta) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisFilterMeta.ProtoReflect.Descriptor instead.
func (*RedisFilterMeta) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{0}
}

func (x *RedisFilterMeta) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *RedisFilterMeta) GetEntitySourceType() string {
	if x != nil {
		return x.EntitySourceType
	}
	return ""
}

func (x *RedisFilterMeta) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RedisFilterMeta) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *RedisFilterMeta) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type RedisFilterStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exists     bool  `protobuf:"varint,1,opt,name=Exists,proto3" json:"Exists,omitempty"`
	Size       int64 `protobuf:"varint,2,opt,name=Size,proto3" json:"Size,omitempty"`
	TTLSeconds int64 `protobuf:"varint,3,opt,name=TTLSeconds,proto3" json:"TTLSeconds,omitempty"`
}

func (x *RedisFilterStatus) Reset() {
	*x = RedisFilterStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisFilterStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisFilterStatus) ProtoMessage() {}

func (x *RedisFilterStatus) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisFilterStatus.ProtoReflect.Descriptor instead.
func (*RedisFilterStatus) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{1}
}

func (x *RedisFilterStatus) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

func (x *RedisFilterStatus) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *RedisFilterStatus) GetTTLSeconds() int64 {
	if x != nil {
		return x.TTLSeconds
	}
	return 0
}

type SetFilterSetting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MaxTTLSeconds int64 `protobuf:"varint,1,opt,name=MaxTTLSeconds,proto3" json:"MaxTTLSeconds,omitempty"`
}

func (x *SetFilterSetting) Reset() {
	*x = SetFilterSetting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetFilterSetting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetFilterSetting) ProtoMessage() {}

func (x *SetFilterSetting) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetFilterSetting.ProtoReflect.Descriptor instead.
func (*SetFilterSetting) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{2}
}

func (x *SetFilterSetting) GetMaxTTLSeconds() int64 {
	if x != nil {
		return x.MaxTTLSeconds
	}
	return 0
}

type BloomFilterSetting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Capacity      int64   `protobuf:"varint,1,opt,name=Capacity,proto3" json:"Capacity,omitempty"`
	ErrorRate     float64 `protobuf:"fixed64,2,opt,name=ErrorRate,proto3" json:"ErrorRate,omitempty"`
	Expansion     int64   `protobuf:"varint,3,opt,name=Expansion,proto3" json:"Expansion,omitempty"`
	NonScaling    bool    `protobuf:"varint,4,opt,name=NonScaling,proto3" json:"NonScaling,omitempty"`
	MaxTTLSeconds int64   `protobuf:"varint,5,opt,name=MaxTTLSeconds,proto3" json:"MaxTTLSeconds,omitempty"`
}

func (x *BloomFilterSetting) Reset() {
	*x = BloomFilterSetting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BloomFilterSetting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BloomFilterSetting) ProtoMessage() {}

func (x *BloomFilterSetting) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BloomFilterSetting.ProtoReflect.Descriptor instead.
func (*BloomFilterSetting) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{3}
}

func (x *BloomFilterSetting) GetCapacity() int64 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

func (x *BloomFilterSetting) GetErrorRate() float64 {
	if x != nil {
		return x.ErrorRate
	}
	return 0
}

func (x *BloomFilterSetting) GetExpansion() int64 {
	if x != nil {
		return x.Expansion
	}
	return 0
}

func (x *BloomFilterSetting) GetNonScaling() bool {
	if x != nil {
		return x.NonScaling
	}
	return false
}

func (x *BloomFilterSetting) GetMaxTTLSeconds() int64 {
	if x != nil {
		return x.MaxTTLSeconds
	}
	return 0
}

type CuckooFilterSetting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Capacity      int64 `protobuf:"varint,1,opt,name=Capacity,proto3" json:"Capacity,omitempty"`
	Expansion     int64 `protobuf:"varint,2,opt,name=Expansion,proto3" json:"Expansion,omitempty"`
	NonScaling    bool  `protobuf:"varint,3,opt,name=NonScaling,proto3" json:"NonScaling,omitempty"`
	BucketSize    int64 `protobuf:"varint,4,opt,name=BucketSize,proto3" json:"BucketSize,omitempty"`
	MaxIterations int64 `protobuf:"varint,5,opt,name=MaxIterations,proto3" json:"MaxIterations,omitempty"`
	MaxTTLSeconds int64 `protobuf:"varint,6,opt,name=MaxTTLSeconds,proto3" json:"MaxTTLSeconds,omitempty"`
}

func (x *CuckooFilterSetting) Reset() {
	*x = CuckooFilterSetting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CuckooFilterSetting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CuckooFilterSetting) ProtoMessage() {}

func (x *CuckooFilterSetting) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CuckooFilterSetting.ProtoReflect.Descriptor instead.
func (*CuckooFilterSetting) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{4}
}

func (x *CuckooFilterSetting) GetCapacity() int64 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

func (x *CuckooFilterSetting) GetExpansion() int64 {
	if x != nil {
		return x.Expansion
	}
	return 0
}

func (x *CuckooFilterSetting) GetNonScaling() bool {
	if x != nil {
		return x.NonScaling
	}
	return false
}

func (x *CuckooFilterSetting) GetBucketSize() int64 {
	if x != nil {
		return x.BucketSize
	}
	return 0
}

func (x *CuckooFilterSetting) GetMaxIterations() int64 {
	if x != nil {
		return x.MaxIterations
	}
	return 0
}

func (x *CuckooFilterSetting) GetMaxTTLSeconds() int64 {
	if x != nil {
		return x.MaxTTLSeconds
	}
	return 0
}

type RedisFilterSetting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Setting:
	//	*RedisFilterSetting_SetfilterSetting
	//	*RedisFilterSetting_BloomfilterSetting
	//	*RedisFilterSetting_CuckoofilterSetting
	Setting isRedisFilterSetting_Setting `protobuf_oneof:"setting"`
}

func (x *RedisFilterSetting) Reset() {
	*x = RedisFilterSetting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisFilterSetting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisFilterSetting) ProtoMessage() {}

func (x *RedisFilterSetting) ProtoReflect() protoreflect.Message {
	mi := &file_pbschema_itemfilterRPC_universal_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisFilterSetting.ProtoReflect.Descriptor instead.
func (*RedisFilterSetting) Descriptor() ([]byte, []int) {
	return file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP(), []int{5}
}

func (m *RedisFilterSetting) GetSetting() isRedisFilterSetting_Setting {
	if m != nil {
		return m.Setting
	}
	return nil
}

func (x *RedisFilterSetting) GetSetfilterSetting() *SetFilterSetting {
	if x, ok := x.GetSetting().(*RedisFilterSetting_SetfilterSetting); ok {
		return x.SetfilterSetting
	}
	return nil
}

func (x *RedisFilterSetting) GetBloomfilterSetting() *BloomFilterSetting {
	if x, ok := x.GetSetting().(*RedisFilterSetting_BloomfilterSetting); ok {
		return x.BloomfilterSetting
	}
	return nil
}

func (x *RedisFilterSetting) GetCuckoofilterSetting() *CuckooFilterSetting {
	if x, ok := x.GetSetting().(*RedisFilterSetting_CuckoofilterSetting); ok {
		return x.CuckoofilterSetting
	}
	return nil
}

type isRedisFilterSetting_Setting interface {
	isRedisFilterSetting_Setting()
}

type RedisFilterSetting_SetfilterSetting struct {
	SetfilterSetting *SetFilterSetting `protobuf:"bytes,1,opt,name=setfilter_setting,json=setfilterSetting,proto3,oneof"`
}

type RedisFilterSetting_BloomfilterSetting struct {
	BloomfilterSetting *BloomFilterSetting `protobuf:"bytes,2,opt,name=BloomfilterSetting,proto3,oneof"`
}

type RedisFilterSetting_CuckoofilterSetting struct {
	CuckoofilterSetting *CuckooFilterSetting `protobuf:"bytes,3,opt,name=CuckoofilterSetting,proto3,oneof"`
}

func (*RedisFilterSetting_SetfilterSetting) isRedisFilterSetting_Setting() {}

func (*RedisFilterSetting_BloomfilterSetting) isRedisFilterSetting_Setting() {}

func (*RedisFilterSetting_CuckoofilterSetting) isRedisFilterSetting_Setting() {}

var File_pbschema_itemfilterRPC_universal_proto protoreflect.FileDescriptor

var file_pbschema_itemfilterRPC_universal_proto_rawDesc = []byte{
	0x0a, 0x26, 0x70, 0x62, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x69, 0x74, 0x65, 0x6d, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x50, 0x43, 0x2f, 0x75, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x69, 0x74, 0x65, 0x6d, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x52, 0x50, 0x43, 0x1a, 0x37, 0x70, 0x62, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65,
	0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xf3, 0x02, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x69, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x4d, 0x65, 0x74, 0x61, 0x12, 0x3c, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x2c, 0x92, 0x41, 0x29, 0x32, 0x27, 0xe5, 0xaf, 0xb9, 0xe8, 0xb1, 0xa1, 0x49, 0x44, 0x2c,
	0xe8, 0xaf, 0xb7, 0xe6, 0xb1, 0x82, 0xe6, 0x97, 0xb6, 0xe4, 0xb8, 0x8d, 0xe5, 0xa1, 0xab, 0xe5,
	0x88, 0x99, 0xe8, 0x87, 0xaa, 0xe5, 0x8a, 0xa8, 0xe7, 0x94, 0x9f, 0xe6, 0x88, 0x90, 0x52, 0x02,
	0x49, 0x44, 0x12, 0x4c, 0x0a, 0x10, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x20, 0x92, 0x41,
	0x1d, 0x32, 0x1b, 0xe5, 0xaf, 0xb9, 0xe8, 0xb1, 0xa1, 0xe9, 0x92, 0x88, 0xe5, 0xaf, 0xb9, 0xe7,
	0x9a, 0x84, 0xe5, 0xae, 0x9e, 0xe4, 0xbd, 0x93, 0xe7, 0xb1, 0xbb, 0xe5, 0x9e, 0x8b, 0x52, 0x10,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x47, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x33,
	0x92, 0x41, 0x30, 0x32, 0x2e, 0xe5, 0xaf, 0xb9, 0xe8, 0xb1, 0xa1, 0xe5, 0x90, 0x8d, 0x2c, 0xe5,
	0x88, 0x9b, 0xe5, 0xbb, 0xba, 0xe4, 0xb8, 0x8a, 0xe4, 0xb8, 0x8b, 0xe6, 0x96, 0x87, 0xe6, 0x88,
	0x96, 0x70, 0x69, 0x63, 0x6b, 0x65, 0x72, 0xe6, 0x97, 0xb6, 0xe4, 0xb8, 0x8d, 0xe7, 0x94, 0xa8,
	0xe5, 0xa1, 0xab, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x4a, 0x0a, 0x04, 0x44, 0x65, 0x73,
	0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x36, 0x92, 0x41, 0x33, 0x32, 0x31, 0xe5, 0xaf,
	0xb9, 0xe8, 0xb1, 0xa1, 0xe6, 0x8f, 0x8f, 0xe8, 0xbf, 0xb0, 0x2c, 0xe5, 0x88, 0x9b, 0xe5, 0xbb,
	0xba, 0xe4, 0xb8, 0x8a, 0xe4, 0xb8, 0x8b, 0xe6, 0x96, 0x87, 0xe6, 0x88, 0x96, 0x70, 0x69, 0x63,
	0x6b, 0x65, 0x72, 0xe6, 0x97, 0xb6, 0xe4, 0xb8, 0x8d, 0xe7, 0x94, 0xa8, 0xe5, 0xa1, 0xab, 0x52,
	0x04, 0x44, 0x65, 0x73, 0x63, 0x12, 0x3f, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x2d, 0x92, 0x41, 0x2a, 0x32, 0x28, 0xe5, 0xaf, 0xb9, 0xe8, 0xb1, 0xa1, 0xe4,
	0xbf, 0x9d, 0xe5, 0xad, 0x98, 0xe7, 0x9a, 0x84, 0x6b, 0x65, 0x79, 0x2c, 0xe8, 0xaf, 0xb7, 0xe6,
	0xb1, 0x82, 0xe5, 0x88, 0x9b, 0xe5, 0xbb, 0xba, 0xe6, 0x97, 0xb6, 0xe4, 0xb8, 0x8d, 0xe5, 0xa1,
	0xab, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0xf6, 0x01, 0x0a, 0x11, 0x52, 0x65, 0x64, 0x69, 0x73,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x06,
	0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x42, 0x17, 0x92, 0x41,
	0x14, 0x32, 0x12, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0xe6, 0x98, 0xaf, 0xe5, 0x90, 0xa6, 0xe5,
	0xad, 0x98, 0xe5, 0x9c, 0xa8, 0x52, 0x06, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x25, 0x0a,
	0x04, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x92, 0x41, 0x0e,
	0x32, 0x0c, 0xe5, 0xbd, 0x93, 0xe5, 0x89, 0x8d, 0xe5, 0xae, 0xb9, 0xe9, 0x87, 0x8f, 0x52, 0x04,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x54, 0x54, 0x4c, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x1f, 0x92, 0x41, 0x1c, 0x32, 0x1a, 0xe5,
	0xbd, 0x93, 0xe5, 0x89, 0x8d, 0xe5, 0x89, 0xa9, 0xe4, 0xbd, 0x99, 0xe6, 0x97, 0xb6, 0xe9, 0x97,
	0xb4, 0x2c, 0xe5, 0x8d, 0x95, 0xe4, 0xbd, 0x8d, 0x73, 0x52, 0x0a, 0x54, 0x54, 0x4c, 0x53, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x73, 0x3a, 0x48, 0x92, 0x41, 0x45, 0x0a, 0x43, 0x2a, 0x0f, 0x52, 0x65,
	0x64, 0x69, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x32, 0x15, 0xe8,
	0xbf, 0x87, 0xe6, 0xbb, 0xa4, 0xe5, 0x99, 0xa8, 0xe9, 0x85, 0x8d, 0xe7, 0xbd, 0xae, 0xe4, 0xbf,
	0xa1, 0xe6, 0x81, 0xaf, 0xd2, 0x01, 0x06, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0xd2, 0x01, 0x08,
	0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0xd2, 0x01, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x22,
	0x5d, 0x0a, 0x10, 0x53, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x49, 0x0a, 0x0d, 0x4d, 0x61, 0x78, 0x54, 0x54, 0x4c, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x23, 0x92, 0x41, 0x20, 0x32,
	0x1e, 0xe4, 0xbf, 0x9d, 0xe5, 0xad, 0x98, 0x6b, 0x65, 0x79, 0xe7, 0x9a, 0x84, 0xe6, 0x9c, 0x80,
	0xe5, 0xa4, 0xa7, 0xe8, 0xbf, 0x87, 0xe6, 0x9c, 0x9f, 0xe7, 0xa7, 0x92, 0xe6, 0x95, 0xb0, 0x52,
	0x0d, 0x4d, 0x61, 0x78, 0x54, 0x54, 0x4c, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x91,
	0x03, 0x0a, 0x12, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x2d, 0x0a, 0x08, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x92, 0x41, 0x0e, 0x32, 0x0c, 0xe6, 0x9c,
	0x80, 0xe5, 0xa4, 0xa7, 0xe5, 0xae, 0xb9, 0xe9, 0x87, 0x8f, 0x52, 0x08, 0x43, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x12, 0x8d, 0x01, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x42, 0x6f, 0x92, 0x41, 0x6c, 0x32, 0x6a, 0xe7,
	0xa2, 0xb0, 0xe6, 0x92, 0x9e, 0xe7, 0x8e, 0x87, 0x2c, 0xe7, 0xa2, 0xb0, 0xe6, 0x92, 0x9e, 0xe7,
	0x8e, 0x87, 0xe8, 0xae, 0xbe, 0xe7, 0xbd, 0xae, 0xe7, 0x9a, 0x84, 0xe8, 0xb6, 0x8a, 0xe4, 0xbd,
	0x8e, 0xe4, 0xbd, 0xbf, 0xe7, 0x94, 0xa8, 0xe7, 0x9a, 0x84, 0x68, 0x61, 0x73, 0x68, 0xe5, 0x87,
	0xbd, 0xe6, 0x95, 0xb0, 0xe8, 0xb6, 0x8a, 0xe5, 0xa4, 0x9a, 0x2c, 0xe4, 0xbd, 0xbf, 0xe7, 0x94,
	0xa8, 0xe7, 0x9a, 0x84, 0xe7, 0xa9, 0xba, 0xe9, 0x97, 0xb4, 0xe4, 0xb9, 0x9f, 0xe8, 0xb6, 0x8a,
	0xe5, 0xa4, 0xa7, 0x2c, 0xe6, 0xa3, 0x80, 0xe7, 0xb4, 0xa2, 0xe6, 0x95, 0x88, 0xe7, 0x8e, 0x87,
	0xe4, 0xb9, 0x9f, 0xe8, 0xb6, 0x8a, 0xe4, 0xbd, 0x8e, 0x52, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x2f, 0x0a, 0x09, 0x45, 0x78, 0x70, 0x61, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x92, 0x41, 0x0e, 0x32, 0x0c, 0xe6, 0x89,
	0xa9, 0xe5, 0xae, 0xb9, 0xe5, 0x80, 0x8d, 0xe6, 0x95, 0xb0, 0x52, 0x09, 0x45, 0x78, 0x70, 0x61,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x0a, 0x4e, 0x6f, 0x6e, 0x53, 0x63, 0x61, 0x6c,
	0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x42, 0x20, 0x92, 0x41, 0x1d, 0x32, 0x1b,
	0xe6, 0x98, 0xaf, 0xe5, 0x90, 0xa6, 0xe5, 0x8f, 0xaf, 0xe4, 0xbb, 0xa5, 0xe6, 0x89, 0xa9, 0xe5,
	0xae, 0xb9, 0xe8, 0xbf, 0x87, 0xe6, 0xbb, 0xa4, 0xe5, 0x99, 0xa8, 0x52, 0x0a, 0x4e, 0x6f, 0x6e,
	0x53, 0x63, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x12, 0x49, 0x0a, 0x0d, 0x4d, 0x61, 0x78, 0x54, 0x54,
	0x4c, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x23,
	0x92, 0x41, 0x20, 0x32, 0x1e, 0xe4, 0xbf, 0x9d, 0xe5, 0xad, 0x98, 0x6b, 0x65, 0x79, 0xe7, 0x9a,
	0x84, 0xe6, 0x9c, 0x80, 0xe5, 0xa4, 0xa7, 0xe8, 0xbf, 0x87, 0xe6, 0x9c, 0x9f, 0xe7, 0xa7, 0x92,
	0xe6, 0x95, 0xb0, 0x52, 0x0d, 0x4d, 0x61, 0x78, 0x54, 0x54, 0x4c, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x22, 0xaf, 0x05, 0x0a, 0x13, 0x43, 0x75, 0x63, 0x6b, 0x6f, 0x6f, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x2d, 0x0a, 0x08, 0x43, 0x61,
	0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x92, 0x41,
	0x0e, 0x32, 0x0c, 0xe6, 0x9c, 0x80, 0xe5, 0xa4, 0xa7, 0xe5, 0xae, 0xb9, 0xe9, 0x87, 0x8f, 0x52,
	0x08, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x2f, 0x0a, 0x09, 0x45, 0x78, 0x70,
	0x61, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x42, 0x11, 0x92, 0x41,
	0x0e, 0x32, 0x0c, 0xe6, 0x89, 0xa9, 0xe5, 0xae, 0xb9, 0xe5, 0x80, 0x8d, 0xe6, 0x95, 0xb0, 0x52,
	0x09, 0x45, 0x78, 0x70, 0x61, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x0a, 0x4e, 0x6f,
	0x6e, 0x53, 0x63, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x42, 0x20,
	0x92, 0x41, 0x1d, 0x32, 0x1b, 0xe6, 0x98, 0xaf, 0xe5, 0x90, 0xa6, 0xe5, 0x8f, 0xaf, 0xe4, 0xbb,
	0xa5, 0xe6, 0x89, 0xa9, 0xe5, 0xae, 0xb9, 0xe8, 0xbf, 0x87, 0xe6, 0xbb, 0xa4, 0xe5, 0x99, 0xa8,
	0x52, 0x0a, 0x4e, 0x6f, 0x6e, 0x53, 0x63, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x12, 0xbf, 0x01, 0x0a,
	0x0a, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x9e, 0x01, 0x92, 0x41, 0x9a, 0x01, 0x32, 0x97, 0x01, 0xe5, 0x8f, 0xaa, 0xe5, 0xaf,
	0xb9, 0x43, 0x55, 0x43, 0x4b, 0x4f, 0x4f, 0xe6, 0x9c, 0x89, 0xe6, 0x95, 0x88, 0x2c, 0xe6, 0xa1,
	0xb6, 0xe5, 0xa4, 0xa7, 0xe5, 0xb0, 0x8f, 0x2c, 0xe6, 0xaf, 0x8f, 0xe4, 0xb8, 0xaa, 0xe5, 0xad,
	0x98, 0xe5, 0x82, 0xa8, 0xe6, 0xa1, 0xb6, 0xe4, 0xb8, 0xad, 0xe7, 0x9a, 0x84, 0xe9, 0xa1, 0xb9,
	0xe7, 0x9b, 0xae, 0xe6, 0x95, 0xb0, 0x2e, 0xe8, 0xbe, 0x83, 0xe9, 0xab, 0x98, 0xe7, 0x9a, 0x84,
	0xe6, 0xa1, 0xb6, 0xe5, 0xa4, 0xa7, 0xe5, 0xb0, 0x8f, 0xe5, 0x80, 0xbc, 0xe4, 0xbc, 0x9a, 0xe6,
	0x8f, 0x90, 0xe9, 0xab, 0x98, 0xe5, 0xa1, 0xab, 0xe5, 0x85, 0x85, 0xe7, 0x8e, 0x87, 0x2c, 0xe4,
	0xbd, 0x86, 0xe4, 0xb9, 0x9f, 0xe4, 0xbc, 0x9a, 0xe5, 0xaf, 0xbc, 0xe8, 0x87, 0xb4, 0xe6, 0x9b,
	0xb4, 0xe9, 0xab, 0x98, 0xe7, 0x9a, 0x84, 0xe9, 0x94, 0x99, 0xe8, 0xaf, 0xaf, 0xe7, 0x8e, 0x87,
	0xe5, 0x92, 0x8c, 0xe7, 0xa8, 0x8d, 0xe6, 0x85, 0xa2, 0xe7, 0x9a, 0x84, 0xe6, 0x80, 0xa7, 0xe8,
	0x83, 0xbd, 0x52, 0x0a, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0xe8,
	0x01, 0x0a, 0x0d, 0x4d, 0x61, 0x78, 0x49, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0xc1, 0x01, 0x92, 0x41, 0xbd, 0x01, 0x32, 0xba, 0x01,
	0xe5, 0x8f, 0xaa, 0xe5, 0xaf, 0xb9, 0x43, 0x55, 0x43, 0x4b, 0x4f, 0x4f, 0xe6, 0x9c, 0x89, 0xe6,
	0x95, 0x88, 0x2c, 0xe5, 0x9c, 0xa8, 0xe5, 0xa3, 0xb0, 0xe6, 0x98, 0x8e, 0xe8, 0xbf, 0x87, 0xe6,
	0xbb, 0xa4, 0xe5, 0x99, 0xa8, 0xe5, 0xb7, 0xb2, 0xe6, 0xbb, 0xa1, 0xe5, 0xb9, 0xb6, 0xe5, 0x88,
	0x9b, 0xe5, 0xbb, 0xba, 0xe9, 0x99, 0x84, 0xe5, 0x8a, 0xa0, 0xe8, 0xbf, 0x87, 0xe6, 0xbb, 0xa4,
	0xe5, 0x99, 0xa8, 0xe4, 0xb9, 0x8b, 0xe5, 0x89, 0x8d, 0xe5, 0xb0, 0x9d, 0xe8, 0xaf, 0x95, 0xe5,
	0x9c, 0xa8, 0xe5, 0xad, 0x98, 0xe5, 0x82, 0xa8, 0xe6, 0xa1, 0xb6, 0xe4, 0xb9, 0x8b, 0xe9, 0x97,
	0xb4, 0xe4, 0xba, 0xa4, 0xe6, 0x8d, 0xa2, 0xe9, 0xa1, 0xb9, 0xe7, 0x9b, 0xae, 0xe7, 0x9a, 0x84,
	0xe6, 0xac, 0xa1, 0xe6, 0x95, 0xb0, 0x2e, 0xe8, 0xbe, 0x83, 0xe4, 0xbd, 0x8e, 0xe7, 0x9a, 0x84,
	0xe5, 0x80, 0xbc, 0xe5, 0xaf, 0xb9, 0xe6, 0x80, 0xa7, 0xe8, 0x83, 0xbd, 0xe6, 0x9b, 0xb4, 0xe5,
	0xa5, 0xbd, 0x2c, 0xe8, 0xbe, 0x83, 0xe9, 0xab, 0x98, 0xe7, 0x9a, 0x84, 0xe5, 0x80, 0xbc, 0xe5,
	0xaf, 0xb9, 0xe8, 0xbf, 0x87, 0xe6, 0xbb, 0xa4, 0xe5, 0x99, 0xa8, 0xe5, 0xa1, 0xab, 0xe5, 0x85,
	0x85, 0xe7, 0x8e, 0x87, 0xe6, 0x9b, 0xb4, 0xe5, 0xa5, 0xbd, 0x52, 0x0d, 0x4d, 0x61, 0x78, 0x49,
	0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x49, 0x0a, 0x0d, 0x4d, 0x61, 0x78,
	0x54, 0x54, 0x4c, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x23, 0x92, 0x41, 0x20, 0x32, 0x1e, 0xe4, 0xbf, 0x9d, 0xe5, 0xad, 0x98, 0x6b, 0x65, 0x79,
	0xe7, 0x9a, 0x84, 0xe6, 0x9c, 0x80, 0xe5, 0xa4, 0xa7, 0xe8, 0xbf, 0x87, 0xe6, 0x9c, 0x9f, 0xe7,
	0xa7, 0x92, 0xe6, 0x95, 0xb0, 0x52, 0x0d, 0x4d, 0x61, 0x78, 0x54, 0x54, 0x4c, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x22, 0xf5, 0x02, 0x0a, 0x12, 0x52, 0x65, 0x64, 0x69, 0x73, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x6a, 0x0a, 0x11, 0x73,
	0x65, 0x74, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x74, 0x65, 0x6d, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x52, 0x50, 0x43, 0x2e, 0x53, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x1a, 0x92, 0x41, 0x17, 0x32, 0x15, 0x73, 0x65,
	0x74, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0xe7, 0x9a, 0x84, 0xe8, 0xae, 0xbe, 0xe7, 0xbd, 0xae,
	0xe9, 0xa1, 0xb9, 0x48, 0x00, 0x52, 0x10, 0x73, 0x65, 0x74, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x71, 0x0a, 0x12, 0x42, 0x6c, 0x6f, 0x6f, 0x6d,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x74, 0x65, 0x6d, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x52, 0x50, 0x43, 0x2e, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x1c, 0x92, 0x41, 0x19, 0x32, 0x17, 0x62, 0x6c, 0x6f,
	0x6f, 0x6d, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0xe7, 0x9a, 0x84, 0xe8, 0xae, 0xbe, 0xe7, 0xbd,
	0xae, 0xe9, 0xa1, 0xb9, 0x48, 0x00, 0x52, 0x12, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x75, 0x0a, 0x13, 0x43, 0x75,
	0x63, 0x6b, 0x6f, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x74, 0x65, 0x6d, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x52, 0x50, 0x43, 0x2e, 0x43, 0x75, 0x63, 0x6b, 0x6f, 0x6f, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x1d, 0x92, 0x41, 0x1a,
	0x32, 0x18, 0x63, 0x75, 0x63, 0x6b, 0x6f, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0xe7, 0x9a,
	0x84, 0xe8, 0xae, 0xbe, 0xe7, 0xbd, 0xae, 0xe9, 0xa1, 0xb9, 0x48, 0x00, 0x52, 0x13, 0x43, 0x75,
	0x63, 0x6b, 0x6f, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x42, 0x09, 0x0a, 0x07, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x2a, 0x31, 0x0a, 0x0f,
	0x52, 0x65, 0x64, 0x69, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x07, 0x0a, 0x03, 0x53, 0x45, 0x54, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x4c, 0x4f, 0x4f,
	0x4d, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x55, 0x43, 0x4b, 0x4f, 0x4f, 0x10, 0x02, 0x42,
	0x14, 0x5a, 0x12, 0x2e, 0x2f, 0x69, 0x74, 0x65, 0x6d, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52,
	0x50, 0x43, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pbschema_itemfilterRPC_universal_proto_rawDescOnce sync.Once
	file_pbschema_itemfilterRPC_universal_proto_rawDescData = file_pbschema_itemfilterRPC_universal_proto_rawDesc
)

func file_pbschema_itemfilterRPC_universal_proto_rawDescGZIP() []byte {
	file_pbschema_itemfilterRPC_universal_proto_rawDescOnce.Do(func() {
		file_pbschema_itemfilterRPC_universal_proto_rawDescData = protoimpl.X.CompressGZIP(file_pbschema_itemfilterRPC_universal_proto_rawDescData)
	})
	return file_pbschema_itemfilterRPC_universal_proto_rawDescData
}

var file_pbschema_itemfilterRPC_universal_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pbschema_itemfilterRPC_universal_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pbschema_itemfilterRPC_universal_proto_goTypes = []interface{}{
	(RedisFilterType)(0),        // 0: itemfilterRPC.RedisFilterType
	(*RedisFilterMeta)(nil),     // 1: itemfilterRPC.RedisFilterMeta
	(*RedisFilterStatus)(nil),   // 2: itemfilterRPC.RedisFilterStatus
	(*SetFilterSetting)(nil),    // 3: itemfilterRPC.SetFilterSetting
	(*BloomFilterSetting)(nil),  // 4: itemfilterRPC.BloomFilterSetting
	(*CuckooFilterSetting)(nil), // 5: itemfilterRPC.CuckooFilterSetting
	(*RedisFilterSetting)(nil),  // 6: itemfilterRPC.RedisFilterSetting
}
var file_pbschema_itemfilterRPC_universal_proto_depIdxs = []int32{
	3, // 0: itemfilterRPC.RedisFilterSetting.setfilter_setting:type_name -> itemfilterRPC.SetFilterSetting
	4, // 1: itemfilterRPC.RedisFilterSetting.BloomfilterSetting:type_name -> itemfilterRPC.BloomFilterSetting
	5, // 2: itemfilterRPC.RedisFilterSetting.CuckoofilterSetting:type_name -> itemfilterRPC.CuckooFilterSetting
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_pbschema_itemfilterRPC_universal_proto_init() }
func file_pbschema_itemfilterRPC_universal_proto_init() {
	if File_pbschema_itemfilterRPC_universal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisFilterMeta); i {
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
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisFilterStatus); i {
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
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetFilterSetting); i {
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
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BloomFilterSetting); i {
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
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CuckooFilterSetting); i {
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
		file_pbschema_itemfilterRPC_universal_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisFilterSetting); i {
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
	file_pbschema_itemfilterRPC_universal_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*RedisFilterSetting_SetfilterSetting)(nil),
		(*RedisFilterSetting_BloomfilterSetting)(nil),
		(*RedisFilterSetting_CuckoofilterSetting)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pbschema_itemfilterRPC_universal_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pbschema_itemfilterRPC_universal_proto_goTypes,
		DependencyIndexes: file_pbschema_itemfilterRPC_universal_proto_depIdxs,
		EnumInfos:         file_pbschema_itemfilterRPC_universal_proto_enumTypes,
		MessageInfos:      file_pbschema_itemfilterRPC_universal_proto_msgTypes,
	}.Build()
	File_pbschema_itemfilterRPC_universal_proto = out.File
	file_pbschema_itemfilterRPC_universal_proto_rawDesc = nil
	file_pbschema_itemfilterRPC_universal_proto_goTypes = nil
	file_pbschema_itemfilterRPC_universal_proto_depIdxs = nil
}
