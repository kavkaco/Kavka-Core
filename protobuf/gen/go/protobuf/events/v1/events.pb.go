// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: protobuf/events/v1/events.proto

package eventsv1

import (
	v1 "github.com/kavkaco/Kavka-Core/protobuf/gen/go/protobuf/model/chat/v1"
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

type SubscribeEventsStreamResponse_Type int32

const (
	SubscribeEventsStreamResponse_TYPE_UNSPECIFIED    SubscribeEventsStreamResponse_Type = 0
	SubscribeEventsStreamResponse_TYPE_ADD_CHAT       SubscribeEventsStreamResponse_Type = 1
	SubscribeEventsStreamResponse_TYPE_REMOVE_CHAT    SubscribeEventsStreamResponse_Type = 2
	SubscribeEventsStreamResponse_TYPE_UPDATE_CHAT    SubscribeEventsStreamResponse_Type = 3
	SubscribeEventsStreamResponse_TYPE_ADD_MESSAGE    SubscribeEventsStreamResponse_Type = 4
	SubscribeEventsStreamResponse_TYPE_REMOVE_MESSAGE SubscribeEventsStreamResponse_Type = 5
	SubscribeEventsStreamResponse_TYPE_UPDATE_MESSAGE SubscribeEventsStreamResponse_Type = 6
	SubscribeEventsStreamResponse_TYPE_CLEAR_CHAT     SubscribeEventsStreamResponse_Type = 7
	SubscribeEventsStreamResponse_TYPE_MESSAGE_SEEN   SubscribeEventsStreamResponse_Type = 8
)

// Enum value maps for SubscribeEventsStreamResponse_Type.
var (
	SubscribeEventsStreamResponse_Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "TYPE_ADD_CHAT",
		2: "TYPE_REMOVE_CHAT",
		3: "TYPE_UPDATE_CHAT",
		4: "TYPE_ADD_MESSAGE",
		5: "TYPE_REMOVE_MESSAGE",
		6: "TYPE_UPDATE_MESSAGE",
		7: "TYPE_CLEAR_CHAT",
		8: "TYPE_MESSAGE_SEEN",
	}
	SubscribeEventsStreamResponse_Type_value = map[string]int32{
		"TYPE_UNSPECIFIED":    0,
		"TYPE_ADD_CHAT":       1,
		"TYPE_REMOVE_CHAT":    2,
		"TYPE_UPDATE_CHAT":    3,
		"TYPE_ADD_MESSAGE":    4,
		"TYPE_REMOVE_MESSAGE": 5,
		"TYPE_UPDATE_MESSAGE": 6,
		"TYPE_CLEAR_CHAT":     7,
		"TYPE_MESSAGE_SEEN":   8,
	}
)

func (x SubscribeEventsStreamResponse_Type) Enum() *SubscribeEventsStreamResponse_Type {
	p := new(SubscribeEventsStreamResponse_Type)
	*p = x
	return p
}

func (x SubscribeEventsStreamResponse_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SubscribeEventsStreamResponse_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_events_v1_events_proto_enumTypes[0].Descriptor()
}

func (SubscribeEventsStreamResponse_Type) Type() protoreflect.EnumType {
	return &file_protobuf_events_v1_events_proto_enumTypes[0]
}

func (x SubscribeEventsStreamResponse_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SubscribeEventsStreamResponse_Type.Descriptor instead.
func (SubscribeEventsStreamResponse_Type) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_events_v1_events_proto_rawDescGZIP(), []int{2, 0}
}

type StreamEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderUserId    string   `protobuf:"bytes,1,opt,name=sender_user_id,json=senderUserId,proto3" json:"sender_user_id,omitempty"`
	ReceiversUserId []string `protobuf:"bytes,2,rep,name=receivers_user_id,json=receiversUserId,proto3" json:"receivers_user_id,omitempty"`
	Payload         []byte   `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *StreamEvent) Reset() {
	*x = StreamEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_events_v1_events_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamEvent) ProtoMessage() {}

func (x *StreamEvent) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_events_v1_events_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamEvent.ProtoReflect.Descriptor instead.
func (*StreamEvent) Descriptor() ([]byte, []int) {
	return file_protobuf_events_v1_events_proto_rawDescGZIP(), []int{0}
}

func (x *StreamEvent) GetSenderUserId() string {
	if x != nil {
		return x.SenderUserId
	}
	return ""
}

func (x *StreamEvent) GetReceiversUserId() []string {
	if x != nil {
		return x.ReceiversUserId
	}
	return nil
}

func (x *StreamEvent) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type SubscribeEventsStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubscribeEventsStreamRequest) Reset() {
	*x = SubscribeEventsStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_events_v1_events_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeEventsStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeEventsStreamRequest) ProtoMessage() {}

func (x *SubscribeEventsStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_events_v1_events_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeEventsStreamRequest.ProtoReflect.Descriptor instead.
func (*SubscribeEventsStreamRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_events_v1_events_proto_rawDescGZIP(), []int{1}
}

type SubscribeEventsStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string                             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type SubscribeEventsStreamResponse_Type `protobuf:"varint,2,opt,name=type,proto3,enum=protobuf.events.v1.SubscribeEventsStreamResponse_Type" json:"type,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*SubscribeEventsStreamResponse_AddChat
	Payload isSubscribeEventsStreamResponse_Payload `protobuf_oneof:"payload"`
}

func (x *SubscribeEventsStreamResponse) Reset() {
	*x = SubscribeEventsStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_events_v1_events_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeEventsStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeEventsStreamResponse) ProtoMessage() {}

func (x *SubscribeEventsStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_events_v1_events_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeEventsStreamResponse.ProtoReflect.Descriptor instead.
func (*SubscribeEventsStreamResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_events_v1_events_proto_rawDescGZIP(), []int{2}
}

func (x *SubscribeEventsStreamResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SubscribeEventsStreamResponse) GetType() SubscribeEventsStreamResponse_Type {
	if x != nil {
		return x.Type
	}
	return SubscribeEventsStreamResponse_TYPE_UNSPECIFIED
}

func (m *SubscribeEventsStreamResponse) GetPayload() isSubscribeEventsStreamResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *SubscribeEventsStreamResponse) GetAddChat() *AddChat {
	if x, ok := x.GetPayload().(*SubscribeEventsStreamResponse_AddChat); ok {
		return x.AddChat
	}
	return nil
}

type isSubscribeEventsStreamResponse_Payload interface {
	isSubscribeEventsStreamResponse_Payload()
}

type SubscribeEventsStreamResponse_AddChat struct {
	AddChat *AddChat `protobuf:"bytes,3,opt,name=add_chat,json=addChat,proto3,oneof"`
}

func (*SubscribeEventsStreamResponse_AddChat) isSubscribeEventsStreamResponse_Payload() {}

type AddChat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chat *v1.Chat `protobuf:"bytes,1,opt,name=chat,proto3" json:"chat,omitempty"`
}

func (x *AddChat) Reset() {
	*x = AddChat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_events_v1_events_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddChat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddChat) ProtoMessage() {}

func (x *AddChat) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_events_v1_events_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddChat.ProtoReflect.Descriptor instead.
func (*AddChat) Descriptor() ([]byte, []int) {
	return file_protobuf_events_v1_events_proto_rawDescGZIP(), []int{3}
}

func (x *AddChat) GetChat() *v1.Chat {
	if x != nil {
		return x.Chat
	}
	return nil
}

var File_protobuf_events_v1_events_proto protoreflect.FileDescriptor

