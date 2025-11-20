# Coremail 账户密码验证工具

这是一个使用 Go 语言编写的，用于快速验证 Coremail 邮箱账户密码有效性的命令行工具。它支持高并发，并可以通过 SOCKS5 代理（包括需要用户名密码验证的代理）进行测试。

## 主要特性

- **高并发验证**: 利用 Go 协程实现高并发，快速测试大量凭据。
- **灵活的字典支持**: 支持将用户名和密码分别存储在不同文件中，工具会自动进行组合测试。
- **SOCKS5 代理支持**:
    - 可通过 SOCKS5 代理（如快代理的动态隧道）发起请求，隐藏真实IP。
    - 支持需要用户名和密码认证的 SOCKS5 代理。
    - 代理为可选参数，不设置则直接连接。
- **增强的成功判断**: 综合检查 Session Cookie 和响应页面内容，有效降低误报率。
- **跨平台**: 基于 Go 语言，可轻松编译运行在 Windows, macOS, Linux 等多个平台。

## 安装与编译

1.  确保您已经安装了 Go 语言环境 (版本 >= 1.18)。
2.  将项目克隆或下载到本地。
3.  在项目根目录 (`coremail-brut`) 下打开终端。

### 本地编译

运行以下命令为当前系统编译：
```bash
go build .
```
编译成功后，在当前目录下会生成一个名为 `coremail-brute` (Windows下为 `coremail-brute.exe`) 的可执行文件。

### 交叉编译 (Cross-Compilation)

Go 语言的强大之处在于其内置的交叉编译能力，让您可以轻松地在单一开发环境（如 macOS 或 Linux）上，为其他多种操作系统和 CPU 架构生成可执行文件。

**命令格式说明：**

您需要在项目根目录（`coremail-brut` 目录）下打开终端，然后输入并执行以下命令。命令的基本结构是：

```
GOOS=<目标操作系统> GOARCH=<目标架构> go build -o <输出文件名> .
```

-   `GOOS=<目标操作系统>`: 这里的 `GOOS` 是一个环境变量，用于指定目标操作系统。常见的值有 `linux`, `windows`, `darwin` (macOS)。
-   `GOARCH=<目标架构>`: 这里的 `GOARCH` 是一个环境变量，用于指定目标 CPU 架构。最常用的值是 `amd64` (适用于绝大多数64位电脑)。
-   `go build`: 这是标准的 Go 语言编译命令。
-   `-o <输出文件名>`: 这个标志（flag）用来指定编译后生成的可执行文件的名字。如果不指定，Go 会使用默认名称。
-   `.`: 表示编译当前目录下的所有 Go 源文件。

**具体示例：**

以下是在您的 macOS 或 Linux 电脑上，为不同平台生成可执行文件的具体命令。直接复制粘贴到您的终端里执行即可。

-   **编译为 Windows (64-bit):**
    这条命令会生成一个名为 `coremail-brute.exe` 的文件，可以直接在 64 位的 Windows 系统上运行。
    ```bash
    GOOS=windows GOARCH=amd64 go build -o coremail-brute.exe .
    ```

-   **编译为 Linux (64-bit):**
    这条命令会生成一个名为 `coremail-brute-linux` 的文件，可以在 64 位的 Linux 系统上运行。
    ```bash
    GOOS=linux GOARCH=amd64 go build -o coremail-brute-linux .
    ```

-   **编译为 macOS (Intel 64-bit):**
    这条命令会生成一个名为 `coremail-brute-macos-intel` 的文件，可以在使用 Intel 芯片的 64 位 macOS 系统上运行。
    ```bash
    GOOS=darwin GOARCH=amd64 go build -o coremail-brute-macos-intel .
    ```

-   **编译为 macOS (Apple Silicon):**
    这条命令会生成一个名为 `coremail-brute-macos-arm64` 的文件，可以在使用 Apple M系列芯片 (M1, M2, etc.) 的 macOS 系统上运行。
    ```bash
    GOOS=darwin GOARCH=arm64 go build -o coremail-brute-macos-arm64 .
    ```

编译完成后，您会在当前目录下找到指定名称的可执行文件。

## 使用方法

```
./coremail-brute -url <target_url> [options]
```

默认情况下，工具会读取程序所在目录下的 `users.txt` 和 `passwords.txt` 文件作为用户和密码字典。

### 参数说明

- `-url` (必需): 目标 Coremail 服务器地址。例如: `https://mail.example.com`
- `-c` (可选): 并发线程数。默认为 `50`。
- `-proxy` (可选): SOCKS5 代理地址。例如: `127.0.0.1:10808`
- `-proxyuser` (可选): SOCKS5 代理的用户名（如果代理需要验证）。
- `-proxypass` (可选): SOCKS5 代理的密码（如果代理需要验证）。

### 使用示例

**1. 不使用代理**

确保 `users.txt` 和 `passwords.txt` 文件与可执行文件在同一目录下。

```bash
./coremail-brute -url https://mail.my-company.com -c 100
```

**2. 使用无需验证的 SOCKS5 代理**

```bash
./coremail-brute -url https://mail.my-company.com -proxy 127.0.0.1:10808
```

**3. 使用需要用户名密码验证的 SOCKS5 代理**

```bash
./coremail-brute -url https://mail.my-company.com -proxy proxy.example.com:12345 -proxyuser your_user -proxypass your_pass
```

## 免责声明

本工具仅用于合法的安全测试和研究目的。请确保您已经获得目标系统的明确授权。任何未经授权的测试行为都可能违反法律。工具开发者不对任何非法使用本工具所造成的后果负责。
