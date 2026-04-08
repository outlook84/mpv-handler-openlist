# mpv-handler-openlist

[中文](./README_zh.md)

A URL protocol handler (`mpv://`) for the [mpv](https://mpv.io/) or [mpv.net](https://github.com/mpvnet-player/mpv.net) media players on Windows. This tool is designed to be used with the [OpenList](https://github.com/OpenListTeam/OpenList) web application to open video links in mpv or mpv.net player.

## Features

- **`mpv://` Protocol**: Handles `mpv://` URLs to open videos in mpv or mpv.net.
- **Easy Setup**: Register or remove the protocol from the GUI, or use the command-line flags.
- **Configurable**: Configure the path to `mpv.exe` or `mpvnet.exe`.
- **Extra mpv Arguments**: Pass optional startup flags such as `--fs`.
- **Logging**: Optional logging for troubleshooting.
- **Custom User-Agent**: Allows setting custom User-Agents for specific URL paths.

## Installation

1.  **Download**: Go to the [Releases page](https://github.com/outlook84/mpv-handler-openlist/releases) and download the latest archive for your system, such as `mpv-handler_v1.2.3_windows_amd64.zip`.
2.  **Extract Files**: Unzip the archive to a permanent location on your computer (for example, inside your mpv or mpv.net folder).
3.  **Windows Warning Note**: Release binaries are currently unsigned. Windows may show a SmartScreen or "file came from another computer" warning for files downloaded from GitHub. If needed, open the file properties and click `Unblock`, or run `Unblock-File .\mpv-handler.exe` in PowerShell before launching it.
4.  **Open the GUI**: Double-click `mpv-handler.exe`.
5.  **Choose your player**: Pick `mpv.exe` or `mpvnet.exe` with the `Browse...` button.
6.  **Optional**: Add extra mpv arguments, for example `--fs`.
7.  **Register Protocol**: Click `Register Protocol`. On systems where registry writes need elevation, run the app as administrator.

## Command-Line Usage

- Install:

    ```shell
    .\mpv-handler.exe --install "C:\path\to\your\mpv.exe"
    ```

- Uninstall:

    ```shell
    .\mpv-handler.exe --uninstall
    ```

## Configuration

The tool uses a configuration file named `mpv-handler.ini`, located in the same directory as the executable.

```ini
[mpv-handler]
mpvPath   = C:\path\to\your\mpv.exe
extraArgs = --fs
enableLog = false
logPath   = mpv-handler.log
[UserAgents]
aaa/bbb = "pan.baidu.com"
bbb/ccc = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0"
```

- `mpvPath`: Path to mpv.exe or mpvnet.exe
- `extraArgs`: Optional extra command-line arguments passed to mpv before the URL
- `enableLog`: Set to true to enable logging
- `logPath`: Path for the log file

### Custom User-Agent

You can specify a custom User-Agent for video sources under specific paths.

The key is a path prefix that will be matched against the part of the URL after `/d/`. For example, for the URL `https://.../d/aaa/bbb/ccc`, the keys used for matching could be `aaa` or `aaa/bbb/`.

## License
This project is licensed under the GNU General Public License v2.0. See the [LICENSE](./LICENSE) file for details.
