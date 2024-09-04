package dnsovertor

import (
	"fmt"
	"os"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/logger"
	"gopkg.in/yaml.v3"
)

// domain_config配置文件的两个值
type Rule struct {
	Type  string
	Value string
}

// Config 结构体用于映射 YAML 文件内容
type Config struct {
	Whitelist  []Rule `yaml:"whitelist"`
	Blacklist  []Rule `yaml:"blacklist"`
	Directlist []Rule `yaml:"directlist"`
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func init() { plugin.Register("dnsovertor", setup) }

func setup(c *caddy.Controller) error {

	var torsocksport string
	var hsaddr string
	var config Config
	c.Next()
	// Default (The most common) tor socks port is 9050
	torsocksport = "9050"
	// config := ""
	configFilePath := ""
	for c.NextBlock() {
		option := c.Val()
		switch option {
		case "hsaddr":
			if !c.NextArg() {
				return c.ArgErr()
			}
			hsaddr = c.Val()
			logger.Debugf("Setting hidden service address to %s", hsaddr)
			logger.Infof("Setting hidden service address to %s", hsaddr)

		case "torsocksport":
			if !c.NextArg() {
				return c.ArgErr()
			}
			torsocksport = c.Val()

		case "domain_config":
			if !c.NextArg() {
				return c.ArgErr()
			}
			configFilePath := c.Val()
			// 域名配置，设置黑白名单
			config, err := ReadConfig(configFilePath)
			if err != nil {
				logger.Fatalf("Error reading config file: %v", err)
			}
		default:
			return plugin.Error("dnsovertor", c.Errf("unexpected '%v' command", option))
		}
	}
	if configFilePath == "" {
		config := &Config{}
	}
	logger.Debugf("Setting Tor socks port to %s", torsocksport)
	logger.Infof("Setting Tor socks port to %s", torsocksport)
	if c.NextArg() {
		return plugin.Error("dnsovertor", c.ArgErr())
	}

	// 'dnsovertor'
	dnsserver.GetConfig(c).
		AddPlugin(func(next plugin.Handler) plugin.Handler {
			return NewDnsovertorPlugin(next, hsaddr, torsocksport, config) //'Dnsovertor'
		})

	return nil
}

func ReadConfig(filePath string) (*Config, error) {
	// 读取 YAML 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 创建 Config 实例
	var config Config

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
