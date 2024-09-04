package dnsovertor

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"golang.org/x/net/proxy"

	"github.com/coredns/coredns/plugin"

	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

// default tor socksport is 9050   
// func NewDnsovertorPlugin(next plugin.Handler, Hsaddr string, TorCtrlPort string, Blacklist string, Whitelist string, Directlist string) *Dnsovertor {
	func NewDnsovertorPlugin(next plugin.Handler, Hsaddr string, TorCtrlPort string, config Config) *Dnsovertor {
		log.Debugf(
			"Creating dnsovertor plugin with HsDNS:%s, TorCtrlPort:%s, %d long Blacklist, %d long Whitelist, %d long Directlist.",
			Hsaddr,
			TorCtrlPort,
			len(Blacklist),
			len(Whitelist),
			len(Directlist)
		)
	
		p := &Dnsovertor{
			Hsaddr:            Hsaddr,
			Next:              next,
			TorCtrlPort:       TorCtrlPort,
			DomainBlacklistMap:  make(map[string]bool),
			SuffixBlacklistTrie: NewTrie(),
			DomainWhitelistMap:  make(map[string]bool),
			SuffixWhitelistTrie: NewTrie(),
			DomainDirectlistMap: make(map[string]bool),
			SuffixDirectlistTrie: NewTrie(),
			Enabled:             false
		}

		// 初始化黑名单
		if config.Blacklist != nil {
			p.Enabled = true
			for _, rule := range config.Blacklist {
				if rule.Type == "DOMAIN" {
					p.DomainBlacklistMap[rule.Value] = true
				} else if rule.Type == "DOMAIN-SUFFIX" {
					p.SuffixBlacklistTrie.Insert(reverse(rule.Value), nil)
				}
			}
		}
	
		// 初始化白名单
		if config.Whitelist != nil {
			p.Enabled = true
			for _, rule := range config.Whitelist {
				if rule.Type == "DOMAIN" {
					p.DomainWhitelistMap[rule.Value] = true
				} else if rule.Type == "DOMAIN-SUFFIX" {
					p.SuffixWhitelistTrie.Insert(reverse(rule.Value), nil)
				}
			}
		}
	
		// 初始化直通名单
		if config.Directlist != nil {
			p.Enabled = true
			for _, rule := range config.Directlist {
				if rule.Type == "DOMAIN" {
					p.DomainDirectlistMap[rule.Value] = true
				} else if rule.Type == "DOMAIN-SUFFIX" {
					p.SuffixDirectlistTrie.Insert(reverse(rule.Value), nil)
				}
			}
		}
		return p
	}