package check

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"log/slog"

	"github.com/li5bo5/subs-check/check/platfrom"
	"github.com/li5bo5/subs-check/config"
	proxyutils "github.com/li5bo5/subs-check/proxy"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/constant"
)

// Result 存储节点检测结果
type Result struct {
	Proxy      map[string]any
	Google     bool
	Cloudflare bool
}

// ProxyChecker 处理代理检测的主要结构体
type ProxyChecker struct {
	results     []Result
	proxyCount  int
	threadCount int
	progress    int32
	available   int32
	resultChan  chan Result
	tasks       chan map[string]any
	mu          sync.Mutex  // 添加互斥锁保护 results
}

// NewProxyChecker 创建新的检测器实例
func NewProxyChecker(proxies []map[string]any) *ProxyChecker {
	proxyCount := len(proxies)
	threadCount := config.GlobalConfig.Concurrent
	if proxyCount < threadCount {
		threadCount = proxyCount
	}

	return &ProxyChecker{
		results:     make([]Result, 0),
		proxyCount:  proxyCount,
		threadCount: threadCount,
		resultChan:  make(chan Result),
		tasks:       make(chan map[string]any, proxyCount),
	}
}

// Check 执行代理检测的主函数
func Check() ([]Result, error) {
	proxies, err := proxyutils.GetProxies()
	if err != nil {
		return nil, fmt.Errorf("获取节点失败: %w", err)
	}
	slog.Info(fmt.Sprintf("获取节点数量: %d", len(proxies)))

	proxies = proxyutils.DeduplicateProxies(proxies)
	slog.Info(fmt.Sprintf("去重后节点数量: %d", len(proxies)))

	checker := NewProxyChecker(proxies)
	return checker.run(proxies)
}

// Run 运行检测流程
func (pc *ProxyChecker) run(proxies []map[string]any) ([]Result, error) {
	slog.Info("开始检测节点")
	slog.Info(fmt.Sprintf("启动工作线程: %d", pc.threadCount))

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GlobalConfig.Timeout*len(proxies)) * time.Millisecond)
	defer cancel()

	progressDone := make(chan bool)
	defer close(progressDone)

	if config.GlobalConfig.PrintProgress {
		go pc.showProgress(progressDone)
	}

	var wg sync.WaitGroup
	// 启动工作线程
	for i := 0; i < pc.threadCount; i++ {
		wg.Add(1)
		go pc.worker(ctx, &wg)
	}

	// 发送任务
	go pc.distributeProxies(proxies)
	slog.Debug(fmt.Sprintf("发送任务: %d", len(proxies)))

	// 收集结果
	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		pc.collectResults()
		collectWg.Done()
	}()

	// 等待所有工作线程完成或超时
	workersDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(pc.resultChan)
		close(workersDone)
	}()

	select {
	case <-ctx.Done():
		slog.Warn("检测超时，正在终止...")
		return pc.results, ctx.Err()
	case <-workersDone:
		// 等待结果收集完成
		collectWg.Wait()
		// 等待进度条显示完成
		time.Sleep(100 * time.Millisecond)

		if config.GlobalConfig.PrintProgress {
			progressDone <- true
		}
		slog.Info(fmt.Sprintf("可用节点数量: %d", len(pc.results)))
		return pc.results, nil
	}
}

// worker 处理单个代理检测的工作线程
func (pc *ProxyChecker) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case proxy, ok := <-pc.tasks:
			if !ok {
				return
			}
			if result := pc.checkProxy(proxy); result != nil {
				pc.resultChan <- *result
			}
			pc.incrementProgress()
		}
	}
}

