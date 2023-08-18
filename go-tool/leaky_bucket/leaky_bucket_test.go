package leaky_bucket

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// 漏桶限流：https://lailin.xyz/post/go-training-week6-4-leaky-bucket.html
// 原理：一个漏水的桶，当有突发流量来临的时候，会先到桶里面，桶下有一个洞，可以以固定的速率向外流水，
//		如果水的从桶中外溢了出来，那么这个请求就会被拒绝掉。
// 适用场景：控制数据注入到网络的速率，平滑网络上的突发流量。漏桶算法提供了一种机制，
//	       通过它，突发流量可以被整形以便为网络提供一个稳定的流量。
// 实现： go.uber.org/ratelimit

// Test_IP_LeakyLimiter 实现基于IP的gin限流中间件（漏桶算法）
func Test_IP_LeakyLimiter(t *testing.T) {
	e := gin.Default()
	// 新建一个限速器，允许突发 3 个并发
	e.Use(NewLeakyLimiter(3))
	e.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	e.Run(":8080")
}

func NewLeakyLimiter(rps int) gin.HandlerFunc {
	limiters := &sync.Map{}
	return func(c *gin.Context) {
		// 获取限速器
		// key 除了 ip 之外也可以是其他的，例如 header，user name 等
		key := c.ClientIP()
		l, _ := limiters.LoadOrStore(key, ratelimit.New(rps))
		now := l.(ratelimit.Limiter).Take()
		fmt.Printf("now: %s\n", now)
		c.Next()
	}
}