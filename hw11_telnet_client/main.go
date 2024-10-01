package main

import (
	"context"
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
		usageAndExit("Invalid args :empty host")
	}
	port = pflag.Arg(1)
	_, err = strconv.Atoi(port)
	if err != nil {
		usageAndExit("Invalid args :bad port")
	}
	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	println("...Connecting to %s", address)
	err = client.Connect()
	if err != nil {
		println("Failed to connect to %s: %s", address, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		println("Interrupted by SIGINT")
		cancel()
	}()

	go func() {
		err := client.Send()
		defer cancel()
		if err != nil {
			println("Failed to send to %s: %s", address, err)
			return
		}
		println("...EOF")
	}()

	go func() {
		err := client.Receive()
		defer cancel()
		if err != nil {
			println("Failed to receive from %s: %s", address, err)
			return
		}
		println("...Connection was closed by peer")
	}()

	<-ctx.Done()
}

func usageAndExit(msg string) {
	println("Usage: %s --timeout duration host port", os.Args[0])
	pflag.PrintDefaults()
	if msg != "" {
		println(msg)
	}
	os.Exit(1)
}
