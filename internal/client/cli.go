package client

import (
	"errors"
	"fmt"
	"regexp"
)

func handleInput() (string, error) {
	ip_regex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	number_regex := regexp.MustCompile(`^\d+$`)
	var ip string
	var mask string
	var amount string

	fmt.Print("IP: ")
	fmt.Scan(&ip)
	if !ip_regex.MatchString(ip) {
		return "", errors.New("not able to parse IP")
	}

	fmt.Print("Mask: ")
	fmt.Scan(&mask)
	if !number_regex.MatchString(mask) {
		return "", errors.New("not able to parse mask")
	}

	fmt.Print("Sub-network amount:")
	fmt.Scan(&amount)
	if !number_regex.MatchString(amount) {
		return "", errors.New("not able to parse amount")
	}

	return ip + "\n" + mask + "\n" + amount, nil
}
