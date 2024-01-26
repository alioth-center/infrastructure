package concurrency

import (
	"fmt"
	"github.com/alioth-center/infrastructure/errors"
	"github.com/alioth-center/infrastructure/utils/values"
)

type ConcurrencyResult[T any] struct {
	err    error
	result T
}

type Promise[T any] chan ConcurrencyResult[T]

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

func Async[out any](fn func() out) (promise Promise[out]) {
	ch := make(chan ConcurrencyResult[out])
	result := ConcurrencyResult[out]{}
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

func Await[out any](promise Promise[out]) (result out, err error) {
	res, open := <-promise
	if !open {
		return values.Nil[out](), errors.NewPromiseCompletedError()
	}
	return res.result, res.err
}
