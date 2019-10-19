package alignbinary

/*
#include <stdio.h>
#define arrayLen 4
//#pragma pack(4)
typedef struct{
	unsigned char Bool;
	unsigned char BoolArray[arrayLen];

	char Int8;
	char Int8Array[arrayLen];

	short Int16;
	short Int16Array[arrayLen];

	int Int32;
	int Int32Array[arrayLen];

	long long Int64;
	long long Int64Array[arrayLen];

	unsigned char Uint8;
	unsigned char Uint8Array[arrayLen];

	unsigned short Uint16;
	unsigned short Uint16Array[arrayLen];

	unsigned int Uint32;
	unsigned int Uint32Array[arrayLen];

	unsigned long long Uint64;
	unsigned long long Uint64Array[arrayLen];

	float Float32;
	float Float32Array[arrayLen];

	double Float64;
	double Float64Array[arrayLen];

	float Complex64[2];
	float Complex64Array[arrayLen][2];

	double Complex128[2];
	double Complex128Array[arrayLen][2];
}Struct;

 */
import "C"
import (
	"math"
	"testing"
	"encoding/binary"
	"reflect"
)

const arrayLen = 4
const sliceLen = 3

type Struct struct {
	Bool      bool
	BoolArray [arrayLen]bool

	Int8      int8
	Int8Array [arrayLen]int8

	Int16      int16
	Int16Array [arrayLen]int16

	Int32      int32
	Int32Array [arrayLen]int32

	Int64      int64
	Int64Array [arrayLen]int64

	Uint8      uint8
	Uint8Array [arrayLen]uint8

	Uint16      uint16
	Uint16Array [arrayLen]uint16

	Uint32      uint32
	Uint32Array [arrayLen]uint32

	Uint64      uint64
	Uint64Array [arrayLen]uint64

	Float32      float32
	Float32Array [arrayLen]float32

	Float64      float64
	Float64Array [arrayLen]float64

	Complex64      complex64
	Complex64Array [arrayLen]complex64

	Complex128      complex128
	Complex128Array [arrayLen]complex128
}

var goStruct = Struct{
	true,
	[arrayLen]bool{true, false, true, false},

	0x01,
	[arrayLen]int8{0x02, 0x03, 0x04, 0x05},

	0x0102,
	[arrayLen]int16{0x0103, 0x0104, 0x0105, 0x0106},

	0x01020304,
	[arrayLen]int32{0x01020305, 0x01020306, 0x01020307, 0x01020308},

	0x0102030405060708,
	[arrayLen]int64{0x0102030405060709, 0x010203040506070a, 0x010203040506070b, 0x010203040506070c},

	0x01,
	[arrayLen]uint8{0x02, 0x03, 0x04, 0x05},

	0x0102,
	[arrayLen]uint16{0x0103, 0x0104, 0x0105, 0x0106},

	0x01020304,
	[arrayLen]uint32{0x01020305, 0x01020306, 0x01020307, 0x01020308},

	0x0102030405060708,
	[arrayLen]uint64{0x0102030405060709, 0x010203040506070a, 0x010203040506070b, 0x010203040506070c},

	math.Float32frombits(0x01020304),
	[arrayLen]float32{
		math.Float32frombits(0x01020305),
		math.Float32frombits(0x01020306),
		math.Float32frombits(0x01020307),
		math.Float32frombits(0x01020308)},

	math.Float64frombits(0x0102030405060708),
	[arrayLen]float64{
		math.Float64frombits(0x0102030405060709),
		math.Float64frombits(0x010203040506070a),
		math.Float64frombits(0x010203040506070b),
		math.Float64frombits(0x010203040506070c)},

	complex(math.Float32frombits(0x01020304), math.Float32frombits(0x01020305)),
	[arrayLen]complex64{
		complex(math.Float32frombits(0x01020306), math.Float32frombits(0x01020307)),
		complex(math.Float32frombits(0x01020308), math.Float32frombits(0x01020309)),
		complex(math.Float32frombits(0x0102030a), math.Float32frombits(0x0102030b)),
		complex(math.Float32frombits(0x0102030c), math.Float32frombits(0x0102030d))},

	complex(math.Float64frombits(0x0102030405060708), math.Float64frombits(0x0102030405060709)),
	[arrayLen]complex128{
		complex(math.Float64frombits(0x010203040506070a), math.Float64frombits(0x010203040506070b)),
		complex(math.Float64frombits(0x010203040506070c), math.Float64frombits(0x010203040506070d)),
		complex(math.Float64frombits(0x010203040506070e), math.Float64frombits(0x010203040506070f)),
		complex(math.Float64frombits(0x0102030405060710), math.Float64frombits(0x0102030405060711))},
}

