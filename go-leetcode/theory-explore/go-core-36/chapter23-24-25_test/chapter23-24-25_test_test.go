package chapter23_24_25_test

import (
	"math"
	"testing"
	"time"
)

// TestTest_1
// (1)测试文件分类:功能测试（test）、基准测试（benchmark，也称性能测试），以及示例测试（example）
// (2)Q:Go 语言对测试函数的名称和签名都有哪些规定?
//    A:
//		a.对于功能测试函数来说，其名称必须以Test为前缀，并且参数列表中只应有一个*testing.T类型的参数声明
//			*testing.T 类型有哪些方法?testing.T的部分功能有(判定失败接口，打印信息接口)
//		b.对于性能测试函数来说，其名称必须以Benchmark为前缀，并且唯一参数的类型必须是*testing.B类型的()
//			*testing.B 类型有哪些方法?testing.B 拥有testing.T的全部接口，同时还可以统计内存消耗，指定并行数目和操作计时器等
//		c.对于示例测试函数来说，其名称必须以Example为前缀，但对函数的参数列表没有强制规定
// (3)Q:go test命令执行的主要测试流程是什么?
//	  A:准备工作(检验代码包或源码文件的有效性,判断标记是否合法)-->针对每个被测代码包，依次地进行构建、执行包中符合要求的测试函数，清理临时文件，打印测试结果
//		需要注意的是：对于每个被测代码包，go test命令会串行地执行测试流程中的每个步骤，这也上上句中【依次】的含义
// (4)基本知识
//	a.测试代码没有任何变动,那么go test命令直接把刚刚缓存测试成功的结果打印出来;一旦有任何变动，缓存数据就会失效，go test命令就会再次真正地执行操作
//	b.go 命令会定期地删除最近未使用的缓存数据;但是，如果你想手动删除所有的缓存数据，运行一下go clean -cache命令就好了
//	c.对于测试成功的结果，go 命令也是会缓存的。运行go clean -testcache将会删除所有的测试结果缓存
//	d.设置环境变量GODEBUG的值也可以稍稍地改变 go 命令的缓存行为。 比如，设置值为gocacheverify=1将会导致 go 命令绕过任何的缓存数据，
// 	  而真正地执行操作并重新生成所有结果，然后再去检查新的结果与现有的缓存数据是否一致。
//	e.对于失败测试的结果，go test命令并不会进行缓存
//	f.t.Log方法以及t.Logf方法的作用，就是打印常规的测试日志，只不过当测试成功的时候，go test命令就不会打印这类日志了
//	g.如果你想在测试结果中看到所有的常规测试日志，那么可以在运行go test命令的时候加入标记-v
//	h.如果你想在测试失败的同时打印失败测试日志，那么可以直接调用t.Error方法或者t.Errorf方法
//	i.常用的性能测试的命令：go test -bench=. -run=^$ puzzlers/article20/q3 表示只针对q3包下名称为空的功能测试函数（即不执行任何功能测试函数）
//	相关文档：https://pkg.go.dev/cmd/go#hdr-Testing_flags
//			https://golang.google.cn/cmd/go/#hdr-Testing_flags
//	示例：go test -bench=. -v -cpu=1,2,4 -count=4 -benchmem -run=^BenchmarkGetPrimes$
//		-bench=.	表示进行性能测试
//		-run=^$		表示需要执行哪些功能测试函数（正则匹配）
//		-cpu=8		表示指定最大P的数量(代表着Go语言运行时系统同时运行goroutine的能力),作用同runtime.GOMAXPROCS函数
//		-benchmem	表示在性能测试完成后打印内存分配统计信息
//		-benchtime	表示设定测试函数的执行时间上限
//		-count		表示用于重复执行测试函数的。它的值必须大于或等于0，并且默认值为1
//		-parallel	表示同一个被测代码包中的功能测试函数的最大并发执行数,默认值是测试运行时的最大P数量（这可以通过调用表达式runtime.GOMAXPROCS(0)获得）。
//					包含了t.Parallel方法调用的功能测试函数就会被go test命令并发地执行，而并发执行的最大数量正是由-parallel标记值决定的。
//					-parallel对性能测试时无效的；不过，b.RunParallel方法、b.SetParallelism方法和-cpu标记的联合运用可以使性能测试程序并发进行。

