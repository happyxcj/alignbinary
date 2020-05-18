package alignbinary

import (
	"encoding/binary"
	"reflect"
	"unsafe"
	"fmt"
	"io"
	"sync"
)

type EncoderGroup struct {
	af AlignFactor
	structInfos sync.Map
	ptrInfo     encodePtrInfo
	msgInfo     encodeMsgInfo
}

func NewEncoderGroup(af AlignFactor) *EncoderGroup {

	checkAlignFactor(af)

	return &EncoderGroup{
		af: af,
	}
}

// Encode writes the binary representation of msg into w.
// It can be called like 'binary.Write'
func (eg *EncoderGroup) Write(w io.Writer, order binary.ByteOrder, msg interface{}) error {
	data, _ := defaultEG.Encode(order, msg)
	_, err := w.Write(data)
	return err
}

// Encode encodes the msg and returns the binary representation of msg.
//
// Msg must be a fixed-size value, a pointer to a fixed-size value,
// or a slice of fixed-size values.
//
// Boolean values encode as one byte: 1 for true, and 0 for false.
//
// Bytes to be returned are encoded using the specified byte order
// and read from successive fields of the msg.
//
// When encoding structs, zero values are encoded for unexported fields or
// fields with blank (_) field names.
func (eg *EncoderGroup) Encode(order binary.ByteOrder, msg interface{}) ([]byte, error) {
	if encoder, size := eg.assertMsg(msg); size != -1 {
		// Fast path for a basic type value, or a slice of basic type values.
		buf := make([]byte, size)
		encoder(msg, buf, order)
		return buf, nil
	}
	// Encode by reflecting the msg.
	ptr, encoder, size := eg.reflectMsg(msg)
	buf := make([]byte, size)
	encoder(ptr, buf, order)
	return buf, nil
}

// assertMsg returns the message encoder and size by asserting the given msg.
// The type of msg must be a basic type, a pointer of basic type,
// or a slice of basic type, if not, return nil and -1.
func (eg *EncoderGroup) assertMsg(msg interface{}) (msgEncoder, int) {
	msgInfo := eg.msgInfo
	switch v := msg.(type) {
	case bool:
		return msgInfo.bool, 1
	case *bool:
		return msgInfo.boolPtr, 1
	case []bool:
		return msgInfo.boolSlice, len(v)
	case int8:
		return msgInfo.int8, 1
	case *int8:
		return msgInfo.int8Ptr, 1
	case []int8:
		return msgInfo.int8Slice, len(v)
	case uint8:
		return msgInfo.uint8, 1
	case *uint8:
		return msgInfo.uint8Ptr, 1
	case []uint8:
		return msgInfo.uint8Slice, len(v)

	case int16:
		return msgInfo.int16, 2
	case *int16:
		return msgInfo.int16Ptr, 2
	case []int16:
		return msgInfo.int16Slice, 2 * len(v)
	case uint16:
		return msgInfo.uint16, 2
	case *uint16:
		return msgInfo.uint16Ptr, 2
	case []uint16:
		return msgInfo.uint16Slice, 2 * len(v)

	case int32:
		return msgInfo.int32, 4
	case *int32:
		return msgInfo.int32Ptr, 4
	case []int32:
		return msgInfo.int32Slice, 4 * len(v)
	case uint32:
		return msgInfo.uint32, 4
	case *uint32:
		return msgInfo.uint32Ptr, 4
	case []uint32:
		return msgInfo.uint32Slice, 4 * len(v)

	case int64:
		return msgInfo.int64, 8
	case *int64:
		return msgInfo.int64Ptr, 8
	case []int64:
		return msgInfo.int64Slice, 8 * len(v)
	case uint64:
		return msgInfo.uint64, 8
	case *uint64:
		return msgInfo.uint64Ptr, 8
	case []uint64:
		return msgInfo.uint64Slice, 8 * len(v)

	case float32:
		return msgInfo.float32, 4
	case *float32:
		return msgInfo.float32Ptr, 4
	case []float32:
		return msgInfo.float32Slice, 4 * len(v)

	case float64:
		return msgInfo.float64, 8
	case *float64:
		return msgInfo.float64Ptr, 8
	case []float64:
		return msgInfo.float64Slice, 8 * len(v)

	case complex64:
		return msgInfo.complex64, 8
	case *complex64:
		return msgInfo.complex64Ptr, 8
	case []complex64:
		return msgInfo.complex64Slice, 8 * len(v)

	case complex128:
		return msgInfo.complex128, 16
	case *complex128:
		return msgInfo.complex128Ptr, 16
	case []complex128:
		return msgInfo.complex128Slice, 16 * len(v)

	}
	return nil, -1
}

