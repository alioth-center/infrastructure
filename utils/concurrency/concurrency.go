package concurrency

import (
	"fmt"
	"github.com/alioth-center/infrastructure/errors"
	"github.com/alioth-center/infrastructure/utils/values"
)

type ConcurrentResult[T any] struct {
	err    error
	result T
}

type Promise[T any] chan ConcurrentResult[T]

// RecoverErr recover panic to error, must be used in defer
// example:
//
//	func main() {
//		defer func() {
//			if err := RecoverErr(recover()); err != nil {
//				fmt.Println(err)
//			}
//		}()
//
//		panic("test")
//	}
//
// then
// output: test
func RecoverErr(e any) error {
	if e == nil {
		return nil
	}

	switch e.(type) {
	case error:
		return e.(error)
	default:
		return fmt.Errorf("%v", e) // nolint:goerr
	}
}

// Async async execute function
// example:
//
//	func main() {
//		promise := Async(func() string {
//			return "test"
//		})
//
//		result, err := Await(promise)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		fmt.Println(result)
//	}
//
// then
// output: test
func Async[out any](fn func() out) (promise Promise[out]) {
	ch := make(chan ConcurrentResult[out])
	result := ConcurrentResult[out]{}
	go func() {
		defer func() {
			result.err = RecoverErr(recover())
			ch <- result
			close(ch)
		}()
		result.result = fn()
	}()
	return ch
}

// Await await promise
// example:
//
//	func main() {
//		promise := Async(func() string {
//			return "test"
//		})
//
//		result, err := Await(promise)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		fmt.Println(result)
//	}
//
// then
// output: test
func Await[out any](promise Promise[out]) (result out, err error) {
	res, open := <-promise
	if !open {
		return values.Nil[out](), errors.NewPromiseCompletedError()
	}
	return res.result, res.err
}
