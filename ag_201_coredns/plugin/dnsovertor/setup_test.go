// 注意：测试文件Torsocksport时需要相应端口的tor运行
package dnsovertor

import (
	"reflect"
	"testing"

	"github.com/coredns/caddy"
	"github.com/stretchr/testify/assert"
)

func TestValidHsaddr(t *testing.T) {
	cfg := `dnsovertor {
		hsaddr dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion
    }`
	c := caddy.NewTestController("dns", cfg)

	err := setup(c)
	assert.NoError(t, err)
}

func TestValidTorsocksport(t *testing.T) {
	cfg := `dnsovertor {
		hsaddr dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion
        torsocksport 9050
    }`
	c := caddy.NewTestController("dns", cfg)

	err := setup(c)
	assert.NoError(t, err)
}

func TestEmptyHsaddr(t *testing.T) {
	cfg := `dnsovertor {
		hsaddr 
    }`
	c := caddy.NewTestController("dns", cfg)

	err := setup(c)
	t.Fatalf("Emplty Hsaddr is not supported: %v", err)
}

func TestEmptyTorsocksport(t *testing.T) {
	cfg := `dnsovertor {
		hsaddr dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion
		torsocksport
    }`
	c := caddy.NewTestController("dns", cfg)

	err := setup(c)
	t.Fatalf("Emplty torsocksport is not supported: %v", err)
}
