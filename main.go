package main

import "fmt"

func main() {
	print(greet("John"))
}

func greet(s string) string {
	return fmt.Sprintf("Hello %s!", s)
}
