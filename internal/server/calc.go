package server

import (
	"fmt"
	"errors"
	"math"
)


func Calculate(ip []uint16, mask int, amount int) (ips [][3][]uint16, err error) {
	var elSize int
	if len(ip) == 4 {
		elSize = 8
	} else if len(ip) == 8 {
		elSize = 16
	} else {
		err = errors.New("invalid ip")
		return
	}

	firstIndexHost := mask/elSize
	firstBitHost := elSize - 1 - (mask%elSize)

	if (ip[firstIndexHost]&((1<<(firstBitHost+1))-1)!=0) {
		err = errors.New("invalid host ip")
		return
	}

	for i:=firstIndexHost+1; i<len(ip); i++ {
		if ip[i]&((1<<elSize)-1) != 0 {
			err = errors.New("invalid host ip")
			return
		}
	}

	firstHost := make([]uint16, len(ip))
	copy(firstHost, ip)
	lastHost := make([]uint16, len(ip))
	copy(lastHost, ip)

	submask := mask + int(math.Log2(float64(amount)))

	firstIndexSubhost := submask/elSize
	firstBitSubhost := elSize-1-submask%elSize

	firstHost[len(firstHost)-1] |= 1

	lastHost[firstIndexSubhost] |= (1<<(firstBitSubhost+1))-1

	log.Printf("Debug: ip=%v, firstIndexHost=%d, firstBitHost=%d, firstIndexSubhost=%d, firstBitSubhost=%d", ip, firstIndexHost, firstBitHost, firstIndexSubhost, firstBitSubhost)
	log.Printf("Debug: ip[firstIndexHost]=%d, condition=%v", ip[firstIndexHost], ip[firstIndexHost]&((1<<(firstBitHost+1))-1) != 0)

	for i:=firstIndexSubhost+1; i<len(lastHost); i++ {
		lastHost[i] |= (1<<elSize) - 1
	}

	ips = make([][3][]uint16, amount)
	
	for i:=0; i<amount; i++ {
		for j:=0; j<3; j++ {
			ips[i][j] = make([]uint16, len(ip))
		}

		copy(ips[i][0], ip)
		copy(ips[i][1], firstHost)
		copy(ips[i][2], lastHost)

		n := i

		temp := uint16((n<<(firstBitSubhost+1))&((1<<elSize)-1))

		ips[i][0][firstIndexSubhost] |= temp
		ips[i][1][firstIndexSubhost] |= temp
		ips[i][2][firstIndexSubhost] |= temp
		n >>= elSize - firstBitSubhost - 1

		k := firstIndexSubhost-1
		for n>0 && k>firstIndexHost {
			temp = uint16(n&((1<<elSize)-1))
			ips[i][0][k] |= temp
			ips[i][1][k] |= temp
			ips[i][2][k] |= temp
			n >>= elSize
			k--
		}

		if n>0 {
			temp = uint16(n)
			ips[i][0][k] |= temp
			ips[i][1][k] |= temp
			ips[i][2][k] |= temp
		}
	}
	
	return
}
