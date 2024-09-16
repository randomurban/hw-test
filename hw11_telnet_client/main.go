package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think if it is useful with blocking operation?
	var host string
	var port int
	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout")
	pflag.Parse()
	var err error
	host = pflag.Arg(0)
	if host == "" {
		usageAndExit("Empty host")
	}
	port, err = strconv.Atoi(pflag.Arg(1))
	if err != nil {
		usageAndExit("Bad port")
	}
	fmt.Println("host:", host)
	fmt.Println("port:", port)
}

func usageAndExit(msg string) {
	fmt.Printf("Usage: %s --timeout duration host port\n", os.Args[0])
	pflag.PrintDefaults()
	if msg != "" {
		fmt.Println(msg)
	}
	os.Exit(1)
}
