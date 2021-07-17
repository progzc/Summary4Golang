package go_concurrency

import (
	"fmt"
	"github.com/petermattis/goid"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// --------------------------------------------------------
// TestMutex 误区一: Lock/Unlock 不是成对出现
func TestMutex(t *testing.T) {
	var mu sync.Mutex
	defer mu.Unlock()
	fmt.Println("hello world!")
}

type Counter struct {
	sync.Mutex
	Count int
}

// --------------------------------------------------------
// TestMutex2 误区二: Copy已使用的Mutex
func TestMutex2(t *testing.T) {
	var c Counter
	c.Lock()
	defer c.Unlock()
	c.Count++
	foo(c) // 复制锁
}

// 这里Counter的参数是通过复制的方式传入的
func foo(c Counter) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("in foo")
}

// --------------------------------------------------------
// TestMutex3 误区三: Mutex是可重入锁
func TestMutex3(t *testing.T) {
	l := &sync.Mutex{}
	foo2(l)
}

// GoID 获取goroutine id
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// --------------------------------------------------------
// TestMutex4 方案一: 通过获取goroutine id实现可重入锁
func TestMutex4(t *testing.T) {
	l := &RecursiveMutex{}
	foo2(l)
}

func foo2(l sync.Locker) {
	fmt.Println("in foo")
	l.Lock()
	bar2(l)
	l.Unlock()
}
func bar2(l sync.Locker) {
	l.Lock()
	fmt.Println("in bar")
	l.Unlock()
}

// RecursiveMutex 包装一个Mutex,实现可重入
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 这个goroutine 重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调用次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

// --------------------------------------------------------
// TestMutex5 方案二: 通过token实现可重入锁
func TestMutex5(t *testing.T) {
	l := &TokenRecursiveMutex{}
	foo3(l)
}

func foo3(l *TokenRecursiveMutex) {
	fmt.Println("in foo")
	token := rand.Int63()
	l.Lock(token)
	bar3(l, token)
	l.Unlock(token)
}
func bar3(l *TokenRecursiveMutex, token int64) {
	l.Lock(token)
	fmt.Println("in bar")
	l.Unlock(token)
}

// Token方式的递归锁
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// 请求锁，需要传入token
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token { //如果传入的token和持有锁的token一致，说明是递归调用
		m.recursion++
		return
	}
	m.Mutex.Lock() // 传入的token不一致，说明不是递归调用
	// 抢到锁之后记录这个token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--         // 当前持有这个锁的token释放锁
	if m.recursion != 0 { // 还没有回退到最初的递归调用
		return
	}
	atomic.StoreInt64(&m.token, 0) // 没有递归调用了，释放锁
	m.Mutex.Unlock()
}

// --------------------------------------------------------
// TestMutex6 误区四: 死锁
func TestMutex6(t *testing.T) {
	// 派出所证明
	var psCertificate sync.Mutex
	// 物业证明
	var propertyCertificate sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2) // 需要派出所和物业都处理

	// 派出所处理goroutine
	go func() {
		defer wg.Done() // 派出所处理完成

		psCertificate.Lock()
		defer psCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求物业的证明
		propertyCertificate.Lock()
		propertyCertificate.Unlock()
	}()

	// 物业处理goroutine
	go func() {
		defer wg.Done() // 物业处理完成

		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求派出所的证明
		psCertificate.Lock()
		psCertificate.Unlock()
	}()
	wg.Wait()
	fmt.Println("成功完成")
}
