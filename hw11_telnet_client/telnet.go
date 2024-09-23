package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &simpleTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type simpleTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *simpleTelnetClient) Connect() error {
	if t.conn != nil {
		return fmt.Errorf("already connected")
	}
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *simpleTelnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *simpleTelnetClient) Send() error {
	if t.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *simpleTelnetClient) Receive() error {
	if t.conn == nil {
		return fmt.Errorf("not connected")
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
