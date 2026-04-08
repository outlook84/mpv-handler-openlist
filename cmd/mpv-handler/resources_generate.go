package main

//go:generate go run github.com/akavel/rsrc@latest -manifest mpv-handler.exe.manifest -arch amd64 -o rsrc_windows_amd64.syso
//go:generate go run github.com/akavel/rsrc@latest -manifest mpv-handler.exe.manifest -arch arm64 -o rsrc_windows_arm64.syso
