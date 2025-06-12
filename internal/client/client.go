package client

import (
	"fmt"
	"net"
)

func Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	input, err := handleInput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Fprint(conn, input)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	handleOutput(string(buf[:n]))
}
