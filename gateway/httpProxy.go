package gateway

import (
	"github.com/youpenglai/apix/http"

)

// HTTP服务转发
// 正常转发，对参数进行处理后向后转发
// 透传转发，将所有参数包括Header所有内容直接透传向后转发



func HttpPassthrough(ctx *http.Context, target *string) error {

}

