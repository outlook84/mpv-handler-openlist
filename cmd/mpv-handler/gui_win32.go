package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	wsOverlappedWindow = 0x00CF0000
	wsVisible          = 0x10000000
	wsChild            = 0x40000000
	wsClipChildren     = 0x02000000
	wsTabStop          = 0x00010000
	wsVScroll          = 0x00200000

	wsExClientEdge = 0x00000200

	esAutoHScroll = 0x0080
	esMultiline   = 0x0004
	esReadOnly    = 0x0800

	bsPushButton    = 0x00000000
	bsAutoCheckbox  = 0x00000003
	bmGetCheck      = 0x00F0
	bmSetCheck      = 0x00F1
	bstChecked      = 0x0001
	wmCreate        = 0x0001
	wmDestroy       = 0x0002
	wmCommand       = 0x0111
	wmCtlColorStatic = 0x0138
	wmSetFont       = 0x0030
	swShow          = 5
	ofnPathMustExist = 0x00000800
	ofnFileMustExist = 0x00001000
	ofnExplorer     = 0x00080000
	idcArrow        = 32512
	colorWindow     = 5
	defaultGUIFont  = 17
	transparentBkMode = 1

	idEditMpvPath      = 1001
	idButtonBrowse     = 1002
	idEditExtraArgs    = 1003
	idCheckEnableLog   = 1004
	idButtonSave       = 1005
	idButtonRegister   = 1006
	idButtonUnregister = 1007
	idStatusLabel      = 1008
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	comdlg32                 = syscall.NewLazyDLL("comdlg32.dll")
	gdi32                    = syscall.NewLazyDLL("gdi32.dll")
	procMessageBoxW          = user32.NewProc("MessageBoxW")
	procRegisterClassExW     = user32.NewProc("RegisterClassExW")
	procCreateWindowExW      = user32.NewProc("CreateWindowExW")
	procDefWindowProcW       = user32.NewProc("DefWindowProcW")
	procGetSysColorBrush     = user32.NewProc("GetSysColorBrush")
	procShowWindow           = user32.NewProc("ShowWindow")
	procUpdateWindow         = user32.NewProc("UpdateWindow")
	procGetMessageW          = user32.NewProc("GetMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessageW     = user32.NewProc("DispatchMessageW")
	procPostQuitMessage      = user32.NewProc("PostQuitMessage")
	procLoadCursorW          = user32.NewProc("LoadCursorW")
	procSendMessageW         = user32.NewProc("SendMessageW")
	procSetWindowTextW       = user32.NewProc("SetWindowTextW")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procGetModuleHandleW     = kernel32.NewProc("GetModuleHandleW")
	procGetOpenFileNameW     = comdlg32.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
	procGetStockObject       = gdi32.NewProc("GetStockObject")
	procCreateFontW          = gdi32.NewProc("CreateFontW")
	procDeleteObject         = gdi32.NewProc("DeleteObject")
	procSetBkMode            = gdi32.NewProc("SetBkMode")

	mainWindowClassName = syscall.StringToUTF16Ptr("MpvHandlerMainWindow")
	currentAppState     *AppState
)

type AppState struct {
	hwnd           uintptr
	editMpvPath    uintptr
	editExtraArgs  uintptr
	checkEnableLog uintptr
	statusLabel    uintptr
	titleFont      uintptr
	labelFont      uintptr
	bodyFont       uintptr
	smallFont      uintptr
	cfg            *Config
}

type point struct {
	X int32
	Y int32
}

type msg struct {
	HWnd     uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       point
	LPrivate uint32
}

type wndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type openFilename struct {
	LStructSize       uint32
	HwndOwner         uintptr
	HInstance         uintptr
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        unsafe.Pointer
	DwReserved        uint32
	FlagsEx           uint32
}

func mustUTF16Ptr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

func createWindow(exStyle uint32, className, windowName string, style uint32, x, y, width, height int32, parent, menu, instance uintptr) uintptr {
	classPtr := mustUTF16Ptr(className)
	windowPtr := mustUTF16Ptr(windowName)
	hwnd, _, _ := procCreateWindowExW.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(classPtr)),
		uintptr(unsafe.Pointer(windowPtr)),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		parent,
		menu,
		instance,
		0,
	)
	return hwnd
}

func setWindowText(hwnd uintptr, text string) {
	textPtr := mustUTF16Ptr(text)
	procSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(textPtr)))
}

