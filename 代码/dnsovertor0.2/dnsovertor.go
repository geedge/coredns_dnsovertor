// Package dnsovertor implements a plugin.
// 实现了： 1.通过domain_config文件配置blacklist,whitelist,directlist.可配置DOMAIN全域名,DOMAIN-SUFFIX域名后缀匹配。
// 预计添加：从mysql数据库中读取记录，和将DNS结果缓存到数据库。
package dnsovertor

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"golang.org/x/net/proxy"

	"github.com/coredns/coredns/plugin"

	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/logger"
)

var logger = clog.NewWithPlugin("dnsovertor")

type StringBoolMap map[string]bool

// Dnsovertor is a plugin in CoreDNS. Details are in setup.go.
type Dnsovertor struct {
	Next                 plugin.Handler
	Hsaddr               string
	TorCtrlPort          string
	BlacklistMap         StringBoolMap
	SuffixBlacklistTrie  *Trie
	WhitelistMap         StringBoolMap
	SuffixWhitelistTrie  *Trie
	DirectlistMap        StringBoolMap
	SuffixDirectlistTrie *Trie
	//Enabled              bool
}

// ServeDNS implements the plugin.Handler interface.
func (p Dnsovertor) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	dotorCount.WithLabelValues(metrics.WithServer(ctx), state.IP(), state.Name()).Inc()

	t0 := time.Now()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("dnsovertor解析\n接收请求时间：" + timeStr)
	fmt.Printf("接收请求域名：" + qname + "\n")

	if p.isBlacklisted(qname) {
		logger.Debugf("Domain %s is Blacklisted, return NameEorror.\n", qname)
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeNameError)
		m.SetReply(r)
		w.WriteMsg(m)

		// Warn and increment the counter for the hit
		blacklistCount.WithLabelValues(metrics.WithServer(ctx), state.IP(), state.Name()).Inc()
		logger.Warning("host ", state.IP(), " requested blacklist domain: ", state.Name())

		return dns.RcodeNameError, nil
	} else {
		fmt.Printf("Domain %s is not blacklisted\n", qname)
	}

	if p.isWhitelisted(qname) {
		logger.Debugf("Domain %s is Whitelisted, continue.\n", qname)
		whitelistCount.WithLabelValues(metrics.WithServer(ctx), state.IP(), state.Name()).Inc()
		logger.Warning("host ", state.IP(), " requested blacklist domain: ", state.Name())
	}

	if p.isDirectlisted(qname) {
		// 简单地将请求传递给下一个插件，不需要在响应返回时做任何处理，直接调用下一个插件。
		logger.Debugf("Domain %s is Directlisted, use forward plugin.\n", qname)
		directlistCount.WithLabelValues(metrics.WithServer(ctx), state.IP(), state.Name()).Inc()
		logger.Warning("host ", state.IP(), " requested blacklist domain: ", state.Name())
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+p.TorCtrlPort, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpTransport.DisableKeepAlives = false
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}

	// set our socks5 as the dialer
	// httpTransport.Dial = dialer.Dial
	httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	query := dns.Msg{}
	query.SetQuestion(qname, dns.TypeA)
	timeStr1 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("发送请求时间：" + timeStr1)
	fmt.Printf("发送请求域名：" + qname + "\n")
	msg, _ := query.Pack()
	b64 := base64.RawURLEncoding.EncodeToString(msg)

	resp, err := httpClient.Get("https://" + p.Hsaddr + "/dns-query?dns=" + b64)
	if err != nil {
		fmt.Printf("Send query error, err:%v\n", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	response := dns.Msg{}
	response.Unpack(bodyBytes)
	timeStr2 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("接收响应时间：" + timeStr2)
	fmt.Println("接收响应信息:")
	for i := 0; i < len(response.Answer); i++ {

		b := response.Answer[i].String()
		//a := response.Answer[i].Header()
		//fmt.Printf("Dns answer name is :" + a.Name + "\n")
		fmt.Printf(b + "\t")
	}
	fmt.Printf("\n")

	//m := new(dns.Msg)
	response.SetReply(r) //m是r的reply
	//m.Authoritative = true
	//m.Answer = answers
	var msg22 *dns.Msg
	msg22 = &response
	timeStr3 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("发送响应时间：" + timeStr3)
	fmt.Printf("=====发送响应信息:")
	for i := 0; i < len(response.Answer); i++ {
		b := response.Answer[i].String()
		fmt.Printf(b + "\t")
	}
	t1 := time.Now()
	d := t1.Sub(t0)
	fmt.Printf("处理时间：")
	fmt.Println(d)

	// 实际写入结果
	w.WriteMsg(msg22)
	dotorCount.WithLabelValues(metrics.WithServer(ctx)).Inc()
	return dns.RcodeSuccess, nil
	// 没返回常规的plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (p Dnsovertor) Name() string { return "dnsovertor" }

func (p Dnsovertor) isBlacklisted(domain string) bool {
	if !p.Enabled {
		return false
	}

	// 1. 直接匹配全域名
	if p.BlacklistMap != nil {
		if p.BlacklistMap[domain] {
			return true
		}
	}

	// 2. 后缀匹配
	if p.SuffixBlacklistTrie != nil {
		parts := strings.Split(domain, ".")
		for i := 0; i < len(parts); i++ {
			suffix := strings.Join(parts[i:], ".")
			if p.SuffixBlacklistTrie.Search(reverse(suffix)) {
				return true
			}
		}
	}

	return false
}

func (p Dnsovertor) isWhitelisted(domain string) bool {
	if !p.Enabled {
		return false
	}

	// 1. 直接匹配全域名
	if p.WhitelistMap != nil {
		if p.WhitelistMap[domain] {
			return true
		}
	}

	// 2. 后缀匹配
	if p.SuffixWhitelistTrie != nil {
		parts := strings.Split(domain, ".")
		for i := 0; i < len(parts); i++ {
			suffix := strings.Join(parts[i:], ".")
			if p.SuffixWhitelistTrie.Search(reverse(suffix)) {
				return true
			}
		}
	}

	return false
}

func (p Dnsovertor) isDirectlisted(domain string) bool {
	if !p.Enabled {
		return false
	}

	// 1. 直接匹配全域名
	if p.DirectlistMap != nil {
		if p.DirectlistMap[domain] {
			return true
		}
	}

	// 2. 后缀匹配
	if p.SuffixDirectlistTrie != nil {
		parts := strings.Split(domain, ".")
		for i := 0; i < len(parts); i++ {
			suffix := strings.Join(parts[i:], ".")
			if p.SuffixDirectlistTrie.Search(reverse(suffix)) {
				return true
			}
		}
	}

	return false
}
