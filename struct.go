package main

import (
	"fmt"
)

type person struct {
	name    string
	surname string
	age     int
}

type student struct {
	person
	isStudent bool
}

func main() {

	p1 := student{
		person: person{
			name:    `Şahan`,
			surname: `ÇAVA`,
			age:     29},
		isStudent: true,
	}

	for i, v := range []interface{}{p1} {
		fmt.Println("Position: ", i)
		for _, val := range []interface{}{v} {
			fmt.Println(val)
		}
	}
}
