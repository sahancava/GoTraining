package main

import "fmt"

type kisi struct {
	name    string
	surname string
	age     int
}

type ogrenci struct {
	kisi
	isStudent bool
}

func main() {
	person1 := ogrenci{
		kisi: kisi{
			name:    "John",
			surname: "Doe",
			age:     29,
		},
		isStudent: true,
	}

	for i, v := range []interface{}{person1} {
		fmt.Println("Index: ", i)
		for j, val := range []interface{}{v} {
			fmt.Printf("Column: %v\t, Value: %v\n", j, val)
		}
	}
}
