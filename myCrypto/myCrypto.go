package myCrypto

import (
	"crypto/md5"
	"encoding/hex"
)

// Simple md5 Checksum as Placeholder for Oneway Hashing Function
func Encode(s []byte) string {
	byteSlice := md5.Sum(s)
	return hex.EncodeToString(byteSlice[:])
}
