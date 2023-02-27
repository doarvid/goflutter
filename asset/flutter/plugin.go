package flutter

import "C"
import (
	"unsafe"

	"github.com/go-flutter-desktop/go-flutter/plugin"
)

type FlutterPlugin interface {
	Name() string
	Handle(data interface{}) error
}

//export pluginMethodCall
func pluginMethodCall(message unsafe.Pointer, msglen C.int, pluginptr unsafe.Pointer) {
	msgdata := C.GoBytes(message, msglen)
	codec := plugin.StandardMessageCodec{}
	msg, err := codec.DecodeMessage(msgdata)
	if err != nil {
		return
	}
	po := *(*FlutterPlugin)(pluginptr)
	po.Handle(msg)

}
