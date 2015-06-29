// 29 june 2015
package main

import (
	"fmt"
	"os"
)

func alert(method string, romname string) {
	fmt.Printf("%-10s %s\n", method, romname)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error ")
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s datfile folder\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	datfile := os.Args[1]
	folder := os.Args[2]

	err := collectFilenames(folder)
	if err != nil {
		die("collecting filenames: %v", err)
	}

	printLeftovers()
}