var file_protobuf_events_v1_events_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68,
	0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a,
	0x11, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x73, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0x1e, 0x0a, 0x1c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x96, 0x03, 0x0a, 0x1d, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4a, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x36, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x61, 0x64, 0x64, 0x5f, 0x63, 0x68, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64,
	0x43, 0x68, 0x61, 0x74, 0x48, 0x00, 0x52, 0x07, 0x61, 0x64, 0x64, 0x43, 0x68, 0x61, 0x74, 0x22,
	0xcf, 0x01, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11,
	0x0a, 0x0d, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x44, 0x44, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x10,
	0x01, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45,
	0x5f, 0x43, 0x48, 0x41, 0x54, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x10, 0x03, 0x12, 0x14, 0x0a,
	0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x44, 0x44, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x10, 0x04, 0x12, 0x17, 0x0a, 0x13, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x4d, 0x4f,
	0x56, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x05, 0x12, 0x17, 0x0a, 0x13,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53,
	0x41, 0x47, 0x45, 0x10, 0x06, 0x12, 0x13, 0x0a, 0x0f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4c,
	0x45, 0x41, 0x52, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x45, 0x4e, 0x10,
	0x08, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x3b, 0x0a, 0x07,
	0x41, 0x64, 0x64, 0x43, 0x68, 0x61, 0x74, 0x12, 0x30, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x68, 0x61, 0x74, 0x52, 0x04, 0x63, 0x68, 0x61, 0x74, 0x32, 0x92, 0x01, 0x0a, 0x0d, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x15,
	0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0xda,
	0x01, 0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x76, 0x6b, 0x61, 0x63, 0x6f, 0x2f, 0x4b, 0x61, 0x76,
	0x6b, 0x61, 0x2d, 0x43, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x45, 0x58, 0xaa, 0x02, 0x12, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x12, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x5c, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x14, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x3a,
	0x3a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_protobuf_events_v1_events_proto_rawDescOnce sync.Once
	file_protobuf_events_v1_events_proto_rawDescData = file_protobuf_events_v1_events_proto_rawDesc
)

func file_protobuf_events_v1_events_proto_rawDescGZIP() []byte {
	file_protobuf_events_v1_events_proto_rawDescOnce.Do(func() {
		file_protobuf_events_v1_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_events_v1_events_proto_rawDescData)
	})
	return file_protobuf_events_v1_events_proto_rawDescData
}

var file_protobuf_events_v1_events_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protobuf_events_v1_events_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protobuf_events_v1_events_proto_goTypes = []any{
	(SubscribeEventsStreamResponse_Type)(0), // 0: protobuf.events.v1.SubscribeEventsStreamResponse.Type
	(*StreamEvent)(nil),                     // 1: protobuf.events.v1.StreamEvent
	(*SubscribeEventsStreamRequest)(nil),    // 2: protobuf.events.v1.SubscribeEventsStreamRequest
	(*SubscribeEventsStreamResponse)(nil),   // 3: protobuf.events.v1.SubscribeEventsStreamResponse
	(*AddChat)(nil),                         // 4: protobuf.events.v1.AddChat
	(*v1.Chat)(nil),                         // 5: protobuf.model.chat.v1.Chat
}
var file_protobuf_events_v1_events_proto_depIdxs = []int32{
	0, // 0: protobuf.events.v1.SubscribeEventsStreamResponse.type:type_name -> protobuf.events.v1.SubscribeEventsStreamResponse.Type
	4, // 1: protobuf.events.v1.SubscribeEventsStreamResponse.add_chat:type_name -> protobuf.events.v1.AddChat
	5, // 2: protobuf.events.v1.AddChat.chat:type_name -> protobuf.model.chat.v1.Chat
	2, // 3: protobuf.events.v1.EventsService.SubscribeEventsStream:input_type -> protobuf.events.v1.SubscribeEventsStreamRequest
	3, // 4: protobuf.events.v1.EventsService.SubscribeEventsStream:output_type -> protobuf.events.v1.SubscribeEventsStreamResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protobuf_events_v1_events_proto_init() }
func file_protobuf_events_v1_events_proto_init() {
	if File_protobuf_events_v1_events_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_events_v1_events_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*StreamEvent); i {
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
		file_protobuf_events_v1_events_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeEventsStreamRequest); i {
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
		file_protobuf_events_v1_events_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeEventsStreamResponse); i {
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
		file_protobuf_events_v1_events_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AddChat); i {
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
	file_protobuf_events_v1_events_proto_msgTypes[2].OneofWrappers = []any{
		(*SubscribeEventsStreamResponse_AddChat)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_events_v1_events_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_events_v1_events_proto_goTypes,
		DependencyIndexes: file_protobuf_events_v1_events_proto_depIdxs,
		EnumInfos:         file_protobuf_events_v1_events_proto_enumTypes,
		MessageInfos:      file_protobuf_events_v1_events_proto_msgTypes,
	}.Build()
	File_protobuf_events_v1_events_proto = out.File
	file_protobuf_events_v1_events_proto_rawDesc = nil
	file_protobuf_events_v1_events_proto_goTypes = nil
	file_protobuf_events_v1_events_proto_depIdxs = nil
}
