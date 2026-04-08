# mpv-handler-openlist (中文)

[English](./README.md)

这是一个为 Windows 平台上的 [mpv](https://mpv.io/) 或 [mpv.net](https://github.com/mpvnet-player/mpv.net) 媒体播放器设计的 URL 协议注册器 (`mpv://`)。该工具用于在 [OpenList](https://github.com/OpenListTeam/OpenList) Web 网页上调用 mpv 或 mpv.net 播放器来打开视频链接。

## 功能特性

- **`mpv://` 协议**: 处理 `mpv://` 格式的 URL，用 mpv 或 mpv.net 打开视频。
- **简易安装**: 可以直接在 GUI 里注册/清除协议，也保留了命令行方式。
- **可配置**: `mpv.exe` 或 `mpvnet.exe` 的路径是可配置的。
- **额外 mpv 参数**: 可以追加诸如 `--fs` 这样的启动参数。
- **日志记录**: 提供可选的日志功能，方便排查问题。
- **自定义 User-Agent**: 允许为特定的 URL 路径设置自定义的 User-Agent。

## 安装步骤

1.  **下载**: 前往 [Releases 页面](https://github.com/outlook84/mpv-handler-openlist/releases) 下载适合你系统的最新压缩包，例如 `mpv-handler_v1.2.3_windows_amd64.zip`。
2.  **解压文件**: 将压缩包解压到你电脑上的一个固定位置（例如，mpv 或 mpv.net 播放器的文件夹内）。
3.  **Windows 警告说明**: 当前发布的二进制文件未做代码签名。Windows 对从 GitHub 下载的文件可能会弹出 SmartScreen 或“此文件来自其他计算机”的警告。如有需要，可先在文件属性里点击“解除锁定”，或在 PowerShell 中运行 `Unblock-File .\mpv-handler.exe` 后再启动。
4.  **打开 GUI**: 双击 `mpv-handler.exe`。
5.  **选择播放器**: 用 `浏览...` 按钮选择 `mpv.exe` 或 `mpvnet.exe`。
6.  **可选**: 填写额外 mpv 参数，例如 `--fs`。
7.  **注册协议**: 点击 `注册协议`。如果系统要求管理员权限，请以管理员身份运行本程序。

## 命令行用法

- 安装：

    ```shell
    .\mpv-handler.exe --install "C:\你的路径\mpv.exe"
    ```

- 卸载：

    ```shell
    .\mpv-handler.exe --uninstall
    ```

## 配置文件

本工具使用一个名为 `mpv-handler.ini` 的配置文件，它位于和主程序相同的目录下。

```ini
[mpv-handler]
mpvPath   = C:\你的路径\mpv.exe
extraArgs = --fs
enableLog = false
logPath   = mpv-handler.log
[UserAgents]
aaa/bbb = "pan.baidu.com"
bbb/ccc = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0"
```

- `mpvPath`: mpv.exe 或 mpvnet.exe 的路径
- `extraArgs`: 在视频 URL 之前附加给 mpv 的额外命令行参数
- `enableLog`: 设置为 true 来启用日志记录
- `logPath`: 日志文件的路径

### 自定义 User-Agent

您可以为特定路径下的视频源指定自定义的 User-Agent。

配置的键是一个路径前缀，它将与 URL 中 `/d/` 之后的部分进行匹配。例如，对于 URL `https://.../d/aaa/bbb/ccc`，用于匹配的路径是 `aaa` 或 `aaa/bbb/`。

## 许可证

本项目基于 GNU General Public License v2.0 许可证。详情请参阅 [LICENSE](./LICENSE) 文件。
