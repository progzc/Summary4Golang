package chapter48_49_pprof

import (
	"log"
	"net/http"
	"testing"

	//a.在import中导入net/http/pprof
	_ "net/http/pprof"
)

// TestPProf_5
// (7)Q:如何为基于 HTTP 协议的网络服务添加性能分析接口？
//	  A:
//		a.在import中导入net/http/pprof
//		b.启动网络服务并开始监听
//		c.在网络浏览器中访问http://localhost:8082/debug/pprof

//		  在/debug/pprof/这个 URL 路径下还有很多可用的子路径，这一点你通过点选网页中的链接就可以了解到。
//		  像allocs、block、goroutine、heap、mutex、threadcreate这 6 个子路径，在底层其实都是通过Lookup函数来处理的。

//		  这些子路径都可以接受查询参数debug。它用于控制概要信息的格式和详细程度。至于它的可选值，我就不再赘述了。它的缺省值是0。
//		  另外，还有一个名叫gc的查询参数。它用于控制是否在获取概要信息之前强制地执行一次垃圾回收。只要它的值大于0，程序就会这样做。
//		  不过，这个参数仅在/debug/pprof/heap路径下有效。

//		  一旦/debug/pprof/profile路径被访问，程序就会去执行对 CPU 概要信息的采样。它接受一个名为seconds的查询参数。
//		  该参数的含义是，采样工作需要持续多少秒。如果这个参数未被显式地指定，那么采样工作会持续30秒。

// d.可以通过go tool pprof工具直接读取这样的 HTTP 响应,例如：
//
//	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=60
//
// e./debug/pprof/trace：在这个路径下，程序主要会利用runtime/trace代码包中的 API 来处理我们的请求。
//
//	更具体地说，程序会先调用trace.Start函数，然后在查询参数seconds指定的持续时间之后再调用trace.Stop函数。这里的seconds的缺省值是1秒。
func TestPProf_5(t *testing.T) {
	//b.启动网络服务并开始监听
	log.Println(http.ListenAndServe("localhost:8082", nil))
}
