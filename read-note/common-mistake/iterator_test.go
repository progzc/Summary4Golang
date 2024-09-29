package common_mistake

import (
	"fmt"
	"sync"
	"testing"
)

// link: https://github.com/golang/go/wiki/CommonMistakes
// link: https://golang.google.cn/doc/faq#closures_and_goroutines
// TestIterator Using reference to loop iterator variable (wrong usage)
func TestIterator(t *testing.T) {
	var out []*int
	for i := 0; i < 3; i++ {
		out = append(out, &i)
	}
	fmt.Println("Values:", *out[0], *out[1], *out[2]) // Values: 3 3 3
	fmt.Println("Addresses:", out[0], out[1], out[2]) // Addresses: 0xc00000a358 0xc00000a358 0xc00000a358
}

// TestIterator2 Using reference to loop iterator variable (right usage)
func TestIterator2(t *testing.T) {
	var out []*int
	for i := 0; i < 3; i++ {
		i := i
		out = append(out, &i)
	}
	fmt.Println("Values:", *out[0], *out[1], *out[2]) // Values: 0 1 2
	fmt.Println("Addresses:", out[0], out[1], out[2]) // Addresses: 0xc00000a358 0xc00000a360 0xc00000a368
}

// TestIterator3 Using reference to loop iterator variable (wrong usage)
func TestIterator3(t *testing.T) {
	var out [][]int
	for _, i := range [][1]int{{1}, {2}, {3}} {
		out = append(out, i[:])
	}
	fmt.Println("Values:", out) // Values: [[3] [3] [3]]
}

// TestIterator4 Using reference to loop iterator variable (right usage)
func TestIterator4(t *testing.T) {
	var out [][]int
	for _, i := range [][1]int{{1}, {2}, {3}} {
		i := i
		out = append(out, i[:])
	}
	fmt.Println("Values:", out) // Values: [[1] [2] [3]]
}

// TestIterator5 Using goroutines on loop iterator variables (wrong usage)
func TestIterator5(t *testing.T) {
	var wg sync.WaitGroup
	nums := []int{1, 2, 3}
	wg.Add(len(nums))
	for _, val := range nums {
		go func() {
			defer wg.Done()
			fmt.Println(val)
		}()
	}
	wg.Wait()
}

// TestIterator6 Using goroutines on loop iterator variables (right usage)
func TestIterator6(t *testing.T) {
	var wg sync.WaitGroup
	nums := []int{1, 2, 3}
	wg.Add(len(nums))
	for _, val := range nums {
		go func(val interface{}) {
			defer wg.Done()
			fmt.Println(val)
		}(val)
	}
	wg.Wait()
}

// TestIterator7 Using goroutines on loop iterator variables (right usage)
func TestIterator7(t *testing.T) {
	var wg sync.WaitGroup
	nums := []int{1, 2, 3}
	wg.Add(len(nums))
	for _, val := range nums {
		val := val
		go func() {
			defer wg.Done()
			fmt.Println(val)
		}()
	}
	wg.Wait()
}

// TestIterator8 Some confusion may arise when using closures with concurrency (wrong usage)
func TestIterator8(t *testing.T) {
	done := make(chan bool)
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
	// wait for all goroutines to complete before exiting
	for _ = range values {
		<-done
	}
}

// TestIterator9 Some confusion may arise when using closures with concurrency (right usage)
func TestIterator9(t *testing.T) {
	done := make(chan bool)
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func(v string) {
			fmt.Println(v)
			done <- true
		}(v)
	}
	// wait for all goroutines to complete before exiting
	for _ = range values {
		<-done
	}
}

// TestIterator10 Some confusion may arise when using closures with concurrency (right usage)
func TestIterator10(t *testing.T) {
	done := make(chan bool)
	values := []string{"a", "b", "c"}
	for _, v := range values {
		v := v // create a new 'v'.
		go func() {
			fmt.Println(v)
			done <- true
		}()
	}
	// wait for all goroutines to complete before exiting
	for _ = range values {
		<-done
	}
}
