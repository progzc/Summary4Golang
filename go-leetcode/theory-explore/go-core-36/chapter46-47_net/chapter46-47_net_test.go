package chapter46_47_net

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"runtime"
	"syscall"
	"testing"
	"time"
)

// TestNet_1
// (1)Q:进程间相互通信的方法?
//	  A:系统信号（signal）、管道（pipe）、套接字 （socket）、文件锁（file lock）、消息队列（message queue）、信号灯（semaphore，有的地方也称之为信号量）等。
// (2)系统调用
//    你可以理解为特殊的 C 语言函数。它们是连接应用程序和操作系统内核的桥梁，也是应用程序使用操作系统功能的唯一渠道。
// (3)Q:net.Dial 函数的第一个参数network有哪些可选值?
//	  A:net.Dial函数会接受两个参数，分别名为network和address，都是string类型的；参数network常用的可选值一共有 9 个。
//	    这些值分别代表了程序底层创建的 socket 实例可使用的不同通信协议，罗列如下：
//		a."tcp"：代表 TCP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
//		b."tcp4"：代表基于 IP 协议第四版的 TCP 协议。
//		c."tcp6"：代表基于 IP 协议第六版的 TCP 协议。
//		d."udp"：代表 UDP 协议，其基于的 IP 协议的版本根据参数address的值自适应。
//		e."udp4"：代表基于 IP 协议第四版的 UDP 协议。
//		f."udp6"：代表基于 IP 协议第六版的 UDP 协议。
//		g."unix"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_STREAM 为 socket 类型。
//		h."unixgram"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_DGRAM 为 socket 类型。
//		i."unixpacket"：代表 Unix 通信域下的一种内部 socket 协议，以 SOCK_SEQPACKET 为 socket 类型。
// (4)关于 syscall.Socket 函数
//	a.domain参数代表通信域（分别为IPv4域、IPv6域、Unix域）,即 syscall.AF_INET、syscall.AF_INET6、syscall.AF_UNIX
//	b.type参数代表数据类型（数据报文(如UDP)、数据流(如TCP)）,即 syscall.SOCK_DGRAM、syscall.SOCK_STREAM、syscall.SOCK_SEQPACKET、syscall.SOCK_RAW
//	c.proto参数代表使用的协议，通常，只要明确指定了前两个参数的值，我们就无需再去确定第三个参数值了，一般把它置为0就可以了。这时，内核程序会自行选择最合适的协议。
//    比如，当前两个参数值分别为syscall.AF_INET和syscall.SOCK_DGRAM的时候，内核程序会选择 UDP 作为协议。
//	  又比如，在前两个参数值分别为syscall.AF_INET6和syscall.SOCK_STREAM时，内核程序可能会选择 TCP 作为协议。
// (5)Q:用net.DialTimeout函数时给定的超时时间意味着什么？
//	  A:这里的超时时间，代表着函数为网络连接建立完成而等待的最长时间。这是一个相对的时间。它会由这个函数的参数timeout的值表示。
//		a.不论执行到哪一步，只要在绝对的超时时间达到的那一刻，网络连接还没有建立完成，该函数就会返回一个代表了 I/O 操作超时的错误值。
//		b.如果解析出的 IP 地址有多个，那么函数会串行或并发地尝试建立连接。但无论用什么样的方式尝试，函数总会以最先建立成功的那个连接为准。
//		c.它还会根据超时前的剩余时间，去设定针对每次连接尝试的超时时间，以便让它们都有适当的时间执行。
//		d.在net包中还有一个名为Dialer的结构体类型。该类型有一个名叫Timeout的字段，它与上述的timeout参数的含义是完全一致的。
//		  实际上，net.DialTimeout函数正是利用了这个类型的值才得以实现功能的。
// (6)Q:怎样在 net.Conn 类型的值上正确地设定针对读操作和写操作的超时时间？
//	  A:net.Conn类型有 3 个可用于设置超时时间的方法，分别是：SetDeadline、SetReadDeadline和SetWriteDeadline。
//		a.这三个方法的签名是一模一样的，只是名称不同罢了。它们都接受一个time.Time类型的参数，并都会返回一个error类型的结果。
//		  其中的SetDeadline方法是用来同时设置读操作超时和写操作超时的。
//		b.这三个方法都会针对任何正在进行以及未来将要进行的相应操作进行超时设定。因此，如果你要在一个循环中进行读操作或写操作的话，最好在每次迭代中都进行一次超时设定。
//		  否则，靠后的操作就有可能因触达超时时间而直接失败。另外，如果有必要，你应该再次调用它们并传入time.Time类型的零值来表达不再限定超时时间。
func TestNet_1(t *testing.T) {
	fd1, err := syscall.Socket(
		syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket error: %v\n", err)
		return
	}
	defer syscall.Close(fd1)
	fmt.Printf("The file descriptor of socket：%d\n", fd1)
}

