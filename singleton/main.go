package main

import (
	"fmt"
	"github.com/py-patterns/singleton"
)

func main() {
	fmt.Println(GetSingleton())

	signle := NewSingleton("redis://localhost")
	fmt.Println(signle)
	fmt.Println(GetSingleton())

	ResetSingleton()

	fmt.Println(GetSingleton())
}
