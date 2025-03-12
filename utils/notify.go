package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log/slog"
	
	"github.com/li5bo5/subs-check/config"
)

// NotifyType 表示通知的类型
type NotifyType int

const (
	// TYPE_NotifyStart 程序启动通知
	TYPE_NotifyStart NotifyType = iota
	// TYPE_NotifyResult 结果通知
	TYPE_NotifyResult
	// TYPE_NotifyError 错误通知
	TYPE_NotifyError
)

// NotifyMessage Apprise API消息结构
type NotifyMessage struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Type    string `json:"type"`
	Tag     string `json:"tag,omitempty"`
	Format  string `json:"format,omitempty"`
	Attach  string `json:"attach,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

// SendNotification 发送通知
// 参数:
//   - notifyType: 通知类型 (启动/结果/错误)
//   - title: 通知标题
//   - body: 通知内容
//   - attach: 可选的附件URL (可为空)
// 返回值:
//   - error: 发送过程中的错误信息
func SendNotification(notifyType NotifyType, title string, body string, attach string) error {
	// 检查是否配置了Apprise URL
	if config.GlobalConfig.AppriseURL == "" {
		slog.Debug("未配置Apprise URL，跳过通知")
		return nil
	}

	// 根据通知类型检查是否需要发送
	switch notifyType {
	case TYPE_NotifyStart:
		if !config.GlobalConfig.NotifyOnStart {
			return nil
		}
	case TYPE_NotifyResult:
		if !config.GlobalConfig.NotifyOnResult {
			return nil
		}
	case TYPE_NotifyError:
		if !config.GlobalConfig.NotifyOnError {
			return nil
		}
	}

	// 准备通知消息
	message := NotifyMessage{
		Title:   title,
		Body:    body,
		Type:    getNotifyTypeString(notifyType),
		Tag:     config.GlobalConfig.AppriseTag,
		Format:  "text", // 使用纯文本格式
		Attach:  attach,
		Timeout: 10000, // 10秒超时
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化通知消息失败: %w", err)
	}

	// 发送请求到Apprise API
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	
	req, err := http.NewRequest("POST", config.GlobalConfig.AppriseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建通知请求失败: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送通知请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("通知服务返回错误: %s", resp.Status)
	}
	
	slog.Info(fmt.Sprintf("成功发送 %s 类型通知", getNotifyTypeString(notifyType)))
	return nil
}

// getNotifyTypeString 返回通知类型对应的字符串
func getNotifyTypeString(notifyType NotifyType) string {
	switch notifyType {
	case TYPE_NotifyStart:
		return "info"
	case TYPE_NotifyResult:
		return "success"
	case TYPE_NotifyError:
		return "failure"
	default:
		return "info"
	}
}

// NotifyStart 发送程序启动通知
func NotifyStart() {
	hostname, _ := GetHostname()
	title := fmt.Sprintf("🚀 订阅检测工具启动通知 [%s]", hostname)
	body := fmt.Sprintf("订阅检测工具已启动\n时间: %s\n节点数量: %d", 
		time.Now().Format("2006-01-02 15:04:05"),
		len(config.GlobalConfig.SubUrls))
	
	if err := SendNotification(TYPE_NotifyStart, title, body, ""); err != nil {
		slog.Error(fmt.Sprintf("发送启动通知失败: %v", err))
	}
}

// NotifyResult 发送处理结果通知
func NotifyResult(totalCount int, availableCount int, imageURL string) {
	hostname, _ := GetHostname()
	title := fmt.Sprintf("✅ 订阅检测完成 [%s]", hostname)
	body := fmt.Sprintf("订阅处理完成\n时间: %s\n总节点: %d\n可用节点: %d\n可用率: %.2f%%", 
		time.Now().Format("2006-01-02 15:04:05"),
		totalCount,
		availableCount,
		float64(availableCount)/float64(totalCount)*100)
	
	if err := SendNotification(TYPE_NotifyResult, title, body, imageURL); err != nil {
		slog.Error(fmt.Sprintf("发送结果通知失败: %v", err))
	}
}

// NotifyError 发送错误通知
func NotifyError(errorMsg string) {
	hostname, _ := GetHostname()
	title := fmt.Sprintf("❌ 订阅检测错误 [%s]", hostname)
	body := fmt.Sprintf("订阅处理出错\n时间: %s\n错误信息: %s", 
		time.Now().Format("2006-01-02 15:04:05"),
		errorMsg)
	
	if err := SendNotification(TYPE_NotifyError, title, body, ""); err != nil {
		slog.Error(fmt.Sprintf("发送错误通知失败: %v", err))
	}
}

// GetHostname 获取主机名
func GetHostname() (string, error) {
	// 这里实现获取主机名的逻辑
	// 如果获取失败，返回默认主机名
	return "订阅服务器", nil
} 