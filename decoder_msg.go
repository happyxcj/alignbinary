package alignbinary

import (
	"encoding/binary"
)

// msgDecoder describes how to decode the msg from the buf.
// The type of msg must be a basic type pointer, or a basic type slice.
type msgDecoder func(msg interface{}, buf []byte, order binary.ByteOrder)

type decodeMsgInfo struct{}

func (decodeMsgInfo) boolPtr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*bool)
	*v = uint8ToBool(buf[0])
}

func (decodeMsgInfo) boolSlice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]bool)
	for i := range vs {
		vs[i] = uint8ToBool(buf[i])
	}
}

func (decodeMsgInfo) int8Ptr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*int8)
	*v = int8(buf[0])
}

func (decodeMsgInfo) int8Slice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]int8)
	for i := range vs {
		vs[i] = int8(buf[i])
	}
}

func (decodeMsgInfo) uint8Ptr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*uint8)
	*v = buf[0]
}

func (decodeMsgInfo) uint8Slice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]uint8)
	for i := range vs {
		vs[i] = buf[i]
	}
}

func (decodeMsgInfo) int16Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int16)
	*v = int16(order.Uint16(buf))
}

func (decodeMsgInfo) int16Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int16)
	for i := range vs {
		vs[i] = int16(order.Uint16(buf[2*i:]))
	}
}

func (decodeMsgInfo) uint16Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint16)
	*v = order.Uint16(buf)
}

func (decodeMsgInfo) uint16Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint16)
	for i := range vs {
		vs[i] = order.Uint16(buf[2*i:])
	}
}

func (decodeMsgInfo) int32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int32)
	*v = int32(order.Uint32(buf))
}

func (decodeMsgInfo) int32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int32)
	for i := range vs {
		vs[i] = int32(order.Uint32(buf[4*i:]))
	}
}

func (decodeMsgInfo) uint32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint32)
	*v = order.Uint32(buf)
}

func (decodeMsgInfo) uint32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint32)
	for i := range vs {
		vs[i] = order.Uint32(buf[4*i:])
	}
}

func (decodeMsgInfo) int64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int64)
	*v = int64(order.Uint64(buf))
}

func (decodeMsgInfo) int64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int64)
	for i := range vs {
		vs[i] = int64(order.Uint64(buf[8*i:]))
	}
}

func (decodeMsgInfo) uint64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint64)
	*v = order.Uint64(buf)
}

func (decodeMsgInfo) uint64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint64)
	for i := range vs {
		vs[i] = order.Uint64(buf[8*i:])
	}
}

func (decodeMsgInfo) float32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*float32)
	*v = uint32ToFloat32(order.Uint32(buf))
}

func (decodeMsgInfo) float32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]float32)
	for i := range vs {
		vs[i] = uint32ToFloat32(order.Uint32(buf[4*i:]))
	}
}

func (decodeMsgInfo) float64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*float64)
	*v = uint64ToFloat64(order.Uint64(buf))
}

func (decodeMsgInfo) float64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]float64)
	for i := range vs {
		vs[i] = uint64ToFloat64(order.Uint64(buf[8*i:]))
	}
}

func (decodeMsgInfo) complex64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*complex64)
	x := order.Uint32(buf)
	y := order.Uint32(buf[4:])
	*v = uint32sToComplex64(x, y)
}

func (decodeMsgInfo) complex64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]complex64)
	var start int
	var x, y uint32
	for i := range vs {
		start = 8 * i
		x = order.Uint32(buf[start:])
		y = order.Uint32(buf[start+4:])
		vs[i] = uint32sToComplex64(x, y)
	}
}

func (decodeMsgInfo) complex128Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*complex128)
	x := order.Uint64(buf)
	y := order.Uint64(buf[8:])
	*v = uint64sToComplex128(x, y)
}

func (decodeMsgInfo) complex128Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]complex128)
	var start int
	var x, y uint64
	for i := range vs {
		start = 16 * i
		x = order.Uint64(buf[start:])
		y = order.Uint64(buf[start+8:])
		vs[i] = uint64sToComplex128(x, y)
	}
}

