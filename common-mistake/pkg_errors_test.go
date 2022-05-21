package common_mistake

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestErrWrap(t *testing.T) {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)       // oh noes: whoops
	fmt.Printf("%+v", err) // 类似于TestErrStack
}

func TestErrStack(t *testing.T) {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Printf("%+v", err)
}

func TestWithMessage(t *testing.T) {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Println(err) // oh noes: whoops
}
