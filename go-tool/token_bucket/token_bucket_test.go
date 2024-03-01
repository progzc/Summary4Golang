package token_bucket

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 令牌桶限流：https://lailin.xyz/post/go-training-week6-2-token-bucket-1.html
// 原理：
//	a.我们以 r/s  的速度向桶内放置令牌，桶的容量为 b , 如果桶满了令牌将会丢弃。
//	b.当请求到达时，我们向桶内获取令牌，如果令牌足够，我们就通过转发请求。
//	c.如果桶内的令牌数量不够，那么这个请求会被缓存等待令牌足够时转发，或者是被直接丢弃掉。
// 适用场景：
//	a.令牌桶算法不仅可以限流还可以应对突发流量的情况
// 特殊情况：
//	a.如果桶的容量为 0，那么相当于禁止请求，因为所有的令牌都被丢弃了。
//	b.如果令牌放置速率为无穷大，那么相当于没有限制。
// 实现: golang.org/x/time/rate

// Test_IP_TokenLimiter 实现基于IP的gin限流中间件（令牌桶算法）
func Test_IP_TokenLimiter(t *testing.T) {
	e := gin.Default()
	// 新建一个限速器，允许突发 10 个并发，限速 3rps，超过 500ms 就不再等待
	e.Use(NewTokenLimiter(3, 10, 500*time.Millisecond))
	e.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	e.Run(":8080")
}

func NewTokenLimiter(r rate.Limit, b int, t time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		// 获取限速器
		// key 除了 ip 之外也可以是其他的，例如 header，user name 等
		key := c.ClientIP()
		l, _ := limiters.LoadOrStore(key, rate.NewLimiter(r, b))

		// 这里注意不要直接使用 gin 的 context 默认是没有超时时间的
		ctx, cancel := context.WithTimeout(c, t)
		defer cancel()

		if err := l.(*rate.Limiter).Wait(ctx); err != nil {
			// 这里先不处理日志了，如果返回错误就直接 429
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": err})
		}
		c.Next()
	}
}