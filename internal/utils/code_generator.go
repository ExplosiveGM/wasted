package utils

import "crypto/rand"

func GenerateCode(min, max int) int {
	var b [8]byte
	rand.Read(b[:])
	n := int(b[0])<<56 | int(b[1])<<48 | int(b[2])<<40 | int(b[3])<<32 |
		int(b[4])<<24 | int(b[5])<<16 | int(b[6])<<8 | int(b[7])
	if n < 0 {
		n = -n
	}
	return min + n%(max-min+1)
}