// reflectMsg returns the message pointer, pointer encoder and message size by reflecting the given msg.
// It panics if msg isn't a fixed-size value, a pointer to a fixed-size value,
// or a slice of fixed-size values.
func (eg *EncoderGroup) reflectMsg(msg interface{}) (unsafe.Pointer, ptrEncoder, int) {
	var encoder ptrEncoder
	var size int
	v := reflect.ValueOf(msg)
	kind := v.Kind()
	if kind == reflect.Slice {
		info := new(encodeListInfo)
		info.init(v, eg)
		encoder, size = info.encode, info.eleSize*info.num
	} else {
		if kind != reflect.Ptr {
			// Convert to a pointer that points to the data of msg interface.
			u := reflect.New(v.Type())
			u.Elem().Set(v)
			v = u
		}
		encoder, size = eg.typePtrEncoder(v.Elem())
	}
	return unsafe.Pointer(v.Pointer()), encoder, size
}

// typePtrEncoder returns the pointer encoder and message size based on the given v.
// It panics if v's Kind is not Array, Struct, Bool,
// Int8, Uint8, Int16, Uint16, Int32, Uint32, Int64, Uint64,
// Float32, Float64, Complex64, or Complex128.
func (eg *EncoderGroup) typePtrEncoder(v reflect.Value) (ptrEncoder, int) {
	switch v.Kind() {
	case reflect.Array:
		info := new(encodeListInfo)
		info.init(v, eg)
		return info.encode, info.num * info.eleSize
	case reflect.Struct:
		info := eg.getEncodeStructInfo(v)
		return info.encode, info.size
	case reflect.Bool:
		return eg.ptrInfo.bool, 1
	case reflect.Int8:
		return eg.ptrInfo.int8, 1
	case reflect.Uint8:
		return eg.ptrInfo.uint8, 1

	case reflect.Int16:
		return eg.ptrInfo.int16, 2
	case reflect.Uint16:
		return eg.ptrInfo.uint16, 2

	case reflect.Int32:
		return eg.ptrInfo.int32, 4
	case reflect.Uint32:
		return eg.ptrInfo.uint32, 4

	case reflect.Int64:
		return eg.ptrInfo.int64, 8
	case reflect.Uint64:
		return eg.ptrInfo.uint64, 8

	case reflect.Float32:
		return eg.ptrInfo.float32, 4
	case reflect.Float64:
		return eg.ptrInfo.float64, 8

	case reflect.Complex64:
		return eg.ptrInfo.complex64, 8
	case reflect.Complex128:
		return eg.ptrInfo.complex128, 16
	}
	panic(fmt.Sprintf("alignbinary: call typePtrEncoder on invalid kind %v", v.Kind()))
}

func (eg *EncoderGroup) getEncodeStructInfo(v reflect.Value) *encodeStructInfo {
	val, ok := eg.structInfos.Load(v.Type())
	if ok {
		return val.(*encodeStructInfo)
	}
	info := new(encodeStructInfo)
	info.init(v, eg)
	eg.structInfos.Store(v.Type(), info)
	return info
}

//func (eg *EncoderGroup) getEncodeStructInfo(v reflect.Value) *encodeStructInfo {
//	addr := (*unsafe.Pointer)(unsafe.Pointer(&eg.structInfos))
//	val := atomic.LoadPointer(addr)
//	if val == nil {
//		// Create a new one and initialize the cache.
//		info := new(encodeStructInfo)
//		info.init(v, eg)
//		newInfos := map[reflect.Type]*encodeStructInfo{v.Type(): info}
//		atomic.StorePointer(addr, unsafe.Pointer(&newInfos))
//		return info
//	}
//	infos := *(*map[reflect.Type]*encodeStructInfo)(val)
//	if info, ok := infos[v.Type()]; ok {
//		// Get from the cache.
//		return info
//	}
//	// Create a new one when there is no cache for the specified type.
//	newInfos := make(map[reflect.Type]*encodeStructInfo, len(infos)+1)
//	// Copy the exist cache.
//	for key, info := range infos {
//		newInfos[key] = info
//	}
//	info := new(encodeStructInfo)
//	info.init(v, eg)
//	newInfos[v.Type()] = info
//	atomic.StorePointer(addr, unsafe.Pointer(&newInfos))
//	return info
//}
