package main

import "fmt"

type persona struct {
	name    string
	surname []string
	age     int
}

type talebe struct {
	persona
	isTalebe bool
}

func main() {
	ogrenci := talebe{
		persona: persona{
			name:    "John",
			surname: []string{"Doe","Second Surname"},
			age:     29,
		},
		isTalebe: true,
	}
	for _, v := range []interface{}{ogrenci} {
		for i, val := range []interface{}{v} {
			fmt.Println(i, val)
		}
	}
}
