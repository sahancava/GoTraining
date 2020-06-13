package main

import "fmt"

func main() {
	i := 1
	for {
		fmt.Println(i)
		i++
		if i > 9 {
			break
		}
	}
}
