package alignbinary

import (
	"io"
	"encoding/binary"
	"reflect"
	"unsafe"
	"fmt"
	"sync"
)

type DecoderGroup struct {
	af       AlignFactor
	structInfos sync.Map
	ptrInfo     decodePtrInfo
	msgInfo     decodeMsgInfo
}

func NewDecoderGroup(af AlignFactor) *DecoderGroup {
	checkAlignFactor(af)
	return &DecoderGroup{
		af: af,
	}
}

// Read reads structured binary data from r into msg.
// It can be called like 'binary.Read'
//
// Msg must be a pointer to a fixed-size value or a slice of fixed-size values.
//
// When decoding boolean values, a zero byte is decoded as false,
// and any other non-zero byte is decoded as true.
//
// When decoding into structs, the field data for unexported fields or
// fields with blank (_) field names is skipped.
func (dg *DecoderGroup) Read(r io.Reader, order binary.ByteOrder, msg interface{}) error {
	if decoder, size := dg.assertMsg(msg); size != -1 {
		// Fast path for a pointer to a basic type value, or a slice of basic type values.
		buf := make([]byte, size)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		decoder(msg, buf, order)
		return nil
	}
	// Decode by reflecting the msg.
	ptr, decoder, size := dg.reflectMsg(msg)
	buf := make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	decoder(ptr, buf, order)
	return nil
}

// Decode decodes the msg using the specified byte order and 
// the given structured binary data.
//
// Msg must be a pointer to a fixed-size value or a slice of fixed-size values.
//
// When decoding boolean values, a zero byte is decoded as false,
// and any other non-zero byte is decoded as true.
//
// When decoding into structs, the field data for unexported fields or
// fields with blank (_) field names is skipped.
func (dg *DecoderGroup) Decode(data []byte, order binary.ByteOrder, msg interface{}) error {
	if decoder, size := dg.assertMsg(msg); size != -1 {
		// Fast path for a pointer to a basic type value, or a slice of basic type values.
		if len(data) < size {
			return io.ErrUnexpectedEOF
		}
		decoder(msg, data, order)
		return nil
	}
	// Decode by reflecting the msg.
	ptr, decoder, size := dg.reflectMsg(msg)
	if len(data) < size {
		return io.ErrUnexpectedEOF
	}
	decoder(ptr, data, order)
	return nil
}

// assertMsg returns the message decoder and size by asserting the given msg.
// The type of msg must be a basic type pointer, or a basic type slice,
// if not, return nil and -1.
func (dg *DecoderGroup) assertMsg(msg interface{}) (msgDecoder, int) {
	switch v := msg.(type) {
	case *bool:
		return dg.msgInfo.boolPtr, 1
	case []bool:
		return dg.msgInfo.boolSlice, len(v)

	case *int8:
		return dg.msgInfo.int8Ptr, 1
	case []int8:
		return dg.msgInfo.int8Slice, len(v)
	case *uint8:
		return dg.msgInfo.uint8Ptr, 1
	case []uint8:
		return dg.msgInfo.uint8Slice, len(v)

	case *int16:
		return dg.msgInfo.int16Ptr, 2
	case []int16:
		return dg.msgInfo.int16Slice, 2 * len(v)
	case *uint16:
		return dg.msgInfo.uint16Ptr, 2
	case []uint16:
		return dg.msgInfo.uint16Slice, 2 * len(v)

	case *int32:
		return dg.msgInfo.int32Ptr, 4
	case []int32:
		return dg.msgInfo.int32Slice, 4 * len(v)
	case *uint32:
		return dg.msgInfo.uint32Ptr, 4
	case []uint32:
		return dg.msgInfo.uint32Slice, 4 * len(v)

	case *int64:
		return dg.msgInfo.int64Ptr, 8
	case []int64:
		return dg.msgInfo.int64Slice, 8 * len(v)
	case *uint64:
		return dg.msgInfo.uint64Ptr, 8
	case []uint64:
		return dg.msgInfo.uint64Slice, 8 * len(v)

	case *float32:
		return dg.msgInfo.float32Ptr, 4
	case []float32:
		return dg.msgInfo.float32Slice, 4 * len(v)

	case *float64:
		return dg.msgInfo.float64Ptr, 8
	case []float64:
		return dg.msgInfo.float64Slice, 8 * len(v)

	case *complex64:
		return dg.msgInfo.complex64Ptr, 8
	case []complex64:
		return dg.msgInfo.complex64Slice, 8 * len(v)

	case *complex128:
		return dg.msgInfo.complex128Ptr, 16
	case []complex128:
		return dg.msgInfo.complex128Slice, 16 * len(v)

	}
	return nil, -1
}

