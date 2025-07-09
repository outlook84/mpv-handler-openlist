# mpv-handler-openlist

[中文](./README_zh.md)

A URL protocol handler (`mpv://`) for the [mpv](https://mpv.io/) or [mpv.net](https://github.com/mpvnet-player/mpv.net) media players on Windows. This tool is designed to be used with the [OpenList](https://github.com/OpenListTeam/OpenList) web application to open video links in mpv or mpv.net player.

## Features

- **`mpv://` Protocol**: Handles `mpv://` URLs to open videos in mpv or mpv.net.
- **Easy Setup**: Simple command-line installation and uninstallation.
- **Configurable**: The path to `mpv.exe` or `mpvnet.exe` is configurable.
- **Logging**: Optional logging for troubleshooting.
- **Custom User-Agent**: Allows setting custom User-Agents for specific URL paths.

## Installation

1.  **Download**: Go to the [Releases page](https://github.com/outlook84/mpv-handler-openlist/releases) and download the latest `mpv-handler.exe`.
2.  **Place Executable**: Move `mpv-handler.exe` to a permanent location on your computer (e.g., inside your mpv or mpv.net folder).
3.  **Register Protocol**: Open a Command Prompt or PowerShell **as an administrator** in the directory where you placed the executable and run the following command. **Remember to replace the path with the actual path to your `mpv.exe` or `mpvnet.exe`**.

    ```shell
    .\mpv-handler.exe --install "C:\path\to\your\mpv.exe"
    ```

    If successful, you will see the message "Protocol installed and mpv path saved." A configuration file named `mpv-handler.ini` will also be created in the same directory.

## Usage

Once installed, simply click the `mpv` icon on the [OpenList](https://github.com/OpenListTeam/OpenList) web video playback page, and it will automatically call the player to play the current video.

## Uninstallation

To remove the URL protocol from your system, open a Command Prompt or PowerShell **as an administrator** in the tool's directory and run:

```shell
.\mpv-handler.exe --uninstall
```

## Configuration

The tool uses a configuration file named `mpv-handler.ini`, located in the same directory as the executable.

```ini
[mpv-handler]
mpvPath   = C:\path\to\your\mpv.exe ; Path to mpv.exe or mpvnet.exe
enableLog = false                   ; Set to true to enable logging
logPath   = mpv-handler.log         ; Path for the log file
[UserAgents]
aaa/bbb = "pan.baidu.com"
bbb/ccc = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0"
```

### Custom User-Agent

You can specify a custom User-Agent for video sources under specific paths. To use this feature, add a `[UserAgents]` section to your `mpv-handler.ini` file.

## License
This project is licensed under the GNU General Public License v2.0. See the [LICENSE](./LICENSE) file for details.
