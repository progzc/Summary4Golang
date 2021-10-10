package common_mistake

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// defer 函数的执行时机
// (1) 对 defer 延迟执行的函数，会在调用它的函数结束时执行，而不是在调用它的语句块结束时执行，注意区分开
//     例：比如在一个长时间执行的函数里，内部 for 循环中使用 defer 来清理每次迭代产生的资源调用，就会出现问题
func TestMistake_049(t *testing.T) {
	wrong049()
	right049()
}

func wrong049() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	dir := os.Args[1]
	start, err := os.Stat(dir)
	if err != nil || !start.IsDir() {
		os.Exit(2)
	}

	var targets []string
	filepath.Walk(dir, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fInfo.Mode().IsRegular() {
			return nil
		}
		targets = append(targets, fPath)
		return nil
	})

	for _, target := range targets {
		f, err := os.Open(target)
		if err != nil {
			fmt.Println("bad target:", target, "error:", err) //error:too many open files
			break
		}
		// (1) 对 defer 延迟执行的函数，会在调用它的函数结束时执行，而不是在调用它的语句块结束时执行，注意区分开
		defer f.Close() // 注意：在每次 for 语句块结束时，不会关闭文件资源
		// 使用 f 资源
	}
}

func right049() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	dir := os.Args[1]
	start, err := os.Stat(dir)
	if err != nil || !start.IsDir() {
		os.Exit(2)
	}

	var targets []string
	filepath.Walk(dir, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fInfo.Mode().IsRegular() {
			return nil
		}
		targets = append(targets, fPath)
		return nil
	})

	for _, target := range targets {
		func() {
			f, err := os.Open(target)
			if err != nil {
				fmt.Println("bad target:", target, "error:", err)
				return // 在匿名函数内使用 return 代替 break 即可
			}
			defer f.Close() // 匿名函数执行结束，调用关闭文件资源
			// 使用 f 资源
		}()
	}
}
