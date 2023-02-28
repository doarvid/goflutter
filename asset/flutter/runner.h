// Copyright 2013 The Flutter Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef RUNNER_HH
#define RUNNER_HH


#ifdef FLUTTER_DESKTOP_LIBRARY

// Add visibility/export annotations when building the library.
#ifdef _WIN32
#define RUNNER_EXPORT __declspec(dllexport)
#else
#define RUNNER_EXPORT __attribute__((visibility("default")))
#endif

#else  // FLUTTER_DESKTOP_LIBRARY

// Add import annotations when consuming the library.
#ifdef _WIN32
#define RUNNER_EXPORT __declspec(dllimport)
#else
#define RUNNER_EXPORT
#endif

#endif  // FLUTTER_DESKTOP_LIBRARY


#if defined(__cplusplus)
extern "C" {
#endif  // defined(__cplusplus)

typedef void (*FlutterPluginMethodCallback)(unsigned char* message,int size,void*plugin,void*resultptr);


RUNNER_EXPORT void *NewFlutterDartProject(char *path);
RUNNER_EXPORT void FlutterDartProjectSetEntrypointArgs(void *project, int argc, char **argv);
RUNNER_EXPORT void *NewFlutterWindow(void *dartProject);
RUNNER_EXPORT void FlutterWindowRegisterPlugin(void *flutterwnd, char* channel,FlutterPluginMethodCallback callback,void*plugin);
RUNNER_EXPORT void FlutterPluginMethodReply(int code,unsigned char* message,int size,void*resultptr);
RUNNER_EXPORT int FlutterWindowCreate(void *flutterwnd, int pos_x, int pos_y, int size_height, int size_width);
RUNNER_EXPORT void FlutterWindowSetQuitOnClose(void *flutterwnd, int quit_and_close);
RUNNER_EXPORT void FlutterRun();
RUNNER_EXPORT void FlutterStartup();
RUNNER_EXPORT void FlutterCleanup();
#if defined(__cplusplus)
}
#endif  // defined(__cplusplus)

#endif  // FLUTTER_SHELL_PLATFORM_COMMON_PUBLIC_FLUTTER_EXPORT_H_
