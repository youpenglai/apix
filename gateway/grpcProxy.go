package gateway

import "github.com/youpenglai/apix/grpc"

var (
	CallGRPCService = grpc.CallService
)

// 将请求转发到GRPC服务上
func init() {
	grpc.LoadAllProxy()
}

