package chapter34_35_sync_Map

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// TestSyncMap_1
// (1)关于同步工具的选择:
//	a.与单纯使用原生map和互斥锁的方案相比，使用sync.Map可以显著地减少锁的争用。sync.Map本身虽然也用到了锁，但是，它其实在尽可能地避免使用锁。
//	b.能用原子操作就不要用锁，不过这很有局限性，毕竟原子只能对一些基本的数据类型提供支持。
// (2)关于sync.Map的使用:
//	所有的方法涉及的键和值的类型都是interface{},Go语言的编译器并不会对它的键和值，进行特殊的类型检查。
// (3)Q:并发安全字典对键的类型有要求吗?
//	  A:有要求。键的实际类型不能是函数类型、字典类型和切片类型。而Go语言编译器是无法在编译期对它们进行检查的，不正确的键值实际类型肯定会引发panic。
//		所以在实际使用sync.Map时,必须保证键的类型是可判等的。如果你实在拿不准，那么可以先通过调用reflect.TypeOf函数得到一个键值对应的反射类型值
//		（即：reflect.Type类型的值），然后再调用这个值的Comparable方法，得到确切的判断结果。
func TestSyncMap_1(t *testing.T) {
	// 使用map+sync.Mutex实现sync.Map的全部功能
	pairs := []struct {
		k int
		v string
	}{
		{k: 1, v: "a"},
		{k: 2, v: "b"},
		{k: 3, v: "c"},
		{k: 4, v: "d"},
	}

	// 示例1。
	{
		cMap := NewConcurrentMap()
		cMap.Store(pairs[0].k, pairs[0].v)
		cMap.Store(pairs[1].k, pairs[1].v)
		cMap.Store(pairs[2].k, pairs[2].v)
		fmt.Println("[Three pairs have been stored in the ConcurrentMap instance]")

		cMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %v, %v\n",
				key, value)
			return true
		})

		k0 := pairs[0].k
		v0, ok := cMap.Load(k0)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v0, ok, k0)

		k3 := pairs[3].k
		v3, ok := cMap.Load(k3)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v3, ok, k3)

		k2, v2 := pairs[2].k, pairs[2].v
		actual2, loaded2 := cMap.LoadOrStore(k2, v2)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual2, loaded2, k2, v2)
		v3 = pairs[3].v
		actual3, loaded3 := cMap.LoadOrStore(k3, v3)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual3, loaded3, k3, v3)

		k1 := pairs[1].k
		cMap.Delete(k1)
		fmt.Printf("[The pair with the key of %v has been removed from the ConcurrentMap instance]\n",
			k1)
		v1, ok := cMap.Load(k1)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v1, ok, k1)
		v1 = pairs[1].v
		actual1, loaded1 := cMap.LoadOrStore(k1, v1)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual1, loaded1, k1, v1)

		cMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %v, %v\n",
				key, value)
			return true
		})
	}
	fmt.Println()

	// -------------------------------------------------------------------------------------
	// sync.Map的基本功能
	// 示例2。
	{
		var sMap sync.Map
		sMap.Store(pairs[0].k, pairs[0].v)
		sMap.Store(pairs[1].k, pairs[1].v)
		sMap.Store(pairs[2].k, pairs[2].v)
		fmt.Println("[Three pairs have been stored in the sync.Map instance]")

		sMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %v, %v\n",
				key, value)
			return true
		})

		k0 := pairs[0].k
		v0, ok := sMap.Load(k0)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v0, ok, k0)

		k3 := pairs[3].k
		v3, ok := sMap.Load(k3)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v3, ok, k3)

		k2, v2 := pairs[2].k, pairs[2].v
		actual2, loaded2 := sMap.LoadOrStore(k2, v2)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual2, loaded2, k2, v2)
		v3 = pairs[3].v
		actual3, loaded3 := sMap.LoadOrStore(k3, v3)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual3, loaded3, k3, v3)

		k1 := pairs[1].k
		sMap.Delete(k1)
		fmt.Printf("[The pair with the key of %v has been removed from the sync.Map instance]\n",
			k1)
		v1, ok := sMap.Load(k1)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v1, ok, k1)
		v1 = pairs[1].v
		actual1, loaded1 := sMap.LoadOrStore(k1, v1)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual1, loaded1, k1, v1)

		sMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %v, %v\n",
				key, value)
			return true
		})
	}
}

