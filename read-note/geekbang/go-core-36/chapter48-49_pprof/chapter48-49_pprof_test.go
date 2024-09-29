package chapter48_49_pprof

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/progzc/Summary4Golang/read-note/geekbang/go-core-36/chapter48-49_pprof/common"
	"github.com/progzc/Summary4Golang/read-note/geekbang/go-core-36/chapter48-49_pprof/common/op"
)

// TestPProf_1
// (1)性能分析工具包：
//
//	代码包（用来收集指标,生成概要文件）：runtime/pprof、net/http/pprof、runtime/trace
//	工具（用来解析概要文件）：go tool pprof、go tool trace
//	此外：go test命令也可以在程序测试完成后生成概要文件。
//
// (2)概要文件：在某一段时间内，对Go程序的相关指标进行多次采样后得到的概要信息。
//
//	a.CPU概要文件（CPU Profile）：其中的每一段独立的概要信息都记录着，在进行某一次采样的那个时刻，CPU 上正在执行的 Go 代码。
//	b.内存概要文件（Mem Profile）：其中的每一段概要信息都记载着，在某个采样时刻，正在执行的 Go 代码以及堆内存的使用情况，
//								这里包含已分配和已释放的字节数量和对象数量。
//	c.阻塞概要文件（Block Profile）：其中的每一段概要信息，都代表着 Go 程序中的一个 goroutine 阻塞事件。
//
// (3)Q:怎样让程序对 CPU 概要信息进行采样?
//
//	  A:这需要用到runtime/pprof包中的 API。更具体地说，在我们想让程序开始对 CPU 概要信息进行采样的时候，需要调用这个代码包中的StartCPUProfile函数，
//	    而在停止采样的时候则需要调用该包中的StopCPUProfile函数。
//		注意事项：
//		a.StartCPUProfile函数设定的采样频率总是固定的，即：100Hz（默认值，设置值不能超过1MHz）。也就是说，每秒采样100次，或者说每10毫秒采样一次。
//		b.在StartCPUProfile函数执行之后，一个新启用的 goroutine 将会负责执行 CPU 概要信息的收集和输出，直到runtime/pprof包中的StopCPUProfile函数被成功调用。
//		c.StopCPUProfile函数也会调用runtime.SetCPUProfileRate函数，并把参数值（也就是采样频率）设为0。这会让针对 CPU 概要信息的采样工作停止。
//	查看概要文件的信息：go tool pprof cpuprofile.out；然后进入交互界面,可按照提示操作。
func TestPProf_1(t *testing.T) {
	var (
		profileName = "blockprofile.out"
	)
	f, err := common.CreateFile("", profileName)
	if err != nil {
		fmt.Printf("CPU profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	if err = startCPUProfile(f); err != nil {
		fmt.Printf("CPU profile start error: %v\n", err)
		return
	}
	if err = common.Execute(op.BlockProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	stopCPUProfile()
}

// TestPProf_2
// (4)Q:怎样设定内存概要信息的采样频率？
//	  A:
//	    a.设定内存概要信息采样频率的方法很简单，只要为runtime.MemProfileRate变量赋值即可。

//		  这个变量的含义是，平均每分配多少个字节，就对堆内存的使用情况进行一次采样。如果把该变量的值设为0，那么，
//	 	  Go 语言运行时系统就会完全停止对内存概要信息的采样。该变量的缺省值是512 KB，也就是512千字节。

//		  如果你要设定这个采样频率，那么越早设定越好，并且只应该设定一次，否则就可能会对 Go 语言运行时系统的采样工作，造成不良影响。
//	   	  比如，只在main函数的开始处设定一次。
//		b.当我们想获取内存概要信息的时候，还需要调用runtime/pprof包中的WriteHeapProfile函数。该函数会把收集好的内存概要信息，
//	 	  写到我们指定的写入器中。WriteHeapProfile函数得到的内存概要信息并不是实时的，它是一个快照，是在最近一次的内存垃圾收集工作完成时产生的。
//		  如果你想要实时的信息，那么可以调用runtime.ReadMemStats函数。不过要特别注意，该函数会引起 Go 语言调度器的短暂停顿。
func TestPProf_2(t *testing.T) {
	var (
		profileName    = "memprofile.out"
		memProfileRate = 8
	)
	f, err := common.CreateFile("", profileName)
	if err != nil {
		fmt.Printf("memory profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	startMemProfile(memProfileRate)
	if err = common.Execute(op.MemProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	if err = stopMemProfile(f); err != nil {
		fmt.Printf("memory profile stop error: %v\n", err)
		return
	}
}

// TestPProf_3
// (5)Q:样获取到阻塞概要信息？
//	  A:
//		a.调用runtime包中的SetBlockProfileRate函数，即可对阻塞概要信息的采样频率进行设定。该函数有一个名叫rate的参数，它是int类型的。
//		  这个参数的含义是，只要发现一个阻塞事件的持续时间达到了多少个纳秒，就可以对其进行采样。如果这个参数的值小于或等于0，
//	 	  那么就意味着 Go 语言运行时系统将会完全停止对阻塞概要信息的采样。

//		  在runtime包中，还有一个名叫blockprofilerate的包级私有变量，它是uint64类型的。
//		  这个变量的含义是，只要发现一个阻塞事件的持续时间跨越了多少个 CPU 时钟周期，就可以对其进行采样。它的含义与我们刚刚提到的rate参数的含义非常相似。
//		  这两者的区别仅仅在于单位不同。runtime.SetBlockProfileRate函数会先对参数rate的值进行单位换算和必要的类型转换，然后，它会把换算结果用原子操作赋给blockprofilerate变量。

//		  由于此变量的缺省值是0，所以 Go 语言运行时系统在默认情况下并不会记录任何在程序中发生的阻塞事件。
//		b.当我们需要获取阻塞概要信息的时候，需要先调用runtime/pprof包中的Lookup函数并传入参数值"block"，从而得到一个*runtime/pprof.Profile类型的值（以下简称Profile值）。
//	  	  在这之后，我们还需要调用这个Profile值的WriteTo方法，以驱使它把概要信息写进我们指定的写入器中。
//
//		  debug参数主要的可选值有两个，即：0和1。当debug的值为0时，通过WriteTo方法写进写入器的概要信息仅会包含go tool pprof工具所需的内存地址，
//		  这些内存地址会以十六进制的形式展现出来。当该值为1时，相应的包名、函数名、源码文件路径、代码行号等信息就都会作为注释被加入进去。
//		  另外，debug为0时的概要信息，会经由 protocol buffers 转换为字节流。而在debug为1的时候，WriteTo方法输出的这些概要信息就是我们可以读懂的普通文本了。
//		  debug的值也可以是2。这时，被输出的概要信息也会是普通的文本，并且通常会包含更多的细节。至于这些细节都包含了哪些内容，
//		  那就要看我们调用runtime/pprof.Lookup函数的时候传入的是什么样的参数值了。
func TestPProf_3(t *testing.T) {
	var (
		profileName      = "blockprofile.out"
		blockProfileRate = 2
		debug            = 0
	)
	f, err := common.CreateFile("", profileName)
	if err != nil {
		fmt.Printf("block profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	startBlockProfile(blockProfileRate)
	if err = common.Execute(op.BlockProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	if err = stopBlockProfile(f, debug); err != nil {
		fmt.Printf("block profile stop error: %v\n", err)
		return
	}
}

// TestPProf_4
// (6)Q:runtime/pprof.Lookup函数的正确调用方式是什么？
//
//	  A:
//		a.runtime/pprof.Lookup函数（以下简称Lookup函数）的功能是，提供与给定的名称相对应的概要信息。这个概要信息会由一个Profile值代表。
//		  如果该函数返回了一个nil，那么就说明不存在与给定名称对应的概要信息。
//		b.runtime/pprof包已经为我们预先定义了 6 个概要名称。它们对应的概要信息收集方法和输出方法也都已经准备好了。我们直接拿来使用就可以了。
//		  它们是：goroutine、heap、allocs、threadcreate、block和mutex。
//			goroutine：该函数会利用相应的方法，收集到当前正在使用的所有 goroutine 的堆栈跟踪信息。注意，这样的收集会引起 Go 语言调度器的短暂停顿。
//					   当调用该函数返回的Profile值的WriteTo方法时，如果参数debug的值大于或等于2，那么该方法就会输出所有 goroutine 的堆栈跟踪信息。
//					   这些信息可能会非常多。如果它们占用的空间超过了64 MB（也就是64兆字节），那么相应的方法就会将超出的部分截掉。
//			heap/allocs：会收集与堆内存的分配和释放有关的采样信息。这实际上就是我们在前面讨论过的内存概要信息。在我们传入"allocs"的时候，后续的操作会与之非常的相似。
//						"heap"会使得被输出的内存概要信息默认以“在用空间”（inuse_space）的视角呈现，而"allocs"对应的默认视角则是“已分配空间”（alloc_space）。
//						“在用空间”是指，已经被分配但还未被释放的内存空间。在这个视角下，go tool pprof工具并不会去理会与已释放空间有关的那部分信息。
//						而在“已分配空间”的视角下，所有的内存分配信息都会被展现出来，无论这些内存空间在采样时是否已被释放。
//			threadcreate：去收集一些堆栈跟踪信息。这些堆栈跟踪信息中的每一个都会描绘出一个代码调用链，这些调用链上的代码都导致新的操作系统线程产生。
//			block：因争用同步原语而被阻塞的那些代码的堆栈跟踪信息（即阻塞概要信息）。
//			mutex：曾经作为同步原语持（例如：异步信号量和原子操作，而非 通道、互斥锁、条件变量、”WaitGroup“）有者的那些代码，它们的堆栈跟踪信息。
func TestPProf_4(t *testing.T) {
	prepare()
	dir, err := createDir()
	if err != nil {
		fmt.Printf("dir creation error: %v\n", err)
		return
	}
	for _, name := range profileNames {
		for _, debug := range debugOpts {
			err = genProfile(dir, name, debug)
			if err != nil {
				return
			}
			time.Sleep(time.Millisecond)
		}
	}
}

// --------------------------------------------------------------------------------------
func startCPUProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.StartCPUProfile(f)
}

func stopCPUProfile() {
	pprof.StopCPUProfile()
}

// --------------------------------------------------------------------------------------
func startMemProfile(memProfileRate int) {
	runtime.MemProfileRate = memProfileRate
}

func stopMemProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.WriteHeapProfile(f)
}

// --------------------------------------------------------------------------------------
func startBlockProfile(blockProfileRate int) {
	runtime.SetBlockProfileRate(blockProfileRate)
}

func stopBlockProfile(f *os.File, debug int) error {
	if f == nil {
		return errors.New("nil file")
	}
	//debug参数主要的可选值有两个，即：0和1。当debug的值为0时，通过WriteTo方法写进写入器的概要信息仅会包含go tool pprof工具所需的内存地址，
	//这些内存地址会以十六进制的形式展现出来。当该值为1时，相应的包名、函数名、源码文件路径、代码行号等信息就都会作为注释被加入进去。
	//另外，debug为0时的概要信息，会经由 protocol buffers 转换为字节流。而在debug为1的时候，WriteTo方法输出的这些概要信息就是我们可以读懂的普通文本了。
	return pprof.Lookup("block").WriteTo(f, debug)
}

// --------------------------------------------------------------------------------------
// profileNames 代表概要信息名称的列表。
var profileNames = []string{
	"goroutine",
	"heap",
	"allocs",
	"threadcreate",
	"block",
	"mutex",
}

// profileOps 代表为了生成不同的概要信息而准备的负载函数的字典。
var profileOps = map[string]common.OpFunc{
	"goroutine":    op.BlockProfile,
	"heap":         op.MemProfile,
	"allocs":       op.MemProfile,
	"threadcreate": op.BlockProfile,
	"block":        op.BlockProfile,
	"mutex":        op.BlockProfile,
}

// debugOpts 代表debug参数的可选值列表。
var debugOpts = []int{
	0,
	1,
	2,
}

func genProfile(dir string, name string, debug int) error {
	fmt.Printf("Generate %s profile (debug: %d) ...\n", name, debug)
	fileName := fmt.Sprintf("%s_%d.out", name, debug)
	f, err := common.CreateFile(dir, fileName)
	if err != nil {
		fmt.Printf("create error: %v (%s)\n", err, fileName)
		return err
	}
	defer f.Close()
	if err = common.Execute(profileOps[name], 10); err != nil {
		fmt.Printf("execute error: %v (%s)\n", err, fileName)
		return err
	}
	profile := pprof.Lookup(name)
	err = profile.WriteTo(f, debug)
	if err != nil {
		fmt.Printf("write error: %v (%s)\n", err, fileName)
		return err
	}
	return nil
}

func createDir() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(currDir, "profiles")
	err = os.Mkdir(path, 0766)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return path, nil
}

func prepare() {
	runtime.MemProfileRate = 8
	runtime.SetBlockProfileRate(2)
}
