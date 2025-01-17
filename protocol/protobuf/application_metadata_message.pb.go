// Code generated by protoc-gen-go. DO NOT EDIT.
// source: application_metadata_message.proto

package protobuf

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

type ApplicationMetadataMessage_Type int32

const (
	ApplicationMetadataMessage_UNKNOWN                                 ApplicationMetadataMessage_Type = 0
	ApplicationMetadataMessage_CHAT_MESSAGE                            ApplicationMetadataMessage_Type = 1
	ApplicationMetadataMessage_CONTACT_UPDATE                          ApplicationMetadataMessage_Type = 2
	ApplicationMetadataMessage_MEMBERSHIP_UPDATE_MESSAGE               ApplicationMetadataMessage_Type = 3
	ApplicationMetadataMessage_PAIR_INSTALLATION                       ApplicationMetadataMessage_Type = 4
	ApplicationMetadataMessage_SYNC_INSTALLATION                       ApplicationMetadataMessage_Type = 5
	ApplicationMetadataMessage_REQUEST_ADDRESS_FOR_TRANSACTION         ApplicationMetadataMessage_Type = 6
	ApplicationMetadataMessage_ACCEPT_REQUEST_ADDRESS_FOR_TRANSACTION  ApplicationMetadataMessage_Type = 7
	ApplicationMetadataMessage_DECLINE_REQUEST_ADDRESS_FOR_TRANSACTION ApplicationMetadataMessage_Type = 8
	ApplicationMetadataMessage_REQUEST_TRANSACTION                     ApplicationMetadataMessage_Type = 9
	ApplicationMetadataMessage_SEND_TRANSACTION                        ApplicationMetadataMessage_Type = 10
	ApplicationMetadataMessage_DECLINE_REQUEST_TRANSACTION             ApplicationMetadataMessage_Type = 11
	ApplicationMetadataMessage_SYNC_INSTALLATION_CONTACT               ApplicationMetadataMessage_Type = 12
	ApplicationMetadataMessage_SYNC_INSTALLATION_ACCOUNT               ApplicationMetadataMessage_Type = 13
	ApplicationMetadataMessage_SYNC_INSTALLATION_PUBLIC_CHAT           ApplicationMetadataMessage_Type = 14
	ApplicationMetadataMessage_CONTACT_CODE_ADVERTISEMENT              ApplicationMetadataMessage_Type = 15
	ApplicationMetadataMessage_PUSH_NOTIFICATION_REGISTRATION          ApplicationMetadataMessage_Type = 16
	ApplicationMetadataMessage_PUSH_NOTIFICATION_REGISTRATION_RESPONSE ApplicationMetadataMessage_Type = 17
	ApplicationMetadataMessage_PUSH_NOTIFICATION_QUERY                 ApplicationMetadataMessage_Type = 18
	ApplicationMetadataMessage_PUSH_NOTIFICATION_QUERY_RESPONSE        ApplicationMetadataMessage_Type = 19
	ApplicationMetadataMessage_PUSH_NOTIFICATION_REQUEST               ApplicationMetadataMessage_Type = 20
	ApplicationMetadataMessage_PUSH_NOTIFICATION_RESPONSE              ApplicationMetadataMessage_Type = 21
	ApplicationMetadataMessage_EMOJI_REACTION                          ApplicationMetadataMessage_Type = 22
	ApplicationMetadataMessage_GROUP_CHAT_INVITATION                   ApplicationMetadataMessage_Type = 23
	ApplicationMetadataMessage_CHAT_IDENTITY                           ApplicationMetadataMessage_Type = 24
	ApplicationMetadataMessage_COMMUNITY_DESCRIPTION                   ApplicationMetadataMessage_Type = 25
	ApplicationMetadataMessage_COMMUNITY_INVITATION                    ApplicationMetadataMessage_Type = 26
	ApplicationMetadataMessage_COMMUNITY_REQUEST_TO_JOIN               ApplicationMetadataMessage_Type = 27
	ApplicationMetadataMessage_PIN_MESSAGE                             ApplicationMetadataMessage_Type = 28
	ApplicationMetadataMessage_EDIT_MESSAGE                            ApplicationMetadataMessage_Type = 29
	ApplicationMetadataMessage_STATUS_UPDATE                           ApplicationMetadataMessage_Type = 30
	ApplicationMetadataMessage_DELETE_MESSAGE                          ApplicationMetadataMessage_Type = 31
	ApplicationMetadataMessage_SYNC_INSTALLATION_COMMUNITY             ApplicationMetadataMessage_Type = 32
	ApplicationMetadataMessage_ANONYMOUS_METRIC_BATCH                  ApplicationMetadataMessage_Type = 33
	ApplicationMetadataMessage_SYNC_CHAT_REMOVED                       ApplicationMetadataMessage_Type = 34
	ApplicationMetadataMessage_SYNC_CHAT_MESSAGES_READ                 ApplicationMetadataMessage_Type = 35
	ApplicationMetadataMessage_BACKUP                                  ApplicationMetadataMessage_Type = 36
	ApplicationMetadataMessage_SYNC_ACTIVITY_CENTER_READ               ApplicationMetadataMessage_Type = 37
	ApplicationMetadataMessage_SYNC_ACTIVITY_CENTER_ACCEPTED           ApplicationMetadataMessage_Type = 38
	ApplicationMetadataMessage_SYNC_ACTIVITY_CENTER_DISMISSED          ApplicationMetadataMessage_Type = 39
	ApplicationMetadataMessage_SYNC_BOOKMARK                           ApplicationMetadataMessage_Type = 40
	ApplicationMetadataMessage_SYNC_CLEAR_HISTORY                      ApplicationMetadataMessage_Type = 41
	ApplicationMetadataMessage_SYNC_SETTING                            ApplicationMetadataMessage_Type = 42
	ApplicationMetadataMessage_COMMUNITY_ARCHIVE_MAGNETLINK            ApplicationMetadataMessage_Type = 43
	ApplicationMetadataMessage_SYNC_PROFILE_PICTURE                    ApplicationMetadataMessage_Type = 44
	ApplicationMetadataMessage_SYNC_WALLET_ACCOUNT                     ApplicationMetadataMessage_Type = 45
	ApplicationMetadataMessage_ACCEPT_CONTACT_REQUEST                  ApplicationMetadataMessage_Type = 46
	ApplicationMetadataMessage_RETRACT_CONTACT_REQUEST                 ApplicationMetadataMessage_Type = 47
	ApplicationMetadataMessage_COMMUNITY_REQUEST_TO_JOIN_RESPONSE      ApplicationMetadataMessage_Type = 48
	ApplicationMetadataMessage_SYNC_COMMUNITY_SETTINGS                 ApplicationMetadataMessage_Type = 49
	ApplicationMetadataMessage_REQUEST_CONTACT_VERIFICATION            ApplicationMetadataMessage_Type = 50
	ApplicationMetadataMessage_ACCEPT_CONTACT_VERIFICATION             ApplicationMetadataMessage_Type = 51
	ApplicationMetadataMessage_DECLINE_CONTACT_VERIFICATION            ApplicationMetadataMessage_Type = 52
	ApplicationMetadataMessage_SYNC_TRUSTED_USER                       ApplicationMetadataMessage_Type = 53
	ApplicationMetadataMessage_SYNC_VERIFICATION_REQUEST               ApplicationMetadataMessage_Type = 54
	ApplicationMetadataMessage_SYNC_CONTACT_REQUEST_DECISION           ApplicationMetadataMessage_Type = 56
	ApplicationMetadataMessage_COMMUNITY_REQUEST_TO_LEAVE              ApplicationMetadataMessage_Type = 57
	ApplicationMetadataMessage_SYNC_DELETE_FOR_ME_MESSAGE              ApplicationMetadataMessage_Type = 58
	ApplicationMetadataMessage_SYNC_SAVED_ADDRESS                      ApplicationMetadataMessage_Type = 59
	ApplicationMetadataMessage_COMMUNITY_CANCEL_REQUEST_TO_JOIN        ApplicationMetadataMessage_Type = 60
	ApplicationMetadataMessage_CANCEL_CONTACT_VERIFICATION             ApplicationMetadataMessage_Type = 61
)

