package main

import (
	"fmt"
	"github.com/youpenglai/apix/proxy"
	"os"
	"os/signal"
	"syscall"
)

func wait()(c chan os.Signal) {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return
}

func main() {
	proxyInst := proxy.InitServiceProxy()
	proxy.RegisterService(proxyInst, "my-service")

	proxy.HandleServiceCall(proxyInst, func(call *proxy.ProxyServiceCall) (data []byte, err error) {
		//fmt.Println("Call Service:", call.ServiceName)
		//fmt.Println("Call method:", call.Method)
		t := `{"name":"wang", "age": 9}`
		//data, err = json.Marshal(map[string]interface{}{"success": true, "data": })
		data = []byte(t)
		return
	})

	<- wait()
	fmt.Println("Exit.")
}
