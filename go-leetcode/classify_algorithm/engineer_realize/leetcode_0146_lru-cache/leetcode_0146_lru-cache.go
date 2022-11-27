package leetcode_0146_lru_cache

import "math"

// 146. LRU 缓存
// https://leetcode.cn/problems/lru-cache/

type LRUCache struct {
	cap  int
	size int
	m    map[int]*Node
	head *Node
	tail *Node
}

type Node struct {
	Key  int
	Val  int
	Pre  *Node
	Next *Node
}

func Constructor(capacity int) LRUCache {
	lRUCache := LRUCache{
		cap:  capacity,
		size: 0,
		m:    make(map[int]*Node),
		head: &Node{Key: math.MinInt32, Val: math.MinInt32},
		tail: &Node{Key: math.MinInt32, Val: math.MinInt32},
	}
	lRUCache.head.Next = lRUCache.tail
	lRUCache.tail.Pre = lRUCache.head
	return lRUCache
}

// Get 查询操作
// 时间复杂度: O(1)
// 空间复杂度: O(1)
func (this *LRUCache) Get(key int) int {
	v, ok := this.m[key]
	if !ok {
		return -1
	}

	if v.Pre != this.head {
		v.Pre.Next = v.Next
		v.Next.Pre = v.Pre

		v.Next = this.head.Next
		v.Pre = this.head
		this.head.Next.Pre = v
		this.head.Next = v
	}
	return v.Val
}

// Put 存放操作
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func (this *LRUCache) Put(key int, value int) {
	v, ok := this.m[key]
	if ok {
		v.Val = value
		if v.Pre != this.head {
			v.Pre.Next = v.Next
			v.Next.Pre = v.Pre

			v.Next = this.head.Next
			v.Pre = this.head
			this.head.Next.Pre = v
			this.head.Next = v
		}
		return
	}

	this.size++
	node := &Node{
		Key: key,
		Val: value,
	}
	node.Next = this.head.Next
	node.Pre = this.head
	this.head.Next.Pre = node
	this.head.Next = node
	this.m[key] = node

	if this.size > this.cap {
		tmp := this.tail.Pre
		this.tail.Pre.Pre.Next = this.tail
		this.tail.Pre = this.tail.Pre.Pre
		tmp.Pre = nil
		tmp.Next = nil
		delete(this.m, tmp.Key)
		this.size--
	}
}
