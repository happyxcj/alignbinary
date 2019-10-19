package alignbinary

import (
	"reflect"
)

const (
	AlignDefault = 0
	Align1Byte   = 1 << 0
	Align2Byte   = 1 << 1
	Align4Byte   = 1 << 2
	Align8Byte   = 1 << 3
)

type AlignFactor uint8

type structTyp struct {
	// size is the size in bytes of a specified struct.
	size uintptr
	// fields is the offsets of all fields within struct
	// based on a specified alignment factor, in bytes.
	fields []uintptr
}

// init calculates the size and fields of st by the given t.
// It panics if t'Kind is not Struct.
func (st *structTyp) init(t reflect.Type, af AlignFactor) {
	n := t.NumField()
	fields := make([]uintptr, n)
	if af == AlignDefault {
		// Fast path to initialize the information of struct fields.
		// It can avoid the repeated calculation of struct fields.
		for i := 0; i < n; i++ {
			fields[i] = t.Field(i).Offset
		}
		st.size = t.Size()
		st.fields = fields
		return
	}
	// Calculate the information of all fields based on the af.
	st.size, _, st.fields = calcStructSizeAlign(t, af)
}

// calcStructSizeAlign calculates and returns the size, alignment and offsets of fields
// for a struct type based on the given af.
//
// The Calculation refers to the:
// 1. 'StructOf' method in file ../reflect/type.go.
// 2. 'widstruct' method in file ../cmd/compile/internal/gc/align.go.
func calcStructSizeAlign(t reflect.Type, af AlignFactor) (uintptr, uint8, []uintptr) {
	n := t.NumField()
	fields := make([]uintptr, n)
	var f reflect.StructField
	var size uintptr
	// Minimum alignment for a struct is 1 byte.
	var typeAlign uint8 = 1
	lastZero := uintptr(0)
	for i := 0; i < n; i++ {
		f = t.Field(i)
		fSize, fAlign := calcSizeAlign(f.Type, af)
		if fAlign > typeAlign {
			// Reset the alignment for of the t.
			typeAlign = fAlign
		}
		var offset uintptr
		if fAlign > 0 {
			// Calculate the offset for the field.
			offset = align(size, uintptr(fAlign))
		} else {
			offset = size
		}
		size = offset + fSize // Reset the size of the t.
		if fSize == 0 {
			lastZero = size
		}
		fields[i] = offset
	}
	if size > 0 && lastZero == size {
		// This is a non-zero sized struct that ends in a
		// zero-sized field. We add an extra byte of padding,
		// to ensure that taking the address of the final
		// zero-sized field can't manufacture a pointer to the
		// next object in the heap. See issue 9401.
		size++
	}
	// Round the size up to be a multiple of the alignment.
	size = align(size, uintptr(typeAlign))
	return size, typeAlign, fields
}

// calcSizeAlign calculates and returns the size and alignment for t based on the given af.
func calcSizeAlign(t reflect.Type, af AlignFactor) (uintptr, uint8) {
	switch t.Kind() {
	case reflect.Array:
		size, align := calcSizeAlign(t.Elem(), af)
		return size * uintptr(t.Len()), align
	case reflect.Struct:
		size, align, _ := calcStructSizeAlign(t, af)
		return size, align
	default:
		// Calculate the valid alignment for the type.
		// We ignore the other types (i.e. Slice, Map) besides the basic types,
		// because the process will panic when decoding or encoding the message
		// if the type is invalid for this library.
		size := t.Size() // The size must be a power of two.
		var align uint8
		if size <= uintptr(af) {
			align = uint8(size)
		} else {
			align = uint8(af)
		}
		return size, align
	}
}

// align returns the result of rounding x up to a multiple of n.
// n must be a power of two.
func align(x, n uintptr) uintptr {
	return (x + n - 1) &^ (n - 1)
}