func (dg *DecoderGroup) reflectMsg(msg interface{}) (unsafe.Pointer, ptrDecoder, int) {
	v := reflect.ValueOf(msg)
	kind := v.Kind()
	var decoder ptrDecoder
	var size int
	if kind == reflect.Slice {
		info := new(decodeListInfo)
		info.init(v, dg)
		decoder, size = info.decode, info.num*info.eleSize
	} else {
		if kind != reflect.Ptr {
			panic(fmt.Sprintf("alignbinary: call reflectMsg on invalid kind %v when decoding", kind))
		}
		decoder, size = dg.typePtrDecoder(v.Elem())
	}
	return unsafe.Pointer(v.Pointer()), decoder, size
}

// typePtrDecoder returns the size and decoder based on the given t under the align.
// It panics if t's Kind is not Array, Struct, Bool,
// Int8, Uint8, Int16, Uint16, Int32, Uint32, Int64, Uint64,
// Float32, Float64, Complex64, or Complex128.
func (dg *DecoderGroup) typePtrDecoder(v reflect.Value) (ptrDecoder, int) {
	switch v.Kind() {
	case reflect.Array:
		info := new(decodeListInfo)
		info.init(v, dg)
		return info.decode, info.num * info.eleSize
	case reflect.Struct:
		info := dg.getDecodeStructInfo(v)
		return info.decode, info.size
	case reflect.Bool:
		return dg.ptrInfo.bool, 1
	case reflect.Int8:
		return dg.ptrInfo.int8, 1
	case reflect.Uint8:
		return dg.ptrInfo.uint8, 1

	case reflect.Int16:
		return dg.ptrInfo.int16, 2
	case reflect.Uint16:
		return dg.ptrInfo.uint16, 2

	case reflect.Int32:
		return dg.ptrInfo.int32, 4
	case reflect.Uint32:
		return dg.ptrInfo.uint32, 4

	case reflect.Int64:
		return dg.ptrInfo.uint64, 8
	case reflect.Uint64:
		return dg.ptrInfo.uint64, 8

	case reflect.Float32:
		return dg.ptrInfo.float32, 4
	case reflect.Float64:
		return dg.ptrInfo.float64, 8
	case reflect.Complex64:
		return dg.ptrInfo.complex64, 8
	case reflect.Complex128:
		return dg.ptrInfo.complex128, 16
	}
	panic(fmt.Sprintf("alignbinary: call typePtrDecoder on invalid kind %v", v.Kind()))
}

func (dg *DecoderGroup) getDecodeStructInfo(v reflect.Value) *decodeStructInfo {
	val, ok := dg.structInfos.Load(v.Type())
	if ok {
		return val.(*decodeStructInfo)
	}
	info := new(decodeStructInfo)
	info.init(v, dg)
	dg.structInfos.Store(v.Type(), info)
	return info
}

//func (dg *DecoderGroup) getDecodeStructInfo(v reflect.Value) *decodeStructInfo {
//	addr := (*unsafe.Pointer)(unsafe.Pointer(&dg.structInfos))
//	val := atomic.LoadPointer(addr)
//	if val == nil {
//		// Create a new one and initialize the cache.
//		info := new(decodeStructInfo)
//		info.init(v, dg)
//		newInfos := map[reflect.Type]*decodeStructInfo{v.Type(): info}
//		atomic.StorePointer(addr, unsafe.Pointer(&newInfos))
//		return info
//	}
//	infos := *(*map[reflect.Type]*decodeStructInfo)(val)
//	if info, ok := infos[v.Type()]; ok {
//		// Get from the cache.
//		return info
//	}
//	// Create a new one when there is no cache for the specified type.
//	newInfos := make(map[reflect.Type]*decodeStructInfo, len(infos)+1)
//	// Copy the exist cache.
//	for key, info := range infos {
//		newInfos[key] = info
//	}
//	info := new(decodeStructInfo)
//	info.init(v, dg)
//	newInfos[v.Type()] = info
//	atomic.StorePointer(addr, unsafe.Pointer(&newInfos))
//	return info
//}
