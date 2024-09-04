package dnsovertor

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 向dns隐藏服务转发的请求的数量
var dotorCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "dnsovertor",
	Name:      "dnsovertor_hits_total",
	Help:      "Counter of the number of requests made to warnlisted domains.",
}, []string{"server", "requestor", "domain"})

