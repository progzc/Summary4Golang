package chapter44_45_os

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
)

// TestOs_1
// (1)os包中的API都是平台不想关的,主要可以帮助我们使用操作系统中的文件系统、权限系统、环境变量、系统进程以及系统信号。
// (2)Q:os.File 类型都实现了哪些io包中的接口?
//
//	  A:os.File 的指针类型则实现了很多io代码包中的接口?
//		a.io包中最核心的 3 个简单接口 io.Reader 、 io.Writer 和 io.Closer。
//		b.另外的 3 个简单接口，即：io.ReaderAt、io.Seeker 和 io.WriterAt。
//		c.实现了io包的 9 个扩展接口中的 7 个。
//		d.没有实现简单接口io.ByteReader和io.RuneReader。
//		e.没有实现分别作为这两者的扩展接口的io.ByteScanner和io.RuneScanner。
//
// (3)Q:怎样才能获得一个os.File类型的指针值?
//
//		  A:
//			a.os.Create 函数用于根据给定的路径创建一个新的文件；可以在该函数返回的File值之上，对相应的文件进行读操作和写操作。
//			  如果我们给定的路径上的某一级父目录并不存在，那么该函数就会返回一个*os.PathError类型的错误值，以表示“不存在的文件或目录”。
//			b.os.NewFile 函数，需要接受一个代表文件描述符的、uintptr类型的值，以及一个用于表示文件名的字符串值。
//			  注意，不要被这个函数的名称误导了，它的功能并不是创建一个新的文件，而是依据一个已经存在的文件的描述符（如标准错误输出），来新建一个包装了该文件的File值。
//			c.os.Open 函数会打开一个文件并返回包装了该文件的File值。 然而，该函数只能以只读模式打开文件。
//			d.os.OpenFile 函数, 其实是os.Create函数和os.Open函数的底层支持，它最为灵活。这个函数有 3 个参数，分别名为name、flag和perm。
//			  name指代的就是文件的路径。
//		      flag参数指操作模式，限定了操作文件的方式，如只读模式O_RDONLY、只写模式O_WRONLY、读写模式O_RDWR。
//	       perm参数指权限模式，可以控制文件的访问权限。
//				需要注意的是：perm参数是os.FileMode类型，而os.FileMode类型可远不只权限模式，它还可以代表文件模式（也可以称之为文件种类）
func TestOs_1(t *testing.T) {
	// (2)Q:os.File类型都实现了哪些io包中的接口?
	file1 := (*os.File)(nil)
	fileType := reflect.TypeOf(file1)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Type %T implements\n", file1)
	for _, t := range ioTypes {
		if fileType.Implements(t) {
			buf.WriteString(t.String())
			buf.WriteByte(',')
			buf.WriteByte('\n')
		}
	}
	output := buf.Bytes()
	output[len(output)-2] = '.'
	fmt.Printf("%s\n", output)

	// a.os.Create函数用于根据给定的路径创建一个新的文件；可以在该函数返回的File值之上，对相应的文件进行读操作和写操作。
	// 如果我们给定的路径上的某一级父目录并不存在，那么该函数就会返回一个*os.PathError类型的错误值，以表示“不存在的文件或目录”。
	fileName1 := "something1.txt"
	filePath1 := filepath.Join(os.TempDir(), fileName1)
	var paths []string
	paths = append(paths, filePath1)
	dir, _ := os.Getwd()
	paths = append(paths, filepath.Join(dir[:len(dir)-1], fileName1))
	for _, path := range paths {
		fmt.Printf("Create a file with path %s ...\n", path)
		_, err := os.Create(path)
		if err != nil {
			var underlyingErr string
			if _, ok := err.(*os.PathError); ok {
				underlyingErr = "(path error)"
			}
			fmt.Printf("error: %v %s\n", err, underlyingErr)
			continue
		}
		fmt.Println("The file has been created.")
	}
	fmt.Println()

	// b.os.NewFile函数，需要接受一个代表文件描述符的、uintptr类型的值，以及一个用于表示文件名的字符串值。
	// 注意，不要被这个函数的名称误导了，它的功能并不是创建一个新的文件，而是依据一个已经存在的文件的描述符（如标准错误输出），来新建一个包装了该文件的File值。
	fmt.Println("New a file associated with stderr ...")
	file3 := os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
	if file3 != nil {
		file3.WriteString(
			"The Go language program writes something to stderr.\n")
	}
	fmt.Println()

	// c.os.Open函数会打开一个文件并返回包装了该文件的File值。 然而，该函数只能以只读模式打开文件。
	fmt.Printf("Open a file with path %s ...\n", filePath1)
	file4, err := os.Open(filePath1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("Write something to the file ...")
	_, err = file4.WriteString("something")
	var underlyingErr string
	if _, ok := err.(*os.PathError); ok {
		underlyingErr = "(path error)"
	}
	fmt.Printf("error: %v %s\n", err, underlyingErr)
	fmt.Println()

	fmt.Printf("Open a file with path %s ...\n", filePath1)
	file5a, err := os.Open(filePath1)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf(
		"Is there only one file descriptor for the same file in the same process? %v\n",
		file5a.Fd() == file4.Fd())
	file5b := os.NewFile(file5a.Fd(), filePath1)
	fmt.Printf("Can the same file descriptor represent the same file? %v\n",
		file5b.Name() == file5a.Name())
	fmt.Println()

	// d.os.OpenFile函数, 其实是os.Create函数和os.Open函数的底层支持，它最为灵活。这个函数有 3 个参数，分别名为name、flag和perm。
	// name指代的就是文件的路径。
	// flag参数指操作模式，限定了操作文件的方式，如只读模式O_RDONLY、只写模式O_WRONLY、读写模式O_RDWR。
	// perm参数指权限模式，可以控制文件的访问权限。
	fmt.Printf("Reuse a file on path %s ...\n", filePath1)
	file6, err := os.OpenFile(filePath1, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	contents := "something"
	fmt.Printf("Write %q to the file ...\n", contents)
	n, err := file6.WriteString(contents)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Printf("The number of bytes written is %d.\n", n)
	}
}

// TestOs_2
// (4)Q:可应用于File值的操作模式都有哪些?
//
//	  A:针对File值的操作模式主要有只读模式、只写模式和读写模式。
//	    这些模式分别由常量os.O_RDONLY、os.O_WRONLY和os.O_RDWR代表,在我们新建或打开一个文件的时候，必须把这三个模式中的一个设定为此文件的操作模式。
//		除此之外，我们还可以为这里的文件设置额外的操作模式：
//			os.O_APPEND：当向文件中写入内容时，把新内容追加到现有内容的后边。
//			os.O_CREATE：当给定路径上的文件不存在时，创建一个新文件。
//			os.O_EXCL：需要与os.O_CREATE一同使用，表示在给定的路径上不能有已存在的文件。
//			os.O_SYNC：在打开的文件之上实施同步 I/O。它会保证读写的内容总会与硬盘上的数据保持同步。
//			os.O_TRUNC：如果文件已存在，并且是常规的文件，那么就先清空其中已经存在的任何内容。
func TestOs_2(t *testing.T) {
	fileName1 := "something2.txt"
	filePath1 := filepath.Join(os.TempDir(), fileName1)
	fmt.Printf("The file path: %s\n", filePath1)
	fmt.Println()

	// 示例1。
	contents0 := "OpenFile is the generalized open call."
	flagDescList := []flagDesc{
		{
			os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
			"os.O_WRONLY|os.O_CREATE|os.O_TRUNC",
		},
		{
			os.O_WRONLY,
			"os.O_WRONLY",
		},
		{
			os.O_WRONLY | os.O_APPEND,
			"os.O_WRONLY|os.O_APPEND",
		},
	}

	for i, v := range flagDescList {
		fmt.Printf("Open the file with flag %s ...\n", v.desc)
		file1a, err := os.OpenFile(filePath1, v.flag, 0666)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("The file descriptor: %d\n", file1a.Fd())

		contents1 := fmt.Sprintf("[%d]: %s ", i+1, contents0)
		fmt.Printf("Write %q to the file ...\n", contents1)
		n, err := file1a.WriteString(contents1)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("The number of bytes written is %d.\n", n)

		file1b, err := os.Open(filePath1)
		fmt.Println("Read bytes from the file ...")
		bytes, err := ioutil.ReadAll(file1b)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("Read(%d): %q\n", len(bytes), bytes)
		fmt.Println()
	}

	// 示例2。
	fmt.Println("Try to create an existing file with flag os.O_TRUNC ...")
	file2, err := os.OpenFile(filePath1, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("The file descriptor: %d\n", file2.Fd())

	fmt.Println("Try to create an existing file with flag os.O_EXCL ...")
	_, err = os.OpenFile(filePath1, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	fmt.Printf("error: %v\n", err)
}

// TestOs_3
// (5)Q:怎样设定常规文件的访问权限？
//
//	  A:perm参数是os.FileMode类型，而os.FileMode类型可远不只权限模式，它还可以代表文件模式（也可以称之为文件种类）
//		a.文件种类，如 os.ModeDir 、 os.ModeNamedPipe ......等等
//		b.perm最低的9个比特位才表示文件的权限，当我们拿到一个此类型的值时，可以把它和 os.ModePerm 常量的值做按位与操作。
//		  9位的权限含义：
//			从高到低，这 3 组分别表示的是文件所有者（也就是创建这个文件的那个用户）、文件所有者所属的用户组，以及其他用户对该文件的访问权限。
//			而对于每个组，其中的 3 个比特位从高到低分别表示读权限、写权限和执行权限。
//			如：八进制整数0777就表示，操作系统中的所有用户都对当前的文件有读、写和执行的权限；
//			   八进制整数0666则表示：所有用户都对当前文件有读和写的权限，但都没有执行的权限。
//
// (6)Q:怎样通过os包中的 API 创建和操纵一个系统进程？
//
//		  A:你可以从os包的FindProcess函数和StartProcess函数开始。前者用于通过进程 ID（pid）查找进程，后者用来基于某个程序启动一个进程。
//	     这两者都会返回一个*os.Process类型的值。该类型提供了一些方法：
//	    		比如，用于杀掉当前进程的Kill方法，又比如，可以给当前进程发送系统信号的Signal方法，以及会等待当前进程结束的Wait方法。
//			与此相关的还有os.ProcAttr类型、os.ProcessState类型、os.Signal类型，等等。你可以通过积极的实践去探索更多的玩法。
func TestOs_3(t *testing.T) {
	// 示例1。
	fmt.Printf("The mode for dir:\n%32b\n", os.ModeDir)
	fmt.Printf("The mode for named pipe:\n%32b\n", os.ModeNamedPipe)
	fmt.Printf("The mode for all of the irregular files:\n%32b\n", os.ModeType)
	fmt.Printf("The mode for permissions:\n%32b\n", os.ModePerm)
	fmt.Println()

	// 示例2。
	fileName1 := "something3.txt"
	filePath1 := filepath.Join(os.TempDir(), fileName1)
	fmt.Printf("The file path: %s\n", filePath1)

	argDescList := []argDesc{
		{
			"Create",
			os.O_RDWR | os.O_CREATE,
			0644,
		},
		{
			"Reuse",
			os.O_RDWR | os.O_TRUNC,
			0666,
		},
		{
			"Open",
			os.O_RDWR | os.O_APPEND,
			0777,
		},
	}

	defer os.Remove(filePath1)
	for _, v := range argDescList {
		fmt.Printf("%s the file with perm %o ...\n", v.action, v.perm)
		file1, err := os.OpenFile(filePath1, v.flag, v.perm)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		info1, err := file1.Stat()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		fmt.Printf("The file permissions: %o\n", info1.Mode().Perm())
	}
}

// ioTypes 代表了io代码包中的所有接口的反射类型。
var ioTypes = []reflect.Type{
	reflect.TypeOf((*io.Reader)(nil)).Elem(),
	reflect.TypeOf((*io.Writer)(nil)).Elem(),
	reflect.TypeOf((*io.Closer)(nil)).Elem(),

	reflect.TypeOf((*io.ByteReader)(nil)).Elem(),
	reflect.TypeOf((*io.RuneReader)(nil)).Elem(),
	reflect.TypeOf((*io.ReaderAt)(nil)).Elem(),
	reflect.TypeOf((*io.Seeker)(nil)).Elem(),
	reflect.TypeOf((*io.WriterTo)(nil)).Elem(),
	reflect.TypeOf((*io.ByteWriter)(nil)).Elem(),
	reflect.TypeOf((*io.WriterAt)(nil)).Elem(),
	reflect.TypeOf((*io.ReaderFrom)(nil)).Elem(),

	reflect.TypeOf((*io.ByteScanner)(nil)).Elem(),
	reflect.TypeOf((*io.RuneScanner)(nil)).Elem(),
	reflect.TypeOf((*io.ReadSeeker)(nil)).Elem(),
	reflect.TypeOf((*io.ReadCloser)(nil)).Elem(),
	reflect.TypeOf((*io.WriteCloser)(nil)).Elem(),
	reflect.TypeOf((*io.WriteSeeker)(nil)).Elem(),
	reflect.TypeOf((*io.ReadWriter)(nil)).Elem(),
	reflect.TypeOf((*io.ReadWriteSeeker)(nil)).Elem(),
	reflect.TypeOf((*io.ReadWriteCloser)(nil)).Elem(),
}

type flagDesc struct {
	flag int
	desc string
}

type argDesc struct {
	action string
	flag   int
	perm   os.FileMode
}
