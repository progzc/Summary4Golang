package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError1 struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError1(err error, messagef string, msgArgs ...interface{}) MyError1 {
	return MyError1{
		Inner:      err, //<1>
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),        // <2>
		Misc:       make(map[string]interface{}), // <3>
	}
}

func (err MyError) Error1() string {
	return err.Message
}

// "lowlevel" module

type LowLevelErr1 struct {
	error
}

func isGloballyExec1(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr1{wrapError(err, err.Error())} // <1>
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

// "intermediate" module

type IntermediateErr1 struct {
	error
}

func runJob1(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec1(jobBinPath)
	if err != nil {
		return err // <1>
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run() // <1>
}

func handleError1(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) // <3>
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob1("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr1); ok { // <1>
			msg = err.Error()
		}
		handleError1(1, err, msg) // <2>
	}
}
