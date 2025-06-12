package server

import (
	"fmt"
	"net"
)

func Start(port string) {
	address := ":8080"
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
}
