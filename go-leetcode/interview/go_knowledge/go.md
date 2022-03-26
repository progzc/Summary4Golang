# 1. Go语法

## 1.1 slice类型nil和空切片的区别?

问题：var a int[]和a:=make([]int,0)或者a := []int{}的区别？

区别：

- 相同点
  - append：作用是相同的
  - 打印：都是[]
- 不同点
  - json序列化：前者是null，后者是[]

```go
func TestQuestion_1_1(t *testing.T) {
	var a1 []int
	a2 := []int{}
	fmt.Println(a1, a2) // [] []
	b1, _ := json.Marshal(a1)
	b2, _ := json.Marshal(a2)
	fmt.Println(string(b1), string(b2)) // null []
}
```

使用建议：

- 针对判断问题，这两者是不相等的。
- 在判断是否位空时，建议统一使用len(a)==0做判断。

注意事项：

- slice的底层数据结构

  ```go
  // runtime/slice.go
  type slice struct {
  	array unsafe.Pointer
  	len   int
  	cap   int
  }
  
  func TestQuestion_1_1_2(t *testing.T) {
  	var a1 []int
  	a2 := make([]int, 0) // 其slice.array的初始值其实是一个全局变量的值 var zerobase uintptr（位于runtime/malloc.go）
  	a3 := []string{}     // 其slice.array的初始值其实是一个全局变量的值 var zerobase uintptr（位于runtime/malloc.go）
  	p1 := (*slice)(unsafe.Pointer(&a1))
  	p2 := (*slice)(unsafe.Pointer(&a2))
  	p3 := (*slice)(unsafe.Pointer(&a3))
  	// a1, array=0, len=0, cap=0
  	fmt.Printf("a1, array=%x, len=%d, cap=%d\n", uintptr(p1.array), uintptr(p1.len), uintptr(p1.cap))
  	// a2, array=c000051e80, len=0, cap=0
  	fmt.Printf("a2, array=%x, len=%d, cap=%d\n", uintptr(p2.array), uintptr(p2.len), uintptr(p2.cap))
  	// a3, array=c000051e80, len=0, cap=0
  	fmt.Printf("a3, array=%x, len=%d, cap=%d\n", uintptr(p3.array), uintptr(p3.len), uintptr(p3.cap))
  }
  ```

## 1.2 map类型nil和空切片的区别？

问题：var a map[int]int和a:=make(map[int]int,0)或者a:=map[int]int{}的区别？

区别：

- 相同点
  - 查询元素：效果相同
  - 删除元素：效果相同
  - 打印：都是map[]
- 不同点
  - 添加元素：前者会panic，后者不会
  - json序列化：前者是null，后者是{}

```go
func TestQuestion_1_2(t *testing.T) {
	var a1 map[int]int
	a2 := map[int]int{}
	//a1[2] = 2 // 这里会panic
	a2[3] = 3
	delete(a1, 1)
	delete(a2, 1)
	fmt.Println(a1, a2) // map[] map[]
	b1, _ := json.Marshal(a1)
	b2, _ := json.Marshal(a2)
	fmt.Println(string(b1), string(b2)) // null {}
}
```

使用建议：

- 在判断是否位空时，建议统一使用len(a)==0做判断。

注意事项：

- map的底层数据结构

  ```go
  // runtime/map.go
  type hmap struct {
  	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
  	// Make sure this stays in sync with the compiler's definition.
  	count     int // # live cells == size of map.  Must be first (used by len() builtin)
  	flags     uint8
  	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
  	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
  	hash0     uint32 // hash seed
  
  	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
  	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
  	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)
  
  	extra *mapextra // optional fields
  }
  ```

# 2. 网络及通信协议

# 3. MySQL

# 4. Redis

# 5. Kafka

# 6. ElasticSearch

# 7. Kubernates







