package go_explore

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	randlist = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func BenchmarkString_1_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateRandomString1(16)
	}
}

func BenchmarkString_2_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateRandomString2(16)
	}
}

func BenchmarkString_3_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateRandomString3(16)
	}
}

func BenchmarkString_1_2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CreateRandomString1(16)
		}
	})
}

func BenchmarkString_2_2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CreateRandomString2(16)
		}
	})
}

func BenchmarkString_3_2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CreateRandomString3(16)
		}
	})
}

func CreateRandomString1(len int) string {
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		b[i] = randlist[rand.Intn(62)]
	}
	return string(b)
}

func CreateRandomString2(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i := 0; i < len; i++ {
		b[i] = randlist[b[i]%(62)]
	}
	return string(b)
}

func CreateRandomString3(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i := 0; i < len; i++ {
		b[i] = randlist[b[i]%(62)]
	}
	return *(*string)(unsafe.Pointer(&b))
}
