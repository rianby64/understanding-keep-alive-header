package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		data := make([]byte, 1000)
		n, err := conn.Read(data)
		if err != nil {
			panic(err)
		}

		fmt.Println(n)
		fmt.Println(string(data[:n]))

		conn.Close()
	}
}
