package alignbinary

import (
	"io"
	"encoding/binary"
)

var defaultEG = NewEncoderGroup(AlignDefault)
var defaultDG = NewDecoderGroup(AlignDefault)

func ReplaceEncoderGroup(eg *EncoderGroup) {
	defaultEG = eg
}

func ReplaceDecoderGroup(dg *DecoderGroup) {
	defaultDG = dg
}

func Write(w io.Writer, order binary.ByteOrder, msg interface{}) error {
	return defaultEG.Write(w, order, msg)
}

func Encode(order binary.ByteOrder, msg interface{}) ([]byte, error) {
	return defaultEG.Encode(order, msg)
}

func Read(r io.Reader, order binary.ByteOrder, msg interface{}) error {
	return defaultDG.Read(r, order, msg)
}

func Decode(data []byte, order binary.ByteOrder, msg interface{}) error {
	return defaultDG.Decode(data, order, msg)
}
