package chapter08_container

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
	"testing"
)

// TestContainer_1 List的设计原理
// container包中的List:
// (1) 底层数据结构: 双向链表（内部具体实现其实是双向循环链表）
// (2) 常用api: 参见https://pkg.go.dev/container/list
//  a.修改方法: MoveBefore/MoveAfter/MoveToFront/MoveToBack针对形参中的e *Element(给的的元素)只能是链表中已存在的元素,而不是自己生成的元素（否则不会对链表做出任何改动）
// 	b.插入方法: InsertBefore/InsertAfter（没有延迟初始化,但是会判断指定节点是不是在所属链表上,通过在Element中嵌入*List达到O(1)的时间复杂度）
//	           PushFront/PushBack（有延迟初始化）可以将自己生成的interface元素插入到链表中
//	c.查询方法: Front/Back
// (3) var l list.List 声明的变量l的零值是一个长度为0的链表, 这样的链表可以"开箱即用",可以这么做的原因是"延迟初始化"（如下）,"延迟初始化"的缺点是当方法被频繁调用时会影响程序性能
// func (l *List) PushFront(v interface{}) *Element {
//	l.lazyInit()
//	return l.insertValue(v, &l.root)
// }
// (4) 切片其实也使用到了延迟初始化
func TestContainer_1(t *testing.T) {
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println()

	var l2 list.List
	l2.PushFront(1)
	l2.PushFront(2)
	for e := l2.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

// -----------------------------------------------------------------------------------------------
// TestContainer_2 Ring和List的区别
// container包中的Ring:
// (1) Ring的底层数据结构: 循环链表
// (2) Ring和List的区别
//	a.Ring类型的数据结构仅由它自身即可代表，而List类型则需要由它以及Element类型联合表示。这是表示方式上的不同，也是结构复杂度上的不同
//  b.一个Ring类型的值严格来讲，只代表了其所属的循环链表中的一个元素，而一个List类型的值则代表了一个完整的链表。这是表示维度上的不同
//  c.在创建并初始化一个Ring值的时候，我们可以指定它包含的元素的数量，但是对于一个List值来说却不能这样做（也没有必要这样做）。
// 循环链表一旦被创建，其长度是不可变的。这是两个代码包中的New函数在功能上的不同，也是两个类型在初始化值方面的第一个不同。
//  d.仅通过var r ring.Ring语句声明的r将会是一个长度为1的循环链表，而List类型的零值则是一个长度为0的链表(即本质上都支持延迟初始化)。
// 别忘了List中的根元素不会持有实际元素值，因此计算长度时不会包含它。这是两个类型在初始化值方面的第二个不同。
//  e.Ring值的Len方法的算法复杂度是 O(N) 的，而List值的Len方法的算法复杂度则是 O(1) 的。这是两者在性能方面最显而易见的差别。
func TestContainer_2(t *testing.T) {
	// Create a new ring of size 5
	r := ring.New(5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p interface{}) {
		fmt.Println(p.(int))
	})
}

// -----------------------------------------------------------------------------------------------
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // 说明是小顶堆
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TestContainer_3 小顶堆的使用
func TestContainer_3(t *testing.T) {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0]) // minimum: 1
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h)) // 1 2 3 5
	}
}

// -----------------------------------------------------------------------------------------------
// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// TestContainer_4 优先队列的使用
func TestContainer_4(t *testing.T) {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value) // 05:orange 04:pear 03:banana 02:apple
	}
}
