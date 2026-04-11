# MySQL Q&A

## 目录

- [一、MySQL 事务及实现原理](#一mysql-事务及实现原理)
- [二、Redo Log、Undo Log、Binlog 的作用及区别](#二redo-logundo-logbinlog-的作用及区别)
- [三、MVCC 的实现原理](#三mvcc-的实现原理)
- [四、MySQL 中的锁有哪几种及其作用](#四mysql-中的锁有哪几种及其作用)
- [五、OLTP 与 OLAP 数据库及存储引擎原理](#五oltp-与-olap-数据库及存储引擎原理)
- [六、什么是索引下推？作用是什么？](#六什么是索引下推作用是什么)
- [七、MySQL 索引实现原理、B 树与 B+ 树区别、索引使用注意事项](#七mysql-索引实现原理b-树与-b-树区别索引使用注意事项)
- [八、一条 SQL 的执行流程](#八一条-sql-的执行流程)
- [九、慢查询排查与优化](#九慢查询排查与优化)
- [十、分库分表方案与常见问题](#十分库分表方案与常见问题)
- [十一、主从复制原理与延迟处理](#十一主从复制原理与延迟处理)
- [十二、MySQL 高可用方案](#十二mysql-高可用方案)
- [十三、Buffer Pool 工作原理](#十三buffer-pool-工作原理)
- [十四、大表优化方案](#十四大表优化方案)
- [十五、utf8 与 utf8mb4 的区别](#十五utf8-与-utf8mb4-的区别)
- [十六、数据库三大范式与反范式](#十六数据库三大范式与反范式)
- [十七、数据库连接池原理与配置](#十七数据库连接池原理与配置)
- [十八、Online DDL 原理](#十八online-ddl-原理)
- [十九、Change Buffer 写缓冲](#十九change-buffer-写缓冲)

---

## 一、MySQL 事务及实现原理

### 什么是事务

```
事务 = 一组 SQL 操作，要么全部成功，要么全部失败，不存在中间状态。

经典场景 — 转账:
  A 账户 -100 元
  B 账户 +100 元

  如果 A 扣了钱，B 没加上 → 钱凭空消失
  事务保证: 两条 SQL 要么都成功，要么都回滚
```

### ACID 四大特性

| 特性 | 含义 | InnoDB 实现方式 |
|------|------|----------------|
| **A — 原子性 (Atomicity)** | 事务中的操作要么全做，要么全不做 | **Undo Log（回滚日志）** |
| **C — 一致性 (Consistency)** | 事务前后数据库从一个一致状态到另一个一致状态 | 由 A + I + D 共同保证 |
| **I — 隔离性 (Isolation)** | 并发事务之间互不干扰 | **MVCC + 锁机制** |
| **D — 持久性 (Durability)** | 事务提交后数据永久保存，即使宕机也不丢 | **Redo Log（重做日志）** |

```
ACID 的依赖关系:

  Undo Log ──→ 原子性 (A) ──┐
  MVCC + 锁 ──→ 隔离性 (I) ──┼──→ 一致性 (C)
  Redo Log ──→ 持久性 (D) ──┘

  一致性是最终目标，原子性、隔离性、持久性是实现手段。
```

---

### 一、原子性 — Undo Log 实现

#### Undo Log 是什么

```
Undo Log = 回滚日志，记录数据修改前的旧值

  执行 UPDATE users SET balance = 200 WHERE id = 1;  (原来 balance = 300)

  Undo Log 记录:
    "将 id=1 的 balance 从 200 改回 300"（逆操作）

  如果事务回滚:
    读 Undo Log → 执行逆操作 → balance 恢复为 300
```

#### 三种操作的 Undo Log

```
INSERT:
  Undo Log 记录: DELETE 这条记录的主键
  回滚时: 根据主键删除

UPDATE:
  Undo Log 记录: 修改前的旧值
  回滚时: 用旧值覆盖

DELETE:
  Undo Log 记录: 删除前的完整行数据（标记删除，不是真删）
  回滚时: 重新插入
```

#### Undo Log 链（版本链）

```
同一行数据被多个事务修改时，形成 Undo Log 版本链:

  当前数据:  id=1, name="C", trx_id=300
                │ roll_pointer
                ▼
  Undo Log:  id=1, name="B", trx_id=200
                │ roll_pointer
                ▼
  Undo Log:  id=1, name="A", trx_id=100

  每行数据有两个隐藏字段:
    trx_id:        最后修改这行的事务 ID
    roll_pointer:  指向上一个版本的 Undo Log

  这条链是 MVCC 的核心数据结构（后面详解）
```

---

### 二、持久性 — Redo Log 实现

#### 为什么需要 Redo Log

```
问题: InnoDB 的数据存在磁盘上（.ibd 文件），但每次修改都直接写磁盘太慢（随机 I/O）

  没有 Redo Log 的方案:
    UPDATE → 修改 Buffer Pool 中的页 → 立即刷到磁盘
    → 每次修改都是随机 I/O → 极慢
    → 如果刷盘前宕机 → 数据丢失

  有 Redo Log 的方案 (WAL: Write-Ahead Logging):
    UPDATE → 修改 Buffer Pool → 写 Redo Log（顺序 I/O，极快）→ 返回成功
    → 后台异步将 Buffer Pool 脏页刷到磁盘
    → 即使宕机，恢复时重放 Redo Log 即可

  核心: 用顺序写（Redo Log）替代随机写（数据文件），大幅提升性能
```

#### Redo Log 的写入流程

```
事务提交:

  1. 修改 Buffer Pool 中的数据页（内存）
  2. 生成 Redo Log 写入 Redo Log Buffer（内存）
  3. 事务提交时 → Redo Log Buffer 刷到磁盘（fsync）
  4. 返回客户端"提交成功"
  5. 后台线程异步将 Buffer Pool 脏页刷到数据文件

  ┌──────────────┐     ┌──────────────────┐     ┌────────────┐
  │ Buffer Pool  │     │ Redo Log Buffer  │     │   磁盘      │
  │  (数据页)    │     │   (内存)         │     │            │
  │  ┌────────┐  │     │  ┌────────────┐  │     │ redo log   │
  │  │ 脏页   │──┼──2──┼→│ redo 记录   │──┼──3──┼→(ib_logfile)│
  │  └────────┘  │     │  └────────────┘  │     │            │
  └──────┬───────┘     └──────────────────┘     │ 数据文件   │
         │                                      │ (.ibd)     │
         └────────── 5. 后台异步刷盘 ───────────→│            │
                                                └────────────┘
```

#### innodb_flush_log_at_trx_commit

```
控制 Redo Log 的刷盘策略:

  = 0: 每秒刷盘（不是每次提交时）
       → 最快，但宕机可能丢 1 秒数据
       → 适合：日志表、非关键数据

  = 1: 每次提交都 fsync（默认）
       → 最安全，不丢数据
       → 适合：金融、交易等核心业务

  = 2: 每次提交写到 OS 缓存，每秒 fsync
       → 折中，MySQL 崩溃不丢，OS 崩溃可能丢 1 秒
       → 适合：大多数业务

  性能: 0 > 2 > 1
  安全: 1 > 2 > 0
```

#### Redo Log 的物理结构

```
Redo Log 是固定大小的环形文件组:

  ┌─────────────┐  ┌─────────────┐
  │ ib_logfile0  │→│ ib_logfile1  │→ (循环)
  └─────────────┘  └─────────────┘
       ↑ write pos          ↑ checkpoint
       │                    │
       └── 新日志写入位置     └── 已刷到数据文件的位置

  write pos 追上 checkpoint:
    → Redo Log 写满了
    → 必须等 checkpoint 推进（刷脏页）才能继续写
    → 这就是 "Redo Log 满了导致 MySQL 卡顿" 的原因

  默认: 2 个文件 × 48MB = 96MB (可配置更大)
```

---

### 三、隔离性 — MVCC + 锁机制

#### 四种隔离级别

| 隔离级别 | 脏读 | 不可重复读 | 幻读 | InnoDB 实现 |
|---------|------|-----------|------|------------|
| READ UNCOMMITTED | 可能 | 可能 | 可能 | 直接读最新数据 |
| READ COMMITTED (RC) | 不会 | 可能 | 可能 | **MVCC: 每条 SQL 创建新 ReadView** |
| **REPEATABLE READ (RR)** | 不会 | 不会 | InnoDB 不会 | **MVCC: 事务开始时创建 ReadView** |
| SERIALIZABLE | 不会 | 不会 | 不会 | 加锁串行执行 |

**InnoDB 默认隔离级别 = RR（可重复读）**

#### 三种并发问题

```
脏读: 读到其他事务未提交的数据（事务回滚后数据就不存在了）
  T1: UPDATE balance = 200 (未提交)
  T2: SELECT balance → 200（读到了未提交的脏数据）
  T1: ROLLBACK → balance 回到 300
  T2 读到的 200 是"从未存在"的数据

不可重复读: 同一事务内两次读同一行，结果不同（被其他事务修改了）
  T1: SELECT balance → 300
  T2: UPDATE balance = 200; COMMIT;
  T1: SELECT balance → 200（同一事务内两次读不一样）

幻读: 同一事务内两次查询，行数不同（被其他事务插入/删除了）
  T1: SELECT COUNT(*) WHERE age > 20 → 5 行
  T2: INSERT INTO users(age) VALUES(25); COMMIT;
  T1: SELECT COUNT(*) WHERE age > 20 → 6 行（多了一行"幻影"）
```

#### MVCC 实现原理

```
MVCC (Multi-Version Concurrency Control) = 多版本并发控制

核心组件:
  1. 隐藏字段: 每行数据有 trx_id（事务ID）和 roll_pointer（版本链指针）
  2. Undo Log 版本链: 同一行的历史版本串成链表
  3. ReadView: 事务的"快照"，决定能看到哪些版本
```

**ReadView 的结构：**

```
ReadView 包含四个字段:

  creator_trx_id:  创建该 ReadView 的事务 ID
  m_ids:           创建时所有活跃（未提交）的事务 ID 列表
  min_trx_id:      m_ids 中的最小值
  max_trx_id:      下一个将要分配的事务 ID（当前最大 + 1）

可见性判断规则（对版本链中每个版本的 trx_id）:

  ┌─────────────────────────────────────────────────────────────┐
  │  if trx_id == creator_trx_id → 可见（自己修改的）             │
  │  if trx_id < min_trx_id → 可见（在 ReadView 创建前已提交）    │
  │  if trx_id >= max_trx_id → 不可见（在 ReadView 创建后才开始）  │
  │  if min_trx_id <= trx_id < max_trx_id:                      │
  │    if trx_id in m_ids → 不可见（创建时还未提交）               │
  │    if trx_id not in m_ids → 可见（创建时已经提交了）           │
  └─────────────────────────────────────────────────────────────┘

  沿版本链从新到旧逐个判断，找到第一个"可见"的版本 → 返回该版本的数据
```

**RC vs RR 的区别就在 ReadView 的创建时机：**

```
READ COMMITTED (RC):
  每条 SELECT 都创建新的 ReadView
  → 其他事务提交后，下一次 SELECT 就能看到 → 不可重复读

REPEATABLE READ (RR):
  整个事务只在第一条 SELECT 时创建 ReadView，后续复用
  → 事务期间 ReadView 不变 → 始终看到相同的数据快照 → 可重复读
```

**MVCC 示例：**

```
初始: id=1, name="A", trx_id=100（已提交）

  T1 (trx_id=200) 开始
  T2 (trx_id=300) 开始

  T2: UPDATE SET name="B" WHERE id=1  → trx_id=300, Undo 指向旧版本
  T2: COMMIT

  T1: SELECT name WHERE id=1

  T1 的 ReadView (RR 模式, 事务开始时创建):
    creator_trx_id = 200
    m_ids = [200, 300]   ← 创建时 T1 和 T2 都活跃
    min_trx_id = 200
    max_trx_id = 301

  版本链:
    name="B", trx_id=300 → name="A", trx_id=100

  判断 trx_id=300:
    300 在 m_ids 中 → 不可见（T2 在 ReadView 创建时未提交）
    虽然 T2 现在已提交，但 ReadView 是事务开始时的快照

  判断 trx_id=100:
    100 < min_trx_id(200) → 可见

  结果: T1 读到 name="A" ✓ （可重复读！即使 T2 已提交也看不到）
```

#### InnoDB 的锁机制

```
按粒度分:

  行锁 (Row Lock):
    Record Lock:  锁定单条索引记录
    Gap Lock:     锁定索引记录之间的"间隙"（防止幻读插入）
    Next-Key Lock: Record Lock + Gap Lock（左开右闭区间）

  表锁 (Table Lock):
    意向锁 (IS/IX): 表明事务想获取行级共享/排他锁
    AUTO-INC Lock: 自增列插入时的表级锁

按模式分:

  共享锁 (S Lock): SELECT ... LOCK IN SHARE MODE / FOR SHARE
    多个事务可同时持有（读读兼容）

  排他锁 (X Lock): SELECT ... FOR UPDATE / INSERT / UPDATE / DELETE
    只有一个事务能持有（读写、写写互斥）

兼容矩阵:
          │ S Lock │ X Lock │
  S Lock  │  兼容  │  互斥  │
  X Lock  │  互斥  │  互斥  │
```

**Next-Key Lock 防幻读：**

```
RR 级别下，SELECT ... FOR UPDATE 使用 Next-Key Lock:

  表数据 id: 5, 10, 15, 20

  SELECT * WHERE id > 10 AND id < 18 FOR UPDATE;

  加锁范围 (Next-Key Lock):
    (10, 15] + (15, 20]   ← 锁住间隙 + 记录

  其他事务:
    INSERT id=12 → 阻塞（在 (10,15] 间隙中）
    INSERT id=16 → 阻塞（在 (15,20] 间隙中）
    INSERT id=25 → 成功（不在锁范围内）

  → 防止了幻读：T1 再次查询时不会多出新行
```

---

### 四、Redo Log vs Undo Log vs Binlog

| 维度 | Redo Log | Undo Log | Binlog |
|------|----------|----------|--------|
| **所属** | InnoDB 引擎层 | InnoDB 引擎层 | MySQL Server 层 |
| **内容** | 物理日志（数据页的修改） | 逻辑日志（修改前的旧值） | 逻辑日志（SQL 语句/行变更） |
| **作用** | 崩溃恢复（保证持久性） | 事务回滚 + MVCC | 主从复制 + 数据恢复 |
| **写入时机** | 事务执行中持续写入 | 数据修改前写入 | 事务提交时写入 |
| **大小** | 固定大小环形文件 | 不固定，随事务生命周期 | 追加写，可持续增长 |
| **是否循环** | 是（空间复用） | 是（purge 线程清理） | 否（归档保留） |

#### 两阶段提交（保证 Redo Log 和 Binlog 一致）

```
问题: Redo Log 和 Binlog 是两个独立的日志，如何保证一致？

  如果 Redo 写了但 Binlog 没写:
    → 主库恢复后有这条数据，从库没有 → 主从不一致

解决: 两阶段提交 (2PC)

  1. Prepare 阶段:
     → 写 Redo Log，标记为 prepare 状态
     → 写 Binlog

  2. Commit 阶段:
     → 将 Redo Log 标记为 commit 状态

  崩溃恢复:
    Redo=prepare, Binlog=有  → 提交（数据完整）
    Redo=prepare, Binlog=无  → 回滚（数据不完整）
    Redo=commit             → 已提交，无需处理
```

---

### 五、事务的使用

```sql
-- 显式事务
BEGIN;  -- 或 START TRANSACTION
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;  -- 或 ROLLBACK

-- 隐式事务
-- autocommit=1 时（默认），每条 SQL 自动提交
-- autocommit=0 时，需要手动 COMMIT

-- 保存点（部分回滚）
BEGIN;
UPDATE ...; -- 操作 1
SAVEPOINT sp1;
UPDATE ...; -- 操作 2（出错了）
ROLLBACK TO sp1;  -- 只回滚操作 2，操作 1 保留
COMMIT;
```

### 事务问题排查

```
-- 查看当前运行的事务
SELECT * FROM information_schema.INNODB_TRX;

-- 查看锁等待
SELECT * FROM information_schema.INNODB_LOCK_WAITS;  -- MySQL 5.7
SELECT * FROM performance_schema.data_lock_waits;     -- MySQL 8.0

-- 查看死锁信息
SHOW ENGINE INNODB STATUS\G  -- 看 LATEST DETECTED DEADLOCK 部分

-- 长事务: 运行超过 60 秒的事务
SELECT * FROM information_schema.INNODB_TRX
WHERE TIME_TO_SEC(TIMEDIFF(NOW(), trx_started)) > 60;
```

### 一句话总结

**MySQL 事务 ACID 的实现：原子性靠 Undo Log（记录旧值用于回滚），持久性靠 Redo Log（WAL 顺序写+fsync 防宕机丢数据），隔离性靠 MVCC（ReadView + Undo 版本链实现快照读，无锁并发）+ 锁（当前读用 Next-Key Lock 防幻读），一致性是前三者共同保证的最终目标。** RC 和 RR 的唯一区别：RC 每条 SQL 创建新 ReadView，RR 整个事务复用一个 ReadView。Redo Log 和 Binlog 通过两阶段提交保证一致性。

---

## 二、Redo Log、Undo Log、Binlog 的作用及区别

### 三种日志在 ACID 中的角色

```
  ┌──────────────────────────────────────────────────────┐
  │                      ACID                            │
  │                                                      │
  │  Undo Log ──→ 原子性 (Atomicity)  ──┐               │
  │                + 隔离性（MVCC 版本链）  │               │
  │                                     ├──→ 一致性 (C)  │
  │  Redo Log ──→ 持久性 (Durability) ──┘               │
  │                                                      │
  │  Binlog ──→ 不属于 ACID，属于 MySQL Server 层         │
  │             用于：主从复制 + 数据恢复（Point-in-Time） │
  └──────────────────────────────────────────────────────┘
```

### 核心对比表

| 维度 | Redo Log | Undo Log | Binlog |
|------|----------|----------|--------|
| **层级** | InnoDB 引擎层 | InnoDB 引擎层 | **MySQL Server 层** |
| **ACID 角色** | **持久性 (D)** | **原子性 (A) + 隔离性 (I)** | 不属于 ACID |
| **记录内容** | 物理日志：某页某偏移处的字节改成了什么 | 逻辑日志：修改前的旧值（逆操作） | 逻辑日志：SQL 语句或行变更 |
| **核心作用** | 崩溃恢复（redo，重做） | 事务回滚 + MVCC 快照读 | 主从复制 + 数据恢复 |
| **写入时机** | 事务执行中**持续**写入 | 数据修改**前**写入 | 事务**提交时**一次性写入 |
| **文件形式** | 固定大小，环形覆写（ib_logfile0/1） | 共享表空间或独立 undo 表空间 | 追加写，不覆盖（mysql-bin.000001...） |
| **是否循环** | 是（空间复用，checkpoint 推进后覆写） | 是（purge 线程清理不再需要的版本） | **否**（持续增长，按策略归档/删除） |
| **所有引擎** | 仅 InnoDB | 仅 InnoDB | **所有引擎**（MyISAM 也有） |
| **是否可关闭** | 不可关闭 | 不可关闭 | 可关闭（`log_bin=OFF`） |

---

### 一、Redo Log — 保证持久性

```
问题: 事务提交后，数据还在 Buffer Pool（内存）中，没刷到磁盘
     此时如果 MySQL 宕机 → 内存丢失 → 数据丢了？

Redo Log 的解决:
  事务提交时，先把"做了什么修改"写到 Redo Log 文件（磁盘）
  → 即使 Buffer Pool 丢了
  → 重启后重放 Redo Log → 数据恢复

本质: "先记日记，后做事" (WAL: Write-Ahead Logging)
```

```
记录格式（物理日志）:
  "表空间号=5, 页号=100, 偏移量=200, 写入值=0x1234"

  不是记录 SQL 语句，而是记录数据页的物理修改
  → 恢复时直接在数据页上应用修改，速度极快（幂等操作）

  特点: 同一个修改多次重放结果相同（幂等性）
        比逻辑日志恢复更快（不需要重新执行 SQL）
```

```
写入流程:

  UPDATE balance = 200 WHERE id = 1

  ① 从磁盘加载 id=1 所在页到 Buffer Pool（如果不在内存）
  ② 在 Buffer Pool 中修改 balance = 200
  ③ 生成 Redo Log 记录，写入 Redo Log Buffer（内存）
  ④ 事务提交 → Redo Log Buffer 刷到 Redo Log 文件（fsync）
  ⑤ 返回客户端 "OK"
  ⑥ 后台线程择机将 Buffer Pool 脏页刷到 .ibd 数据文件

  ③④ 是顺序写（追加写一个文件，极快，~10μs）
  ⑥   是随机写（修改数据文件不同位置，较慢）
  → 用快的顺序写保证提交时安全，慢的随机写异步完成
```

---

### 二、Undo Log — 保证原子性 + 支撑隔离性

#### 作用 1: 事务回滚（原子性）

```
  BEGIN;
  UPDATE users SET name='B' WHERE id=1;  -- 原来 name='A'
  -- Undo Log 记录: "id=1, name 从 'B' 改回 'A'"

  UPDATE users SET age=30 WHERE id=1;    -- 原来 age=25
  -- Undo Log 记录: "id=1, age 从 30 改回 25"

  ROLLBACK;
  -- 读 Undo Log，逆序执行:
  --   age 改回 25
  --   name 改回 'A'
  -- 数据完全恢复，就像事务从未执行过
```

#### 作用 2: MVCC 快照读（隔离性）

```
  Undo Log 形成版本链，MVCC 通过 ReadView 在版本链上找到"该事务能看到的版本":

  当前数据: id=1, name="C", trx_id=300
               │ roll_pointer
               ▼
  Undo Log: id=1, name="B", trx_id=200
               │ roll_pointer
               ▼
  Undo Log: id=1, name="A", trx_id=100

  事务 T1 (RR 级别, ReadView 中 300 不可见):
    → 版本 name="C" (trx_id=300) 不可见
    → 沿链找到 name="B" (trx_id=200) 不可见
    → 沿链找到 name="A" (trx_id=100) 可见 → 返回 "A"

  没有加任何锁！→ 读写不互相阻塞 → 高并发
```

#### 记录格式（逻辑日志）

```
INSERT 的 Undo Log:
  记录主键 → 回滚时按主键 DELETE

UPDATE 的 Undo Log:
  记录旧值 → 回滚时用旧值覆盖
  包含: trx_id, roll_pointer, 修改前的各列值

DELETE 的 Undo Log:
  标记删除（delete_flag=1）+ 记录完整行
  → 回滚时取消标记
  → MVCC 中仍然能看到这行（按 ReadView 判断可见性）
```

#### Undo Log 的清理（Purge）

```
Undo Log 不会无限增长:

  purge 线程定期清理不再需要的 Undo Log:
    → 没有任何活跃事务需要读取的旧版本 → 可以清理
    → INSERT 的 Undo Log: 事务提交后立即可清理（不需要给 MVCC 用）
    → UPDATE/DELETE 的 Undo Log: 等所有可能用到它的 ReadView 都关闭后才清理

  长事务的危害:
    长事务持有旧的 ReadView → purge 无法清理旧版本
    → Undo Log 不断增长 → 磁盘空间膨胀 + 查询变慢（版本链太长）
```

---

### 三、Binlog — 主从复制 + 数据恢复

```
Binlog 不属于 InnoDB，是 MySQL Server 层的日志
→ 所有存储引擎都有 Binlog
→ 即使用 MyISAM 也能做主从复制
```

#### 三种格式

```
1. STATEMENT 格式:
   记录 SQL 语句原文
   Binlog: "UPDATE users SET balance=balance-100 WHERE id=1"
   优点: 日志量小
   缺点: 某些函数（NOW(), UUID()）在从库执行结果不同 → 主从不一致

2. ROW 格式（推荐）:
   记录行的变更前后值
   Binlog: "id=1, balance: 300→200"
   优点: 精确，不会因函数导致不一致
   缺点: 日志量大（UPDATE 10 万行 = 记录 10 万条变更）

3. MIXED 格式:
   默认用 STATEMENT，遇到不确定的函数自动切 ROW
   折中方案
```

#### 主从复制流程

```
  Master                           Slave
    │                                │
    │  事务提交 → 写 Binlog           │
    │                                │
    │←── IO Thread 请求 Binlog ──────│
    │                                │
    │─── 发送 Binlog Events ────────→│  Relay Log (中继日志)
    │                                │
    │                                │  SQL Thread 读取 Relay Log
    │                                │  → 重放 SQL/行变更
    │                                │  → 从库数据与主库一致
```

#### 数据恢复（Point-in-Time Recovery）

```
场景: 误删了表，需要恢复到删除前

  1. 找到最近的全量备份（如昨晚 2:00）
  2. 从全量备份恢复
  3. 用 mysqlbinlog 工具重放 2:00 到误删前的 Binlog
  4. 数据恢复到误删前的状态

  这就是为什么生产必须开 Binlog + 定期全量备份
```

---

### 四、三种日志的写入时序

```
一条 UPDATE 语句的完整日志流程:

  BEGIN;
  UPDATE users SET balance=200 WHERE id=1;  (原 balance=300)
  COMMIT;

  ┌─── 执行阶段 ───────────────────────────────────────────────┐
  │                                                            │
  │  ① 写 Undo Log: "id=1, balance 从 200 改回 300"            │
  │     → 写入 Undo Log（回滚用 + MVCC 版本链）                 │
  │                                                            │
  │  ② 修改 Buffer Pool 中的数据页: balance = 200               │
  │                                                            │
  │  ③ 写 Redo Log (prepare): 记录数据页的物理修改              │
  │     → 写入 Redo Log Buffer                                 │
  │                                                            │
  └────────────────────────────────────────────────────────────┘

  ┌─── 提交阶段（两阶段提交）─────────────────────────────────┐
  │                                                            │
  │  ④ Redo Log 刷盘，标记 prepare                             │
  │     → Redo Log Buffer → fsync → ib_logfile                │
  │                                                            │
  │  ⑤ 写 Binlog 并刷盘                                       │
  │     → Binlog Cache → fsync → mysql-bin.xxxxx              │
  │                                                            │
  │  ⑥ Redo Log 标记 commit                                   │
  │     → 事务正式提交完成                                      │
  │                                                            │
  └────────────────────────────────────────────────────────────┘

  ⑦ 后台异步: Buffer Pool 脏页刷到 .ibd 数据文件
```

```
时序图:

  时间 →

  Undo Log:   ██ (修改前写入)
  Buffer Pool:  ██ (修改数据页)
  Redo Log:       ██ prepare ──────────── commit
  Binlog:                      ██ (提交时写入)
  数据文件:                                        ██ (异步刷盘)

  两阶段提交保证 Redo Log 和 Binlog 的一致性:
    ④⑤⑥ 是原子的：
      Redo=prepare + Binlog=有 → 恢复时提交
      Redo=prepare + Binlog=无 → 恢复时回滚
```

---

### 五、崩溃恢复场景分析

```
场景1: ④ 之前崩溃（Redo Log 未 prepare）
  → Redo Log 没有该事务记录
  → Undo Log 存在 → 回滚
  → 结果: 事务未提交，数据回滚 ✓

场景2: ④ 之后、⑤ 之前崩溃（Redo prepare, Binlog 未写）
  → Redo Log 有 prepare，但 Binlog 无
  → 回滚该事务
  → 结果: 主从一致（都没有这个事务）✓

场景3: ⑤ 之后、⑥ 之前崩溃（Redo prepare, Binlog 已写）
  → Redo Log 有 prepare，Binlog 有完整记录
  → 提交该事务（Redo 标记为 commit）
  → 结果: 主从一致（都有这个事务）✓

场景4: ⑥ 之后崩溃
  → Redo Log 已 commit
  → 正常恢复，重放 Redo Log 即可
  → 结果: 数据不丢 ✓
```

---

### 六、Redo Log 和 Binlog 为什么不能合并

```
为什么不能只用一个日志？

  只用 Redo Log，不用 Binlog？
    ✗ Redo Log 是 InnoDB 专属，MyISAM 等引擎没有
    ✗ Redo Log 是环形覆写，历史记录会被覆盖，不能做全量恢复
    ✗ Redo Log 是物理日志，从库无法重放（页结构可能不同）

  只用 Binlog，不用 Redo Log？
    ✗ Binlog 是提交时才写，没有 WAL 能力
    ✗ 如果 Buffer Pool 脏页未刷盘就崩溃，Binlog 无法恢复
    ✗ Binlog 不记录未提交事务的中间状态，无法做崩溃恢复

  结论: 两者职责不同，无法合并
    Redo Log: 引擎层崩溃恢复（WAL）
    Binlog:   Server 层复制和归档
```

### 一句话总结

| 日志 | ACID 角色 | 一句话 |
|------|----------|--------|
| **Redo Log** | **持久性 (D)** | "做了什么"的物理记录，WAL 保证提交后不丢，崩溃后重放恢复 |
| **Undo Log** | **原子性 (A) + 隔离性 (I)** | "改之前是什么"的逆操作，回滚靠它，MVCC 版本链也靠它 |
| **Binlog** | **不属于 ACID** | Server 层逻辑日志，用于主从复制和 Point-in-Time 数据恢复 |

**三者协作：** Undo 在修改前写入（保护旧值）→ Redo 在执行中写入（记录新值）→ 提交时两阶段提交（Redo prepare → Binlog 刷盘 → Redo commit）保证 Redo 和 Binlog 一致 → 后台异步刷数据页。

---

## 三、MVCC 的实现原理

### MVCC 是什么

```
MVCC (Multi-Version Concurrency Control) = 多版本并发控制

核心思想:
  对同一行数据保留多个历史版本
  读操作不加锁，而是根据"快照"读取对应版本的数据
  → 读写不互相阻塞 → 高并发

没有 MVCC:
  T1: SELECT balance WHERE id=1  → 加读锁，阻塞 T2
  T2: UPDATE balance WHERE id=1  → 等 T1 的读锁释放
  → 读写串行，并发低

有 MVCC:
  T1: SELECT balance WHERE id=1  → 读历史版本（无锁）
  T2: UPDATE balance WHERE id=1  → 修改最新版本（不影响 T1）
  → 读写并行，互不阻塞
```

**MVCC 仅在 READ COMMITTED（RC）和 REPEATABLE READ（RR）两个隔离级别下工作。**

### MVCC 的三个核心组件

```
  ┌─────────────────────────────────────────────────┐
  │  1. 隐藏字段 — 每行数据的元信息                    │
  │  2. Undo Log 版本链 — 同一行的所有历史版本          │
  │  3. ReadView — 事务的可见性快照                    │
  └─────────────────────────────────────────────────┘
```

---

### 一、隐藏字段

```
InnoDB 为每行数据自动添加三个隐藏字段:

  ┌──────────┬───────────┬─────────────┬──────────────────┐
  │   id     │  name     │  trx_id     │  roll_pointer    │
  │  (用户列) │  (用户列)  │  (隐藏)      │  (隐藏)           │
  └──────────┴───────────┴─────────────┴──────────────────┘

  trx_id (6 字节):
    最近一次修改该行的事务 ID
    每个事务有一个全局递增的唯一 ID

  roll_pointer (7 字节):
    回滚指针，指向 Undo Log 中该行的上一个版本
    形成版本链

  row_id (6 字节，仅在无主键时自动生成):
    隐含的自增行 ID
```

---

### 二、Undo Log 版本链

```
同一行数据每次被修改，旧版本通过 roll_pointer 串成链表:

  假设 id=1 这行被 3 个事务依次修改:

  trx_id=100: INSERT (id=1, name="Alice", age=20)
  trx_id=200: UPDATE name="Bob"
  trx_id=300: UPDATE age=25

  当前数据页中 (最新版本):
  ┌────┬───────┬─────┬──────────┬───────────────┐
  │ id │ name  │ age │ trx_id   │ roll_pointer  │
  │ 1  │ Bob   │ 25  │ 300      │  ──→ Undo #2  │
  └────┴───────┴─────┴──────────┴───────────────┘
                                        │
                                        ▼
  Undo Log #2 (trx_id=200 修改前的值):
  ┌────┬───────┬─────┬──────────┬───────────────┐
  │ 1  │ Bob   │ 20  │ 200      │  ──→ Undo #1  │
  └────┴───────┴─────┴──────────┴───────────────┘
                                        │
                                        ▼
  Undo Log #1 (trx_id=100 插入时的原始值):
  ┌────┬───────┬─────┬──────────┬───────────────┐
  │ 1  │ Alice │ 20  │ 100      │  ──→ NULL     │
  └────┴───────┴─────┴──────────┴───────────────┘

  SELECT 时沿着这条链，找到当前事务"能看到"的版本
```

---

### 三、ReadView（核心中的核心）

#### ReadView 的结构

```go
// 伪代码，实际是 C++ 实现
type ReadView struct {
    creator_trx_id  uint64    // 创建该 ReadView 的事务 ID
    m_ids           []uint64  // 创建时所有活跃（未提交）的事务 ID 列表
    min_trx_id      uint64    // m_ids 中的最小值
    max_trx_id      uint64    // 创建时系统分配的下一个事务 ID（当前最大+1）
}
```

```
  活跃事务: 已开始但尚未提交的事务

  举例:
    系统中事务情况: T100(已提交), T200(活跃), T300(活跃), T400(活跃)
    下一个将分配的事务 ID = 500

    此时创建 ReadView:
      creator_trx_id = 200 (假设是 T200 创建的)
      m_ids = [200, 300, 400]
      min_trx_id = 200
      max_trx_id = 500
```

#### 可见性判断规则

```
对版本链中某个版本的 trx_id 进行判断:

                  可见区域              不确定区域              不可见区域
  ◄──────────────────────┤────────────────────────────┤──────────────►
  0                  min_trx_id                   max_trx_id
                        (200)                        (500)

  规则（按优先级顺序）:

  ① trx_id == creator_trx_id
     → 可见 ✓ (自己修改的，当然能看到)

  ② trx_id < min_trx_id
     → 可见 ✓ (在 ReadView 创建前已提交)

  ③ trx_id >= max_trx_id
     → 不可见 ✗ (在 ReadView 创建后才开始的事务)

  ④ min_trx_id <= trx_id < max_trx_id
     → 看 m_ids:
        trx_id 在 m_ids 中 → 不可见 ✗ (创建时该事务还没提交)
        trx_id 不在 m_ids 中 → 可见 ✓ (创建时该事务已经提交了)
```

#### 完整查询流程

```
SELECT name FROM users WHERE id = 1;

  1. 找到 id=1 的最新行: name="Bob", trx_id=300
  2. 用 ReadView 判断 trx_id=300 是否可见
     → 不可见 → 沿 roll_pointer 找下一个版本
  3. Undo Log #2: name="Bob", trx_id=200
     → 不可见 → 继续沿链找
  4. Undo Log #1: name="Alice", trx_id=100
     → 可见 ✓ → 返回 name="Alice"

  整个过程没有加任何锁！
```

---

### 四、RC 和 RR 的区别：ReadView 创建时机

**RC 和 RR 使用完全相同的 MVCC 机制，唯一区别是 ReadView 的创建时机。**

#### 完整对比示例

```
初始数据: id=1, name="Alice", trx_id=80 (已提交)

事务分配顺序: T100, T200, T300

时间线:
  t1: T100 BEGIN
  t2: T200 BEGIN
  t3: T300 BEGIN
  t4: T300: UPDATE name="Bob" WHERE id=1    → trx_id=300
  t5: T300: COMMIT
  t6: T100: SELECT name WHERE id=1          → 第一次查询
  t7: T200: UPDATE name="Carol" WHERE id=1  → trx_id=200
  t8: T200: COMMIT
  t9: T100: SELECT name WHERE id=1          → 第二次查询
```

#### RR 模式（可重复读）

```
t6 时 T100 第一次 SELECT → 创建 ReadView（整个事务只创建一次）:
  creator_trx_id = 100
  m_ids = [100, 200]     ← T300 已提交不在列表中
  min_trx_id = 100
  max_trx_id = 301

  版本链: name="Bob"(trx_id=300) → name="Alice"(trx_id=80)

  判断 trx_id=300: 不在 m_ids 中，且 300 < 301 → 可见 ✓
  结果: name = "Bob" ✓

t9 时 T100 第二次 SELECT → 复用 t6 的 ReadView（不重新创建）:
  同一个 ReadView:
    m_ids = [100, 200]
    min_trx_id = 100
    max_trx_id = 301

  版本链: name="Carol"(trx_id=200) → name="Bob"(trx_id=300) → name="Alice"(trx_id=80)

  判断 trx_id=200: 200 在 m_ids 中 → 不可见 ✗
  判断 trx_id=300: 不在 m_ids 中，300 < 301 → 可见 ✓
  结果: name = "Bob" ✓ (和第一次一样！可重复读！)
```

#### RC 模式（读已提交）

```
t6 时 T100 第一次 SELECT → 创建 ReadView #1:
  creator_trx_id = 100
  m_ids = [100, 200]     ← T300 已提交
  min_trx_id = 100
  max_trx_id = 301

  判断 trx_id=300: 不在 m_ids 中 → 可见 ✓
  结果: name = "Bob" ✓

t9 时 T100 第二次 SELECT → 创建新的 ReadView #2（每次 SQL 都重新创建！）:
  creator_trx_id = 100
  m_ids = [100]           ← T200 和 T300 都已提交，不在列表中了
  min_trx_id = 100
  max_trx_id = 301

  判断 trx_id=200: 不在 m_ids 中 → 可见 ✓
  结果: name = "Carol" ✓ (看到了 T200 提交的新值！不可重复读！)
```

```
总结:
  ┌──────────────────────────────────────────────────────┐
  │  RR: 事务第一条 SELECT 创建 ReadView → 后续复用        │
  │      → 同一事务内多次读结果相同 → 可重复读              │
  │                                                      │
  │  RC: 每条 SELECT 创建新 ReadView                      │
  │      → 每次读都能看到最新已提交数据 → 不可重复读        │
  └──────────────────────────────────────────────────────┘
```

---

### 五、快照读 vs 当前读

```
MVCC 只作用于"快照读"，"当前读"走锁机制:

快照读（Snapshot Read）— 走 MVCC，无锁:
  SELECT ... （普通 SELECT，不加任何锁修饰）
  → 读 ReadView 对应的历史版本
  → 不加锁，不阻塞

当前读（Current Read）— 走锁，读最新版本:
  SELECT ... FOR UPDATE           (加排他锁 X)
  SELECT ... FOR SHARE            (加共享锁 S，MySQL 8.0)
  SELECT ... LOCK IN SHARE MODE   (加共享锁 S，MySQL 5.7)
  INSERT / UPDATE / DELETE        (隐含加排他锁 X)
  → 读取数据的最新版本（不走 MVCC）
  → 加锁，阻塞其他写操作

为什么 UPDATE/DELETE 必须当前读？
  → 如果 UPDATE 读的是历史版本，可能覆盖其他事务的修改
  → 必须读最新值，加锁修改，保证数据正确

示例:
  T1 (RR):
    SELECT balance WHERE id=1;              → 快照读 → 300（历史版本）
    SELECT balance WHERE id=1 FOR UPDATE;   → 当前读 → 200（最新版本）

    同一事务内两种读的结果可能不同！
```

---

### 六、MVCC 和幻读的关系

```
RR 级别下 MVCC 能防幻读吗？

快照读: 能防。
  因为 ReadView 固定，新插入的行 trx_id 一定 >= max_trx_id → 不可见
  T1: SELECT * WHERE age > 20 → 5 行
  T2: INSERT (age=25); COMMIT;
  T1: SELECT * WHERE age > 20 → 还是 5 行 ✓（快照读看不到 T2 的插入）

当前读: 不能完全防，需要配合 Next-Key Lock。
  T1: SELECT * WHERE age > 20 FOR UPDATE → 5 行 + 加 Next-Key Lock
  T2: INSERT (age=25) → 阻塞！被 Gap Lock 挡住
  T1: SELECT * WHERE age > 20 FOR UPDATE → 还是 5 行 ✓

  如果 T1 第一次是快照读（没加锁），T2 的 INSERT 不会阻塞:
  T1: SELECT * WHERE age > 20 → 5 行（快照读，无锁）
  T2: INSERT (age=25); COMMIT;
  T1: UPDATE users SET name='X' WHERE age > 20 → 更新了 6 行！
      → UPDATE 是当前读，能看到 T2 插入的行 → 幻读发生！

  所以: MVCC 的快照读能防幻读，但混用快照读和当前读时仍可能出现
```

---

### MVCC 总结图

```
                        ┌──────────────┐
  普通 SELECT  ────────→│  快照读       │──→ MVCC（无锁）
                        │              │
                        │  ReadView    │──→ RR: 事务内复用
                        │  +           │    RC: 每次 SQL 新建
                        │  版本链       │
                        └──────────────┘

  SELECT FOR UPDATE ──→ ┌──────────────┐
  INSERT/UPDATE/DELETE → │  当前读       │──→ 锁机制（X/S Lock）
                        │              │    + Next-Key Lock 防幻读
                        │  最新版本     │
                        └──────────────┘
```

### 一句话总结

**MVCC = 隐藏字段(trx_id + roll_pointer) + Undo Log 版本链 + ReadView 可见性判断。** 读操作不加锁，沿版本链找到 ReadView 认为"可见"的第一个版本返回。RC 每条 SQL 新建 ReadView（能看到最新已提交），RR 整个事务复用一个 ReadView（永远看到事务开始时的快照）——这就是"可重复读"的全部秘密。快照读走 MVCC 无锁高并发，当前读走锁保证数据正确性，二者配合实现了 InnoDB 的高性能事务。

---

## 四、MySQL 中的锁有哪几种及其作用

### 锁的全景分类

```
MySQL 锁
├── 按粒度分
│   ├── 全局锁 (Global Lock)
│   ├── 表级锁 (Table Lock)
│   │   ├── 表锁 (Table Lock)
│   │   ├── 元数据锁 (MDL, Metadata Lock)
│   │   ├── 意向锁 (Intention Lock)
│   │   └── AUTO-INC 锁
│   └── 行级锁 (Row Lock) — InnoDB 独有
│       ├── 记录锁 (Record Lock)
│       ├── 间隙锁 (Gap Lock)
│       └── 临键锁 (Next-Key Lock)
│
├── 按模式分
│   ├── 共享锁 (S Lock / 读锁)
│   └── 排他锁 (X Lock / 写锁)
│
└── 按思想分
    ├── 悲观锁 (Pessimistic Lock)
    └── 乐观锁 (Optimistic Lock)
```

---

### 一、全局锁（Global Lock）

```sql
FLUSH TABLES WITH READ LOCK;  -- 加全局读锁（FTWRL）
-- 整个数据库变为只读，所有写操作阻塞
-- 做全库逻辑备份时使用 (mysqldump)
UNLOCK TABLES;                -- 释放
```

```
作用: 保证备份时数据一致性（所有表同一时间点的快照）

问题: 全库只读，业务完全阻塞

更好的替代:
  mysqldump --single-transaction  (InnoDB 专用)
  → 利用 MVCC 在事务内获得一致性快照
  → 不需要全局锁，不影响业务写入
```

---

### 二、表级锁

#### 1. 表锁（Table Lock）

```sql
LOCK TABLES users READ;    -- 表级共享锁（只读）
LOCK TABLES users WRITE;   -- 表级排他锁（读写都锁）
UNLOCK TABLES;
```

```
特点:
  - MyISAM 默认使用表锁（MyISAM 没有行锁）
  - InnoDB 一般不需要手动加表锁（有行锁更细粒度）
  - 粒度大，并发低，但开销小、不会死锁
```

#### 2. 元数据锁（MDL, Metadata Lock）

```
MDL 不需要手动加，MySQL 自动管理:

  对表做 DML (SELECT/INSERT/UPDATE/DELETE):
    → 自动加 MDL 读锁（多个 DML 可并发）

  对表做 DDL (ALTER TABLE/DROP TABLE):
    → 自动加 MDL 写锁（和所有 DML/DDL 互斥）

作用: 防止 DML 和 DDL 并发冲突
  → 不能一边查表一边改表结构

经典问题 — DDL 导致全库阻塞:
  T1: SELECT * FROM users (长查询，持有 MDL 读锁)
  T2: ALTER TABLE users ADD COLUMN (需要 MDL 写锁 → 阻塞等 T1)
  T3: SELECT * FROM users (需要 MDL 读锁 → 被 T2 阻塞！)
  T4: SELECT * FROM users (也被阻塞！)
  → T2 阻塞了后续所有对该表的操作，引发雪崩

  解决: DDL 加超时 (ALTER TABLE ... WAIT N / NOWAIT)
       或在低峰期执行
```

#### 3. 意向锁（Intention Lock）

```
意向锁是 InnoDB 自动加的表级锁，用于快速判断"表内是否有行锁":

  IS (Intention Shared):  事务准备给某行加 S 锁前，先在表上加 IS
  IX (Intention Exclusive): 事务准备给某行加 X 锁前，先在表上加 IX

为什么需要？
  假设没有意向锁:
    T1 对 id=1 加了行级 X 锁
    T2 想加表级 S 锁 → 需要检查表中每一行是否有 X 锁 → 遍历全表！

  有了意向锁:
    T1 加行级 X 锁前，先在表上加 IX
    T2 想加表级 S 锁 → 发现表上有 IX → 直接互斥，无需遍历行

兼容矩阵:
           │ IS   │ IX   │  S   │  X   │
  IS       │ 兼容 │ 兼容 │ 兼容 │ 互斥 │
  IX       │ 兼容 │ 兼容 │ 互斥 │ 互斥 │
  S        │ 兼容 │ 互斥 │ 兼容 │ 互斥 │
  X        │ 互斥 │ 互斥 │ 互斥 │ 互斥 │

  注意: IS 和 IX 之间互相兼容（因为行锁之间的冲突在行级判断）
       意向锁只和表级 S/X 锁互斥
```

#### 4. AUTO-INC 锁

```
对自增列插入时的特殊锁:

  innodb_autoinc_lock_mode:

  = 0 (传统模式):
    每次 INSERT 加表级 AUTO-INC 锁，语句结束释放
    → 保证连续自增，但并发差

  = 1 (连续模式，MySQL 5.x 默认):
    简单 INSERT: 轻量级互斥量（快速分配，不加表锁）
    批量 INSERT (INSERT...SELECT): 加表级 AUTO-INC 锁

  = 2 (交错模式，MySQL 8.0 默认):
    所有 INSERT 都用轻量级互斥量
    → 最高并发，但批量插入的 ID 可能不连续
    → 基于 STATEMENT 的 Binlog 复制不安全（推荐 ROW 格式）
```

---

### 三、行级锁（InnoDB 独有，核心重点）

**行锁是加在索引上的，不是加在行数据上。如果查询没走索引，行锁会退化为表锁。**

#### 1. 记录锁（Record Lock）

```
锁定索引中的一条记录:

  SELECT * FROM users WHERE id = 5 FOR UPDATE;

  → 在 id=5 的索引记录上加 X 锁
  → 其他事务不能修改/删除 id=5 这行
  → 但不影响 id=4 或 id=6

  索引:  ... [3] [5] [8] [10] ...
               ↑
          Record Lock (仅锁这一条)
```

#### 2. 间隙锁（Gap Lock）

```
锁定索引记录之间的"间隙"，防止其他事务在间隙中插入新记录:

  表数据 id: 3, 5, 8, 10

  SELECT * FROM users WHERE id = 7 FOR UPDATE;
  → id=7 不存在，但锁住了 (5, 8) 这个间隙

  索引:  ... [3] [5]  (  间隙  )  [8] [10] ...
                      ↑ Gap Lock ↑
                  锁住了 5 和 8 之间的间隙

  其他事务:
    INSERT id=6 → 阻塞（在间隙内）
    INSERT id=7 → 阻塞（在间隙内）
    UPDATE id=5 → 不阻塞（Gap Lock 不锁记录本身）
    UPDATE id=8 → 不阻塞

作用: 防止幻读（阻止在范围内插入新行）

特点:
  - Gap Lock 之间互不冲突（两个事务可以同时锁同一间隙）
  - Gap Lock 只阻止 INSERT，不阻止 SELECT/UPDATE/DELETE
  - 只在 RR 隔离级别下存在，RC 级别没有 Gap Lock
```

#### 3. 临键锁（Next-Key Lock）

```
Next-Key Lock = Record Lock + Gap Lock（左开右闭区间）

  表数据 id: 3, 5, 8, 10

  InnoDB 默认的行锁就是 Next-Key Lock:
    对 id=8 加 Next-Key Lock → 锁 (5, 8]

    索引:  ... [3] [5]  (  间隙  )  [8] [10] ...
                        ↑ Gap Lock ↑↑ Record Lock
                           (5, 8]

  SELECT * FROM users WHERE id BETWEEN 5 AND 8 FOR UPDATE;
  加锁范围:
    (3, 5] → Next-Key Lock
    (5, 8] → Next-Key Lock
    (8, 10] → Next-Key Lock（防止在 8 后面插入）

  其他事务:
    INSERT id=4 → 阻塞
    INSERT id=6 → 阻塞
    INSERT id=9 → 阻塞
    INSERT id=11 → 不阻塞
```

#### 加锁规则总结（RR 级别下）

```
InnoDB 的加锁规则（两个原则 + 两个优化 + 一个 bug）:

原则1: 加锁的基本单位是 Next-Key Lock (左开右闭)
原则2: 只有访问到的对象才会加锁

优化1: 唯一索引上的等值查询，命中记录 → 退化为 Record Lock
优化2: 等值查询，向右遍历到不满足条件的记录 → Gap Lock 退化（右边界不锁记录）

示例 — 等值查询命中唯一索引:
  SELECT * FROM users WHERE id = 5 FOR UPDATE;  (id 是主键)
  → 只加 Record Lock (锁 id=5 这条记录)
  → 不是 Next-Key Lock (优化1: 唯一索引等值命中退化)

示例 — 等值查询未命中:
  SELECT * FROM users WHERE id = 7 FOR UPDATE;  (id=7 不存在)
  表数据: 5, 8
  → 加 Gap Lock (5, 8)
  → 不锁 5 也不锁 8（优化2: 等值查不存在，退化为纯间隙锁）

示例 — 范围查询:
  SELECT * FROM users WHERE id >= 5 AND id < 8 FOR UPDATE;
  表数据: 3, 5, 8, 10
  → id=5: Record Lock（唯一索引等值命中，优化1）
  → (5, 8]: Next-Key Lock
  → 总加锁: id=5 的 Record Lock + (5, 8] 的 Next-Key Lock
```

---

### 四、共享锁 vs 排他锁

```
共享锁 (S Lock / 读锁):
  SELECT ... FOR SHARE;          (MySQL 8.0)
  SELECT ... LOCK IN SHARE MODE; (MySQL 5.7)
  → 允许其他事务也加 S 锁（读读兼容）
  → 阻止其他事务加 X 锁

排他锁 (X Lock / 写锁):
  SELECT ... FOR UPDATE;
  INSERT / UPDATE / DELETE;  (自动加 X 锁)
  → 阻止其他事务加任何锁（S 和 X 都互斥）

兼容矩阵:
           │ S Lock │ X Lock │
  S Lock   │  兼容  │  互斥  │  (读读兼容)
  X Lock   │  互斥  │  互斥  │  (读写、写写互斥)
```

---

### 五、乐观锁 vs 悲观锁

```
这两个不是 MySQL 内置的锁类型，而是并发控制的思想:

悲观锁: "我觉得别人一定会来改数据" → 先加锁再操作
  实现: SELECT ... FOR UPDATE (数据库行锁)
  适用: 写多读少，冲突频繁

  BEGIN;
  SELECT balance FROM accounts WHERE id=1 FOR UPDATE;  -- 加锁
  UPDATE accounts SET balance = balance - 100 WHERE id=1;
  COMMIT;

乐观锁: "我觉得别人大概率不会改" → 不加锁，提交时检查是否冲突
  实现: 版本号(version) 或 CAS
  适用: 读多写少，冲突罕见

  -- 读取时记住版本号
  SELECT balance, version FROM accounts WHERE id=1;
  -- balance=300, version=5

  -- 更新时检查版本号是否变化
  UPDATE accounts SET balance=200, version=6
  WHERE id=1 AND version=5;

  -- affected_rows=1 → 成功
  -- affected_rows=0 → 被别人改过了 → 重试
```

---

### 六、死锁

```
两个事务互相等待对方持有的锁 → 永远等下去:

  T1: BEGIN;
      UPDATE users SET name='A' WHERE id=1;  -- 持有 id=1 的 X 锁
      UPDATE users SET name='B' WHERE id=2;  -- 等待 id=2 的 X 锁 (T2 持有)

  T2: BEGIN;
      UPDATE users SET name='C' WHERE id=2;  -- 持有 id=2 的 X 锁
      UPDATE users SET name='D' WHERE id=1;  -- 等待 id=1 的 X 锁 (T1 持有)

  → T1 等 T2，T2 等 T1 → 死锁！

InnoDB 的死锁处理:
  1. 检测: wait-for graph（等待图），检测到环 → 死锁
  2. 回滚: 选择一个"代价最小"的事务回滚（通常是 undo log 量最少的）
  3. 报错: ERROR 1213 (40001): Deadlock found

预防死锁:
  - 按固定顺序访问表和行（如始终按 id 升序加锁）
  - 事务尽量短小（持锁时间短）
  - 合理使用索引（减少锁范围）
  - 避免大事务
```

---

### 七、查看锁信息

```sql
-- MySQL 8.0 查看当前持有的锁
SELECT * FROM performance_schema.data_locks;

-- 查看锁等待关系
SELECT * FROM performance_schema.data_lock_waits;

-- 查看最近的死锁信息
SHOW ENGINE INNODB STATUS\G
-- 找到 LATEST DETECTED DEADLOCK 部分

-- MySQL 5.7
SELECT * FROM information_schema.INNODB_LOCKS;
SELECT * FROM information_schema.INNODB_LOCK_WAITS;
```

---

### 全景总结表

| 锁类型 | 粒度 | 作用 | 加锁方式 |
|--------|------|------|---------|
| **全局锁** | 整个数据库 | 全库备份时保证一致性 | `FLUSH TABLES WITH READ LOCK` |
| **表锁** | 整张表 | MyISAM 默认锁，粒度大 | `LOCK TABLES ... READ/WRITE` |
| **MDL** | 整张表 | 防止 DML 和 DDL 并发冲突 | 自动加，对表操作即触发 |
| **意向锁 IS/IX** | 表级标记 | 快速判断表内是否有行锁 | 加行锁前自动加 |
| **AUTO-INC** | 表级 | 自增 ID 的并发分配 | INSERT 时自动加 |
| **Record Lock** | 单条索引记录 | 锁住一行，防止修改/删除 | 唯一索引等值命中 |
| **Gap Lock** | 索引间隙 | 锁住间隙，防止插入（防幻读） | 等值查不存在 / RR 级别 |
| **Next-Key Lock** | 记录+前间隙 | Record + Gap，InnoDB 默认行锁 | 范围查询 / 非唯一索引 |
| **S Lock** | 行/表 | 读锁，读读兼容 | `FOR SHARE` |
| **X Lock** | 行/表 | 写锁，与所有锁互斥 | `FOR UPDATE` / DML |

### 一句话总结

**InnoDB 的锁体系 = 表级锁（MDL 防 DDL 冲突 + 意向锁加速判断）+ 行级锁（Record/Gap/Next-Key 三种组合）。** 行锁加在索引上，不走索引则退化为表锁。RR 级别默认用 Next-Key Lock（Record + Gap）防幻读；唯一索引等值命中退化为 Record Lock（最小粒度）；等值查不存在退化为 Gap Lock。死锁通过 wait-for graph 检测并自动回滚代价最小的事务。预防死锁的关键：按固定顺序加锁 + 事务尽量短 + 合理索引。

---

## 五、OLTP 与 OLAP 数据库及存储引擎原理

### OLTP vs OLAP 基本概念

```
OLTP (Online Transaction Processing) — 在线事务处理
  → 面向日常业务操作：增删改查单条/少量数据
  → 关键词：低延迟、高并发、事务、行级操作

OLAP (Online Analytical Processing) — 在线分析处理
  → 面向数据分析：聚合、统计、报表
  → 关键词：大数据量扫描、列聚合、复杂查询、高吞吐
```

### 核心对比

| 维度 | OLTP | OLAP |
|------|------|------|
| **典型操作** | INSERT/UPDATE/DELETE + 点查 | SELECT + GROUP BY + JOIN（大范围扫描） |
| **数据量** | 单次操作几行~几十行 | 单次查询扫描百万~亿行 |
| **延迟要求** | 毫秒级 | 秒级~分钟级可接受 |
| **并发模式** | 高并发短事务（万级 QPS） | 低并发长查询（几十个并发大查询） |
| **数据模型** | 范式化（3NF，减少冗余） | 反范式化/星型模型（宽表，减少 JOIN） |
| **存储模型** | **行存储** | **列存储** |
| **索引需求** | B+ 树索引（快速定位单行） | 稀疏索引/位图索引/Zone Map |
| **事务要求** | 强 ACID | 弱事务或无事务 |
| **典型场景** | 电商下单、银行转账、用户注册 | 销售报表、用户画像、实时大盘 |

---

### 常见数据库分类

#### OLTP 数据库

| 数据库 | 存储引擎 | 核心特点 | 典型场景 |
|--------|---------|---------|---------|
| **MySQL (InnoDB)** | B+ 树 | 最流行的关系型 DB，ACID 事务 | Web 应用、电商、金融 |
| **PostgreSQL** | B+ 树 | 功能最强的开源 RDBMS，支持 JSON/GIS | 企业应用、地理信息 |
| **Oracle** | B+ 树 | 企业级商业数据库 | 金融、电信、政府 |
| **TiDB** | LSM 树(TiKV) | 分布式 NewSQL，兼容 MySQL 协议 | 大规模 OLTP + 轻量 OLAP |
| **CockroachDB** | LSM 树(Pebble) | 分布式 SQL，强一致，兼容 PostgreSQL | 全球部署、云原生 |
| **MongoDB** | B 树 / WiredTiger | 文档型 NoSQL，Schema 灵活 | 内容管理、IoT、游戏 |
| **Redis** | 内存(Hash/SkipList) | 内存 KV，极速 | 缓存、会话、排行榜 |

#### OLAP 数据库

| 数据库 | 存储引擎 | 核心特点 | 典型场景 |
|--------|---------|---------|---------|
| **ClickHouse** | MergeTree（列存 + LSM 思想） | 单表分析极快，百亿行秒级 | 日志分析、实时报表、用户行为 |
| **Doris/StarRocks** | 列存 + MPP | 兼顾实时导入和查询，MPP 架构 | 实时数仓、多维分析 |
| **Hive** | 行/列存(ORC/Parquet) on HDFS | Hadoop 生态 SQL 引擎 | 离线数仓、批处理 |
| **Spark SQL** | 列存(Parquet) | 内存计算，批+流处理 | ETL、机器学习 |
| **Greenplum** | 行/列存 + MPP | PostgreSQL 内核，分布式分析 | 传统数仓 |
| **Elasticsearch** | 倒排索引 + 列存(Doc Values) | 全文搜索 + 聚合分析 | 日志搜索、监控 |
| **Druid** | 列存 + 位图索引 | 实时 OLAP，亚秒级聚合 | 实时监控、广告分析 |

#### HTAP 数据库（混合）

| 数据库 | 存储引擎 | 特点 |
|--------|---------|------|
| **TiDB** | TiKV(行存LSM) + TiFlash(列存) | OLTP 走 TiKV，OLAP 走 TiFlash，同一份数据 |
| **PolarDB-X** | InnoDB + 列存引擎 | 阿里云，OLTP+OLAP 融合 |
| **CockroachDB** | LSM + 列存扩展 | 分布式 SQL，轻量分析 |

---

### 存储引擎原理详解

### 一、B+ 树 — OLTP 的基石

#### 结构

```
B+ 树是一种多路平衡搜索树:

  ┌──────────────────────────────────────────────┐
  │                  根节点                        │
  │              [30 | 60]                        │
  │             /    |    \                       │
  │            /     |     \                      │
  │   ┌──────┐  ┌──────┐  ┌──────┐               │
  │   │10|20 │  │40|50 │  │70|80 │   ← 非叶子节点  │
  │   └┬─┬─┬┘  └┬─┬─┬┘  └┬─┬─┬┘     (只存索引)   │
  │    │ │ │    │ │ │    │ │ │                    │
  │   ┌▼┐┌▼┐┌▼┐┌▼┐┌▼┐┌▼┐┌▼┐┌▼┐┌▼┐               │
  │   │ ││ ││ ││ ││ ││ ││ ││ ││ │ ← 叶子节点      │
  │   │ ││ ││ ││ ││ ││ ││ ││ ││ │   (存数据/指针) │
  │   └─┘└─┘└─┘└─┘└─┘└─┘└─┘└─┘└─┘               │
  │    ←──→←──→←──→←──→←──→←──→←──→              │
  │       叶子节点之间用双向链表相连                  │
  └──────────────────────────────────────────────┘

关键特点:
  1. 所有数据都在叶子节点（非叶子只存索引 key）
  2. 叶子节点用双向链表相连（范围查询高效）
  3. 树高度很矮: 3~4 层可存千万~亿行
  4. 每个节点 = 一个磁盘页（16KB in InnoDB）
```

#### 为什么 OLTP 用 B+ 树

```
  点查 (WHERE id=5):
    根 → 非叶子 → 叶子 → 找到记录
    3~4 次磁盘 I/O → 毫秒级 ✓

  范围查询 (WHERE id BETWEEN 10 AND 50):
    找到 id=10 的叶子 → 沿链表顺序扫描到 id=50
    顺序 I/O，高效 ✓

  插入/更新:
    找到位置 → 原地修改或页分裂
    随机 I/O 但次数少（树很矮）✓

  不适合:
    全表扫描（百万行聚合）→ 每个叶子节点都要读 → 慢
    → 这就是 OLAP 不用 B+ 树的原因
```

#### InnoDB 的 B+ 树存储

```
聚簇索引 (Clustered Index):
  叶子节点存完整行数据
  按主键排序

  主键 B+ 树:
    [PK=1] → {id:1, name:"Alice", age:20, ...}  ← 完整行
    [PK=2] → {id:2, name:"Bob", age:25, ...}

二级索引 (Secondary Index):
  叶子节点存主键值（不存完整行）
  查到主键后"回表"到聚簇索引取完整数据

  name 索引 B+ 树:
    [name="Alice"] → PK=1  ← 只存主键
    [name="Bob"]   → PK=2
    → 回表: 用 PK=1 去聚簇索引查完整行
```

---

### 二、LSM 树 — 写密集型 OLTP 和 OLAP 的核心

#### 问题：B+ 树的写放大

```
B+ 树的写入问题:
  修改一条记录 → 读取整个 16KB 的页 → 修改 → 写回整个页
  即使只改 100 字节，也要写 16KB → 写放大 160 倍
  随机写（不同页散布在磁盘各处）→ HDD 上极慢

LSM 树的解决思路:
  所有写入先到内存（顺序写！）→ 攒一批后批量刷盘（也是顺序写！）
  → 把随机写转换为顺序写 → 写入吞吐量极高
```

#### LSM 树结构

```
  ┌──────────────────────────────────────────────────────┐
  │  写入                                                 │
  │   ↓                                                   │
  │  ┌──────────┐                                         │
  │  │ MemTable │  内存（跳表/红黑树），有序写入             │
  │  │ (活跃)    │  写入速度极快，容量有限（几 MB ~ 几十 MB） │
  │  └────┬─────┘                                         │
  │       │ 满了                                           │
  │       ▼                                               │
  │  ┌───────────┐                                        │
  │  │ Immutable │  冻结的 MemTable，只读，等待刷盘          │
  │  │ MemTable  │                                        │
  │  └────┬──────┘                                        │
  │       │ 刷盘（flush）                                  │
  │       ▼                                               │
  │  ┌──────────────────────────────────────────┐         │
  │  │              磁盘（SSTable）               │         │
  │  │                                          │         │
  │  │  Level 0:  [SST] [SST] [SST]   ← 可能重叠│         │
  │  │              ↓  Compaction                │         │
  │  │  Level 1:  [SST] [SST] [SST] [SST]       │         │
  │  │              ↓  Compaction                │         │
  │  │  Level 2:  [SST] [SST] [SST] [SST] [SST] │         │
  │  │              ↓                            │         │
  │  │  Level N:  ...（越底层数据量越大）          │         │
  │  └──────────────────────────────────────────┘         │
  └──────────────────────────────────────────────────────┘

SSTable (Sorted String Table):
  有序的、不可变的键值对文件
  内部有索引块 + 数据块 + 布隆过滤器
```

#### 读写流程

```
写入:
  1. 写 WAL（Write-Ahead Log，防宕机丢内存数据）
  2. 写 MemTable（内存，O(logN)）
  3. MemTable 满了 → 变成 Immutable → 后台刷盘成 SSTable (Level 0)
  4. Level 0 的 SSTable 积累到阈值 → 合并到 Level 1（Compaction）
  → 全程顺序写，写入极快

读取:
  1. 先查 MemTable（最新数据在内存）
  2. 再查 Immutable MemTable
  3. 再从 Level 0 → Level 1 → ... → Level N 依次查 SSTable
  4. 每层用布隆过滤器快速排除不存在的 key
  → 读取可能需要查多层 → 读放大

Compaction（合并压实）:
  将多个 SSTable 合并为一个更大的有序 SSTable
  → 消除重复 key，清理已删除的数据
  → 减少层数，降低读放大
  → 但 Compaction 本身消耗 I/O 和 CPU（写放大的来源）
```

#### LSM 树的三种放大

```
  写放大 (Write Amplification):
    一条数据经过多次 Compaction 被重复写入
    Level 0 → Level 1 → Level 2 ... 每次合并都要写一次
    典型写放大: 10~30 倍

  读放大 (Read Amplification):
    点查一个 key 可能需要查多层 SSTable
    用布隆过滤器和索引缓解

  空间放大 (Space Amplification):
    同一个 key 的旧版本在 Compaction 前仍占用空间
    → 实际占用 > 有效数据量

  B+ 树 vs LSM 树的放大对比:
         │ 写放大 │ 读放大 │ 空间放大 │
  B+ 树  │  中等  │  低    │  低      │  适合读多写少
  LSM 树 │  高    │  中等  │  中等    │  适合写多读少
```

#### 使用 LSM 树的数据库

```
  TiKV (TiDB 的存储层)   — RocksDB (LSM)
  CockroachDB            — Pebble (LSM)
  Cassandra              — LSM
  HBase                  — LSM
  LevelDB / RocksDB      — LSM（被很多系统嵌入使用）
  MongoDB WiredTiger      — 支持 B 树和 LSM 两种模式
```

---

### 三、列存储 — OLAP 的基石

#### 行存 vs 列存

```
原始数据:
  id │ name  │ age │ city
  1  │ Alice │ 20  │ Beijing
  2  │ Bob   │ 25  │ Shanghai
  3  │ Carol │ 30  │ Beijing

行存储（Row-oriented）:
  磁盘上连续存储每一行:
  [1, Alice, 20, Beijing] [2, Bob, 25, Shanghai] [3, Carol, 30, Beijing]

  ✓ 读取完整行快（一次 I/O 拿到一行所有列）
  ✓ 适合 OLTP: SELECT * WHERE id=1
  ✗ 聚合慢: SELECT AVG(age) → 读每行只为了取 age 字段，但把 name/city 也读了

列存储（Column-oriented）:
  磁盘上连续存储每一列:
  id列:   [1, 2, 3]
  name列: [Alice, Bob, Carol]
  age列:  [20, 25, 30]
  city列: [Beijing, Shanghai, Beijing]

  ✓ SELECT AVG(age) → 只读 age 列，I/O 量极小
  ✓ 同列数据类型相同 → 压缩率极高（10:1 ~ 100:1）
  ✓ 适合 OLAP: 大范围扫描 + 聚合
  ✗ 读取完整行需要拼装多个列文件 → 点查慢
```

#### 列存的关键优化

```
1. 压缩 — 同列数据类型相同，压缩效果极好:

   字典编码:
     city列: [Beijing, Shanghai, Beijing, Beijing, Shanghai]
     → 字典: {0: Beijing, 1: Shanghai}
     → 压缩: [0, 1, 0, 0, 1]  (几个字节代替长字符串)

   游程编码 (RLE):
     排序后 city列: [Beijing, Beijing, Beijing, Shanghai, Shanghai]
     → 压缩: [(Beijing, 3), (Shanghai, 2)]

   位图编码:
     city=Beijing: [1, 0, 1, 1, 0]
     city=Shanghai: [0, 1, 0, 0, 1]
     → WHERE city='Beijing' → 直接用位图过滤

   压缩率: 列存通常 10:1 ~ 100:1（行存通常 2:1 ~ 5:1）

2. 向量化执行 — 批量处理而非逐行处理:

   行存引擎:
     for each row:
       if row.age > 20: sum += row.age   → 每行一次虚函数调用

   列存引擎:
     age_column = load_column("age")
     mask = age_column > 20              → SIMD 批量比较 1024 个值
     result = sum(age_column[mask])      → SIMD 批量求和

   → 利用 CPU SIMD 指令（AVX2/AVX-512），一次处理 8~16 个值
   → 减少分支预测失败，CPU 缓存友好

3. Zone Map / Min-Max 索引:

   每个列的数据块记录最小值和最大值:
     age 块1: min=18, max=25
     age 块2: min=26, max=40
     age 块3: min=35, max=60

   WHERE age > 50 → 块1 和块2 的 max < 50 → 跳过（数据裁剪）
   → 只需要读块3
```

---

### 四、MergeTree — ClickHouse 的存储引擎

```
ClickHouse 的 MergeTree 融合了列存 + LSM 的思想:

  ┌──────────────────────────────────────────────┐
  │  INSERT (批量写入)                             │
  │   ↓                                           │
  │  ┌──────────┐                                 │
  │  │  Part 1   │  每次 INSERT 生成一个 Part       │
  │  │ (列文件)  │  包含: 列数据 + 索引 + 标记文件   │
  │  └──────────┘                                 │
  │  ┌──────────┐                                 │
  │  │  Part 2   │                                │
  │  └──────────┘                                 │
  │       ↓ 后台合并（类似 LSM Compaction）         │
  │  ┌──────────────────┐                         │
  │  │  Merged Part      │  合并后数据按主键排序     │
  │  │ age.bin  name.bin │  每列一个文件，独立压缩   │
  │  │ age.mrk  name.mrk │  标记文件（稀疏索引）     │
  │  │ primary.idx       │  主键稀疏索引            │
  │  └──────────────────────────────────────────┘ │
  └──────────────────────────────────────────────┘

特点:
  - 列存储: 每列独立文件，独立压缩
  - 稀疏索引: 每 8192 行记录一个索引点（不是每行一条）
  - 后台合并: Part 合并消除重复，按排序键合并
  - 不支持事务/UPDATE/DELETE（或非常有限）
  - 面向批量写入，不适合高频单条 INSERT
```

---

### 五、倒排索引 — Elasticsearch

```
倒排索引: 从"词"到"文档"的映射（与正排相反）

  正排索引（行存）:
    doc1 → "Go is fast and efficient"
    doc2 → "Java is popular"
    doc3 → "Go and Java are both great"

  倒排索引:
    "Go"        → [doc1, doc3]
    "is"        → [doc1, doc2]
    "fast"      → [doc1]
    "Java"      → [doc2, doc3]
    "popular"   → [doc2]
    "great"     → [doc3]
    ...

  搜索 "Go fast":
    "Go"  → {doc1, doc3}
    "fast" → {doc1}
    交集 → doc1 ✓

适合: 全文搜索、关键词匹配
不适合: 数值范围聚合（ES 用 Doc Values 列存补充）
```

---

### 存储引擎对比总结

| 存储引擎 | 数据结构 | 写入特点 | 读取特点 | 代表数据库 | 适用类型 |
|---------|---------|---------|---------|-----------|---------|
| **B+ 树** | 多路平衡树 | 随机写，原地更新 | 点查极快 O(logN) | MySQL, PostgreSQL | OLTP |
| **LSM 树** | 内存表+分层SST | 顺序写，批量刷盘 | 点查需查多层 | TiKV, RocksDB, HBase | 写密集 OLTP |
| **列存** | 列式文件+压缩 | 批量追加 | 大范围扫描+聚合极快 | ClickHouse, Parquet | OLAP |
| **MergeTree** | 列存+LSM 合并 | 批量写入，后台合并 | 列扫描+稀疏索引 | ClickHouse | OLAP |
| **倒排索引** | 词→文档映射 | 写入需分词建索引 | 关键词搜索极快 | Elasticsearch | 搜索 |

```
选型决策:

  点查为主 + 事务 + 中等写入 → B+ 树 (MySQL/PostgreSQL)
  写入极高 + 点查 + 分布式 → LSM 树 (TiDB/Cassandra)
  大表聚合分析 + 报表 → 列存 (ClickHouse/Doris)
  全文搜索 + 日志分析 → 倒排索引 (Elasticsearch)
  OLTP + OLAP 都要 → HTAP (TiDB: TiKV 行存 + TiFlash 列存)
```

### 一句话总结

**OLTP 用行存 B+ 树（点查快、事务强），OLAP 用列存（扫描快、压缩好、向量化计算）。** B+ 树适合随机读写少量行（MySQL/PG）；LSM 树把随机写转为顺序写，适合写密集场景（TiKV/RocksDB）；列存将同列数据连续存放，聚合查询只读需要的列，压缩率 10~100 倍（ClickHouse/Doris）。现代 HTAP 数据库（如 TiDB）同时维护行存（TiKV）和列存（TiFlash）副本，一份数据同时服务 OLTP 和 OLAP。

---

## 六、什么是索引下推？作用是什么？

### 索引下推（ICP, Index Condition Pushdown）

```
MySQL 5.6 引入的优化，将部分 WHERE 条件的判断从 Server 层"下推"到存储引擎层（InnoDB）执行。
```

### 先理解 MySQL 的两层架构

```
  SQL 请求
      │
  ┌───▼────────────────┐
  │    Server 层        │  SQL 解析、优化、WHERE 过滤、排序、聚合
  └───┬────────────────┘
      │  调用存储引擎接口: "给我下一行数据"
  ┌───▼────────────────┐
  │  存储引擎层 (InnoDB) │  B+ 树索引查找、读取数据页、返回行
  └────────────────────┘

  传统流程: 存储引擎负责用索引找行 → 返回给 Server 层 → Server 层再用 WHERE 过滤
  索引下推: 把部分 WHERE 过滤也放到存储引擎层做 → 减少返回给 Server 层的行数
```

### 具体示例

```sql
-- 表结构
CREATE TABLE users (
    id       INT PRIMARY KEY,
    name     VARCHAR(50),
    age      INT,
    city     VARCHAR(50),
    INDEX idx_name_age (name, age)  -- 联合索引
);

-- 查询
SELECT * FROM users WHERE name LIKE '张%' AND age = 25;
```

#### 没有索引下推（MySQL 5.6 之前）

```
联合索引 idx_name_age (name, age):
  索引 B+ 树中存储: name + age + 主键id

  ① 存储引擎用索引匹配 name LIKE '张%'
     → 在索引中找到所有 name 以"张"开头的记录
     → age=25 这个条件虽然在索引中有，但存储引擎不管！
     → name LIKE '张%' 匹配到 5 条

  ② 拿到 5 条记录的主键 → 回表查完整行数据
     → 5 次回表（随机 I/O）

  ③ 将 5 条完整行返回给 Server 层

  ④ Server 层用 age=25 过滤
     → 只有 2 条满足 → 返回 2 条结果

  流程:
  存储引擎:  索引匹配 name LIKE '张%' → 找到 5 条 → 回表 5 次 → 返回 5 行
  Server层:  过滤 age=25 → 保留 2 行
                               ↑
                          3 次回表白做了！
```

#### 有索引下推（MySQL 5.6+）

```
  ① 存储引擎用索引匹配 name LIKE '张%'
     → 在索引中找到 5 条候选

  ② 在索引中直接检查 age=25（age 存在索引中！不需要回表）
     → 5 条中只有 2 条 age=25 → 淘汰 3 条
     → 这一步就是"下推"——把 age=25 的判断下推到存储引擎层

  ③ 只对 2 条记录回表
     → 2 次回表（减少了 3 次无效回表！）

  ④ 返回 2 条完整行给 Server 层

  流程:
  存储引擎:  索引匹配 name LIKE '张%' → 5 条候选
            索引中检查 age=25 → 过滤到 2 条（ICP！）
            回表 2 次 → 返回 2 行
  Server层:  直接返回 2 行
```

### 对比图

```
没有 ICP:
  索引 [张三,20] [张四,25] [张五,30] [张六,25] [张七,28]
       ↓ 全部匹配 name LIKE '张%'
  回表  回表       回表      回表       回表       → 5 次回表
       ↓          ↓         ↓         ↓          ↓
  Server 过滤 age=25:
       ✗          ✓         ✗         ✓          ✗  → 返回 2 行

有 ICP:
  索引 [张三,20] [张四,25] [张五,30] [张六,25] [张七,28]
       ↓ 匹配 name LIKE '张%' AND age=25（在索引中直接判断）
       ✗         ✓         ✗         ✓         ✗
                 ↓                    ↓
                回表                  回表        → 2 次回表
                 ↓                    ↓
                返回 2 行（无需 Server 层再过滤）
```

### ICP 的作用

```
核心作用: 减少回表次数

  回表 = 用主键去聚簇索引查完整行数据 = 随机 I/O = 最贵的操作

  没有 ICP: 只要索引前缀匹配就回表 → 可能大量无效回表
  有了 ICP: 在索引中能判断的条件先判断完，不满足的直接跳过 → 减少回表

收益:
  1. 减少回表次数 → 减少随机 I/O → 查询更快
  2. 减少 Server 层和存储引擎层之间的数据传输量
  3. 对大表查询效果显著（索引过滤掉的行越多，收益越大）
```

### ICP 触发条件

```
必须满足:
  1. 使用了二级索引（非主键索引）的范围查询或前缀匹配
  2. WHERE 条件中包含索引列，但这些列不能被索引完全覆盖
     （如果能覆盖索引就不需要回表了，ICP 无意义）
  3. InnoDB 或 MyISAM 引擎

不触发 ICP 的情况:
  - 主键/聚簇索引查询（不存在回表）
  - 覆盖索引（不需要回表，无需下推）
  - 子查询中的条件
  - 存储过程/触发器中的条件
```

### 如何判断是否使用了 ICP

```sql
EXPLAIN SELECT * FROM users WHERE name LIKE '张%' AND age = 25;

-- Extra 列出现 "Using index condition" → 使用了 ICP

-- 对比:
-- "Using where"           → Server 层过滤（没用 ICP）
-- "Using index"           → 覆盖索引（不需要回表，也不需要 ICP）
-- "Using index condition" → 索引下推 ✓
```

### 开启/关闭 ICP

```sql
-- 默认开启
SET optimizer_switch = 'index_condition_pushdown=on';  -- 开启（默认）
SET optimizer_switch = 'index_condition_pushdown=off'; -- 关闭

-- 关闭后对比执行计划，可以看到 Extra 从 "Using index condition" 变为 "Using where"
```

### ICP 与其他索引优化的关系

```
┌──────────────────────────────────────────────────────────────┐
│                  查询: WHERE a=1 AND b>5 AND c=10            │
│                  联合索引: (a, b, c)                          │
│                                                              │
│  最左前缀匹配:                                                │
│    a=1 → 用索引定位                                           │
│    b>5 → 用索引范围扫描                                       │
│    c=10 → 索引中 b 是范围查询，c 无法走索引精确匹配             │
│           (最左前缀在范围查询处"断开")                          │
│                                                              │
│  没有 ICP:                                                    │
│    索引找到 a=1 AND b>5 的所有行 → 全部回表 → Server 过滤 c=10 │
│                                                              │
│  有 ICP:                                                      │
│    索引找到 a=1 AND b>5 的候选 → 在索引中检查 c=10 → 过滤后回表│
│    c 虽然不能参与索引的树搜索，但 c 的值存在索引中               │
│    → ICP 可以利用索引中存储的 c 值来过滤！                      │
└──────────────────────────────────────────────────────────────┘

ICP 弥补了最左前缀"断开"后的不足:
  最左前缀决定"能用索引搜索到哪一步"
  ICP 决定"索引中已有的列能不能在回表前先过滤一轮"
```

### 一句话总结

**索引下推（ICP）= 把原本在 Server 层做的 WHERE 过滤，下推到存储引擎层在索引中直接做。** 核心收益是减少回表次数（随机 I/O）。联合索引中最左前缀"断开"后的列虽然不能参与 B+ 树搜索，但它们的值存在索引中，ICP 能利用这些值在回表前先过滤掉不满足条件的行。EXPLAIN 中 Extra 显示 `Using index condition` 即表示使用了 ICP。

---

## 七、MySQL 索引实现原理、B 树与 B+ 树区别、索引使用注意事项

### 一、为什么需要索引

```
没有索引: SELECT * FROM users WHERE name = 'Alice'
  → 全表扫描，逐行比较，1000 万行 = 1000 万次比较 = 几秒~几十秒

有索引: 同样的查询
  → B+ 树查找，3~4 次磁盘 I/O = 几毫秒

索引的本质:
  一种额外的数据结构，以空间换时间
  让数据库能快速定位到目标数据，避免全表扫描
```

---

### 二、B 树 vs B+ 树

#### B 树（B-Tree）

```
B 树: 多路平衡搜索树，所有节点都存储数据

              [30, 60]  ← 根节点（存数据+子指针）
             /    |    \
        [10,20] [40,50] [70,80]  ← 非叶子节点（也存数据）
        / | \   / | \   / | \
       叶子们  叶子们  叶子们  ← 叶子节点（也存数据）

  特点:
  - 所有节点（根、内部、叶子）都存数据
  - 查找可能在任意层命中（不一定到叶子）
  - 叶子节点之间没有链表
```

#### B+ 树

```
B+ 树: B 树的变种，数据只存在叶子节点

              [30 | 60]  ← 根节点（只存索引 key + 子指针）
             /    |    \
        [10|20] [40|50] [70|80]  ← 非叶子节点（只存索引 key）
         ↓  ↓    ↓  ↓    ↓  ↓
        [数据][数据][数据][数据][数据][数据]  ← 叶子节点（存数据）
         ←──→←──→←──→←──→←──→←──→
              叶子之间用双向链表相连

  特点:
  - 非叶子节点只存索引 key，不存数据
  - 所有数据都在叶子节点
  - 叶子节点之间用双向链表相连
  - 查找必须到叶子节点才能拿到数据
```

#### 核心区别对比

| 维度 | B 树 | B+ 树 |
|------|------|-------|
| **数据存储位置** | 所有节点都存数据 | **只有叶子节点存数据** |
| **非叶子节点** | 存索引 key + 数据 + 子指针 | **只存索引 key + 子指针** |
| **叶子节点链表** | 无 | **双向链表相连** |
| **查找路径** | 可能在中间节点命中 | **必须到叶子节点** |
| **单节点容纳 key 数** | 少（数据占空间） | **多（不存数据，空间全给 key）** |
| **树高度** | 相对高 | **更矮（同数据量）** |
| **范围查询** | 需要中序遍历（跨层） | **叶子链表顺序扫描（极快）** |
| **查询稳定性** | 不稳定（层数不同） | **稳定（都是叶子层）** |

#### 为什么 MySQL 选 B+ 树而非 B 树

```
原因 1: 非叶子节点能存更多 key → 树更矮 → 磁盘 I/O 更少

  InnoDB 页大小 = 16KB
  假设 key = 8B, 指针 = 6B, 数据行 = 1KB

  B 树非叶子节点: 每个 key 附带 1KB 数据
    16KB / (8 + 6 + 1000) ≈ 16 个 key/节点

  B+ 树非叶子节点: 只存 key + 指针
    16KB / (8 + 6) ≈ 1170 个 key/节点

  B+ 树 3 层能索引: 1170 × 1170 × 16 ≈ 2200 万行 ← 3 次 I/O
  B 树 3 层能索引:  16 × 16 × 16 ≈ 4096 行 ← 同样 3 次 I/O 但数据量差千倍

原因 2: 范围查询

  SELECT * FROM users WHERE age BETWEEN 20 AND 30;

  B 树: 找到 age=20 后，需要中序遍历整棵树 → 频繁上下跳层 → 随机 I/O
  B+ 树: 找到 age=20 的叶子 → 沿链表顺序扫描到 age=30 → 顺序 I/O

原因 3: 查询性能稳定

  B 树: 有的 key 在根节点命中（1次 I/O），有的在叶子（3次 I/O）→ 不稳定
  B+ 树: 所有查询都到叶子节点 → 恒定次数的 I/O → 可预测
```

---

### 三、InnoDB 索引实现

#### 聚簇索引（Clustered Index）

```
聚簇索引 = 数据和索引存在一起，叶子节点就是数据行

  主键 B+ 树:
           [30 | 60]
          /    |    \
     [10|20] [40|50] [70|80]
      ↓  ↓    ↓  ↓    ↓  ↓
  ┌───────┐ ┌───────┐ ┌───────┐
  │id=10  │ │id=40  │ │id=70  │
  │name=..│ │name=..│ │name=..│  ← 叶子节点存完整行数据
  │age=.. │ │age=.. │ │age=.. │
  └───────┘ └───────┘ └───────┘

  特点:
  - InnoDB 的主键索引就是聚簇索引
  - 表数据按主键顺序物理存储
  - 每张表只能有一个聚簇索引
  - 没有定义主键 → InnoDB 选第一个唯一非空索引 → 都没有则自动生成 row_id
```

#### 二级索引（Secondary Index）

```
二级索引 = 叶子节点存主键值，不存完整行

  name 索引 B+ 树:
           ["D" | "M"]
          /      |      \
     ["A"|"C"] ["F"|"K"] ["P"|"Z"]
       ↓   ↓     ↓   ↓     ↓   ↓
  ┌──────┐ ┌──────┐ ┌──────┐
  │name=A│ │name=F│ │name=P│
  │PK=5  │ │PK=12 │ │PK=3  │  ← 叶子节点只存 name + 主键 PK
  └──────┘ └──────┘ └──────┘

  查询 SELECT * FROM users WHERE name = 'Alice':
  ① 在 name 索引 B+ 树找到 name='Alice' → 拿到 PK=5
  ② 回表: 用 PK=5 去聚簇索引查完整行（第二次 B+ 树查找）

  回表 = 二级索引查主键 → 聚簇索引查完整行（多一次 B+ 树查找）
```

#### 覆盖索引（Covering Index）

```
如果查询的列全部包含在索引中 → 不需要回表 → 性能最优

  联合索引: idx_name_age (name, age)

  SELECT name, age FROM users WHERE name = 'Alice';
  → name 和 age 都在索引中 → 直接从索引返回 → 不回表！

  EXPLAIN 中 Extra: "Using index" → 覆盖索引 ✓

  实践:
    把高频查询需要的列都放进联合索引
    SELECT id, name, age 也算覆盖（主键 id 在所有索引中都有）
```

#### 联合索引与最左前缀原则

```
联合索引 idx_abc (a, b, c) 的 B+ 树按 (a, b, c) 排序:

  索引中的数据排列:
  (1, 1, 1)
  (1, 1, 2)
  (1, 2, 1)
  (1, 2, 3)
  (2, 1, 1)   ← 先按 a 排序，a 相同按 b，b 相同按 c
  (2, 1, 2)
  (2, 2, 1)

最左前缀原则:
  WHERE a=1                    → 走索引 ✓（用到 a）
  WHERE a=1 AND b=2            → 走索引 ✓（用到 a, b）
  WHERE a=1 AND b=2 AND c=3    → 走索引 ✓（用到 a, b, c）
  WHERE b=2                    → 不走索引 ✗（缺少最左列 a）
  WHERE b=2 AND c=3            → 不走索引 ✗（缺少最左列 a）
  WHERE a=1 AND c=3            → 走索引用到 a ✓，c 走 ICP ✓ 但不能索引精确查找

范围查询导致"断开":
  WHERE a=1 AND b>2 AND c=3
  → a=1 走索引 ✓
  → b>2 走索引范围扫描 ✓
  → c=3 不能走索引树搜索 ✗（b 是范围，c 在 b 范围内无序）
  → 但 c=3 可以走索引下推 ICP ✓（在索引中过滤，减少回表）
```

---

### 四、其他索引类型

```
哈希索引 (Hash Index):
  - Memory 引擎支持，InnoDB 有自适应哈希索引（AHI，自动构建）
  - 等值查询 O(1)，极快
  - 不支持范围查询、排序、最左前缀
  - 有哈希冲突问题

全文索引 (FULLTEXT):
  - InnoDB 5.6+ 支持
  - 用于全文搜索 MATCH AGAINST
  - 基于倒排索引实现
  - 实际项目中通常用 Elasticsearch 替代

前缀索引:
  - 对字符串的前 N 个字符建索引
  - ALTER TABLE users ADD INDEX idx_name (name(10));
  - 节省空间，但不能用于排序和覆盖索引
```

---

### 五、索引使用常见注意事项

#### 索引失效的常见场景

```
1. 对索引列使用函数/表达式:
  WHERE YEAR(created_at) = 2025       ✗ 索引失效
  WHERE created_at >= '2025-01-01'    ✓ 改写为范围查询

2. 隐式类型转换:
  phone 是 VARCHAR 类型
  WHERE phone = 13800138000           ✗ 数值和字符串比较，触发隐式转换
  WHERE phone = '13800138000'         ✓ 类型匹配

3. LIKE 左模糊:
  WHERE name LIKE '%Alice'            ✗ 左模糊无法走索引
  WHERE name LIKE 'Alice%'            ✓ 右模糊走索引

4. OR 连接非索引列:
  WHERE a = 1 OR b = 2               ✗ 如果 b 没有索引，全表扫描
  → 改用 UNION ALL 或给 b 加索引

5. NOT IN / NOT EXISTS / !=:
  WHERE status != 1                   通常不走索引（优化器判断全表更快）
  WHERE status IN (0, 2, 3)           ✓ 改写为 IN

6. 联合索引不满足最左前缀:
  INDEX (a, b, c)
  WHERE b = 1 AND c = 2              ✗ 缺少最左列 a

7. IS NULL / IS NOT NULL:
  取决于列的 NULL 比例，优化器自动判断
  大量 NULL → IS NOT NULL 可能走索引
  大量非 NULL → IS NULL 可能走索引
```

#### 索引设计原则

```
1. 选择性高的列优先建索引:
   选择性 = 不同值数量 / 总行数
   性别(男/女) 选择性 ≈ 0.5 → 不适合单独建索引
   手机号 选择性 ≈ 1.0 → 非常适合建索引

2. 联合索引的列顺序:
   ① 等值查询的列放前面（比范围查询的列靠前）
   ② 选择性高的列放前面
   ③ 频繁查询的列组合用一个联合索引覆盖

   WHERE status = 1 AND created_at > '2025-01-01'
   → INDEX (status, created_at) ✓  （等值在前，范围在后）
   → INDEX (created_at, status) ✗  （范围在前，status 走不了索引树搜索）

3. 覆盖索引减少回表:
   高频查询: SELECT id, name, status FROM orders WHERE status = 1;
   → INDEX (status, name)  ← 覆盖了 status + name + id(主键自带)

4. 前缀索引节省空间:
   长字符串列（如 URL、email）
   → INDEX (url(20))  ← 只索引前 20 个字符
   → 但不能用于 ORDER BY 和覆盖索引

5. 避免过多索引:
   每个索引 = 一棵 B+ 树 = 额外存储空间 + INSERT/UPDATE/DELETE 维护开销
   一般单表索引不超过 5~6 个
```

#### EXPLAIN 关键字段

```sql
EXPLAIN SELECT * FROM users WHERE name = 'Alice' AND age = 25;

重点关注:

  type — 访问类型（性能由好到差）:
    system > const > eq_ref > ref > range > index > ALL
    const:  主键/唯一索引等值（最好）
    ref:    非唯一索引等值
    range:  索引范围扫描
    index:  全索引扫描（遍历索引树，比全表好一点）
    ALL:    全表扫描（最差，需要优化！）

  key — 实际使用的索引（NULL = 没走索引）

  rows — 预估扫描行数（越少越好）

  Extra — 额外信息:
    Using index           → 覆盖索引 ✓
    Using index condition → 索引下推 ICP ✓
    Using where           → Server 层过滤（可能需要优化）
    Using filesort        → 额外排序（没用到索引排序）
    Using temporary       → 使用临时表（GROUP BY / DISTINCT 优化）
```

---

### 全景总结

```
InnoDB 索引体系:

  ┌─────────────────────────────────────────────────────────┐
  │                     聚簇索引                              │
  │                  （主键 B+ 树）                            │
  │              叶子节点 = 完整行数据                          │
  │              表数据按主键物理有序                            │
  ├─────────────────────────────────────────────────────────┤
  │   二级索引 1          二级索引 2          二级索引 3        │
  │  (name B+ 树)        (age B+ 树)       (name,age B+ 树) │
  │  叶子 = name+PK      叶子 = age+PK     叶子 = name+age+PK│
  │      │                   │                   │           │
  │      └── 回表 ──→ 聚簇索引                                │
  │                (覆盖索引则不需要回表)                       │
  └─────────────────────────────────────────────────────────┘

  查询路径:
    WHERE PK=5        → 聚簇索引直达（最快）
    WHERE name='A'    → 二级索引 → 回表聚簇索引
    WHERE name='A'    → 联合索引覆盖 → 不回表（最优化）
    无索引列          → 全表扫描（最慢）
```

### 一句话总结

**InnoDB 用 B+ 树而非 B 树：非叶子节点不存数据，单节点能放上千个 key（B 树只能放十几个），3 层即可索引 2000 万行；叶子节点双向链表范围查询顺序 I/O；所有查询到叶子层性能稳定可预测。** 聚簇索引叶子存完整行，二级索引叶子存主键（需回表），覆盖索引不回表性能最优。索引设计核心：等值列在前范围列在后、高选择性列优先、用联合索引实现覆盖、避免函数/隐式转换/左模糊导致索引失效。

---

## 八、一条 SQL 的执行流程

```
SELECT * FROM users WHERE id = 1;

  客户端
    │
    ▼
┌──────────────────── Server 层 ────────────────────┐
│                                                    │
│  ① 连接器 (Connector)                              │
│     验证用户名密码 → 查权限 → 建立连接              │
│     长连接/短连接、连接池                            │
│                                                    │
│  ② 查询缓存 (Query Cache) ← MySQL 8.0 已移除       │
│     SQL 文本完全匹配 → 直接返回缓存结果              │
│     命中率极低（表任何写操作都会失效）→ 已废弃         │
│                                                    │
│  ③ 解析器 (Parser)                                  │
│     词法分析: SQL 文本 → token (SELECT, FROM, WHERE) │
│     语法分析: token → AST 语法树                     │
│     语义检查: 表/列是否存在                           │
│                                                    │
│  ④ 优化器 (Optimizer)                               │
│     选择最优执行计划:                                │
│     - 选择用哪个索引（可能有多个候选）                 │
│     - 多表 JOIN 的连接顺序                           │
│     - 子查询优化（改写为 JOIN）                       │
│     → 生成执行计划（EXPLAIN 看到的就是这个）          │
│                                                    │
│  ⑤ 执行器 (Executor)                               │
│     检查权限 → 调用存储引擎接口逐行获取数据           │
│                                                    │
└──────────────┬─────────────────────────────────────┘
               │ 调用存储引擎 API
┌──────────────▼─────────────────────────────────────┐
│  ⑥ 存储引擎 (InnoDB)                               │
│     B+ 树索引查找 → Buffer Pool 读取数据页           │
│     → 不在内存则从磁盘加载                           │
│     → 返回满足条件的行给执行器                       │
└────────────────────────────────────────────────────┘
```

```
UPDATE 语句的额外流程:
  在⑥的基础上:
  → 写 Undo Log（旧值）
  → 修改 Buffer Pool 数据页
  → 写 Redo Log (prepare)
  → 写 Binlog
  → Redo Log (commit)
  → 返回 affected rows
```

---

## 九、慢查询排查与优化

### 排查流程

```
1. 开启慢查询日志:
   SET GLOBAL slow_query_log = ON;
   SET GLOBAL long_query_time = 1;  -- 超过 1 秒记录

2. 分析慢查询日志:
   mysqldumpslow -s t -t 10 /var/log/mysql/slow.log
   → 按耗时排序，取 Top 10

3. EXPLAIN 分析执行计划:
   重点看 type(ALL/index 需优化)、key(是否走索引)、rows(扫描行数)、Extra

4. SHOW PROFILE 分析时间分布:
   SET profiling = 1;
   SELECT ...;
   SHOW PROFILE FOR QUERY 1;
   → 看时间花在哪个阶段（Sending data / Sorting / Creating tmp table）
```

### 常见优化手段

```
索引优化（最有效）:
  - 给 WHERE/JOIN/ORDER BY/GROUP BY 的列加索引
  - 用 EXPLAIN 确认走了索引
  - 用联合索引覆盖查询，避免回表

SQL 改写:
  - 避免 SELECT *，只查需要的列
  - 小表驱动大表（IN vs EXISTS 选择）
  - LIMIT 分页深翻优化: LIMIT 100000, 10
    → 改写: WHERE id > last_id ORDER BY id LIMIT 10（游标分页）
  - OR 改 UNION ALL
  - 子查询改 JOIN

架构优化:
  - 读写分离（写主库读从库）
  - 缓存（Redis 缓存热数据）
  - 分库分表（数据量过大时）
```

---

## 十、分库分表方案与常见问题

### 什么时候需要分库分表

```
单表行数 > 2000 万 或 单表大小 > 10GB → 考虑分表
单库 QPS > 5000 或连接数不够 → 考虑分库

分库: 解决数据库连接数/QPS 瓶颈
分表: 解决单表数据量过大导致查询慢
```

### 分片策略

```
水平分表（最常用）:
  按某个字段（分片键）将数据分散到多张表
  user_0, user_1, user_2, ... user_15

  Hash 分片: table_id = hash(user_id) % 16
    → 均匀分布，但范围查询需扫描所有分片
  Range 分片: 按时间/ID 范围
    → 范围查询友好，但可能数据倾斜（热点）

垂直分表:
  将大表按列拆分: 常用列一张表，不常用列一张表
  user_base(id,name,phone) + user_detail(id,avatar,bio,address)
```

### 常见问题

```
1. 分布式 ID: 主键不能用自增（会冲突）→ Snowflake / Leaf

2. 跨分片查询:
   WHERE user_id = 123 → 路由到具体分片（简单）
   WHERE name = 'Alice' → 不知道在哪个分片 → 扫描所有分片（慢）
   → 解决: 建立 name→user_id 的映射表或 ES 辅助

3. 跨分片 JOIN: 无法直接 JOIN 不同分片
   → 应用层拼装 / 冗余数据 / 宽表

4. 跨分片事务: 分布式事务复杂
   → 最终一致性（消息队列）/ TCC / Seata

5. 扩容难: 16 分片扩到 32 分片 → 数据迁移
   → 一致性哈希减少迁移量 / 提前规划充足分片数

常见中间件: ShardingSphere、Vitess、MyCat
```

---

## 十一、主从复制原理与延迟处理

### 复制原理

```
Master                              Slave
  │                                   │
  │ 事务提交 → 写 Binlog               │
  │                                   │
  │ ←── IO Thread 拉取 Binlog ────────│ (Slave 的 IO 线程)
  │                                   │
  │ ──── Binlog Events ──────────────→│ → 写入 Relay Log
  │                                   │
  │                                   │ SQL Thread 读 Relay Log
  │                                   │ → 重放 SQL → 数据同步
  │                                   │

三种复制模式:
  异步复制（默认）:  Master 提交后立即返回，不等 Slave
    → 最快，但主宕机可能丢数据

  半同步复制:       Master 等至少一个 Slave 收到 Binlog 后才返回
    → 兼顾性能和安全

  全同步复制 (MGR): Master 等所有 Slave 都应用完才返回
    → 最安全，性能最差
```

### 主从延迟原因及解决

```
原因:
  1. Slave 单线程重放（5.6 前）→ Master 并发写，Slave 串行放
  2. 大事务（如大表 DDL / 批量 UPDATE）
  3. Slave 机器性能差
  4. 网络延迟

解决:
  1. 并行复制: slave_parallel_workers > 1（MySQL 5.7+ 基于组提交）
  2. 关键读走主库: 写后立即读的场景强制读主库
  3. 监控 Seconds_Behind_Master，延迟大时报警
  4. 半同步复制减少数据丢失
```

---

## 十二、MySQL 高可用方案

```
┌────────────────────────────────────────────────────────────┐
│ 方案            │ 原理                │ RPO  │ 切换时间     │
├────────────────────────────────────────────────────────────┤
│ 主从 + VIP 漂移 │ keepalived 检测主库  │ 可能丢│ 秒级        │
│ MHA             │ 自动选主+补齐 Binlog │ 少丢  │ 10~30s      │
│ MGR (组复制)    │ Paxos 多数派共识     │ 不丢  │ 秒级        │
│ InnoDB Cluster  │ MGR + MySQL Router  │ 不丢  │ 秒级(自动)   │
│ 云 RDS          │ 云厂商托管          │ 不丢  │ 30s(自动)    │
└────────────────────────────────────────────────────────────┘

MHA (Master High Availability):
  监控主库 → 主库宕机 → 从 Slave 中选一个提升为新主
  → 补齐差异 Binlog → 其他 Slave 指向新主 → VIP 切换
  → 最成熟的开源方案

MGR (MySQL Group Replication):
  多节点组成集群，Paxos 协议保证多数派写入
  → 数据不丢，自动故障检测和切换
  → MySQL 5.7.17+ 原生支持
```

---

## 十三、Buffer Pool 工作原理

```
Buffer Pool = InnoDB 的内存缓存，缓存数据页和索引页

  磁盘: [页1][页2][页3]...[页N]  ← .ibd 文件中的 16KB 页
                ↕ 按需加载
  Buffer Pool: [页2][页5][页8]...  ← 内存中的热数据页
                ↑ 所有读写都先经过 Buffer Pool

大小: innodb_buffer_pool_size（推荐物理内存的 60%~80%）
```

### 改进的 LRU 淘汰算法

```
InnoDB 的 LRU 不是普通 LRU，而是分成两段:

  ┌──── young 区（热数据，5/8）────┬──── old 区（冷数据，3/8）────┐
  │  [A] [B] [C] [D] [E]         │  [F] [G] [H] [I] [J]        │
  │  ← 频繁访问的页               │  ← 新加载的页先放这里          │
  └──── head ─────────────────────┴──── midpoint ──── tail ──────┘

  新页加载: 先放到 old 区头部（不是 LRU 头部！）
  页再次被访问且在 old 区停留超过 1s → 提升到 young 区头部
  淘汰: 从 tail 端淘汰

  为什么分两段？
    全表扫描（如 mysqldump）会加载大量页
    如果直接放 LRU 头 → 冲走所有热数据 → Buffer Pool 命中率暴跌
    分段后: 扫描的页在 old 区快速淘汰，不影响 young 区的热数据
```

---

## 十四、大表优化方案

```
场景: 单表 5000 万行，查询越来越慢

优化层次（由轻到重）:

  1. SQL + 索引优化（最先做）
     → 加合适索引、优化 SQL 写法、消除慢查询

  2. 读写分离
     → 写主库，读从库，分散压力

  3. 缓存
     → Redis 缓存热数据，减少 DB 查询

  4. 垂直分表
     → 大字段（TEXT/BLOB）拆到单独的表

  5. 冷热分离
     → 历史数据归档到另一张表/另一个库
     → 主表只保留近期数据

  6. 水平分表
     → 按 user_id hash 或按时间 range 分到多张表

  7. 分库
     → 分散到不同的 MySQL 实例

  8. 换引擎
     → 分析场景考虑 ClickHouse / TiDB / ES
```

---

## 十五、utf8 与 utf8mb4 的区别

```
MySQL 的 utf8 ≠ 真正的 UTF-8！

  utf8 (MySQL):
    最多 3 字节/字符
    不支持 4 字节的 Unicode 字符
    → Emoji (😀)、部分生僻汉字 → 存储失败或乱码！

  utf8mb4 (MySQL):
    真正的 UTF-8，最多 4 字节/字符
    完全兼容 Unicode
    → 支持 Emoji、所有汉字

  结论: 永远使用 utf8mb4，不要用 utf8
    CREATE TABLE ... DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

  排序规则（Collation）:
    utf8mb4_general_ci:  不区分大小写，速度快，精度略低
    utf8mb4_unicode_ci:  不区分大小写，更精确（按 Unicode 标准）
    utf8mb4_0900_ai_ci:  MySQL 8.0 默认，基于 Unicode 9.0，性能好
    utf8mb4_bin:         二进制比较，区分大小写
```

---

## 十六、数据库三大范式与反范式

```
第一范式 (1NF): 列不可再分（原子性）
  ✗ 地址: "北京市海淀区中关村"（一个字段包含省/市/区）
  ✓ 省: 北京市, 区: 海淀区, 街道: 中关村

第二范式 (2NF): 非主键列完全依赖主键（消除部分依赖）
  ✗ 订单表(订单ID, 商品ID, 商品名称, 数量)
    → 商品名称依赖商品ID，不依赖订单ID → 部分依赖
  ✓ 订单表(订单ID, 商品ID, 数量) + 商品表(商品ID, 商品名称)

第三范式 (3NF): 非主键列不传递依赖于主键（消除传递依赖）
  ✗ 员工表(员工ID, 部门ID, 部门名称)
    → 部门名称依赖部门ID，部门ID依赖员工ID → 传递依赖
  ✓ 员工表(员工ID, 部门ID) + 部门表(部门ID, 部门名称)

反范式: 故意冗余数据，减少 JOIN，提升查询性能
  典型场景:
    订单表中冗余商品名称 → 查订单列表不需要 JOIN 商品表
    用户表中冗余部门名称 → 显示时直接取，不 JOIN

  OLTP: 适度范式化（3NF），减少冗余保证一致性
  OLAP: 反范式化（宽表/星型模型），减少 JOIN 提升分析性能
```

---

## 十七、数据库连接池原理与配置

```
为什么需要连接池？
  创建 MySQL 连接: TCP 三次握手 + 认证 ≈ 10~50ms
  如果每次请求都新建连接 → 100 QPS × 30ms = 3s 浪费在建连上
  连接池: 预先创建一批连接复用 → 省去建连开销

工作原理:
  ┌─────────────────────────────────┐
  │          连接池                   │
  │  [conn1] [conn2] [conn3] [conn4]│  ← 预建好的连接
  └───┬────────┬────────┬───────────┘
      │        │        │
  goroutine1  goroutine2  goroutine3
  借出conn1   借出conn2   借出conn3
  用完归还    用完归还    用完归还
```

```go
// Go 标准库 database/sql 连接池配置
db, _ := sql.Open("mysql", dsn)

db.SetMaxOpenConns(100)    // 最大打开连接数（含空闲+活跃）
db.SetMaxIdleConns(20)     // 最大空闲连接数
db.SetConnMaxLifetime(time.Hour)    // 连接最大存活时间
db.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大存活时间

// 配置建议:
// MaxOpenConns: 根据 MySQL max_connections 和实例数量计算
//   如 MySQL max_connections=500, 5 个实例 → 每个 100
// MaxIdleConns: MaxOpenConns 的 20%~50%
// ConnMaxLifetime: 小于 MySQL wait_timeout（默认 8h）
```

---

## 十八、Online DDL 原理

```
问题: ALTER TABLE 大表（千万行）时，传统方式会锁表几十分钟甚至几小时

MySQL 5.6+ Online DDL 流程:
  1. 获取 MDL 写锁（极短暂）
  2. 降级为 MDL 读锁（允许 DML 继续）
  3. 在后台创建新表结构，拷贝数据
     → 期间的 DML 变更记录到 Row Log
  4. 数据拷贝完成 → 重新获取 MDL 写锁（极短暂）
  5. 应用 Row Log 中积累的变更
  6. 新表替换旧表 → 释放锁

支持 Online 的操作:
  ✓ 添加/删除索引
  ✓ 添加列（MySQL 8.0 Instant DDL 更快）
  ✓ 修改列默认值
  ✗ 修改列类型（需要重建表）
  ✗ 删除主键

MySQL 8.0 Instant DDL:
  ALTER TABLE users ADD COLUMN remark VARCHAR(255), ALGORITHM=INSTANT;
  → 只修改元数据，不拷贝数据，毫秒级完成！
  → 仅限在表末尾加列

第三方工具（大表推荐）:
  gh-ost (GitHub): 基于 Binlog 的在线 DDL，不用触发器
  pt-online-schema-change (Percona): 基于触发器
```

---

## 十九、Change Buffer 写缓冲

```
问题: 对非唯一二级索引的写操作（INSERT/UPDATE/DELETE）
     如果目标索引页不在 Buffer Pool 中
     → 需要从磁盘读取该页 → 随机 I/O → 慢

Change Buffer 的优化:
  如果目标索引页不在 Buffer Pool 中
  → 不立即读磁盘，而是将修改"缓存"在 Change Buffer 中
  → 下次该索引页被读取到 Buffer Pool 时，合并（merge）变更
  → 将多次随机 I/O 合并为一次

┌──────────────────────────────────────────────────┐
│  INSERT INTO users(name) VALUES ('Alice')         │
│                                                  │
│  name 索引页不在 Buffer Pool 中:                  │
│                                                  │
│  没有 Change Buffer:                              │
│    从磁盘读索引页 → 修改 → 写回 ← 随机 I/O!       │
│                                                  │
│  有 Change Buffer:                                │
│    记录到 Change Buffer → 不读磁盘 ← 省了随机 I/O! │
│    后续有人查询读到该页时 → merge 变更              │
└──────────────────────────────────────────────────┘

限制:
  - 只对非唯一二级索引有效
  - 唯一索引必须立即读取索引页判断是否冲突 → 无法延迟
  - Change Buffer 占用 Buffer Pool 空间（默认 25%，可调）

配置:
  innodb_change_buffer_max_size = 25   -- 占 Buffer Pool 的百分比
  innodb_change_buffering = all        -- 缓冲哪些操作（inserts/deletes/changes/all）

适用场景:
  写多读少 + 大量非唯一二级索引 → Change Buffer 收益大
  写后立即读 → Change Buffer 没用（马上就要 merge）
```

---

### MySQL 面试知识全景图

```
┌─────────────────────── MySQL 面试知识体系 ───────────────────────┐
│                                                                  │
│  基础架构                   事务与并发                              │
│  ├─ SQL 执行流程(八)        ├─ ACID 原理(一)                       │
│  ├─ Buffer Pool(十三)       ├─ Redo/Undo/Binlog(二)               │
│  └─ Change Buffer(十九)     ├─ MVCC(三)                           │
│                             └─ 锁机制(四)                          │
│                                                                  │
│  索引体系                   高可用与扩展                            │
│  ├─ B+树原理(七)            ├─ 主从复制(十一)                       │
│  ├─ 索引下推(六)            ├─ 高可用方案(十二)                      │
│  └─ 使用注意事项(七)        ├─ 分库分表(十)                         │
│                             └─ Online DDL(十八)                   │
│                                                                  │
│  性能优化                   其他                                   │
│  ├─ 慢查询优化(九)          ├─ 字符集(十五)                         │
│  ├─ 大表优化(十四)          ├─ 范式与反范式(十六)                    │
│  └─ OLTP/OLAP(五)          └─ 连接池(十七)                        │
└──────────────────────────────────────────────────────────────────┘
```