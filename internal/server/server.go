package server

import (
	"fmt"
	"net"
)

func Start(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server listening on", port)

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

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(buf[:n]))
	fmt.Fprintf(conn, "Echo - "+string(buf[:n]))
}