//		b.N			表示执行时间刚刚超过1s,被测试函数执行的次数
//		ns/op		表是单次执行被测函数的平均耗时
//		B/op		当添加-benchmem时测试结果中会出现，即 每次操作分配的字节数
//		allocs/op	当添加-benchmem时测试结果中会出现，即 每次操作进行了几次分配
// (5)Q:在编写示例测试函数的时候，我们怎样指定预期的打印内容?
//	  A:这个问题的答案就在testing代码包的文档中
// (6)Q:怎样设置-cpu标记的值，以及它会对测试流程产生什么样的影响?
//	  A:标记-cpu的值应该是一个正整数的列表,该列表的表现形式为:以英文半角逗号分隔的多个整数字面量,比如1,2,4
//		针对于此值中的每一个正整数，go test命令都会先设置最大 P 数量为该数，然后再执行测试函数。
// (7) 几个公式
//	对于性能测试函数:
//	性能测试函数的执行次数 = `-cpu`标记的值中正整数的个数 x `-count`标记的值 x 探索式执行中测试函数的实际执行次数
//	对于功能测试函数:
//	功能测试函数的执行次数 = `-cpu`标记的值中正整数的个数 x `-count`标记的值
// (8)Q:性能测试函数中的计时器是做什么用的?
//	  A:testing.B类型有这么几个指针方法：StartTimer、StopTimer和ResetTimer。
//		通过对b.StartTimer和b.StopTimer方法的联合运用，再去除掉任何一段代码的执行时间;
//		相比之下，b.ResetTimer方法的灵活性就要差一些了，它只能用于：去除在调用它之前那些代码的执行时间。
// (9)Q:怎样在测试的时候开启测试覆盖度分析？如果开启，会有什么副作用吗？
//	  A:go test命令可以接受-cover标记。该标记的作用就是开启测试覆盖度分析。
//	    不过，由于覆盖度分析开启之后go test命令可能会在程序被编译之前注释掉一部分源代码，
//	    所以，若程序编译或测试失败，那么错误报告可能会记录下与原始的源代码不对应的行号。
func TestTest_1(t *testing.T) {

}

func TestFail(t *testing.T) {
	//t.Fail() // 会让函数测试失败,但会继续向下执行
	t.FailNow() // 此调用会让当前的测试立即失败,但不会继续向下执行
	t.Log("Failed.")
}

func BenchmarkGetPrimes2(b *testing.B) {
	// 	b.StopTimer()与	b.StartTimer()之间的耗时不应计算在GetPrimes的性能测试中,所以这里通过计时器排除了这部分耗时
	// 默认性能测试一启动就会自动创建并启动定时器,所以这里要先停掉
	b.StopTimer()
	time.Sleep(time.Millisecond * 500) // 模拟某个耗时但与被测程序关系不大的操作。
	max := 10000
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		GetPrimes(max)
	}
}

func BenchmarkGetPrimes(b *testing.B) { // 执行时间 上限默认为1秒
	for i := 0; i < b.N; i++ {
		GetPrimes(1000)
	}
}

// GetPrimes 用于获取小于或等于参数max的所有质数。
// 本函数使用的是爱拉托逊斯筛选法（Sieve Of Eratosthenes）。
func GetPrimes(max int) []int {
	if max <= 1 {
		return []int{}
	}
	marks := make([]bool, max)
	var count int
	squareRoot := int(math.Sqrt(float64(max)))
	for i := 2; i <= squareRoot; i++ {
		if marks[i] == false {
			for j := i * i; j < max; j += i {
				if marks[j] == false {
					marks[j] = true
					count++
				}
			}
		}
	}
	primes := make([]int, 0, max-count)
	for i := 2; i < max; i++ {
		if marks[i] == false {
			primes = append(primes, i)
		}
	}
	return primes
}