var ApplicationMetadataMessage_Type_name = map[int32]string{
	0:  "UNKNOWN",
	1:  "CHAT_MESSAGE",
	2:  "CONTACT_UPDATE",
	3:  "MEMBERSHIP_UPDATE_MESSAGE",
	4:  "PAIR_INSTALLATION",
	5:  "SYNC_INSTALLATION",
	6:  "REQUEST_ADDRESS_FOR_TRANSACTION",
	7:  "ACCEPT_REQUEST_ADDRESS_FOR_TRANSACTION",
	8:  "DECLINE_REQUEST_ADDRESS_FOR_TRANSACTION",
	9:  "REQUEST_TRANSACTION",
	10: "SEND_TRANSACTION",
	11: "DECLINE_REQUEST_TRANSACTION",
	12: "SYNC_INSTALLATION_CONTACT",
	13: "SYNC_INSTALLATION_ACCOUNT",
	14: "SYNC_INSTALLATION_PUBLIC_CHAT",
	15: "CONTACT_CODE_ADVERTISEMENT",
	16: "PUSH_NOTIFICATION_REGISTRATION",
	17: "PUSH_NOTIFICATION_REGISTRATION_RESPONSE",
	18: "PUSH_NOTIFICATION_QUERY",
	19: "PUSH_NOTIFICATION_QUERY_RESPONSE",
	20: "PUSH_NOTIFICATION_REQUEST",
	21: "PUSH_NOTIFICATION_RESPONSE",
	22: "EMOJI_REACTION",
	23: "GROUP_CHAT_INVITATION",
	24: "CHAT_IDENTITY",
	25: "COMMUNITY_DESCRIPTION",
	26: "COMMUNITY_INVITATION",
	27: "COMMUNITY_REQUEST_TO_JOIN",
	28: "PIN_MESSAGE",
	29: "EDIT_MESSAGE",
	30: "STATUS_UPDATE",
	31: "DELETE_MESSAGE",
	32: "SYNC_INSTALLATION_COMMUNITY",
	33: "ANONYMOUS_METRIC_BATCH",
	34: "SYNC_CHAT_REMOVED",
	35: "SYNC_CHAT_MESSAGES_READ",
	36: "BACKUP",
	37: "SYNC_ACTIVITY_CENTER_READ",
	38: "SYNC_ACTIVITY_CENTER_ACCEPTED",
	39: "SYNC_ACTIVITY_CENTER_DISMISSED",
	40: "SYNC_BOOKMARK",
	41: "SYNC_CLEAR_HISTORY",
	42: "SYNC_SETTING",
	43: "COMMUNITY_ARCHIVE_MAGNETLINK",
	44: "SYNC_PROFILE_PICTURE",
	45: "SYNC_WALLET_ACCOUNT",
	46: "ACCEPT_CONTACT_REQUEST",
	47: "RETRACT_CONTACT_REQUEST",
	48: "COMMUNITY_REQUEST_TO_JOIN_RESPONSE",
	49: "SYNC_COMMUNITY_SETTINGS",
	50: "REQUEST_CONTACT_VERIFICATION",
	51: "ACCEPT_CONTACT_VERIFICATION",
	52: "DECLINE_CONTACT_VERIFICATION",
	53: "SYNC_TRUSTED_USER",
	54: "SYNC_VERIFICATION_REQUEST",
	56: "SYNC_CONTACT_REQUEST_DECISION",
	57: "COMMUNITY_REQUEST_TO_LEAVE",
	58: "SYNC_DELETE_FOR_ME_MESSAGE",
	59: "SYNC_SAVED_ADDRESS",
	60: "COMMUNITY_CANCEL_REQUEST_TO_JOIN",
	61: "CANCEL_CONTACT_VERIFICATION",
}

