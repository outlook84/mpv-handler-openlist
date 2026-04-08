package main

var (
	procGetUserDefaultUILanguage = kernel32.NewProc("GetUserDefaultUILanguage")
	ui                           = pickStrings()
)

type UIStrings struct {
	AppTitle                      string
	MainTitle                     string
	MainSubtitle                  string
	MpvExecutableLabel            string
	BrowseButton                  string
	ChooseMpvExecutableTitle      string
	ExtraArgsLabel                string
	ExtraArgsHint                 string
	EnableLog                     string
	SaveConfig                    string
	RegisterProtocol              string
	ClearRegistration             string
	CurrentStatus                 string
	ProtocolRegistered            string
	ProtocolNotRegistered         string
	ConfiguredPlayerPrefix        string
	ConfiguredExtraArgsPrefix     string
	ConfiguredExtraArgsNone       string
	ErrorTitle                    string
	UnsupportedLink               string
	DecodeFailed                  string
	MpvNotConfiguredTitle         string
	MpvNotConfiguredMessage       string
	MpvPathInvalidTitle           string
	MpvPathInvalidMessage         string
	MpvArgsInvalidTitle           string
	MpvArgsInvalidMessage         string
	MpvLaunchFailedTitle          string
	MpvLaunchFailedMessage        string
	SaveFailed                    string
	SaveFailedMessage             string
	ConfigSavedTitle              string
	ConfigSavedMessage            string
	MpvPathRequiredTitle          string
	MpvPathRequiredMessage        string
	SelectedPathInvalidTitle      string
	SelectedPathInvalidMessage    string
	RegisterFailedTitle           string
	RegisterFailedSaveMessage     string
	RegisterFailedExeMessage      string
	RegisterFailedProtocolMessage string
	ProtocolRegisteredTitle       string
	ProtocolRegisteredMessage     string
	UnregisterFailedTitle         string
	UnregisterFailedMessage       string
	ProtocolRemovedTitle          string
	ProtocolRemovedMessage        string
	LoadFailedTitle               string
	LoadFailedMessage             string
	BrowseFailedTitle             string
	BrowseFailedMessage           string
	WindowClassRegisterFailed     string
	WindowCreateFailed            string
	MessageLoopFailed             string
	InstallFailedTitle            string
	InstallPathNotFound           string
	InstallSaveFailed             string
	InstallExeFailed              string
	InstallFailed                 string
	InstalledTitle                string
	InstalledMessage              string
	InstallUsage                  string
	UninstalledTitle              string
	UninstalledMessage            string
	ConfigLoadError               string
	FailedToOpenInMpv             string
}

