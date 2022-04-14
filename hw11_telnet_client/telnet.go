package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type telnetClientInstance struct {
	TelnetClient
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClientInstance) Close() error {
	return t.conn.Close()
}

func (t *telnetClientInstance) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClientInstance) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnetClientInstance) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tc := &telnetClientInstance{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}

	return tc
}
