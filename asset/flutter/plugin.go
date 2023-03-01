package flutter

//#cgo CFLAGS:-I./flutter
//#cgo CFLAGS:-I./ -std=c17
//#cgo LDFLAGS: -L./ -L../build/windows/runner/Release -l{{projectname}}
/*
#include <stdio.h>
#include <runner.h>
void pluginMethodCall(void* message,int size,void*plugin,void*resultptr);
*/
import "C"
import (
	"unsafe"

	"github.com/go-flutter-desktop/go-flutter/plugin"
)

type FlutterPlugin interface {
	Name() string
	Handle(data *plugin.MethodCall) (interface{}, error)
}

type ReplyCallInf interface {
	Handle(C.int, unsafe.Pointer, C.int, unsafe.Pointer)
}

func replyMethodCall(code int, result []byte, replyptr unsafe.Pointer) {
	if result != nil {
		creslut := C.CBytes(result)
		C.FlutterPluginMethodReply(C.int(code), (*C.uchar)(creslut), C.int(len(result)), replyptr)
	} else {
		C.FlutterPluginMethodReply(C.int(code), nil, C.int(0), replyptr)
	}
}

//export pluginMethodCall
func pluginMethodCall(message unsafe.Pointer, msglen C.int, pluginptr unsafe.Pointer, replyptr unsafe.Pointer) {
	msgdata := C.GoBytes(message, msglen)
	codec := plugin.StandardMethodCodec{}
	msg, err := codec.DecodeMethodCall(msgdata)
	if err != nil {
		replyMethodCall(-1, []byte("decode methodcall error"), replyptr)
		return
	}
	po := *(*FlutterPlugin)(pluginptr)
	_, err = po.Handle(&msg)
	if err != nil {
		replyMethodCall(-2, []byte(err.Error()), replyptr)
	} else {
		replyMethodCall(0, []byte("decode methodcall error"), replyptr)
	}
	return
}
