package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	ErrNotConnected     = errors.New("Connection not established")
	ErrConnectFailed    = errors.New("Connect failed")
	ErrConnectionClosed = errors.New("Connection closed")
)

var (
	logger = log.New(os.Stderr, "", 0)

	timeout *time.Duration
	address string
)

func printUsage() {
	fmt.Println("go-telnet\nUsage:")
	fmt.Println("  go-telnet --timeout=10s <host> <port>")
	fmt.Println("  go-telnet mysite.ru 8080")
	fmt.Println("  go-telnet --timeout=3s 1.1.1.1 123")
	fmt.Println("Named arguments:")
	flag.PrintDefaults()
}

func main() {
	timeout = flag.Duration("timeout", time.Duration(10*time.Second), "Connection establish timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		printUsage()
	}
	address = net.JoinHostPort(args[0], args[1])

	process()
}

func process() {
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		logger.Println("...Cannot connect: ", err)
		os.Exit(1)
	} else {
		logger.Println("...Connected to: ", address)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		client.Send()
		logger.Println("...EOF")
		os.Exit(0)
	}()

	go func() {
		client.Receive()
		logger.Println("...Server closed connection")
		os.Exit(0)
	}()

	<-ctx.Done()
}
