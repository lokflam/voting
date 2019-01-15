// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/events_pb2/events.proto

package events_pb2

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EventFilter_FilterType int32

const (
	EventFilter_FILTER_TYPE_UNSET EventFilter_FilterType = 0
	EventFilter_SIMPLE_ANY        EventFilter_FilterType = 1
	EventFilter_SIMPLE_ALL        EventFilter_FilterType = 2
	EventFilter_REGEX_ANY         EventFilter_FilterType = 3
	EventFilter_REGEX_ALL         EventFilter_FilterType = 4
)

var EventFilter_FilterType_name = map[int32]string{
	0: "FILTER_TYPE_UNSET",
	1: "SIMPLE_ANY",
	2: "SIMPLE_ALL",
	3: "REGEX_ANY",
	4: "REGEX_ALL",
}

var EventFilter_FilterType_value = map[string]int32{
	"FILTER_TYPE_UNSET": 0,
	"SIMPLE_ANY":        1,
	"SIMPLE_ALL":        2,
	"REGEX_ANY":         3,
	"REGEX_ALL":         4,
}

func (x EventFilter_FilterType) String() string {
	return proto.EnumName(EventFilter_FilterType_name, int32(x))
}

func (EventFilter_FilterType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{2, 0}
}

type Event struct {
	// Used to subscribe to events and servers as a hint for how to deserialize
	// event_data and what pairs to expect in attributes.
	EventType  string             `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Attributes []*Event_Attribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
	// Opaque data defined by the event_type.
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *Event) GetAttributes() []*Event_Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *Event) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// Transparent data defined by the event_type.
type Event_Attribute struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_Attribute) Reset()         { *m = Event_Attribute{} }
func (m *Event_Attribute) String() string { return proto.CompactTextString(m) }
func (*Event_Attribute) ProtoMessage()    {}
func (*Event_Attribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{0, 0}
}

func (m *Event_Attribute) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event_Attribute.Unmarshal(m, b)
}
func (m *Event_Attribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event_Attribute.Marshal(b, m, deterministic)
}
func (m *Event_Attribute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_Attribute.Merge(m, src)
}
func (m *Event_Attribute) XXX_Size() int {
	return xxx_messageInfo_Event_Attribute.Size(m)
}
func (m *Event_Attribute) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_Attribute.DiscardUnknown(m)
}

var xxx_messageInfo_Event_Attribute proto.InternalMessageInfo

func (m *Event_Attribute) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Event_Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type EventList struct {
	Events               []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventList) Reset()         { *m = EventList{} }
func (m *EventList) String() string { return proto.CompactTextString(m) }
func (*EventList) ProtoMessage()    {}
func (*EventList) Descriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{1}
}

func (m *EventList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventList.Unmarshal(m, b)
}
func (m *EventList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventList.Marshal(b, m, deterministic)
}
func (m *EventList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventList.Merge(m, src)
}
func (m *EventList) XXX_Size() int {
	return xxx_messageInfo_EventList.Size(m)
}
func (m *EventList) XXX_DiscardUnknown() {
	xxx_messageInfo_EventList.DiscardUnknown(m)
}

var xxx_messageInfo_EventList proto.InternalMessageInfo

func (m *EventList) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type EventFilter struct {
	// EventFilter is used when subscribing to events to limit the events
	// received within a given event type. See
	// validator/server/events/subscription.py for further explanation.
	Key                  string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	MatchString          string                 `protobuf:"bytes,2,opt,name=match_string,json=matchString,proto3" json:"match_string,omitempty"`
	FilterType           EventFilter_FilterType `protobuf:"varint,3,opt,name=filter_type,json=filterType,proto3,enum=EventFilter_FilterType" json:"filter_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *EventFilter) Reset()         { *m = EventFilter{} }
func (m *EventFilter) String() string { return proto.CompactTextString(m) }
func (*EventFilter) ProtoMessage()    {}
func (*EventFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{2}
}

func (m *EventFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventFilter.Unmarshal(m, b)
}
func (m *EventFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventFilter.Marshal(b, m, deterministic)
}
func (m *EventFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventFilter.Merge(m, src)
}
func (m *EventFilter) XXX_Size() int {
	return xxx_messageInfo_EventFilter.Size(m)
}
func (m *EventFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_EventFilter.DiscardUnknown(m)
}

var xxx_messageInfo_EventFilter proto.InternalMessageInfo

func (m *EventFilter) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *EventFilter) GetMatchString() string {
	if m != nil {
		return m.MatchString
	}
	return ""
}

