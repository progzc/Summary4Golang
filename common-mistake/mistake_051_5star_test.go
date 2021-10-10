package common_mistake

import (
	"testing"
)

// 阻塞的 gorutinue 与资源泄露
func TestMistake_051(t *testing.T) {
	wrong051()
	right051_1()
}

func wrong051() {
	type Result string
	type Search func(query string) Result
	first := func(query string, replicas []Search) Result {
		c := make(chan Result)
		replicaSearch := func(i int) {
			c <- replicas[i](query)
		}
		for i := range replicas {
			go replicaSearch(i)
		}
		return <-c
	}
	_ = first
}

func right051_1() {
	type Result string
	type Search func(query string) Result
	first := func(query string, replicas []Search) Result {
		// a. 使用带缓冲的 channel，确保能接收全部 goroutine 的返回结果
		c := make(chan Result, len(replicas))
		searchReplica := func(i int) {
			c <- replicas[i](query)
		}
		for i := range replicas {
			go searchReplica(i)
		}
		return <-c
	}
	_ = first
}

func right051_2() {
	type Result string
	type Search func(query string) Result
	first := func(query string, replicas []Search) Result {
		// b. 使用 select 语句，配合能保存一个缓冲值的 channel default 语句
		//    default 的缓冲 channel 保证了即使结果 channel 收不到数据，也不会阻塞 goroutine
		c := make(chan Result, 1)
		searchReplica := func(i int) {
			select {
			case c <- replicas[i](query):
			default:
			}
		}
		for i := range replicas {
			go searchReplica(i)
		}
		return <-c
	}
	_ = first
}

func right051_3() {
	type Result string
	type Search func(query string) Result
	first := func(query string, replicas []Search) Result {
		// c. 使用特殊的废弃（cancellation） channel 来中断剩余 goroutine 的执行：
		c := make(chan Result)
		done := make(chan struct{})
		defer close(done)
		searchReplica := func(i int) {
			select {
			case c <- replicas[i](query):
			case <-done:
			}
		}
		for i := range replicas {
			go searchReplica(i)
		}
		return <-c
	}
	_ = first
}
