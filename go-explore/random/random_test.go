package random1

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	randlist = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func TestRandomString1(t *testing.T) {
	fmt.Println(CreateRandomString1(6))
}

func CreateRandomString1(len int) string {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		b[i] = randlist[rand.Intn(62)]
	}
	return string(b)
}
