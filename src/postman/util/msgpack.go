package util

import (
	"github.com/ugorji/go/codec"
)

var mh codec.MsgpackHandle

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

// encode struct to msgpack format
func MsgEncode(v interface{}) string {
	msg := []byte{}
	enc := codec.NewEncoderBytes(&msg, &mh)
	enc.Encode(v)
	return string(msg)
}

// decode struct from msgpack format
func MsgDecode(msg string, v interface{}) error {
	dec := codec.NewDecoderBytes([]byte(msg), &mh)
	return dec.Decode(&v)
}
