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
	MpvPath   string
	EnableLog bool
	LogPath   string
}

// loadConfig reads the configuration file from the executable's directory.
// The config file is named after the executable with a .ini extension.
func loadConfig() (*Config, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not determine executable path: %w", err)
	}
	dir := filepath.Dir(exe)
	iniPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))+".ini")
	cfgFile, err := ini.LoadSources(ini.LoadOptions{Insensitive: true}, iniPath)
	if err != nil {
		return &Config{MpvPath: "", EnableLog: false, LogPath: filepath.Join(dir, "mpv-handler.log")}, nil
	}
	sec := cfgFile.Section("mpv-handler")
	return &Config{
		MpvPath:   sec.Key("mpvPath").MustString(""),
		EnableLog: sec.Key("enableLog").MustBool(false),
		LogPath:   sec.Key("logPath").MustString(filepath.Join(dir, "mpv-handler.log")),
	}, nil
}

// saveConfig writes the configuration file.
func saveConfig(cfg *Config) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}
	dir := filepath.Dir(exe)
	iniPath := filepath.Join(dir, strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))+".ini")
	file := ini.Empty()
	sec := file.Section("mpv-handler")
	sec.Key("mpvPath").SetValue(cfg.MpvPath)
	sec.Key("enableLog").SetValue(fmt.Sprintf("%v", cfg.EnableLog))
	sec.Key("logPath").SetValue(cfg.LogPath)
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
	// Keys must be deleted from deepest to shallowest.
	// We ignore "not found" errors (syscall.ENOENT), as the key may have
	// already been deleted, but we return any other error.
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
	return nil // Success
}

// handleURL processes the URL
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
	return exec.Command(cfg.MpvPath, decoded).Start()
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

	// Command-line handling
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
				// The error is already logged by handleURL if logging is enabled.
				// Avoid printing to stderr to prevent OS error dialogs when called by a browser.
				os.Exit(2)
			}
			return
		}
	}

	// If no arguments are provided, show usage information.
	fmt.Println("mpv-handler: A protocol handler for mpv.")
	fmt.Println("Usage:")
	fmt.Println("  mpv-handler --install \"<full-path-to-mpv.exe>\"   : Register the mpv:// protocol.")
	fmt.Println("  mpv-handler --uninstall                          : Unregister the mpv:// protocol.")
	fmt.Println("\nThis program is usually not called by users directly, but by a web browser via the protocol.")
}
