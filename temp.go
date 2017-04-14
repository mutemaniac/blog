package main

import "fmt"

func main() {
	func() {
		i := 0
		defer func(j int) {
			fmt.Println("i = ", i)
			fmt.Println("j = ", j)
		}(i + 1)
		i = 3
	}()
	fmt.Println("Hello, 世界")
}
