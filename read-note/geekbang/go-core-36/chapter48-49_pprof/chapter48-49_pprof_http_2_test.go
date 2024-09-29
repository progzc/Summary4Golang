package chapter48_49_pprof

import (
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
	"testing"
)

// TestPProf_6
// (8)Q:如何为基于 HTTP 协议的网络服务定制性能分析接口?
//	  A: 改变URL路径，详情如下示例。
// (9)Q:runtime/trace代码包的功用是什么？
//	  A:
//		简单来说，这个代码包是用来帮助 Go 程序实现内部跟踪操作的。其中的程序实体可以帮助我们记录程序中各个 goroutine 的状态、
//		各种系统调用的状态，与 GC 有关的各种事件，以及内存相关和 CPU 相关的变化，等等。

//		通过它们生成的跟踪记录可以通过go tool trace命令来查看。更具体的说明可以参看runtime/trace代码包的文档。

// 有了runtime/trace代码包，我们就可以为 Go 程序加装上可以满足个性化需求的跟踪器了。
// Go 语言标准库中有的代码包正是通过使用该包实现了自身的功能，例如net/http/pprof包。
func TestPProf_6(t *testing.T) {

	mux := http.NewServeMux()
	pathPrefix := "/d/pprof/"
	mux.HandleFunc(pathPrefix,
		func(w http.ResponseWriter, r *http.Request) {
			name := strings.TrimPrefix(r.URL.Path, pathPrefix)
			if name != "" {
				pprof.Handler(name).ServeHTTP(w, r)
				return
			}
			pprof.Index(w, r)
		})
	mux.HandleFunc(pathPrefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(pathPrefix+"profile", pprof.Profile)
	mux.HandleFunc(pathPrefix+"symbol", pprof.Symbol)
	mux.HandleFunc(pathPrefix+"trace", pprof.Trace)

	server := http.Server{
		Addr:    "localhost:8083",
		Handler: mux,
	}
	log.Println(server.ListenAndServe())
}
