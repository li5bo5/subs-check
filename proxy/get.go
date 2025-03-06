// Package proxies 提供代理服务相关的功能，包括订阅获取、解析和处理
package proxies

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"log/slog"

	"github.com/li5bo5/subs-check/config"
	"github.com/li5bo5/subs-check/proxy/parser"
	"github.com/li5bo5/subs-check/utils"
	"gopkg.in/yaml.v3"
)

// GetProxies 获取并处理所有订阅链接中的代理信息
// 使用并发处理以提高效率，支持多种订阅格式
// 返回值:
//   - []map[string]any: 解析后的代理列表
//   - error: 处理过程中的错误信息
func GetProxies() ([]map[string]any, error) {
	// 记录开始处理的订阅数量
	slog.Info(fmt.Sprintf("当前设置订阅链接数量: %d", len(config.GlobalConfig.SubUrls)))

	// 准备订阅链接列表，转换为 interface{} 类型以适配线程池
	subUrls := make([]interface{}, len(config.GlobalConfig.SubUrls))
	for i, url := range config.GlobalConfig.SubUrls {
		subUrls[i] = url
	}

	// 计算最优并发数：配置的并发数和订阅数量的较小值
	numWorkers := min(len(subUrls), config.GlobalConfig.Concurrent)
	slog.Info(fmt.Sprintf("使用并发数: %d", numWorkers))

	// 创建并配置线程池
	pool := utils.NewThreadPool(numWorkers, taskGetProxies)
	pool.Start()                // 启动工作线程
	pool.AddTaskArgs(subUrls)   // 添加订阅处理任务
	pool.Wait()                 // 等待所有任务完成

	// 收集处理结果
	results := pool.GetResults()
	var mihomoProxies []map[string]any

	// 遍历处理结果，合并代理信息
	for _, result := range results {
		if result.Err != nil {
			slog.Error(fmt.Sprintf("处理订阅链接错误: %v", result.Err))
			continue
		}
		if result.Result != nil {
			proxies := result.Result.([]map[string]any)
			mihomoProxies = append(mihomoProxies, proxies...)
		}
	}

	// 记录处理完成的代理数量
	slog.Info(fmt.Sprintf("共获取到 %d 个代理", len(mihomoProxies)))
	return mihomoProxies, nil
}

// taskGetProxies 处理单个订阅链接的任务函数
// 支持处理多种格式的订阅内容，包括 YAML 格式和 URI 格式
// 参数:
//   - args: 订阅链接字符串，类型为 interface{}
// 返回值:
//   - interface{}: 解析后的代理列表
//   - error: 处理过程中的错误
func taskGetProxies(args interface{}) (interface{}, error) {
	subUrl := args.(string)
	var mihomoProxies []map[string]any

	// 获取订阅内容
	data, err := GetDateFromSubs(subUrl)
	if err != nil {
		slog.Error(fmt.Sprintf("获取订阅链接错误: %v", err))
		return nil, err
	}
	slog.Debug(fmt.Sprintf("获取订阅链接: %s，数据长度: %d", subUrl, len(data)))

	// 尝试解析为 YAML 格式
	var config map[string]any
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// YAML 解析失败，尝试处理 URI 格式
		reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|vless|hysteria|hy2|hysteria2)://")
		// 如果不匹配 URI 格式，尝试 base64 解码
		if !reg.Match(data) {
			data = []byte(parser.DecodeBase64(string(data)))
		}
		// 如果是 URI 格式，按行分割处理
		if reg.Match(data) {
			proxies := strings.Split(string(data), "\n")

			// 逐行解析代理配置
			for _, proxy := range proxies {
				parseProxy, err := ParseProxy(proxy)
				if err != nil {
					slog.Debug(fmt.Sprintf("解析proxy错误: %s , %v", proxy, err))
					continue
				}
				// 跳过空代理
				if parseProxy == nil {
					continue
				}
				mihomoProxies = append(mihomoProxies, parseProxy)
			}
			return mihomoProxies, nil
		}
	}

	// 处理 YAML 格式的订阅内容
	proxyInterface, ok := config["proxies"]
	if !ok || proxyInterface == nil {
		return nil, fmt.Errorf("订阅链接没有proxies: %s", subUrl)
	}

	// 转换代理列表
	proxyList, ok := proxyInterface.([]any)
	if !ok {
		return nil, fmt.Errorf("proxies 格式错误: %s", subUrl)
	}

	// 处理每个代理配置
	for _, proxy := range proxyList {
		proxyMap, ok := proxy.(map[string]any)
		if !ok {
			continue
		}
		mihomoProxies = append(mihomoProxies, proxyMap)
	}

	return mihomoProxies, nil
}

// GetDateFromSubs 从订阅链接获取数据内容
// 支持自动重试，使用自定义 User-Agent 避免被屏蔽
// 参数:
//   - subUrl: 订阅链接地址
// 返回值:
//   - []byte: 获取到的数据内容
//   - error: 获取过程中的错误
func GetDateFromSubs(subUrl string) ([]byte, error) {
	// 从配置获取重试次数
	maxRetries := config.GlobalConfig.SubUrlsReTry
	var lastErr error

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 重试循环
	for i := 0; i < maxRetries; i++ {
		// 重试间隔
		if i > 0 {
			time.Sleep(time.Second)
		}

		// 创建 HTTP 请求
		req, err := http.NewRequest("GET", subUrl, nil)
		if err != nil {
			lastErr = err
			continue
		}
		// 如果走clash，那么输出base64的时候还要更改每个类型的key，所以不能走，以后都走URI
		// 如果用户想使用clash源，那可以在订阅链接结尾加上 &flag=clash.meta
		// 设置 User-Agent 模拟浏览器访问，防止被屏蔽
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("订阅链接: %s 返回状态码: %d", subUrl, resp.StatusCode)
			continue
		}

		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}
		return body, nil
	}

	// 所有重试都失败后返回错误
	return nil, fmt.Errorf("重试%d次后失败: %v", maxRetries, lastErr)
}

// min 返回两个整数中的较小值
// 用于优化并发数量，避免创建过多无用的工作线程
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
