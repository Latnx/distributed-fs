package main

import "fmt"

func Hello(name string) string {
	return "hello" + name
}

// test test
func main() {
	fmt.Println(Hello(" world"))
}
