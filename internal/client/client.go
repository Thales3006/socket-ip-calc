package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"socket-ip-calc/internal/utils"
	"strings"
)

type status int32

const (
	EXIT status = iota
	NOT_LOGGED
	LOGGED
)

func Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	handleCli(conn)
}

func handleServerMessage(serverReader *bufio.Reader) (status, error) {
	message, err := utils.Read(serverReader)
	if err != nil {
		fmt.Println("Error:", err)
		return EXIT, err
	}

	parts := strings.Split(message, ";")

	switch parts[0] {
	case "LOGIN_OK":
		showOutput("Login was successful!")
		return LOGGED, nil
	case "LOGIN_FAIL":
		showOutput("Wrong username or password!")
		return NOT_LOGGED, nil
	case "LOGIN_FALSE":
		showOutput("You are still not logged in!")
		return NOT_LOGGED, nil
	case "CALC":
		result := "Resulting calculation: \n"
		for i:=2; i+2<len(parts); i+=3 {
			result += parts[i] + "/" + parts[1] + "\t" + parts[i+1] + " - " + parts[i+2] + "\n"
		}
		showOutput(result)
		return LOGGED, nil
	default:
		showOutput("Error!")
		return EXIT, errors.New("invalid protocol")
	}
}