func TestNet_2(t *testing.T) {
	network := "tcp"
	host := "google.cn"
	reqStrTpl := `HEAD / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: %s
User-Agent: Dialer/%s



`

	// 示例1。
	network1 := network + "4"
	address1 := host + ":80"
	fmt.Printf("Dial %q with network %q ...\n", address1, network1)
	conn1, err := net.Dial(network1, address1)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn1.Close()

	reqStr1 := fmt.Sprintf(reqStrTpl, host, runtime.Version())
	fmt.Printf("The request:\n%s\n", reqStr1)
	_, err = io.WriteString(conn1, reqStr1)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return
	}
	fmt.Println()

	reader1 := bufio.NewReader(conn1)
	line1, err := reader1.ReadString('\n')
	if err != nil {
		fmt.Printf("read error: %v\n", err)
		return
	}
	fmt.Printf("The first line of response:\n%s\n", line1)
	fmt.Println()

	// 示例2。
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}
	network2 := network
	address2 := host + ":443"
	fmt.Printf("Dial %q with network %q ...\n", address2, network2)
	conn2, err := tls.Dial(network2, address2, tlsConf)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn2.Close()

	reqStr2 := fmt.Sprintf(reqStrTpl, host, runtime.Version())
	fmt.Printf("The request:\n%s\n", reqStr2)
	_, err = io.WriteString(conn2, reqStr2)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return
	}

	reader2 := bufio.NewReader(conn2)
	line2, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Printf("read error: %v\n", err)
		return
	}
	fmt.Printf("The first line of response:\n%s\n", line2)
	fmt.Println()
}

func TestIo_3(t *testing.T) {
	type dailArgs struct {
		network string
		address string
		timeout time.Duration
	}
	dialArgsList := []dailArgs{
		{
			"tcp",
			"google.cn:80",
			time.Millisecond * 500,
		},
		{
			"tcp",
			"google.com:80",
			time.Second * 2,
		},
		{
			// 如果在这种情况下发生的错误是：
			// "connect: operation timed out"，
			// 那么代表着什么呢？
			//
			// 简单来说，此错误表示底层的socket在连接网络服务的时候先超时了。
			// 这时抛出的其实是'syscall.ETIMEDOUT'常量代表的错误值。
			"tcp",
			"google.com:80",
			time.Minute * 4,
		},
	}
	for _, args := range dialArgsList {
		fmt.Printf("Dial %q with network %q and timeout %s ...\n",
			args.address, args.network, args.timeout)
		ts1 := time.Now()
		conn, err := net.DialTimeout(args.network, args.address, args.timeout)
		ts2 := time.Now()
		fmt.Printf("Elapsed time: %s\n", time.Duration(ts2.Sub(ts1)))
		if err != nil {
			fmt.Printf("dial error: %v\n", err)
			fmt.Println()
			continue
		}
		defer conn.Close()
		fmt.Printf("The local address: %s\n", conn.LocalAddr())
		fmt.Printf("The remote address: %s\n", conn.RemoteAddr())
		fmt.Println()
	}
}
