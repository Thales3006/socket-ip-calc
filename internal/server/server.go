package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"socket-ip-calc/internal/utils"
	"strings"
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
	fmt.Println("Client connected.")
	defer conn.Close()
	reader := bufio.NewReader(conn)
	isLoged := false

	for {
		message, err := utils.Read(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client has disconnected.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}
		fmt.Println("[Client]: " + message + ".")

		parts := strings.Split(message, ";")
		switch parts[0] {
		case "LOGIN":
			switch {
			case isLoged:
				conn.Write([]byte("LOGIN_OK\n"))
			case parts[1] == "Admin" && parts[2] == "12345":
				isLoged = true
				conn.Write([]byte("LOGIN_OK\n"))
			default:
				conn.Write([]byte("LOGIN_FAIL\n"))
			}
		case "CALC":
			if !isLoged {
				conn.Write([]byte("LOGIN_FALSE\n"))
				continue
			}
			conn.Write([]byte("CALC; calculation" + message + "\n"))
		default:
			conn.Write([]byte("INVALID\n"))
		}
	}

}
