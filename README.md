# HDU 校园网客户端

这是一个跨平台校园网客户端，面向 macOS、Windows 和 Linux。

## 功能特性

- 登录学校校园网
- 断线自动重连
- 开机自动连接
- 托盘常驻
- 关闭窗口时最小化到托盘
- Linux 用户级 `systemd --user` 常驻服务

## 平台说明

- macOS：提供 Wails 桌面端，发布 `.dmg`
- Windows：提供 Wails 桌面端，发布 `.msi`
- Linux：只提供 CLI 二进制和一键安装脚本

## Linux 一键安装

Linux 只需要复制一条命令，执行完成后就会自动安装 `hdu-cli`。
下面这条命令使用了国内常见的 CDN 加速，复制即可：

```bash
wget -qO- https://cdn.jsdelivr.net/gh/hduhelp/hdu-cli@main/scripts/install-linux.sh | bash
```

如果加速地址不可用，也可以切换到 GitHub 原始地址：

```bash
wget -qO- https://raw.githubusercontent.com/hduhelp/hdu-cli/main/scripts/install-linux.sh | bash
```

安装完成后，按顺序执行下面 3 条命令即可：

```bash
hdu-cli config init --username 你的学号 --password 你的密码
hdu-cli service enable
hdu-cli net status
```

常用命令：

```bash
hdu-cli config show
hdu-cli net login
hdu-cli net status
hdu-cli service status
hdu-cli service disable
```

卸载命令：

```bash
wget -qO- https://cdn.jsdelivr.net/gh/hduhelp/hdu-cli@main/scripts/uninstall-linux.sh | bash
```

## CLI 用法

直接登录：

```bash
hdu-cli net login --username 20230001 --password 你的密码
```

前台守护并自动重连：

```bash
hdu-cli net daemon --username 20230001 --password 你的密码 --interval 60
```

初始化配置：

```bash
hdu-cli config init --username 20230001 --password 你的密码
```

## 桌面端

桌面端代码位于 `desktop/wails`，基于 Wails v3 构建。

- macOS：托盘常驻、设置窗口、开机自启
- Windows：托盘常驻、设置窗口、开机自启

关闭窗口不会退出程序，而是自动隐藏到托盘。

## 开发说明

当前仓库要求：

- Go `1.26.4`
- Node `22+`

查看 Go 版本：

```bash
go version
```

运行测试：

```bash
go test ./...
```

构建桌面端前端：

```bash
cd desktop/wails/frontend
npm install
npm run build
```

## 自动发布

GitHub Actions 会自动完成以下工作：

- 运行 Go 测试
- 打包 Linux CLI
- 构建 macOS `.dmg`
- 构建 Windows `.msi`
- 上传 GitHub Release 附件


