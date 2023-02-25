#include <flutter/dart_project.h>
#include <flutter/flutter_view_controller.h>
#include <windows.h>

#include "flutter_window.h"
#include "utils.h"

#include "runner.h"

void *NewFlutterDartProject(char *path)
{
    int bufSize = MultiByteToWideChar(CP_ACP, 0, path, -1, NULL, 0);
    wchar_t* wp = new wchar_t[bufSize]; //leak mem
    MultiByteToWideChar(CP_ACP, 0, path, -1, wp, bufSize);

    return new flutter::DartProject(wp);
}

void FlutterDartProjectSetEntrypointArgs(void *project, int argc, char **argv)
{
    flutter::DartProject *p = (flutter::DartProject *)project;

    std::vector<std::string> command_line_arguments;
    for (int i = 0; i < argc; ++i)
    {
        command_line_arguments.push_back(argv[i]);
    }
    p->set_dart_entrypoint_arguments(std::move(command_line_arguments));
}

void *NewFlutterWindow(void *dartProject)
{
    flutter::DartProject *p = (flutter::DartProject *)dartProject;
    return new FlutterWindow(*p);
}

int FlutterWindowCreate(void *flutterwnd, int pos_x, int pos_y, int size_height, int size_width)
{
    FlutterWindow *window = (FlutterWindow *)flutterwnd;
    Win32Window::Point origin(pos_x, pos_y);
    Win32Window::Size size(size_width, size_height);
    return window->Create(L"example", origin, size) != true;
}

void FlutterWindowSetQuitOnClose(void *flutterwnd, int quit_and_close)
{
    FlutterWindow *window = (FlutterWindow *)flutterwnd;
    window->SetQuitOnClose(quit_and_close >0);
}

void FlutterRun() 
{
    ::MSG msg;
    while (::GetMessage(&msg, nullptr, 0, 0)) {
        ::TranslateMessage(&msg);
        ::DispatchMessage(&msg);
    }
}

RUNNER_EXPORT void FlutterStartup(){
    if (!::AttachConsole(ATTACH_PARENT_PROCESS) && ::IsDebuggerPresent()) {
        CreateAndAttachConsole();
        printf("CreateAndAttachConsole\n");
    }
    HRESULT res =   ::CoInitializeEx(nullptr, COINIT_APARTMENTTHREADED);
    printf("CoInitializeEx:%d\n",res);

}
RUNNER_EXPORT void FlutterCleanup(){
  ::CoUninitialize();
}