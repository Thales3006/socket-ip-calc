package server

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net"
	"socket-ip-calc/internal/utils"
	"strings"
	"strconv"
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

			ip, err := utils.EvalIp(parts[1])
			isIpv4 := len(ip) == 4

			if err != nil {
				fmt.Println("Error (ip): ", err)
				conn.Write([]byte("invalid ip"))
				continue
			}

			fmt.Printf("ip %s is valid\n", parts[1])


			var mask int
			mask, err = strconv.Atoi(parts[2])

			if err != nil {
				fmt.Println("Error (mask): ", err)
				conn.Write([]byte("invalid mask"))
				continue
			}

			if !utils.ValidateMask(mask, isIpv4) {
				fmt.Println("Error: invalid mask")
				conn.Write([]byte("invalid mask"))
				continue
			}

			fmt.Printf("mask %s is valid\n", parts[2])

			var amount int
			amount, err = strconv.Atoi(parts[3])

			if err != nil {
				fmt.Println("Error (mask): ", err)
				conn.Write([]byte("invalid amount"))
				continue
			}

			if !utils.ValidateAmount(amount, mask, isIpv4) {
				fmt.Println("Error (amount): invalid amount")
				conn.Write([]byte("invalid amount"))
				continue
			}

			fmt.Printf("amount %s is valid\n", parts[3])

			var networks [][3][]uint16
			networks, err = Calculate(ip, mask, amount)
		 	
			if err != nil {
				fmt.Println("Error (calculation): ", err)
				conn.Write([]byte("invalid calculation"))
				continue
			}

			result := ";" + strconv.FormatInt(int64(mask) + int64(math.Log2(float64(amount))), 10)
			var ipStr string

			for _, subnetwork := range networks {

				for i:=0; i<3; i++ {
					ipStr, err = utils.StringifyIp(subnetwork[i])
					if err != nil {
						break
					}
					result += ";" + ipStr
				}

				if err != nil {
					break
				}
			}

			if err != nil {
				fmt.Println("Error (string calculation): ", err)
				conn.Write([]byte("invalid calculation"))
				continue
			}

			conn.Write([]byte("CALC" + result + "\n"))
		default:
			conn.Write([]byte("INVALID\n"))
		}
	}

}
