package token_bucket

// 限流算法之令牌桶算法

import (
	// uber实现的漏桶算法
	_ "github.com/juju/ratelimit"
	"sync"
	"sync/atomic"
	"time"
)

// TokenBucketLimiter 令牌桶限频器
type TokenBucketLimiter struct {
	pushMutex     sync.Mutex    // 加入等待队列需要获得的锁
	acquiredMutex sync.Mutex    // 从令牌桶取令牌需要获得的锁
	queueLen      int64         // 等待队列容量
	queueUsed     int64         // 队列已被使用的大小
	capacity      int64         // 令牌桶的容量
	avail         int64         // 桶中令牌的数量
	last          time.Time     // 上一次计算令牌数的时间
	interval      time.Duration // 放令牌的时间间隔
}

// New 创建令牌桶限频器
func New(capacity int64) *TokenBucketLimiter {
	// 计算放令牌的时间间隔
	interval := 1e9 / capacity

	return &TokenBucketLimiter{
		queueLen:  capacity,
		queueUsed: 0,
		capacity:  capacity,
		avail:     0,
		last:      time.Now(),
		interval:  time.Duration(interval) * time.Nanosecond,
	}
}

// Push 加入等待队列
func Push(limiter *TokenBucketLimiter) bool {
	limiter.pushMutex.Lock()
	defer limiter.pushMutex.Unlock()

	// 等待队列已满
	if limiter.queueUsed == limiter.queueLen {
		return false
	}

	// 加入等待队列
	atomic.AddInt64(&limiter.queueUsed, 1)
	return true
}

// Acquired 获取令牌
func Acquired(limiter *TokenBucketLimiter) {
	limiter.acquiredMutex.Lock()
	defer limiter.acquiredMutex.Unlock()

	// 离开等待队列
	defer atomic.AddInt64(&limiter.queueUsed, -1)

	// 桶里有令牌
	if limiter.avail > 0 {
		limiter.avail--
		return
	}

	// 上次计算令牌的时间到现在，匀速放了多少令牌
	now := time.Now()
	add := limiter.capacity * now.Sub(limiter.last).Nanoseconds() / 1e9
	// 桶满了令牌溢出，令牌数量最多就是桶的容量大小
	if add > limiter.capacity {
		add = limiter.capacity
	}

	// 计算到有放了令牌
	if add > 0 {
		limiter.last = now
		limiter.avail += add - 1
		return
	}
	// 还没有放令牌，等待一个时间间隔
	time.Sleep(limiter.last.Add(limiter.interval).Sub(now))
	limiter.last = limiter.last.Add(limiter.interval)
}
