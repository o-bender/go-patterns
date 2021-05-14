package main

import (
	"fmt"
	"time"
)

func simpleFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}
	return simpleFactorial(n-1) * n
}

type SUBFN func() (float64, SUBFN)
type FN func(float64) (float64, SUBFN)

func trampoline(fn FN) func(float64) float64 {
	return func(n float64) float64 {
		result, call := fn(n)
		for call != nil {
			result, call = call()
		}
		return result
	}
}

func factorial(f float64, n float64) (float64, SUBFN) {
	if n == 0 {
		return f, nil
	}
	return 0, func() (float64, SUBFN) {
		return factorial(f*n, n-1)
	}
}

func iterableFactorial2(n float64) (float64, SUBFN) {
	return factorial(1, n)
}

func main() {
	start := time.Now()
	fmt.Println(simpleFactorial(100))
	fmt.Println(time.Now().Sub(start))
	start = time.Now()
	fmt.Println(trampoline(iterableFactorial2)(100))
	fmt.Println(time.Now().Sub(start))
}