var goStructSlice = []Struct{goStruct, goStruct, goStruct}

var cgoStruct = C.Struct{
	Bool:      C.uchar(1),
	BoolArray: [arrayLen]C.uchar{C.uchar(1), C.uchar(0), C.uchar(1), C.uchar(0)},

	Int8:      0x01,
	Int8Array: [arrayLen]C.char{0x02, 0x03, 0x04, 0x05},

	Int16:      0x0102,
	Int16Array: [arrayLen]C.short{0x0103, 0x0104, 0x0105, 0x0106},

	Int32:      0x01020304,
	Int32Array: [arrayLen]C.int{0x01020305, 0x01020306, 0x01020307, 0x01020308},

	Int64:      0x0102030405060708,
	Int64Array: [arrayLen]C.longlong{0x0102030405060709, 0x010203040506070a, 0x010203040506070b, 0x010203040506070c},

	Uint8:      0x01,
	Uint8Array: [arrayLen]C.uchar{0x02, 0x03, 0x04, 0x05},

	Uint16:      0x0102,
	Uint16Array: [arrayLen]C.ushort{0x0103, 0x0104, 0x0105, 0x0106},

	Uint32:      0x01020304,
	Uint32Array: [arrayLen]C.uint{0x01020305, 0x01020306, 0x01020307, 0x01020308},

	Uint64:      0x0102030405060708,
	Uint64Array: [arrayLen]C.ulonglong{0x0102030405060709, 0x010203040506070a, 0x010203040506070b, 0x010203040506070c},

	Float32: C.float(math.Float32frombits(0x01020304)),
	Float32Array: [arrayLen]C.float{
		C.float(math.Float32frombits(0x01020305)),
		C.float(math.Float32frombits(0x01020306)),
		C.float(math.Float32frombits(0x01020307)),
		C.float(math.Float32frombits(0x01020308))},

	Float64: C.double(math.Float64frombits(0x0102030405060708)),
	Float64Array: [arrayLen]C.double{
		C.double(math.Float64frombits(0x0102030405060709)),
		C.double(math.Float64frombits(0x010203040506070a)),
		C.double(math.Float64frombits(0x010203040506070b)),
		C.double(math.Float64frombits(0x010203040506070c))},

	Complex64: [2]C.float{C.float(math.Float32frombits(0x01020304)), C.float(math.Float32frombits(0x01020305))},
	Complex64Array: [arrayLen][2]C.float{
		{C.float(math.Float32frombits(0x01020306)), C.float(math.Float32frombits(0x01020307))},
		{C.float(math.Float32frombits(0x01020308)), C.float(math.Float32frombits(0x01020309))},
		{C.float(math.Float32frombits(0x0102030a)), C.float(math.Float32frombits(0x0102030b))},
		{C.float(math.Float32frombits(0x0102030c)), C.float(math.Float32frombits(0x0102030d))}},

	Complex128: [2]C.double{C.double(math.Float64frombits(0x0102030405060708)), C.double(math.Float64frombits(0x0102030405060709))},
	Complex128Array: [arrayLen][2]C.double{
		{C.double(math.Float64frombits(0x010203040506070a)), C.double(math.Float64frombits(0x010203040506070b))},
		{C.double(math.Float64frombits(0x010203040506070c)), C.double(math.Float64frombits(0x010203040506070d))},
		{C.double(math.Float64frombits(0x010203040506070e)), C.double(math.Float64frombits(0x010203040506070f))},
		{C.double(math.Float64frombits(0x0102030405060710)), C.double(math.Float64frombits(0x0102030405060711))}},
}

var cgoStructSlice = []C.Struct{cgoStruct, cgoStruct, cgoStruct}

func checkResult(t *testing.T, method string, order binary.ByteOrder, err error, have, want interface{}) {
	if err != nil {
		t.Errorf("%v %v: %v", method, order, err)
		return
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("%v %v:\n\thave %+v\n\twant %+v", method, order, have, want)
	}
}
