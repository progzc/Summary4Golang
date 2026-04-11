# 系统设计 Q&A

## 目录

- [一、发号器系统如何设计？](#一发号器系统如何设计)
- [二、Snowflake 时钟回拨问题详解及解决方案](#二snowflake-时钟回拨问题详解及解决方案)
- [三、秒级精度延时任务系统如何设计](#三秒级精度延时任务系统如何设计)
- [四、秒杀系统如何设计](#四秒杀系统如何设计)
- [五、短链接服务如何设计](#五短链接服务如何设计)
- [六、分布式锁设计与实现](#六分布式锁设计与实现)
- [七、实时排行榜系统设计](#七实时排行榜系统设计)
- [八、Feed 流系统设计](#八feed-流系统设计)
- [九、IM 即时通讯系统设计](#九im-即时通讯系统设计)
- [十、微服务治理](#十微服务治理)
- [十一、单点登录（SSO）系统设计](#十一单点登录sso系统设计)
- [十二、接口幂等性设计](#十二接口幂等性设计)
- [十三、文件/对象存储系统设计](#十三文件对象存储系统设计)
- [十四、搜索系统设计](#十四搜索系统设计)
- [十五、抢红包系统设计](#十五抢红包系统设计)
- [十六、评论系统设计](#十六评论系统设计)
- [十七、通知中心设计](#十七通知中心设计)
- [十八、灰度发布与 AB 测试](#十八灰度发布与-ab-测试)
- [十九、多方式登录验证系统设计](#十九多方式登录验证系统设计)
---

## 一、发号器系统如何设计？

### 发号器的核心需求

```
发号器 = 分布式唯一 ID 生成器

核心要求:
  1. 全局唯一: 绝对不能重复（多机房、多实例）
  2. 趋势递增: 数据库做主键时 B+ 树插入性能好（避免页分裂）
  3. 高可用:   发号器挂了 → 整个系统无法创建新数据
  4. 高性能:   10w+ QPS，延迟 <1ms
  5. 信息安全: 不能让外部猜到业务量（如订单量/用户量）
```

### 方案总览

| 方案 | 唯一性 | 有序性 | 性能 | 可用性 | 复杂度 | 适用场景 |
|------|--------|--------|------|--------|--------|---------|
| UUID | 全局唯一 | 无序 | 极高（本地） | 极高 | 极低 | 日志、临时标识 |
| 数据库自增 | 单库唯一 | 严格递增 | 低（依赖DB） | 低 | 低 | 小规模、单库 |
| 数据库号段模式 | 全局唯一 | 趋势递增 | 高 | 高 | 中 | 美团 Leaf-Segment |
| Redis INCR | 全局唯一 | 严格递增 | 高 | 中 | 低 | 中等规模 |
| **Snowflake** | 全局唯一 | 趋势递增 | 极高（本地） | 高 | 中 | **大多数场景首选** |
| Snowflake 变种 | 全局唯一 | 趋势递增 | 极高 | 高 | 中高 | 百度 UidGenerator、美团 Leaf |

---

### 方案一：UUID

```
格式: 550e8400-e29b-41d4-a716-446655440000 (128位，32个十六进制字符)

版本:
  UUID v1: 时间戳 + MAC 地址（有序但泄露 MAC）
  UUID v4: 完全随机（最常用）
  UUID v7: 时间戳前缀 + 随机（2022 年新标准，有序）
```

```go
// Go 标准库没有 UUID，但生成很简单
import "crypto/rand"

func GenUUID() string {
    b := make([]byte, 16)
    rand.Read(b)
    b[6] = (b[6] & 0x0f) | 0x40 // version 4
    b[8] = (b[8] & 0x3f) | 0x80 // variant
    return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
        b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
```

| 优点 | 缺点 |
|------|------|
| 本地生成，无网络开销 | 128 位太长，存储和索引开销大 |
| 无中心节点，可用性极高 | 无序，B+ 树插入随机 → 页分裂严重 |
| 实现极简 | 不可读，不含业务信息 |

**结论：** 不适合做数据库主键，适合日志 traceID、临时文件名等不要求有序的场景。

---

### 方案二：数据库自增 ID

```sql
CREATE TABLE id_generator (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    stub CHAR(1) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    UNIQUE KEY uk_stub (stub)
) ENGINE=InnoDB;

-- 每次获取新 ID:
REPLACE INTO id_generator (stub) VALUES ('a');
SELECT LAST_INSERT_ID();  -- 返回自增 ID
```

**多实例避免冲突：设置不同的起始值和步长**

```
实例1: auto_increment_offset=1, auto_increment_increment=3  → 1, 4, 7, 10...
实例2: auto_increment_offset=2, auto_increment_increment=3  → 2, 5, 8, 11...
实例3: auto_increment_offset=3, auto_increment_increment=3  → 3, 6, 9, 12...
```

| 优点 | 缺点 |
|------|------|
| 严格递增 | 性能瓶颈在 DB（每次都要写库） |
| 实现简单 | DB 是单点，挂了就没 ID |
| | 步长固定，扩缩容非常麻烦 |

**结论：** 小规模系统可用，大规模不推荐。

---

### 方案三：数据库号段模式（Leaf-Segment）

#### 原理

```
核心思想: 不是每次都去 DB 取一个 ID，而是一次取一批（号段）

数据库表:
  biz_tag    | max_id  | step  | update_time
  -----------|---------|-------|------------
  order      | 3000    | 1000  | 2026-04-10 10:00:00
  user       | 5000    | 2000  | 2026-04-10 10:00:01

应用启动:
  SELECT max_id, step FROM id_alloc WHERE biz_tag = 'order' FOR UPDATE;
  UPDATE id_alloc SET max_id = max_id + step WHERE biz_tag = 'order';
  → 获得号段 [2001, 3000]

应用内存:
  从 2001 开始发，本地自增到 3000
  → 这 1000 个 ID 完全不需要访问 DB！

号段快用完时（如剩余 < 20%）:
  提前异步去 DB 取下一个号段 [3001, 4000]
  → 双 Buffer 无缝切换，DB 完全不影响发号

  ┌─────────────────────────┐
  │  Buffer 1: [2001, 3000] │ ← 当前使用，发到 2800 时触发预加载
  │  Buffer 2: [3001, 4000] │ ← 提前从 DB 加载好，无缝切换
  └─────────────────────────┘
```

#### Go 实现

```go
package idgen

import (
	"database/sql"
	"sync"
	"sync/atomic"
)

// Segment 号段
type Segment struct {
	currentID int64 // 当前 ID（原子操作）
	maxID     int64 // 号段最大值
}

// LeafSegment 号段模式发号器
type LeafSegment struct {
	mu      sync.Mutex
	db      *sql.DB
	bizTag  string
	step    int64
	current *Segment   // 当前号段
	next    *Segment   // 预加载的下一个号段
	loading int32      // 是否正在加载下一个号段
}

func NewLeafSegment(db *sql.DB, bizTag string, step int64) (*LeafSegment, error) {
	l := &LeafSegment{db: db, bizTag: bizTag, step: step}
	// 启动时加载第一个号段
	seg, err := l.loadSegmentFromDB()
	if err != nil {
		return nil, err
	}
	l.current = seg
	return l, nil
}

// loadSegmentFromDB 从数据库获取一个号段
func (l *LeafSegment) loadSegmentFromDB() (*Segment, error) {
	tx, err := l.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var maxID int64
	err = tx.QueryRow(
		"SELECT max_id FROM id_alloc WHERE biz_tag = ? FOR UPDATE", l.bizTag,
	).Scan(&maxID)
	if err != nil {
		return nil, err
	}

	newMaxID := maxID + l.step
	_, err = tx.Exec(
		"UPDATE id_alloc SET max_id = ? WHERE biz_tag = ?", newMaxID, l.bizTag,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Segment{currentID: maxID, maxID: newMaxID}, nil
}

// NextID 获取下一个 ID
func (l *LeafSegment) NextID() (int64, error) {
	for {
		// 从当前号段取 ID（原子操作，无锁）
		id := atomic.AddInt64(&l.current.currentID, 1)
		if id <= l.current.maxID {
			// 剩余不足 20% 时，异步预加载下一个号段
			remaining := l.current.maxID - id
			if remaining < l.step/5 && atomic.CompareAndSwapInt32(&l.loading, 0, 1) {
				go func() {
					seg, err := l.loadSegmentFromDB()
					if err == nil {
						l.mu.Lock()
						l.next = seg
						l.mu.Unlock()
					}
					atomic.StoreInt32(&l.loading, 0)
				}()
			}
			return id, nil
		}

		// 当前号段用完，切换到下一个
		l.mu.Lock()
		// 双重检查：可能其他 goroutine 已经切换了
		if atomic.LoadInt64(&l.current.currentID) > l.current.maxID {
			if l.next != nil {
				l.current = l.next
				l.next = nil
			} else {
				// 下一个号段还没加载好，同步加载
				seg, err := l.loadSegmentFromDB()
				if err != nil {
					l.mu.Unlock()
					return 0, err
				}
				l.current = seg
			}
		}
		l.mu.Unlock()
	}
}
```

| 优点 | 缺点 |
|------|------|
| DB 访问极少（每 step 次才一次） | 多实例时 ID 不全局有序（各实例号段交错） |
| 双 Buffer 预加载，DB 短暂不可用也能发号 | 号段用完前 DB 必须恢复 |
| 实现相对简单 | 步长需要合理设置 |

---

### 方案四：Redis INCR

```
原理:
  INCR order_id  → 1
  INCR order_id  → 2
  INCR order_id  → 3
  Redis 单线程，INCR 天然原子，天然递增

集群部署（类似数据库步长方案）:
  Redis-1: INCRBY order_id 3 → 1, 4, 7, 10...  (初始值 1)
  Redis-2: INCRBY order_id 3 → 2, 5, 8, 11...  (初始值 2)
  Redis-3: INCRBY order_id 3 → 3, 6, 9, 12...  (初始值 3)
```

| 优点 | 缺点 |
|------|------|
| 性能高（10w+ QPS） | 依赖 Redis，增加运维复杂度 |
| 严格递增 | Redis 持久化有丢数据风险（RDB 间隔内宕机） |
| 实现简单 | 扩缩容改步长麻烦 |

**结论：** 适合中等规模，对严格递增有要求的场景。但要注意 Redis 宕机的 ID 回退问题。

---

### 方案五：Snowflake 雪花算法（最推荐）

#### 原理

```
64 位 ID 结构:

  0 | 00000000 00000000 00000000 00000000 00000000 0 | 00000 00000 | 000000000000
  │ │←─────────────── 41 位 ──────────────────────→│ │←─ 10位 ─→│ │←── 12位 ──→│
  │ │               时间戳                          │ │  机器ID    │ │   序列号    │
  │ │          (毫秒，可用 69 年)                    │ │(1024节点)  │ │(4096/ms)   │
  符号位(0)

各字段说明:
  - 符号位(1位):  固定为 0，保证 ID 为正数
  - 时间戳(41位): 当前毫秒 - 起始时间(epoch)，可用约 69 年
  - 机器ID(10位): 5位数据中心 + 5位机器编号，共 1024 个节点
  - 序列号(12位): 同一毫秒内的自增序号，每毫秒可生成 4096 个 ID

性能:
  单节点: 4096 × 1000 = 409.6 万/秒
  1024 节点: 409.6 万 × 1024 ≈ 41.9 亿/秒
```

#### Go 实现

```go
package idgen

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch             = int64(1704067200000) // 起始时间: 2024-01-01 00:00:00 UTC (毫秒)
	workerIDBits      = 10                   // 机器ID位数 (最大 1024)
	sequenceBits      = 12                   // 序列号位数 (每毫秒最大 4096)
	maxWorkerID       = (1 << workerIDBits) - 1  // 1023
	maxSequence       = (1 << sequenceBits) - 1   // 4095
	workerIDShift     = sequenceBits              // 机器ID左移位数
	timestampShift    = sequenceBits + workerIDBits // 时间戳左移位数
)

// Snowflake 雪花算法发号器
type Snowflake struct {
	mu        sync.Mutex
	workerID  int64 // 机器 ID
	sequence  int64 // 当前毫秒内序列号
	lastMilli int64 // 上一次生成 ID 的毫秒时间戳
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("workerID must be between 0 and %d", maxWorkerID)
	}
	return &Snowflake{workerID: workerID}, nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// 时钟回拨检测
	if now < s.lastMilli {
		diff := s.lastMilli - now
		if diff <= 5 {
			// 小幅回拨（≤5ms），等一等
			time.Sleep(time.Duration(diff) * time.Millisecond)
			now = time.Now().UnixMilli()
		} else {
			return 0, fmt.Errorf("clock moved backwards by %dms", diff)
		}
	}

	if now == s.lastMilli {
		// 同一毫秒内，序列号自增
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 当前毫秒的序列号用完了（4096个），等到下一毫秒
			now = s.waitNextMilli(now)
		}
	} else {
		// 新的一毫秒，序列号归零
		s.sequence = 0
	}

	s.lastMilli = now

	// 组装 64 位 ID
	id := ((now - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}

// waitNextMilli 等待到下一毫秒
func (s *Snowflake) waitNextMilli(current int64) int64 {
	for {
		now := time.Now().UnixMilli()
		if now > current {
			return now
		}
	}
}

// ParseID 解析 Snowflake ID 的各个字段（调试用）
func ParseID(id int64) (timestamp int64, workerID int64, sequence int64) {
	sequence = id & maxSequence
	workerID = (id >> workerIDShift) & maxWorkerID
	timestamp = (id >> timestampShift) + epoch
	return
}
```

#### 时钟回拨问题及解决

```
问题: NTP 时间同步可能导致系统时钟往回跳
  → 新生成的 ID 时间戳 < 之前的 → ID 不再递增，甚至可能重复！

解决方案:

1. 小幅回拨（≤5ms）: 等一等就好
   time.Sleep(diff)

2. 中幅回拨（5ms~几百ms）: 使用扩展位
   从序列号中借几位做"回拨计数"
   每次检测到回拨，回拨计数+1 → ID 仍然唯一

3. 大幅回拨: 报警 + 拒绝服务
   这种情况通常是人工误操作，应该修复时钟而非绕过

4. 美团 Leaf 方案: 用 ZooKeeper 的有序节点记录每个 worker 的最后时间戳
   重启时对比，发现回拨就拒绝启动
```

#### workerID 怎么分配

```
方案1: 手动配置（简单但不灵活）
  每台机器配置文件写死 workerID

方案2: 数据库注册
  启动时 INSERT 一条记录获取自增 ID 作为 workerID
  配合心跳续期，机器下线后回收

方案3: ZooKeeper/etcd 有序节点
  每个实例启动时创建有序临时节点
  节点序号即为 workerID
  实例下线 → 临时节点删除 → ID 可回收

方案4: IP + 端口哈希（无中心依赖）
  workerID = hash(IP + Port) % 1024
  简单但有小概率冲突
```

---

### 方案六：Snowflake 变种——百度 UidGenerator（预生成）

```
思路: Snowflake 在高并发下可能阻塞（等下一毫秒）
     → 提前生成一批 ID 放入 RingBuffer，请求直接取

  ┌──────────────────────────────────────┐
  │       RingBuffer (环形数组)           │
  │  [ID] [ID] [ID] [ID] [  ] [  ] [  ] │
  │   ↑ 消费指针      ↑ 生产指针          │
  │   取 ID 直接读     后台线程持续填充      │
  └──────────────────────────────────────┘

  优势: 消费完全无锁（CAS），生产与消费解耦
  单机 600w+ QPS
```

---

### 各方案选型决策

```
是否需要趋势递增？
  ├── 不需要 → UUID（最简单，traceID/日志/临时标识）
  └── 需要 → 规模多大？
              ├── 小（单库单表）→ 数据库自增
              ├── 中（万级 QPS）→ Redis INCR 或 号段模式
              └── 大（十万级+）→ Snowflake ✓
                                  └── 是否要严格递增？
                                       ├── 趋势递增就行 → 标准 Snowflake
                                       └── 严格递增 → 号段模式（Leaf-Segment）
```

### 使用示例

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// 创建 Snowflake 发号器（workerID = 1）
	sf, _ := NewSnowflake(1)

	// 生成 ID
	for i := 0; i < 5; i++ {
		id, _ := sf.NextID()
		ts, wid, seq := ParseID(id)
		fmt.Printf("ID: %d\n", id)
		fmt.Printf("  时间: %s, 机器: %d, 序号: %d\n\n",
			time.UnixMilli(ts).Format("2006-01-02 15:04:05.000"), wid, seq)
	}
}

// 输出:
// ID: 386752102408192001
//   时间: 2026-04-10 14:30:00.123, 机器: 1, 序号: 1
//
// ID: 386752102408192002
//   时间: 2026-04-10 14:30:00.123, 机器: 1, 序号: 2
// ...
```

### 一句话总结

- **小规模**用数据库自增，**中规模**用号段模式（Leaf-Segment）或 Redis INCR，**大规模首选 Snowflake**
- Snowflake 核心：41 位时间戳 + 10 位机器 ID + 12 位序列号 = 64 位趋势递增 ID，单节点 409 万/秒
- 生产关注三个问题：**时钟回拨**（等待/拒绝/扩展位）、**workerID 分配**（ZK/etcd/DB 注册）、**号段模式的双 Buffer 预加载**（DB 短暂不可用也能发号）

---

## 二、Snowflake 时钟回拨问题详解及解决方案

### 什么是时钟回拨

```
正常情况下系统时钟单调递增:
  t1=1000ms → t2=1001ms → t3=1002ms → ...

时钟回拨: 系统时钟突然往回跳了
  t1=1000ms → t2=1001ms → t3=998ms ← 回拨了 3ms！
```

### 为什么会发生时钟回拨

| 原因 | 说明 | 回拨幅度 |
|------|------|---------|
| **NTP 时间同步** | 机器时钟与 NTP 服务器有偏差，同步时被校正回去 | 通常 < 几百 ms |
| **闰秒调整** | UTC 闰秒（极罕见，2016 年最后一次） | 1 秒 |
| **虚拟机迁移** | VM 热迁移到时钟不同的宿主机 | 可能几秒~几分钟 |
| **人工修改** | 运维手动调整时钟（误操作） | 任意 |
| **硬件故障** | 主板电池没电，BIOS 时钟重置 | 可能几年 |

### 时钟回拨带来什么后果

```
Snowflake ID = 时间戳(41位) | 机器ID(10位) | 序列号(12位)

正常生成:
  t=1000ms: ID = 1000 | 001 | 0001  → 419432448001
  t=1000ms: ID = 1000 | 001 | 0002  → 419432448002
  t=1001ms: ID = 1001 | 001 | 0001  → 419432452097

时钟回拨到 t=1000ms:
  t=1000ms: ID = 1000 | 001 | 0001  → 419432448001  ← 和之前的 ID 重复了！

后果:
  1. ID 重复 → 数据库主键冲突 → 插入失败
  2. ID 不再递增 → 打破"趋势递增"保证 → B+树索引混乱
  3. 业务逻辑依赖 ID 大小判断先后顺序 → 逻辑错误
```

### 解决方案总览

| 方案 | 适用场景 | 复杂度 | 可用性影响 |
|------|---------|--------|-----------|
| 等待追平 | 小幅回拨（≤5ms） | 低 | 短暂阻塞 |
| 拒绝服务 + 报警 | 大幅回拨 | 低 | 拒绝发号 |
| 回拨位扩展 | 中小幅回拨 | 中 | 无影响 |
| 历史最大时间戳 | 所有幅度 | 中 | 无影响 |
| ZK/etcd 记录时间戳 | 所有幅度（含重启） | 高 | 启动可能阻塞 |
| 关闭 NTP 跳变 + slew | 根源预防 | 低 | 无影响 |

---

### 方案一：等待追平（小幅回拨）

```
检测到时钟回拨 ≤ 阈值（如 5ms）:
  → sleep 等差值的时间
  → 醒来后时钟已追平，继续正常发号

时间线:
  正常: t=1000 → t=1001 → t=1002
  回拨: t=1002 → t=1000 (回拨 2ms)
  等待: sleep(2ms)
  恢复: t=1002 → 继续
```

```go
func (s *Snowflake) NextID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    if now < s.lastMilli {
        diff := s.lastMilli - now
        if diff <= 5 {
            // 小幅回拨，等一等
            time.Sleep(time.Duration(diff) * time.Millisecond)
            now = time.Now().UnixMilli()
        } else {
            return 0, fmt.Errorf("clock moved backwards by %dms", diff)
        }
    }
    // ... 正常逻辑
}
```

**优点：** 简单有效，NTP 造成的小幅回拨几乎都能覆盖
**缺点：** 阻塞发号线程；大幅回拨无法处理

---

### 方案二：拒绝服务 + 报警（大幅回拨的兜底）

```
检测到时钟回拨 > 阈值:
  → 拒绝发号，返回错误
  → 同时触发报警，通知运维排查

调用方处理:
  → 重试（等待时钟追平）
  → 或路由到其他健康节点
```

```go
const maxTolerateMs = 5

func (s *Snowflake) NextID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    if now < s.lastMilli {
        diff := s.lastMilli - now
        if diff <= maxTolerateMs {
            time.Sleep(time.Duration(diff) * time.Millisecond)
            now = time.Now().UnixMilli()
        } else {
            // 大幅回拨，拒绝并报警
            alarm(fmt.Sprintf("clock backwards %dms on worker %d", diff, s.workerID))
            return 0, ErrClockBackwards
        }
    }
    // ...
}
```

**注意：** 这是兜底策略，不能单独使用。生产中应配合其他方案减少触发拒绝的概率。

---

### 方案三：回拨位扩展（从序列号中借位）

#### 原理

```
标准 Snowflake:
  时间戳(41) | 机器ID(10) | 序列号(12)

改造后:
  时间戳(41) | 机器ID(10) | 回拨计数(2) | 序列号(10)
                            ↑
                         从序列号借 2 位做回拨计数器

每次检测到时钟回拨: 回拨计数 +1（0→1→2→3）
回拨计数不同 → 即使时间戳和序列号相同，生成的 ID 也不同 → 不重复！

代价: 序列号从 12 位(4096/ms) 缩减为 10 位(1024/ms)
     但可以容忍 3 次回拨（2 位 = 0~3）
```

#### Go 实现

```go
const (
    epochMs           = int64(1704067200000)
    workerIDBits      = 10
    clockBackBits     = 2   // 回拨计数位
    sequenceBits      = 10  // 序列号缩减为 10 位
    maxClockBack      = (1 << clockBackBits) - 1  // 最大回拨次数 3
    maxSequence       = (1 << sequenceBits) - 1    // 1023
    sequenceShift     = 0
    clockBackShift    = sequenceBits
    workerIDShift     = sequenceBits + clockBackBits
    timestampShift    = sequenceBits + clockBackBits + workerIDBits
)

type SnowflakeWithClockBack struct {
    mu            sync.Mutex
    workerID      int64
    sequence      int64
    lastMilli     int64
    clockBackCount int64  // 回拨计数器
}

func (s *SnowflakeWithClockBack) NextID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    if now < s.lastMilli {
        // 检测到回拨
        s.clockBackCount++
        if s.clockBackCount > maxClockBack {
            return 0, fmt.Errorf("clock back count exceeded, backed %dms",
                s.lastMilli-now)
        }
        // 不等待，不拒绝，直接用新的回拨计数继续发号
        s.sequence = 0
        // 注意: 使用 lastMilli 而非 now，保证时间戳不回退
        // now 回拨了，但我们用 lastMilli 继续往前走
    } else if now == s.lastMilli {
        s.sequence = (s.sequence + 1) & maxSequence
        if s.sequence == 0 {
            now = s.waitNextMilli(now)
        }
    } else {
        s.sequence = 0
    }

    s.lastMilli = now

    id := ((now - epochMs) << timestampShift) |
        (s.workerID << workerIDShift) |
        (s.clockBackCount << clockBackShift) |
        s.sequence

    return id, nil
}
```

**优点：** 不等待、不拒绝，回拨时立即可用
**缺点：** 牺牲序列号位数；回拨次数有限（2 位最多 3 次）

---

### 方案四：记录历史最大时间戳（不依赖系统时钟单调性）

#### 原理

```
核心思想: 不信任系统时钟的"当前时间"，只信任"曾经到达过的最大时间"

  维护一个变量 maxTimestamp，只增不减:
    如果 now > maxTimestamp → maxTimestamp = now（正常推进）
    如果 now <= maxTimestamp → 不用 now，继续用 maxTimestamp 发号

  时间线:
    now=1000 → max=1000, 用 1000 发号
    now=1001 → max=1001, 用 1001 发号
    now=999  → max=1001, 继续用 1001 发号（忽略回拨的 now）
    now=1002 → max=1002, 用 1002 发号（恢复正常）

  相当于: 时钟回拨时，ID 的时间戳不回退，序列号继续自增
         直到系统时钟追上 maxTimestamp，恢复正常

  风险: 如果回拨时间很长，会长时间停留在同一个毫秒
       → 序列号 4096/ms 可能不够用 → 需要等待
```

#### Go 实现

```go
type SnowflakeMonotonic struct {
    mu           sync.Mutex
    workerID     int64
    sequence     int64
    maxTimestamp  int64  // 历史最大时间戳，只增不减
}

func (s *SnowflakeMonotonic) NextID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    // 取 now 和 maxTimestamp 的较大值
    if now > s.maxTimestamp {
        s.maxTimestamp = now
        s.sequence = 0
    } else {
        // now <= maxTimestamp: 时钟回拨或同一毫秒
        // 不管回拨多少，继续用 maxTimestamp 发号
        s.sequence++
        if s.sequence > maxSequence {
            // 当前毫秒的序列号用完了，推进到下一毫秒
            s.maxTimestamp++
            s.sequence = 0
        }
    }

    id := ((s.maxTimestamp - epoch) << timestampShift) |
        (s.workerID << workerIDShift) |
        s.sequence

    return id, nil
}
```

**优点：**
- 永不等待、永不拒绝，任何幅度的回拨都能处理
- 不需要借位，保留完整 12 位序列号
- 实现简洁

**缺点：**
- 大幅回拨时 maxTimestamp 会"超前"于真实时间 → ID 中的时间戳不准确
- 长时间回拨 + 高并发可能导致 maxTimestamp 大幅领先真实时间

---

### 方案五：ZK/etcd 持久化时间戳（美团 Leaf 方案）

#### 原理

```
问题: 方案一~四都在内存中记录 lastMilli
     如果进程重启，lastMilli 丢失 → 重启后用回拨的时钟发号 → 可能重复

美团 Leaf 的解决方案:

  1. 每个 worker 在 ZooKeeper 中记录自己的最后发号时间戳
     /snowflake/workers/{workerID}/timestamp = 1712745600123

  2. 定期上报（如每 3 秒）:
     worker → ZK: "我的最新时间戳是 1712745600123"

  3. 进程重启时:
     → 从 ZK 读取上次记录的时间戳
     → 和当前系统时钟比较
     → 如果系统时钟 < ZK 中的时间戳 → 说明有回拨 → 拒绝启动

  4. 启动前还会检查其他所有 worker 的时间戳:
     → 如果自己的时钟比大多数 worker 慢很多 → 说明自己时钟有问题 → 拒绝启动

  流程:
  ┌──────────┐     ┌──────────┐     ┌──────────┐
  │ Worker 1 │     │ Worker 2 │     │ Worker 3 │
  └────┬─────┘     └────┬─────┘     └────┬─────┘
       │                │                │
       │ 上报 ts=1000   │ 上报 ts=1001   │ 上报 ts=999
       ↓                ↓                ↓
  ┌──────────────────────────────────────────┐
  │              ZooKeeper / etcd            │
  │  /worker/1/ts = 1000                    │
  │  /worker/2/ts = 1001                    │
  │  /worker/3/ts = 999                     │
  └──────────────────────────────────────────┘

  Worker 3 重启时:
    系统时钟 = 995 < ZK 记录的 999 → 检测到回拨 → 拒绝启动!
```

#### Go 实现（简化版，用 Redis 代替 ZK）

```go
type SnowflakeWithPersist struct {
    mu        sync.Mutex
    workerID  int64
    sequence  int64
    lastMilli int64
    rdb       *redis.Client
    persistKey string
}

func NewSnowflakeWithPersist(rdb *redis.Client, workerID int64) (*SnowflakeWithPersist, error) {
    s := &SnowflakeWithPersist{
        workerID:   workerID,
        rdb:        rdb,
        persistKey: fmt.Sprintf("snowflake:worker:%d:last_ts", workerID),
    }

    ctx := context.Background()

    // 启动时检查: 从 Redis 读取上次记录的时间戳
    lastTS, err := rdb.Get(ctx, s.persistKey).Int64()
    if err == nil && lastTS > 0 {
        now := time.Now().UnixMilli()
        if now < lastTS {
            return nil, fmt.Errorf(
                "clock backwards detected on startup: now=%d, last=%d, diff=%dms",
                now, lastTS, lastTS-now,
            )
        }
        s.lastMilli = lastTS
    }

    // 启动后台上报
    go s.persistLoop(ctx)

    return s, nil
}

// persistLoop 定期将最新时间戳写入 Redis
func (s *SnowflakeWithPersist) persistLoop(ctx context.Context) {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        s.mu.Lock()
        ts := s.lastMilli
        s.mu.Unlock()

        if ts > 0 {
            s.rdb.Set(ctx, s.persistKey, ts, 24*time.Hour)
        }
    }
}

func (s *SnowflakeWithPersist) NextID() (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    if now < s.lastMilli {
        diff := s.lastMilli - now
        if diff <= 5 {
            time.Sleep(time.Duration(diff) * time.Millisecond)
            now = time.Now().UnixMilli()
        } else {
            return 0, fmt.Errorf("clock backwards %dms", diff)
        }
    }

    if now == s.lastMilli {
        s.sequence = (s.sequence + 1) & maxSequence
        if s.sequence == 0 {
            now = s.waitNextMilli(now)
        }
    } else {
        s.sequence = 0
    }

    s.lastMilli = now

    id := ((now - epoch) << timestampShift) |
        (s.workerID << workerIDShift) |
        s.sequence
    return id, nil
}
```

---

### 方案六：从根源预防——关闭 NTP 跳变

```
NTP 同步有两种模式:

1. Step 模式（跳变）— 默认:
   时钟偏差大时直接跳到正确时间
   → 可能往前跳（回拨）或往后跳
   → 这就是时钟回拨的主要来源！

2. Slew 模式（渐变）— 推荐:
   不跳变，而是微调时钟走速（走快或走慢）
   → 逐渐追上正确时间，过程可能需要几分钟
   → 时钟永远单调递增，不会回拨！

Linux 配置:
  # /etc/ntp.conf 或 /etc/chrony.conf
  # 设置 chrony 使用 slew 模式
  makestep 0 0    # 禁止跳变（参数：阈值 次数，0 0 表示永不跳变）

  # 或使用 ntpd
  tinker panic 0  # 偏差再大也不跳变，用 slew 逐渐追
```

| 模式 | 行为 | Snowflake 影响 |
|------|------|---------------|
| Step（跳变） | 偏差 >128ms 直接跳到正确时间 | 可能时钟回拨 |
| **Slew（渐变）** | 调整走速，逐渐校正 | **不会回拨** |

**注意：** slew 模式下如果偏差很大，追平需要很长时间（1s 偏差 ≈ 2000s 追平）。所以初次部署时要先用 step 对齐一次，之后切 slew。

---

### 生产环境推荐组合策略

```
┌─────────────────────────────────────────────────────┐
│             Snowflake 时钟回拨完整方案                 │
├─────────────────────────────────────────────────────┤
│                                                     │
│  第一层: 根源预防                                     │
│    → NTP 配置 slew 模式，从源头避免回拨               │
│                                                     │
│  第二层: 启动时检查                                   │
│    → 从 ZK/Redis 读取上次时间戳                       │
│    → 系统时钟 < 记录时间戳 → 拒绝启动                 │
│                                                     │
│  第三层: 运行时处理                                   │
│    → 回拨 ≤ 5ms: 等待追平（方案一）                    │
│    → 回拨 > 5ms: 使用历史最大时间戳继续发号（方案四）    │
│    → 回拨 > 500ms: 拒绝发号 + 报警（方案二兜底）       │
│                                                     │
│  第四层: 定期上报                                     │
│    → 每 3s 将 lastMilli 写入 ZK/Redis                │
│    → 进程重启时有据可查                               │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 各方案对比总结

| 方案 | 回拨容忍范围 | 是否阻塞 | 是否拒绝 | 复杂度 | 适用场景 |
|------|------------|---------|---------|--------|---------|
| 等待追平 | ≤ 几ms | 短暂阻塞 | 否 | 低 | NTP 微调导致的小幅回拨 |
| 拒绝 + 报警 | 任意（兜底） | 否 | 是 | 低 | 兜底策略 |
| 回拨位扩展 | ≤ 3 次回拨 | 否 | 否 | 中 | 回拨频繁但幅度小 |
| 历史最大时间戳 | 任意 | 否 | 否 | 中 | **推荐，覆盖面最广** |
| ZK/Redis 持久化 | 含进程重启 | 启动时可能 | 启动时可能 | 高 | 对唯一性要求极高 |
| NTP slew 模式 | 根源预防 | 否 | 否 | 低 | **必须配置** |

### 一句话总结

**生产最佳实践 = NTP slew 模式（根源预防）+ 历史最大时间戳（运行时兜底）+ 持久化上报（重启保护）。** NTP slew 避免 99% 的回拨；历史最大时间戳保证回拨时不等待不拒绝；持久化时间戳防止进程重启后用回拨的时钟发出重复 ID。三层防线叠加，确保万无一失。

---

## 三、秒级精度延时任务系统如何设计

### 什么是延时任务

```
延时任务 = 在指定的未来时间点触发执行的任务

常见场景:
  - 订单超时关闭: 下单后 30 分钟未支付 → 自动取消
  - 延迟消息推送: 用户注册后 24 小时 → 发送引导通知
  - 红包过期回退: 发红包后 24 小时未领取 → 自动退款
  - 会议提醒:     开会前 15 分钟 → 发送提醒
  - 限时活动:     活动到指定时间 → 自动上下架
  - 重试退避:     请求失败后 5s/10s/30s 后重试
```

### 方案总览

| 方案 | 精度 | 吞吐量 | 可靠性 | 复杂度 | 适用场景 |
|------|------|--------|--------|--------|---------|
| 数据库轮询 | 秒级 | 低 | 高（持久化） | 低 | 小规模、简单场景 |
| Redis ZSET | 秒级 | 中高 | 中 | 低 | 中等规模、最常用 |
| Redis Keyspace 通知 | 秒级 | 中 | **低（不可靠）** | 低 | 不推荐生产使用 |
| RocketMQ 延迟消息 | 秒级（固定级别） | 高 | 高 | 低 | 已有 RocketMQ |
| **时间轮（Timing Wheel）** | **ms~秒级** | **极高** | 需持久化 | 中 | **高性能、大规模** |
| Kafka + 延迟 Topic | 秒级 | 高 | 高 | 中高 | 已有 Kafka |

---

### 方案一：数据库轮询（最简单）

#### 原理

```
表结构:
  CREATE TABLE delay_task (
      id         BIGINT PRIMARY KEY,
      task_type  VARCHAR(64),
      task_data  TEXT,
      execute_at DATETIME     NOT NULL,  -- 执行时间
      status     TINYINT      DEFAULT 0, -- 0=待执行 1=执行中 2=完成 3=失败
      created_at DATETIME,
      INDEX idx_execute (status, execute_at)
  );

定时轮询:
  每秒执行:
  SELECT * FROM delay_task
  WHERE status = 0 AND execute_at <= NOW()
  ORDER BY execute_at ASC
  LIMIT 100;

  取出后: status → 1(执行中) → 执行任务 → status → 2(完成)
```

| 优点 | 缺点 |
|------|------|
| 实现极简 | 轮询对 DB 有压力 |
| 天然持久化，不丢任务 | 精度取决于轮询间隔 |
| 支持任意延迟时间 | 大量任务时性能差 |

**适用：** 任务量小（<1w/天）、对精度要求不高（秒级）的场景。

---

### 方案二：Redis ZSET（最常用的中等规模方案）

#### 原理

```
利用 Redis Sorted Set，score 存执行时间戳:

  添加延时任务:
    ZADD delay_queue <execute_timestamp> <task_id>
    例: ZADD delay_queue 1712745600 "order_cancel:12345"

  轮询到期任务（每秒）:
    ZRANGEBYSCORE delay_queue 0 <now_timestamp> LIMIT 0 100
    → 取出所有 score ≤ 当前时间的任务（已到期）

  取出后删除:
    ZREM delay_queue <task_id>

  时间线:
    score: 1000  1005  1010  1015  1020  1025
    task:  [T1]  [T2]  [T3]  [T4]  [T5]  [T6]
                        ↑
                  now = 1010
           ZRANGEBYSCORE 0 1010 → T1, T2, T3（到期了）
```

#### 多消费者防重复消费

```
问题: 多个消费者同时 ZRANGEBYSCORE，取到同一批任务 → 重复执行

解决: Lua 脚本原子 "取出+删除"

local tasks = redis.call('ZRANGEBYSCORE', KEYS[1], '0', ARGV[1], 'LIMIT', '0', '100')
if #tasks > 0 then
    redis.call('ZREM', KEYS[1], unpack(tasks))
end
return tasks

或者用 ZPOPMIN（Redis 5.0+）:
  ZPOPMIN delay_queue 100  → 原子取出并删除 score 最小的 100 个
  但需要在应用层判断是否到期（score ≤ now）
```

#### Go 实现

```go
package delaytask

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Data     string `json:"data"`
	ExecAt   int64  `json:"exec_at"` // 执行时间戳（秒）
}

type DelayQueue struct {
	rdb  *redis.Client
	key  string
}

func NewDelayQueue(rdb *redis.Client, key string) *DelayQueue {
	return &DelayQueue{rdb: rdb, key: key}
}

// Add 添加延时任务
func (q *DelayQueue) Add(ctx context.Context, task Task) error {
	data, _ := json.Marshal(task)
	return q.rdb.ZAdd(ctx, q.key, redis.Z{
		Score:  float64(task.ExecAt),
		Member: string(data),
	}).Err()
}

// 原子取出到期任务的 Lua 脚本
var fetchScript = redis.NewScript(`
local tasks = redis.call('ZRANGEBYSCORE', KEYS[1], '0', ARGV[1], 'LIMIT', '0', ARGV[2])
if #tasks > 0 then
    redis.call('ZREM', KEYS[1], unpack(tasks))
end
return tasks
`)

// Consume 消费循环
func (q *DelayQueue) Consume(ctx context.Context, batchSize int, handler func(Task) error) {
	ticker := time.NewTicker(1 * time.Second) // 每秒轮询
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.poll(ctx, batchSize, handler)
		}
	}
}

func (q *DelayQueue) poll(ctx context.Context, batchSize int, handler func(Task) error) {
	now := time.Now().Unix()
	results, err := fetchScript.Run(ctx, q.rdb,
		[]string{q.key},
		now,
		batchSize,
	).StringSlice()
	if err != nil && err != redis.Nil {
		log.Printf("poll error: %v", err)
		return
	}

	for _, raw := range results {
		var task Task
		if err := json.Unmarshal([]byte(raw), &task); err != nil {
			continue
		}
		if err := handler(task); err != nil {
			// 执行失败，重新入队（可加重试次数限制）
			task.ExecAt = time.Now().Add(5 * time.Second).Unix()
			q.Add(ctx, task)
		}
	}
}

// 使用示例
func Example() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	queue := NewDelayQueue(rdb, "delay:order_cancel")

	ctx := context.Background()

	// 添加任务: 30 分钟后取消订单
	queue.Add(ctx, Task{
		ID:     "order_12345",
		Type:   "order_cancel",
		Data:   `{"order_id": 12345}`,
		ExecAt: time.Now().Add(30 * time.Minute).Unix(),
	})

	// 消费任务
	queue.Consume(ctx, 100, func(task Task) error {
		fmt.Printf("执行任务: %s, data: %s\n", task.ID, task.Data)
		return nil
	})
}
```

| 优点 | 缺点 |
|------|------|
| 实现简单，性能好 | 依赖轮询（1s 间隔 = 最大 1s 误差） |
| Redis 本身高可用 | Redis 宕机可能丢任务（需持久化） |
| 天然排序，查询 O(logN) | 任务量极大时 ZSET 内存开销大 |

**适用：** 中等规模（百万级任务），秒级精度，最常见的生产方案。

---

### 方案三：时间轮（Timing Wheel）— 高性能方案

#### 原理

```
时间轮 = 一个环形数组 + 指针，类似钟表:

  简单时间轮（秒级，60 个槽）:

         0
      59   1
    58       2
   57         3
   56    ↑    4      指针每秒前进一格
   55  指针   5      指向的槽位中的所有任务 → 触发执行
    54       6
      53   7
        ...

  每个槽是一个链表，存放在这一秒到期的任务:
    slot[0] → [Task_A] → [Task_B]
    slot[1] → [Task_C]
    slot[2] → (空)
    slot[3] → [Task_D] → [Task_E] → [Task_F]
    ...

  添加任务: 延迟 5s → 放入 (当前指针 + 5) % 60 的槽位
  触发任务: 指针每秒走一格，执行当前槽位的所有任务
  时间复杂度: 添加 O(1)，触发 O(1)（不需要排序！）
```

#### 多层时间轮（支持更大延迟范围）

```
问题: 单层 60 个槽只能延迟 60 秒，怎么支持几小时/几天？

方案: 多层时间轮（类似时钟的秒/分/时）

  第1层（秒轮）: 60 格，每格 1 秒，覆盖 60 秒
  第2层（分轮）: 60 格，每格 1 分钟，覆盖 60 分钟
  第3层（时轮）: 24 格，每格 1 小时，覆盖 24 小时
  第4层（天轮）: 30 格，每格 1 天，覆盖 30 天

  添加延迟 90 分钟的任务:
    → 放入第2层(分轮)的第 30 格（90-60=30）

  第2层指针走到第 30 格时:
    → 将该格的任务"降级"到第1层
    → 根据剩余秒数放到第1层对应的槽位
    → 第1层指针到达时触发

  ┌─────────┐    降级    ┌─────────┐    降级    ┌─────────┐
  │ 天轮(30)│ ────────→ │ 时轮(24)│ ────────→ │ 分轮(60)│
  └─────────┘           └─────────┘           └─────────┘
                                                   │ 降级
                                              ┌────▼────┐
                                              │ 秒轮(60)│ → 触发执行
                                              └─────────┘
```

#### Go 实现（单层时间轮核心）

```go
package timewheel

import (
	"container/list"
	"sync"
	"time"
)

type TaskFunc func()

type taskEntry struct {
	delay    time.Duration
	circle   int       // 剩余圈数（支持超过一圈的延迟）
	fn       TaskFunc
}

type TimeWheel struct {
	mu         sync.Mutex
	interval   time.Duration   // 每个槽的时间间隔
	slotNum    int             // 槽数量
	slots      []*list.List    // 环形数组，每个槽是一个链表
	currentPos int             // 当前指针位置
	ticker     *time.Ticker
	stopCh     chan struct{}
}

func New(interval time.Duration, slotNum int) *TimeWheel {
	tw := &TimeWheel{
		interval:   interval,
		slotNum:    slotNum,
		slots:      make([]*list.List, slotNum),
		currentPos: 0,
		stopCh:     make(chan struct{}),
	}
	for i := range tw.slots {
		tw.slots[i] = list.New()
	}
	return tw
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *TimeWheel) Stop() {
	close(tw.stopCh)
	tw.ticker.Stop()
}

// AddTask 添加延时任务
func (tw *TimeWheel) AddTask(delay time.Duration, fn TaskFunc) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	// 计算应该放在哪个槽、转几圈
	ticks := int(delay / tw.interval)
	pos := (tw.currentPos + ticks) % tw.slotNum
	circle := ticks / tw.slotNum

	tw.slots[pos].PushBack(&taskEntry{
		delay:  delay,
		circle: circle,
		fn:     fn,
	})
}

func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.stopCh:
			return
		case <-tw.ticker.C:
			tw.tickHandler()
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	slot := tw.slots[tw.currentPos]
	var next *list.Element
	for e := slot.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(*taskEntry)

		if entry.circle > 0 {
			// 还没到，等下一圈
			entry.circle--
			continue
		}

		// 到期了，异步执行任务
		go entry.fn()
		slot.Remove(e)
	}

	// 指针前进
	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
}
```

| 优点 | 缺点 |
|------|------|
| 添加和触发都是 O(1) | 纯内存，重启丢失（需配合持久化） |
| 性能极高，百万级任务无压力 | 实现相对复杂 |
| 精度高（取决于 tick 间隔） | 单机方案，需额外做分布式 |

**适用：** 高性能、大规模场景。Netty 的 HashedWheelTimer、Kafka 的延迟操作都用时间轮。

---

### 方案四：RocketMQ 延迟消息

```
RocketMQ 原生支持延迟消息:

  // 发送延迟消息
  msg.setDelayTimeLevel(3)  // 级别3 = 延迟 10s

  18 个延迟级别:
  1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h

  RocketMQ 5.x 支持任意时间延迟（秒级精度）

内部实现:
  1. 延迟消息先存入 SCHEDULE_TOPIC_XXXX 内部 Topic
  2. 每个延迟级别一个队列
  3. Broker 内部定时扫描到期的消息
  4. 到期后转投到目标 Topic → 消费者可见

优点: 不需要自己实现延迟逻辑，高可靠（Broker 持久化）
缺点: 延迟级别固定（5.x 前），依赖 RocketMQ
```

---

### 方案五：Redis Keyspace 通知（不推荐）

```
原理: 给 Key 设置 TTL，过期时 Redis 发通知

  SET task:order_12345 "data" EX 1800   // 30 分钟后过期
  → 订阅 __keyevent@0__:expired 频道
  → 收到通知 → 执行任务

致命问题:
  1. 不可靠: Redis 的过期通知是"尽力而为"(at-most-once)
     消息可能丢失，没有 ACK 机制
  2. 不精确: 惰性删除 + 定期删除，可能有几秒甚至更大延迟
  3. 不扩展: 大量 Key 过期时通知可能堆积

结论: 不适合生产环境的延时任务。
```

---

### 生产方案：Redis ZSET + 时间轮组合架构

```
大规模延时任务系统架构:

  ┌──────────┐    添加任务    ┌──────────────┐
  │  业务方   │──────────────→│  API Server  │
  └──────────┘               └──────┬───────┘
                                    │ ZADD
                              ┌─────▼──────┐
                              │  Redis ZSET │  ← 持久化存储（保证不丢）
                              │ (多分片)    │
                              └─────┬──────┘
                                    │ 每秒 ZPOPMIN 取出到期任务
                          ┌─────────┼─────────┐
                    ┌─────▼─────┐  ┌▼─────────┐  ┌▼─────────┐
                    │ Worker 1  │  │ Worker 2  │  │ Worker 3  │
                    │ 内存时间轮 │  │ 内存时间轮 │  │ 内存时间轮 │
                    └─────┬─────┘  └─────┬─────┘  └─────┬─────┘
                          │              │              │
                    ┌─────▼──────────────▼──────────────▼─────┐
                    │              执行任务                     │
                    │  发 MQ / 调 HTTP / 执行回调               │
                    └─────────────────────────────────────────┘

分层设计:
  Redis ZSET:  负责持久化 + 粗粒度调度（秒级）
               多分片: delay_queue:{0~15}，hash(taskID) % 16
  内存时间轮:   负责精确触发（ms 级）
               Worker 从 Redis 取出"未来 5 秒内到期"的任务放入时间轮
               时间轮在内存中精确触发

  两层结合:
    Redis 保证不丢（宕机恢复后重新加载）
    时间轮保证精度（不依赖 Redis 轮询频率）
```

### 关键问题及解决

#### 1. 任务不丢（可靠性）

```
  添加任务:
    ZADD + 同时写 MySQL（异步，用于灾备恢复）

  执行确认:
    任务取出后不立即删除，而是标记为"执行中"
    执行成功后才删除
    超时未确认的任务自动重新入队

  状态机:
    待执行 → 执行中 → 完成
                ↓
              超时 → 重试（最多 N 次）→ 失败 → 死信队列
```

#### 2. 多实例防重复消费

```
  方案 A: Redis 分片 + 实例绑定
    16 个分片，每个 Worker 负责固定几个分片
    Worker 挂了 → 其他 Worker 接管其分片（类似 Kafka Consumer Group Rebalance）

  方案 B: Lua 原子取出
    ZPOPMIN 或 Lua 脚本保证取出和删除原子化
    同一个任务只会被一个 Worker 取到
```

#### 3. 大量任务的性能优化

```
  问题: 单个 ZSET 百万任务 → ZRANGEBYSCORE 变慢

  优化:
    1. 分片: 按 taskID hash 分散到 16/64 个 ZSET
    2. 批量: 每次 ZPOPMIN 取 100~500 个，减少 Redis RT
    3. 预取: Worker 提前取出未来 5~10s 的任务放入内存时间轮
    4. Pipeline: 批量操作用 Redis Pipeline 减少网络往返
```

### 各方案选型建议

```
任务量 < 1w/天，延迟精度 > 10s:
  → 数据库轮询（最简单，够用就行）

任务量 1w~100w/天，延迟精度 = 秒级:
  → Redis ZSET（最常用，推荐）

任务量 > 100w/天，延迟精度 = 秒级:
  → Redis ZSET 分片 + 时间轮（大规模）

已有 RocketMQ，延迟级别固定:
  → RocketMQ 延迟消息（零开发成本）

超高性能，单机百万任务:
  → 纯时间轮（但需要自己做持久化）
```

### 一句话总结

**秒级延时任务 = 持久化存储（保证不丢）+ 高效触发（保证精度）。** 中小规模用 Redis ZSET + 每秒轮询最简洁实用；大规模用 Redis ZSET（持久化层）+ 内存时间轮（精确触发层）两层架构。关键保障：Lua 原子取出防重复消费、执行确认机制防任务丢失、分片打散防热点。时间轮的添加和触发都是 O(1)，是高性能延时的核心数据结构。

---

## 四、秒杀系统如何设计

### 秒杀的核心挑战

```
特点: 极短时间内大量用户抢购少量商品

  10 万人同时抢 100 件商品:
    → 瞬时并发 10 万 QPS（平时可能只有 1000）
    → 100 件库存不能超卖（卖了 101 件就亏了）
    → 用户体验要好（不能全部超时/白屏）

三大核心问题:
  ① 高并发: 瞬时流量是平时的 100 倍
  ② 防超卖: 库存扣减必须准确
  ③ 流量削峰: 不能让洪峰直接打到数据库
```

### 整体架构（层层拦截）

```
核心思想: 尽量在前面的环节把请求挡掉，越少请求到数据库越好

  10 万请求
      │
  ┌───▼────────────────┐
  │ ① CDN + 前端限制    │  静态资源走 CDN；前端按钮置灰/倒计时/验证码
  │   过滤 30% 无效请求  │
  └───┬────────────────┘
  7 万请求
      │
  ┌───▼────────────────┐
  │ ② Nginx 限流        │  令牌桶限流/IP 频率限制/负载均衡
  │   过滤 50%           │
  └───┬────────────────┘
  3.5 万请求
      │
  ┌───▼────────────────┐
  │ ③ 网关/接入层        │  用户身份校验/活动资格检查/接口限流
  │   过滤 50%           │
  └───┬────────────────┘
  1.5 万请求
      │
  ┌───▼────────────────┐
  │ ④ Redis 预扣库存     │  原子扣减库存（Lua 脚本），扣完就拒绝
  │   100 件卖完后全拒绝  │  ← 核心！
  └───┬────────────────┘
  100 个请求（抢到的）
      │
  ┌───▼────────────────┐
  │ ⑤ MQ 异步下单        │  抢到的请求发消息，异步创建订单
  │   削峰填谷           │
  └───┬────────────────┘
      │
  ┌───▼────────────────┐
  │ ⑥ 数据库创建订单      │  最终只有 100 个写入请求到 DB
  └────────────────────┘

  10 万 → 7 万 → 3.5 万 → 1.5 万 → 100 → DB
  层层拦截，到数据库的只有 100 个！
```

---

### 一、前端层

```
1. 静态资源 CDN:
   秒杀页面 HTML/CSS/JS/图片 全部走 CDN
   → 秒杀开始瞬间的页面请求不打到服务器

2. 前端限制:
   - 按钮倒计时（秒杀开始前不可点击）
   - 点击后按钮立即置灰（防重复点击）
   - 动态 URL（秒杀开始时才下发真实接口地址，防提前刷接口）
   - 验证码/滑块（削峰 + 防机器人）

3. 请求合并:
   前端随机延迟 0~3s 发送请求（分散流量）
   → 把 0.1s 的洪峰拉长到 3s
```

---

### 二、接入层（Nginx + 网关）

```
Nginx 限流:
  # 令牌桶限流: 每秒 1000 个请求，允许突发 2000
  limit_req_zone $binary_remote_addr zone=seckill:10m rate=10r/s;
  location /api/seckill {
      limit_req zone=seckill burst=20 nodelay;
  }

  # 单 IP 限制: 每个 IP 每秒最多 10 次
  # 超过直接返回 503

网关层:
  - 用户身份校验（未登录直接拒绝）
  - 活动资格检查（黑名单/限购/会员等级）
  - 接口限流（令牌桶/滑动窗口）
  - 风控拦截（异常行为检测：同设备/同 IP 大量请求）
```

---

### 三、Redis 预扣库存（核心环节）

#### 为什么不直接扣数据库

```
数据库扣库存:
  UPDATE inventory SET stock = stock - 1 WHERE product_id = 1 AND stock > 0;

  问题:
    1 万并发打到这条 SQL → 行锁排队等待 → DB 扛不住
    即使 DB 扛住了，RT 也会从 5ms 涨到 500ms → 用户超时

Redis 预扣:
  原子操作，单线程，10 万 QPS 轻松
  100 件库存 0.01s 就扣完 → 之后全部快速拒绝
  DB 只处理 100 个成功订单
```

#### Lua 脚本原子扣库存

```go
// Redis Lua 脚本: 原子检查+扣减库存
var seckillScript = redis.NewScript(`
local stock = tonumber(redis.call('GET', KEYS[1]))
if stock == nil then
    return -1  -- 活动不存在
end
if stock <= 0 then
    return 0   -- 库存不足
end
redis.call('DECR', KEYS[1])
return 1       -- 扣减成功
`)

func TrySeckill(ctx context.Context, rdb *redis.Client, productID string, userID string) (bool, error) {
    // ① 防重复: 检查用户是否已抢过
    dedupKey := fmt.Sprintf("seckill:dedup:%s:%s", productID, userID)
    ok, _ := rdb.SetNX(ctx, dedupKey, "1", 10*time.Minute).Result()
    if !ok {
        return false, nil // 已经抢过了
    }

    // ② 原子扣库存
    stockKey := fmt.Sprintf("seckill:stock:%s", productID)
    result, err := seckillScript.Run(ctx, rdb, []string{stockKey}).Int()
    if err != nil {
        rdb.Del(ctx, dedupKey) // 扣库存失败，删除去重标记
        return false, err
    }

    switch result {
    case 1:
        return true, nil  // 抢到了
    case 0:
        rdb.Del(ctx, dedupKey)
        return false, nil // 库存不足
    default:
        rdb.Del(ctx, dedupKey)
        return false, fmt.Errorf("activity not found")
    }
}
```

#### 库存预热

```
秒杀开始前: 将库存从数据库加载到 Redis

  活动管理后台:
    1. 创建秒杀活动 → DB: product_id=1, stock=100
    2. 活动开始前 5 分钟 → SET seckill:stock:1 100
    3. 秒杀开始 → 前端下发接口地址 → 用户开抢

  防止 Redis 和 DB 库存不一致:
    Redis 是"预判"，最终以 DB 为准
    Redis 扣成功 → 发 MQ 异步创建订单 → 订单服务在 DB 中真正扣库存
    DB 扣库存失败（极端情况）→ 回滚 Redis + 通知用户
```

---

### 四、MQ 异步下单（削峰填谷）

```
Redis 扣库存成功的 100 个请求 → 不直接写 DB → 发到 MQ

  为什么？
    100 个请求如果同时写 DB → 100 并发 → DB 还行
    但如果是 10000 件商品的大促 → 10000 并发写 DB → 可能扛不住
    MQ 削峰: 1 万请求进 MQ → 消费者以 500/s 匀速写 DB → DB 稳如泰山

  用户体验:
    Redis 扣库存成功 → 立即返回"抢购成功，订单生成中..."
    → 前端轮询订单状态 / WebSocket 推送
    → 后台 MQ 消费创建订单 → 订单生成后通知用户支付
```

```go
// 抢购成功后发 MQ
func handleSeckillSuccess(ctx context.Context, productID, userID string) error {
    msg := SeckillOrderMsg{
        ProductID: productID,
        UserID:    userID,
        CreateAt:  time.Now().Unix(),
    }
    body, _ := json.Marshal(msg)
    return producer.Send(&kafka.Message{
        Topic: "seckill_orders",
        Key:   []byte(userID),
        Value: body,
    })
}

// MQ 消费者: 创建订单 + DB 扣库存
func consumeSeckillOrder(msg SeckillOrderMsg) error {
    tx, _ := db.Begin()
    defer tx.Rollback()

    // DB 扣库存（最终一致性的保障）
    result, _ := tx.Exec(`UPDATE inventory SET stock = stock - 1
        WHERE product_id = ? AND stock > 0`, msg.ProductID)
    rows, _ := result.RowsAffected()
    if rows == 0 {
        // 极端: DB 库存已扣完（Redis 和 DB 有微小差异）
        // → 回滚 Redis + 通知用户抢购失败
        rollbackRedis(msg)
        return nil
    }

    // 创建订单
    _, err := tx.Exec(`INSERT INTO orders (id, user_id, product_id, status)
        VALUES (?, ?, ?, 'pending_payment')`,
        generateOrderID(), msg.UserID, msg.ProductID)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

### 五、防超卖的三层保障

```
第一层: Redis 原子扣减（Lua 脚本，GET+判断+DECR 一气呵成）
  → 99.9% 的超卖在这里被挡住

第二层: 数据库乐观锁（WHERE stock > 0）
  → 兜底，Redis 有微小误差时 DB 兜住

第三层: 数据库唯一键约束
  → 同一用户+同一商品只能有一个订单（防一人多次下单）
  → UNIQUE KEY uk_user_product (user_id, product_id, activity_id)

  三层叠加: Redis 挡流量 + DB 保数据正确 + 唯一键防重复
```

---

### 六、其他关键设计

#### 限购

```go
// 每人限购 1 件
// 方案 1: Redis SET NX（在 TrySeckill 中已实现）
// 方案 2: 数据库唯一键

// 扣库存前检查:
dedupKey := fmt.Sprintf("seckill:dedup:%s:%s", productID, userID)
ok, _ := rdb.SetNX(ctx, dedupKey, "1", 10*time.Minute).Result()
if !ok {
    return "您已经抢过了"
}
```

#### 超时未支付回滚

```
用户抢到后 30 分钟未支付 → 自动释放库存

  方案: 延迟消息（参考 mq-qa.md 第九题 + design-qa.md 第三题）
    抢到后发 30min 延迟消息 → 到期检查支付状态
    → 未支付: 取消订单 + Redis INCR 恢复库存 + DB 回滚

  Redis 库存恢复:
    INCR seckill:stock:{productID}  // 原子加回
```

#### 热点数据隔离

```
秒杀商品的 Redis key 是热点 key（所有请求打到同一个 key）:

  问题: 单个 Redis 节点承受 10 万 QPS → 可能成为瓶颈

  解决:
    1. 本地缓存: 库存为 0 后在应用本地标记"已售罄"
       → 后续请求直接返回，不打 Redis
       var soldOut sync.Map
       if _, ok := soldOut.Load(productID); ok {
           return "已售罄"  // 本地缓存拦截
       }

    2. Redis 分片库存:
       100 件库存 → 拆成 10 个 key 各 10 件
       seckill:stock:1:shard_0 = 10
       seckill:stock:1:shard_1 = 10
       ...
       seckill:stock:1:shard_9 = 10
       → 请求按 hash(userID) % 10 分散到不同 key
       → 10 个 Redis 节点分摊压力
```

---

### 完整流程时序

```
用户                前端           网关         Redis          MQ         DB
 │                  │              │             │             │          │
 │ 点击抢购          │              │             │             │          │
 │─────────────────→│              │             │             │          │
 │                  │ 验证码校验    │             │             │          │
 │                  │ 按钮置灰      │             │             │          │
 │                  │─ POST /seckill ────────────→│             │          │
 │                  │              │ 身份校验     │             │          │
 │                  │              │ 限流检查     │             │          │
 │                  │              │──── 扣库存 ──→│             │          │
 │                  │              │             │ Lua 原子扣减 │          │
 │                  │              │             │ stock:100→99 │          │
 │                  │              │←── 成功 ─────│             │          │
 │                  │              │──── 发MQ ───────────────→│          │
 │                  │←── "抢购成功,订单生成中" ────│             │          │
 │←── 展示结果 ──────│              │             │             │          │
 │                  │              │             │             │ 消费消息  │
 │                  │              │             │             │─ 扣库存 ─→│
 │                  │              │             │             │─ 建订单 ─→│
 │                  │              │             │             │          │
 │←── 推送: 订单已生成，请支付 ──────────────────────────────────│          │
```

### 关键数字参考

```
10 万并发秒杀 100 件商品:

  前端层:   过滤 30% → 7 万到 Nginx
  Nginx:    限流过滤 50% → 3.5 万到网关
  网关:     校验+限流 50% → 1.5 万到 Redis
  Redis:    100 件扣完后全拒绝 → 100 个到 MQ
  MQ:       异步消费 → 100 个到 DB
  DB:       只写 100 条订单

  用户体验:
    抢到: 200ms 内返回（Redis RT ~1ms + 网络）
    没抢到: 100ms 内返回"已售罄"（Redis 库存=0 快速拒绝）

  关键耗时:
    Redis Lua 扣库存: ~1ms
    MQ 发送: ~5ms
    DB 创建订单: ~10ms（MQ 消费者异步，用户无感）
```

### 一句话总结

**秒杀系统 = 层层拦截（CDN→Nginx 限流→网关校验→Redis 预扣→MQ 异步下单→DB）+ Redis Lua 原子扣库存（核心防超卖）+ MQ 削峰（保护 DB）。** 10 万请求经过前端/Nginx/网关层层过滤后只有 1.5 万到 Redis，Redis Lua 脚本原子扣减（GET+判断+DECR），100 件扣完后全部快速拒绝，成功的 100 个请求发 MQ 异步创建订单。防超卖三层保障：Redis 原子扣减 + DB 乐观锁（WHERE stock>0）+ 唯一键约束。热点 key 用本地缓存（售罄标记）+ 分片库存分散压力。

---

## 五、短链接服务如何设计

### 需求分析

```
短链服务的核心功能:
  1. 生成: 输入长 URL → 返回短链 (如 https://t.cn/Ab3Cd5)
  2. 跳转: 用户访问短链 → 302 重定向到原始长 URL
  3. 统计: 记录点击次数、来源、设备等

关键指标:
  - 短链越短越好（6~8 位字符）
  - 同一个长 URL 可以生成不同的短链（不同活动需要独立统计）
  - 读远多于写: 读(跳转) QPS 可能是写(生成) 的 100 倍
  - 跳转必须快（< 50ms）
```

### 整体架构

```
  ┌─────────────────────────────────────────────────────────────┐
  │                                                             │
  │   生成短链:                                                  │
  │   Client → API Server → 发号器生成 ID → Base62 编码 → 存 DB  │
  │                                          → 写 Redis 缓存    │
  │                                          → 返回短链           │
  │                                                             │
  │   短链跳转:                                                  │
  │   Browser → Nginx → API Server → 查 Redis(缓存)             │
  │                                    ↓ miss                   │
  │                                  查 DB → 写 Redis            │
  │                                    ↓                        │
  │                                  302 重定向到长 URL           │
  │                                    ↓                        │
  │                                  异步发 MQ 记录点击统计       │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘
```

### 核心问题一：如何生成短码

#### 6 位 Base62 能支持多少短链

```
Base62 字符集: [0-9a-zA-Z] = 62 个字符

  6 位: 62^6 = 568 亿 → 足够（全球短链服务一年约几十亿）
  7 位: 62^7 = 3.5 万亿
  8 位: 62^8 = 218 万亿

结论: 6 位 Base62 足够大多数场景
     追求保险用 7 位
```

#### 方案一：自增 ID + Base62（推荐）

```
流程:
  ① 发号器生成全局唯一自增 ID（如 Snowflake / DB 自增 / 号段模式）
  ② 将 ID 转为 Base62 编码
  ③ 存储映射关系: short_code → long_url

示例:
  ID = 1000000 → Base62("1000000") = "4c92" (4位)
  ID = 56800235584 → Base62 = "Ab3Cd5" (6位)
```

```go
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(num uint64) string {
    if num == 0 {
        return string(base62Chars[0])
    }
    var result []byte
    for num > 0 {
        result = append([]byte{base62Chars[num%62]}, result...)
        num /= 62
    }
    return string(result)
}

func decodeBase62(s string) uint64 {
    var num uint64
    for _, c := range s {
        num *= 62
        switch {
        case c >= '0' && c <= '9':
            num += uint64(c - '0')
        case c >= 'a' && c <= 'z':
            num += uint64(c - 'a' + 10)
        case c >= 'A' && c <= 'Z':
            num += uint64(c - 'A' + 36)
        }
    }
    return num
}
```

```
优点:
  - 短码递增有序 → 数据库写入友好（B+ 树不页分裂）
  - 不会冲突（ID 全局唯一 → 短码全局唯一）
  - 实现简单

缺点:
  - 可预测（竞争对手通过短码推算你的链接量）
  → 解决: ID 生成后做一次简单混淆（如 XOR 一个固定值再 Base62）
```

#### 方案二：哈希取前 N 位

```
流程:
  ① 对长 URL 做 MD5/MurmurHash
  ② 取前 6~8 位 Base62 字符

示例:
  MD5("https://example.com/very/long/path") = "a1b2c3d4e5f6..."
  取前 6 位 → "a1b2c3"

问题: 哈希冲突!
  不同长 URL 可能得到相同的短码
  → 必须检测冲突: 查 DB 是否已存在 → 冲突则加盐重试

  while (db.exists(shortCode)) {
      longURL += randomSalt  // 加盐
      shortCode = hash(longURL)[:6]
  }

  冲突率随数据量增长而上升 → 不如方案一简洁
```

#### 方案三：预生成短码池

```
提前生成大量短码存入 Redis/DB:
  短码池: [Ab3Cd5, Xk9Lm2, Pq7Rs4, ...]

  生成短链时:
    从池中取一个未使用的短码 → 绑定长 URL → 标记已使用

  优点: 短码完全随机不可预测
  缺点: 需要维护短码池（预生成+补充+去重）
```

#### 方案对比

| 方案 | 冲突 | 可预测 | 复杂度 | 推荐 |
|------|------|--------|--------|------|
| **自增 ID + Base62** | 无 | 是（可混淆） | 低 | **首选** |
| 哈希取前 N 位 | 有（需处理） | 否 | 中 | 一般 |
| 预生成短码池 | 无 | 否 | 高 | 安全要求高时 |

---

### 核心问题二：存储设计

#### 数据库表

```sql
CREATE TABLE short_url (
    id          BIGINT UNSIGNED PRIMARY KEY,   -- 自增 ID / 发号器 ID
    short_code  VARCHAR(8) NOT NULL,           -- 短码 (Base62)
    long_url    VARCHAR(2048) NOT NULL,        -- 原始长 URL
    user_id     BIGINT DEFAULT 0,              -- 创建者
    expire_at   DATETIME DEFAULT NULL,         -- 过期时间（可选）
    created_at  DATETIME NOT NULL,
    UNIQUE KEY uk_short_code (short_code),
    KEY idx_long_url (long_url(128))           -- 前缀索引（可选，查重用）
) ENGINE=InnoDB;
```

```
是否需要"同一长 URL 返回同一短码"？

  场景 A: 不需要（推荐）
    每次生成都是新短码 → 不同活动/不同用户的同一 URL 有独立的统计数据
    实现简单，不需要查重

  场景 B: 需要
    先查 long_url 是否已存在 → 存在则返回已有短码
    → 需要对 long_url 建索引（长字段索引开销大）
    → 或用布隆过滤器先判断
```

#### 缓存设计

```
短链跳转是读密集型（读:写 ≈ 100:1），必须用缓存:

  Redis 缓存:
    key:   short:{code}
    value: {long_url}
    TTL:   24h（热门链接命中缓存，冷门链接回源 DB）

  跳转流程:
    ① 查 Redis: GET short:Ab3Cd5
       → 命中: 直接 302 重定向（~1ms）
    ② Redis miss: 查 DB
       → 查到: 写入 Redis + 302 重定向
       → 没查到: 返回 404

  缓存预热:
    新生成的短链立即写入 Redis（写穿透）
    → 首次访问就能命中缓存

  热点链接:
    某个短链被大量分享（如微博热门）→ 单 key QPS 极高
    → 本地缓存（进程内缓存 1 分钟）+ Redis → 两级缓存
```

---

### 核心问题三：跳转实现

#### 301 vs 302

```
301 (Moved Permanently): 永久重定向
  浏览器会缓存 → 下次直接跳转不经过服务器
  → 服务器压力小，但无法统计点击次数！

302 (Found / Temporary Redirect): 临时重定向
  浏览器不缓存 → 每次都经过服务器
  → 每次跳转都能统计 → 绝大多数短链服务用 302

结论: 需要统计 → 用 302
     不需要统计且追求性能 → 用 301
```

```go
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Path[1:] // 去掉前导 /

    // ① 查缓存
    longURL, err := rdb.Get(ctx, "short:"+code).Result()
    if err == redis.Nil {
        // ② 缓存未命中，查 DB
        longURL, err = queryDB(code)
        if err != nil {
            http.NotFound(w, r)
            return
        }
        // 回写缓存
        rdb.Set(ctx, "short:"+code, longURL, 24*time.Hour)
    }

    // ③ 异步记录点击统计（不阻塞跳转）
    go recordClick(code, r)

    // ④ 302 重定向
    http.Redirect(w, r, longURL, http.StatusFound)
}
```

---

### 核心问题四：点击统计

```
统计需求:
  - 总点击次数（PV）
  - 独立访客数（UV）
  - 来源（Referer）
  - 设备/浏览器/地域
  - 时间维度（每小时/每天）
```

#### 异步统计（不阻塞跳转）

```
跳转时同步统计 → 增加延迟 → 用户体验差
正确做法: 异步

  跳转请求 → 302 重定向（立即返回）
               ↓ 异步
          发 MQ / 写 Redis → 后台消费/聚合 → 存统计 DB

方案 A: Redis 计数（实时，简单）
  // PV: 原子自增
  INCR stat:pv:{code}
  INCR stat:pv:{code}:{date}

  // UV: HyperLogLog（近似去重，误差 0.81%，每个 key 仅 12KB）
  PFADD stat:uv:{code} {user_ip}
  PFCOUNT stat:uv:{code}  → 独立访客数

方案 B: MQ + 离线聚合（精确，大规模）
  跳转时发 Kafka 消息: {code, ip, ua, referer, timestamp}
  → Flink/Spark 消费聚合 → 写入 ClickHouse/ES
  → 支持多维度复杂查询
```

```go
func recordClick(code string, r *http.Request) {
    // 快速统计: Redis
    pipe := rdb.Pipeline()
    today := time.Now().Format("20060102")
    pipe.Incr(ctx, "stat:pv:"+code)
    pipe.Incr(ctx, "stat:pv:"+code+":"+today)
    pipe.PFAdd(ctx, "stat:uv:"+code, getClientIP(r))
    pipe.Exec(ctx)

    // 详细统计: MQ 异步
    clickEvent := ClickEvent{
        ShortCode: code,
        IP:        getClientIP(r),
        UserAgent: r.UserAgent(),
        Referer:   r.Referer(),
        Timestamp: time.Now().Unix(),
    }
    producer.Send("click_events", clickEvent)
}
```

---

### 核心问题五：过期与清理

```
短链不能无限存储 → 需要过期机制

  方案 1: 创建时设置 expire_at
    到期后: 访问返回 404
    定时任务: 每天清理过期数据（DELETE WHERE expire_at < NOW()）

  方案 2: 默认永久 + LRU 淘汰
    长时间无人访问的短链 → 从缓存淘汰
    DB 中保留（存储便宜）

  方案 3: 分库归档
    近期数据: 热库（SSD）
    历史数据: 冷库（HDD）/ 对象存储
```

---

### 完整 Go 实现

```go
type ShortURLService struct {
    db       *sql.DB
    rdb      *redis.Client
    idGen    *Snowflake      // 发号器（参考 design-qa.md 第一题）
    producer *kafka.Producer // 统计用
}

// 生成短链
func (s *ShortURLService) Create(ctx context.Context, longURL string) (string, error) {
    // ① 发号器生成唯一 ID
    id, err := s.idGen.NextID()
    if err != nil {
        return "", err
    }

    // ② ID → Base62 短码（可选: XOR 混淆后再编码）
    code := encodeBase62(uint64(id) ^ 0x5A5A5A5A) // 简单混淆

    // ③ 存 DB
    _, err = s.db.ExecContext(ctx,
        `INSERT INTO short_url (id, short_code, long_url, created_at) VALUES (?, ?, ?, NOW())`,
        id, code, longURL)
    if err != nil {
        return "", err
    }

    // ④ 写缓存（写穿透，首次访问即命中）
    s.rdb.Set(ctx, "short:"+code, longURL, 24*time.Hour)

    return fmt.Sprintf("https://t.cn/%s", code), nil
}

// 跳转
func (s *ShortURLService) Redirect(ctx context.Context, code string) (string, error) {
    // ① 查 Redis
    longURL, err := s.rdb.Get(ctx, "short:"+code).Result()
    if err == nil {
        return longURL, nil
    }

    // ② Redis miss → 查 DB
    err = s.db.QueryRowContext(ctx,
        `SELECT long_url FROM short_url WHERE short_code = ? AND (expire_at IS NULL OR expire_at > NOW())`,
        code).Scan(&longURL)
    if err != nil {
        return "", fmt.Errorf("short url not found")
    }

    // ③ 回写缓存
    s.rdb.Set(ctx, "short:"+code, longURL, 24*time.Hour)
    return longURL, nil
}
```

---

### 高可用与扩展

```
读扩展（跳转 QPS 高）:
  Nginx → 多个 API Server（无状态）→ Redis 集群（缓存）→ DB 读从库
  → 线性水平扩展

写扩展（生成 QPS 高）:
  发号器: Snowflake（本地生成，不依赖 DB）→ 无瓶颈
  DB: 按 ID 范围分库分表
      short_url_0 ~ short_url_15（ID % 16）

Redis 热点 key:
  爆款短链被大量转发 → 单 key QPS 极高
  → 本地缓存 + Redis → 两级缓存
  → 或 Redis 读副本（读写分离）

容灾:
  短链服务不可用 → 用户无法打开链接 → 影响大
  → Redis Cluster（3 主 3 从）
  → DB 主从 + 自动切换
  → 多机房部署
```

### 一句话总结

**短链服务 = 发号器生成唯一 ID → Base62 编码为 6 位短码 → 存 DB + 写 Redis 缓存。跳转时先查 Redis（~1ms），miss 查 DB 回写缓存，302 重定向 + 异步记录点击统计（Redis PV/UV + MQ 详细日志）。** 核心设计决策：用自增 ID+Base62 而非哈希（无冲突、有序写入友好）；用 302 而非 301（每次经过服务器才能统计）；读密集型必须缓存（读:写≈100:1）；统计用 HyperLogLog 算 UV（12KB/key，误差 0.81%）+ Kafka 异步聚合到 ClickHouse 做多维分析。

---

## 六、分布式锁设计与实现

### 为什么需要分布式锁

```
单进程: sync.Mutex 就够了（goroutine 之间互斥）
多进程/多实例: Mutex 只在进程内有效，跨进程需要"分布式锁"

场景:
  扣库存: 3 个实例同时读到 stock=1 → 都认为有货 → 都扣减 → 超卖
  定时任务: 3 个实例都部署了同一个 cron → 同一时刻执行 3 次
  幂等控制: 同一个请求打到不同实例 → 各自检查"未处理"→ 执行两次
```

### 分布式锁的三个核心要求

```
① 互斥性: 同一时刻只有一个客户端持有锁
② 无死锁: 即使持锁客户端崩溃，锁也能自动释放（TTL 过期）
③ 解铃还须系铃人: 只有加锁的客户端能解锁（防止误删别人的锁）
```

### 方案总览

| 方案 | 原理 | 性能 | 可靠性 | 适用场景 |
|------|------|------|--------|---------|
| **Redis SET NX** | 单节点原子操作 | 极高 | 中（主从切换可能丢锁） | **大多数场景首选** |
| Redis Redlock | 多节点多数派 | 高 | 高 | 对可靠性要求高 |
| **etcd/ZooKeeper** | 分布式一致性协议 | 中 | **极高** | 强一致性要求 |
| MySQL | 行锁/唯一键 | 低 | 高 | 已有 MySQL，QPS 低 |

---

### 一、Redis 分布式锁（最常用）

#### 基本实现：SET NX EX

```
加锁:
  SET lock_key unique_value NX EX 30

  NX: 只在 key 不存在时才设置（互斥）
  EX 30: 30 秒自动过期（防死锁）
  unique_value: 唯一标识（如 UUID），用于安全解锁

  返回 OK → 加锁成功
  返回 nil → 锁被别人持有 → 加锁失败

解锁（必须用 Lua 脚本保证原子性）:
  // 先检查是不是自己的锁，再删除
  if redis.call("GET", KEYS[1]) == ARGV[1] then
      return redis.call("DEL", KEYS[1])
  else
      return 0
  end

  为什么不能直接 DEL？
    A 加锁(value=A) → A 超时自动释放 → B 加锁(value=B)
    → A 执行 DEL → 删掉了 B 的锁！
    → 必须检查 value 是不是自己的
```

#### Go 完整实现

```go
package distlock

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
)

var (
    ErrLockFailed  = errors.New("failed to acquire lock")
    ErrLockNotHeld = errors.New("lock not held")
)

type RedisLock struct {
    rdb    *redis.Client
    key    string
    value  string        // 唯一标识（防误删）
    expiry time.Duration // 锁过期时间
}

func NewRedisLock(rdb *redis.Client, key string, expiry time.Duration) *RedisLock {
    return &RedisLock{
        rdb:    rdb,
        key:    "lock:" + key,
        value:  uuid.New().String(),
        expiry: expiry,
    }
}

// Lock 加锁（带重试）
func (l *RedisLock) Lock(ctx context.Context) error {
    // 重试 3 次，每次间隔 100ms
    for i := 0; i < 3; i++ {
        ok, err := l.rdb.SetNX(ctx, l.key, l.value, l.expiry).Result()
        if err != nil {
            return err
        }
        if ok {
            return nil // 加锁成功
        }
        time.Sleep(100 * time.Millisecond)
    }
    return ErrLockFailed
}

// Lua 脚本: 原子检查+删除
var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`)

// Unlock 解锁
func (l *RedisLock) Unlock(ctx context.Context) error {
    result, err := unlockScript.Run(ctx, l.rdb, []string{l.key}, l.value).Int()
    if err != nil {
        return err
    }
    if result == 0 {
        return ErrLockNotHeld // 锁不是自己的（已过期被别人拿了）
    }
    return nil
}

// 使用
func deductStock(ctx context.Context, rdb *redis.Client, productID string) error {
    lock := NewRedisLock(rdb, "stock:"+productID, 10*time.Second)

    if err := lock.Lock(ctx); err != nil {
        return err // 获取锁失败
    }
    defer lock.Unlock(ctx) // 确保释放

    // 临界区: 扣库存
    return doDeductStock(productID)
}
```

---

### 二、锁续期（Watchdog 看门狗）

#### 问题：锁过期但业务没执行完

```
A 加锁(TTL=10s) → 业务执行需要 15s → 10s 时锁自动过期
→ B 拿到锁 → A 和 B 同时执行临界区 → 互斥性被破坏！

  时间线:
  0s          10s            15s
  │── A 加锁 ──│── A 锁过期 ──│── A 执行完 ──
  │           │── B 加锁 ──│── B 执行中 ──
  │           │  A 和 B 同时在临界区!
```

#### 解决：Watchdog 自动续期

```
加锁后启动一个后台 goroutine（看门狗）:
  每隔 TTL/3 检查一次 → 如果锁还是自己的 → 续期 TTL

  A 加锁(TTL=30s)
  │ 10s: watchdog 续期 → TTL 重置为 30s
  │ 20s: watchdog 续期 → TTL 重置为 30s
  │ 25s: 业务执行完 → 解锁 → watchdog 停止
  → 整个过程锁不会过期

  如果 A 崩溃:
  → watchdog 也随之停止 → 不再续期 → 锁到期自动释放
  → 不会死锁 ✓
```

```go
// Lua 脚本: 续期（只续自己的锁）
var renewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
    return 0
end
`)

// LockWithWatchdog 加锁 + 自动续期
func (l *RedisLock) LockWithWatchdog(ctx context.Context) (context.CancelFunc, error) {
    if err := l.Lock(ctx); err != nil {
        return nil, err
    }

    // 启动看门狗
    watchCtx, cancel := context.WithCancel(ctx)
    go l.watchdog(watchCtx)

    // 返回 cancel 函数，解锁时调用
    return func() {
        cancel()             // 停止看门狗
        l.Unlock(context.Background())
    }, nil
}

func (l *RedisLock) watchdog(ctx context.Context) {
    ticker := time.NewTicker(l.expiry / 3) // 每 TTL/3 续期一次
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return // 业务完成，停止续期
        case <-ticker.C:
            renewScript.Run(ctx, l.rdb, []string{l.key},
                l.value, l.expiry.Milliseconds())
        }
    }
}

// 使用
func criticalSection(ctx context.Context, rdb *redis.Client) error {
    lock := NewRedisLock(rdb, "my-resource", 30*time.Second)
    unlock, err := lock.LockWithWatchdog(ctx)
    if err != nil {
        return err
    }
    defer unlock() // 停止看门狗 + 解锁

    // 业务逻辑（执行多久都不怕锁过期）
    return slowBusinessLogic()
}
```

---

### 三、Redis 单节点锁的问题

```
问题: Redis 主从切换时可能丢锁

  ① A 在 Master 上加锁成功
  ② Master 宕机，锁数据还没同步到 Slave（异步复制）
  ③ Slave 提升为新 Master
  ④ B 在新 Master 上加锁成功（因为锁数据丢了）
  → A 和 B 同时持有锁！

  这是 Redis 异步复制的本质问题，单节点方案无法完全解决
```

### 四、Redlock 算法（多节点方案）

```
Redis 作者 antirez 提出的多节点分布式锁:

  前提: 部署 N 个独立的 Redis 实例（奇数个，如 5 个，非主从!）

  加锁流程:
    ① 记录开始时间 T1
    ② 依次向 5 个 Redis 实例发送 SET NX EX（相同 key 和 value）
    ③ 如果在 >= N/2+1 = 3 个实例上加锁成功 → 加锁成功
       且 总耗时 < 锁的 TTL → 锁有效
    ④ 否则 → 加锁失败 → 向所有实例发送 DEL 释放

  解锁:
    向所有 5 个实例发送 Lua 脚本解锁

  为什么安全？
    即使 1~2 个 Redis 宕机 → 另外 3 个仍然持有锁 → 不会被别人获取
    多数派成功 = 锁有效

  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │ Redis 1 │ │ Redis 2 │ │ Redis 3 │ │ Redis 4 │ │ Redis 5 │
  │  加锁✓  │ │  加锁✓  │ │  超时✗  │ │  加锁✓  │ │  宕机✗  │
  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
  3 个成功 >= N/2+1=3 → 加锁成功 ✓
```

```
Redlock 的争议（Martin Kleppmann vs antirez）:

  反对方（Kleppmann）:
    - 依赖时钟: 如果某个 Redis 时钟跳跃 → TTL 提前过期 → 锁失效
    - GC 停顿: 客户端拿到锁后 GC 暂停 → 锁过期 → 另一个客户端拿到锁
    - 不如用 ZooKeeper/etcd（有共识算法保证）

  支持方（antirez）:
    - 实际中时钟偏差极小
    - 比单节点 Redis 锁可靠得多
    - 性能远高于 ZooKeeper

  结论:
    大多数业务场景 → 单节点 Redis + 看门狗 就够了
    极端可靠性要求 → etcd / ZooKeeper
    Redlock 是折中方案，实际生产中用得不多
```

---

### 五、etcd 分布式锁（强一致性方案）

```
etcd 基于 Raft 共识算法 → 天然强一致

  原理:
    1. 创建一个带租约（Lease）的 key（类似 Redis 的 TTL）
    2. 多个客户端竞争创建同一个 key
    3. Revision 最小的客户端获得锁（etcd 事务保证）
    4. 其他客户端 Watch 前一个 Revision 的 key → 释放后通知

  优势:
    - Raft 保证强一致（不存在 Redis 的主从丢锁问题）
    - Watch 机制（不需要轮询）
    - Lease 自动续期

  劣势:
    - 性能低于 Redis（Raft 写入需要多数派确认）
    - 部署复杂（需要 etcd 集群）
```

```go
// Go etcd 分布式锁（使用官方 concurrency 包）
import (
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

func etcdLockExample() error {
    cli, _ := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
    defer cli.Close()

    // 创建会话（带 TTL 的租约，自动续期）
    session, _ := concurrency.NewSession(cli, concurrency.WithTTL(10))
    defer session.Close()

    // 创建锁
    mutex := concurrency.NewMutex(session, "/locks/my-resource")

    // 加锁（阻塞直到获取）
    if err := mutex.Lock(context.Background()); err != nil {
        return err
    }
    defer mutex.Unlock(context.Background())

    // 临界区
    return doBusiness()
}
```

---

### 六、MySQL 分布式锁（兜底方案）

```sql
-- 方案 1: 唯一键
CREATE TABLE distributed_lock (
    lock_key VARCHAR(64) PRIMARY KEY,
    owner VARCHAR(64) NOT NULL,
    expire_at DATETIME NOT NULL
);

-- 加锁
INSERT INTO distributed_lock (lock_key, owner, expire_at)
VALUES ('stock:1001', 'instance-A', NOW() + INTERVAL 30 SECOND);
-- 成功 → 获得锁; 重复键失败 → 锁被占用

-- 解锁
DELETE FROM distributed_lock WHERE lock_key = 'stock:1001' AND owner = 'instance-A';

-- 清理过期锁（定时任务）
DELETE FROM distributed_lock WHERE expire_at < NOW();

-- 方案 2: SELECT FOR UPDATE
BEGIN;
SELECT * FROM distributed_lock WHERE lock_key = 'stock:1001' FOR UPDATE;
-- 获得行锁 → 其他事务阻塞等待
-- 执行业务...
COMMIT; -- 释放锁
```

```
MySQL 锁的问题:
  - 性能差（DB 写入 ~5ms vs Redis ~1ms）
  - DB 是瓶颈（锁请求多了 DB 扛不住）
  - 没有 Watch/通知机制，只能轮询
  适用: QPS 低 + 不想引入 Redis/etcd 的简单场景
```

---

### 各方案选型

```
大多数场景（95%）:
  → Redis SET NX + Lua 解锁 + Watchdog 续期
  → 简单高效，单节点就够

对可靠性要求高（金融/支付）:
  → etcd 分布式锁（Raft 强一致，不丢锁）
  → 或 ZooKeeper（临时顺序节点）

已有 Redis 集群 + 担心主从切换:
  → Redlock（5 节点多数派，比单节点可靠）
  → 但不如 etcd 可靠

QPS 低 + 不想加中间件:
  → MySQL 行锁 / 唯一键

性能排序:  Redis > Redlock > etcd > MySQL
可靠性排序: etcd > Redlock > Redis > MySQL
```

### 一句话总结

**Redis 分布式锁 = SET key value NX EX（原子加锁）+ Lua 脚本检查 value 后 DEL（安全解锁）+ Watchdog 每 TTL/3 续期（防业务未完成锁过期）。** 单节点 Redis 锁的问题是主从切换可能丢锁（异步复制），Redlock 用 5 个独立 Redis 多数派解决但有争议，etcd 基于 Raft 强一致是最可靠方案但性能低于 Redis。生产中 95% 场景用 Redis 单节点 + Watchdog 就够了，金融级强一致用 etcd。核心三要素：互斥（NX）、防死锁（EX 过期）、防误删（unique value + Lua 检查）。

---

## 七、实时排行榜系统设计

### 核心方案：Redis ZSET

```
Redis Sorted Set 天然就是排行榜数据结构:

  底层: 跳表(SkipList) + 哈希表(HashTable)
    跳表: 按 score 有序排列 → 范围查询 O(logN)
    哈希表: member→score 映射 → 点查 O(1)

  关键命令:
    ZADD    key score member    → 添加/更新分数    O(logN)
    ZSCORE  key member          → 查分数           O(1)
    ZREVRANK key member         → 查排名(降序)      O(logN)
    ZREVRANGE key start stop    → 查 Top N(降序)    O(logN + M)
    ZINCRBY key increment member → 增加分数         O(logN)
    ZCARD   key                  → 总人数           O(1)

  百万用户: logN = log(1000000) ≈ 20 次比较 → 每个操作 < 1ms
```

### Go 完整实现

```go
type Leaderboard struct {
    rdb *redis.Client
    key string
}

func New(rdb *redis.Client, name string) *Leaderboard {
    return &Leaderboard{rdb: rdb, key: "leaderboard:" + name}
}

// 更新分数（覆盖式）
func (lb *Leaderboard) UpdateScore(ctx context.Context, userID string, score float64) error {
    return lb.rdb.ZAdd(ctx, lb.key, redis.Z{Score: score, Member: userID}).Err()
}

// 增加分数（累加式）
func (lb *Leaderboard) IncrScore(ctx context.Context, userID string, delta float64) (float64, error) {
    return lb.rdb.ZIncrBy(ctx, lb.key, delta, userID).Result()
}

// 查排名 + 分数（Pipeline 一次 RT）
func (lb *Leaderboard) GetRank(ctx context.Context, userID string) (rank int64, score float64, err error) {
    pipe := lb.rdb.Pipeline()
    rankCmd := pipe.ZRevRank(ctx, lb.key, userID)
    scoreCmd := pipe.ZScore(ctx, lb.key, userID)
    _, err = pipe.Exec(ctx)
    if err != nil {
        return -1, 0, err
    }
    return rankCmd.Val() + 1, scoreCmd.Val(), nil
}

type RankItem struct {
    UserID string
    Score  float64
    Rank   int64
}

// 查 Top N
func (lb *Leaderboard) GetTopN(ctx context.Context, n int64) ([]RankItem, error) {
    results, err := lb.rdb.ZRevRangeWithScores(ctx, lb.key, 0, n-1).Result()
    if err != nil {
        return nil, err
    }
    items := make([]RankItem, len(results))
    for i, z := range results {
        items[i] = RankItem{UserID: z.Member.(string), Score: z.Score, Rank: int64(i) + 1}
    }
    return items, nil
}

// 查看我前后 N 名
func (lb *Leaderboard) GetAroundMe(ctx context.Context, userID string, n int64) ([]RankItem, error) {
    rank, _, err := lb.GetRank(ctx, userID)
    if err != nil {
        return nil, err
    }
    start := rank - n
    if start < 1 {
        start = 1
    }
    results, err := lb.rdb.ZRevRangeWithScores(ctx, lb.key, start-1, rank+n-1).Result()
    if err != nil {
        return nil, err
    }
    items := make([]RankItem, len(results))
    for i, z := range results {
        items[i] = RankItem{UserID: z.Member.(string), Score: z.Score, Rank: start + int64(i)}
    }
    return items, nil
}
```

### 进阶问题

#### 问题一：分数相同怎么排？

```
ZSET score 相同时按 member 字典序（不是按时间）
→ 需求: 同分时先达到的排前面

方案: 时间戳编码进 score 小数部分
  score = float64(points) + (1.0 - float64(timestamp)/1e16)

  用户A: 1000 分, 10:00:00 → score = 1000.999827...
  用户B: 1000 分, 10:05:00 → score = 1000.999826...
  → A 的 score > B → A 排在 B 前面 ✓
```

```go
func makeScore(points int64, t time.Time) float64 {
    return float64(points) + (1.0 - float64(t.UnixMilli())/1e16)
}
```

#### 问题二：多维度排行榜（日榜/周榜/月榜/总榜）

```
每个维度一个 ZSET key，Pipeline 批量更新:

  pipe := rdb.Pipeline()
  pipe.ZIncrBy(ctx, "leaderboard:daily:20260411", delta, userID)
  pipe.ZIncrBy(ctx, "leaderboard:weekly:2026W15", delta, userID)
  pipe.ZIncrBy(ctx, "leaderboard:monthly:202604", delta, userID)
  pipe.ZIncrBy(ctx, "leaderboard:total", delta, userID)
  pipe.Exec(ctx)

  日榜设 TTL 自动过期: rdb.Expire(ctx, key, 48*time.Hour)
```

#### 问题三：百万/千万级用户排名

```
百万级:
  单 ZSET ~90MB，ZADD/ZRANK ~1ms → 完全够用

千万级:
  单 ZSET ~900MB，性能仍 OK → 注意大 key 监控

亿级（极端）:
  → 分桶计数近似排名:
    分数 0~10000 分成 10000 个桶
    bucket[score] = 该分数的人数
    排名 = 比你分高的所有桶人数之和
  → O(1)，适合"你超过了 XX% 的用户"
```

#### 问题四：持久化

```
方案 1: Redis RDB + AOF（可能丢几秒，可接受）
方案 2: 定期同步到 DB（灾备恢复用）
方案 3: 双写 Redis + 异步写 DB
```

### 一句话总结

**实时排行榜用 Redis ZSET 天然实现——底层跳表保证 ZADD/ZRANK 都是 O(logN)，百万用户每个操作 <1ms。** 核心操作：ZADD 更新分数、ZREVRANK 查排名、ZREVRANGE 查 Top N、ZINCRBY 累加分数。同分排序把时间戳编码进 score 小数部分。多维度榜单（日/周/月/总）每个维度一个 ZSET key，Pipeline 批量更新。百万级单 ZSET 足够（~90MB），千万级仍可（~900MB），亿级用分桶计数近似排名。

---

## 八、Feed 流系统设计

### 什么是 Feed 流

```
Feed 流 = 按时间线展示关注对象动态的信息流

  微信朋友圈: 你看到的是好友发的动态（按时间倒序）
  微博首页:   你看到的是关注的人发的微博
  抖音关注页: 你看到的是关注的博主发的视频
  GitHub:    你看到的是关注仓库/用户的动态

核心问题:
  用户 A 关注了 500 人，每人每天发 5 条动态
  → A 的 Feed 流 = 从 2500 条动态中按时间排序取最新的 20 条
  → 如何高效生成这个列表？
```

### 三种模型：推、拉、推拉结合

#### 模型一：拉模式（读扩散）

```
发布时: 只写自己的发件箱（outbox）
读取时: 实时从所有关注人的发件箱拉取 → 合并排序

  发布动态:
    用户 B 发动态 → 写入 B 的发件箱: outbox:B = [post_1, post_2, ...]

  读取 Feed:
    用户 A 关注了 B, C, D
    → 查 outbox:B 最新 N 条
    → 查 outbox:C 最新 N 条
    → 查 outbox:D 最新 N 条
    → 内存中合并排序取 Top 20

  ┌───────┐       ┌─────────────┐
  │ 用户A  │─读──→│ outbox:B    │──┐
  │ Feed  │─读──→│ outbox:C    │──┼→ 合并排序 → Feed
  │       │─读──→│ outbox:D    │──┘
  └───────┘       └─────────────┘
```

```
优点:
  - 写入极简（只写一份到自己的发件箱）
  - 存储省（动态只存一份）
  - 取关/删除动态立即生效（下次拉取就没了）

缺点:
  - 读取时 N 次查询（N = 关注数）→ 关注 500 人 = 500 次查询
  - 每次刷新都要实时合并排序 → 读延迟高
  - 关注数越多读越慢

适用:
  关注数少（< 100）或对实时性要求不高的场景
  如 GitHub 动态、RSS 阅读器
```

#### 模型二：推模式（写扩散）

```
发布时: 将动态推送到每个粉丝的收件箱（inbox）
读取时: 直接读自己的收件箱（已排好序）

  发布动态:
    用户 B 发动态（B 有 10 万粉丝）
    → 遍历 B 的粉丝列表
    → 写入每个粉丝的收件箱:
      inbox:粉丝1 ← post_id
      inbox:粉丝2 ← post_id
      ...
      inbox:粉丝10万 ← post_id

  读取 Feed:
    用户 A 刷 Feed
    → 直接读 inbox:A → 已按时间排好序 → 直接返回

  ┌───────┐  发布  ┌──→ inbox:粉丝1
  │ 用户B  │──────→├──→ inbox:粉丝2
  │ 发动态 │       ├──→ inbox:粉丝3
  └───────┘       └──→ ... (10万个)
```

```
优点:
  - 读取极快（直接读 inbox，一次查询）
  - 用户体验好（Feed 已经准备好了）

缺点:
  - 写放大严重: 10 万粉丝 = 写 10 万次
  - 大 V 发动态 → 需要写入百万/千万粉丝的 inbox → 延迟高
  - 存储爆炸: 每个粉丝 inbox 都存一份（冗余）
  - 取关后需要清理 inbox 中的历史数据

适用:
  粉丝数少（< 1000）的场景
  如微信朋友圈（好友上限 5000，实际平均 200~300）
```

#### 模型三：推拉结合（最佳方案）

```
核心思想: 按用户活跃度和粉丝数分别处理

  普通用户（粉丝 < 阈值，如 1 万）: 推模式
    发动态 → 推送到所有粉丝 inbox（粉丝少，写入快）

  大 V（粉丝 > 阈值）: 拉模式
    发动态 → 只写自己的 outbox（不推送）

  读取 Feed 时:
    ① 读自己的 inbox（普通关注人的动态，已推送过来）
    ② 查关注的大 V 列表 → 拉取各大 V 的 outbox
    ③ 两部分合并排序 → 返回 Feed

  ┌───────┐  读 inbox(普通人动态已在里面)
  │ 用户A  │──→ inbox:A = [post_3, post_5, post_7]
  │ Feed  │
  │       │  拉 outbox(大V动态实时拉)
  │       │──→ outbox:大V_X 最新N条
  │       │──→ outbox:大V_Y 最新N条
  └───────┘
       ↓
    合并排序 → Feed
```

```
大 V 阈值怎么定？
  微博: ~50 万粉丝以上走拉模式
  通常: 1 万~10 万粉丝是分界线（根据写入延迟和存储成本调整）

为什么微信朋友圈可以用纯推？
  好友上限 5000 人 → 每人发一条最多推 5000 次 → 写入量可控
  且朋友圈是双向好友 → 粉丝数 = 好友数 ≤ 5000
```

### 三种模型对比

| 维度 | 拉模式 | 推模式 | 推拉结合 |
|------|--------|--------|---------|
| **写入开销** | 低（写一份） | 高（写 N 份） | 中 |
| **读取开销** | 高（读 N 份合并） | 低（读 inbox） | 中 |
| **存储开销** | 低 | 高（冗余 N 份） | 中 |
| **大 V 友好** | 是 | **否（写爆炸）** | 是 |
| **实时性** | 低（读时才聚合） | 高（预先推送） | 高 |
| **取关生效** | 立即 | 需清理 inbox | 立即（大 V）/需清理（普通） |
| **典型产品** | GitHub | 微信朋友圈 | **微博、抖音、Twitter** |

---

### 数据存储设计

#### 发件箱（Outbox）

```sql
-- 每个用户的发件箱（所有用户都有）
CREATE TABLE user_post (
    post_id     BIGINT PRIMARY KEY,      -- Snowflake ID（含时间信息）
    user_id     BIGINT NOT NULL,
    content     TEXT,
    media_urls  JSON,
    created_at  DATETIME NOT NULL,
    INDEX idx_user_time (user_id, created_at DESC)
);

-- 或 Redis ZSET (score = 时间戳):
-- outbox:{user_id} → ZSET [post_id1(score=ts1), post_id2(score=ts2)]
```

#### 收件箱（Inbox）— 推模式用

```
方案 A: Redis ZSET（推荐，高频读取快）
  key:   inbox:{user_id}
  score: 动态的时间戳
  member: post_id

  ZADD inbox:user_A 1712745600 "post_12345"
  ZREVRANGE inbox:user_A 0 19  → 取最新 20 条 post_id
  → 再批量查 post 详情

  inbox 大小限制: 只保留最近 1000 条（ZREMRANGEBYRANK 裁剪）

方案 B: 数据库表
  CREATE TABLE feed_inbox (
      user_id    BIGINT NOT NULL,
      post_id    BIGINT NOT NULL,
      created_at DATETIME NOT NULL,
      PRIMARY KEY (user_id, post_id),
      INDEX idx_user_time (user_id, created_at DESC)
  );
  → 读取: SELECT post_id FROM feed_inbox WHERE user_id=? ORDER BY created_at DESC LIMIT 20
```

#### 关注关系

```sql
CREATE TABLE user_follow (
    follower_id  BIGINT NOT NULL,  -- 粉丝
    followee_id  BIGINT NOT NULL,  -- 被关注者
    created_at   DATETIME,
    PRIMARY KEY (follower_id, followee_id),
    INDEX idx_followee (followee_id)  -- 查某人的粉丝列表
);

-- Redis 缓存（高频查询）:
-- following:{user_id} → SET [followee_1, followee_2, ...]  关注列表
-- followers:{user_id} → SET [follower_1, follower_2, ...]  粉丝列表
-- follower_count:{user_id} → 粉丝数（判断推/拉的阈值）
```

---

### 推送流程（写扩散部分）

```go
const bigVThreshold = 10000 // 大 V 粉丝阈值

func PublishPost(ctx context.Context, userID int64, post *Post) error {
    // ① 写入发件箱（所有人都写）
    if err := savePost(ctx, post); err != nil {
        return err
    }

    // ② 判断是否大 V
    followerCount, _ := getFollowerCount(ctx, userID)
    if followerCount >= bigVThreshold {
        // 大 V: 只写 outbox，不推送
        rdb.ZAdd(ctx, fmt.Sprintf("outbox:%d", userID),
            redis.Z{Score: float64(post.CreatedAt.Unix()), Member: post.ID})
        return nil
    }

    // ③ 普通用户: 异步推送到所有粉丝的 inbox
    msg := FanoutMsg{UserID: userID, PostID: post.ID, Timestamp: post.CreatedAt.Unix()}
    return producer.Send("feed_fanout", msg)
}

// MQ 消费者: 异步扇出推送
func consumeFanout(msg FanoutMsg) error {
    // 分批获取粉丝列表
    var cursor uint64
    for {
        followers, nextCursor, _ := rdb.SScan(ctx,
            fmt.Sprintf("followers:%d", msg.UserID), cursor, "", 500).Result()

        pipe := rdb.Pipeline()
        for _, followerID := range followers {
            key := fmt.Sprintf("inbox:%s", followerID)
            pipe.ZAdd(ctx, key, redis.Z{
                Score:  float64(msg.Timestamp),
                Member: msg.PostID,
            })
            // 限制 inbox 大小: 只保留最新 1000 条
            pipe.ZRemRangeByRank(ctx, key, 0, -1001)
        }
        pipe.Exec(ctx)

        if nextCursor == 0 {
            break
        }
        cursor = nextCursor
    }
    return nil
}
```

### 读取 Feed 流程

```go
func GetFeed(ctx context.Context, userID int64, cursor int64, limit int64) ([]Post, int64, error) {
    // ① 从 inbox 读取（普通关注人的动态）
    maxScore := "+inf"
    if cursor > 0 {
        maxScore = fmt.Sprintf("%d", cursor-1) // 游标分页
    }
    inboxKey := fmt.Sprintf("inbox:%d", userID)
    inboxPosts, _ := rdb.ZRevRangeByScoreWithScores(ctx, inboxKey, &redis.ZRangeBy{
        Max: maxScore, Min: "-inf", Count: limit,
    }).Result()

    // ② 获取关注的大 V 列表
    bigVList := getBigVFollowing(ctx, userID) // 从关注列表中筛选大 V

    // ③ 拉取大 V 的 outbox
    var bigVPosts []redis.Z
    if len(bigVList) > 0 {
        pipe := rdb.Pipeline()
        cmds := make([]*redis.ZSliceCmd, len(bigVList))
        for i, vid := range bigVList {
            cmds[i] = pipe.ZRevRangeByScoreWithScores(ctx,
                fmt.Sprintf("outbox:%d", vid),
                &redis.ZRangeBy{Max: maxScore, Min: "-inf", Count: limit})
        }
        pipe.Exec(ctx)
        for _, cmd := range cmds {
            bigVPosts = append(bigVPosts, cmd.Val()...)
        }
    }

    // ④ 合并排序（inbox + 大 V outbox）
    allPosts := mergeAndSort(inboxPosts, bigVPosts, limit)

    // ⑤ 批量查询 post 详情
    postIDs := extractPostIDs(allPosts)
    posts, _ := batchGetPosts(ctx, postIDs)

    // ⑥ 返回结果 + 下一页游标
    var nextCursor int64
    if len(allPosts) > 0 {
        nextCursor = int64(allPosts[len(allPosts)-1].Score)
    }
    return posts, nextCursor, nil
}
```

---

### 关键优化

#### 1. 分页方案：游标分页（非 offset）

```
为什么不用 LIMIT offset, count？
  offset=10000 时 → 跳过 10000 条 → 性能差

游标分页: 用最后一条的时间戳做游标
  第一页: ZREVRANGEBYSCORE inbox +inf -inf LIMIT 0 20
  第二页: ZREVRANGEBYSCORE inbox (上一页最后时间戳) -inf LIMIT 0 20
  → 每页查询都是 O(logN + M)，不受偏移量影响
```

#### 2. 缓存热点 Feed

```
活跃用户的 Feed 缓存:
  首次加载: 合并生成 Feed → 缓存到 Redis（TTL 5min）
  后续刷新: 直接读缓存
  有新动态: 增量更新缓存（而非全量重建）
```

#### 3. 删除/取关处理

```
删除动态:
  推模式: 发 MQ → 异步从所有粉丝 inbox 中删除该 post_id
  拉模式: 直接删 outbox 中的记录 → 下次拉取自然没了

取关:
  推模式: 需要清理 inbox 中该用户的所有动态（或懒清理: 读取时过滤）
  拉模式: 不需要额外操作
  推拉结合: 大 V 不需要；普通用户懒清理（读取时检查关注关系，已取关的跳过）
```

#### 4. 新用户冷启动

```
新注册用户 inbox 为空 → Feed 是空的

冷启动策略:
  ① 默认关注一些官方/推荐账号 → inbox 有内容
  ② 首次加载用拉模式 → 从推荐账号的 outbox 拉取
  ③ 推荐系统补充热门内容
```

---

### 各产品的方案选择

```
微信朋友圈: 纯推模式
  好友上限 5000 → 写入量可控
  刷朋友圈频率高 → 读必须快 → inbox 直接读

微博: 推拉结合
  大 V 粉丝千万 → 不能推 → 拉
  普通用户粉丝少 → 推到 inbox
  刷首页: inbox + 拉大 V 合并

Twitter: 推拉结合（类似微博）
  + 智能排序（非纯时间线，加入推荐算法）

抖音关注页: 推拉结合
  + 已看去重 + 推荐穿插

GitHub: 纯拉模式
  关注数少 + 刷新频率低 → 每次拉取合并即可
```

### 一句话总结

**Feed 流的核心是推拉结合：普通用户发动态走推模式（异步写入粉丝 inbox，读取直接查 inbox 极快），大 V 发动态走拉模式（只写 outbox，粉丝读取时实时拉取合并）。** 推模式用 Redis ZSET 做 inbox（score=时间戳，ZREVRANGEBYSCORE 分页），拉模式查 outbox 后多路合并排序。分界线通常为 1 万~10 万粉丝。微信朋友圈（好友上限 5000）可以纯推，微博/Twitter（大 V 千万粉丝）必须推拉结合。关键优化：游标分页（非 offset）、inbox 大小限制（ZREMRANGEBYRANK 裁剪）、异步扇出（MQ 消费批量推送）、取关懒清理（读取时过滤）。

---

## 九、IM 即时通讯系统设计

### 核心需求

```
基础功能:
  ① 单聊: A → B 一对一消息
  ② 群聊: A → 群(100~500人) 消息
  ③ 消息类型: 文本/图片/语音/视频/文件
  ④ 已读回执: "对方已读"
  ⑤ 多端同步: 手机/电脑/平板同时在线，消息同步
  ⑥ 离线消息: 用户不在线时收到的消息，上线后拉取
  ⑦ 历史消息: 拉取聊天记录

关键指标:
  消息延迟: < 200ms（用户感知实时）
  消息不丢: 发出去的消息必须到达（可靠投递）
  消息有序: 同一会话内消息按发送顺序显示
  高并发:   百万在线连接，十万级 msg/s
```

### 整体架构

```
  ┌────────┐     ┌────────┐     ┌────────┐
  │ Client │     │ Client │     │ Client │
  │ (手机)  │     │ (电脑)  │     │ (Web)  │
  └───┬────┘     └───┬────┘     └───┬────┘
      │ WebSocket    │ WebSocket    │ WebSocket
      └──────────┬───┴──────────┬──┘
                 │              │
           ┌─────▼──────────────▼─────┐
           │     长连接网关层            │  持有 WebSocket 连接
           │  Gateway 1  Gateway 2 ... │  维护 用户→连接 映射
           └─────────┬────────────────┘
                     │ gRPC / 内部协议
           ┌─────────▼────────────────┐
           │       消息逻辑层           │  消息处理、路由、持久化
           │  MsgServer 1  MsgServer 2 │
           └────┬──────┬──────┬───────┘
                │      │      │
           ┌────▼──┐ ┌─▼────┐ ┌▼──────┐
           │ MySQL │ │Redis │ │ Kafka │
           │ 消息存储│ │路由表 │ │ 消息总线│
           └───────┘ └──────┘ └───────┘
```

---

### 一、消息发送流程（单聊）

```
A 发消息给 B:

  ① A 的客户端 → WebSocket → Gateway 1
  ② Gateway 1 → MsgServer (消息逻辑服务)
  ③ MsgServer:
     a. 生成全局唯一消息 ID（Snowflake）
     b. 持久化到消息存储（MySQL / HBase）
     c. 查 Redis 路由表: B 连在哪个 Gateway？
        → B 在线: 在 Gateway 2 上
        → B 不在线: 存离线消息
     d. 投递: 通过 Kafka/gRPC → Gateway 2 → B 的 WebSocket
  ④ B 收到消息 → 发送 ACK → 更新已送达状态

  A (Gateway1)          MsgServer          B (Gateway2)
    │                      │                    │
    │── 发送消息 ──────────→│                    │
    │                      │ 生成ID+持久化       │
    │                      │ 查路由: B在Gateway2 │
    │                      │─── 投递消息 ───────→│
    │                      │                    │── 推送给B
    │                      │←── ACK ────────────│
    │←── 服务端ACK ────────│                    │
```

### 二、消息存储设计

#### 读扩散 vs 写扩散

```
IM 消息存储有两种经典模型:

  写扩散（微信方案）:
    每个参与者都存一份消息
    A→B: 写入 A 的消息表 + 写入 B 的消息表
    优点: 读取直接查自己的表，极快；删除/撤回不影响对方
    缺点: 存储翻倍；群聊 500 人 = 写 500 份

  读扩散（钉钉/Discord 方案）:
    消息只存一份（以会话为维度）
    A→B: 写入 conversation_123 的消息表
    优点: 存储省（只一份）；群聊写入 O(1)
    缺点: 读取需按会话查；删除需要"对我不可见"标记

  大多数系统采用: 消息表（读扩散） + 会话信箱（写扩散辅助）
```

#### 数据库表设计

```sql
-- 消息表（一条消息一条记录，读扩散，以会话为维度）
CREATE TABLE message (
    msg_id         BIGINT PRIMARY KEY,        -- Snowflake ID
    conversation_id VARCHAR(64) NOT NULL,      -- 会话 ID
    sender_id      BIGINT NOT NULL,
    msg_type       TINYINT DEFAULT 1,          -- 1=文本 2=图片 3=语音 ...
    content        TEXT,
    send_time      BIGINT NOT NULL,            -- 毫秒时间戳
    INDEX idx_conv_time (conversation_id, send_time DESC)
) ENGINE=InnoDB;

-- 会话信箱（写扩散，每个用户的会话列表 + 最新消息位置）
CREATE TABLE user_conversation (
    user_id         BIGINT NOT NULL,
    conversation_id VARCHAR(64) NOT NULL,
    last_msg_id     BIGINT DEFAULT 0,          -- 该会话最新消息 ID
    last_read_id    BIGINT DEFAULT 0,          -- 已读到的消息 ID
    unread_count    INT DEFAULT 0,             -- 未读数
    updated_at      BIGINT NOT NULL,           -- 最后更新时间（排序用）
    PRIMARY KEY (user_id, conversation_id),
    INDEX idx_user_update (user_id, updated_at DESC)
);

-- conversation_id 生成规则:
--   单聊: "single:{min(A,B)}:{max(A,B)}"  → 保证 A→B 和 B→A 是同一个会话
--   群聊: "group:{group_id}"
```

#### Redis 缓存

```
用户在线状态 + 路由:
  online:{user_id} → "{gateway_id}:{conn_id}"   (用户连在哪个网关)
  TTL = 连接心跳超时（如 60s，每次心跳续期）

最近会话列表:
  recent_conv:{user_id} → ZSET [conv_id(score=update_time)]
  → 打开 APP 首先展示最近会话列表

未读消息数:
  unread:{user_id}:{conv_id} → INT（原子 INCR/重置）
  unread_total:{user_id} → INT（APP 角标数字）
```

---

### 三、消息可靠投递

```
消息不能丢，但网络不可靠 → 需要确认机制

三次握手确认（类似 TCP）:

  A                    Server                   B
  │── ① 发送消息 ──────→│                        │
  │                     │── ② 投递消息 ──────────→│
  │                     │←── ③ B 回 ACK ─────────│
  │←── ④ 服务端ACK ─────│                        │
  │   (告诉A: 消息已送达)│                        │

  如果 B 不在线:
    ② 消息存入离线消息表 → B 上线后拉取
    ④ 服务端告诉 A: "消息已存储，对方离线"

  如果 ③ ACK 丢失（B 收到了但 ACK 网络丢了）:
    Server 超时重传 → B 收到重复消息 → B 用 msg_id 去重

  如果 ④ ACK 丢失:
    A 超时重传 → Server 用 msg_id 去重 → 幂等
```

```go
// 消息确认状态机
const (
    MsgStatusSending   = 0 // 发送中（客户端本地）
    MsgStatusServerACK = 1 // 服务端已收到
    MsgStatusDelivered = 2 // 已送达对方设备
    MsgStatusRead      = 3 // 已读
)
```

---

### 四、已读回执

```
单聊已读:
  B 打开和 A 的聊天窗口 → 客户端上报已读位置:
    "我已读到 conversation_123 的 msg_id = 5000"
  → Server 更新 user_conversation.last_read_id = 5000
  → 推送给 A: "B 已读到 msg_5000"
  → A 的界面: msg_5000 及之前的消息显示"已读"

  注意: 上报的是"最后已读消息 ID"，不是每条单独标记
  → 之前所有消息自动标记为已读（比逐条标记高效得多）

群聊已读:
  方案 A: 存每个成员的已读位置
    group_read:{group_id}:{user_id} = last_read_msg_id
    → 查某条消息被多少人已读:
      遍历群成员，统计 last_read_id >= msg_id 的人数
    → 群人数多时开销大

  方案 B: 只显示已读人数（不显示具体谁读了）
    每条消息维护一个已读计数器
    read_count:{msg_id} → INT
    → 成员已读时: INCR read_count:{msg_id}
    → 展示: "10人已读"（不需要知道具体谁）
```

---

### 五、多端同步

```
同一用户在手机、电脑、平板同时在线:
  → 三个设备三条 WebSocket 连接（可能在不同 Gateway）

  Redis 路由表:
    online:{user_id}:phone  → gateway_1:conn_001
    online:{user_id}:pc     → gateway_2:conn_002
    online:{user_id}:pad    → gateway_3:conn_003

  消息投递时:
    查路由表 → 用户有 3 个在线设备 → 向 3 个 Gateway 都投递
    → 3 个设备同时收到消息

  消息同步（设备新上线/切换时）:
    每个设备记录自己同步到的位置: device_sync_id
    新设备上线 → 拉取 device_sync_id 之后的所有消息 → 追上最新状态

    手机 sync_id = 4800
    新开电脑 → 拉取 msg_4801 ~ msg_5000 → 同步完成
```

---

### 六、群聊消息

```
群聊的核心挑战: 一条消息 → 推送给 N 个群成员（N 可能是 500）

方案: 写扩散到 inbox（类似 Feed 流推模式）

  A 在群里发消息:
  ① 消息写入 message 表（一份）
  ② 查群成员列表（500 人）
  ③ 更新每个成员的 user_conversation（last_msg_id, unread_count++）
  ④ 在线成员: 通过各自的 Gateway 实时推送
     离线成员: 推送消息通知（APNs/FCM）

  优化:
  - ③④ 异步执行（发 Kafka → 消费者批量处理）
  - 群成员列表缓存在 Redis SET 中
  - 在线成员批量推送（同一 Gateway 上的合并为一次 RPC）
```

```
群消息扩散的写入量:
  500 人群 × 每天 200 条消息 = 10 万次 user_conversation 更新/天
  → 用 Kafka 异步 + 批量更新（不在消息发送的主链路上）

  万人大群（如 Discord）:
  → 不再写扩散（太慢）
  → 纯读扩散: 成员打开群时实时拉取最新消息
  → 未读数改为"有新消息"标记（不精确计数）
```

---

### 七、离线消息

```
用户 B 不在线时 A 发的消息:

方案 1: 离线消息表
  B 不在线 → 消息存入 offline_msg 表
  B 上线  → 拉取 offline_msg → 删除
  问题: 大群+长时间离线 → 离线消息量巨大

方案 2: 基于同步位的拉取（推荐）
  每个用户+每个会话记录 last_read_id
  B 上线 → 遍历 user_conversation → 找到 last_read_id < last_msg_id 的会话
  → 拉取差量消息（从 last_read_id 到 last_msg_id）
  → 不需要额外的离线消息表

  优势:
    多端同步、离线消息、历史消息用同一套机制
    只需要 message 表 + user_conversation 表
```

---

### 八、消息有序性

```
问题: 网络延迟不同，先发的消息可能后到

  A 发 msg_1 (10:00:00.000)
  A 发 msg_2 (10:00:00.100)
  → msg_2 先到服务端 → msg_1 后到
  → B 看到的顺序是 msg_2, msg_1 → 乱序！

解决:
  ① 服务端分配序号: MsgServer 按会话维度分配递增 seq
     msg_1 → seq=101
     msg_2 → seq=102（即使 msg_2 先到，seq 也按到达顺序分配）
  ② 客户端按 seq 排序显示（不是按客户端发送时间）

  会话级别的 seq（不是全局 seq）:
    conversation_123 的 seq: 101, 102, 103, ...
    conversation_456 的 seq: 501, 502, 503, ...
    → 每个会话独立递增 → 并发量分散 → 用 Redis INCR 实现

    seq = INCR conv_seq:{conversation_id}
```

---

### 关键数字参考

```
微信级别（10 亿用户）:
  日消息量: 百亿条
  峰值 QPS: 百万级
  长连接数: 数亿

面试级别（中等规模）:
  用户数:    100 万
  在线连接:  10 万
  消息 QPS:  5 万
  单 Gateway: 5 万连接（Go 轻松）
  Gateway 数: 2~4 台
  MsgServer:  2~4 台

消息延迟目标:
  发送→对方收到: < 200ms（同机房）
  发送→持久化:   < 50ms
  离线→上线同步: < 2s（增量拉取）
```

### 一句话总结

**IM 系统 = 长连接网关层（持有 WebSocket，维护用户路由表）+ 消息逻辑层（持久化+路由+投递）+ 存储层（MySQL 消息表 + Redis 路由/缓存 + Kafka 异步扩散）。** 消息发送：客户端→Gateway→MsgServer（生成 ID+持久化+查路由+投递）→目标 Gateway→接收方。消息可靠：三次握手确认（发送→服务端 ACK→接收方 ACK）+ msg_id 去重。多端同步：每设备独立 sync_id，上线拉差量。已读回执：上报"最后已读 msg_id"而非逐条标记。群聊：小群写扩散（更新每个成员的 inbox）+异步 Kafka 批量推送；万人大群改为读扩散（打开时实时拉）。消息有序：服务端按会话维度分配递增 seq（Redis INCR），客户端按 seq 排序。

---

## 十、微服务治理

### 微服务治理全景图

```
  ┌──────────────────────────────────────────────────────────────┐
  │                     微服务治理体系                             │
  │                                                              │
  │  流量入口        服务通信          稳定性           可观测性     │
  │  ├─ API 网关    ├─ 服务注册发现    ├─ 熔断          ├─ 链路追踪 │
  │  ├─ 负载均衡    ├─ RPC 框架       ├─ 降级          ├─ 指标监控 │
  │  └─ 限流        └─ 服务网格       ├─ 重试          └─ 日志聚合 │
  │                                  └─ 超时控制                 │
  └──────────────────────────────────────────────────────────────┘
```

---

### 一、服务注册与发现

#### 为什么需要

```
单体: 所有功能在一个进程 → 直接函数调用
微服务: 功能拆分到多个服务 → 网络调用
  → 问题: 订单服务要调用户服务，用户服务的 IP:Port 是什么？
  → 服务实例可能有 10 个，IP 随时变（扩缩容/重启/迁移）
  → 不能写死配置 → 需要动态的服务注册与发现
```

#### 工作原理

```
  ┌──────────┐  ① 注册  ┌──────────────┐
  │ 用户服务  │─────────→│  注册中心      │
  │ 10.0.0.1 │  心跳续期 │ (etcd/Consul) │
  │ 10.0.0.2 │─────────→│              │
  │ 10.0.0.3 │─────────→│  用户服务:     │
  └──────────┘          │  [10.0.0.1,   │
                        │   10.0.0.2,   │
  ┌──────────┐  ② 发现  │   10.0.0.3]  │
  │ 订单服务  │←─────────│              │
  │          │  ③ 调用   └──────────────┘
  │          │──────────→ 10.0.0.2 (负载均衡选一个)
  └──────────┘

  ① 注册: 服务启动时向注册中心注册自己的地址 + 定期心跳
  ② 发现: 调用方从注册中心获取目标服务的实例列表
  ③ 调用: 负载均衡选择一个实例发起 RPC/HTTP 调用
  ④ 下线: 服务停止 → 主动注销或心跳超时 → 注册中心摘除
```

#### 注册中心对比

| 组件 | 一致性协议 | 健康检查 | 适用场景 | Go 生态 |
|------|-----------|---------|---------|---------|
| **etcd** | **Raft（强一致）** | Lease 租约 | K8s 生态、Go 微服务 | go-micro, Kratos |
| **Consul** | Raft + Gossip | HTTP/gRPC/TCP 多种 | 多语言、服务网格 | go-micro, Consul SDK |
| **ZooKeeper** | ZAB（强一致） | 临时节点 | Java 生态（Dubbo） | 较少直接用 |
| **Nacos** | Raft + Distro | 心跳 | 阿里生态、Java 为主 | Go SDK 有 |
| **K8s Service** | etcd（内置） | Pod 探针 | 容器化环境 | 原生支持 |

```
Go 微服务常见选择:
  K8s 环境: 直接用 K8s Service + CoreDNS（零额外组件）
  非 K8s:   etcd 或 Consul
```

---

### 二、熔断（Circuit Breaker）

#### 为什么需要

```
服务链路: A → B → C → D

  D 出问题（响应慢/超时）:
  → C 调 D 超时 → C 的线程/goroutine 全部阻塞等待
  → B 调 C 也超时 → B 也阻塞
  → A 调 B 也超时 → A 也阻塞
  → 整条链路雪崩！一个下游拖垮所有上游

  熔断: 像电路的保险丝，检测到下游故障后自动"断开"
  → 不再调用故障服务 → 快速失败/降级返回
  → 保护自己不被拖垮
```

#### 熔断器状态机

```
  ┌─────────┐  失败率 > 阈值   ┌─────────┐
  │  Closed  │────────────────→│  Open    │
  │  (正常)  │                 │ (熔断中)  │
  └────┬────┘                 └────┬────┘
       ↑                           │
  成功率恢复                  超时后放几个试探请求
       │                           │
  ┌────┴────────┐                  │
  │ Half-Open   │←─────────────────┘
  │ (半开/试探)  │
  └─────────────┘
  试探成功 → Closed（恢复）
  试探失败 → Open（继续熔断）

Closed（关闭 = 正常通行）:
  所有请求正常通过
  统计失败率 → 超过阈值（如 50%）→ 切换到 Open

Open（打开 = 熔断中）:
  所有请求直接失败（不调用下游）→ 快速返回错误或降级数据
  等待超时时间（如 10s）→ 切换到 Half-Open

Half-Open（半开 = 试探中）:
  放行少量请求（如 1~3 个）试探下游是否恢复
  试探成功 → 切换到 Closed（恢复正常）
  试探失败 → 切换到 Open（继续熔断）
```

#### Go 熔断实现（使用 sony/gobreaker）

```go
import "github.com/sony/gobreaker"

var cb *gobreaker.CircuitBreaker

func init() {
    cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "user-service",
        MaxRequests: 3,                    // Half-Open 时允许 3 个试探请求
        Interval:    10 * time.Second,     // Closed 状态下统计窗口
        Timeout:     30 * time.Second,     // Open → Half-Open 的等待时间
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            // 触发熔断的条件: 失败率 > 60% 且请求数 > 10
            return counts.Requests > 10 &&
                float64(counts.TotalFailures)/float64(counts.Requests) > 0.6
        },
        OnStateChange: func(name string, from, to gobreaker.State) {
            log.Printf("breaker %s: %s → %s", name, from, to)
        },
    })
}

func CallUserService(userID int64) (*User, error) {
    result, err := cb.Execute(func() (interface{}, error) {
        return httpClient.Get(fmt.Sprintf("http://user-service/users/%d", userID))
    })
    if err != nil {
        // 熔断中 → 返回降级数据
        if errors.Is(err, gobreaker.ErrOpenState) {
            return getDefaultUser(userID), nil // 降级: 返回缓存/默认值
        }
        return nil, err
    }
    return result.(*User), nil
}
```

---

### 三、降级（Degradation）

```
降级 = 系统压力大或部分功能不可用时，关闭非核心功能，保证核心功能可用

  电商下单链路:
    核心: 创建订单 + 扣库存 + 扣款
    非核心: 发优惠券 + 加积分 + 发通知 + 更新推荐

  降级策略:
    ① 超时降级: 调用超时 → 返回默认值
    ② 熔断降级: 下游熔断 → 返回缓存数据
    ③ 限流降级: 超出限流 → 返回"系统繁忙"
    ④ 手动降级: 运营后台一键关闭非核心功能（如大促前关闭推荐）
```

```go
// 降级开关（通过配置中心动态控制）
func CreateOrder(ctx context.Context, req OrderReq) (*Order, error) {
    // 核心逻辑（不降级）
    order, err := createOrderCore(ctx, req)
    if err != nil {
        return nil, err
    }

    // 非核心逻辑（可降级）
    if !config.GetBool("degrade.coupon") {
        go sendCoupon(ctx, order) // 异步 + 降级开关
    }
    if !config.GetBool("degrade.points") {
        go addPoints(ctx, order)
    }
    if !config.GetBool("degrade.notification") {
        go sendNotification(ctx, order)
    }

    return order, nil
}
```

```
降级 vs 熔断的区别:
  熔断: 自动的，检测到下游故障后自动触发
  降级: 可以是自动的（超时/限流触发）也可以是手动的（运营开关）
  熔断是降级的一种触发方式
```

---

### 四、重试与超时

```
重试:
  调用失败 → 自动重试（可能是网络抖动，重试就好了）

  重试策略:
    固定间隔:    1s → 1s → 1s
    指数退避:    1s → 2s → 4s → 8s（推荐）
    指数退避+抖动: 1s±random → 2s±random → 4s±random（防止重试风暴）

  重试注意:
    只有幂等操作才能重试（GET 可以，POST 扣款不行）
    设置最大重试次数（3~5 次）
    重试风暴: 下游故障 → 所有调用方都在重试 → 下游更扛不住
    → 配合熔断: 熔断后不再重试

超时:
  每层设置合理超时（参考 go-qa.md 第九题 Context 超时传播）
  共享 deadline 而非各层独立计时
  API 总超时 5s → 下游 RPC 超时 3s → DB 查询超时 1s
```

```go
// Go 重试 + 指数退避
func retryWithBackoff(ctx context.Context, maxRetries int, fn func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        // 指数退避 + 随机抖动
        backoff := time.Duration(1<<uint(i)) * 100 * time.Millisecond
        jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff + jitter):
        }
    }
    return err
}
```

---

### 五、链路追踪（Distributed Tracing）

#### 为什么需要

```
用户请求 → API Gateway → 订单服务 → 用户服务 → Redis
                                  → 库存服务 → MySQL
                                  → 支付服务 → 第三方 API

  请求慢了，是哪个环节慢？
  请求失败了，是哪个服务出错？
  → 需要全链路追踪，把一次请求经过的所有服务串起来
```

#### 核心概念

```
Trace:  一次完整请求的调用链（从入口到结束）
Span:   链路中的一个操作（一次 RPC/DB 查询/HTTP 调用）
TraceID: 全局唯一，标识一次请求
SpanID:  标识一个 Span
ParentSpanID: 父 Span（形成调用树）

  Trace: abc-123
  ├── Span 1: API Gateway (parent=null)
  │   ├── Span 2: 订单服务 (parent=1)
  │   │   ├── Span 3: 用户服务 (parent=2)
  │   │   ├── Span 4: 库存服务 (parent=2)
  │   │   │   └── Span 5: MySQL查询 (parent=4)
  │   │   └── Span 6: 支付服务 (parent=2)
  │   │       └── Span 7: 第三方API (parent=6)

  → 一棵调用树，每个节点记录: 开始时间、耗时、状态、标签
  → 可以看到: 哪个 Span 最慢、哪个 Span 出错
```

#### TraceID 传播

```
TraceID 必须在整个调用链中传递:

  HTTP: 通过 Header 传递
    X-Trace-Id: abc-123
    X-Span-Id: span-001

  gRPC: 通过 Metadata 传递
    md := metadata.Pairs("trace-id", "abc-123")
    ctx = metadata.NewOutgoingContext(ctx, md)

  MQ: 通过消息 Header 传递
    msg.Header["trace-id"] = "abc-123"

  Go 中: 存在 context.Context 中传递
    ctx = context.WithValue(ctx, "traceID", "abc-123")
```

#### 常见链路追踪系统

| 系统 | 开发者 | 特点 | Go 集成 |
|------|--------|------|---------|
| **Jaeger** | Uber | CNCF 毕业，K8s 原生 | OpenTelemetry SDK |
| **Zipkin** | Twitter | 最早的开源方案 | zipkin-go |
| SkyWalking | Apache | Java 生态强，自动探针 | Go 手动埋点 |
| **OpenTelemetry** | CNCF | **统一标准**（Trace+Metric+Log） | otel SDK |

```
推荐: OpenTelemetry（OTel）— 统一的可观测性标准

  OpenTelemetry = OpenTracing + OpenCensus 合并
  → 一套 SDK 同时采集 Trace + Metric + Log
  → 数据导出到 Jaeger/Zipkin/Prometheus/任何后端
  → CNCF 最活跃的项目之一，已成为事实标准
```

```go
// Go OpenTelemetry 简单示例
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func CreateOrder(ctx context.Context, req OrderReq) error {
    // 创建 Span
    ctx, span := otel.Tracer("order-service").Start(ctx, "CreateOrder")
    defer span.End()

    // 添加标签
    span.SetAttributes(attribute.Int64("order.user_id", req.UserID))

    // 调用下游（ctx 自动携带 TraceID）
    user, err := userClient.GetUser(ctx, req.UserID)
    if err != nil {
        span.RecordError(err) // 记录错误
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    // ... 后续逻辑
    return nil
}
```

---

### 六、负载均衡

```
服务有多个实例，选哪个？

  常见策略:
    轮询 (Round Robin):      1→2→3→1→2→3（均匀分配）
    加权轮询:                 权重高的多分配
    随机:                     随机选一个
    最少连接:                 选当前连接数最少的实例
    一致性哈希:               同一 key 总是路由到同一实例（缓存友好）
    P2C (Power of Two Choices): 随机选两个，取负载低的（gRPC 默认）

  Go 微服务常用:
    gRPC: 内置 round_robin / pick_first，可自定义
    Kratos: 支持 WRR (Weighted Round Robin) + P2C
    go-zero: 自适应负载均衡
```

---

### 七、API 网关

```
API 网关 = 微服务的统一入口，聚合多种治理能力

  客户端 → API 网关 → 各微服务

  网关职责:
    路由:      /api/users → 用户服务, /api/orders → 订单服务
    限流:      令牌桶/滑动窗口
    认证鉴权:   JWT 验证, OAuth2
    协议转换:   HTTP → gRPC
    灰度发布:   1% 流量到新版本
    监控:      请求量/延迟/错误率

  常见网关:
    Kong:     基于 Nginx + Lua，插件丰富
    APISIX:   基于 Nginx + Lua，性能好（国产）
    Envoy:    C++ 实现，Service Mesh 标配
    Traefik:  Go 实现，K8s 原生
    自研:     Go 写自定义网关（大厂常见）
```

---

### Go 微服务框架对比

```
Go 微服务框架（内置治理能力）:

  ┌────────────┬───────────────────────────────────────────────┐
  │ 框架        │ 内置治理能力                                   │
  ├────────────┼───────────────────────────────────────────────┤
  │ go-zero    │ 服务发现(etcd) + 限流 + 熔断 + 超时 + 链路追踪  │
  │            │ 自适应负载均衡 + 自适应熔断                      │
  ├────────────┼───────────────────────────────────────────────┤
  │ Kratos     │ 服务发现(多种) + 限流 + 熔断 + OTel + 中间件链  │
  │ (B站开源)  │ gRPC 原生 + HTTP 转换                          │
  ├────────────┼───────────────────────────────────────────────┤
  │ go-micro   │ 服务发现 + 负载均衡 + 消息驱动 + 插件化          │
  ├────────────┼───────────────────────────────────────────────┤
  │ Kitex      │ 服务发现 + 熔断 + 重试 + 负载均衡 + 链路追踪    │
  │ (字节开源)  │ 高性能 Thrift/gRPC                            │
  └────────────┴───────────────────────────────────────────────┘

  不用框架也行: 用 gRPC + 各组件单独集成
    服务发现: etcd + go SDK
    熔断: sony/gobreaker
    限流: golang.org/x/time/rate 或 Redis 分布式限流
    链路追踪: OpenTelemetry SDK
    → 组合式治理，灵活但集成工作量大
```

---

### 微服务治理面试常见追问

```
Q: 服务注册发现用 AP 还是 CP？
A: 注册中心推荐 CP（强一致，etcd/ZK/Consul）
   → 宁可暂时不可用，也不能返回错误的服务列表
   → 配合客户端缓存: 注册中心不可用时用本地缓存兜底

Q: 熔断阈值怎么设？
A: 没有万能值，需要根据业务调:
   错误率阈值: 50%~80%（太低容易误熔断）
   统计窗口:   10s~30s
   恢复超时:   10s~60s
   试探请求数: 3~5 个

Q: 微服务拆分原则？
A: - 单一职责（一个服务做一件事）
   - 按业务领域拆（DDD 的 Bounded Context）
   - 独立部署、独立数据库
   - 不要拆太细（2 个人维护一个服务比较合理）
   - 先单体 → 业务稳定后再拆 → 避免过早微服务化
```

### 一句话总结

**微服务治理六大核心：服务注册发现（etcd/Consul，启动注册+心跳续期+调用方拉取实例列表）、熔断（Closed→Open→Half-Open 状态机，失败率超阈值自动断路+快速失败）、降级（关闭非核心功能保核心，配置中心动态开关）、重试（指数退避+抖动，仅幂等操作+配合熔断防风暴）、链路追踪（OpenTelemetry 标准，TraceID 全链路传递，Span 记录每段耗时）、负载均衡（P2C/加权轮询/一致性哈希）。** Go 微服务推荐组合：K8s Service 或 etcd 服务发现 + gobreaker 熔断 + OpenTelemetry 链路追踪 + 分布式限流。框架推荐 go-zero / Kratos（内置全套治理能力）。

---

## 十一、单点登录（SSO）系统设计

> 多系统统一登录怎么实现？JWT vs Session？OAuth2？

### 1. 什么是 SSO

```
SSO = Single Sign-On（单点登录）

核心诉求:
  用户在 A 系统登录后，访问 B、C、D 系统时 无需再次登录

典型场景:
  - 企业内部: OA、邮箱、CRM、HR 系统 → 一套账号登录
  - 互联网: 淘宝登录后 → 天猫、支付宝免登
  - SaaS 平台: Google 登录后 → Gmail、Drive、Calendar 免登
```

### 2. SSO 核心架构

```
┌────────────────────────────────────────────────────────┐
│                    SSO 认证中心                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │
│  │ 登录页面  │  │ 用户存储  │  │ Token/Ticket 管理    │  │
│  └──────────┘  └──────────┘  └──────────────────────┘  │
└──────────┬──────────┬──────────────┬───────────────────┘
           │          │              │
    ┌──────▼──┐  ┌────▼────┐  ┌─────▼─────┐
    │ App A   │  │ App B   │  │ App C     │
    │(子系统) │  │(子系统) │  │(子系统)   │
    └─────────┘  └─────────┘  └───────────┘

关键角色:
  1. SSO Server（认证中心）: 统一登录入口，签发凭证
  2. SSO Client（子系统）:    拦截请求，校验凭证，获取用户信息
  3. Browser（用户浏览器）:   携带 Cookie / Token 跳转
```

### 3. 三种主流方案对比

| 维度 | Session + Cookie（同域） | CAS（跨域 SSO） | OAuth2 / OIDC |
|------|-------------------------|-----------------|---------------|
| **原理** | 共享 Session 存储 | 中心签发 Ticket，子系统验证 | 授权码 → Access Token → 资源 |
| **适用场景** | 同一主域名下的子系统 | 企业内部跨域多系统 | 开放平台、第三方登录 |
| **跨域支持** | 仅同域（.example.com） | 支持完全跨域 | 支持完全跨域 |
| **复杂度** | 低 | 中 | 高 |
| **标准化** | 无标准 | CAS 协议 | RFC 6749 / OIDC |
| **典型实现** | Spring Session + Redis | Apereo CAS | Keycloak、Auth0、Hydra |
| **移动端友好** | 差（依赖 Cookie） | 一般 | 好（Bearer Token） |

---

### 4. 方案一：Session + Cookie（同域 SSO）

```
适用: 所有子系统在同一主域名下
例如: login.example.com / app1.example.com / app2.example.com

原理:
  Cookie 的 domain 设为 .example.com → 所有子域名共享
  Session 存到 Redis → 所有后端服务共享

┌──────────────────────────────────────────────────┐
│                  用户浏览器                        │
│  Cookie: session_id=abc123 (domain=.example.com) │
└────────┬────────────────┬────────────────┬───────┘
         │                │                │
  ┌──────▼──────┐  ┌──────▼──────┐  ┌─────▼───────┐
  │ App A       │  │ App B       │  │ App C       │
  │ :8001       │  │ :8002       │  │ :8003       │
  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
         │                │                │
         └────────────────┼────────────────┘
                   ┌──────▼──────┐
                   │   Redis     │
                   │ session存储  │
                   └─────────────┘

流程:
  1. 用户访问 App A → 无 session_id → 跳转 login.example.com
  2. 用户登录成功 → 生成 session → 存 Redis → Set-Cookie(domain=.example.com)
  3. 用户访问 App B → 自动携带 Cookie → 后端从 Redis 取 session → 已登录
```

#### Go 实现（同域 Session SSO）

```go
package sso

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionSSO 同域 SSO 实现
type SessionSSO struct {
	rdb        *redis.Client
	domain     string        // 如 ".example.com"
	sessionTTL time.Duration // Session 过期时间
}

type UserSession struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	LoginAt  int64  `json:"login_at"`
}

func NewSessionSSO(rdb *redis.Client, domain string) *SessionSSO {
	return &SessionSSO{
		rdb:        rdb,
		domain:     domain,
		sessionTTL: 24 * time.Hour,
	}
}

// Login 登录成功后调用
func (s *SessionSSO) Login(w http.ResponseWriter, userID int64, username string) error {
	// 1. 生成 session ID
	sessionID := generateSessionID()

	// 2. 存入 Redis
	sess := &UserSession{
		UserID:   userID,
		Username: username,
		LoginAt:  time.Now().Unix(),
	}
	data, _ := json.Marshal(sess)
	ctx := context.Background()
	if err := s.rdb.Set(ctx, "sso:session:"+sessionID, data, s.sessionTTL).Err(); err != nil {
		return err
	}

	// 3. 设置 Cookie（domain 设为主域名，所有子域名共享）
	http.SetCookie(w, &http.Cookie{
		Name:     "sso_session",
		Value:    sessionID,
		Domain:   s.domain, // .example.com
		Path:     "/",
		MaxAge:   int(s.sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

// AuthMiddleware 认证中间件（每个子系统使用）
func (s *SessionSSO) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sso_session")
		if err != nil {
			http.Redirect(w, r, "https://login.example.com/login?redirect="+r.URL.String(), http.StatusFound)
			return
		}

		// 从 Redis 获取 Session
		ctx := r.Context()
		data, err := s.rdb.Get(ctx, "sso:session:"+cookie.Value).Bytes()
		if err != nil {
			http.Redirect(w, r, "https://login.example.com/login?redirect="+r.URL.String(), http.StatusFound)
			return
		}

		var sess UserSession
		if err := json.Unmarshal(data, &sess); err != nil {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}

		// 将用户信息注入上下文
		ctx = context.WithValue(ctx, "user", &sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logout 登出（所有子系统同时登出）
func (s *SessionSSO) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sso_session")
	if err == nil {
		// 删除 Redis 中的 Session
		s.rdb.Del(r.Context(), "sso:session:"+cookie.Value)
	}
	// 清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "sso_session",
		Value:  "",
		Domain: s.domain,
		Path:   "/",
		MaxAge: -1,
	})
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
```

---

### 5. 方案二：CAS 协议（跨域 SSO）

```
CAS = Central Authentication Service
适用: 跨域名的多系统（如 a.com / b.com / c.com）

核心思想:
  1. 所有子系统都信任同一个认证中心（SSO Server）
  2. SSO Server 签发一次性 Ticket
  3. 子系统拿 Ticket 去 SSO Server 换用户信息

┌────────┐                ┌──────────┐               ┌────────┐
│ 浏览器  │                │SSO Server│               │ App A  │
└───┬────┘                └────┬─────┘               └───┬────┘
    │  1.访问 App A             │                        │
    │───────────────────────────────────────────────────>│
    │                          │    2.未登录,重定向到SSO   │
    │<──────────────────────────────────────── 302 ──────│
    │  3.跳转SSO登录页          │                        │
    │─────────────────────────>│                        │
    │  4.用户输入账号密码        │                        │
    │─────────────────────────>│                        │
    │                          │ 5.登录成功               │
    │                          │   a.创建全局Session(TGT) │
    │                          │   b.生成一次性Ticket(ST)  │
    │  6.Set-Cookie(TGC)       │                        │
    │  + 302 重定向到 App A     │                        │
    │  (带上 ?ticket=ST-xxx)    │                        │
    │<─────────────────────────│                        │
    │  7.携带 ticket 访问 App A │                        │
    │───────────────────────────────────────────────────>│
    │                          │  8.App A 后端验证ticket  │
    │                          │<───────────────────────│
    │                          │  9.返回用户信息          │
    │                          │───────────────────────>│
    │  10.App A 创建局部Session │                        │
    │<──────────────────────── 登录成功 ─────────────────│
    │                          │                        │
    │  --- 访问 App B ---      │               ┌────────┐
    │                          │               │ App B  │
    │  11.访问 App B            │               └───┬────┘
    │──────────────────────────────────────────────>│
    │                          │  12.未登录,重定向SSO│
    │<──────────────────────────────────── 302 ────│
    │  13.跳转SSO(自动带TGC Cookie)                 │
    │─────────────────────────>│                   │
    │  14.SSO检测到TGC有效      │                   │
    │     直接签发新Ticket      │                   │
    │  15.302重定向App B+ticket │                   │
    │<─────────────────────────│                   │
    │  16.携带ticket访问App B   │                   │
    │──────────────────────────────────────────────>│
    │                          │  17.验证ticket     │
    │                          │<──────────────────│
    │                          │  18.返回用户信息    │
    │                          │──────────────────>│
    │  19.App B 登录成功（免登）  │                   │
    │<───────────────────────── 登录成功 ───────────│

关键概念:
  TGT (Ticket Granting Ticket): 全局 Session，存在 SSO Server
  TGC (Ticket Granting Cookie): 浏览器存的 Cookie，关联 TGT
  ST  (Service Ticket):        一次性票据，用于子系统验证
```

#### Go 实现（CAS 风格 SSO Server）

```go
package sso

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// SSOServer CAS 风格的 SSO 认证中心
type SSOServer struct {
	rdb *redis.Client
}

type TicketInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Service  string `json:"service"` // 请求方的回调地址
}

func NewSSOServer(rdb *redis.Client) *SSOServer {
	return &SSOServer{rdb: rdb}
}

// HandleLogin 处理登录请求
func (s *SSOServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service") // 子系统回调地址

	if r.Method == http.MethodGet {
		// 检查是否已有全局 Session（TGC Cookie）
		tgcCookie, err := r.Cookie("TGC")
		if err == nil {
			// 已登录 → 直接签发 ST
			userJSON, err := s.rdb.Get(r.Context(), "sso:tgt:"+tgcCookie.Value).Bytes()
			if err == nil {
				var user TicketInfo
				json.Unmarshal(userJSON, &user)
				st := s.issueServiceTicket(r.Context(), user.UserID, user.Username, service)
				http.Redirect(w, r, service+"?ticket="+st, http.StatusFound)
				return
			}
		}
		// 未登录 → 展示登录页面
		fmt.Fprintf(w, `<form method="POST"><input name="username"/><input name="password" type="password"/><input type="hidden" name="service" value="%s"/><button>登录</button></form>`, service)
		return
	}

	// POST 登录验证
	username := r.FormValue("username")
	password := r.FormValue("password")
	service = r.FormValue("service")

	userID, ok := authenticate(username, password)
	if !ok {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 1. 创建全局 Session (TGT)
	tgtID := generateToken("TGT")
	tgtData, _ := json.Marshal(&TicketInfo{UserID: userID, Username: username})
	s.rdb.Set(r.Context(), "sso:tgt:"+tgtID, tgtData, 8*time.Hour)

	// 2. 设置 TGC Cookie（仅 SSO Server 域名）
	http.SetCookie(w, &http.Cookie{
		Name:     "TGC",
		Value:    tgtID,
		Path:     "/",
		MaxAge:   int((8 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// 3. 签发一次性 Service Ticket
	st := s.issueServiceTicket(r.Context(), userID, username, service)

	// 4. 重定向回子系统
	http.Redirect(w, r, service+"?ticket="+st, http.StatusFound)
}

// issueServiceTicket 签发一次性 ST
func (s *SSOServer) issueServiceTicket(ctx context.Context, userID int64, username, service string) string {
	st := generateToken("ST")
	data, _ := json.Marshal(&TicketInfo{
		UserID:   userID,
		Username: username,
		Service:  service,
	})
	// ST 只有 5 分钟有效期，且用一次就删
	s.rdb.Set(ctx, "sso:st:"+st, data, 5*time.Minute)
	return st
}

// HandleValidate 子系统调用，验证 ST 并获取用户信息
func (s *SSOServer) HandleValidate(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	service := r.URL.Query().Get("service")

	ctx := r.Context()
	data, err := s.rdb.Get(ctx, "sso:st:"+ticket).Bytes()
	if err != nil {
		http.Error(w, "ticket 无效或已过期", http.StatusUnauthorized)
		return
	}

	// 一次性使用 → 立即删除
	s.rdb.Del(ctx, "sso:st:"+ticket)

	var info TicketInfo
	json.Unmarshal(data, &info)

	// 校验 service 是否匹配（防止 Ticket 被窃取后用于其他系统）
	if info.Service != service {
		http.Error(w, "service 不匹配", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  info.UserID,
		"username": info.Username,
	})
}

// HandleLogout 全局登出
func (s *SSOServer) HandleLogout(w http.ResponseWriter, r *http.Request) {
	tgcCookie, err := r.Cookie("TGC")
	if err == nil {
		s.rdb.Del(r.Context(), "sso:tgt:"+tgcCookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "TGC",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	fmt.Fprint(w, "已登出")
}

// --- SSO Client（子系统侧） ---

// SSOClient 子系统的 SSO 客户端
type SSOClient struct {
	ssoServerURL string // SSO Server 地址
	serviceURL   string // 本系统的地址
	rdb          *redis.Client
}

func NewSSOClient(ssoServerURL, serviceURL string, rdb *redis.Client) *SSOClient {
	return &SSOClient{
		ssoServerURL: ssoServerURL,
		serviceURL:   serviceURL,
		rdb:          rdb,
	}
}

// AuthMiddleware 子系统认证中间件
func (c *SSOClient) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 检查是否携带 ticket（从 SSO 回调）
		ticket := r.URL.Query().Get("ticket")
		if ticket != "" {
			// 去 SSO Server 验证 ticket
			user, err := c.validateTicket(r.Context(), ticket)
			if err != nil {
				http.Error(w, "ticket 验证失败", http.StatusUnauthorized)
				return
			}
			// 创建局部 Session
			localSID := generateToken("LOCAL")
			data, _ := json.Marshal(user)
			c.rdb.Set(r.Context(), "app:session:"+localSID, data, 2*time.Hour)
			http.SetCookie(w, &http.Cookie{
				Name:     "app_session",
				Value:    localSID,
				Path:     "/",
				MaxAge:   7200,
				HttpOnly: true,
			})
			// 去掉 URL 中的 ticket 参数，避免刷新时重复验证
			cleanURL := r.URL.Path
			http.Redirect(w, r, cleanURL, http.StatusFound)
			return
		}

		// 2. 检查局部 Session
		cookie, err := r.Cookie("app_session")
		if err == nil {
			data, err := c.rdb.Get(r.Context(), "app:session:"+cookie.Value).Bytes()
			if err == nil {
				var user TicketInfo
				json.Unmarshal(data, &user)
				ctx := context.WithValue(r.Context(), "user", &user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// 3. 未登录 → 重定向到 SSO Server
		loginURL := fmt.Sprintf("%s/login?service=%s", c.ssoServerURL, c.serviceURL+r.URL.Path)
		http.Redirect(w, r, loginURL, http.StatusFound)
	})
}

// validateTicket 向 SSO Server 验证 Ticket
func (c *SSOClient) validateTicket(ctx context.Context, ticket string) (*TicketInfo, error) {
	url := fmt.Sprintf("%s/validate?ticket=%s&service=%s", c.ssoServerURL, ticket, c.serviceURL)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ticket 验证失败: %d", resp.StatusCode)
	}

	var info TicketInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

func generateToken(prefix string) string {
	b := make([]byte, 24)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}

func authenticate(username, password string) (int64, bool) {
	// 实际实现: 查数据库、校验密码哈希
	if username == "admin" && password == "123456" {
		return 1001, true
	}
	return 0, false
}
```

---

### 6. 方案三：OAuth2 / OIDC（开放平台登录）

```
OAuth2 = Open Authorization 2.0（授权协议，RFC 6749）
OIDC   = OpenID Connect（在 OAuth2 上层增加身份认证）

适用:
  - 第三方登录（微信登录、Google 登录、GitHub 登录）
  - 开放平台（允许第三方应用访问用户资源）
  - 大型企业内部系统统一认证

OAuth2 四种授权模式:
  ┌──────────────────────────┬────────────────┬──────────────────────┐
  │ 模式                     │ 安全级别        │ 适用场景              │
  ├──────────────────────────┼────────────────┼──────────────────────┤
  │ 授权码模式(最常用)         │ 最高           │ 有后端的 Web 应用      │
  │ PKCE 增强授权码            │ 高             │ SPA / 移动端          │
  │ 客户端凭证模式             │ 中             │ 服务间调用(M2M)       │
  │ 密码模式(已弃用)           │ 低             │ 高度信任的第一方应用    │
  └──────────────────────────┴────────────────┴──────────────────────┘

授权码模式流程（最常用）:

┌────────┐             ┌──────────────┐          ┌──────────────┐
│ 浏览器  │             │ 授权服务器     │          │ 子系统后端    │
└───┬────┘             └──────┬───────┘          └──────┬───────┘
    │ 1.点击"微信登录"         │                        │
    │─────────────────────────────────────────────────>│
    │                         │  2.构造授权URL并重定向    │
    │<──── 302 ───────────────────────────────────────│
    │ 3.跳转到授权服务器        │                        │
    │────────────────────────>│                        │
    │ 4.用户登录并授权          │                        │
    │────────────────────────>│                        │
    │                         │ 5.生成授权码 code        │
    │<── 302 redirect_uri ────│  (redirect_uri?code=xxx)│
    │ 6.携带code回调子系统      │                        │
    │─────────────────────────────────────────────────>│
    │                         │  7.后端用code换token     │
    │                         │<───────────────────────│
    │                         │  8.返回access_token     │
    │                         │───────────────────────>│
    │                         │  9.用token获取用户信息   │
    │                         │<───────────────────────│
    │                         │  10.返回用户信息         │
    │                         │───────────────────────>│
    │  11.登录成功              │                        │
    │<─────────────────────── 设置本地Session ──────────│

安全要点:
  - code 一次性使用，有效期极短（通常 10 分钟）
  - access_token 仅在后端传输，不暴露给前端
  - 必须校验 state 参数（防 CSRF）
  - redirect_uri 必须白名单校验
```

#### OAuth2 vs OIDC

```
OAuth2:
  - 只管 "授权"（Authorization）→ 允许第三方访问你的资源
  - 返回 access_token
  - 不提供标准化的用户信息格式

OIDC（OpenID Connect）:
  - 在 OAuth2 基础上增加 "认证"（Authentication）→ 确认你是谁
  - 额外返回 id_token（JWT 格式，包含用户身份信息）
  - 标准化的 UserInfo Endpoint
  - 标准化的 scope: openid / profile / email

实际应用中:
  - 做 SSO → 用 OIDC（需要知道用户是谁）
  - 做第三方 API 授权 → 用 OAuth2（只需要权限令牌）
```

---

### 7. JWT vs Session 对比

| 维度 | Session（服务端有状态） | JWT（服务端无状态） |
|------|----------------------|-------------------|
| **存储位置** | 服务端（Redis/内存） | 客户端（Cookie/localStorage） |
| **扩展性** | 需要共享 Session 存储 | 天然支持分布式（无需共享） |
| **性能** | 每次请求查 Redis | 本地解密验签，无网络开销 |
| **安全性** | 服务端可控，可随时踢出 | 签发后无法主动失效（除非黑名单） |
| **注销** | 删除 Session 即可 | 需要维护黑名单或等待过期 |
| **Token 大小** | Cookie 仅存 SessionID（几十字节） | JWT 较大（几百字节~1KB） |
| **跨域** | 需要特殊处理 | Bearer Token 天然跨域 |
| **移动端** | 不适合（依赖 Cookie） | 适合（Header 传递） |
| **适用场景** | 传统 Web、需要踢人、会话管理 | 微服务、API 网关、移动端 |

#### JWT 结构

```
JWT = Header.Payload.Signature（三段 Base64URL 编码，用 . 分隔）

Header:                        Payload:                    Signature:
{                              {                           HMACSHA256(
  "alg": "RS256",                "sub": "1001",              base64UrlEncode(header)+"."+
  "typ": "JWT"                   "name": "张三",              base64UrlEncode(payload),
}                                "role": "admin",            secret
                                 "iat": 1710000000,        )
                                 "exp": 1710086400
                               }

安全注意事项:
  1. 永远不要在 Payload 中放密码等敏感信息（Base64 可解码）
  2. 使用 RS256（非对称）而非 HS256（对称），防止密钥泄露
  3. 设置合理的过期时间（access_token: 15min，refresh_token: 7d）
  4. 使用 refresh_token 机制实现续期
```

#### JWT + Refresh Token 续期流程

```
┌────────┐              ┌──────────┐              ┌────────┐
│  客户端  │              │ API 网关  │              │ 认证服务 │
└───┬────┘              └────┬─────┘              └───┬────┘
    │ 1.携带access_token请求  │                       │
    │────────────────────────>│                       │
    │                        │ 2.验证JWT(本地验签)     │
    │                        │ → 未过期 → 放行         │
    │<── 正常响应 ────────────│                       │
    │                        │                       │
    │ 3.access_token过期      │                       │
    │────────────────────────>│                       │
    │                        │ 4.JWT过期，返回401      │
    │<── 401 ────────────────│                       │
    │                        │                       │
    │ 5.用refresh_token续期   │                       │
    │─────────────────────────────────────────────────>│
    │                        │  6.验证refresh_token   │
    │                        │  7.签发新access_token   │
    │<── 新access_token ──────────────────────────────│
    │                        │                       │
    │ 8.用新token重新请求      │                       │
    │────────────────────────>│                       │

Refresh Token 安全策略:
  - 存储在 HttpOnly Cookie 中（防 XSS）
  - 每次使用后轮换（Refresh Token Rotation）
  - 检测到重复使用 → 立即吊销所有 Token（说明被窃取）
```

#### HS256 vs RS256 签名算法对比

```
HS256（HMAC-SHA256）— 对称签名:
  签发方和验证方使用同一个 secret
  JWT签名 = HMAC-SHA256(header.payload, secret)

  ┌──────────┐    secret     ┌──────────┐
  │  认证服务  │ ──────────── │  API 网关  │
  │  (签发)   │  共享密钥     │  (验证)   │
  └──────────┘              └──────────┘

  ❌ 问题: 所有验证方都需要 secret → 任何一个泄露全部完蛋
  ❌ 验证方也能伪造 Token（因为它有 secret）
  ✅ 适合: 单体应用 / 签发和验证是同一个服务

RS256（RSA-SHA256）— 非对称签名:
  签发方用私钥签名，验证方用公钥验证
  JWT签名 = RSA-SHA256(header.payload, privateKey)
  验证 = RSA-Verify(header.payload, signature, publicKey)

  ┌──────────┐              ┌──────────┐
  │  认证服务  │  privateKey  │  API 网关  │  publicKey
  │  (签发)   │  (只有自己有) │  (验证)   │  (可公开)
  └──────────┘              └──────────┘
                               │
                            ┌──▼───────┐
                            │ 子系统 A  │  publicKey
                            │ (验证)   │  (可公开)
                            └──────────┘

  ✅ 验证方无法伪造 Token（没有私钥）
  ✅ 公钥可以随意分发，泄露也无所谓
  ✅ 适合: 微服务架构 / 多系统验证

  ┌──────────┬──────────────┬──────────────┐
  │ 维度      │ HS256        │ RS256        │
  ├──────────┼──────────────┼──────────────┤
  │ 密钥      │ 对称(一个)    │ 非对称(公私钥)│
  │ 签发      │ secret       │ privateKey   │
  │ 验证      │ secret       │ publicKey    │
  │ 性能      │ 快           │ 稍慢(RSA运算) │
  │ 安全性    │ 低(共享密钥)  │ 高(私钥隔离)  │
  │ 适用场景  │ 单体应用      │ 微服务/SSO   │
  └──────────┴──────────────┴──────────────┘
```

#### JWT 验证完整流程

```
收到请求 Bearer Token 后的完整验证链:

  ┌────────────────────────────────────────────────────────────┐
  │                    JWT 验证流程                              │
  │                                                            │
  │  ① 提取 Token                                              │
  │     Authorization: Bearer eyJhbGciOi...                    │
  │     ↓                                                      │
  │  ② 解码 Header（Base64URL Decode，不需要密钥）              │
  │     → 获取 alg 字段 → 确认使用 RS256（防止 alg:none 攻击）  │
  │     ↓                                                      │
  │  ③ 验证签名（用公钥验证）                                    │
  │     publicKey.Verify(header + "." + payload, signature)    │
  │     → 失败 → 401 Unauthorized（Token 被篡改）               │
  │     ↓                                                      │
  │  ④ 验证过期时间                                              │
  │     payload.exp > now()                                     │
  │     → 过期 → 401 Unauthorized（需要 Refresh）               │
  │     ↓                                                      │
  │  ⑤ 验证签发者（iss）                                        │
  │     payload.iss == "https://sso.example.com"               │
  │     → 不匹配 → 403 Forbidden（非法签发者）                  │
  │     ↓                                                      │
  │  ⑥ 验证受众（aud）                                          │
  │     payload.aud 包含当前服务标识                              │
  │     → 不匹配 → 403 Forbidden（Token 不是签发给本服务的）     │
  │     ↓                                                      │
  │  ⑦ 检查黑名单（Redis）                                      │
  │     EXISTS blacklist:jwt:{jti}                              │
  │     → 存在 → 401 Unauthorized（Token 已被吊销）             │
  │     ↓                                                      │
  │  ⑧ 验证通过 → 从 payload 提取用户信息 → 注入 Context        │
  │                                                            │
  └────────────────────────────────────────────────────────────┘

  常见攻击与防御:
  ┌───────────────────┬────────────────────────────────────────┐
  │ 攻击               │ 防御                                   │
  ├───────────────────┼────────────────────────────────────────┤
  │ alg:none 攻击      │ 服务端强制指定算法，忽略 Header 中的 alg│
  │ HS256 降级攻击     │ 用 RS256 公钥当 HS256 secret 验证      │
  │                   │ → 强制只接受 RS256                      │
  │ Token 泄露        │ 短过期(15min) + HTTPS + HttpOnly Cookie │
  │ 重放攻击          │ jti(唯一ID) + 短过期 + 黑名单           │
  │ Payload 篡改      │ 签名验证（任何修改都会导致签名不匹配）   │
  └───────────────────┴────────────────────────────────────────┘
```

#### JWT 完整 Go 实现（RS256 签发 + 验证 + 黑名单 + 中间件）

```go
package jwtauth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenMissing  = errors.New("缺少 Authorization Header")
	ErrTokenInvalid  = errors.New("Token 无效或已篡改")
	ErrTokenExpired  = errors.New("Token 已过期")
	ErrTokenRevoked  = errors.New("Token 已被吊销")
	ErrInvalidIssuer = errors.New("签发者不合法")
)

// ---------- Claims 定义 ----------

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ---------- JWT Manager ----------

// JWTManager JWT 管理器（签发 + 验证 + 黑名单）
type JWTManager struct {
	privateKey       *rsa.PrivateKey // 私钥（仅认证服务持有，用于签发）
	publicKey        *rsa.PublicKey  // 公钥（所有服务持有，用于验证）
	rdb              *redis.Client  // Redis（黑名单）
	issuer           string         // 签发者标识
	accessTokenTTL   time.Duration  // access_token 有效期
	refreshTokenTTL  time.Duration  // refresh_token 有效期
}

type JWTConfig struct {
	PrivateKeyPath  string        // 私钥文件路径（认证服务需要）
	PublicKeyPath   string        // 公钥文件路径（所有服务需要）
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWTManager(cfg *JWTConfig, rdb *redis.Client) (*JWTManager, error) {
	mgr := &JWTManager{
		rdb:             rdb,
		issuer:          cfg.Issuer,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}

	// 加载公钥（所有服务都需要）
	pubKeyData, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取公钥失败: %w", err)
	}
	pubBlock, _ := pem.Decode(pubKeyData)
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}
	mgr.publicKey = pubKey.(*rsa.PublicKey)

	// 加载私钥（仅认证服务需要，其他服务可不配置）
	if cfg.PrivateKeyPath != "" {
		privKeyData, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("读取私钥失败: %w", err)
		}
		privBlock, _ := pem.Decode(privKeyData)
		privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
		mgr.privateKey = privKey
	}

	return mgr, nil
}

// ---------- 签发 Token（仅认证服务调用） ----------

// TokenPair access_token + refresh_token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access_token 剩余秒数
	TokenType    string `json:"token_type"` // "Bearer"
}

// GenerateTokenPair 签发 Token 对
func (m *JWTManager) GenerateTokenPair(ctx context.Context, userID int64, username, role string) (*TokenPair, error) {
	if m.privateKey == nil {
		return nil, errors.New("私钥未配置，无法签发 Token")
	}

	now := time.Now()

	// 1. 签发 access_token
	accessJTI := uuid.New().String() // 唯一 ID（用于黑名单吊销）
	accessClaims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessJTI,
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:   userID,
		Username: username,
		Role:     role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessStr, err := accessToken.SignedString(m.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签发 access_token 失败: %w", err)
	}

	// 2. 签发 refresh_token（Payload 精简，只包含必要信息）
	refreshJTI := uuid.New().String()
	refreshClaims := &jwt.RegisteredClaims{
		ID:        refreshJTI,
		Issuer:    m.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenTTL)),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(m.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签发 refresh_token 失败: %w", err)
	}

	// 3. 将 refresh_token 的 JTI 存入 Redis（用于 Rotation 检测）
	// value = userID，用于续期时查用户信息
	m.rdb.Set(ctx, "jwt:refresh:"+refreshJTI, userID, m.refreshTokenTTL)

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresIn:    int64(m.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ---------- 验证 Token（所有服务调用） ----------

// VerifyAccessToken 验证 access_token 并返回 Claims
func (m *JWTManager) VerifyAccessToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	// 1. 解析并验证签名 + 过期时间
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 防御 alg 切换攻击：强制只接受 RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 2. 验证签发者
	if claims.Issuer != m.issuer {
		return nil, ErrInvalidIssuer
	}

	// 3. 检查黑名单（Token 是否已被吊销）
	exists, err := m.rdb.Exists(ctx, "jwt:blacklist:"+claims.ID).Result()
	if err == nil && exists > 0 {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

// ---------- Refresh Token 续期 ----------

// RefreshTokens 用 refresh_token 换取新的 Token 对
func (m *JWTManager) RefreshTokens(ctx context.Context, refreshTokenStr string) (*TokenPair, error) {
	// 1. 解析 refresh_token
	token, err := jwt.ParseWithClaims(refreshTokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("refresh_token 无效: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("refresh_token 无效")
	}

	// 2. Refresh Token Rotation 检测
	//    正常: Redis 中存在该 JTI → 删除旧的，签发新的
	//    异常: Redis 中不存在 → 说明已被使用过（可能被窃取重放）→ 吊销该用户所有 Token
	userIDStr, err := m.rdb.Get(ctx, "jwt:refresh:"+claims.ID).Result()
	if err != nil {
		// refresh_token 已被使用或不存在 → 可能被窃取
		// 安全策略: 吊销该用户的所有 Token
		m.RevokeUserTokens(ctx, claims.Subject)
		return nil, errors.New("refresh_token 已失效（疑似被盗用，已吊销所有会话）")
	}

	// 3. 删除旧 refresh_token（一次性使用）
	m.rdb.Del(ctx, "jwt:refresh:"+claims.ID)

	// 4. 查询用户信息并签发新 Token 对
	//    生产环境从 DB 或缓存查最新用户信息（角色可能已变更）
	var userID int64
	fmt.Sscanf(userIDStr, "%d", &userID)

	// 此处简化，实际应查 DB 获取最新的 username 和 role
	return m.GenerateTokenPair(ctx, userID, "", "")
}

// ---------- Token 吊销（登出/踢人） ----------

// RevokeAccessToken 吊销单个 access_token（加入黑名单）
func (m *JWTManager) RevokeAccessToken(ctx context.Context, tokenStr string) error {
	// 解析 Token 获取 jti 和过期时间（不做完整验证，因为可能已过期）
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, &CustomClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return errors.New("解析 Claims 失败")
	}

	// 黑名单 TTL = Token 剩余有效时间（过期后自动从黑名单移除，节省空间）
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil // 已过期，无需加黑名单
	}

	return m.rdb.Set(ctx, "jwt:blacklist:"+claims.ID, "1", ttl).Err()
}

// RevokeUserTokens 吊销某用户的所有 Token（账号被盗/修改密码时使用）
func (m *JWTManager) RevokeUserTokens(ctx context.Context, userIDStr string) {
	// 方案: 在 Redis 中记录该用户的"Token 全部失效时间戳"
	// 验证时: 如果 Token 的 iat < 失效时间戳 → 拒绝
	m.rdb.Set(ctx, "jwt:user_revoke:"+userIDStr,
		time.Now().Unix(), 24*time.Hour)
}

// IsUserTokenRevoked 检查用户的 Token 是否在全局吊销范围内
func (m *JWTManager) IsUserTokenRevoked(ctx context.Context, userID string, issuedAt time.Time) bool {
	revokeTimeStr, err := m.rdb.Get(ctx, "jwt:user_revoke:"+userID).Result()
	if err != nil {
		return false // 无吊销记录
	}
	var revokeTime int64
	fmt.Sscanf(revokeTimeStr, "%d", &revokeTime)
	return issuedAt.Unix() < revokeTime
}

// ---------- HTTP 中间件 ----------

type contextKey string

const UserClaimsKey contextKey = "user_claims"

// AuthMiddleware JWT 认证中间件（用于所有需要登录的接口）
func (m *JWTManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 从 Header 提取 Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpUnauthorized(w, "缺少 Authorization Header")
			return
		}

		// 格式: "Bearer eyJhbGciOi..."
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpUnauthorized(w, "Authorization 格式错误，需要 Bearer Token")
			return
		}
		tokenStr := parts[1]

		// 2. 验证 Token
		claims, err := m.VerifyAccessToken(r.Context(), tokenStr)
		if err != nil {
			switch {
			case errors.Is(err, ErrTokenExpired):
				httpUnauthorized(w, "Token 已过期，请刷新")
			case errors.Is(err, ErrTokenRevoked):
				httpUnauthorized(w, "Token 已被吊销，请重新登录")
			case errors.Is(err, ErrInvalidIssuer):
				httpForbidden(w, "非法的 Token 签发者")
			default:
				httpUnauthorized(w, "Token 无效")
			}
			return
		}

		// 3. 检查用户级别吊销（修改密码/账号被盗场景）
		if m.IsUserTokenRevoked(r.Context(), claims.Subject, claims.IssuedAt.Time) {
			httpUnauthorized(w, "登录已失效，请重新登录")
			return
		}

		// 4. 将用户信息注入 Context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware 角色权限中间件（RBAC）
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserClaimsKey).(*CustomClaims)
			if !ok {
				httpUnauthorized(w, "未认证")
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			httpForbidden(w, "权限不足")
		})
	}
}

// GetCurrentUser 从 Context 获取当前用户
func GetCurrentUser(ctx context.Context) *CustomClaims {
	claims, _ := ctx.Value(UserClaimsKey).(*CustomClaims)
	return claims
}

func httpUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"code":401,"msg":"%s"}`, msg)
}

func httpForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, `{"code":403,"msg":"%s"}`, msg)
}
```

#### JWT 认证服务 HTTP 接口示例

```go
package jwtauth

import (
	"encoding/json"
	"net/http"
)

// AuthHandler 认证服务的 HTTP 接口
type AuthHandler struct {
	jwtMgr *JWTManager
}

// HandleLogin 登录接口 → 签发 Token 对
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// 验证用户名密码（实际查 DB + bcrypt 比对）
	userID, role, ok := authenticate(req.Username, req.Password)
	if !ok {
		httpUnauthorized(w, "用户名或密码错误")
		return
	}

	// 签发 Token 对
	tokenPair, err := h.jwtMgr.GenerateTokenPair(r.Context(), userID, req.Username, role)
	if err != nil {
		http.Error(w, "签发 Token 失败", http.StatusInternalServerError)
		return
	}

	// refresh_token 通过 HttpOnly Cookie 返回（防 XSS）
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/api/auth/refresh",          // 只在刷新接口发送
		MaxAge:   int(7 * 24 * 3600),           // 7天
		HttpOnly: true,                          // JavaScript 无法读取
		Secure:   true,                          // 仅 HTTPS
		SameSite: http.SameSiteStrictMode,       // 防 CSRF
	})

	// access_token 通过 JSON Body 返回（前端存内存，不存 localStorage）
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   tokenPair.TokenType,
	})
}

// HandleRefresh 刷新 Token 接口
func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// 从 HttpOnly Cookie 获取 refresh_token
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpUnauthorized(w, "缺少 refresh_token")
		return
	}

	// 用 refresh_token 换取新 Token 对
	tokenPair, err := h.jwtMgr.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		httpUnauthorized(w, err.Error())
		return
	}

	// 设置新的 refresh_token Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/api/auth/refresh",
		MaxAge:   int(7 * 24 * 3600),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   tokenPair.TokenType,
	})
}

// HandleLogout 登出接口
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// 1. 吊销 access_token（加入黑名单）
	authHeader := r.Header.Get("Authorization")
	if parts := splitBearer(authHeader); parts != "" {
		h.jwtMgr.RevokeAccessToken(r.Context(), parts)
	}

	// 2. 清除 refresh_token Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"已登出"}`))
}

func splitBearer(header string) string {
	parts := make([]byte, 0, len(header))
	if len(header) > 7 && header[:7] == "Bearer " {
		parts = append(parts, header[7:]...)
		return string(parts)
	}
	return ""
}

func authenticate(username, password string) (int64, string, bool) {
	// 实际: 查 DB + bcrypt.CompareHashAndPassword
	if username == "admin" && password == "123456" {
		return 1001, "admin", true
	}
	return 0, "", false
}
```

#### 路由注册与使用示例

```go
func SetupRoutes(mux *http.ServeMux, jwtMgr *JWTManager, authHandler *AuthHandler) {
	// 公开接口（无需 Token）
	mux.HandleFunc("/api/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("/api/auth/refresh", authHandler.HandleRefresh)

	// 需要登录的接口
	protected := jwtMgr.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetCurrentUser(r.Context())
		fmt.Fprintf(w, `{"msg":"hello %s, role=%s"}`, user.Username, user.Role)
	}))
	mux.Handle("/api/user/profile", protected)

	// 需要登录 + 需要 admin 角色
	adminOnly := jwtMgr.AuthMiddleware(
		RoleMiddleware("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"msg":"admin area"}`))
		})),
	)
	mux.Handle("/api/admin/dashboard", adminOnly)

	// 登出
	mux.Handle("/api/auth/logout", jwtMgr.AuthMiddleware(
		http.HandlerFunc(authHandler.HandleLogout),
	))
}
```

#### JWT Token 存储安全总结

```
前端 Token 存储方案对比:
  ┌──────────────────┬────────────┬────────────┬──────────────────┐
  │ 存储位置          │ XSS 风险    │ CSRF 风险   │ 推荐度            │
  ├──────────────────┼────────────┼────────────┼──────────────────┤
  │ localStorage     │ ❌ 高       │ ✅ 无       │ ❌ 不推荐         │
  │ (JS可读取)       │ XSS可偷Token│            │                  │
  ├──────────────────┼────────────┼────────────┼──────────────────┤
  │ sessionStorage   │ ❌ 高       │ ✅ 无       │ ❌ 不推荐         │
  │ (JS可读取)       │ XSS可偷Token│ 标签页不共享│                  │
  ├──────────────────┼────────────┼────────────┼──────────────────┤
  │ HttpOnly Cookie  │ ✅ 安全     │ ❌ 有       │ ⚠️ 需配合 CSRF   │
  │ (JS不可读取)     │            │ 自动携带    │  Token 防护       │
  ├──────────────────┼────────────┼────────────┼──────────────────┤
  │ 内存变量         │ ✅ 安全     │ ✅ 无       │ ✅ 推荐           │
  │ (JS闭包/状态管理)│ XSS拿不到  │ 手动携带    │ 刷新页面需重新获取 │
  └──────────────────┴────────────┴────────────┴──────────────────┘

  最佳实践:
    access_token  → 存内存变量（Vue/React 状态管理）
                    → 通过 Authorization Header 手动携带
                    → 页面刷新后用 refresh_token 重新获取

    refresh_token → 存 HttpOnly + Secure + SameSite=Strict Cookie
                    → 只在 /api/auth/refresh 路径发送
                    → 每次使用后轮换（Rotation）
```

#### JWT 签名密钥是怎么保证安全的？

```
问题: RS256 的私钥 / HS256 的 secret 如果泄露，任何人都能伪造 Token
     → 密钥安全 = JWT 安全的根基

一、密钥存储方案对比:

  ┌───────────────────┬───────────┬──────────────────────────────────┐
  │ 存储方式           │ 安全级别   │ 说明                              │
  ├───────────────────┼───────────┼──────────────────────────────────┤
  │ 代码硬编码         │ ❌ 极低    │ 代码泄露 = 密钥泄露               │
  │ 配置文件(明文)     │ ❌ 低      │ 文件泄露即暴露                    │
  │ 环境变量           │ ⚠️ 中     │ ps/env 可查看，容器环境较安全      │
  │ 配置中心(加密)     │ ✅ 高     │ Apollo/Vault 加密存储+访问控制    │
  │ KMS 密钥管理服务   │ ✅ 极高   │ AWS KMS / 阿里云 KMS / HashiCorp │
  │                   │           │ Vault，密钥从不离开 KMS           │
  │ HSM 硬件安全模块   │ ✅ 最高   │ 物理硬件存储密钥，金融级安全       │
  └───────────────────┴───────────┴──────────────────────────────────┘

二、生产环境推荐架构（配置中心 + KMS）:

  ┌─────────────┐
  │ 认证服务     │  启动时从 Vault 获取私钥 → 存内存
  │ (签发JWT)   │  私钥不落盘、不进环境变量
  └──────┬──────┘
         │ 启动时拉取 (TLS + Token 认证)
         ▼
  ┌─────────────┐     ┌─────────────┐
  │ HashiCorp   │     │ 云 KMS       │
  │ Vault       │ 或  │(AWS/阿里云)  │
  │ 密钥存储     │     │ 托管密钥     │
  └─────────────┘     └─────────────┘

  流程:
    1. 密钥由 KMS/Vault 生成并管理（人看不到原文）
    2. 认证服务启动时通过 API + 短期凭证获取私钥
    3. 私钥只存在进程内存中，不写磁盘、不写日志
    4. 公钥通过 JWKS Endpoint 公开发布（安全，可随意分发）

三、JWKS（JSON Web Key Set）公钥分发:

  认证服务暴露标准 JWKS 端点:
    GET https://sso.example.com/.well-known/jwks.json

    响应:
    {
      "keys": [{
        "kty": "RSA",
        "kid": "2024-key-1",          ← 密钥ID（用于密钥轮换）
        "use": "sig",
        "alg": "RS256",
        "n":   "0vx7agoebGc...",      ← 公钥 modulus (Base64URL)
        "e":   "AQAB"                 ← 公钥 exponent
      }]
    }

  各子系统/API网关:
    → 启动时从 JWKS 端点拉取公钥
    → 定期刷新（如每 5 分钟）
    → 验证 JWT 时用 kid 匹配对应公钥

  优点:
    - 公钥可公开，无安全风险
    - 密钥轮换时只需更新 JWKS，子系统自动拉取新公钥
    - OIDC 标准规范，Keycloak/Auth0 等默认支持

四、密钥轮换（Key Rotation）:

  为什么需要轮换？
    - 密钥使用时间越长泄露风险越大
    - 合规要求（PCI-DSS 等要求定期轮换）

  轮换流程:
    ① 生成新密钥对 (kid=2025-key-1)
    ② JWKS 中同时发布新旧两个公钥
    ③ 认证服务开始用新私钥签发 Token
    ④ 旧 Token 仍然可以用旧公钥验证（kid 匹配）
    ⑤ 等旧 Token 全部过期后（如 24 小时）→ 从 JWKS 移除旧公钥

    时间线:
    ──────┬──────────────────┬─────────────────┬──────────
          T0               T1                T2
       发布新公钥        切换到新私钥签发    移除旧公钥
       JWKS=[旧,新]      JWKS=[旧,新]       JWKS=[新]
```

```go
// JWKS 端点示例（认证服务提供）
func HandleJWKS(w http.ResponseWriter, r *http.Request) {
    // 从内存中获取当前活跃的公钥列表
    keys := getActivePublicKeys() // []*rsa.PublicKey with kid

    jwks := map[string]interface{}{
        "keys": formatAsJWKS(keys),
    }
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "public, max-age=300") // 缓存5分钟
    json.NewEncoder(w).Encode(jwks)
}

// API 网关 / 子系统: 从 JWKS 端点获取公钥
type JWKSVerifier struct {
    jwksURL    string
    keys       map[string]*rsa.PublicKey // kid → publicKey
    mu         sync.RWMutex
}

func (v *JWKSVerifier) RefreshKeys(ctx context.Context) error {
    resp, err := http.Get(v.jwksURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // 解析 JWKS JSON → 提取公钥 → 更新 v.keys
    // ...
    return nil
}

// 验证 JWT 时根据 Header 中的 kid 选择公钥
func (v *JWKSVerifier) GetPublicKey(kid string) *rsa.PublicKey {
    v.mu.RLock()
    defer v.mu.RUnlock()
    return v.keys[kid]
}
```

#### JWT 信息是服务端下发给客户端的吗？

```
完整回答: 是的，JWT 由服务端签发并下发给客户端，但客户端可以"看"不能"改"。

一、JWT 的签发与传递流程:

  ┌────────┐                          ┌──────────┐
  │  客户端  │                          │ 认证服务   │
  └───┬────┘                          └────┬─────┘
      │ 1. POST /login                     │
      │    {username, password}             │
      │───────────────────────────────────>│
      │                                    │ 2. 验证密码
      │                                    │ 3. 生成 JWT:
      │                                    │    Header = {"alg":"RS256"}
      │                                    │    Payload = {"uid":1001,"role":"admin","exp":...}
      │                                    │    Signature = RSA-Sign(Header.Payload, 私钥)
      │                                    │ 4. 拼接: Header.Payload.Signature
      │ 5. 返回 JWT                         │
      │<───────────────────────────────────│
      │                                    │
      │ 6. 后续请求携带 JWT                  │
      │    Authorization: Bearer eyJ...     │
      │───────────────────────────────────>│ 或其他业务服务
      │                                    │ 7. 用公钥验证签名
      │                                    │ 8. 从 Payload 取用户信息

  关键点:
    - JWT 完全由服务端生成（客户端不参与生成过程）
    - 客户端只负责 "存储" 和 "携带" JWT
    - 客户端不需要任何密钥

二、客户端能解码 JWT 吗？

  ┌─────────────────────────────────────────────────────────────┐
  │ JWT = Base64URL(Header) . Base64URL(Payload) . Signature   │
  │                                                             │
  │ Base64URL 是编码，不是加密！                                 │
  │ 任何人（包括客户端）都可以 Base64 解码看到 Header 和 Payload  │
  │                                                             │
  │ 示例: eyJhbGciOiJSUzI1NiJ9.eyJ1aWQiOjEwMDEsInJvbGUiOiJhZG1pbiJ9.xxx │
  │       ↓ Base64 Decode                                       │
  │       {"alg":"RS256"}  {"uid":1001,"role":"admin"}          │
  │                                                             │
  │ 所以: ❌ 绝对不要在 Payload 中放密码、密钥、身份证号等敏感信息│
  │       ✅ 只放用户ID、角色、过期时间等业务必要信息             │
  └─────────────────────────────────────────────────────────────┘

三、客户端能篡改 JWT 吗？

  ❌ 不能。即使客户端修改了 Payload 中的 role:"admin"→"superadmin"
  服务端用公钥验证签名时会发现 Signature 不匹配 → 拒绝

  因为 Signature = RSA-Sign(Header + "." + Payload, 私钥)
  Payload 被改 → 重新计算签名需要私钥 → 客户端没有私钥 → 无法伪造

四、JWT vs 传统加密 Token 的区别:

  ┌──────────────────┬─────────────────────┬─────────────────────┐
  │ 维度              │ JWT (签名型)          │ 加密 Token (不透明)  │
  ├──────────────────┼─────────────────────┼─────────────────────┤
  │ 客户端可读取信息  │ ✅ 可以解码Payload   │ ❌ 不可以            │
  │ 客户端可篡改     │ ❌ 签名验证会失败     │ ❌ 解密会失败        │
  │ 服务端需查DB     │ ❌ 自包含(Payload)   │ ✅ 需查DB/Redis     │
  │ 信息透明度       │ 透明（明文编码）      │ 不透明（加密）       │
  │ 适用场景         │ 服务端无状态          │ 需要严格保密Payload  │
  └──────────────────┴─────────────────────┴─────────────────────┘

  JWT 的哲学:
    不追求"信息保密" → 追求"信息不可篡改"
    靠签名(Signature)保证完整性，不靠加密保证机密性

    如果确实需要加密 Payload（如包含手机号等PII）:
    → 使用 JWE (JSON Web Encryption) 而不是 JWS (JSON Web Signature)
    → 或者 Payload 中只放 userID，敏感信息查 DB

五、前端为什么要解码 JWT？

  常见用途（合法且推荐）:
    1. 读取 exp 字段 → 判断 Token 是否即将过期 → 提前静默刷新
    2. 读取用户信息(uid/name/role) → 前端展示，避免额外 API 调用
    3. 读取 role → 前端路由守卫（仅做 UI 显示控制，权限校验在后端）

  JavaScript 解码 JWT:
    function parseJWT(token) {
        const payload = token.split('.')[1];
        return JSON.parse(atob(payload));
    }
    const user = parseJWT(accessToken);
    // { uid: 1001, username: "张三", role: "admin", exp: 1710086400 }

  ⚠️ 注意: 前端读取的信息仅用于 UI 展示
           真正的权限校验必须在服务端（因为客户端数据不可信）
```

---

### 8. SSO 全局登出方案

```
难点: 用户在 SSO 登出后，如何通知所有已登录的子系统？

方案一: 前端跳转登出（简单但不彻底）
  1. 前端跳转到 SSO Server 的 /logout
  2. SSO Server 清除 TGC Cookie
  3. 返回页面中嵌入多个 <iframe src="app-a.com/logout"> 通知子系统

方案二: 后端回调登出（CAS 标准方案）
  1. SSO Server 维护每个 TGT 关联的所有子系统列表
  2. 登出时，SSO Server 向所有子系统发送 HTTP 回调
  3. 子系统收到回调 → 清除局部 Session

  ┌────────┐       ┌──────────┐
  │ 浏览器  │──────>│SSO Server│
  │ 登出    │       └────┬─────┘
  └────────┘            │
                        ├── POST app-a.com/sso/logout {session_id}
                        ├── POST app-b.com/sso/logout {session_id}
                        └── POST app-c.com/sso/logout {session_id}

方案三: Token 短过期 + 黑名单（JWT 场景）
  1. access_token 设置较短过期时间（如 15 分钟）
  2. 登出时将 token 加入 Redis 黑名单
  3. API 网关验证时先查黑名单

  优点: 不需要通知子系统
  缺点: 最多有 15 分钟延迟（或需要维护黑名单）
```

---

### 9. 生产环境 SSO 架构建议

```
推荐架构: OAuth2/OIDC + JWT + API 网关

┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Web 前端   │     │  移动端 App  │     │  第三方应用   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │  API 网关    │  ← JWT 验签（本地，无网络开销）
                    │ (Kong/APISIX)│  ← Token 黑名单检查（Redis）
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────▼──────┐ ┌────▼──────┐ ┌─────▼─────┐
     │  用户服务    │ │ 订单服务   │ │ 商品服务   │
     └─────────────┘ └───────────┘ └───────────┘
                           │
                    ┌──────▼──────┐
                    │  认证中心    │  ← OIDC Provider
                    │ (Keycloak)  │  ← 登录/注册/Token签发
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  用户数据库   │
                    │ (MySQL+Redis)│
                    └─────────────┘

技术选型:
  │ 组件             │ 推荐方案                          │
  │ OIDC Provider   │ Keycloak / Hydra / Casdoor        │
  │ API 网关         │ Kong / APISIX（内置 JWT 插件）     │
  │ Token 格式       │ JWT (RS256) + Refresh Token       │
  │ Session 存储     │ Redis Cluster                     │
  │ 通信安全         │ TLS 1.3 + HTTPS                   │
```

### 一句话总结

**SSO 三种主流方案：同域用 Session+Cookie（Cookie domain 设为主域名+Redis 共享 Session，最简单）；跨域用 CAS 协议（SSO Server 签发一次性 Ticket，子系统拿 Ticket 换用户信息，浏览器通过 TGC Cookie 实现免登）；开放平台用 OAuth2/OIDC（授权码模式：前端跳转授权→获取 code→后端用 code 换 token→获取用户信息，OIDC 在 OAuth2 上增加 id_token 实现身份认证）。** Session 服务端有状态可随时踢出但需共享存储；JWT 无状态天然分布式但无法主动失效需配合 Refresh Token+黑名单。生产推荐：OIDC Provider（Keycloak）+ JWT(RS256) + API 网关验签 + Redis 黑名单。

---

## 十二、接口幂等性设计

> 如何保证接口幂等？Token 机制？去重表？与 MQ 消费端幂等有何区别？

### 1. 什么是幂等性

```
幂等性 = 同一个操作执行一次和执行多次，产生的效果完全相同

数学定义: f(f(x)) = f(x)

HTTP 方法的幂等性:
  ┌──────────┬────────┬──────────────────────────────────────┐
  │ 方法      │ 幂等？  │ 说明                                 │
  ├──────────┼────────┼──────────────────────────────────────┤
  │ GET      │ ✅ 是   │ 查询操作，天然幂等                     │
  │ PUT      │ ✅ 是   │ 全量更新（覆盖），结果一致              │
  │ DELETE   │ ✅ 是   │ 删除同一资源，结果一致（已删则无操作）   │
  │ HEAD     │ ✅ 是   │ 同 GET，不返回 body                   │
  │ OPTIONS  │ ✅ 是   │ 查询支持的方法                         │
  │ POST     │ ❌ 否   │ 创建资源，每次调用可能创建新资源         │
  │ PATCH    │ ❌ 否   │ 部分更新，看具体实现（如 x=x+1 非幂等） │
  └──────────┴────────┴──────────────────────────────────────┘

需要保证幂等的典型场景:
  1. 支付/扣款:   网络超时 → 用户重试 → 不能扣两次钱
  2. 创建订单:   前端按钮抖动 → 不能创建两个订单
  3. 发送通知:   重试机制 → 不能发两条短信
  4. 库存扣减:   并发请求 → 不能多扣库存
  5. 积分发放:   MQ 重复消费 → 不能重复加积分
```

### 2. 产生重复请求的原因

```
┌──────────────────────────────────────────────────────────┐
│                  重复请求的来源                            │
├───────────────┬──────────────────────────────────────────┤
│ 用户行为       │ 快速双击按钮、刷新页面、浏览器回退后重试    │
│ 网络抖动       │ 请求超时→客户端自动重试，但服务端已处理     │
│ 前端重试       │ axios/fetch 的 retry 机制               │
│ 网关重试       │ Nginx/Kong upstream 超时后自动 retry     │
│ MQ 重复投递    │ Broker 未收到 ACK → 重新投递              │
│ 微服务重试     │ gRPC/HTTP retry policy                  │
│ 定时任务重跑   │ 任务调度器故障恢复后重新执行               │
└───────────────┴──────────────────────────────────────────┘

核心矛盾:
  - 不重试 → 请求丢失（可靠性差）
  - 重试   → 可能重复执行（幂等性差）
  - 解决方案 = "允许重试 + 保证幂等"
```

### 3. 六种幂等方案总览

| 方案 | 原理 | 适用场景 | 优点 | 缺点 |
|------|------|---------|------|------|
| **Token 机制** | 预先获取令牌，请求时携带，用后即删 | 表单提交、创建类接口 | 通用性强 | 需要额外一次请求获取 Token |
| **去重表** | 唯一键约束，插入前检查 | 支付、订单、业务流水 | 可靠性高，可追溯 | 需要额外存储，清理策略 |
| **数据库唯一索引** | DB 层面 UNIQUE 约束 | 创建类操作 | 最简单，DB 兜底 | 只适用于 INSERT 场景 |
| **乐观锁（版本号）** | WHERE version=x 更新 | 更新类操作 | 无锁高性能 | 只适用于 UPDATE 场景 |
| **状态机** | 状态只能单向流转 | 订单、工单等有状态的业务 | 业务语义清晰 | 仅适用于有状态的场景 |
| **Redis SETNX** | 分布式锁/去重 | 高并发短时去重 | 性能高 | 可靠性不如 DB |

---

### 4. 方案一：Token 机制（预防式）

```
核心思想:
  请求前先申请一个一次性 Token，请求时携带 Token
  服务端验证 Token 有效后执行业务，同时删除 Token
  重复请求携带同一 Token → 已被删除 → 拒绝

流程:
  ┌────────┐              ┌──────────┐              ┌────────┐
  │  客户端  │              │  服务端    │              │ Redis  │
  └───┬────┘              └────┬─────┘              └───┬────┘
      │ 1.GET /token            │                       │
      │────────────────────────>│                       │
      │                        │ 2.生成UUID存Redis      │
      │                        │───────────────────────>│
      │                        │   SET token:uuid 1 EX 600
      │ 3.返回token             │                       │
      │<────────────────────────│                       │
      │                        │                       │
      │ 4.POST /order           │                       │
      │   Header: X-Token:uuid  │                       │
      │────────────────────────>│                       │
      │                        │ 5.DEL token:uuid       │
      │                        │───────────────────────>│
      │                        │ 6.DEL 返回1(删除成功)   │
      │                        │<───────────────────────│
      │                        │ 7.执行业务逻辑          │
      │ 8.返回成功              │                       │
      │<────────────────────────│                       │
      │                        │                       │
      │ 9.重复请求(同一token)    │                       │
      │────────────────────────>│                       │
      │                        │ 10.DEL token:uuid      │
      │                        │───────────────────────>│
      │                        │ 11.DEL 返回0(不存在)    │
      │                        │<───────────────────────│
      │ 12.返回"重复请求"        │                       │
      │<────────────────────────│                       │

关键点:
  - 必须用 DEL 的返回值判断（原子操作），不能先 GET 再 DEL（有并发问题）
  - 或者用 Lua 脚本: GET + 比较 + DEL 三步原子化
```

#### Go 实现

```go
package idempotent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenInvalid   = errors.New("幂等 Token 无效或已使用")
	ErrTokenMissing   = errors.New("缺少幂等 Token")
)

// TokenIdempotent Token 机制幂等实现
type TokenIdempotent struct {
	rdb      *redis.Client
	tokenTTL time.Duration
}

func NewTokenIdempotent(rdb *redis.Client) *TokenIdempotent {
	return &TokenIdempotent{
		rdb:      rdb,
		tokenTTL: 10 * time.Minute,
	}
}

// GenerateToken 生成幂等 Token（客户端在提交前调用）
func (t *TokenIdempotent) GenerateToken(ctx context.Context) (string, error) {
	token := generateRandomToken()
	key := "idempotent:token:" + token
	if err := t.rdb.Set(ctx, key, "1", t.tokenTTL).Err(); err != nil {
		return "", err
	}
	return token, nil
}

// VerifyAndConsume 验证并消费 Token（原子操作）
// 返回 nil 表示 Token 有效且已消费，可以执行业务
// 返回 ErrTokenInvalid 表示 Token 已被使用（重复请求）
func (t *TokenIdempotent) VerifyAndConsume(ctx context.Context, token string) error {
	if token == "" {
		return ErrTokenMissing
	}

	key := "idempotent:token:" + token

	// 方式一：直接用 DEL 的返回值（最简单，原子性）
	deleted, err := t.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ErrTokenInvalid // Token 不存在或已被消费
	}
	return nil
}

// verifyAndConsumeLua Lua 脚本版本（支持 Token 值校验，更安全）
var verifyLua = redis.NewScript(`
	local val = redis.call('GET', KEYS[1])
	if val == false then
		return -1  -- Token 不存在
	end
	if val ~= ARGV[1] then
		return -2  -- Token 值不匹配
	end
	redis.call('DEL', KEYS[1])
	return 1       -- 验证通过并已删除
`)

func (t *TokenIdempotent) VerifyAndConsumeLua(ctx context.Context, token, expectedVal string) error {
	key := "idempotent:token:" + token
	result, err := verifyLua.Run(ctx, t.rdb, []string{key}, expectedVal).Int()
	if err != nil {
		return err
	}
	switch result {
	case 1:
		return nil
	case -1:
		return ErrTokenInvalid
	case -2:
		return errors.New("Token 值不匹配")
	default:
		return errors.New("未知错误")
	}
}

// TokenMiddleware 幂等中间件
func (t *TokenIdempotent) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 只对写操作做幂等校验
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("X-Idempotent-Token")
		if err := t.VerifyAndConsume(r.Context(), token); err != nil {
			if errors.Is(err, ErrTokenInvalid) {
				http.Error(w, `{"code":409,"msg":"重复请求，请勿重复提交"}`, http.StatusConflict)
				return
			}
			if errors.Is(err, ErrTokenMissing) {
				http.Error(w, `{"code":400,"msg":"缺少幂等Token"}`, http.StatusBadRequest)
				return
			}
			http.Error(w, `{"code":500,"msg":"服务内部错误"}`, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateRandomToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
```

---

### 5. 方案二：去重表（业务级幂等）

```
核心思想:
  为每个业务请求生成全局唯一的业务 ID（如订单号、流水号）
  在去重表中记录已处理的请求
  收到重复请求时，通过查询去重表直接返回上一次的处理结果

方案 A: 独立去重表

  CREATE TABLE idempotent_record (
      id           BIGINT AUTO_INCREMENT PRIMARY KEY,
      biz_type     VARCHAR(32)  NOT NULL COMMENT '业务类型',
      biz_id       VARCHAR(128) NOT NULL COMMENT '业务唯一标识',
      request_hash VARCHAR(64)  NULL     COMMENT '请求参数摘要',
      response     TEXT         NULL     COMMENT '首次处理结果',
      status       TINYINT      NOT NULL DEFAULT 1 COMMENT '1处理中 2成功 3失败',
      created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      UNIQUE KEY uk_biz (biz_type, biz_id)
  ) ENGINE=InnoDB;

方案 B: 利用业务表唯一索引

  -- 订单表本身的订单号就是唯一的
  CREATE TABLE orders (
      id        BIGINT AUTO_INCREMENT PRIMARY KEY,
      order_no  VARCHAR(64)  NOT NULL,
      ...
      UNIQUE KEY uk_order_no (order_no)
  );
  -- INSERT 时如果 order_no 重复 → Duplicate entry → 判定为重复请求

流程:
  ┌────────┐         ┌──────────┐         ┌──────────┐    ┌────────┐
  │  客户端  │         │  服务端    │         │ 去重表    │    │ 业务表  │
  └───┬────┘         └────┬─────┘         └────┬─────┘    └───┬────┘
      │ 1.POST /pay       │                    │              │
      │  biz_id=PAY123    │                    │              │
      │──────────────────>│                    │              │
      │                   │ 2.INSERT 去重记录   │              │
      │                   │  (status=处理中)    │              │
      │                   │───────────────────>│              │
      │                   │ 3.插入成功(首次请求) │              │
      │                   │<───────────────────│              │
      │                   │ 4.执行支付业务       │              │
      │                   │──────────────────────────────────>│
      │                   │ 5.更新去重记录       │              │
      │                   │  (status=成功,      │              │
      │                   │   response=结果)    │              │
      │                   │───────────────────>│              │
      │ 6.返回成功         │                    │              │
      │<──────────────────│                    │              │
      │                   │                    │              │
      │ 7.重复请求         │                    │              │
      │  biz_id=PAY123    │                    │              │
      │──────────────────>│                    │              │
      │                   │ 8.INSERT 去重记录   │              │
      │                   │───────────────────>│              │
      │                   │ 9.Duplicate entry! │              │
      │                   │<───────────────────│              │
      │                   │ 10.查询首次结果      │              │
      │                   │───────────────────>│              │
      │                   │ 11.返回response     │              │
      │                   │<───────────────────│              │
      │ 12.返回首次结果     │                    │              │
      │<──────────────────│                    │              │
```

#### Go 实现

```go
package idempotent

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrDuplicateRequest = errors.New("重复请求")
	ErrProcessing       = errors.New("请求正在处理中，请稍后")
)

const (
	StatusProcessing = 1
	StatusSuccess    = 2
	StatusFailed     = 3
)

// DedupeTable 去重表幂等实现
type DedupeTable struct {
	db *sql.DB
}

// IdempotentRecord 去重记录
type IdempotentRecord struct {
	BizType  string
	BizID    string
	Response string
	Status   int
}

func NewDedupeTable(db *sql.DB) *DedupeTable {
	return &DedupeTable{db: db}
}

// Execute 幂等执行：传入业务唯一标识 + 业务逻辑，自动去重
func (d *DedupeTable) Execute(ctx context.Context, bizType, bizID string, bizFunc func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	// 1. 尝试插入去重记录（status=处理中）
	inserted, err := d.tryInsert(ctx, bizType, bizID)
	if err != nil {
		return nil, fmt.Errorf("去重表操作失败: %w", err)
	}

	if !inserted {
		// 2. 重复请求 → 查询首次处理结果
		record, err := d.getRecord(ctx, bizType, bizID)
		if err != nil {
			return nil, err
		}
		switch record.Status {
		case StatusProcessing:
			return nil, ErrProcessing // 上一次还在处理中
		case StatusSuccess:
			// 返回首次成功结果
			var result interface{}
			json.Unmarshal([]byte(record.Response), &result)
			return result, ErrDuplicateRequest
		case StatusFailed:
			// 上次失败了，允许重试 → 删除旧记录后重新执行
			d.deleteRecord(ctx, bizType, bizID)
			return d.Execute(ctx, bizType, bizID, bizFunc) // 递归重试
		}
	}

	// 3. 首次请求 → 执行业务逻辑
	result, bizErr := bizFunc(ctx)

	// 4. 更新去重记录
	if bizErr != nil {
		d.updateRecord(ctx, bizType, bizID, StatusFailed, "")
		return nil, bizErr
	}

	respJSON, _ := json.Marshal(result)
	d.updateRecord(ctx, bizType, bizID, StatusSuccess, string(respJSON))
	return result, nil
}

func (d *DedupeTable) tryInsert(ctx context.Context, bizType, bizID string) (bool, error) {
	_, err := d.db.ExecContext(ctx,
		"INSERT INTO idempotent_record (biz_type, biz_id, status) VALUES (?, ?, ?)",
		bizType, bizID, StatusProcessing,
	)
	if err != nil {
		// MySQL duplicate entry error code: 1062
		if isDuplicateKeyError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DedupeTable) getRecord(ctx context.Context, bizType, bizID string) (*IdempotentRecord, error) {
	row := d.db.QueryRowContext(ctx,
		"SELECT biz_type, biz_id, COALESCE(response,''), status FROM idempotent_record WHERE biz_type=? AND biz_id=?",
		bizType, bizID,
	)
	var r IdempotentRecord
	if err := row.Scan(&r.BizType, &r.BizID, &r.Response, &r.Status); err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *DedupeTable) updateRecord(ctx context.Context, bizType, bizID string, status int, response string) {
	d.db.ExecContext(ctx,
		"UPDATE idempotent_record SET status=?, response=?, updated_at=NOW() WHERE biz_type=? AND biz_id=?",
		status, response, bizType, bizID,
	)
}

func (d *DedupeTable) deleteRecord(ctx context.Context, bizType, bizID string) {
	d.db.ExecContext(ctx,
		"DELETE FROM idempotent_record WHERE biz_type=? AND biz_id=?",
		bizType, bizID,
	)
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (contains(err.Error(), "Duplicate entry") || contains(err.Error(), "1062"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchSubstring(s, sub)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// HashRequest 对请求参数做摘要（可选，用于同一 bizID 但参数不同时区分）
func HashRequest(data interface{}) string {
	b, _ := json.Marshal(data)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
```

**使用示例：**

```go
// 支付接口
func (s *PayService) Pay(ctx context.Context, req *PayRequest) (*PayResult, error) {
    // bizID 由客户端生成（如订单号），保证同一笔支付用同一个 bizID
    result, err := s.dedup.Execute(ctx, "pay", req.OrderNo, func(ctx context.Context) (interface{}, error) {
        // 实际支付逻辑（只会执行一次）
        return s.doPayment(ctx, req)
    })

    if errors.Is(err, ErrDuplicateRequest) {
        // 重复请求，返回首次结果（对客户端来说无感知）
        return result.(*PayResult), nil
    }
    if err != nil {
        return nil, err
    }
    return result.(*PayResult), nil
}
```

---

### 6. 方案三：数据库唯一索引（最简单）

```
适用: INSERT 场景，利用 DB 的 UNIQUE 约束天然防重

示例: 创建订单

  -- 订单表
  CREATE TABLE orders (
      id        BIGINT AUTO_INCREMENT PRIMARY KEY,
      order_no  VARCHAR(64) NOT NULL,          -- 客户端生成的唯一订单号
      user_id   BIGINT NOT NULL,
      amount    DECIMAL(10,2) NOT NULL,
      UNIQUE KEY uk_order_no (order_no)        -- 唯一索引保证幂等
  );

Go 代码:
  func CreateOrder(ctx context.Context, db *sql.DB, order *Order) error {
      _, err := db.ExecContext(ctx,
          "INSERT INTO orders (order_no, user_id, amount) VALUES (?, ?, ?)",
          order.OrderNo, order.UserID, order.Amount,
      )
      if err != nil {
          if isDuplicateKeyError(err) {
              return nil  // 重复请求，视为成功（幂等）
          }
          return err
      }
      return nil
  }

优点: 零额外开发成本，DB 层兜底
缺点: 只适用于 INSERT，对 UPDATE/复合操作无效
```

---

### 7. 方案四：乐观锁 / 版本号（更新场景）

```
适用: UPDATE 场景，通过版本号或条件判断防止重复更新

原理:
  UPDATE xxx SET amount=amount-100, version=version+1
  WHERE id=1 AND version=5
  → 影响行数=0 说明版本已变（被其他请求更新过）→ 幂等

场景示例: 扣减库存

  ┌───────────────────────────────────────────────────────────┐
  │ 请求A: UPDATE stock SET qty=qty-1, version=6              │
  │        WHERE sku_id='SKU001' AND version=5                │
  │ → 影响行数=1 → 扣减成功                                    │
  │                                                           │
  │ 请求A重试: UPDATE stock SET qty=qty-1, version=6           │
  │           WHERE sku_id='SKU001' AND version=5             │
  │ → 影响行数=0 → version 已变为6 → 跳过（幂等）              │
  └───────────────────────────────────────────────────────────┘
```

#### Go 实现

```go
// DeductStock 乐观锁扣库存（幂等）
func DeductStock(ctx context.Context, db *sql.DB, skuID string, qty int, version int64) (bool, error) {
	result, err := db.ExecContext(ctx,
		`UPDATE stock
		 SET qty = qty - ?, version = version + 1, updated_at = NOW()
		 WHERE sku_id = ? AND version = ? AND qty >= ?`,
		qty, skuID, version, qty,
	)
	if err != nil {
		return false, err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		// version 不匹配（已被更新）或库存不足
		return false, nil
	}
	return true, nil // 扣减成功
}

// DeductBalance 带请求ID的幂等扣款
func DeductBalance(ctx context.Context, db *sql.DB, userID int64, amount float64, requestID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 去重检查（用流水表的唯一索引）
	_, err = tx.ExecContext(ctx,
		"INSERT INTO balance_log (request_id, user_id, amount) VALUES (?, ?, ?)",
		requestID, userID, -amount,
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil // 重复请求，幂等返回
		}
		return err
	}

	// 2. 乐观锁扣款
	result, err := tx.ExecContext(ctx,
		"UPDATE account SET balance = balance - ? WHERE user_id = ? AND balance >= ?",
		amount, userID, amount,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("余额不足")
	}

	return tx.Commit()
}
```

---

### 8. 方案五：状态机（业务流转幂等）

```
适用: 有明确状态流转的业务（订单、工单、审批等）

原理:
  状态只能单向流转，不允许回退
  UPDATE 时加上当前状态条件 → 状态已变则不会再次执行

  订单状态机:
    待支付(0) ──→ 已支付(1) ──→ 已发货(2) ──→ 已完成(3)
        │                                        │
        └──────────→ 已取消(4) ←─────────────────┘

  SQL:
    UPDATE orders SET status = 1  -- 已支付
    WHERE order_no = 'ORD123'
    AND status = 0                -- 必须是"待支付"状态
    → 第二次执行时 status 已经是1 → WHERE 不命中 → 影响行数=0 → 幂等
```

#### Go 实现

```go
// OrderStatus 订单状态
type OrderStatus int

const (
	OrderPending   OrderStatus = 0 // 待支付
	OrderPaid      OrderStatus = 1 // 已支付
	OrderShipped   OrderStatus = 2 // 已发货
	OrderCompleted OrderStatus = 3 // 已完成
	OrderCancelled OrderStatus = 4 // 已取消
)

// 合法的状态转换
var validTransitions = map[OrderStatus][]OrderStatus{
	OrderPending:   {OrderPaid, OrderCancelled},
	OrderPaid:      {OrderShipped, OrderCancelled},
	OrderShipped:   {OrderCompleted},
	OrderCompleted: {},
	OrderCancelled: {},
}

// TransitOrderStatus 幂等的状态流转
func TransitOrderStatus(ctx context.Context, db *sql.DB, orderNo string, from, to OrderStatus) error {
	// 校验状态转换合法性
	allowed := validTransitions[from]
	valid := false
	for _, s := range allowed {
		if s == to {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("非法状态转换: %d → %d", from, to)
	}

	result, err := db.ExecContext(ctx,
		"UPDATE orders SET status = ?, updated_at = NOW() WHERE order_no = ? AND status = ?",
		to, orderNo, from,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		// 已经被转换过了（幂等）或状态不对
		current := getCurrentStatus(ctx, db, orderNo)
		if current == to {
			return nil // 已经是目标状态 → 幂等成功
		}
		return fmt.Errorf("状态转换失败: 当前状态=%d, 期望=%d", current, from)
	}
	return nil
}

func getCurrentStatus(ctx context.Context, db *sql.DB, orderNo string) OrderStatus {
	var status OrderStatus
	db.QueryRowContext(ctx, "SELECT status FROM orders WHERE order_no=?", orderNo).Scan(&status)
	return status
}
```

---

### 9. 方案六：Redis SETNX（高并发短时去重）

```
适用: 高并发场景，短时间内防重（如秒级内的重复点击）

原理:
  SETNX key value EX ttl
  → 第一次设置成功 → 允许执行
  → 第二次设置失败（key 已存在）→ 拒绝

优点: 性能极高（Redis 单机 10w+ QPS）
缺点: Redis 宕机可能丢失去重记录（不如 DB 可靠）
```

#### Go 实现

```go
// RedisIdempotent Redis SETNX 去重
type RedisIdempotent struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisIdempotent(rdb *redis.Client, ttl time.Duration) *RedisIdempotent {
	return &RedisIdempotent{rdb: rdb, ttl: ttl}
}

// TryAcquire 尝试获取幂等锁
// 返回 true 表示首次请求，可以执行；false 表示重复请求
func (r *RedisIdempotent) TryAcquire(ctx context.Context, bizType, bizID string) (bool, error) {
	key := fmt.Sprintf("idempotent:%s:%s", bizType, bizID)
	ok, err := r.rdb.SetNX(ctx, key, time.Now().Unix(), r.ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Release 主动释放（业务失败时需要释放，允许重试）
func (r *RedisIdempotent) Release(ctx context.Context, bizType, bizID string) error {
	key := fmt.Sprintf("idempotent:%s:%s", bizType, bizID)
	return r.rdb.Del(ctx, key).Err()
}

// ExecuteWithIdempotent 封装幂等执行
func (r *RedisIdempotent) ExecuteWithIdempotent(
	ctx context.Context,
	bizType, bizID string,
	fn func(ctx context.Context) error,
) error {
	acquired, err := r.TryAcquire(ctx, bizType, bizID)
	if err != nil {
		return fmt.Errorf("幂等检查失败: %w", err)
	}
	if !acquired {
		return ErrDuplicateRequest
	}

	// 执行业务
	if err := fn(ctx); err != nil {
		// 业务失败 → 释放幂等锁，允许重试
		r.Release(ctx, bizType, bizID)
		return err
	}
	return nil
}
```

---

### 10. 不同场景的方案选择

```
┌───────────────────────┬────────────────────┬────────────────────────┐
│ 场景                   │ 推荐方案            │ 说明                    │
├───────────────────────┼────────────────────┼────────────────────────┤
│ 表单提交 / 创建订单     │ Token 机制          │ 通用性强，前后端配合      │
│ 支付 / 扣款            │ 去重表 + 唯一索引    │ 可靠性最高，需要留痕      │
│ 库存扣减               │ 乐观锁(版本号)       │ 高并发下性能好            │
│ 订单状态变更            │ 状态机              │ 业务语义清晰              │
│ 短信/通知发送           │ Redis SETNX        │ 高性能短时去重            │
│ MQ 消费端              │ 去重表 / Redis       │ 配合 msgID 去重          │
│ 接口防抖(前端)          │ Token + 按钮禁用    │ 前端为主，后端兜底        │
│ 高并发 + 高可靠         │ Redis SETNX + 去重表│ 两级防护                 │
└───────────────────────┴────────────────────┴────────────────────────┘
```

### 11. 接口幂等 vs MQ 消费幂等

```
┌──────────────────┬──────────────────────┬──────────────────────────┐
│ 维度              │ 接口幂等              │ MQ 消费幂等               │
├──────────────────┼──────────────────────┼──────────────────────────┤
│ 重复来源          │ 用户重复点击/网关重试  │ Broker 重复投递/手动重试   │
│ 唯一标识          │ 业务ID/Token/RequestID│ MessageID / 业务唯一键     │
│ 控制方            │ 调用方(客户端)+服务端  │ 消费方(服务端)             │
│ 时间窗口          │ 短(秒~分钟级)         │ 长(可能跨小时/天)          │
│ 返回值要求        │ 需要返回首次结果       │ 幂等即可，无需返回值       │
│ 典型方案          │ Token/去重表/乐观锁   │ 去重表/Redis/唯一索引      │
│ 异常处理          │ 返回友好错误提示       │ 返回 ACK 避免重复投递      │
└──────────────────┴──────────────────────┴──────────────────────────┘

关键区别:
  接口幂等 → 调用方（客户端/网关）与服务端共同保证，需要返回一致的结果
  MQ 幂等  → 纯消费端保证，重点是"不重复执行业务"，不关心返回值
```

---

### 12. 生产环境最佳实践

```
1. 幂等键的生成策略:
   - 支付: 订单号 + 支付方式 → "PAY:ORD20240101001:WECHAT"
   - 发券: 活动ID + 用户ID   → "COUPON:ACT001:UID10086"
   - 扣库存: SKU + 订单号    → "STOCK:SKU001:ORD20240101001"
   原则: 同一笔业务操作 → 同一个幂等键

2. 幂等窗口期:
   - 太短(如 5s): 可能网络延迟导致重试被误判为新请求
   - 太长(如 7天): 占用大量存储空间
   - 推荐: Token 机制 10 分钟，去重表保留 7 天后定期清理

3. 异常处理:
   - 业务执行失败 → 必须释放幂等锁/删除去重记录，允许重试
   - 幂等组件故障 → 降级为直接执行（宁可重复也不阻塞）

4. 分布式事务中的幂等:
   - 本地事务: 去重记录和业务操作在同一个 DB 事务中
   - 跨库事务: 先写去重记录，再执行业务，最后更新去重状态

5. 多级防护:
   ┌────────┐    ┌────────┐    ┌────────┐    ┌────────┐
   │  前端   │───>│ 网关    │───>│ 服务端  │───>│ DB层   │
   │按钮禁用 │    │限流+去重│    │Token/  │    │唯一索引│
   │防抖节流 │    │Redis   │    │去重表   │    │乐观锁  │
   └────────┘    └────────┘    └────────┘    └────────┘
     第1层           第2层          第3层         第4层
```

### 一句话总结

**接口幂等六种方案：Token 机制（预申请一次性令牌，DEL 原子判断，适合表单提交）、去重表（唯一键记录已处理请求+保存首次结果，适合支付/扣款等高可靠场景）、数据库唯一索引（INSERT 时 DB 层 UNIQUE 约束天然防重，最简单）、乐观锁/版本号（UPDATE WHERE version=x，适合库存扣减等更新场景）、状态机（状态单向流转+UPDATE WHERE status=当前状态，适合订单等有状态业务）、Redis SETNX（高性能短时去重，适合通知/短信防重）。** 生产推荐多级防护：前端按钮防抖→网关 Redis 去重→服务端 Token/去重表→DB 唯一索引兜底。核心原则：同一笔业务操作生成同一个幂等键，业务失败必须释放幂等锁允许重试。

---

## 十三、文件/对象存储系统设计

> 如何设计大文件上传？断点续传？分片上传？秒传？对象存储架构？

### 1. 问题背景

```
大文件上传的挑战:
  ┌─────────────────────────────────────────────────────────────┐
  │ 1. 文件太大 → 单次 HTTP 请求超时（Nginx 默认 60s）           │
  │ 2. 网络不稳定 → 上传到 80% 断了，从头来？                     │
  │ 3. 服务端内存 → 2GB 文件一次性读入内存？OOM                   │
  │ 4. 重复上传 → 同一个文件上传 100 次，存 100 份？              │
  │ 5. 并发上传 → 多用户同时上传，带宽/磁盘 IO 打满               │
  │ 6. 安全问题 → 恶意文件、文件篡改、越权访问                    │
  └─────────────────────────────────────────────────────────────┘

解决方案:
  - 分片上传:   大文件切成小块，逐块上传
  - 断点续传:   断线后从上次位置继续
  - 秒传:       文件指纹相同 → 直接返回已有地址
  - 并发上传:   多个分片并行上传，加快速度
```

### 2. 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端                                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐   │
│  │ 文件切片  │  │ MD5/SHA  │  │ 并发上传  │  │ 断点续传管理   │   │
│  │ (Slice)  │  │ 计算指纹  │  │ (Workers)│  │ (LocalStorage)│   │
│  └──────────┘  └──────────┘  └──────────┘  └───────────────┘   │
└──────────┬──────────┬──────────┬──────────────┬────────────────┘
           │          │          │              │
           ▼          ▼          ▼              ▼
┌──────────────────────────────────────────────────────────────────┐
│                       API 网关 (Nginx/Kong)                      │
│   - 限制单个请求体大小(如 10MB/分片)                               │
│   - 限流(每用户上传并发数)                                        │
└──────────────────────────┬───────────────────────────────────────┘
                           │
┌──────────────────────────▼───────────────────────────────────────┐
│                     上传服务 (Upload Service)                     │
│  ┌───────────────┐  ┌───────────────┐  ┌──────────────────────┐ │
│  │ 秒传检查       │  │ 分片接收      │  │ 分片合并             │ │
│  │ (checkExist)  │  │ (uploadChunk) │  │ (mergeChunks)       │ │
│  └───────┬───────┘  └───────┬───────┘  └──────────┬──────────┘ │
└──────────┼──────────────────┼─────────────────────┼─────────────┘
           │                  │                     │
     ┌─────▼─────┐    ┌──────▼──────┐       ┌──────▼──────┐
     │  MySQL     │    │ 临时存储     │       │ 对象存储     │
     │ 文件元数据  │    │ (本地磁盘/   │       │ (S3/OSS/    │
     │ 分片状态    │    │  MinIO临时)  │       │  MinIO)     │
     └───────────┘    └─────────────┘       └─────────────┘

核心 API:
  POST   /api/upload/check     → 秒传检查 + 获取已上传分片列表
  POST   /api/upload/chunk     → 上传单个分片
  POST   /api/upload/merge     → 通知服务端合并分片
  GET    /api/file/:fileId     → 获取文件（下载/预览）
```

### 3. 秒传原理

```
秒传 = 文件内容指纹（MD5/SHA256）相同 → 直接复用已有文件

流程:
  ┌────────┐                     ┌──────────┐          ┌────────┐
  │  客户端  │                     │ 上传服务   │          │  MySQL  │
  └───┬────┘                     └────┬─────┘          └───┬────┘
      │ 1.前端计算文件 MD5              │                   │
      │   (用 SparkMD5 增量计算)        │                   │
      │                               │                   │
      │ 2.POST /check                  │                   │
      │   {md5, filename, size}        │                   │
      │──────────────────────────────>│                   │
      │                               │ 3.查询文件指纹      │
      │                               │──────────────────>│
      │                               │                   │
      │  ┌─── 情况A: 文件已存在 ────────│                   │
      │  │   4a.返回 {exists:true,     │                   │
      │  │         url:"https://..."}  │                   │
      │<─┘                            │                   │
      │  (秒传成功！无需上传)           │                   │
      │                               │                   │
      │  ┌─── 情况B: 文件不存在 ────────│                   │
      │  │   4b.返回 {exists:false,    │                   │
      │  │     uploadId:"xxx",         │                   │
      │  │     uploadedChunks:[1,3,5]} │  ← 已上传的分片列表│
      │<─┘                            │    （支持断点续传）  │
      │  (需要上传，跳过已有分片)        │                   │

文件指纹方案:
  │ 方案         │ 速度    │ 碰撞率      │ 推荐         │
  │ MD5          │ 快      │ 理论可碰撞   │ 一般场景      │
  │ SHA-256      │ 中      │ 几乎不可能   │ 安全要求高    │
  │ MD5 + 文件大小│ 快      │ 极低        │ 推荐(性价比)  │

注意: 前端计算大文件 MD5 会阻塞 UI → 用 Web Worker + 分块增量计算
```

### 4. 分片上传 + 断点续传原理

```
分片上传:
  将大文件按固定大小（如 5MB）切成多个分片，逐个上传

断点续传:
  上传中断后，服务端已记录哪些分片已收到
  客户端重新上传时，跳过已有分片，只传缺失的

分片大小选择:
  │ 分片大小  │ 优点                  │ 缺点                     │
  │ 1MB      │ 失败重传代价小          │ 分片数多，HTTP 开销大      │
  │ 5MB      │ 平衡（推荐）           │ -                        │
  │ 10MB     │ 请求数少               │ 单片上传时间长，超时风险大  │
  │ 动态调整  │ 根据网速自适应          │ 实现复杂                  │

完整流程:
  ┌────────────────────────────────────────────────────────────┐
  │ 文件: video.mp4 (100MB)     分片大小: 5MB     共 20 个分片  │
  └────────────────────────────────────────────────────────────┘

  第一次上传（上传到第 8 片时断网）:
  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
  │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │ 9 │10 │...│19 │20 │
  │ ✅│ ✅│ ✅│ ✅│ ✅│ ✅│ ✅│ ✅│ ❌│ ❌│ ❌│ ❌│ ❌│
  └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘

  断点续传（重新进入页面）:
  1. POST /check → 返回 uploadedChunks: [1,2,3,4,5,6,7,8]
  2. 客户端只上传 [9,10,...,20] 共 12 个分片
  3. 全部上传完成 → POST /merge → 服务端合并分片
```

### 5. 数据库设计

```sql
-- 文件元数据表
CREATE TABLE file_meta (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    file_md5    VARCHAR(64)  NOT NULL COMMENT '文件内容指纹',
    file_name   VARCHAR(256) NOT NULL COMMENT '原始文件名',
    file_size   BIGINT       NOT NULL COMMENT '文件大小(字节)',
    file_ext    VARCHAR(16)  NOT NULL COMMENT '文件扩展名',
    storage_key VARCHAR(256) NOT NULL COMMENT '对象存储 key',
    storage_url VARCHAR(512) NOT NULL COMMENT '访问 URL',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1正常 2删除',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_md5 (file_md5)
) ENGINE=InnoDB COMMENT='文件元数据';

-- 分片上传任务表
CREATE TABLE upload_task (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    upload_id    VARCHAR(64)  NOT NULL COMMENT '上传任务ID',
    file_md5     VARCHAR(64)  NOT NULL COMMENT '文件指纹',
    file_name    VARCHAR(256) NOT NULL COMMENT '文件名',
    file_size    BIGINT       NOT NULL COMMENT '文件总大小',
    chunk_size   INT          NOT NULL COMMENT '分片大小',
    chunk_count  INT          NOT NULL COMMENT '总分片数',
    status       TINYINT      NOT NULL DEFAULT 0 COMMENT '0上传中 1已完成 2已过期',
    user_id      BIGINT       NOT NULL COMMENT '上传用户',
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_at    DATETIME     NOT NULL COMMENT '任务过期时间',
    UNIQUE KEY uk_upload_id (upload_id),
    INDEX idx_user_status (user_id, status)
) ENGINE=InnoDB COMMENT='分片上传任务';

-- 分片记录表
CREATE TABLE upload_chunk (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    upload_id   VARCHAR(64) NOT NULL COMMENT '上传任务ID',
    chunk_index INT         NOT NULL COMMENT '分片序号(从1开始)',
    chunk_size  INT         NOT NULL COMMENT '分片实际大小',
    chunk_md5   VARCHAR(64) NOT NULL COMMENT '分片MD5',
    storage_key VARCHAR(256) NOT NULL COMMENT '分片临时存储key',
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_upload_chunk (upload_id, chunk_index)
) ENGINE=InnoDB COMMENT='分片记录';
```

### 6. Go 服务端完整实现

```go
package upload

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultChunkSize = 5 * 1024 * 1024 // 5MB
	MaxFileSize      = 5 * 1024 * 1024 * 1024 // 5GB
	UploadExpireTime = 24 * time.Hour
)

const (
	TaskUploading = 0
	TaskDone      = 1
	TaskExpired   = 2
)

// UploadService 上传服务
type UploadService struct {
	db       *sql.DB
	tempDir  string // 分片临时目录
	storeDir string // 最终存储目录（生产环境替换为对象存储）
}

func NewUploadService(db *sql.DB, tempDir, storeDir string) *UploadService {
	os.MkdirAll(tempDir, 0o755)
	os.MkdirAll(storeDir, 0o755)
	return &UploadService{db: db, tempDir: tempDir, storeDir: storeDir}
}

// ---------- API 结构体 ----------

type CheckRequest struct {
	FileMD5  string `json:"file_md5"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

type CheckResponse struct {
	Exists         bool   `json:"exists"`           // true=秒传成功
	URL            string `json:"url,omitempty"`     // 秒传时返回文件URL
	UploadID       string `json:"upload_id,omitempty"`
	ChunkSize      int    `json:"chunk_size,omitempty"`
	ChunkCount     int    `json:"chunk_count,omitempty"`
	UploadedChunks []int  `json:"uploaded_chunks,omitempty"` // 已上传分片列表
}

type MergeRequest struct {
	UploadID string `json:"upload_id"`
	FileMD5  string `json:"file_md5"`
	FileName string `json:"file_name"`
}

type MergeResponse struct {
	URL string `json:"url"`
}

// ---------- 1. 秒传检查 + 初始化上传 ----------

func (s *UploadService) HandleCheck(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "参数错误")
		return
	}

	if req.FileSize > MaxFileSize {
		httpError(w, http.StatusBadRequest, "文件大小超限")
		return
	}

	ctx := r.Context()

	// 1. 秒传检查：文件指纹是否已存在
	var storageURL string
	err := s.db.QueryRowContext(ctx,
		"SELECT storage_url FROM file_meta WHERE file_md5 = ? AND status = 1",
		req.FileMD5,
	).Scan(&storageURL)

	if err == nil {
		// 秒传成功
		jsonResp(w, &CheckResponse{Exists: true, URL: storageURL})
		return
	}

	// 2. 检查是否有未完成的上传任务（支持断点续传）
	var uploadID string
	var chunkSize, chunkCount int
	err = s.db.QueryRowContext(ctx,
		"SELECT upload_id, chunk_size, chunk_count FROM upload_task WHERE file_md5 = ? AND status = ?",
		req.FileMD5, TaskUploading,
	).Scan(&uploadID, &chunkSize, &chunkCount)

	if err == nil {
		// 存在未完成任务 → 查询已上传的分片
		uploaded := s.getUploadedChunks(ctx, uploadID)
		jsonResp(w, &CheckResponse{
			UploadID:       uploadID,
			ChunkSize:      chunkSize,
			ChunkCount:     chunkCount,
			UploadedChunks: uploaded,
		})
		return
	}

	// 3. 新建上传任务
	uploadID = uuid.New().String()
	chunkSize = DefaultChunkSize
	chunkCount = int((req.FileSize + int64(chunkSize) - 1) / int64(chunkSize))

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO upload_task (upload_id, file_md5, file_name, file_size, chunk_size, chunk_count, status, user_id, expire_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		uploadID, req.FileMD5, req.FileName, req.FileSize,
		chunkSize, chunkCount, TaskUploading, getUserID(r),
		time.Now().Add(UploadExpireTime),
	)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "创建上传任务失败")
		return
	}

	// 创建临时分片目录
	os.MkdirAll(filepath.Join(s.tempDir, uploadID), 0o755)

	jsonResp(w, &CheckResponse{
		UploadID:       uploadID,
		ChunkSize:      chunkSize,
		ChunkCount:     chunkCount,
		UploadedChunks: []int{},
	})
}

// ---------- 2. 上传单个分片 ----------

func (s *UploadService) HandleUploadChunk(w http.ResponseWriter, r *http.Request) {
	uploadID := r.FormValue("upload_id")
	chunkIndexStr := r.FormValue("chunk_index")
	chunkMD5 := r.FormValue("chunk_md5")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || uploadID == "" {
		httpError(w, http.StatusBadRequest, "参数错误")
		return
	}

	ctx := r.Context()

	// 1. 验证上传任务存在且未完成
	var taskStatus int
	err = s.db.QueryRowContext(ctx,
		"SELECT status FROM upload_task WHERE upload_id = ?", uploadID,
	).Scan(&taskStatus)
	if err != nil || taskStatus != TaskUploading {
		httpError(w, http.StatusBadRequest, "上传任务不存在或已完成")
		return
	}

	// 2. 检查分片是否已上传（幂等）
	var exists int
	s.db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM upload_chunk WHERE upload_id = ? AND chunk_index = ?",
		uploadID, chunkIndex,
	).Scan(&exists)
	if exists > 0 {
		jsonResp(w, map[string]string{"msg": "分片已存在"})
		return
	}

	// 3. 读取分片数据并保存到临时目录
	file, header, err := r.FormFile("chunk")
	if err != nil {
		httpError(w, http.StatusBadRequest, "读取分片失败")
		return
	}
	defer file.Close()

	// 保存分片到临时文件
	chunkPath := filepath.Join(s.tempDir, uploadID, fmt.Sprintf("chunk_%06d", chunkIndex))
	dst, err := os.Create(chunkPath)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "保存分片失败")
		return
	}

	// 边写入边计算 MD5（验证完整性）
	hash := md5.New()
	written, err := io.Copy(io.MultiWriter(dst, hash), file)
	dst.Close()
	if err != nil {
		os.Remove(chunkPath)
		httpError(w, http.StatusInternalServerError, "写入分片失败")
		return
	}

	// 4. 校验分片 MD5
	actualMD5 := hex.EncodeToString(hash.Sum(nil))
	if chunkMD5 != "" && actualMD5 != chunkMD5 {
		os.Remove(chunkPath)
		httpError(w, http.StatusBadRequest, "分片MD5校验失败")
		return
	}

	// 5. 记录分片信息到数据库
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO upload_chunk (upload_id, chunk_index, chunk_size, chunk_md5, storage_key)
		 VALUES (?, ?, ?, ?, ?)`,
		uploadID, chunkIndex, written, actualMD5, chunkPath,
	)
	if isDuplicateKeyError(err) {
		// 并发上传同一分片，幂等
		os.Remove(chunkPath)
		jsonResp(w, map[string]string{"msg": "分片已存在"})
		return
	}
	if err != nil {
		os.Remove(chunkPath)
		httpError(w, http.StatusInternalServerError, "记录分片失败")
		return
	}

	_ = header // suppress unused
	jsonResp(w, map[string]interface{}{
		"chunk_index": chunkIndex,
		"size":        written,
		"md5":         actualMD5,
	})
}

// ---------- 3. 合并分片 ----------

func (s *UploadService) HandleMerge(w http.ResponseWriter, r *http.Request) {
	var req MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "参数错误")
		return
	}

	ctx := r.Context()

	// 1. 查询上传任务
	var fileSize int64
	var chunkCount int
	err := s.db.QueryRowContext(ctx,
		"SELECT file_size, chunk_count FROM upload_task WHERE upload_id = ? AND status = ?",
		req.UploadID, TaskUploading,
	).Scan(&fileSize, &chunkCount)
	if err != nil {
		httpError(w, http.StatusBadRequest, "上传任务不存在或已完成")
		return
	}

	// 2. 检查所有分片是否已上传
	uploaded := s.getUploadedChunks(ctx, req.UploadID)
	if len(uploaded) != chunkCount {
		httpError(w, http.StatusBadRequest,
			fmt.Sprintf("分片不完整: 已上传 %d/%d", len(uploaded), chunkCount))
		return
	}

	// 3. 合并分片
	ext := filepath.Ext(req.FileName)
	storageKey := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), req.UploadID, ext)
	finalPath := filepath.Join(s.storeDir, storageKey)
	os.MkdirAll(filepath.Dir(finalPath), 0o755)

	if err := s.mergeChunks(req.UploadID, chunkCount, finalPath); err != nil {
		httpError(w, http.StatusInternalServerError, "合并分片失败: "+err.Error())
		return
	}

	// 4. 校验合并后文件 MD5
	finalMD5, err := calcFileMD5(finalPath)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "计算文件MD5失败")
		return
	}
	if finalMD5 != req.FileMD5 {
		os.Remove(finalPath)
		httpError(w, http.StatusBadRequest, "文件MD5不匹配，合并失败")
		return
	}

	// 5. 写入文件元数据
	storageURL := "/files/" + storageKey
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO file_meta (file_md5, file_name, file_size, file_ext, storage_key, storage_url)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE storage_url = VALUES(storage_url)`,
		req.FileMD5, req.FileName, fileSize, ext, storageKey, storageURL,
	)
	if err != nil {
		httpError(w, http.StatusInternalServerError, "保存元数据失败")
		return
	}

	// 6. 更新任务状态 & 清理临时文件
	s.db.ExecContext(ctx,
		"UPDATE upload_task SET status = ? WHERE upload_id = ?",
		TaskDone, req.UploadID,
	)
	go s.cleanupChunks(req.UploadID)

	jsonResp(w, &MergeResponse{URL: storageURL})
}

// mergeChunks 按序号合并所有分片
func (s *UploadService) mergeChunks(uploadID string, chunkCount int, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	for i := 1; i <= chunkCount; i++ {
		chunkPath := filepath.Join(s.tempDir, uploadID, fmt.Sprintf("chunk_%06d", i))
		src, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("打开分片 %d 失败: %w", i, err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			return fmt.Errorf("合并分片 %d 失败: %w", i, err)
		}
		src.Close()
	}
	return nil
}

// cleanupChunks 清理临时分片文件
func (s *UploadService) cleanupChunks(uploadID string) {
	dir := filepath.Join(s.tempDir, uploadID)
	os.RemoveAll(dir)
	// 清理数据库分片记录
	s.db.Exec("DELETE FROM upload_chunk WHERE upload_id = ?", uploadID)
}

// getUploadedChunks 获取已上传的分片列表
func (s *UploadService) getUploadedChunks(ctx context.Context, uploadID string) []int {
	rows, err := s.db.QueryContext(ctx,
		"SELECT chunk_index FROM upload_chunk WHERE upload_id = ? ORDER BY chunk_index",
		uploadID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var chunks []int
	for rows.Next() {
		var idx int
		rows.Scan(&idx)
		chunks = append(chunks, idx)
	}
	return chunks
}

// ---------- 辅助函数 ----------

func calcFileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getUserID(r *http.Request) int64 {
	// 从 context / JWT 中获取用户 ID
	return 0
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	for i := 0; i <= len(s)-len("Duplicate entry"); i++ {
		if s[i:i+len("Duplicate entry")] == "Duplicate entry" {
			return true
		}
	}
	return false
}

// 确保 sort 包被使用
var _ = sort.Ints
```

### 7. 前端分片上传核心逻辑（JavaScript）

```javascript
class ChunkUploader {
    constructor(file, options = {}) {
        this.file = file;
        this.chunkSize = options.chunkSize || 5 * 1024 * 1024; // 5MB
        this.concurrency = options.concurrency || 3;            // 并发数
        this.uploadId = null;
        this.uploadedChunks = new Set();
    }

    // 1. 计算文件 MD5（Web Worker + 分块增量计算）
    async calcMD5() {
        const spark = new SparkMD5.ArrayBuffer();
        const chunks = Math.ceil(this.file.size / this.chunkSize);
        for (let i = 0; i < chunks; i++) {
            const start = i * this.chunkSize;
            const end = Math.min(start + this.chunkSize, this.file.size);
            const buf = await this.file.slice(start, end).arrayBuffer();
            spark.append(buf);
        }
        return spark.end();
    }

    // 2. 秒传检查 + 断点续传
    async check(md5) {
        const resp = await fetch('/api/upload/check', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                file_md5: md5,
                file_name: this.file.name,
                file_size: this.file.size,
            }),
        });
        return resp.json();
    }

    // 3. 并发上传分片
    async uploadChunks(chunkCount) {
        const pending = [];
        for (let i = 1; i <= chunkCount; i++) {
            if (this.uploadedChunks.has(i)) continue; // 跳过已上传分片
            pending.push(i);
        }

        // 并发控制：最多 concurrency 个同时上传
        const pool = [];
        for (const idx of pending) {
            const p = this.uploadOneChunk(idx);
            pool.push(p);
            if (pool.length >= this.concurrency) {
                await Promise.race(pool);
                pool.splice(pool.findIndex(p => p.settled), 1);
            }
        }
        await Promise.all(pool);
    }

    async uploadOneChunk(chunkIndex) {
        const start = (chunkIndex - 1) * this.chunkSize;
        const end = Math.min(start + this.chunkSize, this.file.size);
        const blob = this.file.slice(start, end);

        const formData = new FormData();
        formData.append('upload_id', this.uploadId);
        formData.append('chunk_index', chunkIndex);
        formData.append('chunk', blob);

        const resp = await fetch('/api/upload/chunk', {
            method: 'POST',
            body: formData,
        });
        if (!resp.ok) throw new Error(`分片 ${chunkIndex} 上传失败`);
        this.uploadedChunks.add(chunkIndex);
    }

    // 4. 主流程
    async upload() {
        // Step1: 计算 MD5
        const md5 = await this.calcMD5();

        // Step2: 秒传检查
        const check = await this.check(md5);
        if (check.exists) return check.url; // 秒传成功

        // Step3: 断点续传 - 记录已上传分片
        this.uploadId = check.upload_id;
        (check.uploaded_chunks || []).forEach(i => this.uploadedChunks.add(i));

        // Step4: 上传缺失分片
        await this.uploadChunks(check.chunk_count);

        // Step5: 通知合并
        const mergeResp = await fetch('/api/upload/merge', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                upload_id: this.uploadId,
                file_md5: md5,
                file_name: this.file.name,
            }),
        });
        const result = await mergeResp.json();
        return result.url;
    }
}

// 使用示例
// const uploader = new ChunkUploader(fileInput.files[0]);
// const url = await uploader.upload();
```

### 8. 对象存储架构（生产环境）

```
生产环境不会把文件存本地磁盘，而是使用对象存储:

┌─────────────────────────────────────────────────────────────┐
│                    对象存储选型对比                            │
├──────────────┬──────────────┬──────────────┬───────────────┤
│ 方案          │ AWS S3       │ 阿里云 OSS    │ MinIO(自建)   │
├──────────────┼──────────────┼──────────────┼───────────────┤
│ 部署方式      │ 云服务        │ 云服务        │ 自部署         │
│ 成本          │ 按量付费      │ 按量付费      │ 硬件+运维      │
│ S3 兼容       │ 原生          │ 兼容          │ 完全兼容       │
│ 可用性        │ 99.99%       │ 99.99%       │ 取决于部署      │
│ 分片上传      │ ✅ 原生支持   │ ✅ 原生支持   │ ✅ 原生支持    │
│ 适用场景      │ 海外业务      │ 国内业务      │ 私有化/合规    │
└──────────────┴──────────────┴──────────────┴───────────────┘

S3 分片上传 API（各对象存储通用）:
  1. CreateMultipartUpload  → 初始化分片上传，获取 UploadID
  2. UploadPart             → 上传单个分片，返回 ETag
  3. CompleteMultipartUpload→ 提交所有分片 ETag，触发合并
  4. AbortMultipartUpload   → 取消上传，清理分片
  5. ListParts              → 查询已上传分片（断点续传）

客户端直传（推荐，减轻服务端带宽压力）:

  ┌────────┐              ┌──────────┐              ┌──────────┐
  │  客户端  │              │  服务端    │              │  OSS/S3  │
  └───┬────┘              └────┬─────┘              └────┬─────┘
      │ 1.请求上传凭证          │                        │
      │────────────────────────>│                        │
      │                        │ 2.生成预签名URL          │
      │                        │  (STS临时凭证/          │
      │                        │   Presigned URL)        │
      │ 3.返回上传凭证           │                        │
      │<────────────────────────│                        │
      │                        │                        │
      │ 4.直传文件到OSS ────────────────────────────────>│
      │  (不经过服务端!)        │                        │
      │                        │                        │
      │ 5.上传完成              │                        │
      │<──────────────────────────────────────────────── │
      │ 6.通知服务端            │                        │
      │────────────────────────>│                        │
      │                        │ 7.记录元数据             │

  优点: 文件流量不经过业务服务器，节省带宽和 CPU
  安全: 使用临时凭证(STS)或预签名URL，有效期短
```

#### Go 预签名 URL 示例（S3 兼容）

```go
package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Uploader S3 兼容的上传服务
type S3Uploader struct {
	client       *s3.Client
	presignClient *s3.PresignClient
	bucket       string
}

// GeneratePresignedURL 生成预签名上传 URL（客户端直传）
func (u *S3Uploader) GeneratePresignedURL(ctx context.Context, key, contentType string) (string, error) {
	presigned, err := u.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}
	return presigned.URL, nil
}

// InitMultipartUpload 初始化分片上传（大文件）
func (u *S3Uploader) InitMultipartUpload(ctx context.Context, key, contentType string) (string, error) {
	output, err := u.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return *output.UploadId, nil
}

// GeneratePartPresignedURL 生成单个分片的预签名URL
func (u *S3Uploader) GeneratePartPresignedURL(ctx context.Context, key, uploadID string, partNumber int32) (string, error) {
	presigned, err := u.presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(u.bucket),
		Key:        aws.String(key),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(partNumber),
	}, s3.WithPresignExpires(30*time.Minute))
	if err != nil {
		return "", err
	}
	return presigned.URL, nil
}

// CompleteMultipartUpload 完成分片上传
func (u *S3Uploader) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []CompletedPart) error {
	s3Parts := make([]types.CompletedPart, len(parts))
	for i, p := range parts {
		s3Parts[i] = types.CompletedPart{
			PartNumber: aws.Int32(p.PartNumber),
			ETag:       aws.String(p.ETag),
		}
	}
	_, err := u.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(u.bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: s3Parts,
		},
	})
	return err
}

type CompletedPart struct {
	PartNumber int32
	ETag       string
}
```

### 9. 安全与运维

```
安全策略:
  ┌───────────────────────┬──────────────────────────────────────┐
  │ 风险                   │ 解决方案                              │
  ├───────────────────────┼──────────────────────────────────────┤
  │ 恶意文件上传(木马/病毒) │ 文件类型白名单 + 内容嗅探(magic bytes)│
  │ 超大文件攻击           │ 限制文件大小 + 分片大小 + 总分片数     │
  │ 文件篡改               │ 分片 MD5 校验 + 合并后整体 MD5 校验   │
  │ 越权访问               │ 预签名 URL(短有效期) + 权限校验       │
  │ 并发上传资源耗尽        │ 每用户上传并发限制 + 任务过期清理      │
  │ 临时文件堆积           │ 定时清理过期分片 + 任务超时自动过期     │
  └───────────────────────┴──────────────────────────────────────┘

文件类型校验（magic bytes）:
  JPEG: FF D8 FF
  PNG:  89 50 4E 47
  PDF:  25 50 44 46
  ZIP:  50 4B 03 04
  → 不能只看扩展名，要读取文件头几个字节验证真实类型

运维清理:
  -- 定时清理过期上传任务（每小时执行）
  UPDATE upload_task SET status = 2
  WHERE status = 0 AND expire_at < NOW();

  -- 清理孤立分片文件（每天执行）
  DELETE FROM upload_chunk
  WHERE upload_id IN (
      SELECT upload_id FROM upload_task WHERE status = 2
  );
```

### 一句话总结

**大文件上传三板斧：分片上传（文件按 5MB 切片逐个上传，Nginx 限制单请求体大小即可）、断点续传（服务端记录已上传分片列表，客户端重传时 /check 接口返回缺失分片只传增量）、秒传（前端 SparkMD5 增量计算文件指纹，服务端查询 file_md5 唯一索引，命中直接返回 URL）。** 生产环境推荐客户端直传对象存储（S3/OSS/MinIO）：服务端生成预签名 URL→客户端直传 OSS→上传完成回调服务端记录元数据，文件流量不经过业务服务器。安全要点：文件类型白名单+magic bytes 校验、分片 MD5+合并后整体 MD5 双重校验、预签名 URL 短有效期、每用户上传并发限制+过期任务定时清理。

---

## 十四、搜索系统设计

> 站内搜索怎么设计？ES 倒排索引原理？搜索建议怎么实现？数据同步怎么做？

### 1. 搜索系统面临的问题

```
为什么不用 MySQL LIKE 做搜索？

  SELECT * FROM articles WHERE title LIKE '%分布式锁%';

  问题:
  ┌──────────────────────────────────────────────────────────────┐
  │ 1. 全表扫描:  LIKE '%xxx%' 无法使用索引，百万级数据秒级响应？ │
  │ 2. 无法分词:  搜"分布式 锁" → LIKE 无法拆词匹配              │
  │ 3. 无相关性排序: 哪条结果最匹配？LIKE 没有打分机制             │
  │ 4. 不支持同义词: 搜"手机" 找不到"智能电话"                    │
  │ 5. 不支持拼音/纠错: 搜"fenbudishi" 或 "分步式" 都找不到       │
  │ 6. 性能差:    并发搜索时 MySQL CPU 打满                       │
  └──────────────────────────────────────────────────────────────┘

  解决方案 → Elasticsearch（基于 Lucene 的分布式搜索引擎）
```

### 2. 整体架构

```
┌────────────────────────────────────────────────────────────────────┐
│                           客户端                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────────┐ │
│  │ 搜索框输入     │  │ 搜索建议      │  │ 搜索结果页(高亮+分页+筛选)│ │
│  └──────┬───────┘  └──────┬───────┘  └───────────┬───────────────┘ │
└─────────┼─────────────────┼──────────────────────┼─────────────────┘
          │                 │                      │
          ▼                 ▼                      ▼
┌────────────────────────────────────────────────────────────────────┐
│                        搜索服务 (Search Service)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ 搜索建议      │  │ 全文搜索      │  │ 搜索日志 & 热词统计      │ │
│  │ (Suggest)    │  │ (Search)     │  │ (Analytics)              │ │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────────┘ │
└─────────┼─────────────────┼──────────────────────┼─────────────────┘
          │                 │                      │
          ▼                 ▼                      ▼
┌─────────────────┐ ┌──────────────┐  ┌────────────────────────────┐
│  Redis           │ │Elasticsearch │  │ Kafka / ClickHouse         │
│  - 热词缓存      │ │  - 倒排索引  │  │  - 搜索日志                 │
│  - 搜索建议缓存  │ │  - 分片+副本 │  │  - 热词分析                 │
│  - 搜索结果缓存  │ │  - 聚合统计  │  │  - 搜索质量分析             │
└─────────────────┘ └──────┬───────┘  └────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │    数据同步 (Data Sync)   │
              │  ┌──────┐  ┌──────────┐ │
              │  │Binlog│  │ MQ 异步   │ │
              │  │ 监听  │  │ 双写      │ │
              │  └──────┘  └──────────┘ │
              └────────────┬────────────┘
                           │
                    ┌──────▼──────┐
                    │    MySQL     │
                    │  (主数据源)   │
                    └─────────────┘
```

### 3. Elasticsearch 倒排索引原理

```
正排索引（Forward Index）— MySQL 的 B+ 树索引:
  文档ID → 文档内容
  ┌──────┬──────────────────────────────────────┐
  │ DocID│ Content                               │
  ├──────┼──────────────────────────────────────┤
  │ 1    │ "Elasticsearch 是分布式搜索引擎"       │
  │ 2    │ "Redis 是分布式缓存"                   │
  │ 3    │ "Kafka 是分布式消息队列"                │
  └──────┴──────────────────────────────────────┘
  搜索"分布式" → 遍历每条记录看是否包含 → O(N) 全表扫描


倒排索引（Inverted Index）— ES 的核心:
  词项(Term) → 包含该词的文档列表
  ┌──────────────┬──────────────────────────────────┐
  │ Term          │ Posting List（倒排列表）           │
  ├──────────────┼──────────────────────────────────┤
  │ elasticsearch │ [DocID:1, pos:0]                  │
  │ 分布式        │ [DocID:1, pos:1], [DocID:2, pos:1], [DocID:3, pos:1] │
  │ 搜索引擎      │ [DocID:1, pos:2]                  │
  │ redis         │ [DocID:2, pos:0]                  │
  │ 缓存          │ [DocID:2, pos:2]                  │
  │ kafka         │ [DocID:3, pos:0]                  │
  │ 消息队列      │ [DocID:3, pos:2]                  │
  └──────────────┴──────────────────────────────────┘
  搜索"分布式" → 直接查 Term 字典 → 返回 [Doc1, Doc2, Doc3] → O(1) 查找

倒排索引的三层结构:
  ┌─────────────────────────────────────────────────────────┐
  │ 1. Term Dictionary（词项字典）                           │
  │    所有去重的 Term 集合，按字典序排列                      │
  │    使用 FST (Finite State Transducer) 压缩存储           │
  │    ↓                                                    │
  │ 2. Term Index（词项索引）                                │
  │    Term Dictionary 的索引（类似 B+ 树的非叶节点）          │
  │    前缀树 (Trie) 结构，常驻内存                           │
  │    快速定位 Term 在 Dictionary 中的位置                   │
  │    ↓                                                    │
  │ 3. Posting List（倒排列表）                              │
  │    每个 Term 对应的文档 ID 列表                           │
  │    包含: DocID, 词频(TF), 位置(Position), 偏移(Offset)   │
  │    使用多种压缩: FOR (Frame of Reference)、Roaring Bitmap│
  └─────────────────────────────────────────────────────────┘

  搜索"分布式 缓存" 的过程:
    ① "分布式" → Term Index(内存) → Term Dictionary(磁盘) → Posting: [1,2,3]
    ② "缓存"   → Term Index(内存) → Term Dictionary(磁盘) → Posting: [2]
    ③ AND 交集: [1,2,3] ∩ [2] = [2]  → 返回 Doc2
    ④ OR  并集: [1,2,3] ∪ [2] = [1,2,3] → 返回 Doc1,2,3
```

#### 分词过程（Analyzer）

```
原始文本: "Elasticsearch是一个分布式全文搜索引擎"

分词流水线:
  ┌──────────────────────────────────────────────────────┐
  │ 1. Character Filters（字符过滤器）                     │
  │    - HTML Strip: 去除 <b>xxx</b> 标签                 │
  │    - Mapping: 把 "①" 映射为 "1"                       │
  │    → "Elasticsearch是一个分布式全文搜索引擎"             │
  │                                                      │
  │ 2. Tokenizer（分词器）                                 │
  │    - standard: 按空格/标点分词（英文）                   │
  │    - ik_max_word: 中文最细粒度分词                      │
  │    - ik_smart: 中文最粗粒度分词                         │
  │    → ["elasticsearch", "是", "一个", "分布式",          │
  │       "全文", "搜索", "引擎", "搜索引擎", "全文搜索"]    │
  │                                                      │
  │ 3. Token Filters（词项过滤器）                          │
  │    - lowercase: 转小写                                 │
  │    - stop: 去停用词("是","一个","的")                   │
  │    - synonym: 同义词扩展("手机"→"手机,智能电话")         │
  │    → ["elasticsearch", "分布式", "全文", "搜索",         │
  │       "引擎", "搜索引擎", "全文搜索"]                    │
  └──────────────────────────────────────────────────────┘

中文分词器选型:
  │ 分词器        │ 特点                     │ 适用场景         │
  │ ik_smart     │ 最粗粒度（少切）           │ 搜索查询时使用    │
  │ ik_max_word  │ 最细粒度（多切）           │ 索引写入时使用    │
  │ jieba        │ 支持自定义词典+词频调整     │ NLP/精细分词     │
  │ HanLP        │ NLP 全套（分词/NER/情感）  │ 高级文本分析      │

  最佳实践:
    写入时用 ik_max_word（切更多词 → 召回率高）
    搜索时用 ik_smart  （切更少词 → 精确率高）
```

### 4. ES 核心概念与 MySQL 类比

```
┌──────────────────┬──────────────────┬──────────────────────────┐
│ MySQL             │ Elasticsearch    │ 说明                     │
├──────────────────┼──────────────────┼──────────────────────────┤
│ Database          │ Index            │ 索引（类似数据库）         │
│ Table             │ Type(已废弃→Index)│ ES7+ 一个 Index=一个表   │
│ Row               │ Document         │ 一条文档（JSON 格式）     │
│ Column            │ Field            │ 字段                     │
│ Schema            │ Mapping          │ 字段类型定义              │
│ SQL               │ Query DSL        │ JSON 查询语法             │
│ Primary Key       │ _id              │ 文档唯一标识              │
│ Index(B+ Tree)    │ Inverted Index   │ 倒排索引                 │
│ SELECT            │ GET /_search     │ 搜索                     │
│ INSERT            │ PUT /index/_doc  │ 写入文档                  │
│ Partition          │ Shard           │ 分片（水平拆分）          │
│ Replica            │ Replica         │ 副本（高可用）            │
└──────────────────┴──────────────────┴──────────────────────────┘

ES 集群架构:
  ┌───────────────────────────────────────────────────────────┐
  │ ES Cluster (3 节点)                                       │
  │                                                           │
  │  Node1 (Master)     Node2              Node3              │
  │  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐    │
  │  │ Shard-0 (P) │   │ Shard-1 (P) │   │ Shard-2 (P) │    │
  │  │ Shard-1 (R) │   │ Shard-2 (R) │   │ Shard-0 (R) │    │
  │  └─────────────┘   └─────────────┘   └─────────────┘    │
  │                                                           │
  │  P = Primary Shard（主分片，负责读写）                      │
  │  R = Replica Shard（副本分片，负责读+故障切换）             │
  │  同一分片的 P 和 R 不在同一节点（保证可用性）               │
  └───────────────────────────────────────────────────────────┘
```

### 5. 搜索服务 Go 实现

```go
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// SearchService 搜索服务
type SearchService struct {
	es    *elasticsearch.Client
	index string
}

func NewSearchService(es *elasticsearch.Client, index string) *SearchService {
	return &SearchService{es: es, index: index}
}

// ---------- 索引管理 ----------

// CreateIndex 创建索引（定义 Mapping + Settings）
func (s *SearchService) CreateIndex(ctx context.Context) error {
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					// 自定义分析器：写入时最细粒度 + 同义词 + 拼音
					"my_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "ik_max_word",
						"filter":    []string{"lowercase", "my_synonym"},
					},
					// 搜索时最粗粒度
					"my_search_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "ik_smart",
						"filter":    []string{"lowercase", "my_synonym"},
					},
				},
				"filter": map[string]interface{}{
					"my_synonym": map[string]interface{}{
						"type":          "synonym",
						"synonyms_path": "analysis/synonyms.txt", // 手机,智能电话,mobile
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":            "text",
					"analyzer":        "my_analyzer",        // 写入分词器
					"search_analyzer": "my_search_analyzer",  // 搜索分词器
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword", // 精确匹配/聚合用
							"ignore_above": 256,
						},
					},
				},
				"content": map[string]interface{}{
					"type":            "text",
					"analyzer":        "my_analyzer",
					"search_analyzer": "my_search_analyzer",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"tags": map[string]interface{}{
					"type": "keyword",
				},
				"status": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd HH:mm:ss||epoch_millis",
				},
				"suggest": map[string]interface{}{
					"type": "completion", // 搜索建议专用字段
				},
			},
		},
	}

	body, _ := json.Marshal(mapping)
	resp, err := s.es.Indices.Create(
		s.index,
		s.es.Indices.Create.WithContext(ctx),
		s.es.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建索引失败: %s", string(b))
	}
	return nil
}

// ---------- 文档写入 ----------

// ArticleDoc 文章文档
type ArticleDoc struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	Status    int      `json:"status"`
	CreatedAt string   `json:"created_at"`
	Suggest   []string `json:"suggest"` // 搜索建议词
}

// IndexDoc 写入/更新文档
func (s *SearchService) IndexDoc(ctx context.Context, doc *ArticleDoc) error {
	body, _ := json.Marshal(doc)
	resp, err := s.es.Index(
		s.index,
		bytes.NewReader(body),
		s.es.Index.WithContext(ctx),
		s.es.Index.WithDocumentID(fmt.Sprintf("%d", doc.ID)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("写入文档失败: %s", string(b))
	}
	return nil
}

// BulkIndex 批量写入（高性能）
func (s *SearchService) BulkIndex(ctx context.Context, docs []*ArticleDoc) error {
	var buf bytes.Buffer
	for _, doc := range docs {
		// action 行
		meta := fmt.Sprintf(`{"index":{"_index":"%s","_id":"%d"}}`, s.index, doc.ID)
		buf.WriteString(meta + "\n")
		// 文档行
		body, _ := json.Marshal(doc)
		buf.Write(body)
		buf.WriteByte('\n')
	}

	resp, err := s.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		s.es.Bulk.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("批量写入失败: %s", string(b))
	}
	return nil
}

// DeleteDoc 删除文档
func (s *SearchService) DeleteDoc(ctx context.Context, docID int64) error {
	resp, err := s.es.Delete(
		s.index,
		fmt.Sprintf("%d", docID),
		s.es.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ---------- 搜索 ----------

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword  string   `json:"keyword"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	SortBy   string   `json:"sort_by"` // "relevance" / "time"
}

// SearchResult 搜索结果
type SearchResult struct {
	Total int64         `json:"total"`
	Items []*SearchItem `json:"items"`
	Took  int64         `json:"took"` // 耗时ms
}

type SearchItem struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	Score     float64  `json:"score"`
	Highlight []string `json:"highlight"` // 高亮片段
}

// Search 全文搜索（支持高亮、过滤、排序、分页）
func (s *SearchService) Search(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建 Bool Query
	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":    req.Keyword,
				"fields":   []string{"title^3", "content"}, // title 权重 3 倍
				"type":     "best_fields",
				"analyzer": "ik_smart",
			},
		},
	}

	filter := []map[string]interface{}{
		{"term": map[string]interface{}{"status": 1}}, // 只搜已发布
	}

	if req.Category != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"category": req.Category},
		})
	}
	if len(req.Tags) > 0 {
		filter = append(filter, map[string]interface{}{
			"terms": map[string]interface{}{"tags": req.Tags},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
			"fields": map[string]interface{}{
				"title":   map[string]interface{}{},
				"content": map[string]interface{}{"fragment_size": 150, "number_of_fragments": 3},
			},
		},
		"from": (req.Page - 1) * req.PageSize,
		"size": req.PageSize,
	}

	// 排序
	if req.SortBy == "time" {
		query["sort"] = []map[string]interface{}{
			{"created_at": map[string]interface{}{"order": "desc"}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	}

	body, _ := json.Marshal(query)
	resp, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.index),
		s.es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var esResp ESSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&esResp); err != nil {
		return nil, err
	}

	result := &SearchResult{
		Total: esResp.Hits.Total.Value,
		Took:  esResp.Took,
	}

	for _, hit := range esResp.Hits.Hits {
		item := &SearchItem{
			ID:       hit.Source.ID,
			Title:    hit.Source.Title,
			Content:  hit.Source.Content,
			Category: hit.Source.Category,
			Tags:     hit.Source.Tags,
			Score:    hit.Score,
		}
		// 高亮替换
		if hl, ok := hit.Highlight["title"]; ok && len(hl) > 0 {
			item.Title = hl[0]
		}
		if hl, ok := hit.Highlight["content"]; ok {
			item.Highlight = hl
		}
		result.Items = append(result.Items, item)
	}

	return result, nil
}

// ES 响应结构体
type ESSearchResponse struct {
	Took int64 `json:"took"`
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Score     float64            `json:"_score"`
			Source    ArticleDoc         `json:"_source"`
			Highlight map[string][]string `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}
```

### 6. 搜索建议（Suggest / 自动补全）

```
搜索建议的三种实现:

  ┌──────────────────┬───────────────────────────────────────────┐
  │ 方式              │ 说明                                      │
  ├──────────────────┼───────────────────────────────────────────┤
  │ Completion       │ ES 内置 Completion Suggester               │
  │ Suggester        │ 基于 FST 前缀匹配，性能极高（常驻内存）     │
  │                  │ 适合: 搜索框自动补全                        │
  ├──────────────────┼───────────────────────────────────────────┤
  │ 热词 + 前缀搜索  │ Redis Sorted Set 存热搜词                  │
  │                  │ ZRANGEBYLEX 前缀匹配                       │
  │                  │ 适合: 热搜词推荐                            │
  ├──────────────────┼───────────────────────────────────────────┤
  │ Search-as-you-   │ 用 match_phrase_prefix 或 edge_ngram      │
  │ type             │ 实时搜索并返回文档标题                      │
  │                  │ 适合: 搜索结果即时展示                      │
  └──────────────────┴───────────────────────────────────────────┘

用户输入 "分布" 时:
  ┌──────────────────────────────┐
  │ 🔍 分布                      │
  ├──────────────────────────────┤
  │  分布式系统设计               │  ← Completion Suggest
  │  分布式锁实现                 │
  │  分布式事务                   │
  │  分布式缓存                   │
  │ ─────────────────────────── │
  │  🔥 分布式 (热搜)            │  ← 热词推荐
  └──────────────────────────────┘
```

#### Go 实现搜索建议

```go
// ---------- 搜索建议 ----------

// Suggest 搜索建议（ES Completion Suggester）
func (s *SearchService) Suggest(ctx context.Context, prefix string, size int) ([]string, error) {
	if size <= 0 {
		size = 10
	}

	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"title-suggest": map[string]interface{}{
				"prefix": prefix,
				"completion": map[string]interface{}{
					"field":           "suggest",
					"size":            size,
					"skip_duplicates": true,
					"fuzzy": map[string]interface{}{
						"fuzziness": "AUTO", // 支持容错（打错1-2个字也能匹配）
					},
				},
			},
		},
		"_source": false, // 不需要返回文档内容
	}

	body, _ := json.Marshal(query)
	resp, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.index),
		s.es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var esResp struct {
		Suggest map[string][]struct {
			Options []struct {
				Text string `json:"text"`
			} `json:"options"`
		} `json:"suggest"`
	}
	json.NewDecoder(resp.Body).Decode(&esResp)

	var suggestions []string
	if sg, ok := esResp.Suggest["title-suggest"]; ok && len(sg) > 0 {
		for _, opt := range sg[0].Options {
			suggestions = append(suggestions, opt.Text)
		}
	}
	return suggestions, nil
}

// ---------- 热词搜索建议（Redis 实现） ----------

// HotSuggest 基于 Redis 的热词搜索建议
type HotSuggest struct {
	rdb redis.Client // 假设已引入 redis 包
}

// RecordSearch 记录搜索词（每次搜索调用）
func (h *HotSuggest) RecordSearch(ctx context.Context, keyword string) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return
	}
	// ZINCRBY hot:search:daily 1 "分布式锁"
	h.rdb.ZIncrBy(ctx, "hot:search:daily", 1, keyword)
	// 设置过期时间（每天轮转）
	h.rdb.ExpireNX(ctx, "hot:search:daily", 24*time.Hour)
}

// GetHotWords 获取热搜词 Top N
func (h *HotSuggest) GetHotWords(ctx context.Context, topN int64) ([]string, error) {
	// ZREVRANGE hot:search:daily 0 N-1
	return h.rdb.ZRevRange(ctx, "hot:search:daily", 0, topN-1).Result()
}

// PrefixSuggest 前缀匹配热词
// 利用 Redis Sorted Set 的字典序范围查询
func (h *HotSuggest) PrefixSuggest(ctx context.Context, prefix string, limit int64) ([]string, error) {
	// 方案: 维护一个专门用于前缀搜索的 Sorted Set
	// 插入时把词的所有前缀都加入: "分" "分布" "分布式" "分布式锁"
	// 查询时用 ZRANGEBYLEX
	min := "[" + prefix
	max := "[" + prefix + "\xff"
	return h.rdb.ZRangeByLex(ctx, "suggest:prefix", &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Count:  limit,
	}).Result()
}
```

### 7. MySQL → ES 数据同步方案

```
核心问题: MySQL 是主数据源，ES 是搜索副本，如何保持一致？

方案对比:
  ┌─────────────────┬────────────────────────────────────────────┐
  │ 方案             │ 说明                                       │
  ├─────────────────┼────────────────────────────────────────────┤
  │ 同步双写         │ 业务代码同时写 MySQL + ES                   │
  │                 │ ❌ 强耦合、ES 写失败影响主流程、无法保证一致  │
  ├─────────────────┼────────────────────────────────────────────┤
  │ MQ 异步双写      │ 写 MySQL 后发 MQ → 消费者写 ES              │
  │ (推荐)          │ ✅ 解耦、异步、失败可重试                    │
  │                 │ ⚠️ 有短暂延迟（秒级）                       │
  ├─────────────────┼────────────────────────────────────────────┤
  │ Binlog 监听      │ Canal/Debezium 监听 MySQL Binlog → 写 ES   │
  │ (推荐)          │ ✅ 完全解耦、对业务零侵入                    │
  │                 │ ⚠️ 运维 Canal 组件、Binlog 需开启 ROW 格式   │
  ├─────────────────┼────────────────────────────────────────────┤
  │ 定时全量同步     │ 定时任务全量扫 MySQL → 写 ES                │
  │                 │ ✅ 实现简单，兜底方案                        │
  │                 │ ❌ 延迟大、数据量大时慢                      │
  └─────────────────┴────────────────────────────────────────────┘

推荐组合: Binlog 监听（实时增量） + 定时全量（兜底校正）

MQ 异步双写流程:
  ┌────────┐      ┌────────┐      ┌────────┐      ┌────────┐
  │ 业务服务 │──TX──│ MySQL  │      │ Kafka  │      │ 同步服务 │
  └───┬────┘      └────────┘      └───┬────┘      └───┬────┘
      │ 1.写入MySQL(事务)              │               │
      │ 2.发消息到Kafka                │               │
      │───────────────────────────────>│               │
      │                               │ 3.消费消息      │
      │                               │──────────────>│
      │                               │               │ 4.写入ES
      │                               │               │──→ ES

Binlog 监听流程:
  ┌────────┐      ┌────────┐      ┌────────┐      ┌────────┐
  │ 业务服务 │──TX──│ MySQL  │      │ Canal  │      │  ES    │
  └────────┘      └───┬────┘      └───┬────┘      └───┬────┘
                      │ Binlog        │               │
                      │──────────────>│               │
                      │               │ 解析+转换      │
                      │               │──────────────>│
                      │               │  写入ES文档    │
```

#### Go MQ 异步同步实现

```go
// SyncConsumer 消费 MQ 消息，同步数据到 ES
type SyncConsumer struct {
	searchSvc *SearchService
}

type SyncMessage struct {
	Action string      `json:"action"` // "insert" / "update" / "delete"
	Table  string      `json:"table"`
	Data   ArticleDoc  `json:"data"`
}

func (c *SyncConsumer) HandleMessage(ctx context.Context, msg []byte) error {
	var syncMsg SyncMessage
	if err := json.Unmarshal(msg, &syncMsg); err != nil {
		return fmt.Errorf("解析同步消息失败: %w", err)
	}

	switch syncMsg.Action {
	case "insert", "update":
		// 构建 suggest 字段（用于自动补全）
		syncMsg.Data.Suggest = buildSuggestInput(syncMsg.Data.Title, syncMsg.Data.Tags)
		return c.searchSvc.IndexDoc(ctx, &syncMsg.Data)
	case "delete":
		return c.searchSvc.DeleteDoc(ctx, syncMsg.Data.ID)
	default:
		return fmt.Errorf("未知操作: %s", syncMsg.Action)
	}
}

// buildSuggestInput 构建搜索建议输入
func buildSuggestInput(title string, tags []string) []string {
	suggestions := []string{title}
	suggestions = append(suggestions, tags...)
	return suggestions
}

// 业务层发送同步消息
func PublishSyncMessage(producer KafkaProducer, action string, doc *ArticleDoc) error {
	msg := &SyncMessage{
		Action: action,
		Table:  "articles",
		Data:   *doc,
	}
	data, _ := json.Marshal(msg)
	return producer.Send("topic-search-sync", data)
}

type KafkaProducer interface {
	Send(topic string, data []byte) error
}
```

### 8. 搜索相关性打分（TF-IDF / BM25）

```
ES 默认使用 BM25 算法对搜索结果打分:

TF-IDF（ES 5.x 之前默认）:
  score = TF × IDF

  TF (Term Frequency):  词在当前文档中出现的次数（越多越相关）
  IDF (Inverse Document Frequency): log(总文档数 / 包含该词的文档数)
                                    （越少文档包含 → IDF 越高 → 越重要）

BM25（ES 5.x+ 默认，TF-IDF 的改进版）:
  对 TF 做了饱和处理（词出现 10 次和 100 次差距不大）
  考虑了文档长度（长文档不应因为词多而分高）

  score = IDF × (TF × (k1 + 1)) / (TF + k1 × (1 - b + b × dl/avgdl))

  k1 = 1.2（TF 饱和参数）
  b  = 0.75（文档长度归一化参数）
  dl = 当前文档长度
  avgdl = 平均文档长度

自定义打分（Function Score Query）:
  场景: 新文章加权、热门文章加权、付费内容置顶

  {
    "query": {
      "function_score": {
        "query": { "match": { "title": "分布式" } },
        "functions": [
          {
            "gauss": {                         // 时间衰减：越新分越高
              "created_at": {
                "origin": "now",
                "scale":  "7d",
                "decay":  0.5
              }
            }
          },
          {
            "field_value_factor": {            // 热度加权
              "field":    "view_count",
              "modifier": "log1p",
              "factor":   0.1
            }
          }
        ],
        "score_mode": "multiply",
        "boost_mode": "multiply"
      }
    }
  }
```

### 9. 搜索优化最佳实践

```
性能优化:
  ┌───────────────────────┬──────────────────────────────────────┐
  │ 优化点                 │ 具体措施                              │
  ├───────────────────────┼──────────────────────────────────────┤
  │ 索引设计               │ 合理设置分片数（数据量/30GB≈分片数）   │
  │                       │ 不需要搜索的字段设 index:false         │
  │                       │ 精确匹配用 keyword，全文搜索用 text    │
  ├───────────────────────┼──────────────────────────────────────┤
  │ 查询优化               │ filter 替代 query（filter 可缓存）    │
  │                       │ 避免 wildcard 开头通配（如 *abc）      │
  │                       │ 深分页用 search_after 替代 from+size  │
  ├───────────────────────┼──────────────────────────────────────┤
  │ 写入优化               │ Bulk API 批量写入（每批 5-15MB）      │
  │                       │ 写入密集期加大 refresh_interval(30s)  │
  │                       │ 关闭副本写入再打开（全量重建时）        │
  ├───────────────────────┼──────────────────────────────────────┤
  │ 缓存策略               │ 热门搜索词 → Redis 缓存结果（TTL 60s）│
  │                       │ 搜索建议 → Redis 缓存（TTL 5min）     │
  │                       │ ES 自带 Query Cache + Request Cache  │
  └───────────────────────┴──────────────────────────────────────┘

深分页问题:
  from + size 方式:
    from=10000, size=20 → ES 需要在每个分片上取 10020 条，再合并排序
    分片越多越慢，ES 默认限制 from+size ≤ 10000

  解决方案:
    ① search_after（推荐）: 用上一页最后一条的排序值作为起点
       "search_after": [1618888888, "doc_id_123"]
    ② scroll API: 快照方式遍历，适合导出数据（不适合实时搜索）
    ③ 业务限制: 只允许翻到前 100 页（Google/百度都这么做）

搜索质量优化:
  ┌─────────────────┬──────────────────────────────────────┐
  │ 问题              │ 解决方案                              │
  ├─────────────────┼──────────────────────────────────────┤
  │ 搜"手机"漏掉      │ 同义词词典: 手机 → 手机,智能电话,mobile│
  │ "智能电话"        │                                      │
  ├─────────────────┼──────────────────────────────────────┤
  │ 搜"苹果手机"结果  │ multi_match + best_fields 策略       │
  │ 不好              │ title 字段加权(boost:3)               │
  ├─────────────────┼──────────────────────────────────────┤
  │ 新内容搜不到      │ 缩短 refresh_interval (默认1s)       │
  │                  │ 或写入后手动 refresh                  │
  ├─────────────────┼──────────────────────────────────────┤
  │ 停用词干扰       │ 自定义停用词表(的/了/在/是)            │
  ├─────────────────┼──────────────────────────────────────┤
  │ 拼音搜索         │ pinyin 分析器插件(elasticsearch-      │
  │                  │ analysis-pinyin)                     │
  └─────────────────┴──────────────────────────────────────┘
```

### 一句话总结

**站内搜索核心架构：MySQL（主数据源）→ Binlog/MQ 异步同步 → Elasticsearch（倒排索引+全文搜索）→ Redis（热词缓存+搜索建议缓存）。** 倒排索引三层结构：Term Index（前缀树，常驻内存快速定位）→ Term Dictionary（FST 压缩存储所有词项）→ Posting List（每个词对应的文档 ID 列表+词频+位置，FOR/Roaring Bitmap 压缩）。写入时用 ik_max_word 最细粒度分词提高召回率，搜索时用 ik_smart 最粗粒度提高精确率。搜索建议三种方案：ES Completion Suggester（FST 前缀匹配，性能最高）、Redis ZSET 热词推荐、edge_ngram 实时搜索。数据同步推荐 Binlog 监听（Canal，零侵入实时增量）+ 定时全量兜底。深分页用 search_after 替代 from+size。打分用 BM25 + Function Score 自定义加权（时间衰减+热度）。

---

## 十五、抢红包系统设计

> 微信红包的拆分算法？高并发怎么扛？如何防超发？如何保证公平性？

### 1. 红包系统核心挑战

```
微信红包核心数据（2024 除夕）:
  - 峰值 QPS: 数十万/秒
  - 单个红包: N 个人同时抢同一个红包
  - 资金安全: 绝对不能超发（发 100 元红包，不能领出 101 元）
  - 公平性:   每个人领到的金额应足够随机，不能总是第一个最大

挑战分解:
  ┌──────────────────────────────────────────────────────────┐
  │ 1. 拆分算法:  100 元拆 10 个，怎么拆才公平？               │
  │ 2. 高并发:    10 个人同时抢同一个红包，不能超发              │
  │ 3. 资金安全:  金额精度（分）、总金额校验、防重领              │
  │ 4. 性能:     发红包→抢红包→拆红包→入账 全链路低延迟         │
  │ 5. 体验:     抢到后立刻看到金额，不能让用户等                │
  └──────────────────────────────────────────────────────────┘
```

### 2. 整体架构

```
┌────────────────────────────────────────────────────────────────┐
│                         客户端                                  │
│   [发红包] ──→ [红包消息] ──→ [抢/拆红包] ──→ [查看明细]        │
└──────┬──────────────────────────┬──────────────────────────────┘
       │                          │
       ▼                          ▼
┌──────────────────────────────────────────────────────────────┐
│                       API 网关                                │
│   限流(用户级) + 防刷(频率检测) + 鉴权                         │
└──────┬──────────────────────────┬────────────────────────────┘
       │                          │
       ▼                          ▼
┌──────────────┐          ┌──────────────────────────────────┐
│ 发红包服务    │          │ 抢红包服务                         │
│              │          │                                  │
│ 1.参数校验    │          │ 1.Redis 原子抢占（LPOP）           │
│ 2.预扣余额    │          │ 2.抢到 → 返回金额                 │
│ 3.拆分金额    │          │ 3.异步入账（MQ）                   │
│ 4.存入Redis   │          │                                  │
│ 5.发消息通知  │          │                                  │
└──────┬───────┘          └───────┬────────────┬─────────────┘
       │                          │            │
       ▼                          ▼            ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────────────┐
│   MySQL       │   │    Redis     │   │    Kafka / MQ        │
│ - 红包记录    │   │ - 红包金额列表│   │ - 异步入账消息        │
│ - 领取明细    │   │ - 已领取集合  │   │                      │
│ - 用户余额    │   │ - 红包状态    │   │                      │
└──────────────┘   └──────────────┘   └──────────┬───────────┘
                                                  │
                                           ┌──────▼──────┐
                                           │  入账服务    │
                                           │ - 更新余额   │
                                           │ - 写流水     │
                                           └─────────────┘

两种架构模式:
  ┌──────────────────┬──────────────────────────────────────────┐
  │ 预拆分（推荐）    │ 发红包时就把金额拆好，存入 Redis List      │
  │                  │ 抢红包 = LPOP 弹出一个金额（原子操作）      │
  │                  │ ✅ 简单、高性能、无锁竞争                   │
  ├──────────────────┼──────────────────────────────────────────┤
  │ 实时拆分          │ 抢红包时实时计算金额（二倍均值法）          │
  │                  │ 需要原子操作保证剩余金额/人数一致           │
  │                  │ ⚠️ 需要 Lua 脚本保证原子性                 │
  └──────────────────┴──────────────────────────────────────────┘
```

### 3. 红包拆分算法

#### 算法一：二倍均值法（微信红包采用）

```
核心公式:
  每次领取金额 = random(1, 剩余金额/剩余人数 × 2)

  假设: 100 元（10000 分）分 5 个人

  第1人: random(1, 10000/5*2) = random(1, 4000) → 假设 2500
  第2人: random(1, 7500/4*2)  = random(1, 3750) → 假设 1200
  第3人: random(1, 6300/3*2)  = random(1, 4200) → 假设 3800
  第4人: random(1, 2500/2*2)  = random(1, 2500) → 假设 800
  第5人: 剩余全部 = 2500 - 800 = 1700

  总额: 2500 + 1200 + 3800 + 800 + 1700 = 10000 ✅

特点:
  ✅ 期望值相等: 每个人的期望金额 = 总金额/总人数（数学可证明）
  ✅ 简单高效: 只需要一次随机数生成
  ✅ 不会出现 0: 下限是 1 分
  ⚠️ 先抢的方差大，后抢的方差小（越往后越趋近平均值）
  ⚠️ 不会出现超大红包（被 2 倍均值限制住了）
```

#### 算法二：线段切割法

```
原理:
  把总金额看成一条线段 [0, totalAmount]
  随机生成 N-1 个切割点，排序后相邻两点之间的距离就是每个人的金额

  100 元 5 人 → 在 [0, 10000] 上随机取 4 个点
  排序后: 0 - 1500 - 3700 - 6200 - 8800 - 10000
  金额:   1500,  2200,  2500,  2600,  1200

特点:
  ✅ 完全公平: 每个位置领到大/小红包的概率相同
  ✅ 可能出现超大红包（更刺激）
  ❌ 需要排序，不适合实时拆分
  ❌ 可能出现 0（需要处理最小金额约束）
```

#### 算法三：洗牌法（预拆分推荐）

```
原理:
  1. 先用二倍均值法拆好所有金包
  2. 对结果做 Fisher-Yates 洗牌
  3. 消除"先抢方差大"的问题

这样每个人拿到任何金额的概率都相等
```

#### Go 实现三种算法

```go
package redpacket

import (
	"errors"
	"math/rand"
	"sort"
)

var (
	ErrInvalidArgs = errors.New("参数无效: 金额或人数不合法")
	ErrMinAmount   = errors.New("金额不足: 无法保证每人至少1分")
)

// ---- 算法一: 二倍均值法 ----

func SplitDoubleAvg(totalAmount int64, count int) ([]int64, error) {
	if totalAmount <= 0 || count <= 0 {
		return nil, ErrInvalidArgs
	}
	if totalAmount < int64(count) {
		return nil, ErrMinAmount // 每人至少 1 分
	}

	result := make([]int64, count)
	remain := totalAmount
	remainCount := count

	for i := 0; i < count-1; i++ {
		avg := remain / int64(remainCount)
		// 上限 = min(2*avg, remain - 剩余人数保底)
		maxVal := avg * 2
		if maxVal > remain-int64(remainCount-1) {
			maxVal = remain - int64(remainCount-1)
		}
		if maxVal < 1 {
			maxVal = 1
		}
		amount := rand.Int63n(maxVal) + 1 // [1, maxVal]
		result[i] = amount
		remain -= amount
		remainCount--
	}
	result[count-1] = remain // 最后一个人拿剩余

	return result, nil
}

// ---- 算法二: 线段切割法 ----

func SplitSegment(totalAmount int64, count int) ([]int64, error) {
	if totalAmount <= 0 || count <= 0 {
		return nil, ErrInvalidArgs
	}
	if totalAmount < int64(count) {
		return nil, ErrMinAmount
	}

	if count == 1 {
		return []int64{totalAmount}, nil
	}

	// 预留每人 1 分的保底
	remain := totalAmount - int64(count) // 可分配的金额
	if remain < 0 {
		remain = 0
	}

	// 生成 N-1 个切割点
	cuts := make([]int64, count-1)
	used := make(map[int64]bool) // 避免重复切割点
	for i := 0; i < count-1; i++ {
		for {
			c := rand.Int63n(remain + 1)
			if !used[c] {
				used[c] = true
				cuts[i] = c
				break
			}
		}
	}

	sort.Slice(cuts, func(i, j int) bool { return cuts[i] < cuts[j] })

	result := make([]int64, count)
	prev := int64(0)
	for i, c := range cuts {
		result[i] = (c - prev) + 1 // +1 是保底的 1 分
		prev = c
	}
	result[count-1] = (remain - prev) + 1

	return result, nil
}

// ---- 算法三: 二倍均值 + Fisher-Yates 洗牌 ----

func SplitShuffle(totalAmount int64, count int) ([]int64, error) {
	result, err := SplitDoubleAvg(totalAmount, count)
	if err != nil {
		return nil, err
	}

	// Fisher-Yates 洗牌（消除顺序偏差）
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}
```

### 4. 发红包流程

```
发红包 = 预扣余额 + 拆分金额 + 存 Redis + 发消息

  ┌────────┐        ┌──────────┐      ┌────────┐     ┌────────┐
  │  客户端  │        │ 发红包服务 │      │ MySQL  │     │ Redis  │
  └───┬────┘        └────┬─────┘      └───┬────┘     └───┬────┘
      │ 1.发红包(100元,5人) │              │              │
      │────────────────────>│              │              │
      │                    │ 2.参数校验     │              │
      │                    │ 3.预扣余额     │              │
      │                    │──────────────>│              │
      │                    │ 4.生成红包ID   │              │
      │                    │ 5.拆分金额     │              │
      │                    │  [25,12,38,8,17] (单位:分)    │
      │                    │ 6.写红包记录    │              │
      │                    │──────────────>│              │
      │                    │ 7.存入Redis    │              │
      │                    │  RPUSH rp:{id} 25 12 38 8 17 │
      │                    │  SET rp:info:{id} {总额,人数} │
      │                    │───────────────────────────── >│
      │                    │ 8.发消息到群聊  │              │
      │ 9.返回成功          │              │              │
      │<────────────────────│              │              │
```

#### Go 实现

```go
package redpacket

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedPacketService 红包服务
type RedPacketService struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewRedPacketService(db *sql.DB, rdb *redis.Client) *RedPacketService {
	return &RedPacketService{db: db, rdb: rdb}
}

// RedPacketInfo 红包信息
type RedPacketInfo struct {
	ID          string `json:"id"`
	SenderID    int64  `json:"sender_id"`
	TotalAmount int64  `json:"total_amount"` // 单位: 分
	TotalCount  int    `json:"total_count"`
	RemainCount int    `json:"remain_count"`
	GroupID     string `json:"group_id"`
	Greeting    string `json:"greeting"` // 祝福语
	ExpireAt    int64  `json:"expire_at"`
}

// ---- 发红包 ----

type SendRequest struct {
	SenderID    int64  `json:"sender_id"`
	GroupID     string `json:"group_id"`
	TotalAmount int64  `json:"total_amount"` // 单位: 分
	TotalCount  int    `json:"total_count"`
	Greeting    string `json:"greeting"`
}

func (s *RedPacketService) Send(ctx context.Context, req *SendRequest) (string, error) {
	// 1. 参数校验
	if req.TotalAmount < int64(req.TotalCount) {
		return "", ErrMinAmount
	}
	if req.TotalCount <= 0 || req.TotalCount > 100 {
		return "", fmt.Errorf("红包个数需在 1-100 之间")
	}
	if req.TotalAmount > 20000*100 { // 单个红包上限 200 元
		return "", fmt.Errorf("单个红包金额不能超过200元")
	}

	// 2. 预扣余额（数据库事务）
	rpID := uuid.New().String()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// 扣减余额
	result, err := tx.ExecContext(ctx,
		"UPDATE user_balance SET balance = balance - ? WHERE user_id = ? AND balance >= ?",
		req.TotalAmount, req.SenderID, req.TotalAmount,
	)
	if err != nil {
		return "", err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return "", fmt.Errorf("余额不足")
	}

	// 写红包记录
	_, err = tx.ExecContext(ctx,
		`INSERT INTO red_packet (id, sender_id, group_id, total_amount, total_count, remain_amount, remain_count, greeting, status, created_at, expire_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, NOW(), ?)`,
		rpID, req.SenderID, req.GroupID, req.TotalAmount, req.TotalCount,
		req.TotalAmount, req.TotalCount, req.Greeting,
		time.Now().Add(24*time.Hour),
	)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	// 3. 拆分金额 + 洗牌
	amounts, err := SplitShuffle(req.TotalAmount, req.TotalCount)
	if err != nil {
		return "", err
	}

	// 4. 存入 Redis（Pipeline 批量操作）
	pipe := s.rdb.Pipeline()

	// 红包金额列表（List：RPUSH 存入，LPOP 领取）
	rpKey := "rp:amounts:" + rpID
	amountStrs := make([]interface{}, len(amounts))
	for i, a := range amounts {
		amountStrs[i] = a
	}
	pipe.RPush(ctx, rpKey, amountStrs...)
	pipe.Expire(ctx, rpKey, 25*time.Hour)

	// 红包信息
	infoKey := "rp:info:" + rpID
	info := &RedPacketInfo{
		ID:          rpID,
		SenderID:    req.SenderID,
		TotalAmount: req.TotalAmount,
		TotalCount:  req.TotalCount,
		RemainCount: req.TotalCount,
		GroupID:     req.GroupID,
		Greeting:    req.Greeting,
		ExpireAt:    time.Now().Add(24 * time.Hour).Unix(),
	}
	infoJSON, _ := json.Marshal(info)
	pipe.Set(ctx, infoKey, infoJSON, 25*time.Hour)

	// 已领取用户集合（用于防重领）
	claimedKey := "rp:claimed:" + rpID
	pipe.Expire(ctx, claimedKey, 25*time.Hour) // 提前设置过期时间

	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("Redis 存储失败: %w", err)
	}

	return rpID, nil
}
```

### 5. 抢红包流程（核心高并发环节）

```
抢红包 = Redis 原子操作（防重领 + LPOP 抢金额） + 异步入账

  ┌────────┐         ┌──────────┐            ┌────────┐
  │  客户端  │         │ 抢红包服务 │            │ Redis  │
  └───┬────┘         └────┬─────┘            └───┬────┘
      │ 1.抢红包(rpID)     │                     │
      │───────────────────>│                     │
      │                   │ 2.Lua脚本原子操作     │
      │                   │  a.检查红包是否存在    │
      │                   │  b.检查是否已领过      │
      │                   │  c.LPOP 弹出金额      │
      │                   │  d.SADD 记录已领取     │
      │                   │─────────────────────>│
      │                   │ 3.返回金额(或失败)     │
      │                   │<─────────────────────│
      │                   │                     │
      │                   │ 4.异步: 发MQ入账消息  │
      │ 5.返回抢到的金额    │                     │
      │<───────────────────│                     │
      │                   │                     │
      │  (用户看到:恭喜抢到 X.XX 元)              │

  为什么用 Lua 脚本？
    a. SISMEMBER（检查是否领过）
    b. LPOP（弹出金额）
    c. SADD（标记已领）
    这三步必须原子执行，否则并发下可能重复领取
```

#### Redis Lua 脚本 + Go 实现

```go
// ---- 抢红包 Lua 脚本 ----

// KEYS[1] = rp:amounts:{id}    红包金额列表
// KEYS[2] = rp:claimed:{id}    已领取用户集合
// KEYS[3] = rp:info:{id}       红包信息
// ARGV[1] = userID
// 返回: >0 抢到的金额(分)  -1 已领过  -2 已抢完  -3 红包不存在
var grabScript = redis.NewScript(`
	-- 1. 检查红包是否存在
	local info = redis.call('GET', KEYS[3])
	if not info then
		return -3
	end

	-- 2. 检查是否已领过
	if redis.call('SISMEMBER', KEYS[2], ARGV[1]) == 1 then
		return -1
	end

	-- 3. 从列表中弹出一个金额
	local amount = redis.call('LPOP', KEYS[1])
	if not amount then
		return -2  -- 已抢完
	end

	-- 4. 记录已领取
	redis.call('SADD', KEYS[2], ARGV[1])

	return tonumber(amount)
`)

// GrabResult 抢红包结果
type GrabResult struct {
	Success bool   `json:"success"`
	Amount  int64  `json:"amount"`   // 分
	Message string `json:"message"`
}

// Grab 抢红包
func (s *RedPacketService) Grab(ctx context.Context, rpID string, userID int64) (*GrabResult, error) {
	keys := []string{
		"rp:amounts:" + rpID,
		"rp:claimed:" + rpID,
		"rp:info:" + rpID,
	}

	result, err := grabScript.Run(ctx, s.rdb, keys, userID).Int64()
	if err != nil {
		return nil, fmt.Errorf("抢红包失败: %w", err)
	}

	switch result {
	case -3:
		return &GrabResult{Message: "红包不存在或已过期"}, nil
	case -1:
		return &GrabResult{Message: "你已经领过了"}, nil
	case -2:
		return &GrabResult{Message: "手慢了，红包已抢完"}, nil
	default:
		// 抢到了！异步入账
		go s.asyncSettle(context.Background(), rpID, userID, result)

		return &GrabResult{
			Success: true,
			Amount:  result,
			Message: fmt.Sprintf("恭喜抢到 %.2f 元", float64(result)/100),
		}, nil
	}
}

// asyncSettle 异步入账（生产环境用 MQ）
func (s *RedPacketService) asyncSettle(ctx context.Context, rpID string, userID, amount int64) {
	// 1. 写领取记录到 MySQL
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO red_packet_detail (red_packet_id, user_id, amount, created_at)
		 VALUES (?, ?, ?, NOW())`,
		rpID, userID, amount,
	)
	if err != nil {
		// 写入失败 → 补偿队列重试
		fmt.Printf("写领取记录失败: rpID=%s userID=%d err=%v\n", rpID, userID, err)
		return
	}

	// 2. 更新红包剩余
	s.db.ExecContext(ctx,
		`UPDATE red_packet SET remain_amount = remain_amount - ?, remain_count = remain_count - 1
		 WHERE id = ? AND remain_amount >= ? AND remain_count > 0`,
		amount, rpID, amount,
	)

	// 3. 入账（加到用户余额）
	s.db.ExecContext(ctx,
		"UPDATE user_balance SET balance = balance + ? WHERE user_id = ?",
		amount, userID,
	)
}
```

### 6. 数据库设计

```sql
-- 红包表
CREATE TABLE red_packet (
    id             VARCHAR(64) PRIMARY KEY COMMENT '红包ID',
    sender_id      BIGINT       NOT NULL COMMENT '发送者',
    group_id       VARCHAR(64)  NOT NULL COMMENT '群ID',
    total_amount   BIGINT       NOT NULL COMMENT '总金额(分)',
    total_count    INT          NOT NULL COMMENT '总个数',
    remain_amount  BIGINT       NOT NULL COMMENT '剩余金额(分)',
    remain_count   INT          NOT NULL COMMENT '剩余个数',
    greeting       VARCHAR(128) NOT NULL DEFAULT '恭喜发财' COMMENT '祝福语',
    status         TINYINT      NOT NULL DEFAULT 1 COMMENT '1进行中 2已领完 3已过期',
    created_at     DATETIME     NOT NULL,
    expire_at      DATETIME     NOT NULL,
    INDEX idx_sender (sender_id),
    INDEX idx_group (group_id, created_at)
) ENGINE=InnoDB COMMENT='红包表';

-- 红包领取明细表
CREATE TABLE red_packet_detail (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    red_packet_id   VARCHAR(64)  NOT NULL COMMENT '红包ID',
    user_id         BIGINT       NOT NULL COMMENT '领取用户',
    amount          BIGINT       NOT NULL COMMENT '领取金额(分)',
    created_at      DATETIME     NOT NULL,
    UNIQUE KEY uk_rp_user (red_packet_id, user_id), -- 防重领
    INDEX idx_user (user_id, created_at)
) ENGINE=InnoDB COMMENT='红包领取明细';

-- 用户余额表
CREATE TABLE user_balance (
    user_id    BIGINT PRIMARY KEY,
    balance    BIGINT NOT NULL DEFAULT 0 COMMENT '余额(分)',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='用户余额';
```

### 7. 高并发设计要点

```
1. 所有热操作放 Redis，MySQL 异步兜底:
   ┌──────────────────────────────────────────────────────────┐
   │ 热路径 (Redis, 微秒级):                                   │
   │   抢红包 = Lua 脚本（检查+弹出+标记，单个 Redis ≈ 10w QPS）│
   │                                                          │
   │ 冷路径 (MySQL, 异步):                                     │
   │   入账 = MQ 异步消费 → 写领取记录 + 更新余额               │
   └──────────────────────────────────────────────────────────┘

2. 单个红包的并发瓶颈分析:
   一个 100 人红包，100 人同时抢 → 竞争同一个 Redis Key
   Redis 单线程处理 Lua 脚本 → 天然串行化 → 无并发问题
   单个 Redis 10w+ QPS → 100 人并发毫无压力

3. 海量红包的并发:
   问题: 春节几亿个红包同时进行
   方案: 按红包 ID 哈希到不同 Redis 分片
   rpID → hash(rpID) % N → Redis 节点 i
   不同红包之间完全无竞争

4. 防刷策略:
   ┌─────────────────────────────────────────────────┐
   │ - 接口级限流: 单用户 QPS ≤ 10                    │
   │ - 红包级防重: Redis SISMEMBER（Lua 脚本中已包含）│
   │ - DB 级兜底: UNIQUE KEY (red_packet_id, user_id) │
   │ - 群成员校验: 非群成员不能抢                      │
   │ - 设备指纹: 同一设备频繁抢红包 → 降级             │
   └─────────────────────────────────────────────────┘

5. 资金安全保障:
   ┌─────────────────────────────────────────────────────────┐
   │ 层级             │ 措施                                  │
   ├─────────────────┼───────────────────────────────────────┤
   │ 拆分时           │ 金额总和 == 红包总额（代码校验）       │
   │ 存 Redis 时      │ RPUSH 个数 == 红包总数                │
   │ 抢红包时         │ LPOP 原子操作，弹完即止               │
   │ 入账时           │ DB 唯一索引防重 + 余额变更走事务       │
   │ 对账             │ 定时任务校验: 已领总额 ≤ 红包总额      │
   └─────────────────┴───────────────────────────────────────┘
```

### 8. 红包过期退款

```
红包 24 小时未领完 → 剩余金额退回发送者

方案: 延时任务（参考 设计题三：延时任务系统）

  定时扫描:
    每分钟扫描 expire_at < NOW() AND status = 1 的红包
    退回剩余金额

  或用 Redis 过期回调:
    KEY 过期时触发退款（需要开启 notify-keyspace-events）

Go 实现:
```

```go
// ExpireRefund 过期退款（定时任务调用）
func (s *RedPacketService) ExpireRefund(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, sender_id, remain_amount FROM red_packet
		 WHERE status = 1 AND expire_at < NOW()
		 LIMIT 100`, // 每次处理 100 个
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var rpID string
		var senderID, remainAmount int64
		rows.Scan(&rpID, &senderID, &remainAmount)

		if remainAmount <= 0 {
			// 已领完，更新状态即可
			s.db.ExecContext(ctx, "UPDATE red_packet SET status = 2 WHERE id = ?", rpID)
			continue
		}

		// 退款事务
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			continue
		}

		// 更新红包状态
		tx.ExecContext(ctx, "UPDATE red_packet SET status = 3 WHERE id = ? AND status = 1", rpID)
		// 退回余额
		tx.ExecContext(ctx, "UPDATE user_balance SET balance = balance + ? WHERE user_id = ?",
			remainAmount, senderID)

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			continue
		}

		// 清理 Redis
		pipe := s.rdb.Pipeline()
		pipe.Del(ctx, "rp:amounts:"+rpID)
		pipe.Del(ctx, "rp:claimed:"+rpID)
		pipe.Del(ctx, "rp:info:"+rpID)
		pipe.Exec(ctx)
	}
	return nil
}
```

### 9. 红包拆分算法公平性验证

```
验证方法: 模拟 100 万次拆分，统计每个位置的期望值

  100 元 10 人，模拟 100 万次:
  ┌──────┬──────────────────┬──────────────────┬──────────────────┐
  │ 位置  │ 二倍均值法        │ 线段切割法        │ 二倍均值+洗牌     │
  ├──────┼──────────────────┼──────────────────┼──────────────────┤
  │ 第1人 │ 期望≈10.0 方差大  │ 期望≈10.0 方差大  │ 期望≈10.0 方差中  │
  │ 第2人 │ 期望≈10.0 方差中  │ 期望≈10.0 方差大  │ 期望≈10.0 方差中  │
  │ ...  │ ...              │ ...              │ ...              │
  │ 第10人│ 期望≈10.0 方差小  │ 期望≈10.0 方差大  │ 期望≈10.0 方差中  │
  └──────┴──────────────────┴──────────────────┴──────────────────┘

  结论:
    - 三种算法期望值都相等（公平）
    - 二倍均值法: 越靠后方差越小（不够刺激但够公平）
    - 线段切割法: 所有位置方差相同（完全随机，可能出现极值）
    - 洗牌法: 消除了位置偏差（推荐）
```

### 一句话总结

**抢红包系统核心设计：发红包时预拆分（二倍均值法+Fisher-Yates 洗牌保证公平性，每次 random(1, 剩余/剩余人数×2)），拆好的金额 RPUSH 存入 Redis List；抢红包用 Lua 脚本原子执行三步（SISMEMBER 防重领→LPOP 弹出金额→SADD 标记已领），单个 Redis 10w+ QPS 轻松扛住单红包并发；不同红包按 ID 哈希分散到不同 Redis 分片互不竞争；入账走 MQ 异步（写领取记录+更新余额），DB 唯一索引(red_packet_id, user_id)兜底防重。** 资金安全四层保障：拆分时校验总额→Redis LPOP 弹完即止→DB 唯一索引防重→定时对账。过期未领完的红包由定时任务退款。

---

## 十六、评论系统设计

> 评论的楼中楼怎么存储？热评排序怎么实现？评论列表怎么分页？大V 评论区高并发怎么扛？

### 1. 评论系统核心需求

```
功能需求:
  ┌──────────────────────────────────────────────────────────┐
  │ 1. 发评论:    对文章/视频/动态发布评论                     │
  │ 2. 回复评论:  对某条评论进行回复（楼中楼）                 │
  │ 3. 点赞:      对评论点赞/取消点赞                         │
  │ 4. 排序:      按时间（最新）/ 按热度（最热）               │
  │ 5. 分页:      一级评论分页 + 子评论"展开更多"              │
  │ 6. 删除:      删除评论（级联删除子评论？）                 │
  │ 7. 举报/审核: 敏感词过滤 + 人工审核                        │
  └──────────────────────────────────────────────────────────┘

技术挑战:
  ┌──────────────────────────────────────────────────────────┐
  │ 1. 楼中楼存储: 评论的树状结构如何在 DB 中高效存储和查询    │
  │ 2. 热评排序:   点赞数实时变化，排序如何更新                │
  │ 3. 高并发读:   大V 文章评论区可能有百万条评论              │
  │ 4. 计数准确:   评论数、点赞数在高并发下如何保证准确         │
  │ 5. 深分页:     评论翻到第 100 页性能如何保证               │
  └──────────────────────────────────────────────────────────┘
```

### 2. 评论的层级模型

```
三种常见的评论结构:

┌──────────────────────────────────────────────────────────────────┐
│ 模型一: 二级评论（微博/B站/抖音 模式，推荐）                       │
│                                                                  │
│  ┌─ 一级评论 A（直接对文章的评论）                                 │
│  │   ├─ 回复 A-1: "@张三 说得好"                                  │
│  │   ├─ 回复 A-2: "@A-1 同意"                                    │
│  │   └─ 回复 A-3: "@李四 反对"                                    │
│  │   └─ [展开更多 12 条回复]                                      │
│  │                                                               │
│  ├─ 一级评论 B                                                    │
│  │   ├─ 回复 B-1                                                  │
│  │   └─ 回复 B-2                                                  │
│  ...                                                              │
│                                                                  │
│  特点:                                                            │
│    - 只有两级: 一级评论(root) + 二级回复(reply)                    │
│    - 所有回复平铺在一级评论下，用 @用户名 表示回复谁                │
│    - 查询简单: 先查一级评论，再按 root_id 查子评论                 │
│    - 适合瀑布流场景                                               │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│ 模型二: 多级嵌套（Reddit/Hacker News 模式）                       │
│                                                                  │
│  评论 A                                                           │
│    ├── 回复 A-1                                                   │
│    │     ├── 回复 A-1-1                                           │
│    │     │     └── 回复 A-1-1-1                                   │
│    │     └── 回复 A-1-2                                           │
│    └── 回复 A-2                                                   │
│                                                                  │
│  特点:                                                            │
│    - 无限嵌套，前端递归渲染                                       │
│    - 存储用 parent_id 自引用                                      │
│    - 查询复杂: 需要递归或闭包表                                   │
│    - 适合长讨论/辩论场景                                          │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│ 模型三: 盖楼（早期贴吧/网易新闻模式）                              │
│                                                                  │
│  #1楼 张三: 好文章！                                              │
│  #2楼 李四: 回复 #1 确实不错                                      │
│  #3楼 王五: 回复 #2 +1                                           │
│                                                                  │
│  特点:                                                            │
│    - 纯平铺，所有评论按楼层编号                                   │
│    - 存储最简单（一张表，自增楼层号）                              │
│    - 不适合复杂讨论                                               │
└──────────────────────────────────────────────────────────────────┘

推荐: 二级评论模型（业界主流，兼顾功能和性能）
```

### 3. 数据库设计（二级评论模型）

```sql
-- 评论主表
CREATE TABLE comment (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    obj_type       TINYINT      NOT NULL COMMENT '对象类型: 1文章 2视频 3动态',
    obj_id         BIGINT       NOT NULL COMMENT '对象ID(文章ID/视频ID)',
    user_id        BIGINT       NOT NULL COMMENT '评论者',
    content        VARCHAR(2000) NOT NULL COMMENT '评论内容',

    -- 层级关系（核心字段）
    root_id        BIGINT       NOT NULL DEFAULT 0 COMMENT '根评论ID(一级评论此字段=0)',
    parent_id      BIGINT       NOT NULL DEFAULT 0 COMMENT '直接父评论ID(一级评论此字段=0)',
    reply_user_id  BIGINT       NOT NULL DEFAULT 0 COMMENT '被回复的用户ID',

    -- 计数（冗余，异步更新）
    like_count     INT          NOT NULL DEFAULT 0 COMMENT '点赞数',
    reply_count    INT          NOT NULL DEFAULT 0 COMMENT '子回复数(仅一级评论)',

    -- 状态
    status         TINYINT      NOT NULL DEFAULT 1 COMMENT '0审核中 1正常 2删除 3折叠',
    ip_location    VARCHAR(32)  NULL COMMENT 'IP属地',
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- 索引
    INDEX idx_obj (obj_type, obj_id, status, created_at),       -- 按时间查一级评论
    INDEX idx_obj_hot (obj_type, obj_id, status, like_count),   -- 按热度查一级评论
    INDEX idx_root (root_id, created_at),                       -- 查某一级评论下的子回复
    INDEX idx_user (user_id, created_at)                        -- 查某用户的评论
) ENGINE=InnoDB COMMENT='评论表';

-- 评论点赞表（记录谁赞了哪条评论）
CREATE TABLE comment_like (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    comment_id  BIGINT   NOT NULL COMMENT '评论ID',
    user_id     BIGINT   NOT NULL COMMENT '点赞用户',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comment_user (comment_id, user_id) -- 防重复点赞
) ENGINE=InnoDB COMMENT='评论点赞';

-- 评论计数表（独立计数，减少主表热点更新）
CREATE TABLE comment_count (
    obj_type      TINYINT NOT NULL COMMENT '对象类型',
    obj_id        BIGINT  NOT NULL COMMENT '对象ID',
    comment_count INT     NOT NULL DEFAULT 0 COMMENT '总评论数',
    PRIMARY KEY (obj_type, obj_id)
) ENGINE=InnoDB COMMENT='评论计数';
```

```
字段关系说明:

  文章 X 的评论区:
  ┌──────────────────────────────────────────────────────────────┐
  │ id=1001  root_id=0     parent_id=0     "这篇文章写得真好"    │
  │   ├─ id=1005  root_id=1001  parent_id=1001  "同意楼上"      │
  │   ├─ id=1008  root_id=1001  parent_id=1005  "@1005 +1"     │
  │   └─ id=1012  root_id=1001  parent_id=1001  "确实不错"      │
  │                                                              │
  │ id=1002  root_id=0     parent_id=0     "有个问题想请教"       │
  │   └─ id=1006  root_id=1002  parent_id=1002  "你说"          │
  │                                                              │
  │ id=1003  root_id=0     parent_id=0     "排版很好看"           │
  └──────────────────────────────────────────────────────────────┘

  root_id 的作用:
    所有子回复的 root_id 都指向一级评论
    查子回复: WHERE root_id = 1001 ORDER BY created_at
    → 避免递归查询，一条 SQL 拿到所有子回复（扁平化）

  parent_id 的作用:
    指向直接回复的那条评论（用于显示 "@用户名"）
    一级评论: parent_id = 0
    回复一级评论: parent_id = root_id
    回复某条子回复: parent_id = 那条子回复的 id
```

### 4. Go 服务端实现

```go
package comment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CommentService struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewCommentService(db *sql.DB, rdb *redis.Client) *CommentService {
	return &CommentService{db: db, rdb: rdb}
}

// ---------- 数据结构 ----------

type Comment struct {
	ID            int64  `json:"id"`
	ObjType       int    `json:"obj_type"`
	ObjID         int64  `json:"obj_id"`
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`      // 关联查询
	UserAvatar    string `json:"user_avatar"`     // 关联查询
	Content       string `json:"content"`
	RootID        int64  `json:"root_id"`
	ParentID      int64  `json:"parent_id"`
	ReplyUserID   int64  `json:"reply_user_id"`
	ReplyUserName string `json:"reply_user_name"` // 关联查询
	LikeCount     int    `json:"like_count"`
	ReplyCount    int    `json:"reply_count"`
	IsLiked       bool   `json:"is_liked"`  // 当前用户是否点赞
	IPLocation    string `json:"ip_location"`
	CreatedAt     string `json:"created_at"`
	Replies       []*Comment `json:"replies,omitempty"` // 子回复（前 N 条预览）
}

type ListRequest struct {
	ObjType  int    `json:"obj_type"`
	ObjID    int64  `json:"obj_id"`
	SortBy   string `json:"sort_by"` // "time" / "hot"
	Cursor   int64  `json:"cursor"`  // 游标分页
	PageSize int    `json:"page_size"`
	UserID   int64  `json:"user_id"` // 当前登录用户（判断是否点赞）
}

type ListResponse struct {
	Comments     []*Comment `json:"comments"`
	NextCursor   int64      `json:"next_cursor"`
	HasMore      bool       `json:"has_more"`
	TotalCount   int64      `json:"total_count"`
}

// ---------- 发表评论 ----------

type PostRequest struct {
	ObjType     int    `json:"obj_type"`
	ObjID       int64  `json:"obj_id"`
	UserID      int64  `json:"user_id"`
	Content     string `json:"content"`
	RootID      int64  `json:"root_id"`   // 0=一级评论
	ParentID    int64  `json:"parent_id"` // 0=一级评论
	ReplyUserID int64  `json:"reply_user_id"`
	IPLocation  string `json:"ip_location"`
}

func (s *CommentService) Post(ctx context.Context, req *PostRequest) (int64, error) {
	if len(req.Content) == 0 || len(req.Content) > 2000 {
		return 0, errors.New("评论内容长度需在1-2000之间")
	}

	// TODO: 敏感词过滤、审核

	// 1. 插入评论
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO comment (obj_type, obj_id, user_id, content, root_id, parent_id, reply_user_id, ip_location, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		req.ObjType, req.ObjID, req.UserID, req.Content,
		req.RootID, req.ParentID, req.ReplyUserID, req.IPLocation,
	)
	if err != nil {
		return 0, err
	}
	commentID, _ := result.LastInsertId()

	// 2. 更新计数（异步更好，这里简化为同步）
	// 文章总评论数 +1
	s.db.ExecContext(ctx,
		`INSERT INTO comment_count (obj_type, obj_id, comment_count)
		 VALUES (?, ?, 1) ON DUPLICATE KEY UPDATE comment_count = comment_count + 1`,
		req.ObjType, req.ObjID,
	)

	// 如果是子回复，更新一级评论的 reply_count
	if req.RootID > 0 {
		s.db.ExecContext(ctx,
			"UPDATE comment SET reply_count = reply_count + 1 WHERE id = ?",
			req.RootID,
		)
	}

	// 3. 清除缓存
	s.invalidateCache(ctx, req.ObjType, req.ObjID)

	// 4. 发通知（异步 MQ）
	// producer.Send("comment-notify", ...)

	return commentID, nil
}

// ---------- 查询评论列表 ----------

func (s *CommentService) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 1. 查询一级评论（root_id = 0）
	var query string
	var args []interface{}

	if req.SortBy == "hot" {
		// 按热度排序（点赞数降序，点赞相同按时间降序）
		// 游标分页用 (like_count, id) 组合游标
		query = `SELECT id, user_id, content, root_id, parent_id, reply_user_id,
		                like_count, reply_count, ip_location, created_at
		         FROM comment
		         WHERE obj_type = ? AND obj_id = ? AND root_id = 0 AND status = 1`
		args = append(args, req.ObjType, req.ObjID)

		if req.Cursor > 0 {
			// 用 comment ID 作为游标，先查上一页最后一条的 like_count
			var lastLike int
			s.db.QueryRowContext(ctx, "SELECT like_count FROM comment WHERE id = ?", req.Cursor).Scan(&lastLike)
			query += ` AND (like_count < ? OR (like_count = ? AND id < ?))`
			args = append(args, lastLike, lastLike, req.Cursor)
		}
		query += ` ORDER BY like_count DESC, id DESC LIMIT ?`
		args = append(args, req.PageSize+1)
	} else {
		// 按时间排序（游标 = 上一页最后一条的 ID）
		query = `SELECT id, user_id, content, root_id, parent_id, reply_user_id,
		                like_count, reply_count, ip_location, created_at
		         FROM comment
		         WHERE obj_type = ? AND obj_id = ? AND root_id = 0 AND status = 1`
		args = append(args, req.ObjType, req.ObjID)

		if req.Cursor > 0 {
			query += ` AND id < ?`
			args = append(args, req.Cursor)
		}
		query += ` ORDER BY id DESC LIMIT ?`
		args = append(args, req.PageSize+1)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		c := &Comment{}
		rows.Scan(&c.ID, &c.UserID, &c.Content, &c.RootID, &c.ParentID,
			&c.ReplyUserID, &c.LikeCount, &c.ReplyCount, &c.IPLocation, &c.CreatedAt)
		comments = append(comments, c)
	}

	// 判断是否有下一页
	hasMore := len(comments) > req.PageSize
	if hasMore {
		comments = comments[:req.PageSize]
	}

	// 2. 批量查询每条一级评论的前 3 条子回复（预览）
	if len(comments) > 0 {
		rootIDs := make([]int64, len(comments))
		for i, c := range comments {
			rootIDs[i] = c.ID
		}
		repliesMap := s.batchGetReplies(ctx, rootIDs, 3)
		for _, c := range comments {
			c.Replies = repliesMap[c.ID]
		}
	}

	// 3. 批量查询当前用户是否点赞
	if req.UserID > 0 && len(comments) > 0 {
		s.batchCheckLiked(ctx, req.UserID, comments)
	}

	// 4. 批量查询用户信息
	// s.batchGetUserInfo(ctx, comments)

	resp := &ListResponse{
		Comments: comments,
		HasMore:  hasMore,
	}
	if hasMore && len(comments) > 0 {
		resp.NextCursor = comments[len(comments)-1].ID
	}

	// 总评论数
	s.db.QueryRowContext(ctx,
		"SELECT comment_count FROM comment_count WHERE obj_type = ? AND obj_id = ?",
		req.ObjType, req.ObjID,
	).Scan(&resp.TotalCount)

	return resp, nil
}

// batchGetReplies 批量查询子回复（每个一级评论取前 N 条）
func (s *CommentService) batchGetReplies(ctx context.Context, rootIDs []int64, limit int) map[int64][]*Comment {
	if len(rootIDs) == 0 {
		return nil
	}

	// 构造 IN 查询 + 子查询取前 N 条
	// 简化实现: 逐个查询（生产环境可用 ROW_NUMBER 窗口函数或 UNION ALL）
	result := make(map[int64][]*Comment)
	for _, rootID := range rootIDs {
		rows, err := s.db.QueryContext(ctx,
			`SELECT id, user_id, content, root_id, parent_id, reply_user_id,
			        like_count, ip_location, created_at
			 FROM comment
			 WHERE root_id = ? AND status = 1
			 ORDER BY created_at ASC
			 LIMIT ?`,
			rootID, limit,
		)
		if err != nil {
			continue
		}
		var replies []*Comment
		for rows.Next() {
			c := &Comment{}
			rows.Scan(&c.ID, &c.UserID, &c.Content, &c.RootID, &c.ParentID,
				&c.ReplyUserID, &c.LikeCount, &c.IPLocation, &c.CreatedAt)
			replies = append(replies, c)
		}
		rows.Close()
		result[rootID] = replies
	}
	return result
}

// ListReplies 展开某条一级评论的全部子回复（游标分页）
func (s *CommentService) ListReplies(ctx context.Context, rootID, cursor int64, pageSize int) ([]*Comment, bool, error) {
	if pageSize <= 0 {
		pageSize = 20
	}

	query := `SELECT id, user_id, content, root_id, parent_id, reply_user_id,
	                 like_count, ip_location, created_at
	          FROM comment
	          WHERE root_id = ? AND status = 1`
	args := []interface{}{rootID}

	if cursor > 0 {
		query += " AND id > ?"
		args = append(args, cursor)
	}
	query += " ORDER BY created_at ASC LIMIT ?"
	args = append(args, pageSize+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var replies []*Comment
	for rows.Next() {
		c := &Comment{}
		rows.Scan(&c.ID, &c.UserID, &c.Content, &c.RootID, &c.ParentID,
			&c.ReplyUserID, &c.LikeCount, &c.IPLocation, &c.CreatedAt)
		replies = append(replies, c)
	}

	hasMore := len(replies) > pageSize
	if hasMore {
		replies = replies[:pageSize]
	}
	return replies, hasMore, nil
}
```

### 5. 点赞实现

```
点赞挑战:
  - 大V 一条评论可能有百万点赞 → 高并发写
  - 需要判断当前用户是否已赞 → 高并发读
  - 点赞数需要实时展示 → 计数准确性

方案: Redis 缓存 + 异步落库

  ┌────────┐        ┌──────────┐        ┌────────┐      ┌────────┐
  │  客户端  │        │ 点赞服务  │        │ Redis  │      │ MySQL  │
  └───┬────┘        └────┬─────┘        └───┬────┘      └───┬────┘
      │ 1.点赞           │                  │               │
      │─────────────────>│                  │               │
      │                  │ 2.SADD 已赞集合   │               │
      │                  │─────────────────>│               │
      │                  │ 3.INCR 点赞计数   │               │
      │                  │─────────────────>│               │
      │ 4.返回成功        │                  │               │
      │<─────────────────│                  │               │
      │                  │ 5.异步MQ落库      │               │
      │                  │─────────────────────────────────>│
```

```go
// ---------- 点赞 ----------

func (s *CommentService) Like(ctx context.Context, commentID, userID int64) error {
	likeKey := fmt.Sprintf("comment:like:users:%d", commentID)  // 已赞用户集合
	countKey := fmt.Sprintf("comment:like:count:%d", commentID) // 点赞计数

	// 1. 检查是否已赞 + 加入集合（原子操作）
	added, err := s.rdb.SAdd(ctx, likeKey, userID).Result()
	if err != nil {
		return err
	}
	if added == 0 {
		return errors.New("已经点过赞了")
	}

	// 2. 计数 +1
	s.rdb.Incr(ctx, countKey)

	// 3. 异步落库（MQ 或 goroutine）
	go func() {
		s.db.Exec(
			"INSERT IGNORE INTO comment_like (comment_id, user_id) VALUES (?, ?)",
			commentID, userID,
		)
		s.db.Exec(
			"UPDATE comment SET like_count = like_count + 1 WHERE id = ?",
			commentID,
		)
	}()

	return nil
}

// Unlike 取消点赞
func (s *CommentService) Unlike(ctx context.Context, commentID, userID int64) error {
	likeKey := fmt.Sprintf("comment:like:users:%d", commentID)
	countKey := fmt.Sprintf("comment:like:count:%d", commentID)

	removed, err := s.rdb.SRem(ctx, likeKey, userID).Result()
	if err != nil {
		return err
	}
	if removed == 0 {
		return errors.New("未点赞")
	}

	s.rdb.Decr(ctx, countKey)

	go func() {
		s.db.Exec("DELETE FROM comment_like WHERE comment_id = ? AND user_id = ?", commentID, userID)
		s.db.Exec("UPDATE comment SET like_count = like_count - 1 WHERE id = ? AND like_count > 0", commentID)
	}()

	return nil
}

// batchCheckLiked 批量检查是否点赞
func (s *CommentService) batchCheckLiked(ctx context.Context, userID int64, comments []*Comment) {
	pipe := s.rdb.Pipeline()
	cmds := make([]*redis.BoolCmd, len(comments))

	for i, c := range comments {
		likeKey := fmt.Sprintf("comment:like:users:%d", c.ID)
		cmds[i] = pipe.SIsMember(ctx, likeKey, userID)
	}
	pipe.Exec(ctx)

	for i, cmd := range cmds {
		comments[i].IsLiked, _ = cmd.Result()
		// 子回复也需要检查
		for _, reply := range comments[i].Replies {
			rKey := fmt.Sprintf("comment:like:users:%d", reply.ID)
			reply.IsLiked, _ = s.rdb.SIsMember(ctx, rKey, userID).Result()
		}
	}
}
```

### 6. 热评排序算法

```
热评排序不只是按点赞数排序，需要综合考虑:

方案一: 简单点赞排序
  ORDER BY like_count DESC, id DESC
  ❌ 问题: 老评论积累优势大，新评论没机会

方案二: Wilson 置信区间（Reddit 采用）
  考虑: 点赞数 + 踩数 → 计算"真实好评率"的下界
  score = (p̂ + z²/2n - z√(p̂(1-p̂)/n + z²/4n²)) / (1 + z²/n)
  p̂ = 点赞/(点赞+踩), n = 总票数, z = 1.96(95%置信度)
  ✅ 票数少时分数低（避免"1赞0踩=100%"排第一）

方案三: 热度衰减公式（推荐，类 Hacker News）
  score = (like_count - 1) / (T + 2)^G
  T = 评论发布后经过的小时数
  G = 重力因子（通常 1.5~1.8）
  ✅ 兼顾热度和新鲜度

方案四: 加权综合打分（实际生产常用）
  score = w1 × like_count + w2 × reply_count + w3 × time_decay

  time_decay = 1 / (1 + hours_since_post / 24)  // 24小时半衰

  权重示例:
    w1 = 3（点赞权重最高）
    w2 = 2（回复也说明有讨论价值）
    w3 = 10（新鲜度基础分）
```

#### Go 实现热评排序

```go
import "math"

// HotScore 计算评论热度分数
func HotScore(likeCount, replyCount int, createdAt time.Time) float64 {
	hoursSince := time.Since(createdAt).Hours()
	if hoursSince < 0 {
		hoursSince = 0
	}

	// 时间衰减因子（24小时半衰期）
	timeDecay := 1.0 / (1.0 + hoursSince/24.0)

	// 加权打分
	score := 3.0*float64(likeCount) + 2.0*float64(replyCount) + 10.0*timeDecay

	return score
}

// WilsonScore Wilson 置信区间下界（适合有踩功能的场景）
func WilsonScore(upvotes, downvotes int) float64 {
	n := float64(upvotes + downvotes)
	if n == 0 {
		return 0
	}
	z := 1.96 // 95% 置信度
	p := float64(upvotes) / n

	left := p + z*z/(2*n)
	right := z * math.Sqrt((p*(1-p)+z*z/(4*n))/n)
	under := 1 + z*z/n

	return (left - right) / under
}

// 定时更新热评排序（离线计算）
// 生产环境用定时任务每隔 N 分钟重新计算 top K 热评，存入 Redis
func (s *CommentService) RefreshHotComments(ctx context.Context, objType int, objID int64) error {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, like_count, reply_count, created_at
		 FROM comment
		 WHERE obj_type = ? AND obj_id = ? AND root_id = 0 AND status = 1`,
		objType, objID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	type scored struct {
		id    int64
		score float64
	}
	var items []scored

	for rows.Next() {
		var id int64
		var likeCount, replyCount int
		var createdAt time.Time
		rows.Scan(&id, &likeCount, &replyCount, &createdAt)
		items = append(items, scored{id: id, score: HotScore(likeCount, replyCount, createdAt)})
	}

	// 存入 Redis Sorted Set
	hotKey := fmt.Sprintf("comment:hot:%d:%d", objType, objID)
	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, hotKey)
	for _, item := range items {
		pipe.ZAdd(ctx, hotKey, redis.Z{Score: item.score, Member: item.id})
	}
	pipe.Expire(ctx, hotKey, 10*time.Minute)
	_, err = pipe.Exec(ctx)
	return err
}
```

### 7. 缓存设计

```
评论区缓存策略:

  ┌──────────────────────┬───────────────────────────────────────┐
  │ 缓存内容              │ Redis 结构 & 策略                      │
  ├──────────────────────┼───────────────────────────────────────┤
  │ 文章评论数            │ STRING  comment:count:{obj}:{id}      │
  │                      │ TTL 10min, 异步更新                    │
  ├──────────────────────┼───────────────────────────────────────┤
  │ 一级评论列表(按时间)  │ ZSET    comment:list:{obj}:{id}       │
  │                      │ score=commentID, TTL 5min             │
  │                      │ 只缓存前 N 页（热数据）                │
  ├──────────────────────┼───────────────────────────────────────┤
  │ 热评列表              │ ZSET    comment:hot:{obj}:{id}        │
  │                      │ score=热度分, 定时重算                  │
  ├──────────────────────┼───────────────────────────────────────┤
  │ 评论点赞集合          │ SET     comment:like:users:{cid}      │
  │                      │ 用于判断"是否已赞"                     │
  ├──────────────────────┼───────────────────────────────────────┤
  │ 评论点赞数            │ STRING  comment:like:count:{cid}      │
  │                      │ INCR/DECR 原子操作                     │
  └──────────────────────┴───────────────────────────────────────┘

缓存更新策略:
  发评论 → 清除该文章的评论列表缓存
  点赞   → Redis 实时更新（计数+集合），异步落库
  热评   → 定时任务每 5~10 分钟重算
```

```go
func (s *CommentService) invalidateCache(ctx context.Context, objType int, objID int64) {
	pipe := s.rdb.Pipeline()
	prefix := fmt.Sprintf("comment:list:%d:%d", objType, objID)
	pipe.Del(ctx, prefix)
	pipe.Del(ctx, fmt.Sprintf("comment:count:%d:%d", objType, objID))
	pipe.Exec(ctx)
}
```

### 8. 高并发优化

```
场景: 顶流明星发动态 → 百万评论 → 评论区几十万 QPS

优化策略:
  ┌───────────────────┬──────────────────────────────────────────┐
  │ 层级               │ 措施                                     │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 接入层             │ CDN 缓存静态评论页 + 接口限流              │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 缓存层             │ 热门文章评论列表 → Redis 缓存（前 N 页）   │
  │                   │ 冷门文章 → 穿透到 DB 后回填缓存            │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 读写分离           │ 评论读 → 从库                             │
  │                   │ 评论写 → 主库                             │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 异步化             │ 发评论 → MQ → 异步写入（用户先看到"发送中"）│
  │                   │ 点赞 → Redis 实时，DB 异步                │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 分库分表           │ 按 obj_id 分表（同一篇文章的评论在一张表） │
  │                   │ comment_0 ~ comment_63（64 张表）         │
  ├───────────────────┼──────────────────────────────────────────┤
  │ 预加载             │ 一级评论查出后，批量预加载子回复（避免 N+1） │
  │                   │ 批量查用户信息、批量查点赞状态             │
  └───────────────────┴──────────────────────────────────────────┘

分页优化:
  ❌ OFFSET 分页: SELECT ... LIMIT 10 OFFSET 10000 → 越往后越慢
  ✅ 游标分页:    SELECT ... WHERE id < {cursor} LIMIT 10 → 恒定性能

  按时间排序:  cursor = 上一页最后一条的 ID
  按热度排序:  cursor = 上一页最后一条的 ID（配合 like_count 组合条件）
```

### 9. 评论审核与安全

```
评论安全策略:
  ┌──────────────────────────────────────────────────────────────┐
  │ 1. 发评论前:                                                  │
  │    - 敏感词过滤（DFA / AC 自动机，本地词库+远程词库）           │
  │    - 频率限制（同一用户 N 秒内只能发 1 条）                    │
  │    - 内容长度限制（最大 2000 字）                              │
  │                                                              │
  │ 2. 发评论后:                                                  │
  │    - 机审: 接入 NLP 模型判断（垃圾/广告/色情/政治）            │
  │    - 人审: 机审不确定的进入人工审核队列                        │
  │    - 先发后审（多数平台）: 先展示，审核不通过再折叠/删除        │
  │                                                              │
  │ 3. 防刷:                                                      │
  │    - 同一内容重复发送 → 去重拦截                               │
  │    - IP/设备维度限流                                           │
  │    - 新用户发评论 → 需先过验证码                               │
  └──────────────────────────────────────────────────────────────┘
```

### 10. 多级嵌套评论的存储方案（扩展）

```
如果需要 Reddit 式无限嵌套，有两种经典方案:

方案一: 邻接表（Adjacency List）
  每条评论存 parent_id → 查询时递归/多次查询

  优点: 写入简单
  缺点: 查询需递归（MySQL 8+ 支持 WITH RECURSIVE CTE）

  WITH RECURSIVE comment_tree AS (
      SELECT id, content, parent_id, 0 AS depth
      FROM comment WHERE id = 1001
      UNION ALL
      SELECT c.id, c.content, c.parent_id, ct.depth + 1
      FROM comment c JOIN comment_tree ct ON c.parent_id = ct.id
  )
  SELECT * FROM comment_tree ORDER BY depth, id;

方案二: 物化路径（Materialized Path）
  存一个 path 字段: "/1001/1005/1008"

  查所有子评论: WHERE path LIKE '/1001/%'
  查深度: LENGTH(path) - LENGTH(REPLACE(path,'/','')) - 1
  排序:   ORDER BY path（天然层级排序）

  优点: 查询简单高效
  缺点: 移动节点时需要更新所有子节点的 path

方案三: 闭包表（Closure Table）
  额外一张 comment_tree(ancestor, descendant, depth)
  任意两个有祖先关系的节点都有一条记录

  优点: 查询最灵活
  缺点: 空间开销大（O(N²) 最坏情况）

对比:
  │ 方案       │ 写入  │ 查子树  │ 查深度  │ 移动节点  │ 适用场景      │
  │ 邻接表     │ O(1) │ 递归    │ 递归    │ O(1)     │ 层级少(≤5层) │
  │ 物化路径   │ O(1) │ LIKE   │ O(1)   │ O(N)     │ 读多写少      │
  │ 闭包表     │ O(D) │ O(1)   │ O(1)   │ O(N)     │ 查询灵活      │

  D = 节点深度, N = 子树节点数

  推荐: 二级评论模型（绝大多数场景够用，性能最好）
        如果必须多级嵌套 → 物化路径（简单实用）
```

### 一句话总结

**评论系统推荐二级评论模型（微博/B站模式）：一级评论 root_id=0 + 子回复 root_id 指向一级评论（扁平化，一条 SQL 查所有子回复，避免递归）。** 热评排序用加权公式 score=3×like+2×reply+10×timeDecay（定时重算存 Redis ZSET），而非简单按赞排序（避免老评论永远霸榜）。点赞用 Redis SET（SADD/SISMEMBER 判断已赞）+ STRING（INCR 计数），异步落库。高并发优化：热门评论 Redis 缓存前 N 页 + 游标分页（WHERE id < cursor，避免 OFFSET 深分页）+ 按 obj_id 分表 + 批量预加载子回复和用户信息避免 N+1。如需无限嵌套评论，推荐物化路径（path="/1001/1005/1008"，LIKE 查子树）。

---

## 十七、通知中心设计

> 多渠道通知（APP 推送/短信/邮件/站内信）怎么设计？模板化？防骚扰？重试？

### 1. 通知中心核心需求

```
功能需求:
  ┌──────────────────────────────────────────────────────────────┐
  │ 1. 多渠道发送:  APP Push / 短信 / 邮件 / 站内信 / 企业微信    │
  │ 2. 模板管理:    通知内容模板化，运营可配置，支持多语言          │
  │ 3. 消息分类:    系统通知 / 营销通知 / 交易通知 / 社交通知      │
  │ 4. 用户偏好:    用户可关闭某类通知/某渠道                      │
  │ 5. 防骚扰:      频率限制、免打扰时段、同内容去重               │
  │ 6. 可靠投递:    失败重试、状态追踪、送达回执                   │
  │ 7. 高性能:      百万级用户群发（如全站公告、营销活动）          │
  └──────────────────────────────────────────────────────────────┘

技术挑战:
  ┌──────────────────────────────────────────────────────────────┐
  │ 1. 渠道异构:  每个渠道 API 不同、限流不同、回执不同            │
  │ 2. 群发性能:  100 万用户群发短信，不能串行逐条调 API           │
  │ 3. 可靠性:    短信网关超时、Push 服务抖动 → 需要重试机制       │
  │ 4. 防骚扰:    同一用户一天收到 50 条推送 → 用户卸载 APP        │
  │ 5. 成本控制:  短信/电话有成本，不能重复发送、无效发送           │
  └──────────────────────────────────────────────────────────────┘
```

### 2. 整体架构

```
┌──────────────────────────────────────────────────────────────────┐
│                      业务系统 (调用方)                             │
│  订单服务 / 营销服务 / 用户服务 / 风控服务 / 运营后台              │
└────────┬──────────┬──────────┬──────────┬───────────────────────┘
         │          │          │          │
         ▼          ▼          ▼          ▼
┌──────────────────────────────────────────────────────────────────┐
│                   通知网关 (Notification Gateway)                  │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────────┐ │
│  │参数校验  │ │模板渲染   │ │用户偏好  │ │防骚扰(频率/去重/免打扰)│ │
│  └─────────┘ └──────────┘ └──────────┘ └──────────────────────┘ │
└────────────────────────┬────────────────────────────────────────┘
                         │ 写入 MQ
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                        Kafka / RabbitMQ                           │
│  topic-push / topic-sms / topic-email / topic-inbox               │
└──────┬──────────┬──────────┬──────────┬─────────────────────────┘
       │          │          │          │
       ▼          ▼          ▼          ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│ Push     │ │ SMS      │ │ Email    │ │ Inbox    │
│ Worker   │ │ Worker   │ │ Worker   │ │ Worker   │
│          │ │          │ │          │ │          │
│ APNs    │ │ 阿里云    │ │ SES     │ │ MySQL   │
│ FCM     │ │ 腾讯云    │ │ SendGrid│ │ 站内信   │
│ 华为/小米│ │ 云片     │ │ Mailgun │ │          │
└──────┬───┘ └──────┬───┘ └─────┬───┘ └──────┬───┘
       │            │           │             │
       ▼            ▼           ▼             ▼
┌──────────────────────────────────────────────────────────────────┐
│                     发送记录 & 状态追踪                            │
│  MySQL(发送记录) + Redis(频率限制/去重) + ClickHouse(统计分析)     │
└──────────────────────────────────────────────────────────────────┘

核心设计原则:
  1. 业务系统只管"通知什么"，不管"怎么发"（解耦）
  2. 渠道可插拔: 新增渠道只需实现 Sender 接口 + 注册
  3. MQ 削峰: 群发时百万消息先进队列，Worker 按各渠道限速消费
  4. 失败重试: Worker 消费失败 → 延迟队列重试（最多 3 次）
```

### 3. 数据库设计

```sql
-- 通知模板表
CREATE TABLE notify_template (
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    template_code VARCHAR(64)   NOT NULL COMMENT '模板编码(如 ORDER_PAY_SUCCESS)',
    name          VARCHAR(128)  NOT NULL COMMENT '模板名称',
    category      VARCHAR(32)   NOT NULL COMMENT '分类: system/marketing/transaction/social',
    channels      VARCHAR(128)  NOT NULL COMMENT '支持渠道(逗号分隔): push,sms,email,inbox',
    title_tpl     VARCHAR(256)  NULL     COMMENT '标题模板(支持变量 {{.OrderNo}})',
    content_tpl   TEXT          NOT NULL COMMENT '内容模板',
    sms_tpl_id    VARCHAR(64)   NULL     COMMENT '短信平台模板ID(需审核)',
    status        TINYINT       NOT NULL DEFAULT 1 COMMENT '1启用 0禁用',
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (template_code)
) ENGINE=InnoDB COMMENT='通知模板';

-- 发送记录表（按月分表）
CREATE TABLE notify_record (
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_id    VARCHAR(64)   NOT NULL COMMENT '请求唯一ID(幂等)',
    template_code VARCHAR(64)   NOT NULL COMMENT '模板编码',
    channel       VARCHAR(16)   NOT NULL COMMENT '发送渠道: push/sms/email/inbox',
    user_id       BIGINT        NOT NULL COMMENT '接收用户',
    receiver      VARCHAR(128)  NOT NULL COMMENT '接收地址(手机号/邮箱/deviceToken)',
    title         VARCHAR(256)  NULL     COMMENT '渲染后的标题',
    content       TEXT          NOT NULL COMMENT '渲染后的内容',
    status        TINYINT       NOT NULL DEFAULT 0 COMMENT '0待发送 1发送中 2成功 3失败 4用户关闭',
    fail_reason   VARCHAR(256)  NULL     COMMENT '失败原因',
    retry_count   TINYINT       NOT NULL DEFAULT 0 COMMENT '已重试次数',
    vendor_msg_id VARCHAR(128)  NULL     COMMENT '渠道方消息ID(用于查回执)',
    send_at       DATETIME      NULL     COMMENT '实际发送时间',
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_request_channel (request_id, channel),
    INDEX idx_user (user_id, created_at),
    INDEX idx_status (status, created_at)
) ENGINE=InnoDB COMMENT='通知发送记录';

-- 用户通知偏好表
CREATE TABLE user_notify_preference (
    user_id       BIGINT      NOT NULL,
    category      VARCHAR(32) NOT NULL COMMENT '通知分类',
    channel       VARCHAR(16) NOT NULL COMMENT '渠道',
    enabled       TINYINT     NOT NULL DEFAULT 1 COMMENT '1开启 0关闭',
    updated_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, category, channel)
) ENGINE=InnoDB COMMENT='用户通知偏好';

-- 站内信表
CREATE TABLE inbox_message (
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id       BIGINT       NOT NULL,
    title         VARCHAR(256) NOT NULL,
    content       TEXT         NOT NULL,
    category      VARCHAR(32)  NOT NULL,
    is_read       TINYINT      NOT NULL DEFAULT 0,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_unread (user_id, is_read, created_at)
) ENGINE=InnoDB COMMENT='站内信';
```

### 4. 核心 Go 实现

#### 通知接口与数据结构

```go
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Channel 通知渠道
type Channel string

const (
	ChannelPush  Channel = "push"
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
	ChannelInbox Channel = "inbox"
)

// Category 通知分类
type Category string

const (
	CategorySystem      Category = "system"      // 系统通知(验证码/安全提醒)
	CategoryTransaction Category = "transaction" // 交易通知(支付/发货/退款)
	CategoryMarketing   Category = "marketing"   // 营销通知(优惠券/活动)
	CategorySocial      Category = "social"      // 社交通知(点赞/评论/关注)
)

// NotifyRequest 通知请求（业务方调用）
type NotifyRequest struct {
	RequestID    string                 `json:"request_id"`    // 幂等 ID
	TemplateCode string                 `json:"template_code"` // 模板编码
	UserIDs      []int64                `json:"user_ids"`      // 接收用户列表
	Channels     []Channel              `json:"channels"`      // 指定渠道(为空则按模板默认)
	Variables    map[string]string      `json:"variables"`     // 模板变量
	SendAt       *time.Time             `json:"send_at"`       // 定时发送(为空则立即发送)
	Extra        map[string]interface{} `json:"extra"`         // 扩展(Push的跳转链接等)
}

// NotifyMessage 经处理后的单条发送消息（写入 MQ）
type NotifyMessage struct {
	RequestID    string                 `json:"request_id"`
	TemplateCode string                 `json:"template_code"`
	Category     Category               `json:"category"`
	Channel      Channel                `json:"channel"`
	UserID       int64                  `json:"user_id"`
	Receiver     string                 `json:"receiver"`     // 手机号/邮箱/deviceToken
	Title        string                 `json:"title"`        // 渲染后标题
	Content      string                 `json:"content"`      // 渲染后内容
	Extra        map[string]interface{} `json:"extra"`
	RetryCount   int                    `json:"retry_count"`
}

// Sender 渠道发送接口（策略模式）
type Sender interface {
	Send(ctx context.Context, msg *NotifyMessage) error
	Channel() Channel
}
```

#### 通知网关（核心入口）

```go
// NotifyGateway 通知网关
type NotifyGateway struct {
	templateStore TemplateStore     // 模板存储
	userStore     UserStore         // 用户信息/偏好存储
	producer      MQProducer        // MQ 生产者
	rdb           *redis.Client     // Redis(频率限制/去重)
	antiSpam      *AntiSpamService  // 防骚扰服务
}

type TemplateStore interface {
	GetByCode(ctx context.Context, code string) (*NotifyTemplate, error)
}

type UserStore interface {
	GetUserReceivers(ctx context.Context, userID int64) (*UserReceiverInfo, error)
	IsChannelEnabled(ctx context.Context, userID int64, category Category, channel Channel) (bool, error)
}

type MQProducer interface {
	Send(ctx context.Context, topic string, msg []byte) error
}

type NotifyTemplate struct {
	Code       string
	Name       string
	Category   Category
	Channels   []Channel
	TitleTpl   string
	ContentTpl string
	SMSTplID   string
}

type UserReceiverInfo struct {
	UserID      int64
	Phone       string
	Email       string
	DeviceToken string // APNs/FCM token
}

// Send 发送通知（业务方的统一入口）
func (g *NotifyGateway) Send(ctx context.Context, req *NotifyRequest) error {
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}

	// 1. 查询模板
	tpl, err := g.templateStore.GetByCode(ctx, req.TemplateCode)
	if err != nil {
		return fmt.Errorf("模板不存在: %s", req.TemplateCode)
	}

	// 确定发送渠道
	channels := req.Channels
	if len(channels) == 0 {
		channels = tpl.Channels // 按模板默认渠道
	}

	// 2. 渲染模板
	title, err := renderTemplate(tpl.TitleTpl, req.Variables)
	if err != nil {
		return fmt.Errorf("渲染标题失败: %w", err)
	}
	content, err := renderTemplate(tpl.ContentTpl, req.Variables)
	if err != nil {
		return fmt.Errorf("渲染内容失败: %w", err)
	}

	// 3. 逐用户、逐渠道处理
	for _, userID := range req.UserIDs {
		// 获取用户接收地址
		receiver, err := g.userStore.GetUserReceivers(ctx, userID)
		if err != nil {
			continue // 用户不存在，跳过
		}

		for _, ch := range channels {
			// a. 获取接收地址
			addr := getReceiverAddr(receiver, ch)
			if addr == "" {
				continue // 用户未绑定该渠道（如未绑邮箱）
			}

			// b. 检查用户偏好（是否关闭了该类通知）
			enabled, _ := g.userStore.IsChannelEnabled(ctx, userID, tpl.Category, ch)
			if !enabled {
				continue
			}

			// c. 防骚扰检查
			if blocked, reason := g.antiSpam.Check(ctx, userID, ch, tpl.Category); blocked {
				fmt.Printf("防骚扰拦截: user=%d channel=%s reason=%s\n", userID, ch, reason)
				continue
			}

			// d. 幂等去重（同一 requestID + channel + userID 不重复发送）
			dedupeKey := fmt.Sprintf("notify:dedup:%s:%s:%d", req.RequestID, ch, userID)
			if ok, _ := g.rdb.SetNX(ctx, dedupeKey, 1, 24*time.Hour).Result(); !ok {
				continue // 重复请求
			}

			// e. 构建消息投递到 MQ
			msg := &NotifyMessage{
				RequestID:    req.RequestID,
				TemplateCode: req.TemplateCode,
				Category:     tpl.Category,
				Channel:      ch,
				UserID:       userID,
				Receiver:     addr,
				Title:        title,
				Content:      content,
				Extra:        req.Extra,
			}
			data, _ := json.Marshal(msg)

			topic := "notify-" + string(ch) // notify-push / notify-sms / ...
			if req.SendAt != nil && req.SendAt.After(time.Now()) {
				// 定时发送 → 延迟队列
				topic = "notify-delay"
			}
			g.producer.Send(ctx, topic, data)
		}
	}

	return nil
}

func getReceiverAddr(u *UserReceiverInfo, ch Channel) string {
	switch ch {
	case ChannelSMS:
		return u.Phone
	case ChannelEmail:
		return u.Email
	case ChannelPush:
		return u.DeviceToken
	case ChannelInbox:
		return fmt.Sprintf("%d", u.UserID) // 站内信用 userID
	}
	return ""
}

func renderTemplate(tplStr string, vars map[string]string) (string, error) {
	if tplStr == "" {
		return "", nil
	}
	t, err := template.New("").Parse(tplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
```

#### 渠道 Worker（消费 MQ + 发送 + 重试）

```go
// NotifyWorker 渠道发送 Worker
type NotifyWorker struct {
	senders  map[Channel]Sender // 各渠道发送器
	producer MQProducer         // MQ（用于重试）
	db       interface{}        // DB（写发送记录）
}

func NewNotifyWorker(producer MQProducer) *NotifyWorker {
	return &NotifyWorker{
		senders:  make(map[Channel]Sender),
		producer: producer,
	}
}

// RegisterSender 注册渠道发送器（插件化）
func (w *NotifyWorker) RegisterSender(s Sender) {
	w.senders[s.Channel()] = s
}

// HandleMessage 消费 MQ 消息
func (w *NotifyWorker) HandleMessage(ctx context.Context, data []byte) error {
	var msg NotifyMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("消息反序列化失败: %w", err)
	}

	sender, ok := w.senders[msg.Channel]
	if !ok {
		return fmt.Errorf("未注册的渠道: %s", msg.Channel)
	}

	// 发送
	err := sender.Send(ctx, &msg)
	if err != nil {
		// 失败重试（最多 3 次，指数退避）
		if msg.RetryCount < 3 {
			msg.RetryCount++
			retryData, _ := json.Marshal(msg)
			// 延迟重试: 第1次 30s, 第2次 60s, 第3次 120s
			retryTopic := fmt.Sprintf("notify-%s-retry", msg.Channel)
			w.producer.Send(ctx, retryTopic, retryData)
			return nil // 不返回错误，避免 MQ 再次重投
		}
		// 超过重试次数 → 记录失败
		// updateRecordStatus(msg.RequestID, msg.Channel, StatusFailed, err.Error())
		return fmt.Errorf("发送失败(已重试%d次): %w", msg.RetryCount, err)
	}

	// 发送成功 → 更新记录
	// updateRecordStatus(msg.RequestID, msg.Channel, StatusSuccess, "")
	return nil
}
```

#### 各渠道 Sender 实现

```go
// ---- APP Push 发送器 ----

type PushSender struct {
	// APNs / FCM / 华为 Push / 小米 Push 客户端
}

func (s *PushSender) Channel() Channel { return ChannelPush }

func (s *PushSender) Send(ctx context.Context, msg *NotifyMessage) error {
	// 根据 DeviceToken 前缀或用户设备信息判断走哪个推送通道
	// iOS → APNs, Android → FCM / 华为 / 小米
	payload := map[string]interface{}{
		"title": msg.Title,
		"body":  msg.Content,
		"data":  msg.Extra, // 跳转链接等
	}
	_ = payload
	// return apnsClient.Push(ctx, msg.Receiver, payload)
	fmt.Printf("[Push] to=%s title=%s\n", msg.Receiver, msg.Title)
	return nil
}

// ---- 短信发送器 ----

type SMSSender struct {
	// 阿里云短信 / 腾讯云短信 / 云片 客户端
}

func (s *SMSSender) Channel() Channel { return ChannelSMS }

func (s *SMSSender) Send(ctx context.Context, msg *NotifyMessage) error {
	// 调用短信网关 API
	// return aliyunSMS.SendSms(msg.Receiver, templateID, params)
	fmt.Printf("[SMS] to=%s content=%s\n", msg.Receiver, msg.Content)
	return nil
}

// ---- 邮件发送器 ----

type EmailSender struct {
	// SES / SendGrid / Mailgun / SMTP 客户端
}

func (s *EmailSender) Channel() Channel { return ChannelEmail }

func (s *EmailSender) Send(ctx context.Context, msg *NotifyMessage) error {
	// return ses.SendEmail(msg.Receiver, msg.Title, msg.Content)
	fmt.Printf("[Email] to=%s subject=%s\n", msg.Receiver, msg.Title)
	return nil
}

// ---- 站内信发送器 ----

type InboxSender struct {
	db interface{} // *sql.DB
}

func (s *InboxSender) Channel() Channel { return ChannelInbox }

func (s *InboxSender) Send(ctx context.Context, msg *NotifyMessage) error {
	// INSERT INTO inbox_message (user_id, title, content, category) VALUES (...)
	fmt.Printf("[Inbox] to=%d title=%s\n", msg.UserID, msg.Title)
	return nil
}
```

### 5. 防骚扰服务

```
防骚扰策略:
  ┌──────────────────────┬──────────────────────────────────────┐
  │ 策略                  │ 说明                                 │
  ├──────────────────────┼──────────────────────────────────────┤
  │ 频率限制(每用户)      │ 同一用户每天 Push ≤ 5 条             │
  │                      │ 同一用户每天 SMS ≤ 3 条               │
  │                      │ 营销类每周 ≤ 2 条                     │
  ├──────────────────────┼──────────────────────────────────────┤
  │ 免打扰时段            │ 22:00 ~ 08:00 不发 Push/SMS          │
  │                      │ 系统通知(验证码)例外                   │
  ├──────────────────────┼──────────────────────────────────────┤
  │ 同内容去重            │ 同一用户 1 小时内不重复发同样的内容    │
  ├──────────────────────┼──────────────────────────────────────┤
  │ 全局限流              │ 短信渠道 QPS ≤ 500（受运营商限制）    │
  │                      │ Push 渠道 QPS ≤ 10000                │
  ├──────────────────────┼──────────────────────────────────────┤
  │ 黑名单                │ 用户投诉/退订 → 加入黑名单            │
  └──────────────────────┴──────────────────────────────────────┘
```

#### Go 实现

```go
// AntiSpamService 防骚扰服务
type AntiSpamService struct {
	rdb *redis.Client
}

func NewAntiSpamService(rdb *redis.Client) *AntiSpamService {
	return &AntiSpamService{rdb: rdb}
}

// Check 检查是否被防骚扰拦截，返回 (是否拦截, 拦截原因)
func (a *AntiSpamService) Check(ctx context.Context, userID int64, ch Channel, cat Category) (bool, string) {
	// 1. 免打扰时段（22:00 ~ 08:00，系统通知例外）
	if cat != CategorySystem {
		hour := time.Now().Hour()
		if hour >= 22 || hour < 8 {
			return true, "免打扰时段"
		}
	}

	// 2. 频率限制
	if blocked, reason := a.checkFrequency(ctx, userID, ch, cat); blocked {
		return true, reason
	}

	return false, ""
}

func (a *AntiSpamService) checkFrequency(ctx context.Context, userID int64, ch Channel, cat Category) (bool, string) {
	today := time.Now().Format("20060102")

	// 每用户每天每渠道总量限制
	channelLimits := map[Channel]int64{
		ChannelPush:  8,
		ChannelSMS:   5,
		ChannelEmail: 10,
		ChannelInbox: 50,
	}

	dayKey := fmt.Sprintf("antispam:%d:%s:%s", userID, ch, today)
	count, _ := a.rdb.Incr(ctx, dayKey).Result()
	if count == 1 {
		a.rdb.Expire(ctx, dayKey, 25*time.Hour)
	}
	if limit, ok := channelLimits[ch]; ok && count > limit {
		a.rdb.Decr(ctx, dayKey) // 回退计数
		return true, fmt.Sprintf("%s渠道每日上限%d条", ch, limit)
	}

	// 营销类每周限制
	if cat == CategoryMarketing {
		weekKey := fmt.Sprintf("antispam:mkt:%d:%s", userID, weekNumber())
		mktCount, _ := a.rdb.Incr(ctx, weekKey).Result()
		if mktCount == 1 {
			a.rdb.Expire(ctx, weekKey, 8*24*time.Hour)
		}
		if mktCount > 2 {
			a.rdb.Decr(ctx, weekKey)
			return true, "营销通知每周上限2条"
		}
	}

	return false, ""
}

func weekNumber() string {
	y, w := time.Now().ISOWeek()
	return fmt.Sprintf("%d-W%02d", y, w)
}
```

### 6. 模板管理

```
模板示例:

  模板编码: ORDER_PAY_SUCCESS
  标题模板: "订单支付成功"
  内容模板: "您的订单 {{.OrderNo}} 已支付成功，金额 {{.Amount}} 元。"
  短信模板ID: SMS_20240101001（需提前在短信平台审核）

  调用方式:
  gateway.Send(ctx, &NotifyRequest{
      TemplateCode: "ORDER_PAY_SUCCESS",
      UserIDs:      []int64{10086},
      Channels:     []Channel{ChannelSMS, ChannelPush, ChannelInbox},
      Variables:    map[string]string{
          "OrderNo": "ORD20240101001",
          "Amount":  "99.00",
      },
  })

  渲染结果:
    Push:  标题="订单支付成功"  内容="您的订单 ORD20240101001 已支付成功，金额 99.00 元。"
    SMS:   调用短信平台模板API，传参 OrderNo=ORD20240101001, Amount=99.00
    站内信: 直接写入 inbox_message 表

多渠道内容差异化:
  同一事件不同渠道内容不同是常见需求
  ┌──────────────────┬─────────────────────────────────────────┐
  │ Push             │ "您的订单已支付成功 💰"（简短+emoji）     │
  │ SMS              │ "您的订单ORD001已支付99元。"（受字数限制） │
  │ Email            │ 富文本HTML（含图片、按钮、详细信息）       │
  │ 站内信            │ 完整信息 + 跳转链接                      │
  └──────────────────┴─────────────────────────────────────────┘

  实现: 模板表增加 push_content_tpl / sms_content_tpl / email_content_tpl
       或每个渠道一条模板记录，通过 template_code + channel 唯一定位
```

### 7. 群发方案（百万级）

```
场景: 运营要给 100 万用户群发营销 Push

❌ 错误方式:
  for userID := range allUsers {
      gateway.Send(ctx, &NotifyRequest{UserIDs: []int64{userID}, ...})
  }
  → 100 万次 RPC 调用 → 几小时才能发完

✅ 正确方式: 分批 + 异步 + 限速

  ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
  │ 运营后台  │────>│ 群发任务  │────>│ 用户分片  │────>│ MQ 投递  │
  │ 创建群发  │     │ (MySQL)  │     │ (每批1000)│     │ (Kafka)  │
  └──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                          │
                                                   ┌──────▼──────┐
                                                   │ Worker 集群  │
                                                   │ (限速消费)   │
                                                   └─────────────┘

流程:
  1. 运营创建群发任务 → 记录到 task 表（状态=待执行）
  2. 定时任务捞取待执行群发 → 按条件分页查询用户
  3. 每批 1000 个用户 → 构建消息 → 批量写入 Kafka
  4. Worker 消费时按渠道限速: SMS QPS=500, Push QPS=10000
  5. 实时统计: 总量/成功/失败/过滤 → 运营后台可查看进度
```

```go
// BroadcastTask 群发任务
func (g *NotifyGateway) Broadcast(ctx context.Context, templateCode string, userQuery string, variables map[string]string) (string, error) {
	taskID := uuid.New().String()

	// 1. 创建群发任务记录
	// INSERT INTO broadcast_task (id, template_code, user_query, status) VALUES (...)

	// 2. 异步执行（实际由独立 Worker 处理）
	go func() {
		var offset int64
		batchSize := int64(1000)

		for {
			// 分页查询目标用户
			userIDs := queryUserBatch(userQuery, offset, batchSize)
			if len(userIDs) == 0 {
				break
			}

			// 批量投递
			g.Send(context.Background(), &NotifyRequest{
				RequestID:    fmt.Sprintf("%s-%d", taskID, offset),
				TemplateCode: templateCode,
				UserIDs:      userIDs,
				Variables:    variables,
			})

			offset += batchSize
		}

		// 更新任务状态为完成
		// UPDATE broadcast_task SET status='done' WHERE id=taskID
	}()

	return taskID, nil
}

func queryUserBatch(query string, offset, limit int64) []int64 {
	// SELECT user_id FROM users WHERE {query} LIMIT offset, limit
	return nil
}
```

### 8. 送达率与监控

```
各渠道送达率差异:
  ┌──────────┬───────────┬──────────────────────────────────────┐
  │ 渠道      │ 送达率     │ 关键因素                              │
  ├──────────┼───────────┼──────────────────────────────────────┤
  │ 站内信    │ ~100%     │ 直接写 DB，最可靠                     │
  │ 邮件      │ 85~95%   │ 域名信誉/SPF/DKIM/进垃圾箱            │
  │ Push     │ 70~90%   │ 用户关闭通知/卸载APP/Token失效          │
  │ 短信      │ 95~99%   │ 运营商拦截/空号/停机                   │
  └──────────┴───────────┴──────────────────────────────────────┘

监控指标:
  ┌──────────────────────┬────────────────────────────────┐
  │ 指标                  │ 计算方式                        │
  ├──────────────────────┼────────────────────────────────┤
  │ 发送成功率            │ 成功数 / 总发送数               │
  │ 送达率               │ 已送达 / 总发送数（需回执）      │
  │ 打开率(Push/Email)   │ 用户点击数 / 送达数              │
  │ 退订率(Email/SMS)    │ 退订用户数 / 发送用户数          │
  │ 平均延迟             │ 从请求到实际送达的耗时            │
  │ 防骚扰拦截率          │ 被拦截数 / 总请求数              │
  └──────────────────────┴────────────────────────────────┘

  告警规则:
    - 短信成功率 < 95% → P1 告警
    - Push 送达率 < 70% → P2 告警（可能 Token 批量失效）
    - 单渠道 Worker 消费延迟 > 5min → P1 告警
    - 防骚扰拦截率 > 50% → 排查是否模板/频率配置有误
```

### 9. 渠道降级与路由

```
多渠道降级策略:
  场景: 短信通道故障 → 自动降级到其他渠道

  优先级路由:
    验证码:  SMS → 语音电话 → 站内信
    交易通知: Push → SMS → Email → 站内信
    营销通知: Push → 站内信（不走 SMS，节省成本）
    安全告警: SMS + Push + Email（多通道并发，确保送达）

  ┌──────────┐     ┌──────────────┐     ┌──────────────┐
  │ 发送请求  │────>│ 渠道路由器    │────>│ 主渠道(Push) │ → 成功
  └──────────┘     └──────┬───────┘     └──────────────┘
                          │ 失败
                          ▼
                   ┌──────────────┐
                   │ 备选渠道(SMS) │ → 成功
                   └──────┬───────┘
                          │ 失败
                          ▼
                   ┌──────────────┐
                   │ 兜底(站内信)  │ → 一定成功
                   └──────────────┘

  熔断机制:
    某渠道连续失败 N 次 → 自动熔断该渠道 → 流量切到备选渠道
    定期探测 → 恢复后自动切回
```

### 一句话总结

**通知中心核心架构：业务方调用统一网关（传模板编码+变量+用户列表）→ 网关做模板渲染+用户偏好过滤+防骚扰检查+幂等去重 → 按渠道写入不同 MQ Topic → 各渠道 Worker 独立消费并调用三方 API（Push/SMS/Email/站内信），失败走延迟队列重试（最多 3 次指数退避）。** 防骚扰四层策略：频率限制（每用户每天每渠道上限）+ 免打扰时段（22:00~08:00 系统通知例外）+ 同内容去重 + 营销类每周限额。渠道设计用策略模式（Sender 接口），新增渠道只需实现接口+注册。群发场景分批查询用户（每批 1000）+ 批量投递 MQ + Worker 按渠道限速消费。多渠道降级路由：主渠道失败自动切备选渠道，站内信兜底确保触达。

---

## 十八、灰度发布与 AB 测试

> 怎么实现灰度发布？流量染色？AB 测试如何分流？如何保证同一用户始终命中同一实验组？

### 1. 核心概念

```
┌─────────────────────────────────────────────────────────────────┐
│ 灰度发布 (Gray Release / Canary Release):                       │
│   新版本先让一小部分用户使用，观察无异常后再全量放开               │
│   目的: 降低发布风险                                             │
│                                                                 │
│ 蓝绿发布 (Blue-Green Deployment):                               │
│   两套完全相同的环境(蓝=旧, 绿=新)，切换流量一步到位              │
│   目的: 快速回滚                                                 │
│                                                                 │
│ AB 测试 (A/B Testing):                                          │
│   将用户随机分成两组，A组用方案A，B组用方案B，对比数据指标          │
│   目的: 数据驱动决策                                             │
│                                                                 │
│ 流量染色 (Traffic Coloring):                                     │
│   在请求中打标签，标签沿整条调用链透传，路由到对应版本             │
│   目的: 微服务架构下的全链路灰度                                  │
└─────────────────────────────────────────────────────────────────┘

关系:
  灰度发布  → 偏运维侧（新版本安全上线）
  AB 测试   → 偏产品侧（哪个方案效果好）
  流量染色  → 偏技术侧（微服务全链路灰度的实现手段）
  三者底层能力是共通的: 都需要 "分流 + 路由 + 可观测"
```

### 2. 灰度发布模型

```
┌─────────────────────────────────────────────────────────────┐
│ 模型一: 金丝雀发布 (Canary)                                  │
│                                                             │
│   100% ──────┐                                              │
│              │ 流量                                         │
│   ┌──────────▼──────────┐                                   │
│   │     负载均衡          │                                   │
│   └──┬──────────────┬───┘                                   │
│    95%             5%                                       │
│   ┌──▼──┐      ┌───▼───┐                                   │
│   │ V1  │      │ V2    │  ← 新版本(金丝雀)                  │
│   │ 10台│      │ 1台   │                                    │
│   └─────┘      └───────┘                                   │
│                                                             │
│   流程: 部署1台新版本 → 导入5%流量 → 观察指标                │
│         → 正常 → 逐步扩大(10%→30%→50%→100%)                │
│         → 异常 → 回滚(流量切回V1)                           │
│                                                             │
│   适合: 常规版本发布，按比例逐步放量                          │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 模型二: 定向灰度                                             │
│                                                             │
│   不按比例，而是按用户特征定向放量                            │
│   规则示例:                                                  │
│     - 内部员工 (user_type=staff)   → V2                     │
│     - 城市=北京                    → V2                     │
│     - 用户ID尾号=0                 → V2                     │
│     - iOS + 版本>=7.0              → V2                     │
│     - 白名单用户列表               → V2                     │
│     - 其他                         → V1                     │
│                                                             │
│   适合: 先让内部人员/种子用户体验新功能                       │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ 模型三: 蓝绿发布                                             │
│                                                             │
│   ┌──────────────────────┐                                  │
│   │     负载均衡           │                                  │
│   └──┬──────────────┬───┘                                   │
│    100%            0%                                       │
│   ┌──▼──┐      ┌───▼───┐                                   │
│   │ 蓝  │      │ 绿    │  ← 新版本已部署但无流量             │
│   │(旧) │      │(新)   │                                    │
│   └─────┘      └───────┘                                   │
│                                                             │
│   验证通过 → 流量瞬间切到绿                                  │
│   出问题   → 流量瞬间切回蓝                                  │
│                                                             │
│   优点: 回滚极快（秒级）                                     │
│   缺点: 需要双倍资源                                         │
└─────────────────────────────────────────────────────────────┘
```

### 3. 整体架构

```
┌──────────────────────────────────────────────────────────────────┐
│                          运营后台                                  │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐   │
│   │ 实验管理  │  │ 灰度规则  │  │ 分流配置  │  │ 数据看板       │   │
│   └──────────┘  └──────────┘  └──────────┘  └───────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │ 配置下发 (etcd / Apollo / Nacos)
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                   API 网关 / SDK 客户端                            │
│                                                                  │
│  请求进入 → 提取特征(userID/IP/UA/Header)                        │
│          → 匹配实验/灰度规则                                     │
│          → 计算分组（hash(userID) % 100）                        │
│          → 注入标签到 Header（X-Gray-Tag: v2）                   │
│          → 路由到对应版本                                        │
└──────────┬──────────────────────────────────────┬───────────────┘
           │                                      │
    ┌──────▼──────┐                        ┌──────▼──────┐
    │ Service V1  │                        │ Service V2  │
    │ (稳定版)    │                        │ (灰度版)    │
    └──────┬──────┘                        └──────┬──────┘
           │                                      │
           │   流量染色: Header 透传到下游服务      │
           ▼                                      ▼
    ┌─────────────┐                        ┌─────────────┐
    │ 下游服务 V1  │                        │ 下游服务 V2  │
    └─────────────┘                        └─────────────┘

配置中心存储的灰度规则:
  {
    "experiment_id": "exp_new_checkout",
    "status": "running",
    "rules": [
      {"type": "whitelist", "user_ids": [1001, 1002]},
      {"type": "city", "values": ["北京", "上海"]},
      {"type": "ratio", "percent": 10}
    ],
    "groups": {
      "control": {"weight": 50, "tag": "v1"},
      "treatment": {"weight": 50, "tag": "v2"}
    }
  }
```

### 4. 分流算法

```
核心要求:
  1. 同一用户始终命中同一组（粘性）→ 避免体验不一致
  2. 分布均匀 → A/B 组用户数接近 50:50
  3. 不同实验独立 → 实验X的分组不影响实验Y

方案: Hash 取模

  group = hash(userID + experimentID) % 100

  [0, 49]  → 对照组(A)
  [50, 99] → 实验组(B)

  为什么加 experimentID？
    如果只用 hash(userID)，那同一批用户在所有实验中都被分到同一组
    加上 experimentID → 不同实验的分组结果不同

  为什么用 hash 而不是 userID % 100？
    userID 注册有先后 → 低 ID 是老用户，高 ID 是新用户
    直接取模会导致组间用户特征不均匀
    hash 打散后分布更均匀

多层实验（正交实验）:
  ┌──────────────────────────────────────────────────┐
  │ 层1(UI实验):     hash(uid+"layer1") % 100        │
  │   [0,49]=A(旧UI)  [50,99]=B(新UI)               │
  │                                                  │
  │ 层2(推荐算法实验): hash(uid+"layer2") % 100       │
  │   [0,49]=X(算法1)  [50,99]=Y(算法2)              │
  │                                                  │
  │ 用户可能同时在: UI实验的B组 + 推荐实验的X组        │
  │ 两层实验互不影响（正交）                           │
  └──────────────────────────────────────────────────┘
```

#### Go 实现分流引擎

```go
package grayscale

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ---------- 数据结构 ----------

// Experiment 实验/灰度配置
type Experiment struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Status    string           `json:"status"` // draft/running/paused/finished
	Rules     []Rule           `json:"rules"`  // 定向规则(白名单/城市/设备等)
	Groups    map[string]Group `json:"groups"` // 实验分组
	StartAt   time.Time        `json:"start_at"`
	EndAt     time.Time        `json:"end_at"`
}

// Rule 定向规则
type Rule struct {
	Type   string   `json:"type"`   // whitelist / city / platform / app_version / ratio
	Values []string `json:"values"` // 规则参数
}

// Group 实验分组
type Group struct {
	Tag    string `json:"tag"`    // 标签(v1/v2/control/treatment)
	Weight int    `json:"weight"` // 权重(0-100)
}

// UserContext 用户请求上下文（用于规则匹配）
type UserContext struct {
	UserID     int64  `json:"user_id"`
	City       string `json:"city"`
	Platform   string `json:"platform"`    // ios/android/web
	AppVersion string `json:"app_version"`
	IP         string `json:"ip"`
	DeviceID   string `json:"device_id"`
}

// HitResult 命中结果
type HitResult struct {
	ExperimentID string `json:"experiment_id"`
	GroupTag     string `json:"group_tag"`
	Matched      bool   `json:"matched"`
}

// ---------- 分流引擎 ----------

// Router 灰度/AB分流路由器
type Router struct {
	mu          sync.RWMutex
	experiments map[string]*Experiment // 实验配置缓存
	rdb         *redis.Client
}

func NewRouter(rdb *redis.Client) *Router {
	return &Router{
		experiments: make(map[string]*Experiment),
		rdb:         rdb,
	}
}

// LoadExperiments 从配置中心加载实验配置（定期刷新）
func (r *Router) LoadExperiments(ctx context.Context) error {
	data, err := r.rdb.Get(ctx, "grayscale:experiments").Bytes()
	if err != nil {
		return err
	}

	var experiments []*Experiment
	if err := json.Unmarshal(data, &experiments); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.experiments = make(map[string]*Experiment)
	for _, exp := range experiments {
		r.experiments[exp.ID] = exp
	}
	return nil
}

// Route 计算用户在指定实验中的分组
func (r *Router) Route(ctx context.Context, experimentID string, user *UserContext) *HitResult {
	r.mu.RLock()
	exp, ok := r.experiments[experimentID]
	r.mu.RUnlock()

	if !ok || exp.Status != "running" {
		return &HitResult{ExperimentID: experimentID, Matched: false}
	}

	// 1. 检查定向规则（命中任一规则即可进入实验）
	if !r.matchRules(exp.Rules, user) {
		return &HitResult{ExperimentID: experimentID, Matched: false}
	}

	// 2. Hash 分流 → 确定分组
	groupTag := r.assignGroup(user.UserID, experimentID, exp.Groups)

	return &HitResult{
		ExperimentID: experimentID,
		GroupTag:     groupTag,
		Matched:      true,
	}
}

// RouteAll 计算用户在所有运行中实验的分组
func (r *Router) RouteAll(ctx context.Context, user *UserContext) []*HitResult {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*HitResult
	for _, exp := range r.experiments {
		if exp.Status != "running" {
			continue
		}
		result := r.Route(ctx, exp.ID, user)
		if result.Matched {
			results = append(results, result)
		}
	}
	return results
}

// matchRules 匹配定向规则
func (r *Router) matchRules(rules []Rule, user *UserContext) bool {
	if len(rules) == 0 {
		return true // 无规则限制 → 全量用户参与
	}

	for _, rule := range rules {
		switch rule.Type {
		case "whitelist":
			// 白名单: userID 在列表中
			uid := fmt.Sprintf("%d", user.UserID)
			for _, v := range rule.Values {
				if v == uid {
					return true
				}
			}
		case "city":
			for _, v := range rule.Values {
				if v == user.City {
					return true
				}
			}
		case "platform":
			for _, v := range rule.Values {
				if v == user.Platform {
					return true
				}
			}
		case "ratio":
			// 按比例放量: hash(userID) % 100 < percent
			if len(rule.Values) > 0 {
				var percent int
				fmt.Sscanf(rule.Values[0], "%d", &percent)
				bucket := hashBucket(user.UserID, "ratio")
				if bucket < percent {
					return true
				}
			}
		}
	}
	return false
}

// assignGroup Hash 分流确定分组
func (r *Router) assignGroup(userID int64, experimentID string, groups map[string]Group) string {
	bucket := hashBucket(userID, experimentID)

	// 按权重区间分组
	// 例: control=50, treatment=50 → [0,49]=control, [50,99]=treatment
	cumulative := 0
	for name, g := range groups {
		cumulative += g.Weight
		if bucket < cumulative {
			return g.Tag
			_ = name
		}
	}

	// fallback: 返回第一个分组
	for _, g := range groups {
		return g.Tag
	}
	return "default"
}

// hashBucket 将 userID + salt 映射到 [0, 100) 的桶
func hashBucket(userID int64, salt string) int {
	key := fmt.Sprintf("%d_%s", userID, salt)
	h := md5.Sum([]byte(key))
	num := binary.BigEndian.Uint32(h[:4])
	return int(num % 100)
}
```

### 5. 流量染色（微服务全链路灰度）

```
问题:
  单体应用灰度 → 网关层路由即可
  微服务灰度   → 网关路由到 Service-A v2，但 A 调用 B 时怎么也走 B v2？

  用户请求 → 网关(灰度) → Service-A v2 → Service-B ???
                                          ↑ 应该也走 v2

解决方案: 流量染色 + 标签路由

  ┌────────┐    ┌────────┐    ┌──────────┐    ┌──────────┐
  │  网关   │──>│ A (v2) │──>│ B (v2)   │──>│ C (v2)   │
  │  染色   │   │        │   │          │   │          │
  │X-Gray:v2│   │透传Header│  │透传Header│   │          │
  └────────┘    └────────┘    └──────────┘    └──────────┘
      │
      │         ┌────────┐    ┌──────────┐    ┌──────────┐
      └────────>│ A (v1) │──>│ B (v1)   │──>│ C (v1)   │
                │  正常   │   │          │   │          │
                └────────┘    └──────────┘    └──────────┘

流量染色规则:
  1. 网关层: 根据分流结果在 Header 中注入标签
     X-Gray-Tag: v2
     X-Experiment-ID: exp_new_checkout

  2. 服务间调用: 框架自动透传灰度 Header
     gRPC → metadata 透传
     HTTP → Header 透传

  3. 服务注册时打标签:
     Service-A v1: tag=stable
     Service-A v2: tag=gray

  4. 负载均衡路由:
     请求带 X-Gray-Tag:v2 → 优先路由到 tag=gray 的实例
     无灰度标签 → 路由到 tag=stable 的实例
     灰度实例不存在 → 降级到 stable（不能报错）
```

#### Go 实现流量染色中间件

```go
// ---------- 网关层: 染色中间件 ----------

const (
	HeaderGrayTag      = "X-Gray-Tag"
	HeaderExperimentID = "X-Experiment-ID"
	HeaderUserID       = "X-User-ID"
)

// GrayMiddleware 网关灰度染色中间件
func GrayMiddleware(router *Router) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 提取用户上下文
			user := extractUserContext(r)

			// 2. 计算所有命中的实验
			hits := router.RouteAll(r.Context(), user)

			// 3. 注入灰度标签到 Header
			for _, hit := range hits {
				r.Header.Set(HeaderGrayTag, hit.GroupTag)
				r.Header.Set(HeaderExperimentID, hit.ExperimentID)
			}

			// 4. 写入 Context（方便业务代码读取）
			ctx := context.WithValue(r.Context(), "gray_hits", hits)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractUserContext(r *http.Request) *UserContext {
	var userID int64
	fmt.Sscanf(r.Header.Get(HeaderUserID), "%d", &userID)

	return &UserContext{
		UserID:     userID,
		City:       r.Header.Get("X-City"),
		Platform:   r.Header.Get("X-Platform"),
		AppVersion: r.Header.Get("X-App-Version"),
		IP:         r.RemoteAddr,
	}
}

// ---------- 服务间: gRPC 拦截器透传灰度标签 ----------

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrayUnaryClientInterceptor gRPC 客户端拦截器 → 自动透传灰度标签
func GrayUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// 从当前 Context 提取灰度标签
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			grayTag := md.Get(HeaderGrayTag)
			expID := md.Get(HeaderExperimentID)
			if len(grayTag) > 0 {
				// 透传到下游
				ctx = metadata.AppendToOutgoingContext(ctx,
					HeaderGrayTag, grayTag[0],
					HeaderExperimentID, expID[0],
				)
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// GrayUnaryServerInterceptor gRPC 服务端拦截器 → 读取灰度标签注入 Context
func GrayUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			grayTag := md.Get(HeaderGrayTag)
			if len(grayTag) > 0 {
				ctx = context.WithValue(ctx, "gray_tag", grayTag[0])
			}
		}

		return handler(ctx, req)
	}
}

// ---------- HTTP 客户端: 自动透传 ----------

// GrayRoundTripper HTTP 客户端透传灰度标签
type GrayRoundTripper struct {
	Base http.RoundTripper
}

func (t *GrayRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 从 Context 中提取灰度标签并注入到请求 Header
	if tag, ok := req.Context().Value("gray_tag").(string); ok && tag != "" {
		req.Header.Set(HeaderGrayTag, tag)
	}
	if expID, ok := req.Context().Value("exp_id").(string); ok && expID != "" {
		req.Header.Set(HeaderExperimentID, expID)
	}
	return t.Base.RoundTrip(req)
}
```

### 6. 灰度与 AB 测试在业务代码中的使用

```go
// ---------- 业务代码中使用灰度开关 ----------

// IsGrayUser 判断当前请求是否命中灰度
func IsGrayUser(ctx context.Context, tag string) bool {
	hits, ok := ctx.Value("gray_hits").([]*HitResult)
	if !ok {
		return false
	}
	for _, h := range hits {
		if h.GroupTag == tag {
			return true
		}
	}
	return false
}

// 业务使用示例
func HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if IsGrayUser(r.Context(), "v2") {
		// 新版结算流程
		newCheckoutFlow(w, r)
	} else {
		// 旧版结算流程
		oldCheckoutFlow(w, r)
	}
}

// Feature Flag 功能开关（更通用的方式）
type FeatureFlags struct {
	router *Router
}

func (f *FeatureFlags) IsEnabled(ctx context.Context, featureKey string, user *UserContext) bool {
	result := f.router.Route(ctx, featureKey, user)
	return result.Matched && result.GroupTag == "enabled"
}

// 使用
// if featureFlags.IsEnabled(ctx, "new_payment_page", user) { ... }
```

### 7. K8s + Istio 灰度发布实践

```
K8s 原生灰度: 通过 Deployment + Service 实现

  # 稳定版本
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: my-app-stable
    labels:
      app: my-app
      version: v1
  spec:
    replicas: 9    # 90% 流量

  # 灰度版本
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: my-app-canary
    labels:
      app: my-app
      version: v2
  spec:
    replicas: 1    # 10% 流量

  # Service 同时匹配两个 Deployment
  apiVersion: v1
  kind: Service
  metadata:
    name: my-app
  spec:
    selector:
      app: my-app   # 只按 app 匹配，不区分 version
                     # K8s 按 Pod 数比例自动分流


Istio 精细灰度（推荐）:

  # VirtualService: 按权重/Header 分流
  apiVersion: networking.istio.io/v1beta1
  kind: VirtualService
  metadata:
    name: my-app
  spec:
    hosts:
      - my-app
    http:
      # 规则1: Header 带灰度标签 → 路由到 v2
      - match:
          - headers:
              X-Gray-Tag:
                exact: v2
        route:
          - destination:
              host: my-app
              subset: v2

      # 规则2: 默认按权重分流
      - route:
          - destination:
              host: my-app
              subset: v1
            weight: 90
          - destination:
              host: my-app
              subset: v2
            weight: 10

  # DestinationRule: 定义子集
  apiVersion: networking.istio.io/v1beta1
  kind: DestinationRule
  metadata:
    name: my-app
  spec:
    host: my-app
    subsets:
      - name: v1
        labels:
          version: v1
      - name: v2
        labels:
          version: v2

常用灰度发布工具:
  ┌──────────────┬─────────────────────────────────────────┐
  │ 工具          │ 特点                                    │
  ├──────────────┼─────────────────────────────────────────┤
  │ Istio        │ Service Mesh，VirtualService 精细路由    │
  │ Nginx Ingress│ annotations 配置金丝雀(canary-weight)   │
  │ Argo Rollouts│ K8s 渐进式发布控制器(Canary/BlueGreen)  │
  │ Flagger      │ Istio/Nginx 之上的自动金丝雀分析         │
  │ Apollo/Nacos │ 配置中心 + Feature Flag                 │
  │ LaunchDarkly │ SaaS Feature Flag 平台                  │
  └──────────────┴─────────────────────────────────────────┘
```

### 8. AB 测试数据采集与分析

```
AB 测试的数据闭环:

  ┌──────┐     ┌──────┐     ┌────────┐     ┌──────────┐
  │ 分流  │────>│ 曝光  │────>│ 行为   │────>│ 分析决策  │
  │ 确定  │     │ 埋点  │     │ 埋点   │     │          │
  │ 分组  │     │ 记录  │     │ 记录   │     │          │
  └──────┘     └──────┘     └────────┘     └──────────┘

埋点数据:
  {
    "event": "page_view",
    "user_id": 10086,
    "experiment_id": "exp_new_checkout",
    "group": "treatment",
    "timestamp": 1710000000,
    "properties": {
      "page": "checkout",
      "order_amount": 99.00
    }
  }

  → 写入 Kafka → ClickHouse/Hive → 数据分析

核心指标:
  ┌──────────────────┬────────────────────────────────────┐
  │ 指标              │ 计算方式                            │
  ├──────────────────┼────────────────────────────────────┤
  │ 转化率            │ 下单用户数 / 进入页面用户数         │
  │ 人均 GMV          │ 总交易额 / 实验组用户数             │
  │ 留存率            │ 次日回访用户数 / 当日活跃用户数     │
  │ 点击率(CTR)       │ 点击次数 / 曝光次数                 │
  │ 页面停留时长      │ 实验组平均 vs 对照组平均            │
  └──────────────────┴────────────────────────────────────┘

  统计显著性:
    用 p-value 判断差异是否显著 (通常 p < 0.05)
    样本量不够时 → 不要急着下结论 → 继续跑实验

    最小样本量估算(简化):
      每组需要 ≈ 16 × p × (1-p) / MDE²
      p = 基线转化率
      MDE = 最小可检测效应(如期望提升 2%)

灰度发布/AB 的安全检查:
  ┌─────────────────┬─────────────────────────────────────┐
  │ 检查项            │ 方法                                │
  ├─────────────────┼─────────────────────────────────────┤
  │ 错误率           │ V2 错误率不能比 V1 高 0.1%          │
  │ 延迟 P99         │ V2 延迟不能比 V1 高 20%             │
  │ 核心指标         │ 转化率/GMV 不能低于 V1               │
  │ 自动回滚         │ 触发任一告警 → 自动切回 V1           │
  └─────────────────┴─────────────────────────────────────┘
```

### 一句话总结

**灰度发布三种模型：金丝雀（按比例 5%→10%→50%→100% 逐步放量）、定向灰度（白名单/城市/设备等规则定向放量）、蓝绿（双环境瞬时切换）。分流算法用 hash(userID+experimentID)%100 保证同一用户始终命中同一组且不同实验正交独立。** 微服务全链路灰度靠流量染色：网关层根据分流结果注入 X-Gray-Tag Header → 服务间调用自动透传（gRPC metadata / HTTP Header）→ 负载均衡按标签路由到对应版本实例。K8s 实践推荐 Istio VirtualService（按 Header/权重精细路由）或 Argo Rollouts（自动金丝雀分析+回滚）。AB 测试核心闭环：分流→曝光埋点→行为埋点→ClickHouse 分析→p-value 检验显著性→决策上线或回滚。

---

## 十九、多方式登录验证系统设计

> 如何设计支持用户名密码、邮箱验证码、手机号短信验证码、扫码登录的统一登录系统？

### 1. 核心需求与挑战

```
四种登录方式:
  ┌──────────────────────────────────────────────────────────────┐
  │ ① 用户名 + 密码      最传统，需安全存储密码                  │
  │ ② 邮箱 + 验证码      免密登录，验证码邮件                    │
  │ ③ 手机号 + 短信验证码  国内主流，验证码短信                   │
  │ ④ 扫码登录            PC端展示二维码，APP端扫码确认            │
  └──────────────────────────────────────────────────────────────┘

设计挑战:
  ┌──────────────────────────────────────────────────────────────┐
  │ 1. 统一抽象:  四种方式最终都要走到"确认用户身份→签发Token"    │
  │ 2. 账号绑定:  同一用户可能绑定了手机号+邮箱+用户名            │
  │ 3. 密码安全:  密码不能明文存储（bcrypt），防暴力破解           │
  │ 4. 验证码安全: 防刷（频率限制）、防爆破（错误次数限制）        │
  │ 5. 扫码安全:  二维码时效性、确认前后状态同步（轮询/长连接）   │
  │ 6. 安全防护:  图形验证码（人机验证）、IP限流、异地登录告警     │
  └──────────────────────────────────────────────────────────────┘
```

### 2. 整体架构

```
┌────────────────────────────────────────────────────────────────────┐
│                          客户端                                     │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌───────────────┐   │
│  │用户名+密码  │ │邮箱+验证码  │ │手机+验证码  │ │ 扫码登录       │   │
│  └─────┬──────┘ └─────┬──────┘ └─────┬──────┘ └──────┬────────┘   │
└────────┼──────────────┼──────────────┼───────────────┼─────────────┘
         │              │              │               │
         ▼              ▼              ▼               ▼
┌────────────────────────────────────────────────────────────────────┐
│                      API 网关 (限流 + 人机验证)                      │
└───────────────────────────┬────────────────────────────────────────┘
                            │
┌───────────────────────────▼────────────────────────────────────────┐
│                       认证服务 (Auth Service)                       │
│                                                                    │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │              LoginHandler (统一入口)                        │    │
│  │                                                            │    │
│  │  login_type = "password" / "email_code" / "sms_code" /    │    │
│  │               "qrcode"                                     │    │
│  │                                                            │    │
│  │  → 路由到对应 Authenticator                                │    │
│  └──────┬──────────┬──────────┬──────────┬────────────────────┘    │
│         │          │          │          │                          │
│  ┌──────▼──┐ ┌─────▼────┐ ┌──▼───────┐ ┌▼──────────┐             │
│  │Password │ │EmailCode │ │SMSCode   │ │QRCode    │              │
│  │Authen.  │ │Authen.   │ │Authen.   │ │Authen.   │              │
│  └──────┬──┘ └─────┬────┘ └──┬───────┘ └┬──────────┘             │
│         │          │          │          │                          │
│         └──────────┴──────────┴──────────┘                         │
│                        │                                            │
│                ┌───────▼───────┐                                    │
│                │ 认证成功       │                                    │
│                │ → 签发 JWT    │                                    │
│                │ → 记录登录日志 │                                    │
│                └───────────────┘                                    │
└────────────────────────────────────────────────────────────────────┘
         │                │                │
  ┌──────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
  │   MySQL     │  │   Redis     │  │ 短信/邮件    │
  │ 用户+账号表  │  │ 验证码      │  │ 通知服务     │
  │ 登录日志    │  │ 扫码状态    │  │              │
  └─────────────┘  │ 频率限制    │  └──────────────┘
                   └─────────────┘
```

### 3. 数据库设计

```sql
-- 用户主表（核心身份信息）
CREATE TABLE user (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    nickname    VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '昵称',
    avatar      VARCHAR(256) NOT NULL DEFAULT '' COMMENT '头像',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1正常 2禁用 3注销',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='用户表';

-- 用户认证方式表（一个用户可绑定多种登录方式）
-- 核心设计: 账号和认证方式解耦，方便扩展新登录方式
CREATE TABLE user_auth (
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id       BIGINT       NOT NULL COMMENT '用户ID',
    identity_type VARCHAR(16)  NOT NULL COMMENT '认证类型: username/phone/email/wechat/github',
    identifier    VARCHAR(128) NOT NULL COMMENT '标识: 用户名/手机号/邮箱/openid',
    credential    VARCHAR(256) NOT NULL DEFAULT '' COMMENT '凭证: 密码hash(仅password类型)',
    verified      TINYINT      NOT NULL DEFAULT 0 COMMENT '是否已验证(手机/邮箱)',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_type_identifier (identity_type, identifier), -- 同类型标识唯一
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='用户认证方式';

-- 示例数据:
-- user_id=1001 绑定了3种登录方式:
-- | user_id | identity_type | identifier          | credential        |
-- | 1001    | username      | zhangsan             | $2a$10$xxxx(bcrypt) |
-- | 1001    | phone         | 13800138000          |                   |
-- | 1001    | email         | zhangsan@example.com |                   |

-- 登录日志表（安全审计）
CREATE TABLE login_log (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    login_type  VARCHAR(16) NOT NULL COMMENT '登录方式',
    ip          VARCHAR(64) NOT NULL,
    user_agent  VARCHAR(256) NOT NULL DEFAULT '',
    location    VARCHAR(64) NULL COMMENT 'IP归属地',
    success     TINYINT     NOT NULL COMMENT '1成功 0失败',
    fail_reason VARCHAR(128) NULL COMMENT '失败原因',
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user (user_id, created_at)
) ENGINE=InnoDB COMMENT='登录日志';
```

### 4. 策略模式 — 核心接口与路由

```go
package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// ---------- 核心接口 ----------

// Authenticator 认证器接口（策略模式）
type Authenticator interface {
	// Authenticate 验证用户身份，返回 userID
	Authenticate(ctx context.Context, req *LoginRequest) (int64, error)
	// Type 返回认证类型标识
	Type() string
}

// LoginRequest 统一登录请求
type LoginRequest struct {
	LoginType string `json:"login_type"` // password / sms_code / email_code / qrcode
	// 用户名+密码
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	// 手机号/邮箱 + 验证码
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	Code  string `json:"code,omitempty"`
	// 扫码
	QRToken string `json:"qr_token,omitempty"`
	// 通用
	CaptchaID  string `json:"captcha_id,omitempty"`  // 图形验证码ID
	CaptchaVal string `json:"captcha_val,omitempty"` // 图形验证码值
	IP         string `json:"-"`
	UserAgent  string `json:"-"`
}

// LoginResponse 统一登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	Nickname     string `json:"nickname"`
}

// ---------- 认证路由器 ----------

// AuthRouter 认证路由器（根据 login_type 分发到对应 Authenticator）
type AuthRouter struct {
	authenticators map[string]Authenticator
	db             *sql.DB
	rdb            *redis.Client
	jwtMgr         JWTManager // JWT签发（参考 SSO 章节实现）
	rateLimiter    *RateLimiter
}

type JWTManager interface {
	GenerateTokenPair(ctx context.Context, userID int64, username, role string) (*TokenPair, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func NewAuthRouter(db *sql.DB, rdb *redis.Client, jwtMgr JWTManager) *AuthRouter {
	r := &AuthRouter{
		authenticators: make(map[string]Authenticator),
		db:             db,
		rdb:            rdb,
		jwtMgr:         jwtMgr,
		rateLimiter:    NewRateLimiter(rdb),
	}
	// 注册认证器
	r.Register(&PasswordAuth{db: db})
	r.Register(&SMSCodeAuth{db: db, rdb: rdb})
	r.Register(&EmailCodeAuth{db: db, rdb: rdb})
	r.Register(&QRCodeAuth{db: db, rdb: rdb})
	return r
}

func (r *AuthRouter) Register(a Authenticator) {
	r.authenticators[a.Type()] = a
}

// Login 统一登录入口
func (r *AuthRouter) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. IP 限流（同一 IP 每分钟最多 20 次登录尝试）
	if blocked := r.rateLimiter.CheckIP(ctx, req.IP, 20, time.Minute); blocked {
		return nil, errors.New("登录请求过于频繁，请稍后再试")
	}

	// 2. 路由到对应认证器
	auth, ok := r.authenticators[req.LoginType]
	if !ok {
		return nil, fmt.Errorf("不支持的登录方式: %s", req.LoginType)
	}

	// 3. 执行认证
	userID, err := auth.Authenticate(ctx, req)
	if err != nil {
		// 记录失败日志
		r.logLogin(ctx, 0, req, false, err.Error())
		return nil, err
	}

	// 4. 检查用户状态
	var nickname string
	var status int
	r.db.QueryRowContext(ctx,
		"SELECT nickname, status FROM user WHERE id = ?", userID,
	).Scan(&nickname, &status)
	if status != 1 {
		return nil, errors.New("账号已被禁用或注销")
	}

	// 5. 签发 JWT
	tokenPair, err := r.jwtMgr.GenerateTokenPair(ctx, userID, nickname, "user")
	if err != nil {
		return nil, fmt.Errorf("签发Token失败: %w", err)
	}

	// 6. 记录成功日志
	r.logLogin(ctx, userID, req, true, "")

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		UserID:       userID,
		Nickname:     nickname,
	}, nil
}

func (r *AuthRouter) logLogin(ctx context.Context, userID int64, req *LoginRequest, success bool, reason string) {
	s := 0
	if success {
		s = 1
	}
	r.db.ExecContext(ctx,
		`INSERT INTO login_log (user_id, login_type, ip, user_agent, success, fail_reason) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, req.LoginType, req.IP, req.UserAgent, s, reason,
	)
}
```

### 5. 方式一：用户名 + 密码登录

```
流程:
  ┌────────┐              ┌──────────┐              ┌────────┐
  │  客户端  │              │ 认证服务   │              │ MySQL  │
  └───┬────┘              └────┬─────┘              └───┬────┘
      │ 1.POST /login          │                       │
      │  {type:"password",     │                       │
      │   username:"zhangsan", │                       │
      │   password:"123456"}   │                       │
      │───────────────────────>│                       │
      │                       │ 2.查 user_auth 表      │
      │                       │  WHERE identity_type=  │
      │                       │  'username' AND        │
      │                       │  identifier='zhangsan' │
      │                       │──────────────────────>│
      │                       │ 3.返回 credential      │
      │                       │  (bcrypt hash)         │
      │                       │<──────────────────────│
      │                       │ 4.bcrypt.Compare       │
      │                       │  (hash, "123456")     │
      │                       │ 5.匹配 → 签发JWT       │
      │ 6.返回 Token           │                       │
      │<───────────────────────│                       │

密码安全:
  存储: bcrypt (自带盐值, cost=10~12)
    hash = bcrypt.GenerateFromPassword("123456", 10)
    → "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

  ❌ 错误: MD5/SHA256（无盐、可彩虹表破解）
  ❌ 错误: MD5 + 固定盐（盐泄露后仍可批量破解）
  ✅ 正确: bcrypt/scrypt/argon2（自适应 cost，抗 GPU 暴力破解）
```

```go
// ---------- 用户名+密码认证器 ----------

type PasswordAuth struct {
	db *sql.DB
}

func (a *PasswordAuth) Type() string { return "password" }

func (a *PasswordAuth) Authenticate(ctx context.Context, req *LoginRequest) (int64, error) {
	if req.Username == "" || req.Password == "" {
		return 0, errors.New("用户名和密码不能为空")
	}

	// 1. 查询用户认证记录
	var userID int64
	var credential string
	err := a.db.QueryRowContext(ctx,
		"SELECT user_id, credential FROM user_auth WHERE identity_type = 'username' AND identifier = ?",
		req.Username,
	).Scan(&userID, &credential)
	if err != nil {
		return 0, errors.New("用户名或密码错误") // 统一提示，不暴露用户是否存在
	}

	// 2. bcrypt 比对密码
	if err := bcrypt.CompareHashAndPassword([]byte(credential), []byte(req.Password)); err != nil {
		return 0, errors.New("用户名或密码错误")
	}

	return userID, nil
}

// Register 注册时密码加密存储
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // cost=10
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
```

### 6. 方式二/三：手机短信验证码 / 邮箱验证码

```
两种验证码登录流程一致，只是发送渠道不同:

  ┌────────┐         ┌──────────┐       ┌────────┐     ┌──────────┐
  │  客户端  │         │ 认证服务   │       │ Redis  │     │ 短信/邮件 │
  └───┬────┘         └────┬─────┘       └───┬────┘     └────┬─────┘
      │ 1.请求验证码        │                │               │
      │  POST /code/send   │                │               │
      │  {phone/email}     │                │               │
      │───────────────────>│                │               │
      │                   │ 2.频率检查       │               │
      │                   │  同一号码60s内   │               │
      │                   │  不能重复发送    │               │
      │                   │────────────────>│               │
      │                   │ 3.生成6位随机码  │               │
      │                   │  存Redis(5min)  │               │
      │                   │────────────────>│               │
      │                   │ 4.发送验证码     │               │
      │                   │────────────────────────────────>│
      │ 5.返回"已发送"      │                │               │
      │<───────────────────│                │               │
      │                   │                │               │
      │ 6.POST /login      │                │               │
      │  {type:"sms_code", │                │               │
      │   phone:"138xxx",  │                │               │
      │   code:"382916"}   │                │               │
      │───────────────────>│                │               │
      │                   │ 7.从Redis取验证码│               │
      │                   │────────────────>│               │
      │                   │ 8.比对验证码     │               │
      │                   │ 9.验证成功→删码  │               │
      │                   │────────────────>│               │
      │                   │ 10.查user_auth  │               │
      │                   │  查手机号对应uid │               │
      │                   │ 11.不存在→自动注册│              │
      │ 12.返回 Token       │                │               │
      │<───────────────────│                │               │

验证码安全策略:
  ┌────────────────────────┬────────────────────────────────┐
  │ 策略                    │ 实现方式                        │
  ├────────────────────────┼────────────────────────────────┤
  │ 发送频率限制            │ 同号码 60s 内只能发一次          │
  │ 每日发送上限            │ 同号码每天最多 10 条             │
  │ IP 发送限制             │ 同 IP 每小时最多 20 条           │
  │ 验证码有效期            │ 5 分钟自动过期                  │
  │ 错误次数限制            │ 同号码连续错 5 次 → 锁定 30min  │
  │ 验证后立即删除          │ 防止验证码重复使用               │
  │ 图形验证码前置          │ 发验证码前先通过人机验证          │
  └────────────────────────┴────────────────────────────────┘
```

```go
// ---------- 验证码服务（短信和邮件共用） ----------

type CodeService struct {
	rdb *redis.Client
}

func NewCodeService(rdb *redis.Client) *CodeService {
	return &CodeService{rdb: rdb}
}

// SendCode 发送验证码
func (s *CodeService) SendCode(ctx context.Context, target, purpose string) (string, error) {
	// target = 手机号或邮箱
	// purpose = "login" / "register" / "reset_password"

	// 1. 60秒发送间隔限制
	intervalKey := fmt.Sprintf("code:interval:%s:%s", purpose, target)
	if s.rdb.Exists(ctx, intervalKey).Val() > 0 {
		return "", errors.New("验证码已发送，请60秒后重试")
	}

	// 2. 每日发送上限
	dailyKey := fmt.Sprintf("code:daily:%s:%s", purpose, target)
	count, _ := s.rdb.Get(ctx, dailyKey).Int()
	if count >= 10 {
		return "", errors.New("今日验证码发送次数已达上限")
	}

	// 3. 生成6位验证码
	code := generateCode(6)

	// 4. 存入 Redis（5分钟有效）
	codeKey := fmt.Sprintf("code:val:%s:%s", purpose, target)
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, codeKey, code, 5*time.Minute)
	pipe.Set(ctx, intervalKey, 1, 60*time.Second) // 60s 发送间隔
	pipe.Incr(ctx, dailyKey)
	pipe.Expire(ctx, dailyKey, 24*time.Hour)
	pipe.Exec(ctx)

	return code, nil
}

// VerifyCode 校验验证码
func (s *CodeService) VerifyCode(ctx context.Context, target, purpose, inputCode string) error {
	// 1. 检查错误次数（防爆破）
	errKey := fmt.Sprintf("code:err:%s:%s", purpose, target)
	errCount, _ := s.rdb.Get(ctx, errKey).Int()
	if errCount >= 5 {
		return errors.New("验证码错误次数过多，请30分钟后再试")
	}

	// 2. 获取正确验证码
	codeKey := fmt.Sprintf("code:val:%s:%s", purpose, target)
	storedCode, err := s.rdb.Get(ctx, codeKey).Result()
	if err != nil {
		return errors.New("验证码已过期，请重新获取")
	}

	// 3. 比对
	if storedCode != inputCode {
		// 错误计数+1
		pipe := s.rdb.Pipeline()
		pipe.Incr(ctx, errKey)
		pipe.Expire(ctx, errKey, 30*time.Minute)
		pipe.Exec(ctx)
		return errors.New("验证码错误")
	}

	// 4. 验证成功 → 立即删除（一次性使用）
	s.rdb.Del(ctx, codeKey, errKey)
	return nil
}

func generateCode(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	code := ""
	for _, v := range b {
		code += fmt.Sprintf("%d", v%10)
	}
	return code
}

// ---------- 短信验证码认证器 ----------

type SMSCodeAuth struct {
	db  *sql.DB
	rdb *redis.Client
}

func (a *SMSCodeAuth) Type() string { return "sms_code" }

func (a *SMSCodeAuth) Authenticate(ctx context.Context, req *LoginRequest) (int64, error) {
	if req.Phone == "" || req.Code == "" {
		return 0, errors.New("手机号和验证码不能为空")
	}

	// 1. 校验验证码
	codeSvc := NewCodeService(a.rdb)
	if err := codeSvc.VerifyCode(ctx, req.Phone, "login", req.Code); err != nil {
		return 0, err
	}

	// 2. 查询用户（手机号 → userID）
	var userID int64
	err := a.db.QueryRowContext(ctx,
		"SELECT user_id FROM user_auth WHERE identity_type = 'phone' AND identifier = ?",
		req.Phone,
	).Scan(&userID)

	if err != nil {
		// 3. 用户不存在 → 自动注册（手机号验证码登录常见的设计）
		userID, err = a.autoRegister(ctx, req.Phone)
		if err != nil {
			return 0, fmt.Errorf("自动注册失败: %w", err)
		}
	}

	return userID, nil
}

func (a *SMSCodeAuth) autoRegister(ctx context.Context, phone string) (int64, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 创建用户
	result, err := tx.ExecContext(ctx,
		"INSERT INTO user (nickname, status) VALUES (?, 1)",
		"用户"+phone[len(phone)-4:], // 默认昵称: 用户XXXX
	)
	if err != nil {
		return 0, err
	}
	userID, _ := result.LastInsertId()

	// 绑定手机号
	_, err = tx.ExecContext(ctx,
		"INSERT INTO user_auth (user_id, identity_type, identifier, verified) VALUES (?, 'phone', ?, 1)",
		userID, phone,
	)
	if err != nil {
		return 0, err
	}

	return userID, tx.Commit()
}

// ---------- 邮箱验证码认证器（与短信几乎一样） ----------

type EmailCodeAuth struct {
	db  *sql.DB
	rdb *redis.Client
}

func (a *EmailCodeAuth) Type() string { return "email_code" }

func (a *EmailCodeAuth) Authenticate(ctx context.Context, req *LoginRequest) (int64, error) {
	if req.Email == "" || req.Code == "" {
		return 0, errors.New("邮箱和验证码不能为空")
	}

	codeSvc := NewCodeService(a.rdb)
	if err := codeSvc.VerifyCode(ctx, req.Email, "login", req.Code); err != nil {
		return 0, err
	}

	var userID int64
	err := a.db.QueryRowContext(ctx,
		"SELECT user_id FROM user_auth WHERE identity_type = 'email' AND identifier = ?",
		req.Email,
	).Scan(&userID)
	if err != nil {
		return 0, errors.New("该邮箱未注册")
	}

	return userID, nil
}
```

### 7. 方式四：扫码登录

```
扫码登录原理（PC端展示二维码，手机APP端扫码确认）:

  ┌────────┐     ┌──────────┐     ┌────────┐     ┌────────┐
  │ PC浏览器│     │ 认证服务   │     │ Redis  │     │ 手机APP │
  └───┬────┘     └────┬─────┘     └───┬────┘     └───┬────┘
      │               │               │               │
  ─── 阶段1: 生成二维码 ──────────────────────────────────
      │ 1.GET /qr/generate │              │               │
      │──────────────────>│               │               │
      │                   │ 2.生成唯一     │               │
      │                   │   qr_token    │               │
      │                   │ 3.存Redis      │               │
      │                   │  {status:      │               │
      │                   │   "pending"}   │               │
      │                   │──────────────>│               │
      │ 4.返回二维码URL     │               │               │
      │  (含qr_token)     │               │               │
      │<──────────────────│               │               │
      │ 5.展示二维码        │               │               │
      │ 6.开始轮询状态      │               │               │
      │   GET /qr/check   │               │               │
      │   ?token=xxx       │               │               │
      │──────────────────>│               │               │
      │                   │ 7.查Redis      │               │
      │                   │──────────────>│               │
      │ 8.返回"pending"    │               │               │
      │<──────────────────│               │               │
      │                   │               │               │
  ─── 阶段2: 手机扫码 ───────────────────────────────────
      │                   │               │ 9.扫描二维码    │
      │                   │               │  解析qr_token  │
      │                   │               │               │
      │                   │ 10.POST /qr/scan              │
      │                   │  {qr_token, auth_token}       │
      │                   │<──────────────────────────────│
      │                   │  (auth_token=APP端的登录态)    │
      │                   │ 11.验证APP登录态 │               │
      │                   │ 12.更新Redis    │               │
      │                   │  status→"scanned"              │
      │                   │  user_id→1001  │               │
      │                   │──────────────>│               │
      │                   │ 13.返回"待确认" │               │
      │                   │──────────────────────────────>│
      │                   │               │               │
      │ 14.轮询返回        │               │               │
      │  "scanned"(已扫码) │               │               │
      │  PC端显示头像+     │               │               │
      │  "等待确认"        │               │               │
      │<──────────────────│               │               │
      │                   │               │               │
  ─── 阶段3: 确认登录 ───────────────────────────────────
      │                   │               │ 15.点击"确认   │
      │                   │               │    登录"       │
      │                   │ 16.POST /qr/confirm            │
      │                   │  {qr_token}   │               │
      │                   │<──────────────────────────────│
      │                   │ 17.更新Redis   │               │
      │                   │  status→"confirmed"            │
      │                   │──────────────>│               │
      │                   │ 18.返回"成功"  │               │
      │                   │──────────────────────────────>│
      │                   │               │               │
      │ 19.轮询返回        │               │               │
      │  "confirmed"       │               │               │
      │                   │ 20.签发JWT     │               │
      │ 21.返回Token       │               │               │
      │<──────────────────│               │               │

二维码状态机:
  pending ──→ scanned ──→ confirmed ──→ (JWT已签发)
     │            │           │
     │            │           └── 用户取消 → cancelled
     │            └── 超时 → expired
     └── 超时(5分钟) → expired

轮询 vs WebSocket:
  ┌────────────┬──────────────────────────────────────┐
  │ 轮询(简单)  │ PC端每 2s 调一次 GET /qr/check       │
  │            │ ✅ 实现简单   ❌ 请求多，有延迟        │
  ├────────────┼──────────────────────────────────────┤
  │ 长轮询     │ 服务端 hold 请求直到状态变更或超时     │
  │            │ ✅ 延迟低     ❌ 占连接               │
  ├────────────┼──────────────────────────────────────┤
  │ WebSocket  │ PC端建立WS连接，服务端主动推送状态变更 │
  │ (推荐)     │ ✅ 实时       ✅ 省流量               │
  └────────────┴──────────────────────────────────────┘
```

```go
// ---------- 扫码登录实现 ----------

type QRCodeAuth struct {
	db  *sql.DB
	rdb *redis.Client
}

func (a *QRCodeAuth) Type() string { return "qrcode" }

const qrTTL = 5 * time.Minute // 二维码有效期

// QRState 二维码状态
type QRState struct {
	Status string `json:"status"` // pending/scanned/confirmed/cancelled/expired
	UserID int64  `json:"user_id,omitempty"`
}

// GenerateQR 生成二维码 Token（PC端调用）
func (a *QRCodeAuth) GenerateQR(ctx context.Context) (string, error) {
	token := generateQRToken()
	state := &QRState{Status: "pending"}
	data, _ := json.Marshal(state)

	key := "qr:login:" + token
	if err := a.rdb.Set(ctx, key, data, qrTTL).Err(); err != nil {
		return "", err
	}
	return token, nil
}

// Scan 手机端扫码（标记已扫描 + 记录userID）
func (a *QRCodeAuth) Scan(ctx context.Context, qrToken string, userID int64) error {
	key := "qr:login:" + qrToken
	state, err := a.getState(ctx, key)
	if err != nil {
		return errors.New("二维码已过期")
	}
	if state.Status != "pending" {
		return fmt.Errorf("二维码状态异常: %s", state.Status)
	}

	state.Status = "scanned"
	state.UserID = userID
	data, _ := json.Marshal(state)
	a.rdb.Set(ctx, key, data, qrTTL)
	return nil
}

// Confirm 手机端确认登录
func (a *QRCodeAuth) Confirm(ctx context.Context, qrToken string, userID int64) error {
	key := "qr:login:" + qrToken
	state, err := a.getState(ctx, key)
	if err != nil {
		return errors.New("二维码已过期")
	}
	if state.Status != "scanned" || state.UserID != userID {
		return errors.New("操作无效")
	}

	state.Status = "confirmed"
	data, _ := json.Marshal(state)
	a.rdb.Set(ctx, key, data, qrTTL)
	return nil
}

// Cancel 手机端取消登录
func (a *QRCodeAuth) Cancel(ctx context.Context, qrToken string) error {
	key := "qr:login:" + qrToken
	state, _ := a.getState(ctx, key)
	if state != nil {
		state.Status = "cancelled"
		data, _ := json.Marshal(state)
		a.rdb.Set(ctx, key, data, qrTTL)
	}
	return nil
}

// CheckStatus PC端轮询状态
func (a *QRCodeAuth) CheckStatus(ctx context.Context, qrToken string) (*QRState, error) {
	key := "qr:login:" + qrToken
	return a.getState(ctx, key)
}

// Authenticate 扫码认证（PC端轮询到 confirmed 后调用）
func (a *QRCodeAuth) Authenticate(ctx context.Context, req *LoginRequest) (int64, error) {
	if req.QRToken == "" {
		return 0, errors.New("缺少扫码Token")
	}

	key := "qr:login:" + req.QRToken
	state, err := a.getState(ctx, key)
	if err != nil {
		return 0, errors.New("二维码已过期")
	}

	switch state.Status {
	case "confirmed":
		// 确认登录 → 删除二维码状态 → 返回 userID
		a.rdb.Del(ctx, key)
		return state.UserID, nil
	case "pending":
		return 0, errors.New("等待扫码")
	case "scanned":
		return 0, errors.New("已扫码，等待确认")
	case "cancelled":
		a.rdb.Del(ctx, key)
		return 0, errors.New("登录已取消")
	default:
		return 0, errors.New("二维码状态异常")
	}
}

func (a *QRCodeAuth) getState(ctx context.Context, key string) (*QRState, error) {
	data, err := a.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var state QRState
	json.Unmarshal(data, &state)
	return &state, nil
}

func generateQRToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return "QR-" + hex.EncodeToString(b)
}

// 需要引入 encoding/json
import "encoding/json"
```

### 8. 频率限制与安全防护

```go
// ---------- 频率限制器 ----------

type RateLimiter struct {
	rdb *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{rdb: rdb}
}

// CheckIP IP 级限流
func (r *RateLimiter) CheckIP(ctx context.Context, ip string, limit int, window time.Duration) bool {
	key := fmt.Sprintf("rate:login:ip:%s", ip)
	count, _ := r.rdb.Incr(ctx, key).Result()
	if count == 1 {
		r.rdb.Expire(ctx, key, window)
	}
	return count > int64(limit)
}

// CheckAccount 账号级限流（防暴力破解密码）
func (r *RateLimiter) CheckAccount(ctx context.Context, account string) (bool, int) {
	key := fmt.Sprintf("rate:login:account:%s", account)
	count, _ := r.rdb.Incr(ctx, key).Result()
	if count == 1 {
		r.rdb.Expire(ctx, key, 30*time.Minute)
	}
	if count > 5 {
		return true, int(count) // 锁定30分钟
	}
	return false, int(count)
}

// ResetAccount 登录成功后重置错误计数
func (r *RateLimiter) ResetAccount(ctx context.Context, account string) {
	r.rdb.Del(ctx, fmt.Sprintf("rate:login:account:%s", account))
}
```

```
安全防护全景:

  ┌──────────────────────────────────────────────────────────────────┐
  │ 层级        │ 措施                                               │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 接入层      │ HTTPS 传输加密（密码不可明文传输）                   │
  │             │ WAF 防 SQL 注入 / XSS                              │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 人机验证    │ 图形验证码 / 滑块验证 / 行为验证(极验/腾讯天御)     │
  │             │ 触发条件: 连续错误3次后 / 发送验证码前               │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 频率限制    │ IP级: 同IP每分钟≤20次                               │
  │             │ 账号级: 同账号连错5次→锁定30min                     │
  │             │ 验证码: 同号码60s间隔，每天≤10条                    │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 密码安全    │ 存储: bcrypt(cost≥10)，禁止MD5/SHA                 │
  │             │ 传输: HTTPS + 前端RSA加密(可选)                     │
  │             │ 强度: 至少8位+大小写+数字                           │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 异常检测    │ 异地登录 → 短信/邮件通知                           │
  │             │ 新设备登录 → 要求二次验证                           │
  │             │ 同一账号多地同时登录 → 踢出旧设备                   │
  ├─────────────┼────────────────────────────────────────────────────┤
  │ 审计日志    │ 每次登录记录: IP/设备/UA/时间/结果                  │
  │             │ 异常行为告警: 短时间大量失败登录                    │
  └─────────────┴────────────────────────────────────────────────────┘
```

### 9. HTTP 路由注册

```go
func SetupAuthRoutes(mux *http.ServeMux, router *AuthRouter, qrAuth *QRCodeAuth, codeSvc *CodeService) {
	// 统一登录
	mux.HandleFunc("POST /api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		json.NewDecoder(r.Body).Decode(&req)
		req.IP = r.RemoteAddr
		req.UserAgent = r.UserAgent()

		resp, err := router.Login(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(resp)
	})

	// 发送短信验证码
	mux.HandleFunc("POST /api/auth/sms/send", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Phone string `json:"phone"` }
		json.NewDecoder(r.Body).Decode(&req)
		code, err := codeSvc.SendCode(r.Context(), req.Phone, "login")
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		// 实际生产: 调用短信网关发送 code，此处仅示意
		_ = code // smsSender.Send(req.Phone, code)
		json.NewEncoder(w).Encode(map[string]string{"msg": "验证码已发送"})
	})

	// 发送邮箱验证码
	mux.HandleFunc("POST /api/auth/email/send", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email string `json:"email"` }
		json.NewDecoder(r.Body).Decode(&req)
		code, err := codeSvc.SendCode(r.Context(), req.Email, "login")
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = code // emailSender.Send(req.Email, code)
		json.NewEncoder(w).Encode(map[string]string{"msg": "验证码已发送"})
	})

	// 生成扫码二维码
	mux.HandleFunc("GET /api/auth/qr/generate", func(w http.ResponseWriter, r *http.Request) {
		token, err := qrAuth.GenerateQR(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"qr_token": token,
			"qr_url":   "https://example.com/qr/login?token=" + token,
		})
	})

	// PC端轮询扫码状态
	mux.HandleFunc("GET /api/auth/qr/check", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		state, err := qrAuth.CheckStatus(r.Context(), token)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "expired"})
			return
		}
		json.NewEncoder(w).Encode(state)
	})

	// APP端扫码
	mux.HandleFunc("POST /api/auth/qr/scan", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			QRToken string `json:"qr_token"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		userID := getUserIDFromToken(r) // 从APP的JWT中获取
		if err := qrAuth.Scan(r.Context(), req.QRToken, userID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"msg": "扫码成功，请在手机上确认"})
	})

	// APP端确认登录
	mux.HandleFunc("POST /api/auth/qr/confirm", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			QRToken string `json:"qr_token"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		userID := getUserIDFromToken(r)
		if err := qrAuth.Confirm(r.Context(), req.QRToken, userID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"msg": "登录成功"})
	})
}

func getUserIDFromToken(r *http.Request) int64 {
	// 从请求的 JWT 中提取 userID（APP已登录）
	return 0
}
```

### 10. 扩展：第三方登录（微信/GitHub）

```
第三方登录也可以通过 user_auth 表和 Authenticator 接口无缝扩展:

  新增认证方式:
    identity_type = "wechat"  identifier = openid
    identity_type = "github"  identifier = github_user_id

  流程:
    1. 客户端跳转到第三方授权页
    2. 用户授权后回调携带 code
    3. 服务端用 code 换 access_token → 获取用户信息(openid等)
    4. 查 user_auth 表:
       - 存在 → 返回 userID → 签发 JWT
       - 不存在 → 自动注册(创建 user + user_auth 记录)

  实现: 新增 WechatAuth / GithubAuth 实现 Authenticator 接口
       → 注册到 AuthRouter → 自动生效

  ┌─────────────────────────────────────────────────────────┐
  │ user_auth 表统一管理所有登录方式:                         │
  │                                                         │
  │ identity_type │ identifier               │ credential   │
  │ username      │ zhangsan                 │ bcrypt(pwd)  │
  │ phone         │ 13800138000              │              │
  │ email         │ zs@example.com           │              │
  │ wechat        │ oXk8e5Nf...              │              │
  │ github        │ 12345678                 │              │
  │ apple         │ 001234.abcdef...         │              │
  │                                                         │
  │ → 同一个 user_id 可绑定任意多种登录方式                  │
  │ → 新增第三方只需加一行 identity_type                     │
  │ → 不需要改 user 表结构                                   │
  └─────────────────────────────────────────────────────────┘
```

### 一句话总结

**多方式登录系统的核心设计是 user + user_auth 两表分离（用户主表存身份信息，user_auth 表按 identity_type+identifier 唯一键存储多种登录凭证，同一用户可绑定用户名/手机/邮箱/微信等任意多种方式）+ 策略模式（Authenticator 接口统一抽象四种认证方式，AuthRouter 根据 login_type 路由分发）。** 密码登录用 bcrypt 存储+比对（禁止 MD5/SHA）；验证码登录用 Redis 存验证码（5min 过期+一次性使用+错误 5 次锁定 30min+60s 发送间隔+每日上限 10 条）；扫码登录用 Redis 维护二维码状态机（pending→scanned→confirmed，PC 端轮询/WebSocket 监听状态变更，APP 端分两步：扫码+确认）。安全防护：HTTPS 传输+图形验证码人机校验+IP/账号级频率限制+异地登录检测+登录审计日志。扩展第三方登录只需新增 Authenticator 实现+注册，user_auth 表加一行 identity_type 即可。


