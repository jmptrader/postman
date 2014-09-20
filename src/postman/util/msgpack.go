package util

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

var mh codec.MsgpackHandle

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

// encode struct to msgpack format
func MsgEncode(v interface{}) (msg []byte, err error) {
	msg = []byte{}
	enc := codec.NewEncoderBytes(&msg, &mh)
	err = enc.Encode(v)
	return
}

// decode struct from msgpack format
func MsgDecode(msg []byte, v interface{}) error {
	dec := codec.NewDecoderBytes(msg, &mh)
	return dec.Decode(&v)
}
