package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	result, err := sqrt(-1)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

func sqrt(param float64) (float64, error) {
	if param < 0 {
		return 0, errors.New("Undefined integer value. Should be higher than 0 value.")
	}
	return math.Sqrt(param), nil
}
