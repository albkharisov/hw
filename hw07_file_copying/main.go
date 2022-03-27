package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcopy_file -from <input_file> -to <output_file> -limit <limit> -offset <offset>")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		printUsage()
		return
	}

	err := Copy(from, to, offset, limit)

	if err != nil {
		fmt.Println("Copying failed: ", err)
	}
}
