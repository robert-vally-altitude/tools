package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("Usage: go run migrate.go -from=<email> -to=<email>")
	fmt.Println("")
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	from := flag.String("from", "", "Source email address")
	to := flag.String("to", "", "Destination email address")

	flag.CommandLine.Parse(os.Args[1:])

	fmt.Println("Source GUser,Target GUser")
	fmt.Printf("%s,%s\n", *from, *to)
}
