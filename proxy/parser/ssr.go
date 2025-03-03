package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseSsr(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "ssr://") {
		return nil, fmt.Errorf("不是ssr格式")
	}
	// todo: 这些参数解析应该也有问题
	data = strings.TrimPrefix(data, "ssr://")
	data = DecodeBase64(data)
	serverInfoAndParams := strings.SplitN(data, "/?", 2)
	parts := strings.Split(serverInfoAndParams[0], ":")
	if len(parts) < 6 {
		return nil, fmt.Errorf("ssr 参数错误")
	}
	server := parts[0]
	protocol := parts[2]
	method := parts[3]
	obfs := parts[4]
	password := DecodeBase64(parts[5])
	portStr := parts[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("ssr 端口错误")
	}
	var obfsParam string
	var protoParam string
	var remarks string
	if len(serverInfoAndParams) == 2 {
		params, err := url.ParseQuery(serverInfoAndParams[1])
		if err != nil {
			return nil, fmt.Errorf("ssr 参数错误")
		}
		if params.Get("obfsparam") != "" {
			obfsParam = DecodeBase64(params.Get("obfsparam"))
		}
		if params.Get("protoparam") != "" {
			protoParam = DecodeBase64(params.Get("protoparam"))
		}
		if params.Get("remarks") != "" {
			remarks = DecodeBase64(params.Get("remarks"))
		} else {
			remarks = server + ":" + strconv.Itoa(port)
		}

	}
	return map[string]any{
		"name":           remarks,
		"server":         server,
		"port":           port,
		"password":       password,
		"cipher":         method,
		"obfs":           obfs,
		"obfs-param":     obfsParam,
		"protocol":       protocol,
		"protocol-param": protoParam,
	}, nil
}
