package mgr

import (
	"github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/middlewares"
)

const defaultBindAddr = "127.0.0.1:58081"

// 运行管理端服务
func RunManagerServer(bindAddr ...string) {
	mgrServer := http.NewApiX()

	mgrServer.Use(middlewares.Server())

	addr := defaultBindAddr
	if len(bindAddr) > 0 {
		if bindAddr[0] != "" {
			addr = bindAddr[0]
		}
	}

	mgrServer.Run(addr)
}
