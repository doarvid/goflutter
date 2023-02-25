package main

import "goflutter/flutter"

func main() {
	flutter.FlutterStartup()
	defer flutter.FlutterCleanup()
	dartProj := flutter.NewFlutterDartProject("data")
	win := flutter.NewFlutterWindow(dartProj)
	win.Create()
	win.SetQuitOnClose(true)
	flutter.FlutterMessageLoop()
}