// checkProxy 检测单个代理
func (pc *ProxyChecker) checkProxy(proxy map[string]any) *Result {
	httpClient := CreateClient(proxy)
	if httpClient == nil {
		slog.Debug(fmt.Sprintf("创建代理Client失败: %v", proxy["name"]))
		return nil
	}
	defer func() {
		if transport, ok := httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}()

	res := &Result{
		Proxy: proxy,
	}

	if os.Getenv("SUB_CHECK_SKIP") == "true" {
		slog.Debug(fmt.Sprintf("跳过检测代理: %v", proxy["name"]))
		return res
	}

	cloudflare, err := platfrom.CheckCloudflare(httpClient)
	if err != nil {
		slog.Debug(fmt.Sprintf("Cloudflare检测失败 [%v]: %v", proxy["name"], err))
		return nil
	}
	if !cloudflare {
		slog.Debug(fmt.Sprintf("Cloudflare检测不通过: %v", proxy["name"]))
		return nil
	}

	google, err := platfrom.CheckGoogle(httpClient)
	if err != nil {
		slog.Debug(fmt.Sprintf("Google检测失败 [%v]: %v", proxy["name"], err))
		return nil
	}
	if !google {
		slog.Debug(fmt.Sprintf("Google检测不通过: %v", proxy["name"]))
		return nil
	}

	var speed int
	if config.GlobalConfig.SpeedTestUrl != "" {
		speed, err = platfrom.CheckSpeed(httpClient)
		if err != nil {
			slog.Debug(fmt.Sprintf("速度测试失败 [%v]: %v", proxy["name"], err))
			return nil
		}
		if speed < config.GlobalConfig.MinSpeed {
			slog.Debug(fmt.Sprintf("速度测试不达标 [%v]: %d < %d", proxy["name"], speed, config.GlobalConfig.MinSpeed))
			return nil
		}
	}

	pc.incrementAvailable()
	pc.updateProxyName(proxy, httpClient, speed)

	res.Cloudflare = cloudflare
	res.Google = google
	return res
}

// updateProxyName 更新代理名称
func (pc *ProxyChecker) updateProxyName(proxy map[string]any, client *http.Client, speed int) {
	// 获取速度
	if config.GlobalConfig.SpeedTestUrl != "" {
		var speedStr string
		if speed < 1024 {
			speedStr = fmt.Sprintf("%dKB/s", speed)
		} else {
			speedStr = fmt.Sprintf("%.1fMB/s", float64(speed)/1024)
		}
		proxy["name"] = strings.TrimSpace(proxy["name"].(string)) + " | ⬇️ " + speedStr
	}
}

// showProgress 显示进度条
func (pc *ProxyChecker) showProgress(done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println()
			return
		default:
			current := atomic.LoadInt32(&pc.progress)
			available := atomic.LoadInt32(&pc.available)

			if pc.proxyCount == 0 {
				time.Sleep(100 * time.Millisecond)
				break
			}

			// if 0/0 = NaN ,shoule panic
			percent := float64(current) / float64(pc.proxyCount) * 100
			fmt.Printf("\r进度: [%-50s] %.1f%% (%d/%d) 可用: %d",
				strings.Repeat("=", int(percent/2))+">",
				percent,
				current,
				pc.proxyCount,
				available)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// 辅助方法
func (pc *ProxyChecker) incrementProgress() {
	atomic.AddInt32(&pc.progress, 1)
}

func (pc *ProxyChecker) incrementAvailable() {
	atomic.AddInt32(&pc.available, 1)
}

// distributeProxies 分发代理任务
func (pc *ProxyChecker) distributeProxies(proxies []map[string]any) {
	for _, proxy := range proxies {
		pc.tasks <- proxy
	}
	close(pc.tasks)
}

// collectResults 收集检测结果
func (pc *ProxyChecker) collectResults() {
	for result := range pc.resultChan {
		pc.mu.Lock()
		pc.results = append(pc.results, result)
		pc.mu.Unlock()
	}
}

func CreateClient(mapping map[string]any) *http.Client {
	proxy, err := adapter.ParseProxy(mapping)
	if err != nil {
		slog.Debug(fmt.Sprintf("底层mihomo创建代理Client失败: %v", err))
		return nil
	}

	// 计算动态超时时间
	baseTimeout := time.Duration(config.GlobalConfig.Timeout) * time.Millisecond
	dialTimeout := baseTimeout / 3
	if dialTimeout < 1*time.Second {
		dialTimeout = 1 * time.Second
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 使用动态超时时间
			dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
			defer cancel()
			
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			// 转换端口为 uint16
			portNum, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", port)
			}

			metadata := &constant.Metadata{
				Host:    host,
				DstPort: uint16(portNum),
				Type:    constant.HTTP,
			}

			return proxy.DialContext(dialCtx, metadata)
		},
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 10,
	}

	client := &http.Client{
		Timeout:   baseTimeout,
		Transport: transport,
	}

	return client
}
