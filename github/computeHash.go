package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ComputeHmac256 compute hash
func ComputeHmac256(message []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(message)
	return h.Sum(nil)
}

func DecodeHex(hexHash string) (decoded []byte, hexErr error) {
	// drop the first 7 characters because github prefixes the hash with sha256=
	// ideally code should check which hash is being used before stripping the prefix
	decoded, hexErr = hex.DecodeString(hexHash[7:])
	return
}
