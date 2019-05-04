package gateway

import (
	"github.com/youpenglai/apix/apibuilder"
	apixHttp "github.com/youpenglai/apix/http"
	"sync"
	"errors"
)

var (
	ErrNoApiDocName = errors.New("no api doc name")
	ErrApiContentIsEmpty = errors.New("api content is empty")
)

type ApiGatewayOpts struct {
	bindAddr string
}

var defaultApiGatewayOpts = &ApiGatewayOpts{bindAddr: "127.0.0.1:8080"}

type ApiGateway struct {
	allApiDocs map[string]*apibuilder.ApiDoc
	docMu sync.Mutex
	opts *ApiGatewayOpts

	httpServer *apixHttp.ApiX
}

// 创建新的ApiGateway入口
func NewApiGateWay(opts ...*ApiGatewayOpts) *ApiGateway {
	gatewayOpts := defaultApiGatewayOpts
	if len(opts) > 0 {
		if opts[0].bindAddr != "" {
			gatewayOpts = opts[0]
		}
	}

	return &ApiGateway{
		allApiDocs: make(map[string]*apibuilder.ApiDoc),
		opts: gatewayOpts,
		httpServer: apixHttp.NewApiX(),
	}
}

func (g *ApiGateway) installApis() error {
	return nil
}

// 重新加载ApiGateway
// 当更新ApiDoc后，为了让ApiDoc生效，所以需要对ApiGateWay
func (g *ApiGateway) Reload() error {
	g.Shutdown()
	// reload the server
	if g.httpServer != nil {
		g.httpServer = apixHttp.NewApiX()
	}

	g.Serve()
	return nil
}

// 添加Api文档
// docName, 文档名称，最为唯一标识符，如果遇到相同的名称则覆盖
func (g *ApiGateway) AddApiDoc(docName string, docContent []byte) (err error) {
	g.docMu.Lock()
	defer g.docMu.Unlock()
	if docName == "" {
		err = ErrNoApiDocName
		return
	}
	if len(docContent) == 0 {
		err = ErrApiContentIsEmpty
		return
	}

	doc := apibuilder.NewApiDoc()
	if err = doc.Parse(docContent); err != nil {
		return
	}

	g.allApiDocs[docName] = doc

	return nil
}

// 执行服务
// 该程序会阻塞当前程序直到Shutdown
func (g *ApiGateway) Serve() (err error) {
	if err = g.installApis(); err !=nil {
		return
	}

	err = g.httpServer.Run(g.opts.bindAddr)
	return
}

func (g *ApiGateway) Shutdown() error {
	return g.httpServer.Shutdown()
}