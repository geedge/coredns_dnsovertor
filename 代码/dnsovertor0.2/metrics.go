package dnsovertor

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 向dns隐藏服务转发的请求的数量
// 插件处理的请求数量
var dotorCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "dnsovertor",
	Name:      "dnsovertor_hits_total",
	Help:      "Counter of the number of requests made to dotor domains.",
}, []string{"server", "requestor", "domain"})

// TODO:服务时间

var blacklistCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "blacklist",
	Name:      "blacklist_hits_total",
	Help:      "Counter of the number of requests made to blacklisted domains.",
}, []string{"server", "requestor", "domain"})

var whitelistCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "whitelist",
	Name:      "whitelist_hits_total",
	Help:      "Counter of the number of requests made to whitelisted domains.",
}, []string{"server", "requestor", "domain"})

var directlistCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "directlist",
	Name:      "directlist_hits_total",
	Help:      "Counter of the number of requests made to directlisted domains.",
}, []string{"server", "requestor", "domain"})
