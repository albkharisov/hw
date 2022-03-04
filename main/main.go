package main

import (
	"fmt"
	"github.com/albkharisov/hw/hw03_frequency_analysis"
)

func main() {
	res := "albert" > "albatros"
	fmt.Println("'albert' > 'albatros':", res)

	input := "a b c d d b b c c a a e e f f a f"
	fmt.Println("input: ", input)
	output := hw03frequencyanalysis.Top10(input)
	fmt.Println("output: ", output)
}