func englishStrings() UIStrings {
	return UIStrings{
		AppTitle:                      "mpv-handler",
		MainTitle:                     "OpenList mpv handler",
		MainSubtitle:                  "Choose your player, set optional startup arguments, then register the mpv:// protocol.",
		MpvExecutableLabel:            "mpv executable",
		BrowseButton:                  "Browse...",
		ChooseMpvExecutableTitle:      "Choose mpv executable",
		ExtraArgsLabel:                "Extra mpv arguments",
		ExtraArgsHint:                 "Example: --fs (start in fullscreen)",
		EnableLog:                     "Enable log file",
		SaveConfig:                    "Save Config",
		RegisterProtocol:              "Register Protocol",
		ClearRegistration:             "Clear Registration",
		CurrentStatus:                 "Current status",
		ProtocolRegistered:            "Protocol status: Registered",
		ProtocolNotRegistered:         "Protocol status: Not registered",
		ConfiguredPlayerPrefix:        "Configured player: ",
		ConfiguredExtraArgsPrefix:     "Configured extra args: ",
		ConfiguredExtraArgsNone:       "(none)",
		ErrorTitle:                    "mpv-handler Error",
		UnsupportedLink:               "Unsupported link format.\nOnly mpv:// links can be opened by mpv-handler.",
		DecodeFailed:                  "The link could not be decoded.\nPlease try opening it again from OpenList.",
		MpvNotConfiguredTitle:         "mpv Not Configured",
		MpvNotConfiguredMessage:       "mpv path is not configured yet.\nOpen mpv-handler and choose your mpv.exe or mpvnet.exe first.",
		MpvPathInvalidTitle:           "mpv Path Invalid",
		MpvPathInvalidMessage:         "Configured mpv executable was not found:\n%s\n\nOpen mpv-handler and update the player path.",
		MpvArgsInvalidTitle:           "mpv Arguments Invalid",
		MpvArgsInvalidMessage:         "The extra mpv arguments could not be parsed.\nPlease review them in the mpv-handler window.",
		MpvLaunchFailedTitle:          "mpv Launch Failed",
		MpvLaunchFailedMessage:        "mpv was found, but it could not be started.\nExecutable:\n%s",
		SaveFailed:                    "Save Failed",
		SaveFailedMessage:             "Failed to save configuration:\n%v",
		ConfigSavedTitle:              "Configuration Saved",
		ConfigSavedMessage:            "Your mpv-handler settings were saved successfully.",
		MpvPathRequiredTitle:          "mpv Path Required",
		MpvPathRequiredMessage:        "Choose your mpv.exe or mpvnet.exe before registering the protocol.",
		SelectedPathInvalidTitle:      "mpv Path Invalid",
		SelectedPathInvalidMessage:    "The selected player executable does not exist:\n%s",
		RegisterFailedTitle:           "Register Failed",
		RegisterFailedSaveMessage:     "Failed to save configuration:\n%v",
		RegisterFailedExeMessage:      "Could not determine executable path:\n%v",
		RegisterFailedProtocolMessage: "Failed to register mpv:// protocol.\nYou may need to run this app as administrator.\n\n%v",
		ProtocolRegisteredTitle:       "Protocol Registered",
		ProtocolRegisteredMessage:     "mpv:// protocol was registered successfully.\n\nCurrent player:\n%s",
		UnregisterFailedTitle:         "Unregister Failed",
		UnregisterFailedMessage:       "Failed to remove mpv:// protocol.\nYou may need to run this app as administrator.\n\n%v",
		ProtocolRemovedTitle:          "Protocol Removed",
		ProtocolRemovedMessage:        "mpv:// protocol was removed successfully.",
		LoadFailedTitle:               "Load Failed",
		LoadFailedMessage:             "Failed to load configuration:\n%v",
		BrowseFailedTitle:             "Browse Failed",
		BrowseFailedMessage:           "Failed to open file picker:\n%v",
		WindowClassRegisterFailed:     "failed to register window class: %w",
		WindowCreateFailed:            "failed to create main window",
		MessageLoopFailed:             "message loop failed",
		InstallFailedTitle:            "mpv-handler Install Failed",
		InstallPathNotFound:           "mpv.exe not found at the specified path:\n%s",
		InstallSaveFailed:             "Failed to save config:\n%v",
		InstallExeFailed:              "Could not determine executable path:\n%v",
		InstallFailed:                 "Install failed:\n%v",
		InstalledTitle:                "mpv-handler Installed",
		InstalledMessage:              "mpv:// protocol installed successfully.\n\nCurrent mpv executable:\n%s",
		InstallUsage:                  "Usage: mpv-handler --install \"<path-to-mpv.exe>\"",
		UninstalledTitle:              "mpv-handler Uninstalled",
		UninstalledMessage:            "mpv:// protocol was removed successfully.",
		ConfigLoadError:               "Error loading configuration:\n%v",
		FailedToOpenInMpv:             "Failed to open in mpv:\n%v",
	}
}

