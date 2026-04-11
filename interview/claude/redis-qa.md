# Redis Q&A

## 目录

- [一、缓存穿透、击穿、雪崩及缓存一致性](#一缓存穿透击穿雪崩及缓存一致性)
- [二、Redis 数据类型及使用场景](#二redis-数据类型及使用场景)
- [三、Redis 持久化机制（RDB vs AOF）](#三redis-持久化机制rdb-vs-aof)
- [四、Redis 过期策略与内存淘汰机制](#四redis-过期策略与内存淘汰机制)
- [五、Redis 单线程为什么这么快](#五redis-单线程为什么这么快)
- [六、Redis 分布式锁](#六redis-分布式锁)
- [七、Redis Cluster 集群原理](#七redis-cluster-集群原理)
- [八、Redis 主从复制与哨兵机制](#八redis-主从复制与哨兵机制)
- [九、Redis 大 Key 和热 Key 问题](#九redis-大-key-和热-key-问题)
- [十、Redis 事务与 Pipeline](#十redis-事务与-pipeline)
- [十一、Redis/MySQL/etcd 集群数据同步机制与 CAP 分析](#十一redismysqletcd-集群数据同步机制与-cap-分析)
- [十二、ZSet 实现原理：跳表为什么优于红黑树](#十二zset-实现原理跳表为什么优于红黑树)
- [十三、Redis Cluster MGET 原理与 Slot 迁移](#十三redis-cluster-mget-原理与-slot-迁移)

---

## 一、缓存穿透、击穿、雪崩及缓存一致性

### 缓存的基本读写模式

```
读请求:
  ① 先查 Redis 缓存
  ② 命中 → 直接返回（快，~1ms）
  ③ 未命中 → 查 DB → 写入 Redis → 返回

写请求:
  ④ 更新 DB → 删除/更新 Redis 缓存

  Client → Redis ──命中──→ 返回
             │ miss
             ▼
            DB → 写 Redis → 返回
```

---

### 一、缓存穿透（查不存在的数据）

```
定义: 请求的数据在缓存和数据库中都不存在
     → 每次请求都穿过缓存打到 DB → 缓存形同虚设

场景:
  恶意攻击: 构造大量不存在的 ID（如 id=-1, id=999999999）
  业务 bug: 前端传了非法参数

  请求 id=999999999:
    Redis: 没有 → 查 DB → DB 也没有 → 不写缓存
    下次同样请求 → 又穿到 DB → 反复打 DB
    → 如果是攻击（10 万 QPS）→ DB 被打挂
```

#### 解决方案

**方案 1: 缓存空值**

```go
func GetUser(ctx context.Context, id int64) (*User, error) {
    key := fmt.Sprintf("user:%d", id)

    // ① 查缓存
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        if val == "NULL" {
            return nil, nil // 缓存了空值 → 不查 DB
        }
        var user User
        json.Unmarshal([]byte(val), &user)
        return &user, nil
    }

    // ② 查 DB
    user, err := queryUserFromDB(id)
    if err != nil {
        return nil, err
    }

    if user == nil {
        // DB 也没有 → 缓存空值（短 TTL 防止永久占用）
        rdb.Set(ctx, key, "NULL", 5*time.Minute)
        return nil, nil
    }

    // ③ 缓存真实数据
    data, _ := json.Marshal(user)
    rdb.Set(ctx, key, data, 30*time.Minute)
    return user, nil
}
```

```
优点: 简单有效
缺点: 大量不存在的 key → Redis 内存被空值占满
     → 解决: 空值设短 TTL（5min）
```

**方案 2: 布隆过滤器（大规模防穿透）**

```
布隆过滤器: 用极小的内存判断"一个元素一定不存在"

  特点:
    说"不存在" → 一定不存在 ✓
    说"存在"   → 可能存在（有误判率，如 0.1%）

  流程:
    请求 id → 先查布隆过滤器
    → 不存在: 直接返回 404（不查缓存也不查 DB）
    → 可能存在: 继续查缓存 → DB

  ┌────────┐  不存在   ┌────────┐
  │ 请求    │──────────→│ 返回404 │
  │ id=999 │          └────────┘
  │        │  可能存在
  │        │──────────→ 查 Redis → 查 DB
  └────────┘

  内存占用: 1 亿 key 的布隆过滤器 ≈ 120MB（误判率 0.1%）
  → 比缓存空值省内存得多
```

```go
// Redis 布隆过滤器（需要 RedisBloom 模块）
// BF.ADD user_filter 12345    → 添加已存在的 ID
// BF.EXISTS user_filter 99999 → 0 (不存在)

func GetUserWithBloom(ctx context.Context, id int64) (*User, error) {
    // ① 布隆过滤器拦截
    exists, _ := rdb.Do(ctx, "BF.EXISTS", "user_filter", id).Bool()
    if !exists {
        return nil, nil // 一定不存在，直接返回
    }

    // ② 查缓存 → 查 DB（正常流程）
    return getFromCacheOrDB(ctx, id)
}
```

**方案 3: 参数校验（最前置的拦截）**

```
在接入层直接拦截非法参数:
  id <= 0 → 拒绝
  id 格式不对 → 拒绝
  id 不在合法范围内 → 拒绝
  → 最简单但不能解决所有情况（合法 ID 但不存在的数据）
```

---

### 二、缓存击穿（热点 key 过期）

```
定义: 某个热点 key 过期的瞬间，大量请求同时打到 DB

  热门商品详情页: 10 万 QPS 都命中缓存
  → 缓存过期的瞬间
  → 10 万请求同时 miss → 同时查 DB
  → DB 瞬间承受 10 万并发 → 可能被打挂

  时间线:
  ─── 缓存命中(10万QPS) ───│── 缓存过期 ──│── DB 10万并发 ──│
                           ↑
                    这一瞬间所有请求打到 DB
```

#### 解决方案

**方案 1: 互斥锁（singleflight 模式）**

```go
import "golang.org/x/sync/singleflight"

var sfGroup singleflight.Group

func GetProduct(ctx context.Context, id string) (*Product, error) {
    key := "product:" + id

    // ① 查缓存
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        var p Product
        json.Unmarshal([]byte(val), &p)
        return &p, nil
    }

    // ② singleflight: 多个相同请求只有一个真正查 DB
    result, err, _ := sfGroup.Do(key, func() (interface{}, error) {
        // 这个函数只会被执行一次（其他相同 key 的请求等待结果）
        product, err := queryProductFromDB(id)
        if err != nil {
            return nil, err
        }
        // 写入缓存
        data, _ := json.Marshal(product)
        rdb.Set(ctx, key, data, 30*time.Minute)
        return product, nil
    })
    if err != nil {
        return nil, err
    }
    return result.(*Product), nil
}

// 效果:
// 10 万请求同时 miss → 只有 1 个查 DB → 写入缓存 → 其余 99999 个拿到结果
// DB 只承受 1 次查询！
```

**方案 2: 热点 key 永不过期 + 异步刷新**

```go
type CacheValue struct {
    Data      []byte    `json:"data"`
    LogicTTL  time.Time `json:"logic_ttl"` // 逻辑过期时间
}

// Redis 不设 TTL（永不过期），用逻辑过期时间控制刷新
func GetProductNeverExpire(ctx context.Context, id string) (*Product, error) {
    key := "product:" + id
    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        return nil, err // key 不存在（需要预热）
    }

    var cv CacheValue
    json.Unmarshal([]byte(val), &cv)

    if time.Now().Before(cv.LogicTTL) {
        // 未过期，直接返回
        var p Product
        json.Unmarshal(cv.Data, &p)
        return &p, nil
    }

    // 逻辑过期 → 异步刷新（当前请求仍返回旧数据）
    go func() {
        // singleflight 防止并发刷新
        sfGroup.Do("refresh:"+key, func() (interface{}, error) {
            product, _ := queryProductFromDB(id)
            data, _ := json.Marshal(product)
            newCV := CacheValue{
                Data:     data,
                LogicTTL: time.Now().Add(30 * time.Minute),
            }
            cvData, _ := json.Marshal(newCV)
            rdb.Set(ctx, key, cvData, 0) // 不设 Redis TTL
            return nil, nil
        })
    }()

    // 返回旧数据（不等刷新完成）
    var p Product
    json.Unmarshal(cv.Data, &p)
    return &p, nil
}

// 优点: 用户永远不会等待 → 极低延迟
// 缺点: 逻辑过期到刷新完成之间可能返回旧数据（短暂不一致）
```

---

### 三、缓存雪崩（大量 key 同时过期）

```
定义: 大量 key 在同一时间过期 → 大量请求同时打到 DB

  与击穿的区别:
    击穿: 一个热点 key 过期 → 一个点的洪水
    雪崩: 大量 key 同时过期 → 全面崩溃

  场景:
    缓存预热时 TTL 都设了 30min → 30 分钟后全部同时过期
    Redis 宕机 → 所有缓存失效 → 全部打 DB
```

#### 解决方案

```
方案 1: TTL 加随机值（防止同时过期）
  ttl := baseTTL + rand.Intn(300)  // 基础 30min + 随机 0~5min
  rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second)
  → 过期时间分散在 30~35min → 不会同时过期

方案 2: Redis 高可用（防止宕机导致雪崩）
  → Redis Cluster（3 主 3 从）
  → 哨兵模式自动主从切换
  → 多机房部署

方案 3: 多级缓存
  → 本地缓存(进程内) + Redis + DB
  → Redis 挂了还有本地缓存兜底

  请求 → 本地缓存(bigcache/freecache) → Redis → DB
          命中率 ~30%                命中率 ~95%

方案 4: 限流降级（兜底）
  → 缓存失效时对 DB 查询限流（令牌桶/信号量）
  → DB 撑不住时降级返回默认值/缓存旧数据
```

```go
// 多级缓存 + TTL 随机 + singleflight 完整方案
func GetData(ctx context.Context, key string) ([]byte, error) {
    // ① 本地缓存（进程内，~100ns）
    if val, ok := localCache.Get(key); ok {
        return val, nil
    }

    // ② Redis 缓存（~1ms）
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        localCache.Set(key, []byte(val), 1*time.Minute) // 回写本地
        return []byte(val), nil
    }

    // ③ singleflight 防击穿 + 查 DB
    result, err, _ := sfGroup.Do(key, func() (interface{}, error) {
        data, err := queryDB(key)
        if err != nil {
            return nil, err
        }
        // 写 Redis（TTL 加随机防雪崩）
        ttl := 30*time.Minute + time.Duration(rand.Intn(300))*time.Second
        rdb.Set(ctx, key, data, ttl)
        // 写本地缓存
        localCache.Set(key, data, 1*time.Minute)
        return data, nil
    })
    if err != nil {
        return nil, err
    }
    return result.([]byte), nil
}
```

---

### 四、三种问题对比

| 维度 | 缓存穿透 | 缓存击穿 | 缓存雪崩 |
|------|---------|---------|---------|
| **原因** | 查不存在的数据 | 热点 key 过期 | 大量 key 同时过期/Redis 宕机 |
| **影响范围** | 针对特定不存在的 key | 针对单个热点 key | 全局大面积缓存失效 |
| **DB 压力** | 持续（每次查不到都打 DB） | 瞬间洪峰 | 全面洪峰 |
| **核心解决** | 缓存空值 + 布隆过滤器 | singleflight + 永不过期 | TTL 随机 + 高可用 + 多级缓存 |

---

### 五、缓存一致性（最高频的面试题）

```
核心问题: DB 和缓存是两个独立的存储，如何保证数据一致？

  写操作时:
    先更新 DB 还是先更新缓存？
    更新缓存还是删除缓存？
    → 4 种组合，各有问题
```

#### 四种策略对比

| 策略 | 一致性问题 | 推荐 |
|------|-----------|------|
| 先更新 DB，再更新缓存 | 并发写导致缓存是旧值 | 不推荐 |
| 先更新缓存，再更新 DB | DB 更新失败但缓存已改 | 不推荐 |
| 先更新 DB，再删除缓存 | 极小概率不一致（可接受） | **推荐（Cache Aside）** |
| 先删除缓存，再更新 DB | 并发读写导致缓存是旧值 | 不推荐（需延迟双删） |

#### Cache Aside（旁路缓存）— 推荐方案

```
读: 查缓存 → miss → 查 DB → 写缓存
写: 更新 DB → 删除缓存（不是更新缓存！）

为什么删缓存而不是更新缓存？
  场景: 两个并发写
    A: UPDATE price=100 → SET cache=100
    B: UPDATE price=200 → SET cache=200
    如果执行顺序: A写DB → B写DB → B更新缓存 → A更新缓存
    → DB=200, 缓存=100（不一致！A 的旧值覆盖了 B）

  删缓存就没这个问题:
    A写DB → B写DB → B删缓存 → A删缓存
    → 缓存为空 → 下次读会查 DB 拿到最新值 200 ✓
```

```
"先更新 DB 再删缓存"仍有极小概率不一致:

  时间线:
    ① 读请求: 查缓存 miss → 查 DB 得到旧值 old
    ② 写请求: 更新 DB 为 new → 删除缓存
    ③ 读请求: 将旧值 old 写入缓存
    → 缓存中是旧值 old！

  这个概率极低，因为需要:
    读 DB（慢，~5ms）和写 DB + 删缓存 的执行窗口恰好交错
    且读 DB 必须比写 DB 更早开始但更晚结束
    → 实际中几乎不会发生

  如果要解决:
    → 给缓存设 TTL（几分钟后自动过期修正）
    → 延迟双删（写完后延迟 500ms 再删一次）
```

#### 延迟双删

```go
func UpdateWithDoubleDelete(ctx context.Context, key string, value interface{}) error {
    // ① 先删缓存
    rdb.Del(ctx, key)

    // ② 更新 DB
    if err := updateDB(value); err != nil {
        return err
    }

    // ③ 延迟再删一次（消除读写并发的不一致窗口）
    go func() {
        time.Sleep(500 * time.Millisecond) // 等读请求完成
        rdb.Del(context.Background(), key)
    }()

    return nil
}

// 延迟时间: 略大于一次读 DB + 写缓存的耗时（通常 200ms~1s）
```

#### 基于 Binlog 的异步删除（大规模方案）

```
问题: 应用层删缓存可能失败（网络异常/应用崩溃）→ 缓存永远是旧值

终极方案: 监听 MySQL Binlog，DB 变更时自动删缓存

  应用 → 只写 DB（不操作缓存）
  Canal/Debezium → 监听 Binlog → 解析变更 → 删 Redis 缓存

  ┌────────┐       ┌──────┐      ┌───────┐      ┌───────┐
  │  应用   │──写──→│ MySQL │──Binlog→│ Canal │──删──→│ Redis │
  └────────┘       └──────┘      └───────┘      └───────┘

  优点:
    - 应用层不需要操心缓存删除
    - Binlog 是可靠的（不会丢）
    - 即使 Canal 暂时不可用，恢复后也能追上
    - 解耦: 加新缓存只需加 Canal 规则

  缺点:
    - 架构复杂（多了 Canal 组件）
    - 有延迟（从 Binlog 到删缓存通常 100ms~1s）
```

#### 缓存一致性方案选型

```
简单场景（95% 的业务）:
  Cache Aside: 更新 DB → 删缓存 + TTL 兜底
  → 简单可靠，极小概率短暂不一致（TTL 内自动修正）

一致性要求高:
  Cache Aside + 延迟双删
  → 消除并发读写的不一致窗口

大规模 + 强一致:
  Binlog 异步删除（Canal/Debezium）
  → 自动化、可靠、解耦

实时性要求极高:
  Write Through（写穿透: 写 DB + 写缓存在同一事务/框架中）
  → 但实现复杂，一般框架层面做（如 ORM 插件）
```

---

### 六、常见面试追问

```
Q: 缓存和 DB 的双写怎么保证原子性？
A: 严格来说不能保证原子性（两个独立存储）
   → Cache Aside + TTL 保证最终一致性
   → 对强一致要求的场景，用 Binlog 订阅方案
   → 或 Write Behind（异步写回），但实现复杂

Q: 为什么不先删缓存再更新 DB？
A: 并发下会导致不一致:
   ① A 删缓存
   ② B 读缓存 miss → 查 DB（还是旧值）→ 写缓存（旧值）
   ③ A 更新 DB
   → 缓存是旧值，DB 是新值 → 不一致

Q: 缓存设多长 TTL 合适？
A: 取决于业务对数据实时性的要求:
   实时性高（商品价格）: 5~10min + 变更时主动删
   实时性低（用户头像）: 30min~2h
   极少变化（配置）:     24h+

Q: Redis 和 DB 数据不一致了怎么排查？
A: ① 开启 Redis 操作日志（MONITOR / slowlog）
   ② 检查是否有漏删缓存的代码路径
   ③ 定时对账: 后台任务比对 Redis 和 DB 数据，不一致则修正
```

### 一句话总结

**三种缓存问题：穿透（查不存在的数据→缓存空值+布隆过滤器）、击穿（热点 key 过期→singleflight 只放一个请求查 DB）、雪崩（大量 key 同时过期→TTL 加随机值+多级缓存+Redis 高可用）。缓存一致性推荐 Cache Aside 模式：读时查缓存→miss 查 DB→写缓存；写时更新 DB→删缓存（不是更新缓存！）+TTL 兜底。** 删缓存比更新缓存安全（避免并发写覆盖）。极小概率不一致用延迟双删解决。大规模系统用 Binlog（Canal）监听 DB 变更自动删缓存——最可靠的最终一致性方案。

---

## 二、Redis 数据类型及使用场景

### 五种基础数据类型

| 类型 | 底层结构 | 典型场景 | 命令示例 |
|------|---------|---------|---------|
| **String** | SDS（简单动态字符串） | 缓存、计数器、分布式锁、Session | `SET/GET/INCR/SETNX` |
| **Hash** | ziplist / hashtable | 对象存储（用户信息）、购物车 | `HSET/HGET/HGETALL` |
| **List** | quicklist（ziplist + 双向链表） | 消息队列、最新列表、时间线 | `LPUSH/RPOP/LRANGE` |
| **Set** | intset / hashtable | 去重、交集（共同好友）、抽奖 | `SADD/SMEMBERS/SINTER` |
| **ZSet（Sorted Set）** | ziplist / skiplist + hashtable | 排行榜、延迟队列、滑动窗口限流 | `ZADD/ZRANGE/ZRANGEBYSCORE` |

### 三种特殊数据类型

| 类型 | 用途 | 典型场景 |
|------|------|---------|
| **HyperLogLog** | 基数估算（误差 0.81%） | UV 统计（日活/月活），内存仅 12KB |
| **Bitmap** | 位图操作 | 签到打卡、用户在线状态、布隆过滤器 |
| **GEO** | 地理位置 | 附近的人、门店搜索 |

### 底层数据结构演进（Redis 7.0+）

```
String → SDS
Hash   → listpack（≤ 128个元素且值≤ 64字节） / hashtable
List   → quicklist（listpack + 双向链表）
Set    → intset（全为整数且≤ 512个）/ hashtable
ZSet   → listpack（≤ 128个元素）/ skiplist + hashtable

注: Redis 7.0 用 listpack 替代了 ziplist，解决了连锁更新问题
```

### 使用场景详解

```
1. String - 计数器
   INCR article:1001:views     → 文章阅读量 +1（原子操作）
   SETNX lock:order:123 uuid   → 分布式锁（SET if Not eXists）

2. Hash - 对象存储
   HSET user:1001 name "张三" age 25 city "北京"
   → 比 String 存 JSON 更省内存，支持部分字段更新

3. List - 消息队列（简易版）
   LPUSH queue:order msg1       → 生产者 push
   BRPOP queue:order 30         → 消费者阻塞 pop（30s 超时）

4. Set - 共同好友
   SADD user:1001:friends A B C
   SADD user:1002:friends B C D
   SINTER user:1001:friends user:1002:friends → {B, C}

5. ZSet - 排行榜
   ZADD rank:game 1000 player_A 2000 player_B
   ZREVRANGE rank:game 0 9 WITHSCORES → Top 10
```

---

## 三、Redis 持久化机制（RDB vs AOF）

### RDB（Redis Database）快照

```
原理: 按指定间隔将内存数据生成快照保存到 dump.rdb 文件

触发方式:
  ① 配置自动触发: save 900 1   → 900 秒内至少 1 次写操作
                  save 300 10  → 300 秒内至少 10 次写操作
  ② 手动触发:
     SAVE  → 阻塞主线程（生产环境禁用！）
     BGSAVE → fork 子进程生成快照（推荐）

fork + COW（Copy-On-Write）机制:
  ① 主进程 fork 出子进程（共享内存页）
  ② 子进程遍历内存写入 RDB 文件
  ③ 主进程继续处理请求
  ④ 写操作触发 COW: 修改的页面复制一份给子进程
  → 子进程看到的是 fork 时刻的一致性快照

  ┌──────────┐ fork  ┌──────────┐
  │ 主进程    │──────→│ 子进程    │
  │ 继续读写  │       │ 写 RDB   │
  │ COW 复制  │       │ (一致快照)│
  └──────────┘       └──────────┘
```

### AOF（Append Only File）追加日志

```
原理: 每次写操作追加到 appendonly.aof 文件

同步策略:
  appendfsync always    → 每条命令都 fsync（最安全，最慢）
  appendfsync everysec  → 每秒 fsync 一次（推荐，最多丢 1 秒）
  appendfsync no        → 由 OS 决定何时 fsync（最快，可能丢数据）

AOF 重写（Rewrite）:
  问题: AOF 文件越来越大（SET k1 v1 → SET k1 v2 → 两条都记录）
  重写: BGREWRITEAOF → fork 子进程 → 遍历内存生成最精简的命令集
  → SET k1 v2（只保留最终状态）
  → 文件从 10GB 缩小到 1GB
```

### RDB vs AOF 对比

| 维度 | RDB | AOF |
|------|-----|-----|
| **数据安全** | 可能丢失最后一次快照后的数据 | 最多丢 1 秒（everysec） |
| **文件大小** | 小（二进制压缩） | 大（文本命令日志） |
| **恢复速度** | 快（直接加载二进制） | 慢（重放所有命令） |
| **写性能** | 不影响（BGSAVE 异步） | 略有影响（fsync） |
| **适用场景** | 冷备、灾难恢复 | 数据安全要求高 |

### 混合持久化（Redis 4.0+，推荐）

```
aof-use-rdb-preamble yes

AOF 重写时: RDB 格式写入存量数据 + AOF 格式追加增量数据
→ 恢复时: 先加载 RDB 部分（快）→ 再重放 AOF 增量（安全）
→ 兼顾速度和安全性
```

---

## 四、Redis 过期策略与内存淘汰机制

### 过期策略（删除已过期的 key）

```
1. 惰性删除:
   访问 key 时检查是否过期 → 过期则删除
   优点: CPU 友好（只在访问时检查）
   缺点: 内存不友好（不访问的过期 key 永远占内存）

2. 定期删除:
   每 100ms 随机抽取 20 个设有过期时间的 key
   → 检查是否过期 → 过期则删除
   → 如果过期比例 > 25% → 再抽取 20 个（循环）
   → 控制每次执行时间 ≤ 25ms（避免阻塞）

Redis 两种策略配合使用:
  定期删除 → 大部分过期 key
  惰性删除 → 漏网之鱼
  → 仍有漏网的（定期没抽到 + 一直没被访问）→ 由内存淘汰机制兜底
```

### 内存淘汰策略（内存满了怎么办）

```
maxmemory 4gb            → 设置最大内存
maxmemory-policy allkeys-lru  → 内存满时的淘汰策略

8 种淘汰策略:

┌─────────────────────┬──────────────────────────────────────┐
│ 策略                 │ 说明                                  │
├─────────────────────┼──────────────────────────────────────┤
│ noeviction          │ 不淘汰，写入报错 OOM（默认）            │
│ allkeys-lru         │ 所有 key 中淘汰最近最少使用的 ★推荐     │
│ allkeys-lfu         │ 所有 key 中淘汰最不经常使用的（Redis4.0+）│
│ allkeys-random      │ 所有 key 中随机淘汰                    │
│ volatile-lru        │ 有过期时间的 key 中淘汰 LRU             │
│ volatile-lfu        │ 有过期时间的 key 中淘汰 LFU             │
│ volatile-random     │ 有过期时间的 key 中随机淘汰              │
│ volatile-ttl        │ 有过期时间的 key 中淘汰 TTL 最短的       │
└─────────────────────┴──────────────────────────────────────┘

推荐:
  缓存场景 → allkeys-lru（最通用）
  有明确冷热数据 → allkeys-lfu（更精准）
  有些 key 必须保留（配置数据） → volatile-lru（只淘汰设了过期时间的）
```

### LRU vs LFU

```
LRU（Least Recently Used）:
  淘汰最近最久没访问的 → 基于"时间"
  问题: 偶尔被访问一次的冷数据会挤掉热数据

LFU（Least Frequently Used, Redis 4.0+）:
  淘汰访问频率最低的 → 基于"频率"
  → 更精准识别真正的冷热数据

Redis 的近似 LRU:
  不是精确 LRU（太耗内存），而是随机采样 N 个 key，淘汰其中最久没访问的
  maxmemory-samples 5 → 每次采样 5 个（越大越精确，CPU 开销越大）
```

---

## 五、Redis 单线程为什么这么快

```
核心要点（面试必答）:

1. 纯内存操作
   → 内存读写 ~100ns，磁盘 ~10ms → 差 10 万倍
   → Redis 大部分操作是 O(1) 或 O(logN) 的内存操作

2. IO 多路复用（epoll）
   → 单线程监听多个 socket 连接
   → 有数据可读/可写时才处理 → 不阻塞等待
   → 一个线程处理数万并发连接

   ┌──────────┐
   │  Client1  │──┐
   │  Client2  │──┤  epoll    ┌──────────┐
   │  Client3  │──┼──────────→│ 单线程    │
   │  ...      │──┤  事件通知  │ 串行处理  │
   │  ClientN  │──┘           └──────────┘
   └──────────┘

3. 单线程避免了锁竞争和上下文切换
   → 多线程: 加锁/解锁/线程切换 有开销
   → 单线程: 串行执行，无锁 → 简单高效

4. 高效的数据结构
   → SDS、ziplist/listpack、skiplist、intset 等
   → 针对不同场景做了极致优化

注意: Redis 6.0 引入了多线程 IO
  → 多线程仅用于网络 IO（读请求/写响应）
  → 命令执行仍然是单线程（保证原子性）
  → 解决大量客户端连接时网络 IO 成为瓶颈的问题

性能数据:
  单机 QPS: 10 万~15 万（简单命令）
  延迟: 平均 < 1ms
```

---

## 六、Redis 分布式锁

### 基本实现

```go
// 加锁: SET key uuid NX EX 30
// NX: 只在 key 不存在时才设置（互斥）
// EX 30: 30 秒过期（防死锁）
// uuid: 唯一标识（防止误删别人的锁）

func TryLock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (string, bool) {
    uuid := generateUUID()
    ok, err := rdb.SetNX(ctx, key, uuid, ttl).Result()
    if err != nil || !ok {
        return "", false
    }
    return uuid, true
}

// 解锁: 必须用 Lua 脚本保证原子性（先比较再删除）
const unlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`

func Unlock(ctx context.Context, rdb *redis.Client, key, uuid string) bool {
    result, err := rdb.Eval(ctx, unlockScript, []string{key}, uuid).Int64()
    return err == nil && result == 1
}
```

### 为什么解锁需要 Lua 脚本？

```
如果不用 Lua（非原子操作）:
  ① A: GET lock → 值是自己的 uuid ✓
  ② 此时 A 的锁恰好过期被 Redis 删除
  ③ B: SET lock uuid_B NX → 加锁成功
  ④ A: DEL lock → 把 B 的锁删了！

Lua 脚本在 Redis 中原子执行:
  → GET + DEL 在同一次执行中完成，中间不会被打断
```

### 锁续期（看门狗机制）

```
问题: 锁 TTL 设多长？
  设太短 → 业务没执行完锁就过期了 → 其他人拿到锁 → 并发问题
  设太长 → 持锁进程崩溃 → 锁很久才释放 → 影响可用性

解决: 看门狗（Watchdog）自动续期
  → 加锁后启动后台协程，每 TTL/3 续期一次
  → 业务完成后停止续期 + 释放锁
  → 进程崩溃 → 续期停止 → 锁自动过期

  时间线:
  ─── 加锁(TTL=30s) ─── 续期(10s) ─── 续期(10s) ─── 业务完成 ─── 解锁
```

```go
func LockWithWatchdog(ctx context.Context, rdb *redis.Client, key string) (string, context.CancelFunc) {
    uuid := generateUUID()
    ttl := 30 * time.Second

    // 加锁
    ok, _ := rdb.SetNX(ctx, key, uuid, ttl).Result()
    if !ok {
        return "", nil
    }

    // 看门狗: 每 10 秒续期
    watchCtx, cancel := context.WithCancel(ctx)
    go func() {
        ticker := time.NewTicker(ttl / 3)
        defer ticker.Stop()
        for {
            select {
            case <-watchCtx.Done():
                return
            case <-ticker.C:
                rdb.Expire(watchCtx, key, ttl)
            }
        }
    }()

    // 返回 cancel 用于停止看门狗
    return uuid, func() {
        cancel() // 停止看门狗
        Unlock(context.Background(), rdb, key, uuid)
    }
}
```

### Redlock（多节点分布式锁）

```
问题: 单节点 Redis 锁的问题
  主节点加锁成功 → 主节点宕机 → 从节点提升为主（但锁数据还没同步）
  → 新主节点上没有锁 → 另一个客户端也能加锁 → 互斥性被破坏

Redlock 算法（Redis 作者提出）:
  使用 N 个独立的 Redis 实例（推荐 5 个，不是主从/集群）

  加锁流程:
    ① 获取当前时间戳 T1
    ② 依次向 5 个实例发送 SET key uuid NX EX
    ③ 如果在超过半数（≥3）的实例上加锁成功
       且总耗时 < 锁 TTL → 加锁成功
    ④ 锁的有效时间 = TTL - 加锁耗时
    ⑤ 加锁失败 → 向所有实例发送 DEL 释放

争议:
  Martin Kleppmann 质疑 Redlock:
  → 进程暂停（GC、page fault）可能导致锁失效
  → 时钟跳跃可能破坏算法假设
  → 建议需要强互斥的场景使用 ZooKeeper/etcd

实际建议:
  大部分场景 → 单节点 Redis 锁 + 看门狗（足够了）
  强一致要求 → etcd 分布式锁（Raft 共识，更可靠）
```

---

## 七、Redis Cluster 集群原理

### 哈希槽（Hash Slot）

```
Redis Cluster 将数据分成 16384 个哈希槽（slot）

  key → CRC16(key) % 16384 → 落到某个 slot → slot 属于某个节点

  3 个主节点:
    Node A: slot 0 ~ 5460
    Node B: slot 5461 ~ 10922
    Node C: slot 10923 ~ 16383

  SET user:1001 → CRC16("user:1001") % 16384 = 7234 → Node B

为什么是 16384？
  → 节点间心跳包需要携带 slot 信息的 bitmap
  → 16384 个 slot → bitmap 大小 = 16384/8 = 2KB
  → 既能支持 1000+ 节点，又不会让心跳包太大
```

### 集群架构

```
  ┌────────────┐    ┌────────────┐    ┌────────────┐
  │ Master A    │    │ Master B    │    │ Master C    │
  │ slot 0~5460│    │ slot 5461~  │    │ slot 10923~ │
  │             │    │   10922     │    │   16383     │
  └──────┬─────┘    └──────┬─────┘    └──────┬─────┘
         │                 │                 │
  ┌──────┴─────┐    ┌──────┴─────┐    ┌──────┴─────┐
  │ Slave A'   │    │ Slave B'   │    │ Slave C'   │
  └────────────┘    └────────────┘    └────────────┘

特点:
  - 去中心化（无代理，客户端直连节点）
  - 每个节点都知道 slot 分布 → 请求打错节点会返回 MOVED 重定向
  - 主节点挂了 → 从节点自动提升为主（类似哨兵）
```

### MOVED 重定向 vs ASK 重定向

```
MOVED（永久重定向）:
  客户端 → Node A: GET user:1001
  Node A: 这个 slot 不在我这 → MOVED 7234 Node_B_IP:Port
  客户端: 更新本地 slot 映射表 → 后续直接访问 Node B

ASK（临时重定向，扩容迁移中）:
  slot 正在从 A 迁移到 B
  客户端 → Node A: GET user:1001
  Node A: 这个 key 已经迁走了 → ASK 7234 Node_B_IP:Port
  客户端: 先发 ASKING → 再发 GET（仅本次，不更新映射表）
```

---

## 八、Redis 主从复制与哨兵机制

### 主从复制

```
作用: 数据冗余 + 读写分离（主写从读）

全量同步（首次连接/断线太久）:
  ① Slave → Master: PSYNC ? -1（请求全量）
  ② Master: BGSAVE 生成 RDB
  ③ Master → Slave: 发送 RDB 文件
  ④ Slave: 清空旧数据 → 加载 RDB
  ⑤ Master: 发送 RDB 期间的增量命令（repl buffer）
  ⑥ Slave: 执行增量命令 → 数据一致

增量同步（短暂断线重连）:
  ① Slave → Master: PSYNC runid offset（带上复制偏移量）
  ② Master: 检查 offset 是否在 repl_backlog 中
  ③ 在 → 只发送缺失的增量命令（高效）
  ④ 不在（断线太久，积压缓冲区已覆盖）→ 全量同步

  repl_backlog_size 1mb → 积压缓冲区大小（建议调大）
```

### 哨兵机制（Sentinel）

```
作用: 自动故障转移（主挂了自动切换从为主）

  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │ Sentinel1│  │ Sentinel2│  │ Sentinel3│  (至少 3 个，奇数)
  └────┬─────┘  └────┬─────┘  └────┬─────┘
       │              │              │
       └──────────────┼──────────────┘
                      │ 监控
       ┌──────────────┼──────────────┐
  ┌────┴─────┐  ┌─────┴────┐  ┌─────┴────┐
  │ Master   │  │ Slave1   │  │ Slave2   │
  └──────────┘  └──────────┘  └──────────┘

故障检测:
  ① 主观下线（SDOWN）: 单个 Sentinel 认为 Master 不可达
     → ping 超时（down-after-milliseconds 30000）
  ② 客观下线（ODOWN）: quorum 个 Sentinel 都认为 Master 不可达
     → 触发故障转移

故障转移流程:
  ① Sentinel 选举 Leader（Raft 算法）
  ② Leader 选择最优 Slave（优先级 > 偏移量 > runid）
  ③ 将选中的 Slave 执行 SLAVEOF NO ONE（提升为 Master）
  ④ 其他 Slave 指向新 Master
  ⑤ 旧 Master 恢复后变为 Slave
```

---

## 九、Redis 大 Key 和热 Key 问题

### 大 Key

```
定义:
  String 类型 > 10KB
  Hash/List/Set/ZSet 元素个数 > 5000 或总大小 > 10MB

危害:
  ① 阻塞主线程: DEL 大 key 可能耗时几秒（Redis 是单线程！）
  ② 内存不均: Cluster 模式下导致某个节点内存远大于其他节点
  ③ 网络阻塞: 传输大 value 占满带宽
  ④ 持久化变慢: RDB/AOF 写入大 key 耗时长

排查:
  redis-cli --bigkeys          → 扫描各类型最大 key
  redis-cli --memkeys          → 按内存大小排序
  MEMORY USAGE key             → 查看单个 key 占用内存

解决:
  ① 拆分: 大 Hash 拆成多个小 Hash
     user:1001 → user:1001:basic、user:1001:detail
  ② 压缩: 序列化前压缩（gzip/snappy）
  ③ 异步删除: UNLINK 代替 DEL（Redis 4.0+，后台线程删除）
  ④ lazyfree: lazyfree-lazy-expire yes → 过期 key 异步删除
```

### 热 Key

```
定义: 某个 key 被高频访问（如微博热搜、秒杀商品）

危害:
  → 该 key 所在节点压力过大 → 成为瓶颈/被打挂
  → Cluster 模式下无法通过加节点解决（key 只在一个节点上）

排查:
  redis-cli --hotkeys          → 需要开启 LFU 淘汰策略
  Redis MONITOR 命令           → 实时监控（性能影响大，慎用）
  客户端统计（代码层面计数）

解决:
  ① 本地缓存: 进程内缓存热 key（bigcache/freecache）
     → 不经过网络，~100ns
  ② 读副本: 热 key 复制 N 份（key_1、key_2...key_N）
     → 客户端随机读不同副本 → 分散到不同节点
  ③ 代理层缓存: Proxy（如 Codis）层做本地缓存
```

---

## 十、Redis 事务与 Pipeline

### 事务（MULTI/EXEC）

```
MULTI       → 开始事务
SET k1 v1   → 入队（不立即执行）
SET k2 v2   → 入队
INCR k3     → 入队
EXEC        → 原子执行所有命令

特点:
  ① 命令打包，一次性执行（减少网络 RTT）
  ② 不支持回滚！如果某条命令执行失败，其他命令仍会执行
     → 与 MySQL 事务不同！Redis 事务不保证原子性
  ③ WATCH 实现乐观锁:
     WATCH k1        → 监视 k1
     val = GET k1    → 读取当前值
     MULTI
     SET k1 val+1    → 基于读取的值修改
     EXEC            → 如果 k1 被其他客户端修改 → EXEC 返回 nil（事务取消）
```

### Pipeline

```
批量发送命令，减少网络往返（RTT）

无 Pipeline:
  SET k1 → 等响应 → SET k2 → 等响应 → SET k3 → 等响应
  → 3 次 RTT

有 Pipeline:
  SET k1 + SET k2 + SET k3 → 一次发送 → 一次接收所有响应
  → 1 次 RTT

性能: 10 万次 SET
  无 Pipeline: ~50 秒（网络延迟主导）
  有 Pipeline: ~0.5 秒（100 倍提升）
```

```go
// Go Redis Pipeline 示例
pipe := rdb.Pipeline()
for i := 0; i < 10000; i++ {
    pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, time.Hour)
}
cmds, err := pipe.Exec(ctx)
// cmds 包含所有命令的结果
```

### Pipeline vs 事务 vs Lua 脚本

| 维度 | Pipeline | 事务（MULTI/EXEC） | Lua 脚本 |
|------|---------|-------------------|---------|
| **网络 RTT** | 1 次 | 1 次（EXEC 触发） | 1 次 |
| **原子性** | 无（可被其他命令插入） | 弱（不支持回滚） | 强（单线程原子执行） |
| **条件逻辑** | 不支持 | 不支持 | 支持（if/else/for） |
| **适用场景** | 批量无关命令 | 简单多命令打包 | 需原子性+条件逻辑 |

---

## 十一、Redis/MySQL/etcd 集群数据同步机制与 CAP 分析

### CAP 定理回顾

```
CAP 定理: 分布式系统最多同时满足以下三项中的两项

  C - Consistency（一致性）:
    所有节点在同一时刻看到的数据相同
    → 读操作一定能读到最新写入的数据

  A - Availability（可用性）:
    每个请求都能在合理时间内收到响应（不能超时/报错）
    → 即使部分节点故障，系统仍能提供服务

  P - Partition tolerance（分区容忍性）:
    网络分区（节点间通信中断）时系统仍能运行
    → 分布式系统必须选 P（网络分区无法避免）

  因此实际选择是: CP 或 AP

         C（一致性）
        / \
       /   \
      /     \
     / CP  AP\
    /_________\
   P（分区容忍） A（可用性）

  CP: 保一致性，牺牲可用性 → etcd、ZooKeeper
  AP: 保可用性，牺牲一致性 → Redis Cluster、Eureka、Cassandra
```

---

### 一、Redis Cluster 数据同步机制

```
类型: 异步复制
CAP 选择: AP（高可用优先，最终一致性）

同步机制:
  ① 主从异步复制
     → Master 写完立即返回客户端成功
     → 异步将写命令传播给 Slave
     → Slave 可能有延迟（lag）

  写入流程:
    Client → Master: SET k1 v1
    Master: 写入内存 → 返回 OK（不等 Slave）
    Master → Slave: 异步传播 SET k1 v1（通过 repl_backlog）
    Slave: 执行 SET k1 v1

  ┌────────┐  ① SET k1   ┌────────┐  ③ 异步传播  ┌────────┐
  │ Client │────────────→│ Master │─────────────→│ Slave  │
  │        │←────────────│        │              │        │
  └────────┘  ② OK(立即) └────────┘              └────────┘
                          ↑
                    不等 Slave 确认

复制方式:
  全量同步: Slave 首次连接 → Master BGSAVE 生成 RDB → 发送给 Slave
  增量同步: 断线重连 → 基于 offset 从 repl_backlog 获取增量命令
  命令传播: 正常运行时 → Master 实时将写命令异步传给 Slave

数据丢失场景:
  ① Master 写入成功但还没传给 Slave → Master 宕机
     → Slave 提升为新 Master → 丢失未同步的数据
  ② 网络分区: Master 和 Slave 断开 → Master 继续接受写入
     → 分区恢复后 Slave 可能跟新 Master 同步 → 旧数据丢失

  min-replicas-to-write 1      → 至少 1 个 Slave 确认才返回
  min-replicas-max-lag 10      → Slave 延迟不超过 10 秒
  → 可以减少数据丢失，但牺牲了可用性（趋向 CP）
```

#### Gossip 协议（集群元数据同步）

```
Redis Cluster 节点间通过 Gossip 协议交换集群状态:

  每个节点每秒随机选几个节点发送 PING:
    → 携带: 自己的状态 + 部分其他节点的状态
    → 收到 PING 回复 PONG
    → 信息像"八卦"一样逐渐扩散到所有节点

  ┌───────┐ PING  ┌───────┐ PING  ┌───────┐
  │ NodeA │──────→│ NodeB │──────→│ NodeC │
  │       │←──────│       │←──────│       │
  └───────┘ PONG  └───────┘ PONG  └───────┘

  优点: 去中心化、无单点故障、最终一致
  缺点: 状态收敛有延迟（不是立即一致）
```

---

### 二、MySQL 集群数据同步机制

```
类型: 多种模式可选（异步 / 半同步 / 组复制）
CAP 选择: 取决于部署模式
  - 异步复制: AP
  - 半同步复制: 偏 CP（但不严格）
  - Group Replication(MGR): CP
```

#### 2.1 异步复制（默认）

```
流程:
  ① Master 执行写操作 → 写入 Binlog
  ② Master 返回客户端成功（不等 Slave）
  ③ Binlog Dump 线程将 Binlog 事件发给 Slave
  ④ Slave IO 线程接收 → 写入 Relay Log
  ⑤ Slave SQL 线程重放 Relay Log → 数据一致

  ┌────────┐ ①写入  ┌────────────┐ ③Binlog  ┌────────────┐
  │ Client │──────→│  Master     │────────→│   Slave     │
  │        │←──────│ Binlog 落盘 │         │ Relay Log   │
  └────────┘②OK    │             │         │ → SQL 重放  │
                   └────────────┘         └────────────┘

  问题: Master 宕机 → Binlog 还没传到 Slave → 数据丢失
  延迟: 通常 ms 级，高负载下可能到秒级
```

#### 2.2 半同步复制（Semi-Sync）

```
改进: Master 等至少一个 Slave 确认收到 Binlog 才返回客户端

流程:
  ① Master 写 Binlog
  ② Master 将 Binlog 发给 Slave
  ③ Slave 接收并写入 Relay Log → 返回 ACK
  ④ Master 收到至少 1 个 ACK → 返回客户端成功

  ┌────────┐ ①写入  ┌────────┐ ②Binlog  ┌────────┐
  │ Client │──────→│ Master │────────→│ Slave  │
  │        │       │ 等ACK  │←────────│ ACK    │
  │        │←──────│        │ ④返回    └────────┘
  └────────┘       └────────┘

  两种模式:
    AFTER_COMMIT: 先 commit 再等 ACK（其他会话可能读到未确认数据）
    AFTER_SYNC:   先等 ACK 再 commit（更强一致，MySQL 5.7+推荐）

  超时降级:
    Slave 超时未 ACK → 自动退化为异步复制（保可用性 → AP）
    → 所以半同步不是严格的 CP
```

#### 2.3 组复制 MGR（Group Replication）

```
MySQL 5.7.17+ 引入，基于 Paxos 变体协议

  流程:
    ① 写操作在本地执行
    ② 将事务广播给组内所有成员
    ③ 多数节点（N/2+1）确认 → 事务提交
    ④ 少数节点确认 → 事务回滚

  ┌────────┐    ┌────────┐    ┌────────┐
  │ Node A │←──→│ Node B │←──→│ Node C │
  │        │    │        │    │        │
  └────────┘    └────────┘    └────────┘
       ↕ Paxos（多数派同意才提交）

  两种模式:
    单主模式: 一个写节点 + 多个只读节点
    多主模式: 所有节点都可写（冲突检测+回滚）

  CAP: CP（多数派确认 → 强一致，少数派不可写 → 牺牲可用性）
```

---

### 三、etcd 集群数据同步机制

```
类型: Raft 共识算法
CAP 选择: CP（强一致性，牺牲部分可用性）

Raft 角色:
  Leader:    唯一接受写请求的节点
  Follower:  被动接收 Leader 的日志复制
  Candidate: 选举期间的临时角色
```

#### 3.1 写入流程（Log Replication）

```
  ① Client → Leader: PUT /key "value"
  ② Leader 将操作写入本地 WAL 日志（uncommitted）
  ③ Leader 并行发送 AppendEntries RPC 给所有 Follower
  ④ Follower 写入本地 WAL → 返回 ACK
  ⑤ Leader 收到多数派（N/2+1）ACK → 提交（committed）
  ⑥ Leader 返回客户端成功
  ⑦ 下一次心跳通知 Follower 提交

  ┌────────┐ ①写入 ┌────────┐ ③复制  ┌──────────┐
  │ Client │──────→│ Leader │──────→│Follower1 │ ACK ✓
  │        │       │        │──────→│Follower2 │ ACK ✓  → 多数派(2/3) → 提交
  │        │←──────│ commit │──────→│Follower3 │ 超时 ✗
  └────────┘ ⑥OK   └────────┘       └──────────┘

  关键: 必须多数派确认才能提交
    3 节点: 至少 2 个确认 → 容忍 1 个故障
    5 节点: 至少 3 个确认 → 容忍 2 个故障
```

#### 3.2 Leader 选举

```
触发条件: Follower 超时未收到 Leader 心跳（election timeout）

流程:
  ① Follower → Candidate: 自增 term，投票给自己
  ② Candidate → 其他节点: RequestVote RPC
  ③ 获得多数票 → 成为新 Leader
  ④ 立即发送心跳宣示主权

  选举约束:
    每个 term 每个节点只能投一票（先到先得）
    Candidate 的日志必须 ≥ 投票者的日志（日志完整性检查）
    → 保证新 Leader 拥有所有已提交的日志

  ┌──────────┐  超时  ┌───────────┐  多数票  ┌────────┐
  │ Follower │──────→│ Candidate │────────→│ Leader │
  └──────────┘       └───────────┘         └────────┘
                      ↑ 竞选失败            │ 心跳
                      └────────────────────┘
```

#### 3.3 线性一致性读

```
etcd 默认提供线性一致性读（Linearizable Read）:

  问题: 直接读 Follower 可能读到旧数据（日志未同步完）

  方案 1: ReadIndex（默认）
    ① Leader 确认自己仍是 Leader（发一轮心跳确认多数派响应）
    ② 记录当前 commitIndex
    ③ 等待状态机 apply 到该 index → 返回结果
    → 保证读到的是最新已提交的数据

  方案 2: Serializable Read（非线性，但更快）
    → 直接从本地状态机读取（可能读到旧数据）
    → 适合对一致性要求不高的场景
```

---

### 四、三者对比总结

| 维度 | Redis Cluster | MySQL | etcd |
|------|--------------|-------|------|
| **同步协议** | 异步复制 + Gossip | Binlog 复制（异步/半同步/Paxos） | Raft 共识 |
| **CAP 选择** | **AP** | 异步=AP，半同步≈CP，MGR=**CP** | **CP** |
| **一致性** | 最终一致（异步有延迟） | 取决于复制模式 | 强一致（线性一致性） |
| **写入确认** | Master 写完即返回 | 异步:不等；半同步:等1个ACK | 等多数派 ACK |
| **数据丢失风险** | Master 宕机可能丢数据 | 异步:可能丢；半同步:极少丢 | 已提交数据不丢 |
| **可用性** | 高（主挂从顶） | 中（半同步可能降级） | 中（少数派不可用） |
| **故障容忍** | N 主 N 从，每组容忍 1 故障 | 取决于架构 | 2N+1 节点容忍 N 故障 |
| **节点角色** | Master / Slave | Master / Slave（MGR 多主） | Leader / Follower |
| **元数据同步** | Gossip（去中心化） | 无（依赖外部如 MHA/Orchestrator） | Raft（自管理） |
| **适用场景** | 缓存、高吞吐 | 事务型业务数据 | 配置中心、服务发现、分布式锁 |

---

### 五、同步机制的本质区别

```
1. 异步复制（Redis/MySQL 默认）
   → "写完就走，不管从库"
   → 最快，但可能丢数据
   → 适合: 性能优先，允许短暂不一致

2. 半同步复制（MySQL Semi-Sync）
   → "至少一个从库收到才算完"
   → 折中方案，减少数据丢失风险
   → 适合: 业务数据，需要一定安全保障

3. 多数派共识（etcd Raft / MySQL MGR Paxos）
   → "超过半数节点确认才算完"
   → 最安全，但延迟最高
   → 适合: 元数据、配置、锁等关键数据

     性能            安全性
   ◄────────────────────────────►
   异步复制 → 半同步 → 多数派共识
   (Redis)   (MySQL)   (etcd)
```

---

### 六、面试高频追问

```
Q: 为什么 Redis 选 AP 而不是 CP？
A: Redis 的定位是高性能缓存:
   → 缓存数据允许丢失（DB 是最终数据源）
   → 10 万+ QPS 要求极低延迟 → 不能等 Slave 确认
   → 用最终一致性换取极致性能

Q: 为什么 etcd 选 CP 而不是 AP？
A: etcd 存储的是关键元数据（服务注册、配置、分布式锁）:
   → 配置错误/锁失效的后果远大于短暂不可用
   → 宁可拒绝服务，也不能返回错误数据
   → Kubernetes 就依赖 etcd 的强一致性

Q: MySQL 半同步超时降级为异步，怎么办？
A: ① 监控复制延迟（Seconds_Behind_Master）
   ② 降级告警 + 人工介入
   ③ 用 MGR 替代半同步（真正的 CP）
   ④ 业务层读写分离时注意: 刚写完的数据走主库读

Q: 网络分区时各系统的表现？
A:
   Redis: 少数派 Master 继续服务（数据可能不一致）→ 分区恢复后数据合并/覆盖
   MySQL 半同步: 退化为异步（保可用） → 数据安全性降低
   etcd: 少数派停止服务（保一致） → 多数派正常 → 分区恢复后自动追赶日志

Q: Raft 和 Paxos 的区别？
A:
   Raft: Leader 驱动，更简单易懂
     → 日志必须连续、有序
     → Leader 拥有所有已提交日志（选举约束）
     → etcd、Consul 使用
   Paxos: 更通用，理论更完备，但实现复杂
     → Multi-Paxos 在实践中类似 Raft
     → Google Spanner、MySQL MGR 变体使用
```

### 一句话总结

**Redis Cluster 采用异步复制+Gossip 协议（AP，写完即返回不等从库，可能丢数据但性能极高），MySQL 提供三种模式（异步=AP/半同步≈CP 至少一个从库 ACK/MGR=CP 基于 Paxos 多数派确认），etcd 采用 Raft 共识算法（CP，写入必须多数派确认才返回，已提交数据不丢失）。** 本质区别在于写入确认机制——异步（不等）→半同步（等 1 个）→多数派共识（等 N/2+1 个），从左到右性能递减但安全性递增。选型原则：缓存/高吞吐场景选 AP（Redis），业务数据选半同步/MGR（MySQL），关键元数据/配置/锁选 CP（etcd）。

---

## 十二、ZSet 实现原理：跳表为什么优于红黑树

### ZSet 的底层数据结构

```
ZSet（Sorted Set）有两种编码方式:

① listpack（元素少且值小时）:
   条件: 元素个数 ≤ zset-max-listpack-entries（默认 128）
         且每个元素大小 ≤ zset-max-listpack-value（默认 64 字节）
   → 紧凑的连续内存块，省内存，但增删改查 O(N)

② skiplist + hashtable（超过阈值时）:
   skiplist: 按 score 排序，支持范围查询
   hashtable: member → score 的映射，支持 O(1) 查分数

   ┌──────────────────────────────────────────────────┐
   │                     ZSet                          │
   │                                                    │
   │  ┌─────────────┐        ┌─────────────────┐       │
   │  │  hashtable   │        │    skiplist      │       │
   │  │              │        │                  │       │
   │  │ "Alice"→100  │        │  10─→50─→100─→200│       │
   │  │ "Bob"  →50   │        │  (按score有序)    │       │
   │  │ "Carol"→200  │        │                  │       │
   │  │ "Dave" →10   │        │                  │       │
   │  └─────────────┘        └─────────────────┘       │
   │                                                    │
   │  ZSCORE key member → 查 hashtable → O(1)          │
   │  ZRANGEBYSCORE     → 查 skiplist  → O(logN+M)     │
   └──────────────────────────────────────────────────┘

两者共享 member 和 score 的指针（不重复存储）
```

### 跳表（Skip List）结构详解

```
跳表 = 多层有序链表，通过"索引层"加速查找

普通有序链表: 查找 score=80 → 从头遍历 → O(N)

  Level 0: 10 → 20 → 30 → 50 → 60 → 70 → 80 → 90 → NULL

跳表（多层索引）: 查找 score=80 → 从高层往下 → O(logN)

  Level 3: 10 ──────────────────────────────────→ 80 ────→ NULL
           │                                      │
  Level 2: 10 ────────────→ 50 ─────────────────→ 80 ────→ NULL
           │                │                      │
  Level 1: 10 ────→ 30 ──→ 50 ────→ 70 ────────→ 80 ────→ NULL
           │        │       │        │             │
  Level 0: 10 → 20→ 30 → 50→ 60 → 70 → 80 → 90→ NULL
           (底层包含所有元素)

查找 score=80 的路径:
  Level 3: 10 → 80 ✓ 找到！
  → 最多经过 O(logN) 个节点
```

### Redis 跳表的 C 实现结构

```c
// 跳表节点
typedef struct zskiplistNode {
    sds ele;                         // member（SDS 字符串）
    double score;                    // 分数
    struct zskiplistNode *backward;  // 后退指针（只有 Level 0 有）
    struct zskiplistLevel {
        struct zskiplistNode *forward;  // 前进指针
        unsigned long span;             // 跨度（用于计算排名 ZRANK）
    } level[];                       // 柔性数组，层数动态分配
} zskiplistNode;

// 跳表
typedef struct zskiplist {
    struct zskiplistNode *header, *tail;  // 头尾指针
    unsigned long length;                  // 节点总数
    int level;                             // 当前最高层数
} zskiplist;

关键字段:
  span（跨度）: 记录当前节点到 forward 节点之间跳过了多少个节点
  → ZRANK 命令: 沿查找路径累加 span 即可得到排名 → O(logN)
  backward: 反向指针，支持 ZREVRANGE 反向遍历
```

### 跳表每层高度如何确定？（随机化）

```
核心: 每个节点的层数通过随机算法决定，不需要全局重平衡！

Redis 源码（t_zset.c）:

  #define ZSKIPLIST_MAXLEVEL 32   // 最大层数
  #define ZSKIPLIST_P 0.25        // 晋升概率

  int zslRandomLevel(void) {
      int level = 1;
      while ((random() & 0xFFFF) < (ZSKIPLIST_P * 0xFFFF))
          level += 1;
      return (level < ZSKIPLIST_MAXLEVEL) ? level : ZSKIPLIST_MAXLEVEL;
  }

算法解析:
  ① 初始层数 = 1（所有节点至少在 Level 0）
  ② 以 25% 的概率升到下一层（循环）
  ③ 最大不超过 32 层

概率分布:
  Level 1:  75%      的节点（3/4 的节点只在底层）
  Level 2:  18.75%   的节点（约 1/4 × 1/4 的概率连续两次）
  Level 3:  4.69%    的节点
  Level 4:  1.17%    的节点
  ...
  Level 32: 极小概率

  → 平均每个节点层数 = 1/(1-P) = 1/(1-0.25) = 1.33 层
  → 额外内存开销很小（每个节点平均多 0.33 个指针）

为什么 P=0.25 而不是 0.5？
  P=0.5:  平均层数 2.0，索引更密 → 查找更快，但内存多
  P=0.25: 平均层数 1.33，索引更稀 → 查找稍慢，但省内存
  → Redis 选择 P=0.25 是时间和空间的折中
  → 实测差距不大（O(logN) 常数系数差别），但内存省 ~33%

32 层最多支持多少元素？
  理论: 4^32 ≈ 1.8 × 10^19（远超实际使用量）
  → 即使 2^64 个元素也能 O(logN) 查找
```

### 为什么 Redis 用跳表而不用红黑树？

```
这是面试最核心的问题！Redis 作者 antirez 曾亲自回答过:

┌──────────────┬──────────────────────┬──────────────────────┐
│     维度      │       跳表            │       红黑树          │
├──────────────┼──────────────────────┼──────────────────────┤
│ 范围查询      │ O(logN) 定位起点      │ O(logN) 定位起点      │
│ ZRANGEBYSCORE│ → 沿 Level 0 顺序遍历 │ → 中序遍历（需要回溯   │
│              │ → 天然有序链表，极快    │    父节点，不连续）    │
├──────────────┼──────────────────────┼──────────────────────┤
│ 实现复杂度    │ 简单（~200 行核心代码）│ 复杂（旋转+变色，      │
│              │ 插入/删除只改指针      │  ~500 行，易出 bug）   │
├──────────────┼──────────────────────┼──────────────────────┤
│ 插入/删除     │ 改局部指针即可        │ 可能触发多次旋转+变色  │
│              │ 无需全局重平衡        │ （最多 O(logN) 次）    │
├──────────────┼──────────────────────┼──────────────────────┤
│ 内存局部性    │ 节点内存连续度一般    │ 类似                   │
│              │ 但 Level 0 是链表     │                       │
├──────────────┼──────────────────────┼──────────────────────┤
│ 排名查询      │ span 字段 → ZRANK    │ 需要额外维护 size 字段 │
│ ZRANK        │ O(logN)              │ 或遍历 O(N)           │
├──────────────┼──────────────────────┼──────────────────────┤
│ 并发友好      │ 局部锁即可（改几个指针）│ 旋转涉及多个节点，    │
│              │                       │ 锁粒度大              │
├──────────────┼──────────────────────┼──────────────────────┤
│ 时间复杂度    │ 期望 O(logN)         │ 严格 O(logN)          │
│              │ （概率保证）          │ （确定性保证）         │
└──────────────┴──────────────────────┴──────────────────────┘
```

#### 关键原因总结（面试必答 4 点）

```
① 范围查询天然高效
   ZRANGEBYSCORE min max:
     跳表: 找到 min → 沿 Level 0 链表顺序走到 max → O(logN + M)
     红黑树: 找到 min → 中序遍历需要不断回溯父节点 → 不如链表直观

   这是最重要的原因！ZSet 最常用的就是范围操作。

② 实现简单，易于维护和调试
   跳表核心逻辑 ~200 行 C 代码
   红黑树需要处理: 左旋、右旋、变色、4 种插入情况、6 种删除情况
   Redis 是实用主义项目，简单 = 更少 bug = 更好维护

③ ZRANK 排名查询高效
   跳表: 每个节点的 level[].span 记录跨度
   → 查找路径上累加 span 即可得到排名 → O(logN)
   红黑树: 要支持排名需要额外在每个节点维护子树大小
   → 每次旋转都要更新 size → 增加实现复杂度

④ 插入/删除更简单
   跳表: 找到位置 → 修改前后指针（每层独立）→ 不影响其他节点
   红黑树: 插入/删除后可能破坏平衡 → 旋转+变色可能传播到根节点
```

#### antirez 原话（Redis 作者）

```
"There are a few reasons:

1) They are not very memory intensive. It's up to you basically.
   Changing parameters about the probability of a node to have a
   given number of levels will make then less memory intensive than
   btrees.

2) A sorted set is often target of ZRANGEBYSCORE or similar commands,
   that is, traversing the skip list as a linked list. With this
   operation the cache locality of skip lists is at least as good as
   with other kind of balanced trees.

3) They are simpler to implement, debug, and so forth. For instance
   thanks to the skip list simplicity I received a patch (already in
   Redis master) with augmented skip lists implementing ZRANK in
   O(log(N)). It required little changes to the code."

核心总结:
  → 内存可控（调 P 值）
  → 范围遍历高效（链表特性）
  → 实现简单，易于扩展（ZRANK 用 span 很容易加上）
```

### 跳表 vs 红黑树 vs B+树

```
            查找    范围查询      插入/删除     实现   使用场景
跳表        O(logN) O(logN+M)链表 O(logN)简单   简单   Redis ZSet
红黑树      O(logN) O(logN+M)中序 O(logN)旋转   复杂   Java TreeMap
B+树        O(logN) O(logN+M)叶链 O(logN)分裂合并 复杂   MySQL 索引、文件系统

为什么 MySQL 不用跳表？
  → B+树节点大（16KB 页），一个节点存几百个 key
  → 树高 3~4 层，磁盘 IO 只需 3~4 次
  → 跳表每个节点一个 key，层数 ~20+，磁盘 IO 太多
  → B+树是为磁盘优化的，跳表是为内存优化的

为什么 Java TreeMap 用红黑树不用跳表？
  → TreeMap 不需要范围查询那么频繁
  → 红黑树单点查找有确定性 O(logN) 保证（跳表是概率性的）
  → Java 标准库更倾向于确定性数据结构
```

### 一句话总结

**ZSet 底层由 skiplist（按 score 排序，支持范围查询）+ hashtable（member→score，O(1) 查分数）组成。Redis 选跳表而非红黑树的核心原因：① 范围查询天然高效（Level 0 就是有序链表，ZRANGEBYSCORE 直接顺序遍历）；② 实现简单（~200 行 vs 红黑树的旋转变色）；③ ZRANK 排名 O(logN)（span 字段累加）；④ 插入删除只改局部指针无需全局重平衡。** 跳表每层高度通过随机算法决定（P=0.25 概率晋升），平均每节点 1.33 层，最大 32 层。这是概率平衡而非严格平衡——用随机化避免了红黑树复杂的旋转操作，期望时间复杂度同为 O(logN)。

---

## 十三、Redis Cluster MGET 原理与 Slot 迁移

### 单机 Redis 的 MGET

```
单机模式 MGET 很简单:
  MGET key1 key2 key3 → 一次请求，一次返回，O(N)

  Client ──── MGET k1 k2 k3 ────→ Redis
  Client ←── [v1, v2, v3] ───── Redis

  所有 key 在同一个实例 → 一次网络 RTT → 完美
```

### Redis Cluster 的 MGET 问题

```
核心问题: MGET 的多个 key 可能分布在不同的 slot → 不同的节点！

  MGET user:1001 user:2002 user:3003

  CRC16("user:1001") % 16384 = 7234  → Node B (slot 5461~10922)
  CRC16("user:2002") % 16384 = 12456 → Node C (slot 10923~16383)
  CRC16("user:3003") % 16384 = 3210  → Node A (slot 0~5460)

  三个 key 在三个不同节点！
  → Redis Cluster 不支持跨 slot 的 MGET！
  → 直接执行会报错: (error) CROSSSLOT Keys in request don't hash to the same slot
```

### 客户端如何实现跨节点 MGET

```
方案: 客户端按 slot 分组 → 并行发送到各节点 → 合并结果

  ┌────────┐
  │ Client │
  │        │ ① 计算每个 key 的 slot → 按节点分组
  │        │
  │  Group │ key1,key4 → Node A (slot 0~5460)
  │  Group │ key2,key5 → Node B (slot 5461~10922)
  │  Group │ key3,key6 → Node C (slot 10923~16383)
  │        │
  │        │ ② 并行发送 MGET 到各节点
  │        │
  │  Node A│←── MGET key1 key4 ──→ [v1, v4]
  │  Node B│←── MGET key2 key5 ──→ [v2, v5]
  │  Node C│←── MGET key3 key6 ──→ [v3, v6]
  │        │
  │        │ ③ 合并结果，按原始顺序返回
  │        │ → [v1, v2, v3, v4, v5, v6]
  └────────┘

网络 RTT:
  最好情况: 所有 key 在同一节点 → 1 次 RTT
  最坏情况: N 个 key 在 N 个节点 → N 次 RTT（并行则 1 次 RTT）
```

### Go 代码实现

```go
// go-redis 的 cluster client 内部就是这么做的
func ClusterMGet(ctx context.Context, rdb *redis.ClusterClient, keys ...string) ([]interface{}, error) {
    // ① 按 slot 分组
    slotKeys := make(map[int][]string)  // slot → keys
    keyIndex := make(map[string]int)     // key → 原始位置
    for i, key := range keys {
        slot := hashtag.Slot(key)  // CRC16(key) % 16384
        slotKeys[slot] = append(slotKeys[slot], key)
        keyIndex[key] = i
    }

    results := make([]interface{}, len(keys))
    var mu sync.Mutex
    var wg sync.WaitGroup

    // ② 并行发送到各节点
    for _, skeys := range slotKeys {
        wg.Add(1)
        go func(ks []string) {
            defer wg.Done()
            // 同一 slot 的 key 发到同一节点
            vals, err := rdb.MGet(ctx, ks...).Result()
            if err != nil {
                return
            }
            // ③ 合并结果
            mu.Lock()
            for i, k := range ks {
                results[keyIndex[k]] = vals[i]
            }
            mu.Unlock()
        }(skeys)
    }
    wg.Wait()
    return results, nil
}

// 实际使用: go-redis ClusterClient.MGet() 内部已经实现了上述逻辑
// 直接调用即可:
vals, err := clusterClient.MGet(ctx, "k1", "k2", "k3").Result()
```

### Hash Tag：强制 key 落到同一 slot

```
如果业务上多个 key 需要一起 MGET，可以用 Hash Tag 强制同 slot:

  user:{1001}:name    → CRC16("1001") % 16384
  user:{1001}:age     → CRC16("1001") % 16384
  user:{1001}:avatar  → CRC16("1001") % 16384

  Redis 只对 {} 内的部分做 CRC16 → 三个 key 同一 slot → 同一节点
  → MGET 可以一次完成，不用分组！

  规则: 如果 key 中包含 {xxx}，只用 xxx 计算 slot

  优点: 保证相关 key 共存，MGET/事务/Lua 都可以跨 key 操作
  缺点: 可能导致数据倾斜（大量 key 挤在同一 slot/节点）
```

---

### Slot 迁移对 MGET 的影响

```
场景: 集群扩容/缩容时，slot 正在从 Source 迁移到 Target

  扩容前:                    扩容中:
  Node A: slot 0~5460       Node A: slot 0~5460（slot 3000 迁移中→D）
  Node B: slot 5461~10922   Node B: slot 5461~10922
  Node C: slot 10923~16383  Node C: slot 10923~16383
                             Node D: 新节点（接收 slot 3000）
```

#### Slot 迁移流程

```
① 标记迁移状态:
   Target(Node D): CLUSTER SETSLOT 3000 IMPORTING <source-id>
   Source(Node A): CLUSTER SETSLOT 3000 MIGRATING <target-id>

② 逐个迁移 key:
   Source 上执行: MIGRATE target_ip target_port key 0 timeout
   → 把 key 原子地从 Source 移到 Target（Source 上删除）

   ┌─────────┐  MIGRATE key  ┌─────────┐
   │ Source A │─────────────→│ Target D │
   │ slot 3000│  (原子转移)   │ slot 3000│
   │ MIGRATING│              │ IMPORTING│
   └─────────┘              └─────────┘

③ 完成后更新 slot 归属:
   所有节点: CLUSTER SETSLOT 3000 NODE <target-id>
   → Gossip 传播给全集群
```

#### 迁移中 MGET 遇到的重定向

```
情况 1: key 还在 Source → 正常返回

  Client → Source(A): GET key_in_3000
  Source: key 还在本地 → 返回值 ✓

情况 2: key 已迁到 Target → ASK 重定向

  Client → Source(A): GET key_in_3000
  Source: slot 3000 MIGRATING 且 key 不在本地
        → 返回 -ASK 3000 target_ip:port
  Client → Target(D): ASKING     （先发 ASKING 声明）
  Client → Target(D): GET key_in_3000
  Target: 返回值 ✓

  ASK vs MOVED:
    ASK:   临时重定向（迁移中），不更新本地 slot 映射表
    MOVED: 永久重定向（迁移完），更新本地 slot 映射表

情况 3: 迁移完成 → MOVED 重定向

  Client → Source(A): GET key_in_3000
  Source: slot 3000 已不属于我
        → 返回 -MOVED 3000 target_ip:port
  Client: 更新本地映射表 → 后续直接访问 Target(D)
```

#### MGET 在迁移中的完整流程

```
MGET key1(slot100) key2(slot3000) key3(slot8000)

slot 3000 正在从 A 迁移到 D:

  ① 客户端按 slot 分组:
     slot 100  → Node A: MGET key1      → 正常返回 v1
     slot 3000 → Node A: MGET key2      → 可能 ASK 重定向
     slot 8000 → Node B: MGET key3      → 正常返回 v3

  ② 处理 slot 3000 的 ASK:
     收到 -ASK 3000 NodeD_ip:port
     → Client 向 Node D 发送: ASKING + MGET key2
     → Node D 返回 v2

  ③ 合并结果: [v1, v2, v3]

  关键: 客户端需要处理 ASK/MOVED 重定向的重试逻辑
  go-redis 已内置处理:
    MaxRedirects: 8  （最大重定向次数，默认 8）
```

```
迁移中的并发安全:

  MIGRATE 命令是原子的:
    Source: 发送 key 到 Target → Target 写入 → Source 删除
    → 在同一时刻 key 只存在于一个节点
    → 不会出现两个节点同时有同一个 key 的情况

  但是存在短暂的不可用窗口:
    Source 开始 MIGRATE → key 在网络传输中 → Source 和 Target 都没有
    → 此时访问该 key: Source 返回 ASK → Target 也找不到
    → 客户端重试后可获取（MIGRATE 通常 ms 级完成）
```

### Pipeline/MGET 在迁移中的优化

```go
// go-redis ClusterClient 内部处理迁移的伪代码:
func (c *ClusterClient) processPipeline(cmds []Cmd) {
    for attempt := 0; attempt < c.opt.MaxRedirects; attempt++ {
        // ① 按 slot 分组，发送到各节点
        failedCmds := c.mapCmdsByNode(cmds)

        for node, nodeCmds := range failedCmds {
            results := node.execPipeline(nodeCmds)

            for i, result := range results {
                switch {
                case isASK(result):
                    // ② ASK: 发 ASKING 到 target 重试
                    targetNode := parseRedirect(result)
                    targetNode.exec("ASKING")
                    retryResult := targetNode.exec(nodeCmds[i])

                case isMOVED(result):
                    // ③ MOVED: 更新 slot 映射表 + 重试
                    c.updateSlotMapping()
                    // 下次循环会发到正确节点
                }
            }
        }
    }
}
```

### 面试追问

```
Q: MGET 100 个 key 分布在 10 个节点，性能如何？
A: 客户端并行发送 10 个 MGET → 耗时 = 最慢节点的 1 次 RTT
   网络: 1 次 RTT (~1ms 内网)
   比串行 100 次 GET (100ms) 快 100 倍
   比 10 次串行 MGET (10ms) 快 10 倍

Q: 迁移过程中 MGET 会不会有数据不一致？
A: 不会。MIGRATE 是原子操作，key 在任一时刻只存在于一个节点。
   客户端通过 ASK 重定向机制可以找到 key 的实际位置。
   go-redis 内置 MaxRedirects（默认 8）自动重试。

Q: 大规模 MGET（10 万个 key）有什么问题？
A: ① 分组后某个节点的 MGET 可能很大 → 阻塞该节点
   ② 建议分批: 每批 500~1000 个 key
   ③ 或用 Pipeline 代替 MGET（更灵活控制批次大小）

Q: Hash Tag 的副作用？
A: 数据倾斜！如 user:{shop_001}:* 所有 key 都在一个 slot
   → 大商户的数据全挤在一个节点 → 该节点成为瓶颈
   → 解决: 只对真正需要 MGET/事务的 key 用 Hash Tag
```

### 一句话总结

**Redis Cluster 不支持跨 slot 的 MGET，客户端需按 slot 分组后并行发送到各节点再合并结果（go-redis 已内置实现）。** Slot 迁移（扩缩容）时，Source 节点标记 MIGRATING、Target 标记 IMPORTING，逐个 key 原子 MIGRATE。迁移中 MGET 遇到已迁走的 key 会收到 ASK 临时重定向（不更新映射表），客户端先发 ASKING 再重试到 Target 节点获取数据；迁移完成后改为 MOVED 永久重定向（更新映射表）。可用 Hash Tag `{xxx}` 强制相关 key 落同一 slot 避免跨节点，但要注意数据倾斜风险。

