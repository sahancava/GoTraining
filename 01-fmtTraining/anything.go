package main

import "fmt"

func main() {
	var x string = "deneme\n";
	var i int
	var b bool;
	fmt.Println("\nHere is my first line.\n")
	fmt.Printf("x is a %T while i is an %T and also the variable b is a %T\n",x,i,b)
	foo()
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}
}
func foo() {
	fmt.Println("I'm in foo\n")
}
