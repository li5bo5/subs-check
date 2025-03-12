package config

import _ "embed"

type Config struct {
	PrintProgress      bool     `yaml:"print-progress"`
	Concurrent         int      `yaml:"concurrent"`
	CheckInterval      int      `yaml:"check-interval"`
	SpeedTestUrl       string   `yaml:"speed-test-url"`
	DownloadTimeout    int      `yaml:"download-timeout"`
	MinSpeed           int      `yaml:"min-speed"`
	Timeout            int      `yaml:"timeout"`
	FilterRegex        string   `yaml:"filter-regex"`
	SaveMethod         string   `yaml:"save-method"`
	GithubToken        string   `yaml:"github-token"`
	GithubGistID       string   `yaml:"github-gist-id"`
	GithubAPIMirror    string   `yaml:"github-api-mirror"`
	SubUrlsReTry       int      `yaml:"sub-urls-retry"`
	SubUrls            []string `yaml:"sub-urls"`
	ListenPort         string   `yaml:"listen-port"`
	// 通知相关配置
	AppriseURL         string   `yaml:"apprise-url"`
	AppriseTag         string   `yaml:"apprise-tag"`
	NotifyOnStart      bool     `yaml:"notify-on-start"`
	NotifyOnResult     bool     `yaml:"notify-on-result"`
	NotifyOnError      bool     `yaml:"notify-on-error"`
}

var GlobalConfig = &Config{
	// 新增配置，给未更改配置文件的用户一个默认值
	ListenPort: ":8299",
	// 默认不启用通知
	NotifyOnStart:  false,
	NotifyOnResult: false,
	NotifyOnError:  true,
}

//go:embed config.example.yaml
var DefaultConfigTemplate []byte
