package flutter

//#cgo CFLAGS:-I./flutter
//#cgo CFLAGS:-I./ -std=c17
//#cgo LDFLAGS: -L./ -L../build/windows/runner/Release {{runnerlib}}
//#include <runner.h>
//#include <stdio.h>
import "C"
import (
	"unsafe"
)

type FlutterDartProject struct {
	proj unsafe.Pointer
}
type FlutterWind struct {
	proj       *FlutterDartProject
	flutterWnd unsafe.Pointer
}

func NewFlutterDartProject(project_data string) *FlutterDartProject {
	proj := &FlutterDartProject{}
	cstr2 := C.CString(project_data)
	proj.proj = C.NewFlutterDartProject(cstr2)
	return proj
}

func NewFlutterWindow(dartProj *FlutterDartProject) *FlutterWind {
	win := &FlutterWind{proj: dartProj}
	win.flutterWnd = C.NewFlutterWindow(dartProj.proj)
	return win
}

func (f *FlutterWind) Create() {
	C.FlutterWindowCreate(f.flutterWnd, 10, 10, 1280, 720)
	C.FlutterWindowSetQuitOnClose(f.flutterWnd, 1)
}

func (f *FlutterWind) SetQuitOnClose(quit_on_close bool) {
	flag_ := C.int(0)
	if quit_on_close {
		flag_ = C.int(1)
	}
	C.FlutterWindowSetQuitOnClose(f.flutterWnd, flag_)
}

func FlutterStartup() {
	C.FlutterStartup()
}

func FlutterMessageLoop() {
	C.FlutterRun()
}

func FlutterCleanup() {
	C.FlutterCleanup()
}
