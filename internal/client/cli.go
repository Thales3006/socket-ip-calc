package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"socket-ip-calc/internal/utils"

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
	ip_regex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	number_regex := regexp.MustCompile(`^\d+$`)

	text_color.Print("IP: ")
	ip, _ := utils.Read(reader)
	if !ip_regex.MatchString(ip) {
		return "", errors.New("not able to parse IP")
	}

	text_color.Print("Mask: ")
	mask, _ := utils.Read(reader)
	if !number_regex.MatchString(mask) {
		return "", errors.New("not able to parse mask")
	}

	text_color.Print("Sub-network amount: ")
	amount, _ := utils.Read(reader)
	if !number_regex.MatchString(amount) {
		return "", errors.New("not able to parse amount")
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
