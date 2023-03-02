package main

import (
	"fmt"
	"lecture01_homework/fizzbuzz"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(fizzbuzz.FizzBuzz(i))
	}
	// TODO тут напишите цикл с вызовом FizzBuzz
	// fmt.Println(fizzbuzz.FizzBuzz(10))
}
