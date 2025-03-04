package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"

	"github.com/li5bo5/subs-check/config"
)

// 定义通用的 HTTP 客户端接口
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// API 响应的结构体
type providersResponse struct {
	Providers map[string]struct {
		VehicleType string `json:"vehicleType"`
	} `json:"providers"`
}

// makeRequest 处理通用的 HTTP 请求逻辑
func makeRequest(client httpClient, method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GlobalConfig.MihomoApiSecret))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		}
		return nil, fmt.Errorf("API 返回非 200 状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

// UpdateSubs 更新所有HTTP类型的订阅
func UpdateSubs() {
	slog.Info("开始更新订阅...")
	
	names, err := getNeedUpdateNames(http.DefaultClient)
	if err != nil {
		slog.Error(fmt.Sprintf("获取需要更新的订阅失败: %v", err))
		return
	}

	if len(names) == 0 {
		slog.Info("没有需要更新的HTTP类型订阅")
		return
	}

	if err := updateSubs(http.DefaultClient, names); err != nil {
		slog.Error(fmt.Sprintf("更新订阅失败: %v", err))
		return
	}
	slog.Info(fmt.Sprintf("订阅更新完成，共更新 %d 个订阅", len(names)))
}

func getNeedUpdateNames(client httpClient) ([]string, error) {
	url := fmt.Sprintf("%s/providers/proxies", config.GlobalConfig.MihomoApiUrl)
	body, err := makeRequest(client, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var response providersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析提供者信息失败: %w", err)
	}

	var names []string
	for name, provider := range response.Providers {
		if provider.VehicleType == "HTTP" {
			names = append(names, name)
		}
	}
	return names, nil
}

func updateSubs(client httpClient, names []string) error {
	for _, name := range names {
		url := fmt.Sprintf("%s/providers/proxies/%s", config.GlobalConfig.MihomoApiUrl, name)
		if _, err := makeRequest(client, http.MethodPut, url); err != nil {
			slog.Error(fmt.Sprintf("更新订阅%v失败: %v", name, err))
		}
		slog.Info(fmt.Sprintf("成功更新订阅: %s", name))
	}
	return nil
}
