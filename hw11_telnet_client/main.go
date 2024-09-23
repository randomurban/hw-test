package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think if it is useful with blocking operation?
	var host string
	var port string
	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout")
	pflag.Parse()
	var err error
	host = pflag.Arg(0)
	if host == "" {
		usageAndExit("Illegal args :empty host")
	}
	port = pflag.Arg(1)
	_, err = strconv.Atoi(port)
	if err != nil {
		usageAndExit("Illegal args :bad port")
	}
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	fmt.Printf("...Connecting to %s\n", address)
	err = client.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to %s: %s\n", address, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("Interrupted by SIGINT")
		cancel()
	}()

	go func() {
		err := client.Send()
		defer cancel()
		if err != nil {
			fmt.Printf("Failed to send to %s: %s\n", address, err)
			return
		}
		fmt.Println("...EOF")
	}()

	go func() {
		err := client.Receive()
		defer cancel()
		if err != nil {
			fmt.Printf("Failed to receive from %s: %s\n", address, err)
			return
		}
		fmt.Println("...Connection was closed by peer")
	}()

	<-ctx.Done()
}

func usageAndExit(msg string) {
	fmt.Printf("Usage: %s --timeout duration host port\n", os.Args[0])
	pflag.PrintDefaults()
	if msg != "" {
		fmt.Println(msg)
	}
	os.Exit(1)
}
