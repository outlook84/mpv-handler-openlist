package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
	"gopkg.in/ini.v1"
)

const (
	mbOK        = 0x00000000
	mbIconError = 0x00000010
	mbIconInfo  = 0x00000040
)

var version = "dev"

type Config struct {
	MpvPath      string
	ExtraArgs    string
	EnableLog    bool
	LogPath      string
	UserAgentMap map[string]string
}

type userVisibleError struct {
	Title   string
	Message string
	Err     error
}

func (e *userVisibleError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func newUserVisibleError(title, message string, err error) error {
	return &userVisibleError{Title: title, Message: message, Err: err}
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not determine executable path: %w", err)
	}
	return filepath.Dir(exe), nil
}

func defaultLogPath() string {
	dir, err := executableDir()
	if err != nil {
		return "mpv-handler.log"
	}
	return filepath.Join(dir, "mpv-handler.log")
}

func configPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not determine executable path: %w", err)
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))+".ini"), nil
}

func defaultConfig() *Config {
	return &Config{
		MpvPath:      "",
		ExtraArgs:    "",
		EnableLog:    false,
		LogPath:      defaultLogPath(),
		UserAgentMap: make(map[string]string),
	}
}

func loadConfig() (*Config, error) {
	cfg := defaultConfig()
	iniPath, err := configPath()
	if err != nil {
		return nil, err
	}

	loadOpts := ini.LoadOptions{
		Insensitive:         true,
		IgnoreInlineComment: true,
	}
	cfgFile, err := ini.LoadSources(loadOpts, iniPath)
	if err != nil {
		return cfg, nil
	}

	sec := cfgFile.Section("mpv-handler")
	cfg.MpvPath = sec.Key("mpvPath").MustString("")
	cfg.ExtraArgs = sec.Key("extraArgs").MustString("")
	cfg.EnableLog = sec.Key("enableLog").MustBool(false)
	cfg.LogPath = sec.Key("logPath").MustString(cfg.LogPath)

	secUserAgents := cfgFile.Section("UserAgents")
	if secUserAgents != nil {
		cfg.UserAgentMap = secUserAgents.KeysHash()
	}
	return cfg, nil
}

func saveConfig(cfg *Config) error {
	iniPath, err := configPath()
	if err != nil {
		return err
	}

	loadOpts := ini.LoadOptions{
		Insensitive:         true,
		IgnoreInlineComment: true,
	}
	file, err := ini.LoadSources(loadOpts, iniPath)
	if err != nil {
		file = ini.Empty()
	}

	sec, _ := file.GetSection("mpv-handler")
	if sec == nil {
		sec, _ = file.NewSection("mpv-handler")
	}
	sec.Key("mpvPath").SetValue(cfg.MpvPath)
	sec.Key("extraArgs").SetValue(cfg.ExtraArgs)
	sec.Key("enableLog").SetValue(fmt.Sprintf("%v", cfg.EnableLog))
	sec.Key("logPath").SetValue(cfg.LogPath)

	file.DeleteSection("UserAgents")
	secUserAgents, _ := file.NewSection("UserAgents")
	for pattern, ua := range cfg.UserAgentMap {
		secUserAgents.Key(pattern).SetValue(ua)
	}

	return file.SaveTo(iniPath)
}

