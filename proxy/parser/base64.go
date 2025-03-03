package parser

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"
)

// IsBase64String 检查字符串是否为有效的 base64 编码
func IsBase64String(s string) bool {
	// 检查基本的 base64 特征
	s = strings.TrimSpace(s)

	// 字符串长度不能为空
	if len(s) == 0 {
		return false
	}

	// 先处理 URL 安全的 base64
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	// 添加缺失的填充
	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}

	// 检查是否包含非 base64 字符（不包括填充字符）
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	sWithoutPadding := strings.TrimRight(s, "=")
	for _, c := range sWithoutPadding {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	// 尝试解码
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return false
	}

	// 检查解码后的内容是否为有效的 UTF-8 字符串
	return utf8.Valid(decoded)
}

// DecodeBase64 解码 base64 字符串，如果不是 base64 则返回原始字符串
func DecodeBase64(s string) string {
	if !IsBase64String(s) {
		return s
	}

	// 处理 URL 安全的 base64
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	// 添加缺失的填充
	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}

	// 解码
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}

	// 检查是否为有效的 UTF-8 字符串
	if !utf8.Valid(decoded) {
		return s
	}

	return string(decoded)
}
