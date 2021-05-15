package main

import (
	"fmt"
	"time"
)

func contain(obj error, array []error) bool {
	for _, v := range array {
		if obj == v {
			return true
		}
	}
	return false
}

type CallBack func(...interface{}) (interface{}, error)

func Retry(attemptsCount int, delay time.Duration, exceptions []error) func(fn CallBack) CallBack {
	hasException := func(err error) bool {
		time.Sleep(delay)
		return true
	}
	if len(exceptions) > 0 {
		hasException = func(err error) bool {
			isContain := contain(err, exceptions)
			if isContain {
				time.Sleep(delay)
			}
			return isContain
		}
	}

	return func(fn CallBack) CallBack {
		return func(args ...interface{}) (interface{}, error) {
			var response interface{}
			var err error
			for i := 0; i < attemptsCount; i++ {
				response, err = fn(args)
				if err != nil && hasException(err) {
					continue
				}

				break
			}
			return response, err
		}
	}
}

var ERR = fmt.Errorf("test error")
var ERR2 = fmt.Errorf("best error")

func test(args ...interface{}) (interface{}, error) {
	fmt.Println("test")
	return 0, ERR
}

func test2(args ...interface{}) (interface{}, error) {
	fmt.Println("test2")
	return 0, ERR
}

func main() {
	fmt.Println("Retry for any errors")

	retryFn := Retry(3, time.Duration(1)*time.Second, []error{})
	retryDecoratedFn := retryFn(test)
	retryDecoratedFn2 := retryFn(test2)

	_, _ = retryDecoratedFn2()
	r, err := retryDecoratedFn()
	fmt.Println("Response: ", r)
	fmt.Println("Error:", err)

	fmt.Println("Retry for ERR errors")
	r, err = Retry(3, time.Duration(1)*time.Second, []error{ERR, ERR2})(test)()
	fmt.Println("Response: ", r)
	fmt.Println("Error:", err)

	fmt.Println("Retry for ERR2 errors")
	r, err = Retry(3, time.Duration(1)*time.Second, []error{ERR2})(test)()
	fmt.Println("Response: ", r)
	fmt.Println("Error:", err)
}
