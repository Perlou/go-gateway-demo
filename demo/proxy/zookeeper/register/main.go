package main

import (
	"fmt"
	"github.com/perlou/go-gateway-demo/proxy/zookeeper"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"127.0.0.1:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()
	i := 0
	for {
		zkManager.RegistServerPath("/real_server", fmt.Sprint(i))
		fmt.Println("Register", i)
		time.Sleep(5 * time.Second)
		i++
	}
}