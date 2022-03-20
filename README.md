# Summary4Golang

## 0.Primer
- web resource
  - https://gobyexample.com/
  - https://docs.microsoft.com/zh-cn/learn/paths/go-first-steps/
  - https://go101.org/

## 1.Courses From Geekbang
  - 《go语言核心36讲》 [100%]
    - 总结：https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/theory-explore/go-core-36
  - 《go并发编程实战》 [100%]
    - 总结：https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/theory-explore/go-concurrency-combat
  - 《分布式协议与算法实战》 [30%]
  - 《趣谈网络协议》 [60%]
  - 《深入剖析Kubernetes》 [50%]
  - 《etcd实战课》 [20%]
  - 《Elasticsearch核心技术与实战》 [100%]
  - 《MySQL实战45讲》 [100%]
  - 《Redis核心技术与实战》 [80%]
  - 《Kafka核心技术与实战》 [80%]
  - 《RPC实战与核心原理》 [88%]

## 2.Code Standards
- Go Code Review Comments：
  - https://github.com/golang/go/wiki/CodeReviewComments
- Style Guide：
  - https://github.com/uber-go/guide/blob/master/style.md
  - https://github.com/Tencent/secguide
- CommonMistakes：
  - https://github.com/golang/go/wiki/CommonMistakes
  - http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
    - 总结：https://github.com/progzc/Summary4Golang/tree/main/common-mistake
- Effective go：
  - https://golang.google.cn/doc/effective_go
- proper noun：
  - https://github.com/golang/lint/blob/738671d3881b9731cc63024d5d88cf28db875626/lint.go#L770
- book:
  - https://s1.phpcasts.org/Concurrency-in-Go_Tools-and-Techniques-for-Developers.pdf
  - Concurrent in go：https://github.com/kat-co/concurrency-in-go-src
  - Concurrent pattern：https://github.com/lotusirous/go-concurrency-patterns

## 3.High-Performance
- 池化思想
- 多路IO复用（Reactor编程思想）
- 写时复制（减少持有锁的思想）
- 保证Crash Safe
- 延迟初始化的思想
- 原子包中的原子操作及原子值的使用
- 反射（会降低性能，但能提高程序可用性，原理上可实现泛型的功能）

## 4.High Quality Lib
- 网络:
  - https://github.com/valyala/fasthttp (性能强大)
  - https://github.com/guonaihong/gout (使用便捷)
  - https://github.com/panjf2000/gnet (小巧，高性能、非阻塞的事件驱动Go网络框架)
  - https://github.com/alberliu/gim (IM系统)
- 存储:
  - https://github.com/syndtr/goleveldb (小巧,采用了LSM实现的K-V存储器,融合了其他常见数据库的设计思路。使用到的思想：顺序写、Crash safe、布隆过滤器，版本控制，多编码压缩、缓存、快照、跳表、二分、归并、LSM)
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
  
## 5.Improve Code
- leetcode
  - 分类：https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/classify_algorithm
  - 顺序: https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/order_leetcode
- design mode
  - https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/design_mode
- related blog
  - https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/blog
  
## 6.Classic Arithmetic
详见：https://github.com/progzc/Summary4Golang/tree/main/go-leetcode/tool
- 时间轮
- 协程池
- LRU
- LFU
- Filter
- 一致性Hash
- 红黑树
- 四种限流算法：
  - 固定时间窗口计数
  - 滑动时间窗口计数
  - 令牌桶
  - 漏斗
- 跳表
- 共识算法
  - MIT6.824 Distributed Systems：https://pdos.csail.mit.edu/6.824/
    - 视频：https://www.youtube.com/channel/UC_7WrbZTCODu1o_kfUMq88g
    - 翻译：https://www.bilibili.com/video/av91748150
    - database（包括SQL优化器、执行器、向量化）：https://github.com/pingcap/awesome-database-learning
    - paper：https://tidb.io/archived/events//paper-reading/

## 7.About Work
- Git
  - https://marklodato.github.io/visual-git-guide/index-zh-cn.html?no-svg