func writeLog(enable bool, logPath, msg string) {
	if !enable {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	line := fmt.Sprintf("%s | %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	_, _ = f.WriteString(line)
}

func showMessage(title, text string, isError bool) {
	flags := uintptr(mbOK | mbIconInfo)
	if isError {
		flags = uintptr(mbOK | mbIconError)
	}

	owner := uintptr(0)
	if currentAppState != nil && currentAppState.hwnd != 0 {
		owner = currentAppState.hwnd
	}

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	textPtr, _ := syscall.UTF16PtrFromString(text)
	procMessageBoxW.Call(
		owner,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		flags,
	)
}

func installSelf(exePath string) error {
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, "mpv", registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetStringValue("", "URL:Mpv-OpenList Protocol"); err != nil {
		return err
	}
	if err := key.SetStringValue("URL Protocol", ""); err != nil {
		return err
	}

	iconKey, _, err := registry.CreateKey(key, `DefaultIcon`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer iconKey.Close()
	if err := iconKey.SetStringValue("", exePath+",0"); err != nil {
		return err
	}

	cmdKey, _, err := registry.CreateKey(key, `shell\open\command`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer cmdKey.Close()
	return cmdKey.SetStringValue("", fmt.Sprintf("\"%s\" \"%%1\"", exePath))
}

func uninstallSelf() error {
	keysToDelete := []string{
		`mpv\shell\open\command`,
		`mpv\shell\open`,
		`mpv\shell`,
		`mpv\DefaultIcon`,
		`mpv`,
	}
	for _, keyPath := range keysToDelete {
		err := registry.DeleteKey(registry.CLASSES_ROOT, keyPath)
		if err != nil && err != syscall.ENOENT {
			return fmt.Errorf("failed to delete registry key %q: %w", keyPath, err)
		}
	}
	return nil
}

func isProtocolInstalled() bool {
	key, err := registry.OpenKey(registry.CLASSES_ROOT, `mpv\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue("")
	return err == nil
}

func parseExtraArgs(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	args := make([]string, 0, 4)
	for i := 0; i < len(raw); {
		for i < len(raw) && (raw[i] == ' ' || raw[i] == '\t') {
			i++
		}
		if i >= len(raw) {
			break
		}

		var current strings.Builder
		inQuotes := false

		for i < len(raw) {
			backslashes := 0
			for i < len(raw) && raw[i] == '\\' {
				backslashes++
				i++
			}

			if i < len(raw) && raw[i] == '"' {
				current.WriteString(strings.Repeat("\\", backslashes/2))
				if backslashes%2 == 0 {
					if inQuotes && i+1 < len(raw) && raw[i+1] == '"' {
						current.WriteByte('"')
						i += 2
						continue
					}
					inQuotes = !inQuotes
					i++
					continue
				}

				current.WriteByte('"')
				i++
				continue
			}

			if backslashes > 0 {
				current.WriteString(strings.Repeat("\\", backslashes))
			}

			if i >= len(raw) {
				break
			}
			if !inQuotes && (raw[i] == ' ' || raw[i] == '\t') {
				break
			}

			current.WriteByte(raw[i])
			i++
		}

		if inQuotes {
			return nil, fmt.Errorf("unterminated quoted string in extra arguments")
		}
		args = append(args, current.String())

		for i < len(raw) && (raw[i] == ' ' || raw[i] == '\t') {
			i++
		}
	}

	return args, nil
}

func validateMpvExecutablePath(path string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fmt.Errorf("mpv path is empty")
	}

	info, err := os.Stat(trimmed)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("mpv path points to a directory")
	}
	if !strings.EqualFold(filepath.Ext(trimmed), ".exe") {
		return fmt.Errorf("mpv path must point to an .exe file")
	}
	return nil
}

func handleURL(raw string, cfg *Config) error {
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Raw URL: %s", raw))
	const prefix = "mpv://"
	if !strings.HasPrefix(raw, prefix) {
		writeLog(cfg.EnableLog, cfg.LogPath, "Invalid scheme: "+raw)
		return newUserVisibleError(
			ui.ErrorTitle,
			ui.UnsupportedLink,
			fmt.Errorf("invalid scheme"),
		)
	}

	stripped := raw[len(prefix):]
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Stripped URL: %s", stripped))
	decoded, err := url.QueryUnescape(stripped)
	if err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Decode error: %v", err))
		return newUserVisibleError(
			ui.ErrorTitle,
			ui.DecodeFailed,
			err,
		)
	}

	if strings.TrimSpace(cfg.MpvPath) == "" {
		writeLog(cfg.EnableLog, cfg.LogPath, "mpv path is not configured")
		return newUserVisibleError(
			ui.MpvNotConfiguredTitle,
			ui.MpvNotConfiguredMessage,
			nil,
		)
	}

	if err := validateMpvExecutablePath(cfg.MpvPath); err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, "mpv not found at: "+cfg.MpvPath)
		return newUserVisibleError(
			ui.MpvPathInvalidTitle,
			fmt.Sprintf(ui.MpvPathInvalidMessage, cfg.MpvPath),
			err,
		)
	}

	args, err := parseExtraArgs(cfg.ExtraArgs)
	if err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Failed to parse extra args: %v", err))
		return newUserVisibleError(
			ui.MpvArgsInvalidTitle,
			ui.MpvArgsInvalidMessage,
			err,
		)
	}

	userAgent := ""
	parsedURL, err := url.Parse(decoded)
	if err == nil && parsedURL.Path != "" {
		for pathPattern, ua := range cfg.UserAgentMap {
			if strings.Contains(parsedURL.Path, pathPattern) {
				userAgent = ua
				writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Found matching UA for pattern '%s'. Using UA: %s", pathPattern, userAgent))
				break
			}
		}
	} else {
		writeLog(cfg.EnableLog, cfg.LogPath, "Could not parse URL or URL has no path, skipping UA matching.")
	}

	if userAgent != "" {
		args = append(args, "--user-agent="+userAgent)
	}
	args = append(args, decoded)
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Executing: %s %s", cfg.MpvPath, strings.Join(args, " ")))

	if err := exec.Command(cfg.MpvPath, args...).Start(); err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Failed to start mpv: %v", err))
		return newUserVisibleError(
			ui.MpvLaunchFailedTitle,
			fmt.Sprintf(ui.MpvLaunchFailedMessage, cfg.MpvPath),
			err,
		)
	}
	return nil
}

func handleInstallCLI(argPath string) {
	if err := validateMpvExecutablePath(argPath); err != nil {
		msg := fmt.Sprintf(ui.InstallPathNotFound, argPath)
		showMessage(ui.InstallFailedTitle, msg, true)
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		cfg = defaultConfig()
	}
	cfg.MpvPath = argPath
	if err := saveConfig(cfg); err != nil {
		showMessage(ui.InstallFailedTitle, fmt.Sprintf(ui.InstallSaveFailed, err), true)
		os.Exit(1)
	}

	exePath, err := os.Executable()
	if err != nil {
		showMessage(ui.InstallFailedTitle, fmt.Sprintf(ui.InstallExeFailed, err), true)
		os.Exit(1)
	}
	if err := installSelf(exePath); err != nil {
		showMessage(ui.InstallFailedTitle, fmt.Sprintf(ui.InstallFailed, err), true)
		os.Exit(1)
	}

	showMessage(ui.InstalledTitle, fmt.Sprintf(ui.InstalledMessage, argPath), false)
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--install":
			if len(os.Args) != 3 {
				showMessage(ui.InstallFailedTitle, ui.InstallUsage, true)
				os.Exit(1)
			}
			handleInstallCLI(os.Args[2])
			return
		case "--uninstall":
			if err := uninstallSelf(); err != nil {
				showMessage(ui.UnregisterFailedTitle, fmt.Sprintf(ui.UnregisterFailedMessage, err), true)
				os.Exit(1)
			}
			showMessage(ui.UninstalledTitle, ui.UninstalledMessage, false)
			return
		default:
			cfg, err := loadConfig()
			if err != nil {
				showMessage(ui.ErrorTitle, fmt.Sprintf(ui.ConfigLoadError, err), true)
				os.Exit(1)
			}
			if err := handleURL(arg, cfg); err != nil {
				if visibleErr, ok := err.(*userVisibleError); ok {
					showMessage(visibleErr.Title, visibleErr.Message, true)
				} else {
					showMessage(ui.ErrorTitle, fmt.Sprintf(ui.FailedToOpenInMpv, err), true)
				}
				os.Exit(2)
			}
			return
		}
	}

	if err := runGUI(); err != nil {
		showMessage(ui.ErrorTitle, err.Error(), true)
		os.Exit(1)
	}
}
