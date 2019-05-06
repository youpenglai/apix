package gateway

import (
	"os"
	"runtime"
	"github.com/youpenglai/goutils/pathtool"
	"os/exec"
	"encoding/json"
	"errors"
)

var (
	serviceProxy map[string]*ProxyService

	ErrServiceProxyNotFound = errors.New("service proxy not found")
)

type proxyProcess struct {

}

func startProxyProcess(proxyExe string) {
	proxyProcessInst := exec.Command(proxyExe)
	reader, err := proxyProcessInst.StdoutPipe()
	if err != nil {
		// TODO: add code here
		return
	}
	writer, err := proxyProcessInst.StdinPipe()
	if err != nil {
		// TODO: add code here
		return
	}

	proxySvc := NewServiceProxy()
	proxySvc.Attach(reader, writer)
	proxySvc.OnCall(func(param []byte) (retData []byte, err error) {
		var registerMsg ProxyServiceRegMsg
		if err = json.Unmarshal(param, &registerMsg); err != nil {
			return
		}
		for _, svcName := range registerMsg.ServiceNames {
			serviceProxy[svcName] = proxySvc
		}
		return
	})

	if err = proxyProcessInst.Start(); err != nil {
		// TODO: add code here
	}
}

func loadAllProxy() {
	pwd, err := os.Getwd()
	if err != nil {
		// TODO: 错误处理
		return
	}
	suffix := "-proxy"
	if runtime.GOOS == "windows" {
		suffix += ".exe"
	}
	proxyExeList, err := pathtool.GetDirFilesForSuffixs(pwd, []string{suffix})
	for _, proxyExe := range proxyExeList {
		startProxyProcess(proxyExe)
	}
}

func CallService(serviceName, methodName string, params []byte) ([]byte, error) {
	call := &ProxyServiceCall{
		ServiceName:serviceName,
		Method:methodName,
		Params:params,
	}
	data, err := json.Marshal(call)
	if err != nil {
		return nil, err
	}

	serviceInst, exists := serviceProxy[serviceName]
	if !exists {
		return nil, ErrServiceProxyNotFound
	}

	return serviceInst.CallSync(data)
}

func init() {
	loadAllProxy()
}