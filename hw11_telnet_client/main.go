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
	"sync"
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

	err := client.Connect()
	if err != nil {
		logger.Println("...Cannot connect: ", err)
		os.Exit(1)
	} else {
		logger.Println("...Connected to", address)
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	eof := make(chan struct{}, 1)
	connectionClosed := make(chan struct{}, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := client.Send(); err != nil {
			/* Graceful goroutine shutdown (using sync.WaitGroup{}) forces
			   to show error on receiving, so comment it */
			// logger.Printf("...Send error: %v", err)
		}
		eof <- struct{}{}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := client.Receive(); err != nil {
			/* Graceful goroutine shutdown (using sync.WaitGroup{}) forces
			   to show error on receiving, so comment it */
			// logger.Printf("...Receive error: %v", err)
		}
		connectionClosed <- struct{}{}
		wg.Done()
	}()

	select {
	case <-ctx.Done():
		logger.Println("...Ctrl+C")
		os.Exit(2)
	case <-eof:
		logger.Println("...EOF")
		client.Close()
	case <-connectionClosed:
		logger.Println("...Server closed connection")
		client.Close()
		/* According to task: don't signal client.Send() routine about connection closing,
		* and finish program after next sending attempt (sounds illogical).
		* Send() routine releases sending if connection is closed. */
	}
	wg.Wait()
}
