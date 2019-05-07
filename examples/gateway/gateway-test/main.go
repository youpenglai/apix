package main

import (
	"github.com/youpenglai/apix/grpc"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
)

func wait()(c chan os.Signal) {
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return
}

func main() {
	// wait load completed
	grpc.LoadAllProxy()
	time.Sleep(20 * time.Second)
	data, err := grpc.CallService("my-service", "hello", nil)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}
	fmt.Println("Result:", string(data))
	<- wait()
	fmt.Println("Exit main.")

}