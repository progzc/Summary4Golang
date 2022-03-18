# Summary4Golang

## 1.go-bible 《go语言圣经》

## 2.go-core-36 《go语言核心36讲》 极客时间

## 3.go-concurrency 《go并发编程实战》 极客时间

- 分片map的实现：github.com/orcaman/concurrent-map
- 排序map的实现：github.com/elliotchance/orderedmap
- 过期map的实现：github.com/zekroTJA/timedmap
- 红黑树实现map的key排序：treemap
- 欧长坤（Go语言原本）：github.com/golang-design/under-the-hood
- 《Concurrency-in-Go》：https://s1.phpcasts.org/Concurrency-in-Go_Tools-and-Techniques-for-Developers.pdf
- 池化技术syn.Pool：github.com/fatih/pool

## 4.编码规范
- Go Code Review Comments：https://github.com/golang/go/wiki/CodeReviewComments
- Uber Go Style Guide：https://github.com/uber-go/guide/blob/master/style.md
- CommonMistakes：https://github.com/golang/go/wiki/CommonMistakes
- Effective go：https://golang.google.cn/doc/effective_go
- 专有名词：https://github.com/golang/lint/blob/738671d3881b9731cc63024d5d88cf28db875626/lint.go#L770

## 5.经验总结
- golang开发中的常见错误：http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/

## 6.优质的库
- 网络:
  - https://github.com/valyala/fasthttp (性能强大)
  - https://github.com/guonaihong/gout (使用便捷)
  - https://github.com/panjf2000/gnet (小巧，高性能、非阻塞的事件驱动Go网络框架)
- 存储:
  - https://github.com/syndtr/goleveldb (小巧,采用了LSM实现的K-V存储器,融合了其他常见数据库的设计思路)
- 日志库:
  - https://github.com/sirupsen/logrus (小巧，可以看看其设计实现）
  - https://github.com/uber-go/zap (trpc包装了其作为日志库)
- 工具:
  - https://github.com/panjf2000/ants (协程池)
  - https://github.com/ouqiang/timewheel (时间轮)
- 常用客户端:
  - https://github.com/go-gorm/gorm (操作mysql)
  - https://github.com/gomodule/redigo (操作redis)
  - https://github.com/Shopify/sarama (操作kafka)
