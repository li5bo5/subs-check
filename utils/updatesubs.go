// Package utils 提供了订阅更新相关的工具函数
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"

	"github.com/li5bo5/subs-check/config"
)

// httpClient 定义了一个通用的 HTTP 客户端接口
// 这个接口允许我们在测试时注入模拟的 HTTP 客户端
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// providersResponse 定义了从 API 获取的提供者信息的数据结构
type providersResponse struct {
	// Providers 是一个 map，key 是提供者名称，value 是提供者的详细信息
	Providers map[string]struct {
		// VehicleType 表示提供者的类型（如 "HTTP"）
		VehicleType string `json:"vehicleType"`
	} `json:"providers"`
}

// makeRequest 处理通用的 HTTP 请求逻辑
func makeRequest(client httpClient, method, url string) ([]byte, error) {
	// 创建新的 HTTP 请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 执行 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil // 204 状态码表示成功但没有内容
		}
		return nil, fmt.Errorf("API 返回非 200 状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

// UpdateSubs 是主要的订阅更新函数
func UpdateSubs() {
	// 获取需要更新的订阅名称
	names, err := getNeedUpdateNames(http.DefaultClient)
	if err != nil {
		slog.Error(fmt.Sprintf("获取需要更新的订阅失败: %v", err))
		return
	}

	// 如果没有需要更新的订阅，直接返回
	if len(names) == 0 {
		slog.Info("更新完成，共 0 个节点")
		return
	}

	// 执行订阅更新
	if err := updateSubs(http.DefaultClient, names); err != nil {
		slog.Error(fmt.Sprintf("更新订阅失败: %v", err))
		return
	}
	slog.Info(fmt.Sprintf("更新完成，共 %d 个节点", len(names)))
}

// getNeedUpdateNames 获取需要更新的订阅名称列表
func getNeedUpdateNames(client httpClient) ([]string, error) {
	var names []string
	for _, url := range config.GlobalConfig.SubUrls {
		names = append(names, url)
	}
	return names, nil
}

// updateSubs 更新指定的订阅
func updateSubs(client httpClient, names []string) error {
	for _, name := range names {
		if _, err := makeRequest(client, http.MethodGet, name); err != nil {
			return fmt.Errorf("更新订阅 %s 失败: %w", name, err)
		}
	}
	return nil
}