func getWindowText(hwnd uintptr) string {
	length, _, _ := procGetWindowTextLengthW.Call(hwnd)
	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func setCheckbox(hwnd uintptr, checked bool) {
	value := uintptr(0)
	if checked {
		value = bstChecked
	}
	procSendMessageW.Call(hwnd, bmSetCheck, value, 0)
}

func isCheckboxChecked(hwnd uintptr) bool {
	ret, _, _ := procSendMessageW.Call(hwnd, bmGetCheck, 0, 0)
	return ret == bstChecked
}

func applyDefaultFont(font uintptr, hwnds ...uintptr) {
	for _, hwnd := range hwnds {
		if hwnd != 0 {
			procSendMessageW.Call(hwnd, wmSetFont, font, 1)
		}
	}
}

func openFileDialog(owner uintptr) (string, error) {
	buffer := make([]uint16, 1024)
	filter := []uint16{
		'E', 'x', 'e', 'c', 'u', 't', 'a', 'b', 'l', 'e', 's', 0,
		'*', '.', 'e', 'x', 'e', 0,
		'A', 'l', 'l', ' ', 'F', 'i', 'l', 'e', 's', 0,
		'*', '.', '*', 0,
		0,
	}

	title := mustUTF16Ptr(ui.ChooseMpvExecutableTitle)
	ofn := openFilename{
		LStructSize: uint32(unsafe.Sizeof(openFilename{})),
		HwndOwner:   owner,
		LpstrFilter: &filter[0],
		LpstrFile:   &buffer[0],
		NMaxFile:    uint32(len(buffer)),
		LpstrTitle:  title,
		Flags:       ofnExplorer | ofnPathMustExist | ofnFileMustExist,
	}

	ret, _, err := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if ret == 0 {
		extendedErr, _, _ := procCommDlgExtendedError.Call()
		if extendedErr != 0 {
			return "", fmt.Errorf("GetOpenFileNameW failed with CommDlgExtendedError=0x%X", uint32(extendedErr))
		}
		if err != syscall.Errno(0) {
			return "", err
		}
		return "", nil
	}
	return syscall.UTF16ToString(buffer), nil
}

func getModuleHandle() uintptr {
	handle, _, _ := procGetModuleHandleW.Call(0)
	return handle
}

func getDefaultFont() uintptr {
	font, _, _ := procGetStockObject.Call(defaultGUIFont)
	return font
}

func getWindowBrush() uintptr {
	brush, _, _ := procGetSysColorBrush.Call(colorWindow)
	return brush
}

func createFont(height int32, weight int32, face string) uintptr {
	facePtr := mustUTF16Ptr(face)
	font, _, _ := procCreateFontW.Call(
		uintptr(int32(-height)),
		0,
		0,
		0,
		uintptr(weight),
		0,
		0,
		0,
		1,
		0,
		0,
		5,
		0,
		uintptr(unsafe.Pointer(facePtr)),
	)
	return font
}

func loword(v uintptr) uint16 {
	return uint16(v & 0xFFFF)
}

func currentConfigFromUI() *Config {
	cfg := defaultConfig()
	if currentAppState != nil && currentAppState.cfg != nil {
		*cfg = *currentAppState.cfg
		if currentAppState.cfg.UserAgentMap != nil {
			cfg.UserAgentMap = currentAppState.cfg.UserAgentMap
		}
	}

	cfg.MpvPath = strings.TrimSpace(getWindowText(currentAppState.editMpvPath))
	cfg.ExtraArgs = strings.TrimSpace(getWindowText(currentAppState.editExtraArgs))
	cfg.EnableLog = isCheckboxChecked(currentAppState.checkEnableLog)
	if strings.TrimSpace(cfg.LogPath) == "" {
		cfg.LogPath = defaultLogPath()
	}
	if cfg.UserAgentMap == nil {
		cfg.UserAgentMap = make(map[string]string)
	}
	return cfg
}

func refreshStatus() {
	if currentAppState == nil || currentAppState.statusLabel == 0 {
		return
	}

	cfg := currentConfigFromUI()
	status := ui.ProtocolNotRegistered
	if isProtocolInstalled() {
		status = ui.ProtocolRegistered
	}
	if cfg.MpvPath != "" {
		status += "\r\n" + ui.ConfiguredPlayerPrefix + cfg.MpvPath
	}
	if cfg.ExtraArgs != "" {
		status += "\r\n" + ui.ConfiguredExtraArgsPrefix + cfg.ExtraArgs
	} else {
		status += "\r\n" + ui.ConfiguredExtraArgsPrefix + ui.ConfiguredExtraArgsNone
	}
	setWindowText(currentAppState.statusLabel, status)
}

func handleSaveConfig() {
	cfg := currentConfigFromUI()
	if strings.TrimSpace(cfg.ExtraArgs) != "" {
		if _, err := parseExtraArgs(cfg.ExtraArgs); err != nil {
			showMessage(ui.MpvArgsInvalidTitle, ui.MpvArgsInvalidMessage, true)
			return
		}
	}
	if err := saveConfig(cfg); err != nil {
		showMessage(ui.SaveFailed, fmt.Sprintf(ui.SaveFailedMessage, err), true)
		return
	}
	currentAppState.cfg = cfg
	refreshStatus()
	showMessage(ui.ConfigSavedTitle, ui.ConfigSavedMessage, false)
}

func validateMpvPath(cfg *Config) error {
	if strings.TrimSpace(cfg.MpvPath) == "" {
		return newUserVisibleError(
			ui.MpvPathRequiredTitle,
			ui.MpvPathRequiredMessage,
			nil,
		)
	}
	if err := validateMpvExecutablePath(cfg.MpvPath); err != nil {
		return newUserVisibleError(
			ui.SelectedPathInvalidTitle,
			fmt.Sprintf(ui.SelectedPathInvalidMessage, cfg.MpvPath),
			err,
		)
	}
	if _, err := parseExtraArgs(cfg.ExtraArgs); err != nil {
		return newUserVisibleError(
			ui.MpvArgsInvalidTitle,
			ui.MpvArgsInvalidMessage,
			err,
		)
	}
	return nil
}

func handleRegister() {
	cfg := currentConfigFromUI()
	if err := validateMpvPath(cfg); err != nil {
		if visibleErr, ok := err.(*userVisibleError); ok {
			showMessage(visibleErr.Title, visibleErr.Message, true)
			return
		}
		showMessage(ui.RegisterFailedTitle, err.Error(), true)
		return
	}

	if err := saveConfig(cfg); err != nil {
		showMessage(ui.RegisterFailedTitle, fmt.Sprintf(ui.RegisterFailedSaveMessage, err), true)
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		showMessage(ui.RegisterFailedTitle, fmt.Sprintf(ui.RegisterFailedExeMessage, err), true)
		return
	}
	if err := installSelf(exePath); err != nil {
		showMessage(ui.RegisterFailedTitle, fmt.Sprintf(ui.RegisterFailedProtocolMessage, err), true)
		return
	}

	currentAppState.cfg = cfg
	refreshStatus()
	showMessage(ui.ProtocolRegisteredTitle, fmt.Sprintf(ui.ProtocolRegisteredMessage, cfg.MpvPath), false)
}

func handleUnregister() {
	if err := uninstallSelf(); err != nil {
		showMessage(ui.UnregisterFailedTitle, fmt.Sprintf(ui.UnregisterFailedMessage, err), true)
		return
	}
	refreshStatus()
	showMessage(ui.ProtocolRemovedTitle, ui.ProtocolRemovedMessage, false)
}

func createMainControls(hwnd uintptr) {
	instance := getModuleHandle()
	currentAppState = &AppState{
		hwnd:      hwnd,
		titleFont: createFont(20, 600, "Segoe UI"),
		labelFont: createFont(13, 600, "Segoe UI"),
		bodyFont:  createFont(13, 400, "Segoe UI"),
		smallFont: createFont(12, 400, "Segoe UI"),
	}
	if currentAppState.titleFont == 0 {
		currentAppState.titleFont = getDefaultFont()
	}
	if currentAppState.labelFont == 0 {
		currentAppState.labelFont = getDefaultFont()
	}
	if currentAppState.bodyFont == 0 {
		currentAppState.bodyFont = getDefaultFont()
	}
	if currentAppState.smallFont == 0 {
		currentAppState.smallFont = getDefaultFont()
	}

	titleLabel := createWindow(0, "STATIC", ui.MainTitle, wsChild|wsVisible, 24, 18, 300, 30, hwnd, 0, instance)
	subtitleLabel := createWindow(0, "STATIC", ui.MainSubtitle, wsChild|wsVisible, 24, 46, 620, 20, hwnd, 0, instance)

	mpvLabel := createWindow(0, "STATIC", ui.MpvExecutableLabel, wsChild|wsVisible, 24, 86, 150, 20, hwnd, 0, instance)
	currentAppState.editMpvPath = createWindow(wsExClientEdge, "EDIT", "", wsChild|wsVisible|wsTabStop|esAutoHScroll, 24, 110, 470, 28, hwnd, uintptr(idEditMpvPath), instance)
	browseButton := createWindow(0, "BUTTON", ui.BrowseButton, wsChild|wsVisible|wsTabStop|bsPushButton, 506, 109, 110, 30, hwnd, uintptr(idButtonBrowse), instance)

	argsLabel := createWindow(0, "STATIC", ui.ExtraArgsLabel, wsChild|wsVisible, 24, 154, 170, 20, hwnd, 0, instance)
	argsHintLabel := createWindow(0, "STATIC", ui.ExtraArgsHint, wsChild|wsVisible, 24, 176, 260, 18, hwnd, 0, instance)
	currentAppState.editExtraArgs = createWindow(wsExClientEdge, "EDIT", "", wsChild|wsVisible|wsTabStop|esAutoHScroll, 24, 198, 592, 28, hwnd, uintptr(idEditExtraArgs), instance)

	currentAppState.checkEnableLog = createWindow(0, "BUTTON", ui.EnableLog, wsChild|wsVisible|wsTabStop|bsAutoCheckbox, 24, 244, 140, 24, hwnd, uintptr(idCheckEnableLog), instance)
	saveButton := createWindow(0, "BUTTON", ui.SaveConfig, wsChild|wsVisible|wsTabStop|bsPushButton, 24, 286, 124, 34, hwnd, uintptr(idButtonSave), instance)
	registerButton := createWindow(0, "BUTTON", ui.RegisterProtocol, wsChild|wsVisible|wsTabStop|bsPushButton, 160, 286, 146, 34, hwnd, uintptr(idButtonRegister), instance)
	unregisterButton := createWindow(0, "BUTTON", ui.ClearRegistration, wsChild|wsVisible|wsTabStop|bsPushButton, 318, 286, 146, 34, hwnd, uintptr(idButtonUnregister), instance)

	statusHeader := createWindow(0, "STATIC", ui.CurrentStatus, wsChild|wsVisible, 24, 340, 120, 20, hwnd, 0, instance)
	currentAppState.statusLabel = createWindow(wsExClientEdge, "EDIT", "", wsChild|wsVisible|esMultiline|esReadOnly|wsVScroll, 24, 364, 592, 88, hwnd, uintptr(idStatusLabel), instance)

	applyDefaultFont(currentAppState.titleFont, titleLabel)
	applyDefaultFont(currentAppState.labelFont, mpvLabel, argsLabel, statusHeader)
	applyDefaultFont(currentAppState.bodyFont, currentAppState.editMpvPath, browseButton, currentAppState.editExtraArgs, currentAppState.checkEnableLog, saveButton, registerButton, unregisterButton, currentAppState.statusLabel)
	applyDefaultFont(currentAppState.smallFont, subtitleLabel, argsHintLabel)

	cfg, err := loadConfig()
	if err != nil {
		cfg = defaultConfig()
		showMessage(ui.LoadFailedTitle, fmt.Sprintf(ui.LoadFailedMessage, err), true)
	}
	currentAppState.cfg = cfg
	setWindowText(currentAppState.editMpvPath, cfg.MpvPath)
	setWindowText(currentAppState.editExtraArgs, cfg.ExtraArgs)
	setCheckbox(currentAppState.checkEnableLog, cfg.EnableLog)
	refreshStatus()
}

func wndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmCreate:
		createMainControls(hwnd)
		return 0
	case wmCtlColorStatic:
		procSetBkMode.Call(wParam, transparentBkMode)
		return getWindowBrush()
	case wmCommand:
		switch loword(wParam) {
		case idButtonBrowse:
			path, err := openFileDialog(hwnd)
			if err != nil {
				showMessage(ui.BrowseFailedTitle, fmt.Sprintf(ui.BrowseFailedMessage, err), true)
				return 0
			}
			if path != "" {
				setWindowText(currentAppState.editMpvPath, path)
				refreshStatus()
			}
			return 0
		case idButtonSave:
			handleSaveConfig()
			return 0
		case idButtonRegister:
			handleRegister()
			return 0
		case idButtonUnregister:
			handleUnregister()
			return 0
		case idCheckEnableLog:
			refreshStatus()
			return 0
		}
	case wmDestroy:
		if currentAppState != nil {
			for _, font := range []uintptr{
				currentAppState.titleFont,
				currentAppState.labelFont,
				currentAppState.bodyFont,
				currentAppState.smallFont,
			} {
				if font != 0 && font != getDefaultFont() {
					procDeleteObject.Call(font)
				}
			}
		}
		procPostQuitMessage.Call(0)
		return 0
	}
	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func runGUI() error {
	instance := getModuleHandle()
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)
	class := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     instance,
		HCursor:       cursor,
		HbrBackground: colorWindow + 1,
		LpszClassName: mainWindowClassName,
	}

	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&class)))
	if atom == 0 && err != syscall.Errno(1410) {
		return fmt.Errorf(ui.WindowClassRegisterFailed, err)
	}

	hwnd := createWindow(
		0,
		"MpvHandlerMainWindow",
		ui.AppTitle,
		wsOverlappedWindow|wsVisible|wsClipChildren,
		200,
		200,
		660,
		520,
		0,
		0,
		instance,
	)
	if hwnd == 0 {
		return fmt.Errorf("%s", ui.WindowCreateFailed)
	}

	procShowWindow.Call(hwnd, swShow)
	procUpdateWindow.Call(hwnd)

	var message msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&message)), 0, 0, 0)
		switch int32(ret) {
		case -1:
			return fmt.Errorf("%s", ui.MessageLoopFailed)
		case 0:
			return nil
		default:
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&message)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&message)))
		}
	}
}