func chineseStrings() UIStrings {
	return UIStrings{
		AppTitle:                      "mpv-handler",
		MainTitle:                     "OpenList mpv 调用器",
		MainSubtitle:                  "选择播放器，设置可选启动参数，然后注册 mpv:// 协议。",
		MpvExecutableLabel:            "mpv 可执行文件",
		BrowseButton:                  "浏览...",
		ChooseMpvExecutableTitle:      "选择 mpv 可执行文件",
		ExtraArgsLabel:                "额外 mpv 参数",
		ExtraArgsHint:                 "例如：--fs（启动后直接全屏）",
		EnableLog:                     "启用日志文件",
		SaveConfig:                    "保存配置",
		RegisterProtocol:              "注册协议",
		ClearRegistration:             "清除注册",
		CurrentStatus:                 "当前状态",
		ProtocolRegistered:            "协议状态：已注册",
		ProtocolNotRegistered:         "协议状态：未注册",
		ConfiguredPlayerPrefix:        "当前播放器：",
		ConfiguredExtraArgsPrefix:     "当前额外参数：",
		ConfiguredExtraArgsNone:       "（无）",
		ErrorTitle:                    "mpv-handler 错误",
		UnsupportedLink:               "链接格式不受支持。\nmpv-handler 只能打开 mpv:// 链接。",
		DecodeFailed:                  "链接解码失败。\n请回到 OpenList 再试一次。",
		MpvNotConfiguredTitle:         "未配置 mpv",
		MpvNotConfiguredMessage:       "还没有配置 mpv 路径。\n请先打开 mpv-handler，选择你的 mpv.exe 或 mpvnet.exe。",
		MpvPathInvalidTitle:           "mpv 路径无效",
		MpvPathInvalidMessage:         "找不到已配置的 mpv 可执行文件：\n%s\n\n请打开 mpv-handler 更新播放器路径。",
		MpvArgsInvalidTitle:           "mpv 参数无效",
		MpvArgsInvalidMessage:         "额外 mpv 参数无法解析。\n请在 mpv-handler 窗口中检查它们。",
		MpvLaunchFailedTitle:          "mpv 启动失败",
		MpvLaunchFailedMessage:        "已找到 mpv，但启动失败。\n可执行文件：\n%s",
		SaveFailed:                    "保存失败",
		SaveFailedMessage:             "保存配置失败：\n%v",
		ConfigSavedTitle:              "配置已保存",
		ConfigSavedMessage:            "mpv-handler 配置已成功保存。",
		MpvPathRequiredTitle:          "需要 mpv 路径",
		MpvPathRequiredMessage:        "请先选择 mpv.exe 或 mpvnet.exe，再注册协议。",
		SelectedPathInvalidTitle:      "mpv 路径无效",
		SelectedPathInvalidMessage:    "所选播放器可执行文件不存在：\n%s",
		RegisterFailedTitle:           "注册失败",
		RegisterFailedSaveMessage:     "保存配置失败：\n%v",
		RegisterFailedExeMessage:      "无法确定当前程序路径：\n%v",
		RegisterFailedProtocolMessage: "注册 mpv:// 协议失败。\n你可能需要以管理员身份运行此程序。\n\n%v",
		ProtocolRegisteredTitle:       "协议已注册",
		ProtocolRegisteredMessage:     "mpv:// 协议已成功注册。\n\n当前播放器：\n%s",
		UnregisterFailedTitle:         "清除注册失败",
		UnregisterFailedMessage:       "移除 mpv:// 协议失败。\n你可能需要以管理员身份运行此程序。\n\n%v",
		ProtocolRemovedTitle:          "协议已移除",
		ProtocolRemovedMessage:        "mpv:// 协议已成功移除。",
		LoadFailedTitle:               "加载失败",
		LoadFailedMessage:             "加载配置失败：\n%v",
		BrowseFailedTitle:             "浏览失败",
		BrowseFailedMessage:           "打开文件选择器失败：\n%v",
		WindowClassRegisterFailed:     "注册窗口类失败：%w",
		WindowCreateFailed:            "创建主窗口失败",
		MessageLoopFailed:             "消息循环失败",
		InstallFailedTitle:            "mpv-handler 安装失败",
		InstallPathNotFound:           "在指定路径找不到 mpv.exe：\n%s",
		InstallSaveFailed:             "保存配置失败：\n%v",
		InstallExeFailed:              "无法确定当前程序路径：\n%v",
		InstallFailed:                 "安装失败：\n%v",
		InstalledTitle:                "mpv-handler 已安装",
		InstalledMessage:              "mpv:// 协议已成功安装。\n\n当前 mpv 可执行文件：\n%s",
		InstallUsage:                  "用法：mpv-handler --install \"<path-to-mpv.exe>\"",
		UninstalledTitle:              "mpv-handler 已卸载",
		UninstalledMessage:            "mpv:// 协议已成功移除。",
		ConfigLoadError:               "加载配置失败：\n%v",
		FailedToOpenInMpv:             "调用 mpv 打开失败：\n%v",
	}
}

func pickStrings() UIStrings {
	lang, _, _ := procGetUserDefaultUILanguage.Call()
	primary := uint16(lang) & 0x03ff
	if primary == 0x04 {
		return chineseStrings()
	}
	return englishStrings()
}
