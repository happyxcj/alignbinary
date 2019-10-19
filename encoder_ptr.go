package alignbinary

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

// ptrEncoder describes how to encode the given ptr into the buf.
type ptrEncoder func(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder)

// encodeListInfo contains the information to encode a slice or array.
type encodeListInfo struct {
	// num is the number of all elements.
	num int
	// eleSize is the size of an element.
	eleSize    int
	eleEncoder ptrEncoder
}

func (li *encodeListInfo) init(v reflect.Value, eg *EncoderGroup) {
	li.num = v.Len()
	if li.num == 0 {
		return
	}
	li.eleEncoder,li.eleSize = eg.typePtrEncoder(v.Index(0))
	return
}

func (li *encodeListInfo) encode(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	var elePtr unsafe.Pointer
	var offset uintptr
	for i := 0; i < li.num; i++ {
		offset = uintptr(i * li.eleSize)
		elePtr = offsetPtr(ptr, offset)
		li.eleEncoder(elePtr, buf[offset:], order)
	}
}

type encodeStructInfo struct {
	// size is the size of the struct.
	size   int
	fields []*encodeFieldInfo
}

type encodeFieldInfo struct {
	// offset is the offset within struct, in bytes.
	// It's used to require the pointer that points to the field data.
	offset uintptr
	// start indicates the start index of the buf when decoding the field.
	start   int
	encoder ptrEncoder
}

// init initializes the information to encode the struct.
func (si *encodeStructInfo) init(v reflect.Value, eg *EncoderGroup) {
	t := v.Type()
	n := t.NumField()
	fields := make([]*encodeFieldInfo, 0, n)
	st := structTyp{}
	st.init(t, eg.af)
	// Calculates the decoders for all fields.
	for i := 0; i < n; i++ {
		if f, v := t.Field(i), v.Field(i); f.Name != "_" && v.CanSet() {
			e, _ := eg.typePtrEncoder(v)
			fi := &encodeFieldInfo{offset: f.Offset, start: int(st.fields[i]), encoder: e}
			fields = append(fields, fi)
		}
	}
	si.fields = fields
	si.size = int(st.size)
}

func (si *encodeStructInfo) encode(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	var fieldPtr unsafe.Pointer
	for _, f := range si.fields {
		fieldPtr = offsetPtr(ptr, f.offset)
		f.encoder(fieldPtr, buf[f.start:], order)
	}
}



type encodePtrInfo struct{}

func (encodePtrInfo) bool(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*bool)(ptr)
	buf[0] = boolToUint8(*v)
}

func (encodePtrInfo) int8(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*int8)(ptr)
	buf[0] = byte(*v)
}

func (encodePtrInfo) uint8(ptr unsafe.Pointer, buf []byte, _ binary.ByteOrder) {
	v := (*uint8)(ptr)
	buf[0] = *v
}

func (encodePtrInfo) int16(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int16)(ptr)
	order.PutUint16(buf, uint16(*v))
}

func (encodePtrInfo) uint16(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint16)(ptr)
	order.PutUint16(buf, *v)
}

func (encodePtrInfo) int32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int32)(ptr)
	order.PutUint32(buf, uint32(*v))
}

func (encodePtrInfo) uint32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint32)(ptr)
	order.PutUint32(buf, *v)
}

func (encodePtrInfo) int64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*int64)(ptr)
	order.PutUint64(buf, uint64(*v))
}

func (encodePtrInfo) uint64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*uint64)(ptr)
	order.PutUint64(buf, *v)
}

func (encodePtrInfo) float32(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*float32)(ptr)
	order.PutUint32(buf, float32ToUint32(*v))
}

func (encodePtrInfo) float64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*float64)(ptr)
	order.PutUint64(buf, float64ToUint64(*v))
}

func (encodePtrInfo) complex64(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*complex64)(ptr)
	x, y := complex64ToUint32s(*v)
	order.PutUint32(buf, x)
	order.PutUint32(buf[4:], y)
}

func (encodePtrInfo) complex128(ptr unsafe.Pointer, buf []byte, order binary.ByteOrder) {
	v := (*complex128)(ptr)
	x, y := complex128ToUint64s(*v)
	order.PutUint64(buf, x)
	order.PutUint64(buf[8:], y)
}
