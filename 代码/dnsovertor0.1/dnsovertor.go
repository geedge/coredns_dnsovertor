// Package dnsovertor implements a plugin.
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

	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("dnsovertor")

// Dnsovertor is a plugin in CoreDNS。用于匿名解析网关。
// 在启用cache插件的情况下，初次处理DNS请求后，本地留有缓存（一块儿储存结果的全局变量），第二次将不会调用到这个dnsovertor插件。
// Next为固定格式。Hsaddr表示隐藏服务地址。以.onion结尾。 TorCtrlPort应配置为已经启动的tor进程中torrc配置的ControlPort。
type Dnsovertor struct {
	Next        plugin.Handler
	Hsaddr      string
	TorCtrlPort string
}

// 未经特殊配置的默认tor control port is 9050
func NewDnsovertorPlugin(next plugin.Handler, Hsaddr string, TorCtrlPort string) *Dnsovertor {

	log.Debugf(
		"Creating dnsovertor plugin with %s",
		Hsaddr,
	)

	p := &Dnsovertor{
		Hsaddr:      Hsaddr,
		Next:        next,
		TorCtrlPort: TorCtrlPort,
	}

	return p
}

// ServeDNS implements the plugin.Handler interface.
func (p Dnsovertor) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// 固定写法，获取请求信息
	state := request.Request{W: w, Req: r}
	qname := state.Name()
	// 标准输出中打印日志
	t0 := time.Now()
	timeStr := time.Now().Format("2006-01-02 15:04:05") 
	fmt.Printf("dnsovertor解析\n接收请求时间：" + timeStr)
	fmt.Printf("接收请求域名：" + qname + "\n")

	//创建HTTP客户端
	//创建socks5代理拨号器
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+p.TorCtrlPort, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	//设置http client。禁用TLS验证。启用 HTTP Keep-Alive，允许连接重用。设置超时时间为 10 秒。
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpTransport.DisableKeepAlives = false
	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}
	// 将socks5代理设置为 拨号方法，这样 HTTP 客户端就会通过 SOCKS5 代理发送请求。
	httpTransport.Dial = dialer.Dial

	//创建DNS查询请求，查询A记录（实际上应该查询原请求的记录类型）
	query := dns.Msg{}
	query.SetQuestion(qname, dns.TypeA)
	timeStr1 := time.Now().Format("2006-01-02 15:04:05") 
	fmt.Printf("发送请求时间：" + timeStr1)
	fmt.Printf("发送请求域名：" + qname + "\n")
	msg, _ := query.Pack()
	b64 := base64.RawURLEncoding.EncodeToString(msg)


	// 发送HTTPS请求
	resp, err := httpClient.Get("https://" + p.Hsaddr + "/dns-query?dns=" + b64)
	if err != nil {
		fmt.Printf("Send query error, err:%v\n", err)
	}

	defer resp.Body.Close()// 资源清理
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	response := dns.Msg{}
	response.Unpack(bodyBytes)

	timeStr2 := time.Now().Format("2006-01-02 15:04:05") 
	fmt.Println("接收响应时间：" + timeStr2)
	fmt.Println("接收响应信息:")
	for i := 0; i < len(response.Answer); i++ {
		//a := response.Answer[i].Header()
		b := response.Answer[i].String()
		//fmt.Printf("Dns answer name is :" + a.Name + "\n")
		fmt.Printf(b + "\t")
	}
	fmt.Printf("\n")
	//fmt.Printf("Dns answer is :%v\n", response.String())

	//m := new(dns.Msg)
	response.SetReply(r) //m是r的reply
	//m.Authoritative = true
	//m.Answer = answers
	var msg22 *dns.Msg
	msg22 = &response
	timeStr3 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("发送响应时间：" + timeStr3)
	fmt.Printf("发送响应信息:")
	for i := 0; i < len(response.Answer); i++ {
		//a := response.Answer[i].Header()
		b := response.Answer[i].String()
		//fmt.Printf("Dns answer name is :" + a.Name + "\n")
		fmt.Printf(b + "\t")
	}
	t1 := time.Now()
	d := t1.Sub(t0)
	fmt.Printf("处理时间：")
	fmt.Println(d)

	w.WriteMsg(msg22)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (p Dnsovertor) Name() string { return "dnsovertor" }


