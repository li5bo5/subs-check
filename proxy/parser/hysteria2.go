package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseHysteria2(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "hysteria2://") && !strings.HasPrefix(data, "hy2://") {
		return nil, fmt.Errorf("不是hysteria2格式")
	}

	// 移除 "hysteria2://" 前缀

	link, err := url.Parse(data)
	if err != nil {
		return nil, err
	}

	query := link.Query()
	server := link.Hostname()
	if server == "" {
		return nil, fmt.Errorf("hysteria2 服务器地址错误")
	}
	portStr := link.Port()
	if portStr == "" {
		return nil, fmt.Errorf("hysteria2 端口错误")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("hysteria2 端口错误")
	}
	_, obfs, obfsPassword, _, insecure, sni := query.Get("network"), query.Get("obfs"), query.Get("obfs-password"), query.Get("pinSHA256"), query.Get("insecure"), query.Get("sni")
	insecureBool := insecure == "1"

	return map[string]any{
		"type":             "hysteria2",
		"name":             link.Fragment,
		"server":           server,
		"port":             port,
		"ports":            query.Get("mport"),
		"password":         link.User.String(),
		"obfs":             obfs,
		"obfs-password":    obfsPassword,
		"sni":              sni,
		"skip-cert-verify": insecureBool,
		// 添加原配置
		"insecure": insecure,
		"mport":    query.Get("mport"),
	}, nil
}
