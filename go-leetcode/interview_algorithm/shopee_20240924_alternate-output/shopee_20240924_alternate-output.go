package shopee_20240924_alternate_output

import (
	"fmt"
)

// AlternateOutput
// 面试题：2个协程交替输出foobar 10次
func AlternateOutput() {
	chs := make([]chan struct{}, 2)
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan struct{})
	}

	stop := make(chan struct{})
	times := 10
	go func() {
		for i := 0; i < times; i++ {
			select {
			case <-chs[0]:
				fmt.Printf("foo")
				chs[1] <- struct{}{}
			case <-stop:
				return
			}
		}
	}()
	go func() {
		for i := 0; i < times; i++ {
			select {
			case <-chs[1]:
				fmt.Printf("bar\n")
				if i == times-1 {
					close(stop)
					return
				}
				chs[0] <- struct{}{}
			}
		}
	}()
	chs[0] <- struct{}{}
	<-stop
}
