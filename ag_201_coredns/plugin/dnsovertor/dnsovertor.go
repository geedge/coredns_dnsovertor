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

// Dnsovertor is a plugin in CoreDNS
type Dnsovertor struct {
	Next        plugin.Handler
	Hsaddr      string
	TorCtrlPort string
}

// default tor control port is 9050
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
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	t0 := time.Now()
	timeStr := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fmt.Printf("dnsovertor解析\n接收请求时间：" + timeStr)
	fmt.Printf("接收请求域名：" + qname + "\n")

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
	httpTransport.Dial = dialer.Dial

	query := dns.Msg{}
	query.SetQuestion(qname, dns.TypeA)
	timeStr1 := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fmt.Printf("发送请求时间：" + timeStr1)
	fmt.Printf("发送请求域名：" + qname + "\n")
	msg, _ := query.Pack()
	b64 := base64.RawURLEncoding.EncodeToString(msg)

	resp, err := httpClient.Get("https://" + p.Hsaddr + "/dns-query?dns=" + b64)
	if err != nil {
		fmt.Printf("Send query error, err:%v\n", err)
		//os.Exit(1) 怎么好像没见过print上面的
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	response := dns.Msg{}
	response.Unpack(bodyBytes)
	timeStr2 := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
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
	timeStr3 := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	fmt.Println("发送响应时间：" + timeStr3)
	fmt.Printf("发送响应信息gogogogo:")
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


