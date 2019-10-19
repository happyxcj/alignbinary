package alignbinary

import (
	"math"
	"unsafe"
	"fmt"
)

func boolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

func uint8ToBool(v uint8) bool {
	return v != 0
}

func float32ToUint32(v float32) uint32 {
	return math.Float32bits(v)
}

func uint32ToFloat32(v uint32) float32 {
	return math.Float32frombits(v)
}

func float64ToUint64(v float64) uint64 {
	return math.Float64bits(v)
}

func uint64ToFloat64(v uint64) float64 {
	return math.Float64frombits(v)
}

func complex64ToUint32s(v complex64) (x, y uint32) {
	x = float32ToUint32(float32(real(v)))
	y = float32ToUint32(float32(imag(v)))
	return
}

func uint32sToComplex64(x, y uint32) complex64 {
	return complex(math.Float32frombits(x), math.Float32frombits(y))
}

func complex128ToUint64s(v complex128) (x, y uint64) {
	x = float64ToUint64(float64(real(v)))
	y = float64ToUint64(float64(imag(v)))
	return
}

func uint64sToComplex128(x, y uint64) complex128 {
	return complex(math.Float64frombits(x), math.Float64frombits(y))
}

// offsetPtr returns a unsafe.Pointer points to a field in a struct or an element of an array:
//
//	// equivalent to f := unsafe.Pointer(&s.f)
//	f := unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f))
//
//	// equivalent to e := unsafe.Pointer(&x[i])
//	e := unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + i*unsafe.Sizeof(x[0]))
func offsetPtr(ptr unsafe.Pointer, off uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + off)
}

func checkAlignFactor(af AlignFactor) {
	if af > 8 || af&(af-1) != 0 {
		panic(fmt.Sprintf("alignbinary: invalid alignment factor: %v",af))
	}
}
