# Coremail 密码验证工具

一个简单的 Go 小工具，用来批量验证 Coremail 邮箱的账号密码。支持并发，可以挂 SOCKS5 代理跑。

## 功能

- **高并发**: 用 Go 的协程跑，速度很快。
- **字典方便**: `users.txt` 放用户名，`passwords.txt` 放密码，程序会自动组合。
- **SOCKS5 代理**:
    - 可以通过代理隐藏真实 IP。
    - 支持带用户名和密码的代理。
    - 不设置代理参数就直连。
- **判断准确**: 检查最终跳转的 URL 和页面内容，比只看 Cookie 靠谱多了，不容易误报。
- **跨平台**: Go 写的，在哪都能编译运行。

## 编译

首先，你得有 Go 环境 (版本 >= 1.18)。

### 编译当前平台版本

在项目根目录，直接运行：
```bash
go build .
```
之后目录下会多一个 `coremail-brute` (Windows 是 `.exe`) 的可执行文件。

### 交叉编译

想在 Mac 或 Linux 上编译出其他平台的版本也很简单。

- **编译 Windows (64-bit):**
  ```bash
  GOOS=windows GOARCH=amd64 go build -o coremail-brute.exe .
  ```

- **编译 Linux (64-bit):**
  ```bash
  GOOS=linux GOARCH=amd64 go build -o coremail-brute-linux .
  ```

- **编译 macOS (Intel):**
  ```bash
  GOOS=darwin GOARCH=amd64 go build -o coremail-brute-macos-intel .
  ```

- **编译 macOS (Apple Silicon M1/M2...):**
  ```bash
  GOOS=darwin GOARCH=arm64 go build -o coremail-brute-macos-arm64 .
  ```

## 用法

```
./coremail-brute -url <目标URL> [其他参数]
```

程序会默认读取同目录下的 `users.txt` 和 `passwords.txt`。

### 参数

- `-url` (**必需**): 目标 Coremail 地址。比如: `https://mail.example.com`
- `-c` (可选): 并发数，默认 `50`。
- `-proxy` (可选): SOCKS5 代理地址。比如: `127.0.0.1:10808`
- `-proxyuser` (可选): 代理的用户名。
- `-proxypass` (可选): 代理的密码。

### 示例

**1. 直连测试**

把 `users.txt` 和 `passwords.txt` 放在程序旁边。
```bash
./coremail-brute -url https://mail.my-company.com -c 100
```

**2. 使用普通 SOCKS5 代理**
```bash
./coremail-brute -url https://mail.my-company.com -proxy 127.0.0.1:10808
```

**3. 使用带账号密码的 SOCKS5 代理**
```bash
./coremail-brute -url https://mail.my-company.com -proxy proxy.example.com:12345 -proxyuser your_user -proxypass your_pass
```

## 免责声明

这工具只用于合法的安全测试。用之前请确保你已经拿到了授权。乱搞出了事别来找我。
