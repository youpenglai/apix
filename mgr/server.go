package mgr

import (
	"github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/middlewares"
	"io/ioutil"
	"encoding/json"
	"github.com/youpenglai/apix/gateway"
)

const defaultBindAddr = "127.0.0.1:58081"

type ServiceAddParam struct {
	Name string `json:"name"`
	BindAddr string `json:"bindAddr"`
}

type ServiceAddApiParam struct {
	ApiDocName string `json:"apiDocName"`
	ApiDocContent string `json:"apiDocContent"`
}

type ServiceCommandParam struct {
	Command string `json:"command,omitempty"`
}

func readJSON(ctx *http.Context, acceptor interface{}) error {
	defer ctx.Body().Close()
	reqBody, err := ioutil.ReadAll(ctx.Body())
	if err != nil {
		return err
	}

	return json.Unmarshal(reqBody, acceptor)
}

func addApi(ctx *http.Context) {
	serviceName := ctx.Params().GetStringDefault("serviceName", "")

	var param ServiceAddApiParam

	if err := readJSON(ctx, &param); err != nil {
		// TODO:
		ctx.JSON(400, map[string]interface{}{"errCode": 400, "errMsg": err.Error()})
		return
	}

	gw, err := GetHttpService(serviceName)
	if err != nil {
		ctx.JSON(404, map[string]interface{}{"errCode": 404, "errMsg": err.Error()})
		return
	}

	err = gw.AddApiDoc(param.ApiDocName, []byte(param.ApiDocContent))
	if err == ErrHttpServiceNotExists {
		ctx.JSON(404, map[string]interface{}{"errCode": 404, "errMsg": err.Error()})
		return
	} else if err != nil {
		ctx.JSON(500, map[string]interface{}{"errCode": 500, "errMsg": err.Error()})
		return
	}

	ctx.NoContent()
}

func command(ctx *http.Context) {
	serviceName := ctx.Params().GetStringDefault("serviceName", "")
	var param ServiceCommandParam

	if err := readJSON(ctx, &param); err != nil {
		ctx.JSON(400, map[string]interface{}{"errCode": 400, "errMsg": err.Error()})
		return
	}

	gw, err := GetHttpService(serviceName)
	if err != nil {
		ctx.JSON(404, map[string]interface{}{"errCode": 404, "errMsg": err.Error()})
		return
	}

	switch param.Command {
	case "serve":
		gw.Serve()
	case "reload":
		gw.Reload()
	case "stop", "shutdown":
		gw.Shutdown()
	}
}

func addService(ctx *http.Context) {
	var param ServiceAddParam

	if err := readJSON(ctx, &param); err != nil {
		ctx.JSON(400, map[string]interface{}{"errCode": 400, "errMsg": err.Error()})
		return
	}

	var opts gateway.ApiGatewayOpts
	opts.BindAddr = param.BindAddr
	AddHttpService(param.Name, &opts)
	ctx.JSON(200, map[string]interface{}{"success": true})
}

func getServiceState(ctx *http.Context) {
	// TODO: add code here
}

func installHandles(x *http.ApiX) {
	x.Post("/services/:serviceName/apis", addApi)
	x.Post("/services/:serviceName/cmd", command)
	x.Get("/services/:serviceName/state", getServiceState)
	x.Post("/services", addService)
}

// 运行管理端服务
func RunManagerServer(bindAddr ...string) {
	mgrServer := http.NewApiX()

	mgrServer.Use(middlewares.Server())

	mgrServer.Get("/", func(ctx *http.Context) {
		ctx.WriteString(200, "ApiX manager")
	})
	installHandles(mgrServer)

	addr := defaultBindAddr
	if len(bindAddr) > 0 {
		if bindAddr[0] != "" {
			addr = bindAddr[0]
		}
	}

	mgrServer.Run(addr)
}
