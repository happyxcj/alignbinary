package alignbinary

import (
	"encoding/binary"
)

// msgEncoder describes how to encode the msg into the buf.
// The type of msg must be a basic type, a pointer of basic type,
// or a slice of basic type.
type msgEncoder func(msg interface{}, buf []byte, order binary.ByteOrder)

type encodeMsgInfo struct{}

func (encodeMsgInfo) bool(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(bool)
	buf[0] = boolToUint8(v)
}

func (encodeMsgInfo) boolPtr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*bool)
	buf[0] = boolToUint8(*v)
}

func (encodeMsgInfo) boolSlice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]bool)
	for i, v := range vs {
		buf[i] = boolToUint8(v)
	}
}

func (encodeMsgInfo) int8(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(int8)
	buf[0] = byte(v)
}

func (encodeMsgInfo) int8Ptr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*int8)
	buf[0] = byte(*v)
}

func (encodeMsgInfo) int8Slice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]int8)
	for i, v := range vs {
		buf[i] = byte(v)
	}
}

func (encodeMsgInfo) uint8(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(uint8)
	buf[0] = v
}

func (encodeMsgInfo) uint8Ptr(msg interface{}, buf []byte, _ binary.ByteOrder) {
	v := msg.(*uint8)
	buf[0] = *v
}

func (encodeMsgInfo) uint8Slice(msg interface{}, buf []byte, _ binary.ByteOrder) {
	vs := msg.([]uint8)
	for i, v := range vs {
		buf[i] = v
	}
}

func (encodeMsgInfo) int16(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(int16)
	order.PutUint16(buf, uint16(v))
}

func (encodeMsgInfo) int16Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int16)
	order.PutUint16(buf, uint16(*v))
}

func (encodeMsgInfo) int16Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int16)
	for i, v := range vs {
		order.PutUint16(buf[2*i:], uint16(v))
	}
}

func (encodeMsgInfo) uint16(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(uint16)
	order.PutUint16(buf, v)
}

func (encodeMsgInfo) uint16Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint16)
	order.PutUint16(buf, *v)
}

func (encodeMsgInfo) uint16Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint16)
	for i, v := range vs {
		order.PutUint16(buf[2*i:], v)
	}
}

func (encodeMsgInfo) int32(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(int32)
	order.PutUint32(buf, uint32(v))
}

func (encodeMsgInfo) int32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int32)
	order.PutUint32(buf, uint32(*v))
}

func (encodeMsgInfo) int32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int32)
	for i, v := range vs {
		order.PutUint32(buf[4*i:], uint32(v))
	}
}

func (encodeMsgInfo) uint32(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(uint32)
	order.PutUint32(buf, v)
}

func (encodeMsgInfo) uint32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint32)
	order.PutUint32(buf, *v)
}

func (encodeMsgInfo) uint32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint32)
	for i, v := range vs {
		order.PutUint32(buf[4*i:], v)
	}
}

func (encodeMsgInfo) int64(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(int64)
	order.PutUint64(buf, uint64(v))
}

func (encodeMsgInfo) int64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*int64)
	order.PutUint64(buf, uint64(*v))
}

func (encodeMsgInfo) int64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]int64)
	for i, v := range vs {
		order.PutUint64(buf[8*i:], uint64(v))
	}
}

func (encodeMsgInfo) uint64(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(uint64)
	order.PutUint64(buf, v)
}

func (encodeMsgInfo) uint64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*uint64)
	order.PutUint64(buf, *v)
}

func (encodeMsgInfo) uint64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]uint64)
	for i, v := range vs {
		order.PutUint64(buf[8*i:], v)
	}
}

func (encodeMsgInfo) float32(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(float32)
	order.PutUint32(buf, float32ToUint32(v))
}

func (encodeMsgInfo) float32Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*float32)
	order.PutUint32(buf, float32ToUint32(*v))
}

func (encodeMsgInfo) float32Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]float32)
	for i, v := range vs {
		order.PutUint32(buf[4*i:], float32ToUint32(v))
	}
}

func (encodeMsgInfo) float64(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(float64)
	order.PutUint64(buf, float64ToUint64(v))
}

func (encodeMsgInfo) float64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*float64)
	order.PutUint64(buf, float64ToUint64(*v))
}

func (encodeMsgInfo) float64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]float64)
	for i, v := range vs {
		order.PutUint64(buf[8*i:], float64ToUint64(v))
	}
}

func (encodeMsgInfo) complex64(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(complex64)
	x, y := complex64ToUint32s(v)
	order.PutUint32(buf, x)
	order.PutUint32(buf[4:], y)
}

func (encodeMsgInfo) complex64Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*complex64)
	x, y := complex64ToUint32s(*v)
	order.PutUint32(buf, x)
	order.PutUint32(buf[4:], y)
}

func (encodeMsgInfo) complex64Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]complex64)
	var start int
	var x, y uint32
	for i, v := range vs {
		start = 8 * i
		x, y = complex64ToUint32s(v)
		order.PutUint32(buf[start:], x)
		order.PutUint32(buf[start+4:], y)
	}
}

func (encodeMsgInfo) complex128(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(complex128)
	x, y := complex128ToUint64s(v)
	order.PutUint64(buf, x)
	order.PutUint64(buf[8:], y)
}

func (encodeMsgInfo) complex128Ptr(msg interface{}, buf []byte, order binary.ByteOrder) {
	v := msg.(*complex128)
	x, y := complex128ToUint64s(*v)
	order.PutUint64(buf, x)
	order.PutUint64(buf[8:], y)
}

func (encodeMsgInfo) complex128Slice(msg interface{}, buf []byte, order binary.ByteOrder) {
	vs := msg.([]complex128)
	var start int
	var x, y uint64
	for i, v := range vs {
		start = 16 * i
		x, y = complex128ToUint64s(v)
		order.PutUint64(buf[start:], x)
		order.PutUint64(buf[start+8:], y)
	}
}