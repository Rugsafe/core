// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmwasm/will/types.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// WillInfo is the structure that represents the will
type Will struct {
	// will_id is the unique identifier of the will
	will_id uint64 `protobuf:"varint,1,opt,name=will_id,json=willId,proto3" json:"will_id,omitempty"`
	// WillName is the user_generated name of the will
	Will_name string `protobuf:"bytes,2,opt,name=will_name,json=willName,proto3" json:"will_name,omitempty"`
	// WillBeneficiary is the private key or address, depending on the purpose of
	Will_benefiary string `protobuf:"bytes,3,opt,name=will_beneficiary,json=willBeneficiary,proto3" json:"will_beneficiary,omitempty"`
}

func (m *Will) Reset()         { *m = Will{} }
func (m *Will) String() string { return proto.CompactTextString(m) }
func (*Will) ProtoMessage()    {}
func (*Will) Descriptor() ([]byte, []int) {
	return fileDescriptor_cec37ad7aa1ffe0b, []int{0}
}

func (m *Will) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *Will) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Will.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *Will) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Will.Merge(m, src)
}

func (m *Will) XXX_Size() int {
	return m.Size()
}

func (m *Will) XXX_DiscardUnknown() {
	xxx_messageInfo_Will.DiscardUnknown(m)
}

var xxx_messageInfo_Will proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Will)(nil), "cosmwasm.will.Will")
}

func init() { proto.RegisterFile("cosmwasm/will/types.proto", fileDescriptor_cec37ad7aa1ffe0b) }

var fileDescriptor_cec37ad7aa1ffe0b = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4c, 0xce, 0x2f, 0xce,
	0x2d, 0x4f, 0x2c, 0xce, 0xd5, 0x2f, 0xcf, 0xcc, 0xc9, 0xd1, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x85, 0x49, 0xe9, 0x81, 0xa4, 0xa4, 0x44, 0xd2, 0xf3,
	0xd3, 0xf3, 0xc1, 0x32, 0xfa, 0x20, 0x16, 0x44, 0x91, 0xd2, 0x5c, 0x46, 0x2e, 0x96, 0xf0, 0xcc,
	0x9c, 0x1c, 0x21, 0x15, 0x2e, 0x76, 0x90, 0xb2, 0xf8, 0xcc, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0x16, 0x27, 0xee, 0x47, 0xf7, 0xe4, 0x61, 0x42, 0x41, 0x6c, 0x20, 0x86, 0x67, 0x8a, 0x90, 0x16,
	0x17, 0x27, 0x58, 0x28, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x49, 0x81, 0x51, 0x83, 0xd3, 0x89, 0xf7,
	0xd1, 0x3d, 0x79, 0xce, 0x70, 0x98, 0x60, 0x10, 0x07, 0x48, 0xde, 0x2f, 0x31, 0x37, 0x55, 0xc8,
	0x96, 0x4b, 0x00, 0xac, 0x36, 0x29, 0x35, 0x2f, 0x35, 0x2d, 0x33, 0x39, 0x33, 0xb1, 0xa8, 0x52,
	0x82, 0x19, 0xac, 0x45, 0xe8, 0xd1, 0x3d, 0x79, 0xbe, 0x70, 0x84, 0x5c, 0x62, 0x51, 0x65, 0x10,
	0x3f, 0x48, 0xad, 0x13, 0x42, 0xa9, 0x15, 0xcb, 0x8b, 0x05, 0xf2, 0x8c, 0x4e, 0x1e, 0x27, 0x1e,
	0xca, 0x31, 0xac, 0x78, 0x24, 0xc7, 0x78, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f,
	0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c,
	0x51, 0x6a, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xce, 0xf9, 0xc5,
	0xb9, 0xe1, 0xe0, 0xc0, 0x48, 0x2c, 0xce, 0x4d, 0xd1, 0xaf, 0x40, 0x0a, 0x94, 0x24, 0x36, 0xb0,
	0x87, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7e, 0x98, 0x4e, 0xbd, 0x32, 0x01, 0x00, 0x00,
}

func (this *Will) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Will)
	if !ok {
		that2, ok := that.(Will)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.will_id != that1.will_id {
		return false
	}
	if this.Will_name != that1.Will_name {
		return false
	}
	if this.Will_benefiary != that1.Will_benefiary {
		return false
	}
	return true
}

func (m *Will) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Will) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Will) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Will_benefiary) > 0 {
		i -= len(m.Will_benefiary)
		copy(dAtA[i:], m.Will_benefiary)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Will_benefiary)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Will_name) > 0 {
		i -= len(m.Will_name)
		copy(dAtA[i:], m.Will_name)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Will_name)))
		i--
		dAtA[i] = 0x12
	}
	if m.will_id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.will_id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func (m *Will) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.will_id != 0 {
		n += 1 + sovTypes(uint64(m.will_id))
	}
	l = len(m.Will_name)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Will_benefiary)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func (m *Will) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Will: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Will: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field will_id", wireType)
			}
			m.will_id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.will_id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Will_name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Will_name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Will_benefiary", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Will_benefiary = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
