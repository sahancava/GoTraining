package main

import (
	"encoding/json"
	"fmt"
)

type persona struct {
	name    string `json:"name"`
	surname string `json:"surname"`
	age     int `json:"age"`
}

type talebe struct {
	persona
	isTalebe bool `json:"isTalebe"`
}

/*func (obj *persona) info(){
	if obj.name == "" {
		obj.name = "John Doe"
	}
	if obj.age == 0 {
		obj.age	= 29
	}
}*/

func main() {

	/*json_string := `
	{
	"name": "",
"surname":"",
"age":5,
	}`*/

	ogrenci := talebe{
		persona: persona{
			name:    "asds",
			surname: "Doe",
			age:     29,
		},
		isTalebe: true,
	}
	//ogrenci.info()
	for _, v := range []interface{}{ogrenci} {
		for i, val := range []interface{}{v} {
			fmt.Println(i, val)
		}
	}

	jsonSrt, _ := json.Marshal(ogrenci)
	fmt.Println(jsonSrt)
}
