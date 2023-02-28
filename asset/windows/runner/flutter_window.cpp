#include "flutter_window.h"

#include <optional>

#include "flutter/generated_plugin_registrant.h"



#include <map>

#include <flutter/method_channel.h>
#include <flutter/plugin_registrar_windows.h>
#include <flutter/standard_method_codec.h>

namespace {
    class StandardPlugin : public flutter::Plugin {
    public:
        static void RegisterWithRegistrar(flutter::PluginRegistrarWindows* registrar, FlutterPlugin* pluginCfg);

        StandardPlugin();

        virtual ~StandardPlugin();

        FlutterPlugin* pluginCfg;

    private:
        void HandleMethodCall(const flutter::MethodCall<flutter::EncodableValue>& method_call, std::unique_ptr<flutter::MethodResult<flutter::EncodableValue>> result);
    };

    void StandardPlugin::RegisterWithRegistrar(flutter::PluginRegistrarWindows* registrar, FlutterPlugin* pluginCfg) {
        auto channel = std::make_unique<flutter::MethodChannel<flutter::EncodableValue>>(registrar->messenger(), pluginCfg->channel, &flutter::StandardMethodCodec::GetInstance());

        auto plugin = std::make_unique<StandardPlugin>();
        plugin->pluginCfg = pluginCfg;

        channel->SetMethodCallHandler(
            [plugin_pointer = plugin.get()](const auto& call, auto result) {
            plugin_pointer->HandleMethodCall(call, std::move(result));
        }
        );

        registrar->AddPlugin(std::move(plugin));
    }

    StandardPlugin::StandardPlugin() {}

    StandardPlugin::~StandardPlugin() {}

    void StandardPlugin::HandleMethodCall(const flutter::MethodCall<flutter::EncodableValue>& method_call, std::unique_ptr<flutter::MethodResult<flutter::EncodableValue>> result) {
       static const flutter::StandardMethodCodec& decoder = flutter::StandardMethodCodec::GetInstance();
       std::unique_ptr<std::vector<uint8_t>> data = decoder.EncodeMethodCall(method_call);
       
       OutputDebugStringA("handle method call");

       if (this->pluginCfg!=NULL){
        OutputDebugStringA("cfg != NULL");
       }
       if (this->pluginCfg != NULL && this->pluginCfg->callback!= NULL){
          OutputDebugStringA("handle method call##1");
          this->pluginCfg->callback(data->data(), int(data->size()), this->pluginCfg->plugin,(void*)result.get());
       }
    }
}

void StandardPluginRegisterWithRegistrar(flutter::PluginRegistry* registry, FlutterPlugin* pluginCfg) {
    FlutterDesktopPluginRegistrarRef registrar = registry->GetRegistrarForPlugin(pluginCfg->channel);
    StandardPlugin::RegisterWithRegistrar(
        flutter::PluginRegistrarManager::GetInstance()->GetRegistrar<flutter::PluginRegistrarWindows>(registrar),
        pluginCfg
    );
}





void FlutterPluginMethodReply(int code,unsigned char* message,int size,void*resultptr){
    static const flutter::StandardMethodCodec& decoder = flutter::StandardMethodCodec::GetInstance();
    flutter::MethodResult<flutter::EncodableValue>*result=(flutter::MethodResult<flutter::EncodableValue>*)resultptr;
  if (code == 0) {
      if (message == NULL) {
          result->Success();
      }
      else {
          decoder.DecodeAndProcessResponseEnvelope(message, size, result);
      }
  }
  else {
      char myerrno[32] = { 0 };
      snprintf(myerrno, 32, "%d", code);
      result->Error(myerrno,(char*) message);
  }
}



FlutterWindow::FlutterWindow(const flutter::DartProject& project)
    : project_(project) {}

FlutterWindow::~FlutterWindow() {}

void FlutterWindow::RegisterPlugin(char* channel,FlutterPluginMethodCallback callback,void*plugin)
{
    OutputDebugStringA("======================,1channel:");
    FlutterPlugin* fplugin = new FlutterPlugin;
    fplugin->channel = channel;
    fplugin->callback = callback;
    fplugin->plugin = plugin;
    this->standardplugins.push_back(fplugin);
}

bool FlutterWindow::OnCreate() {
  if (!Win32Window::OnCreate()) {
    return false;
  }

  RECT frame = GetClientArea();

  // The size here must match the window dimensions to avoid unnecessary surface
  // creation / destruction in the startup path.
  flutter_controller_ = std::make_unique<flutter::FlutterViewController>(
      frame.right - frame.left, frame.bottom - frame.top, project_);
  // Ensure that basic setup of the controller was successful.
  if (!flutter_controller_->engine() || !flutter_controller_->view()) {
    return false;
  }
  RegisterPlugins(flutter_controller_->engine());
  printf("================\n");
  OutputDebugStringA("====\n");
  for(int i =0;i < standardplugins.size();++i){
    FlutterPlugin* plugin =standardplugins.at(i);
    OutputDebugStringA(plugin->channel);
    StandardPluginRegisterWithRegistrar(flutter_controller_->engine(), plugin);
  }

  SetChildContent(flutter_controller_->view()->GetNativeWindow());

  flutter_controller_->engine()->SetNextFrameCallback([&]() {
    this->Show();
  });

  return true;
}

void FlutterWindow::OnDestroy() {
  if (flutter_controller_) {
    flutter_controller_ = nullptr;
  }

  Win32Window::OnDestroy();
}

LRESULT
FlutterWindow::MessageHandler(HWND hwnd, UINT const message,
                              WPARAM const wparam,
                              LPARAM const lparam) noexcept {
  // Give Flutter, including plugins, an opportunity to handle window messages.
  if (flutter_controller_) {
    std::optional<LRESULT> result =
        flutter_controller_->HandleTopLevelWindowProc(hwnd, message, wparam,
                                                      lparam);
    if (result) {
      return *result;
    }
  }

  switch (message) {
    case WM_FONTCHANGE:
      flutter_controller_->engine()->ReloadSystemFonts();
      break;
  }

  return Win32Window::MessageHandler(hwnd, message, wparam, lparam);
}
