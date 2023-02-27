package main

import (
	"goflutter/flutter"
)

type TestChannel struct {
}

func (t TestChannel) Name() string {
	return "testchannel"
}

func (t TestChannel) Handle(data interface{}) error {
	println(data)
	return nil
}

var (
	testchannel = TestChannel{}
)

func main() {
	flutter.FlutterStartup()
	defer flutter.FlutterCleanup()
	dartProj := flutter.NewFlutterDartProject("data")
	win := flutter.NewFlutterWindow(dartProj)
	win.RegisterPlugin(testchannel)
	win.Create()
	win.SetQuitOnClose(true)
	flutter.FlutterMessageLoop()
}
