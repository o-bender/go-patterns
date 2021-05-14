package main

import (
	"fmt"
	"github.com/o-bender/go-patterns/singleton/singleton"
)

func main() {
	fmt.Println(singleton.GetSingleton())

	signle := singleton.NewSingleton("redis://localhost")
	fmt.Println(signle)
	fmt.Println(singleton.GetSingleton())

	singleton.ResetSingleton()

	fmt.Println(singleton.GetSingleton())
}
