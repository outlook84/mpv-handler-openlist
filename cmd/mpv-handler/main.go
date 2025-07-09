// mpv-handler.go
package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
	"gopkg.in/ini.v1"
)

// Config holds external settings loaded/saved via config.ini
type Config struct {
	MpvPath      string
	EnableLog    bool
	LogPath      string
	UserAgentMap map[string]string
}

// loadConfig reads the configuration file from the executable's directory.
func loadConfig() (*Config, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not determine executable path: %w", err)
	}
	dir := filepath.Dir(exe)
	iniPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))+".ini")

	defaultLogPath := filepath.Join(dir, "mpv-handler.log")
	defaultConfig := &Config{
		MpvPath:      "",
		EnableLog:    false,
		LogPath:      defaultLogPath,
		UserAgentMap: make(map[string]string),
	}

	loadOpts := ini.LoadOptions{
		Insensitive:         true,
		IgnoreInlineComment: true,
	}

	cfgFile, err := ini.LoadSources(loadOpts, iniPath)
	if err != nil {
		return defaultConfig, nil
	}

	secMpvHandler := cfgFile.Section("mpv-handler")
	defaultConfig.MpvPath = secMpvHandler.Key("mpvPath").MustString("")
	defaultConfig.EnableLog = secMpvHandler.Key("enableLog").MustBool(false)
	defaultConfig.LogPath = secMpvHandler.Key("logPath").MustString(defaultLogPath)

	secUserAgents := cfgFile.Section("UserAgents")
	if secUserAgents != nil {
		defaultConfig.UserAgentMap = secUserAgents.KeysHash()
	}

	return defaultConfig, nil
}

// saveConfig writes the configuration file.
// It now loads the existing file first to preserve formatting.
func saveConfig(cfg *Config) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}
	dir := filepath.Dir(exe)
	iniPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))+".ini")

	// 1. Load the existing file using the SAME options as loadConfig.
	loadOpts := ini.LoadOptions{
		Insensitive:         true,
		IgnoreInlineComment: true,
	}
	file, err := ini.LoadSources(loadOpts, iniPath)
	if err != nil {
		// If the file doesn't exist, create a new empty object.
		file = ini.Empty()
	}

	// 2. Overwrite or create sections and keys based on the current cfg state.
	secMpvHandler, _ := file.GetSection("mpv-handler")
	if secMpvHandler == nil {
		secMpvHandler, _ = file.NewSection("mpv-handler")
	}
	secMpvHandler.Key("mpvPath").SetValue(cfg.MpvPath)
	secMpvHandler.Key("enableLog").SetValue(fmt.Sprintf("%v", cfg.EnableLog))
	secMpvHandler.Key("logPath").SetValue(cfg.LogPath)

	file.DeleteSection("UserAgents")
	secUserAgents, _ := file.NewSection("UserAgents")
	for pattern, ua := range cfg.UserAgentMap {
		secUserAgents.Key(pattern).SetValue(ua)
	}

	// 3. Save the modified file object back to disk.
	return file.SaveTo(iniPath)
}


// writeLog appends a message if enabled
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
	f.WriteString(line)
}

// installSelf registers protocol
func installSelf(exePath string) error {
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, "mpv", registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	key.SetStringValue("", "URL:Mpv-OpenList Protocol")
	key.SetStringValue("URL Protocol", "")
	iconKey, _, err := registry.CreateKey(key, `DefaultIcon`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer iconKey.Close()
	iconKey.SetStringValue("", exePath+",0")
	cmdKey, _, err := registry.CreateKey(key, `shell\open\command`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer cmdKey.Close()
	cmdKey.SetStringValue("", fmt.Sprintf("\"%s\" \"%%1\"", exePath))
	return nil
}

// uninstallSelf removes protocol
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

// handleURL processes the URL and launches mpv with appropriate arguments
func handleURL(raw string, cfg *Config) error {
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Raw URL: %s", raw))
	const prefix = "mpv://"
	if !strings.HasPrefix(raw, prefix) {
		writeLog(cfg.EnableLog, cfg.LogPath, "Invalid scheme: "+raw)
		return fmt.Errorf("invalid scheme")
	}
	stripped := raw[len(prefix):]
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Stripped URL: %s", stripped))
	decoded, err := url.QueryUnescape(stripped)
	if err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Decode error: %v", err))
		return err
	}
	writeLog(cfg.EnableLog, cfg.LogPath, fmt.Sprintf("Decoded URL: %s", decoded))

	if _, err := os.Stat(cfg.MpvPath); err != nil {
		writeLog(cfg.EnableLog, cfg.LogPath, "mpv not found at: "+cfg.MpvPath)
		return err
	}

	args := []string{}
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
	return exec.Command(cfg.MpvPath, args...).Start()
}

func main() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting executable path:", err)
		os.Exit(1)
	}
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading configuration:", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--install":
			if len(os.Args) != 3 {
				fmt.Fprintln(os.Stderr, "Usage: mpv-handler --install \"<path-to-mpv.exe>\"")
				os.Exit(1)
			}
			mpvPath := os.Args[2]
			if _, err := os.Stat(mpvPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error: mpv.exe not found at the specified path: %s\n", mpvPath)
				os.Exit(1)
			}
			cfg.MpvPath = mpvPath
			if err := saveConfig(cfg); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to save config:", err)
				os.Exit(1)
			}
			if err := installSelf(exe); err != nil {
				fmt.Fprintln(os.Stderr, "Install failed:", err)
				os.Exit(1)
			}
			fmt.Println("Protocol installed and mpv path saved.")
			return
		case "--uninstall":
			if err := uninstallSelf(); err != nil {
				fmt.Fprintln(os.Stderr, "Uninstall failed:", err)
				os.Exit(1)
			}
			fmt.Println("Protocol uninstalled.")
			return
		default:
			if err := handleURL(arg, cfg); err != nil {
				os.Exit(2)
			}
			return
		}
	}

	fmt.Println("mpv-handler: A protocol handler for mpv.")
	fmt.Println("Usage:")
	fmt.Println("  mpv-handler --install \"<full-path-to-mpv.exe>\"   : Register the mpv:// protocol.")
	fmt.Println("  mpv-handler --uninstall                          : Unregister the mpv:// protocol.")
	fmt.Println("\nThis program is usually not called by users directly, but by a web browser via the protocol.")
	fmt.Println("\nConfiguration is stored in mpv-handler.ini in the same directory as the executable.")
}