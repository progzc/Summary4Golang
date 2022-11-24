package leaky_bucket

// 限流算法之漏桶算法

import (
	// uber实现的漏桶算法
	_ "go.uber.org/ratelimit"
	"sync"
	"sync/atomic"
	"time"
)

// LeakyBucketLimiter 漏桶限频器
type LeakyBucketLimiter struct {
	pushMutex     sync.Mutex    // 加入漏桶需要获得的锁
	acquiredMutex sync.Mutex    // 从桶漏出需要获得的锁
	capacity      int64         // 漏桶的容量
	used          int64         // 漏桶已经被使用的大小(即漏桶中水的容量)
	last          time.Time     // 上一次水滴漏出的时间
	interval      time.Duration // 水滴匀速漏出的时间间隔
}

// New 创建漏桶限频器
func New(capacity int64) *LeakyBucketLimiter {
	// 计算水滴漏的间隔时间
	interval := 1e9 / capacity

	return &LeakyBucketLimiter{
		capacity: capacity,
		used:     0,
		interval: time.Duration(interval) * time.Nanosecond,
	}
}

// Push 加入漏桶
func Push(limiter *LeakyBucketLimiter) bool {
	limiter.pushMutex.Lock()
	defer limiter.pushMutex.Unlock()

	// 漏桶已满
	if limiter.used == limiter.capacity {
		return false
	}

	// 加入漏桶
	atomic.AddInt64(&limiter.used, 1)
	return true
}

// Acquired 从漏桶漏出
func Acquired(limiter *LeakyBucketLimiter) {
	limiter.acquiredMutex.Lock()
	defer limiter.acquiredMutex.Unlock()

	// 从漏桶中漏出
	defer atomic.AddInt64(&limiter.used, -1)

	now := time.Now()

	// 第一个请求到来
	if limiter.last.IsZero() {
		limiter.last = now
		return
	}

	// 计算休眠时间
	sleepFor := limiter.interval - now.Sub(limiter.last)

	// 休眠
	if sleepFor > 0 {
		time.Sleep(sleepFor)
		limiter.last = now.Add(sleepFor)
	} else {
		limiter.last = now
	}
}
