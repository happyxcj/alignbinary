package alignbinary

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

// ptrDecoder describes how to decode the given ptr from the buf.
type ptrDecoder func(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder)

type decodeListInfo struct {
	// num is the number of all elements.
	num int
	// eleSize is the size of an element.
	eleSize    int
	eleDecoder ptrDecoder
}

// init initializes the information to decode a slice or a array
func (li *decodeListInfo) init(v reflect.Value, dg *DecoderGroup) {
	li.num = v.Len()
	if li.num == 0 {
		return
	}
	li.eleDecoder, li.eleSize = dg.typePtrDecoder(v.Index(0))
}

func (li *decodeListInfo) decode(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	var elePtr unsafe.Pointer
	var offset uintptr
	for i := 0; i < li.num; i++ {
		offset = uintptr(i * li.eleSize)
		elePtr = offsetPtr(ptr, offset)
		li.eleDecoder(elePtr, buf[offset:], order)
	}
}

type decodeStructInfo struct {
	// size is the size of the struct.
	size   int
	fields []*decodeFieldInfo
}

type decodeFieldInfo struct {
	// offset is the offset within struct, in bytes.
	// It's used to require the pointer that points to the field data.
	offset uintptr
	// start indicates the start index of the buf in decoder when decoding the field.
	start   int
	decoder ptrDecoder
}

// init initializes the information to decode a struct.
func (si *decodeStructInfo) init(v reflect.Value,dg *DecoderGroup) {
	t := v.Type()
	n := t.NumField()
	fields := make([]*decodeFieldInfo, 0, n)
	st := structTyp{}
	st.init(t, dg.af)
	// Calculates the decoders for all fields.
	for i := 0; i < n; i++ {
		if f, v := t.Field(i), v.Field(i); f.Name != "_" && v.CanSet() {
			d, _ := dg.typePtrDecoder(v)
			fi := &decodeFieldInfo{offset: f.Offset, start: int(st.fields[i]), decoder: d}
			fields = append(fields, fi)
		}
	}
	si.fields = fields
	si.size = int(st.size)
}

func (si *decodeStructInfo) decode(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	var fieldPtr unsafe.Pointer
	for _, f := range si.fields {
		fieldPtr = offsetPtr(ptr, f.offset)
		f.decoder(fieldPtr, buf[f.start:], order)
	}
}

type decodePtrInfo struct{}

func (decodePtrInfo) bool(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*bool)(ptr)
	*v = uint8ToBool(buf[0])
}

func (decodePtrInfo) int8(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*int8)(ptr)
	*v = int8(buf[0])
}

func (decodePtrInfo) uint8(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*uint8)(ptr)
	*v = buf[0]
}

func (decodePtrInfo) int16(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int16)(ptr)
	*v = int16(order.Uint16(buf))
}

func (decodePtrInfo) uint16(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint16)(ptr)
	*v = order.Uint16(buf)
}

func (decodePtrInfo) int32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int32)(ptr)
	*v = int32(order.Uint32(buf))
}

func (decodePtrInfo) uint32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint32)(ptr)
	*v = order.Uint32(buf)
}

func (decodePtrInfo) int64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int64)(ptr)
	*v = int64(order.Uint64(buf))
}

func (decodePtrInfo) uint64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint64)(ptr)
	*v = order.Uint64(buf)
}

func (decodePtrInfo) float32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*float32)(ptr)
	*v = uint32ToFloat32(order.Uint32(buf))
}

func (decodePtrInfo) float64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*float64)(ptr)
	*v = uint64ToFloat64(order.Uint64(buf))
}

func (decodePtrInfo) complex64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*complex64)(ptr)
	x := order.Uint32(buf)
	y := order.Uint32(buf[4:])
	*v = uint32sToComplex64(x, y)
}

func (decodePtrInfo) complex128(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*complex128)(ptr)
	x := order.Uint64(buf)
	y := order.Uint64(buf[8:])
	*v = uint64sToComplex128(x, y)
}
