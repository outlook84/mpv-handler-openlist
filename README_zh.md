# mpv-handler-openlist (中文)

[English](./README.md)

这是一个为 Windows 平台上的 [mpv](https://mpv.io/) 或 [mpv.net](https://github.com/mpvnet-player/mpv.net) 媒体播放器设计的 URL 协议注册器 (`mpv://`)。该工具用于在 [OpenList](https://github.com/OpenListTeam/OpenList) Web 网页上调用 mpv 或 mpv.net 播放器来打开视频链接。

## 功能特性

- **`mpv://` 协议**: 处理 `mpv://` 格式的 URL，用 mpv 或 mpv.net 打开视频。
- **简易安装**: 通过简单的命令行指令进行安装和卸载。
- **可配置**: `mpv.exe` 或 `mpvnet.exe` 的路径是可配置的。
- **日志记录**: 提供可选的日志功能，方便排查问题。

## 安装步骤

1.  **下载**: 前往 [Releases 页面](https://github.com/outlook84/mpv-handler-openlist/releases) 下载最新的 `mpv-handler.exe` 文件。
2.  **放置程序**: 将 `mpv-handler.exe` 移动到你电脑上的一个固定位置（例如，mpv 或 mpv.net 播放器的文件夹内）。
3.  **注册协议**: 在 `mpv-handler.exe` 所在的目录中 **以管理员身份** 打开命令提示符或 PowerShell，然后运行以下命令。**请务必将路径替换为你自己电脑上 `mpv.exe` 或 `mpvnet.exe` 的实际路径**。

    ```shell
    .\mpv-handler.exe --install "C:\你的路径\mpv.exe"
    ```

    如果成功，你将看到提示信息 "Protocol installed and mpv path saved."。同时，一个名为 `mpv-handler.ini` 的配置文件也会在同目录下被创建。

## 如何使用

安装完成后，只需在 [OpenList](https://github.com/OpenListTeam/OpenList) Web 视频播放页面上点击 `mpv` 图标，就会自动调用播放器播放当前视频。

## 卸载

要从你的系统中移除此 URL 协议，请在工具所在的目录 **以管理员身份** 打开命令提示符或 PowerShell，然后运行：

```shell
.\mpv-handler.exe --uninstall
```

## 配置文件

本工具使用一个名为 `mpv-handler.ini` 的配置文件，它位于和主程序相同的目录下。

```ini
[mpv-handler]
mpvPath   = C:\你的路径\mpv.exe ; mpv.exe 或 mpvnet.exe 的路径
enableLog = false              ; 设置为 true 来启用日志记录
logPath   = mpv-handler.log    ; 日志文件的路径
```

## 许可证

本项目基于 GNU General Public License v2.0 许可证。详情请参阅 [LICENSE](./LICENSE) 文件。
