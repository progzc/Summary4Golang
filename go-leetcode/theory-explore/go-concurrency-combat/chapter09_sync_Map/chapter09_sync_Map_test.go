package chapter09_sync_Map

import (
	"sync"
	"testing"
	"time"
)

// TestMap_1
// (1)常规map的一些扩展
//	a.要想遍历map时有序，可以考虑使用有序的map：https://github.com/elliotchance/orderedmap
//	b.带有过期功能的map：https://github.com/zekroTJA/timedmap
//	c.使用红黑树实现的key有序的map：https://github.com/emirpasic/gods/tree/master/maps/treemap
// (2)使用map的常见错误
//	a.未初始化
//	b.并发读写
func TestMap_1(t *testing.T) {
	// a.未初始化
	type Counter struct {
		Website      string
		Start        time.Time
		PageCounters map[string]int
	}
	var c Counter
	c.Website = "baidu.com"
	c.PageCounters["/"]++

	// b.并发读写问题
	var m = make(map[int]int, 10) // 初始化一个map
	go func() {
		for {
			m[1] = 1 //设置key
		}
	}()

	go func() {
		for {
			_ = m[2] //访问这个map
		}
	}()
	select {}
}

// TestMap_2
// (3)如何实现线程安全的map类型?（性能依次增加，不过实际还是以性能测试为准）
//	a.方法一：使用读写锁sync.RWMutex。
//	b.方法三：使用sync.Map。
//	c.方法二：分片加锁（思想是尽量减少锁的粒度和锁的持有时间）：https://github.com/orcaman/concurrent-map
func TestMap_2(t *testing.T) {
	// 分片map的设计思想
	_ = New()
}

// TestSyncMap_1
// (4)sync.Map的使用场景（sync.Map 在生产环境中很少使用）
//	在以下两个场景中使用 sync.Map，会比使用 map+RWMutex 的方式，性能要好得多：
//	a.只会增长的缓存系统中，一个 key 只写入一次而被读很多次。
//	b.多个 goroutine 为不相交的键集读、写和重写键值对。
//	官方建议你针对自己的场景做性能评测，如果确实能够显著提高性能，再使用 sync.Map。
// (5)sync.Map的数据机构
//		type Map struct {
//			mu Mutex
//			// 基本上你可以把它看成一个安全的只读的map
//			// 它包含的元素其实也是通过原子操作更新的，但是已删除的entry就需要加锁操作了
//			read atomic.Value // readOnly
//
//			// 包含需要加锁才能访问的元素
//			// 包括所有在read字段中但未被expunged（删除）的元素以及新加的元素
//			dirty map[interface{}]*entry
//
//			// 记录从read中读取miss的次数，一旦miss数和dirty长度一样了，就会把dirty提升为read，并把dirty置空
//			misses int
//		}
//
//		type readOnly struct {
//			m       map[interface{}]*entry
//			amended bool // 当dirty中包含read没有的数据时为true，比如新增一条数据
//		}
//
//		// expunged是用来标识此项已经删掉的指针
//		// 当map中的一个项目被删除了，只是把它的值标记为expunged，以后才有机会真正删除此项
//		var expunged = unsafe.Pointer(new(interface{}))
//
//		// entry代表一个值
//		type entry struct {
//			p unsafe.Pointer // *interface{}
//		}
//	a.Store方法：Store 既可以是新增元素，也可以是更新元素。如果运气好的话，更新的是已存在的未被删除的元素，直接更新即可，不会用到锁。
//	            如果运气不好，需要更新（重用）删除的对象、更新还未提升的 dirty 中的对象，或者新增加元素的时候就会使用到了锁，这个时候，性能就会下降。
//	b.Load方法：我们从 read 中读取到了这个 key 对应的值，那么就不需要加锁了，性能会非常好。但是，如果请求的 key 不存在或者是新加的，就需要加锁从 dirty 中读取。
//    	       所以，读取不存在的 key 会因为加锁而导致性能下降，读取还没有提升的新值的情况下也会因为加锁性能下降。
//	c.Delete方法: 如果read中不存在，那么就需要从 dirty 中寻找这个项目。最终，如果项目存在就删除（将它的值标记为 nil）。
//	             如果read中存在，或者没有被标记为expunged，那么就将其值标记为expunged=unsafe.Pointer(new(interface{}))，达到逻辑删除的效果。
// (6)sync.Map的实现原理
//	a.空间换时间。通过冗余的两个数据结构（只读的read字段、可写的dirty），来减少加锁对性能的影响。对只读字段（read）的操作不需要加锁。
//	b.优先从 read 字段读取、更新、删除，因为对 read 字段的读取不需要锁。
//	c.动态调整。miss 次数多了之后，将 dirty 数据提升为 read，避免总是从 dirty 中加锁读取。
//	d.double-checking。加锁之后先还要再检查 read 字段，确定真的不存在才操作 dirty 字段。
//	e.延迟删除。删除一个键值只是打标记，只有在提升 dirty 字段为 read 字段的时候才清理删除的数据。
func TestSyncMap_1(t *testing.T) {
}

// -------------------------分片map的设计思想----------------------------
var SHARD_COUNT = 32

// 分成SHARD_COUNT个分片的map
type ConcurrentMap []*ConcurrentMapShared

// 通过RWMutex保护的线程安全的分片，包含一个map
type ConcurrentMapShared struct {
	items        map[string]interface{}
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// 创建并发map
func New() ConcurrentMap {
	m := make(ConcurrentMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &ConcurrentMapShared{items: make(map[string]interface{})}
	}
	return m
}

// 根据key计算分片索引
func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
