package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/youpenglai/apix/proxy"
)

func getRabbitMQConf() (addr, user, password, vhost string) {
	return
}

func wait() (c chan os.Signal) {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return
}

func main() {
	proxyInst := proxy.InitServiceProxy()
	proxy.RegisterService(proxyInst, "")

	proxy.HandleServiceCall(proxyInst, func(call *proxy.ProxyServiceCall) (ret []byte, err error) {
		return
	})

	<-wait()
}
