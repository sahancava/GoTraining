package main

import "fmt"

func main() {
	x := [5]string{"bir", "iki", "üç", "dört", "beş"}
	for ia, v := range x {
		fmt.Println(ia, v)
	}
}