func (m *EventFilter) GetFilterType() EventFilter_FilterType {
	if m != nil {
		return m.FilterType
	}
	return EventFilter_FILTER_TYPE_UNSET
}

type EventSubscription struct {
	// EventSubscription is used when subscribing to events to specify the type
	// of events being subscribed to, along with any additional filters. See
	// validator/server/events/subscription.py for further explanation.
	EventType            string         `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Filters              []*EventFilter `protobuf:"bytes,2,rep,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *EventSubscription) Reset()         { *m = EventSubscription{} }
func (m *EventSubscription) String() string { return proto.CompactTextString(m) }
func (*EventSubscription) ProtoMessage()    {}
func (*EventSubscription) Descriptor() ([]byte, []int) {
	return fileDescriptor_38ead4ef37522737, []int{3}
}

func (m *EventSubscription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSubscription.Unmarshal(m, b)
}
func (m *EventSubscription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSubscription.Marshal(b, m, deterministic)
}
func (m *EventSubscription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSubscription.Merge(m, src)
}
func (m *EventSubscription) XXX_Size() int {
	return xxx_messageInfo_EventSubscription.Size(m)
}
func (m *EventSubscription) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSubscription.DiscardUnknown(m)
}

var xxx_messageInfo_EventSubscription proto.InternalMessageInfo

func (m *EventSubscription) GetEventType() string {
	if m != nil {
		return m.EventType
	}
	return ""
}

func (m *EventSubscription) GetFilters() []*EventFilter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func init() {
	proto.RegisterEnum("EventFilter_FilterType", EventFilter_FilterType_name, EventFilter_FilterType_value)
	proto.RegisterType((*Event)(nil), "Event")
	proto.RegisterType((*Event_Attribute)(nil), "Event.Attribute")
	proto.RegisterType((*EventList)(nil), "EventList")
	proto.RegisterType((*EventFilter)(nil), "EventFilter")
	proto.RegisterType((*EventSubscription)(nil), "EventSubscription")
}

func init() { proto.RegisterFile("github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2/events.proto", fileDescriptor_38ead4ef37522737) }

