package common_mistake

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	fmt.Println(fmt.Sprintf("%x", buf))
	fmt.Println(hex.EncodeToString(buf))
	fmt.Println(base64.StdEncoding.EncodeToString(buf))
}