var ApplicationMetadataMessage_Type_value = map[string]int32{
	"UNKNOWN":                                 0,
	"CHAT_MESSAGE":                            1,
	"CONTACT_UPDATE":                          2,
	"MEMBERSHIP_UPDATE_MESSAGE":               3,
	"PAIR_INSTALLATION":                       4,
	"SYNC_INSTALLATION":                       5,
	"REQUEST_ADDRESS_FOR_TRANSACTION":         6,
	"ACCEPT_REQUEST_ADDRESS_FOR_TRANSACTION":  7,
	"DECLINE_REQUEST_ADDRESS_FOR_TRANSACTION": 8,
	"REQUEST_TRANSACTION":                     9,
	"SEND_TRANSACTION":                        10,
	"DECLINE_REQUEST_TRANSACTION":             11,
	"SYNC_INSTALLATION_CONTACT":               12,
	"SYNC_INSTALLATION_ACCOUNT":               13,
	"SYNC_INSTALLATION_PUBLIC_CHAT":           14,
	"CONTACT_CODE_ADVERTISEMENT":              15,
	"PUSH_NOTIFICATION_REGISTRATION":          16,
	"PUSH_NOTIFICATION_REGISTRATION_RESPONSE": 17,
	"PUSH_NOTIFICATION_QUERY":                 18,
	"PUSH_NOTIFICATION_QUERY_RESPONSE":        19,
	"PUSH_NOTIFICATION_REQUEST":               20,
	"PUSH_NOTIFICATION_RESPONSE":              21,
	"EMOJI_REACTION":                          22,
	"GROUP_CHAT_INVITATION":                   23,
	"CHAT_IDENTITY":                           24,
	"COMMUNITY_DESCRIPTION":                   25,
	"COMMUNITY_INVITATION":                    26,
	"COMMUNITY_REQUEST_TO_JOIN":               27,
	"PIN_MESSAGE":                             28,
	"EDIT_MESSAGE":                            29,
	"STATUS_UPDATE":                           30,
	"DELETE_MESSAGE":                          31,
	"SYNC_INSTALLATION_COMMUNITY":             32,
	"ANONYMOUS_METRIC_BATCH":                  33,
	"SYNC_CHAT_REMOVED":                       34,
	"SYNC_CHAT_MESSAGES_READ":                 35,
	"BACKUP":                                  36,
	"SYNC_ACTIVITY_CENTER_READ":               37,
	"SYNC_ACTIVITY_CENTER_ACCEPTED":           38,
	"SYNC_ACTIVITY_CENTER_DISMISSED":          39,
	"SYNC_BOOKMARK":                           40,
	"SYNC_CLEAR_HISTORY":                      41,
	"SYNC_SETTING":                            42,
	"COMMUNITY_ARCHIVE_MAGNETLINK":            43,
	"SYNC_PROFILE_PICTURE":                    44,
	"SYNC_WALLET_ACCOUNT":                     45,
	"ACCEPT_CONTACT_REQUEST":                  46,
	"RETRACT_CONTACT_REQUEST":                 47,
	"COMMUNITY_REQUEST_TO_JOIN_RESPONSE":      48,
	"SYNC_COMMUNITY_SETTINGS":                 49,
	"REQUEST_CONTACT_VERIFICATION":            50,
	"ACCEPT_CONTACT_VERIFICATION":             51,
	"DECLINE_CONTACT_VERIFICATION":            52,
	"SYNC_TRUSTED_USER":                       53,
	"SYNC_VERIFICATION_REQUEST":               54,
	"SYNC_CONTACT_REQUEST_DECISION":           56,
	"COMMUNITY_REQUEST_TO_LEAVE":              57,
	"SYNC_DELETE_FOR_ME_MESSAGE":              58,
	"SYNC_SAVED_ADDRESS":                      59,
	"COMMUNITY_CANCEL_REQUEST_TO_JOIN":        60,
	"CANCEL_CONTACT_VERIFICATION":             61,
}