var fileDescriptor_38ead4ef37522737 = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xd1, 0x6a, 0xe2, 0x40,
	0x14, 0x86, 0x77, 0x8c, 0xba, 0xe4, 0xc4, 0x95, 0x38, 0xac, 0x6c, 0x58, 0xd8, 0x25, 0x1b, 0xd8,
	0x45, 0x58, 0x48, 0x8b, 0xde, 0xf4, 0xd6, 0x42, 0x2c, 0x42, 0x2a, 0x32, 0x49, 0xa1, 0x7a, 0x13,
	0x12, 0x1d, 0x6b, 0xd0, 0x9a, 0x90, 0x4c, 0x2c, 0x3e, 0x4e, 0x9f, 0xac, 0xaf, 0x52, 0x72, 0x92,
	0x68, 0x2e, 0x0a, 0xbd, 0xca, 0x39, 0xdf, 0x3f, 0x93, 0xff, 0x3f, 0x33, 0x03, 0x7a, 0x9c, 0x44,
	0x22, 0x0a, 0xb2, 0xcd, 0x15, 0x3f, 0xf2, 0x83, 0x48, 0xbd, 0x38, 0x18, 0x96, 0xa5, 0x89, 0x92,
	0xf1, 0x4a, 0xa0, 0x65, 0xe5, 0x80, 0xfe, 0x02, 0x40, 0xc5, 0x13, 0xa7, 0x98, 0x6b, 0x44, 0x27,
	0x03, 0x99, 0xc9, 0x48, 0xdc, 0x53, 0xcc, 0xe9, 0x35, 0x80, 0x2f, 0x44, 0x12, 0x06, 0x99, 0xe0,
	0xa9, 0xd6, 0xd0, 0xa5, 0x81, 0x32, 0x54, 0x4d, 0xdc, 0x6a, 0x8e, 0x2b, 0x81, 0xd5, 0xd6, 0x50,
	0x0a, 0xcd, 0xb5, 0x2f, 0x7c, 0x4d, 0xd2, 0xc9, 0xa0, 0xc3, 0xb0, 0xfe, 0x39, 0x02, 0xf9, 0xbc,
	0x98, 0xaa, 0x20, 0xed, 0xf8, 0xa9, 0xb4, 0xca, 0x4b, 0xfa, 0x1d, 0x5a, 0x47, 0x7f, 0x9f, 0x71,
	0xad, 0x81, 0xac, 0x68, 0x8c, 0xff, 0x20, 0xa3, 0x8f, 0x1d, 0xa6, 0x82, 0xfe, 0x86, 0x76, 0x31,
	0x80, 0x46, 0x30, 0x43, 0xbb, 0xc8, 0xc0, 0x4a, 0x6a, 0xbc, 0x11, 0x50, 0x90, 0x4c, 0xc2, 0xbd,
	0xe0, 0xc9, 0x07, 0x26, 0x7f, 0xa0, 0xf3, 0xec, 0x8b, 0xd5, 0xd6, 0x4b, 0x45, 0x12, 0x1e, 0x9e,
	0x4a, 0x2f, 0x05, 0x99, 0x83, 0x88, 0xde, 0x80, 0xb2, 0xc1, 0xed, 0xc5, 0x61, 0xe4, 0x13, 0x74,
	0x87, 0x3f, 0xcc, 0xda, 0x7f, 0xcd, 0xe2, 0x93, 0x1f, 0x0d, 0x83, 0xcd, 0xb9, 0x36, 0x7c, 0x80,
	0x8b, 0x42, 0xfb, 0xd0, 0x9b, 0x4c, 0x6d, 0xd7, 0x62, 0x9e, 0xbb, 0x98, 0x5b, 0xde, 0xc3, 0xcc,
	0xb1, 0x5c, 0xf5, 0x0b, 0xed, 0x02, 0x38, 0xd3, 0xfb, 0xb9, 0x6d, 0x79, 0xe3, 0xd9, 0x42, 0x25,
	0xf5, 0xde, 0xb6, 0xd5, 0x06, 0xfd, 0x06, 0x32, 0xb3, 0xee, 0xac, 0x47, 0x94, 0xa5, 0x5a, 0x6b,
	0xdb, 0x6a, 0xd3, 0x58, 0x42, 0x0f, 0x83, 0x38, 0x59, 0x90, 0xae, 0x92, 0x30, 0x16, 0x61, 0x74,
	0xf8, 0xec, 0xf6, 0xfe, 0xc1, 0xd7, 0x22, 0x64, 0x75, 0x75, 0x9d, 0xfa, 0x30, 0xac, 0x12, 0x6f,
	0xff, 0x42, 0x3f, 0xf5, 0x5f, 0x44, 0x14, 0x89, 0xad, 0x99, 0xae, 0x77, 0x66, 0xf5, 0x7e, 0xe6,
	0x64, 0x09, 0x97, 0x27, 0x14, 0xb4, 0x91, 0x8f, 0xde, 0x03, 0x00, 0x00, 0xff, 0xff, 0x35, 0xef,
	0xe7, 0x89, 0x60, 0x02, 0x00, 0x00,
}
