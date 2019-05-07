package main

import (
	"github.com/youpenglai/apix/grpc"
	"fmt"
	"encoding/json"
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
	proxy := grpc.InitServiceProxy()
	grpc.RegisterService(proxy, "my-service")

	grpc.HandleServiceCall(proxy, func(call *grpc.ProxyServiceCall) (data []byte, err error) {
		//fmt.Println("Call Service:", call.ServiceName)
		//fmt.Println("Call method:", call.Method)
		data, err = json.Marshal(map[string]bool{"success": true})
		return
	})

	<- wait()
	fmt.Println("Exit.")
}
