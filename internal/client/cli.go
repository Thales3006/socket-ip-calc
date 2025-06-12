package client

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

func handleInput() (string, error) {
	ip_regex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	number_regex := regexp.MustCompile(`^\d+$`)
	var ip string
	var mask string
	var amount string

	color.RGB(255, 255, 255).Print("IP: ")
	fmt.Scan(&ip)
	if !ip_regex.MatchString(ip) {
		return "", errors.New("not able to parse IP")
	}

	color.RGB(255, 255, 255).Print("Mask: ")
	fmt.Scan(&mask)
	if !number_regex.MatchString(mask) {
		return "", errors.New("not able to parse mask")
	}

	color.RGB(255, 255, 255).Print("Sub-network amount: ")
	fmt.Scan(&amount)
	if !number_regex.MatchString(amount) {
		return "", errors.New("not able to parse amount")
	}

	return ip + "\n" + mask + "\n" + amount, nil
}

func handleOutput(output string) {
	d := color.New(color.FgCyan, color.Bold)
	d.Println("\n" + output)
}
