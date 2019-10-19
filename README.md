# Aligned Binary

alignbinary is a binary codec that implements byte alignment for Go (golang). Just like the package binary in the language standard library, it can only serializing and deserializing the fixed-size values. A fixed-size value is either a fixed-size arithmetic type (bool, int8, uint8, int16, float32, complex64, ...) or an array or struct containing only fixed-size values.

## Features

+ As easy to learn and use as the package binary in the  standard library (See [Quick Start](#quick-start).
+ High efficiency for struct (See [Benchmark](#benchmark)).
+ Optional alignment factor (e.g., 1, 2, 4, 8).
+ Compatible with `binary.Write` / `binary.Read`, just choose the alignment  factor: `alignbinary.Align1Byte`  (See [Examples](#examples)).

## Install

```
go get -u github.com/happyxcj/alignbinary
```

## Quick Start

```go
package main

import (
	"github.com/happyxcj/alignbinary"
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	// Encode
	msg := [3]int16{1, 2, 3}
	buf := new(bytes.Buffer)
	err := alignbinary.Write(buf, binary.LittleEndian, msg)
	if err != nil {
		fmt.Println("alignbinary.Write error: ", err.Error())
		return
	}
	fmt.Println("buf=", buf.Bytes())

	// Decode
	readBuf := bytes.NewReader(buf.Bytes())
	var writeMsg [3]int16
	err = alignbinary.Read(readBuf, binary.LittleEndian, &writeMsg)
	if err != nil {
		fmt.Println("alignbinary.Read error:", err)
		return
	}
	fmt.Println("writeMsg=", writeMsg)
}
```

## Benchmark

Some benchmarks comparing with the package binary in the  standard library.

```
goos: windows
goarch: amd64
BenchmarkBinaryWrite-4            200000              5520 ns/op          62.32 MB/s        1632 B/op         62 allocs/op
BenchmarkWrite-4                 1000000              1333 ns/op         257.98 MB/s        1536 B/op          6 allocs/op
BenchmarkBinaryRead-4            1000000              1437 ns/op         239.30 MB/s         560 B/op         27 allocs/op
BenchmarkRead-4                 10000000               161 ns/op        2125.13 MB/s         368 B/op          2 allocs/op
```

## Examples

The simplest way to use the alignbinary is to use the default entry with the system compiler default alignment:

```go
package main

import (
	"github.com/happyxcj/alignbinary"
	"encoding/binary"
	"fmt"
)

func main() {
	var msg uint16= 3
	data, _ := alignbinary.Encode(binary.BigEndian, msg)
	var newMsg uint16
	alignbinary.Decode(data,binary.BigEndian,&newMsg)
	fmt.Println("newMsg=",newMsg)
}
```

For more advanced usage such as choosing a different aligned factor:

```go
package main

import (
	"github.com/happyxcj/alignbinary"
	"encoding/binary"
	"fmt"
)

type User struct {
	Id   [16]byte
	Name [3]uint32
	Age  byte
}

func main() {
	eg := alignbinary.NewEncoderGroup(alignbinary.Align4Byte)
	alignbinary.ReplaceEncoderGroup(eg)
	dg := alignbinary.NewDecoderGroup(alignbinary.Align4Byte)
	alignbinary.ReplaceDecoderGroup(dg)
	msg := &User{
		[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		[3]uint32{100, 200, 300},
		20,
	}
	data, _ := alignbinary.Encode(binary.BigEndian, msg)
	newMsg:=&User{}
	alignbinary.Decode(data, binary.BigEndian, newMsg)
	fmt.Println("newMsg=", newMsg)
}

```
