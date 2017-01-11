package presence

import (
	"github.com/jurgen-kluft/go-home/presence/netgear"
)

type netgearRouter struct {
	router *netgear.Router
}

func (r *netgearRouter) get(mac map[string]bool) error {
	return r.router.Get(mac)
}

func newNetgearRouter(host, username, password string) provider {
	return &netgearRouter{router: netgear.New(host, username, password)}
}
