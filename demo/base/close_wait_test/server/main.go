package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		go func(conn net.Conn) {
			for {
				var buf [128]byte
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Printf("read from connect failed, err: %v\n", err)
					break
				}
				str := string(buf[:n])
				fmt.Printf("receive from client, data: %v\n", str)

			}
		}(conn)
	}
}