// TestSyncMap_2
// (4)Q:怎样保证并发安全字典中的键和值的类型正确性?
//    A:
//   	a.方案一是让并发安全字典只能存储某个特定类型的键。即在interface{}上包装一层，这样可以让编译器来帮助我们做检查。
//		  缺点是 缺少灵活性,如果我们还需要一个键类型为uint32并发安全字典的话，那就不得不再如法炮制地写一遍代码了。
//		  因此，在需求多样化之后，工作量反而更大，甚至会产生很多雷同的代码。
//		b.方案二是接受动态的类型设置(即参数仍是interface{})，并在程序运行的时候通过反射操作进行检查。
//		  缺点是 反射操作回降低程序的性能。
// (5)Q:sync.Map如何做到尽量避免使用锁? 其原理决定适用场景是读多写少的场景
//	  A:使用原子操作来存取键和值 + 使用两个原生的map作为存储介质(read+dirty) + 读写分离的思想（读不加锁,实现原子操作;写加锁）
//		read字典（只读字典）： 类型是atomic.Value
//		dirty字典（脏字典）：类型是map[interface{}]*entry
//      设计思想:
//			a.这两个字典无论是键还是值都是存的指针。所以脏字典和只读字典如果都存有同一个键值对，那么这里的两个键指的肯定是同一个基本值，对于两个值来说也是如此。
//			b.sync.Map在查找指定的键所对应的值时,会先去只读字典中寻找(不需要加互斥锁);只读字典查不到才会去脏字典查(需要加互斥锁)。
//			c.sync.Map在存储键值对时,若只读字典中已存有这个键，并且该键值对未被标记为“已删除”，就会把新值存到里面并直接返回(不需要加互斥锁);
//			  若只读字典中不存在这个键,则会把键值对存储到脏字典中(加互斥锁)
//			d.当一个键值对应该被删除，但却仍然存在于只读字典中的时候，才会被用标记为“已删除”的方式进行逻辑删除(把该键值对中指向值的那个指针置为nil)，而不会直接被物理删除。
//	    	  对于删除键值对，sync.Map会先去检查只读字典中是否有对应的键。如果没有，脏字典中可能有，那么它就会在锁的保护下，试图从脏字典中删掉该键值对。
//		    e.只读字典和脏字典之间是会互相转换的。在脏字典中查找键值对次数足够多的时候，sync.Map会把脏字典直接作为只读字典，保存在它的read字段中，
//		      然后把代表脏字典的dirty字段的值置为nil。一旦再有新的键值对存入，它就会依据只读字典去重建脏字典。这个时候，它会把只读字典中已被逻辑删除的键值对过滤掉。
//			f.在几个写操作当中，新增键值对的操作对并发安全字典的性能影响是最大的，其次是删除操作，最后才是修改操作。
// (6)Q:关于保证并发安全字典中的键和值的类型正确性，你还能想到其他的方案吗?
//    A:这是一道开放的问题，需要你自己去思考。其实怎样做完全取决于你的应用场景。不过，我们应该尽量避免使用反射，因为它对程序性能还是有一定的影响的。
func TestSyncMap_2(t *testing.T) {
	// 示例1。
	var sMap sync.Map
	//sMap.Store([]int{1, 2, 3}, 4) // 这行代码会引发panic。
	_ = sMap

	// A:方案一是让并发安全字典只能存储某个特定类型的键。即在interface{}上包装一层，这样可以让编译器来帮助我们做检查。
	// 示例2。
	{
		var iMap IntStrMap
		iMap.Store(pairs[0].k, pairs[0].v)
		iMap.Store(pairs[1].k, pairs[1].v)
		iMap.Store(pairs[2].k, pairs[2].v)
		fmt.Println("[Three pairs have been stored in the IntStrMap instance]")

		iMap.Range(func(key int, value string) bool {
			fmt.Printf("The result of an iteration in Range: %d, %s\n",
				key, value)
			return true
		})

		k0 := pairs[0].k
		v0, ok := iMap.Load(k0)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v0, ok, k0)

		k3 := pairs[3].k
		v3, ok := iMap.Load(k3)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v3, ok, k3)

		k2, v2 := pairs[2].k, pairs[2].v
		actual2, loaded2 := iMap.LoadOrStore(k2, v2)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual2, loaded2, k2, v2)
		v3 = pairs[3].v
		actual3, loaded3 := iMap.LoadOrStore(k3, v3)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual3, loaded3, k3, v3)

		k1 := pairs[1].k
		iMap.Delete(k1)
		fmt.Printf("[The pair with the key of %v has been removed from the IntStrMap instance]\n",
			k1)
		v1, ok := iMap.Load(k1)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v1, ok, k1)
		v1 = pairs[1].v
		actual1, loaded1 := iMap.LoadOrStore(k1, v1)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual1, loaded1, k1, v1)

		iMap.Range(func(key int, value string) bool {
			fmt.Printf("The result of an iteration in Range: %d, %s\n",
				key, value)
			return true
		})
	}
	fmt.Println()

	// b.方案二是封装的结构体中方法参数仍是interface{},但是在方法中添加一些做类型检查的代码。
	// 示例2。
	{
		cMap, err := NewConcurrentMap2(
			reflect.TypeOf(pairs[0].k), reflect.TypeOf(pairs[0].v))
		if err != nil {
			fmt.Printf("fatal error: %s", err)
			return
		}
		cMap.Store(pairs[0].k, pairs[0].v)
		cMap.Store(pairs[1].k, pairs[1].v)
		cMap.Store(pairs[2].k, pairs[2].v)
		fmt.Println("[Three pairs have been stored in the ConcurrentMap instance]")

		cMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %d, %s\n",
				key, value)
			return true
		})

		k0 := pairs[0].k
		v0, ok := cMap.Load(k0)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v0, ok, k0)

		k3 := pairs[3].k
		v3, ok := cMap.Load(k3)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v3, ok, k3)

		k2, v2 := pairs[2].k, pairs[2].v
		actual2, loaded2 := cMap.LoadOrStore(k2, v2)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual2, loaded2, k2, v2)
		v3 = pairs[3].v
		actual3, loaded3 := cMap.LoadOrStore(k3, v3)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual3, loaded3, k3, v3)

		k1 := pairs[1].k
		cMap.Delete(k1)
		fmt.Printf("[The pair with the key of %v has been removed from the ConcurrentMap instance]\n",
			k1)
		v1, ok := cMap.Load(k1)
		fmt.Printf("The result of Load: %v, %v (key: %v)\n",
			v1, ok, k1)
		v1 = pairs[1].v
		actual1, loaded1 := cMap.LoadOrStore(k1, v1)
		fmt.Printf("The result of LoadOrStore: %v, %v (key: %v, value: %v)\n",
			actual1, loaded1, k1, v1)

		cMap.Range(func(key, value interface{}) bool {
			fmt.Printf("The result of an iteration in Range: %d, %s\n",
				key, value)
			return true
		})
	}
}

