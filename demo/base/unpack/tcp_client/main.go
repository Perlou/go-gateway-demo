package main

import (
	"fmt"
	"net"

	"github.com/perlou/go-gateway-demo/demo/base/unpack/unpack"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	unpack.Encode(conn, "hello world 0!!!")
}
