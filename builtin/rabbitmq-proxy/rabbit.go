package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/youpenglai/apix/proxy"
)

const (
	RABBITMQ_SERVICE_NAME = "rabbitmq"

	DEFAULT_RABBITMQ_SERVICE_ADDR = "127.0.0.1"
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
	err := proxy.RegisterService(proxyInst, RABBITMQ_SERVICE_NAME)
	if err != nil {
		panic(err)
	}

	proxy.HandleServiceCall(proxyInst, func(call *proxy.ProxyServiceCall) (ret []byte, err error) {
		return
	})

	<-wait()
}
