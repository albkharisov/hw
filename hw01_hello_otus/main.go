package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	revertedString := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(revertedString)
}
