package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"socket-ip-calc/internal/utils"
	"strconv"

	"github.com/fatih/color"
)

var text_color = color.RGB(255, 255, 255)
var reader = bufio.NewReader(os.Stdin)
var err error

func handleCli(conn net.Conn) {
	current_status := NOT_LOGGED
	serverReader := bufio.NewReader(conn)

	for {
		switch current_status {
		case NOT_LOGGED:
			login, err := getLoginInput()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			conn.Write([]byte(login))

		case LOGGED:
			input, err := getCalcInput()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			conn.Write([]byte(input))

		default:
			showOutput("Exiting!")
			return
		}

		current_status, err = handleServerMessage(serverReader)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !getContinue() {
			current_status = EXIT
		}
	}
}

func getCalcInput() (string, error) {
	text_color.Print("IP: ")
	ip, _ := utils.Read(reader)
	ipArr, err := utils.EvalIp(ip)
	if err != nil {
		return "", errors.New("not able to parse IP")
	}
	
	isIpv4 := len(ipArr)==4

	if isIpv4 {
		text_color.Print("IPV4\nMask (16 - 29): ")
	} else {
		text_color.Print("IPV6\nMask (48 - 62): ")
	}

	mask, _ := utils.Read(reader)
	var maskInt int
	maskInt, err = strconv.Atoi(mask)
	if err != nil {
		return "", errors.New("not able to parse mask")
	}

	if !utils.ValidateMask(maskInt, isIpv4) {
		return "", errors.New("invalid mask")
	}

	if isIpv4 {
		text_color.Printf("Sub-network amount (1 - %d): ", 1<<(31-maskInt))
	} else {
		text_color.Printf("Sub-network amount: ")
	}

	amount, _ := utils.Read(reader)
	var amountInt int
	amountInt, err = strconv.Atoi(amount)
	if err!=nil {
		return "", errors.New("not able to parse amount")
	}
	if !utils.ValidateAmount(amountInt, maskInt, isIpv4) {
		return "", errors.New("invalid amount")
	}

	return "CALC;" + ip + ";" + mask + ";" + amount + "\n", nil
}

func getLoginInput() (string, error) {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{4,}$`)

	text_color.Print("Username: ")
	username, _ := utils.Read(reader)
	if !regex.MatchString(username) {
		return "", errors.New("invalid username")
	}

	text_color.Print("Password: ")
	password, _ := utils.Read(reader)
	if !regex.MatchString(password) {
		return "", errors.New("invalid password")
	}

	return "LOGIN;" + username + ";" + password + "\n", nil
}

func getContinue() bool {
	yes := regexp.MustCompile(`^(y|Y|yes|Yes)$`)
	no := regexp.MustCompile(`^(n|N|no|No)$`)

	for {
		text_color.Print("Wish to continue [Yes or No]? ")
		decision, _ := utils.Read(reader)
		if yes.MatchString(decision) {
			return true
		}
		if no.MatchString(decision) {
			return false
		}
		text_color.Println("Please, insert a valid option.")
	}
}

func showOutput(output string) {
	d := color.New(color.FgCyan, color.Bold)
	d.Println("\n" + output)
}