// ConcurrentMap 代表自制的简易并发安全字典。
type ConcurrentMap struct {
	m  map[interface{}]interface{}
	mu sync.RWMutex
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m: make(map[interface{}]interface{}),
	}
}

func (cMap *ConcurrentMap) Delete(key interface{}) {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	delete(cMap.m, key)
}

func (cMap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	cMap.mu.RLock()
	defer cMap.mu.RUnlock()
	value, ok = cMap.m[key]
	return
}

func (cMap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	actual, loaded = cMap.m[key]
	if loaded {
		return
	}
	cMap.m[key] = value
	actual = value
	return
}

func (cMap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	cMap.mu.RLock()
	defer cMap.mu.RUnlock()
	for k, v := range cMap.m {
		if !f(k, v) {
			break
		}
	}
}

func (cMap *ConcurrentMap) Store(key, value interface{}) {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	cMap.m[key] = value
}

// ---------------------------------------------------------------------------
// IntStrMap 代表键类型为int、值类型为string的并发安全字典。
type IntStrMap struct {
	m sync.Map
}

func (iMap *IntStrMap) Delete(key int) {
	iMap.m.Delete(key)
}

func (iMap *IntStrMap) Load(key int) (value string, ok bool) {
	v, ok := iMap.m.Load(key)
	if v != nil {
		value = v.(string)
	}
	return
}

func (iMap *IntStrMap) LoadOrStore(key int, value string) (actual string, loaded bool) {
	a, loaded := iMap.m.LoadOrStore(key, value)
	actual = a.(string)
	return
}

func (iMap *IntStrMap) Range(f func(key int, value string) bool) {
	f1 := func(key, value interface{}) bool {
		return f(key.(int), value.(string))
	}
	iMap.m.Range(f1)
}

func (iMap *IntStrMap) Store(key int, value string) {
	iMap.m.Store(key, value)
}

// ConcurrentMap2 代表可自定义键类型和值类型的并发安全字典。
type ConcurrentMap2 struct {
	m         sync.Map
	keyType   reflect.Type
	valueType reflect.Type
}

func NewConcurrentMap2(keyType, valueType reflect.Type) (*ConcurrentMap2, error) {
	if keyType == nil {
		return nil, errors.New("nil key type")
	}
	if !keyType.Comparable() {
		return nil, fmt.Errorf("incomparable key type: %s", keyType)
	}
	if valueType == nil {
		return nil, errors.New("nil value type")
	}
	cMap := &ConcurrentMap2{
		keyType:   keyType,
		valueType: valueType,
	}
	return cMap, nil
}

func (cMap *ConcurrentMap2) Delete(key interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	cMap.m.Delete(key)
}

func (cMap *ConcurrentMap2) Load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		return
	}
	return cMap.m.Load(key)
}

func (cMap *ConcurrentMap2) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	actual, loaded = cMap.m.LoadOrStore(key, value)
	return
}

func (cMap *ConcurrentMap2) Range(f func(key, value interface{}) bool) {
	cMap.m.Range(f)
}

func (cMap *ConcurrentMap2) Store(key, value interface{}) {
	if reflect.TypeOf(key) != cMap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cMap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cMap.m.Store(key, value)
}

// pairs 代表测试用的键值对列表。
var pairs = []struct {
	k int
	v string
}{
	{k: 1, v: "a"},
	{k: 2, v: "b"},
	{k: 3, v: "c"},
	{k: 4, v: "d"},
}
