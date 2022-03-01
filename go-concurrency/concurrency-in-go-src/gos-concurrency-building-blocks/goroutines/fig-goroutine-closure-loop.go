package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, salutation) // <1>
		}()
	}
	wg.Wait()
}
