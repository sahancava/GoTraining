package main

import (
	"fmt"
)

func main() {
	m := map[string][]string{
		"marmara": []string{"İstanbul", "Bursa", "Kocaeli"},
		"ege":     []string{"İzmir", "Aydın", "Muğla"},
		"akdeniz": []string{"Antalya", "Mersin", "Adana"},
	}

	m["trakya"] = []string{"Kırklareli", "Edirne", "Tekirdağ"}

	delete(m, "marmara")

	for i, v := range m {
		fmt.Println("Record For: ", i)
		for j, val := range v {
			fmt.Printf("\tIndex: %v, Value: %v\n", j, val)
		}
	}
}
