package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Connect implements TelnetClient.
func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("faled to connect: %w", err)
	}
	t.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	return nil
}

// Close implements TelnetClient.
func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

// Send implements TelnetClient.
func (t *telnetClient) Send() error {
	buffer := make([]byte, 1024)
	n, err := t.in.Read(buffer)
	if err != nil {
		return err
	}
	_, writeErr := t.conn.Write(buffer[:n])
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// Receive implements TelnetClient.
func (t *telnetClient) Receive() error {
	buffer := make([]byte, 1024)
	n, err := t.conn.Read(buffer)
	if err != nil {
		return err
	}
	_, writeErr := t.out.Write(buffer[:n])
	if writeErr != nil {
		return writeErr
	}
	return nil
}
