package utils

import (
	"errors"
	// "fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var ipv4Pattern = regexp.MustCompile("^((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))$")
var ipv6NoAbbrevPattern = regexp.MustCompile("^([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4})$")
var ipv6AbbrevEndPattern = regexp.MustCompile("^([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4}))?)?)?)?)?)?:{0,2}$")
var ipv6AbbrevStartPattern = regexp.MustCompile("^:{1,2}([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4}))?)?)?)?)?)?$")
var ipv6AbbrevMidPattern = regexp.MustCompile("^([[:xdigit:]]{1,4}(?::[[:xdigit:]]{1,4}){0,5})::([[:xdigit:]]{1,4}(?::[[:xdigit:]]{1,4}){0,5})$")
var ipv6AbbreverPattern = regexp.MustCompile("(?:^(?:0{1,4}:){2,})|(?:(?::0{1,4}){2,}$)|(?:(?::0{1,4}){2,}:)")


func IsIpv4(str string) bool {
	return ipv4Pattern.MatchString(str)
}

func IsIpv6(str string) bool {
	ip, err := EvalIp(str)
	return err==nil && len(ip)==8
}

func ValidateIp(ip []uint16) bool {
	var limit uint16

	if len(ip) == 4 {
		limit = 0xFF
	} else if len(ip) == 8 {
		limit = 0xFFFF
	} else {
		return false
	}

	for _, v := range ip {
		if v>limit {
			return false
		}
	}

	return true
}

func ValidateMask(mask int, isIpv4 bool) bool {
	return (isIpv4 && 16 <= mask && mask <= 29) || (!isIpv4 && 48 <= mask && mask <= 62)
}

func ValidateAmount(amount int, mask int, isIpv4 bool) bool {
	l := math.Log2(float64(amount))
	return ValidateMask(mask, isIpv4) && 1 <= amount && (l == float64(int(l))) && ((isIpv4 && amount<=(1<<(31-mask))) || !isIpv4)
}

func StringifyIp(ipArr []uint16) (ipStr string, err error) {
	if !ValidateIp(ipArr) {
		err = errors.New("invalid ip")
		return
	}

	ipStr = ""
	isIpv4 := len(ipArr) == 4
	var numberSystem int
	var separator string

	if isIpv4 {
		numberSystem = 10
		separator = "."
	} else {
		numberSystem = 16
		separator = ":"
	}

	for _, v := range ipArr {
		ipStr += strconv.FormatInt(int64(v), numberSystem) + separator
	}

	ipStr, _ = strings.CutSuffix(ipStr, separator)

	abbrev := ipv6AbbreverPattern.FindString(ipStr)

	if abbrev != "" {
		ipStr = strings.Replace(ipStr, abbrev, "::", 1)
	}

	return
}

func EvalIp(ipStr string) (ipArr []uint16, err error) {
	var temp uint64
	if ipv4Pattern.MatchString(ipStr) {
		ipArr = make([]uint16, 4)
		matches := ipv4Pattern.FindStringSubmatch(ipStr)
		for i := 0; i < 4; i++ {
			temp, err = strconv.ParseUint(matches[i+1], 10, 8)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[i] = uint16(temp)
		}
	} else if ipv6NoAbbrevPattern.MatchString(ipStr) {
		ipArr = make([]uint16, 8)
		matches := ipv6NoAbbrevPattern.FindStringSubmatch(ipStr)
		for i := 0; i < 8; i++ {
			temp, err = strconv.ParseUint(matches[i+1], 16, 16)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[i] = uint16(temp)
		}
	} else if ipv6AbbrevStartPattern.MatchString(ipStr) {
		ipArr = make([]uint16, 8)
		matches := ipv6AbbrevStartPattern.FindStringSubmatch(ipStr)
		j := 1
		for i := 1; i < len(matches); i++ {
			if matches[len(matches)-i]=="" {
				continue
			}
			temp, err = strconv.ParseUint(matches[len(matches)-i], 16, 16)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[8-j] = uint16(temp)
			j++
		}
	} else if ipv6AbbrevEndPattern.MatchString(ipStr) {
		ipArr = make([]uint16, 8)
		matches := ipv6AbbrevEndPattern.FindStringSubmatch(ipStr)
		j := 0
		for i := 1; i < len(matches); i++ {
			if matches[i]=="" {
				continue
			}
			temp, err = strconv.ParseUint(matches[i], 16, 16)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[j] = uint16(temp)
			j++
		}
	} else if ipv6AbbrevMidPattern.MatchString(ipStr) {
		matches := ipv6AbbrevMidPattern.FindStringSubmatch(ipStr)
		ipArr, err = EvalIp("::" + matches[2])
		if err != nil {
			ipArr = nil
		} else {
			var tempIp []uint16
			tempIp, err = EvalIp(matches[1] + "::")
			if err != nil {
				ipArr = nil
			} else {
				for i := 0; i < 8; i++ {
					if tempIp[i]*ipArr[i] == 0 {
						ipArr[i] += tempIp[i]
					} else {
						err = errors.New("invalid ip")
						ipArr = nil
						break
					}
				}
			}
		}
	} else {
		err = errors.New("invalid ip")
	}
	return
}
