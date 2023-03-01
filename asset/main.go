package main

import (
	"{{projectname}}/flutter"

	"github.com/go-flutter-desktop/go-flutter/plugin"
)

type TestChannel struct {
}

func (t TestChannel) Name() string {
	return "testchannel"
}

func (t TestChannel) Handle(data *plugin.MethodCall) (interface{}, error) {
	println(data)
	return nil, nil
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
