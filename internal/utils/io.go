package utils

import (
	"bufio"
	"strings"
)

func Read(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}