func (x ApplicationMetadataMessage_Type) String() string {
	return proto.EnumName(ApplicationMetadataMessage_Type_name, int32(x))
}

func (ApplicationMetadataMessage_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ad09a6406fcf24c7, []int{0, 0}
}

type ApplicationMetadataMessage struct {
	// Signature of the payload field
	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// This is the encoded protobuf of the application level message, i.e ChatMessage
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	// The type of protobuf message sent
	Type                 ApplicationMetadataMessage_Type `protobuf:"varint,3,opt,name=type,proto3,enum=protobuf.ApplicationMetadataMessage_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *ApplicationMetadataMessage) Reset()         { *m = ApplicationMetadataMessage{} }
func (m *ApplicationMetadataMessage) String() string { return proto.CompactTextString(m) }
func (*ApplicationMetadataMessage) ProtoMessage()    {}
func (*ApplicationMetadataMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad09a6406fcf24c7, []int{0}
}

func (m *ApplicationMetadataMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplicationMetadataMessage.Unmarshal(m, b)
}
func (m *ApplicationMetadataMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplicationMetadataMessage.Marshal(b, m, deterministic)
}
func (m *ApplicationMetadataMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplicationMetadataMessage.Merge(m, src)
}
func (m *ApplicationMetadataMessage) XXX_Size() int {
	return xxx_messageInfo_ApplicationMetadataMessage.Size(m)
}
func (m *ApplicationMetadataMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplicationMetadataMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ApplicationMetadataMessage proto.InternalMessageInfo

func (m *ApplicationMetadataMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *ApplicationMetadataMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ApplicationMetadataMessage) GetType() ApplicationMetadataMessage_Type {
	if m != nil {
		return m.Type
	}
	return ApplicationMetadataMessage_UNKNOWN
}

func init() {
	proto.RegisterEnum("protobuf.ApplicationMetadataMessage_Type", ApplicationMetadataMessage_Type_name, ApplicationMetadataMessage_Type_value)
	proto.RegisterType((*ApplicationMetadataMessage)(nil), "protobuf.ApplicationMetadataMessage")
}

func init() {
	proto.RegisterFile("application_metadata_message.proto", fileDescriptor_ad09a6406fcf24c7)
}

var fileDescriptor_ad09a6406fcf24c7 = []byte{
	// 908 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0x5b, 0x73, 0x13, 0x37,
	0x14, 0x6e, 0x20, 0x4d, 0x40, 0xb9, 0xa0, 0x88, 0x5c, 0x9c, 0xbb, 0x31, 0x34, 0x04, 0x68, 0x4d,
	0x0b, 0x6d, 0xa7, 0x2d, 0xe5, 0x41, 0x96, 0x4e, 0x6c, 0xe1, 0x5d, 0x69, 0x91, 0xb4, 0x66, 0xdc,
	0x17, 0xcd, 0x52, 0x5c, 0x26, 0x33, 0x40, 0x3c, 0xc4, 0x3c, 0xe4, 0x9f, 0xf6, 0x57, 0xf4, 0x37,
	0x74, 0xb4, 0x57, 0x27, 0xd9, 0x90, 0xa7, 0x64, 0xcf, 0xf7, 0xe9, 0x48, 0xe7, 0x3b, 0xdf, 0x39,
	0x46, 0xad, 0x64, 0x3c, 0xfe, 0x70, 0xfc, 0x77, 0x32, 0x39, 0x3e, 0xf9, 0xe4, 0x3e, 0x8e, 0x26,
	0xc9, 0xbb, 0x64, 0x92, 0xb8, 0x8f, 0xa3, 0xd3, 0xd3, 0xe4, 0xfd, 0xa8, 0x3d, 0xfe, 0x7c, 0x32,
	0x39, 0x21, 0xb7, 0xd2, 0x3f, 0x6f, 0xbf, 0xfc, 0xd3, 0xfa, 0x6f, 0x19, 0x6d, 0xd1, 0xea, 0x40,
	0x98, 0xf3, 0xc3, 0x8c, 0x4e, 0x76, 0xd0, 0xed, 0xd3, 0xe3, 0xf7, 0x9f, 0x92, 0xc9, 0x97, 0xcf,
	0xa3, 0xc6, 0x4c, 0x73, 0xe6, 0x70, 0x51, 0x57, 0x01, 0xd2, 0x40, 0xf3, 0xe3, 0xe4, 0xec, 0xc3,
	0x49, 0xf2, 0xae, 0x71, 0x23, 0xc5, 0x8a, 0x4f, 0xf2, 0x12, 0xcd, 0x4e, 0xce, 0xc6, 0xa3, 0xc6,
	0xcd, 0xe6, 0xcc, 0xe1, 0xf2, 0xb3, 0x47, 0xed, 0xe2, 0xbe, 0xf6, 0xd5, 0x77, 0xb5, 0xed, 0xd9,
	0x78, 0xa4, 0xd3, 0x63, 0xad, 0x7f, 0x97, 0xd0, 0xac, 0xff, 0x24, 0x0b, 0x68, 0x3e, 0x96, 0x7d,
	0xa9, 0xde, 0x48, 0xfc, 0x0d, 0xc1, 0x68, 0x91, 0xf5, 0xa8, 0x75, 0x21, 0x18, 0x43, 0xbb, 0x80,
	0x67, 0x08, 0x41, 0xcb, 0x4c, 0x49, 0x4b, 0x99, 0x75, 0x71, 0xc4, 0xa9, 0x05, 0x7c, 0x83, 0xec,
	0xa2, 0xcd, 0x10, 0xc2, 0x0e, 0x68, 0xd3, 0x13, 0x51, 0x1e, 0x2e, 0x8f, 0xdc, 0x24, 0x6b, 0x68,
	0x25, 0xa2, 0x42, 0x3b, 0x21, 0x8d, 0xa5, 0x41, 0x40, 0xad, 0x50, 0x12, 0xcf, 0xfa, 0xb0, 0x19,
	0x4a, 0x76, 0x3e, 0xfc, 0x2d, 0xb9, 0x8f, 0xf6, 0x35, 0xbc, 0x8e, 0xc1, 0x58, 0x47, 0x39, 0xd7,
	0x60, 0x8c, 0x3b, 0x52, 0xda, 0x59, 0x4d, 0xa5, 0xa1, 0x2c, 0x25, 0xcd, 0x91, 0xc7, 0xe8, 0x80,
	0x32, 0x06, 0x91, 0x75, 0xd7, 0x71, 0xe7, 0xc9, 0x13, 0xf4, 0x90, 0x03, 0x0b, 0x84, 0x84, 0x6b,
	0xc9, 0xb7, 0xc8, 0x06, 0xba, 0x5b, 0x90, 0xa6, 0x81, 0xdb, 0x64, 0x15, 0x61, 0x03, 0x92, 0x9f,
	0x8b, 0x22, 0xb2, 0x8f, 0xb6, 0x2f, 0xe6, 0x9e, 0x26, 0x2c, 0x78, 0x69, 0x2e, 0x15, 0xe9, 0x72,
	0x01, 0xf1, 0x62, 0x3d, 0x4c, 0x19, 0x53, 0xb1, 0xb4, 0x78, 0x89, 0xdc, 0x43, 0xbb, 0x97, 0xe1,
	0x28, 0xee, 0x04, 0x82, 0x39, 0xdf, 0x17, 0xbc, 0x4c, 0xf6, 0xd0, 0x56, 0xd1, 0x0f, 0xa6, 0x38,
	0x38, 0xca, 0x07, 0xa0, 0xad, 0x30, 0x10, 0x82, 0xb4, 0xf8, 0x0e, 0x69, 0xa1, 0xbd, 0x28, 0x36,
	0x3d, 0x27, 0x95, 0x15, 0x47, 0x82, 0x65, 0x29, 0x34, 0x74, 0x85, 0xb1, 0x3a, 0x93, 0x1c, 0x7b,
	0x85, 0xbe, 0xce, 0x71, 0x1a, 0x4c, 0xa4, 0xa4, 0x01, 0xbc, 0x42, 0xb6, 0xd1, 0xc6, 0x65, 0xf2,
	0xeb, 0x18, 0xf4, 0x10, 0x13, 0xf2, 0x00, 0x35, 0xaf, 0x00, 0xab, 0x14, 0x77, 0x7d, 0xd5, 0x75,
	0xf7, 0xa5, 0xfa, 0xe1, 0x55, 0x5f, 0x52, 0x1d, 0x9c, 0x1f, 0x5f, 0xf3, 0x16, 0x84, 0x50, 0xbd,
	0x12, 0x4e, 0x43, 0xae, 0xf3, 0x3a, 0xd9, 0x44, 0x6b, 0x5d, 0xad, 0xe2, 0x28, 0x95, 0xc5, 0x09,
	0x39, 0x10, 0x36, 0xab, 0x6e, 0x83, 0xac, 0xa0, 0xa5, 0x2c, 0xc8, 0x41, 0x5a, 0x61, 0x87, 0xb8,
	0xe1, 0xd9, 0x4c, 0x85, 0x61, 0x2c, 0x85, 0x1d, 0x3a, 0x0e, 0x86, 0x69, 0x11, 0xa5, 0xec, 0x4d,
	0xd2, 0x40, 0xab, 0x15, 0x34, 0x95, 0x67, 0xcb, 0xbf, 0xba, 0x42, 0xca, 0x6e, 0x2b, 0xf7, 0x4a,
	0x09, 0x89, 0xb7, 0xc9, 0x1d, 0xb4, 0x10, 0x09, 0x59, 0xda, 0x7e, 0xc7, 0xcf, 0x0e, 0x70, 0x51,
	0xcd, 0xce, 0xae, 0x7f, 0x89, 0xb1, 0xd4, 0xc6, 0xa6, 0x18, 0x9d, 0x3d, 0x5f, 0x0b, 0x87, 0x00,
	0xa6, 0xe6, 0x65, 0xdf, 0x9b, 0xaa, 0xce, 0x33, 0xf9, 0xd5, 0xb8, 0x49, 0xb6, 0xd0, 0x3a, 0x95,
	0x4a, 0x0e, 0x43, 0x15, 0x1b, 0x17, 0x82, 0xd5, 0x82, 0xb9, 0x0e, 0xb5, 0xac, 0x87, 0xef, 0x95,
	0x53, 0x95, 0x96, 0xac, 0x21, 0x54, 0x03, 0xe0, 0xb8, 0xe5, 0xbb, 0x56, 0x85, 0xf3, 0xab, 0x8c,
	0x17, 0x90, 0xe3, 0xfb, 0x04, 0xa1, 0xb9, 0x0e, 0x65, 0xfd, 0x38, 0xc2, 0x0f, 0x4a, 0x47, 0x7a,
	0x65, 0x07, 0xbe, 0x52, 0x06, 0xd2, 0x82, 0xce, 0xa8, 0xdf, 0x95, 0x8e, 0xbc, 0x08, 0x67, 0xd3,
	0x08, 0x1c, 0x1f, 0x78, 0xc7, 0xd5, 0x52, 0xb8, 0x30, 0xa1, 0x30, 0x06, 0x38, 0x7e, 0x98, 0x2a,
	0xe1, 0x39, 0x1d, 0xa5, 0xfa, 0x21, 0xd5, 0x7d, 0x7c, 0x48, 0xd6, 0x11, 0xc9, 0x5e, 0x18, 0x00,
	0xd5, 0xae, 0x27, 0x8c, 0x55, 0x7a, 0x88, 0x1f, 0x79, 0x19, 0xd3, 0xb8, 0x01, 0x6b, 0x85, 0xec,
	0xe2, 0xc7, 0xa4, 0x89, 0x76, 0xaa, 0x46, 0x50, 0xcd, 0x7a, 0x62, 0x00, 0x2e, 0xa4, 0x5d, 0x09,
	0x36, 0x10, 0xb2, 0x8f, 0x9f, 0xf8, 0x26, 0xa6, 0x67, 0x22, 0xad, 0x8e, 0x44, 0x00, 0x2e, 0x12,
	0xcc, 0xc6, 0x1a, 0xf0, 0xf7, 0x7e, 0xbe, 0x53, 0xe4, 0x0d, 0x0d, 0x02, 0xb0, 0xe5, 0xa8, 0xfd,
	0x90, 0x6a, 0x9a, 0x6d, 0x94, 0x62, 0x9c, 0x0a, 0x43, 0xb6, 0xbd, 0x78, 0x1a, 0xac, 0xce, 0x66,
	0xec, 0x3c, 0xf8, 0x94, 0x1c, 0xa0, 0xd6, 0x95, 0xb6, 0xa8, 0x5c, 0xfb, 0x63, 0xd5, 0x81, 0x92,
	0x9c, 0x57, 0x64, 0xf0, 0x4f, 0xbe, 0xa4, 0xe2, 0x68, 0x71, 0xc3, 0x00, 0x74, 0xe9, 0x7e, 0xfc,
	0xcc, 0x9b, 0xe2, 0xc2, 0xfb, 0xce, 0x11, 0x9e, 0xfb, 0x14, 0xc5, 0x2a, 0xaa, 0x65, 0xfc, 0x5c,
	0x5a, 0xc3, 0xea, 0xd8, 0x58, 0xe0, 0x2e, 0x36, 0xa0, 0xf1, 0x2f, 0x65, 0xc7, 0xa7, 0xd9, 0x65,
	0x7d, 0xbf, 0x96, 0x1d, 0xbf, 0x50, 0xb9, 0xe3, 0xc0, 0x84, 0xf1, 0x89, 0x7f, 0xcb, 0x76, 0x50,
	0x8d, 0x04, 0x01, 0xd0, 0x01, 0xe0, 0xdf, 0x3d, 0x9e, 0xa6, 0xc8, 0x9d, 0xee, 0xb7, 0x6e, 0x58,
	0x19, 0xfe, 0x8f, 0xb2, 0xf5, 0x86, 0x0e, 0x80, 0x17, 0xcb, 0x19, 0xbf, 0xf0, 0xdb, 0xa4, 0xca,
	0xcb, 0xa8, 0x64, 0x10, 0x5c, 0x1a, 0xbc, 0x3f, 0xbd, 0x32, 0x39, 0x56, 0x5b, 0xf7, 0xcb, 0xce,
	0xd2, 0x5f, 0x0b, 0xed, 0xa7, 0x2f, 0x8a, 0xdf, 0xc3, 0xb7, 0x73, 0xe9, 0x7f, 0xcf, 0xff, 0x0f,
	0x00, 0x00, 0xff, 0xff, 0x84, 0x77, 0xdc, 0x48, 0xb6, 0x07, 0x00, 0x00,
}
