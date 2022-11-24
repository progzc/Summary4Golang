package leaky_bucket

import (
	"fmt"
	"go.uber.org/ratelimit"
	"testing"
	"time"
)

func Test_leaky_bucket(t *testing.T) {
	rl := ratelimit.New(100)
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
