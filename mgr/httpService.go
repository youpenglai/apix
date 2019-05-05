package mgr

import (
	"github.com/youpenglai/apix/gateway"
	"sync"
	"errors"
)

var (
	services map[string]*gateway.ApiGateway
	servicesMu sync.Mutex

	ErrHttpServiceNotExists = errors.New("http service not exists")
)

func AddHttpService(serviceName string, opts *gateway.ApiGatewayOpts) {
	gw := gateway.NewApiGateWay(opts)
	servicesMu.Lock()
	defer servicesMu.Unlock()
	services[serviceName] = gw
}

func GetHttpService(serviceName string) (gw *gateway.ApiGateway, err error) {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	var exist bool
	if gw, exist = services[serviceName]; !exist {
		err = ErrHttpServiceNotExists
	}

	return
}

func init() {
	services = make(map[string]*gateway.ApiGateway)
}
