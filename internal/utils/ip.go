package utils

import (
	"errors"
	// "fmt"
	"regexp"
	"strconv"
)

var ipv4Pattern = regexp.MustCompile("^((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))\\.((?:\\d{1,2})|(?:[01]\\d{2})|(?:2[0-4]\\d)|(?:25[0-5]))$")
var ipv6NoAbbrevPattern = regexp.MustCompile("^([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4}):([[:xdigit:]]{1,4})$")
var ipv6AbbrevEndPattern = regexp.MustCompile("^([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4}))?)?)?)?)?)?:{0,2}$")
var ipv6AbbrevStartPattern = regexp.MustCompile("^:{1,2}([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4})(?::([[:xdigit:]]{1,4}))?)?)?)?)?)?$")
var ipv6AbbrevMidPattern = regexp.MustCompile("^([[:xdigit:]]{1,4}(?::[[:xdigit:]]{1,4}){0,5})::([[:xdigit:]]{1,4}(?::[[:xdigit:]]{1,4}){0,5})$")


func IsIpv4(str string) bool {
	return ipv4Pattern.MatchString(str)
}

func IsIpv6(str string) bool {
	ip, err := EvalIp(str)
	return err==nil && len(ip)==8
}

func EvalIp(ipStr string) (ipArr []uint32, err error) {
	var temp uint64
	if ipv4Pattern.MatchString(ipStr) {
		ipArr = make([]uint32, 4)
		matches := ipv4Pattern.FindStringSubmatch(ipStr)
		for i := 0; i < 4; i++ {
			temp, err = strconv.ParseUint(matches[i+1], 10, 8)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[i] = uint32(temp)
		}
	} else if ipv6NoAbbrevPattern.MatchString(ipStr) {
		ipArr = make([]uint32, 8)
		matches := ipv6NoAbbrevPattern.FindStringSubmatch(ipStr)
		for i := 0; i < 8; i++ {
			temp, err = strconv.ParseUint(matches[i+1], 16, 32)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[i] = uint32(temp)
		}
	} else if ipv6AbbrevStartPattern.MatchString(ipStr) {
		ipArr = make([]uint32, 8)
		matches := ipv6AbbrevStartPattern.FindStringSubmatch(ipStr)
		j := 1
		for i := 1; i < len(matches); i++ {
			if matches[len(matches)-i]=="" {
				continue
			}
			temp, err = strconv.ParseUint(matches[len(matches)-i], 16, 32)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[8-j] = uint32(temp)
			j++
		}
	} else if ipv6AbbrevEndPattern.MatchString(ipStr) {
		ipArr = make([]uint32, 8)
		matches := ipv6AbbrevEndPattern.FindStringSubmatch(ipStr)
		for i := 1; i < len(matches); i++ {
			if matches[i]=="" {
				continue
			}
			temp, err = strconv.ParseUint(matches[i], 16, 32)
			if err != nil {
				ipArr = nil
				break
			}
			ipArr[i-1] = uint32(temp)
		}
	} else if ipv6AbbrevMidPattern.MatchString(ipStr) {
		matches := ipv6AbbrevMidPattern.FindStringSubmatch(ipStr)
		ipArr, err = EvalIp("::" + matches[2])
		if err != nil {
			ipArr = nil
		} else {
			var tempIp []uint32
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
