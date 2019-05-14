package gateway

import "github.com/youpenglai/apix/proxy"

var (
	CallService = proxy.CallService
)

// 将请求转发到GRPC服务上
func init() {
	proxy.LoadAllProxy()
}

