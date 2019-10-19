package alignbinary

import (
	"encoding/binary"
	"testing"
	"bytes"
	"reflect"
)

// TODO
//go test -bench=. -benchmem -run=none
var order = binary.LittleEndian
//var order = binary.BigEndian

func TestWrite(t *testing.T) {
	srcBuf := &bytes.Buffer{}
	binary.Write(srcBuf, order, cgoStruct)

	tarBuf := &bytes.Buffer{}
	err := Write(tarBuf, order, goStruct)
	checkResult(t, "TestWrite", order, err, tarBuf.Bytes(), srcBuf.Bytes())
}

func TestWriteSlice(t *testing.T) {
	srcBuf := &bytes.Buffer{}
	binary.Write(srcBuf, order, cgoStructSlice)

	tarBuf := &bytes.Buffer{}
	err := Write(tarBuf, order, goStructSlice)
	checkResult(t, "TestWriteSlice", order, err, tarBuf, srcBuf)
}

func TestRead(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStruct)

	val := Struct{}
	err := Read(bytes.NewReader(buf.Bytes()), order, &val)
	checkResult(t, "TestRead", order, err, val, goStruct)
}

func TestReadSlice(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStructSlice)

	val := make([]Struct, sliceLen)
	err := Read(bytes.NewReader(buf.Bytes()), order, val)
	checkResult(t, "TestReadSlice", order, err, val, goStructSlice)
}

//=========================================== Benchmark =======================

func BenchmarkBinaryWrite(b *testing.B) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStruct)
	b.SetBytes(int64(len(buf.Bytes())))
	b.ResetTimer()
	var tarBuf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		tarBuf = &bytes.Buffer{}
		binary.Write(tarBuf, order, cgoStruct)
	}
	b.StopTimer()
}

func BenchmarkWrite(b *testing.B) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStruct)
	b.SetBytes(int64(len(buf.Bytes())))
	b.ResetTimer()
	var tarBuf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		tarBuf = &bytes.Buffer{}
		Write(tarBuf, order, goStruct)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(buf, tarBuf) {
		b.Fatalf("buf doesn't match:\ngot  %v;\nwant %v", tarBuf, buf)
	}
}

func BenchmarkBinaryRead(b *testing.B) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStruct)
	b.SetBytes(int64(len(buf.Bytes())))
	b.ResetTimer()
	val := Struct{}
	for i := 0; i < b.N; i++ {
		binary.Read(buf, order, &val)
	}
	b.StopTimer()
}

func BenchmarkRead(b *testing.B) {
	buf := &bytes.Buffer{}
	binary.Write(buf, order, cgoStruct)
	b.SetBytes(int64(len(buf.Bytes())))
	b.ResetTimer()
	val := Struct{}
	for i := 0; i < b.N; i++ {
		Read(buf, order, &val)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(goStruct, val) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", val, goStruct)
	}
}
