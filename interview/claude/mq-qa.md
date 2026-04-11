# 消息队列 Q&A

## 目录

- [一、消息队列如何选型](#一消息队列如何选型)
- [二、消息堆积怎么办](#二消息堆积怎么办)
- [三、如何保证消息的有序消费](#三如何保证消息的有序消费)
- [四、如何实现按优先级消费](#四如何实现按优先级消费)
- [五、如何保证消息不丢](#五如何保证消息不丢)
- [六、如何保证消息不重复消费（幂等性）](#六如何保证消息不重复消费幂等性)
- [七、MQ 如何保证分布式事务一致性](#七mq-如何保证分布式事务一致性)
- [八、Kafka 为什么吞吐量这么高](#八kafka-为什么吞吐量这么高)
- [九、消息队列如何实现延迟消息](#九消息队列如何实现延迟消息)
- [十、Kafka Rebalance 机制](#十kafka-rebalance-机制)
- [十一、消息队列的推模式 vs 拉模式](#十一消息队列的推模式-vs-拉模式)

---

## 一、消息队列如何选型

### 为什么需要消息队列

```
消息队列解决的三个核心问题:

1. 异步解耦:
   同步: 用户下单 → 扣库存 → 创建订单 → 发通知 → 加积分 (串行 500ms)
   异步: 用户下单 → 发消息 → 立即返回 (50ms)
         后台: 消费消息 → 扣库存 / 发通知 / 加积分 (并行处理)

2. 流量削峰:
   秒杀: 10w QPS 打到数据库 → 数据库崩
   加 MQ: 10w QPS → MQ 缓冲 → 消费者按 1000 QPS 匀速消费 → 数据库安全

3. 数据分发:
   订单数据需要给: 搜索系统、推荐系统、数据分析、风控系统
   没 MQ: 订单服务对接 4 个下游（强耦合）
   有 MQ: 订单服务发一条消息，4 个系统各自消费（解耦）
```

### 主流消息队列对比

| 维度 | **Kafka** | **RocketMQ** | **RabbitMQ** | **Pulsar** | **Redis Stream** |
|------|-----------|-------------|-------------|-----------|-----------------|
| **开发语言** | Scala/Java | Java | Erlang | Java | C |
| **定位** | 大数据流处理平台 | 业务消息中间件 | 企业级消息代理 | 云原生流平台 | 轻量级消息队列 |
| **吞吐量** | **百万级/s** | 十万级/s | 万级/s | 百万级/s | 十万级/s |
| **延迟** | ms 级 | ms 级 | **μs 级** | ms 级 | ms 级 |
| **消息可靠性** | 高（副本机制） | **极高**（同步刷盘+副本） | 高（镜像队列） | 高（BookKeeper） | 中（取决于持久化配置） |
| **消息顺序** | 分区内有序 | **队列内严格有序** | 不保证 | 分区内有序 | 有序 |
| **延迟消息** | 不原生支持 | **原生支持** | 插件支持（TTL+DLX） | 不原生支持 | 不支持 |
| **事务消息** | 支持（0.11+） | **原生支持（半消息）** | 支持（AMQP 事务） | 支持 | 不支持 |
| **消息回溯** | 支持（按 offset/时间） | 支持（按时间） | 不支持 | 支持 | 有限支持 |
| **死信队列** | 不原生 | **原生支持** | 原生支持（DLX） | 原生支持 | 不支持 |
| **消费模型** | Pull（消费者拉） | Pull + Push | Push（Broker 推） | Pull | Pull |
| **协议** | 自定义 | 自定义 | **AMQP**（标准协议） | 自定义 + 兼容 Kafka | Redis 协议 |
| **运维复杂度** | 高（依赖 ZK/KRaft） | 中（依赖 NameServer） | 低 | 高（依赖 ZK + BK） | **极低** |
| **社区生态** | **极强**（大数据标配） | 强（阿里开源） | 强（历史悠久） | 中等 | 极强（Redis 附属） |
| **典型用户** | LinkedIn, Uber, Netflix | 阿里, 字节, 美团 | 金融, 电商 | 腾讯, Yahoo | 中小项目 |

---

### Kafka 详解

```
架构:
  Producer → Broker Cluster → Consumer Group
                 │
           ┌─────┴──────┐
           │ Topic       │
           │ ├─ Partition 0 → [msg0, msg1, msg2, ...]  Leader
           │ │   Replica  → [msg0, msg1, msg2, ...]  Follower
           │ ├─ Partition 1 → [msg0, msg1, msg2, ...]
           │ └─ Partition 2 → [msg0, msg1, msg2, ...]
           └─────────────┘

核心特点:
  - 分区(Partition)是并行和有序的基本单位
  - 消息追加写入分区（顺序 I/O），性能极高
  - Consumer Group 内每个 Partition 只被一个消费者消费 → 分区内有序
  - 消息持久化到磁盘，支持按 offset/时间回溯
  - 消费者拉模式(Pull)，自己控制消费速度
```

**适用场景：** 大数据实时流处理、日志采集、事件溯源、CDC（数据变更捕获）、用户行为追踪

**不适合：** 需要延迟消息、复杂路由、低延迟（μs 级）的业务场景

---

### RocketMQ 详解

```
架构:
  Producer → NameServer(路由) → Broker Cluster → Consumer Group
                                    │
                              ┌─────┴──────┐
                              │ Topic       │
                              │ ├─ Queue 0  │
                              │ ├─ Queue 1  │
                              │ ├─ Queue 2  │
                              │ └─ Queue 3  │
                              └─────────────┘

核心特点:
  - 专为业务消息设计（阿里双十一验证）
  - 原生事务消息: 半消息 → 本地事务 → 提交/回滚
  - 原生延迟消息: 支持 18 个级别（1s/5s/10s/.../2h）
  - 原生死信队列、消息轨迹、消息过滤(Tag/SQL)
  - 严格顺序消息（同一 Queue 内）
```

**事务消息流程：**

```
  Producer                  Broker                    Consumer
     │                        │                          │
     │── 半消息(Half) ───────→│  暂存，消费者不可见         │
     │                        │                          │
     │  执行本地事务            │                          │
     │  (如: 扣款)             │                          │
     │                        │                          │
     │── Commit ────────────→│  消息可见                  │
     │   或 Rollback ────────→│  删除消息                  │
     │                        │                          │
     │  如果 Producer 超时:    │                          │
     │  Broker 主动回查         │                          │
     │←── 事务状态回查 ────────│                          │
     │── Commit/Rollback ──→│                          │
```

**适用场景：** 电商交易、金融支付、延迟任务、顺序消息、事务一致性

---

### RabbitMQ 详解

```
架构:
  Producer → Exchange → Binding → Queue → Consumer
                │
    ┌───────────┼────────────────────────┐
    │ Exchange Types:                     │
    │   Direct:  精确匹配 routing key     │
    │   Fanout:  广播到所有绑定的 Queue    │
    │   Topic:   通配符匹配 routing key   │
    │   Headers: 根据消息头匹配           │
    └────────────────────────────────────┘

核心特点:
  - 最灵活的路由模型（Exchange + Binding）
  - 遵循 AMQP 标准协议（跨语言互通）
  - 消息确认机制完善（Publisher Confirm + Consumer Ack）
  - 延迟低（μs 级）
  - Erlang 开发，天生高可用（热升级）
  - 吞吐量相对低（万级/s）
```

**适用场景：** 复杂路由、企业应用集成、RPC 通信、低延迟要求、中小规模

---

### Redis Stream 详解

```
架构:
  Producer → XADD stream_key → Consumer Group → XREADGROUP

核心特点:
  - Redis 5.0+ 内置，零部署成本
  - 支持消费组、消息确认(XACK)、Pending 列表
  - 数据存内存，性能高（但受内存限制）
  - 功能相对简单（无延迟消息、无事务、无死信队列）
  - 持久化靠 Redis 自身（RDB/AOF），可靠性不如专业 MQ
```

**适用场景：** 已有 Redis 的轻量场景、中小项目、对可靠性要求不极高的业务

---

### 选型决策树

```
需要消息队列:
│
├── 大数据/日志/事件流/高吞吐?
│   └── Kafka ✓ (百万级吞吐，大数据生态完善)
│
├── 电商/金融/交易类业务?
│   ├── 需要事务消息? → RocketMQ ✓
│   ├── 需要延迟消息? → RocketMQ ✓
│   ├── 需要严格顺序? → RocketMQ ✓
│   └── 以上都需要 → RocketMQ ✓
│
├── 复杂路由/企业集成/低延迟?
│   └── RabbitMQ ✓ (AMQP 标准, 路由灵活, μs 级延迟)
│
├── 云原生/存算分离/超大规模?
│   └── Pulsar ✓ (存算分离, 多租户, 无限存储)
│
├── 轻量/已有 Redis/功能简单?
│   └── Redis Stream ✓ (零部署成本)
│
└── 不确定?
    ├── 大多数互联网公司 → Kafka (日志/数据) + RocketMQ (业务消息)
    └── 中小团队 → RabbitMQ (功能全面, 运维简单)
```

### 按维度选型

#### 按吞吐量

```
百万级/s: Kafka > Pulsar >> RocketMQ > Redis Stream >> RabbitMQ

Kafka 吞吐为什么高？
  1. 顺序 I/O: 消息追加写磁盘，不是随机写
  2. Page Cache: 利用 OS 页缓存，读写都走内存
  3. 零拷贝(sendfile): 数据从磁盘到网卡，不经过用户空间
  4. 批量+压缩: 多条消息打包发送，减少网络 round trip
  5. 分区并行: 多个 Partition 并行读写
```

#### 按可靠性

```
极高: RocketMQ (同步刷盘+同步副本，阿里双十一验证)
  >
高:  Kafka (ISR 副本机制，acks=all)
  ≈  RabbitMQ (镜像队列 + Publisher Confirm)
  ≈  Pulsar (BookKeeper 多副本)
  >
中:  Redis Stream (依赖 Redis 持久化，宕机可能丢数据)
```

#### 按功能完整性

```
RocketMQ > RabbitMQ > Pulsar > Kafka > Redis Stream

RocketMQ "全家桶":
  ✓ 事务消息 (半消息+回查)
  ✓ 延迟消息 (18级/任意时间)
  ✓ 顺序消息 (队列级别严格有序)
  ✓ 死信队列 (消费失败自动进入)
  ✓ 消息过滤 (Tag + SQL92)
  ✓ 消息回溯 (按时间点回溯)
  ✓ 消息轨迹 (全链路追踪)
  ✓ 消费重试 (自动重试+退避)
```

#### 按运维复杂度

```
简单 ←─────────────────────────────────────→ 复杂

Redis Stream < RabbitMQ < RocketMQ < Kafka < Pulsar

Redis Stream: 已有 Redis 就能用
RabbitMQ:     单体部署 + Management 控制台
RocketMQ:     NameServer + Broker + Dashboard
Kafka:        ZooKeeper/KRaft + Broker + 各种 Connect
Pulsar:       ZooKeeper + BookKeeper + Broker
```

### 常见组合方案

```
方案一: Kafka + RocketMQ（大型互联网公司主流）
  Kafka  → 日志采集、用户行为、大数据管道、Flink 实时计算
  RocketMQ → 交易订单、支付通知、延迟任务、业务消息

方案二: Kafka 一把梭（数据密集型公司）
  所有场景都用 Kafka
  延迟消息/事务消息自己在应用层实现
  适合: 有强大基础设施团队的公司

方案三: RabbitMQ（中小型公司/传统企业）
  功能全面、运维简单、生态成熟
  吞吐量够用（万级 QPS 足够大多数业务）

方案四: Redis Stream（微服务轻量通信）
  已有 Redis，不想额外部署 MQ
  可靠性要求不高的内部通信
```

### 一句话总结

- **Kafka**：大数据/日志/高吞吐首选，百万级吞吐，大数据生态无敌
- **RocketMQ**：业务消息首选，事务消息+延迟消息+顺序消息一站式，阿里双十一验证
- **RabbitMQ**：中小规模/复杂路由/低延迟首选，AMQP 标准，运维简单
- **Pulsar**：云原生/存算分离/多租户场景
- **Redis Stream**：已有 Redis 的轻量场景，零部署成本
- **不确定？** 大多数公司 Kafka(数据流) + RocketMQ(业务消息) 是最稳的组合

---

## 二、消息堆积怎么办

### 什么是消息堆积

```
消息堆积 = 生产速度 > 消费速度，未消费的消息越积越多

正常:
  生产: 1000 msg/s → [MQ: 100条] → 消费: 1000 msg/s  (进出平衡)

堆积:
  生产: 1000 msg/s → [MQ: 500万条!!!] → 消费: 200 msg/s  (消费跟不上)

后果:
  - 消息处理延迟越来越大（从 ms 级涨到 小时/天级）
  - MQ 磁盘/内存被撑满 → 新消息写入失败 → 生产者阻塞
  - 上游服务超时、业务功能异常
  - 极端情况: 消息过期被丢弃 → 数据丢失
```

### 堆积原因分析

```
┌─────────────────────────────────────────────────────────┐
│  原因一: 消费者处理太慢                                    │
│    - 消费逻辑中有慢查询（数据库/外部 API）                  │
│    - 消费逻辑有 bug（死循环、锁竞争、内存泄漏）             │
│    - 单条消息处理耗时增加（如数据量增大）                    │
├─────────────────────────────────────────────────────────┤
│  原因二: 消费者数量不足                                    │
│    - 消费者实例太少，并行度不够                             │
│    - 分区/队列数少于消费者数 → 部分消费者空闲               │
│    - 消费者挂了没有自动恢复                                │
├─────────────────────────────────────────────────────────┤
│  原因三: 生产者突增                                       │
│    - 上游业务流量暴涨（大促、活动）                         │
│    - 批量数据导入/全量同步                                 │
│    - 上游 bug 重复发送                                    │
├─────────────────────────────────────────────────────────┤
│  原因四: MQ 自身问题                                      │
│    - Broker 节点故障，消费者 rebalance 中断消费             │
│    - 消费者提交 offset 失败 → 重复消费 → 处理变慢          │
│    - 网络问题导致消费者频繁断连                             │
└─────────────────────────────────────────────────────────┘
```

### 紧急处理（线上救火）

#### 第一步：快速定位原因

```
1. 看监控面板:
   - 消费者 lag（积压量）曲线
   - 消费者 TPS（每秒处理量）
   - 消费者实例数量和状态
   - 单条消息处理耗时

2. 判断是哪类问题:
   消费者 TPS 正常但 lag 增长 → 生产者突增（量的问题）
   消费者 TPS 下降 → 消费者变慢/挂了（消费侧问题）
   消费者 TPS 为 0 → 消费者全挂了/MQ 故障
```

#### 第二步：紧急扩容消费者（最快见效）

```
Kafka 场景:

  堆积前: Topic 有 16 个 Partition，4 个消费者实例
    P0~P3  → Consumer 1
    P4~P7  → Consumer 2
    P8~P11 → Consumer 3
    P12~P15→ Consumer 4

  紧急扩容: 增加到 16 个消费者（每个 Partition 一个消费者，最大并行度）
    P0 → Consumer 1
    P1 → Consumer 2
    ...
    P15→ Consumer 16

  注意: 消费者数 > Partition 数时，多出的消费者空闲
  → 如果 Partition 不够 → 需要扩 Partition（但这个操作有风险，不建议紧急做）

RocketMQ 场景:
  增加消费者实例 + 增大每个消费者的线程数:
    consumer.setConsumeThreadMin(20)
    consumer.setConsumeThreadMax(64)
```

#### 第三步：临时跳过非核心逻辑

```go
// 堆积时临时降级消费逻辑，只做核心操作
func consume(msg Message) error {
    // 核心: 更新订单状态（必须做）
    if err := updateOrderStatus(msg); err != nil {
        return err
    }

    // 非核心: 发通知、更新搜索索引、统计（暂时跳过）
    if !isBackpressure() { // 检查是否在堆积状态
        sendNotification(msg)
        updateSearchIndex(msg)
        updateStatistics(msg)
    }
    return nil
}

func isBackpressure() bool {
    // 通过配置中心动态开关控制
    return config.GetBool("mq.backpressure.enabled")
}
```

#### 第四步：转发到临时队列（终极方案）

```
当堆积量极大（千万级），正常消费来不及:

  原队列（500 万堆积）                 临时队列（扩容 10 倍容量）
  ┌──────────────┐    快速转发      ┌─────────────────────┐
  │ order_topic  │──────────────→│ order_topic_tmp_0   │→ 消费者组1
  │ (16 分区)    │   不做业务处理   │ order_topic_tmp_1   │→ 消费者组2
  │              │   只是搬运      │ ...                 │→ ...
  └──────────────┘                │ order_topic_tmp_9   │→ 消费者组10
                                  └─────────────────────┘

步骤:
  1. 创建临时 Topic（更多分区，如 160 个）
  2. 写一个"转发消费者"：从原队列读 → 写到临时 Topic（不做业务处理，极快）
  3. 部署大量消费者实例消费临时 Topic → 执行真正的业务逻辑
  4. 堆积消化完后 → 恢复正常消费 → 删除临时 Topic
```

### 长期预防方案

#### 1. 消费者优化

```
提升单条消息处理速度:
  - 优化慢查询（加索引、缓存热数据）
  - 批量处理代替逐条处理:
    每次拉取 100 条 → 批量 INSERT → 比逐条 INSERT 快 10 倍
  - 异步化: 消费后先落库确认 → 非核心逻辑异步处理
  - 并行消费: 单个消费者内部用 goroutine pool 并行处理

  // Go 并行消费示例
  func consume(msgs []Message) {
      var wg sync.WaitGroup
      sem := make(chan struct{}, 32) // 限制并发 32
      for _, msg := range msgs {
          wg.Add(1)
          sem <- struct{}{}
          go func(m Message) {
              defer wg.Done()
              defer func() { <-sem }()
              processMessage(m)
          }(msg)
      }
      wg.Wait()
  }
```

#### 2. 合理设置分区/队列数

```
Kafka:
  分区数 = 预期峰值 QPS / 单消费者处理能力
  例: 峰值 10000 QPS，单消费者 1000 QPS → 至少 10 分区（建议 16，留余量）

  分区太少 → 扩容消费者也没用（消费者多于分区时空闲）
  分区太多 → 选举开销大、文件描述符多、端到端延迟增加
  经验值: 单 Topic 16~64 分区覆盖大多数场景
```

#### 3. 监控告警

```
必须监控的指标:
  - Consumer Lag（堆积量）→ 超过阈值告警
  - 消费 TPS → 下降告警
  - 消费者实例数 → 减少告警
  - 单条消息处理 P99 耗时 → 增加告警
  - 消费者错误率 → 上升告警

告警阈值示例:
  Lag > 10000   → Warning
  Lag > 100000  → Critical
  消费 TPS = 0  → Critical（消费者全挂）
  P99 > 5s      → Warning
```

#### 4. 消费者弹性伸缩

```
基于堆积量自动扩缩容（K8s HPA + 自定义指标）:

  if lag > 50000:
      scale_up(消费者副本数 × 2)  # 堆积多了自动扩容
  if lag < 1000 && 持续 10 分钟:
      scale_down(消费者副本数 / 2)  # 堆积消化了自动缩容
```

#### 5. 生产者端限流

```
源头控制: 不让生产者无限制发消息

  - 令牌桶限流生产速率
  - 大批量数据导入用限速（如每秒最多 5000 条）
  - 上游服务降级: 非核心消息暂停发送
```

### 不同 MQ 的堆积特点

| MQ | 堆积能力 | 堆积对性能影响 | 特殊注意 |
|----|---------|-------------|---------|
| **Kafka** | 极强（磁盘存储，TB 级） | 几乎无影响（顺序读磁盘） | 堆积不可怕，慢慢消费即可 |
| **RocketMQ** | 强（磁盘存储，百 GB） | 堆积多时读性能略降 | 消息默认保留 48h，超时会被清理 |
| **RabbitMQ** | 弱（内存为主） | **堆积会严重影响性能** | 内存满 → 触发 flow control → 生产者阻塞 |
| **Redis Stream** | 弱（内存） | 内存不足则 OOM | 需设置 MAXLEN 限制长度 |

```
Kafka 为什么堆积不影响性能?
  → 消息持久化到磁盘（顺序写）
  → 消费者读取时走 Page Cache（如果能命中）或磁盘顺序读
  → 百万级堆积和零堆积的消费速度几乎一样
  → 所以 Kafka 堆积通常不需要紧急处理，扩消费者慢慢追即可

RabbitMQ 为什么堆积很危险?
  → 消息主要在内存中
  → 堆积多 → 内存飙升 → 触发 flow control（限流生产者）
  → 极端情况: 内存耗尽 → Broker 崩溃
  → RabbitMQ 堆积时必须立即处理!
```

### 处理流程总结

```
发现堆积 → 判断原因:
│
├── 消费者挂了/变慢
│   ├── 紧急: 重启消费者 / 修复 bug / 扩容实例
│   └── 长期: 优化消费逻辑 + 监控告警
│
├── 生产者突增
│   ├── 紧急: 扩容消费者到最大并行度
│   ├── 不够: 转发到临时队列 + 大量临时消费者
│   └── 长期: 生产端限流 + 弹性伸缩
│
├── 分区不够
│   ├── 紧急: 临时队列方案
│   └── 长期: 提前规划足够分区数
│
└── MQ 自身故障
    ├── 紧急: 故障恢复 + 消费者重连
    └── 长期: MQ 集群高可用 + 监控
```

### 一句话总结

**消息堆积的本质是消费速度 < 生产速度。紧急处理四步走：定位原因（看 lag+TPS 曲线）→ 扩容消费者到最大并行度 → 降级非核心逻辑 → 极端情况用临时队列方案（转发+大量消费者消化堆积）。** 长期预防：优化消费速度（批量+并行+异步）、合理分区数、监控 lag 告警、弹性伸缩、生产端限流。Kafka 堆积不可怕（磁盘存储不影响性能，慢慢追即可），RabbitMQ 堆积很危险（内存为主，堆积会导致 Broker 崩溃必须立即处理）。

---

## 三、如何保证消息的有序消费

### 为什么消息会乱序

```
消息乱序发生在三个环节:

  生产者         MQ Broker         消费者
    │               │                │
    │  ① 发送乱序    │  ② 存储乱序     │  ③ 消费乱序
    │               │                │

① 生产者发送乱序:
   多线程/多实例同时发送 → 网络延迟不同 → 到达 Broker 顺序不等于发送顺序

② MQ 存储乱序:
   消息被分发到不同的分区/队列 → 各分区独立存储 → 全局无序

③ 消费者消费乱序:
   多个消费者并行消费 → 处理速度不同 → 先收到的不一定先处理完

典型问题:
  订单状态变更: 创建(msg1) → 支付(msg2) → 发货(msg3)
  如果乱序消费: 收到 发货 → 支付 → 创建
  → 发货时订单还不存在 → 业务异常!
```

### 两种有序性

```
全局有序: 所有消息严格按发送顺序消费
  → 只能用一个分区 + 一个消费者 → 吞吐极低
  → 绝大多数场景不需要

分区有序（部分有序）: 同一业务 key 的消息有序，不同 key 之间不要求有序
  → 实际需求: 同一个订单的消息有序，不同订单之间无需有序
  → 这才是生产中真正需要的"有序"
```

### 各 MQ 的有序方案

#### Kafka：分区内有序

```
Kafka 天然保证: 同一 Partition 内的消息严格有序

  关键: 让同一业务 key 的消息发到同一个 Partition

生产者:
  producer.Send(&kafka.Message{
      Topic:     "order_events",
      Key:       []byte(orderID),     // 用 orderID 做 key
      Value:     msgBody,
  })
  // Kafka 默认按 key 的 hash 决定 Partition:
  // partition = hash(key) % numPartitions
  // 同一个 orderID → 同一个 Partition → 有序 ✓

  订单 A 的消息 → Partition 0: [创建A, 支付A, 发货A]  ← 有序 ✓
  订单 B 的消息 → Partition 1: [创建B, 支付B]          ← 有序 ✓
  A 和 B 之间 → 无序（不需要有序）

消费者:
  Kafka Consumer Group 保证: 一个 Partition 只被一个消费者消费
  → 同一 Partition 内单线程顺序消费 → 有序 ✓

  Consumer 1 ← Partition 0 (订单 A 的消息，顺序处理)
  Consumer 2 ← Partition 1 (订单 B 的消息，顺序处理)
```

```
但有一个坑 — 消费者内部多线程:

  Consumer 1 从 Partition 0 拉到: [创建A, 支付A, 发货A]
  如果用线程池并行处理:
    Thread 1: 处理 支付A（先完成）
    Thread 2: 处理 创建A（后完成）
    → 乱序了!

  解决: 消费者内部也要保证有序
    方案 A: 单线程处理（简单但慢）
    方案 B: 按 key 分派到固定线程（见下文）
```

#### 消费者内部有序 + 并行：按 key 分派

```go
// 拉到一批消息后，按 orderID 分组，同一 orderID 串行，不同 orderID 并行
func consumeOrdered(messages []Message, workerNum int) {
    // 每个 worker 一个有序队列
    workers := make([]chan Message, workerNum)
    for i := range workers {
        workers[i] = make(chan Message, 256)
    }

    // 启动 worker（每个 worker 顺序处理自己队列中的消息）
    var wg sync.WaitGroup
    for i := 0; i < workerNum; i++ {
        wg.Add(1)
        go func(ch chan Message) {
            defer wg.Done()
            for msg := range ch {
                processMessage(msg) // 顺序处理
            }
        }(workers[i])
    }

    // 按 key hash 分派到固定 worker
    for _, msg := range messages {
        idx := hash(msg.OrderID) % workerNum
        workers[idx] <- msg
        // 同一 orderID 总是分到同一个 worker → 有序
        // 不同 orderID 可以在不同 worker 并行 → 高吞吐
    }

    for _, ch := range workers {
        close(ch)
    }
    wg.Wait()
}
```

```
效果:
  16 个 worker，1000 个不同订单:
    订单A → worker 3:  [创建A, 支付A, 发货A] 顺序处理 ✓
    订单B → worker 7:  [创建B, 支付B] 顺序处理 ✓
    订单C → worker 3:  [创建C]（和订单A在同一worker，排在A后面）
    → 同一订单严格有序，不同订单并行
    → 吞吐 ≈ 无序消费的 80%~90%（非常好的平衡）
```

#### RocketMQ：顺序消息

```
RocketMQ 原生支持顺序消息:

生产者: MessageQueueSelector 指定队列
  producer.send(msg, new MessageQueueSelector() {
      public MessageQueue select(List<MessageQueue> mqs, Message msg, Object arg) {
          String orderID = (String) arg;
          int idx = orderID.hashCode() % mqs.size();
          return mqs.get(idx); // 同一 orderID 发到同一队列
      }
  }, orderID);

消费者: 使用 MessageListenerOrderly（而非 MessageListenerConcurrently）
  consumer.registerMessageListener(new MessageListenerOrderly() {
      public ConsumeOrderlyStatus consumeMessage(List<MessageExt> msgs, ConsumeOrderlyContext ctx) {
          // RocketMQ 保证同一队列的消息由同一线程顺序消费
          // 内部通过队列级别的锁实现
          processMessage(msgs);
          return ConsumeOrderlyStatus.SUCCESS;
      }
  });

原理:
  MessageListenerOrderly 内部:
  → 消费者对每个 Queue 加分布式锁（Broker 端）
  → 同一 Queue 只有一个消费者的一个线程处理
  → 保证队列内严格有序
```

#### RabbitMQ：天然不保证有序

```
RabbitMQ 的问题:
  - 单队列 + 单消费者 → 有序（但吞吐低）
  - 多消费者竞争消费同一队列 → 无法保证有序
  - 消息重试（NACK + requeue）→ 打乱顺序

解决方案:
  1. 单队列 + 单消费者（简单但吞吐低）
  2. 按 key 路由到不同队列（类似 Kafka 分区方案）:
     exchange → routing_key=hash(orderID) → queue_0 ~ queue_15
     每个队列一个消费者 → 队列内有序
  3. 消费者内部按 key 排序（不依赖 MQ 保证顺序）
```

### 生产者端保证发送有序

```
仅 Partition/Queue 级有序还不够，生产者也要保证有序发送:

Kafka 生产者:
  1. 同一 key 的消息用同一个生产者实例的同一个 Partition 发送
  2. 设置 max.in.flight.requests.per.connection = 1
     → 确保上一条消息 ACK 后才发下一条（防止重试导致乱序）
     → 或 enable.idempotence = true（开启幂等，自动保证顺序）

  // 问题: 如果消息发送失败重试呢？
  // msg1 发送成功，msg2 发送失败 → 重试 msg2 → msg3 先到 → 乱序!
  // 解决: enable.idempotence=true + acks=all
  //   幂等生产者为每条消息分配序列号，Broker 按序列号排序

RocketMQ 生产者:
  使用同步发送（不用异步/oneway）
  → 确认上一条成功后再发下一条
```

### 消费失败时的有序保障

```
问题: 有序消费时，如果某条消息处理失败怎么办？

  队列: [msg1, msg2, msg3]
  msg1 处理失败 → 跳过它继续处理 msg2？
  → 不行! msg2 可能依赖 msg1（如支付依赖创建）

方案 1: 阻塞重试（RocketMQ 默认）
  msg1 失败 → 重试 → 重试 → 成功 → 继续 msg2
  → 严格有序，但失败会阻塞后续所有消息
  → 需设置最大重试次数，超过后进死信队列 + 告警人工处理

方案 2: 跳过 + 补偿
  msg1 失败 → 记录到错误表 → 跳过继续处理 msg2
  → 后台异步重试/人工修复 msg1
  → 顺序"断开"了，但不阻塞其他消息
  → 适合容忍短暂不一致的场景

方案 3: 本地状态机校验
  // 消费前检查前置条件
  func consume(msg OrderEvent) error {
      order := db.GetOrder(msg.OrderID)
      switch msg.Type {
      case "pay":
          if order.Status != "created" {
              // 前置状态不对 → 可能是乱序 → 延迟重试
              return RetryLater
          }
          // 正常处理支付
      case "ship":
          if order.Status != "paid" {
              return RetryLater
          }
      }
  }
  // 即使消息乱序到达，也能通过状态校验保证正确性
```

### 方案选型

```
严格有序 + 高吞吐:
  → Kafka: 按 key hash 到 Partition + 消费者内部按 key 分派到固定 worker
  → RocketMQ: MessageQueueSelector + MessageListenerOrderly

严格有序 + 简单:
  → 单 Partition/Queue + 单消费者（吞吐低但实现简单）

容忍短暂乱序:
  → 消费者端状态机校验 + 延迟重试（不依赖 MQ 顺序，应用层保证）
  → 最灵活，MQ 选型不受限制

大多数场景的最佳实践:
  按业务 key 分区 + 消费者内部 key 分派 + 状态机校验兜底
  → 三层保障：生产→存储→消费 全链路有序
```

### 全链路有序总结

```
                       有序保障全链路

  生产者                    MQ                     消费者
  ─────────                ────                    ──────
  同一 key                 同一 key                 同一 key
  同一生产者实例            路由到同一 Partition      分派到同一 worker
  同步发送/幂等            Partition 内有序          worker 内顺序处理
  max.in.flight=1          ✓ 天然保证               ✓ 应用层保证
  或 enable.idempotence                            + 状态机校验兜底
```

### 一句话总结

**消息有序的核心 = 同一业务 key 的消息走同一条路径（同一 Partition → 同一消费者 → 同一 worker 线程）。** 生产者用 key hash 路由到固定 Partition（Kafka）/Queue（RocketMQ）保证存储有序；消费者内部按 key hash 分派到固定 worker 保证消费有序——同一 key 串行，不同 key 并行，吞吐可达无序消费的 80~90%。消费失败时用阻塞重试（严格有序）或状态机校验+延迟重试（容忍短暂乱序）。全局有序没必要（只能单分区单消费者），分区有序才是正确的设计。

---

## 四、如何实现按优先级消费

### 业务场景

```
场景: 同一个业务有不同优先级的消息，高优先级应该先被消费

  电商订单:
    VIP 用户的订单 > 普通用户的订单
  客服工单:
    投诉工单 > 咨询工单 > 建议工单
  任务调度:
    紧急任务 > 普通任务 > 低优任务
  支付回调:
    大额支付 > 小额支付（大额占用资金高，优先确认）
```

### 方案总览

| 方案 | 原理 | 优先级精度 | 复杂度 | 适用 MQ |
|------|------|-----------|--------|---------|
| **多 Topic/Queue 分级** | 不同优先级发到不同队列，高优队列分配更多消费者 | 粗粒度（3~5 级） | 低 | 所有 MQ |
| **RabbitMQ 原生优先级队列** | Broker 内部按 priority 排序 | 细粒度（0~255） | 低 | 仅 RabbitMQ |
| **消费者端优先级调度** | 消费者拉取后本地排序/优先处理 | 细粒度 | 中 | 所有 MQ |
| **加权消费** | 高优队列多拉几条，低优队列少拉 | 粗粒度 | 低 | 所有 MQ |

---

### 方案一：多 Topic/Queue 分级（最通用，生产首选）

#### 原理

```
按优先级拆分为多个独立的 Topic/Queue:

  生产者:
    高优消息 → topic_high
    中优消息 → topic_medium
    低优消息 → topic_low

  消费者分配策略:
    topic_high   → 4 个消费者（资源多，处理快）
    topic_medium → 2 个消费者
    topic_low    → 1 个消费者

  ┌──────────┐     ┌────────────┐     ┌──────────────────┐
  │ 生产者    │────→│ topic_high │────→│ Consumer 1,2,3,4 │  最快处理
  │          │     ├────────────┤     ├──────────────────┤
  │ 按优先级  │────→│topic_medium│────→│ Consumer 5,6     │  正常处理
  │ 路由      │     ├────────────┤     ├──────────────────┤
  │          │────→│ topic_low  │────→│ Consumer 7       │  慢慢处理
  └──────────┘     └────────────┘     └──────────────────┘
```

#### Go 实现

```go
// 生产者: 按优先级路由到不同 Topic
type Priority int
const (
    PriorityHigh   Priority = 1
    PriorityMedium Priority = 2
    PriorityLow    Priority = 3
)

func topicByPriority(baseTopic string, p Priority) string {
    switch p {
    case PriorityHigh:
        return baseTopic + "_high"
    case PriorityMedium:
        return baseTopic + "_medium"
    default:
        return baseTopic + "_low"
    }
}

func sendMessage(producer *kafka.Producer, baseTopic string, priority Priority, msg []byte) {
    topic := topicByPriority(baseTopic, priority)
    producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic},
        Value:          msg,
    }, nil)
}

// 消费者: 高优先级优先消费
func consumeWithPriority(highCh, mediumCh, lowCh <-chan Message, handler func(Message)) {
    for {
        // 先尝试高优
        select {
        case msg := <-highCh:
            handler(msg)
            continue
        default:
        }
        // 高优没有 → 尝试高优 or 中优
        select {
        case msg := <-highCh:
            handler(msg)
            continue
        case msg := <-mediumCh:
            handler(msg)
            continue
        default:
        }
        // 都没有 → 等任意一个
        select {
        case msg := <-highCh:
            handler(msg)
        case msg := <-mediumCh:
            handler(msg)
        case msg := <-lowCh:
            handler(msg)
        }
    }
}
```

| 优点 | 缺点 |
|------|------|
| 通用，所有 MQ 都支持 | 优先级粒度粗（一般 3~5 级） |
| 实现简单，运维清晰 | Topic 数量增多 |
| 高低优互不影响（低优堆积不影响高优） | 低优可能长时间得不到消费（饿死） |
| 可独立监控每个级别的堆积 | |

**解决低优饿死：** 用加权消费（见方案四）而非绝对优先。

---

### 方案二：RabbitMQ 原生优先级队列

```
RabbitMQ 是唯一原生支持消息优先级的主流 MQ:

创建优先级队列:
  channel.queueDeclare("order_queue", true, false, false,
      Map.of("x-max-priority", 10)  // 最大优先级 10
  );

发送带优先级的消息:
  channel.basicPublish("", "order_queue",
      new AMQP.BasicProperties.Builder()
          .priority(8)    // 优先级 0~10，值越大越优先
          .build(),
      body
  );

内部实现:
  每个优先级对应一个内部子队列
  消费者取消息时，Broker 从最高优先级的子队列开始投递

  ┌── order_queue ──────────────────────┐
  │  priority 10: [msg_a, msg_b]        │ ← 最先投递
  │  priority 8:  [msg_c]               │
  │  priority 5:  [msg_d, msg_e, msg_f] │
  │  priority 1:  [msg_g]               │ ← 最后投递
  └─────────────────────────────────────┘
```

| 优点 | 缺点 |
|------|------|
| 原生支持，无需应用层实现 | **仅 RabbitMQ 支持** |
| 细粒度（0~255） | 优先级级别多时内存开销增大 |
| 单队列管理简单 | 堆积时高优消息也受影响（共享队列资源） |
| | 低优消息可能饿死 |

**注意：** Kafka 和 RocketMQ 不支持原生优先级队列。

---

### 方案三：消费者端优先级调度

```
MQ 不感知优先级，消费者拉取后在本地按优先级处理:

  消费者拉取: 从一个 Topic 批量拉取 N 条消息
    ↓
  本地排序: 按消息中的 priority 字段排序
    ↓
  优先处理: 高优先级的先执行

  适合: 消息量不大、不想拆分 Topic 的简单场景
```

```go
type PriorityMessage struct {
    Priority int    // 优先级（值越小优先级越高）
    Body     []byte
}

// 用优先级队列（堆）管理待处理消息
type PriorityQueue []*PriorityMessage

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*PriorityMessage)) }
func (pq *PriorityQueue) Pop() any {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[:n-1]
    return item
}

// 消费者: 拉取 → 入堆 → 按优先级弹出处理
func consumeWithLocalPriority(consumer KafkaConsumer) {
    pq := &PriorityQueue{}
    heap.Init(pq)

    for {
        // 批量拉取
        msgs := consumer.Poll(100) // 拉取最多 100 条
        for _, msg := range msgs {
            pm := parsePriority(msg)
            heap.Push(pq, pm)
        }

        // 按优先级处理
        for pq.Len() > 0 {
            msg := heap.Pop(pq).(*PriorityMessage)
            processMessage(msg)
        }
    }
}
```

| 优点 | 缺点 |
|------|------|
| 不需要拆 Topic | 只在单批次内排序（批次间无法跨批优先） |
| 任意 MQ 都可用 | 高优消息可能在下一批次才被拉到 |
| 细粒度优先级 | 堆积时效果差（高优消息排在后面还没拉到） |

**适用：** 消息量不大、堆积不严重的场景。堆积严重时应用方案一。

---

### 方案四：加权消费（防低优饿死）

```
问题: 绝对优先级 → 高优消息持续不断时，低优永远不被消费（饿死）

加权消费: 高优多拉、低优少拉，但都能被消费到

  每轮消费:
    从 topic_high 拉 6 条
    从 topic_medium 拉 3 条
    从 topic_low 拉 1 条
    → 比例 6:3:1

  高优消费速度是低优的 6 倍，但低优不会饿死
```

```go
type WeightedConsumer struct {
    queues []struct {
        ch     <-chan Message
        weight int  // 每轮拉取条数
    }
}

func (wc *WeightedConsumer) Consume(handler func(Message)) {
    for {
        for _, q := range wc.queues {
            for i := 0; i < q.weight; i++ {
                select {
                case msg := <-q.ch:
                    handler(msg)
                default:
                    break // 该队列没消息了，跳过
                }
            }
        }
    }
}

// 配置
wc := &WeightedConsumer{
    queues: []struct{ ch <-chan Message; weight int }{
        {highCh, 6},   // 高优: 每轮最多 6 条
        {mediumCh, 3}, // 中优: 每轮最多 3 条
        {lowCh, 1},    // 低优: 每轮最多 1 条
    },
}
```

| 方案 | 高优处理速度 | 低优饿死？ | 适用场景 |
|------|------------|-----------|---------|
| 绝对优先级 | 最快 | 会 | 高优远多于低优，或低优可以等 |
| 加权消费 | 快（按比例） | **不会** | 各级别都有稳定流量 |
| 带保底的优先级 | 快 | 不会 | 低优有最低保障消费速率 |

---

### 各 MQ 的优先级能力对比

| MQ | 原生优先级 | 推荐方案 |
|----|-----------|---------|
| **Kafka** | 不支持 | 多 Topic 分级 + 加权消费 |
| **RocketMQ** | 不支持 | 多 Topic 分级 + 加权消费 |
| **RabbitMQ** | **原生支持（x-max-priority）** | 原生优先级队列 或 多队列分级 |
| **Redis Stream** | 不支持 | 多 Stream 分级 |

### 方案选型

```
优先级粒度少（2~5 级）+ 任意 MQ:
  → 多 Topic 分级 + 加权消费（最通用，生产首选）

优先级粒度多（0~255）+ 用 RabbitMQ:
  → 原生优先级队列

消息量小 + 不想拆 Topic:
  → 消费者端本地优先级队列

低优不能饿死:
  → 加权消费（6:3:1）或 带保底消费速率
```

### 一句话总结

**生产中最通用的优先级方案 = 多 Topic 分级（高/中/低优各一个 Topic）+ 加权消费（高优多拉低优少拉防饿死）。** Kafka/RocketMQ 不支持原生优先级，只能用多 Topic 方案；RabbitMQ 唯一原生支持（x-max-priority）。消费者端可以嵌套 select 实现"先看高优有没有 → 没有再看低优"的绝对优先级，或用加权比例（6:3:1）保证低优不饿死。消息量小的简单场景也可以在消费者本地用堆排序实现细粒度优先级。

---

## 五、如何保证消息不丢

### 消息丢失的三个环节

```
消息从生产到消费，经过三个环节，每个环节都可能丢消息:

  生产者               MQ Broker              消费者
  ┌──────┐   ①发送丢失  ┌──────┐  ③消费丢失  ┌──────┐
  │      │────────────→│      │────────────→│      │
  │ 发送  │             │ 存储  │             │ 消费  │
  └──────┘             └──────┘             └──────┘
                        ②存储丢失

① 生产者 → Broker: 网络抖动/超时 → 消息没到 Broker
② Broker 存储: Broker 宕机 → 内存中未刷盘的消息丢失
③ Broker → 消费者: 消费者拿到消息还没处理完就提交了 offset → 宕机后消息不会重新投递
```

---

### 一、生产者端：确保消息发出去

#### Kafka 生产者

```
关键配置:

  acks = all (或 -1)
    → 等待 ISR 中所有副本都写入成功才返回 ACK
    → acks=0: 发了就走（最快，可能丢）
    → acks=1: Leader 写入就返回（Leader 挂了丢）
    → acks=all: 所有 ISR 副本写入（不丢，最安全）

  retries = 3~5（或 Integer.MAX_VALUE）
    → 发送失败自动重试

  enable.idempotence = true
    → 开启幂等生产者，重试不会产生重复消息
    → Broker 用 PID + 序列号去重

  min.insync.replicas = 2（Broker 端配置）
    → ISR 中最少副本数，低于此值拒绝写入
    → 配合 acks=all，保证至少 2 个副本有数据
```

```go
// Go Kafka 生产者 — 保证不丢
config := sarama.NewConfig()
config.Producer.RequiredAcks = sarama.WaitForAll // acks=all
config.Producer.Retry.Max = 5                    // 重试 5 次
config.Producer.Return.Successes = true          // 返回成功结果
config.Producer.Idempotent = true                // 幂等
config.Net.MaxOpenRequests = 1                   // 配合幂等

producer, _ := sarama.NewSyncProducer(brokers, config) // 同步发送

// 同步发送 + 检查结果
partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
    Topic: "orders",
    Key:   sarama.StringEncoder(orderID),
    Value: sarama.ByteEncoder(body),
})
if err != nil {
    // 发送失败: 重试/写本地日志/入DB兜底
    log.Errorf("send failed: %v", err)
    saveToLocalDB(msg) // 本地持久化，后续补发
}
```

#### RocketMQ 生产者

```
三种发送方式的可靠性:

  同步发送 (send):        等 Broker ACK → 最可靠 ✓
  异步发送 (sendAsync):   回调中检查结果 → 可靠（需处理回调）
  单向发送 (sendOneway):  发了就走 → 不可靠 ✗（日志场景用）

  // 同步发送 + 重试
  producer.setRetryTimesWhenSendFailed(3);
  SendResult result = producer.send(msg); // 失败抛异常，自动重试 3 次
```

#### RabbitMQ 生产者

```
两种确认机制:

  1. Publisher Confirm（推荐）:
     channel.confirmSelect(); // 开启 confirm 模式
     channel.basicPublish(...);
     channel.waitForConfirms(); // 等待 Broker 确认
     // Broker 消息持久化后回 ACK → 生产者确认成功

  2. 事务模式（性能差，不推荐）:
     channel.txSelect();
     channel.basicPublish(...);
     channel.txCommit(); // 性能比 confirm 低 10 倍
```

#### 通用兜底：本地消息表

```
最终兜底方案: 先写本地 DB，再发 MQ

  BEGIN;
    INSERT INTO local_msg_table (id, topic, body, status) VALUES (..., 'pending');
  COMMIT;

  // 异步发送 MQ
  err := producer.Send(msg)
  if err == nil {
      UPDATE local_msg_table SET status='sent' WHERE id=?;
  }

  // 后台定时扫描 status='pending' 的消息重新发送
  // 保证: 即使 MQ 完全不可用，消息也不丢（在本地 DB 中）
```

---

### 二、Broker 端：确保消息存下来

#### Kafka Broker

```
关键配置:

  replication.factor = 3
    → 每个 Partition 3 个副本（1 Leader + 2 Follower）

  min.insync.replicas = 2
    → ISR 最少 2 个副本（配合 acks=all）
    → Leader + 至少 1 个 Follower 都写入才算成功
    → 即使 1 台 Broker 宕机也不丢

  unclean.leader.election.enable = false
    → 禁止非 ISR 副本成为 Leader
    → ISR 外的副本数据可能落后 → 成为 Leader 会丢消息

  刷盘策略:
    Kafka 默认异步刷盘（依赖 OS Page Cache）
    → 极端情况（整个机房断电）可能丢最后几条
    → 但 3 副本在不同机器 → 同时断电概率极低
    → 如果要求绝对不丢 → log.flush.interval.messages=1（同步刷盘，性能降 10 倍）

  三副本容灾:
    Broker 1 (Leader):    [msg1, msg2, msg3]
    Broker 2 (Follower):  [msg1, msg2, msg3]  ← ISR
    Broker 3 (Follower):  [msg1, msg2, msg3]  ← ISR

    Broker 1 宕机 → Broker 2 成为新 Leader → 数据完整 ✓
    Broker 1,2 同时宕机 → Broker 3 成为 Leader → 数据完整 ✓
    Broker 1,2,3 全挂 → 数据丢失（但三台同时挂概率极低）
```

#### RocketMQ Broker

```
两种刷盘模式:

  异步刷盘 (ASYNC_FLUSH): 默认
    → 写入 Page Cache 即返回 → Broker 崩溃可能丢几条
    → 性能高，大多数场景够用

  同步刷盘 (SYNC_FLUSH):
    → 写入磁盘后才返回 → 绝对不丢
    → 性能降 30%~50%
    → 金融/支付场景推荐

  主从同步:
    同步复制 (SYNC_MASTER): Master 等 Slave 写入后才返回 → 不丢
    异步复制 (ASYNC_MASTER): Master 写入即返回 → Master 宕机可能丢
```

#### RabbitMQ Broker

```
消息持久化三要素（缺一不可）:

  1. Exchange 持久化:  durable=true
  2. Queue 持久化:     durable=true
  3. 消息持久化:       deliveryMode=2

  channel.queueDeclare("orders", true, ...);  // durable=true
  channel.basicPublish("", "orders",
      MessageProperties.PERSISTENT_TEXT_PLAIN,  // deliveryMode=2
      body);

  镜像队列（高可用）:
    ha-mode=all → 消息在集群所有节点都有副本
    → 任何节点宕机不丢消息
```

---

### 三、消费者端：确保消息处理完

```
核心问题: 消费者拿到消息 → 处理 → 提交 offset/ACK
         如果处理完之前提交了 offset → 宕机后消息不会重新投递 → 丢了

原则: 先处理，后提交（At-Least-Once）
```

#### Kafka 消费者

```
关键: 手动提交 offset（不要用自动提交）

  enable.auto.commit = false（关闭自动提交！）

  自动提交的问题:
    拉取消息 → 还没处理完 → 定时器自动提交 offset → 消费者宕机
    → 消息丢失（Broker 认为已消费）

  手动提交:
    拉取消息 → 处理完成 → 手动提交 offset → 安全 ✓
```

```go
// Go Kafka 消费者 — 手动提交
config := sarama.NewConfig()
config.Consumer.Offsets.AutoCommit.Enable = false // 关闭自动提交

func consume(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
    // ① 处理消息（业务逻辑）
    err := processMessage(msg)
    if err != nil {
        // 处理失败: 重试/写死信队列，不提交 offset
        log.Errorf("process failed: %v", err)
        retryOrDLQ(msg)
        return
    }

    // ② 处理成功后才提交 offset
    session.MarkMessage(msg, "")
    session.Commit()  // 手动提交
}
```

#### RocketMQ 消费者

```
返回消费状态:

  consumer.registerMessageListener(msgs -> {
      try {
          processMessage(msgs);
          return ConsumeConcurrentlyStatus.CONSUME_SUCCESS;  // 成功 → ACK
      } catch (Exception e) {
          return ConsumeConcurrentlyStatus.RECONSUME_LATER;  // 失败 → 重试
      }
  });

  返回 RECONSUME_LATER:
    → Broker 不更新 offset → 延迟后重新投递
    → 重试 16 次仍失败 → 进入死信队列 → 人工处理
```

#### RabbitMQ 消费者

```
手动 ACK:

  channel.basicConsume("orders", false, ...);  // autoAck=false!

  try {
      processMessage(delivery);
      channel.basicAck(delivery.getEnvelope().getDeliveryTag(), false);  // 手动 ACK
  } catch (Exception e) {
      // 处理失败 → NACK + 重新入队
      channel.basicNack(tag, false, true);  // requeue=true
  }

  autoAck=true 的问题:
    Broker 投递消息后立即从队列删除 → 消费者宕机 → 消息丢了
```

---

### 四、消费失败的兜底：死信队列（DLQ）

```
消息重试多次仍然失败 → 不能无限重试（会阻塞队列）→ 进入死信队列

  正常队列: [msg1, msg2, msg3(失败)]
    → msg3 重试 3 次仍失败
    → msg3 转入死信队列

  死信队列: [msg3]
    → 人工排查/修复后重新消费
    → 或丢弃（通知业务方）

各 MQ 的死信:
  Kafka:     无原生死信，应用层自己实现（发到 xxx_dlq Topic）
  RocketMQ:  原生支持（%DLQ%ConsumerGroup）
  RabbitMQ:  原生支持（x-dead-letter-exchange）
```

```go
// Go 实现死信队列逻辑
const maxRetry = 3

func consumeWithDLQ(msg Message) error {
    retryCount := msg.Header.GetInt("retry_count")

    err := processMessage(msg)
    if err != nil {
        if retryCount >= maxRetry {
            // 超过最大重试次数 → 发到死信队列
            sendToDLQ(msg, err)
            return nil // 不再重试
        }
        // 重试: 带上重试次数，延迟重新投递
        msg.Header.Set("retry_count", retryCount+1)
        sendWithDelay(msg, backoffDuration(retryCount)) // 指数退避
        return nil
    }
    return nil
}

func backoffDuration(retry int) time.Duration {
    // 1s → 5s → 30s
    delays := []time.Duration{time.Second, 5 * time.Second, 30 * time.Second}
    if retry < len(delays) {
        return delays[retry]
    }
    return delays[len(delays)-1]
}
```

---

### 五、At-Least-Once 带来的重复问题

```
"先处理后提交"保证了不丢，但引入了新问题 — 重复消费:

  消费者处理完 msg1 → 准备提交 offset → 宕机了 → offset 没提交
  → 恢复后从上次 offset 重新拉取 → msg1 又被消费一次 → 重复!

  ┌─────────────────────────────────────────────┐
  │  不丢 (At-Least-Once) + 幂等 = 恰好一次效果   │
  └─────────────────────────────────────────────┘

幂等方案:
  1. 数据库唯一键:
     INSERT ... ON DUPLICATE KEY UPDATE → 重复消费自动忽略

  2. 消息去重表:
     消费前查: SELECT 1 FROM msg_dedup WHERE msg_id = ?
     有 → 跳过; 无 → 处理 + INSERT msg_dedup

  3. 业务幂等:
     状态机: 只有 status=created 才能变为 paid → 重复消费不影响

  4. Redis 去重:
     SET msg:{id} EX 86400 NX → 设置成功才处理，失败说明已处理过
```

---

### 全链路不丢总结

```
  生产者                    Broker                    消费者
  ───────                  ──────                    ──────
  ① acks=all              ④ 副本数≥3                 ⑦ 关闭自动提交
  ② 重试+幂等             ⑤ min.insync.replicas≥2    ⑧ 先处理后提交
  ③ 本地消息表兜底         ⑥ 同步刷盘(金融场景)       ⑨ 失败重试+死信队列
                                                     ⑩ 幂等消费防重复
```

| 环节 | 丢消息场景 | 解决方案 |
|------|-----------|---------|
| 生产者→Broker | 网络失败/Broker 未确认 | acks=all + 重试 + 本地消息表 |
| Broker 存储 | Broker 宕机未刷盘 | 多副本 + min.insync.replicas + 同步刷盘 |
| Broker→消费者 | 消费者提交 offset 后宕机 | 关闭自动提交 + 先处理后提交 |
| 消费失败 | 处理异常丢弃 | 重试机制 + 死信队列 + 人工修复 |
| 重复消费（副作用） | At-Least-Once 必然带来 | 幂等（唯一键/去重表/状态机） |

### 一句话总结

**消息不丢 = 三端都要保障：生产者用 acks=all+重试+本地消息表确保发出去；Broker 用多副本(≥3)+min.insync.replicas(≥2)确保存下来；消费者关闭自动提交+先处理后提交确保消费完。** 消费失败用指数退避重试+死信队列兜底。At-Least-Once 不可避免带来重复消费，用幂等（唯一键/去重表/状态机）解决。金融级不丢还需 Broker 同步刷盘（RocketMQ SYNC_FLUSH）。本地消息表是终极兜底——MQ 完全不可用时消息也不丢。

---

## 六、如何保证消息不重复消费（幂等性）

### 为什么会重复消费

```
重复消费是 At-Least-Once 语义的必然副产品:

  场景 1: 消费者处理完但 offset 提交失败
    消费者拉取 msg1 → 处理成功 → 提交 offset → 网络超时/消费者宕机
    → offset 没提交成功 → 恢复后从旧 offset 拉取 → msg1 再次被消费

  场景 2: Kafka Rebalance
    Consumer A 处理 msg1 完成 → 还没提交 → 触发 Rebalance
    → Partition 分配给 Consumer B → B 从上次提交的 offset 重新拉取 → msg1 重复

  场景 3: 生产者重试
    生产者发送 msg1 → Broker 收到并存储 → ACK 返回超时
    → 生产者重试发送 msg1 → Broker 存了两条 msg1 → 消费者收到两次

  场景 4: RabbitMQ NACK + requeue
    消费者处理失败 → NACK requeue=true → 消息重新入队 → 被再次投递

结论: 在分布式系统中，消息重复几乎不可避免
     → 不应该在 MQ 层面解决，而应该在消费者端保证幂等
```

### 三种消息语义

```
At-Most-Once (最多一次):
  消息可能丢，但不会重复
  → 自动提交 offset → 处理失败也不重试
  → 适合: 日志、监控指标（丢几条无所谓）

At-Least-Once (至少一次): ← 大多数系统选这个
  消息不会丢，但可能重复
  → 先处理后提交 + 失败重试
  → 配合幂等 → 等效于 Exactly-Once

Exactly-Once (恰好一次):
  消息不丢不重 → 理论上在分布式系统中不存在
  → 实际用 At-Least-Once + 幂等 来模拟
  → Kafka 的 Exactly-Once 也是通过 事务+幂等生产者 实现的
```

### 幂等性方案

#### 方案一：数据库唯一键（最常用）

```
原理: 利用数据库 UNIQUE KEY 约束，重复消费时自然被拦截

适用: 消费结果是写数据库的场景（插入/更新）
```

```go
// 场景: 消费消息后创建订单
func consumeCreateOrder(msg OrderMessage) error {
    _, err := db.Exec(`
        INSERT INTO orders (order_id, user_id, amount, status)
        VALUES (?, ?, ?, 'created')
        ON DUPLICATE KEY UPDATE order_id = order_id`,  // 重复时什么也不做
        msg.OrderID, msg.UserID, msg.Amount,
    )
    // order_id 有唯一索引
    // 第一次: 插入成功
    // 重复消费: ON DUPLICATE KEY → 忽略 → 不影响数据
    return err
}
```

```
优点: 最简单可靠，数据库保证原子性
缺点: 只适用于插入场景；高并发下唯一键冲突有性能开销
```

#### 方案二：去重表（消息级去重）

```
原理: 维护一张已消费消息的 ID 表，消费前先查是否已处理过
```

```go
func consumeWithDedup(msg Message) error {
    // ① 检查是否已消费（用 msgID 去重）
    var exists bool
    db.QueryRow("SELECT 1 FROM msg_dedup WHERE msg_id = ?", msg.ID).Scan(&exists)
    if exists {
        return nil // 已消费过，跳过
    }

    // ② 处理业务逻辑 + 写去重记录（在同一个事务中！）
    tx, _ := db.Begin()
    defer tx.Rollback()

    if err := processOrder(tx, msg); err != nil {
        return err
    }

    // 插入去重记录
    _, err := tx.Exec("INSERT INTO msg_dedup (msg_id, created_at) VALUES (?, NOW())", msg.ID)
    if err != nil {
        return err // 如果是唯一键冲突 → 说明并发重复，忽略
    }

    return tx.Commit()
}

// 去重表结构
// CREATE TABLE msg_dedup (
//     msg_id VARCHAR(64) PRIMARY KEY,
//     created_at DATETIME
// ) ENGINE=InnoDB;
// 定期清理: DELETE FROM msg_dedup WHERE created_at < DATE_SUB(NOW(), INTERVAL 7 DAY);
```

```
关键: 业务操作和去重记录在同一个事务中
  → 要么都成功，要么都失败
  → 不会出现"处理了但去重记录没写"或"去重记录写了但没处理"
```

#### 方案三：Redis 去重（高性能）

```go
func consumeWithRedisDedup(msg Message) error {
    // SET NX: 只有 key 不存在时才设置成功
    ok, err := rdb.SetNX(ctx, "dedup:"+msg.ID, "1", 24*time.Hour).Result()
    if err != nil {
        return err
    }
    if !ok {
        return nil // key 已存在 → 已消费过 → 跳过
    }

    // 处理业务
    if err := processMessage(msg); err != nil {
        // 处理失败 → 删除去重 key → 下次可以重试
        rdb.Del(ctx, "dedup:"+msg.ID)
        return err
    }
    return nil
}
```

```
优点: 极快（Redis 操作 <1ms），适合高 QPS 场景
缺点:
  - Redis 和业务 DB 不在同一个事务中 → 有极小概率不一致
    (Redis SET 成功但业务处理失败、且 Del 也失败 → 消息被"假消费")
  - Redis 宕机可能丢去重记录
适用: 对一致性要求不极高的场景（大多数互联网业务）
```

#### 方案四：业务状态机（最优雅）

```
原理: 业务本身的状态流转保证幂等，不需要额外的去重机制

  订单状态机: created → paid → shipped → completed
                              ↑
  消费"支付成功"消息时:
    当前状态 = created → 更新为 paid ✓
    当前状态 = paid    → 已经是 paid 了，跳过 ✓（幂等）
    当前状态 = shipped → 状态已经更靠后了，跳过 ✓
```

```go
func consumePaymentSuccess(msg PaymentMsg) error {
    result, err := db.Exec(`
        UPDATE orders
        SET status = 'paid', paid_at = NOW()
        WHERE order_id = ? AND status = 'created'`,  // 关键: WHERE 带上前置状态
        msg.OrderID,
    )
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        // 没有更新到 → 状态不是 'created' → 已经处理过了 → 幂等
        log.Infof("order %s already paid, skip", msg.OrderID)
    }
    return nil
}
```

```
优点:
  - 不需要额外的去重表/Redis
  - 和业务逻辑天然结合
  - 用数据库乐观锁（WHERE status=前置状态）保证
  - 并发安全（两个相同消息同时执行，只有一个能 UPDATE 成功）

适用: 有明确状态流转的业务（订单、支付、工单、审批流）
```

#### 方案五：全局唯一业务 ID + 乐观锁

```go
// 场景: 扣减库存（不是状态机，但需要幂等）
func consumeDeductStock(msg DeductMsg) error {
    // 用 request_id（业务幂等键）做去重
    result, err := db.Exec(`
        UPDATE inventory
        SET stock = stock - ?,
            last_deduct_id = ?
        WHERE product_id = ?
          AND stock >= ?
          AND (last_deduct_id IS NULL OR last_deduct_id != ?)`,
        msg.Quantity,
        msg.RequestID,  // 全局唯一的扣减请求 ID
        msg.ProductID,
        msg.Quantity,
        msg.RequestID,  // 同一个 request_id 只能扣一次
    )
    rows, _ := result.RowsAffected()
    if rows == 0 {
        // 没扣到: 库存不足 或 已经扣过（幂等）
    }
    return err
}
```

### 各方案对比

| 方案 | 可靠性 | 性能 | 复杂度 | 适用场景 |
|------|--------|------|--------|---------|
| 数据库唯一键 | 高 | 中 | 低 | 插入场景 |
| 去重表 | 高 | 中 | 中 | 通用（任何场景） |
| Redis 去重 | 中高 | **高** | 低 | 高 QPS + 一致性要求不极高 |
| **业务状态机** | **高** | **高** | **低** | 有状态流转的业务（首选） |
| 全局 ID + 乐观锁 | 高 | 高 | 中 | 更新场景（扣减/累加） |

### 生产中的推荐组合

```
最佳实践: 分层幂等

  第一层 — 生产者去重（防重复发送）:
    Kafka: enable.idempotence=true（PID+序列号，Broker 自动去重）
    RocketMQ: 消息 key 设置为业务唯一 ID

  第二层 — 消费者去重（防重复消费）:
    优先用业务状态机（UPDATE WHERE status=前置状态）
    没有状态机 → Redis SET NX 快速去重
    金融级要求 → 去重表（和业务在同一事务）

  第三层 — 兜底:
    监控重复消费率（正常应 < 0.1%）
    定期核对数据一致性
```

```
Kafka Exactly-Once 语义（生产者+消费者端到端）:

  生产者: enable.idempotence=true + 事务 API
  消费者: read_committed 隔离级别 + 消费+offset提交在同一事务

  // 仅限 Kafka → Kafka 的场景（如 Stream Processing）
  // 涉及外部系统（DB）时仍需应用层幂等
```

### 一句话总结

**消息重复消费在分布式系统中不可避免（Rebalance/重试/网络超时），正确做法不是防重复，而是保证幂等——消费一次和消费多次的效果一样。** 最优方案是业务状态机（UPDATE WHERE status=前置状态，零额外开销）；通用方案是去重表（业务操作+去重记录在同一事务）；高性能方案是 Redis SET NX。生产中推荐分层幂等：生产者幂等（Kafka idempotent）+ 消费者状态机/去重 + 监控兜底。核心公式：**At-Least-Once + 幂等 = Exactly-Once 效果**。

---

## 七、MQ 如何保证分布式事务一致性

### 问题场景

```
电商下单:
  订单服务: 创建订单（MySQL A）
  库存服务: 扣减库存（MySQL B）

  两个服务各有自己的数据库，无法用本地事务一起提交

  可能出现的不一致:
    订单创建成功 + 扣库存失败 → 超卖
    订单创建失败 + 扣库存成功 → 多扣了库存
    订单创建成功 + 发 MQ 失败 → 库存没扣

  目标: 要么都成功，要么都回滚（最终一致性）
```

### 分布式事务方案总览

| 方案 | 一致性 | 性能 | 复杂度 | 适用场景 |
|------|--------|------|--------|---------|
| 2PC（两阶段提交） | 强一致 | 低（同步阻塞） | 高 | 传统数据库分布式事务 |
| TCC | 强一致 | 中 | 很高（三个接口） | 金融级强一致 |
| **本地消息表** | **最终一致** | **高** | **中** | **最常用，推荐** |
| **RocketMQ 事务消息** | **最终一致** | **高** | **低** | **已有 RocketMQ 时首选** |
| Saga | 最终一致 | 高 | 高（补偿逻辑） | 长事务/多步骤编排 |
| 最大努力通知 | 最终一致 | 高 | 低 | 跨企业/弱一致性 |

---

### 方案一：本地消息表（最通用）

#### 原理

```
核心思想: 将"发 MQ"和"业务操作"放在同一个本地事务中

  订单服务 (同一个数据库事务):
    ① INSERT INTO orders (...) → 创建订单
    ② INSERT INTO local_msg_table (msg_id, topic, body, status='pending')
    ③ COMMIT → 两条记录要么都写入，要么都不写入（本地事务保证）

  后台定时任务:
    ④ SELECT * FROM local_msg_table WHERE status='pending'
    ⑤ 发送 MQ 消息
    ⑥ 发送成功 → UPDATE status='sent'
       发送失败 → 下次继续重试

  库存服务:
    ⑦ 消费 MQ 消息 → 扣减库存（幂等处理）
```

```
流程图:

  订单服务                    MQ                     库存服务
    │                         │                        │
    │ BEGIN                   │                        │
    │ INSERT orders           │                        │
    │ INSERT local_msg_table  │                        │
    │ COMMIT ✓                │                        │
    │                         │                        │
    │ 定时扫描 pending 消息    │                        │
    │─── 发送消息 ───────────→│                        │
    │                         │─── 投递消息 ──────────→│
    │ UPDATE status='sent'    │                        │ 扣减库存
    │                         │                        │ (幂等)
    │                         │                        │

  为什么可靠？
    订单+消息记录在同一事务中 → 原子性
    消息发送失败 → 定时任务重试 → 最终一定能发出去
    库存服务幂等 → 重复投递不会多扣
    → 最终一致性 ✓
```

#### Go 实现

```go
// 订单服务: 创建订单 + 写本地消息表
func CreateOrder(ctx context.Context, req CreateOrderReq) error {
    tx, _ := db.BeginTx(ctx, nil)
    defer tx.Rollback()

    // ① 创建订单
    orderID := generateOrderID()
    _, err := tx.Exec(`INSERT INTO orders (id, user_id, product_id, amount, status)
        VALUES (?, ?, ?, ?, 'created')`, orderID, req.UserID, req.ProductID, req.Amount)
    if err != nil {
        return err
    }

    // ② 写本地消息表（同一个事务！）
    msgBody, _ := json.Marshal(DeductStockMsg{
        OrderID:   orderID,
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
    })
    _, err = tx.Exec(`INSERT INTO local_msg_table (msg_id, topic, body, status, created_at)
        VALUES (?, 'deduct_stock', ?, 'pending', NOW())`, orderID, msgBody)
    if err != nil {
        return err
    }

    // ③ 提交事务 → 订单和消息记录原子写入
    return tx.Commit()
}

// 后台定时任务: 扫描 pending 消息发送到 MQ
func SendPendingMessages() {
    rows, _ := db.Query(`SELECT msg_id, topic, body FROM local_msg_table
        WHERE status = 'pending' AND created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)
        LIMIT 100`)
    defer rows.Close()

    for rows.Next() {
        var msgID, topic string
        var body []byte
        rows.Scan(&msgID, &topic, &body)

        err := producer.Send(topic, msgID, body)
        if err == nil {
            db.Exec(`UPDATE local_msg_table SET status='sent' WHERE msg_id=?`, msgID)
        }
        // 失败不更新 → 下次扫描继续重试
    }
}

// 库存服务: 消费消息扣库存（幂等）
func ConsumeDeductStock(msg DeductStockMsg) error {
    result, _ := db.Exec(`UPDATE inventory SET stock = stock - ?
        WHERE product_id = ? AND stock >= ?
        AND product_id NOT IN (
            SELECT product_id FROM deduct_log WHERE order_id = ?
        )`, msg.Quantity, msg.ProductID, msg.Quantity, msg.OrderID)

    rows, _ := result.RowsAffected()
    if rows > 0 {
        // 扣减成功，记录去重日志
        db.Exec(`INSERT INTO deduct_log (order_id, product_id) VALUES (?, ?)`,
            msg.OrderID, msg.ProductID)
    }
    return nil
}
```

---

### 方案二：RocketMQ 事务消息（最优雅）

#### 原理

```
RocketMQ 事务消息 = 把"本地消息表"的逻辑内置到了 MQ 中

  ┌──────────────────────────────────────────────────────────────┐
  │                                                              │
  │  ① 发送半消息 (Half Message)                                  │
  │     生产者 → Broker: "我有一条消息，先存着但别投递给消费者"       │
  │     Broker 存到内部 Topic（RMQ_SYS_TRANS_HALF_TOPIC）         │
  │                                                              │
  │  ② 执行本地事务                                               │
  │     生产者: 本地执行业务逻辑（如创建订单）                      │
  │     → 成功: 返回 COMMIT_MESSAGE                               │
  │     → 失败: 返回 ROLLBACK_MESSAGE                             │
  │     → 不确定: 返回 UNKNOW（等回查）                            │
  │                                                              │
  │  ③ 二次确认                                                   │
  │     COMMIT: Broker 将半消息转投到目标 Topic → 消费者可见        │
  │     ROLLBACK: Broker 删除半消息 → 消费者永远看不到              │
  │                                                              │
  │  ④ 事务回查（兜底）                                           │
  │     如果 Broker 长时间没收到二次确认（生产者宕机/网络断开）      │
  │     → Broker 主动回查生产者: "你那个本地事务到底成功没？"        │
  │     → 生产者查本地数据库 → 回复 COMMIT 或 ROLLBACK             │
  │                                                              │
  └──────────────────────────────────────────────────────────────┘
```

```
完整时序图:

  生产者(订单服务)               Broker                   消费者(库存服务)
     │                            │                          │
     │ ① 发送半消息               │                          │
     │─── Half Message ─────────→│ 暂存（消费者不可见）       │
     │                            │                          │
     │ ② 执行本地事务              │                          │
     │ (创建订单 INSERT)           │                          │
     │                            │                          │
     │ 成功:                       │                          │
     │ ③ COMMIT ────────────────→│ 半消息 → 目标 Topic       │
     │                            │──── 投递消息 ───────────→│
     │                            │                          │ 扣减库存
     │                            │                          │
     │ 失败:                       │                          │
     │ ③ ROLLBACK ──────────────→│ 删除半消息               │
     │                            │ (消费者永远看不到)        │
     │                            │                          │
     │ 超时/宕机:                  │                          │
     │                            │ ④ 回查本地事务状态        │
     │←── 回查请求 ───────────────│                          │
     │─── COMMIT/ROLLBACK ──────→│                          │
```

#### 为什么比本地消息表更好

```
本地消息表:
  需要额外的 local_msg_table 表
  需要写定时扫描任务
  需要处理消息状态管理
  → 每个项目都要重复实现

RocketMQ 事务消息:
  MQ 内部实现了"本地消息表"的逻辑
  半消息 = 本地消息表的 pending 状态
  事务回查 = 定时扫描任务
  → 开箱即用，不需要额外表和定时任务

本质上是同一个思路，只是实现层面下沉到了 MQ 中
```

---

### 方案三：TCC（Try-Confirm-Cancel）

```
TCC 要求每个服务实现三个接口:

  Try:     预留资源（冻结库存、冻结余额）
  Confirm: 确认操作（扣减冻结的库存）
  Cancel:  取消操作（解冻库存）

  流程:
  ┌──────────────┐
  │  TCC 协调者   │
  └──────┬───────┘
         │
    ┌────▼────┐     ┌────────────┐
    │ Try 阶段 │     │ 订单服务    │ → Try: 创建订单(pending)
    │         │     │ 库存服务    │ → Try: 冻结库存(stock=100, frozen=5)
    │         │     │ 支付服务    │ → Try: 冻结余额
    └────┬────┘     └────────────┘
         │ 全部 Try 成功
    ┌────▼────────┐
    │ Confirm 阶段 │ → 订单: pending→created; 库存: frozen-5; 余额: frozen扣减
    └─────────────┘
         │ 某个 Try 失败
    ┌────▼────────┐
    │ Cancel 阶段  │ → 订单: 删除; 库存: 解冻; 余额: 解冻
    └─────────────┘

优点: 强一致性，实时生效
缺点: 每个服务要写三套逻辑（Try/Confirm/Cancel），极其复杂
      空回滚/悬挂/幂等等异常场景处理繁琐
适用: 金融级强一致性（转账、支付）
```

---

### 方案四：Saga 模式

```
Saga = 将长事务拆分为多个本地事务 + 补偿事务

  正向流程:
    T1(创建订单) → T2(扣库存) → T3(扣余额) → 全部成功 ✓

  某步失败时反向补偿:
    T1(创建订单) → T2(扣库存) → T3(扣余额失败!)
    → C2(恢复库存) → C1(取消订单)

  编排方式:
    编排式 (Choreography): 每个服务监听事件自行决定（去中心化）
    协调式 (Orchestration): 中央协调器统一编排（推荐）

适用: 步骤多、流程长的业务（旅游预订: 机票+酒店+租车）
缺点: 补偿逻辑复杂，中间状态可见（不隔离）
```

---

### 各方案对比

| 维度 | 本地消息表 | RocketMQ 事务消息 | TCC | Saga |
|------|-----------|-----------------|-----|------|
| **一致性** | 最终一致 | 最终一致 | 强一致 | 最终一致 |
| **性能** | 高 | 高 | 中（两阶段） | 高 |
| **侵入性** | 中（加表+定时任务） | 低（MQ 内置） | 高（三套接口） | 高（补偿逻辑） |
| **隔离性** | 无 | 无 | 有（资源冻结） | 无（中间状态可见） |
| **MQ 依赖** | 任意 MQ | 仅 RocketMQ | 不需要 MQ | 可选 |
| **适用场景** | 通用 | 已有 RocketMQ | 金融强一致 | 长流程编排 |

### 下单扣库存的推荐方案

```
大多数互联网公司（最终一致即可）:

  方案 A: 已有 RocketMQ
    → RocketMQ 事务消息（开箱即用，最省事）

  方案 B: 用 Kafka 或其他 MQ
    → 本地消息表（通用方案，任何 MQ 都行）

金融/支付（需要强一致）:
    → TCC（Seata 框架支持）

长流程（5 步以上）:
    → Saga 编排式

选型口诀:
  能最终一致 → MQ 事务消息 / 本地消息表（简单高效）
  必须强一致 → TCC（复杂但可靠）
  步骤很多   → Saga（灵活但需要补偿逻辑）
```

### 一句话总结

**MQ 保证分布式事务一致性的核心思路 = 本地事务 + 可靠消息 + 幂等消费。** 本地消息表方案：业务操作和消息记录在同一个本地事务中原子写入，后台定时发送 MQ，消费者幂等处理——三步保证最终一致。RocketMQ 事务消息是本地消息表的 MQ 内置版：半消息（暂存不投递）→ 本地事务 → 二次确认（COMMIT/ROLLBACK）→ 超时回查兜底。TCC 用于强一致（Try 冻结资源 / Confirm 确认 / Cancel 释放），复杂但实时。**大多数场景用本地消息表或 RocketMQ 事务消息足够，TCC 只在金融级强一致时才需要。**

---

## 八、Kafka 为什么吞吐量这么高

### 性能数据

```
Kafka 单 Broker:
  写入: 数十万~百万 msg/s（取决于消息大小）
  读取: 几十万~百万 msg/s
  吞吐: 几百 MB/s ~ GB/s

对比:
  Kafka:    百万 msg/s
  RocketMQ: 十万 msg/s
  RabbitMQ: 万 msg/s

Kafka 快 10~100 倍的原因 = 六大技术的叠加效应
```

### 一、顺序 I/O（最核心）

```
传统数据库写入: 随机 I/O
  B+ 树索引 → 写入时需要定位到不同的磁盘位置
  → 机械硬盘寻道 ~10ms → 每秒 ~100 次随机写

Kafka 写入: 顺序 I/O
  消息追加写到 Partition 日志文件末尾（Append-Only）
  → 磁盘顺序写 ~600MB/s（HDD）/ 几 GB/s（SSD）
  → 不需要寻道，磁头不用移动

  ┌──────────────────────────────────────────┐
  │          Partition 日志文件                │
  │                                          │
  │  [msg1][msg2][msg3][msg4] ... [msgN][新消息]│
  │                                     ↑    │
  │                               追加写到末尾  │
  └──────────────────────────────────────────┘

  顺序写磁盘甚至比随机写内存还快:
    顺序写 HDD:    ~600 MB/s
    随机写内存:     ~50 MB/s（取决于访问模式）
    随机写 HDD:     ~1 MB/s（寻道开销）

  这就是 Kafka 敢把数据全写磁盘而不是内存的底气
```

### 二、Page Cache（OS 页缓存）

```
Kafka 不自己管理缓存，而是利用 OS 的 Page Cache:

  写入:
    Producer → Kafka → write() 系统调用
    → 数据写入 Page Cache（内存）→ 立即返回
    → OS 后台异步刷到磁盘（pdflush）

  读取（消费）:
    Consumer → Kafka → read() 系统调用
    → 数据在 Page Cache 中 → 直接从内存返回（不读磁盘！）

  为什么有效？
    消费者通常读的是最近写入的数据（几秒~几分钟前）
    → 这些数据大概率还在 Page Cache 中
    → 生产和消费几乎都是纯内存操作

  ┌──────────────────────────────────────┐
  │  Producer → write() → Page Cache     │ ← 写内存
  │                         ↕ (热数据)    │
  │  Consumer → read()  → Page Cache     │ ← 读内存
  │                         ↓ (冷数据)    │
  │                       Disk            │ ← 只有冷数据才读磁盘
  └──────────────────────────────────────┘

  优势:
    1. Kafka 进程重启 → Page Cache 不丢（OS 管理）→ 重启后消费仍然快
    2. 不需要在 JVM 内管理缓存 → 避免 GC 开销
    3. OS 的缓存淘汰算法(LRU)非常成熟
```

### 三、零拷贝（sendfile）

```
传统数据发送（4 次拷贝 + 4 次上下文切换）:

  磁盘 → 内核缓冲区 → 用户空间缓冲区 → Socket 缓冲区 → 网卡
         (DMA拷贝)    (CPU拷贝)         (CPU拷贝)       (DMA拷贝)

  ① read():  磁盘 → 内核 Page Cache (DMA)
  ② read():  Page Cache → 用户空间 (CPU 拷贝)
  ③ write(): 用户空间 → Socket 缓冲区 (CPU 拷贝)
  ④ write(): Socket 缓冲区 → 网卡 (DMA)
  → 4 次拷贝，2 次用户态-内核态切换

零拷贝 sendfile（2 次拷贝 + 0 次用户态拷贝）:

  磁盘 → 内核 Page Cache → 网卡
         (DMA拷贝)          (DMA拷贝)

  sendfile(): Page Cache 直接发到网卡，不经过用户空间
  → 数据不需要拷贝到 JVM 堆内存
  → 2 次 DMA 拷贝（硬件做，CPU 几乎不参与）
  → 0 次 CPU 拷贝

  ┌──────────────────────────────────────────────────────┐
  │ 传统方式:                                             │
  │   磁盘 ─DMA→ 内核缓冲 ─CPU→ 用户空间 ─CPU→ Socket ─DMA→ 网卡│
  │                                                      │
  │ 零拷贝 (sendfile):                                    │
  │   磁盘 ─DMA→ Page Cache ──────────────────DMA→ 网卡   │
  │                  (跳过用户空间！)                       │
  └──────────────────────────────────────────────────────┘

Kafka 消费时:
  Consumer 请求读消息 → Broker 调用 sendfile()
  → 消息从 Page Cache 直接发到网卡 → 发给 Consumer
  → CPU 几乎不参与数据搬运
```

### 四、批量处理 + 压缩

```
批量发送:
  Producer 不是每条消息都发一次网络请求
  → 攒一批（batch.size=16KB 或 linger.ms=5ms）再一起发
  → 一次网络 round trip 发送几十~几百条消息

  不批量: 1000 条 × 1次网络请求 = 1000 次 RT
  批量:   1000 条 ÷ 100条/批 = 10 次 RT → 减少 99%

批量压缩:
  一批消息整体压缩（而非逐条压缩）→ 压缩率更高
  支持: Snappy / LZ4 / GZIP / ZSTD

  ┌──────────────────────────────────────┐
  │  原始 batch: [msg1][msg2]...[msg100] │ 16KB
  │  压缩后:     [compressed data]        │ 4KB  (压缩率 75%)
  └──────────────────────────────────────┘

  → 网络传输量减少 75%
  → 磁盘占用减少 75%
  → 整批压缩比逐条压缩效果好（相邻消息相似度高）

Producer 配置:
  batch.size = 16384          // 批次大小 16KB
  linger.ms = 5               // 最多等 5ms 凑批
  compression.type = lz4      // LZ4: 压缩快、CPU 低
```

### 五、分区并行（Partition）

```
一个 Topic 分成多个 Partition → 天然并行:

  Topic: orders (16 个 Partition)
    P0  → Broker 1 → Consumer 1
    P1  → Broker 1 → Consumer 2
    P2  → Broker 2 → Consumer 3
    ...
    P15 → Broker 4 → Consumer 16

  写入并行: 不同 Partition 的消息可以同时写入不同 Broker
  消费并行: 每个 Partition 由不同消费者并行消费
  存储并行: 不同 Partition 在不同磁盘上

  吞吐 ≈ 单 Partition 吞吐 × Partition 数
  16 个 Partition × 10MB/s/Partition ≈ 160MB/s
```

### 六、高效的日志存储格式

```
每个 Partition 是一个有序的日志目录:

  topic-orders-partition-0/
    ├── 00000000000000000000.log     (数据文件，Segment)
    ├── 00000000000000000000.index   (偏移量索引)
    ├── 00000000000000000000.timeindex (时间索引)
    ├── 00000000000065536000.log     (第二个 Segment)
    ├── 00000000000065536000.index
    └── ...

  日志分段 (Segment):
    每个 Segment 默认 1GB（可配）
    老 Segment 可以直接删除（不影响新消息）→ 清理高效
    按 offset 命名 → 用二分查找快速定位

  稀疏索引:
    不是每条消息都索引，而是每隔 N 条记一个索引点
    → 先用索引定位到大概位置 → 再顺序扫描几条
    → 索引文件极小，常驻内存

  查找 offset=65540 的消息:
    ① 二分查找 Segment 文件: 65540 在 00000000000065536000.log 中
    ② 查 .index 文件: 找到 offset=65536 在文件偏移 0 处
    ③ 从偏移 0 开始顺序读 4 条 → 找到 offset=65540
    → 整个查找只需一次磁盘顺序读（通常在 Page Cache 中）
```

### 六大技术协同效应

```
写入一条消息:
  Producer batch → 压缩 → 一次网络请求发送一批
  Broker 收到 → 顺序追加写到 Partition log → write() 到 Page Cache
  → 返回成功（数据在内存，OS 异步刷盘）

消费一条消息:
  Consumer 请求 → Broker 调用 sendfile()
  → Page Cache 直接到网卡（零拷贝）→ 一批消息一起返回
  → Consumer 批量解压 → 处理

  写入: 批量 + 压缩 + 顺序写 + Page Cache → 几乎是纯内存顺序写速度
  读取: 零拷贝 + Page Cache + 批量 → 几乎是纯内存到网卡的速度
  扩展: 分区并行 × N → 线性扩展
```

### 对比总结

| 技术 | 解决的问题 | 效果 |
|------|-----------|------|
| **顺序 I/O** | 磁盘写入慢 | 随机写 ~1MB/s → 顺序写 ~600MB/s |
| **Page Cache** | 内存管理复杂 | 读写都走 OS 缓存，避免 GC |
| **零拷贝** | 数据搬运消耗 CPU | 4 次拷贝→2 次，CPU 零参与 |
| **批量+压缩** | 网络 RT 多、带宽占用大 | RT 减少 99%，带宽减少 75% |
| **分区并行** | 单机吞吐上限 | 线性水平扩展 |
| **日志分段+稀疏索引** | 查找效率 | O(1) 追加写 + O(logN) 查找 |

### 一句话总结

**Kafka 高吞吐 = 六大技术叠加：顺序 I/O（追加写磁盘比随机写快 600 倍）+ Page Cache（读写都走 OS 缓存不走 JVM）+ 零拷贝 sendfile（消费时数据从 Page Cache 直接到网卡，CPU 不参与）+ 批量+压缩（一次网络请求发几百条，压缩减少 75% 体积）+ 分区并行（N 个 Partition 线性扩展吞吐）+ 日志分段索引（O(1) 写 + O(logN) 查）。** 本质上 Kafka 把磁盘当成了"有持久化能力的内存"来用——顺序写入+Page Cache 让写入接近内存速度，零拷贝让读取绕过 CPU 直达网卡。

---

## 九、消息队列如何实现延迟消息

### 什么是延迟消息

```
延迟消息 = 发送后不立即投递给消费者，到达指定时间后才可被消费

场景:
  订单超时取消:  下单后发 30min 延迟消息 → 到期检查是否支付 → 未支付则取消
  会议提醒:      创建会议时发 15min 前延迟消息 → 到期推送提醒
  延迟重试:      消费失败后发 5s/30s/1min 延迟消息 → 延迟后重试
  定时任务:      定时发布活动 → 活动开始前 1 小时预热通知
```

### 各 MQ 延迟消息能力对比

| MQ | 原生支持 | 精度 | 延迟范围 | 实现方式 |
|----|---------|------|---------|---------|
| **RocketMQ 4.x** | 18 个固定级别 | 秒级 | 1s~2h | 内部延迟 Topic + 定时扫描 |
| **RocketMQ 5.x** | **任意时间** | 秒级 | 任意（最长 40 天） | 时间轮 + TimerLog |
| **RabbitMQ** | 插件/TTL+DLX | 毫秒级 | 任意 | 死信交换 / 延迟插件 |
| **Kafka** | **不支持** | — | — | 需应用层自行实现 |
| **Pulsar** | 原生支持 | 秒级 | 任意 | 延迟投递队列 |

---

### 一、RocketMQ 延迟消息原理

#### RocketMQ 4.x：18 个固定延迟级别

```
支持的延迟级别（不能自定义）:
  1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h

使用:
  msg.setDelayTimeLevel(3);  // 级别 3 = 10s 延迟

内部实现:
  ① 生产者发送消息，指定 delayTimeLevel=3
  ② Broker 收到后不投递到目标 Topic
     → 将消息存入内部 Topic: SCHEDULE_TOPIC_XXXX
     → 每个延迟级别一个 Queue（共 18 个 Queue）

  ③ Broker 内部有 18 个定时任务（每个级别一个）
     → 每秒扫描对应 Queue
     → 检查队首消息的"投递时间"是否到了
     → 到了 → 将消息从 SCHEDULE_TOPIC 转投到目标 Topic → 消费者可见

  SCHEDULE_TOPIC_XXXX:
    Queue 0 (1s延迟):   [msg_a(投递时间=10:00:01), msg_b(10:00:05)]
    Queue 1 (5s延迟):   [msg_c(投递时间=10:00:10)]
    Queue 2 (10s延迟):  [msg_d(投递时间=10:00:15)]
    ...
    Queue 17 (2h延迟):  [msg_e(投递时间=12:00:00)]

  每个 Queue 内消息按投递时间有序（同一级别延迟相同，先进先出）
  → 只需要检查队首 → O(1) 扫描效率
```

#### RocketMQ 5.x：任意延迟时间

```
5.x 新增: 指定任意时间戳延迟（秒级精度）

  msg.setDeliveryTimestamp(System.currentTimeMillis() + 90 * 1000); // 90 秒后

内部实现（时间轮 + TimerLog）:
  ① 消息先写入 TimerLog（持久化，防丢）
  ② 到期前 → 消息在时间轮中管理
  ③ 到期时 → 从 TimerLog 读出 → 转投到目标 Topic

  比 4.x 灵活: 不限于 18 个级别，任意秒级延迟
  但实现更复杂: 时间轮 + 持久化日志
```

---

### 二、RabbitMQ 延迟消息实现

#### 方案 A：TTL + 死信交换（DLX）

```
原理: 利用消息过期 + 死信机制

  ① 创建一个"等待队列"，设置消息 TTL（过期时间）
     该队列绑定一个死信交换（DLX）

  ② 消息发到"等待队列" → 没有消费者（故意的！）
     → 消息等 TTL 时间后过期
     → 过期后自动转发到 DLX
     → DLX 路由到"实际消费队列"
     → 消费者从"实际消费队列"消费

  流程:
  Producer → 等待队列(TTL=30min, 无消费者) → 30min 后过期
    → DLX → 实际消费队列 → Consumer

  ┌──────────────┐  TTL过期  ┌─────────┐  路由   ┌──────────────┐
  │  等待队列      │─────────→│  DLX    │────────→│  消费队列     │
  │ (TTL=30min)  │  (死信)   │(Exchange)│         │             │
  │ (无消费者)    │          └─────────┘         │  Consumer    │
  └──────────────┘                               └──────────────┘

问题: 队列级 TTL → 所有消息延迟一样
     消息级 TTL → 队首消息未过期会阻塞后面已过期的消息（不会跳过）

  例: 队列中 [msg1(TTL=30min), msg2(TTL=5s)]
      msg2 虽然 5s 就该过期，但 msg1 在前面挡着 → msg2 要等 30min
      → 严重问题！

解决: 每种延迟时间创建单独的等待队列
     5s 队列、30s 队列、1min 队列 ... → 队列数可能爆炸
```

#### 方案 B：rabbitmq-delayed-message-exchange 插件（推荐）

```
RabbitMQ 官方延迟消息插件:

  rabbitmq-plugins enable rabbitmq_delayed_message_exchange

  声明延迟交换:
  channel.exchangeDeclare("delayed_exchange", "x-delayed-message",
      true, false, Map.of("x-delayed-type", "direct"));

  发送延迟消息:
  AMQP.BasicProperties props = new AMQP.BasicProperties.Builder()
      .headers(Map.of("x-delay", 30000))  // 延迟 30 秒（毫秒）
      .build();
  channel.basicPublish("delayed_exchange", "routing_key", props, body);

内部实现:
  消息存在 Exchange 上（用 Mnesia 数据库）
  内部定时器到期后投递到绑定的 Queue
  → 不依赖 TTL/DLX → 不会有 TTL 队首阻塞问题 ✓

缺点:
  - 需要安装插件（部分云服务不支持）
  - 大量延迟消息时性能不如 RocketMQ（Mnesia 不适合高吞吐）
  - 集群中消息只存在一个节点上（该节点挂了消息丢失）
```

---

### 三、Kafka 延迟消息实现（应用层方案）

Kafka 不支持原生延迟消息，需要在应用层实现。

#### 方案 A：多级延迟 Topic（最常用）

```
按延迟时间创建多个 Topic:

  delay_topic_5s     → 延迟 5 秒的消息
  delay_topic_30s    → 延迟 30 秒的消息
  delay_topic_1min   → 延迟 1 分钟的消息
  delay_topic_5min   → 延迟 5 分钟的消息
  delay_topic_30min  → 延迟 30 分钟的消息

  生产者:
    按延迟时间选择 Topic → 消息写入时带上"目标 Topic"和"投递时间"

  延迟消费者（每个延迟 Topic 一个）:
    拉取消息 → 检查投递时间
    → 未到期: 暂停消费（pause partition），等一会再拉
    → 已到期: 转发到目标 Topic → 真正的消费者消费
```

```go
// Kafka 延迟消息 — 多级 Topic 方案
type DelayMessage struct {
    TargetTopic string `json:"target_topic"` // 最终目标 Topic
    DeliverAt   int64  `json:"deliver_at"`   // 投递时间戳
    Body        []byte `json:"body"`         // 原始消息体
}

// 生产者: 发到延迟 Topic
func SendDelayMessage(producer sarama.SyncProducer, targetTopic string, body []byte, delay time.Duration) error {
    delayTopic := selectDelayTopic(delay) // 选择对应的延迟 Topic
    msg := DelayMessage{
        TargetTopic: targetTopic,
        DeliverAt:   time.Now().Add(delay).Unix(),
        Body:        body,
    }
    data, _ := json.Marshal(msg)
    _, _, err := producer.SendMessage(&sarama.ProducerMessage{
        Topic: delayTopic,
        Value: sarama.ByteEncoder(data),
    })
    return err
}

func selectDelayTopic(delay time.Duration) string {
    switch {
    case delay <= 5*time.Second:
        return "delay_topic_5s"
    case delay <= 30*time.Second:
        return "delay_topic_30s"
    case delay <= time.Minute:
        return "delay_topic_1min"
    case delay <= 5*time.Minute:
        return "delay_topic_5min"
    default:
        return "delay_topic_30min"
    }
}

// 延迟转发消费者: 到期后转发到目标 Topic
func DelayForwarder(consumer sarama.ConsumerGroup, producer sarama.SyncProducer) {
    // 消费延迟 Topic
    for msg := range messages {
        var dm DelayMessage
        json.Unmarshal(msg.Value, &dm)

        now := time.Now().Unix()
        if dm.DeliverAt > now {
            // 还没到期 → 等待
            time.Sleep(time.Duration(dm.DeliverAt-now) * time.Second)
        }

        // 到期 → 转发到目标 Topic
        producer.SendMessage(&sarama.ProducerMessage{
            Topic: dm.TargetTopic,
            Value: sarama.ByteEncoder(dm.Body),
        })
        // 提交 offset
    }
}
```

#### 方案 B：外部存储 + 时间轮（大规模方案）

```
Kafka 配合 Redis ZSET 或 时间轮（参考 design-qa.md 第三题）:

  ① 消息发到 Kafka delay_topic
  ② 转存消费者: 从 delay_topic 拉取 → 写入 Redis ZSET（score=投递时间）
  ③ 定时轮询 Redis: ZPOPMIN 取到期消息 → 发到目标 Kafka Topic
  ④ 真正的消费者: 从目标 Topic 消费

  Kafka(暂存) → Redis ZSET(排序+到期检测) → Kafka(目标Topic) → Consumer

  优点: 支持任意延迟时间，精度高
  缺点: 多了 Redis 依赖，架构复杂
```

---

### 四、各方案对比

| 方案 | 延迟粒度 | 精度 | 复杂度 | 吞吐 | 可靠性 |
|------|---------|------|--------|------|--------|
| **RocketMQ 4.x** | 18 个级别 | 秒级 | 低（原生） | 高 | 高 |
| **RocketMQ 5.x** | 任意时间 | 秒级 | 低（原生） | 高 | 高 |
| RabbitMQ TTL+DLX | 队列级固定 | 毫秒级 | 中 | 中 | 中 |
| **RabbitMQ 插件** | 任意时间 | 毫秒级 | 低 | 中 | 中 |
| Kafka 多级 Topic | 几种级别 | 秒级 | 中 | 高 | 高 |
| Kafka + Redis ZSET | 任意时间 | 秒级 | 高 | 高 | 高 |

### 选型建议

```
已有 RocketMQ:
  4.x → 18 个固定级别够用（订单超时、延迟重试）
  5.x → 任意延迟时间（灵活度最高）
  → 零开发成本，原生支持

已有 RabbitMQ:
  → 安装 delayed-message-exchange 插件
  → 简单场景用 TTL+DLX（注意队首阻塞问题）

已有 Kafka:
  简单需求（3~5 种延迟时间）→ 多级 delay_topic
  复杂需求（任意延迟时间）→ Kafka + Redis ZSET
  → 都需要应用层开发

不依赖特定 MQ:
  → Redis ZSET + 时间轮（参考 design-qa.md 延时任务系统）
```

### 一句话总结

**RocketMQ 延迟消息原理：消息先存入内部 SCHEDULE_TOPIC，每个延迟级别一个 Queue，Broker 内部定时扫描队首消息的到期时间，到期后转投到目标 Topic（4.x 固定 18 级，5.x 时间轮支持任意延迟）。RabbitMQ 用 TTL+DLX 或延迟插件。Kafka 不支持原生延迟，需应用层实现——最常用多级 delay_topic 方案（按延迟时间路由到不同 Topic，延迟消费者到期后转发到目标 Topic），大规模用 Redis ZSET 辅助。** 选型口诀：有 RocketMQ 用原生（最省事），有 RabbitMQ 装插件，有 Kafka 多级 Topic 或 Redis ZSET。

---

## 十、Kafka Rebalance 机制

### 什么是 Rebalance

```
Rebalance = Consumer Group 内重新分配 Partition 和 Consumer 的映射关系

  Rebalance 前:
    P0,P1 → Consumer A
    P2,P3 → Consumer B

  Consumer C 加入后触发 Rebalance:
    P0 → Consumer A
    P1 → Consumer B
    P2,P3 → Consumer C（新分配）

  本质: "谁消费哪个 Partition" 的重新洗牌
```

### 什么时候触发 Rebalance

```
三种触发条件:

  ① Consumer Group 成员变化
     - 新消费者加入（扩容/新部署）
     - 消费者离开（缩容/下线）
     - 消费者被 Coordinator 判定为死亡（心跳超时）

  ② Topic 的 Partition 数变化
     - 增加 Partition（手动扩分区）

  ③ 订阅的 Topic 变化
     - Consumer 通过正则订阅 Topic（如 "order.*"），新 Topic 匹配上了
```

### Rebalance 的完整流程

```
角色:
  Coordinator: 管理 Consumer Group 的 Broker 节点
  Group Leader: Consumer Group 中的第一个消费者，负责制定分区分配方案

流程 (Eager Rebalance, Kafka 2.3 前):

  ① 触发 Rebalance
     Consumer C 加入 → 发送 JoinGroup 请求给 Coordinator

  ② JoinGroup 阶段
     Coordinator 通知所有 Consumer: "准备 Rebalance"
     → 所有 Consumer 停止消费，提交 offset，释放分区
     → 所有 Consumer 发送 JoinGroup 请求
     → Coordinator 选出 Group Leader（通常是第一个加入的）

  ③ SyncGroup 阶段
     Group Leader 根据分配策略计算新的分区分配方案
     → Leader 把方案通过 SyncGroup 请求发给 Coordinator
     → Coordinator 把方案分发给每个 Consumer

  ④ 恢复消费
     每个 Consumer 按新分配方案消费对应的 Partition

  时序:
  Consumer A  Consumer B  Consumer C    Coordinator
      │           │           │              │
      │           │           │─ JoinGroup ─→│  触发 Rebalance
      │           │           │              │
      │←── "Rebalance!" ──────────────────── │  通知所有 Consumer
      │           │←─────────────────────────│
      │ 停止消费   │ 停止消费                  │
      │ 提交offset │ 提交offset                │
      │           │                          │
      │─ JoinGroup ──────────────────────────│
      │           │─ JoinGroup ──────────────│
      │           │           │─ JoinGroup ─→│
      │           │           │              │  选出 Leader (A)
      │           │           │              │
      │←── 你是Leader ───────────────────────│
      │ 计算分配方案                           │
      │─ SyncGroup(方案) ───────────────────→│
      │           │─ SyncGroup ─────────────→│
      │           │           │─ SyncGroup ─→│
      │           │           │              │  分发方案
      │←── 你消费P0 ─────────────────────────│
      │           │←── 你消费P1,P2 ──────────│
      │           │           │←── 你消费P3 ─│
      │ 恢复消费   │ 恢复消费   │ 开始消费     │
```

### Rebalance 的问题

#### 问题一：Stop The World（全组停止消费）

```
Eager Rebalance 期间:
  所有 Consumer 都停止消费 → 直到 Rebalance 完成
  → 如果 Rebalance 耗时长 → 全组消费暂停 → 消息堆积

  典型耗时:
    3 个 Consumer → ~1-3s
    50 个 Consumer → ~10-30s
    100 个 Consumer → ~30-60s（极端情况更久）

  期间的影响:
    → 消费延迟飙升
    → 消息堆积
    → 上游如果有超时机制 → 可能触发告警
```

#### 问题二：重复消费

```
Consumer A 正在处理 P0 的 msg_100:
  → 触发 Rebalance → A 停止消费 → offset 可能还没提交
  → Rebalance 后 P0 分给 Consumer B
  → B 从上次提交的 offset（比如 msg_95）开始消费
  → msg_95 ~ msg_100 被重复消费！

  解决: 消费者必须做幂等处理（参考第六题）
```

#### 问题三：不必要的 Rebalance（误触发）

```
常见误触发场景:

  场景 1: 消费者处理时间过长
    Consumer 拉了一批消息 → 处理耗时 > max.poll.interval.ms（默认 5min）
    → Coordinator 认为消费者死了 → 踢出 → 触发 Rebalance
    → 处理完后消费者发现自己被踢了 → 重新加入 → 又触发 Rebalance

  场景 2: GC 导致心跳超时
    JVM Full GC 暂停 > session.timeout.ms（默认 45s）
    → Coordinator 收不到心跳 → 认为死了 → Rebalance

  场景 3: 网络抖动
    Consumer 和 Coordinator 之间网络短暂中断
    → 心跳超时 → Rebalance

  场景 4: 滚动部署
    逐台重启消费者 → 每次重启都触发 Rebalance
    10 台消费者 → 重启 10 次 → 触发 10 次 Rebalance
```

### 分配策略

```
Consumer Group 的 Partition 分配有三种策略:

  1. Range（默认）:
     按 Topic 维度，按序号范围分配
     P0,P1 → C0;  P2,P3 → C1;  P4,P5 → C2
     问题: 多 Topic 时第一个 Consumer 总是多分一个 → 不均匀

  2. RoundRobin:
     所有 Topic 的 Partition 混在一起轮询分配
     P0→C0, P1→C1, P2→C2, P3→C0, P4→C1, P5→C2
     更均匀，但 Rebalance 时变动大

  3. Sticky（推荐）:
     尽量保持原有的分配不变，只移动最少的 Partition
     → Rebalance 时大部分 Consumer 不受影响
     → 减少不必要的分区迁移 → 恢复更快
```

### 解决方案

#### 1. Cooperative Incremental Rebalance（Kafka 2.4+）

```
Eager (旧): 全部释放 → 全部重新分配 → 全组停止消费
Cooperative (新): 只迁移需要变动的 Partition → 其他 Consumer 不中断

  Rebalance 前:
    C0: P0, P1, P2
    C1: P3, P4, P5

  新增 C2，Cooperative Rebalance:
    第一轮: 确定哪些 Partition 需要迁移
      C0: 释放 P2（只释放需要迁移的）
      C1: 释放 P5
      C0 继续消费 P0, P1 ← 不中断!
      C1 继续消费 P3, P4 ← 不中断!

    第二轮: 将释放的 Partition 分配给新 Consumer
      C2: 获得 P2, P5

  结果:
    C0: P0, P1      ← P0, P1 从头到尾没有中断消费
    C1: P3, P4      ← P3, P4 也没有中断
    C2: P2, P5      ← 只有被迁移的 P2, P5 短暂中断

  → 大幅减少 Stop The World 的范围

配置:
  partition.assignment.strategy = org.apache.kafka.clients.consumer.CooperativeStickyAssignor
```

#### 2. 避免误触发 Rebalance

```
关键配置调优:

  session.timeout.ms = 45000 (默认)
    → Consumer 心跳超时时间
    → 太小: 网络抖动/GC 就触发 Rebalance
    → 太大: 真正宕机时检测慢
    → 推荐: 25000~45000

  heartbeat.interval.ms = 3000 (默认)
    → 心跳发送间隔
    → 必须 < session.timeout.ms / 3
    → 推荐: session.timeout.ms 的 1/3

  max.poll.interval.ms = 300000 (默认 5min)
    → 两次 poll() 的最大间隔
    → 消费处理时间不能超过这个值
    → 如果单批处理慢 → 调大这个值 或 减小 max.poll.records
    → 推荐: 根据实际处理时间设置（留 2~3 倍余量）

  max.poll.records = 500 (默认)
    → 单次 poll 拉取的最大消息数
    → 减小 → 每批处理更快 → 不容易超 max.poll.interval.ms

避免滚动部署频繁触发:
  group.initial.rebalance.delay.ms = 3000 (默认)
    → Group 启动后等待一段时间再 Rebalance（等所有消费者都加入）
    → 滚动部署时适当调大（如 10000）
```

#### 3. Static Membership（Kafka 2.3+）

```
问题: Consumer 重启 → 被认为是新成员 → 触发 Rebalance
     实际上只是重启，分区分配不用变

Static Membership:
  每个 Consumer 配置固定的 group.instance.id
  → 重启后用相同的 ID 重新加入
  → Coordinator 识别为"回来了"而不是"新来的"
  → 在 session.timeout.ms 内回来 → 不触发 Rebalance!

  group.instance.id = "consumer-host-1"  // 每台机器固定一个 ID

  滚动部署时:
    Consumer 1 停止 → Coordinator 等待 session.timeout.ms
    Consumer 1 重启 → 用同一个 instance.id 加入 → 恢复原来的 Partition
    → 整个过程不触发 Rebalance!（只要在超时内回来）
```

### Rebalance 问题排查

```sql
-- 查看 Consumer Group 状态
kafka-consumer-groups.sh --bootstrap-server localhost:9092 \
    --describe --group my_group

-- 查看 Rebalance 日志
grep "Rebalance" /var/log/kafka/server.log

-- 关键指标监控
  rebalance-latency-avg:   Rebalance 平均耗时
  rebalance-rate:           Rebalance 频率（过高说明有问题）
  last-rebalance-seconds-ago: 距上次 Rebalance 的时间
  records-lag:              消费延迟（Rebalance 时会飙升）
```

### 一句话总结

**Kafka Rebalance 在消费者加入/退出/超时、Partition 数变化时触发，全组重新分配"谁消费哪个 Partition"。核心问题：Eager Rebalance 会导致全组停止消费（STW） + 可能重复消费 + 频繁误触发。** 解决方案三板斧：①用 Cooperative Incremental Rebalance（Kafka 2.4+，只迁移变动的 Partition，其他不中断）；②调优参数避免误触发（增大 session.timeout.ms、减小 max.poll.records、合理设置 max.poll.interval.ms）；③用 Static Membership（固定 group.instance.id，重启不触发 Rebalance）。分配策略推荐 CooperativeStickyAssignor（尽量保持原分配不变）。

---

## 十一、消息队列的推模式 vs 拉模式

### 基本概念

```
Push（推模式）: Broker 主动将消息推给消费者
  Broker: "有新消息了，给你！"
  消费者: 被动接收

Pull（拉模式）: 消费者主动从 Broker 拉取消息
  消费者: "有新消息吗？给我 100 条"
  Broker: 被动响应

  ┌─────────┐   Push    ┌──────────┐
  │ Broker  │──────────→│ Consumer │  Broker 主动推
  └─────────┘           └──────────┘

  ┌─────────┐   Pull    ┌──────────┐
  │ Broker  │←──────────│ Consumer │  Consumer 主动拉
  └─────────┘           └──────────┘
```

### 核心对比

| 维度 | Push（推） | Pull（拉） |
|------|-----------|-----------|
| **主动方** | Broker | Consumer |
| **实时性** | 高（有消息立即推） | 取决于轮询间隔 |
| **消费速率控制** | **Broker 控制**（消费者被动） | **Consumer 控制**（自主调节） |
| **消费者压力** | 可能被压垮（Broker 推太快） | 不会（自己控制拉取速度） |
| **Broker 复杂度** | 高（要维护每个消费者的状态和推送速率） | 低（只需响应拉取请求） |
| **空闲时开销** | 无（没消息就不推） | 有（空轮询浪费资源） |
| **消息堆积感知** | Broker 清楚堆积情况 | Consumer 通过 lag 间接感知 |
| **批量能力** | 弱（来一条推一条） | 强（一次拉一批） |

### Push 模式的问题

```
核心问题: Broker 不知道消费者的处理能力 → 推太快消费者扛不住

  场景:
    Broker 有 10000 msg/s 的消息
    Consumer 处理能力只有 1000 msg/s
    → Broker 以 10000/s 推 → Consumer 处理不过来
    → 内存溢出 / 消息丢弃 / 拒绝连接

  解决（Push 模式内部优化）:
    限流: Consumer 告诉 Broker "我最多接受 N 条/秒"
    背压(Backpressure): Consumer 处理不过来时发信号让 Broker 减速
    预取(Prefetch): 限制未 ACK 的消息数量（RabbitMQ 的 prefetch_count）

  ┌─────────┐  推 10000/s  ┌──────────┐
  │ Broker  │─────────────→│ Consumer │ 只能处理 1000/s → 压垮!
  └─────────┘              └──────────┘

  RabbitMQ 的解决:
    channel.basicQos(prefetchCount=100);
    → 未 ACK 的消息 ≤ 100 条 → 超过就不推了
    → 实际上是"变相的 Pull"（用 ACK 反馈控制推送速度）
```

### Pull 模式的问题

```
核心问题: 消费者不知道什么时候有新消息 → 空轮询浪费资源

  场景:
    Consumer 每 100ms poll 一次
    如果 Topic 很久没有新消息
    → 每次 poll 返回空 → CPU 和网络白白浪费

  解决:
    长轮询 (Long Polling):
      Consumer: poll(timeout=500ms)
      Broker: 如果没有消息 → 不立即返回空 → 等最多 500ms
             → 等待期间有新消息 → 立即返回
             → 超时也没消息 → 返回空

      ┌──────────────────────────────────────┐
      │  短轮询:                              │
      │    C: poll → B: 无消息，返回空 (立即)   │
      │    C: poll → B: 无消息，返回空         │
      │    C: poll → B: 有消息，返回数据       │
      │    → 前两次白跑了                      │
      │                                      │
      │  长轮询 (Kafka):                      │
      │    C: poll(500ms) → B: 无消息，等...   │
      │    → 200ms 后有消息 → 立即返回数据     │
      │    → 只等了 200ms，没有空轮询          │
      └──────────────────────────────────────┘

  Kafka 的长轮询参数:
    fetch.min.bytes = 1         // 最少返回 1 字节数据（有数据就返回）
    fetch.max.wait.ms = 500     // 最多等 500ms（没数据到超时也返回）
    → 兼顾实时性和资源消耗
```

### Kafka 为什么选 Pull

```
Kafka 用 Pull 的五个原因:

  ① 消费速率由消费者控制（最核心）
     每个消费者处理能力不同
     Pull: 能力强的多拉，能力弱的少拉 → 自适应
     Push: Broker 无法知道每个消费者的处理能力 → 可能推太快压垮

  ② 天然支持批量
     Consumer 一次 poll 拉取一批（几十~几百条）
     → 批量处理效率远高于逐条推送
     → 配合 Kafka 的批量+压缩 → 吞吐量极高

  ③ Broker 设计简单
     Push: Broker 要维护每个消费者的状态、速率、推送位置
     Pull: Broker 只需要存消息 + 响应拉取请求 → 无状态（offset 由消费者管理）
     → Broker 专注于存储和复制 → 更高效

  ④ 消息回溯能力
     Pull: Consumer 自己管理 offset → 可以重置到任意位置重新消费
     Push: Broker 推了就推了 → 回溯很难

  ⑤ 长轮询解决了实时性问题
     Pull 的唯一缺点是"不够实时"
     → Kafka 用长轮询（fetch.max.wait.ms）解决
     → 实时性接近 Push（几十~几百 ms 延迟）

  Pull 的缺点:
     消费者需要自己管理 offset → 增加客户端复杂度
     → Kafka 通过 __consumer_offsets 内部 Topic 帮消费者管理
```

### 各 MQ 的模式选择

| MQ | 模式 | 原因 |
|----|------|------|
| **Kafka** | **Pull** | 大数据高吞吐 → 消费者控速 + 批量拉取 |
| **RocketMQ** | **Pull + 伪 Push** | 底层 Pull + 长轮询，封装为 Push 接口 |
| **RabbitMQ** | **Push** | 低延迟场景 → 有消息立即推 + prefetch 限流 |
| **Pulsar** | **Pull** | 类似 Kafka |
| **Redis Stream** | **Pull** | XREADGROUP BLOCK 长轮询 |

```
RocketMQ 的"伪 Push":
  对外接口看起来是 Push（注册 MessageListener，有消息就回调）
  内部实现其实是 Pull（长轮询拉取）+ 回调
  → 兼顾了 Push 的简单 API 和 Pull 的控速能力

  PushConsumer (表面):
    consumer.registerMessageListener(msg -> { process(msg); })
    // 看起来像回调（Push）

  内部实现:
    while(true) {
        msgs = broker.pull(timeout=30s);  // 实际是长轮询 Pull
        for (msg : msgs) {
            listener.consumeMessage(msg);  // 拉到后回调用户代码
        }
    }
```

### 如何选择

```
选 Push:
  ✓ 消息量小、消费者少
  ✓ 需要极低延迟（有消息立即处理）
  ✓ 消费者处理能力固定且 Broker 可预估
  → RabbitMQ 的典型场景

选 Pull:
  ✓ 消息量大、消费者多
  ✓ 消费者处理能力差异大
  ✓ 需要批量处理提高吞吐
  ✓ 需要消息回溯能力
  → Kafka / RocketMQ 的典型场景

大多数场景: Pull + 长轮询 是最佳平衡
  → 兼顾实时性（长轮询）+ 控速（消费者自主）+ 批量（高吞吐）
```

### 一句话总结

**Push 由 Broker 主动推消息（实时性高但可能压垮消费者），Pull 由消费者主动拉消息（自控速率但有空轮询问题）。Kafka 选 Pull 因为五个原因：消费者自主控速不会被压垮、天然支持批量拉取高吞吐、Broker 设计简单无状态、支持消息回溯、长轮询弥补了实时性不足。** RocketMQ 是 "伪 Push"——对外是回调 API，内部是 Pull 长轮询。RabbitMQ 是真 Push——靠 prefetch_count 限流防压垮。大多数高吞吐场景 Pull + 长轮询是最佳选择。