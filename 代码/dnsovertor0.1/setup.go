package dnsovertor

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("dnsovertor", setup) }

func setup(c *caddy.Controller) error {

	var torsocksport string
	var hsaddr string
	torsocksport = "9050"
	hsaddr="dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion"
	for c.Next() {
		for c.NextBlock() {
			option := c.Val()
			switch option {
			case "hsaddr":
				if !c.NextArg() {
					return c.ArgErr()
				}
				hsaddr = c.Val()
				log.Debugf("Setting hidden service address to %s", hsaddr)
				log.Infof("Setting hidden service address to %s", hsaddr)

			case "torsocksport":
				if !c.NextArg() {
					return c.ArgErr()
				}
				torsocksport = c.Val()

			default:
				return plugin.Error("dnsovertor", c.Errf("unexpected '%v' command", option))
			}
		}
		// Default (The most common) tor socks port is 9050
		
		log.Debugf("Setting Tor socks port to %s", torsocksport)
		log.Infof("Setting Tor socks port to  %s", torsocksport)
		if c.NextArg() {
			return plugin.Error("dnsovertor", c.ArgErr())
		}

	} 
	dnsserver.GetConfig(c).
		AddPlugin(func(next plugin.Handler) plugin.Handler {
			return NewDnsovertorPlugin(next, hsaddr, torsocksport) //'Dnsovertor'
		})

	return nil
}
